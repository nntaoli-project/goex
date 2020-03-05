package fcoin

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
	. "github.com/nntaoli-project/goex/internal/logger"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unsafe"
)

const (
	SPOT     = "spot"
	ASSETS   = "assets"
	EXCHANGE = "exchange"
)

type FCoin struct {
	httpClient *http.Client
	baseUrl,
	accessKey,
	secretKey string
	timeoffset    int64
	tradeSymbols  []TradeSymbol
	tradeSymbols2 []TradeSymbol2
}

type TradeSymbol struct {
	Name          string `json:"name"`
	BaseCurrency  string `json:"base_currency"`
	QuoteCurrency string `json:"quote_currency"`
	PriceDecimal  int    `json:"price_decimal"`
	AmountDecimal int    `json:"amount_decimal"`
	Tradeable     bool   `json:"tradable"`
}

type TradeSymbol2 struct {
	TradeSymbol
	Category           string  `json:"category"`
	LeveragedMultiple  int     `json:"leveraged_multiple"`
	MarketOrderEnabled bool    `json:"market_order_enabled"`
	LimitAmountMin     float64 `json:"limit_amount_min"`
	LimitAmountMax     float64 `json:"limit_amount_max"`
	MainTag            string  `json:"main_tag"`
	DailyOpenAt        string  `json:"daily_open_at"`
	DailyCloseAt       string  `json:"daily_close_at"`
}

type Asset struct {
	Currency  Currency
	Avaliable float64
	Frozen    float64
	Finances  float64
	Lock      float64
	Total     float64
}

type RawTicker struct {
	Ticker
	SellAmount   float64
	BuyAmount    float64
	LastTradeVol float64
}

func NewFCoin(client *http.Client, apikey, secretkey string) *FCoin {
	return NewWithConfig(&APIConfig{
		HttpClient:   client,
		Endpoint:     "https://api.fcoin.com/v2/",
		ApiKey:       apikey,
		ApiSecretKey: secretkey,
	})
}

func NewWithConfig(c *APIConfig) *FCoin {
	if c.Endpoint == "" {
		c.Endpoint = "https://api.fcoin.com/v2/"
	}
	if !strings.HasSuffix(c.Endpoint, "/") {
		c.Endpoint = c.Endpoint + "/"
	}
	Log.Debug("endpoint=", c.Endpoint)
	fc := &FCoin{baseUrl: c.Endpoint, accessKey: c.ApiKey, secretKey: c.ApiSecretKey, httpClient: c.HttpClient}
	fc.setTimeOffset()
	var err error
	fc.tradeSymbols, err = fc.GetTradeSymbols()
	if len(fc.tradeSymbols) == 0 || err != nil {
		panic("[fcoin] trade symbol is empty, pls check connection...")
	}

	return fc
}

func (fc *FCoin) GetExchangeName() string {
	return FCOIN
}

func (fc *FCoin) GetServerTime() (int64, error) {
	respmap, err := HttpGet(fc.httpClient, fc.baseUrl+"public/server-time")
	if err != nil {
		return 0, err
	}
	stime := int64(ToInt(respmap["data"]))
	return stime, nil
}

func (fc *FCoin) setTimeOffset() error {
	stime, err := fc.GetServerTime()
	if err != nil {
		return err
	}
	st := time.Unix(stime/1000, 0)
	lt := time.Now()
	offset := st.Sub(lt).Seconds()
	fc.timeoffset = int64(offset)
	return nil
}

