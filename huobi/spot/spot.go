package spot

import (
	. "github.com/nntaoli-project/goex/v2"
)

var (
	currencyPaircacheMap = make(map[string]*CurrencyPair, 6)
)

type Spot struct {
	defaultUriOpts    *UriOptions
	tickerUnmarshaler TickerUnmarshaler
	depthUnmarshaler  DepthUnmarshaler
}

type spotImpl struct {
	*Spot
}

func New() *Spot {
	s := &Spot{}
	s.defaultUriOpts = &UriOptions{
		Endpoint:            "https://api.huobi.pro",
		TickerUri:           "/market/detail/merged",
		DepthUri:            "",
		KlineUri:            "",
		GetOrderUri:         "",
		GetPendingOrdersUri: "",
		CancelOrderUri:      "",
		NewOrderUri:         "",
	}
	return s
}

func (s *Spot) NewMarketApi(opts ...UriOption) IMarketRest {
	for _, opt := range opts {
		opt(s.defaultUriOpts)
	}

	unmarshaler := new(RespUnmarshaler)
	s.tickerUnmarshaler = unmarshaler
	s.depthUnmarshaler = unmarshaler

	imp := new(spotImpl)
	imp.Spot = s

	return imp
}
