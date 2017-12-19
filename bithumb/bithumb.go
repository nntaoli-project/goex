package bithumb

import (
	"fmt"
	"github.com/btcsuite/goleveldb/leveldb/errors"
	. "github.com/nntaoli-project/GoEx"
	"net/http"
	"sort"
)

type Bithumb struct {
	client *http.Client
	accesskey,
	secretkey string
}

var (
	baseUrl = "https://api.bithumb.com"
)

func New(client *http.Client, accesskey, secretkey string) *Bithumb {
	return &Bithumb{client: client, accesskey: accesskey, secretkey: secretkey}
}

func (bit *Bithumb) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (bit *Bithumb) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (bit *Bithumb) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (bit *Bithumb) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (bit *Bithumb) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("not implement")
}
func (bit *Bithumb) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (bit *Bithumb) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implement")
}
func (bit *Bithumb) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implement")
}
func (bit *Bithumb) GetAccount() (*Account, error) {
	panic("not implement")
}

func (bit *Bithumb) GetTicker(currency CurrencyPair) (*Ticker, error) {
	respmap, err := HttpGet(bit.client, fmt.Sprintf("%s/public/ticker/%s", baseUrl, currency.CurrencyA))
	if err != nil {
		return nil, err
	}

	if respmap["status"].(string) != "0000" {
		return nil, errors.New(respmap["status"].(string))
	}

	datamap := respmap["data"].(map[string]interface{})

	return &Ticker{
		Low:  ToFloat64(datamap["min_price"]),
		High: ToFloat64(datamap["max_price"]),
		Last: ToFloat64(datamap["closing_price"]),
		Vol:  ToFloat64(datamap["units_traded"]),
		Buy:  ToFloat64(datamap["buy_price"]),
		Sell: ToFloat64(datamap["sell_price"]),
	}, nil
}

func (bit *Bithumb) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	resp, err := HttpGet(bit.client, fmt.Sprintf("%s/public/orderbook/%s", baseUrl, currency.CurrencyA))
	if err != nil {
		return nil, err
	}

	if resp["status"].(string) != "0000" {
		return nil, errors.New(resp["status"].(string))
	}

	datamap := resp["data"].(map[string]interface{})
	bids := datamap["bids"].([]interface{})
	asks := datamap["asks"].([]interface{})

	dep := new(Depth)

	for _, v := range bids {
		bid := v.(map[string]interface{})
		dep.BidList = append(dep.BidList, DepthRecord{ToFloat64(bid["price"]), ToFloat64(bid["quantity"])})
	}

	for _, v := range asks {
		ask := v.(map[string]interface{})
		dep.AskList = append(dep.AskList, DepthRecord{ToFloat64(ask["price"]), ToFloat64(ask["quantity"])})
	}

	sort.Sort(sort.Reverse(dep.AskList))

	return dep, nil
}

func (bit *Bithumb) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录
func (bit *Bithumb) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}

func (bit *Bithumb) GetExchangeName() string {
	return "gate.io"
}
