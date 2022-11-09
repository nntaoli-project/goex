package futures

import . "github.com/nntaoli-project/goex/v2"

type Futures struct {
}

type USDTFutures struct {
	uriOpts         UriOptions
	apiOpts         ApiOptions
	unmarshalerOpts UnmarshalerOptions
}

func New() *Futures {
	return &Futures{}
}

func NewUSDTFutures(uriOpts ...UriOption) *USDTFutures {
	f := &USDTFutures{
		uriOpts: UriOptions{
			Endpoint:            "https://api.hbdm.com",
			TickerUri:           "/linear-swap-ex/market/detail/merged",
			DepthUri:            "/linear-swap-ex/market/depth",
			KlineUri:            "/linear-swap-ex/market/history/kline",
			GetOrderUri:         "/linear-swap-api/v1/swap_cross_order_info",
			GetPendingOrdersUri: "/linear-swap-api/v1/swap_cross_openorders",
			GetHistoryOrdersUri: "/linear-swap-api/v3/swap_cross_hisorders",
			CancelOrderUri:      "/linear-swap-api/v1/swap_cross_cancel",
			NewOrderUri:         "/linear-swap-api/v1/swap_cross_order",
		},
		unmarshalerOpts: UnmarshalerOptions{
			ResponseUnmarshaler:                 UnmarshalResponse,
			KlineUnmarshaler:                    UnmarshalKline,
			TickerUnmarshaler:                   UnmarshalTicker,
			CancelOrderResponseUnmarshaler:      UnmarshalCancelOrderResponse,
			CreateOrderResponseUnmarshaler:      UnmarshalCreateOrderResponse,
			GetOrderInfoResponseUnmarshaler:     UnmarshalGetOrderInfoResponse,
			GetPendingOrdersResponseUnmarshaler: UnmarshalGetPendingOrdersResponse,
			GetHistoryOrdersResponseUnmarshaler: UnmarshalGetHistoryOrdersResponse,
		},
	}

	for _, opt := range uriOpts {
		opt(&f.uriOpts)
	}

	return f
}

func (f *USDTFutures) WithUnmarshalerOptions(opts ...UnmarshalerOption) {
	for _, opt := range opts {
		opt(&f.unmarshalerOpts)
	}
}

func (f *USDTFutures) NewCrossUdtFuturesTrade(key, secret string) ITradeRest {
	return &usdtFuturesTrade{
		USDTFutures: f,
		apiOpts: ApiOptions{
			Key:    key,
			Secret: secret,
		},
	}
}

func (f *USDTFutures) NewUsdtFuturesMarket() IMarketRest {
	return &usdtFuturesMarket{USDTFutures: f}
}
