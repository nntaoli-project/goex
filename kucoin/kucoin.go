package kucoin

import (
	"github.com/Kucoin/kucoin-go-sdk"
	. "github.com/nntaoli-project/goex"
	log "github.com/sirupsen/logrus"
	"time"
)

func New(apiKey, apiSecret, apiPassphrase string) *KuCoin {
	return NewWithConfig(&APIConfig{
		Endpoint:      "https://api.kcs.top",
		ApiKey:        apiKey,
		ApiSecretKey:  apiSecret,
		ApiPassphrase: apiPassphrase,
	})
}

func NewWithConfig(config *APIConfig) *KuCoin {
	if config.Endpoint == "" {
		config.Endpoint = "https://api.kucoin.com"
	}

	kc := &KuCoin{
		baseUrl:       config.Endpoint,
		apiKey:        config.ApiKey,
		apiSecret:     config.ApiSecretKey,
		apiPassphrase: config.ApiPassphrase,
	}

	kc.service = kucoin.NewApiService(
		kucoin.ApiBaseURIOption(kc.baseUrl),
		kucoin.ApiKeyOption(kc.apiKey),
		kucoin.ApiSecretOption(kc.apiSecret),
		kucoin.ApiPassPhraseOption(kc.apiPassphrase),
	)

	return kc
}

type KuCoin struct {
	apiKey        string
	apiSecret     string
	baseUrl       string
	apiPassphrase string
	service       *kucoin.ApiService
}

var inernalKlinePeriodConverter = map[int]string{
	KLINE_PERIOD_1MIN:  "1min",
	KLINE_PERIOD_3MIN:  "3min",
	KLINE_PERIOD_5MIN:  "5min",
	KLINE_PERIOD_15MIN: "15min",
	KLINE_PERIOD_30MIN: "30min",
	KLINE_PERIOD_60MIN: "1hour",
	KLINE_PERIOD_1H:    "1hour",
	KLINE_PERIOD_2H:    "2hour",
	KLINE_PERIOD_4H:    "4hour",
	KLINE_PERIOD_6H:    "6hour",
	KLINE_PERIOD_8H:    "8hour",
	KLINE_PERIOD_12H:   "12hour",
	KLINE_PERIOD_1DAY:  "1day",
	KLINE_PERIOD_1WEEK: "1week",
}

func (kc *KuCoin) GetExchangeName() string {
	return KUCOIN
}

func (kc *KuCoin) GetTicker(currency CurrencyPair) (*Ticker, error) {
	resp, err := kc.service.TickerLevel1(currency.ToSymbol("-"))
	if err != nil {
		log.Error("KuCoin GetTicker error:", err)
		return nil, err
	}

	var model kucoin.TickerLevel1Model

	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin GetTicker error:", err)
		return nil, err
	}

	var ticker Ticker
	ticker.Pair = currency
	ticker.Date = uint64(model.Time / 1000)
	ticker.Last = ToFloat64(model.Price)
	ticker.Buy = ToFloat64(model.BestBid)
	ticker.Sell = ToFloat64(model.BestAsk)
	return &ticker, nil
}

func (kc *KuCoin) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	clientID := UUID()
	params := map[string]string{
		"clientOid": clientID,
		"side":      "buy",
		"symbol":    currency.ToSymbol("-"),
		"type":      "limit",
		"price":     price,
		"size":      amount,
	}
	resp, err := kc.service.CreateOrder(params)
	if err != nil {
		log.Error("KuCoin LimitBuy error:", err)
		return nil, err
	}

	var model kucoin.OrderModel

	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin LimitBuy error:", err)
		return nil, err
	}

	var order Order
	order.OrderID2 = model.Id
	order.Cid = clientID
	return &order, nil
}

