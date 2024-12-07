package redisdb

import (
	"os"
	"strconv"

	"github.com/agent-auth/common-lib/pkg/logger"
	"github.com/agent-auth/common-lib/pkg/redis_client"
	"github.com/go-redis/redis/v8"
)

// NewRedisClient returns new instance of datastore
func NewRedisClient() *redis.Client {
	l := logger.NewLogger()

	// redis environment variables
	redisQueryTimeout, err := strconv.Atoi(os.Getenv("REDIS_QUERY_TIMEOUT_SECONDS"))
	if err != nil {
		redisQueryTimeout = 30
	}

	redisHost := os.Getenv("REDIS_URI")
	if redisHost == "" {
		l.Fatal("REDIS_URI is not set")
	}

	return redis_client.NewRedisStore(redisHost, redisQueryTimeout).Client()
}
