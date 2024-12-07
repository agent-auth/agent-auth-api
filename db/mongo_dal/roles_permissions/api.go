package roles_permissions_dal

import (
	"github.com/agent-auth/common-lib/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RolesDal interface {
	Create(role *models.Roles) (*models.Roles, error)
	Get(id primitive.ObjectID) (*models.Roles, error)
	Delete(id primitive.ObjectID) error
	GetByProjectID(projectID primitive.ObjectID) ([]*models.Roles, error)
	DeleteByProjectID(projectID primitive.ObjectID) error
	GetByProjectIDAndRole(projectID primitive.ObjectID, role string) (*models.Roles, error)

	UpdatePermission(id primitive.ObjectID, resource string, actions []models.Action) error
}