func (kc *KuCoin) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	clientID := UUID()
	params := map[string]string{
		"clientOid": clientID,
		"side":      "sell",
		"symbol":    currency.ToSymbol("-"),
		"type":      "limit",
		"price":     price,
		"size":      amount,
	}
	resp, err := kc.service.CreateOrder(params)
	if err != nil {
		log.Error("KuCoin LimitSell error:", err)
		return nil, err
	}

	var model kucoin.OrderModel

	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin LimitSell error:", err)
		return nil, err
	}

	var order Order
	order.OrderID2 = model.Id
	order.Cid = clientID
	return &order, nil
}

func (kc *KuCoin) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	clientID := UUID()
	params := map[string]string{
		"clientOid": clientID,
		"side":      "buy",
		"symbol":    currency.ToSymbol("-"),
		"type":      "market",
		"price":     price,
		"size":      amount,
	}
	resp, err := kc.service.CreateOrder(params)
	if err != nil {
		log.Error("KuCoin MarketBuy error:", err)
		return nil, err
	}

	var model kucoin.OrderModel

	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin MarketBuy error:", err)
		return nil, err
	}

	var order Order
	order.OrderID2 = model.Id
	order.Cid = clientID
	return &order, nil
}

func (kc *KuCoin) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	clientID := UUID()
	params := map[string]string{
		"clientOid": clientID,
		"side":      "sell",
		"symbol":    currency.ToSymbol("-"),
		"type":      "market",
		"price":     price,
		"size":      amount,
	}
	resp, err := kc.service.CreateOrder(params)
	if err != nil {
		log.Error("KuCoin MarketSell error:", err)
		return nil, err
	}

	var model kucoin.OrderModel

	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin MarketSell error:", err)
		return nil, err
	}

	var order Order
	order.OrderID2 = model.Id
	order.Cid = clientID
	return &order, nil
}

func (kc *KuCoin) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	var resp *kucoin.ApiResponse
	var err error
	if orderId != "" {
		resp, err = kc.service.CancelOrder(orderId)
	} else {
		resp, err = kc.service.CancelOrder(currency.ToSymbol("-"))
	}

	if err != nil {
		log.Error("KuCoin CancelOrder error:", err)
		return false, err
	}

	var model kucoin.CancelOrderResultModel
	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin CancelOrder error:", err)
		return false, err
	}
	return true, nil
}

func (kc *KuCoin) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	resp, err := kc.service.Order(orderId)
	if err != nil {
		log.Error("KuCoin GetOneOrder error:", err)
		return nil, err
	}

	var model kucoin.OrderModel

	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin GetOneOrder error:", err)
		return nil, err
	}

	var order Order
	order.Price = ToFloat64(model.Price)
	order.Amount = ToFloat64(model.Size)
	order.AvgPrice = ToFloat64(model.DealFunds) / ToFloat64(model.DealSize)
	order.DealAmount = ToFloat64(model.DealSize)
	order.Fee = ToFloat64(model.Fee)
	order.Cid = model.ClientOid
	order.OrderID2 = model.Id
	order.OrderTime = int(model.CreatedAt / 1000)
	order.Type = model.Type
	order.Currency = NewCurrencyPair3(model.Symbol, "-")
	if model.Side == "buy" {
		if model.Type == "limit" {
			order.Side = BUY
		} else {
			order.Side = BUY_MARKET
		}
	} else {
		if model.Type == "limit" {
			order.Side = SELL
		} else {
			order.Side = SELL_MARKET
		}
	}

	if model.CancelExist {
		if model.DealSize == "0" {
			order.Status = ORDER_UNFINISH
		} else {
			order.Status = ORDER_PART_FINISH
		}
	} else {
		if model.DealSize == "0" {
			order.Status = ORDER_CANCEL
		} else if model.DealSize == model.Side {
			order.Status = ORDER_FINISH
		}
	}

	return &order, nil
}

func (kc *KuCoin) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	params := map[string]string{
		"status": "active",
		"symbol": currency.ToSymbol("-"),
	}
	resp, err := kc.service.Orders(params, nil)
	if err != nil {
		log.Error("KuCoin GetUnfinishOrders error:", err)
		return nil, err
	}

	var model kucoin.OrderModel

	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin GetUnfinishOrders error:", err)
		return nil, err
	}
	var orders []Order
	return orders, nil
}

