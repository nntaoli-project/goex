package common

import (
	"encoding/json"
	"errors"
	. "github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/internal/logger"
	"net/http"
	"net/url"
)

type V5 struct {
	uriOpts       UriOptions
	unmarshalOpts UnmarshalerOptions
	marketApi     IMarketRest
}

type BaseResp struct {
	Code int             `json:"code,string"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func New() *V5 {
	unmarshaler := new(RespUnmarshaler)

	f := &V5{
		uriOpts: UriOptions{
			Endpoint:            "https://www.okx.com",
			KlineUri:            "/api/v5/market/candles",
			TickerUri:           "/api/v5/market/ticker",
			DepthUri:            "/api/v5/market/books",
			NewOrderUri:         "/api/v5/trade/order",
			GetOrderUri:         "/api/v5/trade/order",
			GetPendingOrdersUri: "/api/v5/trade/orders-pending",
			CancelOrderUri:      "/api/v5/trade/cancel-order",
			GetAccountUri:       "/api/v5/account/balance",
			GetPositionsUri:     "/api/v5/account/positions",
		},
		unmarshalOpts: UnmarshalerOptions{
			ResponseUnmarshaler:                  unmarshaler.UnmarshalResponse,
			KlineUnmarshaler:                     unmarshaler.UnmarshalGetKlineResponse,
			TickerUnmarshaler:                    unmarshaler.UnmarshalTicker,
			DepthUnmarshaler:                     unmarshaler.UnmarshalDepth,
			CreateOrderResponseUnmarshaler:       unmarshaler.UnmarshalCreateOrderResponse,
			GetPendingOrdersResponseUnmarshaler:  unmarshaler.UnmarshalGetPendingOrdersResponse,
			CancelOrderResponseUnmarshaler:       unmarshaler.UnmarshalCancelOrderResponse,
			GetOrderInfoResponseUnmarshaler:      unmarshaler.UnmarshalGetOrderInfoResponse,
			GetAccountResponseUnmarshaler:        unmarshaler.UnmarshalGetAccountResponse,
			GetPositionsResponseUnmarshaler:      unmarshaler.UnmarshalGetPositionsResponse,
			GetFuturesAccountResponseUnmarshaler: unmarshaler.UnmarshalGetFuturesAccountResponse,
		},
	}
	f.marketApi = &Market{f}
	return f
}

func (f *V5) WithUriOption(opts ...UriOption) *V5 {
	for _, opt := range opts {
		opt(&f.uriOpts)
	}
	return f
}

func (f *V5) WithUnmarshalOption(opts ...UnmarshalerOption) *V5 {
	for _, opt := range opts {
		opt(&f.unmarshalOpts)
	}
	return f
}

func (f *V5) DoAuthRequest(httpMethod, reqUrl string, params *url.Values, apiOpts ApiOptions, headers map[string]string) ([]byte, error) {
	var (
		reqBodyStr string
		reqUri     string
	)

	if http.MethodGet == httpMethod {
		reqUrl += "?" + params.Encode()
	}

	if http.MethodPost == httpMethod {
		reqBody, _ := ValuesToJson(*params)
		reqBodyStr = string(reqBody)
	}

	_url, _ := url.Parse(reqUrl)
	reqUri = _url.RequestURI()
	signStr, timestamp := SignParams(httpMethod, reqUri, apiOpts.Secret, reqBodyStr)
	logger.Debugf("[DoAuthRequest] sign base64: %s, timestamp: %s", signStr, timestamp)

	headers = map[string]string{
		"Content-Type": "application/json; charset=UTF-8",
		//"Accept":               "application/json",
		"OK-ACCESS-KEY":        apiOpts.Key,
		"OK-ACCESS-PASSPHRASE": apiOpts.Passphrase,
		"OK-ACCESS-SIGN":       signStr,
		"OK-ACCESS-TIMESTAMP":  timestamp}

	respBody, err := GetHttpCli().DoRequest(httpMethod, reqUrl, reqBodyStr, headers)
	if err != nil {
		return respBody, err
	}
	logger.Debugf("[DoAuthRequest] response body: %s", string(respBody))

	var baseResp BaseResp
	err = f.unmarshalOpts.ResponseUnmarshaler(respBody, &baseResp)
	if err != nil {
		return respBody, err
	}

	if baseResp.Code != 0 {
		return baseResp.Data, errors.New(baseResp.Msg)
	}

	return baseResp.Data, nil
}

func (f *V5) MarketApi() IMarketRest {
	return f.marketApi
}

func (f *V5) NewTradeApi(opts ...ApiOption) ITradeRest {
	api := NewTrade(opts...)
	api.V5 = f
	return api
}

func (f *V5) NewAccountApi(opts ...ApiOption) IAccount {
	var apiOpts ApiOptions
	for _, opt := range opts {
		opt(&apiOpts)
	}

	api := NewAccountApi(apiOpts)
	api.V5 = f

	return api
}
