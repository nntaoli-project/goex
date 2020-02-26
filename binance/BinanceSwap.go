package binance

import (
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	baseUrl = "https://fapi.binance.com"
)

type BinanceSwap struct {
	Binance
}

func NewBinanceSwap(config *APIConfig) *BinanceSwap {
	if config.Endpoint == "" {
		config.Endpoint = baseUrl
	}
	bs := &BinanceSwap{
		Binance: Binance{
			baseUrl:    config.Endpoint,
			accessKey:  config.ApiKey,
			apiV1:      config.Endpoint + "/fapi/v1/",
			secretKey:  config.ApiSecretKey,
			httpClient: config.HttpClient,
		},
	}
	bs.setTimeOffset()
	return bs
}

func (bs *BinanceSwap) SetBaseUri(uri string) {
	bs.baseUrl = uri
}

func (bs *BinanceSwap) GetExchangeName() string {
	return BINANCE_SWAP
}

func (bs *BinanceSwap) Ping() bool {
	_, err := HttpGet(bs.httpClient, bs.apiV1+"ping")
	if err != nil {
		return false
	}
	return true
}

func (bs *BinanceSwap) setTimeOffset() error {
	respmap, err := HttpGet(bs.httpClient, bs.apiV1+SERVER_TIME_URL)
	if err != nil {
		return err
	}

	stime := int64(ToInt(respmap["serverTime"]))
	st := time.Unix(stime/1000, 1000000*(stime%1000))
	lt := time.Now()
	offset := st.Sub(lt).Nanoseconds()
	bs.timeOffset = int64(offset)
	return nil
}

/**
 *获取交割预估价
 */
func (bs *BinanceSwap) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	panic("not supported.")
}

/**
 * 期货行情
 * @param currency_pair   btc_usd:比特币    ltc_usd :莱特币
 * @param contractType  合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 */
func (bs *BinanceSwap) GetFutureTicker(currency CurrencyPair, contractType string) (*Ticker, error) {
	currency2 := bs.adaptCurrencyPair(currency)
	tickerPriceUri := bs.apiV1 + "ticker/price?symbol=" + currency2.ToSymbol("")
	tickerBookUri := bs.apiV1 + "ticker/bookTicker?symbol=" + currency2.ToSymbol("")
	tickerPriceMap := make(map[string]interface{})
	tickerBookeMap := make(map[string]interface{})
	var err1 error
	var err2 error
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		tickerPriceMap, err1 = HttpGet(bs.httpClient, tickerPriceUri)
	}()
	go func() {
		defer wg.Done()
		tickerBookeMap, err2 = HttpGet(bs.httpClient, tickerBookUri)
	}()
	wg.Wait()
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}

	var ticker Ticker
	ticker.Pair = currency
	ticker.Date = uint64(time.Now().UnixNano() / int64(time.Millisecond))
	ticker.Last = ToFloat64(tickerPriceMap["price"])
	ticker.Buy = ToFloat64(tickerBookeMap["bidPrice"])
	ticker.Sell = ToFloat64(tickerBookeMap["askPrice"])
	return &ticker, nil
}

/**
 * 期货深度
 * @param currencyPair  btc_usd:比特币    ltc_usd :莱特币
 * @param contractType  合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 * @param size 获取深度档数
 * @return
 */

func (bs *BinanceSwap) GetFutureDepth(currency CurrencyPair, contractType string, size int) (*Depth, error) {
	if size <= 5 {
		size = 5
	} else if size <= 10 {
		size = 10
	} else if size <= 20 {
		size = 20
	} else if size <= 50 {
		size = 50
	} else if size <= 100 {
		size = 100
	} else if size <= 500 {
		size = 500
	} else {
		size = 1000
	}
	currencyPair2 := bs.adaptCurrencyPair(currency)

	apiUrl := fmt.Sprintf(bs.apiV1+DEPTH_URI, currencyPair2.ToSymbol(""), size)
	resp, err := HttpGet(bs.httpClient, apiUrl)
	if err != nil {
		return nil, err
	}

	if _, isok := resp["code"]; isok {
		return nil, errors.New(resp["msg"].(string))
	}

	bids := resp["bids"].([]interface{})
	asks := resp["asks"].([]interface{})

	depth := new(Depth)
	depth.Pair = currency
	depth.UTime = time.Now()
	n := 0
	for _, bid := range bids {
		_bid := bid.([]interface{})
		amount := ToFloat64(_bid[1])
		price := ToFloat64(_bid[0])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.BidList = append(depth.BidList, dr)
		n++
		if n == size {
			break
		}
	}

	n = 0
	for _, ask := range asks {
		_ask := ask.([]interface{})
		amount := ToFloat64(_ask[1])
		price := ToFloat64(_ask[0])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.AskList = append(depth.AskList, dr)
		n++
		if n == size {
			break
		}
	}

	return depth, nil
}

