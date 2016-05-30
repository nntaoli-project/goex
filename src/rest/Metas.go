package rest

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
	Last string;
	Buy  string;
	Sell string;
	High string;
	Low  string;
	Vol  string;
	Date string;
}

type DepthRecord struct {
	Price,
	Amount string;
}

type Depth struct {
	AskList,
	BidList []DepthRecord
}
