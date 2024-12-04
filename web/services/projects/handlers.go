package projects

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/agent-auth/agent-auth-api/database/dbmodels"
	"github.com/agent-auth/agent-auth-api/pkg/authz"
	"github.com/agent-auth/agent-auth-api/web/renderers"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type ProjectRequest struct {
	*dbmodels.Project
}

func (p *ProjectRequest) Bind(r *http.Request) error {
	if p.Project == nil {
		return ErrIncompleteDetails
	}
	if p.Name == "" {
		return ErrIncompleteDetails
	}
	return nil
}

type ProjectResponse struct {
	*dbmodels.Project
}

type ProjectsResponse struct {
	Projects []*dbmodels.Project `json:"projects"`
}

type AddMemberRequest struct {
	Email string `json:"email"`
}

func (a *AddMemberRequest) Bind(r *http.Request) error {
	if a.Email == "" {
		return ErrIncompleteDetails
	}
	return nil
}

// @Summary Create project
// @Description Creates a new project
// @Tags projects
// @Accept json
// @Produce json
// @Param project body ProjectRequest true "Project details"
// @Success 200 {object} ProjectResponse
// @Failure 400 {object} errorinterface.ErrorResponse
// @Failure 500 {object} errorinterface.ErrorResponse
// @Router /projects [post]
// @Security BearerAuth
func (ps *projectService) Create(w http.ResponseWriter, r *http.Request) {
	email, _ := authz.GetEmailFromClaims(r)

	project := &ProjectRequest{
		Project: &dbmodels.Project{},
	}

	if err := render.Bind(r, project); err != nil {
		ps.logger.Error("failed to bind project request", zap.Error(err))
		render.Render(w, r, renderers.ErrorBadRequest(errors.New("failed to bind project request")))
		return
	}

	project.Project.OwnerID = email
	project.Project.Members = []string{email}

	if err := project.Project.Validate(); err != nil {
		ps.logger.Error("failed to validate project", zap.Error(err))
		render.Render(w, r, renderers.ErrorBadRequest(errors.New("failed to validate project")))
		return
	}

	resp, err := ps.projectDal.Create(project.Project)
	if err != nil {
		ps.logger.Error("failed to create project", zap.Error(err))
		render.Render(w, r, renderers.ErrorInternalServerError(errors.New("failed to create project")))
		return
	}

	render.Respond(w, r, &ProjectResponse{
		Project: resp,
	})
}

// @Summary Get project
// @Description Gets a project by ID
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 {object} ProjectResponse
// @Failure 400 {object} errorinterface.ErrorResponse
// @Failure 401 {object} errorinterface.ErrorResponse
// @Failure 404 {object} errorinterface.ErrorResponse
// @Router /projects/{project_id} [get]
// @Security BearerAuth
func (ps *projectService) Get(w http.ResponseWriter, r *http.Request) {
	projectID, email, err := ps.hasMemberAccess(r)
	if err != nil {
		ps.logger.Error("unauthorized access attempt", zap.String("userID", email))
		render.Render(w, r, renderers.ErrorUnauthorized(errors.New("unauthorized access attempt")))
		return
	}

	project, err := ps.projectDal.GetByID(projectID)
	if err != nil {
		ps.logger.Error("failed to get project", zap.Error(err))
		render.Render(w, r, renderers.ErrorNotFound(errors.New("failed to get project")))
		return
	}

	render.Respond(w, r, &ProjectResponse{Project: project})
}

// @Summary List projects
// @Description Lists all projects for a workspace with pagination
// @Tags projects
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param skip query integer false "Number of records to skip" default(0)
// @Param limit query integer false "Number of records to return" default(10)
// @Success 200 {object} ProjectsResponse
// @Failure 400 {object} errorinterface.ErrorResponse
// @Failure 500 {object} errorinterface.ErrorResponse
// @Router /projects [get]
// @Security BearerAuth
func (ps *projectService) List(w http.ResponseWriter, r *http.Request) {
	// it must find all the projects user is member of
	email, _ := authz.GetEmailFromClaims(r)

	skip, err := strconv.ParseInt(r.URL.Query().Get("skip"), 10, 64)
	if err != nil {
		skip = 0
	}

	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		limit = 10
	}

	if limit > 100 {
		limit = 100 // Maximum page size
	}
	if limit < 1 {
		limit = 10 // Default page size
	}
	if skip < 0 {
		skip = 0
	}

	projects, err := ps.projectDal.List(email, skip, limit)
	if err != nil {
		ps.logger.Error("failed to list projects", zap.Error(err))
		render.Render(w, r, renderers.ErrorInternalServerError(errors.New("failed to list projects")))
		return
	}

	render.Respond(w, r, &ProjectsResponse{Projects: projects})
}

