/*注意:haobtc已经下线交易所*/
package haobtc

import (
	. "github.com/nntaoli-project/GoEx"
	"net/http"
	"errors"
	"fmt"
	"net/url"
	"encoding/json"
	"strings"
	"strconv"
)

const
(
	EXCHANGE_NAME = "haobtc.com";
	API_BASE_URL = "https://haobtc.com/exchange/api/v1/";
	TICKER_URI = "ticker";
	TRADE_URI = "trade";
	CANCEL_URI = "cancel_order";
	DEPTH_URI = "depth?size=%d";
	ACCOUNT_URI = "account_info";
	ORDER_INFO = "order_info";
	ORDERS_INFO = "orders_info";

)

type HaoBtc struct {
	httpClient *http.Client;
	accessKey,
	secretKey  string
}

func New(httpClient *http.Client, accessKey, secretKey string) *HaoBtc {
	return &HaoBtc{httpClient, accessKey, secretKey};
}

func (ctx *HaoBtc) buildPostForm(postForm *url.Values) error {
	postForm.Set("api_key", ctx.accessKey);
	//postForm.Set("secret_key", ctx.secret_key);

	payload := postForm.Encode();
	payload = payload + "&secret_key=" + ctx.secretKey;

	sign, err := GetParamMD5Sign(ctx.secretKey, payload);
	if err != nil {
		return err;
	}

	postForm.Set("sign", strings.ToUpper(sign));
	//postForm.Del("secret_key")
	return nil;
}

func (ctx *HaoBtc) GetTicker(currency CurrencyPair) (*Ticker, error) {
	if currency != BTC_CNY {
		return nil, errors.New("The HaoBtc Unsupport " + currency.String());
	}

	var tickerMap map[string]interface{};
	var ticker Ticker;

	url := API_BASE_URL + TICKER_URI;

	bodyDataMap, err := HttpGet(ctx.httpClient, url);
	if err != nil {
		return nil, err;
	}

	tickerMap = bodyDataMap["ticker"].(map[string]interface{});
	ticker.Date = uint64(bodyDataMap["date"].(float64));
	ticker.Last = tickerMap["last"].(float64);
	ticker.Buy = tickerMap["buy"].(float64);
	ticker.Sell = tickerMap["sell"].(float64);
	ticker.Low = tickerMap["low"].(float64);
	ticker.High = tickerMap["high"].(float64);
	ticker.Vol = tickerMap["vol"].(float64);

	return &ticker, nil;
}

func (ctx *HaoBtc) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	var depthUri string;

	switch currency {
	case BTC_CNY:
		depthUri = API_BASE_URL + fmt.Sprintf(DEPTH_URI, size);
	default:
		return nil, errors.New("Unsupport The CurrencyPair " + currency.ToSymbol("_"));
	}

	bodyDataMap, err := HttpGet(ctx.httpClient,depthUri);

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

func (ctx *HaoBtc) GetKlineRecords(currency CurrencyPair ,period string, size , since int) ([]Kline , error){
	return nil , errors.New("unimplement the method.");
}

func (ctx *HaoBtc) GetAccount() (*Account, error) {
	postData := url.Values{};
	ctx.buildPostForm(&postData);

	bodyData, err := HttpPostForm(ctx.httpClient, API_BASE_URL + ACCOUNT_URI, postData);
	if err != nil {
		return nil, err;
	}

	//fmt.Println(string(bodyData));

	var bodyDataMap map[string]interface{};
	err = json.Unmarshal(bodyData, &bodyDataMap);
	if err != nil {
		println(string(bodyData));
		return nil, err;
	}

	if bodyDataMap["code"] != nil {
		return nil, errors.New(string(bodyData));
	}

	account := new(Account);
	account.Exchange = ctx.GetExchangeName();

	var btcSubAccount SubAccount;
	var cnySubAccount SubAccount;

	btcSubAccount.Currency = BTC;
	btcSubAccount.Amount = bodyDataMap["exchange_btc"].(float64);
	btcSubAccount.ForzenAmount = bodyDataMap["exchange_frozen_btc"].(float64);

	cnySubAccount.Currency = CNY;
	cnySubAccount.Amount = bodyDataMap["exchange_cny"].(float64);
	cnySubAccount.ForzenAmount = bodyDataMap["exchange_frozen_cny"].(float64);

	account.SubAccounts = make(map[Currency]SubAccount, 2);
	account.SubAccounts[BTC] = btcSubAccount;
	account.SubAccounts[CNY] = cnySubAccount;

	return account, nil;
}

func (ctx *HaoBtc) placeOrder(_type , amount , price string , currency CurrencyPair)(*Order , error){
	postData := url.Values{};
	postData.Set("type" , _type);
	postData.Set("amount" , amount);
	postData.Set("price" , price);

	ctx.buildPostForm(&postData);

	bodyData, err := HttpPostForm(ctx.httpClient, API_BASE_URL + TRADE_URI, postData);
	if err != nil {
		return nil, err;
	}

	fmt.Println(string(bodyData));

	var bodyDataMap map[string]interface{};
	err = json.Unmarshal(bodyData, &bodyDataMap);
	if err != nil {
		println(string(bodyData));
		return nil, err;
	}

	if bodyDataMap["code"] != nil {
		return nil, errors.New(string(bodyData));
	}

	id := int(bodyDataMap["order_id"].(float64));

	if id <= 0 {
		return nil , errors.New("Place Order Fail.");
	}

	order := new(Order);
	order.OrderID = int(bodyDataMap["order_id"].(float64));
	order.Price, _ = strconv.ParseFloat(price, 64);
	order.Amount, _ = strconv.ParseFloat(amount, 64);
	order.Currency = currency;
	order.Status = ORDER_UNFINISH;

	switch _type {
	case "sell" ,"sell_market" , "sell_maker_only":
		order.Side = SELL;
	case "buy" , "buy_market" , "buy_maker_only":
		order.Side = BUY;
	}

	return order , nil;
}


