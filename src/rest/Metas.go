package rest

import "net/http"

type Order struct {
	OrderID,
	Price,
	Amount,
	AvgPrice,
	DealAmount string;
	OrderTime  int32;
	Status     TradeStatus;
	Currency   CurrencyPair;
	side       TradeSide;
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
	NetAsset   float64; //净资产
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
