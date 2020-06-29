package goex

type WalletApi interface {
	//获取钱包资产
	GetAccount() (*Account, error)
	//提币
	Withdrawal(param WithdrawParameter) (withdrawId string, err error)
	//划转资产
	Transfer(param TransferParameter) error
	//获取提币记录
	GetWithDrawHistory(currency *Currency) ([]DepositWithdrawHistory, error)
	//获取充值记录
	GetDepositHistory(currency *Currency) ([]DepositWithdrawHistory, error)
}
