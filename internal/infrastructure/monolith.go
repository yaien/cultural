package infrastructure

import (
	"context"
	"log"
	"log/slog"
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
	WebRouter   *http.ServeMux
}

func NewMonolith() *Monolith {
	config := NewConfig()

	// Create MongoDB client
	database := setupMongoDB(config)

	return &Monolith{
		Config:    config,
		MongoDB:   database,
		Router:    http.DefaultServeMux,
		WebRouter: http.NewServeMux(),
	}
}

func setupMongoDB(config *Config) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(config.MongoDB.URI)
	opts.SetLoggerOptions(options.Logger().SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug).SetSink(&sink{}))

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	log.Println("Connected to MongoDB successfully")

	database := client.Database(config.MongoDB.Database)

	return database
}

var _ options.LogSink = (*sink)(nil)

type sink struct {
}

func (s *sink) Info(level int, msg string, args ...any) {
	slog.With(args...).Debug("MongoDB", "level", level, "msg", msg)
}

func (s *sink) Error(err error, msg string, args ...any) {
	slog.With(args...).Error("MongoDB", "error", err, "msg", msg)
}
