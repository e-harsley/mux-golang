package services

import (
	"mux-crud/config"
	"mux-crud/connectors/payment/tatum"
)

var Tatum = tatum.NewTatum(*config.NewSettings())