func (kc *KuCoin) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	params := map[string]string{
		"status": "done",
		"symbol": currency.ToSymbol("-"),
	}
	pagination := kucoin.PaginationParam{
		CurrentPage: int64(currentPage),
		PageSize:    int64(pageSize),
	}
	resp, err := kc.service.Orders(params, &pagination)
	if err != nil {
		log.Error("KuCoin GetOrderHistorys error:", err)
		return nil, err
	}

	var model kucoin.OrderModel

	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin GetOrderHistorys error:", err)
		return nil, err
	}
	var orders []Order
	return orders, nil
}

func (kc *KuCoin) GetAccount() (*Account, error) {
	var account Account
	return &account, nil
}

func (kc *KuCoin) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	dep := 20
	if size > 20 {
		dep = 100
	}
	resp, err := kc.service.AggregatedPartOrderBook(currency.ToSymbol("-"), int64(dep))
	if err != nil {
		log.Error("KuCoin GetDepth error:", err)
		return nil, err
	}

	var model kucoin.PartOrderBookModel

	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin GetDepth error:", err)
		return nil, err
	}
	var depth Depth
	depth.Pair = currency
	depth.UTime = time.Now()
	for i, ask := range model.Asks {
		if i < size {
			depth.AskList = append(depth.AskList, DepthRecord{
				Price:  ToFloat64(ask[0]),
				Amount: ToFloat64(ask[1]),
			})
		}
	}

	for j, bid := range model.Bids {
		if j < size {
			depth.BidList = append(depth.BidList, DepthRecord{
				Price:  ToFloat64(bid[0]),
				Amount: ToFloat64(bid[1]),
			})
		}
	}

	return &depth, nil
}

func (kc *KuCoin) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	resp, err := kc.service.KLines(currency.ToSymbol("-"), inernalKlinePeriodConverter[period], int64(since), time.Now().UnixNano()/int64(time.Millisecond))
	if err != nil {
		log.Error("KuCoin GetKlineRecords error:", err)
		return nil, err
	}

	var model kucoin.KLinesModel

	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin GetKlineRecords error:", err)
		return nil, err
	}
	var kLines []Kline
	for i, item := range model {
		if i < size {
			kLines = append(kLines, Kline{
				Pair:      currency,
				Timestamp: ToInt64((*item)[0]),
				Open:      ToFloat64((*item)[1]),
				Close:     ToFloat64((*item)[2]),
				High:      ToFloat64((*item)[3]),
				Low:       ToFloat64((*item)[4]),
				Vol:       ToFloat64((*item)[6]),
			})
		}
	}

	return kLines, nil
}

func (kc *KuCoin) GetTrades(currency CurrencyPair, since int64) ([]Trade, error) {
	resp, err := kc.service.TradeHistories(currency.ToSymbol("-"))
	if err != nil {
		log.Error("KuCoin GetTrades error:", err)
		return nil, err
	}

	var model kucoin.TradeHistoriesModel

	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin GetTrades error:", err)
		return nil, err
	}

	var trades []Trade

	for _, item := range model {
		typo := BUY
		if (*item).Side == "sell" {
			typo = SELL
		}
		trades = append(trades, Trade{
			Pair:   currency,
			Tid:    ToInt64((*item).Sequence),
			Type:   typo,
			Amount: ToFloat64((*item).Size),
			Price:  ToFloat64((*item).Price),
			Date:   ToInt64((*item).Time / 1000),
		})
	}

	return trades, nil
}

// Account

// Accounts returns a list of accounts.
func (kc *KuCoin) Accounts(currency, typo string) (kucoin.AccountsModel, error) {
	resp, err := kc.service.Accounts(currency, typo)
	if err != nil {
		log.Error("KuCoin Accounts error:", err)
		return nil, err
	}

	var model kucoin.AccountsModel
	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin Accounts error:", err)
		return nil, err
	}

	return model, nil
}

