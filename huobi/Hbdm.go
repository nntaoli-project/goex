package huobi

import (
	"errors"
	. "github.com/nntaoli-project/GoEx"
	"internal/log"
	"sort"
	"time"
)

type Hbdm struct {
	config *APIConfig
}

var (
	apiUrl = "https://api.hbdm.com"
)

func NewHbdm(conf *APIConfig) *Hbdm {
	return &Hbdm{conf}
}

func (dm *Hbdm) GetExchangeName() string {
	return HBDM
}

func (dm *Hbdm) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	panic("not implement")
}

func (dm *Hbdm) GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error) {
	symbol := dm.adaptSymbol(currencyPair, contractType)
	ret, err := HttpGet(dm.config.HttpClient, apiUrl+"/market/detail/merged?symbol="+symbol)
	if err != nil {
		return nil, err
	}
	//log.Println(ret)
	s := ret["status"].(string)
	if s == "error" {
		return nil, errors.New(ret["err-msg"].(string))
	}

	tick := ret["tick"].(map[string]interface{})
	return &Ticker{
		Pair: currencyPair,
		Last: ToFloat64(tick["close"]),
		Vol:  ToFloat64(tick["amount"]),
		Low:  ToFloat64(tick["low"]),
		High: ToFloat64(tick["high"]),
		Sell: ToFloat64(tick["ask"].([]interface{})[0]),
		Buy:  ToFloat64(tick["bid"].([]interface{})[0]),
		Date: ToUint64(ret["ts"])}, nil
}

func (dm *Hbdm) GetFutureDepth(currencyPair CurrencyPair, contractType string, size int) (*Depth, error) {
	symbol := dm.adaptSymbol(currencyPair, contractType)
	url := apiUrl + "/market/depth?type=step0&symbol=" + symbol
	ret, err := HttpGet(dm.config.HttpClient, url)
	if err != nil {
		return nil, err
	}

	s := ret["status"].(string)
	if s == "error" {
		return nil, errors.New(ret["err-msg"].(string))
	}
	log.Println(ret)
	dep := new(Depth)
	dep.Pair = currencyPair
	dep.ContractType = symbol

	mills := ToUint64(ret["ts"])
	dep.UTime = time.Unix(int64(mills/1000), int64(mills%1000)*int64(time.Millisecond))

	tick := ret["tick"].(map[string]interface{})
	asks := tick["asks"].([]interface{})
	bids := tick["bids"].([]interface{})

	for _, item := range asks {
		askItem := item.([]interface{})
		dep.AskList = append(dep.AskList, DepthRecord{ToFloat64(askItem[0]), ToFloat64(askItem[1])})
	}

	for _, item := range bids {
		bidItem := item.([]interface{})
		dep.BidList = append(dep.BidList, DepthRecord{ToFloat64(bidItem[0]), ToFloat64(bidItem[1])})
	}

	sort.Sort(sort.Reverse(dep.AskList))

	return dep, nil
}

func (dm *Hbdm) adaptSymbol(pair CurrencyPair, contractType string) string {
	symbol := pair.CurrencyA.Symbol + "_"
	switch contractType {
	case THIS_WEEK_CONTRACT:
		symbol += "CW"
	case NEXT_WEEK_CONTRACT:
		symbol += "NW"
	case QUARTER_CONTRACT:
		symbol += "CQ"
	}
	return symbol
}
