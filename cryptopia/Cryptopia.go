package cryptopia

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	. "github.com/nntaoli-project/goex"

	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const (
	API_BASE_URL = "https://www.cryptopia.co.nz/api/"

	TICKERS_URI                = "GetMarkets"
	TICKER_URI                 = "GetMarket/"
	DEPTH_URI                  = "GetMarketOrders"
	CANCEL_URI                 = "CancelTrade"
	SUBMIT_ORDER_URI           = "SubmitTrade"
	DEFAULT_HTTPCLIENT_TIMEOUT = 30 // HTTP client timeout
)

type Cryptopia struct {
	accessKey,
	secretKey string
	httpClient *http.Client
	debug      bool
}
type jsonResponse struct {
	Success bool            `json:"success"`
	Error   string          `json:"error"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"data"`
}

type placedOrderResponse struct {
	OrderId int `json:"OrderId"`
}

type cancelOrderPayload struct {
	OrderId     int `json:"OrderId"`
	TradePairId int `json:"TradePairId"`
	Type        string
}

func New(client *http.Client, accessKey, secretKey string) *Cryptopia {
	debug := false
	return &Cryptopia{accessKey, secretKey, client, debug}
}

func (cta *Cryptopia) GetExchangeName() string {
	return CRYPTOPIA
}

func (cta *Cryptopia) GetTickers(currency CurrencyPair) (*Ticker, error) {
	return cta.GetTicker(currency)

	//tickerUri := API_BASE_URL + TICKERS_URI
	////log.Println("tickerUrl:", tickerUri)
	//bodyDataMap, err := HttpGet(cta.httpClient, tickerUri)
	////log.Println("Cryptopia bodyDataMap:", tickerUri, bodyDataMap)
	//
	//if err != nil {
	//	log.Println(err)
	//	return nil, err
	//}
	//
	//if result, isok := bodyDataMap["success"].(bool); isok == true && result != true {
	//	log.Println("bodyDataMap[\"success\"]", isok, result)
	//	return nil, errors.New("err")
	//}
	////timestamp := time.Now().Unix()
	//
	//panic("not implement")
	//return nil, nil
}

func (cta *Cryptopia) GetTicker(currency CurrencyPair) (*Ticker, error) {
	currency = cta.adaptCurrencyPair(currency)

	tickerUri := API_BASE_URL + TICKER_URI + currency.ToSymbol("_")
	//log.Println("tickerUrl:", tickerUri)
	bodyDataMap, err := HttpGet(cta.httpClient, tickerUri)
	//log.Println("Cryptopia bodyDataMap:", tickerUri, bodyDataMap)

	if err != nil {
		//log.Println(err)
		return nil, err
	}
	tickerMap, isok := bodyDataMap["Data"].(map[string]interface{})
	if isok != true {
		//log.Println("Cryptopia bodyDataMap:", tickerUri, bodyDataMap)
		//log.Println("bodyDataMap[\"Error\"]", bodyDataMap["Error"].(string))
		//return nil, errors.New(bodyDataMap["Error"].(string))
		return nil, errors.New("ERR")
	}
	var ticker Ticker

	timestamp := time.Now().Unix()

	//fmt.Println(bodyDataMap)
	ticker.Date = uint64(timestamp)
	ticker.Last, _ = tickerMap["LastPrice"].(float64)

	ticker.Buy, _ = tickerMap["BidPrice"].(float64)
	ticker.Sell, _ = tickerMap["AskPrice"].(float64)
	ticker.Vol, _ = tickerMap["Volume"].(float64)
	//log.Println("Cryptopia", currency, "ticker:", ticker)
	return &ticker, nil
}

func (cta *Cryptopia) GetTickerInBuf(currency CurrencyPair) (*Ticker, error) {
	return cta.GetTicker(currency)
}

//GetDepth get orderbook
func (cta *Cryptopia) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	currency = cta.adaptCurrencyPair(currency)

	depthURI := fmt.Sprintf("%s%s/%s/%d", API_BASE_URL, DEPTH_URI, currency.ToSymbol("_"), size)
	bodyDataMap, err := HttpGet(cta.httpClient, depthURI)
	if err != nil {
		return nil, err
	}
	depthMap, isok := bodyDataMap["Data"].(map[string]interface{})
	if isok != true {
		return nil, err
	}
	bids := depthMap["Buy"].([]interface{})
	asks := depthMap["Sell"].([]interface{})

	depth := new(Depth)
	for _, bid := range bids {
		_bid := bid.(map[string]interface{})
		amount := ToFloat64(_bid["Volume"])
		price := ToFloat64(_bid["Price"])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.BidList = append(depth.BidList, dr)
	}
	for _, ask := range asks {
		_ask := ask.(map[string]interface{})
		amount := ToFloat64(_ask["Volume"])
		price := ToFloat64(_ask["Price"])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.AskList = append(depth.AskList, dr)
	}

	return depth, nil
}

func (cta *Cryptopia) adaptCurrencyPair(pair CurrencyPair) CurrencyPair {
	var currencyA Currency
	var currencyB Currency

	if pair.CurrencyA == BCC {
		currencyA = BCH
	} else {
		currencyA = pair.CurrencyA
	}
	currencyB = pair.CurrencyB
	//if pair.BaseCurrency == USDT {
	//	currencyB = USD
	//} else {
	//	currencyB = pair.BaseCurrency
	//}

	return NewCurrencyPair(currencyA, currencyB)
}

