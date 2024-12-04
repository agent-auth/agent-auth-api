package projects

import (
	"errors"
	"net/http"
)

// ProjectService interface
type ProjectService interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	AddMember(w http.ResponseWriter, r *http.Request)
	RemoveMember(w http.ResponseWriter, r *http.Request)
}

// The list of error types presented to the end user
var (
	ErrIncompleteDetails   = errors.New("incorrect details provided, please provide correct details")
	ErrNotFound            = errors.New("project not found")
	ErrInternalServerError = errors.New("internal server error")
	ErrUnauthorized        = errors.New("unauthorized access")
)
