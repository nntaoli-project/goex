package okcoin

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
	"github.com/thrasher-/gocryptotrader/exchanges/ticker"
)

const (
	EXCHANGE_NAME_EX = "okex.com"

	API_URL = "https://www.okex.com"
	API_URI_PREFIX = "/v2/"
	API_URL_V2 = API_URL+API_URI_PREFIX
	TICKER_URL = "markets/%s/ticker"
	DEPTH_URL = "markets/%s/depth"
	KLINE_URL = "markets/%s/kline?%s=&%d="

	USER_INFO_URL = ""
	GET_ORDER_API = ""
	DELETE_ORDER_API = ""
	PLACE_ORDER_API = ""

)


type OKExC struct {
	apiKey,
	apiSecretKey string
	client *http.Client
}

func NewOKExC(client *http.Client, api_key, secret_key string) *OKExC {
	ok := new(OKExC)
	ok.apiKey = api_key
	ok.apiSecretKey = secret_key
	ok.client = client
	return ok
}
func (ctx *OKExC) buildPostForm(postForm *url.Values) error {
	panic("unimplements")
}

func (ctx *OKExC) placeOrder(side, amount, price string, currency CurrencyPair) (*Order, error) {
	panic("unimplements")
}

func (ctx *OKExC) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("unimplements")
}

func (ctx *OKExC) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("unimplements")
}

func (ctx *OKExC) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("unimplements")
}

func (ctx *OKExC) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("unimplements")
}

func (ctx *OKExC) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("unimplements")
}

func (ctx *OKExC) getOrders(orderId string, currency CurrencyPair) ([]Order, error) {
	panic("unimplements")
}

func (ctx *OKExC) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("unimplements")
}

func (ctx *OKExC) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("unimplements")
}

func (ctx *OKExC) GetAccount() (*Account, error) {
	panic("unimplements")
}

func (ctx *OKExC) GetTicker(currency CurrencyPair) (*Ticker, error) {
	var tickerMap map[string]interface{}
	var ticker Ticker

	url := API_URL_V2 + fmt.Sprintf(TICKER_URL,strings.ToLower(currency.ToSymbol("_")))
	bodyDataMap, err := HttpGet(ctx.client, url)
	if err != nil {
		return nil, err
	}
	if(bodyDataMap["code"].(uint64) != 0){
		return nil, errors.New("response code is not 0")
	}

	tickerMap = bodyDataMap["data"].(map[string]interface{})
	ticker.Date = tickerMap["createdDate"].(uint64)
	ticker.Last, _ = strconv.ParseFloat(tickerMap["last"].(string), 64)
	ticker.Buy, _ = strconv.ParseFloat(tickerMap["buy"].(string), 64)
	ticker.Sell, _ = strconv.ParseFloat(tickerMap["sell"].(string), 64)
	ticker.Low, _ = strconv.ParseFloat(tickerMap["low"].(string), 64)
	ticker.High, _ = strconv.ParseFloat(tickerMap["high"].(string), 64)
	ticker.Vol, _ = strconv.ParseFloat(tickerMap["volume"].(string), 64)

	return &ticker, nil
}

func (ctx *OKExC) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	var depth Depth
	url := API_URL_V2 + fmt.Sprintf(DEPTH_URL,strings.ToLower(currency.ToSymbol("_"))) + "?size=" + strconv.Itoa(size)
	respDataMap, err := HttpGet(ctx.client, url)
	if err != nil {
		return nil, err
	}
	if(respDataMap["code"].(uint64) != 0){
		return nil, errors.New("response code is not 0")
	}
	bodyDataMap := respDataMap["data"].(map[string]interface{})

	for _, v := range bodyDataMap["asks"].([]interface{}) {
		var dr DepthRecord
		for i, vv := range v.([]interface{}) {
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

func (ctx *OKExC) GetExchangeName() string {
	return EXCHANGE_NAME_EX
}

func (ctx *OKExC) GetKlineRecords(currency CurrencyPair, period string, size, since int) ([]Kline, error) {
	klineUrl := API_URL_V2 + fmt.Sprintf(KLINE_URL,strings.ToLower(currency.ToSymbol("_")), period, since)
	resp, err := http.Get(klineUrl)
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

	var klineRecords []Kline

	for _, record := range klines {
		r := Kline{}
		for i, e := range record {
			switch i {
			case 0:
				r.Timestamp = int64(e.(float64)) / 1000 //to unix timestramp
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
		klineRecords = append(klineRecords, r)
	}

	return klineRecords, nil
}

func (ctx *OKExC) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("unimplements")
}

func (ok *OKExC) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("unimplements")
}
