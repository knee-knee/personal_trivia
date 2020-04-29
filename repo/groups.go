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

type Group struct {
	ID      string `dynamodbav:"ID"`
	Name    string `dynamodbav:"name"`
	Creator string `dynamodbav:"creator"`
}

type GroupScanInput struct {
	Key   string
	Value string
}

func (r *Repo) CreateGroup(in Group) (string, error) {
	log.Printf("repo: Attempting to create a group named; %s \n", in.Name)
	in.ID = uuid.New().String()
	item, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		return "", errors.New("could not marshal created question into dynamo map")
	}

	if _, err = r.svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("personal-triviaGroups"),
		Item:      item,
	}); err != nil {
		return "", errors.New("could not put created question into dynamo")
	}

	log.Printf("repo: Successfully created group %s \n", in.ID)

	return in.ID, nil
}

func (r *Repo) GetGroup(groupID string) (Group, error) {
	log.Printf("Getting group with ID %s. \n", groupID)
	queryOutput, err := r.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("personal-triviaGroups"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(groupID),
			},
		},
	})
	if err != nil {
		return Group{}, errors.New("could not retrieve group from dynamo")
	}

	log.Println("Successfully retrieved group from dynamo.")

	group := Group{}
	if err := dynamodbattribute.UnmarshalMap(queryOutput.Item, &group); err != nil {
		return Group{}, err
	}
	return group, nil
}

// TODO: This is a scan right now because I dont want to pay for a GSI.
// If there is ever a reason that I would want to pay for it make this a query.
// Also this is kind of shit because I assume you are only going to get one record back.
func (r *Repo) GetGroupsScan(in GroupScanInput) (Group, error) {
	log.Printf("repo: Getting group by scanning off the %s \n", in.Key)

	filt := expression.Name(in.Key).Equal(expression.Value(in.Value))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return Group{}, err
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		Limit:                     aws.Int64(1),
		TableName:                 aws.String("personal-triviaGroups"),
	}
	result, err := r.svc.Scan(params)
	if err != nil {
		return Group{}, err
	}

	if result.Count == nil {
		return Group{}, errors.New("count of scan was empty")
	}
	if *result.Count == 0 {
		return Group{}, nil
	}

	var group Group
	if err := dynamodbattribute.UnmarshalMap(result.Items[0], &group); err != nil {
		return Group{}, err
	}

	return group, nil
}
