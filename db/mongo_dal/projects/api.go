package projects_dal

import (
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/agent-auth/common-lib/models"
)

// ProjectsDal defines the interface for project database operations
type ProjectsDal interface {
	Create(project *models.Project) (*models.Project, error)
	Update(project *models.Project) error
	GetByID(id bson.ObjectID) (*models.Project, error)
	Delete(id bson.ObjectID) error
	List(email string, skip, limit int64) ([]*models.Project, error)
	GetBySlug(workspaceID bson.ObjectID, slug string) (*models.Project, error)
	GetByOwnerID(ownerID string) ([]*models.Project, error)
	AddMember(projectID bson.ObjectID, email string) error
	RemoveMember(projectID bson.ObjectID, email string) error
	IsMember(projectID bson.ObjectID, email string) (bool, error)
}
