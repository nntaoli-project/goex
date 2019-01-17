package bitmex

import (
	"fmt"
	"net/http"
	"time"

	"github.com/json-iterator/go"
	apiclient "github.com/nntaoli-project/GoEx/bitmex/client"
	"github.com/nntaoli-project/GoEx/bitmex/client/order"
	"github.com/nntaoli-project/GoEx/bitmex/client/order_book"
	"github.com/nntaoli-project/GoEx/bitmex/client/position"
	apiuser "github.com/nntaoli-project/GoEx/bitmex/client/user"
	"github.com/nntaoli-project/GoEx/bitmex/models"

	. "github.com/nntaoli-project/GoEx"

	"errors"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	BaseURL     = "www.bitmex.com"
	BasePath    = "/api/v1"
	TestBaseURL = "testnet.bitmex.com"
)

const (
	OrderBuy  = "Buy"
	OrderSell = "Sell"

	OrderTypeLimit     = "Limit"
	OrderTypeMarket    = "Market"
	OrderTypeStop      = "Stop"      // stop lose with market price, must set stopPx
	OrderTypeStopLimit = "StopLimit" // stop lose with limit price, must set stopPx
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

//bitmex register link  https://www.bitmex.com/register/0fcQP7

type Bitmex struct {
	accessKey,
	secretKey string
	lever float64
	trans *Transport
	api   *apiclient.APIClient
}

type Info struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Timestamp int64  `json:"timestamp"`
}

func New(client *http.Client, accesskey, secretkey, baseUrl string) *Bitmex {
	b := new(Bitmex)

	cfg := &apiclient.TransportConfig{}
	cfg.Host = baseUrl
	cfg.BasePath = BasePath
	cfg.Schemes = []string{"https"}

	b.api = apiclient.NewHTTPClientWithConfig(nil, cfg)
	b.trans = NewTransport(cfg.Host, cfg.BasePath, accesskey, secretkey, cfg.Schemes)
	b.trans.Runtime.Transport = client.Transport
	b.api.SetTransport(b.trans)
	b.setTimeOffset()
	//b.setDebug(true)
	return b
}

func (b *Bitmex) setDebug(bDebug bool) {
	b.trans.SetDebug(bDebug)
}

func (b *Bitmex) setTimeOffset() error {
	info, err := b.info()
	if err != nil {
		fmt.Println(err)
		return err
	}
	nonce := time.Now().UnixNano()
	b.trans.timeOffset = nonce/1000000 - info.Timestamp
	return nil
}

// Info get server information
func (b *Bitmex) info() (info Info, err error) {
	url := fmt.Sprintf("https://%v%v", b.trans.Host, b.trans.BasePath)
	var response *http.Response
	response, err = http.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()
	var body []byte
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &info)
	return
}

func (b *Bitmex) GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error) {
	panic("not implements")
}

// Before  Get OrderBook Data  ,   Please be sure to add apikey , apisecretkey
func (b *Bitmex) GetFutureDepth(currencyPair CurrencyPair, contractType string, size int) (*Depth, error) {
	nDepth := int32(size)
	ret, err := b.api.OrderBook.OrderBookGetL2(&order_book.OrderBookGetL2Params{Depth: &nDepth, Symbol: b.pairToSymbol(currencyPair)})
	if err != nil {
		return nil, err
	}
	depth := new(Depth)
	depth.Pair = currencyPair
	for _, v := range ret.Payload {
		if *v.Side == "Sell" {
			depth.AskList = append(depth.AskList,
				DepthRecord{Price: float64(v.Price),
					Amount: float64(v.Size)})
		} else {
			depth.BidList = append(depth.BidList,
				DepthRecord{Price: float64(v.Price),
					Amount: float64(v.Size)})
		}
	}
	depth.UTime = time.Now()
	return depth, nil
}

/**
 *获取交易所名字
 */
func (b *Bitmex) GetExchangeName() string {
	return BITMEX
}

/**
 *获取交割预估价
 */
func (b *Bitmex) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	panic("not implements")
}

/**
 * 期货指数
 * @param currencyPair   btc_usd:比特币    ltc_usd :莱特币
 */
func (b *Bitmex) GetFutureIndex(currencyPair CurrencyPair) (float64, error) {
	panic("not implements")
}

/**
 *全仓账户
 */
