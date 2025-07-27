package mongo

import (
	"context"
	"fmt"
	"go-messages/internal/app"
	"go-messages/internal/models"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func New(config app.AppConfig) (*MongoRepository, error) {
	cfg := config.GetConfig()
	clientOptions := options.Client().ApplyURI(cfg.MongoDBURI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Проверяем подключение
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}
	slog.Info("MongoRepository Ping done")

	db := client.Database(cfg.MongoDBDatabase)
	collection := db.Collection(cfg.MongoDBCollection)

	return &MongoRepository{
		client:     client,
		database:   db,
		collection: collection,
	}, nil
}

func (r *MongoRepository) Close(ctx context.Context) error {
	if err := r.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}
	return nil
}

func (r *MongoRepository) SaveMessage(ctx context.Context, message models.MessageDTO) error {
	_, err := r.collection.InsertOne(ctx, bson.M{
		"id":        message.Id,
		"producer":  message.Producer,
		"payload":   message.Payload,
		"createdAt": time.Now(),
	})
	if err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}
	return nil
}

func (r *MongoRepository) GetMessages(ctx context.Context) ([]models.MessageDTO, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find messages: %w", err)
	}
	defer cursor.Close(ctx)

	var messages []models.MessageDTO
	for cursor.Next(ctx) {
		var result struct {
			Id       string `bson:"id"`
			Producer string `bson:"producer"`
			Payload  string `bson:"payload"`
		}
		if err := cursor.Decode(&result); err != nil {
			slog.Error("Failed to decode message", "error", err)
			continue
		}
		messages = append(messages, models.MessageDTO{
			Id:       result.Id,
			Producer: result.Producer,
			Payload:  result.Payload,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return messages, nil
}
