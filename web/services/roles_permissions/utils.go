package roles_permissions

import (
	"net/http"

	"errors"

	"github.com/agent-auth/agent-auth-api/pkg/authz"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Helper function to verify project membership
func (rp *rolesService) hasMemberAccess(r *http.Request) (primitive.ObjectID, string, error) {
	projectID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "project_id"))
	if err != nil {
		return primitive.NilObjectID, "", errors.New("invalid project ID")
	}

	email, err := authz.GetEmailFromClaims(r)
	if err != nil {
		return primitive.NilObjectID, "", errors.New("unauthorized access")
	}

	isMember, err := rp.projectsDal.IsMember(projectID, email)
	if err != nil {
		return primitive.NilObjectID, "", errors.New("failed to verify project membership")
	}

	if !isMember {
		return primitive.NilObjectID, "", errors.New("unauthorized access - not a project member")
	}

	return projectID, email, nil
}

// Helper function to verify role belongs to project
func (rp *rolesService) hasRoleAccess(roleID, projectID primitive.ObjectID) error {
	role, err := rp.rolesDal.Get(roleID)
	if err != nil {
		return errors.New("role not found")
	}
	if role.ProjectID != projectID {
		return errors.New("role does not belong to this project")
	}
	return nil
}
