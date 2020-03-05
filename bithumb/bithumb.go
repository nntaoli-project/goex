package bithumb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

type Bithumb struct {
	client *http.Client
	accesskey,
	secretkey string
}

var (
	baseUrl = "https://api.bithumb.com"
)

func New(client *http.Client, accesskey, secretkey string) *Bithumb {
	return &Bithumb{client: client, accesskey: accesskey, secretkey: secretkey}
}

func (bit *Bithumb) placeOrder(side, amount, price string, pair CurrencyPair) (*Order, error) {
	var retmap map[string]interface{}
	params := fmt.Sprintf("order_currency=%s&units=%s&price=%s&type=%s", pair.CurrencyA.Symbol, amount, price, side)
	log.Println(params)
	err := bit.doAuthenticatedRequest("/trade/place", params, &retmap)
	if err != nil {
		return nil, err
	}
	if retmap["status"].(string) != "0000" {
		log.Println(retmap)
		return nil, errors.New(retmap["status"].(string))
	}

	var tradeSide TradeSide
	switch side {
	case "ask":
		tradeSide = SELL
	case "bid":
		tradeSide = BUY
	}

	log.Println(retmap)
	return &Order{
		OrderID:  ToInt(retmap["order_id"]),
		Amount:   ToFloat64(amount),
		Price:    ToFloat64(price),
		Currency: pair,
		Side:     tradeSide,
		Status:   ORDER_UNFINISH}, nil
}

func (bit *Bithumb) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return bit.placeOrder("bid", amount, price, currency)
}

func (bit *Bithumb) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return bit.placeOrder("ask", amount, price, currency)
}

func (bit *Bithumb) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}

func (bit *Bithumb) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}

func (bit *Bithumb) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("please invoke the CancelOrder2 method.")
}

/*补丁*/
func (bit *Bithumb) CancelOrder2(side, orderId string, currency CurrencyPair) (bool, error) {
	var retmap map[string]interface{}
	params := fmt.Sprintf("type=%s&order_id=%s&currency=%s", side, orderId, currency.CurrencyA.Symbol)
	err := bit.doAuthenticatedRequest("/trade/cancel", params, &retmap)
	if err != nil {
		return false, err
	}
	if retmap["status"].(string) == "0000" {
		return true, nil
	}
	return false, errors.New(retmap["status"].(string))
}

func (bit *Bithumb) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("please invoke the GetOneOrder2 method.")
}

/*补丁*/
func (bit *Bithumb) GetOneOrder2(side, orderId string, currency CurrencyPair) (*Order, error) {
	var retmap map[string]interface{}
	params := fmt.Sprintf("type=%s&order_id=%s&currency=%s", side, orderId, currency.CurrencyA.Symbol)
	err := bit.doAuthenticatedRequest("/info/order_detail", params, &retmap)
	if err != nil {
		return nil, err
	}

	if retmap["status"].(string) != "0000" {
		message := retmap["message"].(string)
		if "거래 체결내역이 존재하지 않습니다." == message {
			return nil, EX_ERR_NOT_FIND_ORDER
		}
		log.Println(retmap)
		return nil, errors.New(retmap["status"].(string))
	}

	order := new(Order)
	total := 0.0
	data := retmap["data"].([]interface{})
	for _, v := range data {
		ord := v.(map[string]interface{})
		switch ord["type"].(string) {
		case "ask":
			order.Side = SELL
		case "bid":
			order.Side = BUY
		}
		order.Amount += ToFloat64(ord["units_traded"])
		order.Fee += ToFloat64(ord["fee"])
		total += ToFloat64(ord["total"])
	}

	order.DealAmount = order.Amount
	avg := total / order.DealAmount
	order.AvgPrice = ToFloat64(fmt.Sprintf("%.2f", avg))
	order.Price = order.AvgPrice
	order.Currency = currency
	order.OrderID = ToInt(orderId)
	order.Status = ORDER_FINISH

	log.Println(retmap)
	return order, nil
}

func (bit *Bithumb) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	var retmap map[string]interface{}
	params := fmt.Sprintf("currency=%s", currency.CurrencyA.Symbol)
	err := bit.doAuthenticatedRequest("/info/orders", params, &retmap)
	if err != nil {
		return nil, err
	}

	if retmap["status"].(string) != "0000" {
		message := retmap["message"].(string)
		if "거래 진행중인 내역이 존재하지 않습니다." == message {
			return []Order{}, nil
		}
		return nil, errors.New(fmt.Sprintf("[%s]%s", retmap["status"].(string), message))
	}

	var orders []Order
	datas := retmap["data"].([]interface{})
	for _, v := range datas {
		orderinfo := v.(map[string]interface{})
		ord := Order{
			OrderID:  ToInt(orderinfo["order_id"]),
			Amount:   ToFloat64(orderinfo["units"]),
			Price:    ToFloat64(orderinfo["price"]),
			Currency: currency,
			Fee:      ToFloat64(orderinfo["fee"])}

		remaining := ToFloat64(orderinfo["units_remaining"])
		total := ToFloat64(orderinfo["total"])
		dealamount := ord.Amount - remaining
		ord.DealAmount = dealamount

		if dealamount > 0 {
			avg := fmt.Sprintf("%.4f", total/dealamount)
			ord.AvgPrice = ToFloat64(avg)
		}

		switch orderinfo["type"].(string) {
		case "ask":
			ord.Side = SELL
		case "bid":
			ord.Side = BUY
		}

		switch orderinfo["status"].(string) {
		case "placed":
			ord.Status = ORDER_UNFINISH
		}

		orders = append(orders, ord)
	}

	log.Println(retmap)
	return orders, nil
}

