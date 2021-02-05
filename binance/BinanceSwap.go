package binance

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	baseUrl = "https://fapi.binance.com"
)

type BinanceSwap struct {
	Binance
	f *BinanceFutures
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
		f: NewBinanceFutures(&APIConfig{
			Endpoint:     strings.ReplaceAll(config.Endpoint, "fapi", "dapi"),
			HttpClient:   config.HttpClient,
			ApiKey:       config.ApiKey,
			ApiSecretKey: config.ApiSecretKey,
			Lever:        config.Lever,
		}),
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

func (bn *BinanceSwap) GetExchangeInfo() (*ExchangeInfo, error) {
	resp, err := HttpGet5(bn.httpClient, bn.apiV1+"exchangeInfo", nil)
	if err != nil {
		return nil, err
	}
	info := &ExchangeInfo{}
	err = json.Unmarshal(resp, info)
	if err != nil {
		return nil, err
	}

	return info, nil
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

func (bs *BinanceSwap) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	panic("not supported.")
}

func (bs *BinanceSwap) GetFutureTicker(currency CurrencyPair, contractType string) (*Ticker, error) {
	if contractType == SWAP_CONTRACT {
		return bs.f.GetFutureTicker(currency.AdaptUsdtToUsd(), SWAP_CONTRACT)
	}

	if contractType != SWAP_USDT_CONTRACT {
		return nil, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}

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

func (bs *BinanceSwap) GetFutureDepth(currency CurrencyPair, contractType string, size int) (*Depth, error) {
	if contractType == SWAP_CONTRACT {
		return bs.f.GetFutureDepth(currency.AdaptUsdtToUsd(), SWAP_CONTRACT, size)
	}

	if contractType != SWAP_USDT_CONTRACT {
		return nil, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}

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

func (bs *BinanceSwap) GetFutureOrderHistory(pair CurrencyPair, contractType string, optional ...OptionalParameter) ([]FutureOrder, error) {
	panic("implement me")
}

func (bs *BinanceSwap) GetTrades(contractType string, currencyPair CurrencyPair, since int64) ([]Trade, error) {
	if contractType == SWAP_CONTRACT {
		return bs.f.GetTrades(SWAP_CONTRACT, currencyPair.AdaptUsdtToUsd(), since)
	}

	if contractType != SWAP_USDT_CONTRACT {
		return nil, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}

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

func (bs *BinanceSwap) GetFutureIndex(currencyPair CurrencyPair) (float64, error) {
	respmap, err := HttpGet(bs.httpClient, bs.apiV1+"premiumIndex?symbol="+bs.adaptCurrencyPair(currencyPair).ToSymbol(""))
	if err != nil {
		return 0.0, err
	}

	return ToFloat64(respmap["markPrice"]), nil
}

func (bs *BinanceSwap) GetFutureUserinfo(currencyPair ...CurrencyPair) (*FutureAccount, error) {
	acc, err := bs.f.GetFutureUserinfo(currencyPair...)
	if err != nil {
		return nil, err
	}

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

	balances := respmap["assets"].([]interface{})
	for _, v := range balances {
		vv := v.(map[string]interface{})
		currency := NewCurrency(vv["asset"].(string), "").AdaptBccToBch()
		acc.FutureSubAccounts[currency] = FutureSubAccount{
			Currency:      currency,
			AccountRights: ToFloat64(vv["marginBalance"]),
			KeepDeposit:   ToFloat64(vv["maintMargin"]),
			ProfitUnreal:  ToFloat64(vv["unrealizedProfit"]),
		}
	}

	return acc, nil
}

//@deprecated please call the Wallet api
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

func (bs *BinanceSwap) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice int, leverRate float64) (string, error) {
	fOrder, err := bs.PlaceFutureOrder2(currencyPair, contractType, price, amount, openType, matchPrice, leverRate)
	return fOrder.OrderID2, err
}

func (bs *BinanceSwap) PlaceFutureOrder2(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice int, leverRate float64) (*FutureOrder, error) {
	if contractType == SWAP_CONTRACT {
		orderId, err := bs.f.PlaceFutureOrder(currencyPair.AdaptUsdtToUsd(), contractType, price, amount, openType, matchPrice, leverRate)
		return &FutureOrder{
			OrderID2:     orderId,
			Price:        ToFloat64(price),
			Amount:       ToFloat64(amount),
			Status:       ORDER_UNFINISH,
			Currency:     currencyPair,
			OType:        openType,
			LeverRate:    leverRate,
			ContractName: contractType,
		}, err
	}

	if contractType != SWAP_USDT_CONTRACT {
		return nil, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}

	fOrder := &FutureOrder{
		Currency:     currencyPair,
		ClientOid:    GenerateOrderClientId(32),
		Price:        ToFloat64(price),
		Amount:       ToFloat64(amount),
		OrderType:    openType,
		LeverRate:    leverRate,
		ContractName: contractType,
	}

	pair := bs.adaptCurrencyPair(currencyPair)
	path := bs.apiV1 + ORDER_URI
	params := url.Values{}
	params.Set("symbol", pair.ToSymbol(""))
	params.Set("quantity", amount)
	params.Set("newClientOrderId", fOrder.ClientOid)

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
		return fOrder, err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return fOrder, err
	}

	orderId := ToInt(respmap["orderId"])
	if orderId <= 0 {
		return fOrder, errors.New(string(resp))
	}
	fOrder.OrderID2 = strconv.Itoa(orderId)

	return fOrder, nil
}

func (bs *BinanceSwap) LimitFuturesOrder(currencyPair CurrencyPair, contractType, price, amount string, openType int, opt ...LimitOrderOptionalParameter) (*FutureOrder, error) {
	return bs.PlaceFutureOrder2(currencyPair, contractType, price, amount, openType, 0, 10)
}

func (bs *BinanceSwap) MarketFuturesOrder(currencyPair CurrencyPair, contractType, amount string, openType int) (*FutureOrder, error) {
	return bs.PlaceFutureOrder2(currencyPair, contractType, "0", amount, openType, 1, 10)
}

func (bs *BinanceSwap) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error) {
	if contractType == SWAP_CONTRACT {
		return bs.f.FutureCancelOrder(currencyPair.AdaptUsdtToUsd(), contractType, orderId)
	}

	if contractType != SWAP_USDT_CONTRACT {
		return false, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}

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
	if contractType == SWAP_CONTRACT {
		return false, errors.New("not support")
	}

	if contractType == SWAP_CONTRACT {
		return false, errors.New("not support")
	}

	if contractType != SWAP_USDT_CONTRACT {
		return false, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}

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
	if contractType != SWAP_USDT_CONTRACT {
		return false, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}

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

func (bs *BinanceSwap) GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	if contractType == SWAP_CONTRACT {
		return bs.f.GetFuturePosition(currencyPair.AdaptUsdtToUsd(), contractType)
	}

	if contractType != SWAP_USDT_CONTRACT {
		return nil, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}

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
			LeverRate:      ToFloat64(cont["leverage"]),
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

func (bs *BinanceSwap) GetFutureOrders(orderIds []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	if contractType == SWAP_CONTRACT {
		return nil, errors.New("not support")
	}

	if contractType != SWAP_USDT_CONTRACT {
		return nil, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}

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

func (bs *BinanceSwap) GetFutureOrder(orderId string, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	if contractType == SWAP_CONTRACT {
		return bs.f.GetFutureOrder(orderId, currencyPair.AdaptUsdtToUsd(), contractType)
	}

	if contractType != SWAP_USDT_CONTRACT {
		return nil, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}

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

func (bs *BinanceSwap) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	if contractType == SWAP_CONTRACT {
		return bs.f.GetUnfinishFutureOrders(currencyPair.AdaptUsdtToUsd(), contractType)
	}

	if contractType != SWAP_USDT_CONTRACT {
		return nil, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}

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

func (bs *BinanceSwap) GetFee() (float64, error) {
	panic("not supported.")
}

func (bs *BinanceSwap) GetContractValue(currencyPair CurrencyPair) (float64, error) {
	panic("not supported.")
}

func (bs *BinanceSwap) GetDeliveryTime() (int, int, int, int) {
	panic("not supported.")
}

func (bs *BinanceSwap) GetKlineRecords(contractType string, currency CurrencyPair, period KlinePeriod, size int, opt ...OptionalParameter) ([]FutureKline, error) {
	if contractType == SWAP_CONTRACT {
		return bs.f.GetKlineRecords(contractType, currency.AdaptUsdtToUsd(), period, size, opt...)
	}

	if contractType != SWAP_USDT_CONTRACT {
		return nil, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}

	currency2 := bs.adaptCurrencyPair(currency)
	params := url.Values{}
	params.Set("symbol", currency2.ToSymbol(""))
	params.Set("interval", _INERNAL_KLINE_PERIOD_CONVERTER[KlinePeriod(period)])
	//params.Set("endTime", strconv.Itoa(int(time.Now().UnixNano()/1000000)))
	params.Set("limit", strconv.Itoa(size))
	MergeOptionalParameter(&params, opt...)

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
