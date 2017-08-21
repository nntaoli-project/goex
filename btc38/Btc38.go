package btc38

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	. "github.com/nntaoli-project/GoEx"
)

//因服务器有防CC攻击策略，每60秒内调用次数不可超过120次，超过部分将被防火墙拦截。
const (
	EXCHANGE_NAME = "btc38.com"

	API_BASE_URL = "http://api.btc38.com/"
	API_V1       = API_BASE_URL + "v1/"

	//	BASE_URL    = "http://api.btc38.com/v1/"
	TICKER_URI             = "ticker.php?c=%s&mk_type=%s"
	DEPTH_URI              = "depth.php?c=%s&mk_type=%s"
	ACCOUNT_URI            = "getMyBalance.php"
	TRADE_URI              = "trades.php?c=%s&mk_type=%s"
	CANCEL_URI             = "cancelOrder.php"
	ORDERS_INFO            = "getMyTradeList.php"
	UNFINISHED_ORDERS_INFO = "getOrderList.php"
)

type Btc38 struct {
	accessKey,
	secretKey,
	accountId string
	httpClient *http.Client
}

func New(client *http.Client, accessKey, secretKey, accountId string) *Btc38 {
	return &Btc38{accessKey, secretKey, accountId, client}
}

func (btc38 *Btc38) GetExchangeName() string {
	return EXCHANGE_NAME
}

func (btc38 *Btc38) buildPostForm(postForm *url.Values) error {
	postForm.Set("created", fmt.Sprintf("%d", time.Now().Unix()))
	postForm.Set("access_key", btc38.accessKey)
	postForm.Set("secret_key", btc38.secretKey)
	sign, err := GetParamMD5Sign(btc38.secretKey, postForm.Encode())
	if err != nil {
		return err
	}
	postForm.Set("sign", sign)
	postForm.Del("secret_key")
	return nil
}

func (btc38 *Btc38) GetTicker(currency CurrencyPair) (*Ticker, error) {
	cur:= currency.CurrencyA.String()
	money := currency.CurrencyB.String()
	if cur == "err" {
		log.Println("Unsupport The CurrencyPair")
		return nil, errors.New("Unsupport The CurrencyPair")
	}
	tickerUri := API_V1 + fmt.Sprintf(TICKER_URI, cur, money)
	timestamp := time.Now().Unix()
	bodyDataMap, err := HttpGet(btc38.httpClient, tickerUri)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	//	log.Println("Btc38 bodyDataMap:", bodyDataMap)
	var tickerMap map[string]interface{}
	var ticker Ticker

	switch bodyDataMap["ticker"].(type) {
	case map[string]interface{}:
		tickerMap = bodyDataMap["ticker"].(map[string]interface{})
	default:
		return nil, errors.New(fmt.Sprintf("Type Convert Error ? \n %s", bodyDataMap))
	}

	ticker.Date = uint64(timestamp)
	ticker.Last = tickerMap["last"].(float64)
	ticker.Buy = tickerMap["buy"].(float64)
	ticker.Sell = tickerMap["sell"].(float64)
	ticker.Low = tickerMap["low"].(float64)
	ticker.High = tickerMap["high"].(float64)
	ticker.Vol = tickerMap["vol"].(float64)

	return &ticker, nil
}

func (btc38 *Btc38) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	var depthUri string
	cur:= currency.CurrencyA.String()
	money := currency.CurrencyB.String()
	if cur == "err" {
		log.Println("Unsupport The CurrencyPair")
		return nil, errors.New("Unsupport The CurrencyPair")
	}
	depthUri = fmt.Sprintf(API_V1+DEPTH_URI, cur, money)

	bodyDataMap, err := HttpGet(btc38.httpClient, depthUri)

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

func (btc38 *Btc38) GetAccount() (*Account, error) {
	postData := url.Values{}
	postData.Set("key", btc38.accessKey)
	timeNow := fmt.Sprintf("%d", time.Now().Unix())
	postData.Set("time", timeNow)

	mdt := fmt.Sprintf("%s_%s_%s_%s", btc38.accessKey, btc38.accountId, btc38.secretKey, timeNow)
	sign, _ := GetParamMD5Sign(btc38.secretKey, mdt)
	postData.Set("md5", sign)
	fmt.Println("postData:", postData)
	accountUri := fmt.Sprintf(API_V1 + ACCOUNT_URI)
	//fmt.Println(accountUri)
	bodyData, err := HttpPostForm(btc38.httpClient, accountUri, postData)
	if err != nil {
		fmt.Println("err:", err)
		return nil, err
	}
	var bodyDataMap map[string]interface{}

	err = json.Unmarshal(bodyData, &bodyDataMap)
	if err != nil {
		println(string(bodyData))
		fmt.Println("err:", err)
		return nil, err
	}
	fmt.Println("bodyDataMap:", bodyDataMap)
	if bodyDataMap["code"] != nil {
		return nil, errors.New(fmt.Sprintf("%s", bodyDataMap))
	}

	account := new(Account)
	account.Exchange = btc38.GetExchangeName()

	account.SubAccounts = make(map[Currency]SubAccount, 50)

	var btcSubAccount SubAccount
	btcSubAccount.Currency = BTC
	btcSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["btc_balance"].(string), 64)
	btcSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["btc_balance_lock"].(string), 64)
	account.SubAccounts[btcSubAccount.Currency] = btcSubAccount

	return account, nil
}

