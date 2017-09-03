package huobi

import (
	"errors"
	"fmt"
	. "github.com/nntaoli-project/GoEx"
	"net/http"
	"strings"
	"net/url"
	"strconv"
	"time"
)

const (
	EXCHANGE_NAME_V2 = "huobi.pro/huobi.com"
	BE_API_BASE_URL  = "https://be.huobi.com"
	PRO_API_BASE_URL = "https://api.huobi.pro"

	V2_TICKER_URI_MERGED = "/market/detail/merged?symbol=%s"
	V2_TICKER_URI        = "/market/trade?symbol=%s"

	DEPTH_URL        = "/market/depth?symbol=%s&type=step1"
	GET_ACCOUNT_API  = "/v1/account/accounts/%s/balance"
	CREATE_ORDER_API = "/v1/order/orders"
	PLACE_ORDER_API  = "/v1/order/orders/%d/place"

	API_URL    = "https://be.huobi.com"
	TICKER_URL = "/market/kline?symbol=%s&period=1min"
)

type HuoBi_V2 struct {
	httpClient *http.Client
	baseUrl,
	accessKey,
	secretKey,
	id string
}

var _V2_INERNAL_KLINE_PERIOD_CONVERTER = map[int]string{
	KLINE_PERIOD_1MIN:   "1min",
	KLINE_PERIOD_5MIN:   "5min",
	KLINE_PERIOD_15MIN:  "15min",
	KLINE_PERIOD_30MIN:  "30min",
	KLINE_PERIOD_60MIN:  "60min",
	KLINE_PERIOD_1DAY:   "1day",
	KLINE_PERIOD_1WEEK:  "1week",
	KLINE_PERIOD_1MONTH: "1mon",
	KLINE_PERIOD_1YEAR:  "1year",
}

func (hbV2 *HuoBi_V2) buildPostForm(reqMethod, path string, postForm *url.Values) error {
	postForm.Set("AccessKeyId", hbV2.accessKey)
	postForm.Set("SignatureMethod", "HmacSHA256")
	postForm.Set("SignatureVersion", "2")
	postForm.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05"))
	url := strings.Split(hbV2.baseUrl, "//")[1]
	domain := strings.Split(url, "/")[0]
	payload := fmt.Sprintf("%s\n%s\n%s\n%s", reqMethod, domain, path, postForm.Encode())
	sign, _ := GetParamHmacSHA256Base64Sign(hbV2.secretKey, payload)
	postForm.Set("Signature", sign)

	return nil
}

func (hbV2 *HuoBi_V2) GetAccountID() (string, error) {
	path := "/v1/account/accounts"
	params := &url.Values{}
	hbV2.buildPostForm("GET", path, params)

	respmap, err := HttpGet(hbV2.httpClient, hbV2.baseUrl+path+"?"+params.Encode())
	if err != nil {
		return "", err
	}
	if respmap["status"].(string) != "ok" {
		return "", errors.New(respmap["err-msg"].(string))
	}
	userdata, isOK := respmap["data"].([]interface{})
	if isOK == false {
		return "", errors.New("No userid")
	}
	u := userdata[0].(map[string]interface{})
	id := strconv.Itoa(int(u["id"].(float64)))
	return id, nil
}

func NewV2(httpClient *http.Client, url, accessKey, secretKey string) *HuoBi_V2 {
	var id string = ""
	hbV2 := &HuoBi_V2{httpClient, url, accessKey, secretKey, id}
	id, err := hbV2.GetAccountID()
	if err == nil {
		hbV2.id = id
		return hbV2
	}
	return nil
}

//func NewV2(httpClient *http.Client, url, accessKey, secretKey string) *HuoBi_V2 {
//	return &HuoBi_V2{httpClient, url, accessKey, secretKey, ""}
//}

