package tatum

import (
	"fmt"
	"mux-crud/config"
	"mux-crud/connectors/payment"
	"mux-crud/happiness"
)

type TatumMixin struct {
	endpoint string
	apiKey   string
}

func NewTatumMixin(settings config.Settings) (tatum TatumMixin) {
	tatum.endpoint = settings.TatumBaseUrl
	tatum.apiKey = settings.TatumApiKey

	return tatum
}

func (cls TatumMixin) BuildUrl(stub string) string {
	return cls.endpoint + stub
}

func (cls TatumMixin) Header() map[string]string {
	return map[string]string{
		"content-type": "application/json",
		"x-api-key":    cls.apiKey,
	}
}

func (cls TatumMixin) Send(url string, data interface{}) (*happiness.HTTPResponse, error) {
	postUrl := cls.BuildUrl(url)
	fmt.Println(postUrl, cls.Header(), data)

	return happiness.MakePostRequest(postUrl, &data, cls.Header())
}

func (cls TatumMixin) Retrieve(url string, params ...map[string]string) (*happiness.HTTPResponse, error) {
	postUrl := cls.BuildUrl(url)
	opt := happiness.RequestOptions{
		Header: cls.Header(),
	}

	if len(params) > 0 {
		opt.Params = params[0]
	}
	return happiness.MakeGetRequest(postUrl, &opt)

}

func (cls TatumMixin) ErrorResponse(err error) payment.ErrorResponse {
	errorMes := payment.ErrorDataResponse{
		Status:  "failed",
		Message: err.Error(),
	}
	return payment.ErrorResponse{
		Status: true,
		Data:   errorMes,
	}
}
