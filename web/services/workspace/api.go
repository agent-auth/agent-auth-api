package workspace

import (
	"errors"
	"net/http"
)

// WorkspaceService interface
type WorkspaceService interface {
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
	ErrIncompleteDetails = errors.New("incorrect details provided, please provide correct details")
	ErrNotFound          = errors.New("workspace not found")
	ErrUnauthorized      = errors.New("unauthorized to perform this action")
)

// List of error codes
var (
	FailedToCreateWorkspace = "Failed-To-Create-Workspace"
	FailedToGetWorkspace    = "Failed-To-Get-Workspace"
	FailedToUpdateWorkspace = "Failed-To-Update-Workspace"
	FailedToDeleteWorkspace = "Failed-To-Delete-Workspace"
	FailedToListWorkspaces  = "Failed-To-List-Workspaces"
	FailedToAddMember       = "Failed-To-Add-Member"
	FailedToRemoveMember    = "Failed-To-Remove-Member"
)
