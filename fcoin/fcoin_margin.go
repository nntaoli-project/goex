package fcoin

import (
	"errors"
	"net/url"
	"strings"

	. "github.com/nntaoli-project/GoEx"
)

/**
杠杆交易区借币，
pair ： 操作的交易对
currency： 需要借的币种
amount : 借的金额
*/
type MarginOrder struct {
	Currency             Currency
	Amount               float64 // borrow amount
	BorrowTime           int64   // borrow time, ms
	RepaymentTime        int64   // repayment time, ms
	LendingFee           float64 // LendingFee
	LoanRate             float64 // loan rate
	LoanFeeStartTime     int64   // charge from start time
	LastRepayTime        int64   // last repay time
	LoanBillId           string
	NextLoanFeeStartTime int64
	State                string //账单状态. submitted 已提交; 2 confirmed 已确认; 5 finished 还款完成
	UnPaidAmount         float64
	UnPaidLendingFee     float64
}

type FCoinMargin struct {
	*FCoin
}

func (fm *FCoinMargin) GetMarginAccount(currency CurrencyPair) (*MarginAccount, error) {
	params := url.Values{}
	params.Set("account_type", strings.ToLower(currency.AdaptUsdToUsdt().ToSymbol("")))

	r, err := fm.doAuthenticatedRequest2("GET", "broker/leveraged_accounts/account", params)
	if err != nil {
		return nil, err
	}
	ok, isOk := r["status"].(string)
	if !isOk || ok != "ok" {
		return nil, errors.New("response status error")
	}
	//fmt.Println(r)
	acc := MarginAccount{}
	acc.Sub = make(map[Currency]MarginSubAccount, 2)
	//	"open": true,                                    #是否已经开通该类型杠杆账户. true:已开通;false:未开通
	//	"leveraged_account_type": "btcusdt",             #杠杆账户类型
	//	"state": "open"                                  #账户状态. close 已关闭;open 已开通-未发生借贷;normal 已借贷-风险率正常;blow_up 已爆仓;overrun 已穿仓", allowableValues = "close,open,normal,blow_up,overrun")
	response, isOk := r["data"].(map[string]interface{})
	if !isOk {
		return nil, errors.New("response data error")
	}

	acc.LiquidationPrice = ToFloat64(response["blow_up_price"])
	acc.RiskRate = ToFloat64(response["risk_rate"])
	//acc.MarginRatio = ToFloat64(response["margin_ratio"])

	c := NewCurrency(response["base"].(string), "")
	acc.Sub[c] = MarginSubAccount{
		//Balance:     ToFloat64(response["balance"]),
		Frozen:      ToFloat64(response["frozen_base_currency_amount"]),
		Available:   ToFloat64(response["available_base_currency_amount"]),
		CanWithdraw: ToFloat64(response["available_base_currency_loan_amount"]),
		//Loan:        ToFloat64(response["borrowed"]),
		//LendingFee:  ToFloat64(response["lending_fee"]),
	}

	c = NewCurrency(response["quote"].(string), "")
	acc.Sub[c] = MarginSubAccount{
		//Balance:     ToFloat64(response["balance"]),
		Frozen:      ToFloat64(response["frozen_quote_currency_amount"]),
		Available:   ToFloat64(response["available_quote_currency_amount"]),
		CanWithdraw: ToFloat64(response["available_quote_currency_loan_amount"]),
		//Loan:        ToFloat64(response["borrowed"]),
		//LendingFee:  ToFloat64(response["lending_fee"]),
	}

	return &acc, nil
}

