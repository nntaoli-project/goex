package coin58

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Coin58 struct {
	client       *http.Client
	apikey       string
	apisecretkey string
	apiurl       string
}

//58coin.com  closed the trade api
func New58Coin(client *http.Client, apikey string, apisecretkey string) *Coin58 {
	return &Coin58{client: client, apikey: apikey, apisecretkey: apisecretkey, apiurl: "https://api.58coin.com/v1/"}
}

func (coin58 *Coin58) placeOrder(t, side, amount, price string, currency CurrencyPair) (*Order, error) {
	var params = url.Values{}
	params.Set("symbol", currency.AdaptUsdToUsdt().ToSymbol("_"))
	params.Set("type", t)
	params.Set("side", side)
	params.Set("amount", amount)
	//params.Set("client_oid" , "1")
	if t == "limit" {
		params.Set("price", price)
	}

	r, err := coin58.doAuthenticatedRequest("spot/my/order/place", params)
	if err != nil {
		return nil, err
	}

	o := r.(map[string]interface{})

	//log.Println(r)

	var tradeSide TradeSide = SELL

	switch side {
	case "sell":
		tradeSide = SELL
	case "buy":
		tradeSide = BUY
	}

	return &Order{
		OrderID2:  o["order_id"].(string),
		Currency:  currency,
		Amount:    ToFloat64(amount),
		Price:     ToFloat64(price),
		Status:    ORDER_UNFINISH,
		Side:      tradeSide,
		OrderTime: ToInt(o["created_time"])}, nil
}

func (coin58 *Coin58) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return coin58.placeOrder("limit", "buy", amount, price, currency)
}

func (coin58 *Coin58) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return coin58.placeOrder("limit", "sell", amount, price, currency)
}

func (coin58 *Coin58) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return coin58.placeOrder("market", "buy", amount, price, currency)
}

func (coin58 *Coin58) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return coin58.placeOrder("market", "sell", amount, price, currency)
}

func (coin58 *Coin58) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	params := url.Values{}
	params.Set("symbol", currency.AdaptUsdToUsdt().ToSymbol("_"))
	params.Set("order_id", orderId)

	_, err := coin58.doAuthenticatedRequest("spot/my/order/cancel", params)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (coin58 *Coin58) toOrder(o map[string]interface{}, pair CurrencyPair) *Order {
	var side TradeSide = SELL
	if o["side"].(string) == "buy" {
		side = BUY
	}

	var oStatus TradeStatus = ORDER_UNFINISH

	switch o["status"].(string) {
	case "Finished", "finished":
		oStatus = ORDER_FINISH
	case "Cancelled", "cancelled":
		oStatus = ORDER_CANCEL
	case "Cancelling", "cancelling":
		oStatus = ORDER_CANCEL_ING
		//case "Active", "active":
		//	oStatus = ORDER_PART_FINISH
	}

	avg := 0.0
	if ToFloat64(o["quote_filled"]) > 0 && ToFloat64(o["base_filled"]) > 0 {
		avg = ToFloat64(o["quote_filled"]) / ToFloat64(o["base_filled"])
	}

	return &Order{
		Currency:   pair,
		OrderID2:   fmt.Sprint(o["order_id"]),
		Price:      ToFloat64(o["price"]),
		Amount:     ToFloat64(o["amount"]) / ToFloat64(o["price"]),
		AvgPrice:   avg,
		DealAmount: ToFloat64(o["base_filled"]),
		Side:       side,
		Status:     oStatus,
		OrderTime:  ToInt(o["created_time"])}
}

func (coin58 *Coin58) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	params := url.Values{}
	params.Set("symbol", currency.AdaptUsdToUsdt().ToSymbol("_"))
	params.Set("order_id", orderId)

	r, err := coin58.doAuthenticatedRequest("spot/my/order", params)
	if err != nil {
		return nil, err
	}

	o := r.(map[string]interface{})
	//	log.Println(o)
	return coin58.toOrder(o, currency), nil
}

func (coin58 *Coin58) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	params := url.Values{}
	params.Set("symbol", currency.AdaptUsdToUsdt().ToSymbol("_"))

	r, err := coin58.doAuthenticatedRequest("spot/my/orders", params)
	if err != nil {
		return nil, err
	}

	omaps, isok := r.([]interface{})
	if !isok {
		return []Order{}, nil
	}
	//log.Println(omaps)
	var orders []Order
	for _, o := range omaps {
		orders = append(orders, *coin58.toOrder(o.(map[string]interface{}), currency))
	}

	return orders, nil
}

func (coin58 *Coin58) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implement")
}

