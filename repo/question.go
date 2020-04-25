package repo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Question struct {
	Question      string `dynamodbav:"Question"`
	AnwserA       string `dynamodbav:"A"`
	AnwserB       string `dynamodbav:"B"`
	AnwserC       string `dynamodbav:"C"`
	AnwserD       string `dynamodbav:"D"`
	CorrectAnwser string `dynamodbav:"CorrectAnwser"`
}

func (r *Repo) GetQuestion(questionID string) (Question, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String("personal-triviaQuestions"),
		Limit:     aws.Int64(1),
		KeyConditions: map[string]*dynamodb.Condition{
			"ID": {
				ComparisonOperator: aws.String("EQ"), // this is how we are comparing so here it is equals
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(questionID),
					},
				},
			},
		},
	}

	queryOutput, err := r.svc.Query(queryInput)
	if err != nil {
		return Question{}, err
	}

	questions := []Question{}
	if err := dynamodbattribute.UnmarshalListOfMaps(queryOutput.Items, &questions); err != nil {
		return Question{}, err
	}
	return questions[0], nil
}
