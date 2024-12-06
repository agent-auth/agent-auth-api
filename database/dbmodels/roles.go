package dbmodels

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Action represents an action associated with a resource
type Action struct {
	Action      string   `json:"action,omitempty" bson:"Action"`
	Description string   `json:"description,omitempty" bson:"Description,omitempty"`
	Actions     []Action `json:"actions,omitempty" bson:"Actions,omitempty"`
}

type Permission struct {
	Actions []Action `json:"actions,omitempty" bson:"Actions,omitempty"`
}

// Roles represents roles for a specific resource
type Roles struct {
	ID                  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OwnerID             string             `json:"owner_id,omitempty" bson:"OwnerID,omitempty"`
	ProjectID           primitive.ObjectID `json:"project_id,omitempty" bson:"ProjectID,omitempty"`
	CreatedTimestampUTC time.Time          `json:"created_timestamp_utc,omitempty" bson:"CreatedTimestampUTC,omitempty"`
	UpdatedTimestampUTC time.Time          `json:"updated_timestamp_utc,omitempty" bson:"UpdatedTimestampUTC,omitempty"`

	Role        string `json:"role,omitempty" bson:"Role,omitempty"`
	Description string `json:"description,omitempty" bson:"Description,omitempty"`

	Permissions map[string]Permission `json:"permissions,omitempty" bson:"Permissions,omitempty"`
	AuditLogs   []AuditLog            `json:"audit_logs,omitempty" bson:"AuditLogs,omitempty"`
	Deleted     bool                  `json:"deleted,omitempty" bson:"Deleted,omitempty"`
}

// Validate validates the workspace struct
func (r Roles) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Role, validation.Required),
		validation.Field(&r.OwnerID, validation.Required),
		validation.Field(&r.ProjectID, validation.Required),
		validation.Field(&r.Description, validation.Required),
	)
}
