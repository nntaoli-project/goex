package bitmex

import (
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"net/http"
	"strings"
	"time"
)

var (
	base_url = "https://www.bitmex.com/api/v1/"
)

//bitmex register link  https://www.bitmex.com/register/0fcQP7

type Bitmex struct {
	httpClient *http.Client
	accessKey,
	secretKey string
}

func New(client *http.Client, accesskey, secretkey string) *Bitmex {
	return &Bitmex{client, accesskey, secretkey}
}

func (Bitmex *Bitmex) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (Bitmex *Bitmex) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (Bitmex *Bitmex) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (Bitmex *Bitmex) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (Bitmex *Bitmex) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("not implements")
}

func (Bitmex *Bitmex) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}
func (Bitmex *Bitmex) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implements")
}

func (Bitmex *Bitmex) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implements")
}

func (Bitmex *Bitmex) GetAccount() (*Account, error) {
	panic("not implements")
}

func (Bitmex *Bitmex) GetTicker(currency CurrencyPair) (*Ticker, error) {
	panic("not implements")
}

func (Bitmex *Bitmex) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	uri := fmt.Sprintf("orderBook/L2?symbol=%s&depth=%d", Bitmex.pairToSymbol(currency), size)
	resp, err := HttpGet3(Bitmex.httpClient, base_url+uri, nil)
	if err != nil {
		return nil, HTTP_ERR_CODE.OriginErr(err.Error())
	}

	//log.Println(resp)

	dep := new(Depth)
	dep.UTime = time.Now()
	dep.Pair = currency

	for _, r := range resp {
		rr := r.(map[string]interface{})
		switch strings.ToLower(rr["side"].(string)) {
		case "sell":
			dep.AskList = append(dep.AskList, DepthRecord{Price: ToFloat64(rr["price"]), Amount: ToFloat64(rr["size"])})
		case "buy":
			dep.BidList = append(dep.BidList, DepthRecord{Price: ToFloat64(rr["price"]), Amount: ToFloat64(rr["size"])})
		}
	}

	return dep, nil
}

func (Bitmex *Bitmex) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implements")
}

//非个人，整个交易所的交易记录
func (Bitmex *Bitmex) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implements")
}

func (Bitmex *Bitmex) GetExchangeName() string {
	return BITMEX
}

func (mex *Bitmex) pairToSymbol(pair CurrencyPair) string {
	if pair.CurrencyA.Symbol == BTC.Symbol {
		return NewCurrencyPair(XBT, USD).ToSymbol("")
	}
	return pair.AdaptUsdtToUsd().ToSymbol("")
}