func (ctx *HaoBtc) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error){
	return ctx.placeOrder("buy" , amount , price ,currency);
}

func (ctx *HaoBtc) LimitSell(amount, price string, currency CurrencyPair) (*Order, error){
	return ctx.placeOrder("sell" , amount , price , currency);
}

func (ctx *HaoBtc) CancelOrder(orderId string, currency CurrencyPair) (bool, error)  {
	postData := url.Values{};
	postData.Set("order_id" , orderId);

	ctx.buildPostForm(&postData);

	bodyData, err := HttpPostForm(ctx.httpClient, API_BASE_URL + CANCEL_URI, postData);
	if err != nil {
		return false, err;
	}

	//fmt.Println(string(bodyData));

	var bodyDataMap map[string]interface{};
	err = json.Unmarshal(bodyData, &bodyDataMap);
	if err != nil {
		//println(string(bodyData));
		return false, err;
	}

	if bodyDataMap["code"] != nil {
		return false , errors.New(string(bodyData));
	}

	id := int(bodyDataMap["order_id"].(float64));

	if id <= 0 {
		return false , errors.New("Cancel Order Fail.");
	}

	return false , nil;
}

func (ctx *HaoBtc) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error){
	postData := url.Values{};
	postData.Set("order_id" , orderId);

	ctx.buildPostForm(&postData);

	bodyData, err := HttpPostForm(ctx.httpClient, API_BASE_URL + ORDER_INFO, postData);
	if err != nil {
		return nil, err;
	}

	//fmt.Println(string(bodyData));

	var bodyDataMap map[string]interface{};
	err = json.Unmarshal(bodyData, &bodyDataMap);
	if err != nil {
		//println(string(bodyData));
		return nil, err;
	}

	if bodyDataMap["code"] != nil {
		return nil , errors.New(string(bodyData));
	}

	/*
	 "status": "CANCELED",
         "amount": 0.1,
  	 "create_date": "2016-03-21T10:06:06.904Z",
  	 "avg_price": 0,
   	 "order_id": 57182,
  	 "price": 3000,
         "deal_size": 0,
  	 "type": "LIMIT",
  	 "side": "BUY"
	 */

	order := new(Order);
	order.OrderID = int(bodyDataMap["order_id"].(float64));
	order.Amount = bodyDataMap["amount"].(float64);
	order.DealAmount = bodyDataMap["deal_size"].(float64);
	order.Price =  bodyDataMap["price"].(float64);
	order.AvgPrice = bodyDataMap["avg_price"].(float64);

	status := bodyDataMap["status"].(string);
	switch status {
	case "CANCELED":
		order.Status = ORDER_CANCEL;
	case "CLOSE":
		order.Status = ORDER_FINISH;
	case "OPEN", "SUBMIT", "PENDING":
		order.Status = ORDER_UNFINISH;
	}

	side := bodyDataMap["side"].(string);
	switch side {
	case "BUY":
		order.Side = BUY;
	case "SELL":
		order.Side = SELL;
	}

	return order , nil;
}

func (ctx *HaoBtc) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	postData := url.Values{};

	ctx.buildPostForm(&postData);

	bodyData, err := HttpPostForm(ctx.httpClient, API_BASE_URL + ORDERS_INFO, postData);
	if err != nil {
		return nil, err;
	}

	fmt.Println(string(bodyData));

	bodyStr := string(bodyData);
	if strings.Contains("code", bodyStr) {
		return nil, errors.New(bodyStr);
	}

	var ordersArray []interface{};
	err = json.Unmarshal(bodyData, &ordersArray);
	if err != nil {
		//println(string(bodyData));
		return nil, err;
	}

	var orders []Order;

	for _, v := range ordersArray {
		_map := v.(map[string]interface{});
		order := Order{};
		order.OrderID = int(_map["order_id"].(float64));
		order.Amount = _map["amount"].(float64);
		order.Price = _map["price"].(float64);
		order.AvgPrice = _map["avg_price"].(float64);

		remainsize := _map["remainsize"].(float64);
		order.DealAmount = order.Amount - remainsize;

		status := _map["status"].(string);
		switch status {
		case "Part deal":
			order.Status = ORDER_PART_FINISH;
		case "All deal":
			order.Status = ORDER_FINISH;
		case "No deal":
			order.Status = ORDER_UNFINISH;
		}

		side := _map["type"].(string);
		switch side {
		case "BUY":
			order.Side = BUY;
		case "SELL":
			order.Side = SELL;
		}

		orders = append(orders, order);
	}

	return orders, nil;
}

func (ctx *HaoBtc) GetOrderHistorys(currency CurrencyPair , currentPage , pageSize int)([]Order , error){
	return nil , errors.New("unimplements the method");
}

func (ctx *HaoBtc) GetExchangeName() string {
	return EXCHANGE_NAME;
}