func (fm *FCoinMargin) Borrow(parameter BorrowParameter) (*MarginOrder, error) {
	params := url.Values{}
	params.Set("account_type", strings.ToLower(parameter.CurrencyPair.AdaptUsdToUsdt().ToSymbol("")))
	params.Set("currency", strings.ToLower(parameter.Currency.String()))
	params.Set("amount", FloatToString(parameter.Amount, 8))
	params.Set("loan_type", "normal")

	r, err := fm.doAuthenticatedRequest("POST", "broker/leveraged/loans", params)
	if err != nil {
		return nil, err
	}
	//fmt.Println(r, err)
	response := r.(map[string]interface{})
	return &MarginOrder{
		Currency:             NewCurrency(response["currency"].(string), ""),
		Amount:               ToFloat64(response["amount"]),
		BorrowTime:           ToInt64(response["created_at"]),
		RepaymentTime:        ToInt64(response["finished_at"]),
		LoanFeeStartTime:     ToInt64(response["interest_start_at"]),
		LastRepayTime:        ToInt64(response["last_repayment_at"]),
		NextLoanFeeStartTime: ToInt64(response["next_interest_at"]),
		LendingFee:           ToFloat64(response["interest"]),
		LoanRate:             ToFloat64(response["interest_rate"]),
		UnPaidAmount:         ToFloat64(response["unpaid_amount"]),
		UnPaidLendingFee:     ToFloat64(response["unpaid_interest"]),
		LoanBillId:           response["loan_bill_id"].(string),
		State:                response["state"].(string),
	}, nil
}

func (fm *FCoinMargin) Repayment(parameter RepaymentParameter) (repaymentId string, err error) {
	params := url.Values{}
	params.Set("amount", FloatToString(parameter.Amount, 8))

	response, err := fm.doAuthenticatedRequest("POST", "broker/leveraged/repayments/"+parameter.BorrowId, params)
	if err != nil {
		return "", err
	}
	//fmt.Println("Repayment", response)
	repaymentId = response.(string)
	return
}

func (fm *FCoinMargin) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return fm.PlaceOrder(ORDER_TYPE_LIMIT, "buy", amount, price, currency, true)
}

func (fm *FCoinMargin) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return fm.PlaceOrder(ORDER_TYPE_LIMIT, "sell", amount, price, currency, true)
}

func (fm *FCoinMargin) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return fm.PlaceOrder(ORDER_TYPE_MARKET, "buy", amount, price, currency, true)
}

func (fm *FCoinMargin) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return fm.PlaceOrder(ORDER_TYPE_MARKET, "sell", amount, price, currency, true)
}

func (fm *FCoinMargin) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToLower(currency.AdaptUsdToUsdt().ToSymbol("")))
	params.Set("states", "submitted,partial_filled")
	params.Set("account_type", "margin")
	//params.Set("before", "1")
	//params.Set("after", "0")
	params.Set("limit", "100")

	r, err := fm.doAuthenticatedRequest("GET", "orders", params)
	if err != nil {
		return nil, err
	}

	var ords []Order

	for _, ord := range r.([]interface{}) {
		ords = append(ords, *fm.toOrder(ord.(map[string]interface{}), currency))
	}

	return ords, nil
}

func (fm *FCoinMargin) GetUnfinishLoans(currency CurrencyPair) ([]*MarginOrder, error) {
	params := url.Values{}
	params.Set("account_type", currency.ToSymbol(""))
	params.Set("skip_finish", "true")
	response, err := fm.doAuthenticatedRequest("GET", "broker/leveraged/loans/", params)
	if err != nil {
		return nil, err
	}

	var loanOrders []*MarginOrder
	data := response.(map[string]interface{})
	ctt := data["content"].([]interface{})
	for _, c := range ctt {
		content := c.(map[string]interface{})
		order := &MarginOrder{
			Currency:             NewCurrency(content["currency"].(string), ""),
			Amount:               ToFloat64(content["amount"]),
			BorrowTime:           ToInt64(content["created_at"]),
			RepaymentTime:        ToInt64(content["finished_at"]),
			LoanFeeStartTime:     ToInt64(content["interest_start_at"]),
			LastRepayTime:        ToInt64(content["last_repayment_at"]),
			NextLoanFeeStartTime: ToInt64(content["next_interest_at"]),
			LendingFee:           ToFloat64(content["interest"]),
			LoanRate:             ToFloat64(content["interest_rate"]),
			UnPaidAmount:         ToFloat64(content["unpaid_amount"]),
			UnPaidLendingFee:     ToFloat64(content["unpaid_interest"]),
			LoanBillId:           content["loan_bill_id"].(string),
			State:                content["state"].(string),
		}
		loanOrders = append(loanOrders, order)
	}
	return loanOrders, nil
}

