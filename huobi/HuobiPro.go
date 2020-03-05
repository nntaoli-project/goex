package huobi

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
	. "github.com/nntaoli-project/goex/internal/logger"
	"math/big"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

var HBPOINT = NewCurrency("HBPOINT", "")

var _INERNAL_KLINE_PERIOD_CONVERTER = map[int]string{
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

const (
	HB_POINT_ACCOUNT = "point"
	HB_SPOT_ACCOUNT  = "spot"
)

type AccountInfo struct {
	Id    string
	Type  string
	State string
}

type HuoBiPro struct {
	httpClient *http.Client
	baseUrl    string
	accountId  string
	accessKey  string
	secretKey  string
	//ECDSAPrivateKey string
}

type HuoBiProSymbol struct {
	BaseCurrency    string
	QuoteCurrency   string
	PricePrecision  float64
	AmountPrecision float64
	SymbolPartition string
	Symbol          string
}

func NewHuobiWithConfig(config *APIConfig) *HuoBiPro {
	hbpro := new(HuoBiPro)
	if config.Endpoint == "" {
		hbpro.baseUrl = "https://api.huobi.pro"
	}
	hbpro.httpClient = config.HttpClient
	hbpro.accessKey = config.ApiKey
	hbpro.secretKey = config.ApiSecretKey

	if config.ApiKey != "" && config.ApiSecretKey != "" {
		accinfo, err := hbpro.GetAccountInfo(HB_SPOT_ACCOUNT)
		if err != nil {
			hbpro.accountId = ""
			//panic(err)
		} else {
			hbpro.accountId = accinfo.Id
			//log.Println("account state :", accinfo.State)
			Log.Info("accountId=", accinfo.Id, ",state=", accinfo.State, ",type=", accinfo.Type)
		}
	}

	return hbpro
}

func NewHuoBiPro(client *http.Client, apikey, secretkey, accountId string) *HuoBiPro {
	hbpro := new(HuoBiPro)
	hbpro.baseUrl = "https://api.huobi.pro"
	hbpro.httpClient = client
	hbpro.accessKey = apikey
	hbpro.secretKey = secretkey
	hbpro.accountId = accountId
	return hbpro
}

/**
 *现货交易
 */
func NewHuoBiProSpot(client *http.Client, apikey, secretkey string) *HuoBiPro {
	hb := NewHuoBiPro(client, apikey, secretkey, "")
	accinfo, err := hb.GetAccountInfo(HB_SPOT_ACCOUNT)
	if err != nil {
		hb.accountId = ""
		//panic(err)
	} else {
		hb.accountId = accinfo.Id
		Log.Info("account state :", accinfo.State)
	}
	return hb
}

/**
 * 点卡账户
 */
func NewHuoBiProPoint(client *http.Client, apikey, secretkey string) *HuoBiPro {
	hb := NewHuoBiPro(client, apikey, secretkey, "")
	accinfo, err := hb.GetAccountInfo(HB_POINT_ACCOUNT)
	if err != nil {
		panic(err)
	}
	hb.accountId = accinfo.Id
	Log.Info("account state :" + accinfo.State)
	return hb
}

func (hbpro *HuoBiPro) GetAccountInfo(acc string) (AccountInfo, error) {
	path := "/v1/account/accounts"
	params := &url.Values{}
	hbpro.buildPostForm("GET", path, params)

	//log.Println(hbpro.baseUrl + path + "?" + params.Encode())

	respmap, err := HttpGet(hbpro.httpClient, hbpro.baseUrl+path+"?"+params.Encode())
	if err != nil {
		return AccountInfo{}, err
	}

	if respmap["status"].(string) != "ok" {
		return AccountInfo{}, errors.New(respmap["err-code"].(string))
	}

	var info AccountInfo

	data := respmap["data"].([]interface{})
	for _, v := range data {
		iddata := v.(map[string]interface{})
		if iddata["type"].(string) == acc {
			info.Id = fmt.Sprintf("%.0f", iddata["id"])
			info.Type = acc
			info.State = iddata["state"].(string)
			break
		}
	}
	//log.Println(respmap)
	return info, nil
}

func (hbpro *HuoBiPro) GetAccount() (*Account, error) {
	path := fmt.Sprintf("/v1/account/accounts/%s/balance", hbpro.accountId)
	params := &url.Values{}
	params.Set("accountId-id", hbpro.accountId)
	hbpro.buildPostForm("GET", path, params)

	urlStr := hbpro.baseUrl + path + "?" + params.Encode()
	//println(urlStr)
	respmap, err := HttpGet(hbpro.httpClient, urlStr)

	if err != nil {
		return nil, err
	}

	//log.Println(respmap)

	if respmap["status"].(string) != "ok" {
		return nil, errors.New(respmap["err-code"].(string))
	}

	datamap := respmap["data"].(map[string]interface{})
	if datamap["state"].(string) != "working" {
		return nil, errors.New(datamap["state"].(string))
	}

	list := datamap["list"].([]interface{})
	acc := new(Account)
	acc.SubAccounts = make(map[Currency]SubAccount, 6)
	acc.Exchange = hbpro.GetExchangeName()

	subAccMap := make(map[Currency]*SubAccount)

	for _, v := range list {
		balancemap := v.(map[string]interface{})
		currencySymbol := balancemap["currency"].(string)
		currency := NewCurrency(currencySymbol, "")
		typeStr := balancemap["type"].(string)
		balance := ToFloat64(balancemap["balance"])
		if subAccMap[currency] == nil {
			subAccMap[currency] = new(SubAccount)
		}
		subAccMap[currency].Currency = currency
		switch typeStr {
		case "trade":
			subAccMap[currency].Amount = balance
		case "frozen":
			subAccMap[currency].ForzenAmount = balance
		}
	}

	for k, v := range subAccMap {
		acc.SubAccounts[k] = *v
	}

	return acc, nil
}

func (hbpro *HuoBiPro) placeOrder(amount, price string, pair CurrencyPair, orderType string) (string, error) {
	path := "/v1/order/orders/place"
	params := url.Values{}
	params.Set("account-id", hbpro.accountId)
	params.Set("amount", amount)
	params.Set("symbol", pair.AdaptUsdToUsdt().ToLower().ToSymbol(""))
	params.Set("type", orderType)

	switch orderType {
	case "buy-limit", "sell-limit":
		params.Set("price", price)
	}

	hbpro.buildPostForm("POST", path, &params)

	resp, err := HttpPostForm3(hbpro.httpClient, hbpro.baseUrl+path+"?"+params.Encode(), hbpro.toJson(params),
		map[string]string{"Content-Type": "application/json", "Accept-Language": "zh-cn"})
	if err != nil {
		return "", err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return "", err
	}

	if respmap["status"].(string) != "ok" {
		return "", errors.New(respmap["err-code"].(string))
	}

	return respmap["data"].(string), nil
}

func (hbpro *HuoBiPro) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	orderId, err := hbpro.placeOrder(amount, price, currency, "buy-limit")
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		OrderID:  ToInt(orderId),
		OrderID2: orderId,
		Amount:   ToFloat64(amount),
		Price:    ToFloat64(price),
		Side:     BUY}, nil
}

