package rest

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