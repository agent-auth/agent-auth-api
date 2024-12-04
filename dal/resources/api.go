package resources_dal

import (
	"github.com/agent-auth/agent-auth-api/database/dbmodels"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ResourcesDal defines the interface for resource database operations
type ResourcesDal interface {
	Create(resource *dbmodels.Resource) (*dbmodels.Resource, error)
	Update(resource *dbmodels.Resource) error
	GetByID(id primitive.ObjectID) (*dbmodels.Resource, error)
	Delete(id primitive.ObjectID) error
	GetByProjectID(projectID primitive.ObjectID) ([]*dbmodels.Resource, error)
}
