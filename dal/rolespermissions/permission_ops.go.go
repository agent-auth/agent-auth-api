package rolespermissions

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// UpdatePermission updates a specific permission attribute using dot notation
func (p *roles) UpdatePermission(id primitive.ObjectID, path string, value interface{}) error {
	collection := p.db.Database().Collection(p.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(p.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	updates := bson.M{
		path:                    value,
		"updated_timestamp_utc": time.Now(),
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updates},
	)
	if err != nil {
		return fmt.Errorf("failed to update permission attribute: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("permission not found with id: %v", id)
	}

	return nil
}

// RemovePermission removes a specific permission attribute
func (p *roles) RemovePermission(id primitive.ObjectID, path string) error {
	collection := p.db.Database().Collection(p.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(p.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	updates := bson.M{
		"$unset": bson.M{path: ""},
		"$set":   bson.M{"updated_timestamp_utc": time.Now()},
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		updates,
	)
	if err != nil {
		return fmt.Errorf("failed to remove permission attribute: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("permission not found with id: %v", id)
	}

	return nil
}
