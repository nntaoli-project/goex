package okex

import (
	"fmt"
	"github.com/go-openapi/errors"
	. "github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/internal/logger"
	"net/url"
	"sort"
	"strings"
	"time"
)

type OKExSpot struct {
	*OKEx
}

// [{
//        "frozen":"0",
//        "hold":"0",
//        "id":"9150707",
//        "currency":"BTC",
//        "balance":"0.0049925",
//        "available":"0.0049925",
//        "holds":"0"
//    },
//    ...]

func (ok *OKExSpot) GetAccount() (*Account, error) {
	urlPath := "/api/spot/v3/accounts"
	var response []struct {
		Frozen    float64 `json:"frozen,string"`
		Hold      float64 `json:"hold,string"`
		Currency  string
		Balance   float64 `json:"balance,string"`
		Available float64 `json:"available,string"`
		Holds     float64 `json:"holds,string"`
	}

	err := ok.OKEx.DoRequest("GET", urlPath, "", &response)
	if err != nil {
		return nil, err
	}

	account := &Account{
		SubAccounts: make(map[Currency]SubAccount, 2)}

	for _, itm := range response {
		currency := NewCurrency(itm.Currency, "")
		account.SubAccounts[currency] = SubAccount{
			Currency:     currency,
			ForzenAmount: itm.Hold,
			Amount:       itm.Available,
		}
	}

	return account, nil
}

type PlaceOrderParam struct {
	ClientOid     string  `json:"client_oid"`
	Type          string  `json:"type"`
	Side          string  `json:"side"`
	InstrumentId  string  `json:"instrument_id"`
	OrderType     int     `json:"order_type"`
	Price         float64 `json:"price"`
	Size          float64 `json:"size"`
	Notional      float64 `json:"notional"`
	MarginTrading string  `json:"margin_trading,omitempty"`
}

