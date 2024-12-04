package migrations

import (
	"context"
	"os"

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func init() {

	// Project collection indexes
	projectIndexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"Slug": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"WorkspaceID": 1},
		},
		{
			Keys: bson.M{"OwnerID": 1},
		},
	}

	// Workspace collection indexes
	workspaceIndexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"Slug": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"OwnerID": 1},
		},
	}

	// Register migrations
	migrate.Register(

		// up
		func(ctx context.Context, db *mongo.Database) error {
			// Create Project indexes
			_, err := db.Collection(os.Getenv("DB_PROJECTS_COLLECTION")).Indexes().CreateMany(ctx, projectIndexes)
			if err != nil {
				return err
			}

			// Create Workspace indexes
			_, err = db.Collection(os.Getenv("DB_WORKSPACES_COLLECTION")).Indexes().CreateMany(ctx, workspaceIndexes)
			return err
		},

		// down
		func(ctx context.Context, db *mongo.Database) error {
			// Drop Project indexes
			err := db.Collection(os.Getenv("DB_PROJECTS_COLLECTION")).Indexes().DropOne(ctx, "Slug_1")
			if err != nil {
				return err
			}
			err = db.Collection(os.Getenv("DB_PROJECTS_COLLECTION")).Indexes().DropOne(ctx, "WorkspaceID_1")
			if err != nil {
				return err
			}
			err = db.Collection(os.Getenv("DB_PROJECTS_COLLECTION")).Indexes().DropOne(ctx, "OwnerID_1")
			if err != nil {
				return err
			}

			// Drop Workspace indexes
			err = db.Collection(os.Getenv("DB_WORKSPACES_COLLECTION")).Indexes().DropOne(ctx, "Slug_1")
			if err != nil {
				return err
			}
			err = db.Collection(os.Getenv("DB_WORKSPACES_COLLECTION")).Indexes().DropOne(ctx, "OwnerID_1")
			if err != nil {
				return err
			}
			return nil
		})
}
