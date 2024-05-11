package models

import (
	"mux-crud/happiness"
)

type (
	Currency struct {
		happiness.BaseModel `bson:",inline"`
		Name                string                  `json:"name" bson:"name"`
		Code                string                  `json:"code" bson:"code"`
		CurrencyType        happiness.CurrencyTypes `json:"currency_type" bson:"currency_type"`
		Description         string                  `json:"description" bson:"description"`
	}

	SupportedNetwork struct {
		happiness.BaseModel `bson:",inline"`
		Name                string `json:"name" bson:"name"`
		Code                string `json:"code" bson:"code"`
		Description         string `json:"description,omitempty" bson:"description,omitempty"`
		Type                string `json:"type,omitempty" bson:"type,omitempty"` // bep 20, bitcoin
	}

	SupportedCurrencyNetwork struct {
		happiness.BaseModel `bson:",inline"`
		CurrencyID          string                  `json:"currency_id" bson:"currency_id"`
		CurrencyCode        string                  `json:"currency_code"`
		CurrencyType        happiness.CurrencyTypes `json:"currency_type" bson:"currency_type"`
		Currency            Currency                `json:"currency" bson:"currency"`
		NetworkCode         string                  `json:"network_code" bson:"network_code,omitempty"`
		NetworkName         string                  `json:"network_name" bson:"network_name,omitempty"`
		NetworkID           string                  `json:"network_id" bson:"network_id,omitempty"`
		Network             SupportedNetwork        `json:"network" bson:"network,omitempty"`
		Active              bool                    `json:"active" bson:"active"`
	}
)

func (c Currency) GetModelName() string {
	return "currency"
}

func (scn SupportedCurrencyNetwork) GetModelName() string {
	return "supported_currency_network"
}

func (sn SupportedNetwork) GetModelName() string {
	return "supported_network"
}
