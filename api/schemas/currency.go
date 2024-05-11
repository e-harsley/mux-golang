package schemas

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"mux-crud/happiness"
)

type (
	CurrencySchema struct {
		Name         string                  `json:"name" validate:"required"`
		Code         string                  `json:"code" validate:"required"`
		CurrencyType happiness.CurrencyTypes `json:"currency_type" validate:"required"`
		Description  string                  `json:"description"`
	}

	SupportNetworkSchema struct {
		Name        string `json:"name" validate:"required"`
		Code        string `json:"code" validate:"required"`
		Description string `json:"description"`
		Type        string `json:"type" validate:"required"` // bep 20, bitcoin
	}
	SupportedCurrencyNetworkSchema struct {
		CurrencyID   string                  `json:"currency_id" validate:"required"`
		CurrencyType happiness.CurrencyTypes `json:"currency_type" validate:"required"`
		Active       bool                    `json:"active" bson:"active"`
		NetworkID    string                  `json:"network_id"`
	}
)

func (cls CurrencySchema) Validate() error {
	validate := validator.New()
	err := validate.Struct(cls)
	if err != nil {
		return happiness.NewValidationError(err.(validator.ValidationErrors))
	}
	return nil
}

func (cls SupportedCurrencyNetworkSchema) Validate() error {
	validate := validator.New()
	err := validate.Struct(cls)
	if err != nil {
		return happiness.NewValidationError(err.(validator.ValidationErrors))
	}

	if cls.CurrencyType == happiness.CRYPTOCURRENCY && cls.NetworkID == "" {
		return fmt.Errorf("for currency of type crypto network id is required")
	}
	return nil
}
