package okcoin

import (
	"encoding/json"
	"errors"
	. "github.com/nntaoli-project/goex"
	"net/http"
	"net/url"
	"strconv"
)

type OKExSpot struct {
	OKCoinCN_API
}

func NewOKExSpot(client *http.Client, accesskey, secretkey string) *OKExSpot {
	return &OKExSpot{
		OKCoinCN_API: OKCoinCN_API{client, accesskey, secretkey, "https://www.okex.com/api/v1/"}}
}

func (ctx *OKExSpot) GetExchangeName() string {
	return OKEX
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

	if errcode, isok := respMap["error_code"].(float64); isok {
		errcodeStr := strconv.FormatFloat(errcode, 'f', 0, 64)
		return nil, errors.New(errcodeStr)
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

	account.SubAccounts = make(map[Currency]SubAccount, 6)
	for k, v := range free {
		currencyKey := NewCurrency(k, "")
		subAcc := SubAccount{
			Currency:     currencyKey,
			Amount:       ToFloat64(v),
			ForzenAmount: ToFloat64(freezed[k])}
		account.SubAccounts[currencyKey] = subAcc
	}

	return account, nil
}
