package c_cex

import (
	. "github.com/nntaoli-project/GoEx"
	//"log"
	"net/http"
	"strings"
	"time"
)

const (
	EXCHANGE_NAME = "c-cex.com"

	API_BASE_URL = "https://c-cex.com/"

	TICKER_URI = "t/"
)

type C_cex struct {
	accessKey,
	secretKey string
	httpClient *http.Client
}

func New(client *http.Client, accessKey, secretKey string) *C_cex {
	return &C_cex{accessKey, secretKey, client}
}

func (ccex *C_cex) GetExchangeName() string {
	return EXCHANGE_NAME
}

func (ccex *C_cex) GetTicker(currency CurrencyPair) (*Ticker, error) {
	currency = ccex.adaptCurrencyPair(currency)

	tickerUri := API_BASE_URL + TICKER_URI + strings.ToLower(currency.ToSymbol("-")) + ".json"
	//log.Println("tickerUrl:", tickerUri)
	bodyDataMap, err := HttpGet(ccex.httpClient, tickerUri)
	//log.Println("C_cex bodyDataMap:", tickerUri, bodyDataMap)

	timestamp := time.Now().Unix()
	if err != nil {
		//log.Println(err)
		return nil, err
	}

	tickerMap := bodyDataMap["ticker"].(map[string]interface{})
	var ticker Ticker

	//fmt.Println(bodyDataMap)
	ticker.Date = uint64(timestamp)
	ticker.Last, _ = tickerMap["lastprice"].(float64)
	ticker.Buy, _ = tickerMap["buy"].(float64)
	ticker.Sell, _ = tickerMap["sell"].(float64)
	//ticker.Vol, _ = tickerMap["Volume"].(float64)
	//log.Println("C_cex", currency, "ticker:", ticker)
	return &ticker, nil
}
func (ccex *C_cex) GetTickers(currency CurrencyPair) (*Ticker, error) {
	return ccex.GetTicker(currency)
}

func (ccex *C_cex) GetTickerInBuf(currency CurrencyPair) (*Ticker, error) {
	return ccex.GetTicker(currency)
}

func (ccex *C_cex) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	panic("not implement")
}
func (ccex *C_cex) adaptCurrencyPair(pair CurrencyPair) CurrencyPair {
	var currencyA Currency
	var currencyB Currency

	//if pair.QuoteCurrency == BCC {
	//	currencyA = BCH
	//} else {
	//	currencyA = pair.QuoteCurrency
	//}
	currencyA = pair.CurrencyA
	if pair.CurrencyA == USDT {
		currencyB = USD
	} else {
		currencyB = pair.CurrencyA
	}

	return NewCurrencyPair(currencyA, currencyB)
}

func (ccex *C_cex) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (ccex *C_cex) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (ccex *C_cex) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (ccex *C_cex) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (ccex *C_cex) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("not implements")
}

func (ccex *C_cex) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}
func (ccex *C_cex) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implements")
}

func (ccex *C_cex) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implements")
}

func (ccex *C_cex) GetAccount() (*Account, error) {
	panic("not implements")
}

func (ccex *C_cex) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implements")
}

//非个人，整个交易所的交易记录
func (ccex *C_cex) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implements")
}
