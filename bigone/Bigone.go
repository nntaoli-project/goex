package bigone

import (
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"net/http"

	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"strconv"
	"time"
)

const (
	API_BASE_URL = "https://api.big.one/"
	TICKER_URI   = "markets/%s"
	DEPTH_URI    = "markets/%s/book"
	ACCOUNT_URI  = "accounts"
	ORDERS_URI   = "orders?market=%s"
	TRADE_URI    = "orders"
)

type Bigone struct {
	accessKey,
	secretKey string
	httpClient *http.Client
	uid        string
}

func New(client *http.Client, api_key, secret_key string) *Bigone {
	return &Bigone{api_key, secret_key, client, uuid.New().String()}
}

func (bo *Bigone) GetExchangeName() string {
	return BIGONE
}

func (bo *Bigone) GetTicker(currency CurrencyPair) (*Ticker, error) {
	tickerUri := API_BASE_URL + fmt.Sprintf(TICKER_URI, currency.ToSymbol("-"))
	bodyDataMap, err := HttpGet(bo.httpClient, tickerUri)

	if err != nil {
		return nil, err
	}

	log.Println("bo uri:", tickerUri)
	log.Println("bo bodyDataMap:", currency, bodyDataMap)

	dataMap := bodyDataMap["data"].(map[string]interface{})
	tickerMap := dataMap["ticker"].(map[string]interface{})
	asksMap := dataMap["asks"].([]interface{})
	bidsMap := dataMap["bids"].([]interface{})
	ask := asksMap[0].(map[string]interface{})
	bid := bidsMap[0].(map[string]interface{})

	log.Println(tickerMap)
	log.Println(asksMap)
	log.Println(bidsMap)

	var ticker Ticker

	ticker.Date = uint64(time.Now().Unix())
	ticker.Last = ToFloat64(tickerMap["price"])
	ticker.Buy = ToFloat64(bid["price"])
	ticker.Sell = ToFloat64(ask["price"])
	ticker.Low = ToFloat64(tickerMap["low"])
	ticker.High = ToFloat64(tickerMap["high"])
	ticker.Vol = ToFloat64(tickerMap["volume"])
	return &ticker, nil
}

func (bo *Bigone) placeOrder(amount, price string, pair CurrencyPair, orderType, orderSide string) (*Order, error) {
	path := API_BASE_URL + TRADE_URI
	params := make(map[string]string)
	params["order_market"] = pair.ToSymbol("-")
	params["order_side"] = orderSide
	params["amount"] = amount
	params["price"] = price

	resp, err := HttpPostForm4(bo.httpClient, path, params, bo.privateHeader())
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
	ordinfo, isok := respmap["data"].(map[string]interface{})
	if !isok {
		return nil, errors.New("no order data")
	}

	orderId2, ok := ordinfo["order_id"].(string)
	if ok != true {
		return nil, errors.New(string(resp))
	}

	side := BUY
	if orderSide == "ASK" {
		side = SELL
	}

	return &Order{
		Currency:   pair,
		OrderID2:   orderId2,
		Price:      ToFloat64(price),
		Amount:     ToFloat64(amount),
		DealAmount: 0,
		AvgPrice:   0,
		Side:       TradeSide(side),
		Status:     ORDER_UNFINISH,
		OrderTime:  int(time.Now().Unix())}, nil
}

func (bo *Bigone) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return bo.placeOrder(amount, price, currency, "LIMIT", "BID")
}

func (bo *Bigone) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return bo.placeOrder(amount, price, currency, "LIMIT", "ASK")
}

func (bo *Bigone) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (bo *Bigone) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (bo *Bigone) privateHeader() map[string]string {
	return map[string]string{"User-Agent": "standard browser user agent format",
		"Authorization": "Bearer " + bo.secretKey,
		"Big-Device-Id": bo.uid}
}
func (bo *Bigone) getOrdersList(currencyPair CurrencyPair, size int, sts TradeStatus) ([]Order, error) {
	apiURL := ""
	if size <= 0 { //unlimited
		apiURL = API_BASE_URL + fmt.Sprintf(ORDERS_URI, currencyPair.ToSymbol("-"))
	} else {
		apiURL = API_BASE_URL + fmt.Sprintf(ORDERS_URI, currencyPair.ToSymbol("-")) + "&limit=" + strconv.Itoa(size)
	}
	if sts == ORDER_FINISH {
		apiURL += "&state=done"
	} else {
		apiURL += "&state=open"
	}
	respmap, err := HttpGet2(bo.httpClient, apiURL, bo.privateHeader())
	log.Println(respmap)
	if err != nil {
		return nil, err
	}
	lists, isok := respmap["data"].([]interface{})
	if !isok {
		return nil, errors.New("no order data")
	}
	orders := make([]Order, 0)
	for _, _order := range lists {
		order, _ := _order.(map[string]interface{})
		status, _ := order["order_state"].(string)
		side, _ := order["order_side"].(string)

		ord := Order{}

		switch status {
		case "open":
			ord.Status = ORDER_UNFINISH
		case "filled":
			ord.Status = ORDER_FINISH
		//case "PARTIALLY_FILLED":
		//	ord.Status = ORDER_PART_FINISH
		case "canceled":
			ord.Status = ORDER_CANCEL
			//case "PENDING_CANCEL":
			//	ord.Status = ORDER_CANCEL_ING
			//case "REJECTED":
			//	ord.Status = ORDER_REJECT
		}
		if ord.Status != sts {
			continue // discard
		}

		ord.Currency = currencyPair
		ord.OrderID2 = order["order_id"].(string)

		if side == "ASK" {
			ord.Side = SELL
		} else {
			ord.Side = BUY
		}

		ord.Amount = ToFloat64(order["amount"].(string))
		ord.Price = ToFloat64(order["price"].(string))
		ord.DealAmount = ToFloat64(order["filled_amount"])
		ord.AvgPrice = ord.Price // response no avg price ， fill price
		orders = append(orders, ord)
	}

	return orders, nil
}

