package rest

// api interface

type API interface {
	LimitBuy(amount, price string, currency CurrencyPair) (*Order, error);
	LimitSell(amount, price string, currency CurrencyPair) (*Order, error);
	CancelOrder(orderId string, currency CurrencyPair) (bool, error);
	GetOneOrder(orderId string, currency CurrencyPair) (*Order, error);
	GetUnfinishOrders(currency CurrencyPair) ([]Order, error);

	GetAccount() (*Account, error);

	GetTicker(currency CurrencyPair) (*Ticker, error);
	GetDepth(size int, currency CurrencyPair) (*Depth, error);

	GetExchangeName() string;
}
