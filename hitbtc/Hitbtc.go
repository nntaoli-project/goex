package hitbtc

import (
	"errors"
	"log"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nntaoli-project/goex"
)

const (
	EXCHANGE_NAME = "hitbtc.com"

	API_BASE_URL = "https://api.hitbtc.com/"
	API_V2       = "api/2/"
	SYMBOLS_URI  = "public/symbol"
	TICKER_URI   = "public/ticker/"
	BALANCE_URI  = "account/balance"
	ORDER_URI    = "order"
	DEPTH_URI    = "public/orderbook"
	TRADES_URI   = "public/trades"
	KLINE_URI    = "public/candles"
)

var (
	YCC     = goex.Currency{"YCC", "Yuan Chain New"}
	BTC     = goex.Currency{"BTC", "Bitcoin"}
	YCC_BTC = goex.CurrencyPair{YCC, BTC}
)

type Hitbtc struct {
	accessKey,
	secretKey string
	httpClient *http.Client
}

func New(client *http.Client, accessKey, secretKey string) *Hitbtc {
	return &Hitbtc{accessKey, secretKey, client}
}

func (hitbtc *Hitbtc) GetExchangeName() string {
	return EXCHANGE_NAME
}

// https://api.hitbtc.com/#symbols
/*
curl "https://api.hitbtc.com/api/2/public/symbol"
[
  {
    "id": "ETHBTC",
    "baseCurrency": "ETH",
    "quoteCurrency": "BTC",
    "quantityIncrement": "0.001",
    "tickSize": "0.000001",
    "takeLiquidityRate": "0.001",
    "provideLiquidityRate": "-0.0001",
    "feeCurrency": "BTC"
  }
]
*/
func (hitbtc *Hitbtc) GetSymbols() ([]goex.CurrencyPair, error) {
	resp := []map[string]interface{}{}
	err := hitbtc.doRequest("GET", SYMBOLS_URI, &resp)
	if err != nil {
		return nil, err
	}

	pairs := []goex.CurrencyPair{}
	for _, e := range resp {
		one := goex.CurrencyPair{
			CurrencyA: goex.Currency{e["baseCurrency"].(string), ""},
			CurrencyB: goex.Currency{e["quoteCurrency"].(string), ""},
		}
		pairs = append(pairs, one)
	}
	return pairs, nil
}

// https://api.hitbtc.com/#tickers

/*
curl "https://api.hitbtc.com/api/2/public/ticker"

[
  {
    "ask": "0.050043",
    "bid": "0.050042",
    "last": "0.050042",
    "open": "0.047800",
    "low": "0.047052",
    "high": "0.051679",
    "volume": "36456.720",
    "volumeQuote": "1782.625000",
    "timestamp": "2017-05-12T14:57:19.999Z",
    "symbol": "ETHBTC"
  }
]
*/
func (hitbtc *Hitbtc) GetTicker(currency goex.CurrencyPair) (*goex.Ticker, error) {
	curr := hitbtc.adaptCurrencyPair(currency).ToSymbol("")
	tickerUri := API_BASE_URL + API_V2 + TICKER_URI + curr
	bodyDataMap, err := goex.HttpGet(hitbtc.httpClient, tickerUri)
	if err != nil {
		return nil, err
	}

	if result, isok := bodyDataMap["error"].(map[string]interface{}); isok == true {
		return nil, errors.New(result["message"].(string) + ", " + result["description"].(string))
	}

	tickerMap := bodyDataMap
	var ticker goex.Ticker
	ticker.Pair = currency

	timestamp := time.Now().Unix()
	ticker.Date = uint64(timestamp)
	ticker.Last = goex.ToFloat64(tickerMap["last"])
	ticker.Buy = goex.ToFloat64(tickerMap["bid"])
	ticker.Sell = goex.ToFloat64(tickerMap["ask"])
	ticker.Low = goex.ToFloat64(tickerMap["low"])
	ticker.High = goex.ToFloat64(tickerMap["high"])
	ticker.Vol = goex.ToFloat64(tickerMap["volume"])

	return &ticker, nil
}

