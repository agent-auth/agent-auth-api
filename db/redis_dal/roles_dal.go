package redis_dal

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/agent-auth/agent-auth-api/db/mongodb"
	"github.com/agent-auth/agent-auth-api/db/redisdb"
	"github.com/agent-auth/common-lib/pkg/logger"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	cursor, err := r.mongo.Collection(r.collectionName).Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to fetch initial data: %v", err)
	}
	defer cursor.Close(ctx)

	var syncErrors []error
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			r.logger.Error("Failed to decode document during initial sync", zap.Error(err))
			syncErrors = append(syncErrors, err)
			continue
		}

		objectID, ok := doc["_id"].(primitive.ObjectID)
		if !ok {
			r.logger.Error("Document _id is not an ObjectID")
			continue
		}
		docID := objectID.Hex()

		// Convert BSON to JSON instead of using bson.Marshal
		docJSON, err := bson.MarshalExtJSON(doc, true, true)
		if err != nil {
			r.logger.Error("Failed to marshal document to JSON during initial sync", zap.Error(err))
			syncErrors = append(syncErrors, err)
			continue
		}

		r.logger.Info("Storing document in Redis during initial sync", zap.String("collection", r.collectionName), zap.String("id", docID))

		err = r.redis.Set(context.Background(), fmt.Sprintf("%s:%s", r.collectionName, docID), docJSON, 0).Err()
		if err != nil {
			r.logger.Error("Failed to store document in Redis during initial sync", zap.Error(err))
			syncErrors = append(syncErrors, err)
		}
	}

	// Check for cursor iteration errors
	if err := cursor.Err(); err != nil {
		return fmt.Errorf("cursor iteration error: %v", err)
	}

	if len(syncErrors) > 0 {
		return fmt.Errorf("encountered %d errors during sync", len(syncErrors))
	}

	r.logger.Info("Initial redis sync completed successfully")

	return nil
}

func (r *redis_roles_dal) SyncRolesCollection(ctx context.Context) {
	if err := r.InitialSync(); err != nil {
		r.logger.Error("Failed to perform initial sync", zap.Error(err))
		return // Consider if you want to continue after initial sync failure
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

			filter := bson.M{
				"UpdatedTimestampUTC": bson.M{"$gt": lastSync},
			}

			cursor, err := r.mongo.Collection(r.collectionName).Find(syncCtx, filter)
			if err != nil {
				r.logger.Error("Failed to poll for changes", zap.Error(err))
				cancel()
				continue
			}

			// Process documents
			for cursor.Next(syncCtx) {
				var doc bson.M
				if err := cursor.Decode(&doc); err != nil {
					r.logger.Error("Failed to decode document during polling", zap.Error(err))
					continue
				}

				objectID, ok := doc["_id"].(primitive.ObjectID)
				if !ok {
					r.logger.Error("Document _id is not an ObjectID")
					continue
				}
				docID := objectID.Hex()

				docJSON, err := bson.MarshalExtJSON(doc, true, true)
				if err != nil {
					r.logger.Error("Failed to marshal document to JSON during polling", zap.Error(err))
					continue
				}

				r.logger.Info("Storing document in Redis during polling", zap.String("collection", r.collectionName), zap.String("id", docID))

				err = r.redis.Set(syncCtx, fmt.Sprintf("%s:%s", r.collectionName, docID), docJSON, 0).Err()
				if err != nil {
					r.logger.Error("Failed to store document in Redis during polling", zap.Error(err))
				}
			}

			// Check for cursor errors
			if err := cursor.Err(); err != nil {
				r.logger.Error("Cursor iteration error", zap.Error(err))
			}

			cursor.Close(syncCtx)
			cancel()
			lastSync = time.Now().UTC()
		}
	}
}
