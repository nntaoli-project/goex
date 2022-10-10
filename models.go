package goex

import (
	"time"
)

// OptionParameter 可选参数
type OptionParameter struct {
	Key   string
	Value string
}

type KlinePeriod int

type CurrencyPair struct {
	Symbol         string  `json:"symbol"`
	Market         string  `json:"market"`
	PricePrecision int     `json:"price_precision"` //价格小数点位数
	QtyPrecision   int     `json:"qty_precision"`   //数量小数点位数
	MinQty         float64 `json:"min_qty"`
	MaxQty         float64 `json:"max_qty"`
	MarketQty      float64 `json:"market_qty"`
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
}

type Order struct {
	Pair        CurrencyPair `json:"pair"`
	Id          string       `json:"id"`   //订单ID
	CId         string       `json:"c_id"` //客户端自定义ID
	Side        int          //交易方向: sell,buy
	Status      int          `json:"status"`     //状态
	OrderType   int          `json:"order_type"` //类型: limit , market , ...
	Price       float64      `json:"price"`
	Qty         float64      `json:"qty"`
	ExecutedQty float64      `json:"executed_qty"`
	PriceAvg    float64      `json:"price_avg"`
	Timestamp   int64        `json:"t"`
	Origin      []byte       `json:"origin"`
}
