package poloniex

import (
	"net/http"
	. "github.com/nntaoli/crypto_coin_api"
	"log"
	"strconv"
	"fmt"
	"net/url"
	"time"
	"encoding/json"
	"strings"
	"errors"
)

const EXCHANGE_NAME = "poloniex";

const (
	BASE_URL = "https://poloniex.com/";
	TRADE_API = BASE_URL + "tradingApi";
	PUBLIC_URL = BASE_URL + "public";
	TICKER_API = "?command=returnTicker";
	ORDER_BOOK_API = "?command=returnOrderBook&currencyPair=%s&depth=%d"
)

var _CURRENCYPAIR_TO_SYMBOL = map[CurrencyPair]string{
	BTC_LTC : "BTC_LTC",
	ETH_BTC : "BTC_ETH",
	ETC_BTC : "BTC_ETC",
	XCN_BTC : "BTC_XCN",
	SYS_BTC : "BTC_SYS"};

type Poloniex struct {
	accessKey,
	secretKey string
	client    *http.Client
}

func New(client *http.Client, accessKey, secretKey string) (*Poloniex) {
	return &Poloniex{accessKey, secretKey, client};
}

func (poloniex *Poloniex) GetExchangeName() string {
	return EXCHANGE_NAME;
}

func (Poloniex *Poloniex) GetTicker(currency CurrencyPair) (*Ticker, error) {
	respmap, err := HttpGet(PUBLIC_URL + TICKER_API);
	if err != nil {
		log.Println(err);
		return nil, err;
	}

	tickermap := respmap[_CURRENCYPAIR_TO_SYMBOL[currency]].(map[string]interface{});

	ticker := new(Ticker);
	ticker.High, _ = strconv.ParseFloat(tickermap["high24hr"].(string), 64);
	ticker.Low, _ = strconv.ParseFloat(tickermap["low24hr"].(string), 64);
	ticker.Last, _ = strconv.ParseFloat(tickermap["last"].(string), 64);
	ticker.Buy, _ = strconv.ParseFloat(tickermap["highestBid"].(string), 64);
	ticker.Sell, _ = strconv.ParseFloat(tickermap["lowestAsk"].(string), 64);
	ticker.Vol, _ = strconv.ParseFloat(tickermap["quoteVolume"].(string), 64);

	log.Println(tickermap);

	return ticker, nil;
}
func (Poloniex *Poloniex) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	respmap, err := HttpGet(PUBLIC_URL + fmt.Sprintf(ORDER_BOOK_API, _CURRENCYPAIR_TO_SYMBOL[currency], size));
	if err != nil {
		log.Println(err);
		return nil, err;
	}

	log.Println(respmap);

	var depth Depth;

	for _, v := range respmap["asks"].([]interface{}) {
		var dr DepthRecord;
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price, _ = strconv.ParseFloat(vv.(string), 64);
			case 1:
				dr.Amount = vv.(float64);
			}
		}
		depth.AskList = append(depth.AskList, dr);
	}

	for _, v := range respmap["bids"].([]interface{}) {
		var dr DepthRecord;
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price, _ = strconv.ParseFloat(vv.(string), 64);
			case 1:
				dr.Amount = vv.(float64);
			}
		}
		depth.BidList = append(depth.BidList, dr);
	}

	return &depth, nil;
}
func (Poloniex *Poloniex) GetKlineRecords(currency CurrencyPair, period string, size, since int) ([]Kline, error) {
	return nil, nil; }

func (poloniex *Poloniex) placeLimitOrder(command, amount, price string, currency CurrencyPair) (*Order, error) {
	postData := url.Values{};
	postData.Set("command", command);
	postData.Set("currencyPair", _CURRENCYPAIR_TO_SYMBOL[currency]);
	postData.Set("rate", price);
	postData.Set("amount", amount);

	sign, _ := poloniex.buildPostForm(&postData);

	headers := map[string]string{
		"Key" : poloniex.accessKey,
		"Sign" : sign};
	resp, err := HttpPostForm2(http.DefaultClient, TRADE_API, postData, headers);
	if err != nil {
		log.Println(err);
		return nil, err;
	}

	respmap := make(map[string]interface{});
	err = json.Unmarshal(resp, &respmap);
	if err != nil || respmap["error"] != nil {
		log.Println(err, string(resp));
		return nil, err;
	}

	orderNumber := respmap["orderNumber"].(string);
	order := new(Order);
	order.OrderTime = int(time.Now().Unix() * 1000);
	order.OrderID, _ = strconv.Atoi(orderNumber);
	order.Amount, _ = strconv.ParseFloat(amount, 64);
	order.Price, _ = strconv.ParseFloat(price, 64);
	order.Status = ORDER_UNFINISH;
	order.Currency = currency;

	switch command {
	case "sell":
		order.Side = SELL;
	case "buy":
		order.Side = BUY;
	}

	log.Println(string(resp));
	return order, nil;
}

func (poloniex *Poloniex) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return poloniex.placeLimitOrder("buy", amount, price, currency);
}

func (poloniex *Poloniex) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return poloniex.placeLimitOrder("sell", amount, price, currency);
}

