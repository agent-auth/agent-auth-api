package connection

import (
	"context"
	"sync"
	"time"

	"github.com/agent-auth/agent-auth-api/web/interfaces/v1/healthinterface"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

var (
	once   sync.Once
	db     *mongo.Database
	client *mongo.Client
)

// MongoStore ...
type mongoStore struct {
	logger *zap.Logger
}

// NewMongoStore returns new instance of datastore
func NewMongoStore() MongoStore {
	logger, _ := zap.NewProduction()

	return &mongoStore{
		logger: logger,
	}
}

// Client returns mongodb client instance
func (s *mongoStore) Client() *mongo.Client {
	once.Do(func() {
		db, client = s.initialize()
	})

	return client
}

// Database returns mongodb database instance
func (s *mongoStore) Database() *mongo.Database {
	once.Do(func() {
		db, client = s.initialize()
	})

	return db
}

func (s *mongoStore) initialize() (a *mongo.Database, b *mongo.Client) {
	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("db.host")))
	if err != nil {
		s.logger.Fatal("Failed to connect to database",
			zap.Error(err),
		)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(viper.GetInt("db.query_timeout_in_sec"))*time.Second,
	)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		s.logger.Fatal("Failed to connect to database",
			zap.Error(err),
		)
	}

	database := viper.GetString("db.database")
	db := client.Database(database)
	s.logger.Info("Successfully connected to database",
		zap.String("database", database),
	)

	return db, client
}

func (s *mongoStore) Health() *healthinterface.OutboundInterface {
	once.Do(func() {
		db, client = s.initialize()
	})

	outbound := healthinterface.OutboundInterface{}
	outbound.TimeStampUTC = time.Now().UTC()
	outbound.ConnectionStatus = healthinterface.ConnectionActive
	outbound.ApplicationName = "MongoDB"
	outbound.URLs = []string{viper.GetString("db.host")}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(viper.GetInt("db.query_timeout_in_sec"))*time.Second,
	)
	defer cancel()

	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		outbound.ConnectionStatus = healthinterface.ConnectionDisconnected
	}

	return &outbound
}
