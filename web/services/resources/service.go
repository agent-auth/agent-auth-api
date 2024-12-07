package resources

import (
	projects_dal "github.com/agent-auth/agent-auth-api/db/mongo_dal/projects"
	resources_dal "github.com/agent-auth/agent-auth-api/db/mongo_dal/resources"
	"github.com/agent-auth/common-lib/pkg/logger"
	"go.uber.org/zap"
)

type resourceService struct {
	logger        *zap.Logger
	resources_dal resources_dal.ResourcesDal
	projects_dal  projects_dal.ProjectsDal
}

// NewResourceService returns service impl
func NewResourceService() ResourceService {
	return &resourceService{
		logger:        logger.NewLogger(),
		resources_dal: resources_dal.NewResourcesDal(),
		projects_dal:  projects_dal.NewProjectsDal(),
	}
}