func (bs *BinanceSwap) GetTrades(contractType string, currencyPair CurrencyPair, since int64) ([]Trade, error) {
	param := url.Values{}
	param.Set("symbol", bs.adaptCurrencyPair(currencyPair).ToSymbol(""))
	param.Set("limit", "500")
	if since > 0 {
		param.Set("fromId", strconv.Itoa(int(since)))
	}
	apiUrl := bs.apiV1 + "historicalTrades?" + param.Encode()
	resp, err := HttpGet3(bs.httpClient, apiUrl, map[string]string{
		"X-MBX-APIKEY": bs.accessKey})
	if err != nil {
		return nil, err
	}

	var trades []Trade
	for _, v := range resp {
		m := v.(map[string]interface{})
		ty := SELL
		if m["isBuyerMaker"].(bool) {
			ty = BUY
		}
		trades = append(trades, Trade{
			Tid:    ToInt64(m["id"]),
			Type:   ty,
			Amount: ToFloat64(m["qty"]),
			Price:  ToFloat64(m["price"]),
			Date:   ToInt64(m["time"]),
			Pair:   currencyPair,
		})
	}

	return trades, nil

}

/**
 * 期货指数
 * @param currencyPair   btc_usd:比特币    ltc_usd :莱特币
 */
func (bs *BinanceSwap) GetFutureIndex(currencyPair CurrencyPair) (float64, error) {
	respmap, err := HttpGet(bs.httpClient, bs.apiV1+"premiumIndex?symbol="+bs.adaptCurrencyPair(currencyPair).ToSymbol(""))
	if err != nil {
		return 0.0, err
	}

	return ToFloat64(respmap["markPrice"]), nil
}

/**
 *全仓账户
 */
func (bs *BinanceSwap) GetFutureUserinfo() (*FutureAccount, error) {
	params := url.Values{}
	bs.buildParamsSigned(&params)
	path := bs.apiV1 + ACCOUNT_URI + params.Encode()
	respmap, err := HttpGet2(bs.httpClient, path, map[string]string{"X-MBX-APIKEY": bs.accessKey})
	if err != nil {
		return nil, err
	}
	if _, isok := respmap["code"]; isok == true {
		return nil, errors.New(respmap["msg"].(string))
	}
	acc := &FutureAccount{}
	acc.FutureSubAccounts = make(map[Currency]FutureSubAccount)

	balances := respmap["assets"].([]interface{})
	for _, v := range balances {
		vv := v.(map[string]interface{})
		currency := NewCurrency(vv["asset"].(string), "").AdaptBccToBch()
		acc.FutureSubAccounts[currency] = FutureSubAccount{
			Currency:      currency,
			AccountRights: ToFloat64(vv["walletBalance"]),
			KeepDeposit:   ToFloat64(vv["marginBalance"]),
			ProfitUnreal:  ToFloat64(vv["unrealizedProfit"]),
			RiskRate:      ToFloat64(vv["unrealizedProfit"]),
		}
	}
	return acc, nil
}