func (hbpro *HuoBiPro) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	orderId, err := hbpro.placeOrder(amount, price, currency, "sell-limit")
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		OrderID:  ToInt(orderId),
		OrderID2: orderId,
		Amount:   ToFloat64(amount),
		Price:    ToFloat64(price),
		Side:     SELL}, nil
}

func (hbpro *HuoBiPro) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	orderId, err := hbpro.placeOrder(amount, price, currency, "buy-market")
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		OrderID:  ToInt(orderId),
		OrderID2: orderId,
		Amount:   ToFloat64(amount),
		Price:    ToFloat64(price),
		Side:     BUY_MARKET}, nil
}

func (hbpro *HuoBiPro) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	orderId, err := hbpro.placeOrder(amount, price, currency, "sell-market")
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		OrderID:  ToInt(orderId),
		OrderID2: orderId,
		Amount:   ToFloat64(amount),
		Price:    ToFloat64(price),
		Side:     SELL_MARKET}, nil
}

func (hbpro *HuoBiPro) parseOrder(ordmap map[string]interface{}) Order {
	ord := Order{
		OrderID:    ToInt(ordmap["id"]),
		OrderID2:   fmt.Sprint(ToInt(ordmap["id"])),
		Amount:     ToFloat64(ordmap["amount"]),
		Price:      ToFloat64(ordmap["price"]),
		DealAmount: ToFloat64(ordmap["field-amount"]),
		Fee:        ToFloat64(ordmap["field-fees"]),
		OrderTime:  ToInt(ordmap["created-at"]),
	}

	state := ordmap["state"].(string)
	switch state {
	case "submitted", "pre-submitted":
		ord.Status = ORDER_UNFINISH
	case "filled":
		ord.Status = ORDER_FINISH
	case "partial-filled":
		ord.Status = ORDER_PART_FINISH
	case "canceled", "partial-canceled":
		ord.Status = ORDER_CANCEL
	default:
		ord.Status = ORDER_UNFINISH
	}

	if ord.DealAmount > 0.0 {
		ord.AvgPrice = ToFloat64(ordmap["field-cash-amount"]) / ord.DealAmount
	}

	typeS := ordmap["type"].(string)
	switch typeS {
	case "buy-limit":
		ord.Side = BUY
	case "buy-market":
		ord.Side = BUY_MARKET
	case "sell-limit":
		ord.Side = SELL
	case "sell-market":
		ord.Side = SELL_MARKET
	}
	return ord
}

