package allcoin

import (
	"encoding/json"
	"errors"
	. "github.com/nntaoli-project/goex"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	API_BASE_URL = "https://www.allcoin.ca/"

	TICKER_URI             = "Api_Market/getCoinTrade"
	TICKER_URI_2           = "Api_Order/ticker"
	TICKERS_URI            = "ticker/allBookTickers"
	DEPTH_URI              = "Api_Order/depth"
	ACCOUNT_URI            = "Api_User/userBalance"
	ORDER_URI              = "Api_Order/coinTrust"
	ORDER_CANCEL_URI       = "Api_Order/cancel"
	ORDER_INFO_URI         = "Api_Order/orderInfo"
	UNFINISHED_ORDERS_INFO = "Api_Order/trustList"
)

type Allcoin struct {
	accessKey,
	secretKey string
	httpClient *http.Client
}

func (ac *Allcoin) buildParamsSigned(postForm *url.Values) error {
	//postForm.Set("api_key", ac.accessKey)
	//postForm.Set("secret_key", ac.secretKey)
	payload := postForm.Encode() + "&secret_key=" + ac.secretKey
	//log.Println("payload:", payload, "postForm:", postForm.Encode())
	sign, _ := GetParamMD5Sign(ac.secretKey, payload)
	postForm.Set("sign", sign)
	return nil
}

func New(client *http.Client, api_key, secret_key string) *Allcoin {
	return &Allcoin{api_key, secret_key, client}
}

func (ac *Allcoin) GetExchangeName() string {
	return "allcoin.com"
}

func (ac *Allcoin) GetTicker(currency CurrencyPair) (*Ticker, error) {
	//wg := sync.WaitGroup{}
	//wg.Add(2)
	//go func() {
	//	defer wg.Done()
	//currency2 := ac.adaptCurrencyPair(currency)
	//params := url.Values{}
	//params.Set("symbol", strings.ToLower(currency2.ToSymbol("2")))
	//path := API_BASE_URL + TICKER_URI_2
	//resp, err := HttpPostForm(ac.httpClient, path, params)
	//log.Println("resp:", string(resp), "err:", err)
	//if err != nil {
	//	//return nil, err
	//}
	//
	//respmap := make(map[string]interface{})
	//err = json.Unmarshal(resp, &respmap)
	//if err != nil {
	//	log.Println(string(resp))
	//	//return nil, err
	//}
	//code := respmap["code"].(float64)
	//msg := respmap["msg"].(string)
	//log.Println("code=", code, "msg:", msg)
	//if code != 0 {
	//	//return nil, errors.New(respmap["msg"].(string))
	//}
	//data := respmap["data"].(map[string]interface{})
	//log.Println("1", data)
	//}()

	//go func() {
	//	defer wg.Done()
	currency2 := ac.adaptCurrencyPair(currency)
	params := url.Values{}
	params.Set("part", strings.ToLower(currency2.CurrencyB.String()))
	params.Set("coin", strings.ToLower(currency2.CurrencyA.String()))
	path := API_BASE_URL + TICKER_URI
	resp, err := HttpPostForm(ac.httpClient, path, params)
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

	//log.Println("2", respmap)
	//}()
	//wg.Wait()
	var ticker Ticker
	ticker.Pair = currency
	ticker.Date = uint64(time.Now().Unix())
	ticker.Last = ToFloat64(respmap["price"])
	ticker.Buy = ToFloat64(respmap["buy"])
	ticker.Sell = ToFloat64(respmap["sale"])
	ticker.Low = ToFloat64(respmap["min"])
	ticker.High = ToFloat64(respmap["max"])
	ticker.Vol = ToFloat64(respmap["volume_24h"])
	return &ticker, nil
}

