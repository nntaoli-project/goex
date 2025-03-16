package spot

import (
	. "github.com/nntaoli-project/goex/v2/model"
	. "github.com/nntaoli-project/goex/v2/options"
)

var (
	currencyPaircacheMap = make(map[string]*CurrencyPair, 6)
)

type Spot struct {
	UnmarshalerOpts UnmarshalerOptions
	UriOpts         UriOptions
}

func New() *Spot {
	unmarshaler := new(RespUnmarshaler)
	s := &Spot{
		UriOpts: UriOptions{
			Endpoint:            "https://api.binance.com",
			TickerUri:           "/api/v3/ticker/24hr",
			DepthUri:            "/api/v3/depth",
			KlineUri:            "/api/v3/klines",
			NewOrderUri:         "/api/v3/order",
			GetPendingOrdersUri: "/api/v3/openOrders",
			CancelOrderUri:      "/api/v3/order",
			GetOrderUri:         "/api/v3/order",
			GetHistoryOrdersUri: "/api/v3/allOrders",
		},
		UnmarshalerOpts: UnmarshalerOptions{
			ResponseUnmarshaler:                 unmarshaler.UnmarshalResponse,
			TickerUnmarshaler:                   unmarshaler.UnmarshalGetTickerResponse,
			DepthUnmarshaler:                    unmarshaler.UnmarshalGetDepthResponse,
			KlineUnmarshaler:                    unmarshaler.UnmarshalGetKlineResponse,
			CreateOrderResponseUnmarshaler:      unmarshaler.UnmarshalCreateOrderResponse,
			GetPendingOrdersResponseUnmarshaler: unmarshaler.UnmarshalGetPendingOrdersResponse,
			CancelOrderResponseUnmarshaler:      unmarshaler.UnmarshalCancelOrderResponse,
		},
	}
	return s
}

func (s *Spot) WithUriOption(uriOpts ...UriOption) {
	for _, opt := range uriOpts {
		opt(&s.UriOpts)
	}
}

func (s *Spot) WithUnmarshalerOptions(opts ...UnmarshalerOption) *Spot {
	for _, opt := range opts {
		opt(&s.UnmarshalerOpts)
	}
	return s
}

func (s *Spot) NewPrvApi(apiOpts ...ApiOption) *PrvApi {
	prv := NewPrvApi(apiOpts...)
	prv.Spot = s
	return prv
}
