package kraken

import (
	"errors"
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	EXCHANGE_NAME = "kraken.com"

	API_BASE_URL = "https://api.kraken.com/0/public/Ticker"

	TICKER_URI = "?pair=%s"
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
	tickerUri := fmt.Sprintf(API_BASE_URL+TICKER_URI, curr)
	//log.Println("Kraken tickerUri:",tickerUri)
	bodyDataMap, err := HttpGet(kraken.httpClient, tickerUri)
	//log.Println("Kraken bodyDataMap:",bodyDataMap)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	//result, isok := bodyDataMap["error"].([]interface{})
	//log.Println("bodyDataMap[\"error\"]",result, isok)
	if result, _ := bodyDataMap["error"].([]interface{}); len(result) != 0 {
		log.Println("bodyDataMap[\"error\"]", result)
		return nil, errors.New("err")
	}
	curr = "X" + currency.ToSymbol("Z")
	tickerMap := bodyDataMap["result"].(map[string]interface{})
	//log.Println("Kraken curr:",curr)
	//log.Println("Kraken tickerMap:",tickerMap)
	var ticker Ticker
	tickers := tickerMap[curr].(map[string]interface{})
	//log.Println("Kraken tickers:",tickers)

	//fmt.Println(bodyDataMap)
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
