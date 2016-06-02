package huobi

import (
	. "rest"
	"io/ioutil"
	"errors"
	"encoding/json"
	"strconv"
	"fmt"
	"net/url"
	"net/http"
	"time"
	"strings"
)

const
(
	EXCHANGE_NAME = "huobi";
	API_BASE_URL = "https://api.huobi.com/";
	TRADE_API_V3 = API_BASE_URL + "apiv3";
	TICKER_URI = "staticmarket/ticker_%s_json.js";
	DEPTH_URI = "staticmarket/depth_%s_%d.js";
)

type HuoBi struct {
	httpClient *http.Client;
	accessKey,
	secretKey  string
}

func New(httpClient *http.Client, accessKey, secretKey string) *HuoBi {
	return &HuoBi{httpClient, accessKey, secretKey};
}

func httpGet(uri string, client *http.Client) (map[string]interface{}, error) {
	url := API_BASE_URL + uri;
	//println(url);
	resp, err := client.Get(url);
	if err != nil {
		return nil, err;
	}

	defer resp.Body.Close();

	body, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		return nil, err;
	}

	var bodyDataMap map[string]interface{};
	json.Unmarshal(body, &bodyDataMap);
	return bodyDataMap, nil;
}

func (hb *HuoBi) buildPostForm(postForm *url.Values) error {
	postForm.Set("created", fmt.Sprintf("%d", time.Now().Unix()));
	postForm.Set("access_key", hb.accessKey);
	postForm.Set("secret_key", hb.secretKey);
	sign, err := GetParamMD5Sign(hb.secretKey, postForm.Encode());
	if err != nil {
		return err;
	}
	postForm.Set("sign", sign);
	postForm.Del("secret_key");
	return nil;
}

func (hb *HuoBi) GetExchangeName() string {
	return EXCHANGE_NAME;
}

func (hb *HuoBi) GetTicker(currency CurrencyPair) (*Ticker, error) {
	var tickerUri string;

	switch currency {
	case BTC_CNY:
		tickerUri = fmt.Sprintf(TICKER_URI, "btc");
	case LTC_CNY:
		tickerUri = fmt.Sprintf(TICKER_URI, "ltc");
	default:
		return nil, errors.New("Unsupport The CurrencyPair");
	}

	bodyDataMap, err := httpGet(tickerUri, hb.httpClient);

	if err != nil {
		return nil, err;
	}

	var tickerMap map[string]interface{};
	var ticker Ticker;

	switch  bodyDataMap["ticker"].(type) {
	case map[string]interface{}:
		tickerMap = bodyDataMap["ticker"].(map[string]interface{});
	default:
		return nil, errors.New(fmt.Sprintf("Type Convert Error ? \n %s", bodyDataMap));
	}

	ticker.Date, _ = strconv.ParseUint(bodyDataMap["time"].(string), 10, 64);
	ticker.Last = tickerMap["last"].(float64);
	ticker.Buy = tickerMap["buy"].(float64);
	ticker.Sell = tickerMap["sell"].(float64);
	ticker.Low = tickerMap["low"].(float64);
	ticker.High = tickerMap["high"].(float64);
	ticker.Vol = tickerMap["vol"].(float64);

	return &ticker, nil;
}

func (hb *HuoBi) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	var depthUri string;

	switch currency {
	case BTC_CNY:
		depthUri = fmt.Sprintf(DEPTH_URI, "btc", size);
	case LTC_CNY:
		depthUri = fmt.Sprintf(DEPTH_URI, "ltc", size);
	default:
		return nil, errors.New("Unsupport The CurrencyPair");
	}

	bodyDataMap, err := httpGet(depthUri, hb.httpClient);

	if err != nil {
		return nil, err;
	}

	var depth Depth;

	for _, v := range bodyDataMap["asks"].([]interface{}) {
		var dr DepthRecord;
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price = vv.(float64);
			case 1:
				dr.Amount = vv.(float64);
			}
		}
		depth.AskList = append(depth.AskList, dr);
	}

	for _, v := range bodyDataMap["bids"].([]interface{}) {
		var dr DepthRecord;
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price = vv.(float64);
			case 1:
				dr.Amount = vv.(float64);
			}
		}
		depth.BidList = append(depth.BidList, dr);
	}

	return &depth, nil;
}

