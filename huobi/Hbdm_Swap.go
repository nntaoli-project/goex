package huobi

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/internal/logger"
	"net/url"
	"sort"
	"time"
)

type HbdmSwap struct {
	base *Hbdm
	c    *APIConfig
}

const (
	getSwapContractInfoApiPath = "/swap-ex/v1/swap_contract_info"
	tickerApiPath              = "/swap-ex/market/detail/merged"
	marketApiPath              = "/swap-ex/market/depth"
	klineApiPath               = "/swap-ex/market/history/kline"
	accountApiPath             = "/swap-api/v1/swap_account_info"
	placeOrderApiPath          = "/swap-api/v1/swap_order"
	getPositionApiPath         = "/swap-api/v1/swap_position_info"
	cancelOrderApiPath         = "/swap-api/v1/swap_cancel"
	getOpenOrdersApiPath       = "/swap-api/v1/swap_openorders"
	getOrderInfoApiPath        = "/swap-api/v1/swap_order_info"
	getHistoryOrderPath        = "/swap-api/v1/swap_hisorders_exact"
)

func NewHbdmSwap(c *APIConfig) *HbdmSwap {
	if c.Lever <= 0 {
		c.Lever = 10
	}

	return &HbdmSwap{
		base: NewHbdm(c),
		c:    c,
	}
}

func (swap *HbdmSwap) GetExchangeName() string {
	return HBDM_SWAP
}

func (swap *HbdmSwap) GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error) {
	tickerUrl := fmt.Sprintf("%s%s?contract_code=%s", swap.base.config.Endpoint, tickerApiPath, currencyPair.ToSymbol("-"))
	responseBody, err := HttpGet5(swap.base.config.HttpClient, tickerUrl, map[string]string{})
	if err != nil {
		return nil, err
	}
	logger.Debugf("response body: %s", string(responseBody))

	var (
		tickResponse struct {
			BaseResponse
			Tick struct {
				Id     int64     `json:"id"`
				Vol    float64   `json:"vol,string"`
				Count  int64     `json:"count"`
				Open   float64   `json:"open,string"`
				Close  float64   `json:"close,string"`
				Low    float64   `json:"low,string"`
				High   float64   `json:"high,string"`
				Amount float64   `json:"amount,string"`
				Ask    []float64 `json:"ask"`
				Bid    []float64 `json:"bid"`
				Ts     int64     `json:"ts"`
			} `json:"tick"`
		}
	)

	err = json.Unmarshal(responseBody, &tickResponse)
	if err != nil {
		return nil, err
	}

	if tickResponse.Status != "ok" {
		return nil, errors.New(string(responseBody))
	}

	return &Ticker{
		Pair: currencyPair,
		Last: 0,
		Buy:  tickResponse.Tick.Bid[0],
		Sell: tickResponse.Tick.Ask[0],
		High: tickResponse.Tick.High,
		Low:  tickResponse.Tick.Low,
		Vol:  tickResponse.Tick.Vol,
		Date: uint64(tickResponse.Tick.Ts),
	}, nil
}

func (swap *HbdmSwap) GetFutureDepth(currencyPair CurrencyPair, contractType string, size int) (*Depth, error) {
	step := 0
	if size <= 20 {
		step = 6
	}
	depthUrl := fmt.Sprintf("%s%s?contract_code=%s&type=step%d", swap.base.config.Endpoint, marketApiPath, currencyPair.ToSymbol("-"), step)
	responseBody, err := HttpGet5(swap.base.config.HttpClient, depthUrl, map[string]string{})
	if err != nil {
		return nil, err
	}
	logger.Debugf("response body: %s", string(responseBody))

	var (
		dep          Depth
		tickResponse struct {
			BaseResponse
			Tick struct {
				Id   int64           `json:"id"`
				Ts   int64           `json:"ts"`
				Bids [][]interface{} `json:"bids"`
				Asks [][]interface{} `json:"asks"`
			} `json:"tick"`
		}
	)

	err = json.Unmarshal(responseBody, &tickResponse)
	if err != nil {
		return nil, err
	}

	if tickResponse.Status != "ok" {
		return nil, errors.New(string(responseBody))
	}

	dep.Pair = currencyPair
	dep.ContractType = contractType
	dep.UTime = time.Unix(0, tickResponse.Ts*int64(time.Millisecond))

	for i, item := range tickResponse.Tick.Bids {
		if i >= size {
			break
		}
		dep.BidList = append(dep.BidList, DepthRecord{
			Price:  ToFloat64(item[0]),
			Amount: ToFloat64(item[1]),
		})
	}

	for i, item := range tickResponse.Tick.Asks {
		if i >= size {
			break
		}
		dep.AskList = append(dep.AskList, DepthRecord{
			Price:  ToFloat64(item[0]),
			Amount: ToFloat64(item[1]),
		})
	}

	sort.Sort(sort.Reverse(dep.AskList))

	return &dep, nil
}

