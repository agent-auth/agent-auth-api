package dbmodels

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Project model represents the project collection in the database
type Project struct {
	ID                  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedTimestampUTC time.Time          `json:"CreatedTimestampUTC,omitempty" bson:"CreatedTimestampUTC,omitempty"`
	UpdatedTimestampUTC time.Time          `json:"UpdatedTimestampUTC,omitempty" bson:"UpdatedTimestampUTC,omitempty"`

	// Project details
	Name        string               `json:"Name,omitempty" bson:"Name,omitempty"`               // Name of the project
	Description string               `json:"Description,omitempty" bson:"Description,omitempty"` // Description of the project
	Slug        string               `json:"Slug,omitempty" bson:"Slug,omitempty"`               // Slug for unique URL identification
	WorkspaceID primitive.ObjectID   `json:"WorkspaceID,omitempty" bson:"WorkspaceID,omitempty"` // Reference to the workspace it belongs to
	OwnerID     primitive.ObjectID   `json:"OwnerID,omitempty" bson:"OwnerID,omitempty"`         // Project owner (user ID)
	Members     []primitive.ObjectID `json:"Members,omitempty" bson:"Members,omitempty"`         // List of member IDs associated with the project
	Status      string               `json:"Status,omitempty" bson:"Status,omitempty"`           // Status of the project (e.g., "Active", "Completed")

	// Project settings
	IsPrivate           bool   `json:"IsPrivate" bson:"IsPrivate"`                                         // Whether the project is private or public
	BillingEnabled      bool   `json:"BillingEnabled" bson:"BillingEnabled"`                               // If billing is enabled for this project
	EstimatedCompletion string `json:"EstimatedCompletion,omitempty" bson:"EstimatedCompletion,omitempty"` // Estimated project completion date
	AllowExternalAccess bool   `json:"AllowExternalAccess" bson:"AllowExternalAccess"`                     // Allow external users to access the project

	// Authentication & access settings
	AllowProjectAuth     bool              `json:"AllowProjectAuth" bson:"AllowProjectAuth"`         // Whether project-specific authentication is enabled
	RequireMFA           bool              `json:"RequireMFA" bson:"RequireMFA"`                     // Flag for multi-factor authentication requirement
	EnableSsoIntegration bool              `json:"EnableSsoIntegration" bson:"EnableSsoIntegration"` // Whether SSO (Single Sign-On) is enabled
	AccessLevels         map[string]string `json:"AccessLevels" bson:"AccessLevels"`                 // Custom access levels (e.g., "Admin", "Member", "Viewer")

	// Email notifications settings
	EnableEmailNotifications bool     `json:"EnableEmailNotifications" bson:"EnableEmailNotifications"`                   // Flag for email notifications
	NotificationPreferences  []string `json:"NotificationPreferences,omitempty" bson:"NotificationPreferences,omitempty"` // Email notification settings for users (e.g., "TaskUpdates", "Milestones")

	// Time tracking & project phases
	Phases []ProjectPhase `json:"Phases,omitempty" bson:"Phases,omitempty"` // Phases of the project, such as planning, development, etc.

	// Audit logs for project actions
	AuditLogs []AuditLog `json:"AuditLogs,omitempty" bson:"AuditLogs,omitempty"` // Logs of actions taken within the project

}

// ProjectPhase represents a phase within a project (e.g., planning, development)
type ProjectPhase struct {
	Name        string    `json:"Name" bson:"Name"`               // Name of the phase
	StartDate   time.Time `json:"StartDate" bson:"StartDate"`     // Start date of the phase
	EndDate     time.Time `json:"EndDate" bson:"EndDate"`         // End date of the phase
	Status      string    `json:"Status" bson:"Status"`           // Status of the phase (e.g., "In Progress", "Completed")
	Description string    `json:"Description" bson:"Description"` // Description of the phase
}

// AuditLog represents a log of an action within the project
type AuditLog struct {
	Timestamp time.Time          `json:"Timestamp" bson:"Timestamp"` // When the action was performed
	Action    string             `json:"Action" bson:"Action"`       // Action taken (e.g., "Created Task", "Updated Milestone")
	UserID    primitive.ObjectID `json:"UserID" bson:"UserID"`       // The user who performed the action
	Details   string             `json:"Details" bson:"Details"`     // Detailed description of the action
}

// Validate validates the project struct
func (p Project) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Slug, validation.Required),
		validation.Field(&p.WorkspaceID, validation.Required),
		validation.Field(&p.OwnerID, validation.Required),
	)
}
