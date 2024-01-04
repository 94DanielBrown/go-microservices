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
	err = app.Models.LogEntry.PutItem(
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

	logs, err := app.Models.LogEntry.AllItems()
	if err != nil {
		fmt.Println(err)
	}

	for _, log := range logs {
		fmt.Printf("UUID: %s, Name: %s, Data: %s, CreatedAt: %d, UpdatedAt: %d\n", log.UUID, log.Name, log.Data, log.CreatedAt, log.UpdatedAt)
	}

	//GetItem
	retrievedLog, err := app.Models.LogEntry.GetItem("0235440a-5617-446d-b8de-27a30b08b711")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the details of the retrieved log entry
	fmt.Printf("Retrieved Log - UUID: %s, Name: %s, Data: %s, CreatedAt: %d, UpdatedAt: %d\n",
		retrievedLog.UUID, retrievedLog.Name, retrievedLog.Data, retrievedLog.CreatedAt, retrievedLog.UpdatedAt)

	// Query
	logs, err = app.Models.LogEntry.Query("0235440a-5617-446d-b8de-27a30b08b711")
	if err != nil {
		fmt.Println(err)
	}

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
