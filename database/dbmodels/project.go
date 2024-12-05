package dbmodels

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Project model represents the project collection in the database
type Project struct {
	ID                  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedTimestampUTC time.Time          `json:"created_timestamp_utc,omitempty" bson:"CreatedTimestampUTC,omitempty"`
	UpdatedTimestampUTC time.Time          `json:"updated_timestamp_utc,omitempty" bson:"UpdatedTimestampUTC,omitempty"`

	// Project details
	Name        string             `json:"name,omitempty" bson:"Name,omitempty"`                // Name of the project
	Description string             `json:"description,omitempty" bson:"Description,omitempty"`  // Description of the project
	Slug        string             `json:"slug,omitempty" bson:"Slug,omitempty"`                // Slug for unique URL identification
	WorkspaceID primitive.ObjectID `json:"workspace_id,omitempty" bson:"WorkspaceID,omitempty"` // Reference to the workspace it belongs to
	OwnerID     string             `json:"owner_id,omitempty" bson:"OwnerID,omitempty"`         // Project owner (user ID)
	Members     []string           `json:"members,omitempty" bson:"Members,omitempty"`          // List of member IDs associated with the project

	// Audit logs for project actions
	AuditLogs []AuditLog `json:"audit_logs,omitempty" bson:"AuditLogs,omitempty"` // Logs of actions taken within the project
	Deleted   bool       `json:"deleted,omitempty" bson:"Deleted,omitempty"`      // Flag to indicate if the project is deleted
}

// AuditLog represents a log of an action within the project
type AuditLog struct {
	Timestamp time.Time          `json:"timestamp,omitempty" bson:"Timestamp"` // When the action was performed
	Action    string             `json:"action,omitempty" bson:"Action"`       // Action taken (e.g., "Created Task", "Updated Milestone")
	UserID    primitive.ObjectID `json:"user_id,omitempty" bson:"UserID"`      // The user who performed the action
	Details   string             `json:"details,omitempty" bson:"Details"`     // Detailed description of the action
}

// Validate validates the project struct
func (p Project) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Slug, validation.Required),
		validation.Field(&p.Description, validation.Required),
		// validation.Field(&p.WorkspaceID, validation.Required),
		validation.Field(&p.OwnerID, validation.Required),
		validation.Field(&p.Members, validation.Required),
	)
}
