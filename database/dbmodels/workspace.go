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
	Name        string               `json:"Name,omitempty" bson:"Name,omitempty"`
	Slug        string               `json:"Slug,omitempty" bson:"Slug,omitempty"` // Slug for URL or unique identification
	Description string               `json:"Description,omitempty" bson:"Description,omitempty"`
	OwnerID     primitive.ObjectID   `json:"OwnerID,omitempty" bson:"OwnerID,omitempty"` // Workspace owner
	Members     []primitive.ObjectID `json:"Members,omitempty" bson:"Members,omitempty"` // List of members belonging to the workspace

	// Authentication settings
	PrimaryAuthEnabled   bool `json:"PrimaryAuthEnabled" bson:"PrimaryAuthEnabled"`     // Flag for primary authentication
	SecondaryAuthEnabled bool `json:"SecondaryAuthEnabled" bson:"SecondaryAuthEnabled"` // Flag for secondary authentication
	RequireMFA           bool `json:"RequireMFA" bson:"RequireMFA"`                     // Flag for multi-factor authentication requirement
	AllowSecondaryAuth   bool `json:"AllowSecondaryAuth" bson:"AllowSecondaryAuth"`     // Flag for secondary authentication method allowance

	// Email domain settings
	AllowEmailInvites       bool     `json:"AllowEmailInvites" bson:"AllowEmailInvites"`               // Flag to allow email invites
	InviteDomainRestriction bool     `json:"InviteDomainRestriction" bson:"InviteDomainRestriction"`   // Restrict invites based on email domain
	AllowedDomains          []string `json:"AllowedDomains,omitempty" bson:"AllowedDomains,omitempty"` // List of allowed email domains for invites

	// Just-in-time (JIT) provisioning
	AllowJITProvisioning bool `json:"AllowJITProvisioning" bson:"AllowJITProvisioning"` // Flag to enable JIT provisioning for allowed domains
}

// Validate validates the workspace struct
func (w Workspace) Validate() error {
	return validation.ValidateStruct(&w,
		validation.Field(&w.Name, validation.Required),
		validation.Field(&w.Slug, validation.Required),
		validation.Field(&w.OwnerID, validation.Required),
		validation.Field(&w.Members, validation.Required),
	)
}