func (fc *FCoin) GetTicker(currencyPair CurrencyPair) (*Ticker, error) {
	respmap, err := HttpGet(fc.httpClient, fc.baseUrl+fmt.Sprintf("market/ticker/%s",
		strings.ToLower(currencyPair.ToSymbol(""))))

	if err != nil {
		return nil, err
	}

	if respmap["status"].(float64) != 0 {
		return nil, errors.New(respmap["msg"].(string))
	}

	tick, ok := respmap["data"].(map[string]interface{})
	if !ok {
		return nil, API_ERR
	}

	tickmap, ok := tick["ticker"].([]interface{})
	if !ok {
		return nil, API_ERR
	}

	ticker := new(RawTicker)
	ticker.Pair = currencyPair
	ticker.Date = uint64(time.Now().UnixNano() / 1000000)
	ticker.Last = ToFloat64(tickmap[0])
	ticker.Vol = ToFloat64(tickmap[9])
	ticker.Low = ToFloat64(tickmap[8])
	ticker.High = ToFloat64(tickmap[7])
	ticker.Buy = ToFloat64(tickmap[2])
	ticker.Sell = ToFloat64(tickmap[4])
	ticker.SellAmount = ToFloat64(tickmap[5])
	ticker.BuyAmount = ToFloat64(tickmap[3])
	ticker.LastTradeVol = ToFloat64(tickmap[1])
	return (*Ticker)(unsafe.Pointer(ticker)), nil
}

func (fc *FCoin) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	var uri string
	if size <= 20 {
		uri = fmt.Sprintf("market/depth/L20/%s", strings.ToLower(currency.ToSymbol("")))
	} else {
		uri = fmt.Sprintf("market/depth/L150/%s", strings.ToLower(currency.ToSymbol("")))
	}

	respmap, err := HttpGet(fc.httpClient, fc.baseUrl+uri)
	if err != nil {
		return nil, err
	}

	if respmap["status"].(float64) != 0 {
		return nil, errors.New(respmap["msg"].(string))
	}

	datamap := respmap["data"].(map[string]interface{})

	bids, ok1 := datamap["bids"].([]interface{})
	asks, ok2 := datamap["asks"].([]interface{})

	if !ok1 || !ok2 {
		return nil, errors.New("depth error")
	}

	ts := ToInt64(datamap["ts"])
	depth := new(Depth)
	depth.UTime = time.Unix(0, ts*1000000)
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

	return depth, nil
}

func (fc *FCoin) doAuthenticatedRequest(method, uri string, params url.Values) (interface{}, error) {

	timestamp := time.Now().Unix()*1000 + fc.timeoffset*1000
	sign := fc.buildSigned(method, fc.baseUrl+uri, timestamp, params)

	header := map[string]string{
		"FC-ACCESS-KEY":       fc.accessKey,
		"FC-ACCESS-SIGNATURE": sign,
		"FC-ACCESS-TIMESTAMP": fmt.Sprint(timestamp)}

	var (
		respmap map[string]interface{}
		err     error
	)

	switch method {
	case "GET":
		respmap, err = HttpGet2(fc.httpClient, fc.baseUrl+uri+"?"+params.Encode(), header)
		if err != nil {
			return nil, err
		}

	case "POST":
		var parammap = make(map[string]string, 1)
		for k, v := range params {
			parammap[k] = v[0]
		}

		respbody, err := HttpPostForm4(fc.httpClient, fc.baseUrl+uri, parammap, header)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(respbody, &respmap)
	}
	if ToInt(respmap["status"]) != 0 {
		return nil, errors.New(respmap["msg"].(string))
	}

	return respmap["data"], err
}

func (fc *FCoin) doAuthenticatedRequest2(method, uri string, params url.Values) (map[string]interface{}, error) {

	timestamp := time.Now().Unix()*1000 + fc.timeoffset*1000
	sign := fc.buildSigned(method, fc.baseUrl+uri, timestamp, params)

	header := map[string]string{
		"FC-ACCESS-KEY":       fc.accessKey,
		"FC-ACCESS-SIGNATURE": sign,
		"FC-ACCESS-TIMESTAMP": fmt.Sprint(timestamp)}

	var (
		respmap map[string]interface{}
		err     error
	)

	switch method {
	case "GET":
		respmap, err = HttpGet2(fc.httpClient, fc.baseUrl+uri+"?"+params.Encode(), header)
		if err != nil {
			return nil, err
		}

	case "POST":
		var parammap = make(map[string]string, 1)
		for k, v := range params {
			parammap[k] = v[0]
		}

		respbody, err1 := HttpPostForm4(fc.httpClient, fc.baseUrl+uri, parammap, header)
		if err1 != nil {
			return nil, err1
		}
		err = json.Unmarshal(respbody, &respmap)
		if err != nil {
			return nil, err
		}
	}

	return respmap, err
}

