package resources

import (
	"net/http"
)

// ResourceService interface
type ResourceService interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	ListByProject(w http.ResponseWriter, r *http.Request)
}