func (hbV2 *HuoBi_V2) GetExchangeName() string {
	return EXCHANGE_NAME_V2
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

func (hbV2 *HuoBi_V2) GetDepth(size int, currencyPair CurrencyPair) (*Depth, error) {
	depthUri := hbV2.baseUrl + "/market/depth?symbol=" + strings.ToLower(currencyPair.ToSymbol("")) + "&type=step0"
	respmap, err := HttpGet(hbV2.httpClient, depthUri)
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
	_period := _V2_INERNAL_KLINE_PERIOD_CONVERTER[period]
	if _period == "" {
		return nil, errors.New("unsupport the kline period")
	}
	klineUri := hbV2.baseUrl + "/market/history/kline?symbol=" + strings.ToLower(currency.ToSymbol("")) + "&period=" + _period + "&size=" + strconv.Itoa(size)

	bodyDataMap, err := HttpGet(hbV2.httpClient, klineUri)
	if err != nil {
		return nil, err
	}
	if "ok" != bodyDataMap["status"].(string) {
		return nil, errors.New(bodyDataMap["err-msg"].(string))
	}

	klines, isOK := bodyDataMap["data"].([]interface{})
	if isOK == false {
		return nil, errors.New("No kline")
	}
	fmt.Println(isOK, klines)
	ts := (int64)(bodyDataMap["ts"].(float64))

	var klineRecords []Kline
	for _, record := range klines {
		r := Kline{}
		r.Timestamp = ts
		rec := record.(map[string]interface{})
		r.Open = ToFloat64(rec["open"])
		r.Close = ToFloat64(rec["close"])
		r.High = ToFloat64(rec["high"])
		r.Low = ToFloat64(rec["low"])
		r.Vol = ToFloat64(rec["vol"])
		klineRecords = append(klineRecords, r)
	}

	return klineRecords, nil
}


func (hbV2 *HuoBi_V2) GetAccount() (*Account, error) {
	id, _ := hbV2.GetAccountID()
	path := fmt.Sprintf("/v1/account/accounts/%s/balance", id)
	params := &url.Values{}
	hbV2.buildPostForm("GET", path, params)

	respmap, err := HttpGet(hbV2.httpClient, hbV2.baseUrl+path+"?"+params.Encode())
	if err != nil {
		return nil, err
	}

	account := new(Account)
	account.Exchange = hbV2.GetExchangeName()
	//account.Asset, _ = strconv.ParseFloat(bodyDataMap["total"].(string), 64)
	//account.NetAsset, _ = strconv.ParseFloat(bodyDataMap["net_asset"].(string), 64)

	data, isOK := respmap["data"].(map[string]interface{})
	if isOK == false {
		return nil, errors.New("No account asset")
	}

	list, isOK := data["list"].([]interface{})
	if isOK == false {
		return nil, errors.New("No account asset")
	}

	var btcSubAccount SubAccount
	var ltcSubAccount SubAccount
	var bccSubAccount SubAccount
	var etcSubAccount SubAccount
	var ethSubAccount SubAccount

	for _, cur := range list {
		balance := cur.(map[string]interface{})
		switch balance["currency"].(string) {
		case "btc":
			btcSubAccount.Currency = BTC
			if balance["type"].(string) == "frozen" {
				btcSubAccount.ForzenAmount, _= strconv.ParseFloat(balance["balance"].(string), 64)
			}
			if balance["type"].(string) == "trade" {
				btcSubAccount.Amount, _= strconv.ParseFloat(balance["balance"].(string), 64)
			}
		case "ltc":
			btcSubAccount.Currency = LTC
			if balance["type"].(string) == "frozen" {
				ltcSubAccount.ForzenAmount, _= strconv.ParseFloat(balance["balance"].(string), 64)
			}
			if balance["type"].(string) == "trade" {
				ltcSubAccount.Amount, _= strconv.ParseFloat(balance["balance"].(string), 64)
			}

		case "etc":
			etcSubAccount.Currency = ETC
			if balance["type"].(string) == "frozen" {
				etcSubAccount.ForzenAmount, _= strconv.ParseFloat(balance["balance"].(string), 64)
			}
			if balance["type"].(string) == "trade" {
				etcSubAccount.Amount, _= strconv.ParseFloat(balance["balance"].(string), 64)
			}

		case "bcc":
			bccSubAccount.Currency = BCC
			if balance["type"].(string) == "frozen" {
				bccSubAccount.ForzenAmount, _= strconv.ParseFloat(balance["balance"].(string), 64)
			}
			if balance["type"].(string) == "trade" {
				bccSubAccount.Amount, _= strconv.ParseFloat(balance["balance"].(string), 64)
			}

		case "eth":
			ethSubAccount.Currency = ETH
			if balance["type"].(string) == "frozen" {
				ethSubAccount.ForzenAmount, _= strconv.ParseFloat(balance["balance"].(string), 64)
			}
			if balance["type"].(string) == "trade" {
				ethSubAccount.Amount, _= strconv.ParseFloat(balance["balance"].(string), 64)
			}
		}
	}
	account.SubAccounts = make(map[Currency]SubAccount, 5)
	account.SubAccounts[BTC] = btcSubAccount
	account.SubAccounts[LTC] = ltcSubAccount
	account.SubAccounts[BCC] = bccSubAccount
	account.SubAccounts[ETC] = etcSubAccount
	account.SubAccounts[ETH] = ethSubAccount
	return account, nil
}


func (hbV2 *HuoBi_V2) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("unimplements")
}

func (hbV2 *HuoBi_V2) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("unimplements")
}

func (hbV2 *HuoBi_V2) placeOrder(method, amount, price string, currency CurrencyPair) (*Order, error) {
	panic("unimplements")
}

func (hbV2 *HuoBi_V2) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("unimplements")
}

func (hbV2 *HuoBi_V2) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {

	panic("unimplements")
}

func (hbV2 *HuoBi_V2) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("unimplements")
}

func (hbV2 *HuoBi_V2) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("unimplements")
}

func (hbV2 *HuoBi_V2) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("unimplements")
}


func (hbV2 *HuoBi_V2) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("unimplements")
}

/**
 * 获取全站最近的交易记录
 */
func (hbV2 *HuoBi_V2) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("unimplements")
}
