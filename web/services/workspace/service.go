package workspace

import (
	projectdal "github.com/agent-auth/agent-auth-api/dal/workspace"
	"go.uber.org/zap"
)

type workspaceService struct {
	logger       *zap.Logger
	workspaceDal projectdal.WorkspaceDal
}

// NewWorkspaceService returns service impl
func NewWorkspaceService() WorkspaceService {
	return &workspaceService{
		logger:       zap.L(),
		workspaceDal: projectdal.NewWorkspaceDal(),
	}
}
