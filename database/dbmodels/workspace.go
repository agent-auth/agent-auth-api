package dbmodels

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Workspace model represents the workspace collection in the database
type Workspace struct {
	ID                  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedTimestampUTC time.Time          `json:"created_timestamp_utc,omitempty" bson:"CreatedTimestampUTC,omitempty"`
	UpdatedTimestampUTC time.Time          `json:"updated_timestamp_utc,omitempty" bson:"UpdatedTimestampUTC,omitempty"`

	// Workspace details
	Name        string   `json:"name,omitempty" bson:"Name,omitempty"`
	Slug        string   `json:"slug,omitempty" bson:"Slug,omitempty"`
	Description string   `json:"description,omitempty" bson:"Description,omitempty"`
	OwnerID     string   `json:"owner_id,omitempty" bson:"OwnerID,omitempty"`
	Members     []string `json:"members,omitempty" bson:"Members,omitempty"`
	Deleted     bool     `json:"deleted,omitempty" bson:"Deleted,omitempty"`

	// Audit logs for project actions
	AuditLogs []AuditLog `json:"audit_logs,omitempty" bson:"AuditLogs,omitempty"` // Logs of actions taken within the project
}

// Validate validates the workspace struct
func (w Workspace) Validate() error {
	return validation.ValidateStruct(&w,
		validation.Field(&w.Name, validation.Required),
		validation.Field(&w.Slug, validation.Required),
		validation.Field(&w.Description, validation.Required),
		validation.Field(&w.OwnerID, validation.Required),
		validation.Field(&w.Members, validation.Required),
	)
}
