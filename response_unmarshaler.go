package goex

type Unmarshal interface {
	Unmarshal([]byte, interface{}) error
}

type UnmarshalTicker interface {
	UnmarshalTicker([]byte) (*Ticker, error)
}

type UnmarshalDepth interface {
	UnmarshalDepth([]byte) (*Depth, error)
}
