package binance

import (
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"net/http"
	"strconv"

	"errors"
	"log"
	"net/url"
	"strings"
	"time"
	// "encoding/json"
	"encoding/json"
)

const (
	EXCHANGE_NAME = "binance.com"

	API_BASE_URL = "https://www.binance.com/"
	API_V1       = API_BASE_URL + "api/v1/"
	API_V3       = API_BASE_URL + "api/v3/"

	TICKER_URI      = "ticker/24hr?symbol=%s"
	DEPTH_URI       = "depth?symbol=%s&limit=%d"
	ACCOUNT_URI     = "account"
	PLACE_ORDER_API = "order"
	//TRADE_URI              = "trades.php?c=%s&mk_type=%s"
	//CANCEL_URI             = "cancelOrder.php"
	//ORDERS_INFO            = "getMyTradeList.php"
	//UNFINISHED_ORDERS_INFO = "getOrderList.php"

)

type Binance struct {
	accessKey,
	secretKey string
	httpClient *http.Client
}

func (ba *Binance) buildPostForm(postForm *url.Values) error {
	postForm.Set("type", "LIMIT")
	postForm.Set("timeInForce", "GTC")
	postForm.Set("recvWindow", "6000000")
	tonce := strconv.FormatInt(time.Now().UnixNano(), 10)[0:13]
	postForm.Set("timestamp", tonce)

	payload := postForm.Encode()
	log.Println("payload", payload)
	sign, _ := GetParamHmacSHA256Sign(ba.secretKey, payload)
	log.Println("sign", sign)

	postForm.Set("signature", sign)
	return nil
	//return payload + "&signature="+sign
}

func New(client *http.Client, api_key, secret_key string) *Binance {
	return &Binance{api_key, secret_key, client}
}

func (ba *Binance) GetExchangeName() string {
	return EXCHANGE_NAME
}

func (ba *Binance) GetTicker(currency CurrencyPair) (*Ticker, error) {
	tickerUri := API_V1 + fmt.Sprintf(TICKER_URI, currency.ToSymbol(""))
	bodyDataMap, err := HttpGet(ba.httpClient, tickerUri)

	if err != nil {
		return nil, err
	}
	fmt.Println("binance bodyDataMap:", currency, bodyDataMap)
	var tickerMap map[string]interface{} = bodyDataMap
	var ticker Ticker

	ticker.Date = uint64(tickerMap["closeTime"].(float64))
	ticker.Last, _ = strconv.ParseFloat(tickerMap["lastPrice"].(string), 10)
	ticker.Buy, _ = strconv.ParseFloat(tickerMap["bidPrice"].(string), 10)
	ticker.Sell, _ = strconv.ParseFloat(tickerMap["askPrice"].(string), 10)
	ticker.Low, _ = strconv.ParseFloat(tickerMap["lowPrice"].(string), 10)
	ticker.High, _ = strconv.ParseFloat(tickerMap["highPrice"].(string), 10)
	ticker.Vol, _ = strconv.ParseFloat(tickerMap["volume"].(string), 10)
	return &ticker, nil
}

func (ba *Binance) GetDepth(size int, currencyPair CurrencyPair) (*Depth, error) {
	if size > 100 {
		size = 100
	} else if size < 0 {
		size = 0
	}

	apiUrl := fmt.Sprintf(API_V1+DEPTH_URI, currencyPair.ToSymbol(""), size)
	resp, err := HttpGet(ba.httpClient, apiUrl)
	log.Println(err, resp)
	if err != nil {
		return nil, err
	}

	if _, isok := resp["code"]; isok {
		return nil, errors.New(resp["msg"].(string))
	}

	bids := resp["bids"].([]interface{})
	asks := resp["asks"].([]interface{})

	log.Println(bids)
	log.Println(asks)

	depth := new(Depth)

	for _, bid := range bids {
		_bid := bid.([]interface{})
		amount := ToFloat64(_bid[1])
		price := ToFloat64(_bid[0])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.BidList = append(depth.BidList, dr)
	}

	for _, ask := range asks {
		_ask := ask.([]interface{})
		amount := ToFloat64(_ask[1])
		price := ToFloat64(_ask[0])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.AskList = append(depth.AskList, dr)
	}

	return depth, nil
}

func (ba *Binance) GetKlineRecords(currencyPair CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录

func (ba *Binance) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}

func (ba *Binance) placeOrder(amount, price string, pair CurrencyPair, orderType, orderSide string) (string, error) {
	path := API_V3 + PLACE_ORDER_API + "/test"
	params := url.Values{}
	params.Set("symbol", pair.ToSymbol(""))
	params.Set("side", orderSide)
	params.Set("type", orderType)

	params.Set("quantity", amount)

	switch orderType {
	case "LIMIT":
		params.Set("price", price)
	}

	ba.buildPostForm(&params)

	resp, err := HttpPostForm2(ba.httpClient, path, params,
		map[string]string{"Content-Type": "application/x-www-form-urlencoded", "X-MBX-APIKEY": ba.accessKey})
	log.Println("resp:", string(resp), "err:", err)

	if err != nil {
		return "", err
	}
	//respmap := make(map[string]interface{})
	//err = json.Unmarshal(resp, &respmap)
	//if err != nil {
	//	return "", err
	//}
	//
	//if respmap["status"].(string) != "ok" {
	//	return "", errors.New(respmap["err-code"].(string))
	//}
	return "", nil
	//return respmap["data"].(string), nil
}

func (ba *Binance) GetAccount() (*Account, error) {
	return nil, nil
}

func (ba *Binance) LimitBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {

	ba.placeOrder(amount, price, currencyPair, "LIMIT", "BUY")

	return nil, nil
}

func (ba *Binance) LimitSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return nil, nil
}

func (ba *Binance) MarketBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	panic("not implement.")
}

func (ba *Binance) MarketSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	panic("not implement.")
}

func (ba *Binance) CancelOrder(orderId string, currencyPair CurrencyPair) (bool, error) {
	return false, nil
}

func (ba *Binance) GetOneOrder(orderId string, currencyPair CurrencyPair) (*Order, error) {
	return nil, nil
}

func (ba *Binance) GetUnfinishOrders(currencyPair CurrencyPair) ([]Order, error) {
	return nil, nil
}

func (ba *Binance) GetOrderHistorys(currencyPair CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	return nil, nil
}

func (ba *Binance) currencyPairToSymbol(currencyPair CurrencyPair) string {
	return strings.ToLower(currencyPair.ToSymbol(""))
}
func (ba *Binance) toJson(params url.Values) string {
	parammap := make(map[string]string)
	for k, v := range params {
		parammap[k] = v[0]
	}
	jsonData, _ := json.Marshal(parammap)
	return string(jsonData)
}