//
/*
curl -X PUT -u "ff20f250a7b3a414781d1abe11cd8cee:fb453577d11294359058a9ae13c94713" \
    "https://api.hitbtc.com/api/2/order/d8574207d9e3b16a4a5511753eeef175" \
    -d 'symbol=ETHBTC&side=sell&quantity=0.063&price=0.046016'

    {
        "id": 0,
        "clientOrderId": "d8574207d9e3b16a4a5511753eeef175",
        "symbol": "ETHBTC",
        "side": "sell",
        "status": "new",
        "type": "limit",
        "timeInForce": "GTC",
        "quantity": "0.063",
        "price": "0.046016",
        "cumQuantity": "0.000",
        "createdAt": "2017-05-15T17:01:05.092Z",
        "updatedAt": "2017-05-15T17:01:05.092Z"
    }
*/
func (hitbtc *Hitbtc) placeOrder(ty goex.TradeSide, amount, price string, currency goex.CurrencyPair) (*goex.Order, error) {
	postData := url.Values{}
	postData.Set("symbol", currency.ToSymbol(""))
	var side string
	var orderType string
	switch ty {
	case goex.BUY:
		side = "buy"
		orderType = "limit"
	case goex.BUY_MARKET:
		side = "buy"
		orderType = "market"
	case goex.SELL:
		side = "sell"
		orderType = "limit"
	case goex.SELL_MARKET:
		side = "sell"
		orderType = "market"
	default:
		panic(ty)
	}
	postData.Set("side", side)
	postData.Set("type", orderType)
	postData.Set("quantity", amount)
	if orderType == "limit" {
		postData.Set("price", price)
	}

	reqUrl := API_BASE_URL + API_V2 + ORDER_URI
	headers := make(map[string]string)
	headers["Content-type"] = "application/x-www-form-urlencoded"
	headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(hitbtc.accessKey+":"+hitbtc.secretKey))

	bytes, err := goex.HttpPostForm3(hitbtc.httpClient, reqUrl, postData.Encode(), headers)
	if err != nil {
		return nil, err
	}

	/*
		{
			"error": {
				"code": 20001,
			    "message": "Insufficient funds",
			    "description": "Check that the funds are sufficient, given commissions"
			}
		}
	*/
	resp := make(map[string]interface{})
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	if errObj, ok := resp["error"]; ok {
		log.Println(errObj)
		return nil, errors.New(errObj.(map[string]string)["message"])
	}

	return hitbtc.toOrder(resp), nil
}

func (hitbtc *Hitbtc) LimitBuy(amount, price string, currency goex.CurrencyPair) (*goex.Order, error) {
	return hitbtc.placeOrder(goex.BUY, amount, price, currency)
}

func (hitbtc *Hitbtc) LimitSell(amount, price string, currency goex.CurrencyPair) (*goex.Order, error) {
	return hitbtc.placeOrder(goex.SELL, amount, price, currency)
}

func (hitbtc *Hitbtc) MarketBuy(amount, price string, currency goex.CurrencyPair) (*goex.Order, error) {
	return hitbtc.placeOrder(goex.BUY_MARKET, amount, price, currency)
}

func (hitbtc *Hitbtc) MarketSell(amount, price string, currency goex.CurrencyPair) (*goex.Order, error) {
	return hitbtc.placeOrder(goex.SELL_MARKET, amount, price, currency)
}

func (hitbtc *Hitbtc) CancelOrder(orderId string, currency goex.CurrencyPair) (bool, error) {
	postData := url.Values{}
	reqUrl := API_BASE_URL + API_V2 + ORDER_URI + "/" + orderId
	headers := make(map[string]string)
	headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(hitbtc.accessKey+":"+hitbtc.secretKey))
	bytes, err := goex.HttpDeleteForm(hitbtc.httpClient, reqUrl, postData, headers)
	if err != nil {
		return false, err
	}

	var resp map[string]interface{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return false, err
	}

	if errObj, ok := resp["error"]; ok {
		log.Println(errObj)
		return false, errors.New(errObj.(map[string]string)["message"])
	}

	return true, nil
}

func (hitbtc *Hitbtc) GetOneOrder(orderId string, currency goex.CurrencyPair) (*goex.Order, error) {
	resp := make(map[string]interface{})
	err := hitbtc.doRequest("GET", ORDER_URI+"/"+orderId, &resp)
	if err != nil {
		return nil, err
	}

	if errObj, ok := resp["error"]; ok {
		return nil, errors.New(errObj.(map[string]string)["message"])
	}

	return hitbtc.toOrder(resp), nil
}

