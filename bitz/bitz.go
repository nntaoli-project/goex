package bitz

import (
	. "github.com/nntaoli-project/GoEx"
	"net/http"
	"net/url"

	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	API_BASE_URL = "https://apiv2.bitz.com"
)

type BitZ struct {
	httpClient *http.Client
	accessKey,
	secretKey string
	tradePwd  string
}

var _INERNAL_KLINE_PERIOD_CONVERTER = map[int]string{
	KLINE_PERIOD_1MIN:  "1min",
	KLINE_PERIOD_5MIN:  "5min",
	KLINE_PERIOD_15MIN: "15min",
	KLINE_PERIOD_30MIN: "30min",
	KLINE_PERIOD_60MIN: "60min",
	KLINE_PERIOD_4H:    "4hour",
	KLINE_PERIOD_1DAY:  "1day",
	KLINE_PERIOD_1WEEK: "1week",
	KLINE_PERIOD_1MONTH:"1mon",
}

func New(client *http.Client, api_key, secret_key string, params map[string]interface{}) *BitZ {
	_tradePwd := ""
	if(params["tradePwd"] != nil){
		_tradePwd = params["tradePwd"].(string)
	}
	return &BitZ{accessKey: api_key, secretKey: secret_key, httpClient: client, tradePwd:_tradePwd}
}

func (ctx *BitZ) GetExchangeName() string {
	return "bitz.com"
}

func (ctx *BitZ) buildSigned(params *url.Values) string {
	s, _ := GetParamMD5Sign(ctx.secretKey, encodeWithoutEscape(*params)+ctx.secretKey)
	params.Set("sign", s)

	return s
}

func (ctx *BitZ) GetAccount() (*Account, error) {
	bodyDataMap, err := ctx.postSignedRequest("/Assets/getUserAssets", url.Values{})
	if err != nil {
		return nil, err
	}

	if bodyDataMap["status"].(float64) != 200 {
		return nil, errors.New(bodyDataMap["msg"].(string))
	}

	balances, isok := bodyDataMap["data"].(map[string]interface{})
	if isok != true {
		return nil, errors.New("No account data!")
	}
	info, isok := balances["info"].([]interface{})
	if isok != true {
		return nil, errors.New("No account data!")
	}
	acc := Account{}
	acc.Exchange = ctx.GetExchangeName()
	acc.SubAccounts = make(map[Currency]SubAccount)

	for _, v := range info {
		vv := v.(map[string]interface{})
		currency := NewCurrency(vv["name"].(string), "")
		sub := SubAccount{}
		sub = acc.SubAccounts[currency]
		sub.Amount = ToFloat64(vv["over"])
		sub.ForzenAmount = ToFloat64(vv["lock"])
		acc.SubAccounts[currency] = sub
	}

	return &acc, nil
}

func (ctx *BitZ) placeOrder(amount, price string, pair CurrencyPair, orderSide int) (*Order, error) {
	params := url.Values{}
	params.Set("tradePwd", ctx.tradePwd)
	params.Set("symbol", strings.ToLower(pair.String()))
	params.Set("price", price)
	params.Set("number", amount)
	params.Set("type", fmt.Sprintf("%d", orderSide))

	bodyDataMap, err := ctx.postSignedRequest("/Trade/addEntrustSheet", params)
	if err != nil {
		return nil, err
	}

	if bodyDataMap["status"].(float64) != 200 {
		return nil, errors.New(bodyDataMap["msg"].(string))
	}

	data, isok := bodyDataMap["data"].(map[string]interface{})
	if isok != true {
		return nil, errors.New("No order information")
	}
	order := ctx.parseOrder(pair, data)

	return &order, nil
}

func (ctx *BitZ) LimitBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return ctx.placeOrder(amount, price, currencyPair, 1)
}

func (ctx *BitZ) LimitSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return ctx.placeOrder(amount, price, currencyPair, 2)
}

