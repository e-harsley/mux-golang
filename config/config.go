package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func getEnv(key, fallback string) string {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

type Settings struct {
	Env              string
	MonnifyEndpoint  string
	MonnifyKey       string
	MonnifyCC        string
	MonnifySecret    string
	MonnifyLoginUrl  string
	JwtSecretKey     string
	JwtAudienceClaim string
	JwtIssuerClaim   string
	TatumApiKey      string
	TatumBaseUrl     string
	TatumWebhookURL  string
}

func NewSettings() *Settings {
	return &Settings{
		Env:              getEnv("ENV", "staging"),
		MonnifyEndpoint:  getEnv("MONNIFY_ENDPOINT", "MONNIFY_ENDPOINT"),
		MonnifyKey:       getEnv("MONNIFY_KEY", "MONNIFY_KEY"),
		MonnifyCC:        getEnv("MONNIFY_KEY", "MONNIFY_CC"),
		MonnifySecret:    getEnv("MONNIFY_SECRET", "MONNIFY_SECRET"),
		JwtSecretKey:     getEnv("JWT_SECRET_KEY", "JWT_SECRET_KEY"),
		JwtAudienceClaim: getEnv("JWT_AUDIENCE_CLAIM", "JWT_AUDIENCE_CLAIM"),
		JwtIssuerClaim:   getEnv("JWT_ISSUER_CLAIM", "JWT_ISSUER_CLAIM"),
		TatumApiKey:      getEnv("TATUM_API_KEY", "t-64f72e1d0c34f3d88deb95f2-2c025722fe304ef69e98a834"),
		TatumBaseUrl:     getEnv("TATUM_BASE_URL", "https://api.tatum.io/v3/"),
		TatumWebhookURL:  getEnv("TATUM_WEBHOOK_URL", "https://dashboard.tatum.io"),
	}
}
