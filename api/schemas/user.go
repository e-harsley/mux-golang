package schemas

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"mux-crud/happiness"
)

type (
	SignupRequest struct {
		FirstName        string                              `json:"first_name" validate:"required"`
		LastName         string                              `json:"last_name" validate:"required"`
		Name             string                              `json:"name"`
		Country          string                              `json:"country" validate:"required"`
		Phone            string                              `json:"phone" validate:"required"`
		Email            string                              `json:"email" validate:"required"`
		BusinessName     string                              `json:"business_name" validate:"required"`
		BusinessType     happiness.BusinessTypes             `json:"business_type" validate:"required"`
		RegistrationType happiness.BusinessRegistrationTypes `json:"registration_type"`
		Password         string                              `json:"password" validate:"required"`
		VerifyPassword   string                              `json:"verify_password" validate:"required"`
	}

	LoginRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
		Otp      string `json:"otp"`
	}

	UserResponse struct {
		ID        string `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Name      string `json:"name"`
		Country   string `json:"country"`
		Phone     string `json:"phone"`
		Email     string `json:"email"`
	}

	AuthResponse struct {
		Status    happiness.AuthStatuses `json:"status"`
		ID        string                 `json:"id"`
		FirstName string                 `json:"first_name"`
		LastName  string                 `json:"last_name"`
		Name      string                 `json:"name"`
		Country   string                 `json:"country"`
		Phone     string                 `json:"phone"`
		Email     string                 `json:"email"`
		Token     string                 `json:"token"`
	}
)

func (cls SignupRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(cls)
	if err != nil {
		return happiness.NewValidationError(err.(validator.ValidationErrors))
	}

	if cls.Password != cls.VerifyPassword {
		return fmt.Errorf("password does not match verify password")
	}

	if cls.BusinessType == happiness.INDIVIDUAL {
		return nil
	}

	if cls.RegistrationType == "" {
		return fmt.Errorf("registration type is required")
	}
	return nil
}

func (cls LoginRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(cls)
	if err != nil {
		return happiness.NewValidationError(err.(validator.ValidationErrors))
	}
	return nil
}
