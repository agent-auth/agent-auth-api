package workspaces_dal

import (
	"github.com/agent-auth/common-lib/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// WorkspaceDal defines the interface for workspace database operations
type WorkspaceDal interface {
	Create(workspace *models.Workspace) (*models.Workspace, error)
	Update(workspace *models.Workspace) error
	GetByID(id bson.ObjectID) (*models.Workspace, error)
	Delete(id bson.ObjectID) error
	List(skip, limit int64) ([]*models.Workspace, error)
	GetBySlug(slug string) (*models.Workspace, error)
	GetByOwnerID(ownerID bson.ObjectID) ([]*models.Workspace, error)
	AddMember(workspaceID string, memberID string) error
	RemoveMember(workspaceID string, memberID string) error
	IsMember(workspaceID, email string) (bool, error)
}
