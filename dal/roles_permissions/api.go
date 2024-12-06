package roles_permissions_dal

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/agent-auth/agent-auth-api/database/dbmodels"
)

type RolesDal interface {
	Create(role *dbmodels.Roles) (*dbmodels.Roles, error)
	Get(id primitive.ObjectID) (*dbmodels.Roles, error)
	Delete(id primitive.ObjectID) error
	GetByProjectID(projectID primitive.ObjectID) ([]*dbmodels.Roles, error)
	DeleteByProjectID(projectID primitive.ObjectID) error
	GetByProjectIDAndName(projectID primitive.ObjectID, name string) (*dbmodels.Roles, error)

	UpdatePermission(id primitive.ObjectID, resource string, actions []dbmodels.Action) error
}
