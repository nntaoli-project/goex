package zb

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	MARKET_URL = "http://api.zb.com/data/v1/"
	TICKER_API = "ticker?market=%s"
	DEPTH_API  = "depth?market=%s&size=%d"

	TRADE_URL                 = "https://trade.zb.com/api/"
	GET_ACCOUNT_API           = "getAccountInfo"
	GET_ORDER_API             = "getOrder"
	GET_UNFINISHED_ORDERS_API = "getUnfinishedOrdersIgnoreTradeType"
	CANCEL_ORDER_API          = "cancelOrder"
	PLACE_ORDER_API           = "order"
	WITHDRAW_API              = "withdraw"
	CANCELWITHDRAW_API        = "cancelWithdraw"
)

type Zb struct {
	httpClient *http.Client
	accessKey,
	secretKey string
}

func New(httpClient *http.Client, accessKey, secretKey string) *Zb {
	return &Zb{httpClient, accessKey, secretKey}
}

func (zb *Zb) GetExchangeName() string {
	return ZB
}

func (zb *Zb) GetTicker(currency CurrencyPair) (*Ticker, error) {
	symbol := currency.AdaptBchToBcc().AdaptUsdToUsdt().ToSymbol("_")
	resp, err := HttpGet(zb.httpClient, MARKET_URL+fmt.Sprintf(TICKER_API, symbol))
	if err != nil {
		return nil, err
	}
	//log.Println(resp)
	tickermap := resp["ticker"].(map[string]interface{})

	ticker := new(Ticker)
	ticker.Pair = currency
	ticker.Date, _ = strconv.ParseUint(resp["date"].(string), 10, 64)
	ticker.Buy, _ = strconv.ParseFloat(tickermap["buy"].(string), 64)
	ticker.Sell, _ = strconv.ParseFloat(tickermap["sell"].(string), 64)
	ticker.Last, _ = strconv.ParseFloat(tickermap["last"].(string), 64)
	ticker.High, _ = strconv.ParseFloat(tickermap["high"].(string), 64)
	ticker.Low, _ = strconv.ParseFloat(tickermap["low"].(string), 64)
	ticker.Vol, _ = strconv.ParseFloat(tickermap["vol"].(string), 64)

	return ticker, nil
}

func (zb *Zb) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	symbol := currency.AdaptBchToBcc().AdaptUsdToUsdt().ToSymbol("_")
	resp, err := HttpGet(zb.httpClient, MARKET_URL+fmt.Sprintf(DEPTH_API, symbol, size))
	if err != nil {
		return nil, err
	}

	//log.Println(resp)

	asks, isok1 := resp["asks"].([]interface{})
	bids, isok2 := resp["bids"].([]interface{})
	
	if isok2 != true || isok1 != true {
		return nil, errors.New("no depth data!")
	}
	//log.Println(asks)
	//log.Println(bids)

	depth := new(Depth)
	depth.Pair = currency

	for _, e := range bids {
		var r DepthRecord
		ee := e.([]interface{})
		r.Amount = ee[1].(float64)
		r.Price = ee[0].(float64)

		depth.BidList = append(depth.BidList, r)
	}

	for _, e := range asks {
		var r DepthRecord
		ee := e.([]interface{})
		r.Amount = ee[1].(float64)
		r.Price = ee[0].(float64)

		depth.AskList = append(depth.AskList, r)
	}

	return depth, nil
}

func (zb *Zb) buildPostForm(postForm *url.Values) error {
	postForm.Set("accesskey", zb.accessKey)

	payload := postForm.Encode()
	secretkeySha, _ := GetSHA(zb.secretKey)

	sign, err := GetParamHmacMD5Sign(secretkeySha, payload)
	if err != nil {
		return err
	}

	postForm.Set("sign", sign)
	//postForm.Del("secret_key")
	postForm.Set("reqTime", fmt.Sprintf("%d", time.Now().UnixNano()/1000000))
	return nil
}

