package roles_permissions

import (
	"errors"
	"fmt"
	"net/http"

	_ "github.com/agent-auth/agent-auth-api/web/interfaces/v1/errorinterface" // docs is generated by Swag CLI, you have to import it.
	"github.com/agent-auth/agent-auth-api/web/renderers"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// @Summary Update permission attribute
// @Description Updates a specific attribute of a permission
// @Tags permissions
// @Accept json
// @Produce json
// @Param role_id path string true "Role ID"
// @Param attribute body UpdateAttributeRequest true "Attribute update details"
// @Success 204 "No Content"
// @Failure 400 {object} errorinterface.ErrorResponse
// @Failure 404 {object} errorinterface.ErrorResponse
// @Router /projects/{project_id}/roles/{role_id}/permissions [put]
// @Security BearerAuth
func (rp *rolesService) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	projectID, _, err := rp.hasMemberAccess(r)
	if err != nil {
		msg := "project membership verification failed"
		rp.logger.Error(msg, zap.Error(err))
		render.Render(w, r, renderers.ErrorForbidden(errors.New(msg)))
		return
	}

	roleID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "role_id"))
	if err != nil {
		msg := "invalid role ID format"
		rp.logger.Error(msg, zap.Error(err))
		render.Render(w, r, renderers.ErrorBadRequest(errors.New(msg)))
		return
	}

	if err := rp.hasRoleAccess(roleID, projectID); err != nil {
		rp.logger.Error("role verification failed", zap.Error(err))
		render.Render(w, r, renderers.ErrorForbidden(fmt.Errorf("role verification failed")))
		return
	}

	var req UpdateAttributeRequest
	if err := render.Bind(r, &req); err != nil {
		msg := "invalid or incomplete request body"
		rp.logger.Error(msg, zap.Error(err))
		render.Render(w, r, renderers.ErrorBadRequest(errors.New(msg)))
		return
	}

	if err := rp.rolesDal.UpdatePermission(roleID, req.Resource, req.Path, req.Value); err != nil {
		msg := "failed to update permission attribute"
		rp.logger.Error(msg, zap.Error(err))
		render.Render(w, r, renderers.ErrorInternalServerError(errors.New(msg)))
		return
	}

	render.Status(r, http.StatusNoContent)
}

// @Summary Remove permission attribute
// @Description Removes a specific attribute from a permission
// @Tags permissions
// @Accept json
// @Produce json
// @Param role_id path string true "Role ID"
// @Param path path string true "Attribute path"
// @Success 204 "No Content"
// @Failure 400 {object} errorinterface.ErrorResponse
// @Failure 404 {object} errorinterface.ErrorResponse
// @Router /projects/{project_id}/roles/{role_id}/permissions [delete]
// @Security BearerAuth
func (rp *rolesService) RemovePermission(w http.ResponseWriter, r *http.Request) {
	projectID, _, err := rp.hasMemberAccess(r)
	if err != nil {
		msg := "project membership verification failed"
		rp.logger.Error(msg, zap.Error(err))
		render.Render(w, r, renderers.ErrorForbidden(errors.New(msg)))
		return
	}

	roleID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "role_id"))
	if err != nil {
		msg := "invalid role ID format"
		rp.logger.Error(msg, zap.Error(err))
		render.Render(w, r, renderers.ErrorBadRequest(errors.New(msg)))
		return
	}

	if err := rp.hasRoleAccess(roleID, projectID); err != nil {
		msg := "role verification failed"
		rp.logger.Error(msg, zap.Error(err))
		render.Render(w, r, renderers.ErrorForbidden(fmt.Errorf(msg)))
		return
	}

	var req UpdateAttributeRequest
	if err := render.Bind(r, &req); err != nil {
		msg := "invalid or incomplete request body"
		rp.logger.Error(msg, zap.Error(err))
		render.Render(w, r, renderers.ErrorBadRequest(errors.New(msg)))
		return
	}

	if err := rp.rolesDal.RemovePermission(roleID, req.Resource, req.Path); err != nil {
		msg := "failed to remove permission attribute"
		rp.logger.Error(msg, zap.Error(err))
		render.Render(w, r, renderers.ErrorNotFound(errors.New(msg)))
		return
	}

	render.Status(r, http.StatusNoContent)
}