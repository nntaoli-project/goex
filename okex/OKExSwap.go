package okex

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	. "github.com/nntaoli-project/goex"
)

const (
	/*
	  http headers
	*/
	OK_ACCESS_KEY        = "OK-ACCESS-KEY"
	OK_ACCESS_SIGN       = "OK-ACCESS-SIGN"
	OK_ACCESS_TIMESTAMP  = "OK-ACCESS-TIMESTAMP"
	OK_ACCESS_PASSPHRASE = "OK-ACCESS-PASSPHRASE"

	/**
	  paging params
	*/
	OK_FROM  = "OK-FROM"
	OK_TO    = "OK-TO"
	OK_LIMIT = "OK-LIMIT"

	CONTENT_TYPE = "Content-Type"
	ACCEPT       = "Accept"
	COOKIE       = "Cookie"
	LOCALE       = "locale="

	APPLICATION_JSON      = "application/json"
	APPLICATION_JSON_UTF8 = "application/json; charset=UTF-8"

	/*
	  i18n: internationalization
	*/
	ENGLISH            = "en_US"
	SIMPLIFIED_CHINESE = "zh_CN"
	//zh_TW || zh_HK
	TRADITIONAL_CHINESE = "zh_HK"

	/*
	  http methods
	*/
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"

	/*
	 others
	*/
	ResultDataJsonString = "resultDataJsonString"
	ResultPageJsonString = "resultPageJsonString"

	BTC_USD_SWAP = "BTC-USD-SWAP"
	LTC_USD_SWAP = "LTC-USD-SWAP"
	ETH_USD_SWAP = "ETH-USD-SWAP"
	ETC_USD_SWAP = "ETC-USD-SWAP"
	BCH_USD_SWAP = "BCH-USD-SWAP"
	BSV_USD_SWAP = "BSV-USD-SWAP"
	EOS_USD_SWAP = "EOS-USD-SWAP"
	XRP_USD_SWAP = "XRP-USD-SWAP"

	/*Rest Endpoint*/
	Endpoint              = "https://www.okex.com"
	GET_ACCOUNTS          = "/api/swap/v3/accounts"
	PLACE_ORDER           = "/api/swap/v3/order"
	CANCEL_ORDER          = "/api/swap/v3/cancel_order/%s/%s"
	GET_ORDER             = "/api/swap/v3/orders/%s/%s"
	GET_POSITION          = "/api/swap/v3/%s/position"
	GET_DEPTH             = "/api/swap/v3/instruments/%s/depth?size=%d"
	GET_TICKER            = "/api/swap/v3/instruments/%s/ticker"
	GET_ALL_TICKER        = "/api/swap/v3/instruments/ticker"
	GET_UNFINISHED_ORDERS = "/api/swap/v3/orders/%s?status=%d&limit=%d"
	PLACE_ALGO_ORDER      = "/api/swap/v3/order_algo"
	CANCEL_ALGO_ORDER     = "/api/swap/v3/cancel_algos"
	GET_ALGO_ORDER        = "/api/swap/v3/order_algo/%s?order_type=%d&"
)

type BaseResponse struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Result       bool   `json:"result,string"`
}

type OKExSwap struct {
	*OKEx
	config *APIConfig
}

func NewOKExSwap(config *APIConfig) *OKExSwap {
	return &OKExSwap{OKEx: &OKEx{config: config}, config: config}
}

func (ok *OKExSwap) GetExchangeName() string {
	return OKEX_SWAP
}

func (ok *OKExSwap) GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error) {
	var resp struct {
		InstrumentId string  `json:"instrument_id"`
		Last         float64 `json:"last,string"`
		High24h      float64 `json:"high_24h,string"`
		Low24h       float64 `json:"low_24h,string"`
		BestBid      float64 `json:"best_bid,string"`
		BestAsk      float64 `json:"best_ask,string"`
		Volume24h    float64 `json:"volume_24h,string"`
		Timestamp    string  `json:"timestamp"`
	}
	contractType = ok.adaptContractType(currencyPair)
	err := ok.DoRequest("GET", fmt.Sprintf(GET_TICKER, contractType), "", &resp)
	if err != nil {
		return nil, err
	}

	date, _ := time.Parse(time.RFC3339, resp.Timestamp)
	return &Ticker{
		Pair: currencyPair,
		Last: resp.Last,
		Low:  resp.Low24h,
		High: resp.High24h,
		Vol:  resp.Volume24h,
		Buy:  resp.BestBid,
		Sell: resp.BestAsk,
		Date: uint64(date.UnixNano() / int64(time.Millisecond))}, nil
}

