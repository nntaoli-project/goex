package common

import (
	"encoding/json"
	. "github.com/nntaoli-project/goex/v2"
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
		},
		unmarshalOpts: UnmarshalerOptions{
			ResponseUnmarshaler:                 unmarshaler.UnmarshalResponse,
			KlineUnmarshaler:                    unmarshaler.UnmarshalGetKlineResponse,
			TickerUnmarshaler:                   unmarshaler.UnmarshalTicker,
			DepthUnmarshaler:                    unmarshaler.UnmarshalDepth,
			CreateOrderResponseUnmarshaler:      unmarshaler.UnmarshalCreateOrderResponse,
			GetPendingOrdersResponseUnmarshaler: unmarshaler.UnmarshalGetPendingOrdersResponse,
			CancelOrderResponseUnmarshaler:      unmarshaler.UnmarshalCancelOrderResponse,
			GetOrderInfoResponseUnmarshaler:     unmarshaler.UnmarshalGetOrderInfoResponse,
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

func (f *V5) MarketApi() IMarketRest {
	return f.marketApi
}

func (f *V5) NewTradeApi(opts ...ApiOption) ITradeRest {
	api := NewTrade(opts...)
	api.V5 = f
	return api
}
