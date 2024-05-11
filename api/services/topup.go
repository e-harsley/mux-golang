package services

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mux-crud/api/models"
	"mux-crud/happiness"
	"strings"
)

var TopupRepository = happiness.NewBaseRepository(models.Topup{})

type TopupService struct {
}

func (as TopupService) GenerateKey() string {
	key := fmt.Sprintf("%s", happiness.GenerateID(10))

	key = strings.ToUpper(key)

	count, _ := TopupRepository.CountDocuments(bson.M{"key": key})

	if count > 0 {
		key = as.GenerateKey()
	}
	return key
}

func (ts TopupService) Initiate(accountID string, amount string, paymentSource string, currency string, network string) (models.Topup, error) {

	data := map[string]interface{}{}

	data["code"] = ts.GenerateKey()
	data["narration"] = fmt.Sprintf("Account topup with reference %s", data["code"])
	accountIdHex, _ := primitive.ObjectIDFromHex(accountID)
	account, err := AccountRepository.FindOne(bson.M{"_id": accountIdHex})

	if err != nil {
		return models.Topup{}, err
	}

	data["action"] = "topup"

	data["reference_code"] = data["code"]
	data["amount"] = amount
	data["payment_source_code"] = paymentSource
	businessIdHex, _ := primitive.ObjectIDFromHex(account.ID)

	business, err := BusinessRepository.FindOne(bson.M{"_id": businessIdHex})
	data["currency"] = account.Currency
	data["currency_code"] = account.CurrencyCode
	data["network"] = account.Network
	data["network_code"] = account.NetworkCode
	data["name"] = business.Name
	data["email"] = business.User.Email
	data["phone"] = business.User.Phone
	data["account_id"] = account.ID
	data["account"] = account
	data["user_id"] = account.UserID
	data["business_id"] = account.BusinessID

	topUp, err := TopupRepository.BindDataOperationStruct(data)
	if err != nil {
		return models.Topup{}, err
	}

	return TopupRepository.InsertOne(topUp)
}
