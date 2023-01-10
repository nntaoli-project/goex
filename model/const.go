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
