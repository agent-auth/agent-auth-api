package roles_permissions

import (
	projects_dal "github.com/agent-auth/agent-auth-api/db/mongo_dal/projects"
	resources_dal "github.com/agent-auth/agent-auth-api/db/mongo_dal/resources"
	roles_dal "github.com/agent-auth/agent-auth-api/db/mongo_dal/roles_permissions"
	"github.com/agent-auth/common-lib/pkg/logger"
	"go.uber.org/zap"
)

type rolesService struct {
	logger       *zap.Logger
	rolesDal     roles_dal.RolesDal
	resourcesDal resources_dal.ResourcesDal
	projectsDal  projects_dal.ProjectsDal
}

// NewRolesService returns service impl
func NewRolesService() RolesService {
	return &rolesService{
		logger:       logger.NewLogger(),
		rolesDal:     roles_dal.NewRolesDal(),
		resourcesDal: resources_dal.NewResourcesDal(),
		projectsDal:  projects_dal.NewProjectsDal(),
	}
}
