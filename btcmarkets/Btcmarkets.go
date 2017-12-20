package btcmarkets

import (
	"errors"
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"net/http"
	"time"
)

const (
	EXCHANGE_NAME = "btcmarkets.net"

	API_BASE_URL = "https://api.btcmarkets.net/"

	TICKER_URI = "market/%s/%s/tick"
)

type Btcmarkets struct {
	accessKey,
	secretKey string
	httpClient *http.Client
}

func New(client *http.Client, accessKey, secretKey string) *Btcmarkets {
	return &Btcmarkets{accessKey, secretKey, client}
}

func (btcm *Btcmarkets) GetExchangeName() string {
	return EXCHANGE_NAME
}

func (btcm *Btcmarkets) GetTicker(currency CurrencyPair) (*Ticker, error) {
	tickerUri := fmt.Sprintf(API_BASE_URL+TICKER_URI, currency.CurrencyA.String(), currency.CurrencyB.String())
	//log.Println("tickerUrl:", tickerUri)
	bodyDataMap, err := HttpGet(btcm.httpClient, tickerUri)
	//log.Println("Btcmarkets bodyDataMap:", tickerUri, bodyDataMap)

	timestamp := time.Now().Unix()
	if err != nil {
		//log.Println(err)
		return nil, err
	}

	if result, isok := bodyDataMap["success"].(bool); isok == true && result != true {
		//log.Println("bodyDataMap[\"success\"]", isok, result)
		return nil, errors.New("err")
	}

	var tickerMap map[string]interface{} = bodyDataMap
	var ticker Ticker

	//fmt.Println(bodyDataMap)
	ticker.Date = uint64(timestamp)
	ticker.Last, _ = tickerMap["lastPrice"].(float64)

	ticker.Buy, _ = tickerMap["bestBid"].(float64)
	ticker.Sell, _ = tickerMap["bestAsk"].(float64)
	ticker.Vol, _ = tickerMap["volume24h"].(float64)
	//log.Println("Btcmarkets", currency, "ticker:", ticker)
	return &ticker, nil
}
func (btcm *Btcmarkets) GetTickers(currency CurrencyPair) (*Ticker, error) {
	return btcm.GetTicker(currency)
}

func (btcm *Btcmarkets) GetTickerInBuf(currency CurrencyPair) (*Ticker, error) {
	return btcm.GetTicker(currency)
}

func (btcm *Btcmarkets) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	panic("not implement")
}

func (btcm *Btcmarkets) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (btcm *Btcmarkets) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (btcm *Btcmarkets) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (btcm *Btcmarkets) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (btcm *Btcmarkets) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("not implements")
}

func (btcm *Btcmarkets) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}
func (btcm *Btcmarkets) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implements")
}

func (btcm *Btcmarkets) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implements")
}

func (btcm *Btcmarkets) GetAccount() (*Account, error) {
	panic("not implements")
}

func (btcm *Btcmarkets) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implements")
}

//非个人，整个交易所的交易记录
func (btcm *Btcmarkets) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implements")
}
