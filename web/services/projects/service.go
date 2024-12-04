package projects

import (
	projects_dal "github.com/agent-auth/agent-auth-api/dal/projects"
	"github.com/agent-auth/agent-auth-api/pkg/logger"
	"go.uber.org/zap"
)

type projectService struct {
	logger     *zap.Logger
	projectDal projects_dal.ProjectsDal
}

// NewProjectService returns service impl
func NewProjectService() ProjectService {
	return &projectService{
		logger:     logger.NewLogger(),
		projectDal: projects_dal.NewProjectsDal(),
	}
}
