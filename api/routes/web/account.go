package web

import (
	"github.com/gorilla/mux"
	"mux-crud/api/services"
	"mux-crud/happiness"
)

func AccountRoute(r *mux.Router) {
	happiness.Crud(r, services.UserRepository, "/account", actions,
		[]mux.MiddlewareFunc{privateMiddleware.AuthDeps},
		happiness.BindContext("user_context"),
		happiness.WithoutCreate, happiness.WithoutFetch,
		happiness.WithoutUpdate)
}
