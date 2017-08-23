package yunbi

import (
	"net/http"
	. "github.com/nntaoli-project/GoEx"
	"fmt"
	"log"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"time"
	"net/url"
	"errors"
	"strings"
)

const _EXCHANGE_NAME = "yunbi.com"

var (
	API_URL = "https://yunbi.com"
	API_URI_PREFIX = "/api/v2/"
	TICKER_URL = "tickers/%s.json"
	DEPTH_URL = "depth.json?market=%s&limit=%d"
	USER_INFO_URL = "members/me.json"
	GET_ORDER_API = "order.json"
	DELETE_ORDER_API = "order/delete.json"
	PLACE_ORDER_API = "orders.json"
)

type YunBi struct {
	accessKey,
	secretKey string
	client    *http.Client
}

func New(client *http.Client, apikey, secretkey string) *YunBi {
	return &YunBi{apikey, secretkey, client}
}

func (yunbi *YunBi) GetExchangeName() string {
	return _EXCHANGE_NAME
}

type _TickerResponse struct {
	At     uint64 `json:"at"`
	Ticker *struct {
		Buy  float64 `json:"buy,string"`
		Sell float64 `json:"sell,string"`
		Low  float64 `json:"low,string"`
		High float64 `json:"high,string"`
		Last float64 `json:"last,string"`
		Vol  float64 `json:"vol,string"`
	}`json:"ticker"`
}

func (yunbi *YunBi)GetTicker(currency CurrencyPair) (*Ticker, error) {
	urlStr := fmt.Sprintf(API_URL + API_URI_PREFIX + TICKER_URL, convertCurrencyPair(currency))
	resp, err := yunbi.client.Get(urlStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	
	//println(string(respData))
	tickerResp := new(_TickerResponse)
	err = json.Unmarshal(respData, tickerResp)
	if err != nil {
		return nil , err
	}
	
	ticker := new(Ticker)
	ticker.Date = tickerResp.At
	ticker.Buy = tickerResp.Ticker.Buy
	ticker.Sell = tickerResp.Ticker.Sell
	ticker.Last = tickerResp.Ticker.Last
	ticker.Low = tickerResp.Ticker.Low
	ticker.High = tickerResp.Ticker.High
	ticker.Vol = tickerResp.Ticker.Vol
	
	return ticker, nil
}

func (yunbi *YunBi)GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	urlStr := fmt.Sprintf(API_URL + API_URI_PREFIX + DEPTH_URL, convertCurrencyPair(currency), size)
	respMap, err := HttpGet(yunbi.client, urlStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	
	depth := new(Depth)
	
	for _, v := range respMap["asks"].([]interface{}) {
		var dr DepthRecord;
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price, _ = strconv.ParseFloat(vv.(string), 64);
			case 1:
				dr.Amount, _ = strconv.ParseFloat(vv.(string), 64);
			}
		}
		depth.AskList = append(depth.AskList, dr);
	}
	
	for _, v := range respMap["bids"].([]interface{}) {
		var dr DepthRecord;
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price, _ = strconv.ParseFloat(vv.(string), 64);
			case 1:
				dr.Amount, _ = strconv.ParseFloat(vv.(string), 64);
			}
		}
		depth.BidList = append(depth.BidList, dr);
	}
	
	return depth, nil
}

func (yunbi *YunBi)GetAccount() (*Account, error) {
	urlStr := API_URL + API_URI_PREFIX + USER_INFO_URL;
	postParams := url.Values{}
	yunbi.buildPostForm("GET", API_URI_PREFIX + USER_INFO_URL, &postParams)
	
	resp, err := HttpGet(yunbi.client, urlStr + "?" + postParams.Encode());
	if err != nil {
		log.Println(err)
		return nil, err
	}
	
	//log.Println(resp)
	if resp["error"] != nil{
		errmap := resp["error"].(map[string]interface{})
		errcode := errmap["code"].(float64)
		errmsg := errmap["message"].(string)
		return nil , errors.New(fmt.Sprintf("%.0f:%s" , errcode , errmsg))
	}

	acc := new(Account)
	acc.SubAccounts = make(map[Currency]SubAccount)
	
	accountsMap := resp["accounts"].([]interface{})
	for _, v := range accountsMap {
		vv := v.(map[string]interface{})
		subAcc := SubAccount{}
		subAcc.Amount, _ = strconv.ParseFloat(vv["balance"].(string), 64)
		subAcc.ForzenAmount, _ = strconv.ParseFloat(vv["locked"].(string), 64)
		
		var
		(
			currency Currency
			skip bool = false
		)
		
		switch vv["currency"] {
		case "btc":
			currency = BTC
		case "cny":
			currency = CNY
		case "etc":
			currency = ETC
		case "eth":
			currency = ETH
		case "zec":
			currency = ZEC
		case "sc":
			currency = SC
		case "bts":
			currency = BTS
		case "eos":
			currency = EOS
		default:
			skip = true
		}
		
		if !skip {
			subAcc.Currency = currency
			acc.SubAccounts[currency] = subAcc
		}
	}
	
	return acc, nil
}

func (yunbi *YunBi)LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return yunbi.placeOrder("buy", amount, price, currency)
}

