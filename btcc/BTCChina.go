package btcc

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type BTCChina struct {
	httpClient *http.Client
	accessKey,
	secretKey string
}

type ReqBody struct {
	Id     int           `json:"id"`
	Params []interface{} `json:"params"`
	Method string        `json:"method"`
	Sign   string        `json:"-"`
	Tonce  string        `json:"-"`
}

var (
	_MARKET_API_URL   = "https://data.btcchina.com/data"
	_TRADE_API_V1_URL = "https://api.btcchina.com/api_trade_v1.php"
)

func NewBTCChina(client *http.Client, accessKey, secretKey string) *BTCChina {
	return &BTCChina{client, accessKey, secretKey}
}

func (btch *BTCChina) GetTicker(currency CurrencyPair) (*Ticker, error) {
	tickerResp, err := HttpGet(btch.httpClient, fmt.Sprintf("%s/ticker?market=%s",
		_MARKET_API_URL, strings.ToLower(currency.ToSymbol(""))))

	if err != nil {
		return nil, err
	}

	//log.Println(tickerResp)

	tickermap, ok := tickerResp["ticker"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Get Ticker Error")
	}

	ticker := Ticker{}
	ticker.Buy = ToFloat64(tickermap["buy"])
	ticker.Sell = ToFloat64(tickermap["sell"])
	ticker.High = ToFloat64(tickermap["high"])
	ticker.Low = ToFloat64(tickermap["low"])
	ticker.Last = ToFloat64(tickermap["last"])
	ticker.Vol = ToFloat64(tickermap["vol"])
	date := tickermap["date"].(float64)
	ticker.Date = uint64(date)

	return &ticker, nil
}

func (btch *BTCChina) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	depthresp, err := HttpGet(btch.httpClient, fmt.Sprintf("%s/orderbook?market=%s&limit=%d",
		_MARKET_API_URL, strings.ToLower(currency.ToSymbol("")), size))

	if err != nil {
		return nil, err
	}

	depth := Depth{}
	asks := depthresp["asks"].([]interface{})
	bids := depthresp["bids"].([]interface{})

	//log.Println(asks)
	//log.Println(bids)

	for _, v := range asks {
		var dr DepthRecord
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price = ToFloat64(vv)
			case 1:
				dr.Amount = ToFloat64(vv)
			}
		}
		depth.AskList = append(depth.AskList, dr)
	}

	for _, v := range bids {
		var dr DepthRecord
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price = ToFloat64(vv)
			case 1:
				dr.Amount = ToFloat64(vv)
			}
		}
		depth.BidList = append(depth.BidList, dr)
	}

	return &depth, nil
}

func (btch *BTCChina) GetKlineRecords(currency CurrencyPair, period , size, since int) ([]Kline, error) {
	panic("unimplement")
}

func (btch *BTCChina) GetAccount() (*Account, error) {
	respmap, err := btch.sendAuthorizationRequst("getAccountInfo", []interface{}{})
	if err != nil {
		return nil, err
	}

	result := respmap["result"].(map[string]interface{})
	balance := result["balance"].(map[string]interface{})
	frozen := result["frozen"].(map[string]interface{})

	acc := new(Account)
	acc.SubAccounts = make(map[Currency]SubAccount, 3)
	//log.Println(balance)
	//log.Println(frozen)

	for c, v := range balance {
		vv := v.(map[string]interface{})
		_frozen := frozen[c].(map[string]interface{})

		sub := SubAccount{
			Amount:       ToFloat64(vv["amount"]),
			ForzenAmount: ToFloat64(_frozen["amount"])}
		var currency Currency

		switch c {
		case "cny":
			currency = CNY
		case "btc":
			currency = BTC
		case "ltc":
			currency = LTC
		default:
			currency = UNKNOWN
		}

		acc.SubAccounts[currency] = sub
	}

	return acc, nil
}

func (btch *BTCChina) placeorder(method, amount, price string, currencyPair CurrencyPair) (*Order, error) {
	respmap, err := btch.sendAuthorizationRequst(method, []interface{}{
		price, amount, strings.ToUpper(currencyPair.ToSymbol(""))})

	if err != nil {
		return nil, err
	}

	ord := new(Order)
	ord.OrderID = ToInt(respmap["result"])
	ord.Currency = currencyPair
	ord.Amount = ToFloat64(amount)
	ord.Price = ToFloat64(price)

	return ord, nil
}

func (btch *BTCChina) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	ord, err := btch.placeorder("buyOrder2", amount, price, currency)
	ord.Side = BUY
	return ord, err
}

func (btch *BTCChina) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	ord, err := btch.placeorder("sellOrder2", amount, price, currency)
	ord.Side = SELL
	return ord, err
}

func (btch *BTCChina) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("unimplement")
}

func (btch *BTCChina) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("unimplement")
}

