package services

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"mux-crud/api/models"
	"mux-crud/happiness"
)

var BusinessRepository = happiness.NewBaseRepository(models.Business{})

type BusinessService struct {
	apiCred ApiCredentialService
}

func (bs BusinessService) Register(userID string, businessName string, businessType happiness.BusinessTypes, registrationType interface{}, isActive bool) (*models.Business, error) {
	filter := bson.D{
		{"$or", bson.A{
			bson.D{{"user_id", userID}},
		}},
	}

	count, err := BusinessRepository.CountDocuments(filter)

	if err != nil {
		return nil, err
	}

	if count > 0 {
		fmt.Println("user already has a business created")
		return nil, nil
	}

	data := map[string]interface{}{
		"user_id":  userID,
		"type":     businessType,
		"name":     businessName,
		"isActive": isActive,
	}
	if businessType != happiness.INDIVIDUAL {
		data["registration_type"] = registrationType
	}

	obj, err := BusinessRepository.BindDataOperationStruct(data)

	if err != nil {
		return nil, err
	}
	obj, err = BusinessRepository.InsertOne(obj)
	if err != nil {
		return nil, err
	}

	_, err = bs.apiCred.Register(obj.UserID, obj.ID, false)

	if err != nil {
		return nil, err
	}

	return &obj, err
}
