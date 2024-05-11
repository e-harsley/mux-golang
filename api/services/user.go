package services

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"mux-crud/api/models"
	"mux-crud/api/schemas"
	"mux-crud/happiness"
)

var UserRepository = happiness.NewBaseRepository(models.User{})

type UserService struct {
	business BusinessService
	account  AccountService
}

func (us UserService) Signup(payload schemas.SignupRequest) (interface{}, error) {
	payload.Name = payload.FirstName + " " + payload.LastName
	fmt.Println("payload", payload)
	filter := bson.D{
		{"$or", bson.A{
			bson.D{{"email", payload.Email}},
			bson.D{{"phone", payload.Phone}},
		}},
	}

	count, err := UserRepository.CountDocuments(filter)

	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, fmt.Errorf("user with email address or phone number already exists")
	}

	password := payload.Password

	fmt.Println(payload)

	obj, err := UserRepository.BindDataOperationStruct(payload)

	fmt.Println(obj)

	if err != nil {
		return nil, err
	}

	err = obj.SetPassword(password)
	if err != nil {
		return nil, err
	}
	obj.Generate2faSecret()
	if err != nil {
		return nil, err
	}
	obj.OtpProvider = happiness.Email

	obj, err = UserRepository.InsertOne(obj)

	business, err := us.business.Register(obj.ID, payload.BusinessName, payload.BusinessType, payload.RegistrationType, false)
	if err != nil {
		return obj, err
	}

	_, _ = us.account.CreateAccount(business.ID, business.UserID, business.Name, obj.Email, obj.Phone)

	return obj, err

}

func (us UserService) Login(payload schemas.LoginRequest) (interface{}, error) {
	var (
		status = "success"
	)
	filter := bson.D{
		{"$or", bson.A{
			bson.D{{"email", payload.Username}},
			bson.D{{"phone", payload.Username}},
		}},
	}

	obj, err := UserRepository.FindOne(filter)

	if err != nil {
		return nil, err
	}

	if obj == nil {
		return nil, fmt.Errorf(fmt.Sprintf("user with username %s does not exist", payload.Username))
	}

	check := obj.CheckPassword(payload.Password)

	fmt.Println("obj", obj)

	if !check {
		return nil, fmt.Errorf(fmt.Sprintf("invalid password for user with username %s", payload.Username))
	}

	if obj.IsSuspended != nil && *obj.IsSuspended {
		return nil, fmt.Errorf("user with this username (%s) has been suspended", payload.Username)
	}

	if obj.IsSuspended != nil && *obj.Is2faEnabled && (payload.Otp == "" || obj.ValidateOTP(payload.Otp) == false) {
		status = "needs_otp"

		if payload.Otp == "" && obj.OtpProvider != happiness.Authenticator {
			otp := obj.GenerateOTP()
			fmt.Println("OTP ---> ", otp)
			//todo:: implement sending sms here
		}
		data, err := happiness.MapDump(obj)
		if err != nil {
			return nil, err
		}
		data["status"] = status

		return data, nil
	}

	authToken := happiness.NewAuthToken(obj.ID)

	token, err := authToken.Token()

	if err != nil {
		return nil, err
	}
	data, err := happiness.MapDump(obj)

	data["token"] = token
	data["status"] = status

	return data, nil
}
