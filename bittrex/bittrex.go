package bittrex

import (
	"fmt"
	. "github.com/nntaoli-project/goex"
	"net/http"
	"sort"
	"errors"
)

type Bittrex struct {
	client *http.Client
	baseUrl,
	accesskey,
	secretkey string
}

func New(client *http.Client, accesskey, secretkey string) *Bittrex {
	return &Bittrex{client: client, accesskey: accesskey, secretkey: secretkey, baseUrl: "https://bittrex.com/api/v1.1"}
}

func (bx *Bittrex) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (bx *Bittrex) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (bx *Bittrex) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (bx *Bittrex) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (bx *Bittrex) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("not implement")
}
func (bx *Bittrex) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (bx *Bittrex) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implement")
}
func (bx *Bittrex) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implement")
}
func (bx *Bittrex) GetAccount() (*Account, error) {
	panic("not implement")
}

func (bx *Bittrex) GetTicker(currency CurrencyPair) (*Ticker, error) {
	resp, err := HttpGet(bx.client, fmt.Sprintf("%s/public/getmarketsummary?market=%s", bx.baseUrl, currency.ToSymbol2("-")))
	if err != nil {
		errCode := HTTP_ERR_CODE
		errCode.OriginErrMsg = err.Error()
		return nil, errCode
	}

	result, _ := resp["result"].([]interface{})
	if len(result) <= 0 {
		return nil, API_ERR
	}

	tickermap := result[0].(map[string]interface{})

	return &Ticker{
		Last: ToFloat64(tickermap["Last"]),
		Sell: ToFloat64(tickermap["Ask"]),
		Buy:  ToFloat64(tickermap["Bid"]),
		Low:  ToFloat64(tickermap["Low"]),
		High: ToFloat64(tickermap["High"]),
		Vol:  ToFloat64(tickermap["Volume"]),
	}, nil
}

func (bx *Bittrex) GetDepth(size int, currency CurrencyPair) (*Depth, error) {

	resp, err := HttpGet(bx.client, fmt.Sprintf("%s/public/getorderbook?market=%s&type=both", bx.baseUrl, currency.ToSymbol2("-")))
	if err != nil {
		errCode := HTTP_ERR_CODE
		errCode.OriginErrMsg = err.Error()
		return nil, errCode
	}

	result, err2 := resp["result"].(map[string]interface{})
	if err2 != true {
		return nil, errors.New(resp["message"].(string))
	}
	bids, _ := result["buy"].([]interface{})
	asks, _ := result["sell"].([]interface{})

	dep := new(Depth)

	for _, v := range bids {
		r := v.(map[string]interface{})
		dep.BidList = append(dep.BidList, DepthRecord{ToFloat64(r["Rate"]), ToFloat64(r["Quantity"])})
	}

	for _, v := range asks {
		r := v.(map[string]interface{})
		dep.AskList = append(dep.AskList, DepthRecord{ToFloat64(r["Rate"]), ToFloat64(r["Quantity"])})
	}

	sort.Sort(sort.Reverse(dep.AskList))

	return dep, nil
}

func (bx *Bittrex) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录
func (bx *Bittrex) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}

func (bx *Bittrex) GetExchangeName() string {
	return BITTREX
}
