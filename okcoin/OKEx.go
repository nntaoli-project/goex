package okcoin

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	FUTURE_API_BASE_URL    = "https://www.okex.com/api/v1/"
	FUTURE_TICKER_URI      = "future_ticker.do?symbol=%s&contract_type=%s"
	FUTURE_DEPTH_URI       = "future_depth.do?symbol=%s&contract_type=%s"
	FUTURE_INDEX_PRICE     = "future_index.do?symbol=%s"
	FUTURE_USERINFO_URI    = "future_userinfo.do"
	FUTURE_CANCEL_URI      = "future_cancel.do"
	FUTURE_ORDER_INFO_URI  = "future_order_info.do"
	FUTURE_ORDERS_INFO_URI = "future_orders_info.do"
	FUTURE_POSITION_URI    = "future_position.do"
	FUTURE_TRADE_URI       = "future_trade.do"
	TRADES_URI             = "future_trades.do"
	FUTURE_ESTIMATED_PRICE = "future_estimated_price.do?symbol=%s"
	_EXCHANGE_RATE_URI     = "exchange_rate.do"
	_GET_KLINE_URI         = "future_kline.do"
)

type OKEx struct {
	apiKey,
	apiSecretKey string
	client *http.Client
}

func NewOKEx(client *http.Client, api_key, secret_key string) *OKEx {
	ok := new(OKEx)
	ok.apiKey = api_key
	ok.apiSecretKey = secret_key
	ok.client = client
	return ok
}

func (ok *OKEx) buildPostForm(postForm *url.Values) error {
	postForm.Set("api_key", ok.apiKey)
	//postForm.Set("secret_key", ctx.secret_key);

	payload := postForm.Encode()
	payload = payload + "&secret_key=" + ok.apiSecretKey
	payload2, _ := url.QueryUnescape(payload) // can't escape for sign
	//
	sign, err := GetParamMD5Sign(ok.apiSecretKey, payload2)
	if err != nil {
		return err
	}

	postForm.Set("sign", strings.ToUpper(sign))
	//postForm.Del("secret_key")
	//fmt.Println(postForm)
	return nil
}

func (ok *OKEx) GetExchangeName() string {
	return OKEX_FUTURE
}

func (ok *OKEx) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	resp, err := ok.client.Get(fmt.Sprintf(FUTURE_API_BASE_URL+FUTURE_ESTIMATED_PRICE, strings.ToLower(currencyPair.ToSymbol("_"))))
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return 0, err
	}

	bodyMap := make(map[string]interface{})

	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		return 0, err
	}

	//println(string(body))
	return bodyMap["forecast_price"].(float64), nil
}

func (ok *OKEx) GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error) {
	url := FUTURE_API_BASE_URL + FUTURE_TICKER_URI
	//fmt.Println(fmt.Sprintf(url, strings.ToLower(currencyPair.ToSymbol("_")), contractType));
	resp, err := ok.client.Get(fmt.Sprintf(url, strings.ToLower(currencyPair.ToSymbol("_")), contractType))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	//println(string(body))

	bodyMap := make(map[string]interface{})

	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		return nil, err
	}

	if bodyMap["result"] != nil && !bodyMap["result"].(bool) {
		return nil, errors.New(string(body))
	}

	tickerMap := bodyMap["ticker"].(map[string]interface{})

	ticker := new(Ticker)
	ticker.Pair = currencyPair
	ticker.Date, _ = strconv.ParseUint(bodyMap["date"].(string), 10, 64)
	ticker.Buy = tickerMap["buy"].(float64)
	ticker.Sell = tickerMap["sell"].(float64)
	ticker.Last = tickerMap["last"].(float64)
	ticker.High = tickerMap["high"].(float64)
	ticker.Low = tickerMap["low"].(float64)
	ticker.Vol = tickerMap["vol"].(float64)

	//fmt.Println(bodyMap)
	return ticker, nil
}

