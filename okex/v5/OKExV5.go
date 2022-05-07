package okex

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	. "github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/internal/logger"
)

const (
	v5RestBaseUrl = "https://www.okex.com"
	v5WsBaseUrl   = "wss://ws.okex.com:8443/ws/v5"

	CONTENT_TYPE          = "Content-Type"
	ACCEPT                = "Accept"
	APPLICATION_JSON_UTF8 = "application/json; charset=UTF-8"
	APPLICATION_JSON      = "application/json"
	OK_ACCESS_KEY         = "OK-ACCESS-KEY"
	OK_ACCESS_SIGN        = "OK-ACCESS-SIGN"
	OK_ACCESS_TIMESTAMP   = "OK-ACCESS-TIMESTAMP"
	OK_ACCESS_PASSPHRASE  = "OK-ACCESS-PASSPHRASE"
)

// base interface for okex v5
type OKExV5 struct {
	config        *APIConfig
	customCIDFunc func() string
}

func NewOKExV5(config *APIConfig) *OKExV5 {
	if config.Endpoint == "" {
		config.Endpoint = v5RestBaseUrl
	}
	okex := &OKExV5{config: config}
	return okex
}

func (ok *OKExV5) ExchangeName() string {
	return OKEX
}

func (ok *OKExV5) SetCustomCID(f func() string) {
	ok.customCIDFunc = f
}

//获取所有产品行情信息
//产品类型instType
// SPOT：币币
// SWAP：永续合约
// FUTURES：交割合约
// OPTION：期权
// func (ok *OKExV5) GetTickersV5(instType, uly string) ([]Ticker, error) {
// 	urlPath := fmt.Sprintf("/api/v5/market/tickers?instType=%s", instType)
// 	if instType == "SWAP" || instType == "FUTURES" || instType == "OPTION" {
// 		urlPath = fmt.Sprintf("%s&uly=%s", urlPath, uly)
// 	}
// 	var response spotTickerResponse
// 	err := ok.OKEx.DoAuthorRequest("GET", urlPath, "", &response)
// 	if err != nil {
// 		return nil, err
// 	}

// 	date, _ := time.Parse(time.RFC3339, response.Timestamp)
// 	return &Ticker{
// 		Pair: currency,
// 		Last: response.Last,
// 		High: response.High24h,
// 		Low:  response.Low24h,
// 		Sell: response.BestAsk,
// 		Buy:  response.BestBid,
// 		Vol:  response.BaseVolume24h,
// 		Date: uint64(time.Duration(date.UnixNano() / int64(time.Millisecond)))}, nil

// }

type TickerV5 struct {
	InstId    string  `json:"instId"`
	Last      float64 `json:"last,string"`
	BuyPrice  float64 `json:"bidPx,string"`
	BuySize   float64 `json:"bidSz,string"`
	SellPrice float64 `json:"askPx,string"`
	SellSize  float64 `json:"askSz,string"`
	Open      float64 `json:"open24h,string"`
	High      float64 `json:"high24h,string"`
	Low       float64 `json:"low24h,string"`
	Vol       float64 `json:"volCcy24h,string"`
	VolQuote  float64 `json:"vol24h,string"`
	Timestamp uint64  `json:"ts,string"` // 单位:ms
}

func (ok *OKExV5) GetTickerV5(instId string) (*TickerV5, error) {
	urlPath := fmt.Sprintf("%s/api/v5/market/ticker?instId=%s", ok.config.Endpoint, instId)
	type TickerV5Response struct {
		Code int        `json:"code,string"`
		Msg  string     `json:"msg"`
		Data []TickerV5 `json:"data"`
	}
	var response TickerV5Response
	err := HttpGet4(ok.config.HttpClient, urlPath, nil, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("GetTickerV5 error:%s", response.Msg)
	}
	return &response.Data[0], nil
}

type DepthV5 struct {
	Asks      [][]string `json:"asks,string"`
	Bids      [][]string `json:"bids,string"`
	Timestamp uint64     `json:"ts,string"` // 单位:ms
}

func (ok *OKExV5) GetDepthV5(instId string, size int) (*DepthV5, error) {

	urlPath := fmt.Sprintf("%s/api/v5/market/books?instId=%s", ok.config.Endpoint, instId)
	if size > 0 {
		urlPath = fmt.Sprintf("%s&sz=%d", urlPath, size)
	}
	type DepthV5Response struct {
		Code int       `json:"code,string"`
		Msg  string    `json:"msg"`
		Data []DepthV5 `json:"data"`
	}
	var response DepthV5Response
	err := HttpGet4(ok.config.HttpClient, urlPath, nil, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("GetDepthV5 error:%s", response.Msg)
	}
	return &response.Data[0], nil
}