type PlaceOrderResponse struct {
	OrderId      string `json:"order_id"`
	ClientOid    string `json:"client_oid"`
	Result       bool   `json:"result"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

/**
Must Set Client Oid
*/
func (ok *OKExSpot) BatchPlaceOrders(orders []Order) ([]PlaceOrderResponse, error) {
	var param []PlaceOrderParam
	var response map[string][]PlaceOrderResponse

	for _, ord := range orders {
		param = append(param, PlaceOrderParam{
			InstrumentId: ord.Currency.AdaptUsdToUsdt().ToSymbol("-"),
			ClientOid:    ord.Cid,
			Side:         strings.ToLower(ord.Side.String()),
			Size:         ord.Amount,
			Price:        ord.Price,
			Type:         "limit",
			OrderType:    ord.OrderType})
	}
	reqBody, _, _ := ok.BuildRequestBody(param)
	err := ok.DoRequest("POST", "/api/spot/v3/batch_orders", reqBody, &response)
	if err != nil {
		return nil, err
	}

	var ret []PlaceOrderResponse

	for _, v := range response {
		ret = append(ret, v...)
	}

	return ret, nil
}

func (ok *OKExSpot) PlaceOrder(ty string, ord *Order) (*Order, error) {
	urlPath := "/api/spot/v3/orders"
	param := PlaceOrderParam{
		ClientOid:    GenerateOrderClientId(32),
		InstrumentId: ord.Currency.AdaptUsdToUsdt().ToLower().ToSymbol("-"),
	}

	var response PlaceOrderResponse

	switch ord.Side {
	case BUY, SELL:
		param.Side = strings.ToLower(ord.Side.String())
		param.Price = ord.Price
		param.Size = ord.Amount
	case SELL_MARKET:
		param.Side = "sell"
		param.Size = ord.Amount
	case BUY_MARKET:
		param.Side = "buy"
		param.Notional = ord.Price
	default:
		param.Size = ord.Amount
		param.Price = ord.Price
	}

	switch ty {
	case "limit":
		param.Type = "limit"
	case "market":
		param.Type = "market"
	case "post_only":
		param.OrderType = ORDER_FEATURE_POST_ONLY
	case "fok":
		param.OrderType = ORDER_FEATURE_FOK
	case "ioc":
		param.OrderType = ORDER_FEATURE_IOC
	}

	jsonStr, _, _ := ok.OKEx.BuildRequestBody(param)
	err := ok.OKEx.DoRequest("POST", urlPath, jsonStr, &response)
	if err != nil {
		return nil, err
	}

	if !response.Result {
		return nil, errors.New(int32(ToInt(response.ErrorCode)), response.ErrorMessage)
	}

	ord.Cid = response.ClientOid
	ord.OrderID2 = response.OrderId

	return ord, nil
}

func (ok *OKExSpot) LimitBuy(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	ty := "limit"
	if len(opt) > 0 {
		ty = opt[0].String()
	}
	return ok.PlaceOrder(ty, &Order{
		Price:    ToFloat64(price),
		Amount:   ToFloat64(amount),
		Currency: currency,
		Side:     BUY,
	})
}

func (ok *OKExSpot) LimitSell(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	ty := "limit"
	if len(opt) > 0 {
		ty = opt[0].String()
	}
	return ok.PlaceOrder(ty, &Order{
		Price:    ToFloat64(price),
		Amount:   ToFloat64(amount),
		Currency: currency,
		Side:     SELL,
	})
}

func (ok *OKExSpot) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return ok.PlaceOrder("market", &Order{
		Price:    ToFloat64(price),
		Amount:   ToFloat64(amount),
		Currency: currency,
		Side:     BUY_MARKET,
	})
}

func (ok *OKExSpot) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return ok.PlaceOrder("market", &Order{
		Price:    ToFloat64(price),
		Amount:   ToFloat64(amount),
		Currency: currency,
		Side:     SELL_MARKET,
	})
}

//orderId can set client oid or orderId
func (ok *OKExSpot) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	urlPath := "/api/spot/v3/cancel_orders/" + orderId
	param := struct {
		InstrumentId string `json:"instrument_id"`
	}{currency.AdaptUsdToUsdt().ToLower().ToSymbol("-")}
	reqBody, _, _ := ok.BuildRequestBody(param)
	var response struct {
		ClientOid    string `json:"client_oid"`
		OrderId      string `json:"order_id"`
		Result       bool   `json:"result"`
		ErrorCode    string `json:"error_code"`
		ErrorMessage string `json:"error_message"`
	}
	err := ok.OKEx.DoRequest("POST", urlPath, reqBody, &response)
	if err != nil {
		return false, err
	}
	if response.Result {
		return true, nil
	}
	return false, errors.New(400, fmt.Sprintf("cancel fail, %s", response.ErrorMessage))
}

type OrderResponse struct {
	InstrumentId   string  `json:"instrument_id"`
	ClientOid      string  `json:"client_oid"`
	OrderId        string  `json:"order_id"`
	Price          string  `json:"price,omitempty"`
	Size           float64 `json:"size,string"`
	Notional       string  `json:"notional"`
	Side           string  `json:"side"`
	Type           string  `json:"type"`
	FilledSize     string  `json:"filled_size"`
	FilledNotional string  `json:"filled_notional"`
	PriceAvg       string  `json:"price_avg"`
	State          int     `json:"state,string"`
	Fee            string  `json:"fee"`
	OrderType      int     `json:"order_type,string"`
	Timestamp      string  `json:"timestamp"`
}

func (ok *OKExSpot) adaptOrder(response OrderResponse) *Order {
	ordInfo := &Order{
		Cid:        response.ClientOid,
		OrderID2:   response.OrderId,
		Price:      ToFloat64(response.Price),
		Amount:     response.Size,
		AvgPrice:   ToFloat64(response.PriceAvg),
		DealAmount: ToFloat64(response.FilledSize),
		Status:     ok.adaptOrderState(response.State),
		Fee:        ToFloat64(response.Fee)}

	switch response.Side {
	case "buy":
		if response.Type == "market" {
			ordInfo.Side = BUY_MARKET
			ordInfo.DealAmount = ToFloat64(response.Notional) //成交金额
		} else {
			ordInfo.Side = BUY
		}
	case "sell":
		if response.Type == "market" {
			ordInfo.Side = SELL_MARKET
			ordInfo.DealAmount = ToFloat64(response.Notional) //成交数量
		} else {
			ordInfo.Side = SELL
		}
	}

	date, err := time.Parse(time.RFC3339, response.Timestamp)
	//log.Println(date.Local().UnixNano()/int64(time.Millisecond))
	if err != nil {
		logger.Error("parse timestamp err=", err, ",timestamp=", response.Timestamp)
	} else {
		ordInfo.OrderTime = int(date.UnixNano() / int64(time.Millisecond))
	}

	return ordInfo
}

//orderId can set client oid or orderId
func (ok *OKExSpot) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	urlPath := "/api/spot/v3/orders/" + orderId + "?instrument_id=" + currency.AdaptUsdToUsdt().ToSymbol("-")
	//param := struct {
	//	InstrumentId string `json:"instrument_id"`
	//}{currency.AdaptUsdToUsdt().ToLower().ToSymbol("-")}
	//reqBody, _, _ := ok.BuildRequestBody(param)
	var response OrderResponse
	err := ok.OKEx.DoRequest("GET", urlPath, "", &response)
	if err != nil {
		return nil, err
	}

	ordInfo := ok.adaptOrder(response)
	ordInfo.Currency = currency

	return ordInfo, nil
}

func (ok *OKExSpot) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	urlPath := fmt.Sprintf("/api/spot/v3/orders_pending?instrument_id=%s", currency.AdaptUsdToUsdt().ToSymbol("-"))
	var response []OrderResponse
	err := ok.OKEx.DoRequest("GET", urlPath, "", &response)
	if err != nil {
		return nil, err
	}

	var ords []Order
	for _, itm := range response {
		ord := ok.adaptOrder(itm)
		ord.Currency = currency
		ords = append(ords, *ord)
	}

	return ords, nil
}

func (ok *OKExSpot) GetOrderHistorys(currency CurrencyPair, optional ...OptionalParameter) ([]Order, error) {
	urlPath := "/api/spot/v3/orders"

	param := url.Values{}
	param.Set("instrument_id", currency.AdaptUsdToUsdt().ToSymbol("-"))
	param.Set("state", "7")
	MergeOptionalParameter(&param, optional...)

	urlPath += "?" + param.Encode()

	var response []OrderResponse
	err := ok.OKEx.DoRequest("GET", urlPath, "", &response)
	if err != nil {
		return nil, err
	}

	var orders []Order
	for _, itm := range response {
		ord := ok.adaptOrder(itm)
		ord.Currency = currency
		orders = append(orders, *ord)
	}

	return orders, nil
}

func (ok *OKExSpot) GetExchangeName() string {
	return OKEX
}

type spotTickerResponse struct {
	InstrumentId  string  `json:"instrument_id"`
	Last          float64 `json:"last,string"`
	High24h       float64 `json:"high_24h,string"`
	Low24h        float64 `json:"low_24h,string"`
	BestBid       float64 `json:"best_bid,string"`
	BestAsk       float64 `json:"best_ask,string"`
	BaseVolume24h float64 `json:"base_volume_24h,string"`
	Timestamp     string  `json:"timestamp"`
}

func (ok *OKExSpot) GetTicker(currency CurrencyPair) (*Ticker, error) {
	urlPath := fmt.Sprintf("/api/spot/v3/instruments/%s/ticker", currency.AdaptUsdToUsdt().ToSymbol("-"))
	var response spotTickerResponse
	err := ok.OKEx.DoRequest("GET", urlPath, "", &response)
	if err != nil {
		return nil, err
	}

	date, _ := time.Parse(time.RFC3339, response.Timestamp)
	return &Ticker{
		Pair: currency,
		Last: response.Last,
		High: response.High24h,
		Low:  response.Low24h,
		Sell: response.BestAsk,
		Buy:  response.BestBid,
		Vol:  response.BaseVolume24h,
		Date: uint64(time.Duration(date.UnixNano() / int64(time.Millisecond)))}, nil
}

func (ok *OKExSpot) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	urlPath := fmt.Sprintf("/api/spot/v3/instruments/%s/book?size=%d", currency.AdaptUsdToUsdt().ToSymbol("-"), size)

	var response struct {
		Asks      [][]interface{} `json:"asks"`
		Bids      [][]interface{} `json:"bids"`
		Timestamp string          `json:"timestamp"`
	}

	err := ok.OKEx.DoRequest("GET", urlPath, "", &response)
	if err != nil {
		return nil, err
	}

	dep := new(Depth)
	dep.Pair = currency
	dep.UTime, _ = time.Parse(time.RFC3339, response.Timestamp)

	for _, itm := range response.Asks {
		dep.AskList = append(dep.AskList, DepthRecord{
			Price:  ToFloat64(itm[0]),
			Amount: ToFloat64(itm[1]),
		})
	}

	for _, itm := range response.Bids {
		dep.BidList = append(dep.BidList, DepthRecord{
			Price:  ToFloat64(itm[0]),
			Amount: ToFloat64(itm[1]),
		})
	}

	sort.Sort(sort.Reverse(dep.AskList))

	return dep, nil
}

func (ok *OKExSpot) GetKlineRecords(currency CurrencyPair, period KlinePeriod, size int, optional ...OptionalParameter) ([]Kline, error) {
	urlPath := "/api/spot/v3/instruments/%s/candles?granularity=%d"

	optParam := url.Values{}
	MergeOptionalParameter(&optParam, optional...)
	urlPath += "&" + optParam.Encode()

	granularity := 60
	switch period {
	case KLINE_PERIOD_1MIN:
		granularity = 60
	case KLINE_PERIOD_3MIN:
		granularity = 180
	case KLINE_PERIOD_5MIN:
		granularity = 300
	case KLINE_PERIOD_15MIN:
		granularity = 900
	case KLINE_PERIOD_30MIN:
		granularity = 1800
	case KLINE_PERIOD_1H, KLINE_PERIOD_60MIN:
		granularity = 3600
	case KLINE_PERIOD_2H:
		granularity = 7200
	case KLINE_PERIOD_4H:
		granularity = 14400
	case KLINE_PERIOD_6H:
		granularity = 21600
	case KLINE_PERIOD_1DAY:
		granularity = 86400
	case KLINE_PERIOD_1WEEK:
		granularity = 604800
	default:
		granularity = 1800
	}

	var response [][]interface{}
	err := ok.DoRequest("GET", fmt.Sprintf(urlPath, currency.AdaptUsdToUsdt().ToSymbol("-"), granularity), "", &response)
	if err != nil {
		return nil, err
	}

	var klines []Kline
	for _, itm := range response {
		t, _ := time.Parse(time.RFC3339, fmt.Sprint(itm[0]))
		klines = append(klines, Kline{
			Timestamp: t.Unix(),
			Pair:      currency,
			Open:      ToFloat64(itm[1]),
			High:      ToFloat64(itm[2]),
			Low:       ToFloat64(itm[3]),
			Close:     ToFloat64(itm[4]),
			Vol:       ToFloat64(itm[5])})
	}

	return klines, nil
}

//非个人，整个交易所的交易记录
func (ok *OKExSpot) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	urlPath := fmt.Sprintf("/api/spot/v3/instruments/%s/trades?limit=%d", currencyPair.AdaptUsdToUsdt().ToSymbol("-"), since)

	var response []struct {
		Timestamp string  `json:"timestamp"`
		TradeId   int64   `json:"trade_id,string"`
		Price     float64 `json:"price,string"`
		Size      float64 `json:"size,string"`
		Side      string  `json:"side"`
	}
	err := ok.DoRequest("GET", urlPath, "", &response)
	if err != nil {
		return nil, err
	}

	var trades []Trade
	for _, item := range response {
		t, _ := time.Parse(time.RFC3339, item.Timestamp)
		trades = append(trades, Trade{
			Tid:    item.TradeId,
			Type:   AdaptTradeSide(item.Side),
			Amount: item.Size,
			Price:  item.Price,
			Date:   t.Unix(),
			Pair:   currencyPair,
		})
	}

	return trades, nil
}

type OKExSpotSymbol struct {
	BaseCurrency    string
	QuoteCurrency   string
	PricePrecision  float64
	AmountPrecision float64
	MinAmount       float64
	MinValue        float64
	SymbolPartition string
	Symbol          string
}

func (ok *OKExSpot) GetCurrenciesPrecision() ([]OKExSpotSymbol, error) {
	var response []struct {
		InstrumentId  string  `json:"instrument_id"`
		BaseCurrency  string  `json:"base_currency"`
		QuoteCurrency string  `json:"quote_currency"`
		MinSize       float64 `json:"min_size,string"`
		SizeIncrement string  `json:"size_increment"`
		TickSize      string  `json:"tick_size"`
	}
	err := ok.DoRequest("GET", "/api/spot/v3/instruments", "", &response)
	if err != nil {
		return nil, err
	}

	var Symbols []OKExSpotSymbol
	for _, v := range response {
		var sym OKExSpotSymbol
		sym.BaseCurrency = v.BaseCurrency
		sym.QuoteCurrency = v.QuoteCurrency
		sym.Symbol = v.InstrumentId
		sym.MinAmount = v.MinSize

		pres := strings.Split(v.TickSize, ".")
		if len(pres) == 1 {
			sym.PricePrecision = 0
		} else {
			sym.PricePrecision = float64(len(pres[1]))
		}

		pres = strings.Split(v.SizeIncrement, ".")
		if len(pres) == 1 {
			sym.AmountPrecision = 0
		} else {
			sym.AmountPrecision = float64(len(pres[1]))
		}

		Symbols = append(Symbols, sym)
	}
	return Symbols, nil
}