func (ac *Allcoin) GetDepth(size int, currencyPair CurrencyPair) (*Depth, error) {

	currency2 := ac.adaptCurrencyPair(currencyPair)
	params := url.Values{}
	params.Set("symbol", strings.ToLower(currency2.ToSymbol("2")))
	path := API_BASE_URL + DEPTH_URI
	resp, err := HttpPostForm(ac.httpClient, path, params)
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

func (ac *Allcoin) placeOrder(amount, price string, pair CurrencyPair, orderType, orderSide string) (*Order, error) {
	pair = ac.adaptCurrencyPair(pair)
	path := API_BASE_URL + ORDER_URI
	params := url.Values{}
	params.Set("api_key", ac.accessKey)
	params.Set("symbol", strings.ToLower(pair.ToSymbol("2")))
	params.Set("type", orderSide)

	params.Set("price", price)
	params.Set("number", amount)

	ac.buildParamsSigned(&params)

	resp, err := HttpPostForm(ac.httpClient, path, params)
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

func (ac *Allcoin) GetAccount() (*Account, error) {
	params := url.Values{}

	params.Set("api_key", ac.accessKey)
	ac.buildParamsSigned(&params)
	//log.Println("params=", params)
	path := API_BASE_URL + ACCOUNT_URI
	resp, err := HttpPostForm(ac.httpClient, path, params)
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

	acc := Account{}
	acc.Exchange = ac.GetExchangeName()
	acc.SubAccounts = make(map[Currency]SubAccount)

	for k, v := range data {
		s := strings.Split(k, "_")
		if len(s) == 2 {
			cur := NewCurrency(s[0], "")
			if s[1] == "over" {
				sub := SubAccount{}
				sub = acc.SubAccounts[cur]
				sub.Amount = ToFloat64(v)
				acc.SubAccounts[cur] = sub
			} else if s[1] == "lock" {
				sub := SubAccount{}
				sub = acc.SubAccounts[cur]
				sub.ForzenAmount = ToFloat64(v)
				acc.SubAccounts[cur] = sub
			}
		}
	}
	log.Println(acc)
	return &acc, nil
}

func (ac *Allcoin) LimitBuy(amount, price string, currencyPair CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	return ac.placeOrder(amount, price, currencyPair, "LIMIT", "buy")
}

func (ac *Allcoin) LimitSell(amount, price string, currencyPair CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	return ac.placeOrder(amount, price, currencyPair, "LIMIT", "sale")
}

func (ac *Allcoin) MarketBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return ac.placeOrder(amount, price, currencyPair, "MARKET", "buy")
}

func (ac *Allcoin) MarketSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return ac.placeOrder(amount, price, currencyPair, "MARKET", "sale")
}

func (ac *Allcoin) CancelOrder(orderId string, currencyPair CurrencyPair) (bool, error) {
	currencyPair = ac.adaptCurrencyPair(currencyPair)
	path := API_BASE_URL + ORDER_CANCEL_URI
	params := url.Values{}
	params.Set("api_key", ac.accessKey)
	params.Set("symbol", strings.ToLower(currencyPair.ToSymbol("2")))
	params.Set("order_id", orderId)

	ac.buildParamsSigned(&params)

	resp, err := HttpPostForm(ac.httpClient, path, params)

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

func (ac *Allcoin) GetOneOrder(orderId string, currencyPair CurrencyPair) (*Order, error) {
	currencyPair = ac.adaptCurrencyPair(currencyPair)
	path := API_BASE_URL + ORDER_INFO_URI
	params := url.Values{}
	params.Set("api_key", ac.accessKey)
	params.Set("symbol", strings.ToLower(currencyPair.ToSymbol("2")))
	params.Set("trust_id", orderId)

	ac.buildParamsSigned(&params)

	resp, err := HttpPostForm(ac.httpClient, path, params)

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

func (ac *Allcoin) GetUnfinishOrders(currencyPair CurrencyPair) ([]Order, error) {
	currencyPair = ac.adaptCurrencyPair(currencyPair)
	path := API_BASE_URL + UNFINISHED_ORDERS_INFO
	params := url.Values{}
	params.Set("api_key", ac.accessKey)
	params.Set("symbol", strings.ToLower(currencyPair.ToSymbol("2")))
	params.Set("type", "open")

	ac.buildParamsSigned(&params)

	resp, err := HttpPostForm(ac.httpClient, path, params)

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
		//ord := v.(map[string]interface{})

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

func (ac *Allcoin) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implements")
}

//非个人，整个交易所的交易记录
func (ac *Allcoin) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implements")
}

func (ac *Allcoin) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implements")
}
func (ba *Allcoin) adaptCurrencyPair(pair CurrencyPair) CurrencyPair {
	return pair
}
