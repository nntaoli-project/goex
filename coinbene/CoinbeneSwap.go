package coinbene

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
	"net/http"
	"sort"
	"strings"
	"time"
)

type baseResp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type CoinbeneOrder struct {
	OrderId        string    `json:"orderId"`
	Direction      string    `json:"direction"`
	Leverage       int       `json:"leverage,string"`
	Symbol         string    `json:"symbol"`
	OrderType      string    `json:"orderType"`
	Quantity       float64   `json:"quantity,string"`
	OrderPrice     float64   `json:"orderPrice,string"`
	OrderValue     float64   `json:"orderValue,string"`
	Fee            float64   `json:"fee,string"`
	FilledQuantity float64   `json:"filledQuantity,string"`
	AveragePrice   float64   `json:"averagePrice"`
	OrderTime      time.Time `json:"orderTime"`
	Status         string    `json:"status"`
}

type CoinbeneSwap struct {
	config APIConfig
}

func NewCoinbeneSwap(config APIConfig) *CoinbeneSwap {
	if config.Endpoint == "" {
		config.Endpoint = "http://openapi-contract.coinbene.com"
	}
	if strings.HasSuffix(config.Endpoint, "/") {
		config.Endpoint = config.Endpoint[0 : len(config.Endpoint)-1]
	}
	if config.HttpClient == nil {
		config.HttpClient = http.DefaultClient
	}
	return &CoinbeneSwap{config: config}
}

func (swap *CoinbeneSwap) GetExchangeName() string {
	return COINBENE
}

func (swap *CoinbeneSwap) GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error) {
	var data map[string]struct {
		LastPrice    string `json:"lastPrice"`
		BestAskPrice string `json:"bestAskPrice"`
		BestBidPrice string `json:"bestBidPrice"`
		High24h      string `json:"high24h"`
		Low24h       string `json:"low24h"`
		Volume24h    string `json:"volume24h"`
	}

	resp, err := swap.doAuthRequest("GET", "/api/swap/v2/market/tickers", nil)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(resp.Data, &data)
	tick := data[currencyPair.AdaptUsdToUsdt().ToSymbol("")]
	return &Ticker{
		Pair: currencyPair,
		Last: ToFloat64(tick.LastPrice),
		Sell: ToFloat64(tick.BestAskPrice),
		Buy:  ToFloat64(tick.BestBidPrice),
		Low:  ToFloat64(tick.Low24h),
		High: ToFloat64(tick.High24h),
		Vol:  ToFloat64(tick.Volume24h)}, nil
}

func (swap *CoinbeneSwap) GetFutureDepth(currencyPair CurrencyPair, contractType string, size int) (*Depth, error) {
	//adapt size
	if size <= 5 {
		size = 5
	} else if size > 5 && size <= 10 {
		size = 10
	} else if size > 10 && size <= 50 {
		size = 50
	} else if size > 50 {
		size = 100
	} else {
		size = 5
	}

	uri := fmt.Sprintf("/api/swap/v2/market/orderBook?symbol=%s&size=%d", currencyPair.AdaptUsdToUsdt().ToSymbol(""), size)
	var data struct {
		Symbol string          `json:"symbol"`
		Asks   [][]interface{} `json:"asks"`
		Bids   [][]interface{} `json:"bids"`
	}
	resp, err := swap.doAuthRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(resp.Data, &data)
	dep := new(Depth)
	dep.Pair = currencyPair
	dep.ContractType = contractType

	for _, ask := range data.Asks {
		dep.AskList = append(dep.AskList, DepthRecord{
			Price:  ToFloat64(ask[0]),
			Amount: ToFloat64(ask[1]),
		})
	}

	for _, bid := range data.Bids {
		dep.BidList = append(dep.BidList, DepthRecord{
			Price:  ToFloat64(bid[0]),
			Amount: ToFloat64(bid[1]),
		})
	}

	sort.Sort(sort.Reverse(dep.AskList))

	return dep, err
}

func (swap *CoinbeneSwap) GetFutureIndex(currencyPair CurrencyPair) (float64, error) { panic("") }