func (b *Bitmex) GetFutureUserinfo() (*FutureAccount, error) {
	wallet, err := b.api.User.UserGetMargin(&apiuser.UserGetMarginParams{}, nil)
	if err != nil {
		return nil, err
	}
	balances := make(map[Currency]FutureSubAccount)
	balance := FutureSubAccount{
		Currency:      XBT,
		AccountRights: float64(wallet.Payload.AvailableMargin), //账户权益
		KeepDeposit:   float64(wallet.Payload.InitMargin),
		ProfitReal:    float64(wallet.Payload.RealisedPnl),
		ProfitUnreal:  float64(wallet.Payload.UnrealisedPnl),
		RiskRate:      float64(wallet.Payload.RiskValue), //保证金率

	}
	balances[XBT] = balance
	acc := new(FutureAccount)
	acc.FutureSubAccounts = balances
	return acc, nil

}

func (b *Bitmex) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice, leverRate int) (string, error) {
	ord := new(FutureOrder)
	var err error
	switch openType {
	case 1, 4:
		if matchPrice == 1 {
			ord, err = b.MarketBuy(amount, price, currencyPair)
		} else {
			ord, err = b.LimitBuy(amount, price, currencyPair)
		}

	case 2, 3:
		if matchPrice == 1 {
			ord, err = b.MarketSell(amount, price, currencyPair)
		} else {
			ord, err = b.LimitSell(amount, price, currencyPair)
		}
	}
	return ord.OrderID2, err
}

// orderType. Valid options: Market, Limit, Stop, StopLimit, MarketIfTouched, LimitIfTouched, MarketWithLeftOverAsLimit, Pegged.
// Defaults to 'Limit' when `price` is specified. Defaults to 'Stop' when `stopPx` is specified.
// Defaults to 'StopLimit' when `price` and `stopPx` are specified.
func (b *Bitmex) placeOrder(currencyPair CurrencyPair, stopPrice float64, price float64, amount int32, side, orderType, comment string, execInst string, timeInForce string) (newOrder *models.Order, err error) {
	symbols := b.pairToSymbol(currencyPair)
	params := order.OrderNewParams{
		Side:     &side,
		Symbol:   symbols,
		Text:     &comment,
		OrderQty: &amount,
		OrdType:  &orderType,
	}
	if stopPrice != 0 {
		params.StopPx = &stopPrice
	}
	if price != 0 {
		params.Price = &price
	}
	if timeInForce != "" {
		params.TimeInForce = &timeInForce
	}
	if execInst != "" {
		params.ExecInst = &execInst
	}
	orderInfo, err := b.api.Order.OrderNew(&params, nil)
	if err != nil {
		return
	}
	newOrder = orderInfo.Payload
	return
}

func (b *Bitmex) createOrder(currencyPair CurrencyPair, price float64, amount int32, side, orderType, comment string, postOnly bool, timeInForce string) (newOrder *models.Order, err error) {
	symbols := b.pairToSymbol(currencyPair)

	params := order.OrderNewParams{
		Side:     &side,
		Symbol:   symbols,
		Text:     &comment,
		OrderQty: &amount,
		OrdType:  &orderType,
	}
	if price != 0 {
		params.Price = &price
	}
	if timeInForce != "" {
		timeInForceString := timeInForce
		params.TimeInForce = &timeInForceString
	}
	if postOnly {
		execInst := "ParticipateDoNotInitiate"
		params.ExecInst = &execInst
	}
	orderInfo, err := b.api.Order.OrderNew(&params, nil)
	if err != nil {
		return
	}
	newOrder = orderInfo.Payload
	return
}

func (b *Bitmex) LimitBuy(amount, price string, currency CurrencyPair) (*FutureOrder, error) {
	comment := "open long with bitmex api"
	side := "Buy"
	orderType := "Limit"
	num, _ := strconv.Atoi(amount)
	nAmount := num
	postOnly := false
	timeInForce := ""
	p, _ := strconv.ParseFloat(price, 64)
	newOrder, err := b.createOrder(currency, p, int32(nAmount), side, orderType, comment, postOnly, timeInForce)
	if err != nil {
		return nil, err
	}
	ret := b.transOrder(currency, newOrder)
	return ret, nil
}

func (b *Bitmex) LimitSell(amount, price string, currency CurrencyPair) (*FutureOrder, error) {
	comment := "open short with bitmex api"
	num, _ := strconv.Atoi(amount)
	nAmount := 0 - num
	postOnly := false
	timeInForce := ""
	p, _ := strconv.ParseFloat(price, 64)
	newOrder, err := b.createOrder(currency, p, int32(nAmount), OrderSell, OrderTypeLimit, comment, postOnly, timeInForce)
	if err != nil {
		return nil, err
	}
	ret := b.transOrder(currency, newOrder)
	return ret, nil
}

