package models

import "mux-crud/happiness"

type (
	Account struct {
		happiness.BaseModel        `bson:",inline"`
		Key                        string                 `json:"key" bson:"key,omitempty"`
		CurrencyCode               string                 `json:"currency_code" bson:"currency_code,omitempty"`
		Currency                   Currency               `json:"currency" bson:"currency,omitempty"`
		Network                    SupportedNetwork       `json:"network" bson:"network,omitempty"`
		NetworkCode                string                 `json:"network_code" bson:"network_code,omitempty"`
		Type                       happiness.AccountTypes `json:"type" bson:"type,omitempty"`
		AccountingCurrency         string                 `json:"accounting_currency,omitempty" bson:"accounting_currency,omitempty"`
		Balance                    float64                `json:"balance" bson:"balance,omitempty"`
		UserID                     string                 `json:"user_id" bson:"user_id,omitempty"`
		BusinessID                 string                 `json:"business_id" bson:"business_id,omitempty"`
		Locked                     bool                   `json:"locked" bson:"locked,omitempty"`
		LockKey                    string                 `json:"lock_key" bson:"lock_key,omitempty"`
		Name                       string                 `json:"name" bson:"name,omitempty"`
		SupportedCurrencyNetworkID string                 `json:"supported_currency_network_id,omitempty" bson:"supported_currency_network_id,omitempty"`
		TatumSubscriptionID        string                 `json:"tatum_subscription_id" bson:"tatum_subscription_id"`
		TatumAccountID             string                 `json:"tatum_account_id" bson:"tatum_account_id,omitempty"`
		CryptoConfigID             string                 `json:"crypto_config_id" bson:"crypto_config_id,omitempty"`
	}

	CryptoConfig struct {
		happiness.BaseModel        `bson:",inline"`
		Code                       string                   `json:"code" bson:"code"`
		UserID                     string                   `json:"user_id" bson:"user_id"`
		BusinessID                 string                   `json:"business_id" bson:"business_id"`
		SupportedCurrencyNetworkID string                   `json:"supported_currency_network_id,omitempty" bson:"supported_currency_network_id,omitempty"`
		SupportedCurrencyNetwork   SupportedCurrencyNetwork `json:"supported_currency_network,omitempty" bson:"supported_currency_network,omitempty"`
		NetworkCode                string                   `json:"network_code,omitempty" bson:"network_code,omitempty"`
		Network                    SupportedNetwork         `json:"network" bson:"network"`
		CurrencyCode               string                   `json:"currency_code" bson:"currency_code"`
		Currency                   Currency                 `json:"currency" bson:"currency"`
		Mnemonic                   string                   `json:"mnemonic" bson:"mnemonic"`
		Xpub                       string                   `json:"xpub" bson:"xpub"`
	}

	Wallet struct {
		happiness.BaseModel `bson:",inline"`
		AccountID           string `json:"account_id"`
		Address             string `json:"address"`
		SubscriptionID      string `json:"subscription_id"`
	}
)

func (cls Account) GetModelName() string {
	return "account"
}

func (cls CryptoConfig) GetModelName() string {
	return "crypto_config"
}

func (cls Wallet) GetModelName() string {
	return "wallet"
}
