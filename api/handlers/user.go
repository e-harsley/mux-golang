package handlers

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mux-crud/api/schemas"
	"mux-crud/api/services"
	"mux-crud/happiness"
	"net/http"
)

type AuthHandler struct {
	service services.UserService
}

func (ah AuthHandler) LoginHandler(req schemas.LoginRequest, c happiness.C) *happiness.Response {
	user, err := ah.service.Login(req)

	if err != nil {
		return happiness.Err(err, http.StatusBadRequest)
	}

	return happiness.ResWithBinding(user, schemas.AuthResponse{})
}

func (ah AuthHandler) SignupHandler(req schemas.SignupRequest, c happiness.C) *happiness.Response {
	user, err := ah.service.Signup(req)

	if err != nil {
		return happiness.Err(err, http.StatusBadRequest)
	}

	return happiness.ResWithBinding(user, schemas.AuthResponse{})
}

func (ah AuthHandler) Me(req happiness.ExtraParameters, c happiness.C) *happiness.Response {
	userId, _ := primitive.ObjectIDFromHex(req.UserID)

	user, err := services.UserRepository.FindOne(bson.M{"_id": userId})

	if err != nil {
		return happiness.Err(err, http.StatusBadRequest)
	}
	return happiness.ResWithBinding(user, schemas.UserResponse{})
}
