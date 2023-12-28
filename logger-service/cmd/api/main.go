package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"logger/cmd/infrastructure"
	"logger/cmd/initializers"
	"logger/cmd/utils"
)

var client *mongo.Client

type Config struct {
}

func main() {
	initializers.LoadEnvVariables()
	client, err := connectToDynamo()
	if err != nil {
		log.Fatal(err)
	}

	// List tables
	var tableNames []string
	tables, err := utils.ListTables(client)
	if err != nil {
		fmt.Println(err)
	} else {
		tableNames = tables.TableNames
	}
	fmt.Println(tableNames)

}

func connectToDynamo() (*dynamodb.Client, error) {
	config, err := infrastructure.NewAwsConfig()
	if err != nil {
		log.Fatal(err)
	}

	client := infrastructure.NewDynamoDBClient(config)
	return client, nil
}
