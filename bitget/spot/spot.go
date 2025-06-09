package spot

import (
	. "github.com/nntaoli-project/goex/v2/model"
	. "github.com/nntaoli-project/goex/v2/options"
)

type Spot struct {
	UnmarshalerOpts UnmarshalerOptions
	UriOpts         UriOptions
	currencyPairM   map[string]CurrencyPair
}

func New() *Spot {
	currencyPairCacheMap := make(map[string]CurrencyPair, 64)
	unmarshaler := &RespUnmarshaler{}
	return &Spot{
		UriOpts: UriOptions{
			Endpoint:           "https://api.bitget.com",
			GetExchangeInfoUri: "/api/v2/spot/public/symbols",
			TickerUri:          "/api/v2/spot/market/tickers",
			DepthUri:           "/api/v2/spot/market/orderbook",
			KlineUri:           "/api/v2/spot/market/candles",
		},
		UnmarshalerOpts: UnmarshalerOptions{
			GetExchangeInfoResponseUnmarshaler: unmarshaler.UnmarshalGetExchangeInfoResponse,
		},
		currencyPairM: currencyPairCacheMap,
	}
}