// @Summary Update project
// @Description Updates an existing project (owner only)
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param project body ProjectRequest true "Updated project details"
// @Success 200 {object} ProjectResponse
// @Failure 400 {object} errorinterface.ErrorResponse
// @Failure 401 {object} errorinterface.ErrorResponse
// @Failure 404 {object} errorinterface.ErrorResponse
// @Router /projects/{project_id} [put]
// @Security BearerAuth
func (ps *projectService) Update(w http.ResponseWriter, r *http.Request) {
	projectID, email, err := ps.hasMemberAccess(r)
	if err != nil {
		ps.logger.Error("unauthorized access attempt", zap.String("userID", email))
		render.Render(w, r, renderers.ErrorUnauthorized(errors.New("unauthorized access attempt")))
		return
	}

	// Get existing project
	existing, err := ps.projectDal.GetByID(projectID)
	if err != nil {
		ps.logger.Error("failed to get project", zap.Error(err))
		render.Render(w, r, renderers.ErrorNotFound(errors.New("failed to get project")))
		return
	}

	// Bind and validate update request
	updateReq := &ProjectRequest{Project: &dbmodels.Project{}}
	if err := render.Bind(r, updateReq); err != nil {
		ps.logger.Error("failed to bind project request", zap.Error(err))
		render.Render(w, r, renderers.ErrorBadRequest(errors.New("failed to bind project request")))
		return
	}

	// Update only mutable fields
	existing.Name = updateReq.Name
	existing.Description = updateReq.Description
	existing.UpdatedTimestampUTC = time.Now()

	if err := existing.Validate(); err != nil {
		ps.logger.Error("failed to validate project", zap.Error(err))
		render.Render(w, r, renderers.ErrorBadRequest(errors.New("failed to validate project")))
		return
	}

	if err := ps.projectDal.Update(existing); err != nil {
		ps.logger.Error("failed to update project", zap.Error(err))
		render.Render(w, r, renderers.ErrorInternalServerError(errors.New("failed to update project")))
		return
	}

	render.Respond(w, r, &ProjectResponse{Project: existing})
}

// @Summary Delete project
// @Description Deletes a project (owner only)
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 204 "No Content"
// @Failure 400 {object} errorinterface.ErrorResponse
// @Failure 401 {object} errorinterface.ErrorResponse
// @Failure 404 {object} errorinterface.ErrorResponse
// @Router /projects/{project_id} [delete]
// @Security BearerAuth
func (ps *projectService) Delete(w http.ResponseWriter, r *http.Request) {
	projectID, email, err := ps.hasMemberAccess(r)
	if err != nil {
		ps.logger.Error("unauthorized access attempt", zap.String("userID", email))
		render.Render(w, r, renderers.ErrorUnauthorized(errors.New("unauthorized access attempt")))
		return
	}

	existing, err := ps.projectDal.GetByID(projectID)
	if err != nil {
		ps.logger.Error("failed to get project", zap.Error(err))
		render.Render(w, r, renderers.ErrorNotFound(errors.New("failed to get project")))
		return
	}

	// only owner can delete project
	if email != existing.OwnerID {
		ps.logger.Error("unauthorized access attempt", zap.String("userID", email))
		render.Render(w, r, renderers.ErrorUnauthorized(errors.New("unauthorized access attempt")))
		return
	}

	if err := ps.projectDal.Delete(projectID); err != nil {
		ps.logger.Error("failed to delete project", zap.Error(err))
		render.Render(w, r, renderers.ErrorInternalServerError(errors.New("failed to delete project")))
		return
	}

	render.Status(r, http.StatusNoContent)
}

