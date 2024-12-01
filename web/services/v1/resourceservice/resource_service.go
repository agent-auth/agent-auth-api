package resourceservice

import (
	"errors"
	"net/http"
	"strings"

	"github.com/agent-auth/agent-auth-api/pkg/keycloak"
	"github.com/agent-auth/agent-auth-api/web/renderers"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

// ResourceService interface
type ResourceService interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
}

type resourceService struct {
	client    *keycloak.Client
	resources *keycloak.ResourceService
	logger    *zap.Logger
}

// Error types
var (
	ErrIncompleteDetails = errors.New("incorrect details provided, please provide correct details")
	ErrResourceNotFound  = errors.New("resource not found")
)

// Error codes
const (
	FailedToCreateResource = "Failed-To-Create-Resource"
	FailedToGetResource    = "Failed-To-Get-Resource"
	FailedToUpdateResource = "Failed-To-Update-Resource"
	FailedToDeleteResource = "Failed-To-Delete-Resource"
	FailedToListResources  = "Failed-To-List-Resources"
	ErrMissingWorkspaceID  = "workspace_id is required"
	ErrMissingProjectID    = "project_id is required"
	ErrMissingResourceID   = "resource_id is required"
	ErrQueryParamTooLong   = "query parameter exceeds maximum length"
	MaxNameLength          = 255
	MaxURILength           = 1024
)

// ResourceRequest represents the incoming resource request
type ResourceRequest struct {
	*keycloak.ResourceRepresentation
}

func (r *ResourceRequest) Bind(req *http.Request) error {
	if r.ResourceRepresentation == nil {
		return ErrIncompleteDetails
	}
	if r.Name == "" || len(r.Name) > MaxNameLength {
		return errors.New("invalid resource name")
	}

	return nil
}

// ResourceResponse represents the outgoing resource response
type ResourceResponse struct {
	*keycloak.ResourceRepresentation
}

func NewResourceService(keycloakURL, token string) (ResourceService, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	client, err := keycloak.NewClient(keycloakURL, token)
	if err != nil {
		return nil, err
	}

	return &resourceService{
		client:    client,
		resources: keycloak.NewResourceService(client),
		logger:    logger,
	}, nil
}

func validateCommonParams(_ http.ResponseWriter, r *http.Request) (string, string, error) {
	realmName := chi.URLParam(r, "workspace_id")
	clientID := chi.URLParam(r, "project_id")

	if realmName == "" {
		return "", "", errors.New(ErrMissingWorkspaceID)
	}
	if clientID == "" {
		return "", "", errors.New(ErrMissingProjectID)
	}

	return realmName, clientID, nil
}

// @Summary Create a new resource
// @Description Creates a new Keycloak resource
// @Tags resources
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param resource body ResourceRequest true "Resource Details"
// @Success 200 {object} ResourceResponse
// @Failure 400 {object} errorinterface.ErrorResponse
// @Failure 500 {object} errorinterface.ErrorResponse
// @Router /{workspace_id}/{project_id}/oauth/resources [POST]
func (rs *resourceService) Create(w http.ResponseWriter, r *http.Request) {
	realmName, clientID, err := validateCommonParams(w, r)
	if err != nil {
		render.Render(w, r, renderers.ErrorBadRequest(err))
		return
	}

	resource := &ResourceRequest{}
	if err := render.Bind(r, resource); err != nil {
		rs.logger.Error("failed to bind resource request",
			zap.Error(err),
		)
		render.Render(w, r, renderers.ErrorBadRequest(ErrIncompleteDetails))
		return
	}

	created, err := rs.resources.CreateAuthzResource(realmName, clientID, resource.ResourceRepresentation)
	if err != nil {
		rs.logger.Error("failed to create resource",
			zap.Error(err),
			zap.String("realm", realmName),
			zap.String("clientID", clientID),
		)
		render.Render(w, r, renderers.ErrorInternalServerError(err))
		return
	}

	render.Respond(w, r, &ResourceResponse{ResourceRepresentation: created})
}

// @Summary Get a resource
// @Description Retrieves a Keycloak resource by ID
// @Tags resources
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param project_id path string true "Project ID"
// @Param resource_id path string true "Resource ID"
// @Success 200 {object} ResourceResponse
// @Failure 404 {object} errorinterface.ErrorResponse
// @Failure 500 {object} errorinterface.ErrorResponse
// @Router /{workspace_id}/{project_id}/oauth/resources/{resource_id} [GET]
func (rs *resourceService) Get(w http.ResponseWriter, r *http.Request) {
	realmName, clientID, err := validateCommonParams(w, r)
	if err != nil {
		render.Render(w, r, renderers.ErrorBadRequest(err))
		return
	}

	resourceID := chi.URLParam(r, "resource_id")
	if resourceID == "" {
		render.Render(w, r, renderers.ErrorBadRequest(errors.New(ErrMissingResourceID)))
		return
	}

	resource, err := rs.resources.GetAuthzResource(realmName, clientID, resourceID)
	if err != nil {
		rs.logger.Error("failed to get resource",
			zap.Error(err),
			zap.String("realm", realmName),
			zap.String("clientID", clientID),
			zap.String("resourceID", resourceID),
		)
		render.Render(w, r, renderers.ErrorNotFound(ErrResourceNotFound))
		return
	}

	render.Respond(w, r, &ResourceResponse{ResourceRepresentation: resource})
}

