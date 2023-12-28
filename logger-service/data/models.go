package data

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
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

func (l *LogEntry) Put(entry LogEntry) error {
	av, err := attributevalue.MarshalMap(entry)
	if err != nil {
		fmt.Printf("Got error marshalling data: %s\n", err)
		return err
	}
	// save chat to db
	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("logs"), Item: av,
	})
	if err != nil {
		log.Println("Couldn't add item to table.: %v\n", err)
	}
	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	// Define a slice to hold log entries
	var logs []*LogEntry

	// Create input parameters for the Scan operation
	input := &dynamodb.ScanInput{
		TableName: aws.String("logs"),
	}

	// Perform the Scan operation
	paginator := dynamodb.NewScanPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}

		// Unmarshal items from the page to LogEntry objects
		for _, item := range page.Items {
			var logEntry LogEntry
			if err := attributevalue.UnmarshalMap(item, &logEntry); err != nil {
				return nil, err
			}
			logs = append(logs, &logEntry)
		}
	}

	return logs, nil
}
