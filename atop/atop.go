package atop

/**

LimitBuy(amount, price string, currency CurrencyPair) (*Order, error)
LimitSell(amount, price string, currency CurrencyPair) (*Order, error)
MarketBuy(amount, price string, currency CurrencyPair) (*Order, error)
MarketSell(amount, price string, currency CurrencyPair) (*Order, error)
CancelOrder(orderId string, currency CurrencyPair) (bool, error)
GetOneOrder(orderId string, currency CurrencyPair) (*Order, error)
GetUnfinishOrders(currency CurrencyPair) ([]Order, error)
GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error)
GetAccount() (*Account, error)

GetTicker(currency CurrencyPair) (*Ticker, error)
GetDepth(size int, currency CurrencyPair) (*Depth, error)
GetKlineRecords(currency CurrencyPair, period , size, since int) ([]Kline, error)
//Non-individual, transaction record of the entire exchange
GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error)

GetExchangeName() string
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"

	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	//test  https://testapi.a.top
	//product   https://api.a.top
	ApiBaseUrl = "https://testapi.a.top"
	//ApiBaseUrl = "https://api.a.top"

	////market data

	//Trading market configuration
	GetMarketConfig = "/data/api/v1/getMarketConfig"

	//K line data
	GetKLine = "/data/api/v1/getKLine"

	//Aggregate market
	GetTicker = "/data/api/v1/getTicker?market=%s"

	//The latest Ticker for all markets
	GetTickers = "/data/api/v1/getTickers"

	//Market depth data
	GetDepth = "/data/api/v1/getDepth?market=%s"

	//Recent market record
	GetTrades = "/data/api/v1/getTrades"

	////trading

	//Get server time (no signature required)
	GetServerTime = "/trade/api/v1/getServerTime"

	//Get atcount balance
	GetBalance = "/trade/api/v1/getBalance"

	//Plate the order
	PlateOrder = "/trade/api/v1/order"

	//Commissioned by batch
	BatchOrder = "/trade/api/v1/batchOrder"

	//cancellations
	CancelOrder = "/trade/api/v1/cancel"

	//From a single batch
	BatchCancel = "/trade/api/v1/batchCancel"

	//The order information
	GetOrder = "/trade/api/v1/getOrder"

	//Gets an outstanding order
	GetOpenOrders = "/trade/api/v1/getOpenOrders"

	//Get orders history
	GetHistorys = "/trade/api/v1/getHistorys"

	//Gets multiple order information
	GetBatchOrders = "/trade/api/v1/getBatchOrders"

	//Gets the recharge address
	GetPayInAddress = "/trade/api/v1/getPayInAddress"

	//Get the withdrawal address
	GetPayOutAddress = "/trade/api/v1/getPayOutAddress"

	//Gets the recharge record
	GetPayInRecord = "/trade/api/v1/getPayInRecord"

	//Get the withdrawal record
	GetPayOutRecord = "/trade/api/v1/getPayOutRecord"

	//Withdrawal configuration
	GetWithdrawConfig = "/trade/api/v1/getWithdrawConfig"

	//withdraw
	Withdrawal = "/trade/api/v1/withdraw"
)

var KlinePeriodConverter = map[int]string{
	KLINE_PERIOD_1MIN:   "1min",
	KLINE_PERIOD_3MIN:   "3min",
	KLINE_PERIOD_5MIN:   "5min",
	KLINE_PERIOD_15MIN:  "15min",
	KLINE_PERIOD_30MIN:  "30min",
	KLINE_PERIOD_60MIN:  "1hour",
	KLINE_PERIOD_1H:     "1hour",
	KLINE_PERIOD_2H:     "2hour",
	KLINE_PERIOD_4H:     "4hour",
	KLINE_PERIOD_6H:     "6hour",
	KLINE_PERIOD_8H:     "8hour",
	KLINE_PERIOD_12H:    "12hour",
	KLINE_PERIOD_1DAY:   "1day",
	KLINE_PERIOD_3DAY:   "3day",
	KLINE_PERIOD_1WEEK:  "7day",
	KLINE_PERIOD_1MONTH: "30day",
}

type Atop struct {
	accessKey,
	secretKey string
	httpClient *http.Client
}

//hao
func (at *Atop) buildPostForm(postForm *url.Values) error {
	postForm.Set("accesskey", at.accessKey)
	nonce := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	postForm.Set("nonce", nonce)
	payload := postForm.Encode()
	//fmt.Println("payload", payload)
	sign, _ := GetParamHmacSHA256Sign(at.secretKey, payload)
	postForm.Set("signature", sign)
	return nil
}

func New(client *http.Client, apiKey, secretKey string) *Atop {
	return &Atop{apiKey, secretKey, client}
}

func (at *Atop) GetExchangeName() string {
	return "atop.com"
}

//hao
func (at *Atop) GetTicker(currency CurrencyPair) (*Ticker, error) {
	market := strings.ToLower(currency.String())
	tickerUrl := ApiBaseUrl + fmt.Sprintf(GetTicker, market)
	resp, err := HttpGet(at.httpClient, tickerUrl)
	if err != nil {
		return nil, err
	}
	respMap := resp
	var ticker Ticker
	ticker.Pair = currency
	ticker.Date = uint64(time.Now().Unix())
	ticker.Last = ToFloat64(respMap["price"])
	ticker.Buy = ToFloat64(respMap["bid"])
	ticker.Sell = ToFloat64(respMap["ask"])
	ticker.Low = ToFloat64(respMap["low"])
	ticker.High = ToFloat64(respMap["high"])
	ticker.Vol = ToFloat64(respMap["coinVol"])
	return &ticker, nil
}

//hao
func (at *Atop) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	market := strings.ToLower(currency.String())
	depthUrl := ApiBaseUrl + fmt.Sprintf(GetDepth, market)
	resp, err := HttpGet(at.httpClient, depthUrl)
	if err != nil {
		return nil, err
	}
	respMap := resp

	bids := respMap["bids"].([]interface{})
	asks := respMap["asks"].([]interface{})

	depth := new(Depth)
	depth.Pair = currency
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
	sort.Sort(depth.AskList)
	return depth, nil
}

//hao
func (at *Atop) plateOrder(amount, price string, currencyPair CurrencyPair, orderType, orderSide string) (*Order, error) {
	pair := at.adaptCurrencyPair(currencyPair)
	path := ApiBaseUrl + PlateOrder
	params := url.Values{}
	params.Set("market", pair.ToLower().String()) //btc_usdt eth_usdt
	if orderSide == "buy" {
		params.Set("type", strconv.Itoa(1))
	} else {
		params.Set("type", strconv.Itoa(0))
	}
	//params.Set("type", orderSide)//Transaction Type  1、buy 0、sell
	params.Set("price", price)
	params.Set("number", amount)
	if orderType == "market" {
		params.Set("entrustType", strconv.Itoa(1))
	} else {
		params.Set("entrustType", strconv.Itoa(0))
	}
	//params.Set("entrustType", orderType)//Delegate type  0、limit，1、market
	at.buildPostForm(&params)
	resp, err := HttpPostForm(at.httpClient, path, params)
	//log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return nil, err
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		log.Println(string(resp))
		return nil, err
	}

	code := respMap["code"].(float64)
	if code != 200 {
		return nil, errors.New(respMap["info"].(string))
	}

	//return &Order{}, nil
	data := respMap["data"].(map[string]interface{})

	orderId := data["id"].(float64)
	side := BUY
	if orderSide == "sale" {
		side = SELL
	}

	return &Order{
		Currency: pair,
		//OrderID:
		OrderID2:   strconv.FormatFloat(orderId, 'f', 0, 64),
		Price:      ToFloat64(price),
		Amount:     ToFloat64(amount),
		DealAmount: 0,
		AvgPrice:   0,
		Side:       TradeSide(side),
		Status:     ORDER_UNFINISH,
		OrderTime:  int(time.Now().Unix())}, nil
}

//hao
func (at *Atop) GetAccount() (*Account, error) {
	params := url.Values{}
	at.buildPostForm(&params)
	path := ApiBaseUrl + GetBalance
	//fmt.Println("GetBalance", path)
	resp, err := HttpPostForm(at.httpClient, path, params)
	if err != nil {
		return nil, err
	}
	respMap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return nil, err
	}
	data := respMap["data"].(map[string]interface{})
	atc := Account{}
	atc.Exchange = at.GetExchangeName()
	atc.SubAccounts = make(map[Currency]SubAccount)
	for k, v := range data {
		cur := NewCurrency(k, "")
		vv := v.(map[string]interface{})
		sub := SubAccount{}
		sub.Currency = cur
		sub.Amount = ToFloat64(vv["available"]) + ToFloat64(vv["freeze"])
		sub.ForzenAmount = ToFloat64(vv["freeze"])
		atc.SubAccounts[cur] = sub
	}
	return &atc, nil
}

//hao
func (at *Atop) LimitBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return at.plateOrder(amount, price, currencyPair, "limit", "buy")
}

//hao
func (at *Atop) LimitSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return at.plateOrder(amount, price, currencyPair, "limit", "sale")
}

//hao
func (at *Atop) MarketBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return at.plateOrder(amount, price, currencyPair, "market", "buy")
}

//hao
func (at *Atop) MarketSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return at.plateOrder(amount, price, currencyPair, "market", "sale")
}

func (at *Atop) CancelOrder(orderId string, currencyPair CurrencyPair) (bool, error) {
	currencyPair = at.adaptCurrencyPair(currencyPair)
	path := ApiBaseUrl + CancelOrder
	params := url.Values{}
	params.Set("api_key", at.accessKey)
	params.Set("market", currencyPair.ToLower().String())
	params.Set("id", orderId)

	at.buildPostForm(&params)

	resp, err := HttpPostForm(at.httpClient, path, params)

	if err != nil {
		return false, err
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		log.Println(string(resp))
		return false, err
	}
	code := respMap["code"].(float64)
	if code != 200 {
		return false, errors.New(respMap["info"].(string))
	}

	//orderIdCanceled := ToInt(respmap["orderId"])
	//if orderIdCanceled <= 0 {
	//	return false, errors.New(string(resp))
	//}

	return true, nil
}

//hao？
func (at *Atop) GetOneOrder(orderId string, currencyPair CurrencyPair) (*Order, error) {
	currencyPair = at.adaptCurrencyPair(currencyPair)
	path := ApiBaseUrl + GetOrder
	log.Println(path)
	params := url.Values{}
	params.Set("api_key", at.accessKey)
	params.Set("market", currencyPair.ToLower().String())
	params.Set("id", orderId)
	at.buildPostForm(&params)
	resp, err := HttpPostForm(at.httpClient, path, params)

	if err != nil {
		return nil, err
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		log.Println(string(resp))
		return nil, err
	}
	code := respMap["code"].(float64)

	if code != 200 {
		return nil, errors.New(respMap["info"].(string))
	}

	data := respMap["data"].(map[string]interface{})

	status := data["status"].(float64)
	side := data["flag"]

	ord := Order{}
	ord.Currency = currencyPair
	//ord.OrderID = ToInt(orderId)
	ord.OrderID2 = orderId

	if side == "sale" {
		ord.Side = SELL
	} else {
		ord.Side = BUY
	}

	switch status {
	case 0:
		ord.Status = ORDER_UNFINISH
	case 1:
		ord.Status = ORDER_PART_FINISH
	case 2:
		ord.Status = ORDER_FINISH
	case 3:
		ord.Status = ORDER_CANCEL
		//case 4:
		//	ord.Status = new(TradeStatus)//settle
		//case "PENDING_CANCEL":
		//	ord.Status = ORDER_CANCEL_ING
		//case "REJECTED":
		//	ord.Status = ORDER_REJECT
	}
	ord.Amount = ToFloat64(data["number"])
	ord.Price = ToFloat64(data["price"])
	ord.DealAmount = ord.Amount - ToFloat64(data["completeNumber"]) //？
	ord.AvgPrice = ToFloat64(data["avg_price"])                     // response no avg price ， fill price

	return &ord, nil
}

//hao
func (at *Atop) GetUnfinishOrders(currencyPair CurrencyPair) ([]Order, error) {
	pair := at.adaptCurrencyPair(currencyPair)
	path := ApiBaseUrl + GetOpenOrders
	params := url.Values{}
	params.Set("market", pair.ToLower().String())
	params.Set("page", "1")
	params.Set("pageSize", "10000")
	at.buildPostForm(&params)

	resp, err := HttpPostForm(at.httpClient, path, params)
	if err != nil {
		return nil, err
	}
	respMap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		log.Println(string(resp))
		return nil, err
	}

	code := respMap["code"].(float64)
	if code != 200 {
		return nil, errors.New(respMap["info"].(string))
	}
	data := respMap["data"].([]interface{})
	orders := make([]Order, 0)
	for _, ord := range data {
		ordData := ord.(map[string]interface{})
		orderId := strconv.FormatFloat(ordData["id"].(float64), 'f', 0, 64)
		orders = append(orders, Order{
			OrderID:   0,
			OrderID2:  orderId,
			Currency:  currencyPair,
			Price:     ToFloat64(ordData["price"]),
			Amount:    ToFloat64(ordData["number"]),
			Side:      TradeSide(ToInt(ordData["type"])),
			Status:    TradeStatus(ToInt(ordData["status"])),
			OrderTime: ToInt(ordData["time"])})
	}
	return orders, nil
}

//hao
func (at *Atop) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	pair := at.adaptCurrencyPair(currency)
	params := url.Values{}
	params.Set("market", pair.ToLower().String())
	//params.Set("type", "1min") //1min,5min,15min,30min,1hour,6hour,1day,7day,30day
	params.Set("type", KlinePeriodConverter[period]) //1min,5min,15min,30min,1hour,6hour,1day,7day,30day
	params.Set("since", fmt.Sprintf("%d", size))     //The first time is 0, followed by the value of the response since

	klineUrl := ApiBaseUrl + GetKLine + "?" + params.Encode()
	kLines, err := HttpGet(at.httpClient, klineUrl)
	if err != nil {
		return nil, err
	}
	var klineRecords []Kline
	for _, _record := range kLines["datas"].([]interface{}) {
		r := Kline{Pair: currency}
		record := _record.([]interface{})
		for i, e := range record {
			switch i {
			case 0:
				r.Timestamp = int64(e.(float64)) //to unix timestramp
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

// hao Non-individual, transaction record of the entire exchange
func (at *Atop) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	pair := at.adaptCurrencyPair(currencyPair)
	params := url.Values{}
	params.Set("market", pair.ToLower().String())

	apiUrl := ApiBaseUrl + GetTrades + "?" + params.Encode()

	resp, err := HttpGet(at.httpClient, apiUrl)
	if err != nil {
		return nil, err
	}

	var trades []Trade
	for _, v := range resp {
		m := v.(map[string]interface{})
		ty := SELL
		if m["isBuyerMaker"].(bool) {
			ty = BUY
		}
		trades = append(trades, Trade{
			Tid:    ToInt64(m["id"]),
			Type:   ty,
			Amount: ToFloat64(m["qty"]),
			Price:  ToFloat64(m["price"]),
			Date:   ToInt64(m["time"]),
			Pair:   currencyPair,
		})
	}

	return trades, nil
}

func (at *Atop) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	//panic("not support")
	pair := at.adaptCurrencyPair(currency)
	path := ApiBaseUrl + GetHistorys
	params := url.Values{}
	params.Set("market", pair.ToLower().String())
	//params.Set("type", "1")
	//params.Set("status", "0")
	params.Set("page", fmt.Sprint(currentPage))
	params.Set("pageSize", fmt.Sprint(pageSize))

	at.buildPostForm(&params)
	resp, err := HttpPostForm(at.httpClient, path, params)
	if err != nil {
		return nil, err
	}
	respMap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		log.Println(string(resp))
		return nil, err
	}

	code := respMap["code"].(float64)
	if code != 200 {
		return nil, errors.New(respMap["info"].(string))
	}
	data := respMap["data"].(map[string]interface{})
	records := data["record"].([]interface{})
	orders := make([]Order, 0)
	for _, ord := range records {
		ordData := ord.(map[string]interface{})
		orderId := strconv.FormatFloat(ordData["id"].(float64), 'f', 0, 64)
		orders = append(orders, Order{
			OrderID:   0,
			OrderID2:  orderId,
			Currency:  currency,
			Price:     ToFloat64(ordData["price"]),
			Amount:    ToFloat64(ordData["number"]),
			Side:      TradeSide(ToInt(ordData["type"])),
			Status:    TradeStatus(ToInt(ordData["status"])),
			OrderTime: ToInt(ordData["time"])})
	}
	return orders, nil

}

// hao
func (at *Atop) Withdraw(amount, memo string, currency Currency, fees, receiveAddr, safePwd string) (string, error) {
	params := url.Values{}
	coin := strings.ToLower(currency.Symbol)
	path := ApiBaseUrl + Withdrawal
	params.Set("coin", coin)
	params.Set("address", receiveAddr)
	params.Set("amount", amount)
	params.Set("receiveAddr", receiveAddr)
	params.Set("safePwd", safePwd)
	//params.Set("memo", memo)
	at.buildPostForm(&params)

	resp, err := HttpPostForm(at.httpClient, path, params)

	if err != nil {
		return "", err
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return "", err
	}

	if respMap["code"].(float64) == 200 {
		return respMap["id"].(string), nil
	}
	return "", errors.New(string(resp))
}

func (at *Atop) CancelWithdraw(id string, currency Currency, safePwd string) (bool, error) {
	panic("not support")
}
func (at *Atop) adaptCurrencyPair(pair CurrencyPair) CurrencyPair {
	return pair
}
