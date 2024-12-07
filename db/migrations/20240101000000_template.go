package migrations

import (
	"context"

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func init() {
	migrate.MustRegister(
		func(ctx context.Context, db *mongo.Database) error {

			return nil
		}, func(ctx context.Context, db *mongo.Database) error {
			return nil
		})
}
