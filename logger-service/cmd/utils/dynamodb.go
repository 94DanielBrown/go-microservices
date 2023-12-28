package utils

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func ListTables(client *dynamodb.Client) (*dynamodb.ListTablesOutput, error) {
	tables, err := client.ListTables(
		context.TODO(),
		&dynamodb.ListTablesInput{},
	)

	return tables, err
}