func (bo *Bigone) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	path := API_BASE_URL + TRADE_URI + "/" + orderId

	resp, err := HttpDeleteForm(bo.httpClient, path, nil, bo.privateHeader())

	log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (bo *Bigone) GetOneOrder(orderId string, currencyPair CurrencyPair) (*Order, error) {
	path := API_BASE_URL + TRADE_URI + "/" + orderId

	respmap, err := HttpGet2(bo.httpClient, path, bo.privateHeader())
	log.Println(respmap)
	if err != nil {
		return nil, err
	}
	_order, isok := respmap["data"].(map[string]interface{})
	if !isok {
		return nil, errors.New("no order data")
	}
	status, _ := _order["order_state"].(string)
	side, _ := _order["order_side"].(string)

	ord := Order{}
	ord.Currency = currencyPair
	ord.OrderID2 = orderId

	if side == "ASK" {
		ord.Side = SELL
	} else {
		ord.Side = BUY
	}

	switch status {
	case "open":
		ord.Status = ORDER_UNFINISH
	case "filled":
		ord.Status = ORDER_FINISH
		//case "PARTIALLY_FILLED":
		//	ord.Status = ORDER_PART_FINISH
	case "canceled":
		ord.Status = ORDER_CANCEL
		//case "PENDING_CANCEL":
		//	ord.Status = ORDER_CANCEL_ING
		//case "REJECTED":
		//	ord.Status = ORDER_REJECT
	}

	ord.Amount = ToFloat64(_order["amount"].(string))
	ord.Price = ToFloat64(_order["price"].(string))
	ord.DealAmount = ToFloat64(_order["filled_amount"])
	ord.AvgPrice = ord.Price // response no avg price ， fill price

	return &ord, nil

}
func (bo *Bigone) GetUnfinishOrders(currencyPair CurrencyPair) ([]Order, error) {
	return bo.getOrdersList(currencyPair, -1, ORDER_UNFINISH)
}
func (bo *Bigone) GetOrderHistorys(currencyPair CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	return bo.getOrdersList(currencyPair, -1, ORDER_FINISH)
}

func (bo *Bigone) GetAccount() (*Account, error) {
	apiUrl := API_BASE_URL + ACCOUNT_URI

	resp, err := HttpGet2(bo.httpClient, apiUrl, bo.privateHeader())
	if err != nil {
		log.Println("GetAccount error:", err)
		return nil, err
	}
	balances, isok := resp["data"].([]interface{})
	if isok != true {
		return nil, errors.New("No account data!")
	}
	acc := Account{}
	acc.Exchange = bo.GetExchangeName()
	acc.SubAccounts = make(map[Currency]SubAccount)

	for _, v := range balances {
		//log.Println(v)
		vv := v.(map[string]interface{})
		currency := NewCurrency(vv["account_type"].(string), "")
		acc.SubAccounts[currency] = SubAccount{
			Currency:     currency,
			Amount:       ToFloat64(vv["active_balance"]),
			ForzenAmount: ToFloat64(vv["frozen_balance"]),
		}
	}

	return &acc, nil
}

func (bo *Bigone) GetDepth(size int, currencyPair CurrencyPair) (*Depth, error) {

	apiUrl := API_BASE_URL + fmt.Sprintf(DEPTH_URI, currencyPair.ToSymbol("-"))
	resp, err := HttpGet(bo.httpClient, apiUrl)
	if err != nil {
		log.Println("GetDepth error:", err)
		return nil, err
	}
	//log.Println(resp)

	dep, isok := resp["data"].(map[string]interface{})
	if !isok {
		return nil, errors.New("no depth data")
	}

	bids, _ := dep["bids"].([]interface{})
	asks, _ := dep["asks"].([]interface{})

	depth := new(Depth)

	for _, bid := range bids {
		_bid := bid.(map[string]interface{})
		amount := ToFloat64(_bid["amount"])
		price := ToFloat64(_bid["price"])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.BidList = append(depth.BidList, dr)
	}

	for _, ask := range asks {
		_ask := ask.(map[string]interface{})
		amount := ToFloat64(_ask["amount"])
		price := ToFloat64(_ask["price"])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.AskList = append(depth.AskList, dr)
	}

	return depth, nil
}

func (bo *Bigone) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implements")
}

//非个人，整个交易所的交易记录
func (bo *Bigone) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implements")
}
