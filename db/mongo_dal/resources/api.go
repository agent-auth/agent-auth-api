package resources_dal

import (
	"github.com/agent-auth/common-lib/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ResourcesDal defines the interface for resource database operations
type ResourcesDal interface {
	Create(resource *models.Resource) (*models.Resource, error)
	Update(resource *models.Resource) error
	GetByID(id primitive.ObjectID) (*models.Resource, error)
	Delete(id primitive.ObjectID) error
	GetByProjectID(projectID primitive.ObjectID) ([]*models.Resource, error)
	GetByURNAndProjectID(urn string, projectID primitive.ObjectID) (*models.Resource, error)
}
