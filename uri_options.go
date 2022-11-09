package goex

type UriOptions struct {
	Endpoint            string
	TickerUri           string
	DepthUri            string
	KlineUri            string
	GetOrderUri         string
	GetPendingOrdersUri string
	GetHistoryOrdersUri string
	CancelOrderUri      string
	NewOrderUri         string
}

type UriOption func(*UriOptions)

func WithEndpoint(endpoint string) UriOption {
	return func(c *UriOptions) {
		c.Endpoint = endpoint
	}
}

func WithTickerUri(uri string) UriOption {
	return func(c *UriOptions) {
		c.TickerUri = uri
	}
}

func WithDepthUri(uri string) UriOption {
	return func(c *UriOptions) {
		c.DepthUri = uri
	}
}

func WithKlineUri(uri string) UriOption {
	return func(c *UriOptions) {
		c.KlineUri = uri
	}
}

func WithGetOrderUri(uri string) UriOption {
	return func(c *UriOptions) {
		c.GetOrderUri = uri
	}
}

func WithGetPendingOrdersUri(uri string) UriOption {
	return func(c *UriOptions) {
		c.GetPendingOrdersUri = uri
	}
}

func WithCancelOrderUri(uri string) UriOption {
	return func(c *UriOptions) {
		c.CancelOrderUri = uri
	}
}

func WithNewOrderUri(uri string) UriOption {
	return func(c *UriOptions) {
		c.NewOrderUri = uri
	}
}

func WithGetHistoryOrdersUri(uri string) UriOption {
	return func(c *UriOptions) {
		c.GetHistoryOrdersUri = uri
	}
}
