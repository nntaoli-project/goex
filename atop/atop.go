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
//非个人，整个交易所的交易记录
GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error)

GetExchangeName() string
*/

import (
	"encoding/json"
	"errors"
	. "github.com/nntaoli-project/GoEx"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	//test  https://testapi.a.top
	//product   https://api.a.top
	ApiBaseUrl = "https://api.a.top"

	////market data

	//Trading market configuration
	GetMarketConfig = "/data/api/v1/getMarketConfig"

	//K line data
	GetKLine = "/data/api/v1/getKLine"

	//Aggregate market
	GetTicker = "/data/api/v1/getTicker"

	//The latest Ticker for all markets
	GetTickers = "/data/api/v1/getTickers"

	//Market depth data
	GetDepth = "/data/api/v1/getDepth"

	//Recent market record
	GetTrades = "/data/api/v1/getTrades"

	////trading

	//Get server time (no signature required)
	GetServerTime = "/trade/api/v1/getServerTime"

	//Get atcount balance
	GetBalance = "trade/api/v1/getBalance"

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
	Withdrawal = "trade/api/v1/withdraw"
)

type Atop struct {
	accessKey,
	secretKey string
	httpClient *http.Client
}

func (at *Atop) buildParamsSigned(postForm *url.Values) error {
	//postForm.Set("api_key", at.atcessKey)
	//postForm.Set("secret_key", at.secretKey)
	payload := postForm.Encode() + "&secret_key=" + at.secretKey
	//log.Println("payload:", payload, "postForm:", postForm.Encode())
	sign, _ := GetParamMD5Sign(at.secretKey, payload)
	postForm.Set("sign", sign)
	return nil
}

func New(client *http.Client, apiKey, secretKey string) *Atop {
	return &Atop{apiKey, secretKey, client}
}

func (at *Atop) GetExchangeName() string {
	return "atop.com"
}

func (at *Atop) GetTicker(currency CurrencyPair) (*Ticker, error) {
	currencyPair := at.adaptCurrencyPair(currency)
	params := url.Values{}
	params.Set("market", strings.ToLower(currencyPair.CurrencyA.String()))
	path := ApiBaseUrl + GetTicker
	resp, err := HttpGet(at.httpClient, path)
	//log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return nil, err
	}

	var ticker Ticker
	ticker.Pair = currency
	ticker.Date = uint64(time.Now().Unix())
	ticker.Last = ToFloat64(resp["price"])
	ticker.Buy = ToFloat64(resp["bid"])
	ticker.Sell = ToFloat64(resp["ask"])
	ticker.Low = ToFloat64(resp["low"])
	ticker.High = ToFloat64(resp["high"])
	ticker.Vol = ToFloat64(resp["coinVol"])
	return &ticker, nil
}

func (at *Atop) GetDepth(size int, currencyPair CurrencyPair) (*Depth, error) {

	currency2 := at.adaptCurrencyPair(currencyPair)
	params := url.Values{}
	params.Set("market", strings.ToLower(currency2.ToSymbol("2")))
	path := ApiBaseUrl + GetDepth
	resp, err := HttpPostForm(at.httpClient, path, params)
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
	code := respmap["code"].(float64)
	msg := respmap["msg"].(string)
	log.Println("code=", code, "msg:", msg)
	if code != 0 {
		return nil, errors.New(respmap["msg"].(string))
	}
	data := respmap["data"].(map[string]interface{})
	log.Println("1", data)

	bids := data["bids"].([]interface{})
	asks := data["asks"].([]interface{})

	//log.Println("len bids", len(bids))
	//log.Println("len asks", len(asks))
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
	sort.Sort(depth.AskList)
	return depth, nil
}

func (at *Atop) plateOrder(amount, price string, pair CurrencyPair, orderType, orderSide string) (*Order, error) {
	pair = at.adaptCurrencyPair(pair)
	path := ApiBaseUrl + PlateOrder
	params := url.Values{}
	params.Set("api_key", at.accessKey)
	params.Set("symbol", strings.ToLower(pair.ToSymbol("2")))
	params.Set("type", orderSide)

	params.Set("price", price)
	params.Set("number", amount)

	at.buildParamsSigned(&params)

	resp, err := HttpPostForm(at.httpClient, path, params)
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
	code := respmap["code"].(float64)
	if code != 0 {
		return nil, errors.New(respmap["msg"].(string))
	}
	data := respmap["data"].(map[string]interface{})

	orderId := data["order_id"].(string)
	side := BUY
	if orderSide == "sale" {
		side = SELL
	}

	return &Order{
		Currency: pair,
		//OrderID:
		OrderID2:   orderId,
		Price:      ToFloat64(price),
		Amount:     ToFloat64(amount),
		DealAmount: 0,
		AvgPrice:   0,
		Side:       TradeSide(side),
		Status:     ORDER_UNFINISH,
		OrderTime:  int(time.Now().Unix())}, nil
}

