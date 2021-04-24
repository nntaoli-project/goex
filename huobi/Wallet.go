package huobi

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/Jameslu041/goex"
	"github.com/Jameslu041/goex/internal/logger"
	"net/url"
	"strings"
)

type Wallet struct {
	pro *HuoBiPro
}

func NewWallet(c *APIConfig) *Wallet {
	return &Wallet{pro: NewHuobiWithConfig(c)}
}

//获取钱包资产
func (w *Wallet) GetAccount() (*Account, error) {
	return nil, errors.New("not implement")
}

func (w *Wallet) Withdrawal(param WithdrawParameter) (withdrawId string, err error) {
	return "", errors.New("not implement")
}

func (w *Wallet) Transfer(param TransferParameter) error {
	if param.From == SUB_ACCOUNT || param.To == SUB_ACCOUNT ||
		param.From == SPOT_MARGIN || param.To == SPOT_MARGIN {
		return errors.New("not implements")
	}

	httpParam := url.Values{}
	httpParam.Set("currency", strings.ToLower(param.Currency))
	httpParam.Set("amount", FloatToString(param.Amount, 8))

	path := ""

	if (param.From == SPOT && param.To == FUTURE) ||
		(param.From == FUTURE && param.To == SPOT) {
		path = "/v1/futures/transfer"
	}

	if param.From == SWAP || param.From == SWAP_USDT ||
		param.To == SWAP || param.To == SWAP_USDT {
		path = "/v2/account/transfer"
	}

	if param.From == SPOT && param.To == FUTURE {
		httpParam.Set("type", "pro-to-futures")
	}

	if param.From == FUTURE && param.To == SPOT {
		httpParam.Set("type", "futures-to-pro")
	}

	if param.From == SPOT && param.To == SWAP {
		httpParam.Set("from", "spot")
		httpParam.Set("to", "swap")
	}

	if param.From == SPOT && param.To == SWAP_USDT {
		httpParam.Set("currency", "usdt")
		httpParam.Set("from", "spot")
		httpParam.Set("to", "linear-swap")
		httpParam.Set("margin-account", fmt.Sprintf("%s-usdt", strings.ToLower(param.Currency)))
	}

	if param.From == SWAP && param.To == SPOT {
		httpParam.Set("from", "swap")
		httpParam.Set("to", "spot")
	}

	if param.From == SWAP_USDT && param.To == SPOT {
		httpParam.Set("currency", "usdt")
		httpParam.Set("from", "linear-swap")
		httpParam.Set("to", "spot")
		httpParam.Set("margin-account",
			fmt.Sprintf("%s-usdt", strings.ToLower(param.Currency)))
	}

	w.pro.buildPostForm("POST", path, &httpParam)

	postJsonParam, _ := ValuesToJson(httpParam)
	responseBody, err := HttpPostForm3(w.pro.httpClient,
		fmt.Sprintf("%s%s?%s", w.pro.baseUrl, path, httpParam.Encode()),
		string(postJsonParam),
		map[string]string{"Content-Type": "application/json", "Accept-Language": "zh-cn"})

	if err != nil {
		return err
	}

	logger.Debugf("[response body] %s", string(responseBody))

	var responseRet map[string]interface{}

	err = json.Unmarshal(responseBody, &responseRet)
	if err != nil {
		return err
	}

	if responseRet["status"] != nil &&
		responseRet["status"].(string) == "ok" {
		return nil
	}

	if responseRet["code"] != nil && responseRet["code"].(float64) == 200 {
		return nil
	}

	return errors.New(string(responseBody))
}

func (w *Wallet) GetWithDrawHistory(currency *Currency) ([]DepositWithdrawHistory, error) {
	return nil, errors.New("not implement")
}

func (w *Wallet) GetDepositHistory(currency *Currency) ([]DepositWithdrawHistory, error) {
	return nil, errors.New("not implement")
}
