package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func main() {
	fmt.Println("Now strarting the web server")
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/questions", question)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func helloWorld(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func question(w http.ResponseWriter, req *http.Request) {
	question, err := getQuestion()
	if err != nil {
		fmt.Println(err)
	}
	if question != nil {
		fmt.Println("***********************************")
		fmt.Println(aws.String(question.CorrectAnwser))
		fmt.Println("***********************************")
	}
	fmt.Fprintf(w, "got the question")
}

// this will return a question from dyanmo
func getQuestion() (*Question, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String("us-east-2"),
			Credentials: crednetials.NewStaticCredentials("", "", ""),
		},
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)

	queryInput := &dynamodb.QueryInput{
		TableName: aws.String("personal-triviaQuestions"),
		Limit:     aws.Int64(1),
		KeyConditions: map[string]*dynamodb.Condition{
			"ID": {
				ComparisonOperator: aws.String("EQ"), // this is how we are comparing so here it is equals
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String("55ed5152-00ac-4356-92bd-eb54c03e31ec"),
					},
				},
			},
		},
	}

	queryOutput, err := svc.Query(queryInput)
	if err != nil {
		return nil, err
	}

	questions := []Question{}
	if err := dynamodbattribute.UnmarshalListOfMaps(queryOutput.Items, &questions); err != nil {
		return nil, err
	}
	fmt.Println(questions)

	return &questions[0], nil
}

type Question struct {
	Question      string `json:"question"`
	AnwserA       string `json:"anwserA"`
	AnwserB       string `json:"anwserB"`
	AnwserC       string `json:"anwserC"`
	AnwserD       string `json:"anwserD"`
	CorrectAnwser string `json:"correct_anwser"`
}