func (fm *FCoinMargin) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToLower(currency.AdaptUsdToUsdt().ToSymbol("")))
	params.Set("states", "partial_canceled,filled")
	//params.Set("before", "1")
	//params.Set("after", "0")
	params.Set("limit", "100")
	params.Set("account_type", "margin")

	r, err := fm.doAuthenticatedRequest("GET", "orders", params)
	if err != nil {
		return nil, err
	}
	var ords []Order

	for _, ord := range r.([]interface{}) {
		ords = append(ords, *fm.toOrder(ord.(map[string]interface{}), currency))
	}

	return ords, nil

}

func (fm *FCoinMargin) GetOneLoan(borrowId string) (*MarginOrder, error) {
	params := url.Values{}
	params.Set("leveraged_loan_id", borrowId)

	response, err := fm.doAuthenticatedRequest2("GET", "broker/leveraged/loans/"+borrowId, params)
	if err != nil {
		return nil, err
	}
	return &MarginOrder{
		Currency:             NewCurrency(response["currency"].(string), ""),
		Amount:               ToFloat64(response["amount"]),
		BorrowTime:           ToInt64(response["created_at"]),
		RepaymentTime:        ToInt64(response["finished_at"]),
		LoanFeeStartTime:     ToInt64(response["interest_start_at"]),
		LastRepayTime:        ToInt64(response["last_repayment_at"]),
		NextLoanFeeStartTime: ToInt64(response["next_interest_at"]),
		LendingFee:           ToFloat64(response["interest"]),
		LoanRate:             ToFloat64(response["interest_rate"]),
		UnPaidAmount:         ToFloat64(response["unpaid_amount"]),
		UnPaidLendingFee:     ToFloat64(response["unpaid_interest"]),
		LoanBillId:           response["loan_bill_id"].(string),
		State:                response["state"].(string),
	}, nil

}

// currency币种名称：usdt、btc、eth
// from资产来源账户类型: exchange: 交易账户; assets: 资产账户
// to目标账户类型: leveraged_btcusdt、leveraged_ethusdt、leveraged_eosusdt、leveraged_xrpusdt
func (fm *FCoinMargin) AssetTransferIn(currency Currency, amount, from string, to CurrencyPair) (bool, error) {
	params := url.Values{}
	params.Set("currency", strings.ToLower(currency.String()))
	params.Set("amount", amount)
	params.Set("source_account_type", from)
	params.Set("target_account_type", "leveraged_"+strings.ToLower(to.ToSymbol("")))
	r, err := fm.doAuthenticatedRequest2("POST", "broker/leveraged/assets/transfer/in", params)
	if err != nil {
		return false, err
	}
	ok, isOk := r["status"].(string)
	if !isOk || ok != "ok" {
		return false, errors.New("response status error")
	}

	return true, nil
}

// currency币种名称：usdt、btc、eth
// from资产来源账户类型: exchange: 交易账户; assets: 资产账户
// to目标账户类型: leveraged_btcusdt、leveraged_ethusdt、leveraged_eosusdt、leveraged_xrpusdt
func (fm *FCoinMargin) AssetTransferOut(currency Currency, amount string, from CurrencyPair, to string) (bool, error) {
	params := url.Values{}
	params.Set("currency", strings.ToLower(currency.String()))
	params.Set("amount", amount)
	params.Set("source_account_type", "leveraged_"+strings.ToLower(from.ToSymbol("")))
	params.Set("target_account_type", to)
	r, err := fm.doAuthenticatedRequest2("POST", "broker/leveraged/assets/transfer/out", params)
	if err != nil {
		return false, err
	}
	ok, isOk := r["status"].(string)
	if !isOk || ok != "ok" {
		return false, errors.New("response status error")
	}

	return true, nil
}