func (hbpro *HuoBiPro) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	path := "/v1/order/orders/" + orderId
	params := url.Values{}
	hbpro.buildPostForm("GET", path, &params)
	respmap, err := HttpGet(hbpro.httpClient, hbpro.baseUrl+path+"?"+params.Encode())
	if err != nil {
		return nil, err
	}

	if respmap["status"].(string) != "ok" {
		return nil, errors.New(respmap["err-code"].(string))
	}

	datamap := respmap["data"].(map[string]interface{})
	order := hbpro.parseOrder(datamap)
	order.Currency = currency

	return &order, nil
}

func (hbpro *HuoBiPro) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	return hbpro.getOrders(queryOrdersParams{
		pair:   currency,
		states: "pre-submitted,submitted,partial-filled",
		size:   100,
		//direct:""
	})
}

func (hbpro *HuoBiPro) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	path := fmt.Sprintf("/v1/order/orders/%s/submitcancel", orderId)
	params := url.Values{}
	hbpro.buildPostForm("POST", path, &params)
	resp, err := HttpPostForm3(hbpro.httpClient, hbpro.baseUrl+path+"?"+params.Encode(), hbpro.toJson(params),
		map[string]string{"Content-Type": "application/json", "Accept-Language": "zh-cn"})
	if err != nil {
		return false, err
	}

	var respmap map[string]interface{}
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return false, err
	}

	if respmap["status"].(string) != "ok" {
		return false, errors.New(string(resp))
	}

	return true, nil
}

func (hbpro *HuoBiPro) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	return hbpro.getOrders(queryOrdersParams{
		pair:   currency,
		size:   pageSize,
		states: "partial-canceled,filled",
		direct: "next",
	})
}

type queryOrdersParams struct {
	types,
	startDate,
	endDate,
	states,
	from,
	direct string
	size int
	pair CurrencyPair
}

func (hbpro *HuoBiPro) getOrders(queryparams queryOrdersParams) ([]Order, error) {
	path := "/v1/order/orders"
	params := url.Values{}
	params.Set("symbol", strings.ToLower(queryparams.pair.AdaptUsdToUsdt().ToSymbol("")))
	params.Set("states", queryparams.states)

	if queryparams.direct != "" {
		params.Set("direct", queryparams.direct)
	}

	if queryparams.size > 0 {
		params.Set("size", fmt.Sprint(queryparams.size))
	}

	hbpro.buildPostForm("GET", path, &params)
	respmap, err := HttpGet(hbpro.httpClient, fmt.Sprintf("%s%s?%s", hbpro.baseUrl, path, params.Encode()))
	if err != nil {
		return nil, err
	}

	if respmap["status"].(string) != "ok" {
		return nil, errors.New(respmap["err-code"].(string))
	}

	datamap := respmap["data"].([]interface{})
	var orders []Order
	for _, v := range datamap {
		ordmap := v.(map[string]interface{})
		ord := hbpro.parseOrder(ordmap)
		ord.Currency = queryparams.pair
		orders = append(orders, ord)
	}

	return orders, nil
}

