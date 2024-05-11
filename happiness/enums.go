package happiness

const AccountingCurrency = "USD"

type EndpointActionTypes string

const (
	POST   EndpointActionTypes = "post"
	PUT    EndpointActionTypes = "put"
	DELETE EndpointActionTypes = "DELETE"
	GET    EndpointActionTypes = "get"
)

type OTPProviders string

const (
	Authenticator OTPProviders = "authenticator"
	Email         OTPProviders = "email"
	Phone         OTPProviders = "phone"
	Default       OTPProviders = "default"
)

type BusinessTypes string

const (
	INDIVIDUAL          BusinessTypes = "individual"
	REGISTERED_BUSINESS BusinessTypes = "registered_business"
)

type BusinessRegistrationTypes string

const (
	SOLE_PROPRIETORSHIP BusinessRegistrationTypes = "sole_proprietorship"
	LIMITED_COMPANY     BusinessRegistrationTypes = "limited_company"
)

type BusinessEstimatedTransactions string

const (
	LESS_THAN_1000    BusinessRegistrationTypes = "NGN 1 - 1,000"
	LESS_THAN_10000   BusinessRegistrationTypes = "NGN 1001 - 10,000"
	LESS_THAN_50000   BusinessRegistrationTypes = "NGN 10,001 - 50,000"
	LESS_THAN_100000  BusinessRegistrationTypes = "NGN 50,001 - 100,000"
	LESS_THAN_500000  BusinessRegistrationTypes = "NGN 100,001 - 500,000"
	LESS_THAN_1000000 BusinessRegistrationTypes = "NGN 500,001 - 1,000,000"
	ABOVE_1000000     BusinessRegistrationTypes = "NGN 1,000,001 - Above"
)

type CurrencyTypes string

const (
	FIAT           CurrencyTypes = "fiat"
	CRYPTOCURRENCY CurrencyTypes = "cryptocurrency"
)

type CryptoCurrencyNetworks string

const (
	BITCOIN CryptoCurrencyNetworks = "bitcoin"
)

var AccountTypeList = []AccountTypes{PURCHASE, EARNING, RESERVE, MAIN}

type AccountTypes string

const (
	PURCHASE AccountTypes = "purchase"
	EARNING  AccountTypes = "earning"
	RESERVE  AccountTypes = "reserve"
	MAIN     AccountTypes = "main"
)

type ApplicationError struct {
	Code    string
	Message string
}

const (
	TOKEN_REQUIRED            = "TOKEN.REQUIRED"
	TOKEN_INVALID             = "TOKEN.INVALID"
	TOKEN_DENIED              = "TOKEN.DENIED"
	TOKEN_EXPIRED             = "TOKEN.EXPIRED"
	HEADER_ATTRIBUTE_MISSING  = "HEADER.ATTRIBUTE.MISSING"
	HEADER_ATTRIBUTE_INVALID  = "HEADER.ATTRIBUTE.INVALID"
	ENDPOINT_NOT_IMPLEMENTED  = "ENDPOINT.NOT.IMPLEMENTED"
	ENDPOINT_OPERATION_DENIED = "ENDPOINT.OPERATION.RESTRICTED"
)

var ErrorMessages = map[string]string{
	TOKEN_REQUIRED:            "Access token is required. Modify your request and try again.",
	TOKEN_INVALID:             "Access token is invalid. Please provide a valid token and try again.",
	TOKEN_DENIED:              "Access token is denied. The user does not exist or is suspended.",
	TOKEN_EXPIRED:             "Your access token has expired. Try again with a valid token.",
	HEADER_ATTRIBUTE_MISSING:  "A required header parameter is missing.",
	HEADER_ATTRIBUTE_INVALID:  "The header parameter provided is invalid.",
	ENDPOINT_NOT_IMPLEMENTED:  "The requested endpoint is not yet implemented.",
	ENDPOINT_OPERATION_DENIED: "This activity cannot be carried out because the endpoint is restricted",
}

func NewApplicationError(code string) ApplicationError {
	return ApplicationError{
		Code:    code,
		Message: ErrorMessages[code],
	}
}

type AuthStatuses string

const (
	SUCCESS  AuthStatuses = "success"
	PENDING  AuthStatuses = "pending"
	NeedsOtp AuthStatuses = "need_otp"
	ERROR    AuthStatuses = "error"
	FAILED   AuthStatuses = "failed"
)

type PaymentStatuses string

const (
	PaymentSuccess PaymentStatuses = "success"
	PendingPayment PaymentStatuses = "pending"
)
