package goex

// api interface

type API interface {
	LimitBuy(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error)
	LimitSell(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error)
	MarketBuy(amount, price string, currency CurrencyPair) (*Order, error)
	MarketSell(amount, price string, currency CurrencyPair) (*Order, error)
	CancelOrder(orderId string, currency CurrencyPair) (bool, error)
	GetOneOrder(orderId string, currency CurrencyPair) (*Order, error)
	GetUnfinishOrders(currency CurrencyPair) ([]Order, error)
	GetOrderHistorys(currency CurrencyPair, opt ...OptionalParameter) ([]Order, error)
	GetAccount() (*Account, error)

	GetTicker(currency CurrencyPair) (*Ticker, error)
	GetDepth(size int, currency CurrencyPair) (*Depth, error)
	GetKlineRecords(currency CurrencyPair, period KlinePeriod, size int, optional ...OptionalParameter) ([]Kline, error)
	//非个人，整个交易所的交易记录
	GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error)

	GetExchangeName() string
}
