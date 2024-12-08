package projects

import (
	"net/http"

	"errors"

	"github.com/agent-auth/agent-auth-api/pkg/authz"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Helper function to check if a user has access to a project
func (rp *projectService) hasMemberAccess(r *http.Request) (bson.ObjectID, string, error) {
	projectID, err := bson.ObjectIDFromHex(chi.URLParam(r, "project_id"))
	if err != nil {
		return bson.NilObjectID, "", errors.New("invalid project ID format")
	}

	email, err := authz.GetEmailFromClaims(r)
	if err != nil {
		return bson.NilObjectID, "", ErrUnauthorized
	}

	isMember, err := rp.projectDal.IsMember(projectID, email)
	if err != nil {
		return bson.NilObjectID, "", ErrInternalServerError
	}

	if !isMember {
		return bson.NilObjectID, "", ErrUnauthorized
	}

	return projectID, email, nil
}