func (fc *FCoin) buildSigned(httpmethod string, apiurl string, timestamp int64, para url.Values) string {

	var (
		param = ""
		err   error
	)

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

	signStr2, err := url.QueryUnescape(signStr) // 不需要编码
	if err != nil {
		signStr2 = signStr
	}

	sign := base64.StdEncoding.EncodeToString([]byte(signStr2))

	mac := hmac.New(sha1.New, []byte(fc.secretKey))

	mac.Write([]byte(sign))
	sum := mac.Sum(nil)

	s := base64.StdEncoding.EncodeToString(sum)
	return s
}

// orderType:limit, market, fok, ioc

func (fc *FCoin) PlaceOrder(orderType, orderSide, amount, price string, pair CurrencyPair, isMargin bool) (*Order, error) {
	params := url.Values{}

	params.Set("side", orderSide)
	params.Set("amount", amount)
	params.Set("symbol", strings.ToLower(pair.AdaptUsdToUsdt().ToSymbol("")))
	if "market" != orderType {
		params.Set("price", price)
	}

	if isMargin {
		params.Set("account_type", "margin")
	}
	params.Set("type", strings.ToLower(orderType))
	params.Set("exchange", "main")

	r, err := fc.doAuthenticatedRequest("POST", "orders", params)
	if err != nil {
		return nil, err
	}

	side := SELL
	if orderSide == "buy" {
		side = BUY
	}

	return &Order{
		Currency:  pair,
		OrderID2:  r.(string),
		Amount:    ToFloat64(amount),
		Price:     ToFloat64(price),
		Side:      TradeSide(side),
		Status:    ORDER_UNFINISH,
		OrderTime: int(time.Now().UnixNano() / 1000000),
		//OrderType: int(orderType),
	}, nil
}

func (fc *FCoin) IocBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return fc.PlaceOrder("ioc", "buy", amount, price, currency, false)
}

func (fc *FCoin) IocSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return fc.PlaceOrder("ioc", "sell", amount, price, currency, false)
}

func (fc *FCoin) FokBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return fc.PlaceOrder("fok", "buy", amount, price, currency, false)
}

func (fc *FCoin) FokSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return fc.PlaceOrder("fok", "sell", amount, price, currency, false)
}

func (fc *FCoin) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return fc.PlaceOrder("limit", "buy", amount, price, currency, false)
}

func (fc *FCoin) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return fc.PlaceOrder("limit", "sell", amount, price, currency, false)
}

func (fc *FCoin) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return fc.PlaceOrder("market", "buy", amount, price, currency, false)
}

func (fc *FCoin) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return fc.PlaceOrder("market", "sell", amount, price, currency, false)
}

func (fc *FCoin) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	uri := fmt.Sprintf("orders/%s/submit-cancel", orderId)
	_, err := fc.doAuthenticatedRequest("POST", uri, url.Values{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (fc *FCoin) toOrder(o map[string]interface{}, pair CurrencyPair) *Order {
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
	case "pending_cancel":
		orderStatus = ORDER_CANCEL_ING
	case "canceled", "partial_canceled":
		orderStatus = ORDER_CANCEL
	}
	var fees float64
	refund := ToFloat64(o["fees_income"])
	fee := ToFloat64(o["fill_fees"])
	if fee == 0 {
		fees = -refund
	} else {
		fees = fee
	}
	avgPrice := 0.0
	dealAmount := ToFloat64(o["filled_amount"])
	executeValue := ToFloat64(o["executed_value"])
	if dealAmount != 0.0 {
		avgPrice = executeValue / dealAmount
	}
	return &Order{
		Currency:   pair,
		Side:       TradeSide(side),
		OrderID2:   o["id"].(string),
		Amount:     ToFloat64(o["amount"]),
		Price:      ToFloat64(o["price"]),
		DealAmount: dealAmount,
		Status:     TradeStatus(orderStatus),
		Type:       o["type"].(string),
		Fee:        fees,
		AvgPrice:   avgPrice,
		OrderTime:  ToInt(o["created_at"])}
}

func (fc *FCoin) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	uri := fmt.Sprintf("orders/%s", orderId)
	r, err := fc.doAuthenticatedRequest("GET", uri, url.Values{})

	if err != nil {
		return nil, err
	}

	return fc.toOrder(r.(map[string]interface{}), currency), nil

}