func (ctx *BitZ) MarketBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (ctx *BitZ) MarketSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (ctx *BitZ) CancelOrder(orderId string, currencyPair CurrencyPair) (bool, error) {
	params := url.Values{}
	params.Set("entrustSheetId", orderId)

	bodyDataMap, err := ctx.postSignedRequest("/Trade/cancelEntrustSheet", params)
	if err != nil {
		return false, err
	}

	if bodyDataMap["status"].(float64) != 200 {
		return false, errors.New(bodyDataMap["msg"].(string))
	}

	return true, nil
}

func (ctx *BitZ) GetOneOrder(orderId string, currencyPair CurrencyPair) (*Order, error) {
	params := url.Values{}
	params.Set("entrustSheetId", orderId)

	bodyDataMap, err := ctx.postSignedRequest("/Trade/getEntrustSheetInfo", params)
	if err != nil {
		return nil, err
	}

	if bodyDataMap["status"].(float64) != 200 {
		return nil, errors.New(bodyDataMap["msg"].(string))
	}

	data, _ := bodyDataMap["data"].(map[string]interface{})
	order := ctx.parseOrder(currencyPair, data)
	return &order, nil
}
func (ctx *BitZ) GetUnfinishOrders(currencyPair CurrencyPair) ([]Order, error) {
	params := url.Values{}
	params.Set("coinFrom", strings.ToLower(currencyPair.CurrencyA.String()))
	params.Set("coinTo", strings.ToLower(currencyPair.CurrencyB.String()))

	bodyDataMap, err := ctx.postSignedRequest("/Trade/getUserNowEntrustSheet", params)
	if err != nil {
		return nil, err
	}

	if bodyDataMap["status"].(float64) != 200 {
		return nil, errors.New(bodyDataMap["msg"].(string))
	}

	data, _ := bodyDataMap["data"].(map[string]interface{})
	list, _ := data["data"].([]interface{})

	var orders []Order
	for _, v := range list {
		vv := v.(map[string]interface{})
		orders = append(orders, ctx.parseOrder(currencyPair, vv))
	}
	return orders, nil
}
/**
发送签名请求
 */
func (ctx *BitZ) postSignedRequest(path string, params url.Values) (map[string]interface{}, error){
	url := API_BASE_URL + path
	params.Set("apiKey", ctx.accessKey)
	params.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano() % 1000000))
	params.Set("timeStamp", strconv.Itoa(int(time.Now().UnixNano()/1000000000)))
	ctx.buildSigned(&params)

	resp, err := HttpPostForm2(ctx.httpClient, url, params,
		map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36"})

	if err != nil {
		return nil, err
	}

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(resp, &bodyDataMap)
	if err != nil {
		return nil, err
	}

	return bodyDataMap, nil
}
func (ctx *BitZ) parseOrder(currencyPair CurrencyPair, data map[string]interface{}) Order{
	var orderId string
	orderId, ok := data["id"].(string)
	if(!ok){
		orderId = fmt.Sprintf("%d", int(ToFloat64(data["id"])))
	}

	return Order{
		Currency:   currencyPair,
		OrderID:    int(ToFloat64(data["id"])),
		OrderID2:   orderId,
		Price:      ToFloat64(data["price"]),
		Amount:     ToFloat64(data["number"]),
		DealAmount: ToFloat64(data["numberDeal"]),
		Side:       orderTypeAdapter(data["flag"].(string)),
		Status:     orderStatusAdapter(ToInt(data["status"])),
		OrderTime:  ToInt(data["created"]) * 1000,
	}
}
func (ctx *BitZ) GetOrderHistorys(currencyPair CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	return nil, nil
}
func (ctx *BitZ) GetDepth(size int, currencyPair CurrencyPair) (*Depth, error) {
	path := API_BASE_URL + "/Market/depth"
	path += fmt.Sprintf("?symbol=%s", strings.ToLower(currencyPair.String()))
	bodyDataMap, err := HttpGet(ctx.httpClient, path)

	if err != nil {
		return nil, err
	}
	if bodyDataMap["status"].(float64) != 200 {
		return nil, errors.New(bodyDataMap["msg"].(string))
	}

	data := bodyDataMap["data"].(map[string]interface{})

	bids, _ := data["bids"].([]interface{})
	asks, _ := data["asks"].([]interface{})

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
	depth.AskList = depth.AskList[0:size]
	depth.BidList = depth.BidList[0:size]

	return depth, nil
}