func (zb *Zb) GetAccount() (*Account, error) {
	params := url.Values{}
	params.Set("method", "getAccountInfo")
	zb.buildPostForm(&params)
	//log.Println(params.Encode())
	resp, err := HttpPostForm(zb.httpClient, TRADE_URL+GET_ACCOUNT_API, params)
	if err != nil {
		return nil, err
	}

	var respmap map[string]interface{}
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		log.Println("json unmarshal error")
		return nil, err
	}

	if respmap["code"] != nil && respmap["code"].(float64) != 1000 {
		return nil, errors.New(string(resp))
	}

	acc := new(Account)
	acc.Exchange = zb.GetExchangeName()
	acc.SubAccounts = make(map[Currency]SubAccount)

	resultmap := respmap["result"].(map[string]interface{})
	coins := resultmap["coins"].([]interface{})

	acc.NetAsset = ToFloat64(resultmap["netAssets"])
	acc.Asset = ToFloat64(resultmap["totalAssets"])

	for _, v := range coins {
		vv := v.(map[string]interface{})
		subAcc := SubAccount{}
		subAcc.Amount = ToFloat64(vv["available"])
		subAcc.ForzenAmount = ToFloat64(vv["freez"])
		subAcc.Currency = NewCurrency(vv["key"].(string), "").AdaptBchToBcc()
		acc.SubAccounts[subAcc.Currency] = subAcc
	}

	//log.Println(string(resp))
	//log.Println(acc)

	return acc, nil
}

func (zb *Zb) placeOrder(amount, price string, currency CurrencyPair, tradeType int) (*Order, error) {
	symbol := currency.AdaptBchToBcc().AdaptUsdToUsdt().ToSymbol("_")
	params := url.Values{}
	params.Set("method", "order")
	params.Set("price", price)
	params.Set("amount", amount)
	params.Set("currency", symbol)
	params.Set("tradeType", fmt.Sprintf("%d", tradeType))
	zb.buildPostForm(&params)

	resp, err := HttpPostForm(zb.httpClient, TRADE_URL+PLACE_ORDER_API, params)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	//log.Println(string(resp));

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	code := respmap["code"].(float64)
	if code != 1000 {
		log.Println(string(resp))
		return nil, errors.New(fmt.Sprintf("%.0f", code))
	}

	orid := respmap["id"].(string)

	order := new(Order)
	order.Amount, _ = strconv.ParseFloat(amount, 64)
	order.Price, _ = strconv.ParseFloat(price, 64)
	order.Status = ORDER_UNFINISH
	order.Currency = currency
	order.OrderTime = int(time.Now().UnixNano() / 1000000)
	order.OrderID, _ = strconv.Atoi(orid)

	switch tradeType {
	case 0:
		order.Side = SELL
	case 1:
		order.Side = BUY
	}

	return order, nil
}

func (zb *Zb) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return zb.placeOrder(amount, price, currency, 1)
}

func (zb *Zb) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return zb.placeOrder(amount, price, currency, 0)
}

func (zb *Zb) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	symbol := currency.AdaptBchToBcc().AdaptUsdToUsdt().ToSymbol("_")
	params := url.Values{}
	params.Set("method", "cancelOrder")
	params.Set("id", orderId)
	params.Set("currency", symbol)
	zb.buildPostForm(&params)

	resp, err := HttpPostForm(zb.httpClient, TRADE_URL+CANCEL_ORDER_API, params)
	if err != nil {
		log.Println(err)
		return false, err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		log.Println(err)
		return false, err
	}

	code := respmap["code"].(float64)

	if code == 1000 {
		return true, nil
	}

	//log.Println(respmap)
	return false, errors.New(fmt.Sprintf("%.0f", code))
}

func parseOrder(order *Order, ordermap map[string]interface{}) {
	//log.Println(ordermap)
	//order.Currency = currency;
	order.OrderID, _ = strconv.Atoi(ordermap["id"].(string))
	order.OrderID2 = ordermap["id"].(string)
	order.Amount = ordermap["total_amount"].(float64)
	order.DealAmount = ordermap["trade_amount"].(float64)
	order.Price = ordermap["price"].(float64)
	//	order.Fee = ordermap["fees"].(float64)
	if order.DealAmount > 0 {
		order.AvgPrice = ToFloat64(ordermap["trade_money"]) / order.DealAmount
	} else {
		order.AvgPrice = 0
	}

	order.OrderTime = int(ordermap["trade_date"].(float64))

	orType := ordermap["type"].(float64)
	switch orType {
	case 0:
		order.Side = SELL
	case 1:
		order.Side = BUY
	default:
		log.Printf("unknown order type %f", orType)
	}

	_status := TradeStatus(ordermap["status"].(float64))
	switch _status {
	case 0:
		order.Status = ORDER_UNFINISH
	case 1:
		order.Status = ORDER_CANCEL
	case 2:
		order.Status = ORDER_FINISH
	case 3:
		order.Status = ORDER_UNFINISH
	}

}

