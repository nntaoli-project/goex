package ocx

import (
	"crypto"
	"crypto/hmac"
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

var TimeOffset int64 = 0

const (
	API_BASE_URL = "https://openapi.ocx.com"
	V2           = "/api/v2/"
	//API_URL      = API_BASE_URL + V2

	TICKER_API  = "market/ticker/%s"
	TICKERS_API = "market/tickers"
	DEPTH_API   = "depth?market_code=%s"
	SERVER_TIME = "timestamp"

	TRADE_URL       = "orders"
	GET_ACCOUNT_API = "accounts"
	GET_ORDER_API   = "orders/%s"

	CANCEL_ORDER_API     = "orders/:%s/cancel"
	CANCEL_ALL_ORDER_API = "orders/clear"

	PLACE_ORDER_API = "orders"
)

type Ocx struct {
	httpClient *http.Client
	accessKey,
	secretKey string
}

func New(client *http.Client, apikey, secretkey string) *Ocx {
	return &Ocx{accessKey: apikey, secretKey: secretkey, httpClient: client}
}

func (o *Ocx) GetExchangeName() string {
	return "ocx.com"
}

func (o *Ocx) GetServerTime() int64 {
	url := API_BASE_URL + V2 + SERVER_TIME
	respmap, err := HttpGet(o.httpClient, url)
	if err != nil {
		return 0
	}
	data := respmap["data"].(interface{})
	d := data.(map[string]interface{})
	ts := d["timestamp"].(float64)
	servertime := int64(ts)
	now := time.Now().Unix()
	TimeOffset = servertime - now
	return servertime
}

func (o *Ocx) GetTicker(currencyPair CurrencyPair) (*Ticker, error) {
	url := API_BASE_URL + V2 + fmt.Sprintf(TICKER_API, strings.ToLower(currencyPair.ToSymbol("")))
	respmap, err := HttpGet(o.httpClient, url)
	if err != nil {
		return nil, err
	}
	//log.Println("ticker respmap:", respmap)

	tickmap, ok := respmap["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("tick assert error")
	}
	ticker := new(Ticker)
	ticker.Pair = currencyPair
	ticker.Date = ToUint64(tickmap["timestamp"])
	ticker.Last = ToFloat64(tickmap["last"])
	ticker.Vol = ToFloat64(tickmap["volume"])
	ticker.Low = ToFloat64(tickmap["low"])
	ticker.High = ToFloat64(tickmap["high"])
	ticker.Buy = ToFloat64(tickmap["buy"])
	ticker.Sell = ToFloat64(tickmap["sell"])

	return ticker, nil

}

func (o *Ocx) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	url := API_BASE_URL + V2 + fmt.Sprintf(DEPTH_API, strings.ToLower(currency.ToSymbol("")))
	resp, err := HttpGet(o.httpClient, url)
	if err != nil {
		return nil, err
	}
	respmap, _ := resp["status"].(map[string]interface{})
	bids, ok1 := respmap["bids"].([]interface{})
	asks, ok2 := respmap["asks"].([]interface{})

	if !ok1 || !ok2 {
		return nil, errors.New("tick assert error")
	}

	depth := new(Depth)

	for _, r := range asks {
		var dr DepthRecord
		rr := r.([]interface{})
		dr.Price = ToFloat64(rr[0])
		dr.Amount = ToFloat64(rr[1])
		depth.AskList = append(depth.AskList, dr)
	}

	for _, r := range bids {
		var dr DepthRecord
		rr := r.([]interface{})
		dr.Price = ToFloat64(rr[0])
		dr.Amount = ToFloat64(rr[1])
		depth.BidList = append(depth.BidList, dr)
	}

	return depth, nil
}

func (o *Ocx) buildSigned(method, apiurl string, para *url.Values) string {
	param := ""
	if para != nil {
		param = para.Encode()
	}

	log.Println("param", param)

	sig_str := ""
	if para != nil {
		sig_str = method + "|" + apiurl + "|" + param
	} else {
		sig_str = method + "|" + apiurl
	}

	log.Println("sig_str", sig_str)

	mac := hmac.New(crypto.SHA256.New, []byte("abc"))

	mac.Write([]byte(sig_str))

	sum := mac.Sum(nil)

	return fmt.Sprintf("%x", sum)
}

