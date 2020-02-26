package huobi

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Hbdm struct {
	config *APIConfig
}

type OrderInfo struct {
	Symbol         string  `json:"symbol"`
	ContractType   string  `json:"contract_type"`
	ContractCode   string  `json:"contract_code"`
	Volume         float64 `json:"volume"`
	Price          float64 `json:"price"`
	OrderPriceType string  `json:"order_price_type"`
	Direction      string  `json:"direction"`
	Offset         string  `json:"offset"`
	LeverRate      int     `json:"lever_rate"`
	OrderId        int64   `json:"order_id"`
	ClientOrderId  int64   `json:"client_order_id"`
	OrderSource    string  `json:"order_source"`
	CreatedAt      int64   `json:"created_at"`
	TradeVolume    float64 `json:"trade_volume"`
	TradeTurnover  float64 `json:"trade_turnover"`
	Fee            float64 `json:"fee"`
	TradeAvgPrice  float64 `json:"trade_avg_price"`
	MarginFrozen   float64 `json:"margin_frozen"`
	Status         int     `json:"status"`
}

type BaseResponse struct {
	Status  string          `json:"status"`
	Ch      string          `json:"ch"`
	Ts      int64           `json:"ts"`
	ErrCode int             `json:"err_code"`
	ErrMsg  string          `json:"err_msg"`
	Data    json.RawMessage `json:"data"`
}

const (
	defaultBaseUrl = "https://api.hbdm.com"
)

func NewHbdm(conf *APIConfig) *Hbdm {
	if conf.Endpoint == "" {
		conf.Endpoint = defaultBaseUrl
	}
	return &Hbdm{conf}
}

func (dm *Hbdm) GetExchangeName() string {
	return HBDM
}

func (dm *Hbdm) GetFutureUserinfo() (*FutureAccount, error) {
	path := "/api/v1/contract_account_info"
	var data []struct {
		Symbol            string  `json:"symbol"`
		MarginBalance     float64 `json:"margin_balance"`
		MarginPosition    float64 `json:"margin_position"`
		MarginFrozen      float64 `json:"margin_frozen"`
		MarginAvailable   float64 `json:"margin_available"`
		ProfitReal        float64 `json:"profit_real"`
		ProfitUnreal      float64 `json:"profit_unreal"`
		RiskRate          float64 `json:"risk_rate"`
		LiquidationPrice  float64 `json:"liquidation_price"`
		WithdrawAvailable float64 `json:"withdraw_available"`
		LeverRate         float64 `json:"lever_rate"`
	}

	params := &url.Values{}
	err := dm.doRequest(path, params, &data)
	if err != nil {
		return nil, err
	}

	acc := new(FutureAccount)
	acc.FutureSubAccounts = make(map[Currency]FutureSubAccount, 4)
	for _, sub := range data {
		subAcc := FutureSubAccount{
			Currency:      NewCurrency(sub.Symbol, ""),
			AccountRights: sub.MarginBalance,
			KeepDeposit:   sub.MarginPosition,
			ProfitReal:    sub.ProfitReal,
			ProfitUnreal:  sub.ProfitUnreal,
			RiskRate:      sub.RiskRate}
		acc.FutureSubAccounts[subAcc.Currency] = subAcc
	}

	return acc, nil
}