func (swap *CoinbeneSwap) GetFutureUserinfo() (*FutureAccount, error) {
	var data struct {
		AvailableBalance float64 `json:"availableBalance,string"`
		FrozenBalance    float64 `json:"frozenBalance,string"`
		Balance          float64 `json:"balance,string"`
		MarginRate       float64 `json:"marginRate,string"`
		UnrealisedPnl    float64 `json:"unrealisedPnl,string"`
	}
	resp, err := swap.doAuthRequest("GET", "/api/swap/v2/account/info", nil)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(resp.Data, &data)

	acc := new(FutureAccount)
	acc.FutureSubAccounts = make(map[Currency]FutureSubAccount, 1)
	acc.FutureSubAccounts[BTC] = FutureSubAccount{
		Currency:      BTC,
		AccountRights: data.Balance,
		KeepDeposit:   data.FrozenBalance,
		ProfitReal:    0,
		ProfitUnreal:  data.UnrealisedPnl,
		RiskRate:      data.MarginRate,
	}
	return acc, nil
}

func (swap *CoinbeneSwap) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice, leverRate int) (string, error) {
	var param struct {
		Symbol     string `json:"symbol"`
		Leverage   string `json:"leverage"`
		OrderType  string `json:"orderType"`
		OrderPrice string `json:"orderPrice"`
		Quantity   string `json:"quantity"`
		Direction  string `json:"direction"`
	}

	var data struct {
		orderId  string `json:"orderId"`
		clientId string `json:"clientId"`
	}

	param.Symbol = currencyPair.AdaptUsdToUsdt().ToSymbol("")
	param.Leverage = fmt.Sprint(leverRate)
	param.OrderType = "limit"
	param.OrderPrice = price
	param.Quantity = amount
	param.Direction = swap.adaptOpenType(openType)

	resp, err := swap.doAuthRequest("POST", "/api/swap/v2/order/place", param)

	if err != nil {
		return "", err
	}

	json.Unmarshal(resp.Data, &data)
	return data.orderId, nil
}

func (swap *CoinbeneSwap) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error) {
	var param struct {
		OrderId string `json:"orderId"`
	}
	_, err := swap.doAuthRequest("POST", "/api/swap/v2/order/cancel", param)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (swap *CoinbeneSwap) GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	uri := fmt.Sprintf("/api/swap/v2/position/list?symbol=%s", currencyPair.ToSymbol(""))
	var data []struct {
		Quantity          float64   `json:"quantity,string"`
		AvailableQuantity float64   `json:"availableQuantity,string"`
		AveragePrice      float64   `json:"averagePrice,string"`
		CreateTime        time.Time `json:"createTime"`
		Leverage          int       `json:"leverage,string"`
		LiquidationPrice  float64   `json:"liquidationPrice,string"`
		RealisedPnl       float64   `json:"realisedPnl,string"`
		UnrealisedPnl     float64   `json:"unrealisedPnl,string"`
		Side              string    `json:"side"`
		Symbol            string    `json:"symbol"`
	}
	resp, err := swap.doAuthRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(resp.Data, &data)

	var positions []FuturePosition
	var position FuturePosition
	for _, pos := range data {
		if pos.Side == "long" {
			position.BuyAmount = pos.Quantity
			position.BuyAvailable = pos.AvailableQuantity
			position.BuyPriceAvg = pos.AveragePrice
			position.BuyProfitReal = pos.RealisedPnl
			position.CreateDate = pos.CreateTime.Unix()
		} else {
			position.SellAmount = pos.Quantity
			position.SellAvailable = pos.AvailableQuantity
			position.SellPriceAvg = pos.AveragePrice
			position.SellPriceAvg = pos.RealisedPnl
			position.CreateDate = pos.CreateTime.Unix()
		}
		position.ForceLiquPrice = pos.LiquidationPrice
		position.LeverRate = pos.Leverage
	}

	positions = append(positions, position)

	return positions, nil
}

func (swap *CoinbeneSwap) GetFutureOrders(orderIds []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	panic("")
}

func (swap *CoinbeneSwap) GetFutureOrder(orderId string, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	resp, err := swap.doAuthRequest("GET", "/api/swap/v2/order/info?orderId="+orderId, nil)
	if err != nil {
		return nil, err
	}
	var data CoinbeneOrder
	json.Unmarshal(resp.Data, &data)

	if data.OrderId == "" {
		return nil, errors.New(fmt.Sprintf("not fund order[%s]", orderId))
	}

	return &FutureOrder{
		OrderID2:     data.OrderId,
		Price:        data.OrderPrice,
		Amount:       data.Quantity,
		AvgPrice:     data.AveragePrice,
		DealAmount:   data.FilledQuantity,
		OrderTime:    data.OrderTime.Unix(),
		Status:       swap.adaptOrderStatus(data.Status),
		Currency:     currencyPair,
		OrderType:    0,
		OType:        swap.adaptOType(data.Direction),
		LeverRate:    data.Leverage,
		Fee:          data.Fee,
		ContractName: data.Symbol,
	}, nil
}

