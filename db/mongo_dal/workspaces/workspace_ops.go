package workspaces_dal

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
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type workspaces struct {
	db                  *mongo.Database
	collectionName      string
	queryTimeoutSeconds int
}

// NewWorkspaceDal ...
func NewWorkspaceDal() WorkspaceDal {
	timeoutStr := os.Getenv("DB_QUERY_TIMEOUT_SECONDS")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		timeout = 30 // default timeout
	}

	return &workspaces{
		db:                  mongodb.NewMongoClient(),
		collectionName:      os.Getenv("DB_WORKSPACES_COLLECTION"),
		queryTimeoutSeconds: timeout,
	}
}

// Delete soft deletes a workspace by setting its deleted flag to true
func (w *workspaces) Delete(id bson.ObjectID) error {
	collection := w.db.Collection(w.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(w.queryTimeoutSeconds)*time.Second,
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
		return fmt.Errorf("failed to soft delete workspace: %v", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("workspace not found with id: %v", id)
	}

	return nil
}

// List retrieves workspaces with pagination
func (w *workspaces) List(skip, limit int64) ([]*models.Workspace, error) {
	collection := w.db.Collection(w.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(w.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.M{"CreatedTimestampUTC": -1}) // Sort by creation date, newest first

	cursor, err := collection.Find(ctx, bson.M{"Deleted": false}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list workspaces: %v", err)
	}
	defer cursor.Close(ctx)

	var workspaces []*models.Workspace
	if err = cursor.All(ctx, &workspaces); err != nil {
		return nil, fmt.Errorf("failed to decode workspaces: %v", err)
	}

	return workspaces, nil
}

// GetBySlug retrieves a workspace by its slug
func (w *workspaces) GetBySlug(slug string) (*models.Workspace, error) {
	collection := w.db.Collection(w.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(w.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	var workspace models.Workspace
	if err := collection.FindOne(ctx, bson.M{"Slug": slug, "Deleted": false}).Decode(&workspace); err != nil {
		return nil, fmt.Errorf("failed to find workspace by slug: %v", err)
	}

	return &workspace, nil
}

// GetByOwnerID retrieves all workspaces owned by a specific user
func (w *workspaces) GetByOwnerID(ownerID bson.ObjectID) ([]*models.Workspace, error) {
	collection := w.db.Collection(w.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(w.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"OwnerID": ownerID, "Deleted": false})
	if err != nil {
		return nil, fmt.Errorf("failed to find workspaces by owner: %v", err)
	}
	defer cursor.Close(ctx)

	var workspaces []*models.Workspace
	if err = cursor.All(ctx, &workspaces); err != nil {
		return nil, fmt.Errorf("failed to decode workspaces: %v", err)
	}

	return workspaces, nil
}

// AddMember adds a member to a workspace
func (w *workspaces) AddMember(workspaceID string, memberID string) error {
	collection := w.db.Collection(w.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(w.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	update := bson.M{
		"$addToSet": bson.M{"Members": memberID},
		"$set": bson.M{
			"UpdatedTimestampUTC": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": workspaceID}, update)
	if err != nil {
		return fmt.Errorf("failed to add member to workspace: %v", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("workspace not found with id: %v", workspaceID)
	}

	return nil
}

// RemoveMember removes a member from a workspace
func (w *workspaces) RemoveMember(workspaceID string, memberID string) error {
	collection := w.db.Collection(w.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(w.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	update := bson.M{
		"$pull": bson.M{"Members": memberID},
		"$set": bson.M{
			"UpdatedTimestampUTC": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": workspaceID}, update)
	if err != nil {
		return fmt.Errorf("failed to remove member from workspace: %v", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("workspace not found with id: %v", workspaceID)
	}

	return nil
}

// Update method needs to be modified to handle more fields
func (w *workspaces) Update(workspace *models.Workspace) error {
	collection := w.db.Collection(w.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(w.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	// Fields that can be updated
	updateDoc := bson.M{
		"Name":                workspace.Name,
		"Description":         workspace.Description,
		"UpdatedTimestampUTC": time.Now(),
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": workspace.ID},
		bson.M{"$set": updateDoc},
	)
	if err != nil {
		return fmt.Errorf("failed to update workspace: %v", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("workspace not found with id: %v", workspace.ID)
	}

	return nil
}

// Create creates a new workspace
func (w *workspaces) Create(workspace *models.Workspace) (*models.Workspace, error) {
	collection := w.db.Collection(w.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(w.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	workspace.CreatedTimestampUTC = time.Now()
	workspace.UpdatedTimestampUTC = workspace.CreatedTimestampUTC

	result, err := collection.InsertOne(ctx, workspace)
	if err != nil {
		return nil, fmt.Errorf("failed to create workspace: %w", err)
	}

	// Set the ID from the insertion result
	workspace.ID = result.InsertedID.(bson.ObjectID)

	return workspace, nil
}

// GetByID retrieves a workspace by its ID
func (w *workspaces) GetByID(id bson.ObjectID) (*models.Workspace, error) {
	collection := w.db.Collection(w.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(w.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	var workspace models.Workspace
	if err := collection.FindOne(ctx, bson.M{"_id": id, "Deleted": false}).Decode(&workspace); err != nil {
		return nil, fmt.Errorf("failed to find workspace: %w", err)
	}

	return &workspace, nil
}

// IsMember checks if the given email is a member of the specified project
func (w *workspaces) IsMember(workspaceID, email string) (bool, error) {
	collection := w.db.Collection(w.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(w.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	// Find project members where email matches
	count, err := collection.CountDocuments(ctx, bson.M{
		"_id":     workspaceID,
		"Members": email,
	})
	if err != nil {
		return false, fmt.Errorf("error checking workspace membership: %w", err)
	}

	return count > 0, nil
}
