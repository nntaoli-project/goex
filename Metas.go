package coinapi

import "net/http"

type Order struct {
	Price,
	Amount,
	AvgPrice,
	DealAmount float64;
	OrderID    int;
	OrderTime  int;
	Status     TradeStatus;
	Currency   CurrencyPair;
	Side       TradeSide;
}

type SubAccount struct {
	Currency   Currency;
	Amount,
	ForzenAmount,
	LoanAmount float64;
}

type Account struct {
	Exchange    string;
	Asset       float64; //总资产
	NetAsset    float64; //净资产
	SubAccounts map[Currency]SubAccount;
}

type Ticker struct {
	Last float64;
	Buy  float64;
	Sell float64;
	High float64;
	Low  float64;
	Vol  float64;
	Date uint64;
}

type DepthRecord struct {
	Price,
	Amount float64;
}

type Depth struct {
	AskList,
	BidList []DepthRecord
}

type APIConfig struct {
	HttpClient *http.Client;
	ApiUrl,
	AccessKey,
	SecretKey  string;
}

type Kline struct  {
	Timestamp int64;
	Open,
	Close,
	High,
	Low ,
	Vol float64;
}

type FutureSubAccount struct {
	Currency      Currency;
	AccountRights float64; //账户权益
	KeepDeposit   float64; //保证金
	ProfitReal    float64; //已实现盈亏
	ProfitUnreal  float64;
	RiskRate      float64; //保证金率
}

type FutureAccount struct {
	FutureSubAccounts map[Currency]FutureSubAccount;
}

type FutureOrder struct {
	Price      float64;
	Amount     float64;
	AvgPrice   float64;
	DealAmount float64;
	OrderID    int64;
	OrderTime  int64;
	Status     TradeStatus;
	Currency   CurrencyPair;
	OType      int;//1：开多 2：开空 3：平多 4： 平空
	LeverRate  int;     //倍数
	Fee        float64; //手续费
	ContractName string;
}

type FuturePosition struct {
	BuyAmount      float64;
	BuyAvailable   float64;
	BuyPriceAvg    float64;
	BuyPriceCost   float64;
	BuyProfitReal  float64;
	CreateDate     int64;
	LeverRate      int;
	SellAmount     float64;
	SellAvailable  float64;
	SellPriceAvg   float64;
	SellPriceCost  float64;
	SellProfitReal float64;
	Symbol         CurrencyPair; //btc_usd:比特币,ltc_usd:莱特币
	ContractType   string;
	ContractId     int64;
	ForceLiquPrice float64; //预估爆仓价
}