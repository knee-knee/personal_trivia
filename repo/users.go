package repo

import (
	"errors"
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

func (r *Repo) GetUserByEmail(email string) (User, error) {
	log.Printf("Getting user with email %s. \n", email)

	filt := expression.Name("email").Equal(expression.Value(email))
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

	var user User
	if err := dynamodbattribute.UnmarshalMap(result.Items[0], &user); err != nil {
		return User{}, err
	}

	return user, nil
}