func (dm *Hbdm) GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	var data []struct {
		Symbol         string  `json:"symbol"`
		ContractCode   string  `json:"contract_code"`
		ContractType   string  `json:"contract_type"`
		Volume         float64 `json:"volume"`
		Available      float64 `json:"available"`
		Frozen         float64 `json:"frozen"`
		CostOpen       float64 `json:"cost_open"`
		CostHold       float64 `json:"cost_hold"`
		ProfitUnreal   float64 `json:"profit_unreal"`
		ProfitRate     float64 `json:"profit_rate"`
		Profit         float64 `json:"profit"`
		PositionMargin float64 `json:"position_margin"`
		LeverRate      int     `json:"lever_rate"`
		Direction      string  `json:"direction"`
	}

	path := "/api/v1/contract_position_info"
	params := &url.Values{}
	params.Add("symbol", currencyPair.CurrencyA.Symbol)

	err := dm.doRequest(path, params, &data)
	if err != nil {
		return nil, err
	}

	//	log.Println(data)

	var positions []FuturePosition
	for _, d := range data {
		switch d.Direction {
		case "buy":
			positions = append(positions, FuturePosition{
				ContractType:  contractType,
				ContractId:    int64(ToInt(d.ContractCode[3:])),
				Symbol:        currencyPair,
				BuyAmount:     d.Volume,
				BuyAvailable:  d.Available,
				BuyPriceAvg:   d.CostOpen,
				BuyPriceCost:  d.CostHold,
				BuyProfitReal: d.ProfitRate,
				LeverRate:     d.LeverRate})
		case "sell":
			positions = append(positions, FuturePosition{
				ContractType:   contractType,
				ContractId:     int64(ToInt(d.ContractCode[3:])),
				Symbol:         currencyPair,
				SellAmount:     d.Volume,
				SellAvailable:  d.Available,
				SellPriceAvg:   d.CostOpen,
				SellPriceCost:  d.CostHold,
				SellProfitReal: d.ProfitRate,
				LeverRate:      d.LeverRate})
		}
	}

	return positions, nil
}

func (dm *Hbdm) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice, leverRate int) (string, error) {
	var data struct {
		OrderId  int64 `json:"order_id"`
		COrderId int64 `json:"client_order_id"`
	}

	params := &url.Values{}
	path := "/api/v1/contract_order"

	params.Add("contract_type", contractType)
	params.Add("symbol", currencyPair.CurrencyA.Symbol)
	params.Add("volume", amount)
	params.Add("lever_rate", fmt.Sprint(leverRate))
	params.Add("contract_code", "")

	if matchPrice == 1 {
		params.Set("order_price_type" , "opponent") //对手价下单
	}else{
		params.Set("order_price_type", "limit")
		params.Add("price", price)
	}

	direction, offset := dm.adaptOpenType(openType)
	params.Add("offset", offset)
	params.Add("direction", direction)

	err := dm.doRequest(path, params, &data)

	return fmt.Sprint(data.OrderId), err
}

func (dm *Hbdm) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error) {
	var data struct {
		Successes string `json:"successes"`
		Errors    []struct {
			OrderID string `json:"order_id"`
			ErrCode int    `json:"err_code"`
			ErrMsg  string `json:"err_msg"`
		} `json:"errors"`
	}
	path := "/api/v1/contract_cancel"
	params := &url.Values{}

	params.Add("order_id", orderId)
	params.Add("symbol", currencyPair.CurrencyA.Symbol)

	err := dm.doRequest(path, params, &data)
	if err != nil {
		return false, err
	}

	if len(data.Errors) > 0 {
		return false, errors.New(fmt.Sprintf("%d:[%s]", data.Errors[0].ErrCode, data.Errors[0].ErrMsg))
	} else {
		return true, nil
	}
}

func (dm *Hbdm) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	var data struct {
		Orders      []OrderInfo `json:"orders"`
		TotalPage   int         `json:"total_page"`
		CurrentPage int         `json:"current_page"`
		TotalSize   int         `json:"total_size"`
	}

	path := "/api/v1/contract_openorders"
	params := &url.Values{}
	params.Add("symbol", currencyPair.CurrencyA.Symbol)

	err := dm.doRequest(path, params, &data)
	if err != nil {
		return nil, err
	}
	//log.Println(data)

	var ords []FutureOrder
	for _, ord := range data.Orders {

		ords = append(ords, FutureOrder{
			ContractName: contractType,
			Currency:     currencyPair,
			OType:        dm.adaptOffsetDirectionToOpenType(ord.Offset, ord.Direction),
			OrderID2:     fmt.Sprint(ord.OrderId),
			OrderID:      ord.OrderId,
			Amount:       ord.Volume,
			Price:        ord.Price,
			AvgPrice:     ord.TradeAvgPrice,
			DealAmount:   ord.TradeVolume,
			Status:       dm.adaptOrderStatus(ord.Status),
			Fee:          ord.Fee,
			LeverRate:    ord.LeverRate,
		})
	}

	return ords, err
}

func (dm *Hbdm) GetFutureOrder(orderId string, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	ords, err := dm.GetFutureOrders([]string{orderId}, currencyPair, contractType)
	if err != nil {
		return nil, err
	}

	if len(ords) == 1 {
		return &ords[0], nil
	}
	return nil, errors.New("not found order")
}

