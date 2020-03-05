package fmex

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
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

	baseUrl = "https://api.fmex.com"
)

type FMexSwap struct {
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

type OrderFilter struct {
	Range    string
	Symbol   string
	OffsetId string
	Limit    string
}

func NewFMexSwap(config *APIConfig) *FMexSwap {
	if config.Endpoint == "" {
		config.Endpoint = baseUrl
	}
	fm := &FMexSwap{baseUrl: config.Endpoint, accessKey: config.ApiKey, secretKey: config.ApiSecretKey, httpClient: config.HttpClient}
	fm.setTimeOffset()
	return fm
}

func (fm *FMexSwap) SetBaseUri(uri string) {
	fm.baseUrl = uri
}

func (fm *FMexSwap) GetExchangeName() string {
	return FMEX
}

/**
 *获取交割预估价
 */
func (fm *FMexSwap) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	panic("not supported.")
}

/**
 * 期货行情
 * @param currency_pair   btc_usd:比特币    ltc_usd :莱特币
 * @param contractType  合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 */
func (fm *FMexSwap) GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error) {
	respmap, err := HttpGet(fm.httpClient, fm.baseUrl+fmt.Sprintf("/v2/market/ticker/%s",
		adaptContractType(currencyPair)))

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

/**
 * 期货深度
 * @param currencyPair  btc_usd:比特币    ltc_usd :莱特币
 * @param contractType  合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 * @param size 获取深度档数
 * @return
 */

func (fm *FMexSwap) GetFutureDepth(currency CurrencyPair, contractType string, size int) (*Depth, error) {
	var uri string
	if size <= 20 {
		uri = fmt.Sprintf("/v2/market/depth/L20/%s", adaptContractType(currency))
	} else {
		uri = fmt.Sprintf("/v2/market/depth/L150/%s", adaptContractType(currency))
	}
	respmap, err := HttpGet(fm.httpClient, fm.baseUrl+uri)
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
	depth.Pair = currency
	depth.UTime = time.Unix(0, ts*1000000)

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

func (fm *FMexSwap) GetTrades(contract_type string, currencyPair CurrencyPair, since int64) ([]Trade, error) {
	var uri = "/v2/market/trades/" + adaptContractType(currencyPair)
	respmap, err := HttpGet(fm.httpClient, fm.baseUrl+uri)
	if err != nil {
		return nil, err
	}

	if respmap["status"].(float64) != 0 {
		return nil, errors.New(respmap["msg"].(string))
	}

	datamap := respmap["data"].([]interface{})

	trades := make([]Trade, 0)
	for _, v := range datamap {
		vv := v.(map[string]interface{})
		side := BUY
		if vv["side"] == "sell" {
			side = SELL
		}
		trades = append(trades, Trade{
			Tid:    ToInt64(vv["id"]),
			Type:   side,
			Amount: ToFloat64(vv["amount"]),
			Price:  ToFloat64(vv["price"]),
			Date:   ToInt64(vv["ts"]),
			Pair:   currencyPair,
		})
	}

	return trades, nil

}

/**
 * 期货指数
 * @param currencyPair   btc_usd:比特币    ltc_usd :莱特币
 */
func (fm *FMexSwap) GetFutureIndex(currencyPair CurrencyPair) (float64, error) {
	var uri = "/v2/market/indexes"
	respmap, err := HttpGet(fm.httpClient, fm.baseUrl+uri)
	if err != nil {
		return 0.0, err
	}

	if respmap["status"].(float64) != 0 {
		return 0.0, errors.New(respmap["msg"].(string))
	}
	pair := "." + adaptCurrencyPair(currencyPair).ToLower().ToSymbol("") + "_spot"
	datamap := respmap["data"].(map[string]interface{})
	spot := datamap[pair].([]interface{})
	return ToFloat64(spot[1]), nil

}

/**
 *全仓账户
 */
func (fm *FMexSwap) GetFutureUserinfo() (*FutureAccount, error) {
	r, err := fm.doAuthenticatedRequest("GET", "/v3/contracts/accounts", url.Values{})
	if err != nil {
		return nil, err
	}
	acc := new(FutureAccount)
	acc.FutureSubAccounts = make(map[Currency]FutureSubAccount)

	balances := r.(map[string]interface{})
	for k, v := range balances {
		vv := v.([]interface{})
		currency := NewCurrency(k, "")
		acc.FutureSubAccounts[currency] = FutureSubAccount{
			Currency:      currency,
			AccountRights: ToFloat64(vv[0]),
			KeepDeposit:   ToFloat64(vv[2]),
			ProfitReal:    0,
			ProfitUnreal:  0,
			RiskRate:      0,
		}
	}
	return acc, nil

}

func (fm *FMexSwap) MarginTransferOut(currency Currency, amount float64) (bool, error) {
	params := url.Values{}
	params.Set("currency", strings.ToLower(currency.String()))
	params.Set("amount", fmt.Sprint(amount))
	params.Set("transferFrom", "CONTRACTS")
	params.Set("transferFrom", "WALLET")
	_, err := fm.doAuthenticatedRequest("POST", "/v3/contracts/transfer/out/request", params)
	if err != nil {
		return false, err
	}

	return true, nil
}

/**
 * @deprecated
 * 期货下单
 * @param currencyPair   btc_usd:比特币    ltc_usd :莱特币
 * @param contractType   合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 * @param price  价格
 * @param amount  委托数量
 * @param openType   1:开多   2:开空   3:平多   4:平空
 * @param matchPrice  是否为对手价 0:不是    1:是   ,当取值为1时,price无效
 */
func (fm *FMexSwap) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice, leverRate int) (string, error) {
	params := url.Values{}

	params.Set("source", "goex")
	params.Set("symbol", adaptContractType(currencyPair))

	switch openType {
	case OPEN_BUY, CLOSE_SELL:
		params.Set("direction", "LONG")
	case OPEN_SELL, CLOSE_BUY:
		params.Set("direction", "SHORT")
	}
	if matchPrice == 0 {
		params.Set("type", "LIMIT")
		params.Set("price", price)
	} else {
		params.Set("type", "MARKET")
	}
	params.Set("quantity", amount)

	r, err := fm.doAuthenticatedRequest("POST", "/v3/contracts/orders", params)
	if err != nil {
		return "", err
	}
	data := r.(map[string]interface{})
	return fmt.Sprintf("%d", int64(data["id"].(float64))), nil
}

type OrderParam struct {
	Currency         CurrencyPair
	Type             OrderType
	Direction        int
	Price            string
	Amount           string
	OType            int    //1：开多 2：开空 3：平多 4： 平空
	TriggerOn        string //订单触发价格，如果不填，则立刻执行
	TrailingDistance string //止盈止损订单触发距离，如果不填，则不会按止盈止损执行
	IsFOK            bool   //是否设置FOK订单
	IsIOC            bool   //是否设置IOC订单
	IsPostOnly       bool   //是否设置PostOnly订单
	IsHidden         bool   //是否设置Hidden订单
	IsReduceOnly     bool   //是否设置ReduceOnly订单
}

func (fm *FMexSwap) orderParamCheck(param *OrderParam) error {
	if param.TrailingDistance != "" && param.TriggerOn != "" {
		return errors.New("triggerOn and trailingDistance couldn't be filled at the same time")
	}

	if param.IsFOK && (param.IsIOC || param.IsPostOnly || param.IsHidden || param.IsReduceOnly) {
		return errors.New("IsIOC is true, coundn't set IsFOK | IsPostOnly | IsHidden | IsReduceOnly")
	}
	if param.IsIOC && (param.IsFOK || param.IsPostOnly || param.IsHidden || param.IsReduceOnly) {
		return errors.New("IsIOC is true, coundn't set IsFOK | IsPostOnly | IsHidden | IsReduceOnly")
	}
	if param.IsPostOnly && (param.IsFOK || param.IsIOC || param.IsReduceOnly) {
		return errors.New("IsPostOnly is true, coundn't set IsFOK | IsIOC | IsReduceOnly")
	}
	if param.IsHidden && (param.IsFOK || param.IsIOC) {
		return errors.New("IsHidden is true, coundn't set IsFOK | IsIOC ")
	}
	//	只有MARKET订单与IOC订单可以设置reduceOnly=true。

	return nil
}

func (fm *FMexSwap) PlaceFutureOrder2(ord *OrderParam) (string, error) {
	err := fm.orderParamCheck(ord)
	if err != nil {
		return "", err
	}

	params := url.Values{}

	params.Set("source", "goex")
	params.Set("symbol", adaptContractType(ord.Currency))

	switch ord.Direction {
	case OPEN_BUY, CLOSE_SELL:
		params.Set("direction", "LONG")
	case OPEN_SELL, CLOSE_BUY:
		params.Set("direction", "SHORT")
	}
	if ord.Type == ORDER_TYPE_LIMIT {
		params.Set("type", "LIMIT")
		params.Set("price", ord.Price)
	} else {
		params.Set("type", "MARKET")
	}
	params.Set("quantity", ord.Amount)

	if ord.TriggerOn != "" {
		params.Set("trigger_on", ord.TriggerOn)
	}
	if ord.TrailingDistance != "" {
		params.Set("trailing_distance", ord.TrailingDistance)
	}

	params.Set("fill_or_kill", fmt.Sprint(ord.IsFOK))
	params.Set("immediate_or_cancel", fmt.Sprint(ord.IsIOC))
	params.Set("post_only", fmt.Sprint(ord.IsPostOnly))
	params.Set("hidden", fmt.Sprint(ord.IsHidden))
	params.Set("reduce_only", fmt.Sprint(ord.IsReduceOnly))

	r, err := fm.doAuthenticatedRequest("POST", "/v3/contracts/orders", params)
	if err != nil {
		return "", err
	}
	data, isOk := r.(map[string]interface{})
	if !isOk {
		return "", errors.New(fmt.Sprintf("PlaceFutureOrder2 UNKNOW:%v", r))
	}

	return fmt.Sprintf("%d", int64(data["id"].(float64))), nil
}

/**
 * 取消订单
 * @param symbol   btc_usd:比特币    ltc_usd :莱特币
 * @param contractType    合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 * @param orderId   订单ID

 */
func (fm *FMexSwap) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error) {
	uri := fmt.Sprintf("/v3/contracts/orders/%s/cancel", orderId)
	_, err := fm.doAuthenticatedRequest("POST", uri, url.Values{})
	if err != nil {
		return false, err
	}
	return true, nil
}