func (at *Atop) GetAccount() (*Account, error) {
	params := url.Values{}

	params.Set("api_key", at.accessKey)
	at.buildParamsSigned(&params)
	//log.Println("params=", params)
	path := ApiBaseUrl + GetBalance
	resp, err := HttpPostForm(at.httpClient, path, params)
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
	code := respmap["code"].(float64)
	//msg := respmap["msg"].(string)
	//log.Println("code=", code, "msg:", msg)
	if code != 0 {
		return nil, errors.New(respmap["msg"].(string))
	}
	data := respmap["data"].(map[string]interface{})

	atc := Account{}
	atc.Exchange = at.GetExchangeName()
	//map[Currency]Subatcount
	atc.SubAccounts = make(map[Currency]SubAccount)

	for k, v := range data {
		s := strings.Split(k, "_")
		if len(s) == 2 {
			cur := NewCurrency(s[0], "")
			if s[1] == "over" {
				sub := SubAccount{}
				sub = atc.SubAccounts[cur]
				sub.Amount = ToFloat64(v)
				atc.SubAccounts[cur] = sub
			} else if s[1] == "lock" {
				sub := SubAccount{}
				sub = atc.SubAccounts[cur]
				sub.ForzenAmount = ToFloat64(v)
				atc.SubAccounts[cur] = sub
			}
		}
	}
	log.Println(atc)
	return &atc, nil
}

func (at *Atop) LimitBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return at.plateOrder(amount, price, currencyPair, "LIMIT", "buy")
}

func (at *Atop) LimitSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return at.plateOrder(amount, price, currencyPair, "LIMIT", "sale")
}

func (at *Atop) MarketBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return at.plateOrder(amount, price, currencyPair, "MARKET", "buy")
}

func (at *Atop) MarketSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return at.plateOrder(amount, price, currencyPair, "MARKET", "sale")
}

func (at *Atop) CancelOrder(orderId string, currencyPair CurrencyPair) (bool, error) {
	currencyPair = at.adaptCurrencyPair(currencyPair)
	path := ApiBaseUrl + CancelOrder
	params := url.Values{}
	params.Set("api_key", at.accessKey)
	params.Set("symbol", strings.ToLower(currencyPair.ToSymbol("2")))
	params.Set("order_id", orderId)

	at.buildParamsSigned(&params)

	resp, err := HttpPostForm(at.httpClient, path, params)

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
	code := respmap["code"].(int)
	if code != 0 {
		return false, errors.New(respmap["msg"].(string))
	}

	//orderIdCanceled := ToInt(respmap["orderId"])
	//if orderIdCanceled <= 0 {
	//	return false, errors.New(string(resp))
	//}

	return true, nil
}

func (at *Atop) GetOneOrder(orderId string, currencyPair CurrencyPair) (*Order, error) {
	currencyPair = at.adaptCurrencyPair(currencyPair)
	path := ApiBaseUrl + GetOrder
	params := url.Values{}
	params.Set("api_key", at.accessKey)
	params.Set("symbol", strings.ToLower(currencyPair.ToSymbol("2")))
	params.Set("trust_id", orderId)

	at.buildParamsSigned(&params)

	resp, err := HttpPostForm(at.httpClient, path, params)

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
	code := respmap["code"].(float64)
	if code != 0 {
		return nil, errors.New(respmap["msg"].(string))
	}

	data := respmap["data"].(map[string]interface{})

	status := data["status"]
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
	case "1": ///////////////////////////////////////////////////////////////////////////////////#################TODO
		ord.Status = ORDER_FINISH
	case "2":
		ord.Status = ORDER_PART_FINISH
	case "3":
		ord.Status = ORDER_CANCEL
	case "PENDING_CANCEL":
		ord.Status = ORDER_CANCEL_ING
	case "REJECTED":
		ord.Status = ORDER_REJECT
	}

	ord.Amount = ToFloat64(data["number"])
	ord.Price = ToFloat64(data["price"])
	ord.DealAmount = ord.Amount - ToFloat64(data["numberover"])
	ord.AvgPrice = ToFloat64(data["avg_price"]) // response no avg price ， fill price

	return &ord, nil
}

func (at *Atop) GetUnfinishOrders(currencyPair CurrencyPair) ([]Order, error) {
	currencyPair = at.adaptCurrencyPair(currencyPair)
	path := ApiBaseUrl + GetOpenOrders
	params := url.Values{}
	params.Set("api_key", at.accessKey)
	params.Set("symbol", strings.ToLower(currencyPair.ToSymbol("2")))
	params.Set("type", "open")

	at.buildParamsSigned(&params)

	resp, err := HttpPostForm(at.httpClient, path, params)

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
	code := respmap["code"].(float64)
	//msg := respmap["msg"].(string)
	//log.Println("code=", code, "msg:", msg)
	if code != 0 {
		return nil, errors.New(respmap["msg"].(string))
	}
	data, isok := respmap["data"].([]map[string]interface{})

	orders := make([]Order, 0)
	if isok != true {
		return orders, nil
	}
	for _, ord := range data {
		//ord := v.(map[string]interfate{})

		//side := ord["side"].(string)
		//orderSide := SELL
		//if side == "BUY" {
		//	orderSide = BUY
		//}

		orders = append(orders, Order{
			OrderID:  ToInt(ord["id"]),
			OrderID2: ord["id"].(string),
			Currency: currencyPair,
			Price:    ToFloat64(ord["price"]),
			Amount:   ToFloat64(ord["number"]),
			//Side:      TradeSide(orderSide),
			//Status:    ORDER_UNFINISH,
			OrderTime: ToInt(ord["created"])})
	}
	return orders, nil
}

func (at *Atop) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not support")
}

//非个人，整个交易所的交易记录
func (at *Atop) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not support")
}

func (at *Atop) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not support")
}
func (at *Atop) adaptCurrencyPair(pair CurrencyPair) CurrencyPair {
	return pair
}
