// Bter
package bter

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	. "github.com/nntaoli-project/GoEx"
	"strings"
)

const (
	EXCHANGE_NAME   = "bter.com"
	API_BASE_URL    = "http://data.bter.com/"
	API_V1          = API_BASE_URL + "api2/1/"
	TICKER_URI      = "ticker/%s/"
	ALL_TICKERS_URI = "tickers/"
)

type Bter struct {
	accessKey,
	secretKey string
	httpClient *http.Client
}

func New(client *http.Client, accessKey, secretKey string) *Bter {
	return &Bter{accessKey, secretKey, client}
}

func (bter *Bter) GetExchangeName() string {
	return EXCHANGE_NAME
}

var AllTickerMap map[string]interface{}
var timeStampAllTicker = time.Now().Unix()

func (bter *Bter) GetAllTicker() error {
	var err error
	tickerUri := API_V1 + ALL_TICKERS_URI
	timeStampAllTicker = time.Now().Unix()

	AllTickerMap, err = HttpGet(bter.httpClient, tickerUri)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func (bter *Bter) GetTickerInBuf(currency CurrencyPair) (*Ticker, error) {
	cur := strings.ToLower(currency.ToSymbol("_"))
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
		ticker.Buy, _ = tickerMap["highestBid"].(float64)
		ticker.Sell, _ = tickerMap["lowestAsk"].(float64)
		ticker.Low, _ = tickerMap["low24hr"].(float64)
		ticker.High, _ = tickerMap["high24hr"].(float64)
		ticker.Vol, _ = tickerMap["baseVolume"].(float64)
		return &ticker, nil
	default:
		return nil, errors.New(fmt.Sprintf("Type Convert Error ? \n %s", AllTickerMap))
	}
}

func (bter *Bter) GetTicker(currency CurrencyPair) (*Ticker, error) {
	cur := strings.ToLower(currency.ToSymbol("_"))
	if cur == "nil" {
		log.Println("Unsupport The CurrencyPair")
		return nil, errors.New("Unsupport The CurrencyPair")
	}
	tickerUri := API_V1 + fmt.Sprintf(TICKER_URI, cur)

	timestamp := time.Now().Unix()
	bodyDataMap, err := HttpGet(bter.httpClient, tickerUri)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var tickerMap map[string]interface{} = bodyDataMap
	var ticker Ticker
	fmt.Println("bter tickerMap:", bodyDataMap)
	ticker.Date = uint64(timestamp)
	ticker.Last, _ = tickerMap["last"].(float64)

	ticker.Buy, _ = tickerMap["highestBid"].(float64)
	ticker.Sell, _ = tickerMap["lowestAsk"].(float64)
	ticker.Low, _ = tickerMap["low24hr"].(float64)
	ticker.High, _ = tickerMap["high24hr"].(float64)
	ticker.Vol = tickerMap["baseVolume"].(float64)

	return &ticker, nil
}
