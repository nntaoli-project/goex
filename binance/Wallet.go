package binance

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/Jameslu041/goex"
	"github.com/Jameslu041/goex/internal/logger"
	"net/url"
)

type Wallet struct {
	ba   *Binance
	conf *APIConfig
}

func NewWallet(c *APIConfig) *Wallet {
	return &Wallet{ba: NewWithConfig(c), conf: c}
}

func (w *Wallet) GetAccount() (*Account, error) {
	return nil, errors.New("not implement")
}

func (w *Wallet) Withdrawal(param WithdrawParameter) (withdrawId string, err error) {
	return "", errors.New("not implement")
}

func (w *Wallet) Transfer(param TransferParameter) error {
	transferUrl := w.conf.Endpoint + "/sapi/v1/futures/transfer"

	postParam := url.Values{}
	postParam.Set("asset", param.Currency)
	postParam.Set("amount", fmt.Sprint(param.Amount))

	if param.From == SPOT && param.To == SWAP_USDT {
		postParam.Set("type", "1")
	}

	if param.From == SWAP_USDT && param.To == SPOT {
		postParam.Set("type", "2")
	}

	if param.From == SPOT && param.To == FUTURE {
		postParam.Set("type", "3")
	}

	if param.From == FUTURE && param.To == SPOT {
		postParam.Set("type", "4")
	}

	w.ba.buildParamsSigned(&postParam)

	resp, err := HttpPostForm2(w.ba.httpClient, transferUrl, postParam,
		map[string]string{"X-MBX-APIKEY": w.ba.accessKey})

	if err != nil {
		return err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return err
	}

	if respmap["tranId"] != nil && ToInt64(respmap["tranId"]) > 0 {
		return nil
	}

	return errors.New(string(resp))
}

func (w *Wallet) GetWithDrawHistory(currency *Currency) ([]DepositWithdrawHistory, error) {
	//historyUrl := w.conf.Endpoint + "/wapi/v3/withdrawHistory.html"
	historyUrl := w.conf.Endpoint + "/sapi/v1/accountSnapshot"
	postParam := url.Values{}
	postParam.Set("type", "SPOT")
	w.ba.buildParamsSigned(&postParam)

	resp, err := HttpGet5(w.ba.httpClient, historyUrl+"?"+postParam.Encode(),
		map[string]string{"X-MBX-APIKEY": w.ba.accessKey})

	if err != nil {
		return nil, err
	}
	logger.Debugf("response body: %s", string(resp))
	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (w *Wallet) GetDepositHistory(currency *Currency) ([]DepositWithdrawHistory, error) {
	historyUrl := w.conf.Endpoint + "/wapi/v3/depositHistory.html"
	postParam := url.Values{}
	postParam.Set("asset", currency.Symbol)
	w.ba.buildParamsSigned(&postParam)

	resp, err := HttpGet5(w.ba.httpClient, historyUrl+"?"+postParam.Encode(),
		map[string]string{"X-MBX-APIKEY": w.ba.accessKey})

	if err != nil {
		return nil, err
	}
	logger.Debugf("response body: %s", string(resp))
	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