func (ok *OKEx) GetFutureDepth(currencyPair CurrencyPair, contractType string, size int) (*Depth, error) {
	url := FUTURE_API_BASE_URL + FUTURE_DEPTH_URI
	//fmt.Println(fmt.Sprintf(url, strings.ToLower(currencyPair.ToSymbol("_")), contractType));
	resp, err := ok.client.Get(fmt.Sprintf(url, strings.ToLower(strings.ToLower(currencyPair.ToSymbol("_"))), contractType))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	//println(string(body))

	bodyMap := make(map[string]interface{})

	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		return nil, err
	}

	if bodyMap["error_code"] != nil {
		log.Println(bodyMap)
		return nil, errors.New(string(body))
	}

	depth := new(Depth)
	depth.Pair = currencyPair
	depth.ContractType = contractType

	size2 := len(bodyMap["asks"].([]interface{}))
	skipSize := 0
	if size < size2 {
		skipSize = size2 - size
	}

	for _, v := range bodyMap["asks"].([]interface{}) {
		if skipSize > 0 {
			skipSize--
			continue
		}

		var dr DepthRecord
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price = vv.(float64)
			case 1:
				dr.Amount = vv.(float64)
			}
		}
		depth.AskList = append(depth.AskList, dr)
	}

	for _, v := range bodyMap["bids"].([]interface{}) {
		var dr DepthRecord
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price = vv.(float64)
			case 1:
				dr.Amount = vv.(float64)
			}
		}
		depth.BidList = append(depth.BidList, dr)

		size--
		if size == 0 {
			break
		}
	}

	//fmt.Println(bodyMap)
	return depth, nil
}

func (ok *OKEx) GetFutureIndex(currencyPair CurrencyPair) (float64, error) {
	resp, err := ok.client.Get(fmt.Sprintf(FUTURE_API_BASE_URL+FUTURE_INDEX_PRICE, strings.ToLower(currencyPair.ToSymbol("_"))))
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return 0, err
	}

	bodyMap := make(map[string]interface{})

	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		return 0, err
	}

	//println(string(body))
	return bodyMap["future_index"].(float64), nil
}

type futureUserInfoResponse struct {
	Info struct {
		Btc map[string]float64 `json:btc`
		Ltc map[string]float64 `json:ltc`
		Etc map[string]float64 `json:"etc"`
		Eth map[string]float64 `json:"eth"`
		Bch map[string]float64 `json:"bch"`
		Xrp map[string]float64 `json:"xrp"`
		Eos map[string]float64 `json:"eos"`
		Btg map[string]float64 `json:"btg"`
		Bsv map[string]float64 `json:"bsv"`
	} `json:info`
	Result     bool `json:"result,bool"`
	Error_code int  `json:"error_code"`
}

func (ok *OKEx) GetFutureUserinfo() (*FutureAccount, error) {
	userInfoUrl := FUTURE_API_BASE_URL + FUTURE_USERINFO_URI

	postData := url.Values{}
	ok.buildPostForm(&postData)

	body, err := HttpPostForm(ok.client, userInfoUrl, postData)

	if err != nil {
		return nil, err
	}

	//println(string(body));
	resp := futureUserInfoResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if !resp.Result && resp.Error_code > 0 {
		return nil, ok.errorWrapper(resp.Error_code)
	}

	account := new(FutureAccount)
	account.FutureSubAccounts = make(map[Currency]FutureSubAccount, 2)

	btcMap := resp.Info.Btc
	ltcMap := resp.Info.Ltc
	bchMap := resp.Info.Bch
	ethMap := resp.Info.Eth
	etcMap := resp.Info.Etc
	xrpMap := resp.Info.Xrp
	eosMap := resp.Info.Eos
	btgMap := resp.Info.Btg
	bsvMap := resp.Info.Bsv

	account.FutureSubAccounts[BTC] = FutureSubAccount{BTC, btcMap["account_rights"], btcMap["keep_deposit"], btcMap["profit_real"], btcMap["profit_unreal"], btcMap["risk_rate"]}
	account.FutureSubAccounts[LTC] = FutureSubAccount{LTC, ltcMap["account_rights"], ltcMap["keep_deposit"], ltcMap["profit_real"], ltcMap["profit_unreal"], ltcMap["risk_rate"]}
	account.FutureSubAccounts[BCH] = FutureSubAccount{BCH, bchMap["account_rights"], bchMap["keep_deposit"], bchMap["profit_real"], bchMap["profit_unreal"], bchMap["risk_rate"]}
	account.FutureSubAccounts[ETH] = FutureSubAccount{ETH, ethMap["account_rights"], ethMap["keep_deposit"], ethMap["profit_real"], ethMap["profit_unreal"], ethMap["risk_rate"]}
	account.FutureSubAccounts[ETC] = FutureSubAccount{ETC, etcMap["account_rights"], etcMap["keep_deposit"], etcMap["profit_real"], etcMap["profit_unreal"], etcMap["risk_rate"]}
	account.FutureSubAccounts[XRP] = FutureSubAccount{XRP, xrpMap["account_rights"], xrpMap["keep_deposit"], xrpMap["profit_real"], xrpMap["profit_unreal"], xrpMap["risk_rate"]}
	account.FutureSubAccounts[EOS] = FutureSubAccount{EOS, eosMap["account_rights"], eosMap["keep_deposit"], eosMap["profit_real"], eosMap["profit_unreal"], eosMap["risk_rate"]}
	account.FutureSubAccounts[BTG] = FutureSubAccount{BTG, btgMap["account_rights"], btgMap["keep_deposit"], btgMap["profit_real"], btgMap["profit_unreal"], btgMap["risk_rate"]}
	account.FutureSubAccounts[BSV] = FutureSubAccount{BSV, bsvMap["account_rights"], bsvMap["keep_deposit"], bsvMap["profit_real"], bsvMap["profit_unreal"], bsvMap["risk_rate"]}

	return account, nil
}

