package okcoin

import (
	"encoding/json"
	"errors"
	. "github.com/nntaoli-project/GoEx"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type OKExSpot struct {
	OKCoinCN_API
}

func NewOKExSpot(client *http.Client, accesskey, secretkey string) *OKExSpot {
	return &OKExSpot{
		OKCoinCN_API{client, accesskey, secretkey, "https://www.okex.com/api/v1/"}}
}

func (ctx *OKExSpot) GetExchangeName() string {
	return "okex.com"
}

func (ctx *OKExSpot) GetAccount() (*Account, error) {
	postData := url.Values{}
	err := ctx.buildPostForm(&postData)
	if err != nil {
		return nil, err
	}

	body, err := HttpPostForm(ctx.client, ctx.api_base_url+url_userinfo, postData)
	if err != nil {
		return nil, err
	}

	var respMap map[string]interface{}

	err = json.Unmarshal(body, &respMap)
	if err != nil {
		return nil, err
	}

	log.Println(respMap)

	if !respMap["result"].(bool) {
		errcode := strconv.FormatFloat(respMap["error_code"].(float64), 'f', 0, 64)
		return nil, errors.New(errcode)
	}
	//log.Println(respMap)
	info, ok := respMap["info"].(map[string]interface{})
	if !ok {
		return nil, errors.New(string(body))
	}

	funds := info["funds"].(map[string]interface{})
	free := funds["free"].(map[string]interface{})
	freezed := funds["freezed"].(map[string]interface{})

	account := new(Account)
	account.Exchange = ctx.GetExchangeName()

	var (
		btcSubAccount  SubAccount
		ltcSubAccount  SubAccount
		ethSubAccount  SubAccount
		etcSubAccount  SubAccount
		bccSubAccount  SubAccount
		usdtSubAccount SubAccount
	)

	btcSubAccount.Currency = BTC
	btcSubAccount.Amount = ToFloat64(free["btc"])
	btcSubAccount.LoanAmount = 0
	btcSubAccount.ForzenAmount = ToFloat64(freezed["btc"])

	ltcSubAccount.Currency = LTC
	ltcSubAccount.Amount = ToFloat64(free["ltc"])
	ltcSubAccount.LoanAmount = 0
	ltcSubAccount.ForzenAmount = ToFloat64(freezed["ltc"])

	ethSubAccount.Currency = ETH
	ethSubAccount.Amount = ToFloat64(free["eth"])
	ethSubAccount.LoanAmount = 0
	ethSubAccount.ForzenAmount = ToFloat64(freezed["eth"])

	etcSubAccount.Currency = ETC
	etcSubAccount.Amount = ToFloat64(free["etc"])
	etcSubAccount.LoanAmount = 0
	etcSubAccount.ForzenAmount = ToFloat64(freezed["etc"])

	bccSubAccount.Currency = BCC
	bccSubAccount.Amount = ToFloat64(free["bcc"])
	bccSubAccount.LoanAmount = 0
	bccSubAccount.ForzenAmount = ToFloat64(freezed["bcc"])

	usdtSubAccount.Currency = USDT
	usdtSubAccount.Amount = ToFloat64(free["usdt"])
	usdtSubAccount.LoanAmount = 0
	usdtSubAccount.ForzenAmount = ToFloat64(freezed["usdt"])

	account.SubAccounts = make(map[Currency]SubAccount, 5)
	account.SubAccounts[BTC] = btcSubAccount
	account.SubAccounts[LTC] = ltcSubAccount
	account.SubAccounts[ETH] = ethSubAccount
	account.SubAccounts[ETC] = etcSubAccount
	account.SubAccounts[BCC] = bccSubAccount
	account.SubAccounts[USDT] = usdtSubAccount

	return account, nil
}
