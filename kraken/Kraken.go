package kraken

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type BaseResponse struct {
	Error  []string    `json:"error"`
	Result interface{} `json:"result"`
}

type NewOrderResponse struct {
	Description interface{} `json:"descr"`
	TxIds       []string    `json:"txid"`
}

type Kraken struct {
	httpClient *http.Client
	accessKey,
	secretKey string
}

var (
	BASE_URL   = "https://api.kraken.com"
	API_V0     = "/0/"
	API_DOMAIN = BASE_URL + API_V0
	PUBLIC     = "public/"
	PRIVATE    = "private/"
)

func New(client *http.Client, accesskey, secretkey string) *Kraken {
	return &Kraken{client, accesskey, secretkey}
}

func (k *Kraken) placeOrder(orderType, side, amount, price string, pair CurrencyPair) (*Order, error) {
	apiuri := "private/AddOrder"

	params := url.Values{}
	params.Set("pair", k.convertPair(pair).ToSymbol(""))
	params.Set("type", side)
	params.Set("ordertype", orderType)
	params.Set("price", price)
	params.Set("volume", amount)

	var resp NewOrderResponse
	err := k.doAuthenticatedRequest("POST", apiuri, params, &resp)
	//log.Println
	if err != nil {
		return nil, err
	}

	var tradeSide TradeSide = SELL
	if "buy" == side {
		tradeSide = BUY
	}

	return &Order{
		Currency: pair,
		OrderID2: resp.TxIds[0],
		Amount:   ToFloat64(amount),
		Price:    ToFloat64(price),
		Side:     tradeSide,
		Status:   ORDER_UNFINISH}, nil
}

func (k *Kraken) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return k.placeOrder("limit", "buy", amount, price, currency)
}

func (k *Kraken) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return k.placeOrder("limit", "sell", amount, price, currency)
}

func (k *Kraken) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return k.placeOrder("market", "buy", amount, price, currency)
}

func (k *Kraken) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return k.placeOrder("market", "sell", amount, price, currency)
}

func (k *Kraken) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	params := url.Values{}
	apiuri := "private/CancelOrder"
	params.Set("txid", orderId)

	var respmap map[string]interface{}
	err := k.doAuthenticatedRequest("POST", apiuri, params, &respmap)
	if err != nil {
		return false, err
	}
	//log.Println(respmap)
	return true, nil
}

func (k *Kraken) toOrder(orderinfo interface{}) Order {
	omap := orderinfo.(map[string]interface{})
	descmap := omap["descr"].(map[string]interface{})
	return Order{
		Amount:     ToFloat64(omap["vol"]),
		Price:      ToFloat64(descmap["price"]),
		DealAmount: ToFloat64(omap["vol_exec"]),
		AvgPrice:   ToFloat64(omap["price"]),
		Side:       AdaptTradeSide(descmap["type"].(string)),
		Status:     k.convertOrderStatus(omap["status"].(string)),
		Fee:        ToFloat64(omap["fee"]),
		OrderTime:  ToInt(omap["opentm"]),
	}
}

func (k *Kraken) GetOrderInfos(txids ...string) ([]Order, error) {
	params := url.Values{}
	params.Set("txid", strings.Join(txids, ","))

	var resultmap map[string]interface{}
	err := k.doAuthenticatedRequest("POST", "private/QueryOrders", params, &resultmap)
	if err != nil {
		return nil, err
	}
	//log.Println(resultmap)
	var ords []Order
	for txid, v := range resultmap {
		ord := k.toOrder(v)
		ord.OrderID2 = txid
		ords = append(ords, ord)
	}

	return ords, nil
}

func (k *Kraken) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	orders, err := k.GetOrderInfos(orderId)

	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, errors.New("not fund the order " + orderId)
	}

	ord := &orders[0]
	ord.Currency = currency
	return ord, nil
}

func (k *Kraken) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	var result struct {
		Open map[string]interface{} `json:"open"`
	}

	err := k.doAuthenticatedRequest("POST", "private/OpenOrders", url.Values{}, &result)
	if err != nil {
		return nil, err
	}

	var orders []Order

	for txid, v := range result.Open {
		ord := k.toOrder(v)
		ord.OrderID2 = txid
		ord.Currency = currency
		orders = append(orders, ord)
	}

	return orders, nil
}

func (k *Kraken) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("")
}

func (k *Kraken) GetAccount() (*Account, error) {
	params := url.Values{}
	apiuri := "private/Balance"

	var resustmap map[string]interface{}
	err := k.doAuthenticatedRequest("POST", apiuri, params, &resustmap)
	if err != nil {
		return nil, err
	}

	acc := new(Account)
	acc.Exchange = k.GetExchangeName()
	acc.SubAccounts = make(map[Currency]SubAccount)

	for key, v := range resustmap {
		currency := k.convertCurrency(key)
		amount := ToFloat64(v)
		//log.Println(symbol, amount)
		acc.SubAccounts[currency] = SubAccount{Currency: currency, Amount: amount, ForzenAmount: 0, LoanAmount: 0}

		if currency.Symbol == "XBT" { // adapt to btc
			acc.SubAccounts[BTC] = SubAccount{Currency: BTC, Amount: amount, ForzenAmount: 0, LoanAmount: 0}
		}
	}

	return acc, nil

}

