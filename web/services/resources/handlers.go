package resources

import (
	"errors"
	"net/http"

	"github.com/agent-auth/agent-auth-api/database/dbmodels"
	"github.com/agent-auth/agent-auth-api/web/renderers"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type ResourceRequest struct {
	*dbmodels.Resource
}

func (r *ResourceRequest) Bind(req *http.Request) error {
	return nil
}

type ResourceResponse struct {
	*dbmodels.Resource
}

type ResourcesResponse struct {
	Resources []*dbmodels.Resource `json:"resources"`
}

// @Summary Create resource
// @Description Creates a new resource
// @Tags resources
// @Accept json
// @Produce json
// @Param resource body ResourceRequest true "Resource details"
// @Success 200 {object} ResourceResponse
// @Failure 400,401,500 {object} errorinterface.ErrorResponse
// @Router /projects/{project_id}/resources [post]
// @Security BearerAuth
func (rs *resourceService) Create(w http.ResponseWriter, r *http.Request) {
	// Verify project membership first
	project_id, email, err := rs.hasMemberAccess(r)
	if err != nil {
		rs.logger.Error("unauthorized access attempt", zap.Error(err))
		render.Render(w, r, renderers.ErrorUnauthorized(errors.New("unauthorized access attempt")))
		return
	}

	resource := &ResourceRequest{
		Resource: &dbmodels.Resource{},
	}

	if err := render.Bind(r, resource); err != nil {
		rs.logger.Error("failed to bind resource request", zap.Error(err))
		render.Render(w, r, renderers.ErrorBadRequest(errors.New("invalid resource data")))
		return
	}

	// Ensure the resource is created for the verified project
	resource.Resource.ProjectID = project_id
	resource.Resource.OwnerID = email

	if err := resource.Resource.Validate(); err != nil {
		rs.logger.Error("invalid resource data", zap.Error(err))
		render.Render(w, r, renderers.ErrorBadRequest(errors.New("invalid resource data")))
		return
	}

	resp, err := rs.resources_dal.Create(resource.Resource)
	if err != nil {
		rs.logger.Error("failed to create resource", zap.Error(err))
		render.Render(w, r, renderers.ErrorInternalServerError(errors.New("failed to create resource")))
		return
	}

	render.Respond(w, r, &ResourceResponse{
		Resource: resp,
	})
}

// @Summary Get resource
// @Description Gets a resource by ID
// @Tags resources
// @Accept json
// @Produce json
// @Param resource_id path string true "Resource ID"
// @Success 200 {object} ResourceResponse
// @Failure 400,401,404 {object} errorinterface.ErrorResponse
// @Router /projects/{project_id}/resources/{resource_id} [get]
// @Security BearerAuth
func (rs *resourceService) Get(w http.ResponseWriter, r *http.Request) {
	// Verify project membership first
	_, _, err := rs.hasMemberAccess(r)
	if err != nil {
		rs.logger.Error("unauthorized access attempt", zap.Error(err))
		render.Render(w, r, renderers.ErrorUnauthorized(errors.New("unauthorized access attempt")))
		return
	}

	resource_id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "resource_id"))
	if err != nil {
		rs.logger.Error("invalid resource ID", zap.Error(err))
		render.Render(w, r, renderers.ErrorBadRequest(errors.New("invalid resource ID")))
		return
	}

	resource, err := rs.resources_dal.GetByID(resource_id)
	if err != nil {
		rs.logger.Error("failed to get resource", zap.Error(err))
		render.Render(w, r, renderers.ErrorNotFound(errors.New("resource not found")))
		return
	}

	render.Respond(w, r, &ResourceResponse{Resource: resource})
}

// @Summary List resources by project
// @Description Lists all resources for a project
// @Tags resources
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 {object} ResourcesResponse
// @Failure 400,401,500 {object} errorinterface.ErrorResponse
// @Router /projects/{project_id}/resources [get]
// @Security BearerAuth
func (rs *resourceService) ListByProject(w http.ResponseWriter, r *http.Request) {
	project_id, _, err := rs.hasMemberAccess(r)
	if err != nil {
		rs.logger.Error("unauthorized access attempt", zap.Error(err))
		render.Render(w, r, renderers.ErrorUnauthorized(errors.New("unauthorized access attempt")))
		return
	}

	resources, err := rs.resources_dal.GetByProjectID(project_id)
	if err != nil {
		rs.logger.Error("failed to list resources", zap.Error(err))
		render.Render(w, r, renderers.ErrorInternalServerError(errors.New("failed to list resources")))
		return
	}

	render.Respond(w, r, &ResourcesResponse{Resources: resources})
}