/**
 * 用户持仓查询
 * @param symbol   btc_usd:比特币    ltc_usd :莱特币
 * @param contractType   合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 * @return
 */
func (fm *FMexSwap) GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	r, err := fm.doAuthenticatedRequest("GET", "/v3/contracts/positions", url.Values{})
	if err != nil {
		return nil, err
	}

	data := r.(map[string]interface{})
	result := data["results"].([]interface{})
	var positions []FuturePosition
	for _, info := range result {
		cont := info.(map[string]interface{})
		if !cont["closed"].(bool) {
			p := FuturePosition{
				CreateDate:     int64(cont["updated_at"].(float64)),
				LeverRate:      int(cont["leverage"].(float64)),
				Symbol:         currencyPair,
				ContractId:     int64(cont["user_id"].(float64)),
				ForceLiquPrice: cont["liquidation_price"].(float64),
			}
			if cont["direction"] == "long" {
				p.BuyAmount = cont["quantity"].(float64)
				p.BuyPriceAvg = cont["entry_price"].(float64)
				p.BuyPriceCost = cont["margin"].(float64)
				p.BuyProfitReal = cont["realized_pnl"].(float64)
			} else {
				p.SellAmount = cont["quantity"].(float64)
				p.SellPriceAvg = cont["entry_price"].(float64)
				p.SellPriceCost = cont["margin"].(float64)
				p.SellProfitReal = cont["realized_pnl"].(float64)
			}
			positions = append(positions, p)
		}
	}
	return positions, nil
}

