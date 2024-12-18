package resources_dal

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/agent-auth/agent-auth-api/db/mongodb"
	"github.com/agent-auth/common-lib/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type resources struct {
	db                  *mongo.Database
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
		db:                  mongodb.NewMongoClient(),
		collectionName:      os.Getenv("DB_RESOURCES_COLLECTION"),
		queryTimeoutSeconds: timeout,
	}
}

// Create creates a new resource
func (r *resources) Create(resource *models.Resource) (*models.Resource, error) {
	if resource == nil {
		return nil, fmt.Errorf("resource cannot be nil")
	}
	collection := r.db.Collection(r.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(r.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	// Set timestamps and ensure fields are initialized
	now := time.Now().UTC()
	resource.CreatedTimestampUTC = now
	resource.UpdatedTimestampUTC = now
	resource.Deleted = false // Ensure new resources aren't created as deleted

	result, err := collection.InsertOne(ctx, resource)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	resource.ID = result.InsertedID.(bson.ObjectID)
	return resource, nil
}

// Update updates a resource's mutable fields
func (r *resources) Update(resource *models.Resource) error {
	if resource == nil {
		return fmt.Errorf("resource cannot be nil")
	}
	if resource.ID.IsZero() {
		return fmt.Errorf("resource ID cannot be empty")
	}
	collection := r.db.Collection(r.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(r.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	// Only update mutable fields
	updateDoc := bson.M{
		"$set": bson.M{
			"Description":         resource.Description,
			"Actions":             resource.Actions,
			"UpdatedTimestampUTC": time.Now().UTC(),
		},
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{
			"_id":     resource.ID,
			"Deleted": bson.M{"$ne": true}, // Don't update deleted resources
		},
		updateDoc,
	)
	if err != nil {
		return fmt.Errorf("failed to update resource: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("resource not found or already deleted with id: %v", resource.ID)
	}

	return nil
}

// GetByID retrieves a resource by its ID
func (r *resources) GetByID(id bson.ObjectID) (*models.Resource, error) {
	if id.IsZero() {
		return nil, fmt.Errorf("invalid resource ID")
	}
	collection := r.db.Collection(r.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(r.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	filter := bson.M{
		"_id":     id,
		"Deleted": bson.M{"$ne": true},
	}

	var resource models.Resource
	if err := collection.FindOne(ctx, filter).Decode(&resource); err != nil {
		return nil, fmt.Errorf("failed to find resource: %w", err)
	}

	return &resource, nil
}

// Delete soft-deletes a resource by ID
func (r *resources) Delete(id bson.ObjectID) error {
	if id.IsZero() {
		return fmt.Errorf("invalid resource ID")
	}
	collection := r.db.Collection(r.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(r.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"Deleted":             true,
			"UpdatedTimestampUTC": time.Now().UTC(),
		},
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{
			"_id":     id,
			"Deleted": bson.M{"$ne": true}, // Prevent re-deleting
		},
		update,
	)
	if err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("resource not found or already deleted with id: %v", id)
	}

	return nil
}

// GetByProjectID retrieves all non-deleted resources for a given project ID
func (r *resources) GetByProjectID(projectID bson.ObjectID) ([]*models.Resource, error) {
	if projectID.IsZero() {
		return nil, fmt.Errorf("invalid project ID")
	}
	collection := r.db.Collection(r.collectionName)
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

	var resources []*models.Resource
	if err = cursor.All(ctx, &resources); err != nil {
		return nil, fmt.Errorf("failed to decode resources: %w", err)
	}

	return resources, nil
}

// GetByURNAndProjectID retrieves a resource by URN and project ID
func (r *resources) GetByURNAndProjectID(urn string, projectID bson.ObjectID) (*models.Resource, error) {
	if projectID.IsZero() {
		return nil, fmt.Errorf("invalid project ID")
	}
	collection := r.db.Collection(r.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(r.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	var resource models.Resource

	filter := bson.M{
		"URN":       urn,
		"ProjectID": projectID,
	}

	err := collection.FindOne(ctx, filter).Decode(&resource)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &resource, nil
}
