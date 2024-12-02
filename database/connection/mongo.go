package connection

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/agent-auth/agent-auth-api/web/interfaces/v1/healthinterface"

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

// Add environment variable validation
func validateEnvVars() (string, string, int) {
	dbHost := os.Getenv("MONGODB_URI")
	if dbHost == "" {
		zap.L().Fatal("MONGODB_URI environment variable is required")
	}

	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		zap.L().Fatal("MONGODB_DATABASE environment variable is required")
	}

	queryTimeout, err := strconv.Atoi(os.Getenv("MONGODB_QUERY_TIMEOUT_SECONDS"))
	if err != nil || queryTimeout <= 0 {
		zap.L().Fatal("MONGODB_QUERY_TIMEOUT_SECONDS must be a positive integer",
			zap.Error(err),
		)
	}

	return dbHost, dbName, queryTimeout
}

func (s *mongoStore) initialize() (a *mongo.Database, b *mongo.Client) {
	dbHost, dbName, queryTimeout := validateEnvVars()

	clientOptions := options.Client().ApplyURI(dbHost)
	clientOptions.SetServerSelectionTimeout(time.Duration(queryTimeout) * time.Second)

	var err error
	client, err = mongo.NewClient(clientOptions)
	if err != nil {
		s.logger.Fatal("Failed to create MongoDB client",
			zap.Error(err),
			zap.String("host", dbHost),
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(queryTimeout)*time.Second)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		s.logger.Fatal("Failed to connect to MongoDB",
			zap.Error(err),
			zap.String("host", dbHost),
		)
	}

	// Test the connection
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		s.logger.Fatal("Failed to ping MongoDB",
			zap.Error(err),
			zap.String("host", dbHost),
		)
	}

	db = client.Database(dbName)
	s.logger.Info("Successfully connected to MongoDB",
		zap.String("database", dbName),
		zap.String("host", dbHost),
	)

	return db, client
}

// Update Health method to use environment variables
func (s *mongoStore) Health() *healthinterface.OutboundInterface {
	once.Do(func() {
		db, client = s.initialize()
	})

	dbHost, _, queryTimeout := validateEnvVars()

	outbound := healthinterface.OutboundInterface{
		TimeStampUTC:     time.Now().UTC(),
		ConnectionStatus: healthinterface.ConnectionActive,
		ApplicationName:  "MongoDB",
		URLs:             []string{dbHost},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(queryTimeout)*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		outbound.ConnectionStatus = healthinterface.ConnectionDisconnected
		s.logger.Error("Failed to ping MongoDB during health check",
			zap.Error(err),
			zap.String("host", dbHost),
		)
	}

	return &outbound
}
