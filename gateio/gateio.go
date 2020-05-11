package gateio

import (
	"fmt"
	. "github.com/nntaoli-project/goex"
	"net/http"
	"sort"
	"strings"
)

var (
	marketBaseUrl = "http://data.gate.io/api2/1"
)

type Gate struct {
	client *http.Client
	accesskey,
	secretkey string
}

func New(client *http.Client, accesskey, secretkey string) *Gate {
	return &Gate{client: client, accesskey: accesskey, secretkey: secretkey}
}

func (g *Gate) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (g *Gate) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (g *Gate) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (g *Gate) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (g *Gate) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("not implement")
}
func (g *Gate) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (g *Gate) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implement")
}
func (g *Gate) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implement")
}
func (g *Gate) GetAccount() (*Account, error) {
	panic("not implement")
}

func (g *Gate) GetTicker(currency CurrencyPair) (*Ticker, error) {
	uri := fmt.Sprintf("%s/ticker/%s", marketBaseUrl, strings.ToLower(currency.ToSymbol("_")))

	resp, err := HttpGet(g.client, uri)
	if err != nil {
		errCode := HTTP_ERR_CODE
		errCode.OriginErrMsg = err.Error()
		return nil, errCode
	}

	return &Ticker{
		Last: ToFloat64(resp["last"]),
		Sell: ToFloat64(resp["lowestAsk"]),
		Buy:  ToFloat64(resp["highestBid"]),
		High: ToFloat64(resp["high24hr"]),
		Low:  ToFloat64(resp["low24hr"]),
		Vol:  ToFloat64(resp["quoteVolume"]),
	}, nil
}

func (g *Gate) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	resp, err := HttpGet(g.client, fmt.Sprintf("%s/orderBook/%s", marketBaseUrl, currency.ToSymbol("_")))
	if err != nil {
		errCode := HTTP_ERR_CODE
		errCode.OriginErrMsg = err.Error()
		return nil, errCode
	}

	bids, _ := resp["bids"].([]interface{})
	asks, _ := resp["asks"].([]interface{})

	dep := new(Depth)

	for _, v := range bids {
		r := v.([]interface{})
		dep.BidList = append(dep.BidList, DepthRecord{ToFloat64(r[0]), ToFloat64(r[1])})
	}

	for _, v := range asks {
		r := v.([]interface{})
		dep.AskList = append(dep.AskList, DepthRecord{ToFloat64(r[0]), ToFloat64(r[1])})
	}

	sort.Sort(sort.Reverse(dep.AskList))

	return dep, nil
}

func (g *Gate) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录
func (g *Gate) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}

func (g *Gate) GetExchangeName() string {
	return GATEIO
}
