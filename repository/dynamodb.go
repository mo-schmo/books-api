package repository

import (
	"booksApi/entity"
	"booksApi/util"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gorilla/mux"
)

// Reference: https://dynobase.dev/dynamodb-golang-query-examples/#get-item

func createDynamoDBClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	return dynamodb.NewFromConfig(cfg)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userId"]

	res, err := getUserFromTable(userId)

	if res.Item == nil {
		w.Write([]byte(fmt.Sprintf("No account found with userId: %s", userId)))
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := entity.User{}
	err = attributevalue.UnmarshalMap(res.Item, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ret, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}

func ScanUsers(w http.ResponseWriter, r *http.Request) {
	dynamoDBClient := createDynamoDBClient()

	res, err := dynamoDBClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("Users"),
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var users []entity.User
	err = attributevalue.UnmarshalListOfMaps(res.Items, &users)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	b, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var userPayload struct {
			entity.User
			Password string
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &userPayload)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if userPayload == (struct {
			entity.User
			Password string
		}{}) {
			fmt.Println("Empty body")
			return
		}

		now := time.Now().UTC().Format("2006-01-02T15:04:05Z07:00")
		userPayload.Created_Timestamp = now
		userPayload.Updated_Timestamp = now

		hashedPassword, err := util.HashPassword(userPayload.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		userPayload.Hash = hashedPassword
		_, err = addUserToTable(&userPayload.User)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func ValidateUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		params := r.URL.Query()
		if len(params) < 2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		username := params.Get("username")
		password := params.Get("password")
		if len(password) < 1 || len(username) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		dynamodbClient := createDynamoDBClient()
		res, err := dynamodbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
			TableName: aws.String("Users"),
			Key: map[string]types.AttributeValue{
				"Username": &types.AttributeValueMemberS{Value: username},
			},
			ProjectionExpression: aws.String("Hashkey"),
		})
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var hash string
		err = attributevalue.Unmarshal(res.Item["Hashkey"], &hash)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		isValid := util.ComparePasswordHash(password, hash)
		if !isValid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func addUserToTable(user *entity.User) (*dynamodb.PutItemOutput, error) {
	dynamodbClient := createDynamoDBClient()
	marshaledUser, err := attributevalue.MarshalMap(user)
	if err != nil {
		return nil, err
	}
	res, err := dynamodbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName:           aws.String("Users"),
		Item:                marshaledUser,
		ConditionExpression: aws.String("attribute_not_exists(Username)"),
	})
	return res, err
}

func getUserFromTable(userId string) (*dynamodb.GetItemOutput, error) {
	dynamoDBClient := createDynamoDBClient()
	return dynamoDBClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Users"),
		Key: map[string]types.AttributeValue{
			"Username": &types.AttributeValueMemberS{Value: userId},
		},
	})
}
