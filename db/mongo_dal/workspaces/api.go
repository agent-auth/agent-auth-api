package workspaces_dal

import (
	"github.com/agent-auth/common-lib/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WorkspaceDal defines the interface for workspace database operations
type WorkspaceDal interface {
	Create(workspace *models.Workspace) (*models.Workspace, error)
	Update(workspace *models.Workspace) error
	GetByID(id primitive.ObjectID) (*models.Workspace, error)
	Delete(id primitive.ObjectID) error
	List(skip, limit int64) ([]*models.Workspace, error)
	GetBySlug(slug string) (*models.Workspace, error)
	GetByOwnerID(ownerID primitive.ObjectID) ([]*models.Workspace, error)
	AddMember(workspaceID string, memberID string) error
	RemoveMember(workspaceID string, memberID string) error
	IsMember(workspaceID, email string) (bool, error)
}
