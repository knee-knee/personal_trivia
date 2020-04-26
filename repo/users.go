package repo

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type User struct {
	ID       string `dynamodbav:"ID"`
	Email    string `dynamodbav:"email"`
	Password string `dynamodbav:"password"`
	Username string `dynamodbav:"username"`
	Group    string `dynamodbav:"group"`
}

func (r *Repo) CreateUser(in User) (string, error) {
	log.Printf("repo: Attempting to create user with email %s \n", in.Email)
	in.ID = uuid.New().String()
	item, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		return "", errors.New("could not marshal created question into dynamo map")
	}

	if _, err = r.svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("personal-triviaUsers"),
		Item:      item,
	}); err != nil {
		return "", errors.New("could not put created question into dynamo")
	}

	log.Printf("repo: Successfully created user %s \n", in.ID)

	return in.ID, nil
}
