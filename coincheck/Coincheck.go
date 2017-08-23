package coincheck

import (
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"log"
	"net/http"
	//"strconv"
	"sort"
	"strconv"
)

type Coincheck struct {
	client *http.Client
	baseUrl,
	accessKey,
	secretKey string
}

func New(httpClient *http.Client, accessKey, secretKey string) (coinCheck *Coincheck) {
	cc := new(Coincheck)
	cc.client = httpClient
	cc.accessKey = accessKey
	cc.secretKey = secretKey
	cc.baseUrl = "https://coincheck.com/"
	return cc
}

func (cc *Coincheck) GetExchangeName() string {
	return "coincheck.com"
}

func (cc *Coincheck) GetTicker(currency CurrencyPair) (*Ticker, error) {
	tickerUrl := fmt.Sprintf(cc.baseUrl + "api/ticker")

	//println(tickerUrl)
	resp, err := HttpGet(cc.client, tickerUrl)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	//log.Println(resp)
	ticker := new(Ticker)
	ticker.Buy = resp["bid"].(float64)
	ticker.Sell = resp["ask"].(float64)
	ticker.Last = resp["last"].(float64)
	ticker.High = resp["high"].(float64)
	ticker.Low = resp["low"].(float64)
	ticker.Date = uint64(resp["timestamp"].(float64))
	ticker.Vol = resp["volume"].(float64)
	return ticker, nil
}

func (cc *Coincheck) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	depthUrl := cc.baseUrl + "api/order_books"
	resp, err := HttpGet(cc.client, depthUrl)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//log.Println(resp)
	var depth Depth

	//asks, isOK := resp["asks"].([]interface{})
	//if !isOK {
	//	return nil, errors.New("asks assert error")
	//}
	_sz := size
	for _, v := range resp["asks"].([]interface{}) {
		var dr DepthRecord
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price, _ = strconv.ParseFloat(vv.(string), 64)
			case 1:
				dr.Amount, _ = strconv.ParseFloat(vv.(string), 64)
			}
		}
		depth.AskList = append(depth.AskList, dr)
		_sz--
		if _sz == 0 {
			break
		}
	}

	sort.Sort(sort.Reverse(depth.AskList))

	_sz = size
	for _, v := range resp["bids"].([]interface{}) {
		var dr DepthRecord
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price, _ = strconv.ParseFloat(vv.(string), 64)
			case 1:
				dr.Amount, _ = strconv.ParseFloat(vv.(string), 64)
			}
		}
		depth.BidList = append(depth.BidList, dr)

		_sz--
		if _sz == 0 {
			break
		}
	}

	return &depth, nil
}


func (cc *Coincheck) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return nil, nil
}

func (cc *Coincheck) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return nil, nil
}

func (cc *Coincheck) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return nil, nil
}

func (cc *Coincheck) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return nil, nil
}

func (cc *Coincheck) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	return false, nil
}

func (cc *Coincheck) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	return nil, nil
}

func (cc *Coincheck) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	return nil, nil
}

func (cc *Coincheck) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	return nil, nil
}

func (cc *Coincheck) GetAccount() (*Account, error) {
	return nil, nil
}

func (cc *Coincheck) GetKlineRecords(currency CurrencyPair, period , size, since int) ([]Kline, error) {
	return nil, nil
}

//非个人，整个交易所的交易记录
func (cc *Coincheck) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	return nil, nil
}