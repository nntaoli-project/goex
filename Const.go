package coinapi

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
	ETC
)

const
(
	BTC_CNY = 1 + iota
	BTC_USD
	BTC_LTC

	LTC_CNY
	LTC_USD

	ETH_CNY
	ETH_USD
	ETH_BTC

	ETC_CNY
	ETC_USD
	ETC_BTC

	XCN_BTC
	SYS_BTC
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
	ORDER_CANCEL_ING
)

const
(
	OPEN_BUY = 1 + iota  //开多
	OPEN_SELL              //开空
	CLOSE_BUY             //平多
	CLOSE_SELL           //平空
)

var CurrencyPairSymbol = map[CurrencyPair]string{
	BTC_CNY : "btc_cny",
	BTC_USD : "btc_usd",
	BTC_LTC : "btc_ltc",
	LTC_CNY : "ltc_cny",
	LTC_USD : "ltc_usd",
	ETH_CNY : "eth_cny",
	ETH_USD : "eth_usd",
	ETH_BTC : "eth_btc",
	ETC_CNY : "etc_cny",
	ETC_USD : "etc_usd",
	ETC_BTC : "etc_btc"};

var
(
	THIS_WEEK_CONTRACT = "this_week";//周合约
	NEXT_WEEK_CONTRACT = "next_week";//次周合约
	QUARTER_CONTRACT = "quarter";//季度合约
)

