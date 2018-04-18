package cryptopia

import (
	"errors"
	. "github.com/nntaoli-project/GoEx"
	//"log"
	"net/http"
	"time"
)

const (
	EXCHANGE_NAME = "cryptopia.co.nz"

	API_BASE_URL = "https://www.cryptopia.co.nz/api/"

	TICKERS_URI = "GetMarkets"
	TICKER_URI  = "GetMarket/"
)

type Cryptopia struct {
	accessKey,
	secretKey string
	httpClient *http.Client
}

func New(client *http.Client, accessKey, secretKey string) *Cryptopia {
	return &Cryptopia{accessKey, secretKey, client}
}

func (cta *Cryptopia) GetExchangeName() string {
	return EXCHANGE_NAME
}

func (cta *Cryptopia) GetTickers(currency CurrencyPair) (*Ticker, error) {
	return cta.GetTicker(currency)

	//tickerUri := API_BASE_URL + TICKERS_URI
	////log.Println("tickerUrl:", tickerUri)
	//bodyDataMap, err := HttpGet(cta.httpClient, tickerUri)
	////log.Println("Cryptopia bodyDataMap:", tickerUri, bodyDataMap)
	//
	//if err != nil {
	//	log.Println(err)
	//	return nil, err
	//}
	//
	//if result, isok := bodyDataMap["success"].(bool); isok == true && result != true {
	//	log.Println("bodyDataMap[\"success\"]", isok, result)
	//	return nil, errors.New("err")
	//}
	////timestamp := time.Now().Unix()
	//
	//panic("not implement")
	//return nil, nil
}

func (cta *Cryptopia) GetTicker(currency CurrencyPair) (*Ticker, error) {
	currency = cta.adaptCurrencyPair(currency)

	tickerUri := API_BASE_URL + TICKER_URI + currency.ToSymbol("_")
	//log.Println("tickerUrl:", tickerUri)
	bodyDataMap, err := HttpGet(cta.httpClient, tickerUri)
	//log.Println("Cryptopia bodyDataMap:", tickerUri, bodyDataMap)

	if err != nil {
		//log.Println(err)
		return nil, err
	}
	tickerMap, isok := bodyDataMap["Data"].(map[string]interface{})
	if isok != true {
		//log.Println("Cryptopia bodyDataMap:", tickerUri, bodyDataMap)
		//log.Println("bodyDataMap[\"Error\"]", bodyDataMap["Error"].(string))
		//return nil, errors.New(bodyDataMap["Error"].(string))
		return nil, errors.New("ERR")
	}
	var ticker Ticker

	timestamp := time.Now().Unix()

	//fmt.Println(bodyDataMap)
	ticker.Date = uint64(timestamp)
	ticker.Last, _ = tickerMap["LastPrice"].(float64)

	ticker.Buy, _ = tickerMap["BidPrice"].(float64)
	ticker.Sell, _ = tickerMap["AskPrice"].(float64)
	ticker.Vol, _ = tickerMap["Volume"].(float64)
	//log.Println("Cryptopia", currency, "ticker:", ticker)
	return &ticker, nil
}

func (cta *Cryptopia) GetTickerInBuf(currency CurrencyPair) (*Ticker, error) {
	return cta.GetTicker(currency)
}

func (cta *Cryptopia) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	panic("not implement")
}

func (cta *Cryptopia) adaptCurrencyPair(pair CurrencyPair) CurrencyPair {
	var currencyA Currency
	var currencyB Currency

	if pair.CurrencyA == BCC {
		currencyA = BCH
	} else {
		currencyA = pair.CurrencyA
	}
	currencyB = pair.CurrencyB
	//if pair.BaseCurrency == USDT {
	//	currencyB = USD
	//} else {
	//	currencyB = pair.BaseCurrency
	//}

	return NewCurrencyPair(currencyA, currencyB)
}

func (cta *Cryptopia) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (cta *Cryptopia) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (cta *Cryptopia) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (cta *Cryptopia) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (cta *Cryptopia) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("not implements")
}

func (cta *Cryptopia) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}
func (cta *Cryptopia) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implements")
}

func (cta *Cryptopia) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implements")
}

func (cta *Cryptopia) GetAccount() (*Account, error) {
	panic("not implements")
}
