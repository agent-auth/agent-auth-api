package dbmodels

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Permissions struct {
	Resource map[string]interface{} `json:"resource" bson:"resource"`
}

// Roles represents roles for a specific resource
type Roles struct {
	ID                  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OwnerID             string             `json:"owner_id,omitempty" bson:"owner_id,omitempty"`
	ProjectID           primitive.ObjectID `json:"project_id,omitempty" bson:"project_id,omitempty"`
	CreatedTimestampUTC time.Time          `json:"created_timestamp_utc,omitempty" bson:"created_timestamp_utc,omitempty"`
	UpdatedTimestampUTC time.Time          `json:"updated_timestamp_utc,omitempty" bson:"updated_timestamp_utc,omitempty"`

	Role        string `json:"role,omitempty" bson:"role,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`

	Permissions Permissions `json:"permissions,omitempty" bson:"permissions,omitempty"`
	AuditLogs   []AuditLog  `json:"audit_logs,omitempty" bson:"audit_logs,omitempty"`
	Deleted     bool        `json:"deleted,omitempty" bson:"deleted,omitempty"`
}

// Validate validates the workspace struct
func (r Roles) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Role, validation.Required),
		validation.Field(&r.Permissions, validation.Required),
		validation.Field(&r.OwnerID, validation.Required),
		validation.Field(&r.ProjectID, validation.Required),
		validation.Field(&r.Description, validation.Required),
	)
}
