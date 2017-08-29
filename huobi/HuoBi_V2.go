package huobi

import (
	"errors"
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