func (cta *Cryptopia) currencyPairToSymbol(currencyPair CurrencyPair) string {
	return strings.ToUpper(currencyPair.ToSymbol("_"))
}

func (cta *Cryptopia) placeOrder(orderType, side, amount, price string, pair CurrencyPair) (*Order, error) {
	payload := map[string]interface{}{
		"Market":   cta.currencyPairToSymbol(pair),
		"Amount":   amount,
		"Rate":     price,
		"Type":     orderType,
		"Exchange": "cryptopia"}
	p, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}
	resp, err := cta.do("POST", SUBMIT_ORDER_URI, string(p), true)
	if err != nil {
		return nil, err
	}

	var jsonResp jsonResponse
	err = json.Unmarshal(resp, &jsonResp)
	if err != nil {
		log.Println(string(resp))
		return nil, err
	}
	order := new(Order)

	if jsonResp.Success {
		data := new(placedOrderResponse)
		err = json.Unmarshal(jsonResp.Result, &data)
		if err != nil {
			return nil, err
		}
		order.Currency = pair
		order.OrderID = data.OrderId
		order.Amount = ToFloat64(amount)
		order.Price = ToFloat64(price)
	} else {
		return nil, errors.New(jsonResp.Error)
	}

	switch side {
	case "buy":
		order.Side = BUY
	case "sell":
		order.Side = SELL
	}
	return order, nil
}

//LimitBuy ...
func (cta *Cryptopia) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return cta.placeOrder("Buy", "Buy", amount, price, currency)
}

//LimitSell ...
func (cta *Cryptopia) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return cta.placeOrder("Sell", "Sell", amount, price, currency)
}

func (cta *Cryptopia) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (cta *Cryptopia) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (cta *Cryptopia) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	payload := map[string]interface{}{
		"OrderId": ToInt(orderId),
	}
	p, err := json.Marshal(&payload)
	if err != nil {
		return false, err
	}
	resp, err := cta.do("POST", CANCEL_URI, string(p), true)
	if err != nil {
		return false, err
	}
	var jsonResp jsonResponse
	err = json.Unmarshal(resp, &jsonResp)
	if err != nil {
		log.Println(string(resp))
		return false, err
	}
	return jsonResp.Success, nil
}

func (cta *Cryptopia) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}
func (cta *Cryptopia) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implements")
}

func (cta *Cryptopia) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implements")
}

func (cta *Cryptopia) GetAccount() (*Account, error) {
	panic("not implements")
}
func (cta *Cryptopia) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implements")
}
func (cta *Cryptopia) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implements")
}

func (cta *Cryptopia) do(method string, resource string, payload string, authNeeded bool) (response []byte, err error) {
	connectTimer := time.NewTimer(DEFAULT_HTTPCLIENT_TIMEOUT * time.Second)

	var rawurl string
	if strings.HasPrefix(resource, "http") {
		rawurl = resource
	} else {
		rawurl = fmt.Sprintf("%s%s", API_BASE_URL, resource)
	}

	req, err := http.NewRequest(method, rawurl, strings.NewReader(payload))
	if err != nil {
		return
	}
	if method == "POST" || method == "PUT" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	}
	req.Header.Add("Accept", "application/json")

	// Auth
	if authNeeded {
		if len(cta.accessKey) == 0 || len(cta.secretKey) == 0 {
			err = errors.New("You need to set API Key and API Secret to call this method")
			return
		}
		nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
		md5 := md5.Sum([]byte(payload))
		signature := cta.accessKey + method + strings.ToLower(url.QueryEscape(req.URL.String())) +
			nonce + base64.StdEncoding.EncodeToString(md5[:])
		secret, _ := base64.StdEncoding.DecodeString(cta.secretKey)
		mac := hmac.New(sha256.New, secret)
		_, err = mac.Write([]byte(signature))
		sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))
		auth := "amx " + cta.accessKey + ":" + sig + ":" + nonce
		req.Header.Add("Authorization", auth)
	}

	resp, err := cta.doTimeoutRequest(connectTimer, req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)
	//fmt.Println(fmt.Sprintf("reponse %s", response), err)
	if err != nil {
		return response, err
	}
	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
	}
	return response, err
}

func (cta Cryptopia) dumpRequest(r *http.Request) {
	if r == nil {
		log.Print("dumpReq ok: <nil>")
		return
	}
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Print("dumpReq err:", err)
	} else {
		log.Print("dumpReq ok:", string(dump))
	}
}

func (cta Cryptopia) dumpResponse(r *http.Response) {
	if r == nil {
		log.Print("dumpResponse ok: <nil>")
		return
	}
	dump, err := httputil.DumpResponse(r, true)
	if err != nil {
		log.Print("dumpResponse err:", err)
	} else {
		log.Print("dumpResponse ok:", string(dump))
	}
}

// doTimeoutRequest do a HTTP request with timeout
func (cta *Cryptopia) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	// Do the request in the background so we can check the timeout
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		if cta.debug {
			cta.dumpRequest(req)
		}
		resp, err := cta.httpClient.Do(req)
		if cta.debug {
			cta.dumpResponse(resp)
		}
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, errors.New("timeout on reading data from Bittrex API")
	}
}

func (cta *Cryptopia) SetDebug(enable bool) {
	cta.debug = enable
}
