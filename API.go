package coinapi

// api interface

type API interface {
	LimitBuy(amount, price string, currency CurrencyPair) (*Order, error);
	LimitSell(amount, price string, currency CurrencyPair) (*Order, error);
	MarketBuy(amount, price string, currency CurrencyPair) (*Order, error)
	MarketSell(amount, price string, currency CurrencyPair) (*Order, error)
	CancelOrder(orderId string, currency CurrencyPair) (bool, error);
	GetOneOrder(orderId string, currency CurrencyPair) (*Order, error);
	GetUnfinishOrders(currency CurrencyPair) ([]Order, error);
	GetOrderHistorys(currency CurrencyPair , currentPage , pageSize int)([]Order , error);
	GetAccount() (*Account, error);

	GetTicker(currency CurrencyPair) (*Ticker, error);
	GetDepth(size int, currency CurrencyPair) (*Depth, error);
	GetKlineRecords(currency CurrencyPair ,period string, size , since int) ([]Kline , error);

	GetExchangeName() string;
}