func (b *Bitmex) MarketBuy(amount, price string, currency CurrencyPair) (*FutureOrder, error) {
	comment := "open market long with bitmex api"
	num, _ := strconv.Atoi(amount)
	nAmount := num
	newOrder, err := b.createOrder(currency, 0, int32(nAmount), OrderBuy, OrderTypeMarket, comment, false, "")
	if err != nil {
		return nil, err
	}
	ret := b.transOrder(currency, newOrder)
	return ret, nil
}

func (b *Bitmex) MarketSell(amount, price string, currency CurrencyPair) (*FutureOrder, error) {
	comment := "open market short with bitmex api"
	num, _ := strconv.Atoi(amount)
	nAmount := 0 - num
	newOrder, err := b.createOrder(currency, 0, int32(nAmount), OrderSell, OrderTypeMarket, comment, false, "")
	if err != nil {
		return nil, err
	}
	ret := b.transOrder(currency, newOrder)
	return ret, nil

}

func (b *Bitmex) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error) {
	comment := "cancle order with bitmex api"
	params := order.OrderCancelParams{
		OrderID: &orderId,
		Text:    &comment,
	}
	orderInfo, err := b.api.Order.OrderCancel(&params, nil)
	if err != nil {
		return false, err
	}
	if len(orderInfo.Payload) == 0 {
		return false, errors.New("no such order")
	}
	newOrder := b.transOrder(currencyPair, orderInfo.Payload[0])
	if newOrder.Status == ORDER_CANCEL || newOrder.Status == ORDER_CANCEL_ING {
		return true, nil
	} else {
		return false, errors.New(fmt.Sprint("order status is", newOrder.Status.String()))
	}

}

func (b *Bitmex) GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	symbols := b.pairToSymbol(currencyPair)
	filters := `{"isOpen": true, "symbol": "` + symbols + `"}`
	params := position.PositionGetParams{
		Filter: &filters,
	}
	pos, err := b.api.Position.PositionGet(&params, nil)
	if err != nil {
		return nil, err
	}
	var position *FuturePosition
	var positions []FuturePosition
	for _, v := range pos.Payload {
		position = b.transPosition(currencyPair, v)
		if position == nil {
			continue
		}
		// UnrealisedRoePcnt 是按标记价格计算的盈亏
		// UnrealisedPnl 未实现盈亏
		// UnrealisedPnlPcnt 未实现盈亏%
		// markPrice 标记价
		// avgEntryPrice 开仓均价
		positions = append(positions, *position)
	}
	return positions, nil
}

func (b *Bitmex) transPosition(currency CurrencyPair, v *models.Position) (pos *FuturePosition) {
	var orderType string
	if v.CurrentQty > 0 {
		orderType = "Long"
	} else {
		orderType = "Short"
	}
	if v.CurrentQty == 0 {
		return
	}
	pos = &FuturePosition{
		BuyAmount:    float64(v.OpenOrderBuyQty),
		BuyAvailable: float64(v.ExecBuyQty),
		BuyPriceAvg:  v.AvgCostPrice,
		BuyPriceCost: float64(v.OpenOrderBuyCost),
		//BuyProfitReal
		CreateDate: time.Time(v.OpeningTimestamp).Unix(),
		//LeverRate      int
		SellAmount:    float64(v.OpenOrderSellQty),
		SellAvailable: float64(v.ExecSellQty),
		SellPriceAvg:  v.AvgCostPrice,
		SellPriceCost: float64(v.OpenOrderBuyCost),
		//SellProfitReal float64
		Symbol:       currency, //btc_usd:比特币,ltc_usd:莱特币
		ContractType: orderType,
		//ContractId     int64
		ForceLiquPrice: v.LiquidationPrice, //预估爆仓价

	}
	return
}

/**
 *获取订单信息
 */
func (b *Bitmex) GetFutureOrdersHistory(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	symbols := b.pairToSymbol(currencyPair)
	params := order.OrderGetOrdersParams{
		Symbol: &symbols,
	}
	orderInfo, err := b.api.Order.OrderGetOrders(&params, nil)
	if err != nil {
		return nil, err
	}
	if len(orderInfo.Payload) == 0 {
		return nil, errors.New("no such order")
	}
	var orders []FutureOrder
	for _, v := range orderInfo.Payload {
		orders = append(orders, *b.transOrder(currencyPair, v))
	}

	return orders, nil
}

/**
 *获取订单信息
 */
