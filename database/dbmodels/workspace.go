package dbmodels

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Workspace model represents the workspace collection in the database
type Workspace struct {
	ID                  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedTimestampUTC time.Time          `json:"CreatedTimestampUTC,omitempty" bson:"CreatedTimestampUTC,omitempty"`
	UpdatedTimestampUTC time.Time          `json:"UpdatedTimestampUTC,omitempty" bson:"UpdatedTimestampUTC,omitempty"`

	// Workspace details
	Name        string   `json:"Name,omitempty" bson:"Name,omitempty"`
	Slug        string   `json:"Slug,omitempty" bson:"Slug,omitempty"`
	Description string   `json:"Description,omitempty" bson:"Description,omitempty"`
	OwnerID     string   `json:"OwnerID,omitempty" bson:"OwnerID,omitempty"`
	Members     []string `json:"Members,omitempty" bson:"Members,omitempty"`
	Deleted     bool     `json:"Deleted,omitempty" bson:"Deleted,omitempty"`

	// Audit logs for project actions
	AuditLogs []AuditLog `json:"AuditLogs,omitempty" bson:"AuditLogs,omitempty"` // Logs of actions taken within the project
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
