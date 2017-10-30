package kraken

import (
	"errors"
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	EXCHANGE_NAME = "kraken.com"

	API_BASE_URL = "https://api.kraken.com"
	API_V0       = "/0/"
	PUBLIC       = "public/"
	PRIVATE      = "private/"
	TICKER_URI   = "Ticker?pair=%s"
	ACCOUNT_URI  = "Balance"
)

type Kraken struct {
	accessKey,
	secretKey string
	httpClient *http.Client
}

func New(client *http.Client, accessKey, secretKey string) *Kraken {
	return &Kraken{accessKey, secretKey, client}
}

func (kraken *Kraken) GetExchangeName() string {
	return EXCHANGE_NAME
}

func (kraken *Kraken) GetTicker(currency CurrencyPair) (*Ticker, error) {
	curr := currency.ToSymbol("")
	tickerUri := fmt.Sprintf(API_BASE_URL+API_V0+PUBLIC+TICKER_URI, curr)
	bodyDataMap, err := HttpGet(kraken.httpClient, tickerUri)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if result, _ := bodyDataMap["error"].([]interface{}); len(result) != 0 {
		log.Println("bodyDataMap[\"error\"]", result)
		return nil, errors.New("err")
	}
	switch currency.CurrencyA.Symbol {
	case "XBT", "XRP", "ETC", "XMR":
		curr = "X" + currency.ToSymbol("Z")
	}

	curr = "X" + currency.ToSymbol("Z")
	tickerMap := bodyDataMap["result"].(map[string]interface{})
	var ticker Ticker
	tickers, isOk := tickerMap[curr].(map[string]interface{})
	if isOk != true {
		return nil, errors.New("err")
	}

	timestamp := time.Now().Unix()
	ticker.Date = uint64(timestamp)

	t := tickers["c"].([]interface{})
	ticker.Last, _ = strconv.ParseFloat(t[0].(string), 64)
	t = tickers["b"].([]interface{})
	ticker.Buy, _ = strconv.ParseFloat(t[0].(string), 64)
	t = tickers["a"].([]interface{})
	ticker.Sell, _ = strconv.ParseFloat(t[0].(string), 64)
	t = tickers["l"].([]interface{})
	ticker.Low, _ = strconv.ParseFloat(t[0].(string), 64)
	t = tickers["h"].([]interface{})
	ticker.High, _ = strconv.ParseFloat(t[0].(string), 64)
	t = tickers["v"].([]interface{})
	ticker.Vol, _ = strconv.ParseFloat(t[0].(string), 64)
	log.Println("Kraken", currency, "ticker:", ticker)

	return &ticker, nil
}

func (kraken *Kraken) GetDepth(size int, currencyPair CurrencyPair) (*Depth, error) {
	panic("not implement")

}

func (kraken *Kraken) buildParamsSigned(method string, postForm *url.Values) string {
	//	postForm.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))
	//	urlPath := API_V0+PRIVATE+method
	// Create signature
	//secret, _ := base64.StdEncoding.DecodeString(api.secret)
	//signature := createSignature(urlPath, values, secret)

	return nil
}
func (kraken *Kraken) queryPrivate(method string, values url.Values, typ interface{}) (interface{}, error) {

}
func (kraken *Kraken) placeOrder(amount, price string, pair CurrencyPair, orderType, orderSide string) (*Order, error) {
	panic("not implement")

}

func (kraken *Kraken) GetAccount() (*Account, error) {
	panic("not implement")

}

func (kraken *Kraken) LimitBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	panic("not implement")
}

func (kraken *Kraken) LimitSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	panic("not implement")
}

func (kraken *Kraken) MarketBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	panic("not implement")
}

func (kraken *Kraken) MarketSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	panic("not implement")
}

func (kraken *Kraken) CancelOrder(orderId string, currencyPair CurrencyPair) (bool, error) {
	panic("not implement")

}

func (kraken *Kraken) GetOneOrder(orderId string, currencyPair CurrencyPair) (*Order, error) {
	panic("not implement")

}

func (kraken *Kraken) GetUnfinishOrders(currencyPair CurrencyPair) ([]Order, error) {
	panic("not implement")
}
