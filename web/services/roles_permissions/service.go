package roles_permissions

import (
	projects_dal "github.com/agent-auth/agent-auth-api/dal/projects"
	resources_dal "github.com/agent-auth/agent-auth-api/dal/resources"
	roles_permissions_dal "github.com/agent-auth/agent-auth-api/dal/roles_permissions"
	"github.com/agent-auth/agent-auth-api/pkg/logger"
	"go.uber.org/zap"
)

type rolesService struct {
	logger       *zap.Logger
	rolesDal     roles_permissions_dal.RolesDal
	resourcesDal resources_dal.ResourcesDal
	projectsDal  projects_dal.ProjectsDal
}

// NewRolesService returns service impl
func NewRolesService() RolesService {
	return &rolesService{
		logger:       logger.NewLogger(),
		rolesDal:     roles_permissions_dal.NewRolesDal(),
		resourcesDal: resources_dal.NewResourcesDal(),
		projectsDal:  projects_dal.NewProjectsDal(),
	}
}
