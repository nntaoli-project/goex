package coinbig

import (
	. "github.com/nntaoli-project/goex"
	"net/http"
	"net/url"

	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	API_BASE_URL = "https://www.coinbig.com"
)

type CoinBig struct {
	httpClient *http.Client
	accessKey,
	secretKey string
	timeoffset int64
}

func New(client *http.Client, api_key, secret_key string) *CoinBig {
	return &CoinBig{accessKey: api_key, secretKey: secret_key, httpClient: client}
}

func (cb *CoinBig) GetExchangeName() string {
	return "coinbig.com"
}

func (cb *CoinBig) buildSigned(params *url.Values) string {
	s, _ := GetParamMD5Sign(cb.secretKey, encodeWithoutEscape(*params)+fmt.Sprintf("&secret_key=%s", cb.secretKey))
	params.Set("sign", strings.ToUpper(s))

	return s
}

func (cb *CoinBig) GetAccount() (*Account, error) {
	api_url := API_BASE_URL + "/api/publics/v1/userinfo"

	params := url.Values{}
	params.Set("time", strconv.Itoa(int(time.Now().UnixNano()/1000000)))
	params.Set("apikey", cb.accessKey)
	cb.buildSigned(&params)
	body, err := HttpPostForm(cb.httpClient, api_url, params)
	if err != nil {
		return nil, err
	}
	//log.Println("body", string(body))

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(body, &bodyDataMap)
	if err != nil {
		// log.Println("respData", string(body))
		return nil, err
	}
	if bodyDataMap["code"].(float64) != 0 {
		// log.Println("respData", string(body))
		return nil, errors.New(bodyDataMap["msg"].(string))
	}

	balances, isok := bodyDataMap["data"].(map[string]interface{})
	if isok != true {
		return nil, errors.New("No account data!")
	}
	info, isok := balances["info"].(map[string]interface{})
	if isok != true {
		return nil, errors.New("No account data!")
	}
	acc := Account{}
	acc.Exchange = cb.GetExchangeName()
	acc.SubAccounts = make(map[Currency]SubAccount)

	free := info["free"].(map[string]interface{})
	//log.Println(free)

	for k, v := range free {
		currency := NewCurrency(k, "")
		sub := SubAccount{}
		sub = acc.SubAccounts[currency]
		sub.Amount = ToFloat64(v)
		acc.SubAccounts[currency] = sub
	}

	freezed := info["freezed"].(map[string]interface{})

	for k, v := range freezed {
		currency := NewCurrency(k, "")
		sub := SubAccount{}
		sub = acc.SubAccounts[currency]
		sub.Amount = ToFloat64(v)
		acc.SubAccounts[currency] = sub
	}
	return &acc, nil
}

func (cb *CoinBig) placeOrder(amount, price string, pair CurrencyPair, orderType, orderSide string) (*Order, error) {
	api_url := API_BASE_URL + "/api/publics/v1/trade"

	params := url.Values{}
	params.Set("time", strconv.Itoa(int(time.Now().UnixNano()/1000000)))
	params.Set("apikey", cb.accessKey)
	params.Set("symbol", strings.ToLower(pair.String()))

	ty := orderSide
	if orderType == "market" {
		ty += "_market"
		if orderSide == "buy" {
			params.Set("price", price)
		} else {
			params.Set("amount", amount)
		}
	} else {
		params.Set("price", price)
		params.Set("amount", amount)
	}
	params.Set("type", ty)

	cb.buildSigned(&params)

	body, err := HttpPostForm(cb.httpClient, api_url, params)
	if err != nil {
		return nil, err
	}
	// log.Println("body", string(body))

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(body, &bodyDataMap)
	if err != nil {
		// log.Println("respData", string(body))
		return nil, err
	}
	if bodyDataMap["code"].(float64) != 0 {
		// log.Println("respData", string(body))
		return nil, errors.New(bodyDataMap["msg"].(string))
	}

	data, isok := bodyDataMap["data"].(map[string]interface{})
	if isok != true {
		return nil, errors.New("No order information")
	}

	orderId := ToInt(data["order_id"])
	side := BUY
	if orderSide == "sell" {
		side = SELL
	}

	return &Order{
		Currency:   pair,
		OrderID:    orderId,
		OrderID2:   strconv.Itoa(orderId),
		Price:      ToFloat64(price),
		Amount:     ToFloat64(amount),
		DealAmount: 0,
		AvgPrice:   0,
		Side:       TradeSide(side),
		Status:     ORDER_UNFINISH,
		OrderTime:  int(time.Now().Unix())}, nil
}

func (cb *CoinBig) LimitBuy(amount, price string, currencyPair CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	return cb.placeOrder(amount, price, currencyPair, "limit", "buy")
}