func (hbpro *HuoBiPro) GetTicker(currencyPair CurrencyPair) (*Ticker, error) {
	pair := currencyPair.AdaptUsdToUsdt()
	url := hbpro.baseUrl + "/market/detail/merged?symbol=" + strings.ToLower(pair.ToSymbol(""))
	respmap, err := HttpGet(hbpro.httpClient, url)
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
	ticker.Pair = currencyPair
	ticker.Vol = ToFloat64(tickmap["amount"])
	ticker.Low = ToFloat64(tickmap["low"])
	ticker.High = ToFloat64(tickmap["high"])
	bid, isOk := tickmap["bid"].([]interface{})
	if isOk != true {
		return nil, errors.New("no bid")
	}
	ask, isOk := tickmap["ask"].([]interface{})
	if isOk != true {
		return nil, errors.New("no ask")
	}
	ticker.Buy = ToFloat64(bid[0])
	ticker.Sell = ToFloat64(ask[0])
	ticker.Last = ToFloat64(tickmap["close"])
	ticker.Date = ToUint64(respmap["ts"])

	return ticker, nil
}

func (hbpro *HuoBiPro) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	url := hbpro.baseUrl + "/market/depth?symbol=%s&type=step0&depth=%d"
	n := 5
	pair := currency.AdaptUsdToUsdt()
	if size <= 5 {
		n = 5
	} else if size <= 10 {
		n = 10
	} else if size <= 20 {
		n = 20
	} else {
		url = hbpro.baseUrl + "/market/depth?symbol=%s&type=step0&d=%d"
	}
	respmap, err := HttpGet(hbpro.httpClient, fmt.Sprintf(url, strings.ToLower(pair.ToSymbol("")), n))
	if err != nil {
		return nil, err
	}

	if "ok" != respmap["status"].(string) {
		return nil, errors.New(respmap["err-msg"].(string))
	}

	tick, _ := respmap["tick"].(map[string]interface{})

	dep := hbpro.parseDepthData(tick, size)
	dep.Pair = currency
	mills := ToUint64(tick["ts"])
	dep.UTime = time.Unix(int64(mills/1000), int64(mills%1000)*int64(time.Millisecond))

	return dep, nil
}

//倒序
func (hbpro *HuoBiPro) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	url := hbpro.baseUrl + "/market/history/kline?period=%s&size=%d&symbol=%s"
	symbol := strings.ToLower(currency.AdaptUsdToUsdt().ToSymbol(""))
	periodS, isOk := _INERNAL_KLINE_PERIOD_CONVERTER[period]
	if isOk != true {
		periodS = "1min"
	}

	ret, err := HttpGet(hbpro.httpClient, fmt.Sprintf(url, periodS, size, symbol))
	if err != nil {
		return nil, err
	}

	data, ok := ret["data"].([]interface{})
	if !ok {
		return nil, errors.New("response format error")
	}

	var klines []Kline
	for _, e := range data {
		item := e.(map[string]interface{})
		klines = append(klines, Kline{
			Pair:      currency,
			Open:      ToFloat64(item["open"]),
			Close:     ToFloat64(item["close"]),
			High:      ToFloat64(item["high"]),
			Low:       ToFloat64(item["low"]),
			Vol:       ToFloat64(item["vol"]),
			Timestamp: int64(ToUint64(item["id"]))})
	}

	return klines, nil
}

//非个人，整个交易所的交易记录
//https://github.com/huobiapi/API_Docs/wiki/REST_api_reference#get-markettrade-获取-trade-detail-数据
func (hbpro *HuoBiPro) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	var (
		trades []Trade
		ret    struct {
			Status string
			ErrMsg string `json:"err-msg"`
			Data   []struct {
				Ts   int64
				Data []struct {
					Id        big.Int
					Amount    float64
					Price     float64
					Direction string
					Ts        int64
				}
			}
		}
	)

	url := hbpro.baseUrl + "/market/history/trade?size=2000&symbol=" + currencyPair.AdaptUsdToUsdt().ToLower().ToSymbol("")
	err := HttpGet4(hbpro.httpClient, url, map[string]string{}, &ret)
	if err != nil {
		return nil, err
	}

	if ret.Status != "ok" {
		return nil, errors.New(ret.ErrMsg)
	}

	for _, d := range ret.Data {
		for _, t := range d.Data {

			//fix huobi   Weird rules of tid
			//火币交易ID规定固定23位, 导致超出int64范围，每个交易对有不同的固定填充前缀
			//实际交易ID远远没有到23位数字。
			tid := ToInt64(strings.TrimPrefix(t.Id.String()[4:], "0"))
			if tid == 0 {
				tid = ToInt64(strings.TrimPrefix(t.Id.String()[5:], "0"))
			}
			///

			trades = append(trades, Trade{
				Tid:    ToInt64(tid),
				Pair:   currencyPair,
				Amount: t.Amount,
				Price:  t.Price,
				Type:   AdaptTradeSide(t.Direction),
				Date:   t.Ts})
		}
	}

	return trades, nil
}

