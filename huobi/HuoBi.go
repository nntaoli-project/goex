package huobi

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	EXCHANGE_NAME = "huobi.com"
	API_BASE_URL  = "https://api.huobi.com/"
	TRADE_API_V3  = API_BASE_URL + "apiv3"
	TICKER_URI    = "staticmarket/ticker_%s_json.js"
	DEPTH_URI     = "staticmarket/depth_%s_%d.js"
	KLINE_URI     = "staticmarket/%s_kline_%03s_json.js?length=%d"
	trade_url     = "staticmarket/detail_%s_json.js"
)

type HuoBi struct {
	httpClient *http.Client
	accessKey,
	secretKey string
}

var _INERNAL_KLINE_PERIOD_CONVERTER = map[int]string{
	KLINE_PERIOD_1MIN:  "001",
	KLINE_PERIOD_5MIN:  "005",
	KLINE_PERIOD_15MIN: "015",
	KLINE_PERIOD_30MIN: "030",
	KLINE_PERIOD_60MIN: "060",
	KLINE_PERIOD_1DAY:  "100",
	KLINE_PERIOD_1WEEK: "200",
}

func New(httpClient *http.Client, accessKey, secretKey string) *HuoBi {
	return &HuoBi{httpClient, accessKey, secretKey}
}

func (hb *HuoBi) buildPostForm(postForm *url.Values) error {
	postForm.Set("created", fmt.Sprintf("%d", time.Now().Unix()))
	postForm.Set("access_key", hb.accessKey)
	postForm.Set("secret_key", hb.secretKey)
	sign, err := GetParamMD5Sign(hb.secretKey, postForm.Encode())
	if err != nil {
		return err
	}
	postForm.Set("sign", sign)
	postForm.Del("secret_key")
	return nil
}

func (hb *HuoBi) GetExchangeName() string {
	return EXCHANGE_NAME
}

func (hb *HuoBi) GetTicker(currency CurrencyPair) (*Ticker, error) {
	var tickerUri string

	switch currency {
	case BTC_CNY:
		tickerUri = fmt.Sprintf(API_BASE_URL+TICKER_URI, "btc")
	case LTC_CNY:
		tickerUri = fmt.Sprintf(API_BASE_URL+TICKER_URI, "ltc")
	default:
		return nil, errors.New("Unsupport The CurrencyPair")
	}

	bodyDataMap, err := HttpGet(hb.httpClient, tickerUri)

	if err != nil {
		return nil, err
	}

	var tickerMap map[string]interface{}
	var ticker Ticker

	switch bodyDataMap["ticker"].(type) {
	case map[string]interface{}:
		tickerMap = bodyDataMap["ticker"].(map[string]interface{})
	default:
		return nil, errors.New(fmt.Sprintf("Type Convert Error ? \n %s", bodyDataMap))
	}

	ticker.Date, _ = strconv.ParseUint(bodyDataMap["time"].(string), 10, 64)
	ticker.Last = tickerMap["last"].(float64)
	ticker.Buy = tickerMap["buy"].(float64)
	ticker.Sell = tickerMap["sell"].(float64)
	ticker.Low = tickerMap["low"].(float64)
	ticker.High = tickerMap["high"].(float64)
	ticker.Vol = tickerMap["vol"].(float64)

	return &ticker, nil
}

func (hb *HuoBi) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	var depthUri string

	switch currency {
	case BTC_CNY:
		depthUri = fmt.Sprintf(API_BASE_URL+DEPTH_URI, "btc", size)
	case LTC_CNY:
		depthUri = fmt.Sprintf(API_BASE_URL+DEPTH_URI, "ltc", size)
	default:
		return nil, errors.New("Unsupport The CurrencyPair")
	}

	bodyDataMap, err := HttpGet(hb.httpClient, depthUri)

	if err != nil {
		return nil, err
	}

	if bodyDataMap["code"] != nil {
		log.Println(bodyDataMap)
		return nil, errors.New(fmt.Sprintf("%s", bodyDataMap))
	}

	var depth Depth

	asks, isOK := bodyDataMap["asks"].([]interface{})
	if !isOK {
		return nil, errors.New("asks assert error")
	}

	i := len(asks) - 1

	for ; i >= 0; i-- {
		ask := asks[i]
		var dr DepthRecord
		for i, vv := range ask.([]interface{}) {
			switch i {
			case 0:
				dr.Price = vv.(float64)
			case 1:
				dr.Amount = vv.(float64)
			}
		}
		depth.AskList = append(depth.AskList, dr)
	}

	for _, v := range bodyDataMap["bids"].([]interface{}) {
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
	}

	return &depth, nil
}

