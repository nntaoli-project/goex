package coinpark

import (
	"net/http"
	//"log"
	"fmt"
	. "github.com/nntaoli-project/goex"
	//"net/url"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	API_BASE_URL = "https://api.coinpark.cc/"
	V1           = "v1/"
	API_URL      = API_BASE_URL + V1
	TICKER_API   = "mdata?cmd=ticker&pair=%s"
	ALL_TICKERs  = "mdata?cmd=marketAll"
	Pair_List    = "mdata?cmd=pairList"
	DEPTH_API    = "mdata?cmd=depth&pair=%s&size=%d"

	TRADE_URL       = "orderpending"
	GET_ACCOUNT_API = "transfer"
	GET_ORDER_API   = "orders/%s"
	//GET_ORDERS_LIST_API             = ""
	GET_UNFINISHED_ORDERS_API = "getUnfinishedOrdersIgnoreTradeType"
	CANCEL_ORDER_API          = TRADE_URL + "/%s/submit-cancel"
	PLACE_ORDER_API           = "order"
	WITHDRAW_API              = "withdraw"
	CANCELWITHDRAW_API        = "cancelWithdraw"
)

type Cpk struct {
	httpClient *http.Client
	accessKey,
	secretKey string
}

func New(client *http.Client, apikey, secretkey string) *Cpk {
	return &Cpk{accessKey: apikey, secretKey: secretkey, httpClient: client}
}

func (c *Cpk) buildSigned(cmd string) string {
	signed, _ := GetParamHmacMD5Sign(c.secretKey, cmd)
	return signed
}

func (c *Cpk) GetExchangeName() string {
	return "coinpark.cc"
}

func (c *Cpk) GetServerTime() error {
	return nil
}

func (c *Cpk) GetTicker(currencyPair CurrencyPair) (*Ticker, error) {
	url := API_URL + fmt.Sprintf(TICKER_API, currencyPair.String())
	respmap, err := HttpGet(c.httpClient, url)
	if err != nil {
		return nil, err
	}
	//log.Println("ticker respmap:", respmap)
	errcode, isok := respmap["error"].(map[string]interface{})
	if isok == true {
		return nil, errors.New(errcode["msg"].(string))
	}

	tickmap, ok := respmap["result"].(map[string]interface{})
	if !ok {
		return nil, errors.New("tick assert error")
	}
	ticker := new(Ticker)
	ticker.Pair = currencyPair
	ticker.Date = uint64(time.Now().Nanosecond() / 1000)
	ticker.Last = ToFloat64(tickmap["last"])
	ticker.Vol = ToFloat64(tickmap["vol"])
	ticker.Low = ToFloat64(tickmap["low"])
	ticker.High = ToFloat64(tickmap["high"])
	ticker.Buy = ToFloat64(tickmap["buy"])
	ticker.Sell = ToFloat64(tickmap["sell"])
	return ticker, nil

}

func (c *Cpk) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	url := API_URL + fmt.Sprintf(DEPTH_API, currency.String(), size)
	respmap, err := HttpGet(c.httpClient, url)
	if err != nil {
		return nil, err
	}

	errcode, isok := respmap["error"].(map[string]interface{})
	if isok == true {
		return nil, errors.New(errcode["msg"].(string))
	}

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

