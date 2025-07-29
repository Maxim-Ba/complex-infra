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
    // Устанавливаем текущее время, если оно не задано
    if message.CreatedAt.IsZero() {
			// TODO делать это на стороне монго "createdAt": bson.M{"$currentDate": bson.M{"$type": "date"}}
        message.CreatedAt = time.Now()
    }

    filter := bson.M{fmt.Sprintf("groups.%s", message.Group): bson.M{"$exists": true}}
    update := bson.M{
        "$push": bson.M{
            fmt.Sprintf("groups.%s", message.Group): message,
        },
    }

    // Пытаемся добавить в существующую группу
    result, err := r.collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return fmt.Errorf("failed to update message: %w", err)
    }

    // Создаем новую группу при необходимости
    if result.MatchedCount == 0 {
        update = bson.M{
            "$set": bson.M{
                fmt.Sprintf("groups.%s", message.Group): []models.MessageDTO{message},
            },
        }
        _, err = r.collection.UpdateOne(ctx, bson.M{}, update, options.Update().SetUpsert(true))
        if err != nil {
            return fmt.Errorf("failed to create new group: %w", err)
        }
    }

    return nil
}



func (r *MongoRepository) GetMessagesByGroup(ctx context.Context, req models.RequestMessages) ([]models.MessageDTO, error) {
    var result struct {
        Groups map[string][]models.MessageDTO `bson:"groups"`
    }

    err := r.collection.FindOne(ctx, bson.M{
        fmt.Sprintf("groups.%s", req.GroupiD): bson.M{"$exists": true},
    }).Decode(&result)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return []models.MessageDTO{}, nil
        }
        return nil, fmt.Errorf("failed to find group messages: %w", err)
    }

    messages := result.Groups[req.GroupiD]
    
    start := int(req.Offset)
    if start < 0 {
        start = 0
    }
    
    end := start + int(req.Count)
    if end > len(messages) {
        end = len(messages)
    }
    
    // Если запрошено больше, чем есть, или неверные параметры - возвращаем пустой массив
    if start >= len(messages) || req.Count <= 0 {
        return []models.MessageDTO{}, nil
    }
    
    // Возвращаем срез сообщений в обратном порядке (новые первыми)
    reversed := make([]models.MessageDTO, len(messages))
    for i, j := 0, len(messages)-1; j >= 0; i, j = i+1, j-1 {
        reversed[i] = messages[j]
    }
    
    return reversed[start:end], nil
}
