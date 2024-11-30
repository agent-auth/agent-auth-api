package jobsservice

import (
	"github.com/agent-auth/agent-auth-api/dal/jobsdal"
	"github.com/agent-auth/agent-auth-api/dal/organizationdal"
	"github.com/agent-auth/agent-auth-api/dal/recruiterdal"
	"github.com/agent-auth/agent-auth-api/dal/tokendal"
	"github.com/agent-auth/agent-auth-api/logging"
	"github.com/agent-auth/agent-auth-api/web/middlewares"
)

type jobsservice struct {
	logger logging.Logger

	tokenDal     tokendal.TokenDal
	recruiterDal recruiterdal.RecruiterDal
	orgDal       organizationdal.OrganizationDal
	jobsDal      jobsdal.JobsDal
	middlewares  middlewares.Middlewares
}

// NewJobService returns service impl
func NewJobService() JobsService {
	return &jobsservice{
		logger: logging.NewLogger(),

		tokenDal:     tokendal.NewTokenDal(),
		recruiterDal: recruiterdal.NewRecruiterDal(),
		orgDal:       organizationdal.NewOrganizationDal(),
		jobsDal:      jobsdal.NewJobsDal(),
		middlewares:  middlewares.NewMiddlewares(),
	}
}