func (ok *OKExSwap) GetFutureAllTicker() (*[]FutureTicker, error) {
	var resp SwapTickerList
	err := ok.DoRequest("GET", GET_ALL_TICKER, "", &resp)
	if err != nil {
		return nil, err
	}

	var tickers []FutureTicker
	for _, t := range resp {
		date, _ := time.Parse(time.RFC3339, t.Timestamp)
		tickers = append(tickers, FutureTicker{
			ContractType: t.InstrumentId,
			Ticker: &Ticker{
				Pair: NewCurrencyPair3(t.InstrumentId, "-"),
				Sell: t.BestAsk,
				Buy:  t.BestBid,
				Low:  t.Low24h,
				High: t.High24h,
				Last: t.Last,
				Vol:  t.Volume24h,
				Date: uint64(date.UnixNano() / int64(time.Millisecond))}})
	}

	return &tickers, nil
}

func (ok *OKExSwap) GetFutureDepth(currencyPair CurrencyPair, contractType string, size int) (*Depth, error) {
	var resp SwapInstrumentDepth
	contractType = ok.adaptContractType(currencyPair)

	err := ok.DoRequest("GET", fmt.Sprintf(GET_DEPTH, contractType, size), "", &resp)
	if err != nil {
		return nil, err
	}

	var dep Depth
	dep.ContractType = contractType
	dep.Pair = currencyPair
	dep.UTime, _ = time.Parse(time.RFC3339, resp.Timestamp)

	for _, v := range resp.Bids {
		dep.BidList = append(dep.BidList, DepthRecord{
			Price:  ToFloat64(v[0]),
			Amount: ToFloat64(v[1])})
	}

	for i := len(resp.Asks) - 1; i >= 0; i-- {
		dep.AskList = append(dep.AskList, DepthRecord{
			Price:  ToFloat64(resp.Asks[i][0]),
			Amount: ToFloat64(resp.Asks[i][1])})
	}

	return &dep, nil
}

func (ok *OKExSwap) GetFutureUserinfo(currencyPair ...CurrencyPair) (*FutureAccount, error) {
	var infos SwapAccounts

	err := ok.OKEx.DoRequest("GET", GET_ACCOUNTS, "", &infos)
	if err != nil {
		return nil, err
	}

	//log.Println(infos)
	acc := FutureAccount{}
	acc.FutureSubAccounts = make(map[Currency]FutureSubAccount, 2)

	for _, account := range infos.Info {
		subAcc := FutureSubAccount{AccountRights: account.Equity,
			KeepDeposit: account.Margin, ProfitReal: account.RealizedPnl, ProfitUnreal: account.UnrealizedPnl, RiskRate: account.MarginRatio}
		switch account.InstrumentId {
		case BTC_USD_SWAP:
			subAcc.Currency = BTC
		case LTC_USD_SWAP:
			subAcc.Currency = LTC
		case ETH_USD_SWAP:
			subAcc.Currency = ETH
		case ETC_USD_SWAP:
			subAcc.Currency = ETC
		case BCH_USD_SWAP:
			subAcc.Currency = BCH
		case BSV_USD_SWAP:
			subAcc.Currency = BSV
		case EOS_USD_SWAP:
			subAcc.Currency = EOS
		case XRP_USD_SWAP:
			subAcc.Currency = XRP
		default:
			subAcc.Currency = UNKNOWN
		}
		acc.FutureSubAccounts[subAcc.Currency] = subAcc
	}

	return &acc, nil
}

