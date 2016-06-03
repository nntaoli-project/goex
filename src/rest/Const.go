package rest

type CurrencyPair int;

type Currency int;

type TradeSide int;

type TradeStatus int;

const
(
	CNY = 1 + iota
	USD
	BTC
	LTC
	ETH
)

const
(
	BTC_CNY = 1 + iota
	BTC_USD
	BTC_LTC
	BTC_ETH

	LTC_CNY
	LTC_USD

	ETH_CNY
	ETH_USD
)

const
(
	BUY = 1 + iota
	SELL
)

const
(
	ORDER_UNFINISH = iota
	ORDER_PART_FINISH
	ORDER_FINISH
	ORDER_CANCEL
	ORDER_REJECT
)

var CurrencyPairSymbol = map[CurrencyPair]string{
	BTC_CNY : "btc_cny",
	BTC_USD : "btc_usd",
	BTC_LTC : "btc_ltc",
	BTC_ETH : "btc_eth",
	LTC_CNY : "ltc_cny",
	LTC_USD : "ltc_usd",
	ETH_CNY : "eth_cny",
	ETH_USD : "eth_usd" };