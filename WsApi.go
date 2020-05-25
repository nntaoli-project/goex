package goex

type FuturesWsApi interface {
	DepthCallback(func(depth *Depth))
	TickerCallback(func(ticker *FutureTicker))
	TradeCallback(func(trade *Trade, contract string))
	//OrderCallback(func(order *FutureOrder))
	//PositionCallback(func(position *FuturePosition))
	//AccountCallback(func(account *FutureAccount))

	SubscribeDepth(pair CurrencyPair, contractType string) error
	SubscribeTicker(pair CurrencyPair, contractType string) error
	SubscribeTrade(pair CurrencyPair, contractType string) error

	//Login() error
	//SubscribeOrder(pair CurrencyPair, contractType string) error
	//SubscribePosition(pair CurrencyPair, contractType string) error
	//SubscribeAccount(pair CurrencyPair) error
}

type SpotWsApi interface {
	DepthCallback(func(depth *Depth))
	TickerCallback(func(ticker *Ticker))
	TradeCallback(func(trade *Trade))
	//OrderCallback(func(order *Order))
	//AccountCallback(func(account *Account))

	SubscribeDepth(pair CurrencyPair) error
	SubscribeTicker(pair CurrencyPair) error
	SubscribeTrade(pair CurrencyPair) error

	//Login() error
	//SubscribeOrder(pair CurrencyPair) error
	//SubscribeAccount(pair CurrencyPair) error
}
