package goex

import (
	"time"
)

type CurrencyPair struct {
	Symbol         string
	Market         string
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
