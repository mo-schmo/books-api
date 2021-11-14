package entity

type User struct {
	Username          string `dynamodbav:"Username" json:"Username"`
	Email             string `dynamodbav:"Email" json:"Email"`
	First             string `dynamodbav:"First_Name" json:"First_Name"`
	Last              string `dynamodbav:"Last_Name" json:"Last_Name"`
	Created_Timestamp string `dynamodbav:"Created_Timestamp" json:"Created_Timestamp"`
	Updated_Timestamp string `dynamodbav:"Updated_Timestamp" json:"Updated_Timestamp"`
	Hash              string `dynamodbav:"Hash" json:"Hash"`
}