func (fc *FCoin) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToLower(currency.AdaptUsdToUsdt().ToSymbol("")))
	params.Set("states", "submitted,partial_filled")
	//params.Set("before", "1")
	//params.Set("after", "0")
	params.Set("limit", "100")

	r, err := fc.doAuthenticatedRequest("GET", "orders", params)
	if err != nil {
		return nil, err
	}

	var ords []Order

	for _, ord := range r.([]interface{}) {
		ords = append(ords, *fc.toOrder(ord.(map[string]interface{}), currency))
	}

	return ords, nil
}

func (fc *FCoin) getAfterTimeOrderHistorys(currency CurrencyPair, times time.Time) ([]Order, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToLower(currency.AdaptUsdToUsdt().ToSymbol("")))
	params.Set("states", "filled")
	params.Set("after", fmt.Sprint(times.Unix()*1000))
	params.Set("limit", "100")

	r, err := fc.doAuthenticatedRequest("GET", "orders", params)
	if err != nil {
		return nil, err
	}

	var ords []Order

	for _, ord := range r.([]interface{}) {
		ords = append(ords, *fc.toOrder(ord.(map[string]interface{}), currency))
	}

	return ords, nil
}
func (fc *FCoin) getBeforeTimeOrderHistorys(currency CurrencyPair, times time.Time) ([]Order, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToLower(currency.AdaptUsdToUsdt().ToSymbol("")))
	params.Set("states", "filled")
	params.Set("before", fmt.Sprint(times.Unix()*1000))
	params.Set("limit", "100")

	r, err := fc.doAuthenticatedRequest("GET", "orders", params)
	if err != nil {
		return nil, err
	}

	var ords []Order

	for _, ord := range r.([]interface{}) {
		ords = append(ords, *fc.toOrder(ord.(map[string]interface{}), currency))
	}

	return ords, nil
}

func (fc *FCoin) GetHoursOrderHistorys(currency CurrencyPair, start time.Time, hours int64) ([]Order, error) {
	ord1, _ := fc.getAfterTimeOrderHistorys(currency, start)
	ord2, _ := fc.getBeforeTimeOrderHistorys(currency, start.Add(time.Hour*time.Duration(hours)))
	ords := make([]Order, 0)
	for _, v1 := range ord1 {
		for _, v2 := range ord2 {
			if v1.OrderID2 == v2.OrderID2 {
				ords = append(ords, v1)
			}
		}
	}
	return ords, nil
}

func (fc *FCoin) GetDaysOrderHistorys(currency CurrencyPair, start time.Time, days int64) ([]Order, error) {
	ord1, _ := fc.getAfterTimeOrderHistorys(currency, start)
	ord2, _ := fc.getBeforeTimeOrderHistorys(currency, start.Add(time.Hour*24*time.Duration(days)))
	ords := make([]Order, 0)
	for _, v1 := range ord1 {
		for _, v2 := range ord2 {
			if v1.OrderID2 == v2.OrderID2 {
				ords = append(ords, v1)
			}
		}
	}
	if len(ords) == 0 && len(ord2) > 1 && len(ord1) > 1 {
		return nil, errors.New("more than 100 orders")
	}
	return ords, nil
}