func (swap *HbdmSwap) GetFutureUserinfo(currencyPair ...CurrencyPair) (*FutureAccount, error) {
	var accountInfoResponse []struct {
		Symbol           string  `json:"symbol"`
		MarginBalance    float64 `json:"margin_balance"`
		MarginPosition   float64 `json:"margin_position"`
		MarginFrozen     float64 `json:"margin_frozen"`
		MarginAvailable  float64 `json:"margin_available"`
		ProfitReal       float64 `json:"profit_real"`
		ProfitUnreal     float64 `json:"profit_unreal"`
		RiskRate         float64 `json:"risk_rate"`
		LiquidationPrice float64 `json:"liquidation_price"`
	}

	param := url.Values{}
	if len(currencyPair) > 0 {
		param.Set("contract_code", currencyPair[0].ToSymbol("-"))
	}

	err := swap.base.doRequest(accountApiPath, &param, &accountInfoResponse)
	if err != nil {
		return nil, err
	}

	var futureAccount FutureAccount
	futureAccount.FutureSubAccounts = make(map[Currency]FutureSubAccount, 4)

	for _, acc := range accountInfoResponse {
		currency := NewCurrency(acc.Symbol, "")
		futureAccount.FutureSubAccounts[currency] = FutureSubAccount{
			Currency:      currency,
			AccountRights: acc.MarginBalance,
			KeepDeposit:   acc.MarginPosition,
			ProfitReal:    acc.ProfitReal,
			ProfitUnreal:  acc.ProfitUnreal,
			RiskRate:      acc.RiskRate,
		}
	}

	return &futureAccount, nil
}

func (swap *HbdmSwap) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice int, leverRate float64) (string, error) {
	param := url.Values{}
	param.Set("contract_code", currencyPair.ToSymbol("-"))
	param.Set("client_order_id", fmt.Sprint(time.Now().UnixNano()))
	param.Set("price", price)
	param.Set("volume", amount)
	param.Set("lever_rate", fmt.Sprintf("%.0f", leverRate))

	direction, offset := swap.base.adaptOpenType(openType)
	param.Set("direction", direction)
	param.Set("offset", offset)
	logger.Info(direction, offset)

	if matchPrice == 1 {
		param.Set("order_price_type", "opponent")
	} else {
		param.Set("order_price_type", "limit")
	}

	var orderResponse struct {
		OrderId       string `json:"order_id_str"`
		ClientOrderId int64  `json:"client_order_id"`
	}

	err := swap.base.doRequest(placeOrderApiPath, &param, &orderResponse)
	if err != nil {
		return "", err
	}

	return orderResponse.OrderId, nil
}

func (swap *HbdmSwap) LimitFuturesOrder(currencyPair CurrencyPair, contractType, price, amount string, openType int, opt ...LimitOrderOptionalParameter) (*FutureOrder, error) {
	orderId, err := swap.PlaceFutureOrder(currencyPair, contractType, price, amount, openType, 0, swap.c.Lever)
	return &FutureOrder{
		Currency:     currencyPair,
		OrderID2:     orderId,
		Amount:       ToFloat64(amount),
		Price:        ToFloat64(price),
		OType:        openType,
		ContractName: contractType,
	}, err
}

func (swap *HbdmSwap) MarketFuturesOrder(currencyPair CurrencyPair, contractType, amount string, openType int) (*FutureOrder, error) {
	orderId, err := swap.PlaceFutureOrder(currencyPair, contractType, "", amount, openType, 1, 10)
	return &FutureOrder{
		Currency:     currencyPair,
		OrderID2:     orderId,
		Amount:       ToFloat64(amount),
		OType:        openType,
		ContractName: contractType,
	}, err
}

func (swap *HbdmSwap) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error) {
	param := url.Values{}
	param.Set("order_id", orderId)
	param.Set("contract_code", currencyPair.ToSymbol("-"))

	var cancelResponse struct {
		Errors []struct {
			ErrMsg    string `json:"err_msg"`
			Successes string `json:"successes,omitempty"`
		} `json:"errors"`
	}

	err := swap.base.doRequest(cancelOrderApiPath, &param, &cancelResponse)
	if err != nil {
		return false, err
	}

	if len(cancelResponse.Errors) > 0 {
		return false, errors.New(cancelResponse.Errors[0].ErrMsg)
	}

	return true, nil
}

