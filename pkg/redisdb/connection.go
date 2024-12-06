package redisdb

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/agent-auth/agent-auth-api/database/connection"
	"github.com/agent-auth/agent-auth-api/pkg/logger"
	"github.com/agent-auth/agent-auth-api/web/interfaces/v1/healthinterface"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	once   sync.Once
	client *redis.Client
)

// RedisStore ...
type redisStore struct {
	logger       *zap.Logger
	db           connection.MongoStore
	dbHost       string
	queryTimeout int
}

// NewRedisStore returns new instance of datastore
func NewRedisStore() RedisStore {
	dbHost, queryTimeout := validateEnvVars()

	r := &redisStore{
		logger:       logger.NewLogger(),
		db:           connection.NewMongoStore(),
		dbHost:       dbHost,
		queryTimeout: queryTimeout,
	}

	once.Do(func() {
		client = r.initialize()
	})

	return r
}

// Client returns redis client instance
func (s *redisStore) Client() *redis.Client {
	return client
}

// Add environment variable validation
func validateEnvVars() (string, int) {
	dbHost := os.Getenv("REDIS_URI")
	if dbHost == "" {
		zap.L().Fatal("REDIS_URI environment variable is required")
	}

	queryTimeout, err := strconv.Atoi(os.Getenv("REDIS_QUERY_TIMEOUT_SECONDS"))
	if err != nil || queryTimeout <= 0 {
		zap.L().Fatal("REDIS_QUERY_TIMEOUT_SECONDS must be a positive integer",
			zap.Error(err),
		)
	}
	return dbHost, queryTimeout
}

func (s *redisStore) initialize() *redis.Client {
	s.logger.Info("Initializing Redis client", zap.String("host", s.dbHost), zap.Int("queryTimeout", s.queryTimeout))

	// Redis client setup
	redisClient := redis.NewClient(&redis.Options{
		Addr:        s.dbHost,
		DialTimeout: time.Duration(s.queryTimeout) * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.queryTimeout)*time.Second)
	defer cancel()

	// Test the connection
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		s.logger.Fatal("Failed to ping Redis",
			zap.Error(err),
			zap.String("host", s.dbHost),
		)
	}

	return redisClient
}

// Update Health method to use environment variables
func (s *redisStore) Health() *healthinterface.OutboundInterface {
	outbound := healthinterface.OutboundInterface{
		TimeStampUTC:     time.Now().UTC(),
		ConnectionStatus: healthinterface.ConnectionActive,
		ApplicationName:  "Redis",
		URLs:             []string{s.dbHost},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.queryTimeout)*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		s.logger.Fatal("Failed to ping Redis",
			zap.Error(err),
			zap.String("host", s.dbHost),
		)
	}

	return &outbound
}
