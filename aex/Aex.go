package aex

import (
	"errors"
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	//"log"
	"net/http"
	"time"
)

//因服务器有防CC攻击策略，每60秒内调用次数不可超过120次，超过部分将被防火墙拦截。
const (
	EXCHANGE_NAME = "aex.com"

	API_BASE_URL = "https://api.aex.com/"
	API_V1       = API_BASE_URL

	TICKER_URI             = "ticker.php?c=%s&mk_type=%s"
	DEPTH_URI              = "depth.php?c=%s&mk_type=%s"
	ACCOUNT_URI            = "getMyBalance.php"
	TRADE_URI              = "trades.php?c=%s&mk_type=%s"
	CANCEL_URI             = "cancelOrder.php"
	ORDERS_INFO            = "getMyTradeList.php"
	UNFINISHED_ORDERS_INFO = "getOrderList.php"
)

type Aex struct {
	accessKey,
	secretKey,
	accountId string
	httpClient *http.Client
}

func New(client *http.Client, accessKey, secretKey, accountId string) *Aex {
	return &Aex{accessKey, secretKey, accountId, client}
}

func (aex *Aex) GetExchangeName() string {
	return EXCHANGE_NAME
}

func (aex *Aex) GetTicker(currency CurrencyPair) (*Ticker, error) {
	cur := currency.CurrencyA.String()
	money := currency.CurrencyB.String()
	if cur == "UNKNOWN" {
		//log.Println("Unsupport The CurrencyPair")
		return nil, errors.New("Unsupport The CurrencyPair")
	}
	tickerUri := API_V1 + fmt.Sprintf(TICKER_URI, cur, money)
	timestamp := time.Now().Unix()

	bodyDataMap, err := HttpGet(aex.httpClient, tickerUri)

	if err != nil {
		//log.Println(err)
		return nil, err
	}
	//	log.Println("Aex bodyDataMap:", bodyDataMap)
	var tickerMap map[string]interface{}
	var ticker Ticker

	switch bodyDataMap["ticker"].(type) {
	case map[string]interface{}:
		tickerMap = bodyDataMap["ticker"].(map[string]interface{})
	default:
		return nil, errors.New(fmt.Sprintf("Type Convert Error ? \n %s", bodyDataMap))
	}

	ticker.Date = uint64(timestamp)
	ticker.Last = tickerMap["last"].(float64)
	ticker.Buy = tickerMap["buy"].(float64)
	ticker.Sell = tickerMap["sell"].(float64)
	ticker.Low = tickerMap["low"].(float64)
	ticker.High = tickerMap["high"].(float64)
	ticker.Vol = tickerMap["vol"].(float64)
	//log.Println("Aex", currency, "ticker:", ticker)

	return &ticker, nil
}

func (aex *Aex) GetTickers(currency CurrencyPair) (*Ticker, error) {
	return aex.GetTicker(currency)
}

func (aex *Aex) GetTickerInBuf(currency CurrencyPair) (*Ticker, error) {
	return aex.GetTicker(currency)
}

func (aex *Aex) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	panic("not implement")
}

func (aex *Aex) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (aex *Aex) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (aex *Aex) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (aex *Aex) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (aex *Aex) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("not implements")
}

func (aex *Aex) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}
func (aex *Aex) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implements")
}

func (aex *Aex) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implements")
}

func (aex *Aex) GetAccount() (*Account, error) {
	panic("not implements")
}

func (aex *Aex) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implements")
}

//非个人，整个交易所的交易记录
func (aex *Aex) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implements")
}