type AccountInfo struct {
	Info struct {
		Currency          string  `json:"currency"`
		Equity            float64 `json:"equity,string"`
		FixedBalance      float64 `json:"fixed_balance,string"`
		InstrumentID      string  `json:"instrument_id"`
		MaintMarginRatio  float64 `json:"maint_margin_ratio,string"`
		Margin            float64 `json:"margin,string"`
		MarginFrozen      float64 `json:"margin_frozen,string"`
		MarginMode        string  `json:"margin_mode"`
		MarginRatio       float64 `json:"margin_ratio,string"`
		MaxWithdraw       float64 `json:"max_withdraw,string"`
		RealizedPnl       float64 `json:"realized_pnl,string"`
		Timestamp         string  `json:"timestamp"`
		TotalAvailBalance float64 `json:"total_avail_balance,string"`
		Underlying        string  `json:"underlying"`
		UnrealizedPnl     float64 `json:"unrealized_pnl,string"`
	} `json:"info"`
}

func (ok *OKExSwap) GetFutureAccountInfo(currency CurrencyPair) (*AccountInfo, error) {
	var infos AccountInfo

	err := ok.OKEx.DoRequest("GET", fmt.Sprintf("/api/swap/v3/%s/accounts", ok.adaptContractType(currency)), "", &infos)
	if err != nil {
		return nil, err
	}
	return &infos, nil
}

/*
 OKEX swap api parameter's definition
 @author Lingting Fu
 @date 2018-12-27
 @version 1.0.0
*/

type BasePlaceOrderInfo struct {
	ClientOid  string `json:"client_oid"`
	Price      string `json:"price"`
	MatchPrice string `json:"match_price"`
	Type       string `json:"type"`
	Size       string `json:"size"`
	OrderType  string `json:"order_type"`
}

type PlaceOrderInfo struct {
	BasePlaceOrderInfo
	InstrumentId string `json:"instrument_id"`
}

type PlaceOrdersInfo struct {
	InstrumentId string                `json:"instrument_id"`
	OrderData    []*BasePlaceOrderInfo `json:"order_data"`
}

func (ok *OKExSwap) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice int, leverRate float64) (string, error) {
	fOrder, err := ok.PlaceFutureOrder2(currencyPair, contractType, price, amount, openType, matchPrice)
	return fOrder.OrderID2, err
}

func (ok *OKExSwap) PlaceFutureOrder2(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice int, opt ...LimitOrderOptionalParameter) (*FutureOrder, error) {
	cid := GenerateOrderClientId(32)
	param := PlaceOrderInfo{
		BasePlaceOrderInfo{
			ClientOid:  cid,
			Price:      price,
			MatchPrice: fmt.Sprint(matchPrice),
			Type:       fmt.Sprint(openType),
			Size:       amount},
		ok.adaptContractType(currencyPair),
	}

	if len(opt) > 0 {
		switch opt[0] {
		case PostOnly:
			param.OrderType = "1"
		case Fok:
			param.OrderType = "2"
		case Ioc:
			param.OrderType = "3"
		}
	}

	reqBody, _, _ := ok.OKEx.BuildRequestBody(param)

	fOrder := &FutureOrder{
		ClientOid:    cid,
		Currency:     currencyPair,
		ContractName: contractType,
		OType:        openType,
		Price:        ToFloat64(price),
		Amount:       ToFloat64(amount),
	}

	var resp struct {
		BaseResponse
		OrderID   string `json:"order_id"`
		ClientOid string `json:"client_oid"`
	}

	err := ok.DoRequest("POST", PLACE_ORDER, reqBody, &resp)
	if err != nil {
		return fOrder, err
	}

	if resp.ErrorMessage != "" {
		return fOrder, errors.New(fmt.Sprintf("%s:%s", resp.ErrorCode, resp.ErrorMessage))
	}

	fOrder.OrderID2 = resp.OrderID

	return fOrder, nil
}

func (ok *OKExSwap) LimitFuturesOrder(currencyPair CurrencyPair, contractType, price, amount string, openType int, opt ...LimitOrderOptionalParameter) (*FutureOrder, error) {
	return ok.PlaceFutureOrder2(currencyPair, contractType, price, amount, openType, 0, opt...)
}

func (ok *OKExSwap) MarketFuturesOrder(currencyPair CurrencyPair, contractType, amount string, openType int) (*FutureOrder, error) {
	return ok.PlaceFutureOrder2(currencyPair, contractType, "0", amount, openType, 1)
}