func (fc *FCoin) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToLower(currency.AdaptUsdToUsdt().ToSymbol("")))
	params.Set("states", "partial_canceled,filled")
	//params.Set("before", "1")
	//params.Set("after", "0")
	params.Set("limit", fmt.Sprint(pageSize))

	r, err := fc.doAuthenticatedRequest("GET", "orders", params)
	if err != nil {
		return nil, err
	}
	var ords []Order

	for _, ord := range r.([]interface{}) {
		ords = append(ords, *fc.toOrder(ord.(map[string]interface{}), currency))
	}

	return ords, nil

}

func (fc *FCoin) GetOrderHistorys2(currency CurrencyPair, currentPage, pageSize int, states ...string) ([]Order, error) {
	sts := ""
	for i := 0; i < len(states); i++ {
		sts += states[i] + ","
	}
	sts = sts[:len(sts)-1]
	params := url.Values{}
	params.Set("symbol", strings.ToLower(currency.AdaptUsdToUsdt().ToSymbol("")))
	params.Set("states", sts)
	params.Set("limit", fmt.Sprint(pageSize))

	r, err := fc.doAuthenticatedRequest("GET", "orders", params)
	if err != nil {
		return nil, err
	}
	var ords []Order

	for _, ord := range r.([]interface{}) {
		ords = append(ords, *fc.toOrder(ord.(map[string]interface{}), currency))
	}

	return ords, nil
}

func (fc *FCoin) GetAccount() (*Account, error) {
	r, err := fc.doAuthenticatedRequest("GET", "accounts/balance", url.Values{})
	if err != nil {
		return nil, err
	}
	acc := new(Account)
	acc.SubAccounts = make(map[Currency]SubAccount)
	acc.Exchange = fc.GetExchangeName()

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

func (fc *FCoin) GetAssets() ([]Asset, error) {
	r, err := fc.doAuthenticatedRequest("GET", "assets/accounts/balance", url.Values{})
	if err != nil {
		return nil, err
	}
	assets := make([]Asset, 0)
	balances := r.([]interface{})
	for _, v := range balances {
		vv := v.(map[string]interface{})
		currency := NewCurrency(vv["currency"].(string), "")
		assets = append(assets, Asset{
			Currency:  currency,
			Avaliable: ToFloat64(vv["available"]),
			Frozen:    ToFloat64(vv["frozen"]),
			Finances:  ToFloat64(vv["demand_deposit"]),
			Lock:      ToFloat64(vv["lock_deposit"]),
			Total:     ToFloat64(vv["balance"]),
		})
	}
	return assets, nil
}

// from, to: assets, spot
func (fc *FCoin) AssetTransfer(currency Currency, amount, from, to string) (bool, error) {
	params := url.Values{}
	params.Set("currency", strings.ToLower(currency.String()))
	params.Set("amount", amount)
	_, err := fc.doAuthenticatedRequest("POST", fmt.Sprintf("assets/accounts/%s-to-%s", from, to), params)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (fc *FCoin) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	uri := fmt.Sprintf("market/candles/%s/%s?limit=%d", _INERNAL_KLINE_PERIOD_CONVERTER[period], strings.ToLower(currency.ToSymbol("")), size)

	respmap, err := HttpGet(fc.httpClient, fc.baseUrl+uri)
	if err != nil {
		return nil, err
	}

	if respmap["status"].(float64) != 0 {
		return nil, errors.New(respmap["msg"].(string))
	}

	datamap, isOk := respmap["data"].([]interface{})

	if !isOk {
		return nil, errors.New("kline error")
	}

	var klineRecords []Kline

	for _, record := range datamap {
		r := record.(map[string]interface{})
		klineRecords = append(klineRecords, Kline{
			Pair:      currency,
			Timestamp: int64(ToInt(r["id"])),
			Open:      ToFloat64(r["open"]),
			Close:     ToFloat64(r["close"]),
			High:      ToFloat64(r["high"]),
			Low:       ToFloat64(r["low"]),
			Vol:       ToFloat64(r["quote_vol"]),
		})
	}

	return klineRecords, nil
}

//非个人，整个交易所的交易记录
func (fc *FCoin) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}

