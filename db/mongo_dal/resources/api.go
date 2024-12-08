package resources_dal

import (
	"github.com/agent-auth/common-lib/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ResourcesDal defines the interface for resource database operations
type ResourcesDal interface {
	Create(resource *models.Resource) (*models.Resource, error)
	Update(resource *models.Resource) error
	GetByID(id bson.ObjectID) (*models.Resource, error)
	Delete(id bson.ObjectID) error
	GetByProjectID(projectID bson.ObjectID) ([]*models.Resource, error)
	GetByURNAndProjectID(urn string, projectID bson.ObjectID) (*models.Resource, error)
}