// transferType - 1: 现货账户向合约账户划转 2: 合约账户向现货账户划转
func (bs *BinanceSwap) Transfer(currency Currency, transferType int, amount float64) (int64, error) {
	params := url.Values{}

	params.Set("currency", currency.String())
	params.Set("amount", fmt.Sprint(amount))
	params.Set("type", strconv.Itoa(transferType))
	uri := "https://api.binance.com/sapi/v1/futures/transfer"
	bs.buildParamsSigned(&params)

	resp, err := HttpPostForm2(bs.httpClient, uri, params,
		map[string]string{"X-MBX-APIKEY": bs.accessKey})
	if err != nil {
		return 0, err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return 0, err
	}

	return ToInt64(respmap["tranId"]), nil
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
func (bs *BinanceSwap) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice, leverRate int) (string, error) {

	pair := bs.adaptCurrencyPair(currencyPair)
	path := bs.apiV1 + ORDER_URI
	params := url.Values{}
	params.Set("symbol", pair.ToSymbol(""))
	params.Set("quantity", amount)

	switch openType {
	case OPEN_BUY, CLOSE_SELL:
		params.Set("side", "BUY")
	case OPEN_SELL, CLOSE_BUY:
		params.Set("side", "SELL")
	}
	if matchPrice == 0 {
		params.Set("type", "LIMIT")
		params.Set("price", price)
		params.Set("timeInForce", "GTC")
	} else {
		params.Set("type", "MARKET")
	}

	bs.buildParamsSigned(&params)
	resp, err := HttpPostForm2(bs.httpClient, path, params,
		map[string]string{"X-MBX-APIKEY": bs.accessKey})
	if err != nil {
		return "", err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return "", err
	}

	orderId := ToInt(respmap["orderId"])
	if orderId <= 0 {
		return "", errors.New(string(resp))
	}
	return strconv.Itoa(orderId), nil
}

/**
 * 取消订单
 * @param symbol   btc_usd:比特币    ltc_usd :莱特币
 * @param contractType    合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 * @param orderId   订单ID

 */
func (bs *BinanceSwap) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error) {
	currencyPair = bs.adaptCurrencyPair(currencyPair)
	path := bs.apiV1 + ORDER_URI
	params := url.Values{}
	params.Set("symbol", bs.adaptCurrencyPair(currencyPair).ToSymbol(""))
	params.Set("orderId", orderId)

	bs.buildParamsSigned(&params)

	resp, err := HttpDeleteForm(bs.httpClient, path, params, map[string]string{"X-MBX-APIKEY": bs.accessKey})

	if err != nil {
		return false, err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return false, err
	}

	orderIdCanceled := ToInt(respmap["orderId"])
	if orderIdCanceled <= 0 {
		return false, errors.New(string(resp))
	}

	return true, nil
}

func (bs *BinanceSwap) FutureCancelAllOrders(currencyPair CurrencyPair, contractType string) (bool, error) {
	currencyPair = bs.adaptCurrencyPair(currencyPair)
	path := bs.apiV1 + "allOpenOrders"
	params := url.Values{}
	params.Set("symbol", bs.adaptCurrencyPair(currencyPair).ToSymbol(""))

	bs.buildParamsSigned(&params)

	resp, err := HttpDeleteForm(bs.httpClient, path, params, map[string]string{"X-MBX-APIKEY": bs.accessKey})

	if err != nil {
		return false, err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return false, err
	}

	if ToInt(respmap["code"]) != 200 {
		return false, errors.New(respmap["msg"].(string))
	}

	return true, nil
}

func (bs *BinanceSwap) FutureCancelOrders(currencyPair CurrencyPair, contractType string, orderIdList []string) (bool, error) {
	currencyPair = bs.adaptCurrencyPair(currencyPair)
	path := bs.apiV1 + "batchOrders"

	if len(orderIdList) == 0 {
		return false, errors.New("list is empty, no order will be cancel")
	}
	list, _ := json.Marshal(orderIdList)
	params := url.Values{}

	params.Set("symbol", bs.adaptCurrencyPair(currencyPair).ToSymbol(""))
	params.Set("orderIdList", string(list))

	bs.buildParamsSigned(&params)

	resp, err := HttpDeleteForm(bs.httpClient, path, params, map[string]string{"X-MBX-APIKEY": bs.accessKey})

	if err != nil {
		return false, err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return false, err
	}

	if ToInt(respmap["code"]) != 200 {
		return false, errors.New(respmap["msg"].(string))
	}

	return true, nil
}

/**
 * 用户持仓查询
 * @param symbol   btc_usd:比特币    ltc_usd :莱特币
 * @param contractType   合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 * @return
 */
