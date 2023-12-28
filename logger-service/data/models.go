package data

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"time"
)

var client *dynamodb.Client

func New(dynamo *dynamodb.Client) Models {
	client = dynamo

	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `json:"id,omitempty" dynamodbav:"id"`
	Name      string    `json:"name,omitempty" dynamodbav:"name"`
	Data      string    `json:"data,omitempty" dynamodbav:"data"`
	CreatedAt time.Time `json:"created_at,omitempty" dynamodbav:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" dynamodbav:"updated_at"`
}

func (l *LogEntry) Put(entry LogEntry) (*dynamodb.PutItemOutput, error) {
	av, err := attributevalue.MarshalMap(entry)
	if err != nil {
		fmt.Printf("Got error marshalling data: %s\n", err)
		return nil, err
	}
	// save chat to db
	output, err := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(""), Item: av,
	})
	if err != nil {
		fmt.Printf("Couldn't add item to table.: %v\n", err)
	}
	return output, err
}
