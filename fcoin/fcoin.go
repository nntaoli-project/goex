package fcoin

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DEPTH_API       = "market/depth/%s/%s"
	TRADE_URL       = "orders"
	GET_ACCOUNT_API = "accounts/balance"
	GET_ORDER_API   = "orders/%s"
	//GET_ORDERS_LIST_API             = ""
	GET_UNFINISHED_ORDERS_API = "getUnfinishedOrdersIgnoreTradeType"
	PLACE_ORDER_API           = "order"
	WITHDRAW_API              = "withdraw"
	CANCELWITHDRAW_API        = "cancelWithdraw"
	SERVER_TIME               = "public/server-time"
)

type FCoinTicker struct {
	Ticker
	SellAmount,
	BuyAmount float64
}

type FCoin struct {
	httpClient *http.Client
	baseUrl,
	accessKey,
	secretKey string
	timeoffset int64
}

func NewFCoin(client *http.Client, apikey, secretkey string) *FCoin {
	fc := &FCoin{baseUrl: "https://api.fcoin.com/v2/", accessKey: apikey, secretKey: secretkey, httpClient: client}
	fc.setTimeOffset()
	return fc
}

func (ft *FCoin) GetExchangeName() string {
	return FCOIN
}

func (ft *FCoin) setTimeOffset() error {
	respmap, err := HttpGet(ft.httpClient, ft.baseUrl+"public/server-time")
	if err != nil {
		return err
	}
	stime := int64(ToInt(respmap["data"]))
	st := time.Unix(stime/1000, 0)
	lt := time.Now()
	offset := st.Sub(lt).Seconds()
	ft.timeoffset = int64(offset)
	return nil
}

func (ft *FCoin) GetTicker(currencyPair CurrencyPair) (*Ticker, error) {
	respmap, err := HttpGet(ft.httpClient, ft.baseUrl+fmt.Sprintf("market/ticker/%s",
		strings.ToLower(currencyPair.ToSymbol(""))))

	if err != nil {
		return nil, err
	}

	////log.Println("ticker respmap:", respmap)
	if respmap["status"].(float64) != 0 {
		return nil, errors.New(respmap["err-msg"].(string))
	}

	//
	tick, ok := respmap["data"].(map[string]interface{})
	if !ok {
		return nil, API_ERR
	}

	tickmap, ok := tick["ticker"].([]interface{})
	if !ok {
		return nil, API_ERR
	}

	ticker := new(Ticker)
	ticker.Pair = currencyPair
	ticker.Date = uint64(time.Now().Nanosecond() / 1000)
	ticker.Last = ToFloat64(tickmap[0])
	ticker.Vol = ToFloat64(tickmap[9])
	ticker.Low = ToFloat64(tickmap[8])
	ticker.High = ToFloat64(tickmap[7])
	ticker.Buy = ToFloat64(tickmap[2])
	ticker.Sell = ToFloat64(tickmap[4])
	//ticker.SellAmount = ToFloat64(tickmap[5])
	//ticker.BuyAmount = ToFloat64(tickmap[3])

	return ticker, nil

}

func (ft *FCoin) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	respmap, err := HttpGet(ft.httpClient, ft.baseUrl+fmt.Sprintf("market/depth/L20/%s", strings.ToLower(currency.ToSymbol(""))))
	if err != nil {
		return nil, err
	}

	if respmap["status"].(float64) != 0 {
		return nil, errors.New(respmap["err-msg"].(string))
	}

	datamap := respmap["data"].(map[string]interface{})

	bids, ok1 := datamap["bids"].([]interface{})
	asks, ok2 := datamap["asks"].([]interface{})

	if !ok1 || !ok2 {
		return nil, errors.New("depth error")
	}

	depth := new(Depth)
	depth.Pair = currency

	n := 0
	for i := 0; i < len(bids); {
		depth.BidList = append(depth.BidList, DepthRecord{ToFloat64(bids[i]), ToFloat64(bids[i+1])})
		i += 2
		n++
		if n == size {
			break
		}
	}

	n = 0
	for i := 0; i < len(asks); {
		depth.AskList = append(depth.AskList, DepthRecord{ToFloat64(asks[i]), ToFloat64(asks[i+1])})
		i += 2
		n++
		if n == size {
			break
		}
	}

	//sort.Sort(sort.Reverse(depth.AskList))

	return depth, nil
}
func (ft *FCoin) doAuthenticatedRequest(method, uri string, params url.Values) (interface{}, error) {

	timestamp := time.Now().Unix()*1000 + ft.timeoffset*1000
	sign := ft.buildSigned(method, ft.baseUrl+uri, timestamp, params)

	header := map[string]string{
		"FC-ACCESS-KEY":       ft.accessKey,
		"FC-ACCESS-SIGNATURE": sign,
		"FC-ACCESS-TIMESTAMP": fmt.Sprint(timestamp)}

	var (
		respmap map[string]interface{}
		err     error
	)

	switch method {
	case "GET":
		respmap, err = HttpGet2(ft.httpClient, ft.baseUrl+uri+"?"+params.Encode(), header)
		if err != nil {
			return nil, err
		}

	case "POST":
		var parammap = make(map[string]string, 1)
		for k, v := range params {
			parammap[k] = v[0]
		}

		respbody, err := HttpPostForm4(ft.httpClient, ft.baseUrl+uri, parammap, header)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(respbody, &respmap)
	}
	log.Println(respmap)
	if ToInt(respmap["status"]) != 0 {
		return nil, errors.New(respmap["msg"].(string))
	}

	return respmap["data"], err
}

