package bigone

import (
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"net/http"

	"log"
	"time"
)

const (
	EXCHANGE_NAME = "big.one"

	API_BASE_URL = "https://api.big.one/"

	TICKER_URI = "markets/%s"
	//DEPTH_URI              = "depth.php?c=%s&mk_type=%s"
	//ACCOUNT_URI            = "getMyBalance.php"
	//TRADE_URI              = "trades.php?c=%s&mk_type=%s"
	//CANCEL_URI             = "cancelOrder.php"
	//ORDERS_INFO            = "getMyTradeList.php"
	//UNFINISHED_ORDERS_INFO = "getOrderList.php"

)

type Bigone struct {
	accessKey,
	secretKey string
	httpClient *http.Client
}

func New(client *http.Client, api_key, secret_key string) *Bigone {
	return &Bigone{api_key, secret_key, client}
}

func (bo *Bigone) GetExchangeName() string {
	return EXCHANGE_NAME
}

func (bo *Bigone) GetTicker(currency CurrencyPair) (*Ticker, error) {
	tickerUri := API_BASE_URL + fmt.Sprintf(TICKER_URI, currency.ToSymbol("-"))
	bodyDataMap, err := HttpGet(bo.httpClient, tickerUri)

	if err != nil {
		return nil, err
	}

	log.Println("bo uri:", tickerUri)
	log.Println("bo bodyDataMap:", currency, bodyDataMap)

	dataMap := bodyDataMap["data"].(map[string]interface{})
	tickerMap := dataMap["ticker"].(map[string]interface{})
	asksMap := dataMap["asks"].([]interface{})
	bidsMap := dataMap["bids"].([]interface{})
	ask := asksMap[0].(map[string]interface{})
	bid := bidsMap[0].(map[string]interface{})

	log.Println(tickerMap)
	log.Println(asksMap)
	log.Println(bidsMap)

	var ticker Ticker

	ticker.Date = uint64(time.Now().Unix())
	ticker.Last = ToFloat64(tickerMap["price"])
	ticker.Buy = ToFloat64(bid["price"])
	ticker.Sell = ToFloat64(ask["price"])
	ticker.Low = ToFloat64(tickerMap["low"])
	ticker.High = ToFloat64(tickerMap["high"])
	ticker.Vol = ToFloat64(tickerMap["volume"])
	return &ticker, nil
}
func (bo *Bigone) GetTickers(currency CurrencyPair) (*Ticker, error) {
	return bo.GetTicker(currency)
}
func (bo *Bigone) GetTickerInBuf(currency CurrencyPair) (*Ticker, error) {
	return bo.GetTicker(currency)
}

func (bo *Bigone) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (bo *Bigone) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (bo *Bigone) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (bo *Bigone) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (bo *Bigone) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("not implements")
}

func (bo *Bigone) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}
func (bo *Bigone) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implements")
}

func (bo *Bigone) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implements")
}

func (bo *Bigone) GetAccount() (*Account, error) {
	panic("not implements")
}

func (bo *Bigone) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	return nil, nil
}

func (bo *Bigone) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implements")
}

//非个人，整个交易所的交易记录
func (bo *Bigone) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implements")
}