func (swap *HbdmSwap) GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	param := url.Values{}
	param.Set("contract_code", currencyPair.ToSymbol("-"))

	var (
		tempPositionMap  map[string]*FuturePosition
		futuresPositions []FuturePosition
		positionResponse []struct {
			Symbol         string
			ContractCode   string  `json:"contract_code"`
			Volume         float64 `json:"volume"`
			Available      float64 `json:"available"`
			CostOpen       float64 `json:"cost_open"`
			CostHold       float64 `json:"cost_hold"`
			ProfitUnreal   float64 `json:"profit_unreal"`
			ProfitRate     float64 `json:"profit_rate"`
			Profit         float64 `json:"profit"`
			PositionMargin float64 `json:"position_margin"`
			LeverRate      float64 `json:"lever_rate"`
			Direction      string  `json:"direction"`
		}
	)

	err := swap.base.doRequest(getPositionApiPath, &param, &positionResponse)
	if err != nil {
		return nil, err
	}

	futuresPositions = make([]FuturePosition, 0, 2)
	tempPositionMap = make(map[string]*FuturePosition, 2)

	for _, pos := range positionResponse {
		if tempPositionMap[pos.ContractCode] == nil {
			tempPositionMap[pos.ContractCode] = new(FuturePosition)
		}
		switch pos.Direction {
		case "sell":
			tempPositionMap[pos.ContractCode].ContractType = pos.ContractCode
			tempPositionMap[pos.ContractCode].Symbol = NewCurrencyPair3(pos.ContractCode, "-")
			tempPositionMap[pos.ContractCode].SellAmount = pos.Volume
			tempPositionMap[pos.ContractCode].SellAvailable = pos.Available
			tempPositionMap[pos.ContractCode].SellPriceAvg = pos.CostOpen
			tempPositionMap[pos.ContractCode].SellPriceCost = pos.CostHold
			tempPositionMap[pos.ContractCode].SellProfitReal = pos.ProfitRate
			tempPositionMap[pos.ContractCode].SellProfit = pos.Profit
		case "buy":
			tempPositionMap[pos.ContractCode].ContractType = pos.ContractCode
			tempPositionMap[pos.ContractCode].Symbol = NewCurrencyPair3(pos.ContractCode, "-")
			tempPositionMap[pos.ContractCode].BuyAmount = pos.Volume
			tempPositionMap[pos.ContractCode].BuyAvailable = pos.Available
			tempPositionMap[pos.ContractCode].BuyPriceAvg = pos.CostOpen
			tempPositionMap[pos.ContractCode].BuyPriceCost = pos.CostHold
			tempPositionMap[pos.ContractCode].BuyProfitReal = pos.ProfitRate
			tempPositionMap[pos.ContractCode].BuyProfit = pos.Profit
		}
	}

	for _, pos := range tempPositionMap {
		futuresPositions = append(futuresPositions, *pos)
	}

	return futuresPositions, nil
}

func (swap *HbdmSwap) GetFutureOrders(orderIds []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	return nil, nil
}

func (swap *HbdmSwap) GetFutureOrder(orderId string, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	var (
		orderInfoResponse []OrderInfo
		param             = url.Values{}
	)

	param.Set("contract_code", currencyPair.ToSymbol("-"))
	param.Set("order_id", orderId)

	err := swap.base.doRequest(getOrderInfoApiPath, &param, &orderInfoResponse)
	if err != nil {
		return nil, err
	}

	if len(orderInfoResponse) == 0 {
		return nil, errors.New("not found")
	}

	orderInfo := orderInfoResponse[0]

	return &FutureOrder{
		Currency:     currencyPair,
		ClientOid:    fmt.Sprint(orderInfo.ClientOrderId),
		OrderID2:     fmt.Sprint(orderInfo.OrderId),
		Price:        orderInfo.Price,
		Amount:       orderInfo.Volume,
		AvgPrice:     orderInfo.TradeAvgPrice,
		DealAmount:   orderInfo.TradeVolume,
		OrderID:      orderInfo.OrderId,
		Status:       swap.base.adaptOrderStatus(orderInfo.Status),
		OType:        swap.base.adaptOffsetDirectionToOpenType(orderInfo.Offset, orderInfo.Direction),
		LeverRate:    orderInfo.LeverRate,
		Fee:          orderInfo.Fee,
		ContractName: orderInfo.ContractCode,
		OrderTime:    orderInfo.CreatedAt,
	}, nil
}

