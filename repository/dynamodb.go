package repository

import (
	"booksApi/entity"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	dynamoDBClient := createDynamoDBClient()
	params := mux.Vars(r)
	userId := params["userId"]

	res, err := dynamoDBClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Users"),
		Key: map[string]types.AttributeValue{
			"Username": &types.AttributeValueMemberS{Value: userId},
		},
	})

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