/**
 *获取订单信息
 */
func (fm *FMexSwap) GetFutureOrders(orderIds []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	if len(orderIds) == 0 {
		return nil, errors.New("orderIds is empty")
	}
	orders := make([]FutureOrder, 0)
	for _, orderId := range orderIds {
		ord, err := fm.GetFutureOrder(orderId, currencyPair, contractType)
		if err != nil {
			return orders, err
		}
		orders = append(orders, *ord)
	}
	return orders, nil
}

func (fm *FMexSwap) GetFutureOrderHistory(filter *OrderFilter, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	param := url.Values{}
	if filter != nil {
		if filter.Range != "" {
			param.Set("range", filter.Range)
		}
		if filter.Symbol != "" {
			param.Set("symbol", filter.Symbol)
		}
		if filter.OffsetId != "" {
			param.Set("offsetId", filter.OffsetId)
		}
		if filter.Limit != "" {
			param.Set("limit", filter.Limit)
		}
	}
	r, err := fm.doAuthenticatedRequest("GET", "/v3/contracts/orders/closed", param)
	if err != nil {
		return nil, err
	}
	data := r.(map[string]interface{})
	result := data["results"].([]interface{})
	var orders []FutureOrder
	for _, info := range result {
		ord := fm.parseOrder(info)
		ord.Currency = currencyPair
		ord.ContractName = contractType
		orders = append(orders, ord)
	}
	return orders, nil
}

