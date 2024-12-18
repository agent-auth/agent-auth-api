package resources

import (
	"errors"
	"net/http"

	"github.com/agent-auth/agent-auth-api/pkg/authz"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Helper function to verify project membership
func (rs *resourceService) hasMemberAccess(r *http.Request) (bson.ObjectID, string, error) {
	projectID, err := bson.ObjectIDFromHex(chi.URLParam(r, "project_id"))
	if err != nil {
		return bson.NilObjectID, "", errors.New("invalid project ID")
	}

	email, err := authz.GetEmailFromClaims(r)
	if err != nil {
		return bson.NilObjectID, "", errors.New("unauthorized access")
	}

	isMember, err := rs.projects_dal.IsMember(projectID, email)
	if err != nil {
		return bson.NilObjectID, "", errors.New("failed to verify project membership")
	}

	if !isMember {
		return bson.NilObjectID, "", errors.New("unauthorized access - not a project member")
	}

	return projectID, email, nil
}

// hasRoleAccess checks if the resource belongs to the specified project
func (rs *resourceService) hasRoleAccess(resource_id, project_id bson.ObjectID) error {
	existing, err := rs.resources_dal.GetByID(resource_id)
	if err != nil {
		return errors.New("resource not found")
	}

	if existing.ProjectID != project_id {
		return errors.New("unauthorized access attempt")
	}

	return nil
}