//func (k *Kraken) GetTradeBalance() {
//	var resultmap map[string]interface{}
//	k.doAuthenticatedRequest("POST", "private/TradeBalance", url.Values{}, &resultmap)
//	log.Println(resultmap)
//}

func (k *Kraken) GetTicker(currency CurrencyPair) (*Ticker, error) {
	var resultmap map[string]interface{}
	err := k.doAuthenticatedRequest("GET", "public/Ticker?pair="+k.convertPair(currency).ToSymbol(""), url.Values{}, &resultmap)
	if err != nil {
		return nil, err
	}

	ticker := new(Ticker)
	ticker.Pair = currency
	for _, t := range resultmap {
		tickermap := t.(map[string]interface{})
		ticker.Last = ToFloat64(tickermap["c"].([]interface{})[0])
		ticker.Buy = ToFloat64(tickermap["b"].([]interface{})[0])
		ticker.Sell = ToFloat64(tickermap["a"].([]interface{})[0])
		ticker.Low = ToFloat64(tickermap["l"].([]interface{})[0])
		ticker.High = ToFloat64(tickermap["h"].([]interface{})[0])
		ticker.Vol = ToFloat64(tickermap["v"].([]interface{})[0])
	}

	return ticker, nil
}

func (k *Kraken) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	apiuri := fmt.Sprintf("public/Depth?pair=%s&count=%d", k.convertPair(currency).ToSymbol(""), size)
	var resultmap map[string]interface{}
	err := k.doAuthenticatedRequest("GET", apiuri, url.Values{}, &resultmap)
	if err != nil {
		return nil, err
	}

	//log.Println(respmap)
	dep := Depth{}
	dep.Pair = currency
	for _, d := range resultmap {
		depmap := d.(map[string]interface{})
		asksmap := depmap["asks"].([]interface{})
		bidsmap := depmap["bids"].([]interface{})
		for _, v := range asksmap {
			ask := v.([]interface{})
			dep.AskList = append(dep.AskList, DepthRecord{ToFloat64(ask[0]), ToFloat64(ask[1])})
		}
		for _, v := range bidsmap {
			bid := v.([]interface{})
			dep.BidList = append(dep.BidList, DepthRecord{ToFloat64(bid[0]), ToFloat64(bid[1])})
		}
		break
	}

	sort.Sort(sort.Reverse(dep.AskList)) //reverse

	return &dep, nil
}

func (k *Kraken) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("")
}

//非个人，整个交易所的交易记录
func (k *Kraken) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("")
}

func (k *Kraken) GetExchangeName() string {
	return KRAKEN
}

func (k *Kraken) buildParamsSigned(apiuri string, postForm *url.Values) string {
	postForm.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))
	urlPath := API_V0 + apiuri

	secretByte, _ := base64.StdEncoding.DecodeString(k.secretKey)
	encode := []byte(postForm.Get("nonce") + postForm.Encode())

	sha := sha256.New()
	sha.Write(encode)
	shaSum := sha.Sum(nil)

	pathSha := append([]byte(urlPath), shaSum...)

	mac := hmac.New(sha512.New, secretByte)
	mac.Write(pathSha)
	macSum := mac.Sum(nil)

	sign := base64.StdEncoding.EncodeToString(macSum)

	return sign
}

func (k *Kraken) doAuthenticatedRequest(method, apiuri string, params url.Values, ret interface{}) error {
	headers := map[string]string{}

	if "POST" == method {
		signature := k.buildParamsSigned(apiuri, &params)
		headers = map[string]string{
			"API-Key":  k.accessKey,
			"API-Sign": signature,
		}
	}

	resp, err := NewHttpRequest(k.httpClient, method, API_DOMAIN+apiuri, params.Encode(), headers)
	if err != nil {
		return err
	}
	//println(string(resp))
	var base BaseResponse
	base.Result = ret

	err = json.Unmarshal(resp, &base)
	if err != nil {
		return err
	}

	//println(string(resp))

	if len(base.Error) > 0 {
		return errors.New(base.Error[0])
	}

	return nil
}

func (k *Kraken) convertCurrency(currencySymbol string) Currency {
	if len(currencySymbol) >= 4 {
		currencySymbol = strings.Replace(currencySymbol, "X", "", 1)
		currencySymbol = strings.Replace(currencySymbol, "Z", "", 1)
	}
	return NewCurrency(currencySymbol, "")
}

func (k *Kraken) convertPair(pair CurrencyPair) CurrencyPair {
	if "BTC" == pair.CurrencyA.Symbol {
		return NewCurrencyPair(XBT, pair.CurrencyB)
	}

	if "BTC" == pair.CurrencyB.Symbol {
		return NewCurrencyPair(pair.CurrencyA, XBT)
	}

	return pair
}

func (k *Kraken) convertOrderStatus(status string) TradeStatus {
	switch status {
	case "open", "pending":
		return ORDER_UNFINISH
	case "canceled", "expired":
		return ORDER_CANCEL
	case "filled", "closed":
		return ORDER_FINISH
	case "partialfilled":
		return ORDER_PART_FINISH
	}
	return ORDER_UNFINISH
}
