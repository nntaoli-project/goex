package zaif

import (
	"fmt"
	. "github.com/nntaoli/crypto_coin_api"
	"log"
	"net/http"
	"sort"
)

var ToCurrency = map[CurrencyPair]string{
	BTC_JPY: "btc",
}

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
	tickerUrl := fmt.Sprintf(zf.baseUrl+"1/ticker/%s_jpy", ToCurrency[currency])
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
	depthUrl := fmt.Sprintf(zf.baseUrl+"1/depth/%s_jpy", ToCurrency[currency])
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
