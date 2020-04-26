package repo

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type Question struct {
	ID            string `dynamodbav:"ID"`
	Question      string `dynamodbav:"Question"`
	AnswerA       string `dynamodbav:"A"`
	AnswerB       string `dynamodbav:"B"`
	AnswerC       string `dynamodbav:"C"`
	AnswerD       string `dynamodbav:"D"`
	CorrectAnswer string `dynamodbav:"CorrectAnswer"`
}

func (r *Repo) GetQuestion(questionID string) (Question, error) {
	log.Printf("Getting question with ID %s. \n", questionID)
	queryOutput, err := r.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("personal-triviaQuestions"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(questionID),
			},
		},
	})
	if err != nil {
		return Question{}, errors.New("could not retrieve question from dynamo")
	}

	log.Println("Successfully retrieved question from dynamo.")

	question := Question{}
	if err := dynamodbattribute.UnmarshalMap(queryOutput.Item, &question); err != nil {
		return Question{}, err
	}
	return question, nil
}

func (r *Repo) CreateQuestion(in Question) (string, error) {
	log.Println("Attempting to create a question in dynamo.")
	in.ID = uuid.New().String()
	item, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		return "", errors.New("could not marshal created question into dynamo map")
	}

	if _, err = r.svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("personal-triviaQuestions"),
		Item:      item,
	}); err != nil {
		return "", errors.New("could not put created question into dynamo")
	}

	log.Printf("Successfully created a question in dynamo with the ID %s", in.ID)

	return in.ID, nil
}
