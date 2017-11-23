package cexio

import (
	"encoding/json"
	"errors"
	. "github.com/nntaoli-project/GoEx"
	"net/http"
	"strconv"
)

const (
	EXCHANGE_NAME = "cex.io"
)

func New(client *http.Client, accessKey, secretKey string, clientId string) *Cex {
	return &Cex{client, clientId, accessKey, secretKey}
}

func (cex *Cex) GetExchangeName() string {
	return EXCHANGE_NAME
}

func (cex *Cex) GetTicker(currency CurrencyPair) (*Ticker, error) {
	opt := currency.ToSymbol("/")
	respData := cex.Ticker(opt)
	if respData == nil {
		return nil, errors.New("ticker error")
	}
	resp, err := cex.convert(respData)
	if err != nil {
		return nil, err
	}
	ticker := new(Ticker)
	intTime, _ := strconv.Atoi(resp["timestamp"].(string))
	ticker.Last = ToFloat64(resp["last"])
	ticker.Vol = ToFloat64(resp["volume"])
	ticker.High = ToFloat64(resp["high"])
	ticker.Low = ToFloat64(resp["low"])
	ticker.Sell = ToFloat64(resp["ask"])
	ticker.Buy = ToFloat64(resp["bid"])
	ticker.Date = uint64(intTime)
	return ticker, nil
}

func (cex *Cex) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}

func (cex *Cex) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}

func (cex *Cex) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}

func (cex *Cex) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}

func (cex *Cex) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("not implement")
}

func (cex *Cex) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}

func (cex *Cex) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implement")
}

func (cex *Cex) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implement")
}

func (cex *Cex) GetAccount() (*Account, error) {
	panic("not implement")
}

func (cex *Cex) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	opt := currency.ToSymbol("/")
	respData := cex.OrderBook(size, opt)
	//fmt.Println(string(respData))
	resp, err := cex.convert(respData)
	if err != nil {
		return nil, err
	}
	bids := resp["bids"].([]interface{})
	asks := resp["asks"].([]interface{})

	depth := new(Depth)

	for _, bid := range bids {
		_bid := bid.([]interface{})
		amount := ToFloat64(_bid[1].(float64))
		price := ToFloat64(_bid[0].(float64))
		dr := DepthRecord{Amount: amount, Price: price}
		depth.BidList = append(depth.BidList, dr)
	}

	for _, ask := range asks {
		_ask := ask.([]interface{})
		amount := ToFloat64(_ask[1].(float64))
		price := ToFloat64(_ask[0].(float64))
		dr := DepthRecord{Amount: amount, Price: price}
		depth.AskList = append(depth.AskList, dr)
	}

	return depth, nil
}

func (cex *Cex) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录
func (cex *Cex) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}

func (cex *Cex) convert(data []byte) (map[string]interface{}, error) {
	var resp map[string]interface{}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
