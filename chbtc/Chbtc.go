package chbtc

import (
	"net/http"
	. "github.com/nntaoli/crypto_coin_api"
	"fmt"
	"strconv"
	"log"
	"net/url"
	"time"
	"encoding/json"
	"strings"
	"errors"
)

const
(
	MARKET_URL = "http://api.chbtc.com/data/v1/"
	TICKER_API = "ticker?currency=%s"
	DEPTH_API = "depth?currency=%s&size=%d"

	TRADE_URL = "https://trade.chbtc.com/api/"
	GET_ACCOUNT_API = "getAccountInfo"
	GET_ORDER_API = "getOrder"
	GET_UNFINISHED_ORDERS_API = "getUnfinishedOrdersIgnoreTradeType"
	CANCEL_ORDER_API = "cancelOrder"
	PLACE_ORDER_API = "order"
	WITHDRAW_API = "withdraw"
	CANCELWITHDRAW_API = "cancelWithdraw"
)

type Chbtc struct {
	httpClient *http.Client;
	accessKey,
	secretKey  string
}

func New(httpClient *http.Client, accessKey, secretKey string) *Chbtc {
	return &Chbtc{httpClient, accessKey, secretKey};
}

func (chbtc *Chbtc)GetExchangeName() string {
	return "chbtc";
}

func (chbtc *Chbtc) GetTicker(currency CurrencyPair) (*Ticker, error) {
	resp, err := HttpGet(chbtc.httpClient , MARKET_URL + fmt.Sprintf(TICKER_API, CurrencyPairSymbol[currency]));
	if err != nil {
		return nil, err;
	}
	//log.Println(resp)
	tickermap := resp["ticker"].(map[string]interface{});

	ticker := new(Ticker);
	ticker.Date, _ = strconv.ParseUint(resp["date"].(string), 10, 64);
	ticker.Buy, _ = strconv.ParseFloat(tickermap["buy"].(string), 64);
	ticker.Sell, _ = strconv.ParseFloat(tickermap["sell"].(string), 64);
	ticker.Last, _ = strconv.ParseFloat(tickermap["last"].(string), 64);
	ticker.High, _ = strconv.ParseFloat(tickermap["high"].(string), 64);
	ticker.Low, _ = strconv.ParseFloat(tickermap["low"].(string), 64);
	ticker.Vol, _ = strconv.ParseFloat(tickermap["vol"].(string), 64);

	return ticker, nil;
}

func (chbtc *Chbtc) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	resp, err := HttpGet(chbtc.httpClient , MARKET_URL + fmt.Sprintf(DEPTH_API, CurrencyPairSymbol[currency], size));
	if err != nil {
		return nil, err
	}

	//log.Println(resp);

	asks := resp["asks"].([]interface{});
	bids := resp["bids"].([]interface{});

	log.Println(asks)
	log.Println(bids)

	depth := new(Depth)

	for _, e := range bids {
		var r DepthRecord;
		ee := e.([]interface{});
		r.Amount = ee[0].(float64);
		r.Price = ee[1].(float64);

		depth.BidList = append(depth.BidList, r);
	}

	for _, e := range asks {
		var r DepthRecord;
		ee := e.([]interface{});
		r.Amount = ee[0].(float64);
		r.Price = ee[1].(float64);

		depth.AskList = append(depth.AskList, r);
	}

	return depth, nil;
}

func (chbtc *Chbtc) buildPostForm(postForm *url.Values) error {
	postForm.Set("accesskey", chbtc.accessKey);

	payload := postForm.Encode();
	secretkeySha, _ := GetSHA(chbtc.secretKey);

	sign, err := GetParamHmacMD5Sign(secretkeySha, payload);
	if err != nil {
		return err;
	}

	postForm.Set("sign", sign);
	//postForm.Del("secret_key")
	postForm.Set("reqTime", fmt.Sprintf("%d", time.Now().UnixNano() / 1000000));
	return nil;
}

