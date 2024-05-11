package tatum

import "mux-crud/config"

type Tatum struct {
	Wallet Wallet
}

func NewTatum(settings config.Settings) (tatum Tatum) {

	tatum.Wallet = Wallet{
		NewTatumMixin(settings),
	}

	return tatum
}