func (cb *CoinBig) LimitSell(amount, price string, currencyPair CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	return cb.placeOrder(amount, price, currencyPair, "limit", "sell")
}

func (cb *CoinBig) MarketBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return cb.placeOrder(amount, price, currencyPair, "market", "buy")
}

func (cb *CoinBig) MarketSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return cb.placeOrder(amount, price, currencyPair, "market", "sell")
}

func (cb *CoinBig) CancelOrder(orderId string, currencyPair CurrencyPair) (bool, error) {
	path := API_BASE_URL + "/api/publics/v1/cancel_order"
	params := url.Values{}

	params.Set("apikey", cb.accessKey)

	params.Set("time", strconv.Itoa(int(time.Now().UnixNano()/1000000)))
	params.Set("order_id", orderId)
	cb.buildSigned(&params)

	resp, err := HttpPostForm(cb.httpClient, path, params)

	//log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return false, err
	}

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(resp, &bodyDataMap)
	if err != nil {
		// log.Println("respData", string(resp))
		return false, err
	}
	if bodyDataMap["code"].(float64) != 0 {
		// log.Println("respData", string(body))
		return false, errors.New(bodyDataMap["msg"].(string))
	}
	return true, nil
}

func (cb *CoinBig) CancelOrders(orderId []string) (bool, error) {
	path := API_BASE_URL + "/order/batchCancel"
	params := url.Values{}

	orders := strings.Join(orderId, ",")
	//log.Println("orders", orders)
	params.Set("orderIds", orders)
	cb.buildSigned(&params)

	resp, err := HttpPostForm(cb.httpClient, path, params)

	// log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return false, err
	}

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(resp, &bodyDataMap)
	if err != nil {
		// log.Println("respData", string(resp))
		return false, err
	}

	if bodyDataMap["status"].(string) != "1000" {
		return false, errors.New(bodyDataMap["msg"].(string))
	}

	return true, nil
}

