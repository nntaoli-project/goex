package spot

import (
	. "github.com/nntaoli-project/goex/v2/model"
	. "github.com/nntaoli-project/goex/v2/options"
)

var (
	currencyPaircacheMap = make(map[string]*CurrencyPair, 6)
)

type BaseResponse struct {
	Status  string `json:"status"`
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

type Spot struct {
	uriOpts         UriOptions
	unmarshalerOpts UnmarshalerOptions
}

func New() *Spot {
	s := &Spot{
		uriOpts: UriOptions{
			Endpoint:            "https://api.huobi.pro",
			TickerUri:           "/market/detail/merged",
			DepthUri:            "",
			KlineUri:            "",
			GetOrderUri:         "",
			GetPendingOrdersUri: "",
			GetHistoryOrdersUri: "",
			CancelOrderUri:      "",
			NewOrderUri:         "",
		},
		unmarshalerOpts: UnmarshalerOptions{
			ResponseUnmarshaler: UnmarshalResponse,
			TickerUnmarshaler:   UnmarshalTicker,
			DepthUnmarshaler:    UnmarshalDepth,
		},
	}

	return s
}

func (s *Spot) WithUnmarshalerOptions(opts ...UnmarshalerOption) *Spot {
	for _, opt := range opts {
		opt(&s.unmarshalerOpts)
	}
	return s
}

func (s *Spot) WithUriOptions(opts ...UriOption) *Spot {
	for _, opt := range opts {
		opt(&s.uriOpts)
	}
	return s
}