func (hitbtc *Hitbtc) GetUnfinishOrders(currency goex.CurrencyPair) ([]goex.Order, error) {
	params := url.Values{}
	params.Set("symbol", currency.ToSymbol(""))
	resp := []map[string]interface{}{}
	err := hitbtc.doRequest("GET", ORDER_URI+"?"+params.Encode(), &resp)
	if err != nil {
		return nil, err
	}

	// TODO error

	orders := []goex.Order{}
	for _, e := range resp {
		o := hitbtc.toOrder(e)
		if o.Status == goex.ORDER_UNFINISH || o.Status == goex.ORDER_PART_FINISH {
			orders = append(orders, *o)
		}
	}
	return orders, nil
}

// TODO
// https://api.hitbtc.com/#orders-history
func (hitbtc *Hitbtc) GetOrderHistorys(currency goex.CurrencyPair, currentPage, pageSize int) ([]goex.Order, error) {
	params := url.Values{}
	params.Set("symbol", currency.ToSymbol(""))
	resp := []map[string]interface{}{}
	err := hitbtc.doRequest("GET", ORDER_URI+"?"+params.Encode(), &resp)
	if err != nil {
		return nil, err
	}

	// TODO error

	orders := []goex.Order{}
	for _, e := range resp {
		o := hitbtc.toOrder(e)
		orders = append(orders, *o)
	}
	return orders, nil
}

// https://api.hitbtc.com/#account-balance
func (hitbtc *Hitbtc) GetAccount() (*goex.Account, error) {
	var ret []interface{}
	err := hitbtc.doRequest("GET", BALANCE_URI, &ret)
	if err != nil {
		return nil, err
	}

	acc := new(goex.Account)
	acc.SubAccounts = make(map[goex.Currency]goex.SubAccount, 1)

	for _, v := range ret {
		vv := v.(map[string]interface{})
		currency := goex.NewCurrency(vv["currency"].(string), "")
		acc.SubAccounts[currency] = goex.SubAccount{
			Currency:     currency,
			Amount:       goex.ToFloat64(vv["available"]),
			ForzenAmount: goex.ToFloat64(vv["reserved"])}
	}

	return acc, nil
}

// https://api.hitbtc.com/#orderbook
/*
{
  "ask": [
    {
      "price": "0.046002",
      "size": "0.088"
    },
    {
      "price": "0.046800",
      "size": "0.200"
    }
  ],
  "bid": [
    {
      "price": "0.046001",
      "size": "0.005"
    },
    {
      "price": "0.046000",
      "size": "0.200"
    }
  ]
}
*/
func (hitbtc *Hitbtc) GetDepth(size int, currency goex.CurrencyPair) (*goex.Depth, error) {
	params := url.Values{}
	params.Set("limit", fmt.Sprintf("%v", size))
	resp := map[string]interface{}{}
	err := hitbtc.doRequest("GET", DEPTH_URI+"/"+currency.ToSymbol("")+"?"+params.Encode(), &resp)
	if err != nil {
		return nil, err
	}

	if errObj, ok := resp["error"]; ok {
		return nil, errors.New(errObj.(map[string]string)["message"])
	}

	askList := []goex.DepthRecord{}

	for _, ee := range resp["ask"].([]interface{}) {
		e := ee.(map[string]interface{})
		one := goex.DepthRecord{
			Price:  goex.ToFloat64(e["price"]),
			Amount: goex.ToFloat64(e["size"]),
		}
		askList = append(askList, one)
	}

	bidList := []goex.DepthRecord{}
	for _, ee := range resp["bid"].([]interface{}) {
		e := ee.(map[string]interface{})
		one := goex.DepthRecord{
			Price:  goex.ToFloat64(e["price"]),
			Amount: goex.ToFloat64(e["size"]),
		}
		bidList = append(bidList, one)
	}

	return &goex.Depth{AskList: askList, BidList: bidList}, nil
}

func (hitbtc *Hitbtc) GetKlineRecords(currency goex.CurrencyPair, period, size, since int) ([]goex.Kline, error) {
	panic("not implement")
}

// https://api.hitbtc.com/#candles