func (bs *BinanceSwap) GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	currencyPair1 := bs.adaptCurrencyPair(currencyPair)

	params := url.Values{}
	bs.buildParamsSigned(&params)
	path := bs.apiV1 + "positionRisk?" + params.Encode()

	result, err := HttpGet3(bs.httpClient, path, map[string]string{"X-MBX-APIKEY": bs.accessKey})

	if err != nil {
		return nil, err
	}

	var positions []FuturePosition
	for _, info := range result {

		cont := info.(map[string]interface{})
		if cont["symbol"] != currencyPair1.ToSymbol("") {
			continue
		}
		p := FuturePosition{
			LeverRate:      ToInt(cont["leverage"]),
			Symbol:         currencyPair,
			ForceLiquPrice: ToFloat64(cont["liquidationPrice"]),
		}
		amount := ToFloat64(cont["positionAmt"])
		price := ToFloat64(cont["entryPrice"])
		upnl := ToFloat64(cont["unRealizedProfit"])
		if amount > 0 {
			p.BuyAmount = amount
			p.BuyPriceAvg = price
			p.BuyProfitReal = upnl
		} else if amount < 0 {
			p.SellAmount = amount
			p.SellPriceAvg = price
			p.SellProfitReal = upnl
		}
		positions = append(positions, p)
	}
	return positions, nil
}

/**
 *获取订单信息
 */
func (bs *BinanceSwap) GetFutureOrders(orderIds []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	if len(orderIds) == 0 {
		return nil, errors.New("orderIds is empty")
	}
	currencyPair1 := bs.adaptCurrencyPair(currencyPair)

	params := url.Values{}
	params.Set("symbol", currencyPair1.ToSymbol(""))
	bs.buildParamsSigned(&params)

	path := bs.apiV1 + "allOrders?" + params.Encode()

	result, err := HttpGet3(bs.httpClient, path, map[string]string{"X-MBX-APIKEY": bs.accessKey})

	if err != nil {
		return nil, err
	}

	orders := make([]FutureOrder, 0)
	for _, info := range result {

		_ord := info.(map[string]interface{})
		if _ord["symbol"].(string) != currencyPair1.ToSymbol("") {
			continue
		}
		orderId := ToInt(_ord["orderId"])
		ordId := strconv.Itoa(orderId)

		for _, id := range orderIds {
			if id == ordId {
				order := &FutureOrder{}
				order = bs.parseOrder(_ord)
				order.Currency = currencyPair
				orders = append(orders, *order)
				break
			}
		}
	}
	return orders, nil

}

/**
 *获取单个订单信息
 */
func (bs *BinanceSwap) GetFutureOrder(orderId string, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	currencyPair1 := bs.adaptCurrencyPair(currencyPair)

	params := url.Values{}
	params.Set("symbol", currencyPair1.ToSymbol(""))
	params.Set("orderId", orderId)
	bs.buildParamsSigned(&params)

	path := bs.apiV1 + "allOrders?" + params.Encode()

	result, err := HttpGet3(bs.httpClient, path, map[string]string{"X-MBX-APIKEY": bs.accessKey})

	if err != nil {
		return nil, err
	}

	order := &FutureOrder{}
	ordId, _ := strconv.Atoi(orderId)
	for _, info := range result {

		_ord := info.(map[string]interface{})
		if _ord["symbol"].(string) != currencyPair1.ToSymbol("") {
			continue
		}

		if ToInt(_ord["orderId"]) != ordId {
			continue
		}

		order = bs.parseOrder(_ord)
		order.Currency = currencyPair
		return order, nil
	}
	return nil, errors.New(fmt.Sprintf("not found order:%s", orderId))
}

func (bs *BinanceSwap) parseOrder(rsp map[string]interface{}) *FutureOrder {
	order := &FutureOrder{}
	order.Price = ToFloat64(rsp["price"])
	order.Amount = ToFloat64(rsp["origQty"])
	order.DealAmount = ToFloat64(rsp["executedQty"])
	order.AvgPrice = ToFloat64(rsp["avgPrice"])
	order.OrderTime = ToInt64(rsp["time"])

	status := rsp["status"].(string)
	order.Status = bs.parseOrderStatus(status)
	order.OrderID = ToInt64(rsp["orderId"])
	order.OrderID2 = strconv.Itoa(int(order.OrderID))
	order.OType = OPEN_BUY
	if rsp["side"].(string) == "SELL" {
		order.OType = OPEN_SELL
	}

	//GTC - Good Till Cancel 成交为止
	//IOC - Immediate or Cancel 无法立即成交(吃单)的部分就撤销
	//FOK - Fill or Kill 无法全部立即成交就撤销
	//GTX - Good Till Crossing 无法成为挂单方就撤销
	ot := rsp["timeInForce"].(string)
	switch ot {
	case "GTC":
		order.OrderType = ORDER_FEATURE_LIMIT
	case "IOC":
		order.OrderType = ORDER_FEATURE_IOC
	case "FOK":
		order.OrderType = ORDER_FEATURE_FOK
	case "GTX":
		order.OrderType = ORDER_FEATURE_IOC

	}

	//LIMIT 限价单
	//MARKET 市价单
	//STOP 止损限价单
	//STOP_MARKET 止损市价单
	//TAKE_RPOFIT 止盈限价单
	//TAKE_RPOFIT_MARKET 止盈市价单

	return order
}

