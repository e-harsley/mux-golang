package web

import (
	"github.com/gorilla/mux"
	"mux-crud/api/services"
	"mux-crud/happiness"
)

var actions = happiness.CustomActions{
	happiness.Action{
		Name:        "/me",
		Handler:     happiness.Depend(authHandler.Me, happiness.BindContext("user_context")),
		Method:      happiness.GET,
		Middlewares: []mux.MiddlewareFunc{privateMiddleware.AuthDeps},
	},
}

func ProfileRoute(r *mux.Router) {
	happiness.Crud(r, services.UserRepository, "/user", actions,
		happiness.WithoutCreate, happiness.WithoutList, happiness.WithoutFetch,
		happiness.WithoutUpdate)
}
