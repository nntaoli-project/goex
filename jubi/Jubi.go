package jubi

import (
	//	"encoding/json"

	//"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	. "github.com/nntaoli-project/GoEx"
)

const (
	EXCHANGE_NAME = "jubi.com"

	API_BASE_URL = "https://www.jubi.com"
	API_V1       = API_BASE_URL + "/api/v1/"

	TICKER_URI     = "ticker?coin=%s"
	ALL_TICKER_URI = "allticker/"
	DEPTH_URI      = "depth?coin=%s"
	TRADES_URI     = "orders?coin=%s"

	//	ACCOUNT_URI = "getMyBalance.php"

//	CANCEL_URI  = "cancelOrder.php"
//	ORDERS_INFO = "getOrderList.php"
)

type Jubi struct {
	accessKey,
	secretKey,
	accountId string
	httpClient *http.Client
}

func New(client *http.Client, accessKey, secretKey, accountId string) *Jubi {
	return &Jubi{accessKey, secretKey, accountId, client}
}

func (jubi *Jubi) GetExchangeName() string {
	return EXCHANGE_NAME
}

var AllTickerMap map[string]interface{}
var timeStampAllTicker = time.Now().Unix()

func (jubi *Jubi) GetAllTicker() error {
	var err error
	tickerUri := API_V1 + ALL_TICKER_URI
	timeStampAllTicker = time.Now().Unix()

	AllTickerMap, err = HttpGet(jubi.httpClient, tickerUri)
	//	log.Println("Jubi bodyDataMap:", AllTickerMap)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func (jubi *Jubi) GetTickerInBuf(currency CurrencyPair) (*Ticker, error) {
	cur := strings.ToLower(currency.CurrencyA.String())
	if cur == "nil" {
		log.Println("Unsupport The CurrencyPair")
		return nil, errors.New("Unsupport The CurrencyPair")
	}

	if AllTickerMap == nil {
		return nil, errors.New("Ticker buffer is nil")
	}
	var ticker Ticker

	switch AllTickerMap[cur].(type) {
	case map[string]interface{}:
		tickerMap := AllTickerMap[cur].(map[string]interface{})
		ticker.Date = uint64(timeStampAllTicker)
		ticker.Last, _ = tickerMap["last"].(float64)
		ticker.Buy, _ = tickerMap["buy"].(float64)
		ticker.Sell, _ = tickerMap["sell"].(float64)
		ticker.Low, _ = tickerMap["low"].(float64)
		ticker.High, _ = tickerMap["high"].(float64)
		ticker.Vol, _ = tickerMap["vol"].(float64)
		return &ticker, nil
	default:
		return nil, errors.New(fmt.Sprintf("Type Convert Error ? \n %s", AllTickerMap))
	}

}
func (jubi *Jubi) GetTicker(currency CurrencyPair) (*Ticker, error) {
	cur := strings.ToLower(currency.CurrencyA.String())
	//	money := currency.CurrencyB.String()
	if cur == "nil" {
		log.Println("Unsupport The CurrencyPair")
		return nil, errors.New("Unsupport The CurrencyPair")
	}
	tickerUri := API_V1 + fmt.Sprintf(TICKER_URI, cur)
	//fmt.Println("tickerUrl:",tickerUri)
	timestamp := time.Now().Unix()
	bodyDataMap, err := HttpGet(jubi.httpClient, tickerUri)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var tickerMap map[string]interface{} = bodyDataMap
	var ticker Ticker

	//fmt.Println(bodyDataMap)
	ticker.Date = uint64(timestamp)
	ticker.Last, _ = strconv.ParseFloat(tickerMap["last"].(string), 64)

	ticker.Buy, _ = strconv.ParseFloat(tickerMap["buy"].(string), 64)
	ticker.Sell, _ = strconv.ParseFloat(tickerMap["sell"].(string), 64)
	ticker.Low, _ = strconv.ParseFloat(tickerMap["low"].(string), 64)
	ticker.High, _ = strconv.ParseFloat(tickerMap["high"].(string), 64)
	ticker.Vol, _ = tickerMap["vol"].(float64)

	return &ticker, nil
}

func (jubi *Jubi) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	var depthUri string

	cur := strings.ToLower(currency.CurrencyA.String())
	if cur == "nil" {
		log.Println("Unsupport The CurrencyPair")
		return nil, errors.New("Unsupport The CurrencyPair")
	}
	depthUri = API_V1 + fmt.Sprintf(DEPTH_URI, cur)

	bodyDataMap, err := HttpGet(jubi.httpClient, depthUri)

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

/**
 * 获取全站最近的交易记录
 */
/*
func (jubi *Jubi) GetTrades(currency CurrencyPair, since int64) ([]Trade2, error) {
	cur := currency.CurrencyA.String()
	if cur == "nil" {
		log.Println("Unsupport The CurrencyPair")
		return nil, errors.New("Unsupport The CurrencyPair")
	}
	s := strings.Split(cur, "_")
	tradeUrl := API_V1 + fmt.Sprintf(TRADES_URI, cur)

	bodyDataMap, err := HttpGet2(jubi.httpClient, tradeUrl)
	fmt.Println("bodyDataMap:", string(bodyDataMap))

	var trades []Trade2

	err = json.Unmarshal(bodyDataMap, &trades)
	if err != nil {
		return nil, err
	}
	fmt.Println("trades:", trades)

	return trades, nil

	panic("unimplements")

}
*/
