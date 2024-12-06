package main

import (
	"context"

	"github.com/agent-auth/agent-auth-api/pkg/redisdb"
)

func main() {
	// Create a channel to handle graceful shutdown

	// Initialize the Redis roles DAL
	rolesDal := redisdb.NewRedisRolesDal()

	rolesDal.SyncRolesCollection(context.Background())

}
