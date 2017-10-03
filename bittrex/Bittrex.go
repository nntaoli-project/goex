package bittrex

import (
	"errors"
	. "github.com/nntaoli-project/GoEx"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	API_BASE                   = "https://bittrex.com/api/" // Bittrex API endpoint
	API_VERSION                = "v1.1/"                    // Bittrex API version
	DEFAULT_HTTPCLIENT_TIMEOUT = 30                         // HTTP client timeout
)

type Currencys struct {
	Currency        string  `json:"Currency"`
	CurrencyLong    string  `json:"CurrencyLong"`
	MinConfirmation int     `json:"MinConfirmation"`
	TxFee           float64 `json:"TxFee"`
	IsActive        bool    `json:"IsActive"`
	CoinType        string  `json:"CoinType"`
	BaseAddress     string  `json:"BaseAddress"`
	Notice          string  `json:"Notice"`
}

// bittrex represent a bittrex client
type Bittrex struct {
	httpClient *http.Client
	accessKey,
	secretKey string
}

func New(httpClient *http.Client, accessKey, secretKey string) *Bittrex {
	return &Bittrex{httpClient, accessKey, secretKey}
}

// GetCurrencies is used to get all supported currencies at Bittrex along with other meta data.
func (b *Bittrex) GetCurrencies() ([]Currencys, error) {
	curl := API_BASE + API_VERSION + "/public/getcurrencies"
	//log.Println(curl)
	r, err := HttpGet(b.httpClient, curl)
	if err != nil {
		return nil, err
	}
	if r["success"].(bool) != true {
		return nil, errors.New(r["success"].(string))
	}
	currency := r["result"].([]interface{})
	currencys := make([]Currencys, 0)
	for _, vv := range currency {
		v := vv.(map[string]interface{})
		cur := Currencys{}
		cur.BaseAddress, _ = v["BaseAddress"].(string)
		cur.CoinType, _ = v["CoinType"].(string)
		cur.Currency, _ = v["Currency"].(string)
		cur.CurrencyLong, _ = v["CurrencyLong"].(string)
		cur.IsActive, _ = v["IsActive"].(bool)
		cur.Notice, _ = v["Notice"].(string)
		cur.MinConfirmation, _ = v["MinConfirmation"].(int)
		cur.TxFee, _ = v["TxFee"].(float64)
		currencys = append(currencys, cur)
	}
	return currencys, nil
}

// GetTicker is used to get the current ticker values for a market.
func (b *Bittrex) GetTicker(currency CurrencyPair) (*Ticker, error) {
	r, err := HttpGet(b.httpClient, API_BASE+API_VERSION+"public/getmarketsummary?market="+strings.ToUpper(currency.ToSymbol2("-")))
	if err != nil {
		return nil, err
	}
	if r["success"].(bool) != true {
		return nil, errors.New(r["success"].(string))
	}

	resp := r["result"].([]interface{})
	tickerMap := resp[0].(map[string]interface{})
	ticker := new(Ticker)

	ticker.Date = uint64(time.Now().Unix())

	ticker.Buy, _ = tickerMap["Bid"].(float64)
	ticker.Sell, _ = tickerMap["Ask"].(float64)
	ticker.Last, _ = tickerMap["Last"].(float64)
	ticker.High, _ = tickerMap["Bid"].(float64)
	ticker.Low, _ = tickerMap["Low"].(float64)
	ticker.Vol, _ = tickerMap["Volume"].(float64)

	log.Println(ticker)

	return ticker, nil

}