func (poloniex *Poloniex) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	postData := url.Values{};
	postData.Set("command", "cancelOrder");
	postData.Set("orderNumber", orderId);

	sign, err := poloniex.buildPostForm(&postData);
	if err != nil {
		log.Println(err)
		return false, err;
	}

	headers := map[string]string{
		"Key" : poloniex.accessKey,
		"Sign" : sign};
	resp, err := HttpPostForm2(http.DefaultClient, TRADE_API, postData, headers);
	if err != nil {
		log.Println(err);
		return false, err;
	}

	//log.Println(string(resp));

	respmap := make(map[string]interface{});
	err = json.Unmarshal(resp, &respmap);
	if err != nil || respmap["error"] != nil {
		log.Println(err, string(resp));
		return false, err;
	}

	success := int(respmap["success"].(float64));
	if success != 1 {
		log.Println(respmap);
		return false, nil;
	}

	return true, nil; }

func (poloniex *Poloniex) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	postData := url.Values{};
	postData.Set("command", "returnOrderTrades");
	postData.Set("orderNumber", orderId);

	sign, _ := poloniex.buildPostForm(&postData);

	headers := map[string]string{
		"Key" : poloniex.accessKey,
		"Sign" : sign};
	resp, err := HttpPostForm2(http.DefaultClient, TRADE_API, postData, headers);
	if err != nil {
		log.Println(err);
		return nil, err;
	}

	if strings.Contains(string(resp), "error:") {
		log.Println(string(resp));
		return nil, errors.New(string(resp));
	}

	respmap := make([]interface{}, 0);
	err = json.Unmarshal(resp, &respmap);
	if err != nil {
		log.Println(err, string(resp));
		return nil, err;
	}

	order := new(Order);
	order.OrderID, _ = strconv.Atoi(orderId);
	order.Currency = currency;

	total := 0.0;

	for _, v := range respmap {
		vv := v.(map[string]interface{});
		_amount, _ := strconv.ParseFloat(vv["amount"].(string), 64);
		_rate, _ := strconv.ParseFloat(vv["rate"].(string), 64);
		order.DealAmount += _amount;
		total += (_amount * _rate);
	}

	order.AvgPrice = total / order.DealAmount;

	return order, nil;
}

func (poloniex *Poloniex) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	postData := url.Values{};
	postData.Set("command", "returnOpenOrders");
	postData.Set("currencyPair", _CURRENCYPAIR_TO_SYMBOL[currency]);

	sign, err := poloniex.buildPostForm(&postData);
	if err != nil {
		log.Println(err);
		return nil, err;
	}

	headers := map[string]string{
		"Key" : poloniex.accessKey,
		"Sign" : sign};
	resp, err := HttpPostForm2(http.DefaultClient, TRADE_API, postData, headers);
	if err != nil {
		log.Println(err);
		return nil, err;
	}

	orderAr := make([]interface{}, 1);
	err = json.Unmarshal(resp, &orderAr);
	if err != nil {
		log.Println(err, string(resp));
		return nil, err;
	}

	orders := make([]Order, 0);
	for _, v := range orderAr {
		vv := v.(map[string]interface{});
		order := Order{};
		order.OrderID, _ = strconv.Atoi(vv["orderNumber"].(string));
		order.Amount, _ = strconv.ParseFloat(vv["amount"].(string), 64);
		order.Price, _ = strconv.ParseFloat(vv["rate"].(string), 64);
		order.Status = ORDER_UNFINISH;

		side := vv["type"].(string);
		switch side {
		case "buy":
			order.Side = BUY;
		case "sell":
			order.Side = SELL;
		}

		orders = append(orders, order);
	}

	log.Println(string(resp));
	return orders, nil;
}
func (Poloniex *Poloniex) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	return nil, nil;
}

func (poloniex *Poloniex) GetAccount() (*Account, error) {
	postData := url.Values{};
	postData.Add("command", "returnCompleteBalances");
	sign, err := poloniex.buildPostForm(&postData);
	if err != nil {
		return nil, err;
	}

	headers := map[string]string{
		"Key" : poloniex.accessKey,
		"Sign" : sign};
	resp, err := HttpPostForm2(http.DefaultClient, TRADE_API, postData, headers);

	if err != nil {
		log.Println(err);
		return nil, err;
	}

	respmap := make(map[string]interface{});
	err = json.Unmarshal(resp, &respmap);

	if err != nil || respmap["error"] != nil {
		log.Println(err);
		return nil, err;
	}

	acc := new(Account);
	acc.Exchange = EXCHANGE_NAME;
	acc.SubAccounts = make(map[Currency]SubAccount);

	for k, v := range respmap {
		var currency Currency;

		switch k {
		case "BTC":
			currency = BTC;
		case "LTC":
			currency = LTC;
		case "ETH":
			currency = ETH;
		case "ETC":
			currency = ETC;
		case "USD":
			currency = USD;
		default:
			currency = -1;
		}

		if currency > 0 {
			vv := v.(map[string]interface{});
			subAcc := SubAccount{};
			subAcc.Currency = currency;
			subAcc.Amount, _ = strconv.ParseFloat(vv["available"].(string), 64);
			subAcc.ForzenAmount, _ = strconv.ParseFloat(vv["onOrders"].(string), 64);
			acc.SubAccounts[subAcc.Currency] = subAcc;
		}
	}

	return acc, nil
}

func (poloniex *Poloniex) buildPostForm(postForm *url.Values) (string, error) {
	postForm.Add("nonce", fmt.Sprintf("%d", time.Now().UnixNano() / 1000000));
	payload := postForm.Encode();
	//println(payload)
	sign, err := GetParamHmacSHA512Sign(poloniex.secretKey, payload);
	if err != nil {
		return "", err;
	}
	//log.Println(sign)
	return sign, nil;
}
