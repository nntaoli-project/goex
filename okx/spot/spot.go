package spot

import (
	. "github.com/nntaoli-project/goex/v2"
)

var (
	currencyPaircacheMap = make(map[string]*CurrencyPair, 6)
)

type Spot struct {
	uriOpts         UriOptions
	unmarshalerOpts UnmarshalerOptions
	marketApi       IMarketRest
}

type spotImpl struct {
	*Spot
}

func New(opts ...UriOption) *Spot {
	unmarshal := new(RespUnmarshaler)
	s := &Spot{
		uriOpts: UriOptions{
			Endpoint:            "https://www.okx.com",
			TickerUri:           "/api/v5/market/ticker",
			DepthUri:            "/api/v5/market/books",
			KlineUri:            "",
			GetOrderUri:         "",
			GetPendingOrdersUri: "",
			GetHistoryOrdersUri: "",
			CancelOrderUri:      "",
			NewOrderUri:         "",
		},
		unmarshalerOpts: UnmarshalerOptions{
			ResponseUnmarshaler: unmarshal.UnmarshalResponse,
			TickerUnmarshaler:   unmarshal.UnmarshalTicker,
			DepthUnmarshaler:    unmarshal.UnmarshalDepth,
		},
	}

	for _, opt := range opts {
		opt(&s.uriOpts)
	}

	s.marketApi = &spotImpl{Spot: s}

	return s
}

func (s *Spot) WithUnmarshalerOptions(opts ...UnmarshalerOption) *Spot {
	for _, opt := range opts {
		opt(&s.unmarshalerOpts)
	}
	return s
}

func (s *Spot) MarketApi() IMarketRest {
	return s.marketApi
}
