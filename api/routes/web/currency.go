package web

import (
	"github.com/gorilla/mux"
	"mux-crud/api/services"
	"mux-crud/happiness"
)

func CurrencyRoute(r *mux.Router) {
	happiness.Crud(r, services.CurrencyRepository, "/currency",
		happiness.WithoutCreate, happiness.WithoutUpdate)
	happiness.Crud(r, services.SupportedNetworkRepository, "/supported-network",
		happiness.WithoutCreate, happiness.WithoutUpdate)
	happiness.Crud(r, services.SupportedCurrencyNetworkRepo, "/supported-currency-network",
		happiness.WithoutCreate, happiness.WithoutUpdate)
}
