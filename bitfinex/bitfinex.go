package bitfinex

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Bitfinex struct {
	httpClient *http.Client
	accessKey,
	secretKey string
}

const (
	BASE_URL = "https://api.bitfinex.com/v1"
)

func New(client *http.Client, accessKey, secretKey string) *Bitfinex {
	return &Bitfinex{client, accessKey, secretKey}
}

func (bfx *Bitfinex) GetExchangeName() string {
	return BITFINEX
}

func (bfx *Bitfinex) GetTicker(currencyPair CurrencyPair) (*Ticker, error) {
	//pubticker
	currencyPair = bfx.adaptCurrencyPair(currencyPair)

	apiUrl := fmt.Sprintf("%s/pubticker/%s", BASE_URL, strings.ToLower(currencyPair.ToSymbol("")))
	resp, err := HttpGet(bfx.httpClient, apiUrl)
	if err != nil {
		return nil, err
	}

	if resp["error"] != nil {
		return nil, errors.New(resp["error"].(string))
	}

	//fmt.Println(resp)
	ticker := new(Ticker)
	ticker.Pair = currencyPair
	ticker.Last = ToFloat64(resp["last_price"])
	ticker.Vol = ToFloat64(resp["volume"])
	ticker.High = ToFloat64(resp["high"])
	ticker.Low = ToFloat64(resp["low"])
	ticker.Sell = ToFloat64(resp["ask"])
	ticker.Buy = ToFloat64(resp["bid"])
	ticker.Date = uint64(bfx.adaptTimestamp(resp["timestamp"].(string)))
	return ticker, nil
}

func (bfx *Bitfinex) GetDepth(size int, currencyPair CurrencyPair) (*Depth, error) {
	apiUrl := fmt.Sprintf("%s/book/%s?limit_bids=%d&limit_asks=%d", BASE_URL, bfx.currencyPairToSymbol(currencyPair), size, size)
	resp, err := HttpGet(bfx.httpClient, apiUrl)
	if err != nil {
		return nil, err
	}
	println("resp:", resp)
	bids := resp["bids"].([]interface{})
	asks := resp["asks"].([]interface{})

	depth := new(Depth)

	for _, bid := range bids {
		_bid := bid.(map[string]interface{})
		amount := ToFloat64(_bid["amount"])
		price := ToFloat64(_bid["price"])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.BidList = append(depth.BidList, dr)
	}

	for _, ask := range asks {
		_ask := ask.(map[string]interface{})
		amount := ToFloat64(_ask["amount"])
		price := ToFloat64(_ask["price"])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.AskList = append(depth.AskList, dr)
	}

	return depth, nil
}

func (bfx *Bitfinex) GetKlineRecords(currencyPair CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录

func (bfx *Bitfinex) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}

func (bfx *Bitfinex) GetWalletBalances() (map[string]*Account, error) {
	var respmap []interface{}
	err := bfx.doAuthenticatedRequest("GET", "balances", map[string]interface{}{}, &respmap)
	if err != nil {
		return nil, err
	}
	//log.Println(respmap)

	walletmap := make(map[string]*Account, 1)

	for _, v := range respmap {
		subacc := v.(map[string]interface{})
		typeStr := subacc["type"].(string)

		currency := NewCurrency(subacc["currency"].(string), "")

		if currency == UNKNOWN {
			continue
		}

		//typeS := subacc["type"].(string)
		amount := ToFloat64(subacc["amount"])
		available := ToFloat64(subacc["available"])

		account := walletmap[typeStr]
		if account == nil {
			account = new(Account)
			account.SubAccounts = make(map[Currency]SubAccount, 6)
		}

		account.NetAsset = amount
		account.Asset = amount
		account.SubAccounts[currency] = SubAccount{
			Currency:     currency,
			Amount:       available,
			ForzenAmount: amount - available,
			LoanAmount:   0}

		walletmap[typeStr] = account
	}

	return walletmap, nil
}

/*defalut only return exchange wallet balance*/
func (bfx *Bitfinex) GetAccount() (*Account, error) {
	wallets, err := bfx.GetWalletBalances()
	if err != nil {
		return nil, err
	}
	return wallets["exchange"], nil
}

func (bfx *Bitfinex) placeOrder(orderType, side, amount, price string, pair CurrencyPair) (*Order, error) {
	path := "order/new"
	params := map[string]interface{}{
		"symbol":   bfx.currencyPairToSymbol(pair),
		"amount":   amount,
		"price":    price,
		"side":     side,
		"type":     orderType,
		"exchange": "bitfinex"}

	var respmap map[string]interface{}
	err := bfx.doAuthenticatedRequest("POST", path, params, &respmap)
	if err != nil {
		return nil, err
	}

	order := new(Order)
	order.Currency = pair
	order.OrderID = ToInt(respmap["id"])
	order.OrderID2 = fmt.Sprint(ToInt(respmap["id"]))
	order.Amount = ToFloat64(amount)
	order.Price = ToFloat64(price)
	order.AvgPrice = ToFloat64(respmap["avg_execution_price"])
	order.DealAmount = ToFloat64(respmap["executed_amount"])
	order.Status = ORDER_UNFINISH

	switch side {
	case "buy":
		if orderType == "limit" || orderType == "exchange limit" {
			order.Side = BUY
		} else {
			order.Side = BUY_MARKET
		}
	case "sell":
		if orderType == "limit" || orderType == "exchange limit" {
			order.Side = SELL
		} else {
			order.Side = SELL_MARKET
		}

	}
	return order, nil
}

func (bfx *Bitfinex) LimitBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return bfx.placeOrder("exchange limit", "buy", amount, price, currencyPair)
}

