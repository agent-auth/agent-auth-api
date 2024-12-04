package projects_dal

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

type projects struct {
	db                  connection.MongoStore
	collectionName      string
	queryTimeoutSeconds int
}

// NewProjectsDal creates a new ProjectsDal instance
func NewProjectsDal() ProjectsDal {
	timeoutStr := os.Getenv("DB_QUERY_TIMEOUT_SECONDS")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		timeout = 30 // default timeout
	}

	return &projects{
		db:                  connection.NewMongoStore(),
		collectionName:      os.Getenv("DB_PROJECTS_COLLECTION"),
		queryTimeoutSeconds: timeout,
	}
}

// Create creates a new project
func (p *projects) Create(project *dbmodels.Project) (*dbmodels.Project, error) {
	if project == nil {
		return nil, fmt.Errorf("project cannot be nil")
	}
	collection := p.db.Database().Collection(p.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(p.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	// Set timestamps
	now := time.Now()
	project.CreatedTimestampUTC = now
	project.UpdatedTimestampUTC = now

	// Validate the project
	if err := project.Validate(); err != nil {
		return nil, fmt.Errorf("invalid project data: %w", err)
	}

	result, err := collection.InsertOne(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	project.ID = result.InsertedID.(primitive.ObjectID)
	return project, nil
}

// Update updates a project's mutable fields
func (p *projects) Update(project *dbmodels.Project) error {
	if project == nil {
		return fmt.Errorf("project cannot be nil")
	}
	if project.ID.IsZero() {
		return fmt.Errorf("project ID cannot be empty")
	}
	collection := p.db.Database().Collection(p.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(p.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	// Only include mutable fields in the update
	updateDoc := bson.M{
		"Name":                project.Name,
		"Description":         project.Description,
		"UpdatedTimestampUTC": time.Now(),
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": project.ID},
		bson.M{"$set": updateDoc},
	)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("project not found with id: %v", project.ID)
	}

	return nil
}

// List retrieves projects for a workspace with pagination
func (p *projects) List(workspaceID primitive.ObjectID, skip, limit int64) ([]*dbmodels.Project, error) {
	if workspaceID.IsZero() {
		return nil, fmt.Errorf("workspace ID cannot be empty")
	}
	if limit <= 0 {
		limit = 10 // Set a default limit
	}
	if skip < 0 {
		skip = 0
	}
	collection := p.db.Database().Collection(p.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(p.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.M{"CreatedTimestampUTC": -1})

	filter := bson.M{
		"WorkspaceID": workspaceID,
		"Deleted":     bson.M{"$ne": true},
	}

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	defer cursor.Close(ctx)

	var projects []*dbmodels.Project
	if err = cursor.All(ctx, &projects); err != nil {
		return nil, fmt.Errorf("failed to decode projects: %w", err)
	}

	return projects, nil
}

// GetByID retrieves a project by its ID
func (p *projects) GetByID(id primitive.ObjectID) (*dbmodels.Project, error) {
	collection := p.db.Database().Collection(p.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(p.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	filter := bson.M{
		"_id":     id,
		"Deleted": bson.M{"$ne": true},
	}

	var project dbmodels.Project
	if err := collection.FindOne(ctx, filter).Decode(&project); err != nil {
		return nil, fmt.Errorf("failed to find project: %w", err)
	}

	return &project, nil
}

// Delete soft-deletes a project by ID
func (p *projects) Delete(id primitive.ObjectID) error {
	collection := p.db.Database().Collection(p.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(p.queryTimeoutSeconds)*time.Second,
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
		return fmt.Errorf("failed to delete project: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("project not found with id: %v", id)
	}

	return nil
}

// GetBySlug retrieves a project by its slug within a workspace
func (p *projects) GetBySlug(workspaceID primitive.ObjectID, slug string) (*dbmodels.Project, error) {
	collection := p.db.Database().Collection(p.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(p.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	var project dbmodels.Project
	filter := bson.M{
		"WorkspaceID": workspaceID,
		"Slug":        slug,
		"Deleted":     bson.M{"$ne": true},
	}

	if err := collection.FindOne(ctx, filter).Decode(&project); err != nil {
		return nil, fmt.Errorf("failed to find project by slug: %w", err)
	}

	return &project, nil
}

// GetByOwnerID retrieves all projects owned by a specific user
func (p *projects) GetByOwnerID(ownerID string) ([]*dbmodels.Project, error) {
	collection := p.db.Database().Collection(p.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(p.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	opts := options.Find().SetSort(bson.M{"CreatedTimestampUTC": -1})
	cursor, err := collection.Find(ctx, bson.M{"OwnerID": ownerID, "Deleted": bson.M{"$ne": true}}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find projects by owner: %w", err)
	}
	defer cursor.Close(ctx)

	var projects []*dbmodels.Project
	if err = cursor.All(ctx, &projects); err != nil {
		return nil, fmt.Errorf("failed to decode projects: %w", err)
	}

	return projects, nil
}

// AddMember adds a member to a project
func (p *projects) AddMember(projectID, memberID primitive.ObjectID) error {
	if projectID.IsZero() || memberID.IsZero() {
		return fmt.Errorf("project ID and member ID cannot be empty")
	}
	collection := p.db.Database().Collection(p.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(p.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	update := bson.M{
		"$addToSet": bson.M{"Members": memberID},
		"$set": bson.M{
			"UpdatedTimestampUTC": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": projectID}, update)
	if err != nil {
		return fmt.Errorf("failed to add member to project: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("project not found with id: %v", projectID)
	}

	return nil
}

// RemoveMember removes a member from a project
func (p *projects) RemoveMember(projectID, memberID primitive.ObjectID) error {
	collection := p.db.Database().Collection(p.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(p.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	update := bson.M{
		"$pull": bson.M{"Members": memberID},
		"$set": bson.M{
			"UpdatedTimestampUTC": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": projectID}, update)
	if err != nil {
		return fmt.Errorf("failed to remove member from project: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("project not found with id: %v", projectID)
	}

	return nil
}

// IsMember checks if the given email is a member of the specified project
func (p *projects) IsMember(projectID, email string) (bool, error) {
	if projectID == "" {
		return false, fmt.Errorf("project ID cannot be empty")
	}
	if email == "" {
		return false, fmt.Errorf("email cannot be empty")
	}
	collection := p.db.Database().Collection(p.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(p.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	// Find project members where email matches
	count, err := collection.CountDocuments(ctx, bson.M{
		"_id":     projectID,
		"Members": email,
	})
	if err != nil {
		return false, fmt.Errorf("error checking project membership: %w", err)
	}

	return count > 0, nil
}