func (coin58 *Coin58) GetAccount() (*Account, error) {
	r, err := coin58.doAuthenticatedRequest("spot/my/accounts", url.Values{})
	if err != nil {
		return nil, err
	}

	acc := new(Account)
	acc.Exchange = COIN58
	acc.SubAccounts = make(map[Currency]SubAccount)

	for _, c := range r.([]interface{}) {
		c2 := c.(map[string]interface{})
		balance := ToFloat64(c2["balance"])
		available := ToFloat64(c2["available"])
		currency := NewCurrency(c2["currency"].(string), "")
		acc.SubAccounts[currency] = SubAccount{
			Currency:     currency,
			Amount:       available,
			ForzenAmount: balance - available,
			LoanAmount:   0}
	}

	return acc, err
}

func (coin58 *Coin58) GetTicker(currency CurrencyPair) (*Ticker, error) {
	tickerUrl := coin58.apiurl + "spot/ticker?symbol=" + currency.AdaptUsdToUsdt().ToSymbol("_")
	m, err := HttpGet(coin58.client, tickerUrl)
	if err != nil {
		return nil, err
	}

	error := m["error"]
	if error != nil {
		return nil, coin58.adaptError(error.(map[string]interface{}))
	}

	r, isok := m["result"].([]interface{})
	if !isok || len(r) == 0 {
		return nil, API_ERR
	}

	t := r[0].(map[string]interface{})

	return &Ticker{Pair: currency,
		Last: ToFloat64(t["last"]),
		Low:  ToFloat64(t["low"]),
		High: ToFloat64(t["high"]),
		Sell: ToFloat64(t["ask"]),
		Buy:  ToFloat64(t["bid"]),
		Vol:  ToFloat64(t["volume"]),
		Date: ToUint64(t["time"]) / 1000}, nil
}

func (coin58 *Coin58) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	depurl := coin58.apiurl + "spot/order_book?symbol=" + currency.AdaptUsdToUsdt().ToSymbol("_") + "&limit=" + fmt.Sprint(size)
	m, err := HttpGet(coin58.client, depurl)
	if err != nil {
		return nil, HTTP_ERR_CODE.OriginErr(err.Error())
	}

	error := m["error"]
	if error != nil {
		return nil, coin58.adaptError(error.(map[string]interface{}))
	}
	//	log.Println(m)
	r := m["result"].(map[string]interface{})
	asks := r["asks"].([]interface{})
	bids := r["bids"].([]interface{})

	dep := new(Depth)
	dep.Pair = currency

	for _, ask := range asks {
		ask2 := ask.([]interface{})
		dep.AskList = append(dep.AskList, DepthRecord{ToFloat64(ask2[0]), ToFloat64(ask2[1])})
	}

	for _, bid := range bids {
		bid2 := bid.([]interface{})
		dep.BidList = append(dep.BidList, DepthRecord{ToFloat64(bid2[0]), ToFloat64(bid2[1])})
	}

	sort.Sort(sort.Reverse(dep.AskList))
	sort.Sort(sort.Reverse(dep.BidList))

	return dep, nil
}

func (coin58 *Coin58) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录
func (coin58 *Coin58) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}

func (coin58 *Coin58) GetExchangeName() string {
	return COIN58
}

func (coin58 *Coin58) doAuthenticatedRequest(api string, params url.Values) (interface{}, error) {
	url := coin58.apiurl + api

	var paramsmap = map[string]string{}
	for k, v := range params {
		paramsmap[k] = v[0]
	}

	timestamp := fmt.Sprint(time.Now().Unix() * 1000)

	params.Set("api_key", coin58.apikey)

	signParamString := params.Encode() + "&api_secret=" + coin58.apisecretkey + "&timestamp=" + timestamp

	sign, _ := GetParamHmacSHA256Sign(coin58.apisecretkey, signParamString)

	header := map[string]string{
		"X-58COIN-APIKEY":    coin58.apikey,
		"X-58COIN-SIGNATURE": base64.URLEncoding.EncodeToString([]byte(strings.ToUpper(sign))),
		"X-58COIN-TIMESTAMP": timestamp}

	params.Del("api_key")

	body, err := HttpPostForm2(coin58.client, url, params, header)

	if err != nil {
		return nil, HTTP_ERR_CODE
	}

	var m map[string]interface{}
	json.Unmarshal(body, &m)

	errm := m["error"]
	if errm != nil {
		return nil, coin58.adaptError(errm.(map[string]interface{}))
	}

	return m["result"], nil
}

func (coin58 *Coin58) adaptError(errmap map[string]interface{}) ApiError {
	code := ToInt(errmap["code"])
	switch code {
	case 10006:
		return EX_ERR_NOT_FIND_APIKEY
	case 10007:
		return EX_ERR_SIGN
	case 20000:
		return EX_ERR_SYMBOL_ERR
	case 20015:
		return EX_ERR_NOT_FIND_ORDER
	case 20007:
		return EX_ERR_INSUFFICIENT_BALANCE
	}
	//log.Println(errmap)
	return API_ERR
}
