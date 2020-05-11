package aacoin

import (
	. "github.com/nntaoli-project/goex"
	"net/http"
	"net/url"

	"encoding/json"
	"errors"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	host = "https://api.aacoin.com/v1"
)

type Aacoin struct {
	accessKey,
	secretKey string
	httpClient *http.Client
}

func New(client *http.Client, api_key, secret_key string) *Aacoin {
	return &Aacoin{accessKey: api_key, secretKey: secret_key, httpClient: client}
}

func (aa *Aacoin) GetExchangeName() string {
	return "aacoin.com"
}

func (aa *Aacoin) buildSigned(params *url.Values) string {

	//log.Println("params", params.Encode())
	params.Set("accessKey", aa.accessKey)
	s, _ := GetParamHmacSHA256Sign(aa.secretKey, params.Encode())
	params.Set("sign", s)

	return s
}

func (aa *Aacoin) GetAccount() (*Account, error) {
	api_url := host + "/account/accounts"

	params := url.Values{}
	//params.Set("accessKey", aa.accessKey)
	//params.Set("sign", aa.buildSigned(&params))
	aa.buildSigned(&params)
	body, err := HttpPostForm(aa.httpClient, api_url, params)
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

	if bodyDataMap["status"].(string) != "1000" {
		return nil, errors.New(bodyDataMap["msg"].(string))
	}

	balances, isok := bodyDataMap["data"].([]interface{})
	if isok != true {
		return nil, errors.New("No account data!")
	}

	acc := Account{}
	acc.Exchange = aa.GetExchangeName()
	acc.SubAccounts = make(map[Currency]SubAccount)

	for _, v := range balances {
		vv := v.(map[string]interface{})
		currency := NewCurrency(vv["currencyCode"].(string), "")
		trade := 0.0
		frozen := 0.0
		vvv := vv["accounts"].([]interface{})
		for _, vvvv := range vvv {
			vvvvv := vvvv.(map[string]interface{})
			if vvvvv["type"].(string) == "trade" {
				trade = vvvvv["balance"].(float64)
			} else if vvvvv["type"].(string) == "frozen" {
				frozen = vvvvv["balance"].(float64)
			}
		}

		acc.SubAccounts[currency] = SubAccount{
			Currency:     currency,
			Amount:       trade + frozen,
			ForzenAmount: frozen,
		}
	}
	return &acc, nil
}

func (aa *Aacoin) placeOrder(amount, price string, pair CurrencyPair, orderType, orderSide string) (*Order, error) {
	path := host + "/order/place"
	params := url.Values{}
	params.Set("symbol", pair.String())
	params.Set("type", orderSide+"-"+orderType)
	params.Set("quantity", amount)
	params.Set("price", price)
	aa.buildSigned(&params)

	resp, err := HttpPostForm(aa.httpClient, path, params)
	log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return nil, err
	}

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(resp, &bodyDataMap)
	if err != nil {
		// log.Println("respData", string(resp))
		return nil, err
	}

	if bodyDataMap["status"].(string) != "1000" {
		return nil, errors.New(bodyDataMap["msg"].(string))
	}

	side := BUY
	if orderSide == "ASK" {
		side = SELL
	}

	return &Order{
		Currency:   pair,
		OrderID2:   strconv.FormatFloat(ToFloat64(bodyDataMap["data"]), 'f', 30, 64),
		Price:      ToFloat64(price),
		Amount:     ToFloat64(amount),
		DealAmount: 0,
		AvgPrice:   0,
		Side:       TradeSide(side),
		Status:     ORDER_UNFINISH,
		OrderTime:  int(time.Now().Unix())}, nil
}

func (aa *Aacoin) LimitBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return aa.placeOrder(amount, price, currencyPair, "limit", "buy")
}

func (aa *Aacoin) LimitSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return aa.placeOrder(amount, price, currencyPair, "limit", "sell")
}

func (aa *Aacoin) MarketBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return aa.placeOrder(amount, price, currencyPair, "market", "sell")
}

func (aa *Aacoin) MarketSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return aa.placeOrder(amount, price, currencyPair, "market", "sell")
}

