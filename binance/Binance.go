package binance

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	API_BASE_URL = "https://api.binance.com/"
	API_V1       = API_BASE_URL + "api/v1/"
	API_V3       = API_BASE_URL + "api/v3/"

	TICKER_URI             = "ticker/24hr?symbol=%s"
	TICKERS_URI            = "ticker/allBookTickers"
	DEPTH_URI              = "depth?symbol=%s&limit=%d"
	ACCOUNT_URI            = "account?"
	ORDER_URI              = "order?"
	UNFINISHED_ORDERS_INFO = "openOrders?"
	KLINE_URI              = "klines"
	SERVER_TIME_URL        = "api/v1/time"
)

var _INERNAL_KLINE_PERIOD_CONVERTER = map[int]string{
	KLINE_PERIOD_1MIN:   "1m",
	KLINE_PERIOD_3MIN:   "3m",
	KLINE_PERIOD_5MIN:   "5m",
	KLINE_PERIOD_15MIN:  "15m",
	KLINE_PERIOD_30MIN:  "30m",
	KLINE_PERIOD_60MIN:  "1h",
	KLINE_PERIOD_2H:     "2h",
	KLINE_PERIOD_4H:     "4h",
	KLINE_PERIOD_6H:     "6h",
	KLINE_PERIOD_8H:     "8h",
	KLINE_PERIOD_12H:    "12h",
	KLINE_PERIOD_1DAY:   "1d",
	KLINE_PERIOD_3DAY:   "3d",
	KLINE_PERIOD_1WEEK:  "1w",
	KLINE_PERIOD_1MONTH: "1M",
}

type Binance struct {
	accessKey,
	secretKey string
	httpClient *http.Client
	timeoffset int64 //nanosecond
}

func (bn *Binance) buildParamsSigned(postForm *url.Values) error {
	postForm.Set("recvWindow", "6000000")
	tonce := strconv.FormatInt(time.Now().UnixNano()+bn.timeoffset, 10)[0:13]
	postForm.Set("timestamp", tonce)
	payload := postForm.Encode()
	sign, _ := GetParamHmacSHA256Sign(bn.secretKey, payload)
	postForm.Set("signature", sign)
	return nil
}

func New(client *http.Client, api_key, secret_key string) *Binance {
	bn := &Binance{accessKey: api_key, secretKey: secret_key, httpClient: client}
	bn.setTimeOffset()
	return bn
}

func (bn *Binance) GetExchangeName() string {
	return BINANCE
}

func (bn *Binance) setTimeOffset() error {
	respmap, err := HttpGet(bn.httpClient, API_BASE_URL+SERVER_TIME_URL)
	if err != nil {
		return err
	}

	stime := int64(ToInt(respmap["serverTime"]))
	st := time.Unix(stime/1000, 1000000*(stime%1000))
	lt := time.Now()
	offset := st.Sub(lt).Nanoseconds()
	bn.timeoffset = int64(offset)
	return nil
}

func (bn *Binance) GetTicker(currency CurrencyPair) (*Ticker, error) {
	currency2 := bn.adaptCurrencyPair(currency)
	tickerUri := API_V1 + fmt.Sprintf(TICKER_URI, currency2.ToSymbol(""))
	tickerMap, err := HttpGet(bn.httpClient, tickerUri)

	if err != nil {
		log.Println("GetTicker error:", err)
		return nil, err
	}

	var ticker Ticker
	ticker.Pair = currency
	t, _ := tickerMap["closeTime"].(float64)
	ticker.Date = uint64(t / 1000)
	ticker.Last = ToFloat64(tickerMap["lastPrice"])
	ticker.Buy = ToFloat64(tickerMap["bidPrice"])
	ticker.Sell = ToFloat64(tickerMap["askPrice"])
	ticker.Low = ToFloat64(tickerMap["lowPrice"])
	ticker.High = ToFloat64(tickerMap["highPrice"])
	ticker.Vol = ToFloat64(tickerMap["volume"])
	return &ticker, nil
}

