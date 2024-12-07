package mongodb

import (
	"os"
	"strconv"

	"github.com/agent-auth/common-lib/pkg/mongodb_client"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

// NewMongoClient returns new instance of datastore
func NewMongoClient() *mongo.Database {
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

	return mongodb_client.NewMongoStore(dbHost, dbName, queryTimeout).Database()
}
