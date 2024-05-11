package main

import (
	"mux-crud/api/routes/web"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	webInit(r)
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func webInit(r *mux.Router) {
	web.AuthRoute(r)
	web.ProfileRoute(r)
	web.CurrencyRoute(r)
	web.AccountRoute(r)
}
