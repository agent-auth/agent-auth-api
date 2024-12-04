package resources_dal

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/agent-auth/agent-auth-api/database/connection"
	"github.com/agent-auth/agent-auth-api/database/dbmodels"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

type resources struct {
	db                  connection.MongoStore
	collectionName      string
	queryTimeoutSeconds int
}

// NewResourcesDal creates a new ResourcesDal instance
func NewResourcesDal() ResourcesDal {
	timeoutStr := os.Getenv("DB_QUERY_TIMEOUT_SECONDS")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		timeout = 30 // default timeout
	}

	return &resources{
		db:                  connection.NewMongoStore(),
		collectionName:      os.Getenv("DB_RESOURCES_COLLECTION"),
		queryTimeoutSeconds: timeout,
	}
}

// Create creates a new resource
func (r *resources) Create(resource *dbmodels.Resource) (*dbmodels.Resource, error) {
	if resource == nil {
		return nil, fmt.Errorf("resource cannot be nil")
	}
	collection := r.db.Database().Collection(r.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(r.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	// Set timestamps
	now := time.Now()
	resource.CreatedTimestampUTC = now
	resource.UpdatedTimestampUTC = now

	result, err := collection.InsertOne(ctx, resource)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	resource.ID = result.InsertedID.(primitive.ObjectID)
	return resource, nil
}

// Update updates a resource's mutable fields
func (r *resources) Update(resource *dbmodels.Resource) error {
	if resource == nil {
		return fmt.Errorf("resource cannot be nil")
	}
	if resource.ID.IsZero() {
		return fmt.Errorf("resource ID cannot be empty")
	}
	collection := r.db.Database().Collection(r.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(r.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	updateDoc := bson.M{
		"Description":         resource.Description,
		"Actions":             resource.Actions,
		"UpdatedTimestampUTC": time.Now(),
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": resource.ID},
		bson.M{"$set": updateDoc},
	)
	if err != nil {
		return fmt.Errorf("failed to update resource: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("resource not found with id: %v", resource.ID)
	}

	return nil
}

// GetByID retrieves a resource by its ID
func (r *resources) GetByID(id primitive.ObjectID) (*dbmodels.Resource, error) {
	collection := r.db.Database().Collection(r.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(r.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	filter := bson.M{
		"_id":     id,
		"Deleted": bson.M{"$ne": true},
	}

	var resource dbmodels.Resource
	if err := collection.FindOne(ctx, filter).Decode(&resource); err != nil {
		return nil, fmt.Errorf("failed to find resource: %w", err)
	}

	return &resource, nil
}

// Delete soft-deletes a resource by ID
func (r *resources) Delete(id primitive.ObjectID) error {
	collection := r.db.Database().Collection(r.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(r.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"Deleted":             true,
			"UpdatedTimestampUTC": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("resource not found with id: %v", id)
	}

	return nil
}

// GetByProjectID retrieves all non-deleted resources for a given project ID
func (r *resources) GetByProjectID(projectID primitive.ObjectID) ([]*dbmodels.Resource, error) {
	collection := r.db.Database().Collection(r.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(r.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	filter := bson.M{
		"ProjectID": projectID,
		"Deleted":   bson.M{"$ne": true},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find resources by project ID: %w", err)
	}
	defer cursor.Close(ctx)

	var resources []*dbmodels.Resource
	if err = cursor.All(ctx, &resources); err != nil {
		return nil, fmt.Errorf("failed to decode resources: %w", err)
	}

	return resources, nil
}
