package projects_dal

import (
	"github.com/agent-auth/agent-auth-api/database/dbmodels"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProjectsDal defines the interface for project database operations
type ProjectsDal interface {
	Create(txID string, project *dbmodels.Project) (*dbmodels.Project, error)
	Update(project *dbmodels.Project) error
	GetByID(id primitive.ObjectID) (*dbmodels.Project, error)
	Delete(id primitive.ObjectID) error
	List(workspaceID primitive.ObjectID, skip, limit int64) ([]*dbmodels.Project, error)
	GetBySlug(workspaceID primitive.ObjectID, slug string) (*dbmodels.Project, error)
	GetByOwnerID(ownerID primitive.ObjectID) ([]*dbmodels.Project, error)
	AddMember(projectID, memberID primitive.ObjectID) error
	RemoveMember(projectID, memberID primitive.ObjectID) error
	IsMember(projectID, email string) (bool, error)
}