// @Summary Remove member from project
// @Description Removes a member from a project (owner only)
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 204 "No Content"
// @Failure 400 {object} errorinterface.ErrorResponse
// @Failure 401 {object} errorinterface.ErrorResponse
// @Failure 404 {object} errorinterface.ErrorResponse
// @Router /projects/{project_id}/members [delete]
// @Security BearerAuth
func (ps *projectService) RemoveMember(w http.ResponseWriter, r *http.Request) {
	projectID, email, err := ps.hasMemberAccess(r)
	if err != nil {
		ps.logger.Error("unauthorized access attempt", zap.String("userID", email))
		render.Render(w, r, renderers.ErrorUnauthorized(errors.New("unauthorized access attempt")))
		return
	}

	var req AddMemberRequest
	if err := render.Bind(r, &req); err != nil {
		ps.logger.Error("failed to bind project request", zap.Error(err))
		render.Render(w, r, renderers.ErrorBadRequest(errors.New("failed to bind project request")))
		return
	}

	existing, err := ps.projectDal.GetByID(projectID)
	if err != nil {
		ps.logger.Error("failed to get project", zap.Error(err))
		render.Render(w, r, renderers.ErrorNotFound(errors.New("failed to get project")))
		return
	}

	if req.Email == existing.OwnerID {
		ps.logger.Error("attempt to remove owner", zap.String("email", req.Email))
		render.Render(w, r, renderers.ErrorBadRequest(errors.New("cannot remove project owner")))
		return
	}

	if req.Email == email && email != existing.OwnerID {
		ps.logger.Error("non-owner attempting self-removal", zap.String("email", email))
		render.Render(w, r, renderers.ErrorUnauthorized(errors.New("unauthorized access attempt")))
		return
	}

	if err := ps.projectDal.RemoveMember(projectID, req.Email); err != nil {
		ps.logger.Error("failed to remove member", zap.Error(err))
		render.Render(w, r, renderers.ErrorInternalServerError(errors.New("failed to remove member")))
		return
	}

	render.Status(r, http.StatusNoContent)
}

// @Summary Add member to project
// @Description Adds a new member to a project
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param member body AddMemberRequest true "Member ID to add"
// @Success 204 "No Content"
// @Failure 400 {object} errorinterface.ErrorResponse
// @Failure 401 {object} errorinterface.ErrorResponse
// @Failure 404 {object} errorinterface.ErrorResponse
// @Failure 500 {object} errorinterface.ErrorResponse
// @Router /projects/{project_id}/members [post]
// @Security BearerAuth
func (ps *projectService) AddMember(w http.ResponseWriter, r *http.Request) {
	projectID, email, err := ps.hasMemberAccess(r)
	if err != nil {
		ps.logger.Error("unauthorized access attempt", zap.String("userID", email))
		render.Render(w, r, renderers.ErrorUnauthorized(errors.New("unauthorized access attempt")))
		return
	}

	var req AddMemberRequest
	if err := render.Bind(r, &req); err != nil {
		ps.logger.Error("failed to bind project request", zap.Error(err))
		render.Render(w, r, renderers.ErrorBadRequest(errors.New("failed to bind project request")))
		return
	}

	isMember, err := ps.projectDal.IsMember(projectID, req.Email)
	if err != nil {
		ps.logger.Error("failed to check member status", zap.Error(err))
		render.Render(w, r, renderers.ErrorInternalServerError(errors.New("failed to check member status")))
		return
	}
	if isMember {
		ps.logger.Error("attempt to add existing member", zap.String("email", req.Email))
		render.Render(w, r, renderers.ErrorBadRequest(errors.New("user is already a member")))
		return
	}

	if err := ps.projectDal.AddMember(projectID, req.Email); err != nil {
		ps.logger.Error("failed to add member", zap.Error(err))
		render.Render(w, r, renderers.ErrorInternalServerError(errors.New("failed to add member")))
		return
	}

	render.Status(r, http.StatusNoContent)
}