func (btch *BTCChina) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	respmap, err := btch.sendAuthorizationRequst("cancelOrder",
		[]interface{}{ToInt(orderId), strings.ToUpper(currency.ToSymbol(""))})
	if err != nil {
		return false, err
	}
	return respmap["result"].(bool), nil
}

func (btch *BTCChina) toOrder(ordermap map[string]interface{}) Order {
	ord := Order{}
	ord.OrderID = ToInt(ordermap["id"])
	ord.Price = ToFloat64(ordermap["price"])
	ord.Amount = ToFloat64(ordermap["amount_original"])
	ord.DealAmount = ToFloat64(ordermap["amount"])
	ord.AvgPrice = ToFloat64(ordermap["price"])

	switch ordermap["status"].(string) {
	case "closed":
		ord.Status = ORDER_FINISH
	case "cancelled":
		ord.Status = ORDER_CANCEL
	case "open":
		ord.Status = ORDER_UNFINISH
	case "pending":
		ord.Status = ORDER_UNFINISH
	default:
		ord.Status = ORDER_UNFINISH
	}

	switch ordermap["type"].(string) {
	case "ask":
		ord.Side = SELL
	case "bid":
		ord.Side = BUY
	}

	return ord
}

func (btch *BTCChina) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	respmap, err := btch.sendAuthorizationRequst("getOrder",
		[]interface{}{ToInt(orderId), strings.ToUpper(currency.ToSymbol(""))})
	if err != nil {
		return nil, err
	}

	result := respmap["result"].(map[string]interface{})
	ordermap := result["order"].(map[string]interface{})

	ord := btch.toOrder(ordermap)
	return &ord, nil
}

func (btch *BTCChina) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	respmap, err := btch.sendAuthorizationRequst("getOrders", []interface{}{
		true, strings.ToUpper(currency.ToSymbol(""))})
	if err != nil {
		return nil, err
	}

	orders := make([]Order, 0)
	result := respmap["result"].(map[string]interface{})
	ordersmap := result["order"].([]interface{})
	for _, ord := range ordersmap {
		orders = append(orders, btch.toOrder(ord.(map[string]interface{})))
	}

	return orders, nil
}

func (btch *BTCChina) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("unimplement")
}

//非个人，整个交易所的交易记录
func (btch *BTCChina) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("unimplement")
}

func (btch *BTCChina) GetExchangeName() string {
	return "btcchina.com"
}

func (btch *BTCChina) GetBasicAuth(sign string) string {
	authStr := btch.accessKey + ":" + sign
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(authStr))
	return basicAuth
}

func (btch *BTCChina) sendAuthorizationRequst(method string, params []interface{}) (map[string]interface{}, error) {
	reqParams := btch.buildPostForm(method, params)
	reqJsonParams, _ := json.Marshal(reqParams)
	println(string(reqJsonParams))

	resp, err := HttpPostForm3(btch.httpClient,
		_TRADE_API_V1_URL,
		string(reqJsonParams),
		map[string]string{"Json-Rpc-Tonce": reqParams.Tonce,
			"Authorization": btch.GetBasicAuth(reqParams.Sign)})

	if err != nil {
		return nil, err
	}

	var respmap map[string]interface{}
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return nil, err
	}

	if respmap["error"] != nil {
		return nil, errors.New(string(resp))
	}

	return respmap, nil
}

func (btch *BTCChina) buildPostForm(method string, paramAr []interface{}) ReqBody {
	/*
	   tonce=1377743828095093
	      &accesskey=1d87effa-e84d-48c1-a172-0232b86305dd
	      &requestmethod=post
	      &id=1
	      &method=getAccountInfo
	      &params=
	*/

	paramStr := ""

	if paramAr != nil && len(paramAr) > 0 {
		for _, p := range paramAr {
			switch p.(type) {
			case float64:
				paramStr = fmt.Sprintf("%s,%f", paramStr, p)
			case string:
				paramStr = fmt.Sprintf("%s,%s", paramStr, p)
			case int:
				paramStr = fmt.Sprintf("%s,%d", paramStr, p)
			case bool:
				paramStr = fmt.Sprintf("%s,%t", paramStr, p)
			}
		}
		paramStr = paramStr[1:]
		//println(paramStr)
	}

	tonce := strconv.FormatInt(time.Now().UnixNano(), 10)[0:16]
	signPyload := fmt.Sprintf("tonce=%s&accesskey=%s&requestmethod=post&id=1&method=%s&params=%s",
		tonce, btch.accessKey, method, paramStr)
	signStr, _ := GetParamHmacSHA1Sign(btch.secretKey, signPyload)

	reqbody := ReqBody{}
	reqbody.Id = 1
	reqbody.Method = method
	reqbody.Params = paramAr
	reqbody.Tonce = tonce
	reqbody.Sign = signStr

	return reqbody
}