func (ok *OKExSwap) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error) {
	var cancelParam struct {
		OrderId      string `json:"order_id"`
		InstrumentId string `json:"instrument_id"`
	}

	var resp SwapCancelOrderResult

	cancelParam.InstrumentId = contractType
	cancelParam.OrderId = orderId

	//req, _, _ := BuildRequestBody(cancelParam)

	err := ok.DoRequest("POST", fmt.Sprintf(CANCEL_ORDER, ok.adaptContractType(currencyPair), orderId), "", &resp)
	if err != nil {
		return false, err
	}

	return resp.Result, nil
}

func (ok *OKExSwap) parseOrder(ord BaseOrderInfo) FutureOrder {
	oTime, _ := time.Parse(time.RFC3339, ord.Timestamp)
	return FutureOrder{
		ClientOid:  ord.ClientOid,
		OrderID2:   ord.OrderId,
		Amount:     ord.Size,
		Price:      ord.Price,
		DealAmount: ord.FilledQty,
		AvgPrice:   ord.PriceAvg,
		OType:      ord.Type,
		Status:     ok.AdaptTradeStatus(ord.Status),
		Fee:        ord.Fee,
		OrderTime:  oTime.UnixNano() / int64(time.Millisecond)}
}

func (ok *OKExSwap) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	var (
		resp SwapOrdersInfo
	)
	contractType = ok.adaptContractType(currencyPair)
	err := ok.DoRequest("GET", fmt.Sprintf(GET_UNFINISHED_ORDERS, contractType, 6, 100), "", &resp)
	if err != nil {
		return nil, err
	}

	if resp.Message != "" {
		return nil, errors.New(fmt.Sprintf("{\"ErrCode\":%d,\"ErrMessage\":\"%s\"", resp.Code, resp.Message))
	}

	var orders []FutureOrder
	for _, info := range resp.OrderInfo {
		ord := ok.parseOrder(info)
		ord.Currency = currencyPair
		ord.ContractName = contractType
		orders = append(orders, ord)
	}

	//log.Println(len(orders))
	return orders, nil
}

/**
 *获取订单信息
 */
func (ok *OKExSwap) GetFutureOrders(orderIds []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	panic("")
}

/**
 *获取单个订单信息
 */
func (ok *OKExSwap) GetFutureOrder(orderId string, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	var getOrderParam struct {
		OrderId      string `json:"order_id"`
		InstrumentId string `json:"instrument_id"`
	}

	var resp struct {
		BizWarmTips
		BaseOrderInfo
	}

	contractType = ok.adaptContractType(currencyPair)

	getOrderParam.OrderId = orderId
	getOrderParam.InstrumentId = contractType

	//reqBody, _, _ := BuildRequestBody(getOrderParam)

	err := ok.DoRequest("GET", fmt.Sprintf(GET_ORDER, contractType, orderId), "", &resp)
	if err != nil {
		return nil, err
	}

	if resp.Message != "" {
		return nil, errors.New(fmt.Sprintf("{\"ErrCode\":%d,\"ErrMessage\":\"%s\"}", resp.Code, resp.Message))
	}

	oTime, err := time.Parse(time.RFC3339, resp.Timestamp)

	return &FutureOrder{
		ClientOid:    resp.ClientOid,
		Currency:     currencyPair,
		ContractName: contractType,
		OrderID2:     resp.OrderId,
		Amount:       resp.Size,
		Price:        resp.Price,
		DealAmount:   resp.FilledQty,
		AvgPrice:     resp.PriceAvg,
		OType:        resp.Type,
		Fee:          resp.Fee,
		Status:       ok.AdaptTradeStatus(resp.Status),
		OrderTime:    oTime.UnixNano() / int64(time.Millisecond),
	}, nil
}