func (ft *FCoin) buildSigned(httpmethod string, apiurl string, timestamp int64, para url.Values) string {
	param := ""
	if para != nil {
		param = para.Encode()
	}

	if "GET" == httpmethod && param != "" {
		apiurl += "?" + param
	}

	signStr := httpmethod + apiurl + fmt.Sprint(timestamp)
	if "POST" == httpmethod && param != "" {
		signStr += param
	}

	log.Println(signStr)

	sign := base64.StdEncoding.EncodeToString([]byte(signStr))

	mac := hmac.New(sha1.New, []byte(ft.secretKey))

	mac.Write([]byte(sign))
	sum := mac.Sum(nil)

	s := base64.StdEncoding.EncodeToString(sum)
	log.Println(s)
	return s
}

func (ft *FCoin) placeOrder(orderType, orderSide, amount, price string, pair CurrencyPair) (*Order, error) {
	params := url.Values{}

	params.Set("side", orderSide)
	params.Set("amount", amount)
	//params.Set("price", price)
	params.Set("symbol", strings.ToLower(pair.AdaptUsdToUsdt().ToSymbol("")))

	switch orderType {
	case "LIMIT", "limit":
		params.Set("price", price)
		params.Set("type", "limit")
	case "MARKET", "market":
		params.Set("type", "market")
	}

	r, err := ft.doAuthenticatedRequest("POST", "orders", params)
	if err != nil {
		return nil, err
	}

	side := SELL
	if orderSide == "buy" {
		side = BUY
	}

	return &Order{
		Currency: pair,
		OrderID2: r.(string),
		Amount:   ToFloat64(amount),
		Price:    ToFloat64(price),
		Side:     TradeSide(side),
		Status:   ORDER_UNFINISH}, nil
}

func (ft *FCoin) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return ft.placeOrder("limit", "buy", amount, price, currency)
}

func (ft *FCoin) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return ft.placeOrder("limit", "sell", amount, price, currency)
}

func (ft *FCoin) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return ft.placeOrder("market", "buy", amount, price, currency)
}

func (ft *FCoin) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return ft.placeOrder("market", "sell", amount, price, currency)
}

func (ft *FCoin) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	uri := fmt.Sprintf("orders/%s/submit-cancel", orderId)
	_, err := ft.doAuthenticatedRequest("POST", uri, url.Values{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (ft *FCoin) toOrder(o map[string]interface{}, pair CurrencyPair) *Order {
	side := SELL
	if o["side"].(string) == "buy" {
		side = BUY
	}

	orderStatus := ORDER_UNFINISH
	switch o["state"].(string) {
	case "partial_filled":
		orderStatus = ORDER_PART_FINISH
	case "filled":
		orderStatus = ORDER_FINISH
	case "canceled", "partial_canceled":
		orderStatus = ORDER_CANCEL
	}

	return &Order{
		Currency:   pair,
		Side:       TradeSide(side),
		OrderID2:   o["id"].(string),
		Amount:     ToFloat64(o["amount"]),
		Price:      ToFloat64(o["price"]),
		DealAmount: ToFloat64(o["filled_amount"]),
		Status:     TradeStatus(orderStatus),
		Fee:        ToFloat64(o["fill_fees"]),
		OrderTime:  ToInt(o["created_at"])}
}

func (ft *FCoin) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	uri := fmt.Sprintf("orders/%s", orderId)
	r, err := ft.doAuthenticatedRequest("GET", uri, url.Values{})

	if err != nil {
		return nil, err
	}

	return ft.toOrder(r.(map[string]interface{}), currency), nil

}

func (ft *FCoin) GetOrdersList() {
	//path := API_URL + fmt.Sprintf(CANCEL_ORDER_API, strings.ToLower(currency.ToSymbol("")))

}

func (ft *FCoin) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToLower(currency.AdaptUsdToUsdt().ToSymbol("")))
	params.Set("states", "submitted")
	//params.Set("before", "1")
	//params.Set("after", "0")
	params.Set("limit", "100")

	r, err := ft.doAuthenticatedRequest("GET", "orders", params)
	if err != nil {
		return nil, err
	}

	var ords []Order

	for _, ord := range r.([]interface{}) {
		ords = append(ords, *ft.toOrder(ord.(map[string]interface{}), currency))
	}

	return ords, nil
}

func (ft *FCoin) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implement")
}

func (ft *FCoin) GetAccount() (*Account, error) {

	r, err := ft.doAuthenticatedRequest("GET", "accounts/balance", url.Values{})
	if err != nil {
		return nil, err
	}

	acc := new(Account)
	acc.SubAccounts = make(map[Currency]SubAccount)
	acc.Exchange = ft.GetExchangeName()

	balances := r.([]interface{})
	for _, v := range balances {
		vv := v.(map[string]interface{})
		currency := NewCurrency(vv["currency"].(string), "")
		acc.SubAccounts[currency] = SubAccount{
			Currency:     currency,
			Amount:       ToFloat64(vv["available"]),
			ForzenAmount: ToFloat64(vv["frozen"]),
		}
	}

	return acc, nil

}

func (ft *FCoin) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录
func (ft *FCoin) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}