// @Summary Update a resource
// @Description Updates an existing Keycloak resource
// @Tags resources
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param project_id path string true "Project ID"
// @Param resource_id path string true "Resource ID"
// @Param resource body ResourceRequest true "Resource Details"
// @Success 200 {object} ResourceResponse
// @Failure 400,404 {object} errorinterface.ErrorResponse
// @Failure 500 {object} errorinterface.ErrorResponse
// @Router /{workspace_id}/{project_id}/oauth/{resource_id} [PUT]
func (rs *resourceService) Update(w http.ResponseWriter, r *http.Request) {
	resourceID := chi.URLParam(r, "resource_id")
	realmName := chi.URLParam(r, "workspace_id")
	clientID := chi.URLParam(r, "project_id")

	resource := &ResourceRequest{}
	if err := render.Bind(r, resource); err != nil {
		rs.logger.Error("failed to bind resource request",
			zap.Error(err),
		)
		render.Render(w, r, renderers.ErrorBadRequest(ErrIncompleteDetails))
		return
	}

	err := rs.resources.UpdateAuthzResource(realmName, clientID, resourceID, resource.ResourceRepresentation)
	if err != nil {
		rs.logger.Error("failed to update resource",
			zap.Error(err),
			zap.String("realm", realmName),
			zap.String("clientID", clientID),
			zap.String("resourceID", resourceID),
		)
		render.Render(w, r, renderers.ErrorInternalServerError(err))
		return
	}

	render.Respond(w, r, &ResourceResponse{ResourceRepresentation: resource.ResourceRepresentation})
}

// @Summary Delete a resource
// @Description Deletes a Keycloak resource
// @Tags resources
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param project_id path string true "Project ID"
// @Param resource_id path string true "Resource ID"
// @Success 204 "No Content"
// @Failure 404 {object} errorinterface.ErrorResponse
// @Failure 500 {object} errorinterface.ErrorResponse
// @Router /{workspace_id}/{project_id}/oauth/resources/{resource_id} [DELETE]
func (rs *resourceService) Delete(w http.ResponseWriter, r *http.Request) {
	resourceID := chi.URLParam(r, "resource_id")
	realmName := chi.URLParam(r, "workspace_id")
	clientID := chi.URLParam(r, "project_id")

	err := rs.resources.DeleteAuthzResource(realmName, clientID, resourceID)
	if err != nil {
		rs.logger.Error("failed to delete resource",
			zap.Error(err),
			zap.String("realm", realmName),
			zap.String("clientID", clientID),
			zap.String("resourceID", resourceID),
		)
		render.Render(w, r, renderers.ErrorInternalServerError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary List resources
// @Description Lists all Keycloak resources with optional filtering
// @Tags resources
// @Accept json
// @Produce json
// @Param name query string false "Filter by name"
// @Param type query string false "Filter by type"
// @Param uri query string false "Filter by URI"
// @Param owner query string false "Filter by owner"
// @Param scope query string false "Filter by scope"
// @Param workspace_id path string true "Workspace ID"
// @Param project_id path string true "Project ID"
// @Success 200 {array} ResourceResponse
// @Failure 500 {object} errorinterface.ErrorResponse
// @Router /{workspace_id}/{project_id}/oauth/resources [GET]
func (rs *resourceService) List(w http.ResponseWriter, r *http.Request) {
	realmName, clientID, err := validateCommonParams(w, r)
	if err != nil {
		render.Render(w, r, renderers.ErrorBadRequest(err))
		return
	}

	query := &keycloak.ResourceQuery{
		Name:  strings.TrimSpace(r.URL.Query().Get("name")),
		Type:  strings.TrimSpace(r.URL.Query().Get("type")),
		URI:   strings.TrimSpace(r.URL.Query().Get("uri")),
		Owner: strings.TrimSpace(r.URL.Query().Get("owner")),
		Scope: strings.TrimSpace(r.URL.Query().Get("scope")),
	}

	if len(query.Name) > MaxNameLength || len(query.Type) > MaxNameLength ||
		len(query.URI) > MaxURILength || len(query.Owner) > MaxNameLength ||
		len(query.Scope) > MaxNameLength {
		render.Render(w, r, renderers.ErrorBadRequest(errors.New(ErrQueryParamTooLong)))
		return
	}

	resources, err := rs.resources.GetAuthzResources(realmName, clientID, query)
	if err != nil {
		rs.logger.Error("failed to list resources",
			zap.Error(err),
			zap.String("realm", realmName),
			zap.String("clientID", clientID),
			zap.Any("query", query),
		)
		render.Render(w, r, renderers.ErrorInternalServerError(err))
		return
	}

	response := make([]ResourceResponse, len(resources))
	for i, resource := range resources {
		response[i] = ResourceResponse{ResourceRepresentation: &resource}
	}

	render.Respond(w, r, response)
}