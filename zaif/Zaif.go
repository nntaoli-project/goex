package zaif

import (
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"log"
	"net/http"
	"sort"
	"strings"
)

type Zaif struct {
	client *http.Client
	baseUrl,
	accessKey,
	secretKey string
}

func New(httpClient *http.Client, accessKey, secretKey string) *Zaif {
	zaif := new(Zaif)
	zaif.accessKey = accessKey
	zaif.secretKey = secretKey
	zaif.client = httpClient
	zaif.baseUrl = "https://api.zaif.jp/api/"
	return zaif
}

func (zf *Zaif) GetExchangeName() string {
	return "zaif.jp"
}

func (zf *Zaif) GetTicker(currency CurrencyPair) (*Ticker, error) {
	tickerUrl := fmt.Sprintf(zf.baseUrl+"1/ticker/%s_jpy", strings.ToLower(currency.CurrencyA.Symbol))
	//println(tickerUrl)
	resp, err := HttpGet(zf.client, tickerUrl)
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
	ticker.Vol = resp["volume"].(float64)
	return ticker, nil
}

func (zf *Zaif) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	depthUrl := fmt.Sprintf(zf.baseUrl+"1/depth/%s_jpy", strings.ToLower(currency.CurrencyA.Symbol))
	resp, err := HttpGet(zf.client, depthUrl)
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
				dr.Price = vv.(float64)
			case 1:
				dr.Amount = vv.(float64)
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
				dr.Price = vv.(float64)
			case 1:
				dr.Amount = vv.(float64)
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

func (zf *Zaif) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return nil, nil
}

func (zf *Zaif) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return nil, nil
}

func (zf *Zaif) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return nil, nil
}

func (zf *Zaif) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return nil, nil
}

func (zf *Zaif) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	return false, nil
}

func (zf *Zaif) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	return nil, nil
}

func (zf *Zaif) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	return nil, nil
}

func (zf *Zaif) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	return nil, nil
}

func (zf *Zaif) GetAccount() (*Account, error) {
	return nil, nil
}

func (zf *Zaif) GetKlineRecords(currency CurrencyPair , period int, size, since int) ([]Kline, error) {
	return nil, nil
}

//非个人，整个交易所的交易记录
func (zf *Zaif) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	return nil, nil
}
