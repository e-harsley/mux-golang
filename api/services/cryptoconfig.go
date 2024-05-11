package services

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mux-crud/api/models"
	"mux-crud/happiness"
	"strings"
)

var CryptoConfigRepository = happiness.NewBaseRepository(models.CryptoConfig{})

type CryptoConfigService struct {
}

func (cc CryptoConfigService) generateCode(currencyCode string) string {
	key := fmt.Sprintf("%s%s", currencyCode, happiness.GenerateID(8))

	key = strings.ToUpper(key)

	count, _ := CryptoConfigRepository.CountDocuments(bson.M{"code": key})

	if count > 0 {
		key = cc.generateCode(currencyCode)
	}
	return key
}

func (cc CryptoConfigService) CreateConfig(userId, businessId, name string, currency models.SupportedCurrencyNetwork) (models.CryptoConfig, error) {
	code := cc.generateCode(currency.CurrencyCode)

	data := map[string]interface{}{
		"code":                          code,
		"user_id":                       userId,
		"business_id":                   businessId,
		"supported_currency_network_id": currency.ID,
		"supported_currency_network":    currency,
		"currency":                      currency.Currency,
		"currency_code":                 currency.Currency.Code,
		"network":                       currency.Network,
		"network_code":                  currency.NetworkCode,
	}

	obj, err := CryptoConfigRepository.BindDataOperationStruct(data)

	if err != nil {
		return obj, err
	}
	fmt.Println("about to create data", obj)
	ccObj, err := CryptoConfigRepository.InsertOne(obj)

	if err != nil {
		return obj, err
	}

	tatum, resErr := Tatum.Wallet.Wallet(currency.Network.Name)

	if resErr.Status {
		return obj, err
	}

	fmt.Println("tatum response", tatum)
	ccObj.Mnemonic = tatum.Mnemonic
	ccObj.Xpub = tatum.Xpub

	idHex, _ := primitive.ObjectIDFromHex(ccObj.ID)

	ccObj.Xpub = tatum.Xpub
	ccObj.Mnemonic = tatum.Mnemonic

	ccObj, err = CryptoConfigRepository.FindOneAndUpdate(bson.D{{"_id", idHex}}, ccObj)
	if err != nil {
		return obj, err
	}
	ccObj.Xpub = tatum.Xpub
	ccObj.Mnemonic = tatum.Xpub
	return ccObj, nil
}