func (ok *OKExV5) GetKlineRecordsV5(instId string, period KlinePeriod, params *url.Values) ([][]string, error) {
	urlPath := fmt.Sprintf("%s/api/v5/market/candles?instId=%s&bar=%s", ok.config.Endpoint, instId, ok.adaptKLineBar(period))

	if params.Encode() != "" {
		urlPath = fmt.Sprintf("%s&%s", urlPath, params.Encode())
	}

	logger.Debugf("[OKExV5] GetKlineRecordsV5 Url: %s", urlPath)

	type CandleResponse struct {
		Code int        `json:"code,string"`
		Msg  string     `json:"msg"`
		Data [][]string `json:"data"`
	}

	var response CandleResponse
	err := HttpGet4(ok.config.HttpClient, urlPath, nil, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("GetKlineRecordsV5 error:%s", response.Msg)
	}
	return response.Data, nil
}

/*
 Get a iso time
  eg: 2018-03-16T18:02:48.284Z
*/
func IsoTime() string {
	utcTime := time.Now().UTC()
	iso := utcTime.String()
	isoBytes := []byte(iso)
	iso = string(isoBytes[:10]) + "T" + string(isoBytes[11:23]) + "Z"
	return iso
}

/*
 Get a http request body is a json string and a byte array.
*/
func (ok *OKExV5) BuildRequestBody(params interface{}) (string, *bytes.Reader, error) {
	if params == nil {
		return "", nil, errors.New("illegal parameter")
	}
	data, err := json.Marshal(params)
	if err != nil {
		//log.Println(err)
		return "", nil, errors.New("json convert string error")
	}

	jsonBody := string(data)
	binBody := bytes.NewReader(data)

	return jsonBody, binBody, nil
}

func (ok *OKExV5) doParamSign(httpMethod, uri, requestBody string) (string, string) {
	timestamp := IsoTime()
	preText := fmt.Sprintf("%s%s%s%s", timestamp, strings.ToUpper(httpMethod), uri, requestBody)
	//log.Println("preHash", preText)
	sign, _ := GetParamHmacSHA256Base64Sign(ok.config.ApiSecretKey, preText)
	return sign, timestamp
}

func (ok *OKExV5) DoAuthorRequest(httpMethod, uri, reqBody string, response interface{}) error {
	url := ok.config.Endpoint + uri
	sign, timestamp := ok.doParamSign(httpMethod, uri, reqBody)
	//logger.Log.Debug("timestamp=", timestamp, ", sign=", sign)
	resp, err := NewHttpRequest(ok.config.HttpClient, httpMethod, url, reqBody, map[string]string{
		CONTENT_TYPE: APPLICATION_JSON_UTF8,
		ACCEPT:       APPLICATION_JSON,
		//COOKIE:               LOCALE + "en_US",
		OK_ACCESS_KEY:        ok.config.ApiKey,
		OK_ACCESS_PASSPHRASE: ok.config.ApiPassphrase,
		OK_ACCESS_SIGN:       sign,
		OK_ACCESS_TIMESTAMP:  fmt.Sprint(timestamp)})
	if err != nil {
		//log.Println(err)
		return err
	} else {
		logger.Log.Debug(string(resp))
		return json.Unmarshal(resp, &response)
	}
}

type CreateOrderParam struct {
	Symbol    string //产品ID
	TradeMode string //交易模式,	保证金模式：isolated：逐仓 ；cross：全仓,	非保证金模式：cash：非保证金
	Side      string // 订单方向 buy：买 sell：卖
	OrderType string //订单类型
	// market：市价单
	// limit：限价单
	// post_only：只做maker单
	// fok：全部成交或立即取消
	// ioc：立即成交并取消剩余

	Size        string //	委托数量
	PosSide     string //持仓方向 在双向持仓模式下必填，且仅可选择 long 或 short
	Price       string //委托价格，仅适用于限价单
	CCY         string // 保证金币种，仅适用于单币种保证金模式下的全仓杠杆订单
	ClientOrdId string //客户自定义订单ID	字母（区分大小写）与数字的组合，可以是纯字母、纯数字且长度要在1-32位之间。
	Tag         string //订单标签	字母（区分大小写）与数字的组合，可以是纯字母、纯数字，且长度在1-8位之间。
	ReduceOnly  bool   //是否只减仓，true 或 false，默认false	仅适用于币币杠杆订单
}