func (aa *Aacoin) CancelOrder(orderId string, currencyPair CurrencyPair) (bool, error) {
	path := host + "/order/cancel"
	params := url.Values{}
	params.Set("orderId", orderId)
	aa.buildSigned(&params)

	resp, err := HttpPostForm(aa.httpClient, path, params)

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

func (aa *Aacoin) CancelOrders(orderId []string) (bool, error) {
	path := host + "/order/batchCancel"
	params := url.Values{}

	orders := strings.Join(orderId, ",")
	//log.Println("orders", orders)
	params.Set("orderIds", orders)
	aa.buildSigned(&params)

	resp, err := HttpPostForm(aa.httpClient, path, params)

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

func (aa *Aacoin) GetOneOrder(orderId string, currencyPair CurrencyPair) (*Order, error) {
	//path := API_BASE_URL + TRADE_URI + "/" + orderId
	//
	//respmap, err := HttpGet2(aa.httpClient, path, aa.privateHeader())
	////log.Println(respmap)
	//if err != nil {
	//	return nil, err
	//}
	//_order, isok := respmap["data"].(map[string]interface{})
	//if !isok {
	//	return nil, errors.New("no order data")
	//}
	//status, _ := _order["order_state"].(string)
	//side, _ := _order["order_side"].(string)
	//
	//ord := Order{}
	//ord.Currency = currencyPair
	//ord.OrderID2 = orderId
	//
	//if side == "ASK" {
	//	ord.Side = SELL
	//} else {
	//	ord.Side = BUY
	//}
	//
	//switch status {
	//case "open":
	//	ord.Status = ORDER_UNFINISH
	//case "filled":
	//	ord.Status = ORDER_FINISH
	//	//case "PARTIALLY_FILLED":
	//	//	ord.Status = ORDER_PART_FINISH
	//case "canceled":
	//	ord.Status = ORDER_CANCEL
	//	//case "PENDING_CANCEL":
	//	//	ord.Status = ORDER_CANCEL_ING
	//	//case "REJECTED":
	//	//	ord.Status = ORDER_REJECT
	//}
	//
	//ord.Amount = ToFloat64(_order["amount"].(string))
	//ord.Price = ToFloat64(_order["price"].(string))
	//ord.DealAmount = ToFloat64(_order["filled_amount"])
	//ord.AvgPrice = ord.Price // response no avg price ， fill price
	//
	//return &ord, nil
	return nil, nil
}
func (aa *Aacoin) GetUnfinishOrders(currencyPair CurrencyPair) ([]Order, error) {
	path := host + "/order/currentOrders"
	params := url.Values{}

	params.Set("symbol", currencyPair.String())
	aa.buildSigned(&params)

	resp, err := HttpPostForm(aa.httpClient, path, params)

	log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return nil, err
	}

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(resp, &bodyDataMap)
	if err != nil {
		// log.Println("respData", string(resp))
		return nil, err
	}

	if bodyDataMap["status"].(string) != "1000" {
		return nil, errors.New(bodyDataMap["msg"].(string))
	}
	data := bodyDataMap["data"].(map[string]interface{})
	list := data["list"].([]interface{})
	var orders []Order
	for _, v := range list {
		vv := v.(map[string]interface{})
		//yyyy-MM-dd HH:mm:ss
		ordertime, _ := time.Parse("2016-02-02 11:11:11", vv["orderTime"].(string))
		orders = append(orders, Order{
			Currency:   NewCurrencyPair2(vv["symbol"].(string)),
			OrderID2:   vv["orderId"].(string),
			Price:      ToFloat64(vv["price"]),
			Amount:     ToFloat64(vv["quantity"]),
			DealAmount: ToFloat64(vv["filledQuantity"]),
			Side:       orderTypeAdapter(vv["type"].(string)),
			Status:     orderStatusAdapter(vv["status"].(string)),
			OrderTime:  int(ordertime.Unix()),
		})

	}

	return nil, nil
}
func (aa *Aacoin) GetOrderHistorys(currencyPair CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	return nil, nil
}
func (aa *Aacoin) GetDepth(size int, currencyPair CurrencyPair) (*Depth, error) {
	path := host + "/market/depth"
	params := url.Values{}

	params.Set("symbol", currencyPair.String())
	//aa.buildSigned(&params)

	resp, err := HttpPostForm(aa.httpClient, path, params)

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

	if bodyDataMap["status"].(string) != "1000" {
		return nil, errors.New(bodyDataMap["msg"].(string))
	}

	dep := bodyDataMap["data"].(map[string]interface{})

	bids, _ := dep["bids"].([]interface{})
	asks, _ := dep["asks"].([]interface{})

	depth := new(Depth)
	//i := 0
	for _, bid := range bids {
		_bid := bid.([]interface{})
		amount := ToFloat64(_bid[1])
		price := ToFloat64(_bid[0])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.BidList = append(depth.BidList, dr)
		//if i < size {
		//	i++
		//} else {
		//	break
		//}
	}

	//i = 0
	for _, ask := range asks {
		_ask := ask.([]interface{})
		amount := ToFloat64(_ask[1])
		price := ToFloat64(_ask[0])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.AskList = append(depth.AskList, dr)
		//if i < size {
		//	i++
		//} else {
		//	break
		//}
	}
	sort.Sort(depth.AskList)
	sort.Sort(sort.Reverse(depth.BidList))
	depth.AskList = depth.AskList[0:size]
	depth.BidList = depth.BidList[0:size]
	//log.Println(depth)
	return depth, nil
}

func (aa *Aacoin) GetTicker(currencyPair CurrencyPair) (*Ticker, error) {
	path := host + "/market/detail"
	params := url.Values{}

	params.Set("symbol", currencyPair.String())

	resp, err := HttpPostForm(aa.httpClient, path, params)

	// log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return nil, err
	}

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(resp, &bodyDataMap)
	if err != nil {
		// log.Println("respData", string(resp))
		return nil, err
	}

	if bodyDataMap["status"].(string) != "1000" {
		return nil, errors.New(bodyDataMap["msg"].(string))
	}
	tickerMap := bodyDataMap["data"].(map[string]interface{})
	var ticker Ticker

	ticker.Date = uint64(time.Now().Unix())
	ticker.Last = ToFloat64(tickerMap["current"])
	ticker.Buy = ToFloat64(tickerMap["buy"])
	ticker.Sell = ToFloat64(tickerMap["sell"])
	ticker.Low = ToFloat64(tickerMap["lowest"])
	ticker.High = ToFloat64(tickerMap["highest"])
	ticker.Vol = ToFloat64(tickerMap["totalTradeAmount"])
	return &ticker, nil

}

func (aa *Aacoin) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录
func (aa *Aacoin) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}

func orderTypeAdapter(t string) TradeSide {
	switch t {
	case "sell":
		return SELL
	case "buy":
		return BUY
	}
	return BUY
}

func orderStatusAdapter(s string) TradeStatus {
	switch s {
	case "open":
		return ORDER_UNFINISH
	case "partial_filled":
		return ORDER_PART_FINISH
	}
	return ORDER_FINISH
}
