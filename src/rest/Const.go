package rest

type CurrencyPair int;

type TradeSide int;

type TradeStatus int;

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
