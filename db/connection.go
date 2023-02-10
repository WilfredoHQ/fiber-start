package db

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/wilfredohq/fiber-start/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCollection(client *mongo.Client, name string) *mongo.Collection {
	collection := client.Database(config.Config.DBName).Collection(name)
	return collection
}

func connectDB() *mongo.Client {
	uri := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority",
		url.QueryEscape(config.Config.DBUser),
		url.QueryEscape(config.Config.DBPassword),
		config.Config.DBHost,
	)

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	if err := createSuperuser(ctx, client); err != nil {
		log.Fatal(err)
	}

	return client
}

var DB = connectDB()