func (zb *Zb) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	symbol := currency.AdaptBchToBcc().AdaptUsdToUsdt().ToSymbol("_")
	params := url.Values{}
	params.Set("method", "getOrder")
	params.Set("id", orderId)
	params.Set("currency", symbol)
	zb.buildPostForm(&params)

	resp, err := HttpPostForm(zb.httpClient, TRADE_URL+GET_ORDER_API, params)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//println(string(resp))
	ordermap := make(map[string]interface{})
	err = json.Unmarshal(resp, &ordermap)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	order := new(Order)
	order.Currency = currency

	parseOrder(order, ordermap)

	return order, nil
}

func (zb *Zb) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	params := url.Values{}
	symbol := currency.AdaptBchToBcc().AdaptUsdToUsdt().ToSymbol("_")
	params.Set("method", "getUnfinishedOrdersIgnoreTradeType")
	params.Set("currency", symbol)
	params.Set("pageIndex", "1")
	params.Set("pageSize", "100")
	zb.buildPostForm(&params)

	resp, err := HttpPostForm(zb.httpClient, TRADE_URL+GET_UNFINISHED_ORDERS_API, params)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	respstr := string(resp)
	//println(respstr)

	if strings.Contains(respstr, "\"code\":3001") {
		log.Println(respstr)
		return nil, nil
	}

	var resps []interface{}
	err = json.Unmarshal(resp, &resps)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var orders []Order
	for _, v := range resps {
		ordermap := v.(map[string]interface{})
		order := Order{}
		order.Currency = currency
		parseOrder(&order, ordermap)
		orders = append(orders, order)
	}

	return orders, nil
}

func (zb *Zb) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	return nil, nil
}

func (zb *Zb) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	return nil, nil
}

func (zb *Zb) Withdraw(amount string, currency Currency, fees, receiveAddr, safePwd string) (string, error) {
	params := url.Values{}
	params.Set("method", "withdraw")
	params.Set("currency", strings.ToLower(currency.AdaptBchToBcc().String()))
	params.Set("amount", amount)
	params.Set("fees", fees)
	params.Set("receiveAddr", receiveAddr)
	params.Set("safePwd", safePwd)
	zb.buildPostForm(&params)

	resp, err := HttpPostForm(zb.httpClient, TRADE_URL+WITHDRAW_API, params)
	if err != nil {
		log.Println("withdraw fail.", err)
		return "", err
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		log.Println(err, string(resp))
		return "", err
	}

	if respMap["code"].(float64) == 1000 {
		return respMap["id"].(string), nil
	}

	return "", errors.New(string(resp))
}

func (zb *Zb) CancelWithdraw(id string, currency Currency, safePwd string) (bool, error) {
	params := url.Values{}
	params.Set("method", "cancelWithdraw")
	params.Set("currency", strings.ToLower(currency.AdaptBchToBcc().String()))
	params.Set("downloadId", id)
	params.Set("safePwd", safePwd)
	zb.buildPostForm(&params)

	resp, err := HttpPostForm(zb.httpClient, TRADE_URL+CANCELWITHDRAW_API, params)
	if err != nil {
		log.Println("cancel withdraw fail.", err)
		return false, err
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		log.Println(err, string(resp))
		return false, err
	}

	if respMap["code"].(float64) == 1000 {
		return true, nil
	}

	return false, errors.New(string(resp))
}

func (zb *Zb) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("unimplements")
}

func (zb *Zb) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("unsupport the market order")
}

func (zb *Zb) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("unsupport the market order")
}