/**
 *获取单个订单信息
 */
func (fm *FMexSwap) GetFutureOrder(orderId string, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	r, err := fm.doAuthenticatedRequest("GET", "/v3/contracts/orders/"+orderId, url.Values{})
	if err != nil {
		return nil, err
	}
	data := r.(map[string]interface{})
	order := fm.parseOrder(data)
	order.Currency = currencyPair
	order.ContractName = contractType

	return &order, nil
}

func (fm *FMexSwap) parseOrderStatus(sts string) TradeStatus {
	orderStatus := ORDER_UNFINISH
	switch sts {
	case "PARTIAL_FILLED", "partial_filled":
		orderStatus = ORDER_PART_FINISH
	case "FULLY_FILLED", "fully_filled":
		orderStatus = ORDER_FINISH
	//case "STOP_CANCELLED":
	//	orderStatus = ORDER_CANCEL_ING
	case "FULLY_CANCELLED", "PARTIAL_CANCELLED", "STOP_CANCELLED", "fully_cancelled", "partial_cancelled", "stop_cancelled":
		orderStatus = ORDER_CANCEL
	}
	return orderStatus
}

func (fm *FMexSwap) parseOrder(ord interface{}) FutureOrder {
	order := ord.(map[string]interface{})
	return FutureOrder{
		OrderID2:   fmt.Sprintf("%d", ToInt64(order["id"])),
		Amount:     ToFloat64(order["quantity"]),
		DealAmount: ToFloat64(order["quantity"]) - ToFloat64(order["unfilled_quantity"]),
		Price:      ToFloat64(order["price"]),
		Fee:        ToFloat64(order["fee"]),
		OrderTime:  ToInt64(order["created_at"]),
		Status:     fm.parseOrderStatus(order["status"].(string)),
	}
}

/**
 *获取未完成订单信息
 */
func (fm *FMexSwap) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	r, err := fm.doAuthenticatedRequest("GET", "/v3/contracts/orders/open", url.Values{})
	if err != nil {
		return nil, err
	}
	data := r.(map[string]interface{})
	var orders []FutureOrder
	for _, info := range data["results"].([]interface{}) {
		ord := fm.parseOrder(info)
		ord.Currency = currencyPair
		ord.ContractName = contractType
		orders = append(orders, ord)
	}

	return orders, nil
}

/**
 *获取交易费
 */
func (fm *FMexSwap) GetFee() (float64, error) {
	panic("not supported.")
}

/**
 *获取每张合约价值
 */
func (fm *FMexSwap) GetContractValue(currencyPair CurrencyPair) (float64, error) {
	panic("not supported.")
}

/**
 *获取交割时间 星期(0,1,2,3,4,5,6)，小时，分，秒
 */
func (fm *FMexSwap) GetDeliveryTime() (int, int, int, int) {
	panic("not supported.")
}

/**
 * 获取K线数据
 */
func (fm *FMexSwap) GetKlineRecords(contract_type string, currency CurrencyPair, period, size, since int) ([]FutureKline, error) {
	uri := fmt.Sprintf("/v2/market/candles/%s/%s?limit=%d", _INERNAL_KLINE_PERIOD_CONVERTER[period], adaptContractType(currency), size)

	respmap, err := HttpGet(fm.httpClient, fm.baseUrl+uri)
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

	var klineRecords []FutureKline

	for _, record := range datamap {
		r := record.(map[string]interface{})
		klineRecords = append(klineRecords, FutureKline{Kline: &Kline{
			Pair:      currency,
			Timestamp: int64(ToInt(r["id"])),
			Open:      ToFloat64(r["open"]),
			Close:     ToFloat64(r["close"]),
			High:      ToFloat64(r["high"]),
			Low:       ToFloat64(r["low"]),
			Vol:       ToFloat64(r["quote_vol"])}})
	}
	return klineRecords, nil
}

