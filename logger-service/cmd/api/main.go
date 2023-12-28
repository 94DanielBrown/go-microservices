package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"logger/cmd/infrastructure"
	"logger/cmd/initializers"
	"time"
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

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

}

func connectToDynamo() (*dynamodb.Client, error) {
	config, err := infrastructure.NewAwsConfig()
	if err != nil {
		log.Fatal(err)
	}
	client := infrastructure.NewDynamoDBClient(config)
	return client, nil
}
