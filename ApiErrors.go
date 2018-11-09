package goex

type ApiError struct {
	ErrCode,
	ErrMsg,
	OriginErrMsg string
}

func (e ApiError) Error() string {
	return e.ErrMsg
}

func (e ApiError) OriginErr(err string) ApiError {
	e.ErrMsg = err
	return e
}

var (
	API_ERR                      = ApiError{ErrCode: "EX_ERR_0000", ErrMsg: "unknown error"}
	HTTP_ERR_CODE                = ApiError{ErrCode: "HTTP_ERR_0001", ErrMsg: "http request error"}
	EX_ERR_API_LIMIT             = ApiError{ErrCode: "EX_ERR_1000", ErrMsg: "api limited"}
	EX_ERR_SIGN                  = ApiError{ErrCode: "EX_ERR_0001", ErrMsg: "signature error"}
	EX_ERR_NOT_FIND_SECRETKEY    = ApiError{ErrCode: "EX_ERR_0002", ErrMsg: "not find secretkey"}
	EX_ERR_NOT_FIND_APIKEY       = ApiError{ErrCode: "EX_ERR_0003", ErrMsg: "not find apikey"}
	EX_ERR_INSUFFICIENT_BALANCE  = ApiError{ErrCode: "EX_ERR_0004", ErrMsg: "Insufficient Balance"}
	EX_ERR_PLACE_ORDER_FAIL      = ApiError{ErrCode: "EX_ERR_0005", ErrMsg: "place order failure"}
	EX_ERR_CANCEL_ORDER_FAIL     = ApiError{ErrCode: "EX_ERR_0006", ErrMsg: "cancel order failure"}
	EX_ERR_INVALID_CURRENCY_PAIR = ApiError{ErrCode: "EX_ERR_0007", ErrMsg: "invalid currency pair"}
	EX_ERR_NOT_FIND_ORDER        = ApiError{ErrCode: "EX_ERR_0008", ErrMsg: "not find order"}
	EX_ERR_SYMBOL_ERR            = ApiError{ErrCode: "EX_ERR_0009", ErrMsg: "symbol error"}
)