func (c *Cpk) placeOrder(orderType, orderSide, amount, price string, pair CurrencyPair) (*Order, error) {
	path := API_URL + TRADE_URL
	params := make(map[string]interface{})
	params["cmd"] = "orderpending/trade"
	params["index"] = strconv.Itoa(rand.Intn(1000))

	body := make(map[string]interface{})
	body["pair"] = pair.String()
	body["account_type"] = "0"
	if orderType == "limit" {
		body["order_type"] = "2"
	} else if orderType == "market" {
		body["order_type"] = "1"
		body["money"] = ToFloat64(price) * ToFloat64(amount)
	}
	if orderSide == "buy" {
		body["order_side"] = "1"
	} else if orderSide == "sell" {
		body["order_side"] = "2"
	}
	body["price"] = price
	body["amount"] = amount
	params["body"] = body

	cmd, _ := json.Marshal(params)
	cmds := "[" + string(cmd) + "]"
	sign := c.buildSigned(cmds)

	param := make(map[string]string)
	//param["cmds"] =  strconv.Quote("["+string(cmds)+"]")
	param["cmds"] = cmds
	param["apikey"] = c.accessKey
	param["sign"] = sign
	resp, err := HttpPostForm4(c.httpClient, path, param, nil)
	if err != nil {
		return nil, err
	}
	respmap := make(map[string]interface{})
	json.Unmarshal(resp, &respmap)

	errcode, isok := respmap["error"].(map[string]interface{})
	if isok == true {
		return nil, errors.New(errcode["msg"].(string))
	}

	orderId := ToInt(respmap["result"])
	if orderId <= 0 {
		return nil, errors.New(string(resp))
	}

	side := BUY
	if orderSide == "sell" {
		side = SELL
	}

	return &Order{
		Currency:   pair,
		OrderID:    ToInt(respmap["index"]),
		OrderID2:   respmap["result"].(string),
		Price:      ToFloat64(price),
		Amount:     ToFloat64(amount),
		DealAmount: 0,
		AvgPrice:   0,
		Side:       TradeSide(side),
		Status:     ORDER_UNFINISH,
		OrderTime:  int(time.Now().Unix())}, nil
}

func (c *Cpk) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return c.placeOrder("limit", "buy", amount, price, currency)
}

func (c *Cpk) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return c.placeOrder("limit", "sell", amount, price, currency)
}

func (c *Cpk) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return c.placeOrder("market", "buy", amount, price, currency)
}
func (c *Cpk) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return c.placeOrder("market", "sell", amount, price, currency)
}

func (c *Cpk) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	path := API_URL + TRADE_URL
	params := make(map[string]interface{})
	params["cmd"] = "orderpending/cancelTrade"
	params["index"] = strconv.Itoa(rand.Intn(1000))

	body := make(map[string]interface{})
	body["orders_id"] = orderId
	params["body"] = body

	cmd, _ := json.Marshal(params)
	cmds := "[" + string(cmd) + "]"
	sign := c.buildSigned(cmds)

	param := make(map[string]string)
	//param["cmds"] =  strconv.Quote("["+string(cmds)+"]")
	param["cmds"] = cmds
	param["apikey"] = c.accessKey
	param["sign"] = sign
	resp, err := HttpPostForm4(c.httpClient, path, param, nil)
	if err != nil {
		return false, err
	}
	respmap := make(map[string]interface{})
	json.Unmarshal(resp, &respmap)

	errcode, isok := respmap["error"].(map[string]interface{})
	if isok == true {
		return false, errors.New(errcode["msg"].(string))
	}

	status := respmap["result"].(string)
	if strings.Contains(status, "撤销中") {
		return true, nil
	}
	return false, errors.New("fail")
}

func (c *Cpk) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	path := API_URL + TRADE_URL
	params := make(map[string]interface{})
	params["cmd"] = "orderpending/order"

	body := make(map[string]interface{})
	body["id"] = orderId
	params["body"] = body

	cmd, _ := json.Marshal(params)
	cmds := "[" + string(cmd) + "]"
	sign := c.buildSigned(cmds)

	param := make(map[string]string)
	//param["cmds"] =  strconv.Quote("["+string(cmds)+"]")
	param["cmds"] = cmds
	param["apikey"] = c.accessKey
	param["sign"] = sign
	resp, err := HttpPostForm4(c.httpClient, path, param, nil)
	if err != nil {
		return nil, err
	}
	respmap := make(map[string]interface{})
	json.Unmarshal(resp, &respmap)

	errcode, isok := respmap["error"].(map[string]interface{})
	if isok == true {
		return nil, errors.New(errcode["msg"].(string))
	}

	orderInfo := respmap["result"].(map[string]interface{})
	status := ToInt(orderInfo["status"])
	var orderState TradeStatus
	switch status {
	case 1:
		orderState = ORDER_UNFINISH
	case 2:
		orderState = ORDER_PART_FINISH
	case 3:
		orderState = ORDER_FINISH
	case 4, 5:
		orderState = ORDER_CANCEL
	case 6:
		orderState = ORDER_CANCEL_ING
	}
	var side TradeSide
	switch ToInt(orderInfo["order_side"]) {
	case 1:
		side = BUY
	case 2:
		side = SELL
	}
	return &Order{
		Currency:   currency,
		OrderID2:   orderInfo["id"].(string),
		Price:      ToFloat64(orderInfo["price"]),
		Amount:     ToFloat64(orderInfo["amount"]),
		DealAmount: ToFloat64(orderInfo["deal_amount"]),
		AvgPrice:   0,
		Side:       side,
		Status:     orderState,
		OrderTime:  int(time.Now().Unix())}, nil

}