func (cb *CoinBig) GetOneOrder(orderId string, currencyPair CurrencyPair) (*Order, error) {
	path := API_BASE_URL + "/api/publics/v1/order_info"
	params := url.Values{}

	params.Set("apikey", cb.accessKey)
	params.Set("order_id", orderId)
	//params.Set("symbol", strings.ToLower(currencyPair.String()))

	params.Set("time", strconv.Itoa(int(time.Now().UnixNano()/1000000)))
	//params.Set("size", "50")
	//params.Set("type", "1,2")
	cb.buildSigned(&params)

	resp, err := HttpPostForm(cb.httpClient, path, params)

	//log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return nil, err
	}

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(resp, &bodyDataMap)
	if err != nil {
		// log.Println("respData", string(resp))
		return nil, err
	}
	if bodyDataMap["code"].(float64) != 0 {
		// log.Println("respData", string(body))
		return nil, errors.New(bodyDataMap["msg"].(string))
	}

	data, _ := bodyDataMap["data"].(map[string]interface{})
	//list, _ := data["orders"].([]interface{})
	result, _ := data["result"].(bool)
	if result != true {
		return nil, errors.New("no order info")
	}
	vv := data["orders"].(map[string]interface{})
	return &Order{
		Currency:   currencyPair,
		OrderID2:   vv["order_id"].(string),
		Price:      ToFloat64(vv["price"]),
		Amount:     ToFloat64(vv["count"]),
		DealAmount: ToFloat64(vv["count"]) - ToFloat64(vv["leftCount"]),
		Side:       orderTypeAdapter(vv["isLimit"].(int), vv["entrustType"].(int)),
		Status:     orderStatusAdapter(vv["status"].(int)),
		OrderTime:  ToInt(vv["create_date"]),
	}, nil
}
func (cb *CoinBig) GetUnfinishOrders(currencyPair CurrencyPair) ([]Order, error) {
	path := API_BASE_URL + "/api/publics/v1/orders_info"
	params := url.Values{}

	params.Set("apikey", cb.accessKey)
	params.Set("symbol", strings.ToLower(currencyPair.String()))

	params.Set("time", strconv.Itoa(int(time.Now().UnixNano()/1000000)))
	params.Set("size", "50")
	params.Set("type", "1,2")
	cb.buildSigned(&params)

	resp, err := HttpPostForm(cb.httpClient, path, params)

	//log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return nil, err
	}

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(resp, &bodyDataMap)
	if err != nil {
		// log.Println("respData", string(resp))
		return nil, err
	}
	if bodyDataMap["code"].(float64) != 0 {
		// log.Println("respData", string(body))
		return nil, errors.New(bodyDataMap["msg"].(string))
	}

	data, _ := bodyDataMap["data"].(map[string]interface{})
	list, _ := data["orders"].([]interface{})
	result, _ := data["result"].(bool)
	if result != true {
		return nil, errors.New("no order info")
	}
	var orders []Order
	for _, v := range list {
		vv := v.(map[string]interface{})
		//yyyy-MM-dd HH:mm:ss
		orders = append(orders, Order{
			Currency:   currencyPair,
			OrderID2:   vv["order_id"].(string),
			Price:      ToFloat64(vv["price"]),
			Amount:     ToFloat64(vv["count"]),
			DealAmount: ToFloat64(vv["count"]) - ToFloat64(vv["leftCount"]),
			Side:       orderTypeAdapter(vv["isLimit"].(int), vv["entrustType"].(int)),
			Status:     orderStatusAdapter(vv["status"].(int)),
			OrderTime:  ToInt(vv["create_date"]),
		})
	}
	return orders, nil
}
func (cb *CoinBig) GetOrderHistorys(currencyPair CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	return nil, nil
}
func (cb *CoinBig) GetDepth(size int, currencyPair CurrencyPair) (*Depth, error) {
	path := API_BASE_URL + "/api/publics/v1/depth"
	path += fmt.Sprintf("?size=%d&symbol=%s", size, currencyPair.String())
	bodyDataMap, err := HttpGet(cb.httpClient, path)

	//log.Println("resp:", bodyDataMap, "err:", err)
	if err != nil {
		return nil, err
	}
	if bodyDataMap["code"].(float64) != 0 {
		// log.Println("respData", string(body))
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
	//sort.Sort(depth.AskList)
	//sort.Sort(sort.Reverse(depth.BidList))
	depth.AskList = depth.AskList[0:size]
	depth.BidList = depth.BidList[0:size]

	return depth, nil
}

func (cb *CoinBig) GetTicker(currencyPair CurrencyPair) (*Ticker, error) {
	path := API_BASE_URL + "/api/publics/v1/ticker"
	path += fmt.Sprintf("?symbol=%s", currencyPair.String())
	bodyDataMap, err := HttpGet(cb.httpClient, path)

	//log.Println("resp:", bodyDataMap, "err:", err)
	if err != nil {
		return nil, err
	}
	if bodyDataMap["code"].(float64) != 0 {
		// log.Println("respData", string(body))
		return nil, errors.New(bodyDataMap["msg"].(string))
	}

	data := bodyDataMap["data"].(map[string]interface{})
	tickerMap := data["ticker"].(map[string]interface{})
	var ticker Ticker

	ticker.Date = ToUint64(data["date"])
	ticker.Last = ToFloat64(tickerMap["last"])
	ticker.Buy = ToFloat64(tickerMap["buy"])
	ticker.Sell = ToFloat64(tickerMap["sell"])
	ticker.Low = ToFloat64(tickerMap["low"])
	ticker.High = ToFloat64(tickerMap["high"])
	ticker.Vol = ToFloat64(tickerMap["vol"])
	return &ticker, nil

}

func (cb *CoinBig) GetServerSync() error {
	path := API_BASE_URL + "/api/publics/v1/getClientIpAndServerTime"
	bodyDataMap, err := HttpGet(cb.httpClient, path)

	log.Println("GetServerSync resp:", bodyDataMap, "err:", err)
	if err != nil {
		return err
	}
	if bodyDataMap["code"].(float64) != 0 {
		// log.Println("respData", string(body))
		return errors.New(bodyDataMap["msg"].(string))
	}
	data, _ := bodyDataMap["data"].(map[string]interface{})

	stime := int64(ToInt(data["time"]))

	st := time.Unix(stime/1000, (stime%1000)*1000)
	lt := time.Now()
	offset := st.Sub(lt).Seconds()
	cb.timeoffset = int64(offset)

	return nil

}

func (cb *CoinBig) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录
func (cb *CoinBig) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}

func orderTypeAdapter(t1, t2 int) TradeSide {
	if t1 == 0 { //是否是限价单，0为限价单，1市价单
		switch t2 { //类型:0买入,1卖出
		case 1:
			return SELL
		case 0:
			return BUY
		}

	} else if t1 == 1 {
		switch t2 { //类型:0买入,1卖出
		case 1:
			return SELL_MARKET
		case 0:
			return BUY_MARKET
		}
	}
	return BUY

}

func orderStatusAdapter(s int) TradeStatus {
	switch s {
	case 1:
		return ORDER_UNFINISH
	case 2:
		return ORDER_PART_FINISH
	case 3:
		return ORDER_FINISH
	case 4, 5:
		return ORDER_CANCEL
	case 6:
		return ORDER_REJECT
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
