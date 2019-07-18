package okex

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	. "github.com/nntaoli-project/GoEx"
	"sort"
	"strings"
	"sync"
	"time"
)

//合约信息
type FutureContractInfo struct {
	InstrumentID    string  `json:"instrument_id"` //合约ID:如BTC-USD-180213
	UnderlyingIndex string  `json:"underlying_index"`
	QuoteCurrency   string  `json:"quote_currency"`
	TickSize        float64 `json:"tick_size,string"` //下单价格精度
	TradeIncrement  string  `json:"trade_increment"`  //数量精度
	ContractVal     string  `json:"contract_val"`     //合约面值(美元)
	Listing         string  `json:"listing"`
	Delivery        string  `json:"delivery"` //交割日期
	Alias           string  `json:"alias"`    //	本周 this_week 次周 next_week 季度 quarter
}

type AllFutureContractInfo struct {
	contractInfos []FutureContractInfo
	uTime         time.Time
}

type OKExFuture struct {
	*OKEx
	sync.Locker
	allContractInfo AllFutureContractInfo
}

func (ok *OKExFuture) GetExchangeName() string {
	return OKEX_FUTURE
}

//获取法币汇率
func (ok *OKExFuture) GetRate() (float64, error) {
	var response struct {
		Rate         float64   `json:"rate,string"`
		InstrumentId string    `json:"instrument_id"` //USD_CNY
		Timestamp    time.Time `json:"timestamp"`
	}
	err := ok.DoRequest("GET", "/api/futures/v3/rate", "", &response)
	if err != nil {
		return 0, err
	}

	return response.Rate, nil
}

func (ok *OKExFuture) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	urlPath := fmt.Sprintf("/api/futures/v3/instruments/%s/estimated_price", ok.getFutureContractId(currencyPair, QUARTER_CONTRACT))
	var response struct {
		InstrumentId    string  `json:"instrument_id"`
		SettlementPrice float64 `json:"settlement_price,string"`
		Timestamp       string  `json:"timestamp"`
	}
	err := ok.DoRequest("GET", urlPath, "", &response)
	if err != nil {
		return 0, err
	}
	return response.SettlementPrice, nil
}

