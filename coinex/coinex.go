package coinex

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	. "github.com/nntaoli-project/goex"
)

type CoinEx struct {
	httpClient *http.Client
	accessKey,
	secretKey string
}

var (
	baseurl = "https://api.coinex.com/v1/"
)

func New(client *http.Client, accessKey, secretKey string) *CoinEx {
	return &CoinEx{client, accessKey, secretKey}
}

func (coinex *CoinEx) GetExchangeName() string {
	return COINEX
}

func (coinex *CoinEx) GetTicker(currency CurrencyPair) (*Ticker, error) {
	params := url.Values{}
	params.Set("market", currency.ToSymbol(""))
	datamap, err := coinex.doRequest("GET", "market/ticker", &params)
	if err != nil {
		return nil, err
	}

	tickermap := datamap["ticker"].(map[string]interface{})

	return &Ticker{
		Date: ToUint64(datamap["date"]) / 1000,
		Last: ToFloat64(tickermap["last"]),
		Buy:  ToFloat64(tickermap["buy"]),
		Sell: ToFloat64(tickermap["sell"]),
		High: ToFloat64(tickermap["high"]),
		Low:  ToFloat64(tickermap["low"]),
		Vol:  ToFloat64(tickermap["vol"])}, nil
}

func (coinex *CoinEx) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	params := url.Values{}
	params.Set("market", currency.ToSymbol(""))
	params.Set("merge", "0.00000001")
	params.Set("limit", fmt.Sprint(size))

	datamap, err := coinex.doRequest("GET", "market/depth", &params)
	if err != nil {
		return nil, err
	}

	dep := Depth{}
	dep.AskList = make([]DepthRecord, 0, size)
	dep.BidList = make([]DepthRecord, 0, size)

	asks := datamap["asks"].([]interface{})
	bids := datamap["bids"].([]interface{})

	for _, v := range asks {
		r := v.([]interface{})
		dep.AskList = append(dep.AskList, DepthRecord{ToFloat64(r[0]), ToFloat64(r[1])})
	}

	for _, v := range bids {
		r := v.([]interface{})
		dep.BidList = append(dep.BidList, DepthRecord{ToFloat64(r[0]), ToFloat64(r[1])})
	}

	sort.Sort(sort.Reverse(dep.AskList))

	return &dep, nil
}

func (coinex *CoinEx) placeLimitOrder(side, amount, price string, pair CurrencyPair) (*Order, error) {
	params := url.Values{}
	params.Set("market", pair.ToSymbol(""))
	params.Set("type", side)
	params.Set("amount", amount)
	params.Set("price", price)

	retmap, err := coinex.doRequest("POST", "order/limit", &params)
	if err != nil {
		return nil, err
	}

	order := coinex.adaptOrder(retmap, pair)

	return &order, nil
}

func (coinex *CoinEx) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return coinex.placeLimitOrder("buy", amount, price, currency)
}

func (coinex *CoinEx) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return coinex.placeLimitOrder("sell", amount, price, currency)
}

func (coinex *CoinEx) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (coinex *CoinEx) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}

func (coinex *CoinEx) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	params := url.Values{}
	params.Set("id", orderId)
	params.Set("market", currency.ToSymbol(""))
	_, err := coinex.doRequest("DELETE", "order/pending", &params)
	if err != nil {
		return false, err
	}
	//	log.Println(retmap)
	return true, nil
}

func (coinex *CoinEx) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	params := url.Values{}
	params.Set("id", orderId)
	params.Set("market", currency.ToSymbol(""))
	retmap, err := coinex.doRequest("GET", "order", &params)
	if err != nil {
		if "Order not found" == err.Error() {
			return nil, EX_ERR_NOT_FIND_ORDER
		}
		return nil, err
	}
	order := coinex.adaptOrder(retmap, currency)
	return &order, nil
}

func (coinex *CoinEx) GetPendingOrders(page, limit int, pair CurrencyPair) ([]Order, error) {
	params := url.Values{}
	params.Set("page", fmt.Sprint(page))
	params.Set("limit", fmt.Sprint(limit))
	params.Set("market", pair.ToSymbol(""))

	retmap, err := coinex.doRequest("GET", "order/pending", &params)
	if err != nil {
		return nil, err
	}

	//log.Println(retmap)

	datamap, isok := retmap["data"].([]interface{})
	if !isok {
		log.Println(datamap)
		return nil, errors.New("response format error")
	}

	var orders []Order
	for _, v := range datamap {
		vv := v.(map[string]interface{})
		orders = append(orders, coinex.adaptOrder(vv, pair))
	}

	return orders, nil

}

func (coinex *CoinEx) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	return coinex.GetPendingOrders(1, 100, currency)
}

func (coinex *CoinEx) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implement")
}

type coinexDifficulty struct {
	Code int `json:"code"`
	Data struct {
		Difficulty string `json:"difficulty"`
		Prediction string `json:"prediction"`
		UpdateTime int    `json:"update_time"`
	} `json:"data"`
	Message string `json:"message"`
}