func (dm *Hbdm) GetFutureOrders(orderIds []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	var data []OrderInfo
	path := "/api/v1/contract_order_info"
	params := &url.Values{}

	params.Add("order_id", strings.Join(orderIds, ","))
	params.Add("symbol", currencyPair.CurrencyA.Symbol)

	err := dm.doRequest(path, params, &data)
	if err != nil {
		return nil, err
	}
	//	log.Println(data)

	var ords []FutureOrder
	for _, ord := range data {
		ords = append(ords, FutureOrder{
			ContractName: contractType,
			Currency:     currencyPair,
			OType:        dm.adaptOffsetDirectionToOpenType(ord.Offset, ord.Direction),
			OrderID2:     fmt.Sprint(ord.OrderId),
			OrderID:      ord.OrderId,
			Amount:       ord.Volume,
			Price:        ord.Price,
			AvgPrice:     ord.TradeAvgPrice,
			DealAmount:   ord.TradeVolume,
			Status:       dm.adaptOrderStatus(ord.Status),
			Fee:          ord.Fee,
			LeverRate:    ord.LeverRate,
		})
	}

	return ords, nil

}

func (dm *Hbdm) GetContractValue(currencyPair CurrencyPair) (float64, error) {
	switch currencyPair.CurrencyA {
	case BTC:
		return 100, nil
	default:
		return 10, nil
	}
}

func (dm *Hbdm) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	ret, err := HttpGet(dm.config.HttpClient, dm.config.Endpoint+"/api/v1//contract_delivery_price?symbol="+currencyPair.CurrencyA.Symbol)
	if err != nil {
		return -1, err
	}

	if ret["status"].(string) != "ok" {
		return -1, errors.New(fmt.Sprintf("%+v", ret))
	}

	return ToFloat64(ret["data"].(map[string]interface{})["delivery_price"]), nil
}

