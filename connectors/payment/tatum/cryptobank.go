package tatum

import (
	"fmt"
	"mux-crud/connectors/payment"
	"mux-crud/happiness"
	"strings"
)

type (
	Wallet struct {
		TatumMixin
	}

	TatumWalletRequestResponse struct {
		Mnemonic string `json:"mnemonic"`
		Xpub     string `json:"xpub"`
	}

	TatumVirtualAccountRequest struct {
		Currency        string `json:"currency"`
		Xpub            string `json:"xpub"`
		Compliant       bool   `json:"compliant"`
		AccountCode     string `json:"accountCode"`
		AccountCurrency string `json:"accountCurrency"`
		AccountNumber   string `json:"accountNumber"`
	}

	TatumVirtualAccountBalanceResponse struct {
		AccountBalance   string `json:"accountBalance"`
		AvailableBalance string `json:"availableBalance"`
	}

	TatumVirtualAccountResponse struct {
		Currency           string                             `json:"currency"`
		Active             bool                               `json:"active"`
		Balance            TatumVirtualAccountBalanceResponse `json:"balance"`
		Frozen             bool                               `json:"frozen"`
		Xpub               string                             `json:"xpub"`
		AccountingCurrency string                             `json:"accountingCurrency"`
		ID                 string                             `json:"id"`
	}

	TatumSubscriptionAttr struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}

	TatumVirtualAccountSubscription struct {
		Type string                `json:"type"`
		Attr TatumSubscriptionAttr `json:"attr"`
	}
	TatumVirtualAccountSubscriptionResponse struct {
		ID string `json:"id"`
	}

	VirtualAccountAddressResponse struct {
		Xpub          string `json:"xpub"`
		DerivationKey string `json:"derivationKey"`
		Address       string `json:"address"`
		Currency      string `json:"currency"`
	}
)

func (cls Wallet) Wallet(crypto string) (TatumWalletRequestResponse, payment.ErrorResponse) {

	fmt.Println("wallet", crypto)
	var (
		stub     = strings.ToLower(crypto) + "/wallet"
		response TatumWalletRequestResponse
	)
	fmt.Println("wallet", crypto)

	res, err := cls.Retrieve(stub)
	fmt.Println("wallet", crypto)

	if err != nil {
		fmt.Println("error occured here 1", err)

		return response, cls.ErrorResponse(err)
	}

	err = res.JSONBody(&response)
	if err != nil {
		fmt.Println("error occured here 2", err)
		return response, cls.ErrorResponse(err)
	}
	return response, payment.ErrorResponse{Status: false}
}

func (cls Wallet) VirtualAccount(payload interface{}) (TatumVirtualAccountResponse, payment.ErrorResponse) {

	var (
		stub     = "ledger/account"
		request  TatumVirtualAccountRequest
		response TatumVirtualAccountResponse
	)

	fmt.Println("payload >>>> ", payload)
	err := happiness.BindDataOperationStruct(payload, &request)
	fmt.Println(request)
	if err != nil {
		fmt.Println("err >>>>", err)
		return response, cls.ErrorResponse(err)
	}
	res, err := cls.Send(stub, &request)

	if err != nil {
		fmt.Println("err err >>.", err)
		return response, cls.ErrorResponse(err)
	}
	fmt.Println("???? >>>>", res)
	err = res.JSONBody(&response)

	if err != nil {
		fmt.Println("?????????")
		return response, cls.ErrorResponse(err)
	}
	fmt.Println("virtual response ", &response)

	return response, payment.ErrorResponse{Status: false}

}

func (cls Wallet) VirtualAccountSubscription(payload interface{}) (TatumVirtualAccountSubscriptionResponse, payment.ErrorResponse) {
	fmt.Println("subscription subscription", payload)
	var (
		stub     = "subscription"
		request  TatumVirtualAccountSubscription
		response TatumVirtualAccountSubscriptionResponse
	)

	err := happiness.BindDataOperationStruct(payload, &request)

	if err != nil {
		fmt.Println(">>>>>>>> ", err)
		return response, cls.ErrorResponse(err)
	}

	res, err := cls.Send(stub, request)

	if err != nil {
		fmt.Println("subscription ", err)
		return response, cls.ErrorResponse(err)
	}

	err = res.JSONBody(&response)

	if err != nil {
		return response, cls.ErrorResponse(err)
	}

	return response, payment.ErrorResponse{Status: false}

}

func (cls Wallet) VirtualAccountAddress(accountID string) (VirtualAccountAddressResponse, payment.ErrorResponse) {

	var (
		stub     = fmt.Sprintf("offchain/account/%s", accountID)
		response VirtualAccountAddressResponse
	)

	res, err := cls.Retrieve(stub)

	if err != nil {
		return response, cls.ErrorResponse(err)
	}

	fmt.Println(res.StatusCode)
	err = res.JSONBody(&response)
	if err != nil {
		return response, cls.ErrorResponse(err)
	}
	return response, payment.ErrorResponse{Status: false}

}
