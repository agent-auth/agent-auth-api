package projectdal

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/agent-auth/agent-auth-api/database/connection"
	"github.com/agent-auth/agent-auth-api/database/dbmodels"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type workspaces struct {
	db                  connection.MongoStore
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
		db:                  connection.NewMongoStore(),
		collectionName:      os.Getenv("DB_WORKSPACES_COLLECTION"),
		queryTimeoutSeconds: timeout,
	}
}

// Delete removes a workspace by ID
func (w *workspaces) Delete(id primitive.ObjectID) error {
	collection := w.db.Database().Collection(w.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(w.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete workspace: %v", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("workspace not found with id: %v", id)
	}

	return nil
}

// List retrieves workspaces with pagination
func (w *workspaces) List(skip, limit int64) ([]*dbmodels.Workspace, error) {
	collection := w.db.Database().Collection(w.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(w.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.M{"CreatedTimestampUTC": -1}) // Sort by creation date, newest first

	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list workspaces: %v", err)
	}
	defer cursor.Close(ctx)

	var workspaces []*dbmodels.Workspace
	if err = cursor.All(ctx, &workspaces); err != nil {
		return nil, fmt.Errorf("failed to decode workspaces: %v", err)
	}

	return workspaces, nil
}

// GetBySlug retrieves a workspace by its slug
func (w *workspaces) GetBySlug(slug string) (*dbmodels.Workspace, error) {
	collection := w.db.Database().Collection(w.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(w.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	var workspace dbmodels.Workspace
	if err := collection.FindOne(ctx, bson.M{"Slug": slug}).Decode(&workspace); err != nil {
		return nil, fmt.Errorf("failed to find workspace by slug: %v", err)
	}

	return &workspace, nil
}

// GetByOwnerID retrieves all workspaces owned by a specific user
func (w *workspaces) GetByOwnerID(ownerID primitive.ObjectID) ([]*dbmodels.Workspace, error) {
	collection := w.db.Database().Collection(w.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(w.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"OwnerID": ownerID})
	if err != nil {
		return nil, fmt.Errorf("failed to find workspaces by owner: %v", err)
	}
	defer cursor.Close(ctx)

	var workspaces []*dbmodels.Workspace
	if err = cursor.All(ctx, &workspaces); err != nil {
		return nil, fmt.Errorf("failed to decode workspaces: %v", err)
	}

	return workspaces, nil
}

// AddMember adds a member to a workspace
func (w *workspaces) AddMember(workspaceID, memberID primitive.ObjectID) error {
	collection := w.db.Database().Collection(w.collectionName)
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
func (w *workspaces) RemoveMember(workspaceID, memberID primitive.ObjectID) error {
	collection := w.db.Database().Collection(w.collectionName)
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
func (w *workspaces) Update(workspace *dbmodels.Workspace) error {
	collection := w.db.Database().Collection(w.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(w.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	// Fields that can be updated
	updateDoc := bson.M{
		"Name":                 workspace.Name,
		"Description":          workspace.Description,
		"PrimaryAuthEnabled":   workspace.PrimaryAuthEnabled,
		"SecondaryAuthEnabled": workspace.SecondaryAuthEnabled,
		"RequireMFA":           workspace.RequireMFA,
		"AllowSecondaryAuth":   workspace.AllowSecondaryAuth,
		"AllowEmailInvites":    workspace.AllowEmailInvites,
		"AllowedDomains":       workspace.AllowedDomains,
		"UpdatedTimestampUTC":  time.Now(),
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
func (w *workspaces) Create(txID string, workspace *dbmodels.Workspace) (*dbmodels.Workspace, error) {
	collection := w.db.Database().Collection(w.collectionName)
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
	workspace.ID = result.InsertedID.(primitive.ObjectID)

	return workspace, nil
}

// GetByID retrieves a workspace by its ID
func (w *workspaces) GetByID(id primitive.ObjectID) (*dbmodels.Workspace, error) {
	collection := w.db.Database().Collection(w.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(w.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	var workspace dbmodels.Workspace
	if err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&workspace); err != nil {
		return nil, fmt.Errorf("failed to find workspace: %w", err)
	}

	return &workspace, nil
}