func (ok *OKEx) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice, leverRate int) (string, error) {
	postData := url.Values{}
	postData.Set("symbol", strings.ToLower(currencyPair.ToSymbol("_")))
	postData.Set("price", price)
	postData.Set("contract_type", contractType)
	postData.Set("amount", amount)
	postData.Set("type", strconv.Itoa(openType))
	postData.Set("lever_rate", strconv.Itoa(leverRate))
	postData.Set("match_price", strconv.Itoa(matchPrice))

	ok.buildPostForm(&postData)

	placeOrderUrl := FUTURE_API_BASE_URL + FUTURE_TRADE_URI
	body, err := HttpPostForm(ok.client, placeOrderUrl, postData)

	if err != nil {
		return "", err
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(body, &respMap)
	if err != nil {
		return "", err
	}

	//println(string(body));

	if !respMap["result"].(bool) {
		return "", errors.New(string(body))
	}

	return fmt.Sprintf("%.0f", respMap["order_id"].(float64)), nil
}

func (ok *OKEx) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error) {
	postData := url.Values{}
	postData.Set("symbol", strings.ToLower(currencyPair.ToSymbol("_")))
	postData.Set("order_id", orderId)
	postData.Set("contract_type", contractType)

	ok.buildPostForm(&postData)

	cancelUrl := FUTURE_API_BASE_URL + FUTURE_CANCEL_URI

	body, err := HttpPostForm(ok.client, cancelUrl, postData)
	if err != nil {
		return false, err
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(body, &respMap)
	if err != nil {
		return false, err
	}

	if respMap["result"] != nil && !respMap["result"].(bool) {
		return false, errors.New(string(body))
	}

	return true, nil
}

func (ok *OKEx) GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	positionUrl := FUTURE_API_BASE_URL + FUTURE_POSITION_URI

	postData := url.Values{}
	postData.Set("contract_type", contractType)
	postData.Set("symbol", strings.ToLower(currencyPair.ToSymbol("_")))

	ok.buildPostForm(&postData)

	body, err := HttpPostForm(ok.client, positionUrl, postData)

	if err != nil {
		return nil, err
	}

	respMap := make(map[string]interface{})

	err = json.Unmarshal(body, &respMap)
	if err != nil {
		return nil, err
	}

	if !respMap["result"].(bool) {
		return nil, errors.New(string(body))
	}

	//println(string(body))

	var posAr []FuturePosition

	forceLiquPriceStr := respMap["force_liqu_price"].(string)
	forceLiquPriceStr = strings.Replace(forceLiquPriceStr, ",", "", 1)
	forceLiquPrice, err := strconv.ParseFloat(forceLiquPriceStr, 64)

	holdings := respMap["holding"].([]interface{})
	for _, v := range holdings {
		holdingMap := v.(map[string]interface{})

		pos := FuturePosition{}
		pos.ForceLiquPrice = forceLiquPrice
		pos.LeverRate = int(holdingMap["lever_rate"].(float64))
		pos.ContractType = holdingMap["contract_type"].(string)
		pos.ContractId = int64(holdingMap["contract_id"].(float64))
		pos.BuyAmount = holdingMap["buy_amount"].(float64)
		pos.BuyAvailable = holdingMap["buy_available"].(float64)
		pos.BuyPriceAvg = holdingMap["buy_price_avg"].(float64)
		pos.BuyPriceCost = holdingMap["buy_price_cost"].(float64)
		pos.BuyProfitReal = holdingMap["buy_profit_real"].(float64)
		pos.SellAmount = holdingMap["sell_amount"].(float64)
		pos.SellAvailable = holdingMap["sell_available"].(float64)
		pos.SellPriceAvg = holdingMap["sell_price_avg"].(float64)
		pos.SellPriceCost = holdingMap["sell_price_cost"].(float64)
		pos.SellProfitReal = holdingMap["sell_profit_real"].(float64)
		pos.CreateDate = int64(holdingMap["create_date"].(float64))
		pos.Symbol = currencyPair
		posAr = append(posAr, pos)

	}

	return posAr, nil
}