func (yunbi *YunBi) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return yunbi.placeOrder("sell", amount, price, currency)
}

func (yunbi *YunBi) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return nil, nil
}

func (yunbi *YunBi) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return nil, nil
}

func (yunbi *YunBi)placeOrder(side, amount, price string, currencyPair CurrencyPair) (*Order, error) {
	params := url.Values{}
	params.Set("market", convertCurrencyPair(currencyPair))
	params.Set("side", side)
	params.Set("price", price)
	params.Set("volume", amount)
	
	yunbi.buildPostForm("POST", API_URI_PREFIX + PLACE_ORDER_API, &params)
	
	resp, err := HttpPostForm(yunbi.client, API_URL + API_URI_PREFIX + PLACE_ORDER_API, params)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	
	//println(string(resp))
	
	respMap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		log.Println(err, string(resp))
		return nil, err
	}
	
	if respMap["error"] != nil {
		return nil, errors.New(string(resp))
	}
	
	ord := new(Order)
	ord.OrderID = int(respMap["id"].(float64))
	ord.Currency = currencyPair
	ord.Price, _ = strconv.ParseFloat(price, 64);
	ord.Amount, _ = strconv.ParseFloat(amount, 64);
	ord.Status = ORDER_UNFINISH;
	ord.OrderTime = int(time.Now().Unix())
	
	switch side {
	case "buy":
		ord.Side = BUY;
	case "sell":
		ord.Side = SELL;
	}
	
	return ord, nil
}

func (yunbi *YunBi)CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	params := url.Values{}
	params.Set("id", orderId)
	yunbi.buildPostForm("POST", API_URI_PREFIX + DELETE_ORDER_API, &params)
	
	resp, err := HttpPostForm(yunbi.client, API_URL + API_URI_PREFIX + DELETE_ORDER_API, params)
	if err != nil {
		log.Println(err, string(resp))
		return false, err
	}
	
	//println(string(resp))
	
	respMap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		log.Println(err, string(resp))
		return false, err
	}
	
	return true, nil
}

func (yunbi *YunBi)GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	params := url.Values{}
	params.Set("id", orderId)
	yunbi.buildPostForm("GET", API_URI_PREFIX + GET_ORDER_API, &params)

	respMap, err := HttpGet(yunbi.client, API_URL + API_URI_PREFIX + GET_ORDER_API + "?" + params.Encode())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println(respMap)

	ord := yunbi.parseOrder(respMap)
	ord.Currency = currency
	return &ord, err
}

func (yunbi *YunBi) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	params := url.Values{}
	params.Set("market", convertCurrencyPair(currency))
	params.Set("state", "wait")

	yunbi.buildPostForm("GET", API_URI_PREFIX + PLACE_ORDER_API, &params)

	resp, err := yunbi.client.Get(API_URL + API_URI_PREFIX + PLACE_ORDER_API + "?" + params.Encode())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ordersMap := make([]map[string]interface{}, 1)
	err = json.Unmarshal(body, &ordersMap)

	orders := make([]Order, 0)
	for _, v := range ordersMap {
		ord := yunbi.parseOrder(v)
		ord.Currency = currency
		orders = append(orders, ord)
	}

	log.Println(ordersMap)

	return orders, nil
}

func (yunbi *YunBi) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	return nil, nil
}

func (yunbi *YunBi) GetKlineRecords(currency CurrencyPair, period , size, since int) ([]Kline, error) {
	return nil, nil
}

func (yunbi *YunBi) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("unimplements")
}

func (yunbi *YunBi)parseOrder(orderMap map[string]interface{}) Order {
	ord := Order{}
	ord.OrderID = int(orderMap["id"].(float64))
	ord.Amount, _ = strconv.ParseFloat(orderMap["volume"].(string), 64)
	ord.Price, _ = strconv.ParseFloat(orderMap["price"].(string), 64)
	ord.AvgPrice, _ = strconv.ParseFloat(orderMap["avg_price"].(string), 64)
	ord.DealAmount, _ = strconv.ParseFloat(orderMap["executed_volume"].(string), 64)

	switch orderMap["state"].(string) {
	case "wait":
		ord.Status = ORDER_UNFINISH
	case "done":
		ord.Status = ORDER_FINISH
	case "cancel":
		ord.Status = ORDER_CANCEL
	default:
		log.Println("unknow state :", orderMap["state"])
	}

	switch orderMap["side"].(string) {
	case "buy":
		ord.Side = BUY
	case "sell":
		ord.Side = SELL
	}
	return ord
}

func (yunbi *YunBi) buildPostForm(httpMethod, apiURI string, postForm *url.Values) error {
	postForm.Set("access_key", yunbi.accessKey);
	postForm.Set("tonce", fmt.Sprintf("%d", time.Now().UnixNano() / 1000000));
	
	params := postForm.Encode();
	payload := httpMethod + "|" + apiURI + "|" + params
	//println(payload)
	
	sign, err := GetParamHmacSHA256Sign(yunbi.secretKey, payload);
	if err != nil {
		return err;
	}
	
	postForm.Set("signature", sign);
	//postForm.Del("secret_key")
	return nil;
}

func convertCurrencyPair(currencyPair CurrencyPair) string {
	return strings.ToLower(currencyPair.ToSymbol(""))
}