func (ctx *BitZ) GetTicker(currencyPair CurrencyPair) (*Ticker, error) {
	path := API_BASE_URL + "/Market/ticker"
	path += fmt.Sprintf("?symbol=%s", strings.ToLower(currencyPair.String()))
	bodyDataMap, err := HttpGet(ctx.httpClient, path)

	if err != nil {
		return nil, err
	}
	if bodyDataMap["status"].(float64) != 200 {
		return nil, errors.New(bodyDataMap["msg"].(string))
	}

	tickerMap := bodyDataMap["data"].(map[string]interface{})
	var ticker Ticker

	ticker.Pair = currencyPair
	ticker.Date = ToUint64(bodyDataMap["time"])
	ticker.Last = ToFloat64(tickerMap["now"])
	ticker.Buy = ToFloat64(tickerMap["bidPrice"])
	ticker.Sell = ToFloat64(tickerMap["askPrice"])
	ticker.Low = ToFloat64(tickerMap["low"])
	ticker.High = ToFloat64(tickerMap["high"])
	ticker.Vol = ToFloat64(tickerMap["volume"])
	return &ticker, nil

}

func (ctx *BitZ) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	path := API_BASE_URL + "/Market/kline"
	path += fmt.Sprintf("?symbol=%s&resolution=%s&size=%d", strings.ToLower(currency.String()), _INERNAL_KLINE_PERIOD_CONVERTER[period], size)
	bodyDataMap, err := HttpGet(ctx.httpClient, path)

	if err != nil {
		return nil, err
	}
	if bodyDataMap["status"].(float64) != 200 {
		return nil, errors.New(bodyDataMap["msg"].(string))
	}

	data := bodyDataMap["data"].(map[string]interface{})
	list := data["bars"].([]interface{})

	var klineRecords []Kline
	for _, record := range list {
		_record := record.(map[string]interface{})
		r := Kline{
			Pair:currency,
			Timestamp:int64(ToFloat64(_record["time"])),
			Open:ToFloat64(_record["open"]),
			Close:ToFloat64(_record["close"]),
			High:ToFloat64(_record["high"]),
			Low:ToFloat64(_record["low"]),
			Vol:ToFloat64(_record["volume"]),
		}
		klineRecords = append(klineRecords, r)
	}

	return klineRecords, nil
}

//非个人，整个交易所的交易记录
func (ctx *BitZ) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	path := API_BASE_URL + "/Market/order"
	path += fmt.Sprintf("?symbol=%s", strings.ToLower(currencyPair.String()))
	bodyDataMap, err := HttpGet(ctx.httpClient, path)

	if err != nil {
		return nil, err
	}
	if bodyDataMap["status"].(float64) != 200 {
		return nil, errors.New(bodyDataMap["msg"].(string))
	}

	data := bodyDataMap["data"].([]interface{})

	trades := []Trade{}
	for _, trade := range data {
		_trade := trade.(map[string]interface{})
		item := Trade{
			Tid: int64(ToFloat64(_trade["id"])),
			Type:_trade["s"].(string),
			Amount: ToFloat64(_trade["n"]),
			Price: ToFloat64(_trade["p"]),
			Date: int64(ToFloat64(_trade["T"])) * 1000,
		}
		trades = append(trades, item)
	}

	return trades, nil
}

func orderTypeAdapter(side string) TradeSide {
	switch side { //类型:sale和?
	case "sale":
		return SELL
	default:
		return BUY
	}
}

func orderStatusAdapter(s int) TradeStatus {
	switch s {
	case 0:
		return ORDER_UNFINISH
	case 1:
		return ORDER_PART_FINISH
	case 2:
		return ORDER_FINISH
	case 3:
		return ORDER_CANCEL
	}
	return ORDER_REJECT

}

func encodeWithoutEscape(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		prefix := k + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(v)
		}
	}
	return buf.String()
}
