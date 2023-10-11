package fapi

import (
	"github.com/nntaoli-project/goex/v2/model"
	"github.com/nntaoli-project/goex/v2/options"
)

type FApi struct {
	currencyPairM map[string]model.CurrencyPair

	UriOpts       options.UriOptions
	UnmarshalOpts options.UnmarshalerOptions
}

func NewFApi() *FApi {
	f := &FApi{
		UriOpts: options.UriOptions{
			Endpoint:            "https://fapi.binance.com",
			KlineUri:            "/fapi/v1/klines",
			TickerUri:           "/fapi/v1/ticker/24hr",
			DepthUri:            "/fapi/v1/depth",
			NewOrderUri:         "/fapi/v1/order",
			GetOrderUri:         "/fapi/v1/order",
			GetHistoryOrdersUri: "/fapi/v1/allOrders",
			GetPendingOrdersUri: "/fapi/v1/openOrders",
			CancelOrderUri:      "/fapi/v1/order",
			GetAccountUri:       "/fapi/v2/balance",
			GetPositionsUri:     "/fapi/v2/positionRisk",
			GetExchangeInfoUri:  "/fapi/v1/exchangeInfo",
		},
		UnmarshalOpts: options.UnmarshalerOptions{
			GetExchangeInfoResponseUnmarshaler:  UnmarshalGetExchangeInfoResponse,
			DepthUnmarshaler:                    UnmarshalDepthResponse,
			KlineUnmarshaler:                    UnmarshalKlinesResponse,
			GetAccountResponseUnmarshaler:       UnmarshalGetAccountResponse,
			CreateOrderResponseUnmarshaler:      UnmarshalCreateOrderResponse,
			CancelOrderResponseUnmarshaler:      UnmarshalCancelOrderResponse,
			GetOrderInfoResponseUnmarshaler:     UnmarshalGetOrderInfoResponse,
			GetPendingOrdersResponseUnmarshaler: UnmarshalGetPendingOrdersResponse,
			GetHistoryOrdersResponseUnmarshaler: UnmarshalGetHistoryOrdersResponse,
			GetPositionsResponseUnmarshaler:     UnmarshalGetPositionsResponse,
		},
	}

	return f
}

func (f *FApi) WithUriOption(opts ...options.UriOption) *FApi {
	for _, opt := range opts {
		opt(&f.UriOpts)
	}
	return f
}

func (f *FApi) WithUnmarshalOption(opts ...options.UnmarshalerOption) *FApi {
	for _, opt := range opts {
		opt(&f.UnmarshalOpts)
	}
	return f
}

func (f *FApi) NewPrvApi(opts ...options.ApiOption) *Prv {
	api := NewPrvApi(f, opts...)
	return api
}
