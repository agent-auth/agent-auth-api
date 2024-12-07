package projects_dal

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/agent-auth/common-lib/models"
)

// ProjectsDal defines the interface for project database operations
type ProjectsDal interface {
	Create(project *models.Project) (*models.Project, error)
	Update(project *models.Project) error
	GetByID(id primitive.ObjectID) (*models.Project, error)
	Delete(id primitive.ObjectID) error
	List(email string, skip, limit int64) ([]*models.Project, error)
	GetBySlug(workspaceID primitive.ObjectID, slug string) (*models.Project, error)
	GetByOwnerID(ownerID string) ([]*models.Project, error)
	AddMember(projectID primitive.ObjectID, email string) error
	RemoveMember(projectID primitive.ObjectID, email string) error
	IsMember(projectID primitive.ObjectID, email string) (bool, error)
}
