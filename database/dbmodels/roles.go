package dbmodels

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Roles represents roles for a specific resource
type Roles struct {
	ID                  primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	ProjectID           primitive.ObjectID     `json:"project_id" bson:"project_id"`
	CreatedTimestampUTC time.Time              `json:"created_timestamp_utc" bson:"created_timestamp_utc"`
	UpdatedTimestampUTC time.Time              `json:"updated_timestamp_utc" bson:"updated_timestamp_utc"`
	Role                string                 `json:"role" bson:"role"`
	Permissions         map[string]interface{} `json:"permissions" bson:"permissions"`
}
