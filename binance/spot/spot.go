package spot

import (
	. "github.com/nntaoli-project/goex/v2"
)

var (
	currencyPaircacheMap = make(map[string]*CurrencyPair, 6)
)

type Spot struct {
	unmarshalerOpts UnmarshalerOptions
	uriOpts         UriOptions

	marketApi IMarketRest
}

type spotImpl struct {
	*Spot
}

func New(opts ...UriOption) *Spot {
	unmarshaler := new(RespUnmarshaler)
	s := &Spot{
		uriOpts: UriOptions{
			Endpoint:  "https://api.binance.com",
			TickerUri: "/api/v3/ticker/24hr",
			DepthUri:  "/api/v3/depth",
			KlineUri:  "/api/v3/klines",
		},
		unmarshalerOpts: UnmarshalerOptions{
			ResponseUnmarshaler: unmarshaler.UnmarshalResponse,
			TickerUnmarshaler:   unmarshaler.UnmarshalGetTickerResponse,
			DepthUnmarshaler:    unmarshaler.UnmarshalGetDepthResponse,
			KlineUnmarshaler:    unmarshaler.UnmarshalGetKlineResponse,
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
