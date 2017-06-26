package coinapi

type CurrencyPair int;

func (c CurrencyPair) String() string {
	if c == 0 {
		return "nil"
	}
	return currencyPairSymbol[c - 1];
}

type Currency int;

func (c Currency) String() string {
	if c == 0 {
		return "nil"
	}
	return currencySymbol[c - 1];
}

type TradeSide int;

func (ts TradeSide)String() string {
	switch ts {
	case 1:
		return "BUY";
	case 2:
		return "SELL";
	case 3:
		return "BUY_MARKET";
	case 4:
		return "SELL_MARKET";
	default:
		return "UNKNOWN";
	}
}

type TradeStatus int;

func (ts TradeStatus) String() string {
	return orderStatusSymbol[ts];
}

/*currencySymbol array*/
var currencySymbol = [...]string{"cny", "usd", "btc", "ltc", "eth", "etc", "zec", "sc" , "bts"};

const
(
	CNY = 1 + iota
	USD
	BTC
	LTC
	ETH
	ETC
	ZEC
	SC
	BTS
)

var currencyPairSymbol = [...]string{"btc_cny", "btc_usd", "btc_jpy", "fx_btc_jpy", "ltc_cny", "ltc_usd", "eth_cny",
	"eth_usd", "eth_btc", "etc_cny", "etc_usd", "etc_btc", "xcn_btc", "sys_btc", "zec_cny", "zec_usd", "zec_btc" , "bts_cny","bts_btc"};

const
(
	BTC_CNY = 1 + iota
	BTC_USD
	BTC_JPY
	FX_BTC_JPY

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

	ZEC_CNY
	ZEC_USD
	ZEC_BTC

	BTS_CNY
	BTS_BTC
)

const
(
	BUY = 1 + iota
	SELL
	BUY_MARKET
	SELL_MARKET
)

var orderStatusSymbol = [...]string{"UNFINISH", "PART_FINISH", "FINISH", "CANCEL", "REJECT", "CANCEL_ING"}

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
	LTC_CNY : "ltc_cny",
	LTC_USD : "ltc_usd",
	ETH_CNY : "eth_cny",
	ETH_USD : "eth_usd",
	ETH_BTC : "eth_btc",
	ETC_CNY : "etc_cny",
	ETC_USD : "etc_usd",
	ETC_BTC : "etc_btc",
	BTS_CNY : "bts_cny"};

var
(
	THIS_WEEK_CONTRACT = "this_week"; //周合约
	NEXT_WEEK_CONTRACT = "next_week"; //次周合约
	QUARTER_CONTRACT = "quarter"; //季度合约
)

