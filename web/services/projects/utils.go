package projects

import (
	"net/http"

	"errors"

	"github.com/agent-auth/agent-auth-api/pkg/authz"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Helper function to check if a user has access to a project
func (rp *projectService) hasMemberAccess(r *http.Request) (primitive.ObjectID, string, error) {
	projectID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "project_id"))
	if err != nil {
		return primitive.NilObjectID, "", errors.New("invalid project ID format")
	}

	email, err := authz.GetEmailFromClaims(r)
	if err != nil {
		return primitive.NilObjectID, "", ErrUnauthorized
	}

	isMember, err := rp.projectDal.IsMember(projectID, email)
	if err != nil {
		return primitive.NilObjectID, "", ErrInternalServerError
	}

	if !isMember {
		return primitive.NilObjectID, "", ErrUnauthorized
	}

	return projectID, email, nil
}
