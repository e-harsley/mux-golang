package models

import "mux-crud/happiness"

type (
	Topup struct {
		happiness.BaseModel        `bson:",inline"`
		Code                       string                    `json:"code" bson:"code"`
		BankCode                   string                    `json:"bank_code" bson:"bank_code"`
		Source                     string                    `json:"source" bson:"source"`
		Amount                     float64                   `json:"amount" bson:"amount"`
		ToppedAmount               float64                   `json:"topped_amount" bson:"topped_amount"`
		Payable                    float64                   `json:"payable" bson:"payable"`
		AccountID                  string                    `json:"account_id" bson:"account_id"`
		Account                    Account                   `json:"account" bson:"account"`
		Currency                   Currency                  `json:"currency" bson:"currency,omitempty"`
		Network                    SupportedNetwork          `json:"network" bson:"network,omitempty"`
		NetworkCode                string                    `json:"network_code" bson:"network_code,omitempty"`
		SupportedCurrencyNetworkID string                    `json:"supported_currency_network_id,omitempty" bson:"supported_currency_network_id,omitempty"`
		Name                       string                    `json:"name" bson:"name"`
		Email                      string                    `json:"email" bson:"email,required"`
		Phone                      string                    `json:"phone" bson:"phone,required"`
		Narration                  string                    `json:"narration" bson:"narration,required"`
		Status                     happiness.PaymentStatuses `json:"status" bson:"status"`
	}
)

func (cls Topup) GetModelName() string {
	return "account"
}
