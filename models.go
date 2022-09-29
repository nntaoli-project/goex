package goex

type CurrencyPair struct {
	Symbol         string
	PricePrecision int //价格小数点位数
	QtyPrecision   int //数量小数点位数
	MinQty         float64
	MaxQty         float64
	MarketQty      float64
}

type FuturesCurrencyPair struct {
	CurrencyPair
	DeliveryDate int64   //结算日期
	OnboardDate  int64   //上线日期
	MarginAsset  float64 //保证金资产
}