// Account returns an account when you know the accountId.
func (kc *KuCoin) Account(accountId string) (*kucoin.AccountModel, error) {
	resp, err := kc.service.Account(accountId)
	if err != nil {
		log.Error("KuCoin Accounts error:", err)
		return nil, err
	}

	var model *kucoin.AccountModel
	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin Accounts error:", err)
		return nil, err
	}

	return model, nil
}

// SubAccountUsers returns a list of sub-account user.
func (kc *KuCoin) SubAccountUsers() (kucoin.SubAccountUsersModel, error) {
	resp, err := kc.service.SubAccountUsers()
	if err != nil {
		log.Error("KuCoin SubAccountUsers error:", err)
		return nil, err
	}

	var model kucoin.SubAccountUsersModel
	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin SubAccountUsers error:", err)
		return nil, err
	}

	return model, nil
}

// SubAccounts returns the aggregated balance of all sub-accounts of the current user.
func (kc *KuCoin) SubAccounts() (kucoin.SubAccountsModel, error) {
	resp, err := kc.service.SubAccounts()
	if err != nil {
		log.Error("KuCoin SubAccounts error:", err)
		return nil, err
	}

	var model kucoin.SubAccountsModel
	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin SubAccounts error:", err)
		return nil, err
	}

	return model, nil
}

// SubAccount returns the detail of a sub-account.
func (kc *KuCoin) SubAccount(subUserId string) (*kucoin.SubAccountModel, error) {
	resp, err := kc.service.SubAccount(subUserId)
	if err != nil {
		log.Error("KuCoin SubAccount error:", err)
		return nil, err
	}

	var model *kucoin.SubAccountModel
	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin SubAccount error:", err)
		return nil, err
	}

	return model, nil
}

// CreateAccount creates an account according to type(main|trade) and currency
// Parameter #1 typo is type of account.
func (kc *KuCoin) CreateAccount(typo, currency string) (*kucoin.AccountModel, error) {
	resp, err := kc.service.CreateAccount(typo, currency)
	if err != nil {
		log.Error("KuCoin CreateAccount error:", err)
		return nil, err
	}

	var model *kucoin.AccountModel
	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin CreateAccount error:", err)
		return nil, err
	}

	return model, nil
}

// The inner transfer interface is used for transferring assets between the accounts of a user and is free of charges.
// For example, a user could transfer assets from their main account to their trading account on the platform.
func (kc *KuCoin) InnerTransfer(currency, from, to, amount string) (string, error) {
	resp, err := kc.service.InnerTransferV2(UUID(), currency, from, to, amount)
	if err != nil {
		log.Error("KuCoin InnerTransfer error:", err)
		return "", err
	}

	var model *kucoin.InnerTransferResultModel
	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin InnerTransfer error:", err)
		return "", err
	}

	return model.OrderId, nil
}

// SubTransfer transfers between master account and sub-account.
func (kc *KuCoin) SubTransfer(currency, amount, direction, subUserId, accountType, subAccountType string) (string, error) {
	params := map[string]string{
		"clientOid":      UUID(),
		"currency":       currency,
		"amount":         amount,
		"direction":      direction,      // IN or OUT
		"subUserId":      subUserId,      // the user ID of a sub-account
		"accountType":    accountType,    // The account type of the master user: MAIN
		"subAccountType": subAccountType, //The account type of the sub user: MAIN, TRADE or MARGIN
	}
	resp, err := kc.service.SubTransfer(params)
	if err != nil {
		log.Error("KuCoin SubTransfer error:", err)
		return "", err
	}

	var model *kucoin.InnerTransferResultModel
	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin SubTransfer error:", err)
		return "", err
	}

	return model.OrderId, nil
}

// Deposits

