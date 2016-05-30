package rest

// api interface

type API interface {
	LimitBuy(amount, price string, currency CurrencyPair) (string, error); //if success return orderId
	LimitSell(amount, price string, currency CurrencyPair) (string, error); //if success return orderId
	CancelOrder(orderId string, currency CurrencyPair) (string, error);
	GetOneOrder(orderId string, currency CurrencyPair) (Order, error);
	GetUnfinishOrders(currency CurrencyPair) ([]Order, error);

	GetAccount() (Account, error);

	GetTicker(currency CurrencyPair) (Ticker, error);
	GetDepth(currency CurrencyPair) (Depth, error);

	GetExchangeName() string;
}