func (coinex *CoinEx) GetDifficulty() (limit, cur float64, err error) {
	buf, err := coinex.doRequestInner("GET", "order/mining/difficulty", &url.Values{})
	if nil != err {
		log.Printf("GetDifficulty - http.NewRequest failed : %v", err)
		return 0.0, 0.0, err
	}

	var diff coinexDifficulty
	if err = json.Unmarshal(buf, &diff); nil != err {
		log.Printf("GetDifficulty - json.Unmarshal failed : %v", err)
		return 0.0, 0.0, err
	}
	limit, err = strconv.ParseFloat(diff.Data.Difficulty, 64)
	if nil != err {
		log.Printf("GetDifficulty - strconv.ParseFloat failed : %v", err)
		return 0.0, 0.0, err
	}
	cur, err = strconv.ParseFloat(diff.Data.Prediction, 64)
	if nil != err {
		log.Printf("GetDifficulty - strconv.ParseFloat failed : %v", err)
		return 0.0, 0.0, err
	}

	return limit, cur, nil
}

func (coinex *CoinEx) GetAccount() (*Account, error) {
	datamap, err := coinex.doRequest("GET", "balance", &url.Values{})
	if err != nil {
		return nil, err
	}
	//log.Println(datamap)
	acc := new(Account)
	acc.SubAccounts = make(map[Currency]SubAccount, 2)
	acc.Exchange = coinex.GetExchangeName()
	for c, v := range datamap {
		vv := v.(map[string]interface{})
		currency := NewCurrency(c, "")
		acc.SubAccounts[currency] = SubAccount{
			Currency:     currency,
			Amount:       ToFloat64(vv["available"]),
			ForzenAmount: ToFloat64(vv["frozen"])}
	}
	return acc, nil
}

func (coinex *CoinEx) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录
func (coinex *CoinEx) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}

func (coinex *CoinEx) doRequestInner(method, uri string, params *url.Values) (buf []byte, err error) {
	reqUrl := baseurl + uri

	headermap := map[string]string{
		"Content-Type": "application/json; charset=utf-8"}

	if !strings.HasPrefix(uri, "market") {
		params.Set("access_id", coinex.accessKey)
		params.Set("tonce", fmt.Sprint(time.Now().UnixNano()/int64(time.Millisecond)))
		//	println(params.Encode() + "&secret_key=" + coinex.secretKey)
		sign, _ := GetParamMD5Sign("", params.Encode()+"&secret_key="+coinex.secretKey)
		headermap["authorization"] = strings.ToUpper(sign)
	}

	if ("GET" == method || "DELETE" == method) && len(params.Encode()) > 0 {
		reqUrl += "?" + params.Encode()
	}

	var paramStr string = ""
	if "POST" == method {
		//to json
		paramStr = params.Encode()
		var parammap map[string]string = make(map[string]string, 2)
		for _, v := range strings.Split(paramStr, "&") {
			vv := strings.Split(v, "=")
			parammap[vv[0]] = vv[1]
		}
		jsonData, _ := json.Marshal(parammap)
		paramStr = string(jsonData)
	}

	return NewHttpRequest(coinex.httpClient, method, reqUrl, paramStr, headermap)
}

func (coinex *CoinEx) doRequest(method, uri string, params *url.Values) (map[string]interface{}, error) {
	resp, err := coinex.doRequestInner(method, uri, params)

	if err != nil {
		return nil, err
	}

	retmap := make(map[string]interface{}, 1)
	err = json.Unmarshal(resp, &retmap)
	if err != nil {
		return nil, err
	}

	if ToInt(retmap["code"]) != 0 {
		return nil, errors.New(retmap["message"].(string))
	}

	//	log.Println(retmap)
	datamap := retmap["data"].(map[string]interface{})

	return datamap, nil
}

func (coinex *CoinEx) adaptTradeSide(side string) TradeSide {
	switch side {
	case "sell":
		return SELL
	case "buy":
		return BUY
	}
	return BUY
}

func (coinex *CoinEx) adaptTradeStatus(status string) TradeStatus {
	var tradeStatus TradeStatus = ORDER_UNFINISH
	switch status {
	case "not_deal":
		tradeStatus = ORDER_UNFINISH
	case "done":
		tradeStatus = ORDER_FINISH
	case "partly":
		tradeStatus = ORDER_PART_FINISH
	case "cancel":
		tradeStatus = ORDER_CANCEL
	}
	return tradeStatus
}

func (coinex *CoinEx) adaptOrder(ordermap map[string]interface{}, pair CurrencyPair) Order {
	return Order{
		Currency:   pair,
		OrderID:    ToInt(ordermap["id"]),
		OrderID2:   fmt.Sprint(ToInt(ordermap["id"])),
		Amount:     ToFloat64(ordermap["amount"]),
		Price:      ToFloat64(ordermap["price"]),
		DealAmount: ToFloat64(ordermap["deal_amount"]),
		AvgPrice:   ToFloat64(ordermap["avg_price"]),
		Status:     coinex.adaptTradeStatus(ordermap["status"].(string)),
		Side:       coinex.adaptTradeSide(ordermap["type"].(string)),
		Fee:        ToFloat64(ordermap["deal_fee"]),
		OrderTime:  ToInt(ordermap["create_time"])}
}
