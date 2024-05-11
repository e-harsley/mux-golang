package models

import "mux-crud/happiness"

type (
	Business struct {
		happiness.BaseModel  `bson:",inline"`
		Name                 string                                  `json:"name" bson:"name"`
		Address              Address                                 `json:"address,omitempty" bson:"address,omitempty"`
		UserID               string                                  `json:"user_id" bson:"user_id"`
		User                 UserMeta                                `json:"user" bson:"user"`
		Type                 happiness.BusinessTypes                 `json:"type" bson:"type"`
		IsActive             *bool                                   `json:"is_active" bson:"is_active"`
		RegistrationType     happiness.BusinessRegistrationTypes     `json:"registration_type,omitempty" bson:"registration_type,omitempty"`
		EstimatedTransaction happiness.BusinessEstimatedTransactions `json:"estimated_transaction,omitempty" bson:"estimated_transaction,omitempty"`
	}
)

func (u Business) GetModelName() string {
	return "business"
}
