package binance

import (
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
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
	return nil, errors.New("not implement")
}

func (w *Wallet) GetDepositHistory(currency *Currency) ([]DepositWithdrawHistory, error) {
	return nil, errors.New("not implement")
}
