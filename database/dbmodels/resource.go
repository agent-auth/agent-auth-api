package dbmodels

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Resource represents the resource model
type Resource struct {
	ID                  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedTimestampUTC time.Time          `json:"CreatedTimestampUTC,omitempty" bson:"CreatedTimestampUTC,omitempty"`
	UpdatedTimestampUTC time.Time          `json:"UpdatedTimestampUTC,omitempty" bson:"UpdatedTimestampUTC,omitempty"`

	OwnerID   string             `json:"OwnerID,omitempty" bson:"OwnerID,omitempty"` // Project owner (user ID)
	ProjectID primitive.ObjectID `json:"project_id,omitempty" bson:"project_id,omitempty"`

	ResourceName string   `bson:"resourceId" json:"resourceId"`
	Description  string   `bson:"description" json:"description"`
	Actions      []Action `bson:"actions" json:"actions"`
	URN          string   `bson:"urn" json:"urn"`

	// Audit logs for project actions
	AuditLogs []AuditLog `json:"AuditLogs,omitempty" bson:"AuditLogs,omitempty"` // Logs of actions taken within the project
	Deleted   bool       `json:"Deleted,omitempty" bson:"Deleted,omitempty"`     // Flag to indicate if the project is deleted
}

// Action represents an action associated with a resource
type Action struct {
	Action      string `bson:"action" json:"action"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`
}

// Validate validates the workspace struct
func (r Resource) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ResourceName, validation.Required),
		validation.Field(&r.Description, validation.Required),
		validation.Field(&r.Actions, validation.Required),
		validation.Field(&r.URN, validation.Required),
		validation.Field(&r.OwnerID, validation.Required),
		validation.Field(&r.ProjectID, validation.Required),
	)
}
