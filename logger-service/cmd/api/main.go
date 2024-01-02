package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"

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

	uid := uuid.New()
	err = app.Models.LogEntry.Put(
		data.LogEntry{
			UUID:      uid.String(),
			Name:      "test",
			Data:      "test",
			CreatedAt: int(time.Now().Unix()),
			UpdatedAt: int(time.Now().Unix()),
		})
	if err != nil {
		fmt.Println(err)
	}

	logs, err := app.Models.LogEntry.All()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", logs)

	for _, log := range logs {
		fmt.Printf("UUID: %s, Name: %s, Data: %s, CreatedAt: %d, UpdatedAt: %d\n", log.UUID, log.Name, log.Data, log.CreatedAt, log.UpdatedAt)
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