// CreateDepositAddress creates a deposit address.
func (kc *KuCoin) CreateDepositAddress(currency, chain string) (*kucoin.DepositAddressModel, error) {
	resp, err := kc.service.CreateDepositAddress(currency, chain)
	if err != nil {
		log.Error("KuCoin CreateDepositAddress error:", err)
		return nil, err
	}

	var model *kucoin.DepositAddressModel
	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin CreateDepositAddress error:", err)
		return nil, err
	}

	return model, nil
}

// DepositAddresses returns the deposit address of currency for deposit.
// If return data is empty, you may need create a deposit address first.
func (kc *KuCoin) DepositAddresses(currency, chain string) (*kucoin.DepositAddressModel, error) {
	resp, err := kc.service.DepositAddresses(currency, chain)
	if err != nil {
		log.Error("KuCoin DepositAddresses error:", err)
		return nil, err
	}

	var model *kucoin.DepositAddressModel
	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin DepositAddresses error:", err)
		return nil, err
	}

	return model, nil
}

// Deposits returns a list of deposit.
func (kc *KuCoin) Deposits(currency, startAt, endAt, status string) (*kucoin.DepositsModel, error) {
	params := map[string]string{
		"currency": currency,
		"startAt":  startAt,
		"endAt":    endAt,
		"status":   status,
	}
	resp, err := kc.service.Deposits(params, nil)
	if err != nil {
		log.Error("KuCoin Deposits error:", err)
		return nil, err
	}

	var model *kucoin.DepositsModel
	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin Deposits error:", err)
		return nil, err
	}

	return model, nil
}

// Withdrawals

// Deposits returns a list of deposit.
func (kc *KuCoin) Withdrawals(currency, startAt, endAt, status string) (*kucoin.WithdrawalsModel, error) {
	params := map[string]string{
		"currency": currency,
		"startAt":  startAt,
		"endAt":    endAt,
		"status":   status,
	}
	resp, err := kc.service.Withdrawals(params, nil)
	if err != nil {
		log.Error("KuCoin Withdrawals error:", err)
		return nil, err
	}

	var model *kucoin.WithdrawalsModel
	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin Withdrawals error:", err)
		return nil, err
	}

	return model, nil
}

// ApplyWithdrawal applies a withdrawal.
func (kc *KuCoin) ApplyWithdrawal(currency, address, amount, memo, isInner, remark, chain string) (string, error) {
	resp, err := kc.service.ApplyWithdrawal(currency, address, amount, map[string]string{
		"memo":    memo,
		"remark":  remark,
		"chain":   chain,
		"isInner": isInner,
	})
	if err != nil {
		log.Error("KuCoin ApplyWithdrawal error:", err)
		return "", err
	}

	var model *kucoin.ApplyWithdrawalResultModel
	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin ApplyWithdrawal error:", err)
		return "", err
	}

	return model.WithdrawalId, nil
}

// WithdrawalQuotas returns the quotas of withdrawal.
func (kc *KuCoin) WithdrawalQuotas(currency, chain string) (*kucoin.WithdrawalQuotasModel, error) {
	resp, err := kc.service.WithdrawalQuotas(currency, chain)
	if err != nil {
		log.Error("KuCoin WithdrawalQuotas error:", err)
		return nil, err
	}

	var model *kucoin.WithdrawalQuotasModel
	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin WithdrawalQuotas error:", err)
		return nil, err
	}

	return model, nil
}

// CancelWithdrawal cancels a withdrawal by withdrawalId.
func (kc *KuCoin) CancelWithdrawal(withdrawalId string) (*kucoin.CancelWithdrawalResultModel, error) {
	resp, err := kc.service.CancelWithdrawal(withdrawalId)
	if err != nil {
		log.Error("KuCoin CancelWithdrawal error:", err)
		return nil, err
	}

	var model *kucoin.CancelWithdrawalResultModel
	err = resp.ReadData(&model)
	if err != nil {
		log.Error("KuCoin CancelWithdrawal error:", err)
		return nil, err
	}

	return model, nil
}
