package binance

import (
	//"encoding/json"
	//"errors"
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	//"io/ioutil"
	//"log"
	"net/http"
	//"net/url"
	"strconv"

)

const (
	EXCHANGE_NAME = "binance.com"

	API_BASE_URL = "http://binance.com/"
	API_V1       = API_BASE_URL + "api/v1/"

	TICKER_URI             = "/ticker/24hr?symbol=%s"
	//DEPTH_URI              = "depth.php?c=%s&mk_type=%s"
	//ACCOUNT_URI            = "getMyBalance.php"
	//TRADE_URI              = "trades.php?c=%s&mk_type=%s"
	//CANCEL_URI             = "cancelOrder.php"
	//ORDERS_INFO            = "getMyTradeList.php"
	//UNFINISHED_ORDERS_INFO = "getOrderList.php"

)

type Binance struct {
	accessKey,
	secretKey string
	httpClient *http.Client
}

func New(client *http.Client, api_key, secret_key string) *Binance {
	return &Binance{api_key, secret_key,client }
}

func (ba *Binance) GetExchangeName() string {
	return EXCHANGE_NAME
}

func (ba *Binance) GetTicker(currency CurrencyPair) (*Ticker, error) {
	tickerUri := API_V1 + fmt.Sprintf(TICKER_URI, currency.ToSymbol(""))
	bodyDataMap, err := HttpGet(ba.httpClient, tickerUri)

	if err != nil {
		return nil, err
	}
	fmt.Println("binance bodyDataMap:", currency, bodyDataMap)
	var tickerMap map[string]interface{} = bodyDataMap
	var ticker Ticker

	ticker.Date = uint64(tickerMap["closeTime"].(float64))
	ticker.Last,_ = strconv.ParseFloat(tickerMap["lastPrice"].(string), 10)
	ticker.Buy,_ = strconv.ParseFloat(tickerMap["bidPrice"].(string), 10)
	ticker.Sell,_ = strconv.ParseFloat(tickerMap["askPrice"].(string), 10)
	ticker.Low,_ = strconv.ParseFloat(tickerMap["lowPrice"].(string), 10)
	ticker.High,_ = strconv.ParseFloat(tickerMap["highPrice"].(string), 10)
	ticker.Vol,_ = strconv.ParseFloat(tickerMap["volume"].(string), 10)
	return &ticker, nil
}