func (ok *OKExFuture) GetFutureContractInfo() ([]FutureContractInfo, error) {
	urlPath := "/api/futures/v3/instruments"
	var response []FutureContractInfo
	err := ok.DoRequest("GET", urlPath, "", &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (ok *OKExFuture) getFutureContractId(pair CurrencyPair, contractAlias string) string {
	if contractAlias != QUARTER_CONTRACT && contractAlias != NEXT_WEEK_CONTRACT && contractAlias != THIS_WEEK_CONTRACT { //传Alias，需要转为具体ContractId
		return contractAlias
	}

	now := time.Now()
	hour := now.Hour()
	mintue := now.Minute()

	if ok.allContractInfo.uTime.IsZero() || (hour == 16 && mintue <= 11) {
		ok.Lock()
		defer ok.Unlock()

		contractInfo, err := ok.GetFutureContractInfo()
		if err == nil {
			ok.allContractInfo.uTime = time.Now()
			ok.allContractInfo.contractInfos = contractInfo
		} else {
			panic(err)
		}
	}

	contractId := ""
	for _, itm := range ok.allContractInfo.contractInfos {
		if itm.Alias == contractAlias && itm.UnderlyingIndex == pair.CurrencyA.Symbol && itm.QuoteCurrency == pair.CurrencyB.Symbol {
			contractId = itm.InstrumentID
			break
		}
	}

	return contractId
}

func (ok *OKExFuture) GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error) {
	var urlPath = fmt.Sprintf("/api/futures/v3/instruments/%s/ticker", ok.getFutureContractId(currencyPair, contractType))

	var response struct {
		InstrumentId string  `json:"instrument_id"`
		Last         float64 `json:"last,string"`
		High24h      float64 `json:"high_24h,string"`
		Low24h       float64 `json:"low_24h,string"`
		BestBid      float64 `json:"best_bid,string"`
		BestAsk      float64 `json:"best_ask,string"`
		Volume24h    float64 `json:"volume_24h,string"`
		Timestamp    string  `json:"timestamp"`
	}
	err := ok.DoRequest("GET", urlPath, "", &response)
	if err != nil {
		return nil, err
	}

	date, _ := time.Parse(time.RFC3339, response.Timestamp)

	return &Ticker{
		Pair: currencyPair,
		Sell: response.BestAsk,
		Buy:  response.BestBid,
		Low:  response.Low24h,
		High: response.High24h,
		Last: response.Last,
		Vol:  response.Volume24h,
		Date: uint64(date.UnixNano() / int64(time.Millisecond))}, nil
}

func (ok *OKExFuture) GetFutureDepth(currencyPair CurrencyPair, contractType string, size int) (*Depth, error) {
	urlPath := fmt.Sprintf("/api/futures/v3/instruments/%s/book?size=%d", ok.getFutureContractId(currencyPair, contractType), size)
	var response struct {
		Bids      [][4]interface{} `json:"bids"`
		Asks      [][4]interface{} `json:"asks"`
		Timestamp string           `json:"timestamp"`
	}
	err := ok.DoRequest("GET", urlPath, "", &response)
	if err != nil {
		return nil, err
	}

	var dep Depth
	dep.Pair = currencyPair
	dep.ContractType = contractType
	dep.UTime, _ = time.Parse(time.RFC3339, response.Timestamp)
	for _, itm := range response.Asks {
		dep.AskList = append(dep.AskList, DepthRecord{
			Price:  ToFloat64(itm[0]),
			Amount: ToFloat64(itm[1])})
	}
	for _, itm := range response.Bids {
		dep.BidList = append(dep.BidList, DepthRecord{
			Price:  ToFloat64(itm[0]),
			Amount: ToFloat64(itm[1])})
	}

	sort.Sort(sort.Reverse(dep.AskList))
	return &dep, nil
}

func (ok *OKExFuture) GetFutureIndex(currencyPair CurrencyPair) (float64, error) {
	//统一交易对，当周，次周，季度指数一样的
	urlPath := fmt.Sprintf("/api/futures/v3/instruments/%s/index", ok.getFutureContractId(currencyPair, QUARTER_CONTRACT))
	var response struct {
		InstrumentId string  `json:"instrument_id"`
		Index        float64 `json:"index,string"`
		Timestamp    string  `json:"timestamp"`
	}
	err := ok.DoRequest("GET", urlPath, "", &response)
	if err != nil {
		return 0, nil
	}
	return response.Index, nil
}

type CrossedAccountInfo struct {
	MarginMode string  `json:"margin_mode"`
	Equity     float64 `json:"equity,string"`
}

func (ok *OKExFuture) GetFutureUserinfo() (*FutureAccount, error) {
	urlPath := "/api/futures/v3/accounts"
	var response struct {
		Info map[string]map[string]interface{}
	}

	err := ok.DoRequest("GET", urlPath, "", &response)
	if err != nil {
		return nil, err
	}

	acc := new(FutureAccount)
	acc.FutureSubAccounts = make(map[Currency]FutureSubAccount, 2)
	for c, info := range response.Info {
		if info["margin_mode"] == "crossed" {
			currency := NewCurrency(c, "")
			acc.FutureSubAccounts[currency] = FutureSubAccount{
				Currency:      currency,
				AccountRights: ToFloat64(info["equity"]),
				ProfitReal:    ToFloat64(info["realized_pnl"]),
				ProfitUnreal:  ToFloat64(info["unrealized_pnl"]),
				KeepDeposit:   ToFloat64(info["margin_frozen"]),
				RiskRate:      ToFloat64(info["margin_ratio"]),
			}
		} else {
			//todo 逐仓模式
		}
	}
	return acc, nil
}

func (ok *OKExFuture) normalizePrice(price float64, pair CurrencyPair) string {
	for _, info := range ok.allContractInfo.contractInfos {
		if info.UnderlyingIndex == pair.CurrencyA.Symbol && info.QuoteCurrency == pair.CurrencyB.Symbol {
			var bit = 0
			for info.TickSize < 1 {
				bit++
				info.TickSize *= 10
			}
			return FloatToString(price, bit)
		}
	}
	return FloatToString(price, 2)
}

//matchPrice:是否以对手价下单(0:不是 1:是)，默认为0;当取值为1时,price字段无效，当以对手价下单，order_type只能选择0:普通委托
func (ok *OKExFuture) PlaceFutureOrder2(matchPrice int, ord *FutureOrder) (*FutureOrder, error) {
	urlPath := "/api/futures/v3/order"
	var param struct {
		ClientOid    string `json:"client_oid"`
		InstrumentId string `json:"instrument_id"`
		Type         int    `json:"type"`
		OrderType    int    `json:"order_type"`
		Price        string `json:"price"`
		Size         string `json:"size"`
		MatchPrice   int    `json:"match_price"`
		Leverage     int    `json:"leverage"`
	}

	var response struct {
		Result       bool   `json:"result"`
		ErrorMessage string `json:"error_message"`
		ErrorCode    string `json:"error_code"`
		ClientOid    string `json:"client_oid"`
		OrderId      string `json:"order_id"`
	}
	if ord == nil {
		return nil, errors.New("ord param is nil")
	}
	param.InstrumentId = ok.getFutureContractId(ord.Currency, ord.ContractName)
	param.ClientOid = strings.Replace(uuid.New().String(), "-", "", 32)
	param.Type = ord.OType
	param.OrderType = ord.OrderType
	param.Price = ok.normalizePrice(ord.Price, ord.Currency)
	param.Size = fmt.Sprint(ord.Amount)
	param.Leverage = ord.LeverRate
	param.MatchPrice = matchPrice

	//当matchPrice=1以对手价下单，order_type只能选择0:普通委托
	if param.MatchPrice == 1 {
		println("注意:当matchPrice=1以对手价下单时，order_type只能选择0:普通委托")
		param.OrderType = ORDINARY
	}

	reqBody, _, _ := ok.BuildRequestBody(param)
	err := ok.DoRequest("POST", urlPath, reqBody, &response)

	if err != nil {
		return ord, err
	}

	ord.ClientOid = response.ClientOid
	ord.OrderID2 = response.OrderId
	ord.OrderTime = time.Now().UnixNano() / int64(time.Millisecond)

	return ord, nil
}

//deprecated
func (ok *OKExFuture) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice, leverRate int) (string, error) {
	urlPath := "/api/futures/v3/order"
	var param struct {
		ClientOid    string `json:"client_oid"`
		InstrumentId string `json:"instrument_id"`
		Type         string `json:"type"`
		OrderType    string `json:"order_type"`
		Price        string `json:"price"`
		Size         string `json:"size"`
		MatchPrice   string `json:"match_price"`
		Leverage     string `json:"leverage"`
	}

	var response struct {
		Result       bool   `json:"result"`
		ErrorMessage string `json:"error_message"`
		ErrorCode    string `json:"error_code"`
		ClientOid    string `json:"client_oid"`
		OrderId      string `json:"order_id"`
	}

	param.InstrumentId = ok.getFutureContractId(currencyPair, contractType)
	param.ClientOid = strings.Replace(uuid.New().String(), "-", "", 32)
	param.Type = fmt.Sprint(openType)
	param.OrderType = "0"
	param.Price = price
	param.Size = amount
	param.MatchPrice = fmt.Sprint(matchPrice)
	param.Leverage = fmt.Sprint(leverRate)

	reqBody, _, _ := ok.BuildRequestBody(param)
	err := ok.DoRequest("POST", urlPath, reqBody, &response)
	if err != nil {
		return "", err
	}

	return response.OrderId, nil
}

