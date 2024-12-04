package rolespermissions

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

type roles struct {
	db                  connection.MongoStore
	collectionName      string
	queryTimeoutSeconds int
}

// NewRolesDal ...
func NewRolesDal() RolesDal {
	timeoutStr := os.Getenv("DB_QUERY_TIMEOUT_SECONDS")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		timeout = 30 // default timeout
	}

	return &roles{
		db:                  connection.NewMongoStore(),
		collectionName:      os.Getenv("DB_ROLES_COLLECTION"),
		queryTimeoutSeconds: timeout,
	}
}

// Create creates a new role record
func (p *roles) Create(txID string, role *dbmodels.Roles) (*dbmodels.Roles, error) {
	collection := p.db.Database().Collection(p.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(p.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	role.CreatedTimestampUTC = time.Now()
	role.UpdatedTimestampUTC = role.CreatedTimestampUTC

	result, err := collection.InsertOne(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	role.ID = result.InsertedID.(primitive.ObjectID)
	return role, nil
}

// Delete removes a role record by ID
func (p *roles) Delete(id primitive.ObjectID) error {
	collection := p.db.Database().Collection(p.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(p.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete role: %v", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("role not found with id: %v", id)
	}

	return nil
}

// Get retrieves a role by ID
func (p *roles) Get(id primitive.ObjectID) (*dbmodels.Roles, error) {
	collection := p.db.Database().Collection(p.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(p.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	var role dbmodels.Roles
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&role)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return &role, nil
}

// DeleteByProjectID removes all roles for a specific project
func (p *roles) DeleteByProjectID(projectID primitive.ObjectID) error {
	collection := p.db.Database().Collection(p.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(p.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	result, err := collection.DeleteMany(ctx, bson.M{"project_id": projectID})
	if err != nil {
		return fmt.Errorf("failed to delete roles for project: %v", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("no roles found for project id: %v", projectID)
	}

	return nil
}

// GetByProjectID retrieves all roles for a specific project
func (p *roles) GetByProjectID(projectID primitive.ObjectID) ([]*dbmodels.Roles, error) {
	collection := p.db.Database().Collection(p.collectionName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(p.queryTimeoutSeconds)*time.Second,
	)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"project_id": projectID})
	if err != nil {
		return nil, fmt.Errorf("failed to get roles for project: %w", err)
	}
	defer cursor.Close(ctx)

	var roles []*dbmodels.Roles
	if err = cursor.All(ctx, &roles); err != nil {
		return nil, fmt.Errorf("failed to decode roles: %w", err)
	}

	return roles, nil
}
