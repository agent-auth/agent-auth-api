package migrations

import (
	"context"
	"os"

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func init() {
	// Project collection indexes
	roleIndexes := []mongo.IndexModel{
		{
			Keys: bson.M{"ProjectID": 1},
		},
	}

	migrate.MustRegister(
		// up
		func(ctx context.Context, db *mongo.Database) error {
			// Create Project indexes
			_, err := db.Collection(os.Getenv("DB_ROLES_COLLECTION")).Indexes().CreateOne(ctx, roleIndexes[0])
			if err != nil {
				return err
			}

			return nil
		},

		// down
		func(ctx context.Context, db *mongo.Database) error {
			// Drop Project indexes
			err := db.Collection(os.Getenv("DB_ROLES_COLLECTION")).Indexes().DropOne(ctx, "ProjectID_1")
			if err != nil {
				return err
			}

			return nil
		},
	)
}
