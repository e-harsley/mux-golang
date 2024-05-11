package models

import "mux-crud/happiness"

type (
	Checkout struct {
		happiness.BaseModel        `bson:",inline"`
		CurrencyCode               string                 `json:"currency_code" bson:"currency_code,omitempty"`
		ReferenceCode              string                 `json:"reference_code" bson:"reference_code,omitempty"`
		Currency                   Currency               `json:"currency" bson:"currency,omitempty"`
		Network                    SupportedNetwork       `json:"network" bson:"network,omitempty"`
		NetworkCode                string                 `json:"network_code" bson:"network_code,omitempty"`
		AccountID                  string                 `json:"account_id" bson:"account_id"`
		SupportedCurrencyNetworkID string                 `json:"supported_currency_network_id,omitempty" bson:"supported_currency_network_id,omitempty"`
		BusinessID                 string                 `json:"business_id" bson:"business_id"`
		UserID                     string                 `json:"user_id" bson:"user_id"`
		Amount                     float64                `json:"amount" bson:"amount"`
		TransactionID              string                 `json:"transaction_id" bson:"transaction_id"`
		Transaction                map[string]interface{} `json:"transaction" bson:"transaction"`
		ReferenceID                string                 `json:"reference_id" bson:"reference_id"`
		ProductType                string                 `json:"product_type" bson:"product_type"`
		ExchangeData               map[string]interface{} `json:"exchange_data" bson:"exchange_data"`
		Source                     string                 `json:"source" bson:"source"`
		Payer                      map[string]interface{} `json:"payer" bson:"payer"`
	}
)