type ecdsaSignature struct {
	R, S *big.Int
}

func (hbpro *HuoBiPro) buildPostForm(reqMethod, path string, postForm *url.Values) error {
	postForm.Set("AccessKeyId", hbpro.accessKey)
	postForm.Set("SignatureMethod", "HmacSHA256")
	postForm.Set("SignatureVersion", "2")
	postForm.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05"))
	domain := strings.Replace(hbpro.baseUrl, "https://", "", len(hbpro.baseUrl))
	payload := fmt.Sprintf("%s\n%s\n%s\n%s", reqMethod, domain, path, postForm.Encode())
	sign, _ := GetParamHmacSHA256Base64Sign(hbpro.secretKey, payload)
	postForm.Set("Signature", sign)

	/**
	p, _ := pem.Decode([]byte(hbpro.ECDSAPrivateKey))
	pri, _ := secp256k1_go.PrivKeyFromBytes(secp256k1_go.S256(), p.Bytes)
	signer, _ := pri.Sign([]byte(sign))
	signAsn, _ := asn1.Marshal(signer)
	priSign := base64.StdEncoding.EncodeToString(signAsn)
	postForm.Set("PrivateSignature", priSign)
	*/

	return nil
}

func (hbpro *HuoBiPro) toJson(params url.Values) string {
	parammap := make(map[string]string)
	for k, v := range params {
		parammap[k] = v[0]
	}
	jsonData, _ := json.Marshal(parammap)
	return string(jsonData)
}

func (hbpro *HuoBiPro) parseDepthData(tick map[string]interface{}, size int) *Depth {
	bids, _ := tick["bids"].([]interface{})
	asks, _ := tick["asks"].([]interface{})

	depth := new(Depth)
	n := 0
	for _, r := range asks {
		var dr DepthRecord
		rr := r.([]interface{})
		dr.Price = ToFloat64(rr[0])
		dr.Amount = ToFloat64(rr[1])
		depth.AskList = append(depth.AskList, dr)
		n++
		if n == size {
			break
		}
	}

	n = 0
	for _, r := range bids {
		var dr DepthRecord
		rr := r.([]interface{})
		dr.Price = ToFloat64(rr[0])
		dr.Amount = ToFloat64(rr[1])
		depth.BidList = append(depth.BidList, dr)
		n++
		if n == size {
			break
		}
	}

	sort.Sort(sort.Reverse(depth.AskList))

	return depth
}

func (hbpro *HuoBiPro) GetExchangeName() string {
	return HUOBI_PRO
}

func (hbpro *HuoBiPro) GetCurrenciesList() ([]string, error) {
	url := hbpro.baseUrl + "/v1/common/currencys"

	ret, err := HttpGet(hbpro.httpClient, url)
	if err != nil {
		return nil, err
	}

	data, ok := ret["data"].([]interface{})
	if !ok {
		return nil, errors.New("response format error")
	}
	fmt.Println(data)
	return nil, nil
}

func (hbpro *HuoBiPro) GetCurrenciesPrecision() ([]HuoBiProSymbol, error) {
	url := hbpro.baseUrl + "/v1/common/symbols"

	ret, err := HttpGet(hbpro.httpClient, url)
	if err != nil {
		return nil, err
	}

	data, ok := ret["data"].([]interface{})
	if !ok {
		return nil, errors.New("response format error")
	}
	var Symbols []HuoBiProSymbol
	for _, v := range data {
		_sym := v.(map[string]interface{})
		var sym HuoBiProSymbol
		sym.BaseCurrency = _sym["base-currency"].(string)
		sym.QuoteCurrency = _sym["quote-currency"].(string)
		sym.PricePrecision = _sym["price-precision"].(float64)
		sym.AmountPrecision = _sym["amount-precision"].(float64)
		sym.SymbolPartition = _sym["symbol-partition"].(string)
		sym.Symbol = _sym["symbol"].(string)
		Symbols = append(Symbols, sym)
	}
	//fmt.Println(Symbols)
	return Symbols, nil
}
