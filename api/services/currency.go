package services

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mux-crud/api/models"
	"mux-crud/api/schemas"
	"mux-crud/happiness"
)

var CurrencyRepository = happiness.NewBaseRepository(models.Currency{})
var SupportedNetworkRepository = happiness.NewBaseRepository(models.SupportedNetwork{})
var SupportedCurrencyNetworkRepo = happiness.NewBaseRepository(models.SupportedCurrencyNetwork{})

type SupportedCurrencyNetworkService struct{}

func (cls SupportedCurrencyNetworkService) Create(schema schemas.SupportedCurrencyNetworkSchema) (interface{}, error) {
	fmt.Println("honest feedback", schema)
	currencyID := schema.CurrencyID

	idHex, _ := primitive.ObjectIDFromHex(currencyID)
	fmt.Println("honnest >>>", idHex)
	currencyObj, err := CurrencyRepository.FindOne(bson.M{"_id": idHex})

	if err != nil {
		return nil, err
	}

	fmt.Println(currencyObj)
	currency := *currencyObj

	payload := map[string]interface{}{
		"currency_id":   currency.ID,
		"currency_code": currency.Code,
		"currency_type": currency.CurrencyType,
		"currency":      currency,
		"active":        schema.Active,
	}

	if currency.CurrencyType == happiness.CRYPTOCURRENCY && schema.NetworkID == "" {
		return nil, fmt.Errorf("network id is required")
	}

	if currency.CurrencyType == happiness.CRYPTOCURRENCY {
		idHex, _ := primitive.ObjectIDFromHex(schema.NetworkID)
		networkPtr, err := SupportedNetworkRepository.FindOne(bson.M{"_id": idHex})
		if err != nil {
			return nil, err
		}
		network := *networkPtr

		payload["network"] = network
		payload["network_code"] = network.Code
		payload["network_name"] = network.Name
		payload["network_id"] = network.ID
	}

	fmt.Println("payload", payload)
	obj, err := SupportedCurrencyNetworkRepo.BindDataOperationStruct(&payload)
	fmt.Println("payload", payload)
	if err != nil {
		fmt.Println(">>>e err")
		return nil, err
	}
	obj, err = SupportedCurrencyNetworkRepo.InsertOne(obj)

	if err != nil {
		return nil, err
	}

	return obj, nil
}
