package mongodb

import (
	"context"
	"log/slog"

	"github.com/LashkaPashka/EventConsumerService/internal/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName = "eventHub"
	collectionName = "events"
)


type MongoDb struct {
	mongoCli *mongo.Client
	cfg *configs.Configs
	logger *slog.Logger
}

func New(cfg *configs.Configs, logger *slog.Logger) *MongoDb {
	var cli = MongoDb{
		cfg: cfg,
		logger: logger,
	}
	
	cli.mongoCli = checkConnection(cli.cfg, cli.logger)
	if cli.mongoCli == nil {
		return nil
	}
	
	return &cli
}

func checkConnection(cfg *configs.Configs, logger *slog.Logger) *mongo.Client {
	const op = "EventConsumerService.storage.mongo.CheckConnection"

	clientOption := options.Client().ApplyURI(cfg.MongoUri)
	client, err := mongo.Connect(context.Background(), clientOption)
	if err != nil {
		logger.Error("Invalid connection to MongoDb", slog.String("Error: ", op))
		return nil
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		logger.Error("Invalid ping to MongoDb", slog.String("Error: ", op))
		return nil
	}
	
	return client
}

func (m *MongoDb) InsertDoc(eventModel any) any {
	const op = "EventConsumerService.storage.mongo.InserOne"

	collection := m.mongoCli.Database(databaseName).Collection(collectionName)

	res, err := collection.InsertOne(context.Background(), eventModel)
	if err != nil {
		m.logger.Error("Invalid insert Docs in collection", slog.String("Error: ", op))
		return nil
	}

	return res.InsertedID
}