package bitfinex

import (
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"net/http"
	"strconv"
	"strings"
)

type Bitfinex struct {
	httpClient *http.Client
	accessKey,
	secretKey string
}

const (
	BASE_URL = "https://api.bitfinex.com/v1"
)

var CURRENCYPAIR_TO_SYMBOL = map[CurrencyPair]string{
	BTC_USD: "btcusd",
	LTC_USD: "ltcusd",
	LTC_BTC: "ltcbtc",
	ETH_BTC: "ethbtc",
	ETC_BTC: "etcbtc",
	ETC_USD: "etcusd",
}

func New(client *http.Client, accessKey, secretKey string) *Bitfinex {
	return &Bitfinex{client, accessKey, secretKey}
}

func (bfx *Bitfinex) GetExchangeName() string {
	return "bitfinex.com"
}

func (bfx *Bitfinex) GetTicker(currency CurrencyPair) (*Ticker, error) {
	//pubticker
	apiUrl := fmt.Sprintf("%s/pubticker/%s", BASE_URL, CURRENCYPAIR_TO_SYMBOL[currency])
	resp, err := HttpGet(bfx.httpClient, apiUrl)
	if err != nil {
		return nil, err
	}
	//fmt.Println(resp)
	ticker := new(Ticker)
	ticker.Last, _ = strconv.ParseFloat(resp["last_price"].(string), 64)
	ticker.Vol, _ = strconv.ParseFloat(resp["volume"].(string), 64)
	ticker.High, _ = strconv.ParseFloat(resp["high"].(string), 64)
	ticker.Low, _ = strconv.ParseFloat(resp["low"].(string), 64)
	ticker.Sell, _ = strconv.ParseFloat(resp["ask"].(string), 64)
	ticker.Buy, _ = strconv.ParseFloat(resp["bid"].(string), 64)
	dateStr := resp["timestamp"].(string)
	dataMeta := strings.Split(dateStr, ".")
	ticker.Date, _ = strconv.ParseUint(dataMeta[0], 10, 64)
	return ticker, nil
}

func (bfx *Bitfinex) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	apiUrl := fmt.Sprintf("%s/book/%s?limit_bids=%d&limit_asks=%d", BASE_URL, CURRENCYPAIR_TO_SYMBOL[currency], size, size)
	resp, err := HttpGet(bfx.httpClient, apiUrl)
	if err != nil {
		return nil, err
	}

	bids := resp["bids"].([]interface{})
	asks := resp["asks"].([]interface{})

	depth := new(Depth)

	for _, bid := range bids {
		_bid := bid.(map[string]interface{})
		amount, _ := strconv.ParseFloat(_bid["amount"].(string), 64)
		price, _ := strconv.ParseFloat(_bid["price"].(string), 64)
		dr := DepthRecord{Amount: amount, Price: price}
		depth.BidList = append(depth.BidList, dr)
	}

	for _, ask := range asks {
		_ask := ask.(map[string]interface{})
		amount, _ := strconv.ParseFloat(_ask["amount"].(string), 64)
		price, _ := strconv.ParseFloat(_ask["price"].(string), 64)
		dr := DepthRecord{Amount: amount, Price: price}
		depth.AskList = append(depth.AskList, dr)
	}

	return depth, nil
}

func (bfx *Bitfinex) GetKlineRecords(currency CurrencyPair, period string, size, since int) ([]Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录

func (bfx *Bitfinex) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}

func (bfx *Bitfinex) GetAccount() (*Account, error) {
	return nil, nil
}

func (bfx *Bitfinex) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return nil, nil
}

func (bfx *Bitfinex) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return nil, nil
}

func (bfx *Bitfinex) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement.")
}

func (bfx *Bitfinex) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement.")
}

func (bfx *Bitfinex) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	return false, nil
}

func (bfx *Bitfinex) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	return nil, nil
}

func (bfx *Bitfinex) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	return nil, nil
}

func (bfx *Bitfinex) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	return nil, nil
}
