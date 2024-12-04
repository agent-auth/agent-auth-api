package rolespermissions

import (
	"net/http"

	"github.com/agent-auth/agent-auth-api/pkg/authz"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Helper function to verify project membership
func (rp *rolesService) verifyProjectMembership(r *http.Request) (primitive.ObjectID, string, error) {
	projectID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "project_id"))
	if err != nil {
		return primitive.NilObjectID, "", ErrIncompleteDetails
	}

	email, err := authz.GetEmailFromClaims(r)
	if err != nil {
		return primitive.NilObjectID, "", ErrUnauthorized
	}

	isMember, err := rp.projectsDal.IsMember(projectID.String(), email)
	if err != nil {
		return primitive.NilObjectID, "", ErrInternalServerError
	}

	if !isMember {
		return primitive.NilObjectID, "", ErrUnauthorized
	}

	return projectID, email, nil
}

// Helper function to verify role belongs to project
func (rp *rolesService) verifyRoleInProject(roleID, projectID primitive.ObjectID) error {
	role, err := rp.rolesDal.Get(roleID)
	if err != nil {
		return ErrNotFound
	}
	if role.ProjectID != projectID {
		return ErrInvalidRole
	}
	return nil
}