type OrderSummaryV5 struct {
	OrdId       string `json:"ordId"`
	ClientOrdId string `json:"clOrdId"` //客户自定义订单ID	字母（区分大小写）与数字的组合，可以是纯字母、纯数字且长度要在1-32位之间。
	Tag         string `json:"tag"`
	SCode       string `json:"sCode"`
	SMsg        string `json:"sMsg"`
}

func (ok *OKExV5) CreateOrder(param *CreateOrderParam) (*OrderSummaryV5, error) {

	reqBody := make(map[string]interface{})

	reqBody["instId"] = param.Symbol
	reqBody["tdMode"] = param.TradeMode
	reqBody["side"] = param.Side
	reqBody["ordType"] = param.OrderType
	reqBody["sz"] = param.Size

	if param.CCY != "" {
		reqBody["ccy"] = param.CCY
	}
	if param.ClientOrdId != "" {
		reqBody["clOrdId"] = param.ClientOrdId
	} else {
		reqBody["clOrdId"] = param.ClientOrdId
	}
	if param.Tag != "" {
		reqBody["tag"] = param.Tag
	}
	if param.PosSide != "" {
		reqBody["posSide"] = param.PosSide
	}
	if param.Price != "" {
		reqBody["px"] = param.Price
	}
	if param.ReduceOnly != false {
		reqBody["reduceOnly"] = param.ReduceOnly
	}

	type OrderResponse struct {
		Code int              `json:"code,string"`
		Msg  string           `json:"msg"`
		Data []OrderSummaryV5 `json:"data"`
	}
	var response OrderResponse

	uri := "/api/v5/trade/order"

	jsonStr, _, _ := ok.BuildRequestBody(reqBody)
	err := ok.DoAuthorRequest(http.MethodPost, uri, jsonStr, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		msg := response.Msg
		if msg == "" {
			if len(response.Data) > 0 {
				msg = fmt.Sprintf("code:%d, scode:%s, smsg:%s", response.Code, response.Data[0].SCode, response.Data[0].SMsg)
			} else {
				msg = fmt.Sprintf("code:%d", response.Code)
			}
		}
		return nil, fmt.Errorf("CreateOrder error:%s", msg)
	}
	return &response.Data[0], nil
}

func (ok *OKExV5) CancelOrderV5(instId, ordId, clOrdId string) (*OrderSummaryV5, error) {

	reqBody := make(map[string]interface{})

	reqBody["instId"] = instId
	if ordId != "" {
		reqBody["ordId"] = ordId
	}
	if clOrdId != "" {
		reqBody["clOrdId"] = clOrdId
	}

	type OrderResponse struct {
		Code int              `json:"code,string"`
		Msg  string           `json:"msg"`
		Data []OrderSummaryV5 `json:"data"`
	}
	var response OrderResponse

	uri := "/api/v5/trade/cancel-order"

	jsonStr, _, _ := ok.BuildRequestBody(reqBody)
	err := ok.DoAuthorRequest(http.MethodPost, uri, jsonStr, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		msg := response.Msg
		if msg == "" {
			if len(response.Data) > 0 {
				msg = fmt.Sprintf("code:%d, scode:%s, smsg:%s", response.Code, response.Data[0].SCode, response.Data[0].SMsg)
			} else {
				msg = fmt.Sprintf("code:%d", response.Code)
			}
		}
		return nil, fmt.Errorf("CancelOrderV5 error:%s", msg)
	}
	return &response.Data[0], nil
}

type PendingOrderParam struct {
	InstType string
	Uly      string
	InstId   string //产品ID
	OrdType  string
	State    string
	After    string
	Before   string
	Limit    string
}

