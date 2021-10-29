package repository

import (
	"booksApi/entity"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func createDynamoDBClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	return dynamodb.NewFromConfig(cfg)
}

func GetItem(tableName string, id string) (*dynamodb.GetItemOutput, error) {
	dynamoDBClient := createDynamoDBClient()

	out, err := dynamoDBClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Users"),
		Key: map[string]types.AttributeValue{
			"Username": &types.AttributeValueMemberS{Value: "admin"},
		},
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return out, nil
}

func ScanUsers(w http.ResponseWriter, r *http.Request) {
	dynamoDBClient := createDynamoDBClient()

	res, err := dynamoDBClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("Users"),
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var users []entity.Users
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
