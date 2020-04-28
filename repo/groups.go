package repo

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type Group struct {
	ID      string `dynamodbav:"ID"`
	Name    string `dynamodbav:"name"`
	Creator string `dynamodbav:"creator"`
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
