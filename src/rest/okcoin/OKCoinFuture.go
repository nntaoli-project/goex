package okcoin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	. "rest"
	"strconv"
)

const (
	FUTURE_API_BASE_URL    = "https://www.okcoin.com/api/v1/"
	FUTURE_TICKER_URI      = "future_ticker.do?symbol=%s&contract_type=%s"
	FUTURE_DEPTH_URI       = "future_depth.do?symbol=%s&contract_type=%s"
	FUTURE_USERINFO_URI    = "future_userinfo.do"
	FUTURE_CANCEL_URI      = "future_cancel.do"
	FUTURE_ORDER_INFO_URI  = "future_order_info.do"
	FUTURE_ORDERS_INFO_URI = "future_orders_info.do"
	FUTURE_POSITION_URI    = "future_position.do"
)

type OKCoinFuture struct {
	apiKey,
	apiSecretKey string
	client *http.Client
}

func NewFuture(client *http.Client, api_key, secret_key string) *OKCoinFuture {
	ok := new(OKCoinFuture)
	ok.apiKey = api_key
	ok.apiSecretKey = secret_key
	ok.client = client
	return ok
}

func (ok *OKCoinFuture) GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error) {
	url := FUTURE_API_BASE_URL + FUTURE_TICKER_URI
	//fmt.Println(fmt.Sprintf(url, CurrencyPairSymbol[currencyPair], contractType));
	resp, err := ok.client.Get(fmt.Sprintf(url, CurrencyPairSymbol[currencyPair], contractType))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	//println(string(body))

	bodyMap := make(map[string]interface{})

	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		return nil, err
	}

	tickerMap := bodyMap["ticker"].(map[string]interface{})

	ticker := new(Ticker)
	ticker.Date, _ = strconv.ParseUint(bodyMap["date"].(string), 10, 64)
	ticker.Buy = tickerMap["buy"].(float64)
	ticker.Sell = tickerMap["sell"].(float64)
	ticker.Last = tickerMap["last"].(float64)
	ticker.High = tickerMap["high"].(float64)
	ticker.Low = tickerMap["low"].(float64)
	ticker.Vol = tickerMap["vol"].(float64)

	//fmt.Println(bodyMap)
	return ticker, nil
}

func (ok *OKCoinFuture) GetFutureDepth(currencyPair CurrencyPair, contractType string) (*Depth, error) {
	url := FUTURE_API_BASE_URL + FUTURE_DEPTH_URI
	//fmt.Println(fmt.Sprintf(url, CurrencyPairSymbol[currencyPair], contractType));
	resp, err := ok.client.Get(fmt.Sprintf(url, CurrencyPairSymbol[currencyPair], contractType))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	//println(string(body))

	bodyMap := make(map[string]interface{})

	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		return nil, err
	}

	depth := new(Depth)

	for _, v := range bodyMap["asks"].([]interface{}) {
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

	for _, v := range bodyMap["bids"].([]interface{}) {
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

	//fmt.Println(bodyMap)
	return depth, nil
}

func (ok *OKCoinFuture) GetFutureIndex(currencyPair CurrencyPair) (float64, error) {
	return 0, nil
}

func (ok *OKCoinFuture) GetFutureUserinfo() (*FutureAccount, error) {
	return nil, nil
}

func (ok *OKCoinFuture) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount, openType, matchPrice string) (string, error) {
	return "", nil
}

func (ok *OKCoinFuture) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error) {
	return false, nil
}

func (ok *OKCoinFuture) GetFuturePosition(currencyPair CurrencyPair, contractType string) (*FuturePosition, error) {
	return nil, nil
}

func (ok *OKCoinFuture) GetFutureOrders(orderId int64, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	return nil, nil
}

func (ok *OKCoinFuture) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	return nil, nil
}
