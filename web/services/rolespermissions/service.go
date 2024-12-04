package rolespermissions

import (
	projectsdal "github.com/agent-auth/agent-auth-api/dal/project"
	rolespermissionsdal "github.com/agent-auth/agent-auth-api/dal/rolespermissions"
	"github.com/agent-auth/agent-auth-api/pkg/logger"
	"go.uber.org/zap"
)

type rolesService struct {
	logger      *zap.Logger
	rolesDal    rolespermissionsdal.RolesDal
	projectsDal projectsdal.ProjectsDal
}

// NewRolesService returns service impl
func NewRolesService() RolesService {
	return &rolesService{
		logger:      logger.NewLogger(),
		rolesDal:    rolespermissionsdal.NewRolesDal(),
		projectsDal: projectsdal.NewProjectsDal(),
	}
}
