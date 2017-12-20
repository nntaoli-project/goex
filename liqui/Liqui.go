package liqui

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	. "github.com/nntaoli-project/GoEx"
	"strings"
)

const (
	EXCHANGE_NAME = "liqui.io"

	API_BASE_URL = "https://api.liqui.io/"
	API_V1       = API_BASE_URL + "api/3/"
	TICKER_URI   = "ticker/%s"
)

type Liqui struct {
	accessKey,
	secretKey string
	httpClient *http.Client
}

func New(client *http.Client, accessKey, secretKey string) *Liqui {
	return &Liqui{accessKey, secretKey, client}
}

func (liqui *Liqui) GetExchangeName() string {
	return EXCHANGE_NAME
}

func (liqui *Liqui) GetTicker(currency CurrencyPair) (*Ticker, error) {
	cur := strings.ToLower(currency.ToSymbol("_"))
	if cur == "nil" {
		log.Println("Unsupport The CurrencyPair")
		return nil, errors.New("Unsupport The CurrencyPair")
	}
	tickerUri := API_V1 + fmt.Sprintf(TICKER_URI, cur)
	bodyDataMap, err := HttpGet(liqui.httpClient, tickerUri)
	//fmt.Println("tickerUri:", tickerUri)
	//fmt.Println("Liqui bodyDataMap:", bodyDataMap)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var tickerMap map[string]interface{}
	var ticker Ticker

	switch bodyDataMap[cur].(type) {
	case map[string]interface{}:
		tickerMap = bodyDataMap[cur].(map[string]interface{})
	default:
		return nil, errors.New(fmt.Sprintf("Type Convert Error ? \n %s", bodyDataMap))
	}

	//	fmt.Println(cur, " tickerMap:", tickerMap)
	date := tickerMap["updated"].(float64)
	ticker.Date = uint64(date)
	ticker.Last = tickerMap["last"].(float64)
	ticker.Buy = tickerMap["buy"].(float64)
	ticker.Sell = tickerMap["sell"].(float64)
	ticker.Low = tickerMap["low"].(float64)
	ticker.High = tickerMap["high"].(float64)
	ticker.Vol = tickerMap["vol"].(float64)

	return &ticker, nil
}

func (liqui *Liqui) GetDepth() string {
	panic("not implements")
}
func (liqui *Liqui) GetAccount() string {
	panic("not implements")
}
func (liqui *Liqui) LimitBuy() string {
	panic("not implements")
}
func (liqui *Liqui) LimitSell() string {
	panic("not implements")
}
func (liqui *Liqui) MarketBuy() string {
	panic("not implements")
}
func (liqui *Liqui) MarketSell() string {
	panic("not implements")
}
func (liqui *Liqui) CancelOrder() string {
	panic("not implements")
}
func (liqui *Liqui) GetOneOrder() string {
	panic("not implements")
}
func (liqui *Liqui) GetUnfinishOrders() string {
	panic("not implements")
}

func (bn *Binance) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implements")
}

//非个人，整个交易所的交易记录
func (bn *Binance) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implements")
}

func (bn *Binance) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implements")
}
