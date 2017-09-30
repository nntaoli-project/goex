package wex

import (
	. "comparewebs/goex"
	"net/http"
)

const (
	EXCHANGE_NAME = "wex.nz"

	API_BASE_URL = "https://wex.nz/api/3/"
	TICKER_URI   = "ticker/"
)

type Wex struct {
	accessKey,
	secretKey string
	httpClient *http.Client
}

func New(client *http.Client, api_key, secret_key string) *Wex {
	return &Wex{api_key, secret_key, client}
}

func (wex *Wex) GetExchangeName() string {
	return EXCHANGE_NAME
}

func (wex *Wex) GetTicker(currency CurrencyPair) (*Ticker, error) {
	tickerUri := API_BASE_URL + TICKER_URI + currency.String()

	resp, err := HttpGet(wex.httpClient, tickerUri)
	//fmt.Println(err,resp)
	if err != nil {
		return nil, err
	}
	tickerMap := resp[currency.String()].(map[string]interface{})
	ticker := new(Ticker)

	ticker.Last, _ = tickerMap["last"].(float64)
	ticker.Date = uint64(tickerMap["updated"].(float64))
	ticker.Buy, _ = tickerMap["buy"].(float64)
	ticker.Sell, _ = tickerMap["sell"].(float64)
	ticker.Low, _ = tickerMap["low"].(float64)
	ticker.High, _ = tickerMap["high"].(float64)
	ticker.Vol, _ = tickerMap["vol"].(float64)

	return ticker, nil
}