func (hb *HuoBi) GetAccount() (*Account, error) {
	postData := url.Values{};
	postData.Set("method", "get_account_info");
	postData.Set("created", fmt.Sprintf("%d", time.Now().Unix()));
	postData.Set("access_key", hb.accessKey);
	postData.Set("secret_key", hb.secretKey);

	sign, _ := GetParamMD5Sign(hb.secretKey, postData.Encode());
	postData.Set("sign", sign);
	postData.Del("secret_key");

	bodyData, err := HttpPostForm(hb.httpClient, TRADE_API_V3, postData);
	if err != nil {
		return nil, err;
	}

	var bodyDataMap map[string]interface{};
	err = json.Unmarshal(bodyData, &bodyDataMap);
	if err != nil {
		println(string(bodyData));
		return nil, err;
	}

	//fmt.Println(bodyDataMap);

	account := new(Account);
	account.Exchange = hb.GetExchangeName();
	account.Asset, _ = strconv.ParseFloat(bodyDataMap["total"].(string), 64);
	account.NetAsset, _ = strconv.ParseFloat(bodyDataMap["net_asset"].(string), 64);

	var btcSubAccount SubAccount;
	var ltcSubAccount SubAccount;
	var cnySubAccount SubAccount;

	btcSubAccount.Currency = BTC;
	btcSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["available_btc_display"].(string), 64);
	btcSubAccount.LoanAmount, _ = strconv.ParseFloat(bodyDataMap["loan_btc_display"].(string), 64);
	btcSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["frozen_btc_display"].(string), 64);

	ltcSubAccount.Currency = LTC;
	ltcSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["available_ltc_display"].(string), 64);
	ltcSubAccount.LoanAmount, _ = strconv.ParseFloat(bodyDataMap["loan_ltc_display"].(string), 64);
	ltcSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["frozen_ltc_display"].(string), 64);

	cnySubAccount.Currency = CNY;
	cnySubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["available_cny_display"].(string), 64);
	cnySubAccount.LoanAmount, _ = strconv.ParseFloat(bodyDataMap["loan_cny_display"].(string), 64);
	cnySubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["frozen_cny_display"].(string), 64);

	account.SubAccounts = make(map[Currency]SubAccount, 3);
	account.SubAccounts[BTC] = btcSubAccount;
	account.SubAccounts[LTC] = ltcSubAccount;
	account.SubAccounts[CNY] = cnySubAccount;

	return account, nil;
}

func (hb *HuoBi) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	postData := url.Values{};
	postData.Set("method", "order_info");
	postData.Set("id", orderId);

	switch currency {
	case BTC_CNY:
		postData.Set("coin_type", "1");
	case LTC_CNY:
		postData.Set("coin_type", "2");
	}

	hb.buildPostForm(&postData);

	bodyData, _ := HttpPostForm(hb.httpClient, TRADE_API_V3, postData);

	var bodyDataMap map[string]interface{};
	err := json.Unmarshal(bodyData, &bodyDataMap);
	if err != nil {
		println(string(bodyData));
		return nil, err;
	}

	//fmt.Println(bodyDataMap);
	order := new(Order);
	order.OrderID, _ = strconv.Atoi(orderId);
	order.Side = TradeSide(bodyDataMap["type"].(float64));
	order.Amount, _ = strconv.ParseFloat(bodyDataMap["order_amount"].(string), 64);
	order.DealAmount, _ = strconv.ParseFloat(bodyDataMap["processed_amount"].(string), 64);
	order.Price, _ = strconv.ParseFloat(bodyDataMap["order_price"].(string), 64);
	order.AvgPrice, _ = strconv.ParseFloat(bodyDataMap["processed_price"].(string), 64);

	tradeStatus := TradeStatus(bodyDataMap["status"].(float64));
	switch tradeStatus {
	case 0:
		order.Status = ORDER_UNFINISH;
	case 1:
		order.Status = ORDER_PART_FINISH;
	case 2:
		order.Status = ORDER_FINISH;
	case 3:
		order.Status = ORDER_CANCEL;
	}

	return order, nil;
}

