package goex

type ResponseUnmarshaler interface {
	UnmarshalResponse([]byte, interface{}) error
}

type TickerUnmarshaler interface {
	UnmarshalTicker([]byte) (*Ticker, error)
}

type DepthUnmarshaler interface {
	UnmarshalDepth([]byte) (*Depth, error)
}
