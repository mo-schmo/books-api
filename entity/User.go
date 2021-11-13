package entity

type User struct {
	Username     string `dynamodbav:"Username" json:"Username"`
	Email        string `dynamodbav:"Email" json:"Email"`
	First        string `dynamodbav:"First_Name" json:"First_Name"`
	Last         string `dynamodbav:"Last_Name" json:"Last_Name"`
	Date_Created string `dynamodbav:"Date_Created" json:"Date_Created"`
}
