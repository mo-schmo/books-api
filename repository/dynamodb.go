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
		var user entity.User
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if user == (entity.User{}) {
			fmt.Println("Empty body")
			return
		}

		// Check if user already exists in table
		res, _ := getUserFromTable(user.Username)
		if res.Item != nil {
			fmt.Println("User already exists and/or could not be created")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Continue if user does not already exists
		now := time.Now().UTC().Format("2006-01-02T15:04:05Z07:00")
		user.Created_Timestamp = now
		user.Updated_Timestamp = now

		fmt.Println("Password: ", user.Password)
		hashedPassword, err := util.HashPassword(user.Password)
		if err != nil {
			fmt.Println("Error hashing password: ", err)
		}
		user.Hash = hashedPassword
		fmt.Printf("%+v\n", user)
		_, err = addUserToTable(&user)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func addUserToTable(user *entity.User) (*dynamodb.PutItemOutput, error) {
	dynamodbClient := createDynamoDBClient()
	marshaledUser, err := attributevalue.MarshalMap(user)
	if err != nil {
		return nil, err
	}
	res, err := dynamodbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Users"),
		Item:      marshaledUser,
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
