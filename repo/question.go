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
	AnwserA       string `dynamodbav:"A"`
	AnwserB       string `dynamodbav:"B"`
	AnwserC       string `dynamodbav:"C"`
	AnwserD       string `dynamodbav:"D"`
	CorrectAnwser string `dynamodbav:"CorrectAnswer"`
}

func (r *Repo) GetQuestion(questionID string) (Question, error) {
	log.Printf("Getting question with ID %s. \n", questionID)

	// queryOutput, err := r.svc.Query(&dynamodb.QueryInput{
	// 	TableName: aws.String("personal-triviaQuestions"),
	// 	Limit:     aws.Int64(1),
	// 	KeyConditions: map[string]*dynamodb.Condition{
	// 		"ID": {
	// 			ComparisonOperator: aws.String("EQ"), // this is how we are comparing so here it is equals
	// 			AttributeValueList: []*dynamodb.AttributeValue{
	// 				{
	// 					S: aws.String(questionID),
	// 				},
	// 			},
	// 		},
	// 	},
	// })
	queryOutput, err := r.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("personal-triviaQuestions"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(questionID),
			},
		},
		// KeyConditions: map[string]*dynamodb.Condition{
		// 	"ID": {
		// 		ComparisonOperator: aws.String("EQ"), // this is how we are comparing so here it is equals
		// 		AttributeValueList: []*dynamodb.AttributeValue{
		// 			{
		// 				S: aws.String(questionID),
		// 			},
		// 		},
		// 	},
		// },
	})
	if err != nil {
		return Question{}, errors.New("could not retrieve question from dynamo")
	}

	log.Println("Succesfully retrieved question from dyanmo.")

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
		return "", errors.New("could not put created question into dyanmo")
	}

	log.Printf("Succesfully created a question in dynamo with the ID %s", in.ID)

	return in.ID, nil
}