func (chbtc *Chbtc) GetAccount() (*Account, error) {
	params := url.Values{};
	params.Set("method", "getAccountInfo");
	chbtc.buildPostForm(&params);
	//log.Println(params.Encode())
	resp, err := HttpPostForm(chbtc.httpClient, TRADE_URL + GET_ACCOUNT_API, params);
	if err != nil {
		return nil, err;
	}

	var respmap map[string]interface{};
	err = json.Unmarshal(resp, &respmap);
	if err != nil {
		log.Println("json unmarshal error");
		return nil, err;
	}

	acc := new(Account);
	acc.Exchange = "chbtc";
	acc.SubAccounts = make(map[Currency]SubAccount);

	resultmap := respmap["result"].(map[string]interface{});
	balancemap := resultmap["balance"].(map[string]interface{});
	frozenmap := resultmap["frozen"].(map[string]interface{});
	p2pmap := resultmap["p2p"].(map[string]interface{});
	netAssets := resultmap["netAssets"].(float64);
	asset := resultmap["totalAssets"].(float64);

	acc.NetAsset = netAssets;
	acc.Asset = asset;

	for t, v := range balancemap {
		vv := v.(map[string]interface{});
		subAcc := SubAccount{};
		subAcc.Amount = vv["amount"].(float64);

		switch t {
		case "CNY":
			subAcc.Currency = CNY;
			cnyfrozen := frozenmap["CNY"].(map[string]interface{});
			subAcc.ForzenAmount = cnyfrozen["amount"].(float64);
			subAcc.LoanAmount = p2pmap["inCNY"].(float64);
		case "BTC":
			subAcc.Currency = BTC;
			btcfrozen := frozenmap["BTC"].(map[string]interface{});
			subAcc.ForzenAmount = btcfrozen["amount"].(float64);
			subAcc.LoanAmount = p2pmap["inBTC"].(float64);
		case "LTC":
			subAcc.Currency = LTC;
			ltcfrozen := frozenmap["LTC"].(map[string]interface{});
			subAcc.ForzenAmount = ltcfrozen["amount"].(float64);
			subAcc.LoanAmount = p2pmap["inLTC"].(float64);
		case "ETH":
			subAcc.Currency = ETH;
			ethfrozen := frozenmap["ETH"].(map[string]interface{});
			subAcc.ForzenAmount = ethfrozen["amount"].(float64);
			subAcc.LoanAmount = p2pmap["inETH"].(float64);
		case "ETC":
			subAcc.Currency = ETC;
			etcfrozen := frozenmap["ETC"].(map[string]interface{});
			subAcc.ForzenAmount = etcfrozen["amount"].(float64);
			subAcc.LoanAmount = p2pmap["inETC"].(float64);
		default:
			log.Println("unknown ", t);

		}
		acc.SubAccounts[subAcc.Currency] = subAcc;
	}

	//log.Println(string(resp))
	//log.Println(acc)

	return acc, nil;
}

func (chbtc *Chbtc) placeOrder(amount, price string, currency CurrencyPair, tradeType int) (*Order, error) {
	params := url.Values{};
	params.Set("method", "order");
	params.Set("price", price);
	params.Set("amount", amount);
	params.Set("currency", CurrencyPairSymbol[currency]);
	params.Set("tradeType", fmt.Sprintf("%d", tradeType));
	chbtc.buildPostForm(&params);

	resp, err := HttpPostForm(chbtc.httpClient, TRADE_URL + PLACE_ORDER_API, params);
	if err != nil {
		log.Println(err)
		return nil, err;
	}

	log.Println(string(resp));

	respmap := make(map[string]interface{});
	err = json.Unmarshal(resp, &respmap);
	if err != nil {
		log.Println(err);
		return nil, err;
	}

	code := respmap["code"].(float64);
	if code != 1000 {
		log.Println(string(resp));
		return nil, errors.New(fmt.Sprintf("%.0f", code));
	}

	orid := respmap["id"].(string);

	order := new(Order);
	order.Amount, _ = strconv.ParseFloat(amount, 64);
	order.Price, _ = strconv.ParseFloat(price, 64);
	order.Status = ORDER_UNFINISH;
	order.Currency = currency;
	order.OrderTime = int(time.Now().UnixNano() / 1000000);
	order.OrderID, _ = strconv.Atoi(orid);

	switch tradeType {
	case 0:
		order.Side = SELL;
	case 1:
		order.Side = BUY;
	}

	return order, nil
}

func (chbtc *Chbtc) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return chbtc.placeOrder(amount, price, currency, 1);
}

func (chbtc *Chbtc) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return chbtc.placeOrder(amount, price, currency, 0);
}

func (chbtc *Chbtc) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	params := url.Values{};
	params.Set("method", "cancelOrder");
	params.Set("id", orderId);
	params.Set("currency", CurrencyPairSymbol[currency]);
	chbtc.buildPostForm(&params);

	resp, err := HttpPostForm(chbtc.httpClient, TRADE_URL + CANCEL_ORDER_API, params);
	if err != nil {
		log.Println(err);
		return false, err;
	}

	respmap := make(map[string]interface{});
	err = json.Unmarshal(resp, &respmap);
	if err != nil {
		log.Println(err);
		return false, err;
	}

	code := respmap["code"].(float64);

	if code == 1000 {
		return true, nil;
	}

	log.Println(respmap)
	return false, errors.New(fmt.Sprintf("%.0f", code));
}

