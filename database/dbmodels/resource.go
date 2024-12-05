package dbmodels

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Resource represents the resource model
type Resource struct {
	ID                  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedTimestampUTC time.Time          `json:"created_timestamp_utc,omitempty" bson:"CreatedTimestampUTC,omitempty"`
	UpdatedTimestampUTC time.Time          `json:"updated_timestamp_utc,omitempty" bson:"UpdatedTimestampUTC,omitempty"`

	OwnerID   string             `json:"owner_id,omitempty" bson:"OwnerID,omitempty"` // Project owner (user ID)
	ProjectID primitive.ObjectID `json:"project_id,omitempty" bson:"ProjectID,omitempty"`

	Type    string `json:"type,omitempty" bson:"Type,omitempty"`
	Version string `json:"version,omitempty" bson:"Version,omitempty"`

	Name        string   `json:"name,omitempty" bson:"Name"`
	Description string   `json:"description,omitempty" bson:"Description"`
	Actions     []Action `json:"actions,omitempty" bson:"Actions"`
	URN         string   `json:"urn,omitempty" bson:"URN"`

	// Audit logs for project actions
	AuditLogs []AuditLog `json:"audit_logs,omitempty" bson:"AuditLogs,omitempty"` // Logs of actions taken within the project
	Deleted   bool       `json:"deleted,omitempty" bson:"Deleted,omitempty"`      // Flag to indicate if the project is deleted
}

// Action represents an action associated with a resource
type Action struct {
	Action      string `json:"action,omitempty" bson:"Action"`
	Description string `json:"description,omitempty" bson:"Description,omitempty"`
}

// Validate validates the workspace struct
func (r Resource) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Description, validation.Required),
		validation.Field(&r.Actions, validation.Required),
		validation.Field(&r.URN, validation.Required),
		validation.Field(&r.OwnerID, validation.Required),
		validation.Field(&r.ProjectID, validation.Required),
	)
}
