package bitfinex

import (
	"errors"
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"log"
	"net/http"
)

type Bitfinex struct {
	httpClient *http.Client
	accessKey,
	secretKey string
}

const (
	EXCHANGE_NAME = "bitfinex.com"

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
	return EXCHANGE_NAME
}

func (bfx *Bitfinex) GetTicker(currency CurrencyPair) (*Ticker, error) {
	//pubticker
	cur := currency.DeleteUnderLineString()
	if cur == "nil" {
		log.Println("Unsupport The CurrencyPair")
		return nil, errors.New("Unsupport The CurrencyPair")
	}

	apiUrl := fmt.Sprintf("%s/pubticker/%s", BASE_URL, cur)
	resp, err := HttpGet(bfx.httpClient, apiUrl)
	if err != nil {
		return nil, err
	}

	if resp["error"] != nil {
		return nil, errors.New(resp["error"].(string))
	}

	//fmt.Println(resp)
	ticker := new(Ticker)
	ticker.Last = ToFloat64(resp["last_price"])
	ticker.Vol = ToFloat64(resp["volume"])
	ticker.High = ToFloat64(resp["high"])
	ticker.Low = ToFloat64(resp["low"])
	ticker.Sell = ToFloat64(resp["ask"])
	ticker.Buy = ToFloat64(resp["bid"])
	ticker.Date = uint64(ToFloat64(resp["timestamp"]))
	//dateStr := resp["timestamp"].(string)
	//dataMeta := strings.Split(dateStr, ".")
	//ticker.Date, _ = strconv.ParseUint(dataMeta[0], 10, 64)
	return ticker, nil
}

func (bfx *Bitfinex) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	apiUrl := fmt.Sprintf("%s/book/%s?limit_bids=%d&limit_asks=%d", BASE_URL, currency.DeleteUnderLineString(), size, size)
	resp, err := HttpGet(bfx.httpClient, apiUrl)
	if err != nil {
		return nil, err
	}

	bids := resp["bids"].([]interface{})
	asks := resp["asks"].([]interface{})

	depth := new(Depth)

	for _, bid := range bids {
		_bid := bid.(map[string]interface{})
		amount := ToFloat64(_bid["amount"])
		price := ToFloat64(_bid["price"])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.BidList = append(depth.BidList, dr)
	}

	for _, ask := range asks {
		_ask := ask.(map[string]interface{})
		amount := ToFloat64(_ask["amount"])
		price := ToFloat64(_ask["price"])
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