func (ok *OKEx) parseOrders(body []byte, currencyPair CurrencyPair) ([]FutureOrder, error) {
	respMap := make(map[string]interface{})

	err := json.Unmarshal(body, &respMap)
	if err != nil {
		return nil, err
	}

	if !respMap["result"].(bool) {
		return nil, errors.New(string(body))
	}

	var orders []interface{}
	orders = respMap["orders"].([]interface{})

	var futureOrders []FutureOrder

	for _, v := range orders {
		vv := v.(map[string]interface{})
		futureOrder := FutureOrder{}
		futureOrder.OrderID = int64(vv["order_id"].(float64))
		futureOrder.OrderID2 = fmt.Sprint(int64(vv["order_id"].(float64)))
		futureOrder.Amount = vv["amount"].(float64)
		futureOrder.Price = vv["price"].(float64)
		futureOrder.AvgPrice = vv["price_avg"].(float64)
		futureOrder.DealAmount = vv["deal_amount"].(float64)
		futureOrder.Fee = vv["fee"].(float64)
		futureOrder.OType = int(vv["type"].(float64))
		futureOrder.OrderTime = int64(vv["create_date"].(float64))
		futureOrder.LeverRate = int(vv["lever_rate"].(float64))
		futureOrder.ContractName = vv["contract_name"].(string)
		futureOrder.Currency = currencyPair

		switch s := int(vv["status"].(float64)); s {
		case 0:
			futureOrder.Status = ORDER_UNFINISH
		case 1:
			futureOrder.Status = ORDER_PART_FINISH
		case 2:
			futureOrder.Status = ORDER_FINISH
		case 4:
			futureOrder.Status = ORDER_CANCEL_ING
		case -1:
			futureOrder.Status = ORDER_CANCEL

		}

		futureOrders = append(futureOrders, futureOrder)
	}
	return futureOrders, nil
}

func (ok *OKEx) GetFutureOrders(orderIds []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	postData := url.Values{}
	postData.Set("order_id", strings.Join(orderIds, ","))
	postData.Set("contract_type", contractType)
	postData.Set("symbol", strings.ToLower(currencyPair.ToSymbol("_")))
	ok.buildPostForm(&postData)

	body, err := HttpPostForm(ok.client, FUTURE_API_BASE_URL+FUTURE_ORDERS_INFO_URI, postData)
	if err != nil {
		return nil, err
	}

	return ok.parseOrders(body, currencyPair)
}

func (ok *OKEx) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	postData := url.Values{}
	postData.Set("order_id", "-1")
	postData.Set("contract_type", contractType)
	postData.Set("symbol", strings.ToLower(currencyPair.ToSymbol("_")))
	postData.Set("status", "1")
	postData.Set("current_page", "1")
	postData.Set("page_length", "50")

	ok.buildPostForm(&postData)

	body, err := HttpPostForm(ok.client, FUTURE_API_BASE_URL+FUTURE_ORDER_INFO_URI, postData)
	if err != nil {
		return nil, err
	}

	//println(string(body))

	return ok.parseOrders(body, currencyPair)
}

func (ok *OKEx) GetFutureOrder(orderId string, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	postData := url.Values{}
	postData.Set("order_id", orderId)
	postData.Set("contract_type", contractType)
	postData.Set("symbol", strings.ToLower(currencyPair.ToSymbol("_")))
	//postData.Set("status", "1")
	postData.Set("current_page", "1")
	postData.Set("page_length", "2")

	ok.buildPostForm(&postData)

	body, err := HttpPostForm(ok.client, FUTURE_API_BASE_URL+FUTURE_ORDER_INFO_URI, postData)
	if err != nil {
		return nil, err
	}

	//println(string(body))
	orders, err := ok.parseOrders(body, currencyPair)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, errors.New("找不到订单")
	}

	return &orders[0], nil
}

