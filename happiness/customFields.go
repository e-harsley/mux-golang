package happiness

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"mux-crud/config"
	"time"
)

const bindPrefix = "context_bind:::"

var setting = config.NewSettings()

type AuthToken struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func NewAuthToken(userId string) *AuthToken {
	return &AuthToken{
		UserID: userId,
	}
}

func (cls AuthToken) Token() (string, error) {

	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = uuid.New().String()
	atClaims["user_id"] = cls.UserID
	atClaims["aud"] = setting.JwtAudienceClaim
	atClaims["iss"] = setting.JwtIssuerClaim
	atClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	return at.SignedString([]byte(setting.JwtSecretKey))
}

func (cls AuthToken) ParseAuthToken(tokenString string) (*AuthToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &cls, func(token *jwt.Token) (interface{}, error) {
		return []byte(setting.JwtSecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AuthToken); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