type OrderV5 struct {
	AccFillSz   string  `json:"accFillSz"`
	AvgPx       string  `json:"avgPx"`
	CTime       int     `json:"cTime,string"`
	Category    string  `json:"category"`
	Ccy         string  `json:"ccy"`
	ClOrdID     string  `json:"clOrdId"`
	Fee         float64 `json:"fee,string"`
	FeeCcy      string  `json:"feeCcy"`
	FillPx      string  `json:"fillPx"`
	FillSz      string  `json:"fillSz"`
	FillTime    string  `json:"fillTime"`
	InstID      string  `json:"instId"`
	InstType    string  `json:"instType"`
	Lever       string  `json:"lever"`
	OrdID       string  `json:"ordId"`
	OrdType     string  `json:"ordType"`
	Pnl         string  `json:"pnl"`
	PosSide     string  `json:"posSide"`
	Px          float64 `json:"px,string"`
	Rebate      string  `json:"rebate"`
	RebateCcy   string  `json:"rebateCcy"`
	Side        string  `json:"side"`
	SlOrdPx     string  `json:"slOrdPx"`
	SlTriggerPx string  `json:"slTriggerPx"`
	State       string  `json:"state"`
	Sz          float64 `json:"sz,string"`
	Tag         string  `json:"tag"`
	TdMode      string  `json:"tdMode"`
	TpOrdPx     string  `json:"tpOrdPx"`
	TpTriggerPx string  `json:"tpTriggerPx"`
	TradeID     string  `json:"tradeId"`
	UTime       int64   `json:"uTime,string"`
}

func (ok *OKExV5) GetPendingOrders(param *PendingOrderParam) ([]OrderV5, error) {

	reqBody := make(map[string]string)

	if param.InstType != "" {
		reqBody["instType"] = param.InstType
	}
	if param.Uly != "" {
		reqBody["uly"] = param.Uly
	}
	if param.InstId != "" {
		reqBody["instId"] = param.InstId
	}
	if param.OrdType != "" {
		reqBody["ordType"] = param.OrdType
	}
	if param.State != "" {
		reqBody["state"] = param.State
	}
	if param.Before != "" {
		reqBody["before"] = param.Before
	}
	if param.After != "" {
		reqBody["after"] = param.After
	}
	if param.Limit != "" {
		reqBody["limit"] = param.Limit
	}

	type OrderResponse struct {
		Code int       `json:"code,string"`
		Msg  string    `json:"msg"`
		Data []OrderV5 `json:"data"`
	}
	var response OrderResponse

	uri := url.Values{}
	for k, v := range reqBody {
		uri.Set(k, v)
	}
	path := "/api/v5/trade/orders-pending"
	if len(reqBody) > 0 {
		path = fmt.Sprintf("%s?%s", path, uri.Encode())
	}

	jsonStr, _, _ := ok.BuildRequestBody(reqBody)
	err := ok.DoAuthorRequest(http.MethodGet, path, jsonStr, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("GetPendingOrders error:%s", response.Msg)
	}
	return response.Data, nil
}

func (ok *OKExV5) GetOrderV5(instId, ordId, clOrdId string) (*OrderV5, error) {

	reqBody := make(map[string]string)

	reqBody["instId"] = instId
	if ordId != "" {
		reqBody["ordId"] = ordId
	}
	if clOrdId != "" {
		reqBody["clOrdId"] = clOrdId
	}

	type OrderResponse struct {
		Code int       `json:"code,string"`
		Msg  string    `json:"msg"`
		Data []OrderV5 `json:"data"`
	}
	var response OrderResponse

	uri := url.Values{}
	for k, v := range reqBody {
		uri.Set(k, v)
	}
	path := "/api/v5/trade/order"
	if len(reqBody) > 0 {
		path = fmt.Sprintf("%s?%s", path, uri.Encode())
	}

	jsonStr, _, _ := ok.BuildRequestBody(reqBody)
	err := ok.DoAuthorRequest(http.MethodGet, path, jsonStr, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("GetOrderV5 error:%s", response.Msg)
	}
	return &response.Data[0], nil
}

func (ok *OKExV5) GetOrderHistory(instType, instId, ordType, state, afterID, beforeID string) ([]OrderV5, error) {

	reqBody := make(map[string]string)

	reqBody["instType"] = instType
	if instId != "" {
		reqBody["instId"] = instId
	}
	if ordType != "" {
		reqBody["ordType"] = ordType
	}
	if state != "" {
		reqBody["state"] = state
	}
	if afterID != "" {
		reqBody["after"] = afterID
	}
	if beforeID != "" {
		reqBody["before"] = beforeID
	}
	// reqBody["limit"] = "100"

	type OrderResponse struct {
		Code int       `json:"code,string"`
		Msg  string    `json:"msg"`
		Data []OrderV5 `json:"data"`
	}
	var response OrderResponse

	uri := url.Values{}
	for k, v := range reqBody {
		uri.Set(k, v)
	}
	path := "/api/v5/trade/orders-history-archive"
	if len(reqBody) > 0 {
		path = fmt.Sprintf("%s?%s", path, uri.Encode())
	}

	jsonStr, _, _ := ok.BuildRequestBody(reqBody)
	err := ok.DoAuthorRequest(http.MethodGet, path, jsonStr, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("GetOrderV5 error:%s", response.Msg)
	}
	return response.Data, nil
}