func (hb *HuoBi) GetAccount() (*Account, error) {
	postData := url.Values{}
	postData.Set("method", "get_account_info")
	postData.Set("created", fmt.Sprintf("%d", time.Now().Unix()))
	postData.Set("access_key", hb.accessKey)
	postData.Set("secret_key", hb.secretKey)

	sign, _ := GetParamMD5Sign(hb.secretKey, postData.Encode())
	postData.Set("sign", sign)
	postData.Del("secret_key")

	bodyData, err := HttpPostForm(hb.httpClient, TRADE_API_V3, postData)
	if err != nil {
		return nil, err
	}

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(bodyData, &bodyDataMap)
	if err != nil {
		println(string(bodyData))
		return nil, err
	}

	if bodyDataMap["code"] != nil {
		return nil, errors.New(fmt.Sprintf("%s", bodyDataMap))
	}

	account := new(Account)
	account.Exchange = hb.GetExchangeName()
	account.Asset, _ = strconv.ParseFloat(bodyDataMap["total"].(string), 64)
	account.NetAsset, _ = strconv.ParseFloat(bodyDataMap["net_asset"].(string), 64)

	var btcSubAccount SubAccount
	var ltcSubAccount SubAccount
	var cnySubAccount SubAccount

	btcSubAccount.Currency = BTC
	btcSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["available_btc_display"].(string), 64)
	btcSubAccount.LoanAmount, _ = strconv.ParseFloat(bodyDataMap["loan_btc_display"].(string), 64)
	btcSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["frozen_btc_display"].(string), 64)

	ltcSubAccount.Currency = LTC
	ltcSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["available_ltc_display"].(string), 64)
	ltcSubAccount.LoanAmount, _ = strconv.ParseFloat(bodyDataMap["loan_ltc_display"].(string), 64)
	ltcSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["frozen_ltc_display"].(string), 64)

	cnySubAccount.Currency = CNY
	cnySubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["available_cny_display"].(string), 64)
	cnySubAccount.LoanAmount, _ = strconv.ParseFloat(bodyDataMap["loan_cny_display"].(string), 64)
	cnySubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["frozen_cny_display"].(string), 64)

	account.SubAccounts = make(map[Currency]SubAccount, 3)
	account.SubAccounts[BTC] = btcSubAccount
	account.SubAccounts[LTC] = ltcSubAccount
	account.SubAccounts[CNY] = cnySubAccount

	return account, nil
}

func (hb *HuoBi) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	postData := url.Values{}
	postData.Set("method", "order_info")
	postData.Set("id", orderId)

	switch currency {
	case BTC_CNY:
		postData.Set("coin_type", "1")
	case LTC_CNY:
		postData.Set("coin_type", "2")
	}

	hb.buildPostForm(&postData)

	bodyData, err := HttpPostForm(hb.httpClient, TRADE_API_V3, postData)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(bodyData, &bodyDataMap)
	if err != nil {
		println(string(bodyData))
		return nil, err
	}

	if bodyDataMap["code"] != nil {
		return nil, errors.New(string(bodyData))
	}

	//fmt.Println(bodyDataMap);
	order := new(Order)
	order.Currency = currency
	order.OrderID, _ = strconv.Atoi(orderId)
	order.Side = TradeSide(bodyDataMap["type"].(float64))
	order.Amount, _ = strconv.ParseFloat(bodyDataMap["order_amount"].(string), 64)
	order.DealAmount, _ = strconv.ParseFloat(bodyDataMap["processed_amount"].(string), 64)
	order.Price, _ = strconv.ParseFloat(bodyDataMap["order_price"].(string), 64)
	order.AvgPrice, _ = strconv.ParseFloat(bodyDataMap["processed_price"].(string), 64)
	order.Fee, _ = strconv.ParseFloat(bodyDataMap["fee"].(string), 64)

	tradeStatus := TradeStatus(bodyDataMap["status"].(float64))
	switch tradeStatus {
	case 0:
		order.Status = ORDER_UNFINISH
	case 1:
		order.Status = ORDER_PART_FINISH
	case 2:
		order.Status = ORDER_FINISH
	case 3:
		order.Status = ORDER_CANCEL
	}
	//fmt.Println(order)
	return order, nil
}

