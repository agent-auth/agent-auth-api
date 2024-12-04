package rolespermissions

import (
	"errors"
	"net/http"
)

// RolesService interface
type RolesService interface {
	CreateRole(w http.ResponseWriter, r *http.Request)
	GetRole(w http.ResponseWriter, r *http.Request)
	DeleteRole(w http.ResponseWriter, r *http.Request)
	GetRolesByProject(w http.ResponseWriter, r *http.Request)
	DeleteRolesByProject(w http.ResponseWriter, r *http.Request)

	UpdatePermission(w http.ResponseWriter, r *http.Request)
	RemovePermission(w http.ResponseWriter, r *http.Request)
}

// The list of error types presented to the end user
var (
	ErrIncompleteDetails = errors.New("incorrect details provided, please provide correct details")
	ErrNotFound          = errors.New("role not found")
)
