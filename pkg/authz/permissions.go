package authz

// Role represents a user role in the system
type Role string

const (
	SystemAdmin    Role = "system_admin"    // Can manage everything
	WorkspaceAdmin Role = "workspace_admin" // Can manage a specific workspace
	AppAdmin       Role = "app_admin"       // Can manage specific apps
	AppDeveloper   Role = "app_developer"   // Can develop and deploy apps
	AppViewer      Role = "app_viewer"      // Can view app details
)

// Permission represents a specific action that can be performed
type Permission string

const (
	// Workspace permissions
	WorkspaceCreate Permission = "workspace:create"
	WorkspaceRead   Permission = "workspace:read"
	WorkspaceUpdate Permission = "workspace:update"
	WorkspaceDelete Permission = "workspace:delete"
	WorkspaceList   Permission = "workspace:list"

	// App permissions
	AppCreate Permission = "app:create"
	AppRead   Permission = "app:read"
	AppUpdate Permission = "app:update"
	AppDelete Permission = "app:delete"
	AppList   Permission = "app:list"
	AppDeploy Permission = "app:deploy"
)

// RolePermissions maps roles to their allowed permissions
var RolePermissions = map[Role][]Permission{
	SystemAdmin: {
		WorkspaceCreate, WorkspaceRead, WorkspaceUpdate, WorkspaceDelete, WorkspaceList,
		AppCreate, AppRead, AppUpdate, AppDelete, AppList, AppDeploy,
	},
	WorkspaceAdmin: {
		WorkspaceRead, WorkspaceUpdate,
		AppCreate, AppRead, AppUpdate, AppDelete, AppList, AppDeploy,
	},
	AppAdmin: {
		WorkspaceRead,
		AppRead, AppUpdate, AppDelete, AppList, AppDeploy,
	},
	AppDeveloper: {
		WorkspaceRead,
		AppRead, AppUpdate, AppList, AppDeploy,
	},
	AppViewer: {
		WorkspaceRead,
		AppRead, AppList,
	},
}