func (bn *Binance) GetDepth(size int, currencyPair CurrencyPair) (*Depth, error) {
	if size > 100 {
		size = 100
	} else if size < 5 {
		size = 5
	}
	currencyPair2 := bn.adaptCurrencyPair(currencyPair)

	apiUrl := fmt.Sprintf(API_V1+DEPTH_URI, currencyPair2.ToSymbol(""), size)
	resp, err := HttpGet(bn.httpClient, apiUrl)
	if err != nil {
		log.Println("GetDepth error:", err)
		return nil, err
	}

	if _, isok := resp["code"]; isok {
		return nil, errors.New(resp["msg"].(string))
	}

	bids := resp["bids"].([]interface{})
	asks := resp["asks"].([]interface{})

	//log.Println(bids)
	//log.Println(asks)

	depth := new(Depth)
	depth.Pair = currencyPair
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

func (bn *Binance) placeOrder(amount, price string, pair CurrencyPair, orderType, orderSide string) (*Order, error) {
	pair = bn.adaptCurrencyPair(pair)
	path := API_V3 + ORDER_URI
	params := url.Values{}
	params.Set("symbol", pair.ToSymbol(""))
	params.Set("side", orderSide)
	params.Set("type", orderType)

	params.Set("quantity", amount)
	params.Set("timeInForce", "GTC")

	switch orderType {
	case "LIMIT":
		params.Set("price", price)
	}

	bn.buildParamsSigned(&params)

	resp, err := HttpPostForm2(bn.httpClient, path, params,
		map[string]string{"X-MBX-APIKEY": bn.accessKey})
	//log.Println("resp:", string(resp), "err:", err)
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
}

func (bn *Binance) GetAccount() (*Account, error) {
	params := url.Values{}
	bn.buildParamsSigned(&params)
	path := API_V3 + ACCOUNT_URI + params.Encode()
	respmap, err := HttpGet2(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//log.Println("respmap:", respmap)
	if _, isok := respmap["code"]; isok == true {
		return nil, errors.New(respmap["msg"].(string))
	}
	acc := Account{}
	acc.Exchange = bn.GetExchangeName()
	acc.SubAccounts = make(map[Currency]SubAccount)

	balances := respmap["balances"].([]interface{})
	for _, v := range balances {
		//log.Println(v)
		vv := v.(map[string]interface{})
		currency := NewCurrency(vv["asset"].(string), "").AdaptBccToBch()
		acc.SubAccounts[currency] = SubAccount{
			Currency:     currency,
			Amount:       ToFloat64(vv["free"]),
			ForzenAmount: ToFloat64(vv["locked"]),
		}
	}

	return &acc, nil
}

func (bn *Binance) LimitBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return bn.placeOrder(amount, price, currencyPair, "LIMIT", "BUY")
}

func (bn *Binance) LimitSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return bn.placeOrder(amount, price, currencyPair, "LIMIT", "SELL")
}

func (bn *Binance) MarketBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return bn.placeOrder(amount, price, currencyPair, "MARKET", "BUY")
}

func (bn *Binance) MarketSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return bn.placeOrder(amount, price, currencyPair, "MARKET", "SELL")
}

func (bn *Binance) CancelOrder(orderId string, currencyPair CurrencyPair) (bool, error) {
	currencyPair = bn.adaptCurrencyPair(currencyPair)
	path := API_V3 + ORDER_URI
	params := url.Values{}
	params.Set("symbol", currencyPair.ToSymbol(""))
	params.Set("orderId", orderId)

	bn.buildParamsSigned(&params)

	resp, err := HttpDeleteForm(bn.httpClient, path, params, map[string]string{"X-MBX-APIKEY": bn.accessKey})

	//log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return false, err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		log.Println(string(resp))
		return false, err
	}

	orderIdCanceled := ToInt(respmap["orderId"])
	if orderIdCanceled <= 0 {
		return false, errors.New(string(resp))
	}

	return true, nil
}

