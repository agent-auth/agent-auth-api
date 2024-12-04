package projects_dal

import (
	"github.com/agent-auth/agent-auth-api/database/dbmodels"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProjectsDal defines the interface for project database operations
type ProjectsDal interface {
	Create(project *dbmodels.Project) (*dbmodels.Project, error)
	Update(project *dbmodels.Project) error
	GetByID(id primitive.ObjectID) (*dbmodels.Project, error)
	Delete(id primitive.ObjectID) error
	List(email string, skip, limit int64) ([]*dbmodels.Project, error)
	GetBySlug(workspaceID primitive.ObjectID, slug string) (*dbmodels.Project, error)
	GetByOwnerID(ownerID string) ([]*dbmodels.Project, error)
	AddMember(projectID primitive.ObjectID, email string) error
	RemoveMember(projectID primitive.ObjectID, email string) error
	IsMember(projectID primitive.ObjectID, email string) (bool, error)
}