func (hb *HuoBi) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	postData := url.Values{}
	postData.Set("method", "get_orders")

	switch currency {
	case BTC_CNY:
		postData.Set("coin_type", "1")
	case LTC_CNY:
		postData.Set("coin_type", "2")
	}

	hb.buildPostForm(&postData)

	bodyData, err := HttpPostForm(hb.httpClient, TRADE_API_V3, postData)
	if err != nil {
		return nil, err
	}

	if strings.Contains(string(bodyData), "code") {
		return nil, errors.New(string(bodyData))
	}
	//println(string(bodyData))

	var bodyDataMap []map[string]interface{}
	err = json.Unmarshal(bodyData, &bodyDataMap)
	if err != nil {
		return nil, err
	}

	var orders []Order

	for _, v := range bodyDataMap {
		order := Order{}
		order.Currency = currency
		order.Amount, _ = strconv.ParseFloat(v["order_amount"].(string), 64)
		order.Price, _ = strconv.ParseFloat(v["order_price"].(string), 64)
		order.DealAmount, _ = strconv.ParseFloat(v["processed_amount"].(string), 64)
		order.OrderTime = int(v["order_time"].(float64))
		order.OrderID = int(v["id"].(float64))
		order.Side = TradeSide(v["type"].(float64))
		orders = append(orders, order)
	}

	return orders, nil
}

func (hb *HuoBi) placeOrder(method, amount, price string, currency CurrencyPair) (*Order, error) {
	postData := url.Values{}
	postData.Set("method", method)

	switch method {
	case "buy", "sell":
		postData.Set("amount", amount)
		postData.Set("price", price)
	case "buy_market":
		postData.Set("amount", price)
	case "sell_market":
		postData.Set("amount", amount)
	}

	switch currency {
	case BTC_CNY:
		postData.Set("coin_type", "1")
	case LTC_CNY:
		postData.Set("coin_type", "2")
	}

	hb.buildPostForm(&postData)

	bodyData, err := HttpPostForm(hb.httpClient, TRADE_API_V3, postData)
	if err != nil {
		return nil, err
	}

	//{"result":"success","id":1321475746}
	//println(string(bodyData))

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(bodyData, &bodyDataMap)

	if err != nil {
		return nil, err
	}

	if bodyDataMap["code"] != nil {
		return nil, errors.New(string(bodyData))
	}

	ret := bodyDataMap["result"].(string)

	if strings.Compare(ret, "success") == 0 {
		order := new(Order)
		order.OrderID = int(bodyDataMap["id"].(float64))
		order.Price, _ = strconv.ParseFloat(price, 64)
		order.Amount, _ = strconv.ParseFloat(amount, 64)
		order.Currency = currency
		order.Status = ORDER_UNFINISH
		return order, nil
	}

	return nil, errors.New(fmt.Sprintf("Place Limit %s Order Fail.", method))
}

func (hb *HuoBi) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	order, err := hb.placeOrder("buy", amount, price, currency)

	if err != nil {
		return nil, err
	}

	order.Side = BUY

	return order, nil
}

func (hb *HuoBi) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	order, err := hb.placeOrder("sell", amount, price, currency)

	if err != nil {
		return nil, err
	}

	order.Side = SELL

	return order, nil
}

func (hb *HuoBi) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	order, err := hb.placeOrder("buy_market", amount, price, currency)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	order.Side = BUY
	return order, nil
}

func (hb *HuoBi) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	order, err := hb.placeOrder("sell_market", amount, price, currency)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	order.Side = SELL
	return order, nil
}

