package model

const (
	Kline_1min  KlinePeriod = "1min"
	Kline_5min              = "5min"
	Kline_15min             = "15min"
	Kline_30min             = "30min"
	Kline_60min             = "60min"
	Kline_1h                = "1h"
	Kline_4h                = "4h"
	Kline_6h                = "6h"
	Kline_1day              = "1day"
	Kline_1week             = "1week"
)

const (
	OrderStatus_Pending      OrderStatus = 1
	OrderStatus_Finished                 = 2
	OrderStatus_Canceled                 = 3
	OrderStatus_PartFinished             = 4
	OrderStatus_Canceling                = 5
)

const (
	Spot_Buy          OrderSide = "buy"
	Spot_Sell         OrderSide = "sell"
	Futures_OpenBuy   OrderSide = "futures_open_buy"
	Futures_OpenSell  OrderSide = "futures_open_sell"
	Futures_CloseBuy  OrderSide = "futures_close_buy"
	Futures_CloseSell OrderSide = "futures_close_sell"
)

const (
	OrderType_Limit    OrderType = "limit"
	OrderType_Market   OrderType = "market"
	OrderType_opponent OrderType = "opponent"
)

//coin const list
//a-z排序
const (
	ADA  = "ADA"
	ATOM = "ATOM"
	AAVE = "AAVE"
	ALGO = "ALGO"
	AR   = "AR"

	BTC  = "BTC"
	BNB  = "BNB"
	BSV  = "BSV"
	BCH  = "BCH"
	BUSD = "BUSD"

	CEL = "CEL"
	CRV = "CRV"

	DAI  = "DAI"
	DCR  = "DCR"
	DOT  = "DOT"
	DOGE = "DOGE"
	DASH = "DASH"
	DYDX = "DYDX"

	ETH  = "ETH"
	ETHW = "ETHW"
	ETC  = "ETC"
	EOS  = "EOS"
	ENJ  = "ENJ"
	ENS  = "ENS"

	FLOW = "FLOW"
	FIL  = "FIL"
	FLM  = "FLM"

	GALA = "GALA"
	GAS  = "GAS"

	HT = "HT"

	IOTA = "IOTA"
	IOST = "IOST"

	KSM = "KSM"

	LTC = "LTC"
	LDO = "LDO"

	MINA = "MINA"
	MEME = "MEME"

	NEO  = "NEO"
	NEAR = "NEAR"

	OP   = "OP"
	OKB  = "OKB"
	OKT  = "OKT"
	ORDI = "ORDI"

	PLG  = "PLG"
	PERP = "PERP"
	PEPE = "PEPE"

	QTUM = "QTUM"

	RACA = "RACA"
	RVN  = "RVN"

	STORJ = "STORJ"
	SOL   = "SOL"
	SHIB  = "SHIB"
	SC    = "SC"
	SAND  = "SAND"
	SUSHI = "SUSHI"
	SUI   = "SUI"

	TRX   = "TRX"
	TRADE = "TRADE"
	TRB   = "TRB"

	USD  = "USD"
	USDT = "USDT"
	USDC = "USDC"
	UNI  = "UNI"

	VELO = "VELO"

	WBTC  = "WBTC"
	WAVES = "WAVES"

	XRP = "XRP"
	XTZ = "XTZ"

	YFI  = "YFI"
	YFII = "YFII"

	ZEC  = "ZEC"
	ZYRO = "ZYRO"
)

//exchange name const list
const (
	OKX     = "okx.com"
	BINANCE = "binance.com"
)

const (
	Order_Client_ID__Opt_Key = "OrderClientID"
)