func (bn *Binance) GetOneOrder(orderId string, currencyPair CurrencyPair) (*Order, error) {
	params := url.Values{}
	currencyPair = bn.adaptCurrencyPair(currencyPair)
	params.Set("symbol", currencyPair.ToSymbol(""))
	if orderId != "" {
		params.Set("orderId", orderId)
	}
	params.Set("orderId", orderId)

	bn.buildParamsSigned(&params)
	path := API_V3 + ORDER_URI + params.Encode()

	respmap, err := HttpGet2(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	//log.Println(respmap)
	if err != nil {
		return nil, err
	}
	status := respmap["status"].(string)
	side := respmap["side"].(string)

	ord := Order{}
	ord.Currency = currencyPair
	ord.OrderID = ToInt(orderId)
	ord.OrderID2 = orderId

	if side == "SELL" {
		ord.Side = SELL
	} else {
		ord.Side = BUY
	}

	switch status {
	case "NEW":
		ord.Status = ORDER_UNFINISH
	case "FILLED":
		ord.Status = ORDER_FINISH
	case "PARTIALLY_FILLED":
		ord.Status = ORDER_PART_FINISH
	case "CANCELED":
		ord.Status = ORDER_CANCEL
	case "PENDING_CANCEL":
		ord.Status = ORDER_CANCEL_ING
	case "REJECTED":
		ord.Status = ORDER_REJECT
	}

	ord.Amount = ToFloat64(respmap["origQty"].(string))
	ord.Price = ToFloat64(respmap["price"].(string))
	ord.DealAmount = ToFloat64(respmap["executedQty"])
	ord.AvgPrice = ord.Price // response no avg price ， fill price

	return &ord, nil
}

func (bn *Binance) GetUnfinishOrders(currencyPair CurrencyPair) ([]Order, error) {
	params := url.Values{}
	currencyPair = bn.adaptCurrencyPair(currencyPair)
	params.Set("symbol", currencyPair.ToSymbol(""))

	bn.buildParamsSigned(&params)
	path := API_V3 + UNFINISHED_ORDERS_INFO + params.Encode()

	respmap, err := HttpGet3(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	//log.Println("respmap", respmap, "err", err)
	if err != nil {
		return nil, err
	}

	orders := make([]Order, 0)
	for _, v := range respmap {
		ord := v.(map[string]interface{})
		side := ord["side"].(string)
		orderSide := SELL
		if side == "BUY" {
			orderSide = BUY
		}

		orders = append(orders, Order{
			OrderID:   ToInt(ord["orderId"]),
			OrderID2:  fmt.Sprint(ToInt(ord["orderId"])),
			Currency:  currencyPair,
			Price:     ToFloat64(ord["price"]),
			Amount:    ToFloat64(ord["origQty"]),
			Side:      TradeSide(orderSide),
			Status:    ORDER_UNFINISH,
			OrderTime: ToInt(ord["time"])})
	}
	return orders, nil
}

func (bn *Binance) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	currency2 := bn.adaptCurrencyPair(currency)
	params := url.Values{}
	params.Set("symbol", currency2.ToSymbol(""))
	params.Set("interval", _INERNAL_KLINE_PERIOD_CONVERTER[period])
	params.Set("startTime", strconv.Itoa(since/1000000))
	params.Set("endTime", strconv.Itoa(int(time.Now().UnixNano()/1000000)))
	params.Set("limit", fmt.Sprintf("%d", size))

	klineUrl := API_V1 + KLINE_URI + "?" + params.Encode()
	fmt.Println(klineUrl)
	klines, err := HttpGet3(bn.httpClient, klineUrl, nil)
	if err != nil {
		return nil, err
	}
	var klineRecords []Kline

	for _, _record := range klines {
		r := Kline{Pair: currency}
		record := _record.([]interface{})
		for i, e := range record {
			switch i {
			case 0:
				r.Timestamp = int64(e.(float64)) / 1000 //to unix timestramp
			case 1:
				r.Open = ToFloat64(e)
			case 2:
				r.High = ToFloat64(e)
			case 3:
				r.Low = ToFloat64(e)
			case 4:
				r.Close = ToFloat64(e)
			case 5:
				r.Vol = ToFloat64(e)
			}
		}
		klineRecords = append(klineRecords, r)
	}

	return klineRecords, nil

}

//非个人，整个交易所的交易记录
func (bn *Binance) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implements")
}

func (bn *Binance) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implements")
}
func (ba *Binance) adaptCurrencyPair(pair CurrencyPair) CurrencyPair {
	if pair.CurrencyA.Eq(BCH) || pair.CurrencyA.Eq(BCC) {
		return NewCurrencyPair(NewCurrency("BCHABC", ""), pair.CurrencyB).AdaptUsdToUsdt()
	}

	if pair.CurrencyA.Symbol == "BSV" {
		return NewCurrencyPair(NewCurrency("BCHSV", ""), pair.CurrencyB).AdaptUsdToUsdt()
	}

	return pair.AdaptUsdToUsdt()
}