func (hb *HuoBi) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	//1321490762
	postData := url.Values{}
	postData.Set("method", "cancel_order")
	postData.Set("id", orderId)

	switch currency {
	case BTC_CNY:
		postData.Set("coin_type", "1")
	case LTC_CNY:
		postData.Set("coin_type", "2")
	}

	hb.buildPostForm(&postData)

	bodyData, err := HttpPostForm(hb.httpClient, TRADE_API_V3, postData)
	if err != nil {
		return false, err
	}

	//{"result":"success"}
	//{"code":42,"msg":"该委托已经取消, 不能取消或修改","message":"该委托已经取消, 不能取消或修改"}
	//println(string(bodyData))

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(bodyData, &bodyDataMap)
	if err != nil {
		return false, err
	}

	if bodyDataMap["code"] != nil {
		return false, errors.New(string(bodyData))
	}

	ret := bodyDataMap["result"].(string)
	return (strings.Compare(ret, "success") == 0), nil
}

/**
 * 具体参数详解: https://github.com/huobiapi/API_Docs/wiki/REST-Interval
 */
func (hb *HuoBi) GetKlineRecords(currency CurrencyPair, period , size, since int) ([]Kline, error) {
	klineUri := API_BASE_URL + KLINE_URI

	_period := _INERNAL_KLINE_PERIOD_CONVERTER[period]
	if _period == "" {
		return nil, errors.New("unsupport the kline period")
	}

	switch currency {
	case BTC_CNY:
		klineUri = fmt.Sprintf(klineUri, "btc", _period, size)
	case LTC_CNY:
		klineUri = fmt.Sprintf(klineUri, "ltc", _period, size)
	default:
		return nil, errors.New("Unsupport " + currency.String())
	}
	//println(klineUri)
	resp, err := http.Get(klineUri)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var klines [][]interface{}

	err = json.Unmarshal(body, &klines)

	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Local") //获取本地时区
	var klineRecords []Kline

	for _, record := range klines {
		r := Kline{}
		for i, e := range record {
			switch i {
			case 0:
				d := e.(string)
				if len(d) >= 12 {
					t, _ := time.ParseInLocation("200601021504", d[0:12], loc)
					r.Timestamp = t.Unix()
				}
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
			}
		}

		if r.Timestamp < int64(since/1000) {
			continue
		}

		klineRecords = append(klineRecords, r)
	}

	return klineRecords, nil
}

func (hb *HuoBi) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	return nil, nil
}

/**
 * 获取全站最近的交易记录
 */
func (hb *HuoBi) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	tradeUrl := API_BASE_URL + trade_url
	switch currencyPair {
	case BTC_CNY:
		tradeUrl = fmt.Sprintf(tradeUrl, "btc")
	case LTC_CNY:
		tradeUrl = fmt.Sprintf(tradeUrl, "ltc")
	default:
		return nil, errors.New("unsupport " + currencyPair.String())
	}

	var respmap map[string]interface{}

	resp, err := http.Get(tradeUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &respmap)
	if err != nil {
		return nil, err
	}

	tradesmap, isOK := respmap["trades"].([]interface{})
	if !isOK {
		return nil, errors.New("assert error")
	}

	now := time.Now()
	var trades []Trade
	for _, t := range tradesmap {
		tr := t.(map[string]interface{})
		trade := Trade{}
		trade.Amount = tr["amount"].(float64)
		trade.Price = tr["price"].(float64)
		trade.Type = tr["type"].(string)
		timeStr := tr["time"].(string)
		timeMeta := strings.Split(timeStr, ":")
		h, _ := strconv.Atoi(timeMeta[0])
		m, _ := strconv.Atoi(timeMeta[1])
		s, _ := strconv.Atoi(timeMeta[2])
		//临界点处理
		if now.Hour() == 0 {
			if h <= 23 && h >= 20 {
				pre := now.AddDate(0, 0, -1)
				trade.Date = time.Date(pre.Year(), pre.Month(), pre.Day(), h, m, s, 0, time.Local).Unix() * 1000
			} else if h == 0 {
				trade.Date = time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, time.Local).Unix() * 1000
			}
		} else {
			trade.Date = time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, time.Local).Unix() * 1000
		}
		//fmt.Println(time.Unix(trade.Date/1000 , 0))
		trades = append(trades, trade)
	}

	//fmt.Println(tradesmap)

	return trades, nil
}
