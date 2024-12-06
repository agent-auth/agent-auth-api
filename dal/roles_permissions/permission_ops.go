package roles_permissions_dal

import (
	"context"
	"fmt"
	"time"

	"github.com/agent-auth/agent-auth-api/database/dbmodels"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// UpdatePermission updates a specific permission attribute using dot notation
func (p *roles) UpdatePermission(id primitive.ObjectID, resource string, actions []dbmodels.Action) error {
	collection := p.db.Database().Collection(p.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(p.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	var result *mongo.UpdateResult
	var err error

	if len(actions) == 0 {
		// Delete the resource key path if actions are empty
		result, err = collection.UpdateOne(
			ctx,
			bson.M{"_id": id},
			bson.M{
				"$unset": bson.M{
					fmt.Sprintf("Permissions.%s", resource): "",
				},
				"$set": bson.M{
					"UpdatedTimestampUTC": time.Now(),
				},
			},
		)
	} else {
		// Update the actions for the specific resource
		result, err = collection.UpdateOne(
			ctx,
			bson.M{"_id": id},
			bson.M{
				"$set": bson.M{
					fmt.Sprintf("Permissions.%s.Actions", resource): actions,
					"UpdatedTimestampUTC":                           time.Now(),
				},
			},
		)
	}

	if err != nil {
		return fmt.Errorf("failed to update permission actions: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("resource '%s' not found in permissions for role id: %v", resource, id)
	}

	return nil
}
