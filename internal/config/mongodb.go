package config

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/hbc/app-hrm-api-go/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var log = utils.NewLogger("MongoDB")

// ConnectMongoDB connects to MongoDB
func ConnectMongoDB() error {
	cfg := Get()
	
	// URL encode password to handle special characters
	encodedPassword := url.QueryEscape(cfg.MongoPass)
	
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		cfg.MongoUser,
		encodedPassword,
		cfg.MongoHost,
		cfg.MongoPort,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Error("Failed to connect to MongoDB: %v", err)
		return err
	}

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Error("Failed to ping MongoDB: %v", err)
		return err
	}

	mongoClient = client
	log.Info("Connect database success")
	return nil
}

// GetMongoDB returns the MongoDB database instance
func GetMongoDB() *mongo.Database {
	cfg := Get()
	return mongoClient.Database(cfg.MongoDB)
}

// GetMongoClient returns the MongoDB client
func GetMongoClient() *mongo.Client {
	return mongoClient
}

// CloseMongoDB closes the MongoDB connection
func CloseMongoDB() error {
	if mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return mongoClient.Disconnect(ctx)
	}
	return nil
}
