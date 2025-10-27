package infrastructure

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Monolith struct {
	Config      *Config
	MongoDB     *mongo.Database
	MongoClient *mongo.Client
	Router      *http.ServeMux
}

func NewMonolith() *Monolith {
	config := NewConfig()

	// Create MongoDB client
	database := setupMongoDB(config.MongoDB)

	return &Monolith{
		Config:  config,
		MongoDB: database,
		Router:  http.DefaultServeMux,
	}
}

func setupMongoDB(config MongoDBConfig) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.URI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	log.Println("Connected to MongoDB successfully")

	database := client.Database(config.Database)

	return database
}
