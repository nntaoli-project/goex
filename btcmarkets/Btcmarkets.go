package btcmarkets

import (
	"errors"
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"log"
	"net/http"
	"time"
)

const (
	EXCHANGE_NAME = "btcmarkets.net"

	API_BASE_URL = "https://api.btcmarkets.net/"

	TICKER_URI = "market/%s/%s/tick"
)

type Btcmarkets struct {
	accessKey,
	secretKey string
	httpClient *http.Client
}

func New(client *http.Client, accessKey, secretKey string) *Btcmarkets {
	return &Btcmarkets{accessKey, secretKey, client}
}

func (btcm *Btcmarkets) GetExchangeName() string {
	return EXCHANGE_NAME
}

func (btcm *Btcmarkets) GetTicker(currency CurrencyPair) (*Ticker, error) {
	tickerUri := fmt.Sprintf(API_BASE_URL+TICKER_URI, currency.CurrencyA.String(), currency.CurrencyB.String())
	//log.Println("tickerUrl:", tickerUri)
	bodyDataMap, err := HttpGet(btcm.httpClient, tickerUri)
	//log.Println("Btcmarkets bodyDataMap:", tickerUri, bodyDataMap)

	timestamp := time.Now().Unix()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if result, isok := bodyDataMap["success"].(bool); isok == true && result != true {
		log.Println("bodyDataMap[\"success\"]", isok, result)
		return nil, errors.New("err")
	}

	var tickerMap map[string]interface{} = bodyDataMap
	var ticker Ticker

	//fmt.Println(bodyDataMap)
	ticker.Date = uint64(timestamp)
	ticker.Last, _ = tickerMap["lastPrice"].(float64)

	ticker.Buy, _ = tickerMap["bestBid"].(float64)
	ticker.Sell, _ = tickerMap["bestAsk"].(float64)
	ticker.Vol, _ = tickerMap["volume24h"].(float64)
	log.Println("Btcmarkets", currency, "ticker:", ticker)
	return &ticker, nil
}