func (ok *OKExFuture) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error) {
	urlPath := fmt.Sprintf("/api/futures/v3/cancel_order/%s/%s", ok.getFutureContractId(currencyPair, contractType), orderId)
	var response struct {
		Result       bool   `json:"result"`
		OrderId      string `json:"order_id"`
		ClientOid    string `json:"client_oid"`
		InstrumentId string `json:"instrument_id"`
	}
	err := ok.DoRequest("POST", urlPath, "", &response)
	if err != nil {
		return false, err
	}
	return response.Result, nil
}

func (ok *OKExFuture) GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	urlPath := fmt.Sprintf("/api/futures/v3/%s/position", ok.getFutureContractId(currencyPair, contractType))
	var response struct {
		Result     bool   `json:"result"`
		MarginMode string `json:"margin_mode"`
		Holding    []struct {
			InstrumentId         string    `json:"instrument_id"`
			LongQty              float64   `json:"long_qty,string"` //多
			LongAvailQty         float64   `json:"long_avail_qty,string"`
			LongAvgCost          float64   `json:"long_avg_cost,string"`
			LongSettlementPrice  float64   `json:"long_settlement_price,string"`
			LongMargin           float64   `json:"long_margin,string"`
			LongPnl              float64   `json:"long_pnl,string"`
			LongPnlRatio         float64   `json:"long_pnl_ratio,string"`
			LongUnrealisedPnl    float64   `json:"long_unrealised_pnl,string"`
			RealisedPnl          float64   `json:"realised_pnl,string"`
			Leverage             int       `json:"leverage,string"`
			ShortQty             float64   `json:"short_qty,string"`
			ShortAvailQty        float64   `json:"short_avail_qty,string"`
			ShortAvgCost         float64   `json:"short_avg_cost,string"`
			ShortSettlementPrice float64   `json:"short_settlement_price,string"`
			ShortMargin          float64   `json:"short_margin,string"`
			ShortPnl             float64   `json:"short_pnl,string"`
			ShortPnlRatio        float64   `json:"short_pnl_ratio,string"`
			ShortUnrealisedPnl   float64   `json:"short_unrealised_pnl,string"`
			LiquidationPrice     float64   `json:"liquidation_price,string"`
			CreatedAt            time.Time `json:"created_at,string"`
			UpdatedAt            time.Time `json:"updated_at"`
		}
	}
	err := ok.DoRequest("GET", urlPath, "", &response)
	if err != nil {
		return nil, err
	}

	var postions []FuturePosition

	if !response.Result {
		return nil, errors.New("unknown error")
	}

	if response.MarginMode == "fixed" {
		panic("not support the fix future")
	}

	for _, pos := range response.Holding {
		postions = append(postions, FuturePosition{
			Symbol:         currencyPair,
			ContractType:   contractType,
			ContractId:     ToInt64(pos.InstrumentId[8:]),
			BuyAmount:      pos.LongQty,
			BuyAvailable:   pos.LongAvailQty,
			BuyPriceAvg:    pos.LongAvgCost,
			BuyPriceCost:   pos.LongAvgCost,
			BuyProfitReal:  pos.LongPnl,
			SellAmount:     pos.ShortQty,
			SellAvailable:  pos.ShortAvailQty,
			SellPriceAvg:   pos.ShortAvgCost,
			SellPriceCost:  pos.ShortAvgCost,
			SellProfitReal: pos.ShortPnl,
			ForceLiquPrice: pos.LiquidationPrice,
			LeverRate:      pos.Leverage,
			CreateDate:     pos.CreatedAt.Unix(),
		})
	}

	return postions, nil
}

