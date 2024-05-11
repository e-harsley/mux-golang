package services

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"mux-crud/api/models"
	"mux-crud/happiness"
)

var ApiCredentialRepository = happiness.NewBaseRepository(models.ApiCredential{})

type ApiCredentialService struct {
}

func (ac ApiCredentialService) Register(userId string, businessId string, isActive bool) (interface{}, error) {

	filter := bson.D{
		{"$or", bson.A{
			bson.D{{"user_id", userId}},
			bson.D{{"business_id", businessId}},
		}},
	}

	count, err := ApiCredentialRepository.CountDocuments(filter)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	if count > 0 {
		fmt.Println("auth cred already exist")
		return nil, nil
	}

	data := map[string]interface{}{
		"user_id":     userId,
		"business_id": businessId,
		"is_active":   isActive,
		"app_id":      happiness.GenerateID(10),
	}

	if !isActive {
		data["test_key"] = happiness.GenerateSecretKey(true)
	}

	fmt.Println("data", data)

	obj, err := ApiCredentialRepository.BindDataOperationStruct(data)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return ApiCredentialRepository.InsertOne(obj)
}
