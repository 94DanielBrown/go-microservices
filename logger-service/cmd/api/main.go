package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.mongodb.org/mongo-driver/mongo"
	"log"

	"logger/cmd/infrastructure"
	"logger/cmd/initializers"
	"logger/data"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	initializers.LoadEnvVariables()
	client, err := connectToDynamo()
	if err != nil {
		log.Fatal(err)
	}

	app := Config{
		Models: data.New(client),
	}

	err = app.Models.LogEntry.Put(
		data.LogEntry{
			ID:   "test",
			Name: "test",
			Data: "test",
		})
	if err != nil {
		fmt.Println(err)
	}

}

func connectToDynamo() (*dynamodb.Client, error) {
	config, err := infrastructure.NewAwsConfig()
	if err != nil {
		log.Fatal(err)
	}
	client := infrastructure.NewDynamoDBClient(config)
	return client, nil
}