func (bs *BinanceSwap) parseOrderStatus(sts string) TradeStatus {
	orderStatus := ORDER_UNFINISH
	switch sts {
	case "PARTIALLY_FILLED", "partially_filled":
		orderStatus = ORDER_PART_FINISH
	case "FILLED", "filled":
		orderStatus = ORDER_FINISH
	case "CANCELED", "REJECTED", "EXPIRED":
		orderStatus = ORDER_CANCEL
	}
	return orderStatus
}

/**
 *获取未完成订单信息
 */
func (bs *BinanceSwap) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	currencyPair1 := bs.adaptCurrencyPair(currencyPair)

	params := url.Values{}
	params.Set("symbol", currencyPair1.ToSymbol(""))
	bs.buildParamsSigned(&params)

	path := bs.apiV1 + "openOrders?" + params.Encode()

	result, err := HttpGet3(bs.httpClient, path, map[string]string{"X-MBX-APIKEY": bs.accessKey})

	if err != nil {
		return nil, err
	}

	orders := make([]FutureOrder, 0)
	for _, info := range result {

		_ord := info.(map[string]interface{})
		if _ord["symbol"].(string) != currencyPair1.ToSymbol("") {
			continue
		}
		order := &FutureOrder{}
		order = bs.parseOrder(_ord)
		order.Currency = currencyPair
		orders = append(orders, *order)
	}
	return orders, nil
}

/**
 *获取交易费
 */
func (bs *BinanceSwap) GetFee() (float64, error) {
	panic("not supported.")
}

/**
 *获取每张合约价值
 */
func (bs *BinanceSwap) GetContractValue(currencyPair CurrencyPair) (float64, error) {
	panic("not supported.")
}

/**
 *获取交割时间 星期(0,1,2,3,4,5,6)，小时，分，秒
 */
func (bs *BinanceSwap) GetDeliveryTime() (int, int, int, int) {
	panic("not supported.")
}

/**
 * 获取K线数据
 */
func (bs *BinanceSwap) GetKlineRecords(contractType string, currency CurrencyPair, period, size, since int) ([]FutureKline, error) {
	currency2 := bs.adaptCurrencyPair(currency)
	params := url.Values{}
	params.Set("symbol", currency2.ToSymbol(""))
	params.Set("interval", _INERNAL_KLINE_PERIOD_CONVERTER[period])
	if since > 0 {
		params.Set("startTime", strconv.Itoa(since))
	}
	//params.Set("endTime", strconv.Itoa(int(time.Now().UnixNano()/1000000)))
	params.Set("limit", strconv.Itoa(size))

	klineUrl := bs.apiV1 + KLINE_URI + "?" + params.Encode()
	klines, err := HttpGet3(bs.httpClient, klineUrl, nil)
	if err != nil {
		return nil, err
	}
	var klineRecords []FutureKline

	for _, _record := range klines {
		r := Kline{Pair: currency}
		record := _record.([]interface{})
		r.Timestamp = int64(record[0].(float64)) / 1000 //to unix timestramp
		r.Open = ToFloat64(record[1])
		r.High = ToFloat64(record[2])
		r.Low = ToFloat64(record[3])
		r.Close = ToFloat64(record[4])
		r.Vol = ToFloat64(record[5])

		klineRecords = append(klineRecords, FutureKline{Kline: &r})
	}

	return klineRecords, nil
}

func (bs *BinanceSwap) GetServerTime() (int64, error) {
	respmap, err := HttpGet(bs.httpClient, bs.apiV1+SERVER_TIME_URL)
	if err != nil {
		return 0, err
	}

	stime := int64(ToInt(respmap["serverTime"]))

	return stime, nil
}

func (bs *BinanceSwap) adaptCurrencyPair(pair CurrencyPair) CurrencyPair {
	return pair.AdaptUsdToUsdt()
}