func (ok *OKExFuture) GetFutureOrders(orderIds []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	panic("")
}

type futureOrderResponse struct {
	InstrumentId string    `json:"instrument_id"`
	ClientOid    string    `json:"client_oid"`
	OrderId      string    `json:"order_id"`
	Size         float64   `json:"size,string"`
	Price        float64   `json:"price,string"`
	FilledQty    float64   `json:"filled_qty,string"`
	PriceAvg     float64   `json:"price_avg,string"`
	Fee          float64   `json:"fee,string"`
	Type         int       `json:"type,string"`
	OrderType    int       `json:"order_type,string"`
	Pnl          float64   `json:"pnl,string"`
	Leverage     int       `json:"leverage,string"`
	ContractVal  float64   `json:"contract_val,string"`
	State        int       `json:"state,string"`
	Timestamp    time.Time `json:"timestamp,string"`
}

func (ok *OKExFuture) adaptOrder(response futureOrderResponse) FutureOrder {
	return FutureOrder{
		ContractName: response.InstrumentId,
		OrderID2:     response.OrderId,
		ClientOid:    response.ClientOid,
		Amount:       response.Size,
		Price:        response.Price,
		DealAmount:   response.FilledQty,
		AvgPrice:     response.PriceAvg,
		OType:        response.Type,
		OrderType:    response.OrderType,
		Status:       ok.adaptOrderState(response.State),
		Fee:          response.Fee,
		OrderTime:    response.Timestamp.UnixNano() / int64(time.Millisecond),
	}
}

