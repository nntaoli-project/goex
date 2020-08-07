package goex

import (
	"net/http"
	"time"
)

type Order struct {
	Price        float64
	Amount       float64
	AvgPrice     float64
	DealAmount   float64
	Fee          float64
	Cid          string //客户端自定义ID
	OrderID2     string
	OrderID      int //deprecated
	Status       TradeStatus
	Currency     CurrencyPair
	Side         TradeSide
	Type         string //limit / market
	OrderType    int    //0:default,1:maker,2:fok,3:ioc
	OrderTime    int    // create  timestamp
	FinishedTime int64  //finished timestamp
}

type Trade struct {
	Tid    int64        `json:"tid"`
	Type   TradeSide    `json:"type"`
	Amount float64      `json:"amount,string"`
	Price  float64      `json:"price,string"`
	Date   int64        `json:"date_ms"`
	Pair   CurrencyPair `json:"omitempty"`
}

type SubAccount struct {
	Currency     Currency
	Amount       float64
	ForzenAmount float64
	LoanAmount   float64
}

type MarginSubAccount struct {
	Balance     float64
	Frozen      float64
	Available   float64
	CanWithdraw float64
	Loan        float64
	LendingFee  float64
}

type MarginAccount struct {
	Sub              map[Currency]MarginSubAccount
	LiquidationPrice float64
	RiskRate         float64
	MarginRatio      float64
}

type Account struct {
	Exchange    string
	Asset       float64 //总资产
	NetAsset    float64 //净资产
	SubAccounts map[Currency]SubAccount
}

type Ticker struct {
	Pair CurrencyPair `json:"omitempty"`
	Last float64      `json:"last,string"`
	Buy  float64      `json:"buy,string"`
	Sell float64      `json:"sell,string"`
	High float64      `json:"high,string"`
	Low  float64      `json:"low,string"`
	Vol  float64      `json:"vol,string"`
	Date uint64       `json:"date"` // 单位:ms
}

type FutureTicker struct {
	*Ticker
	ContractType string  `json:"omitempty"`
	ContractId   string  `json:"contractId"`
	LimitHigh    float64 `json:"limitHigh,string"`
	LimitLow     float64 `json:"limitLow,string"`
	HoldAmount   float64 `json:"hold_amount,string"`
	UnitAmount   float64 `json:"unitAmount,string"`
}

type DepthRecord struct {
	Price  float64
	Amount float64
}

type DepthRecords []DepthRecord

func (dr DepthRecords) Len() int {
	return len(dr)
}

func (dr DepthRecords) Swap(i, j int) {
	dr[i], dr[j] = dr[j], dr[i]
}

func (dr DepthRecords) Less(i, j int) bool {
	return dr[i].Price < dr[j].Price
}

type Depth struct {
	ContractType string `json:"contract_type,omitempty"` //for futures
	ContractId   string `json:"contract_id,omitempty"`   // for futures
	Pair         CurrencyPair
	UTime        time.Time
	AskList      DepthRecords // Descending order
	BidList      DepthRecords // Descending order
}

type APIConfig struct {
	HttpClient    *http.Client
	Endpoint      string
	ApiKey        string
	ApiSecretKey  string
	ApiPassphrase string //for okex.com v3 api
	ClientId      string //for bitstamp.net , huobi.pro

	Lever float64 //杠杆倍数 , for future
}

type Kline struct {
	Pair      CurrencyPair
	Timestamp int64
	Open      float64
	Close     float64
	High      float64
	Low       float64
	Vol       float64
}

type FutureKline struct {
	*Kline
	Vol2 float64 //个数
}

type FutureSubAccount struct {
	Currency      Currency
	AccountRights float64 //账户权益
	KeepDeposit   float64 //保证金
	ProfitReal    float64 //已实现盈亏
	ProfitUnreal  float64
	RiskRate      float64 //保证金率
}

type FutureAccount struct {
	FutureSubAccounts map[Currency]FutureSubAccount
}

type FutureOrder struct {
	ClientOid    string //自定义ID，GoEx内部自动生成
	OrderID2     string //请尽量用这个字段替代OrderID字段
	Price        float64
	Amount       float64
	AvgPrice     float64
	DealAmount   float64
	OrderID      int64 //deprecated
	OrderTime    int64
	Status       TradeStatus
	Currency     CurrencyPair
	OrderType    int     //ORDINARY=0 POST_ONLY=1 FOK= 2 IOC= 3
	OType        int     //1：开多 2：开空 3：平多 4： 平空
	LeverRate    float64 //倍数
	Fee          float64 //手续费
	ContractName string
	FinishedTime int64 // finished timestamp

	//策略委托单
	TriggerPrice float64
	AlgoType     int //1:限价 2:市场价；触发价格类型，默认是限价；为市场价时，委托价格不必填；
}

type FuturePosition struct {
	BuyAmount      float64
	BuyAvailable   float64
	BuyPriceAvg    float64
	BuyPriceCost   float64
	BuyProfitReal  float64
	BuyProfit      float64
	CreateDate     int64
	LeverRate      float64
	SellAmount     float64
	SellAvailable  float64
	SellPriceAvg   float64
	SellPriceCost  float64
	SellProfitReal float64
	SellProfit     float64
	Symbol         CurrencyPair //btc_usd:比特币,ltc_usd:莱特币
	ContractType   string
	ContractId     int64
	ForceLiquPrice float64 //预估爆仓价
	ShortPnlRatio  float64 //空仓收益率
	LongPnlRatio   float64 //多仓收益率
}

type HistoricalFunding struct {
	InstrumentId string    `json:"instrument_id"`
	RealizedRate float64   `json:"realized_rate,string"`
	FundingTime  time.Time `json:"funding_time"`
}

type TickSize struct {
	InstrumentID    string
	UnderlyingIndex string
	QuoteCurrency   string
	PriceTickSize   float64 //下单价格精度
	AmountTickSize  float64 //数量精度
}

type FuturesContractInfo struct {
	*TickSize
	ContractVal  float64 //合约面值(美元)
	Delivery     string  //交割日期
	ContractType string  //	本周 this_week 次周 next_week 季度 quarter
}

//api parameter struct

type BorrowParameter struct {
	CurrencyPair CurrencyPair
	Currency     Currency
	Amount       float64
}

type RepaymentParameter struct {
	BorrowParameter
	BorrowId string
}

type TransferParameter struct {
	Currency       string  `json:"currency"`
	From           int     `json:"from"`
	To             int     `json:"to"`
	Amount         float64 `json:"amount"`
	SubAccount     string  `json:"sub_account"`
	InstrumentId   string  `json:"instrument_id"`
	ToInstrumentId string  `json:"to_instrument_id"`
}

type WithdrawParameter struct {
	Currency    string  `json:"currency"`
	Amount      float64 `json:"amount,string"`
	Destination int     `json:"destination"` //提币到(2:OKCoin国际 3:OKEx 4:数字货币地址)
	ToAddress   string  `json:"to_address"`
	TradePwd    string  `json:"trade_pwd"`
	Fee         string  `json:"fee"`
}

type DepositWithdrawHistory struct {
	WithdrawalId string    `json:"withdrawal_id,omitempty"`
	Currency     string    `json:"currency"`
	Txid         string    `json:"txid"`
	Amount       float64   `json:"amount,string"`
	From         string    `json:"from,omitempty"`
	To           string    `json:"to"`
	Memo         string    `json:"memo,omitempty"`
	Fee          string    `json:"fee"`
	Status       int       `json:"status,string"`
	Timestamp    time.Time `json:"timestamp"`
}
