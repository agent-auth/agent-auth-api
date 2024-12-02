package projectdal

import (
	"github.com/agent-auth/agent-auth-api/database/dbmodels"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WorkspaceDal defines the interface for workspace database operations
type WorkspaceDal interface {
	Create(txID string, workspace *dbmodels.Workspace) (*dbmodels.Workspace, error)
	Update(workspace *dbmodels.Workspace) error
	GetByID(id primitive.ObjectID) (*dbmodels.Workspace, error)
	Delete(id primitive.ObjectID) error
	List(skip, limit int64) ([]*dbmodels.Workspace, error)
	GetBySlug(slug string) (*dbmodels.Workspace, error)
	GetByOwnerID(ownerID primitive.ObjectID) ([]*dbmodels.Workspace, error)
	AddMember(workspaceID, memberID primitive.ObjectID) error
	RemoveMember(workspaceID, memberID primitive.ObjectID) error
}
