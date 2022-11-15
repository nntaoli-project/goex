package goex

type ResponseUnmarshaler func([]byte, interface{}) error
type GetTickerResponseUnmarshaler func([]byte) (*Ticker, error)
type GetDepthResponseUnmarshaler func([]byte) (*Depth, error)
type GetKlineResponseUnmarshaler func([]byte) ([]Kline, error)
type CreateOrderResponseUnmarshaler func([]byte) (*Order, error)
type GetOrderInfoResponseUnmarshaler func([]byte) (*Order, error)
type GetPendingOrdersResponseUnmarshaler func([]byte) ([]Order, error)
type CancelOrderResponseUnmarshaler func([]byte) error
type GetHistoryOrdersResponseUnmarshaler func([]byte) ([]Order, error)
type GetAccountResponseUnmarshaler func([]byte) (map[string]Account, error)

type UnmarshalerOptions struct {
	ResponseUnmarshaler                 ResponseUnmarshaler
	TickerUnmarshaler                   GetTickerResponseUnmarshaler
	DepthUnmarshaler                    GetDepthResponseUnmarshaler
	KlineUnmarshaler                    GetKlineResponseUnmarshaler
	CreateOrderResponseUnmarshaler      CreateOrderResponseUnmarshaler
	GetOrderInfoResponseUnmarshaler     GetOrderInfoResponseUnmarshaler
	GetPendingOrdersResponseUnmarshaler GetPendingOrdersResponseUnmarshaler
	GetHistoryOrdersResponseUnmarshaler GetHistoryOrdersResponseUnmarshaler
	CancelOrderResponseUnmarshaler      CancelOrderResponseUnmarshaler
	GetAccountResponseUnmarshaler       GetAccountResponseUnmarshaler
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