func (dm *Hbdm) GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error) {
	symbol := dm.adaptSymbol(currencyPair, contractType)
	ret, err := HttpGet(dm.config.HttpClient, dm.config.Endpoint+"/market/detail/merged?symbol="+symbol)
	if err != nil {
		return nil, err
	}
	//log.Println(ret)
	s := ret["status"].(string)
	if s == "error" {
		return nil, errors.New(ret["err_msg"].(string))
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
	url := dm.config.Endpoint + "/market/depth?type=step0&symbol=" + symbol
	ret, err := HttpGet(dm.config.HttpClient, url)
	if err != nil {
		return nil, err
	}

	s := ret["status"].(string)
	if s == "error" {
		return nil, errors.New(ret["err_msg"].(string))
	}
	//log.Println(ret)
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

func (dm *Hbdm) GetFutureIndex(currencyPair CurrencyPair) (float64, error) {
	ret, err := HttpGet(dm.config.HttpClient, dm.config.Endpoint+"/api/v1/contract_index?symbol="+currencyPair.CurrencyA.Symbol)
	if err != nil {
		return -1, err
	}

	if ret["status"].(string) != "ok" {
		return -1, errors.New(fmt.Sprintf("%+v", ret))
	}

	datamap := ret["data"].([]interface{})
	index := datamap[0].(map[string]interface{})["index_price"]
	return ToFloat64(index), nil
}

func (dm *Hbdm) GetKlineRecords(contract_type string, currency CurrencyPair, period, size, since int) ([]FutureKline, error) {
	symbol := dm.adaptSymbol(currency, contract_type)
	periodS := dm.adaptKLinePeriod(period)
	url := fmt.Sprintf("%s/market/history/kline?symbol=%s&period=%s&size=%d", dm.config.Endpoint, symbol, periodS, size)

	var ret struct {
		BaseResponse
		Data []struct {
			Id     int64   `json:"id"`
			Amount float64 `json:"amount"`
			Close  float64 `json:"close"`
			High   float64 `json:"high"`
			Low    float64 `json:"low"`
			Open   float64 `json:"open"`
			Vol    float64 `json:"vol"`
		} `json:"data"`
	}

	err := HttpGet4(dm.config.HttpClient, url, nil, &ret)
	if err != nil {
		return nil, err
	}

	if ret.Status != "ok" {
		return nil, errors.New(ret.ErrMsg)
	}

	var klines []FutureKline
	for _, d := range ret.Data {
		klines = append(klines, FutureKline{
			Kline: &Kline{
				Pair:      currency,
				Vol:       d.Vol,
				Open:      d.Open,
				Close:     d.Close,
				High:      d.High,
				Low:       d.Low,
				Timestamp: d.Id},
			Vol2: d.Vol})
	}

	return klines, nil
}

func (dm *Hbdm) GetDeliveryTime() (int, int, int, int) {
	return 0, 4, 0, 0
}

func (dm *Hbdm) GetExchangeRate() (float64, error) {
	panic("not supported.")
}

func (dm *Hbdm) GetFee() (float64, error) {
	return 0.003, nil
}

func (dm *Hbdm) GetTrades(contract_type string, currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not supported.")
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

func (dm *Hbdm) adaptKLinePeriod(period int) string {
	switch period {
	case KLINE_PERIOD_1MIN:
		return "1min"
	case KLINE_PERIOD_5MIN:
		return "5min"
	case KLINE_PERIOD_15MIN:
		return "15min"
	case KLINE_PERIOD_30MIN:
		return "30min"
	case KLINE_PERIOD_60MIN:
		return "60min"
	case KLINE_PERIOD_1H:
		return "1h"
	case KLINE_PERIOD_4H:
		return "4h"
	case KLINE_PERIOD_1DAY:
		return "1day"
	case KLINE_PERIOD_1WEEK:
		return "1week"
	case KLINE_PERIOD_1MONTH:
		return "1mon"
	default:
		return "1day"
	}
}

func (dm *Hbdm) adaptOpenType(openType int) (direction string, offset string) {
	switch openType {
	case OPEN_BUY:
		return "buy", "open"
	case OPEN_SELL:
		return "sell", "open"
	case CLOSE_SELL:
		return "buy", "close"
	case CLOSE_BUY:
		return "sell", "close"
	default:
		return "", ""
	}
}

func (dm *Hbdm) adaptOffsetDirectionToOpenType(offset, direction string) int {
	switch offset {
	case "close":
		if direction == "buy" {
			return CLOSE_BUY
		} else {
			return CLOSE_SELL
		}

	default:
		if direction == "buy" {
			return OPEN_BUY
		} else {
			return OPEN_SELL
		}
	}
}

func (dm *Hbdm) adaptOrderStatus(s int) TradeStatus {
	switch s {
	case 3:
		return ORDER_UNFINISH
	case 4:
		return ORDER_PART_FINISH
	case 5:
		return ORDER_FINISH
	case 6:
		return ORDER_FINISH
	case 7:
		return ORDER_CANCEL
	default:
		return ORDER_UNFINISH
	}
}

func (dm *Hbdm) buildPostForm(reqMethod, path string, postForm *url.Values) error {
	postForm.Set("AccessKeyId", dm.config.ApiKey)
	postForm.Set("SignatureMethod", "HmacSHA256")
	postForm.Set("SignatureVersion", "2")
	postForm.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05"))
	domain := strings.Replace(dm.config.Endpoint, "https://", "", len(dm.config.Endpoint))
	payload := fmt.Sprintf("%s\n%s\n%s\n%s", reqMethod, domain, path, postForm.Encode())
	sign, _ := GetParamHmacSHA256Base64Sign(dm.config.ApiSecretKey, payload)
	postForm.Set("Signature", sign)

	return nil
}

func (dm *Hbdm) doRequest(path string, params *url.Values, data interface{}) error {
	dm.buildPostForm("POST", path, params)
	jsonD, _ := ValuesToJson(*params)
	//log.Println(string(jsonD))

	var ret BaseResponse

	resp, err := HttpPostForm3(dm.config.HttpClient, dm.config.Endpoint+path+"?"+params.Encode(), string(jsonD),
		map[string]string{"Content-Type": "application/json", "Accept-Language": "zh-cn"})

	if err != nil {
		return err
	}

	//log.Println(string(resp))
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return err
	}

	if ret.Status != "ok" {
		return errors.New(fmt.Sprintf("%d:[%s]", ret.ErrCode, ret.ErrMsg))
	}

	return json.Unmarshal(ret.Data, data)
}
