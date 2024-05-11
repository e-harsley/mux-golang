package services

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mux-crud/api/models"
	"mux-crud/config"
	"mux-crud/happiness"
	"strings"
)

var (
	AccountRepository = happiness.NewBaseRepository(models.Account{})
)

type AccountService struct {
	cryptoConfig CryptoConfigService
}

func (as AccountService) GenerateKey(accountType happiness.AccountTypes) string {
	key := fmt.Sprintf("%s%s", accountType[:3], happiness.GenerateID(8))

	key = strings.ToUpper(key)

	count, _ := AccountRepository.CountDocuments(bson.M{"key": key})

	if count > 0 {
		key = as.GenerateKey(accountType)
	}
	return key
}

func (as AccountService) CryptoAccount(businessID, userID, name, email, phone string, accountType happiness.AccountTypes, currency models.SupportedCurrencyNetwork) (models.Account, error) {

	ccObj, err := as.cryptoConfig.CreateConfig(userID, businessID, name, currency)
	if err != nil {
		fmt.Println(err)
		return models.Account{}, err
	}

	key := as.GenerateKey(accountType)

	data := map[string]interface{}{
		"key":                           key,
		"currency_code":                 currency.Currency.Code,
		"currency":                      currency.Currency,
		"network":                       currency.Network,
		"network_code":                  currency.Network.Code,
		"type":                          accountType,
		"balance":                       0.00,
		"user_id":                       userID,
		"business_id":                   businessID,
		"name":                          name,
		"supported_currency_network_id": currency.ID,
		"accounting_currency":           happiness.AccountingCurrency,
		"crypto_config_id":              ccObj.ID,
	}

	obj, err := AccountRepository.BindDataOperationStruct(data)

	if err != nil {
		fmt.Println(err)

	}
	fmt.Println(" >>>> create here....")
	obj, err = AccountRepository.InsertOne(obj)
	fmt.Println(" after >>>> create here....")
	if err != nil {
		fmt.Println("err err", err)
	}
	fmt.Println("xpub", ccObj.Xpub)

	virtualAccount, resErr := Tatum.Wallet.VirtualAccount(map[string]interface{}{
		"xpub":          ccObj.Xpub,
		"currency":      strings.ToUpper(ccObj.Currency.Code),
		"accountCode":   obj.AccountingCurrency,
		"accountNumber": obj.ID,
	})
	fmt.Println("?????????????? >>>>>>>")
	if resErr.Status {
		fmt.Println(resErr.Data)
		return models.Account{}, fmt.Errorf("failed to create virtual account")
	}
	fmt.Println("?????????????? >>>>>>> check")

	subscrip, resErr := Tatum.Wallet.VirtualAccountSubscription(map[string]interface{}{
		"type": "ACCOUNT_INCOMING_BLOCKCHAIN_TRANSACTION",
		"attr": map[string]interface{}{
			"id":  virtualAccount.ID,
			"url": config.NewSettings().TatumWebhookURL,
		},
	})
	if resErr.Status {
		fmt.Println("resErr >>>", resErr)
		return models.Account{}, fmt.Errorf("failed to create subscription")
	}

	fmt.Println("subscrip >>> subscrip", subscrip)

	idHex, _ := primitive.ObjectIDFromHex(obj.ID)

	obj.TatumAccountID = virtualAccount.ID
	obj.TatumSubscriptionID = subscrip.ID

	obj, err = AccountRepository.FindOneAndUpdate(bson.D{{"_id", idHex}}, obj)
	if err != nil {
		return obj, err
	}

	return obj, nil
}

func (as AccountService) FiatAccount(businessID, userID, name, email, phone string, accountType happiness.AccountTypes, currency models.SupportedCurrencyNetwork) (models.Account, error) {

	key := as.GenerateKey(accountType)

	data := map[string]interface{}{
		"key":                        key,
		"currency":                   currency.Currency,
		"currency_code":              currency.Currency.Code,
		"supported_currency_network": currency.ID,
		"network":                    currency.Network,
		"network_code":               currency.Network.Code,
		"type":                       accountType,
		"balance":                    0.00,
		"user_id":                    userID,
		"business_id":                businessID,
		"name":                       name,
	}

	obj, err := AccountRepository.BindDataOperationStruct(data)

	if err != nil {
		return obj, err
	}

	obj, err = AccountRepository.InsertOne(obj)
	if err != nil {
		return obj, err
	}

	return obj, nil

}

func (as AccountService) CreateFiatAccount(businessID, userID, name, email, phone string, currency models.SupportedCurrencyNetwork) {

	accountTypes := happiness.AccountTypeList

	for _, accountType := range accountTypes {
		_, _ = as.FiatAccount(businessID, userID, name, email, phone, accountType, currency)
	}
}

func (as AccountService) CreateCryptoAccount(businessID, userID, name, email, phone string, currency models.SupportedCurrencyNetwork) {

	accountTypes := happiness.AccountTypeList

	for _, accountType := range accountTypes {
		_, _ = as.CryptoAccount(businessID, userID, name, email, phone, accountType, currency)
	}

}

func (as AccountService) CreateAccount(businessID, userID, name, email, phone string) (interface{}, error) {
	currencies, err := SupportedCurrencyNetworkRepo.Find(bson.M{"active": true})

	if err != nil {
		return nil, err
	}

	for _, currency := range currencies {
		if currency.CurrencyType == happiness.FIAT {
			as.CreateFiatAccount(businessID, userID, name, email, phone, currency)
		} else {
			as.CreateCryptoAccount(businessID, userID, name, email, phone, currency)
		}
	}

	return map[string]interface{}{"message": "account setup in progress"}, nil
}