/*
curl "https://api.hitbtc.com/api/2/public/candles/ETHBTC?period=M30"
[
  {
    "timestamp": "2017-10-20T20:00:00.000Z",
    "open": "0.050459",
    "close": "0.050087",
    "min": "0.050000",
    "max": "0.050511",
    "volume": "1326.628",
    "volumeQuote": "66.555987736"
  },
  {
    "timestamp": "2017-10-20T20:30:00.000Z",
    "open": "0.050108",
    "close": "0.050139",
    "min": "0.050068",
    "max": "0.050223",
    "volume": "87.515",
    "volumeQuote": "4.386062831"
  }
]
*/
func (hitbtc *Hitbtc) GetKline(currencyPair goex.CurrencyPair, period string, size, since int64) ([]goex.Kline, error) {
	switch period {
	case "M1", "M3", "M5", "M15", "M30", "H1", "H4", "D1", "D7", "1M":
	default:
		return nil, errors.New("Invalid period")
	}
	if size < 0 {
		return nil, errors.New("Invalid size")
	}

	params := url.Values{}
	params.Set("period", period)
	if size > 0 {
		params.Set("limit", fmt.Sprintf("%v", size))
	}
	resp := []map[string]interface{}{}
	err := hitbtc.doRequest("GET", KLINE_URI+"/"+currencyPair.ToSymbol("")+"?"+params.Encode(), &resp)
	if err != nil {
		return nil, err
	}

	klines := []goex.Kline{}
	for _, e := range resp {
		one := goex.Kline{
			Timestamp: parseTime(e["timestamp"].(string)),
			Open:      goex.ToFloat64(e["open"]),
			Close:     goex.ToFloat64(e["close"]),
			High:      goex.ToFloat64(e["high"]),
			Low:       goex.ToFloat64(e["low"]),
			Vol:       goex.ToFloat64(e["volume"]), // base currency, eg: ETH for pair ETHBTC
		}
		klines = append(klines, one)
	}
	return klines, nil
}

// https://api.hitbtc.com/#trades

/*
curl "https://api.hitbtc.com/api/2/public/trades/ETHBTC?from=2018-05-22T07:22:00&limit=2"
[
   {
      "id" : 297604734,
      "timestamp" : "2018-05-22T07:23:06.556Z",
      "quantity" : "6.551",
      "side" : "sell",
      "price" : "0.083421"
   },
   {
      "side" : "sell",
      "price" : "0.083401",
      "quantity" : "0.021",
      "timestamp" : "2018-05-22T07:23:05.908Z",
      "id" : 297604724
   },
]
*/
func (hitbtc *Hitbtc) GetTrades(currencyPair goex.CurrencyPair, since int64) ([]goex.Trade, error) {
	params := url.Values{}
	timestamp := time.Unix(since, 0).Format("2006-01-02T15:04:05")
	params.Set("from", timestamp)
	resp := []map[string]interface{}{}
	err := hitbtc.doRequest("GET", TRADES_URI+"/"+currencyPair.ToSymbol("")+"?"+params.Encode(), &resp)
	if err != nil {
		return nil, err
	}

	trades := []goex.Trade{}
	for _, e := range resp {
		one := goex.Trade{
			Tid:    int64(goex.ToUint64(e["id"])),
			Type:   goex.AdaptTradeSide(e["side"].(string)),
			Amount: goex.ToFloat64(e["quantity"]),
			Price:  goex.ToFloat64(e["price"]),
			Date:   parseTime(e["timestamp"].(string)),
		}
		trades = append(trades, one)
	}
	return trades, nil
}

func (hitbtc *Hitbtc) doRequest(reqMethod, uri string, ret interface{}) error {
	url := API_BASE_URL + API_V2 + uri
	req, _ := http.NewRequest(reqMethod, url, strings.NewReader(""))
	req.SetBasicAuth(hitbtc.accessKey, hitbtc.secretKey)
	resp, err := hitbtc.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("HttpStatusCode:%d ,Desc:%s", resp.StatusCode, string(bodyData)))
	}

	err = json.Unmarshal(bodyData, ret)

	if err != nil {
		return err
	}

	return nil
}

