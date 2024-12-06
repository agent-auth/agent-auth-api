package redisdb

import (
	"github.com/agent-auth/agent-auth-api/web/interfaces/v1/healthinterface"
	"github.com/go-redis/redis/v8"
)

const (
	// DBConnectionFailed used when failed to create a database client
	DBConnectionFailed = "Database-Connection-Failed"
)

// RedisStore implents the database store
type RedisStore interface {
	Client() *redis.Client
	Health() *healthinterface.OutboundInterface
}