func (btc38 *Btc38) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	postData := url.Values{}
	cur:= currency.CurrencyA.String()
	money := currency.CurrencyB.String()
	postData.Set("key", btc38.accessKey)
	timeNow := fmt.Sprintf("%d", time.Now().Unix())
	postData.Set("time", timeNow)

	mdt := fmt.Sprintf("%s_%s_%s_%s", btc38.accessKey, btc38.accountId, btc38.secretKey, timeNow)
	sign, _ := GetParamMD5Sign(btc38.secretKey, mdt)
	postData.Set("md5", sign)
	postData.Set("mk_type", money)
	postData.Set("coinname", cur)

	orderInfoUri := fmt.Sprintf(API_V1 + ORDERS_INFO)
	bodyData, err := HttpPostForm(btc38.httpClient, orderInfoUri, postData)
	if err != nil {
		fmt.Println("err:", err)
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

func (btc38 *Btc38) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	//cur:= currency.CurrencyA.String()
	money := currency.CurrencyB.String()
	postData := url.Values{}
	postData.Set("key", btc38.accessKey)
	timeNow := fmt.Sprintf("%d", time.Now().Unix())
	postData.Set("time", timeNow)

	mdt := fmt.Sprintf("%s_%s_%s_%s", btc38.accessKey, btc38.accountId, btc38.secretKey, timeNow)
	sign, _ := GetParamMD5Sign(btc38.secretKey, mdt)
	postData.Set("md5", sign)
	postData.Set("mk_type", strings.ToUpper(money))
	postData.Set("coinname", "")

	orderList := fmt.Sprintf(API_V1 + UNFINISHED_ORDERS_INFO)
	bodyData, err := HttpPostForm(btc38.httpClient, orderList, postData)
	if err != nil {
		fmt.Println("err:", err)
		return nil, err
	}

	if strings.Contains(string(bodyData), "code") {
		return nil, errors.New(string(bodyData))
	}
	//	fmt.Println("bodyData:", string(bodyData))

	var bodyDataMap []map[string]interface{}
	err = json.Unmarshal(bodyData, &bodyDataMap)
	if err != nil {
		return nil, err
	}

	var orders []Order

	for _, v := range bodyDataMap {
		order := Order{}
		//cur := fmt.Sprintf("%s_cny", v["coinname"].(string))
		order.Currency = currency
		order.Amount, _ = strconv.ParseFloat(v["amount"].(string), 64)
		order.Price, _ = strconv.ParseFloat(v["price"].(string), 64)
		t1 := v["time"].(string)
		t2, _ := time.Parse("2006-01-02 15:04:05", t1)
		order.OrderTime = (int)(t2.Unix())
		id, _ := strconv.ParseInt(v["id"].(string), 10, 64)
		order.OrderID = (int)(id)

		types, _ := strconv.ParseInt(v["type"].(string), 10, 64)
		order.Side = TradeSide(types)
		orders = append(orders, order)
		fmt.Println("order:", order)
	}

	return orders, nil
}

func (btc38 *Btc38) placeOrder(method, amount, price string, currency CurrencyPair) (*Order, error) {

	return nil, errors.New("unimplements")
}

func (btc38 *Btc38) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return nil, errors.New("unimplements")
}

func (btc38 *Btc38) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return nil, errors.New("unimplements")
}

func (btc38 *Btc38) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return nil, errors.New("unimplements")
}

func (btc38 *Btc38) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	order, err := btc38.placeOrder("sell_market", amount, price, currency)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	order.Side = SELL
	return order, nil
}

func (btc38 *Btc38) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	return false, nil
}

func (btc38 *Btc38) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("unimplements")
}

/**
 * 获取全站最近的交易记录
 */
func (btc38 *Btc38) GetTrades(currency CurrencyPair, since int64) ([]Trade, error) {
	cur := currency.CurrencyA.String()
	money := currency.CurrencyB.String()

	if cur == "err" {
		log.Println("Unsupport The CurrencyPair")
		return nil, errors.New("Unsupport The CurrencyPair")
	}
	tradeUrl := fmt.Sprintf(API_V1+TRADE_URI, cur, money)

	var respmap map[string]interface{}

	resp, err := http.Get(tradeUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("body", string(body))

	err = json.Unmarshal(body, &respmap)
	fmt.Println("err", err)

	if err != nil {
		return nil, err
	}

	var trades []Trade
	for _, t := range respmap {
		tr := t.(map[string]interface{})
		trade := Trade{}
		trade.Amount = tr["amount"].(float64)
		trade.Price = tr["price"].(float64)
		trade.Type = tr["type"].(string)
		trade.Date = tr["date"].(int64)
		trades = append(trades, trade)
	}
	return trades, nil
}
