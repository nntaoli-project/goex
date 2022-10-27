package goex

import (
	"time"
)

// OptionParameter 可选参数
type OptionParameter struct {
	Key   string
	Value string
}

type KlinePeriod string

type CurrencyPair struct {
	Symbol         string  `json:"symbol"`          //交易对
	BaseSymbol     string  `json:"base_symbol"`     //币种
	MarketSymbol   string  `json:"market"`          //交易市场
	PricePrecision int     `json:"price_precision"` //价格小数点位数
	QtyPrecision   int     `json:"qty_precision"`   //数量小数点位数
	MinQty         float64 `json:"min_qty"`
	MaxQty         float64 `json:"max_qty"`
	MarketQty      float64 `json:"market_qty"`
}

func (pair CurrencyPair) String() string {
	return pair.Symbol
}

type FuturesCurrencyPair struct {
	CurrencyPair
	DeliveryDate int64   //结算日期
	OnboardDate  int64   //上线日期
	MarginAsset  float64 //保证金资产
}

type Ticker struct {
	Pair      CurrencyPair `json:"pair"`
	Last      float64      `json:"l"`
	Buy       float64      `json:"b"`
	Sell      float64      `json:"s"`
	High      float64      `json:"h"`
	Low       float64      `json:"lw"`
	Vol       float64      `json:"v"`
	Percent   float64      `json:"percent"`
	Timestamp int64        `json:"t"`
	Origin    []byte       `json:"origin"`
}

type DepthItem struct {
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
}

type DepthItems []DepthItem

func (dr DepthItems) Len() int {
	return len(dr)
}

func (dr DepthItems) Swap(i, j int) {
	dr[i], dr[j] = dr[j], dr[i]
}

func (dr DepthItems) Less(i, j int) bool {
	return dr[i].Price < dr[j].Price
}

type Depth struct {
	Pair   CurrencyPair `json:"pair"`
	UTime  time.Time    `json:"ut"`
	Asks   DepthItems   `json:"asks"`
	Bids   DepthItems   `json:"bids"`
	Origin []byte       `json:"origin"`
}

type Kline struct {
	Pair      CurrencyPair `json:"pair"`
	Timestamp int64        `json:"t"`
	Open      float64      `json:"o"`
	Close     float64      `json:"s"`
	High      float64      `json:"h"`
	Low       float64      `json:"l"`
	Vol       float64      `json:"v"`
	Origin    []byte       `json:"-"`
}

type Order struct {
	Pair        CurrencyPair `json:"pair"`
	Id          string       `json:"id"`       //订单ID
	CId         string       `json:"c_id"`     //客户端自定义ID
	Side        OrderSide    `json:"side"`     //交易方向: sell,buy
	OrderTy     OrderType    `json:"order_ty"` //类型: limit , market , ...
	Status      int          `json:"status"`   //状态
	Price       float64      `json:"price"`
	Qty         float64      `json:"qty"`
	ExecutedQty float64      `json:"executed_qty"`
	PriceAvg    float64      `json:"price_avg"`
	Timestamp   int64        `json:"t"`
	Origin      []byte       `json:"origin"`
}

type OrderType struct {
	Code int
	Type string
}

func (ty OrderType) String() string {
	return ty.Type
}

var (
	OrderType_Limit    = OrderType{Code: 1, Type: "limit"}
	OrderType_Market   = OrderType{Code: 2, Type: "market"}
	OrderType_opponent = OrderType{Code: 3, Type: "opponent"}
)

type OrderSide struct {
	Code int
	Type string
}

func (s OrderSide) String() string {
	return s.Type
}

var (
	Spot_Buy          = OrderSide{Type: "buy", Code: 1}
	Spot_Sell         = OrderSide{Type: "sell", Code: 2}
	Futures_OpenBuy   = OrderSide{Type: "open_buy", Code: 1}
	Futures_OpenSell  = OrderSide{Type: "open_sell", Code: 2}
	Futures_CloseBuy  = OrderSide{Type: "close_buy", Code: 3}
	Futures_CloseSell = OrderSide{Type: "close_sell", Code: 4}
)