func (swap *HbdmSwap) GetFutureOrderHistory(pair CurrencyPair, contractType string, optional ...OptionalParameter) ([]FutureOrder, error) {
	params := url.Values{}
	params.Add("status", "0")     //all
	params.Add("type", "1")       //all
	params.Add("trade_type", "0") //all

	if contractType == "" || contractType == SWAP_CONTRACT {
		params.Add("contract_code", pair.AdaptUsdtToUsd().ToSymbol("-"))
	} else {
		return nil, errors.New("contract type is error")
	}

	MergeOptionalParameter(&params, optional...)

	var historyOrderResp struct {
		Orders     []OrderInfo `json:"orders"`
		RemainSize int64       `json:"remain_size"`
		NextId     int64       `json:"next_id"`
	}

	err := swap.base.doRequest(getHistoryOrderPath, &params, &historyOrderResp)
	if err != nil {
		return nil, err
	}

	var historyOrders []FutureOrder

	for _, ord := range historyOrderResp.Orders {
		historyOrders = append(historyOrders, FutureOrder{
			OrderID:      ord.OrderId,
			OrderID2:     fmt.Sprintf("%d", ord.OrderId),
			Price:        ord.Price,
			Amount:       ord.Volume,
			AvgPrice:     ord.TradeAvgPrice,
			DealAmount:   ord.TradeVolume,
			OrderTime:    ord.CreateDate,
			Status:       swap.base.adaptOrderStatus(ord.Status),
			Currency:     pair,
			OType:        swap.base.adaptOffsetDirectionToOpenType(ord.Offset, ord.Direction),
			LeverRate:    ord.LeverRate,
			Fee:          ord.Fee,
			ContractName: ord.ContractCode,
		})
	}

	return historyOrders, nil
}

func (swap *HbdmSwap) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	param := url.Values{}
	param.Set("contract_code", currencyPair.ToSymbol("-"))
	param.Set("page_size", "50")

	var openOrderResponse struct {
		Orders []OrderInfo
	}

	err := swap.base.doRequest(getOpenOrdersApiPath, &param, &openOrderResponse)
	if err != nil {
		return nil, err
	}

	openOrders := make([]FutureOrder, 0, len(openOrderResponse.Orders))
	for _, ord := range openOrderResponse.Orders {
		openOrders = append(openOrders, FutureOrder{
			Currency:   currencyPair,
			ClientOid:  fmt.Sprint(ord.ClientOrderId),
			OrderID2:   fmt.Sprint(ord.OrderId),
			Price:      ord.Price,
			Amount:     ord.Volume,
			AvgPrice:   ord.TradeAvgPrice,
			DealAmount: ord.TradeVolume,
			OrderID:    ord.OrderId,
			Status:     swap.base.adaptOrderStatus(ord.Status),
			OType:      swap.base.adaptOffsetDirectionToOpenType(ord.Offset, ord.Direction),
			LeverRate:  ord.LeverRate,
			Fee:        ord.Fee,
			OrderTime:  ord.CreatedAt,
		})
	}

	return openOrders, nil
}

func (swap *HbdmSwap) GetContractValue(currencyPair CurrencyPair) (float64, error) {
	switch currencyPair {
	case BTC_USD, BTC_USDT:
		return 100, nil
	default:
		return 0, nil
	}
}

func (swap *HbdmSwap) GetKlineRecords(contractType string, currency CurrencyPair, period KlinePeriod, size int, opt ...OptionalParameter) ([]FutureKline, error) {
	klineType := AdaptKlinePeriodForOKEx(int(period))
	contractCode := currency.ToSymbol("-")
	apiUrl := fmt.Sprintf("%s%s?contract_code=%s&period=%s&size=%d", swap.c.Endpoint, klineApiPath, contractCode, klineType, size)
	responseBody, err := HttpGet5(swap.base.config.HttpClient, apiUrl, map[string]string{})
	if err != nil {
		return nil, err
	}
	logger.Debugf("response body: %s", string(responseBody))

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

	err = json.Unmarshal(responseBody, &ret)
	if err != nil {
		logger.Errorf("[hbdm-swap] err=%s", err.Error())
		return nil, err
	}

	var lines []FutureKline
	for i := len(ret.Data) - 1; i >= 0; i-- {
		d := ret.Data[i]
		lines = append(lines, FutureKline{
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

	return lines, err
}

func (swap *HbdmSwap) GetTrades(contractType string, currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}

func (swap *HbdmSwap) GetFee() (float64, error) {
	panic("not implement")
}

func (swap *HbdmSwap) GetFutureIndex(currencyPair CurrencyPair) (float64, error) {
	panic("not implement")
}

func (swap *HbdmSwap) GetDeliveryTime() (int, int, int, int) {
	panic("not implement")
}

func (swap *HbdmSwap) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	panic("not implement")
}