func (ok *OKEx) GetFee() (float64, error) {
	return 0.03, nil //期货固定0.03%手续费
}

func (ok *OKEx) GetExchangeRate() (float64, error) {
	respMap, err := HttpGet(ok.client, FUTURE_API_BASE_URL+_EXCHANGE_RATE_URI)

	if err != nil {
		log.Println(respMap)
		return -1, err
	}

	if respMap["rate"] == nil {
		log.Println(respMap)
		return -1, errors.New("error")
	}

	return respMap["rate"].(float64), nil
}

/**
 * BTC: 100美元一张合约
 * LTC/ETH/ETC/BCH: 10美元一张合约
 */
func (ok *OKEx) GetContractValue(currencyPair CurrencyPair) (float64, error) {
	switch currencyPair {
	case BTC_USD:
		return 100, nil
	case LTC_USD, ETH_USD, ETC_USD, BCH_USD:
		return 10, nil
	default:
		return 10, nil
	}

	return -1, errors.New("error")
}

func (ok *OKEx) GetDeliveryTime() (int, int, int, int) {
	return 4, 16, 0, 0 //星期五，下午4点交割
}

func (ok *OKEx) GetKlineRecords(contract_type string, currencyPair CurrencyPair, period, size, since int) ([]FutureKline, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToLower(currencyPair.ToSymbol("_")))
	params.Set("type", _INERNAL_KLINE_PERIOD_CONVERTER[period])
	params.Set("contract_type", contract_type)
	params.Set("size", fmt.Sprintf("%d", size))
	params.Set("since", fmt.Sprintf("%d", since))
	//log.Println(params.Encode())
	resp, err := ok.client.Get(FUTURE_API_BASE_URL + _GET_KLINE_URI + "?" + params.Encode())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//log.Println(string(body))

	var klines [][]interface{}
	err = json.Unmarshal(body, &klines)
	if err != nil {
		log.Println(string(body))
		return nil, err
	}

	var klineRecords []FutureKline
	for _, record := range klines {
		r := FutureKline{}
		r.Kline = new(Kline)
		for i, e := range record {
			switch i {
			case 0:
				r.Timestamp = int64(e.(float64)) / 1000 //to unix timestramp
			case 1:
				r.Open = e.(float64)
			case 2:
				r.High = e.(float64)
			case 3:
				r.Low = e.(float64)
			case 4:
				r.Close = e.(float64)
			case 5:
				r.Vol = e.(float64)
			case 6:
				r.Vol2 = e.(float64)
			}
		}
		klineRecords = append(klineRecords, r)
	}

	return klineRecords, nil
}

func (okFuture *OKEx) GetTrades(contract_type string, currencyPair CurrencyPair, since int64) ([]Trade, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToLower(currencyPair.ToSymbol("_")))
	params.Set("contract_type", contract_type)
	//log.Println(params.Encode())

	url := FUTURE_API_BASE_URL + TRADES_URI + "?" + params.Encode()
	body, err := HttpGet5(okFuture.client, url, nil)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))

	var trades []Trade
	var resp []interface{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, errors.New(string(body))
	}

	for _, v := range resp {
		item := v.(map[string]interface{})

		tid := int64(item["tid"].(float64))
		direction := item["type"].(string)
		amount := item["amount"].(float64)
		price := item["price"].(float64)
		time := int64(item["date_ms"].(float64))

		trades = append(trades, Trade{tid, AdaptTradeSide(direction), amount, price, time, currencyPair})
	}

	return trades, nil
}

func (okFuture *OKEx) errorWrapper(errorCode int) ApiError {
	switch errorCode {
	case 20024:
		return EX_ERR_SIGN
	case 20020:
		return EX_ERR_NOT_FIND_SECRETKEY
	case 20015:
		return EX_ERR_NOT_FIND_ORDER
	case 20049:
		return EX_ERR_API_LIMIT
	}
	errmsg := fmt.Sprintf("%d", errorCode)
	return ApiError{ErrCode: errmsg, OriginErrMsg: errmsg}
}