/*

// https://api.hitbtc.com/#order-model
{
	"id": 0,
	"clientOrderId": "d8574207d9e3b16a4a5511753eeef175",
	"symbol": "ETHBTC",
	"side": "sell",
	"status": "new",
	"type": "limit",
	"timeInForce": "GTC",
	"quantity": "0.063",
	"price": "0.046016",
	"cumQuantity": "0.000",
	"createdAt": "2017-05-15T17:01:05.092Z",
	"updatedAt": "2017-05-15T17:01:05.092Z"
}
*/
func (hitbtc *Hitbtc) toOrder(resp map[string]interface{}) *goex.Order {
	return &goex.Order{
		Price:      goex.ToFloat64(resp["price"]),
		Amount:     goex.ToFloat64(resp["quantity"]),
		DealAmount: goex.ToFloat64(resp["cumQuantity"]),
		OrderID2:   resp["clientOrderId"].(string),
		OrderID:    goex.ToInt(resp["id"]),
		OrderTime:  int(parseTime(resp["createdAt"].(string))),
		Status:     parseStatus(resp["status"].(string)),
		Currency:   hitbtc.adaptSymbolToCurrencyPair(resp["symbol"].(string)),
		Side:       parseSide(resp["side"].(string), resp["type"].(string)),
	}
}

func parseStatus(s string) goex.TradeStatus {
	var status goex.TradeStatus
	switch s {
	case "new", "suspended":
		status = goex.ORDER_UNFINISH
	case "partiallyFilled":
		status = goex.ORDER_PART_FINISH
	case "filled":
		status = goex.ORDER_FINISH
	case "canceled":
		status = goex.ORDER_CANCEL
	case "expired":
		// TODO
		status = goex.ORDER_REJECT
	default:
		panic(s)
	}
	return status
}

func (hitbtc *Hitbtc) adaptCurrencyPair(pair goex.CurrencyPair) goex.CurrencyPair {
	return pair.AdaptUsdtToUsd().AdaptBccToBch()
}

func (hitbtc *Hitbtc) adaptSymbolToCurrencyPair(pair string) goex.CurrencyPair {
	pair = strings.ToUpper(pair)
	if strings.HasSuffix(pair, "BTC") {
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_%s", strings.TrimSuffix(pair, "BTC"), "BTC"))
	} else if strings.HasSuffix(pair, "TUSD") {
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_%s", strings.TrimSuffix(pair, "TUSD"), "TUSD"))
	} else if strings.HasSuffix(pair, "GUSD") {
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_%s", strings.TrimSuffix(pair, "GUSD"), "GUSD"))
	} else if strings.HasSuffix(pair, "USDC") {
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_%s", strings.TrimSuffix(pair, "USDC"), "USDC"))
	} else if strings.HasSuffix(pair, "USD") {
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_%s", strings.TrimSuffix(pair, "USD"), "USD"))
	} else if strings.HasSuffix(pair, "ETH") {
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_%s", strings.TrimSuffix(pair, "ETH"), "ETH"))
	} else if strings.HasSuffix(pair, "KRWB") {
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_%s", strings.TrimSuffix(pair, "KRWB"), "KRWB"))
	} else if strings.HasSuffix(pair, "PAX") {
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_%s", strings.TrimSuffix(pair, "PAX"), "PAX"))
	} else if strings.HasSuffix(pair, "DAI") {
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_%s", strings.TrimSuffix(pair, "DAI"), "DAI"))
	} else if strings.HasSuffix(pair, "EOS") {
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_%s", strings.TrimSuffix(pair, "EOS"), "EOS"))
	} else if strings.HasSuffix(pair, "EURS") {
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_%s", strings.TrimSuffix(pair, "EURS"), "EURS"))
	}
	return goex.UNKNOWN_PAIR
}

func parseTime(timeStr string) int64 {
	t, _ := time.Parse(time.RFC3339, timeStr) // UTC
	return t.Unix()
}

func parseSide(side, oType string) goex.TradeSide {
	if side == "buy" && oType == "limit" {
		return goex.BUY
	} else if side == "sell" && oType == "limit" {
		return goex.SELL
	} else if side == "buy" && oType == "market" {
		return goex.BUY_MARKET
	} else if side == "sell" && oType == "market" {
		return goex.SELL_MARKET
	} else {
		panic("Invalid TradeSide:" + side + "&" + oType)
	}
}
