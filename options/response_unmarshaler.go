package options

import "github.com/nntaoli-project/goex/v2/model"

type ResponseUnmarshaler func([]byte, interface{}) error
type GetTickerResponseUnmarshaler func([]byte) (*model.Ticker, error)
type GetDepthResponseUnmarshaler func([]byte) (*model.Depth, error)
type GetKlineResponseUnmarshaler func([]byte) ([]model.Kline, error)
type CreateOrderResponseUnmarshaler func([]byte) (*model.Order, error)
type GetOrderInfoResponseUnmarshaler func([]byte) (*model.Order, error)
type GetPendingOrdersResponseUnmarshaler func([]byte) ([]model.Order, error)
type CancelOrderResponseUnmarshaler func([]byte) error
type GetHistoryOrdersResponseUnmarshaler func([]byte) ([]model.Order, error)
type GetAccountResponseUnmarshaler func([]byte) (map[string]model.Account, error)
type GetPositionsResponseUnmarshaler func([]byte) ([]model.FuturesPosition, error)
type GetFuturesAccountResponseUnmarshaler func([]byte) (map[string]model.FuturesAccount, error)
type GetExchangeInfoResponseUnmarshaler func([]byte) (map[string]model.CurrencyPair, error)
type GetFundingRateResponseUnmarshaler func([]byte) (*model.FundingRate, error)
type GetFundingRateHistoryResponseUnmarshaler func([]byte) ([]model.FundingRate, error)

type UnmarshalerOptions struct {
	ResponseUnmarshaler                      ResponseUnmarshaler
	TickerUnmarshaler                        GetTickerResponseUnmarshaler
	DepthUnmarshaler                         GetDepthResponseUnmarshaler
	KlineUnmarshaler                         GetKlineResponseUnmarshaler
	CreateOrderResponseUnmarshaler           CreateOrderResponseUnmarshaler
	GetOrderInfoResponseUnmarshaler          GetOrderInfoResponseUnmarshaler
	GetPendingOrdersResponseUnmarshaler      GetPendingOrdersResponseUnmarshaler
	GetHistoryOrdersResponseUnmarshaler      GetHistoryOrdersResponseUnmarshaler
	CancelOrderResponseUnmarshaler           CancelOrderResponseUnmarshaler
	GetAccountResponseUnmarshaler            GetAccountResponseUnmarshaler
	GetPositionsResponseUnmarshaler          GetPositionsResponseUnmarshaler
	GetFuturesAccountResponseUnmarshaler     GetFuturesAccountResponseUnmarshaler
	GetExchangeInfoResponseUnmarshaler       GetExchangeInfoResponseUnmarshaler
	GetFundingRateResponseUnmarshaler        GetFundingRateResponseUnmarshaler
	GetFundingRateHistoryResponseUnmarshaler GetFundingRateHistoryResponseUnmarshaler
}

type UnmarshalerOption func(options *UnmarshalerOptions)

func WithResponseUnmarshaler(unmarshaler ResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.ResponseUnmarshaler = unmarshaler
	}
}

func WithTickerUnmarshaler(unmarshaler GetTickerResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.TickerUnmarshaler = unmarshaler
	}
}

func WithDepthUnmarshaler(unmarshaler GetDepthResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.DepthUnmarshaler = unmarshaler
	}
}

func WithKlineUnmarshaler(unmarshaler GetKlineResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.KlineUnmarshaler = unmarshaler
	}
}

func WithGetOrderInfoResponseUnmarshaler(unmarshaler GetOrderInfoResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.GetOrderInfoResponseUnmarshaler = unmarshaler
	}
}

func WithCreateOrderResponseUnmarshaler(unmarshaler CreateOrderResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.CreateOrderResponseUnmarshaler = unmarshaler
	}
}

func WithGetPendingOrdersResponseUnmarshaler(unmarshaler GetPendingOrdersResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.GetPendingOrdersResponseUnmarshaler = unmarshaler
	}
}

func WithCancelOrderResponseUnmarshaler(unmarshaler CancelOrderResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.CancelOrderResponseUnmarshaler = unmarshaler
	}
}

func WithGetHistoryOrdersResponseUnmarshaler(unmarshaler GetHistoryOrdersResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.GetHistoryOrdersResponseUnmarshaler = unmarshaler
	}
}

func WithGetAccountResponseUnmarshaler(unmarshaler GetAccountResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.GetAccountResponseUnmarshaler = unmarshaler
	}
}
func WithGetPositionsResponseUnmarshaler(unmarshaler GetPositionsResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.GetPositionsResponseUnmarshaler = unmarshaler
	}
}

func WithGetFuturesAccountResponseUnmarshaler(unmarshaler GetFuturesAccountResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.GetFuturesAccountResponseUnmarshaler = unmarshaler
	}
}

func WithGetExchangeInfoResponseUnmarshaler(unmarshaler GetExchangeInfoResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.GetExchangeInfoResponseUnmarshaler = unmarshaler
	}
}

func WithGetFundingRateResponseUnmarshaler(unmarshaler GetFundingRateResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.GetFundingRateResponseUnmarshaler = unmarshaler
	}
}

func WithGetFundingRateHistoryResponseUnmarshaler(unmarshaler GetFundingRateHistoryResponseUnmarshaler) UnmarshalerOption {
	return func(options *UnmarshalerOptions) {
		options.GetFundingRateHistoryResponseUnmarshaler = unmarshaler
	}
}
