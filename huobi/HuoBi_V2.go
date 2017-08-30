package huobi

import (
	"errors"
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"net/http"
	"strings"
)

type HuoBi_V2 struct {
	httpClient *http.Client
	baseUrl,
	accessKey,
	secretKey string
}

func NewV2(httpClient *http.Client, accessKey, secretKey string) *HuoBi_V2 {
	return &HuoBi_V2{httpClient, "https://be.huobi.com/", accessKey, secretKey}
}

func (hbV2 *HuoBi_V2) GetExchangeName() string {
	return "huobi.com"
}

func (hbV2 *HuoBi_V2) GetTicker(currencyPair CurrencyPair) (*Ticker, error) {
	url := hbV2.baseUrl + "/market/detail/merged?symbol=" + strings.ToLower(currencyPair.ToSymbol(""))
	respmap, err := HttpGet(hbV2.httpClient, url)
	if err != nil {
		return nil, err
	}

	if respmap["status"].(string) == "error" {
		return nil, errors.New(respmap["err-msg"].(string))
	}

	tickmap, ok := respmap["tick"].(map[string]interface{})
	if !ok {
		return nil, errors.New("tick assert error")
	}

	ticker := new(Ticker)
	ticker.Vol = ToFloat64(tickmap["amount"])
	ticker.Low = ToFloat64(tickmap["low"])
	ticker.High = ToFloat64(tickmap["high"])
	ticker.Buy = ToFloat64((tickmap["bid"].([]interface{}))[0])
	ticker.Sell = ToFloat64((tickmap["ask"].([]interface{}))[0])
	ticker.Last = ToFloat64(tickmap["close"])
	ticker.Date = ToUint64(respmap["ts"])

	return ticker, nil
}

func (hbV2 *HuoBi_V2) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	url := hbV2.baseUrl + "/market/depth?symbol=%s&type=step0"
	respmap, err := HttpGet(hbV2.httpClient, fmt.Sprintf(url, strings.ToLower(currency.ToSymbol(""))))
	if err != nil {
		return nil, err
	}

	if "ok" != respmap["status"].(string) {
		return nil, errors.New(respmap["err-msg"].(string))
	}

	tick, _ := respmap["tick"].(map[string]interface{})
	bids, _ := tick["bids"].([]interface{})
	asks, _ := tick["asks"].([]interface{})

	depth := new(Depth)
	_size := size
	for _, r := range asks {
		var dr DepthRecord
		rr := r.([]interface{})
		dr.Price = ToFloat64(rr[0])
		dr.Amount = ToFloat64(rr[1])
		depth.AskList = append(depth.AskList, dr)

		_size--
		if _size == 0 {
			break
		}
	}

	_size = size
	for _, r := range bids {
		var dr DepthRecord
		rr := r.([]interface{})
		dr.Price = ToFloat64(rr[0])
		dr.Amount = ToFloat64(rr[1])
		depth.BidList = append(depth.BidList, dr)

		_size--
		if _size == 0 {
			break
		}
	}

	return depth, nil
}

func (hbV2 *HuoBi_V2) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录
func (hbV2 *HuoBi_V2) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}