func (fm *FMexSwap) GetServerTime() (int64, error) {
	respmap, err := HttpGet(fm.httpClient, fm.baseUrl+"/v2/public/server-time")
	if err != nil {
		return 0, err
	}
	stime := int64(ToInt(respmap["data"]))
	return stime, nil
}

func (fm *FMexSwap) setTimeOffset() error {
	stime, err := fm.GetServerTime()
	if err != nil {
		return err
	}
	st := time.Unix(stime/1000, 0)
	lt := time.Now()
	offset := st.Sub(lt).Seconds()
	fm.timeoffset = int64(offset)
	return nil
}

func (fm *FMexSwap) doAuthenticatedRequest(method, uri string, params url.Values) (interface{}, error) {

	timestamp := time.Now().Unix()*1000 + fm.timeoffset*1000
	sign := fm.buildSigned(method, fm.baseUrl+uri, timestamp, params)

	header := map[string]string{
		"FC-ACCESS-KEY":       fm.accessKey,
		"FC-ACCESS-SIGNATURE": sign,
		"FC-ACCESS-TIMESTAMP": fmt.Sprint(timestamp)}

	var (
		respmap map[string]interface{}
		err     error
	)

	switch method {
	case "GET":
		if len(params) != 0 {
			respmap, err = HttpGet2(fm.httpClient, fm.baseUrl+uri+"?"+params.Encode(), header)
		} else {
			respmap, err = HttpGet2(fm.httpClient, fm.baseUrl+uri, header)
		}
		if err != nil {
			return nil, err
		}

	case "POST":
		var parammap = make(map[string]string, 1)
		for k, v := range params {
			parammap[k] = v[0]
		}

		respbody, err := HttpPostForm4(fm.httpClient, fm.baseUrl+uri, parammap, header)
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

func (fm *FMexSwap) doAuthenticatedRequest2(method, uri string, params url.Values) (map[string]interface{}, error) {

	timestamp := time.Now().Unix()*1000 + fm.timeoffset*1000
	sign := fm.buildSigned(method, fm.baseUrl+uri, timestamp, params)

	header := map[string]string{
		"FC-ACCESS-KEY":       fm.accessKey,
		"FC-ACCESS-SIGNATURE": sign,
		"FC-ACCESS-TIMESTAMP": fmt.Sprint(timestamp)}

	var (
		respmap map[string]interface{}
		err     error
	)

	switch method {
	case "GET":
		if len(params) != 0 {
			respmap, err = HttpGet2(fm.httpClient, fm.baseUrl+uri+"?"+params.Encode(), header)
		} else {
			respmap, err = HttpGet2(fm.httpClient, fm.baseUrl+uri, header)
		}
		if err != nil {
			return nil, err
		}

	case "POST":
		var parammap = make(map[string]string, 1)
		for k, v := range params {
			parammap[k] = v[0]
		}

		respbody, err1 := HttpPostForm4(fm.httpClient, fm.baseUrl+uri, parammap, header)
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

func (fm *FMexSwap) buildSigned(httpmethod string, apiurl string, timestamp int64, para url.Values) string {

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

	mac := hmac.New(sha1.New, []byte(fm.secretKey))

	mac.Write([]byte(sign))
	sum := mac.Sum(nil)

	s := base64.StdEncoding.EncodeToString(sum)
	return s
}

func adaptCurrencyPair(pair CurrencyPair) CurrencyPair {
	if pair.CurrencyA.Eq(BCH) || pair.CurrencyA.Eq(BCC) {
		return NewCurrencyPair(NewCurrency("BCHABC", ""), pair.CurrencyB).AdaptUsdToUsdt()
	}

	if pair.CurrencyA.Symbol == "BSV" {
		return NewCurrencyPair(NewCurrency("BCHSV", ""), pair.CurrencyB).AdaptUsdToUsdt()
	}

	return pair.AdaptUsdtToUsd()
}

func adaptContractType(currencyPair CurrencyPair) string {
	return strings.ToLower(adaptCurrencyPair(currencyPair).ToSymbol("")) + "_p"
}
