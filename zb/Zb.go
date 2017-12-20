package zb

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"net/http"
	"strconv"
	"strings"
)

const (
	EXCHANGE_NAME = "zb.com"
	MARKET_URL    = "http://api.zb.com/data/v1/"
	TICKER_API    = "ticker?market=%s"
	DEPTH_API     = "depth?currency=%s&size=%d"

	TRADE_URL                 = "https://trade.zb.com/api/"
	GET_ACCOUNT_API           = "getAccountInfo"
	GET_ORDER_API             = "getOrder"
	GET_UNFINISHED_ORDERS_API = "getUnfinishedOrdersIgnoreTradeType"
	CANCEL_ORDER_API          = "cancelOrder"
	PLACE_ORDER_API           = "order"
	WITHDRAW_API              = "withdraw"
	CANCELWITHDRAW_API        = "cancelWithdraw"
)

type ZB struct {
	httpClient *http.Client
	accessKey,
	secretKey string
}

func New(httpClient *http.Client, accessKey, secretKey string) *ZB {
	return &ZB{httpClient, accessKey, secretKey}
}

func (zb *ZB) GetExchangeName() string {
	return EXCHANGE_NAME
}

func (zb *ZB) GetTicker(currency CurrencyPair) (*Ticker, error) {
	//log.Println("ZB###")
	resp, err := HttpGet(zb.httpClient, MARKET_URL+fmt.Sprintf(TICKER_API, strings.ToLower(currency.ToSymbol("_"))))
	if err != nil {
		//log.Println("ZB err", err)
		return nil, err
	}
	//log.Println(resp)
	result, ok := resp["result"].bool
	if ok == true && result == false {
		//log.Println("err:", "{\"message\":\"服务端忙碌\",\"result\":false}")
		return nil, errors.New("server busy")
	}

	tickermap, ok := resp["ticker"].(map[string]interface{})
	if ok != true {
		return nil, errors.New("no ticker")
	}
	ticker := new(Ticker)
	ticker.Date = ToFloat64(resp["date"])
	ticker.Buy = ToFloat64(tickermap["buy"])
	ticker.Sell = ToFloat64(tickermap["sell"])
	ticker.Last = ToFloat64(tickermap["last"])
	ticker.High = ToFloat64(tickermap["high"])
	ticker.Low = ToFloat64(tickermap["low"])
	ticker.Vol = ToFloat64(tickermap["vol"])
	//log.Println("ZB####", ticker)
	return ticker, nil
}
func (zb *ZB) GetTickers(currency CurrencyPair) (*Ticker, error) {
	return zb.GetTicker(currency)
}

func (zb *ZB) GetTickerInBuf(currency CurrencyPair) (*Ticker, error) {
	return zb.GetTicker(currency)
}

func (zb *ZB) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (zb *ZB) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (zb *ZB) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (zb *ZB) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (zb *ZB) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("not implements")
}

func (zb *ZB) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}
func (zb *ZB) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implements")
}

func (zb *ZB) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implements")
}

func (zb *ZB) GetAccount() (*Account, error) {
	panic("not implements")
}

func (zb *ZB) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	return nil, nil
}

func (zb *ZB) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implements")
}

//非个人，整个交易所的交易记录
func (zb *ZB) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implements")
}