func (b *Bitmex) GetFutureOrders(orderIds []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	symbols := b.pairToSymbol(currencyPair)
	params := order.OrderGetOrdersParams{
		Symbol: &symbols,
	}
	orderInfo, err := b.api.Order.OrderGetOrders(&params, nil)
	if err != nil {
		return nil, err
	}
	if len(orderInfo.Payload) == 0 {
		return nil, errors.New("no such order")
	}
	var orders []FutureOrder
	for _, v := range orderInfo.Payload {
		for _, vv := range orderIds {
			if vv == *v.OrderID {
				orders = append(orders, *b.transOrder(currencyPair, v))
				break
			}
		}
		if len(orderIds) == len(orders) {
			break
		}
	}

	return orders, nil
}

/**
 *获取单个订单信息
 */
func (b *Bitmex) GetFutureOrder(orderId string, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	filters := fmt.Sprintf(`{"orderID":"%s"}`, orderId)
	symbols := b.pairToSymbol(currencyPair)
	params := order.OrderGetOrdersParams{
		Symbol: &symbols,
		Filter: &filters,
	}
	orderInfo, err := b.api.Order.OrderGetOrders(&params, nil)
	if err != nil {
		return nil, err
	}
	if len(orderInfo.Payload) == 0 {
		return nil, errors.New("no such order")
	}
	newOrder := b.transOrder(currencyPair, orderInfo.Payload[0])
	return newOrder, nil
}

func (b *Bitmex) tradeStatusAdapt(status string) TradeStatus {
	switch status {
	case "New":
		return ORDER_UNFINISH
	case "Filled":
		return ORDER_FINISH
	case "Canceled":
		return ORDER_CANCEL
	default:
		return ORDER_UNFINISH
	}

}

func (b *Bitmex) tradeTypeAdapt(text string) int {
	if strings.Contains(text, "open long") {
		return 1
	} else if strings.Contains(text, "open short") {
		return 2
	} else if strings.Contains(text, "close long") {
		return 3
	} else if strings.Contains(text, "close short") {
		return 4
	} else {
		return -1
	}
}
func (b *Bitmex) transOrder(currency CurrencyPair, o *models.Order) (ret *FutureOrder) {
	ret = &FutureOrder{OrderID2: *o.OrderID,
		Currency:   currency,
		Amount:     float64(o.OrderQty),
		DealAmount: float64(o.OrderQty - o.LeavesQty),
		Price:      o.Price,
		AvgPrice:   o.AvgPx,
		Status:     b.tradeStatusAdapt(o.OrdStatus),
		OType:      b.tradeTypeAdapt(o.Text),
		OrderTime:  time.Time(o.Timestamp).Unix()}

	return
}

/**
 *获取未完成订单信息
 */
func (b *Bitmex) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	filters := `{"ordStatus":"New"}`
	symbols := b.pairToSymbol(currencyPair)
	params := order.OrderGetOrdersParams{
		Symbol: &symbols,
		Filter: &filters,
	}
	orderInfo, err := b.api.Order.OrderGetOrders(&params, nil)
	if err != nil {
		return nil, err
	}
	var orders []FutureOrder
	for _, v := range orderInfo.Payload {
		orders = append(orders, *b.transOrder(currencyPair, v))
	}
	return orders, nil

}

/**
 *获取交易费
 */
func (b *Bitmex) GetFee() (float64, error) {
	panic("not implements")
}

/**
 *获取交易所的美元人民币汇率
 */
func (b *Bitmex) GetExchangeRate() (float64, error) {
	panic("not implements")
}

/**
 *获取每张合约价值
 */
func (b *Bitmex) GetContractValue(currencyPair CurrencyPair) (float64, error) {
	panic("not implements")
}

/**
 *获取交割时间 星期(0,1,2,3,4,5,6)，小时，分，秒
 */
func (b *Bitmex) GetDeliveryTime() (int, int, int, int) {
	panic("not implements")
}

/**
 * 获取K线数据
 */
func (b *Bitmex) GetKlineRecords(contract_type string, currency CurrencyPair, period, size, since int) ([]FutureKline, error) {
	panic("not implements")
}

/**
 * 获取Trade数据
 *非个人，整个交易所的交易记录
 */
func (b *Bitmex) GetTrades(contract_type string, currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implements")
}

func (b *Bitmex) pairToSymbol(pair CurrencyPair) string {
	if pair.CurrencyA.Symbol == BTC.Symbol {
		pair.CurrencyA = XBT
	}
	if pair.CurrencyB.Symbol == BTC.Symbol {
		pair.CurrencyB = XBT
	}
	return pair.AdaptUsdtToUsd().ToSymbol("")
}