func (swap *CoinbeneSwap) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	uri := "/api/swap/v2/order/openOrders?symbol=%s&pageSize=%d&pageNum=%d"

	var unFinishOrders []FutureOrder

	pageNum := 1

	for {
		resp, err := swap.doAuthRequest("GET", fmt.Sprintf(uri, currencyPair.AdaptUsdToUsdt().ToSymbol(""), 10, pageNum), nil)
		if err != nil {
			return nil, err
		}

		var data []CoinbeneOrder
		json.Unmarshal(resp.Data, &data)

		if len(data) == 0 {
			break
		}

		for _, d := range data {
			unFinishOrders = append(unFinishOrders, FutureOrder{
				OrderID2:     d.OrderId,
				Price:        d.OrderPrice,
				Amount:       d.Quantity,
				AvgPrice:     d.AveragePrice,
				DealAmount:   d.FilledQuantity,
				OrderID:      0,
				OrderTime:    d.OrderTime.Unix(),
				Status:       ORDER_UNFINISH,
				Currency:     currencyPair,
				OrderType:    0,
				OType:        swap.adaptOType(d.Direction),
				LeverRate:    d.Leverage,
				Fee:          d.Fee,
				ContractName: d.Symbol,
			})
		}

		if len(data) < 10 {
			break
		}

		pageNum++
	}

	return unFinishOrders, nil
}
func (swap *CoinbeneSwap) GetFee() (float64, error)                                    { panic("") }
func (swap *CoinbeneSwap) GetContractValue(currencyPair CurrencyPair) (float64, error) { panic("") }
func (swap *CoinbeneSwap) GetDeliveryTime() (int, int, int, int)                       { panic("") }
func (swap *CoinbeneSwap) GetKlineRecords(contract_type string, currency CurrencyPair, period, size, since int) ([]FutureKline, error) {
	panic("")
}
func (swap *CoinbeneSwap) GetTrades(contract_type string, currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("")
}

func (swap *CoinbeneSwap) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	panic("")
}

func (swap *CoinbeneSwap) doAuthRequest(method, uri string, param interface{}) (*baseResp, error) {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	header := map[string]string{
		"Content-Type":     "application/json; charset=UTF-8",
		"ACCESS-TIMESTAMP": timestamp,
		"ACCESS-KEY":       swap.config.ApiKey}
	postBody := ""
	if param != nil {
		postBodyArray, _ := json.Marshal(param)
		postBody = string(postBodyArray)
	}
	payload := fmt.Sprintf("%s%s%s%s", timestamp, method, uri, postBody)
	//println(payload)
	sign, _ := GetParamHmacSHA256Sign(swap.config.ApiSecretKey, payload)
	header["ACCESS-SIGN"] = sign
	resp, err := NewHttpRequest(swap.config.HttpClient, method, swap.config.Endpoint+uri, postBody, header)
	if err != nil {
		return nil, err
	}
	//println(string(resp))
	var ret baseResp
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return nil, err
	}
	if ret.Code != 200 {
		return nil, errors.New(fmt.Sprintf("[%d]%s", ret.Code, ret.Message))
	}
	return &ret, nil
}

func (swap *CoinbeneSwap) adaptOpenType(oType int) string {
	switch oType {
	case OPEN_BUY:
		return "openLong"
	case OPEN_SELL:
		return "openShort"
	case CLOSE_BUY:
		return "closeLong"
	case CLOSE_SELL:
		return "closeShort"
	}
	return ""
}

func (swap *CoinbeneSwap) adaptOType(direction string) int {
	switch direction {
	case "openLong":
		return OPEN_BUY
	case "openShort":
		return OPEN_SELL
	case "closeLong":
		return CLOSE_BUY
	case "closeShort":
		return CLOSE_SELL
	}

	return 0
}

func (swap *CoinbeneSwap) adaptOrderStatus(status string) TradeStatus {
	switch status {
	case "new":
		return ORDER_UNFINISH
	case "filled", "partiallyFilled":
		return ORDER_FINISH
	case "canceled":
		return ORDER_CANCEL
	}

	return ORDER_UNFINISH
}
