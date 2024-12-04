package rolespermissions

import (
	rolespermissionsdal "github.com/agent-auth/agent-auth-api/dal/rolespermissions"
	"go.uber.org/zap"
)

type rolesService struct {
	logger   *zap.Logger
	rolesDal rolespermissionsdal.RolesDal
}

// NewRolesService returns service impl
func NewRolesService() RolesService {
	return &rolesService{
		logger:   zap.L(),
		rolesDal: rolespermissionsdal.NewRolesDal(),
	}
}