func (c *Cpk) GetOrdersList() {
	panic("not implement")
}

func (c *Cpk) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implement")
}

func (c *Cpk) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implement")
}

func (c *Cpk) GetAccount() (*Account, error) {
	path := API_URL + GET_ACCOUNT_API
	//cmds := "[{\"cmd\":\"transfer/assets\",\"body\":{\"select\":1}}]"

	params := make(map[string]interface{})
	params["cmd"] = "transfer/assets"

	body := make(map[string]interface{})
	body["select"] = 1

	params["body"] = body

	cmd, _ := json.Marshal(params)
	cmds := "[" + string(cmd) + "]"
	sign := c.buildSigned(cmds)

	param := make(map[string]string)
	param["cmds"] = cmds
	param["apikey"] = c.accessKey
	param["sign"] = sign
	resp, err := HttpPostForm4(c.httpClient, path, param, nil)

	//cmds := `[{cmd:transfer/assets,body:{}}]`
	//cmds := "[{\"cmd\":\"transfer/assets\",\"body\":{\"select\":1}}]"
	//sign := c.buildSigned(cmds)

	//params := make(map[string]string)
	//params["cmds"] =  cmds
	//params["apikey"]= c.accessKey
	//params["sign"]= sign

	//resp, err := HttpPostForm4(c.httpClient, path, params, nil)

	if err != nil {
		return nil, err
	}
	respmap := make(map[string]interface{})
	json.Unmarshal(resp, &respmap)
	log.Println(string(resp))
	//log.Println(respmap)

	errcode, isok := respmap["error"].(map[string]interface{})
	if isok == true {
		return nil, errors.New(errcode["msg"].(string))
	}
	ba := respmap["result"].([]interface{})
	ba1 := ba[0].(map[string]interface{})
	balances := ba1["result"].(map[string]interface{})
	acc := new(Account)

	acc.Asset = ToFloat64(balances["total_btc"])
	acc.SubAccounts = make(map[Currency]SubAccount)
	acc.Exchange = c.GetExchangeName()
	assets_list, isok := balances["assets_list"].([]interface{})
	if isok != true {
		return acc, nil
	}
	for _, v := range assets_list {
		//log.Println(v)
		vv := v.(map[string]string)
		currency := NewCurrency(vv["coin_symbol"], "")
		acc.SubAccounts[currency] = SubAccount{
			Currency:     currency,
			Amount:       ToFloat64(vv["balance"]),
			ForzenAmount: ToFloat64(vv["freeze"]),
		}
	}

	return acc, nil

}

func (c *Cpk) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录
func (c *Cpk) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}

func (c *Cpk) GetPairList() ([]CurrencyPair, error) {
	url := API_URL + Pair_List
	respmap, err := HttpGet(c.httpClient, url)
	if err != nil {
		return nil, err
	}
	log.Println("respmap:", respmap)
	errcode, isok := respmap["error"].(map[string]interface{})
	if isok == true {
		return nil, errors.New(errcode["msg"].(string))
	}

	list, ok := respmap["result"].([]interface{})
	if !ok {
		return nil, errors.New("tick assert error")
	}
	pairlist := make([]CurrencyPair, 0)
	for _, v := range list {
		vv := v.(map[string]interface{})
		pairlist = append(pairlist, NewCurrencyPair2(vv["pair"].(string)))
	}

	return pairlist, nil

}
