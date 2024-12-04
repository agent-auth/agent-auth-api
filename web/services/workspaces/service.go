package workspaces

import (
	workspaces_dal "github.com/agent-auth/agent-auth-api/dal/workspaces"
	"github.com/agent-auth/agent-auth-api/pkg/logger"
	"go.uber.org/zap"
)

type workspaceService struct {
	logger       *zap.Logger
	workspaceDal workspaces_dal.WorkspaceDal
}

// NewWorkspaceService returns service impl
func NewWorkspaceService() WorkspaceService {
	return &workspaceService{
		logger:       logger.NewLogger(),
		workspaceDal: workspaces_dal.NewWorkspaceDal(),
	}
}