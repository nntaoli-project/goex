package spot

import (
	. "github.com/nntaoli-project/goex/v2"
)

var (
	currencyPaircacheMap = make(map[string]*CurrencyPair, 6)
)

type Spot struct {
	unmarshalerOpts UnmarshalerOptions
}

type spotImpl struct {
	*Spot
	uriOpts UriOptions
}

func New(opts ...UnmarshalerOption) *Spot {
	unmarshaler := new(RespUnmarshaler)
	s := &Spot{
		unmarshalerOpts: UnmarshalerOptions{
			ResponseUnmarshaler: nil,
			TickerUnmarshaler:   unmarshaler.UnmarshalTicker,
			DepthUnmarshaler:    unmarshaler.UnmarshalDepth,
		},
	}
	for _, opt := range opts {
		opt(&s.unmarshalerOpts)
	}
	return s
}

func (s *Spot) NewMarketApi(opts ...UriOption) IMarketRest {
	imp := new(spotImpl)
	imp.Spot = s
	imp.uriOpts = UriOptions{
		Endpoint:  "https://api.huobi.pro",
		TickerUri: "/market/detail/merged",
	}
	for _, opt := range opts {
		opt(&imp.uriOpts)
	}
	return imp
}
