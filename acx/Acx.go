package acx

import (
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"net/http"
	"strings"
)

const (
	EXCHANGE_NAME = "acx.io"

	API_BASE_URL = "https://acx.io/"
	API_V1       = API_BASE_URL + "/api/v2/"

	TICKER_URI = "/tickers/%s.json"
	//DEPTH_URI              = "depth.php?c=%s&mk_type=%s"
	//ACCOUNT_URI            = "getMyBalance.php"
	//TRADE_URI              = "trades.php?c=%s&mk_type=%s"
	//CANCEL_URI             = "cancelOrder.php"
	//ORDERS_INFO            = "getMyTradeList.php"
	//UNFINISHED_ORDERS_INFO = "getOrderList.php"

)

type Acx struct {
	accessKey,
	secretKey string
	httpClient *http.Client
}

func New(client *http.Client, api_key, secret_key string) *Acx {
	return &Acx{api_key, secret_key, client}
}

func (acx *Acx) GetExchangeName() string {
	return EXCHANGE_NAME
}

func (acx *Acx) GetTicker(currency CurrencyPair) (*Ticker, error) {
	tickerUri := API_V1 + fmt.Sprintf(TICKER_URI, strings.ToLower(currency.ToSymbol("")))
	bodyDataMap, err := HttpGet(acx.httpClient, tickerUri)

	if err != nil {
		return nil, err
	}

	//log.Println("acx uri:", tickerUri)
	//log.Println("acx bodyDataMap:", currency, bodyDataMap)

	tickerMap := bodyDataMap["ticker"].(map[string]interface{})
	var ticker Ticker

	ticker.Date = uint64(bodyDataMap["at"].(float64))
	ticker.Last = ToFloat64(tickerMap["last"])
	ticker.Buy = ToFloat64(tickerMap["buy"])
	ticker.Sell = ToFloat64(tickerMap["sell"])
	ticker.Low = ToFloat64(tickerMap["low"])
	ticker.High = ToFloat64(tickerMap["high"])
	ticker.Vol = ToFloat64(tickerMap["vol"])
	return &ticker, nil
}
func (acx *Acx) GetTickers(currency CurrencyPair) (*Ticker, error) {
	return acx.GetTicker(currency)
}
func (acx *Acx) GetTickerInBuf(currency CurrencyPair) (*Ticker, error) {
	return acx.GetTicker(currency)
}

func (acx *Acx) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (acx *Acx) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (acx *Acx) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (acx *Acx) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (acx *Acx) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("not implements")
}

func (acx *Acx) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}
func (acx *Acx) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implements")
}

func (acx *Acx) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implements")
}

func (acx *Acx) GetAccount() (*Account, error) {
	panic("not implements")
}

func (acx *Acx) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	return nil, nil
}

func (acx *Acx) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implements")
}

//非个人，整个交易所的交易记录
func (acx *Acx) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implements")
}
