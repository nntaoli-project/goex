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
	Currency     CurrencyPair;
	Amount       string;
	ForzenAmount string;
	LoanAmount   string;
}

type Account struct {
	Exchange    string;
	SubAccounts []SubAccount
}

type Ticker struct {
	Last float64;
	Buy  float64;
	Sell float64;
	High float64;
	Low  float64;
	Vol  string;
	Date int;
}

type DepthRecord struct {
	Price,
	Amount string;
}

type Depth struct {
	AskList,
	BidList []DepthRecord
}

type APIConfig struct {
	HttpClient *http.Client;
	ApiUrl,
	AccessKey,
	SecretKey string;
}