func (o *Ocx) placeOrder(orderType, orderSide, amount, price string, pair CurrencyPair) (*Order, error) {
	uri := API_BASE_URL + V2 + TRADE_URL
	method := "POST"
	path := V2 + TRADE_URL
	params := url.Values{}
	tonce := strconv.Itoa(int(time.Now().UnixNano() / 1000000))
	params.Set("access_key", o.accessKey)
	params.Set("tonce", tonce)
	//params.Set("foo", "bar")

	signed := o.buildSigned(method, path, &params)
	
	f := fmt.Sprintf("access_key=%s&tonce=%s&signature=%s&market_code=%s&price=%s&side=%s&volume=%s",
		o.accessKey, tonce, signed, strings.ToLower(pair.ToSymbol("")), price, orderSide, amount)
	resp, err := HttpPostForm3(o.httpClient, uri, f, nil)
	//resp, err := HttpPostForm3(o.httpClient, uri, form.Encode(), nil)
	log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return nil, err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		log.Println(string(resp))
		return nil, err
	}

	orderId := ToInt(respmap["orderId"])
	if orderId <= 0 {
		return nil, errors.New(string(resp))
	}

	side := BUY
	if orderSide == "SELL" {
		side = SELL
	}

	return &Order{
		Currency:   pair,
		OrderID:    orderId,
		OrderID2:   fmt.Sprint(orderId),
		Price:      ToFloat64(price),
		Amount:     ToFloat64(amount),
		DealAmount: 0,
		AvgPrice:   0,
		Side:       TradeSide(side),
		Status:     ORDER_UNFINISH,
		OrderTime:  int(time.Now().Unix())}, nil

	panic("1")
}

func (o *Ocx) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return o.placeOrder("limit", "buy", amount, price, currency)
}

func (o *Ocx) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return o.placeOrder("limit", "sell", amount, price, currency)
}

func (o *Ocx) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (o *Ocx) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}

func (o *Ocx) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	path := API_BASE_URL + V2 + fmt.Sprintf(CANCEL_ORDER_API, strings.ToLower(currency.ToSymbol("")))
	params := url.Values{}

	params.Set("order_id", orderId)

	sign := o.buildSigned("POST", path, &params)
	log.Println("path", path, "params", params.Encode(), "sign", sign)

	resp, err := HttpPostForm2(o.httpClient, path, params, nil)
	log.Println("resp:", string(resp), "err:", err)
	return true, nil
}

func (o *Ocx) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	path := API_BASE_URL + V2 + fmt.Sprintf(GET_ORDER_API, orderId)
	para := url.Values{}
	para.Set("order_id", orderId)

	//sign := o.buildSigned("GET", path, &para)

	respmap, err := HttpGet2(o.httpClient, path, nil)

	if err != nil {
		return nil, err
	}

	log.Println(respmap)

	return nil, nil

}

func (o *Ocx) GetOrdersList() {
	panic("not implement")

}

func (o *Ocx) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implement")
}

func (o *Ocx) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implement")
}

func (o *Ocx) GetAccount() (*Account, error) {
	url := API_BASE_URL + V2 + GET_ACCOUNT_API
	//timestamp := strconv.FormatInt((time.Now().UnixNano() / 1000000), 10)

	//sign := o.buildSigned("GET", url, nil)

	respmap, err := HttpGet2(o.httpClient, url, nil)

	if err != nil {
		return nil, err
	}

	log.Println(respmap)

	if respmap["status"].(float64) != 0 {
		return nil, errors.New(respmap["msg"].(string))
	}

	acc := new(Account)
	acc.SubAccounts = make(map[Currency]SubAccount)
	acc.Exchange = o.GetExchangeName()

	balances := respmap["data"].([]interface{})
	for _, v := range balances {
		//log.Println(v)
		vv := v.(map[string]interface{})
		currency := NewCurrency(vv["currency"].(string), "")
		acc.SubAccounts[currency] = SubAccount{
			Currency:     currency,
			Amount:       ToFloat64(vv["available"]),
			ForzenAmount: ToFloat64(vv["frozen"]),
		}
	}

	return acc, nil

}

func (o *Ocx) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录
func (o *Ocx) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}
