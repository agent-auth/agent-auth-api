package roles_permissions

import (
	"fmt"
	"net/http"
	"time"

	"github.com/agent-auth/agent-auth-api/database/dbmodels"
	_ "github.com/agent-auth/agent-auth-api/web/interfaces/v1/errorinterface" // docs is generated by Swag CLI, you have to import it.
	"github.com/agent-auth/agent-auth-api/web/renderers"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// Request/Response models
type RoleRequest struct {
	*dbmodels.Roles
}

func (p *RoleRequest) Bind(r *http.Request) error {
	return nil
}

type RoleResponse struct {
	*dbmodels.Roles
}

type UpdateAttributeRequest struct {
	Resource string      `json:"resource"`
	Path     string      `json:"path"`
	Value    interface{} `json:"value"`
}

func (u *UpdateAttributeRequest) Bind(r *http.Request) error {
	if u.Resource == "" {
		return ErrIncompleteDetails
	}
	if u.Path == "" {
		return ErrIncompleteDetails
	}

	if u.Value == nil {
		return ErrIncompleteDetails
	}

	return nil
}

// @Summary Create role
// @Description Creates a new role
// @Tags roles
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param role body RoleRequest true "Role details"
// @Success 200 {object} RoleResponse
// @Failure 400 {object} errorinterface.ErrorResponse
// @Failure 500 {object} errorinterface.ErrorResponse
// @Router /projects/{project_id}/roles [post]
// @Security BearerAuth
func (rp *rolesService) CreateRole(w http.ResponseWriter, r *http.Request) {
	projectID, email, err := rp.hasMemberAccess(r)
	if err != nil {
		rp.logger.Error("project membership verification failed", zap.Error(err))
		render.Render(w, r, renderers.ErrorForbidden(err))
		return
	}

	var req RoleRequest
	if err := render.Bind(r, &req); err != nil {
		rp.logger.Error("failed to bind request", zap.Error(err))
		render.Render(w, r, renderers.ErrorBadRequest(err))
		return
	}

	// Set timestamps
	req.Roles.CreatedTimestampUTC = time.Now().UTC()
	req.Roles.UpdatedTimestampUTC = req.Roles.CreatedTimestampUTC

	req.Roles.ProjectID = projectID
	req.Roles.OwnerID = email
	req.Roles.Permissions = dbmodels.Permissions{}

	if err := req.Roles.Validate(); err != nil {
		rp.logger.Error("invalid role details", zap.Error(err))
		render.Render(w, r, renderers.ErrorBadRequest(err))
		return
	}

	role, err := rp.rolesDal.Create(req.Roles)
	if err != nil {
		rp.logger.Error("failed to create role", zap.Error(err))
		render.Render(w, r, renderers.ErrorInternalServerError(err))
		return
	}

	render.Respond(w, r, &RoleResponse{
		Roles: role,
	})
}

// @Summary Get role
// @Description Gets a role by ID
// @Tags roles
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param role_id path string true "Role ID"
// @Success 200 {object} RoleResponse
// @Failure 400 {object} errorinterface.ErrorResponse
// @Failure 404 {object} errorinterface.ErrorResponse
// @Router /projects/{project_id}/roles/{role_id} [get]
// @Security BearerAuth
func (rp *rolesService) GetRole(w http.ResponseWriter, r *http.Request) {
	_, _, err := rp.hasMemberAccess(r)
	if err != nil {
		rp.logger.Error("project membership verification failed", zap.Error(err))
		render.Render(w, r, renderers.ErrorForbidden(err))
		return
	}

	roleID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "role_id"))
	if err != nil {
		rp.logger.Error("invalid role ID", zap.Error(err))
		render.Render(w, r, renderers.ErrorBadRequest(err))
		return
	}

	role, err := rp.rolesDal.Get(roleID)
	if err != nil {
		rp.logger.Error("failed to get role", zap.Error(err))
		render.Render(w, r, renderers.ErrorNotFound(err))
		return
	}

	render.Respond(w, r, &RoleResponse{
		Roles: role,
	})
}

// @Summary Delete permission
// @Description Deletes a role by ID
// @Tags roles
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param role_id path string true "Role ID"
// @Success 204 "No Content"
// @Failure 400 {object} errorinterface.ErrorResponse
// @Failure 404 {object} errorinterface.ErrorResponse
// @Router /projects/{project_id}/roles/{role_id} [delete]
// @Security BearerAuth
func (rp *rolesService) DeleteRole(w http.ResponseWriter, r *http.Request) {
	_, _, err := rp.hasMemberAccess(r)
	if err != nil {
		rp.logger.Error("project membership verification failed", zap.Error(err))
		render.Render(w, r, renderers.ErrorForbidden(fmt.Errorf("project membership verification failed")))
		return
	}

	roleID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "role_id"))
	if err != nil {
		rp.logger.Error("invalid role ID", zap.Error(err))
		render.Render(w, r, renderers.ErrorBadRequest(fmt.Errorf("invalid role ID")))
		return
	}

	if err := rp.rolesDal.Delete(roleID); err != nil {
		rp.logger.Error("failed to delete role", zap.Error(err))
		render.Render(w, r, renderers.ErrorNotFound(fmt.Errorf("failed to delete role")))
		return
	}

	render.Status(r, http.StatusNoContent)
}

// @Summary Get roles by project
// @Description Gets all roles for a specific project
// @Tags roles
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 {array} RoleResponse
// @Failure 400 {object} errorinterface.ErrorResponse
// @Failure 404 {object} errorinterface.ErrorResponse
// @Router /projects/{project_id}/roles [get]
// @Security BearerAuth
func (rp *rolesService) GetRolesByProject(w http.ResponseWriter, r *http.Request) {
	projectID, _, err := rp.hasMemberAccess(r)
	if err != nil {
		rp.logger.Error("project membership verification failed", zap.Error(err))
		render.Render(w, r, renderers.ErrorForbidden(fmt.Errorf("project membership verification failed")))
		return
	}

	roles, err := rp.rolesDal.GetByProjectID(projectID)
	if err != nil {
		rp.logger.Error("failed to get project roles", zap.Error(err))
		render.Render(w, r, renderers.ErrorNotFound(fmt.Errorf("failed to get project roles")))
		return
	}

	response := make([]*RoleResponse, len(roles))
	for i, role := range roles {
		response[i] = &RoleResponse{Roles: role}
	}

	render.Respond(w, r, response)
}

// @Summary Delete permissions by project
// @Description Deletes all roles for a specific project
// @Tags roles
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 204 "No Content"
// @Failure 400 {object} errorinterface.ErrorResponse
// @Failure 404 {object} errorinterface.ErrorResponse
// @Router /projects/{project_id}/roles [delete]
// @Security BearerAuth
func (rp *rolesService) DeleteRolesByProject(w http.ResponseWriter, r *http.Request) {
	projectID, _, err := rp.hasMemberAccess(r)
	if err != nil {
		rp.logger.Error("project membership verification failed", zap.Error(err))
		render.Render(w, r, renderers.ErrorForbidden(fmt.Errorf("project membership verification failed")))
		return
	}

	if err := rp.rolesDal.DeleteByProjectID(projectID); err != nil {
		rp.logger.Error("failed to delete project roles", zap.Error(err))
		render.Render(w, r, renderers.ErrorNotFound(fmt.Errorf("failed to delete project roles")))
		return
	}

	render.Status(r, http.StatusNoContent)
}