// @Summary Update resource
// @Description Updates an existing resource
// @Tags resources
// @Accept json
// @Produce json
// @Param resource_id path string true "Resource ID"
// @Param resource body ResourceRequest true "Updated resource details"
// @Success 200 {object} ResourceResponse
// @Failure 400,401,404,500 {object} errorinterface.ErrorResponse
// @Router /projects/{project_id}/resources/{resource_id} [put]
// @Security BearerAuth
func (rs *resourceService) Update(w http.ResponseWriter, r *http.Request) {
	project_id, _, err := rs.hasMemberAccess(r)
	if err != nil {
		rs.logger.Error("unauthorized access attempt", zap.Error(err))
		render.Render(w, r, renderers.ErrorUnauthorized(errors.New("unauthorized access attempt")))
		return
	}

	resource_id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "resource_id"))
	if err != nil {
		render.Render(w, r, renderers.ErrorBadRequest(errors.New("invalid resource ID")))
		return
	}

	if err := rs.hasRoleAccess(resource_id, project_id); err != nil {
		rs.logger.Error("resource verification failed", zap.Error(err))
		render.Render(w, r, renderers.ErrorUnauthorized(err))
		return
	}

	existing, err := rs.resources_dal.GetByID(resource_id)
	if err != nil {
		rs.logger.Error("failed to get resource", zap.Error(err))
		render.Render(w, r, renderers.ErrorNotFound(errors.New("resource not found")))
		return
	}

	resource := &ResourceRequest{Resource: existing}
	if err := render.Bind(r, resource); err != nil {
		rs.logger.Error("failed to bind update request", zap.Error(err))
		render.Render(w, r, renderers.ErrorBadRequest(errors.New("invalid update data")))
		return
	}

	// Update mutable fields
	existing.Description = resource.Description
	existing.Actions = resource.Actions

	if err := existing.Validate(); err != nil {
		rs.logger.Error("invalid resource data", zap.Error(err))
		render.Render(w, r, renderers.ErrorBadRequest(errors.New("invalid resource data")))
		return
	}

	if err := rs.resources_dal.Update(existing); err != nil {
		rs.logger.Error("failed to update resource", zap.Error(err))
		render.Render(w, r, renderers.ErrorInternalServerError(errors.New("failed to update resource")))
		return
	}

	render.Respond(w, r, &ResourceResponse{Resource: existing})
}

// @Summary Delete resource
// @Description Deletes a resource
// @Tags resources
// @Accept json
// @Produce json
// @Param resource_id path string true "Resource ID"
// @Success 204 "No Content"
// @Failure 400,401,404,500 {object} errorinterface.ErrorResponse
// @Router /projects/{project_id}/resources/{resource_id} [delete]
// @Security BearerAuth
func (rs *resourceService) Delete(w http.ResponseWriter, r *http.Request) {
	project_id, _, err := rs.hasMemberAccess(r)
	if err != nil {
		rs.logger.Error("unauthorized access attempt", zap.Error(err))
		render.Render(w, r, renderers.ErrorUnauthorized(errors.New("unauthorized access attempt")))
		return
	}

	resource_id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "resource_id"))
	if err != nil {
		render.Render(w, r, renderers.ErrorBadRequest(errors.New("invalid resource ID")))
		return
	}

	if err := rs.hasRoleAccess(resource_id, project_id); err != nil {
		rs.logger.Error("resource verification failed", zap.Error(err))
		render.Render(w, r, renderers.ErrorUnauthorized(err))
		return
	}

	if err := rs.resources_dal.Delete(resource_id); err != nil {
		rs.logger.Error("failed to delete resource", zap.Error(err))
		render.Render(w, r, renderers.ErrorInternalServerError(errors.New("failed to delete resource")))
		return
	}

	render.Status(r, http.StatusNoContent)
}
