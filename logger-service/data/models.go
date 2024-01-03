package data

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
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
	UUID      string `json:"uuid,omitempty" dynamodbav:"uuid"`
	Name      string `json:"name,omitempty" dynamodbav:"name"`
	Data      string `json:"data,omitempty" dynamodbav:"data"`
	CreatedAt int    `json:"created_at,omitempty" dynamodbav:"created_at"`
	UpdatedAt int    `json:"updated_at,omitempty" dynamodbav:"updated_at"`
}

// PutItem adds an item to the dynamodb table
func (l *LogEntry) PutItem(entry LogEntry) error {
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

// AllItems gets all items from the dynamodb table
func (l *LogEntry) AllItems() ([]*LogEntry, error) {
	var logs []*LogEntry

	input := &dynamodb.ScanInput{
		TableName: aws.String("logs"),
	}

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

// GetItem gets an item by ID
func (l *LogEntry) GetItem(uuid string) (*LogEntry, error) {
	var logEntry LogEntry

	input := &dynamodb.GetItemInput{
		TableName: aws.String("logs"),
		Key: map[string]types.AttributeValue{
			"UUID": &types.AttributeValueMemberS{
				Value: uuid,
			},
		},
	}

	data, err := client.GetItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if data.Item == nil {
		return nil, fmt.Errorf("no data found for UUID %s", uuid)
	}

	err = attributevalue.UnmarshalMap(data.Item, &logEntry)
	return &logEntry, err
}

func (l *LogEntry) Query(uuid string) ([]*LogEntry, error) {
	var logs []*LogEntry

	input := &dynamodb.QueryInput{
		TableName:              aws.String("logs"),
		KeyConditionExpression: aws.String("#uuidKey = :uuidVal"),
		ExpressionAttributeNames: map[string]string{
			"#uuidKey": "uuid",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uuidVal": &types.AttributeValueMemberS{Value: uuid},
		},
	}

	data, err := client.Query(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if data.Items == nil {
		return nil, fmt.Errorf("no data found for UUID %s", uuid)
	}

	// Unmarshal items from the page to LogEntry objects
	for _, item := range data.Items {
		var logEntry LogEntry
		if err := attributevalue.UnmarshalMap(item, &logEntry); err != nil {
			return nil, err
		}
		logs = append(logs, &logEntry)
	}

	return logs, nil
}
