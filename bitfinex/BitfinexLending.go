package bitfinex

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
	"io/ioutil"
	"strconv"
	"strings"
)

type LendBookItem struct {
	Rate      float64 `json:",string"`
	Amount    float64 `json:",string"`
	Period    int     `json:"period"`
	Timestamp string  `json:"timestamp"`
	Frr       string  `json:"frr"`
}

type LendBook struct {
	Bids []LendBookItem `json:"bids"`
	Asks []LendBookItem `json:"asks"`
}

type LendOrder struct {
	Id              int     `json:"id"`
	Currency        string  `json:"currency"`
	Rate            float64 `json:"rate,string"`
	Period          int     `json:"period"`
	Direction       string  `json:"direction"`
	IsLive          bool    `json:"is_live"`
	IsCancelled     bool    `json:"is_cancelled"`
	Amount          float64 `json:"amount,string"`
	ExecutedAmount  float64 `json:"executed_amount,string"`
	RemainingAmount float64 `json:"remaining_amount,string"`
	OriginalAmount  float64 `json:"original_amount,string"`
	Timestamp       string  `json:"timestamp"`
}

type LendTicker struct {
	Ticker
	Coin            Currency
	DailyChangePerc float64
}

func (bfx *Bitfinex) GetLendTickers() ([]LendTicker, error) {

	resp, err := bfx.httpClient.Get("https://api.bitfinex.com/v2/tickers?symbols=ALL")
	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var ret []interface{}
	json.Unmarshal(body, &ret)

	var tickers []LendTicker

	for _, v := range ret {
		vv := v.([]interface{})
		symbol := vv[0].(string)
		if strings.HasPrefix(symbol, "f") {
			tickers = append(tickers, LendTicker{
				Ticker: Ticker{
					Last: ToFloat64(vv[10]) * 100,
					Vol:  ToFloat64(vv[11])},
				DailyChangePerc: ToFloat64(vv[9]) * 100,
				Coin:            NewCurrency(symbol[1:], "")})
		}
	}

	return tickers, nil
}

func (bfx *Bitfinex) GetDepositWalletBalance() (*Account, error) {
	wallets, err := bfx.GetWalletBalances()
	if err != nil {
		return nil, err
	}
	return wallets["deposit"], nil
}

func (bfx *Bitfinex) GetLendBook(currency Currency) (error, *LendBook) {
	path := fmt.Sprintf("/lendbook/%s", currency.Symbol)
	resp, err := bfx.httpClient.Get(BASE_URL + path)
	if err != nil {
		return err, nil
	}

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("HttpCode: %d , errmsg: %s", resp.StatusCode, string(body))), nil
	}
	//println(string(body))
	var lendBook LendBook
	err = json.Unmarshal(body, &lendBook)
	if err != nil {
		return err, nil
	}

	return nil, &lendBook
}

func (bfx *Bitfinex) Transfer(amount float64, currency Currency, fromWallet, toWallet string) error {
	path := "transfer"
	params := map[string]interface{}{
		"amount":     strconv.FormatFloat(amount, 'f', -1, 32),
		"currency":   strings.ToUpper(currency.Symbol),
		"walletfrom": fromWallet,
		"walletto":   toWallet,
	}

	var resp []map[string]interface{}

	err := bfx.doAuthenticatedRequest("POST", path, params, &resp)
	if err != nil {
		return err
	}

	if "success" == resp[0]["status"] {
		return nil
	}

	return errors.New(resp[0]["message"].(string))
}

func (bfx *Bitfinex) newOffer(currency Currency, amount, rate string, period int, direction string) (error, *LendOrder) {
	path := "offer/new"
	params := map[string]interface{}{
		"amount":    amount,
		"currency":  currency.Symbol,
		"rate":      rate,
		"period":    period,
		"direction": direction,
	}

	var lendOrder LendOrder
	err := bfx.doAuthenticatedRequest("POST", path, params, &lendOrder)
	if err != nil {
		return err, nil
	}

	return nil, &lendOrder
}

func (bfx *Bitfinex) NewLendOrder(currency Currency, amount, rate string, period int) (error, *LendOrder) {
	return bfx.newOffer(currency, amount, rate, period, "lend")
}

func (bfx *Bitfinex) NewLoanOrder(currency Currency, amount, rate string, period int) (error, *LendOrder) {
	return bfx.newOffer(currency, amount, rate, period, "loan")
}

func (bfx *Bitfinex) CancelLendOrder(id int) (error, *LendOrder) {
	println("id=", id)
	path := "offer/cancel"
	var lendOrder LendOrder
	err := bfx.doAuthenticatedRequest("POST", path, map[string]interface{}{"offer_id": id}, &lendOrder)
	if err != nil {
		return err, nil
	}
	return nil, &lendOrder
}

func (bfx *Bitfinex) GetLendOrderStatus(id int) (error, *LendOrder) {
	path := "offer/status"
	var lendOrder LendOrder
	err := bfx.doAuthenticatedRequest("POST", path, map[string]interface{}{"offer_id": id}, &lendOrder)
	if err != nil {
		return err, nil
	}
	return nil, &lendOrder
}

func (bfx *Bitfinex) ActiveLendOrders() (error, []LendOrder) {
	var lendOrders []LendOrder
	err := bfx.doAuthenticatedRequest("POST", "offers", map[string]interface{}{}, &lendOrders)
	if err != nil {
		return err, nil
	}
	return nil, lendOrders
}

func (bfx *Bitfinex) OffersHistory(limit int) (error, []LendOrder) {
	var offerOrders []LendOrder
	err := bfx.doAuthenticatedRequest("POST", "offers/hist", map[string]interface{}{"limit": limit}, &offerOrders)
	if err != nil {
		return err, nil
	}
	return nil, offerOrders
}

func (bfx *Bitfinex) ActiveCredits() (error, []LendOrder) {
	var offerOrders []LendOrder
	err := bfx.doAuthenticatedRequest("POST", "credits", map[string]interface{}{}, &offerOrders)
	if err != nil {
		return err, nil
	}
	return nil, offerOrders
}

type TradeFunding struct {
	Rate      string `json:"rate"`
	Period    string `json:"period"`
	Amount    string `json:"amount"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Tid       int64  `json:"tid"`
	OfferId   int64  `json:"offer_id"`
}

func (bfx *Bitfinex) MytradesFunding(currency Currency, limit int) (error, []TradeFunding) {
	var trades []TradeFunding
	err := bfx.doAuthenticatedRequest("POST", "mytrades_funding", map[string]interface{}{"limit_trades": limit, "symbol": currency.Symbol}, &trades)
	if err != nil {
		return err, nil
	}
	return nil, trades
}