func (hb *HuoBi) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	postData := url.Values{};
	postData.Set("method", "get_orders");

	switch currency {
	case BTC_CNY:
		postData.Set("coin_type", "1");
	case LTC_CNY:
		postData.Set("coin_type", "2");
	}

	hb.buildPostForm(&postData);

	bodyData, err := HttpPostForm(hb.httpClient, TRADE_API_V3, postData);
	if err != nil {
		return nil, err;
	}

	//println(string(bodyData))

	var bodyDataMap []map[string]interface{};
	err = json.Unmarshal(bodyData, &bodyDataMap);
	if err != nil {
		return nil, err;
	}

	var orders []Order;

	for _, v := range bodyDataMap {
		order := Order{};
		order.Amount, _ = strconv.ParseFloat(v["order_amount"].(string), 64);
		order.Price, _ = strconv.ParseFloat(v["order_price"].(string), 64);
		order.DealAmount, _ = strconv.ParseFloat(v["processed_amount"].(string), 64);
		order.OrderTime = int(v["order_time"].(float64));
		order.OrderID = int(v["id"].(float64));
		order.Side = TradeSide(v["type"].(float64));
		orders = append(orders, order);
	}

	return orders, nil;
}

func (hb *HuoBi) placeOrder(method, amount, price string, currency CurrencyPair) (*Order, error) {
	postData := url.Values{};
	postData.Set("method", method);
	postData.Set("price", price);
	postData.Set("amount", amount);

	switch currency {
	case BTC_CNY:
		postData.Set("coin_type", "1");
	case LTC_CNY:
		postData.Set("coin_type", "2");
	}

	hb.buildPostForm(&postData);

	bodyData, err := HttpPostForm(hb.httpClient, TRADE_API_V3, postData);
	if err != nil {
		return nil, err;
	}

	//{"result":"success","id":1321475746}
	//println(string(bodyData))

	var bodyDataMap map[string]interface{};
	err = json.Unmarshal(bodyData, &bodyDataMap);

	if err != nil {
		return nil, err;
	}

	if bodyDataMap["code"] != nil {
		return nil, errors.New(string(bodyData));
	}

	ret := bodyDataMap["result"].(string);

	if strings.Compare(ret, "success") == 0 {
		order := new(Order);
		order.OrderID = int(bodyDataMap["id"].(float64));
		order.Price, _ = strconv.ParseFloat(price, 64);
		order.Amount, _ = strconv.ParseFloat(amount, 64);
		order.Currency = currency;
		order.Status = ORDER_UNFINISH;
		return order, nil;
	}

	return nil, errors.New(fmt.Sprintf("Place Limit %s Order Fail.", method));
}

func (hb *HuoBi) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	order, err := hb.placeOrder("buy", amount, price, currency);

	if err != nil {
		return nil, err;
	}

	order.Side = BUY;

	return order, nil;
}

func (hb *HuoBi) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	order, err := hb.placeOrder("sell", amount, price, currency);

	if err != nil {
		return nil, err;
	}

	order.Side = SELL;

	return order, nil;
}

func (hb *HuoBi) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	//1321490762
	postData := url.Values{};
	postData.Set("method", "cancel_order");
	postData.Set("id", orderId);

	switch currency {
	case BTC_CNY:
		postData.Set("coin_type", "1");
	case LTC_CNY:
		postData.Set("coin_type", "2");
	}

	hb.buildPostForm(&postData);

	bodyData, err := HttpPostForm(hb.httpClient, TRADE_API_V3, postData);
	if err != nil {
		return false, err;
	}

	//{"result":"success"}
	//{"code":42,"msg":"该委托已经取消, 不能取消或修改","message":"该委托已经取消, 不能取消或修改"}
	println(string(bodyData))

	var bodyDataMap map[string]interface{};
	err = json.Unmarshal(bodyData, &bodyDataMap);
	if err != nil {
		return false, err;
	}

	if bodyDataMap["code"] != nil {
		return false, errors.New(string(bodyData));
	}

	ret := bodyDataMap["result"].(string);
	return (strings.Compare(ret, "success") == 0), nil;
}