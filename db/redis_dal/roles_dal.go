package redis_dal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/agent-auth/agent-auth-api/db/mongodb"
	"github.com/agent-auth/agent-auth-api/db/redisdb"
	"github.com/agent-auth/common-lib/models"
	"github.com/agent-auth/common-lib/pkg/logger"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

// RedisRolesDal ...
type redis_roles_dal struct {
	logger         *zap.Logger
	mongo          *mongo.Database
	redis          *redis.Client
	collectionName string
	timeoutSeconds int
	syncInterval   int
}

// NewRedisRolesDal returns new instance of datastore
func NewRedisRolesDal() *redis_roles_dal {
	l := logger.NewLogger()

	// redis environment variables
	redisQueryTimeout, err := strconv.Atoi(os.Getenv("REDIS_QUERY_TIMEOUT_SECONDS"))
	if err != nil {
		redisQueryTimeout = 30
	}

	redisSyncInterval, err := strconv.Atoi(os.Getenv("REDIS_SYNC_INTERVAL"))
	if err != nil {
		redisSyncInterval = 10
	}

	rolesCollectionName := os.Getenv("DB_ROLES_COLLECTION")
	if rolesCollectionName == "" {
		l.Fatal("DB_ROLES_COLLECTION is not set")
	}

	return &redis_roles_dal{
		logger:         logger.NewLogger(),
		mongo:          mongodb.NewMongoClient(),
		redis:          redisdb.NewRedisClient(),
		collectionName: rolesCollectionName,
		timeoutSeconds: redisQueryTimeout,
		syncInterval:   redisSyncInterval,
	}
}

func (r *redis_roles_dal) InitialSync() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeoutSeconds)*time.Second)
	defer cancel()

	var roles []models.Roles
	cursor, err := r.mongo.Collection(r.collectionName).Find(ctx, bson.M{
		"Deleted": bson.M{"$ne": true},
	})
	if err != nil {
		return fmt.Errorf("failed to fetch initial data: %v", err)
	}
	defer cursor.Close(ctx)

	// Collect all roles
	for cursor.Next(ctx) {
		var role models.Roles
		if err := cursor.Decode(&role); err != nil {
			r.logger.Error("Failed to decode document during initial sync", zap.Error(err))
			continue
		}
		roles = append(roles, role)
	}

	projectRoles := r.transformRolesToProjectMap(roles)
	if err := r.storeProjectRolesInRedis(ctx, projectRoles); err != nil {
		return fmt.Errorf("failed to store roles in Redis: %v", err)
	}

	r.logger.Info("Initial redis sync completed successfully")
	return nil
}

func (r *redis_roles_dal) SyncRolesCollection(ctx context.Context) {
	if err := r.InitialSync(); err != nil {
		r.logger.Error("Failed to perform initial sync", zap.Error(err))
		return
	}

	ticker := time.NewTicker(time.Duration(r.syncInterval) * time.Second)
	defer ticker.Stop()

	lastSync := time.Now().UTC()

	for {
		select {
		case <-ctx.Done():
			r.logger.Info("Stopping roles collection sync due to context cancellation")
			return
		case <-ticker.C:
			syncCtx, cancel := context.WithTimeout(ctx, time.Duration(r.timeoutSeconds)*time.Second)
			// Query for updated roles
			filter := bson.M{
				"UpdatedTimestampUTC": bson.M{"$gt": lastSync},
				"Deleted":             bson.M{"$ne": true},
			}

			var roles []models.Roles
			cursor, err := r.mongo.Collection(r.collectionName).Find(syncCtx, filter)
			if err != nil {
				r.logger.Error("Failed to poll for changes", zap.Error(err))
				cancel()
				continue
			}

			// Collect all roles
			for cursor.Next(syncCtx) {
				var role models.Roles
				if err := cursor.Decode(&role); err != nil {
					r.logger.Error("Failed to decode document during polling", zap.Error(err))
					continue
				}
				roles = append(roles, role)
			}
			cursor.Close(syncCtx)

			projectRoles := r.transformRolesToProjectMap(roles)
			if err := r.storeProjectRolesInRedis(syncCtx, projectRoles); err != nil {
				r.logger.Error("Failed to store roles in Redis", zap.Error(err))
			} else {
				lastSync = time.Now().UTC()
			}

			cancel()

			r.logger.Info("Roles collection sync completed successfully, looking for more changes")
		}
	}
}

// New helper function to handle the transformation
func (r *redis_roles_dal) transformRolesToProjectMap(roles []models.Roles) map[string]map[string]map[string]models.Permission {
	projectRoles := make(map[string]map[string]map[string]models.Permission)

	for _, role := range roles {
		projectID := role.ProjectID.Hex()

		// Initialize project map if it doesn't exist
		if _, exists := projectRoles[projectID]; !exists {
			projectRoles[projectID] = make(map[string]map[string]models.Permission)
		}

		// Initialize role map if it doesn't exist
		if _, exists := projectRoles[projectID][role.Role]; !exists {
			projectRoles[projectID][role.Role] = make(map[string]models.Permission)
		}

		// Add permissions to the role
		for urn, permission := range role.Permissions {
			projectRoles[projectID][role.Role][urn] = permission
		}
	}

	return projectRoles
}

// Helper function to store transformed data in Redis
func (r *redis_roles_dal) storeProjectRolesInRedis(ctx context.Context, projectRoles map[string]map[string]map[string]models.Permission) error {
	for projectID, roleData := range projectRoles {
		transformedData := struct {
			ProjectID   string                                  `json:"project_id"`
			Permissions map[string]map[string]models.Permission `json:"permissions"`
		}{
			ProjectID:   projectID,
			Permissions: roleData,
		}

		jsonData, err := json.Marshal(transformedData)
		if err != nil {
			r.logger.Error("Failed to marshal transformed data", zap.Error(err))
			return err
		}

		r.logger.Info("Storing transformed data in Redis", zap.String("projectID", fmt.Sprintf("roles:%s", projectID)))

		err = r.redis.Set(ctx, fmt.Sprintf("roles:%s", projectID), jsonData, 0).Err()
		if err != nil {
			r.logger.Error("Failed to store transformed data in Redis",
				zap.String("projectID", projectID),
				zap.Error(err))
			return err
		}
	}
	return nil
}
