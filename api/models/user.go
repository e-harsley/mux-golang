package models

import (
	"github.com/xlzd/gotp"
	"golang.org/x/crypto/bcrypt"
	"mux-crud/happiness"
	"time"
)

type (
	UserMeta struct {
		Name      string `json:"name" bson:"name"`
		FirstName string `json:"first_name" bson:"first_name"`
		LastName  string `json:"last_name" bson:"last_name"`
		Email     string `json:"email" bson:"email"`
		Phone     string `json:"phone" bson:"phone"`
	}

	User struct {
		happiness.BaseModel `bson:",inline"`
		FirstName           string                 `json:"first_name" bson:"first_name,required"`
		LastName            string                 `json:"last_name" bson:"last_name,required"`
		Name                string                 `json:"name" bson:"name"`
		Email               string                 `json:"email" bson:"email,required"`
		Phone               string                 `json:"phone" bson:"phone,required"`
		Otp2faString        string                 `json:"otp2fa_secret" bson:"otp2fa_secret,omitempty"`
		OtpProvider         happiness.OTPProviders `json:"otp_provider" bson:"otp_provider,omitempty"`
		IsEmailVerified     *bool                  `json:"is_email_verified,omitempty" bson:"is_email_verified"`
		IsPhoneVerified     *bool                  `json:"is_phone_verified,omitempty" bson:"is_phone_verified"`
		Country             string                 `json:"country" bson:"country"`
		Password            string                 `json:"password" bson:"password,required"`
		IsSuspended         *bool                  `json:"is_suspended,omitempty" bson:"is_suspended"`
		Is2faEnabled        *bool                  `json:"is_2_fa_enabled,omitempty" bson:"is_2_fa_enabled"`
	}
)

func (u User) GetModelName() string {
	return "user"
}

func (cls *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	cls.Password = string(hashedPassword)

	return nil
}

func (cls *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(cls.Password), []byte(password))
	return err == nil
}

func (cls *User) Generate2faSecret() {
	cls.Otp2faString = gotp.RandomSecret(16)
}

func (cls *User) ValidateOTP(otp string) bool {
	totp := gotp.NewDefaultTOTP(cls.Otp2faString)
	return totp.Verify(otp, time.Now().Unix())
}

func (cls *User) GenerateOTP() string {
	totp := gotp.NewDefaultTOTP(cls.Otp2faString)
	return totp.Now()
}

func (cls *User) UserMeta() (user UserMeta) {
	user.FirstName = cls.FirstName
	user.LastName = cls.LastName
	user.Phone = cls.Phone
	user.Email = cls.Email
	user.Name = cls.Name
	return user
}
