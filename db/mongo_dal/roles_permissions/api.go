package roles_permissions_dal

import (
	"github.com/agent-auth/common-lib/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RolesDal interface {
	Create(role *models.Roles) (*models.Roles, error)
	Get(id bson.ObjectID) (*models.Roles, error)
	Delete(id bson.ObjectID) error
	GetByProjectID(projectID bson.ObjectID) ([]*models.Roles, error)
	DeleteByProjectID(projectID bson.ObjectID) error
	GetByProjectIDAndRole(projectID bson.ObjectID, role string) (*models.Roles, error)

	UpdatePermission(id bson.ObjectID, resource string, actions []models.Action) error
}