func parseOrder(order *Order, ordermap map[string]interface{}) {
	//order.Currency = currency;
	order.OrderID, _ = strconv.Atoi(ordermap["id"].(string));
	order.Amount = ordermap["total_amount"].(float64);
	order.DealAmount = ordermap["trade_amount"].(float64);
	order.Price = ordermap["price"].(float64);
	if order.DealAmount > 0 {
		order.AvgPrice = ordermap["trade_money"].(float64) / order.DealAmount;
	} else {
		order.AvgPrice = 0;
	}
	order.Status = TradeStatus(ordermap["status"].(float64));
	order.OrderTime = int(ordermap["trade_date"].(float64));

	orType := ordermap["type"].(float64);
	switch orType {
	case 0:
		order.Side = SELL;
	case 1:
		order.Side = BUY;
	default:
		log.Printf("unknown order type %f", orType);
	}

}

func (chbtc *Chbtc) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	params := url.Values{};
	params.Set("method", "getOrder");
	params.Set("id", orderId);
	params.Set("currency", CurrencyPairSymbol[currency]);
	chbtc.buildPostForm(&params);

	resp, err := HttpPostForm(chbtc.httpClient, TRADE_URL + GET_ORDER_API, params);
	if err != nil {
		log.Println(err);
		return nil, err;
	}

	ordermap := make(map[string]interface{});
	err = json.Unmarshal(resp, &ordermap);
	if err != nil {
		log.Println(err);
		return nil, err;
	}

	order := new(Order);
	order.Currency = currency;

	parseOrder(order, ordermap);

	return order, nil;
}

func (chbtc *Chbtc) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	params := url.Values{};
	params.Set("method", "getUnfinishedOrdersIgnoreTradeType");
	params.Set("currency", CurrencyPairSymbol[currency]);
	params.Set("pageIndex", "1");
	params.Set("pageSize", "100");
	chbtc.buildPostForm(&params);

	resp, err := HttpPostForm(chbtc.httpClient, TRADE_URL + GET_UNFINISHED_ORDERS_API, params);
	if err != nil {
		log.Println(err);
		return nil, err;
	}

	respstr := string(resp);
	if strings.Contains(respstr, "\"code\":3001") {
		log.Println(respstr);
		return nil, nil;
	}

	var resps []interface{};
	err = json.Unmarshal(resp, &resps);
	if err != nil {
		log.Println(err);
		return nil, err;
	}

	var orders []Order;
	for _, v := range resps {
		ordermap := v.(map[string]interface{});
		order := Order{};
		order.Currency = currency;
		parseOrder(&order, ordermap);
		orders = append(orders, order);
	}

	return orders, nil;
}

func (chbtc *Chbtc) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	return nil, nil;
}

func (chbtc *Chbtc) GetKlineRecords(currency CurrencyPair, period string, size, since int) ([]Kline, error) {
	return nil, nil;
}

func (chbtc *Chbtc) Withdraw(amount string, currency Currency, fees, receiveAddr, safePwd string) (string, error) {
	params := url.Values{};
	params.Set("method", "withdraw");
	params.Set("currency", strings.ToLower(currency.String()));
	params.Set("amount", amount);
	params.Set("fees", fees);
	params.Set("receiveAddr", receiveAddr);
	params.Set("safePwd", safePwd);
	chbtc.buildPostForm(&params);

	resp, err := HttpPostForm(chbtc.httpClient, TRADE_URL + WITHDRAW_API, params)
	if err != nil {
		log.Println("withdraw fail.", err)
		return "", err
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		log.Println(err, string(resp))
		return "", err
	}

	if respMap["code"].(float64) == 1000 {
		return respMap["id"].(string), nil;
	}

	return "", errors.New(string(resp))
}

func (chbtc *Chbtc) CancelWithdraw(id string, currency Currency, safePwd string) (bool, error) {
	params := url.Values{};
	params.Set("method", "cancelWithdraw");
	params.Set("currency", strings.ToLower(currency.String()));
	params.Set("downloadId", id);
	params.Set("safePwd", safePwd);
	chbtc.buildPostForm(&params);

	resp, err := HttpPostForm(chbtc.httpClient, TRADE_URL + CANCELWITHDRAW_API, params)
	if err != nil {
		log.Println("cancel withdraw fail.", err)
		return false, err
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		log.Println(err, string(resp))
		return false, err
	}

	if respMap["code"].(float64) == 1000 {
		return true, nil;
	}

	return false, errors.New(string(resp))
}
