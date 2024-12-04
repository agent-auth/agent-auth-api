package rolespermissions

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

	UpdatePermission(id primitive.ObjectID, resource string, key string, value interface{}) error
	RemovePermission(id primitive.ObjectID, resource string, key string) error
}
