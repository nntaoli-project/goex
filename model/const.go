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