func (bfx *Bitfinex) LimitSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return bfx.placeOrder("exchange limit", "sell", amount, price, currencyPair)
}

func (bfx *Bitfinex) MarketBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return bfx.placeOrder("exchange market", "buy", amount, price, currencyPair)
}

func (bfx *Bitfinex) MarketSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return bfx.placeOrder("exchange market", "sell", amount, price, currencyPair)
}

func (bfx *Bitfinex) CancelOrder(orderId string, currencyPair CurrencyPair) (bool, error) {
	var respmap map[string]interface{}
	path := "order/cancel"
	err := bfx.doAuthenticatedRequest("POST", path, map[string]interface{}{"order_id": ToInt(orderId)}, &respmap)
	if err != nil {
		return false, err
	}
	return respmap["is_cancelled"].(bool), nil
}

func (bfx *Bitfinex) toOrder(respmap map[string]interface{}) *Order {
	order := new(Order)
	order.Currency = bfx.symbolToCurrencyPair(respmap["symbol"].(string))
	order.OrderID = ToInt(respmap["id"])
	order.OrderID2 = fmt.Sprint(ToInt(respmap["id"]))
	order.Amount = ToFloat64(respmap["original_amount"])
	order.Price = ToFloat64(respmap["price"])
	order.DealAmount = ToFloat64(respmap["executed_amount"])
	order.AvgPrice = ToFloat64(respmap["avg_execution_price"])
	order.OrderTime = bfx.adaptTimestamp(respmap["timestamp"].(string))

	if order.DealAmount == order.Amount {
		order.Status = ORDER_FINISH
	} else if order.DealAmount > 0 {
		order.Status = ORDER_PART_FINISH
	}

	side := respmap["side"].(string)
	if side == "sell" {
		order.Side = SELL
	} else if side == "buy" {
		order.Side = BUY
	}

	if respmap["is_cancelled"].(bool) {
		order.Status = ORDER_CANCEL
	}
	return order
}

func (bfx *Bitfinex) GetOneOrder(orderId string, currencyPair CurrencyPair) (*Order, error) {
	var respmap map[string]interface{}
	path := "order/status"
	err := bfx.doAuthenticatedRequest("POST", path, map[string]interface{}{"order_id": ToInt(orderId)}, &respmap)
	if err != nil {
		return nil, err
	}
	return bfx.toOrder(respmap), nil
}

func (bfx *Bitfinex) GetUnfinishOrders(currencyPair CurrencyPair) ([]Order, error) {
	var ordersmap []interface{}
	err := bfx.doAuthenticatedRequest("POST", "orders", map[string]interface{}{}, &ordersmap)
	if err != nil {
		return nil, err
	}

	var orders []Order
	for _, v := range ordersmap {
		ordermap := v.(map[string]interface{})
		orders = append(orders, *bfx.toOrder(ordermap))
	}
	return orders, nil
}

func (bfx *Bitfinex) GetOrderHistorys(currencyPair CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implement")
}

func (bfx *Bitfinex) doAuthenticatedRequest(method, path string, payload map[string]interface{}, ret interface{}) error {
	nonce := time.Now().UnixNano()
	payload["request"] = "/v1/" + path
	payload["nonce"] = fmt.Sprintf("%d.2", nonce)

	//for k, v := range params {
	//	payload[k] = v[0]
	//}

	p, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	//println(string(p))
	encoded := base64.StdEncoding.EncodeToString(p)
	sign, _ := GetParamHmacSha384Sign(bfx.secretKey, encoded)
	//log.Println(BASE_URL + "/" + path)

	resp, err := NewHttpRequest(bfx.httpClient, method, BASE_URL+"/"+path, "", map[string]string{
		"Content-Type":    "application/json",
		"Accept":          "application/json",
		"X-BFX-APIKEY":    bfx.accessKey,
		"X-BFX-PAYLOAD":   encoded,
		"X-BFX-SIGNATURE": sign})

	if err != nil {
		return err
	}
	//print(string(resp))
	err = json.Unmarshal(resp, ret)
	return err
}

func (bfx *Bitfinex) currencyPairToSymbol(currencyPair CurrencyPair) string {
	return strings.ToUpper(currencyPair.ToSymbol(""))
}

func (bfx *Bitfinex) symbolToCurrencyPair(symbol string) CurrencyPair {
	currencyA := strings.ToUpper(symbol[0:3])
	currencyB := strings.ToUpper(symbol[3:])
	return NewCurrencyPair(NewCurrency(currencyA, ""), NewCurrency(currencyB, ""))
}

func (bfx *Bitfinex) adaptTimestamp(timestamp string) int {
	times := strings.Split(timestamp, ".")
	intTime, _ := strconv.Atoi(times[0])
	return intTime
}

func (bfx *Bitfinex) adaptCurrencyPair(pair CurrencyPair) CurrencyPair {
	var currencyA Currency
	var currencyB Currency

	DASH := NewCurrency("DASH", "")
	DSH := NewCurrency("DSH", "")
	QTM := NewCurrency("QTM", "")
	IOTA := NewCurrency("IOTA", "")
	IOT := NewCurrency("IOT", "")

	if pair.CurrencyA == DASH {
		currencyA = DSH
	} else if pair.CurrencyA == QTUM {
		currencyA = QTM
	} else if pair.CurrencyA == IOTA {
		currencyA = IOT
	} else {
		currencyA = pair.CurrencyA
	}

	if pair.CurrencyB == USDT {
		currencyB = USD
	} else {
		currencyB = pair.CurrencyB
	}

	return NewCurrencyPair(currencyA, currencyB)
}