func (ok *OKExSwap) GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	var resp SwapPosition
	contractType = ok.adaptContractType(currencyPair)

	err := ok.DoRequest("GET", fmt.Sprintf(GET_POSITION, contractType), "", &resp)
	if err != nil {
		return nil, err
	}

	var positions []FuturePosition

	positions = append(positions, FuturePosition{
		ContractType: contractType,
		Symbol:       currencyPair})

	var (
		buyPosition  SwapPositionHolding
		sellPosition SwapPositionHolding
	)

	if len(resp.Holding) > 0 {
		if resp.Holding[0].Side == "long" {
			buyPosition = resp.Holding[0]
			if len(resp.Holding) == 2 {
				sellPosition = resp.Holding[1]
			}
		} else {
			sellPosition = resp.Holding[0]
			if len(resp.Holding) == 2 {
				buyPosition = resp.Holding[1]
			}
		}

		positions[0].ForceLiquPrice = buyPosition.LiquidationPrice
		positions[0].BuyAmount = buyPosition.Position
		positions[0].BuyAvailable = buyPosition.AvailPosition
		positions[0].BuyPriceAvg = buyPosition.AvgCost
		positions[0].BuyProfitReal = buyPosition.RealizedPnl
		positions[0].BuyPriceCost = buyPosition.SettlementPrice

		positions[0].ForceLiquPrice = sellPosition.LiquidationPrice
		positions[0].SellAmount = sellPosition.Position
		positions[0].SellAvailable = sellPosition.AvailPosition
		positions[0].SellPriceAvg = sellPosition.AvgCost
		positions[0].SellProfitReal = sellPosition.RealizedPnl
		positions[0].SellPriceCost = sellPosition.SettlementPrice

		positions[0].LeverRate = ToFloat64(sellPosition.Leverage)
	}
	return positions, nil
}

/**
 * BTC: 100美元一张合约
 * LTC/ETH/ETC/BCH: 10美元一张合约
 */
func (ok *OKExSwap) GetContractValue(currencyPair CurrencyPair) (float64, error) {
	if currencyPair.CurrencyA.Eq(BTC) {
		return 100, nil
	}
	return 10, nil
}

func (ok *OKExSwap) GetFee() (float64, error) {
	panic("not support")
}

func (ok *OKExSwap) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	panic("not support")
}

func (ok *OKExSwap) GetFutureIndex(currencyPair CurrencyPair) (float64, error) {
	panic("not support")
}

func (ok *OKExSwap) GetDeliveryTime() (int, int, int, int) {
	panic("not support")
}

func (ok *OKExSwap) GetKlineRecords(contractType string, currency CurrencyPair, period, size, since int) ([]FutureKline, error) {

	sinceTime := time.Unix(int64(since), 0).UTC()

	if since/int(time.Second) != 1 { //如果不为秒，转为秒
		sinceTime = time.Unix(int64(since)/int64(time.Second), 0).UTC()
	}

	granularity := adaptKLinePeriod(KlinePeriod(period))
	if granularity == -1 {
		return nil, errors.New("kline period parameter is error")
	}
	return ok.GetKlineRecords2(contractType, currency, sinceTime.Format(time.RFC3339), "", strconv.Itoa(granularity))
}

/**
  since : 单位秒,开始时间
*/
func (ok *OKExSwap) GetKlineRecords2(contractType string, currency CurrencyPair, start, end, period string) ([]FutureKline, error) {
	urlPath := "/api/swap/v3/instruments/%s/candles?%s"
	params := url.Values{}
	if start != "" {
		params.Set("start", start)
	}
	if end != "" {
		params.Set("end", end)
	}
	if period != "" {
		params.Set("granularity", period)
	}
	contractId := ok.adaptContractType(currency)

	var response [][]interface{}
	err := ok.DoRequest("GET", fmt.Sprintf(urlPath, contractId, params.Encode()), "", &response)
	if err != nil {
		return nil, err
	}

	var kline []FutureKline
	for _, itm := range response {
		t, _ := time.Parse(time.RFC3339, fmt.Sprint(itm[0]))
		kline = append(kline, FutureKline{
			Kline: &Kline{
				Timestamp: t.Unix(),
				Pair:      currency,
				Open:      ToFloat64(itm[1]),
				High:      ToFloat64(itm[2]),
				Low:       ToFloat64(itm[3]),
				Close:     ToFloat64(itm[4]),
				Vol:       ToFloat64(itm[5])},
			Vol2: ToFloat64(itm[6])})
	}

	return kline, nil
}

func (ok *OKExSwap) GetTrades(contractType string, currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not support")
}

func (ok *OKExSwap) GetExchangeRate() (float64, error) {
	panic("not support")
}

