package models

import "mux-crud/happiness"

type (
	Address struct {
		Street  string `json:"street" bson:"street,omitempty"`
		City    string `json:"city" bson:"city,omitempty"`
		State   string `json:"state" bson:"state,omitempty"`
		Country string `json:"country" bson:"country,omitempty"`
	}

	ApiCredential struct {
		happiness.BaseModel `bson:",inline"`
		UserID              string `json:"user_id" bson:"user_id,required"`
		BusinessID          string `json:"business_id" bson:"business_id"`
		AppID               string `json:"app_id" bson:"app_id,required"`
		TestKey             string `json:"test_key" bson:"test_key"`
		TestWebhook         string `json:"test_webhook" bson:"test_webhook"`
		LiveKey             string `json:"live_key" bson:"live_key"`
		LiveWebhook         string `json:"live_webhook" bson:"live_webhook"`
		IsActive            bool   `json:"is_active" bson:"is_active"`
	}
)

func (u ApiCredential) GetModelName() string {
	return "api_credential"
}