type AssetSummary struct {
	Currency  string  `json:"ccy"`
	Total     float64 `json:"bal"`
	Available float64 `json:"availBal,string"`
	Frozen    float64 `json:"frozenBal,string"`
}

func (ok *OKExV5) GetAssetBalances(currency string) ([]AssetSummary, error) {

	reqBody := make(map[string]interface{})

	path := "/api/v5/asset/balances"
	jsonStr := ""
	reqBody["ccy"] = currency
	if currency != "" {
		reqBody["ccy"] = currency
		jsonStr, _, _ = ok.BuildRequestBody(reqBody)
		path = fmt.Sprintf("%s?ccy=%s", path, currency)
	}

	type AssetSummaryResponse struct {
		Code int            `json:"code,string"`
		Msg  string         `json:"msg"`
		Data []AssetSummary `json:"data"`
	}
	var response AssetSummaryResponse

	err := ok.DoAuthorRequest(http.MethodGet, path, jsonStr, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("GetAssetBalances error:%s", response.Msg)
	}
	return response.Data, nil
}

type BalanceV5 struct {
	AdjEq    string          `json:"adjEq"`
	Details  []BalanceDetail `json:"details"`
	Imr      string          `json:"imr"`
	IsoEq    string          `json:"isoEq"`
	MgnRatio string          `json:"mgnRatio"`
	Mmr      string          `json:"mmr"`
	OrdFroz  string          `json:"ordFroz"`
	TotalEq  string          `json:"totalEq"`
	UTime    string          `json:"uTime"`
}

type BalanceDetail struct {
	Available string `json:"availBal"`
	AvailEq   string `json:"availEq"`
	CashBal   string `json:"cashBal"`
	Currency  string `json:"ccy"`
	DisEq     string `json:"disEq"`
	Eq        string `json:"eq"`
	Frozen    string `json:"frozenBal"`
	Interest  string `json:"interest"`
	IsoEq     string `json:"isoEq"`
	Liab      string `json:"liab"`
	MgnRatio  string `json:"mgnRatio"`
	OrdFrozen string `json:"ordFrozen"`
	UTime     string `json:"uTime"`
	Upl       string `json:"upl"`
	UplLiab   string `json:"uplLiab"`
}

func (ok *OKExV5) GetAccountBalances(currency string) (*BalanceV5, error) {

	reqBody := make(map[string]interface{})

	path := "/api/v5/account/balance"
	jsonStr := ""
	reqBody["ccy"] = currency
	if currency != "" {
		reqBody["ccy"] = currency
		jsonStr, _, _ = ok.BuildRequestBody(reqBody)
		path = fmt.Sprintf("%s?ccy=%s", path, currency)
	}

	type BalanceV5Response struct {
		Code int         `json:"code,string"`
		Msg  string      `json:"msg"`
		Data []BalanceV5 `json:"data"`
	}
	var response BalanceV5Response

	err := ok.DoAuthorRequest(http.MethodGet, path, jsonStr, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("GetAccountBalances error:%s", response.Msg)
	}
	return &response.Data[0], nil
}

func (ok *OKExV5) adaptKLineBar(period KlinePeriod) string {
	bar := "1D"
	switch period {
	case KLINE_PERIOD_1MIN:
		bar = "1m"
	case KLINE_PERIOD_3MIN:
		bar = "3m"
	case KLINE_PERIOD_5MIN:
		bar = "5m"
	case KLINE_PERIOD_15MIN:
		bar = "15m"
	case KLINE_PERIOD_30MIN:
		bar = "30m"
	case KLINE_PERIOD_1H, KLINE_PERIOD_60MIN:
		bar = "1H"
	case KLINE_PERIOD_2H:
		bar = "2H"
	case KLINE_PERIOD_4H:
		bar = "4H"
	case KLINE_PERIOD_6H:
		bar = "6H"
	case KLINE_PERIOD_12H:
		bar = "12H"
	case KLINE_PERIOD_1DAY:
		bar = "1D"
	case KLINE_PERIOD_1WEEK:
		bar = "1W"
	default:
		bar = "1D"
	}
	return bar
}
