package payment

type (
	ErrorDataResponse struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}

	ErrorResponse struct {
		Status bool              `json:"status"`
		Data   ErrorDataResponse `json:"data"`
	}
)
