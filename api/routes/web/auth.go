package web

import (
	"github.com/gorilla/mux"
	"mux-crud/api/handlers"
	"mux-crud/api/services"
	"mux-crud/happiness"
)

var privateMiddleware = happiness.Middleware{PublicRoute: false}

var authHandler = handlers.AuthHandler{}

var authActions = happiness.CustomActions{

	happiness.Action{
		Name:    "/signup",
		Handler: happiness.Depend(authHandler.SignupHandler),
		Method:  happiness.POST,
	},
	happiness.Action{
		Name:    "/login",
		Handler: happiness.Depend(authHandler.LoginHandler),
		Method:  happiness.POST,
	},
}

func AuthRoute(r *mux.Router) {
	happiness.Crud(r, services.UserRepository, "/auth", authActions,
		happiness.WithoutCreate, happiness.WithoutList, happiness.WithoutFetch,
		happiness.WithoutUpdate)
}