func (ok *OKExFuture) GetFutureOrder(orderId string, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	urlPath := fmt.Sprintf("/api/futures/v3/orders/%s/%s", ok.getFutureContractId(currencyPair, contractType), orderId)
	var response futureOrderResponse
	err := ok.DoRequest("GET", urlPath, "", &response)
	if err != nil {
		return nil, err
	}

	ord := ok.adaptOrder(response)
	ord.Currency = currencyPair

	return &ord, nil
}

func (ok *OKExFuture) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	urlPath := fmt.Sprintf("/api/futures/v3/orders/%s?state=6&limit=100", ok.getFutureContractId(currencyPair, contractType))
	var response struct {
		Result    bool                  `json:"result"`
		OrderInfo []futureOrderResponse `json:"order_info"`
	}
	err := ok.DoRequest("GET", urlPath, "", &response)
	if err != nil {
		return nil, err
	}
	if !response.Result {
		return nil, errors.New("error")
	}

	var ords []FutureOrder
	for _, itm := range response.OrderInfo {
		ord := ok.adaptOrder(itm)
		ord.Currency = currencyPair
		ords = append(ords, ord)
	}

	return ords, nil
}

func (ok *OKExFuture) GetFee() (float64, error) { panic("") }

func (ok *OKExFuture) GetContractValue(currencyPair CurrencyPair) (float64, error) {
	for _, info := range ok.allContractInfo.contractInfos {
		if info.UnderlyingIndex == currencyPair.CurrencyA.Symbol && info.QuoteCurrency == currencyPair.CurrencyB.Symbol {
			return ToFloat64(info.ContractVal), nil
		}
	}
	return 0, nil
}

func (ok *OKExFuture) GetDeliveryTime() (int, int, int, int) {
	return 4, 16, 0, 0 //星期五，下午4点交割
}

/**
  since : 单位秒,开始时间
*/
func (ok *OKExFuture) GetKlineRecords(contract_type string, currency CurrencyPair, period, size, since int) ([]FutureKline, error) {
	urlPath := "/api/futures/v3/instruments/%s/candles?start=%s&granularity=%d"
	contractId := ok.getFutureContractId(currency, contract_type)
	sinceTime := time.Unix(int64(since), 0).UTC()

	if since/int(time.Second) != 1 { //如果不为秒，转为秒
		sinceTime = time.Unix(int64(since)/int64(time.Second), 0).UTC()
	}

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
	err := ok.DoRequest("GET", fmt.Sprintf(urlPath, contractId, sinceTime.Format(time.RFC3339), granularity), "", &response)
	if err != nil {
		return nil, err
	}

	var klines []FutureKline
	for _, itm := range response {
		t, _ := time.Parse(time.RFC3339, fmt.Sprint(itm[0]))
		klines = append(klines, FutureKline{
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

	return klines, nil
}

func (ok *OKExFuture) GetTrades(contract_type string, currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("")
}

//特殊接口
/*
 市价全平仓
 contract:合约ID
 oType：平仓方向：CLOSE_SELL平空，CLOSE_BUY平多
*/
func (ok *OKExFuture) MarketCloseAllPosition(currency CurrencyPair, contract string, oType int) (bool, error) {
	urlPath := "/api/futures/v3/close_position"
	var response struct {
		InstrumentId string `json:"instrument_id"`
		Result       bool   `json:"result"`
		Message      string `json:"message"`
		Code         int    `json:"code"`
	}

	var param struct {
		InstrumentId string `json:"instrument_id"`
		Direction    string `json:"direction"`
	}

	param.InstrumentId = ok.getFutureContractId(currency, contract)
	if oType == CLOSE_BUY {
		param.Direction = "long"
	} else {
		param.Direction = "short"
	}
	reqBody, _, _ := ok.BuildRequestBody(param)
	err := ok.DoRequest("POST", urlPath, reqBody, &response)
	if err != nil {
		return false, err
	}

	if !response.Result {
		return false, errors.New(response.Message)
	}

	return true, nil
}