func (ok *OKExSwap) GetHistoricalFunding(contractType string, currencyPair CurrencyPair, page int) ([]HistoricalFunding, error) {
	var resp []HistoricalFunding
	uri := fmt.Sprintf("/api/swap/v3/instruments/%s/historical_funding_rate?from=%d", ok.adaptContractType(currencyPair), page)
	err := ok.DoRequest("GET", uri, "", &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (ok *OKExSwap) AdaptTradeStatus(status int) TradeStatus {
	switch status {
	case -1:
		return ORDER_CANCEL
	case 0:
		return ORDER_UNFINISH
	case 1:
		return ORDER_PART_FINISH
	case 2:
		return ORDER_FINISH
	default:
		return ORDER_UNFINISH
	}
}

func (ok *OKExSwap) adaptContractType(currencyPair CurrencyPair) string {
	return fmt.Sprintf("%s-SWAP", currencyPair.ToSymbol("-"))
}

type Instrument struct {
	InstrumentID        string    `json:"instrument_id"`
	UnderlyingIndex     string    `json:"underlying_index"`
	QuoteCurrency       string    `json:"quote_currency"`
	Coin                string    `json:"coin"`
	ContractVal         float64   `json:"contract_val,string"`
	Listing             time.Time `json:"listing"`
	Delivery            time.Time `json:"delivery"`
	SizeIncrement       int       `json:"size_increment,string"`
	TickSize            float64   `json:"tick_size,string"`
	BaseCurrency        string    `json:"base_currency"`
	Underlying          string    `json:"underlying"`
	SettlementCurrency  string    `json:"settlement_currency"`
	IsInverse           bool      `json:"is_inverse,string"`
	ContractValCurrency string    `json:"contract_val_currency"`
}

func (ok *OKExSwap) GetInstruments() ([]Instrument, error) {
	var resp []Instrument
	err := ok.DoRequest("GET", "/api/swap/v3/instruments", "", &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type MarginLeverage struct {
	LongLeverage  float64 `json:"long_leverage,string"`
	MarginMode    string  `json:"margin_mode"`
	ShortLeverage float64 `json:"short_leverage,string"`
	InstrumentId  string  `json:"instrument_id"`
}

func (ok *OKExSwap) GetMarginLevel(currencyPair CurrencyPair) (*MarginLeverage, error) {
	var resp MarginLeverage
	uri := fmt.Sprintf("/api/swap/v3/accounts/%s/settings", ok.adaptContractType(currencyPair))

	err := ok.DoRequest("GET", uri, "", &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil

}

// marginmode
//1:逐仓-多仓
//2:逐仓-空仓
//3:全仓
func (ok *OKExSwap) SetMarginLevel(currencyPair CurrencyPair, level, marginMode int) (*MarginLeverage, error) {
	var resp MarginLeverage
	uri := fmt.Sprintf("/api/swap/v3/accounts/%s/leverage", ok.adaptContractType(currencyPair))

	reqBody := make(map[string]string)
	reqBody["leverage"] = strconv.Itoa(level)
	reqBody["side"] = strconv.Itoa(marginMode)
	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	err = ok.DoRequest("POST", uri, string(data), &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

//委托策略下单 algo_type 1:限价 2:市场价；触发价格类型，默认是限价；为市场价时，委托价格不必填；
func (ok *OKExSwap) PlaceFutureAlgoOrder(ord *FutureOrder) (*FutureOrder, error) {
	var param struct {
		InstrumentId string `json:"instrument_id"`
		Type         int    `json:"type"`
		OrderType    int    `json:"order_type"` //1：止盈止损 2：跟踪委托 3：冰山委托 4：时间加权
		Size         string `json:"size"`
		TriggerPrice string `json:"trigger_price"`
		AlgoPrice    string `json:"algo_price"`
		AlgoType     string `json:"algo_type"`
	}

	var response struct {
		ErrorMessage string `json:"error_message"`
		ErrorCode    string `json:"error_code"`
		DetailMsg    string `json:"detail_msg"`

		Data struct {
			Result       string `json:"result"`
			ErrorMessage string `json:"error_message"`
			ErrorCode    string `json:"error_code"`
			AlgoId       string `json:"algo_id"`
			InstrumentId string `json:"instrument_id"`
			OrderType    int    `json:"order_type"`
		} `json:"data"`
	}
	if ord == nil {
		return nil, errors.New("ord param is nil")
	}
	param.InstrumentId = ok.adaptContractType(ord.Currency)
	param.Type = ord.OType
	param.OrderType = ord.OrderType
	param.AlgoType = fmt.Sprint(ord.AlgoType)
	param.TriggerPrice = fmt.Sprint(ord.TriggerPrice)
	param.AlgoPrice = fmt.Sprint(ToFloat64(ord.Price))
	param.Size = fmt.Sprint(ord.Amount)

	reqBody, _, _ := ok.BuildRequestBody(param)
	err := ok.DoRequest("POST", PLACE_ALGO_ORDER, reqBody, &response)

	if err != nil {
		return ord, err
	}

	ord.OrderID2 = response.Data.AlgoId
	ord.OrderTime = time.Now().UnixNano() / int64(time.Millisecond)

	return ord, nil
}

//委托策略撤单
func (ok *OKExSwap) FutureCancelAlgoOrder(currencyPair CurrencyPair, orderId []string) (bool, error) {
	if len(orderId) == 0 {
		return false, errors.New("invalid order id")
	}
	var cancelParam struct {
		InstrumentId string   `json:"instrument_id"`
		AlgoIds      []string `json:"algo_ids"`
		OrderType    string   `json:"order_type"`
	}

	var resp struct {
		ErrorMessage string `json:"error_message"`
		ErrorCode    string `json:"error_code"`
		DetailMsg    string `json:"detailMsg"`
		Data         struct {
			Result       string `json:"result"`
			AlgoIds      string `json:"algo_ids"`
			InstrumentID string `json:"instrument_id"`
			OrderType    string `json:"order_type"`
		} `json:"data"`
	}

	cancelParam.InstrumentId = ok.adaptContractType(currencyPair)
	cancelParam.OrderType = "1"
	cancelParam.AlgoIds = orderId

	reqBody, _, _ := ok.BuildRequestBody(cancelParam)

	err := ok.DoRequest("POST", CANCEL_ALGO_ORDER, reqBody, &resp)
	if err != nil {
		return false, err
	}

	return resp.Data.Result == "success", nil
}

//获取委托单列表, status和algo_id必填且只能填其一
func (ok *OKExSwap) GetFutureAlgoOrders(algo_id string, status string, currencyPair CurrencyPair) ([]FutureOrder, error) {
	uri := fmt.Sprintf(GET_ALGO_ORDER, ok.adaptContractType(currencyPair), 1)
	if algo_id != "" {
		uri += "algo_id=" + algo_id
	} else if status != "" {
		uri += "status=" + status
	} else {
		return nil, errors.New("status or algo_id is needed")
	}

	var resp struct {
		OrderStrategyVOS []struct {
			AlgoId       string `json:"algo_id"`
			AlgoPrice    string `json:"algo_price"`
			InstrumentId string `json:"instrument_id"`
			Leverage     string `json:"leverage"`
			OrderType    string `json:"order_type"`
			RealAmount   string `json:"real_amount"`
			RealPrice    string `json:"real_price"`
			Size         string `json:"size"`
			Status       string `json:"status"`
			Timestamp    string `json:"timestamp"`
			TriggerPrice string `json:"trigger_price"`
			Type         string `json:"type"`
		} `json:"orderStrategyVOS"`
	}

	err := ok.DoRequest("GET", uri, "", &resp)
	if err != nil {
		return nil, err
	}

	var orders []FutureOrder
	for _, info := range resp.OrderStrategyVOS {
		oTime, _ := time.Parse(time.RFC3339, info.Timestamp)

		ord := FutureOrder{
			OrderID2:     info.AlgoId,
			Price:        ToFloat64(info.AlgoPrice),
			Amount:       ToFloat64(info.Size),
			AvgPrice:     ToFloat64(info.RealPrice),
			DealAmount:   ToFloat64(info.RealAmount),
			OrderTime:    oTime.UnixNano() / int64(time.Millisecond),
			Status:       ok.AdaptTradeStatus(ToInt(info.Status)),
			Currency:     CurrencyPair{},
			OrderType:    ToInt(info.OrderType),
			OType:        ToInt(info.Type),
			TriggerPrice: ToFloat64(info.TriggerPrice),
		}
		ord.Currency = currencyPair
		orders = append(orders, ord)
	}

	//log.Println(len(orders))
	return orders, nil
}