//交易符号
func (fc *FCoin) GetTradeSymbols() ([]TradeSymbol, error) {
	respmap, err := HttpGet(fc.httpClient, fc.baseUrl+"public/symbols")
	if err != nil {
		return nil, err
	}

	if respmap["status"].(float64) != 0 {
		return nil, errors.New(respmap["msg"].(string))
	}

	datamap := respmap["data"].([]interface{})

	tradeSymbols := make([]TradeSymbol, 0)
	for _, v := range datamap {
		vv := v.(map[string]interface{})
		var symbol TradeSymbol
		symbol.Name = vv["name"].(string)
		symbol.BaseCurrency = vv["base_currency"].(string)
		symbol.QuoteCurrency = vv["quote_currency"].(string)
		symbol.PriceDecimal = int(vv["price_decimal"].(float64))
		symbol.AmountDecimal = int(vv["amount_decimal"].(float64))
		symbol.Tradeable = vv["tradable"].(bool)
		if symbol.Tradeable {
			tradeSymbols = append(tradeSymbols, symbol)
		}
	}
	return tradeSymbols, nil
}

func (fc *FCoin) GetTradeSymbol(currencyPair CurrencyPair) (*TradeSymbol, error) {
	if len(fc.tradeSymbols) == 0 {
		var err error
		fc.tradeSymbols, err = fc.GetTradeSymbols()
		if err != nil {
			return nil, err
		}
	}
	for k, v := range fc.tradeSymbols {
		if v.Name == strings.ToLower(currencyPair.ToSymbol("")) {
			return &fc.tradeSymbols[k], nil
		}
	}
	return nil, errors.New("symbol not found")
}

func (fc *FCoin) GetTradeSymbols2() ([]TradeSymbol2, error) {
	respmap, err := HttpGet(fc.httpClient, "https://www.fcoin.com/openapi/v2/symbols")
	if err != nil {
		return nil, err
	}

	if respmap["status"].(string) != "ok" {
		return nil, errors.New(respmap["msg"].(string))
	}

	datamap := respmap["data"].(map[string]interface{})
	symbols := datamap["symbols"].(map[string]interface{})
	tradeSymbols := make([]TradeSymbol2, 0)
	for _, v := range symbols {
		vv := v.(map[string]interface{})
		var symbol TradeSymbol2
		symbol.Name, _ = vv["symbol"].(string)
		symbol.BaseCurrency, _ = vv["base_currency"].(string)
		symbol.QuoteCurrency, _ = vv["quote_currency"].(string)
		symbol.PriceDecimal = ToInt(vv["price_decimal"])
		symbol.AmountDecimal = ToInt(vv["amount_decimal"])
		symbol.Tradeable, _ = vv["tradeable"].(bool)
		symbol.Category, _ = vv["category"].(string)
		symbol.LeveragedMultiple = ToInt(vv["leveraged_multiple"])
		symbol.MarketOrderEnabled, _ = vv["market_order_enabled"].(bool)
		symbol.LimitAmountMin = ToFloat64(vv["limit_amount_min"])
		symbol.LimitAmountMax = ToFloat64(vv["limit_amount_max"])
		symbol.MainTag, _ = vv["main_tag"].(string)
		symbol.DailyOpenAt, _ = vv["daily_open_at"].(string)
		symbol.DailyCloseAt, _ = vv["daily_close_at"].(string)

		if symbol.Tradeable {
			tradeSymbols = append(tradeSymbols, symbol)
		}
	}
	return tradeSymbols, nil
}

func (fc *FCoin) GetTradeSymbol2(currencyPair CurrencyPair) (*TradeSymbol2, error) {
	if len(fc.tradeSymbols2) == 0 {
		var err error
		fc.tradeSymbols2, err = fc.GetTradeSymbols2()
		if err != nil {
			return nil, err
		}
	}
	for k, v := range fc.tradeSymbols2 {
		if v.Name == strings.ToLower(currencyPair.ToSymbol("")) {
			return &fc.tradeSymbols2[k], nil
		}
	}
	return nil, errors.New("symbol not found")
}