func (bit *Bithumb) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implement")
}

func (bit *Bithumb) GetAccount() (*Account, error) {
	var retmap map[string]interface{}
	err := bit.doAuthenticatedRequest("/info/balance", "currency=ALL", &retmap)
	if err != nil {
		return nil, err
	}
	datamap := retmap["data"].(map[string]interface{})
	acc := new(Account)
	acc.SubAccounts = make(map[Currency]SubAccount)
	acc.SubAccounts[LTC] = SubAccount{
		Currency:     LTC,
		Amount:       ToFloat64(datamap["available_ltc"]),
		ForzenAmount: ToFloat64(datamap["in_use_ltc"]),
		LoanAmount:   0}
	acc.SubAccounts[BTC] = SubAccount{
		Currency:     BTC,
		Amount:       ToFloat64(datamap["available_btc"]),
		ForzenAmount: ToFloat64(datamap["in_use_etc"]),
		LoanAmount:   0}
	acc.SubAccounts[ETH] = SubAccount{
		Currency:     ETH,
		Amount:       ToFloat64(datamap["available_eth"]),
		ForzenAmount: ToFloat64(datamap["in_use_eth"]),
		LoanAmount:   0}
	acc.SubAccounts[ETC] = SubAccount{
		Currency:     ETC,
		Amount:       ToFloat64(datamap["available_etc"]),
		ForzenAmount: ToFloat64(datamap["in_use_etc"]),
		LoanAmount:   0}
	acc.SubAccounts[BCH] = SubAccount{
		Currency:     BCH,
		Amount:       ToFloat64(datamap["available_bch"]),
		ForzenAmount: ToFloat64(datamap["in_use_bch"]),
		LoanAmount:   0}
	acc.SubAccounts[KRW] = SubAccount{
		Currency:     KRW,
		Amount:       ToFloat64(datamap["available_krw"]),
		ForzenAmount: ToFloat64(datamap["in_use_krw"]),
		LoanAmount:   0}
	//log.Println(datamap)
	acc.Exchange = bit.GetExchangeName()
	return acc, nil
}

func (bit *Bithumb) doAuthenticatedRequest(uri, params string, ret interface{}) error {
	nonce := time.Now().UnixNano() / int64(time.Millisecond)
	api_nonce := fmt.Sprint(nonce)
	e_endpoint := url.QueryEscape(uri)
	params += "&endpoint=" + e_endpoint

	// Api-Sign information generation.
	hmac_data := uri + string(0) + params + string(0) + api_nonce
	hash_hmac_str := GetParamHmacSHA512Base64Sign(bit.secretkey, hmac_data)
	api_sign := hash_hmac_str
	content_length_str := strconv.Itoa(len(params))

	// Connects to Bithumb API server and returns JSON result value.
	resp, err := NewHttpRequest(bit.client, "POST", baseUrl+uri,
		bytes.NewBufferString(params).String(), map[string]string{
			"Api-Key":        bit.accesskey,
			"Api-Sign":       api_sign,
			"Api-Nonce":      api_nonce,
			"Content-Type":   "application/x-www-form-urlencoded",
			"Content-Length": content_length_str,
		}) // URL-encoded payload

	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, ret)

	return err
}

func (bit *Bithumb) GetTicker(currency CurrencyPair) (*Ticker, error) {
	respmap, err := HttpGet(bit.client, fmt.Sprintf("%s/public/ticker/%s", baseUrl, currency.CurrencyA))
	if err != nil {
		return nil, err
	}
	s, isok := respmap["status"].(string)
	if s != "0000" || isok != true {
		msg := "ticker error"
		if isok {
			msg = s
		}
		return nil, errors.New(msg)
	}

	datamap := respmap["data"].(map[string]interface{})

	return &Ticker{
		Low:  ToFloat64(datamap["min_price"]),
		High: ToFloat64(datamap["max_price"]),
		Last: ToFloat64(datamap["closing_price"]),
		Vol:  ToFloat64(datamap["units_traded"]),
		Buy:  ToFloat64(datamap["buy_price"]),
		Sell: ToFloat64(datamap["sell_price"]),
	}, nil
}

func (bit *Bithumb) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	resp, err := HttpGet(bit.client, fmt.Sprintf("%s/public/orderbook/%s", baseUrl, currency.CurrencyA))
	if err != nil {
		return nil, err
	}

	if resp["status"].(string) != "0000" {
		return nil, errors.New(resp["status"].(string))
	}

	datamap := resp["data"].(map[string]interface{})
	bids := datamap["bids"].([]interface{})
	asks := datamap["asks"].([]interface{})

	dep := new(Depth)

	for _, v := range bids {
		bid := v.(map[string]interface{})
		dep.BidList = append(dep.BidList, DepthRecord{ToFloat64(bid["price"]), ToFloat64(bid["quantity"])})
	}

	for _, v := range asks {
		ask := v.(map[string]interface{})
		dep.AskList = append(dep.AskList, DepthRecord{ToFloat64(ask["price"]), ToFloat64(ask["quantity"])})
	}

	sort.Sort(sort.Reverse(dep.AskList))

	return dep, nil
}

func (bit *Bithumb) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录
func (bit *Bithumb) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}

func (bit *Bithumb) GetExchangeName() string {
	return BITHUMB
}
