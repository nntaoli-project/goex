package futures

import (
	. "github.com/nntaoli-project/goex/v2/options"
)

type Futures struct {
	USDTSwapFutures *USDTSwap
}

type USDTSwap struct {
	uriOpts         UriOptions
	unmarshalerOpts UnmarshalerOptions
}

func New() *Futures {
	return &Futures{
		USDTSwapFutures: NewUSDTSwap(),
	}
}

func NewUSDTSwap() *USDTSwap {
	f := &USDTSwap{
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
	return f
}

func (f *USDTSwap) WithUnmarshalerOptions(opts ...UnmarshalerOption) *USDTSwap {
	for _, opt := range opts {
		opt(&f.unmarshalerOpts)
	}
	return f
}

func (f *USDTSwap) WithUriOptions(uriOpts ...UriOption) *USDTSwap {
	for _, opt := range uriOpts {
		opt(&f.uriOpts)
	}
	return f
}

func (f *USDTSwap) NewUSDTSwapPrvApi(apiOpts ...ApiOption) *USDTSwapPrvApi {
	prv := NewUSDTSwapPrvApi(apiOpts...)
	prv.USDTSwap = f
	return prv
}
