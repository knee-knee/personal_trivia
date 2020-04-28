package repo

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
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

func (r *Repo) GetUser(userID string) (User, error) {
	log.Printf("Getting user with ID %s. \n", userID)
	queryOutput, err := r.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("personal-triviaUsers"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(userID),
			},
		},
	})
	if err != nil {
		return User{}, errors.New("could not retrieve user from dynamo")
	}

	log.Println("Successfully retrieved user from dynamo.")

	user := User{}
	if err := dynamodbattribute.UnmarshalMap(queryOutput.Item, &user); err != nil {
		return User{}, err
	}
	return user, nil
}

// TODO: This is a scan right now because I dont want to pay for a GSI.
// If there is ever a reason that I would want to pay for it make this a query.
// Also this is kind of shit because I assume you are only going to get one record back.
func (r *Repo) GetUserScan(in UserScanInput) (User, error) {
	log.Printf("Getting user by scanning off the %s. \n", in.Key)

	filt := expression.Name(in.Key).Equal(expression.Value(in.Value))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return User{}, err
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		Limit:                     aws.Int64(1),
		TableName:                 aws.String("personal-triviaUsers"),
	}
	result, err := r.svc.Scan(params)
	if err != nil {
		return User{}, err
	}

	if result.Count == nil {
		fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
		return User{}, errors.New("count of scan was empty")
	}
	if *result.Count == 0 {
		return User{}, nil
	}

	var user User
	if err := dynamodbattribute.UnmarshalMap(result.Items[0], &user); err != nil {
		return User{}, err
	}

	return user, nil
}

type UserScanInput struct {
	Key   string
	Value string
}
