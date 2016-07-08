package rest

type FutureRestAPI interface {
	/**
	 *获取交易所名字
	 */
	GetExchangeName() string

	/**
	 *获取交割预估价
	 */
	GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64 , error)

	/**
	 * 期货行情
	 * @param currency_pair   btc_usd:比特币    ltc_usd :莱特币
	 * @param contractType  合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
	 */
	GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error)

	/**
	 * 期货深度
	 * @param currencyPair  btc_usd:比特币    ltc_usd :莱特币
	 * @param contractType  合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
	 * @return
	 */
	GetFutureDepth(currencyPair CurrencyPair, contractType string) (*Depth, error)

	/**
	 * 期货指数
	 * @param currencyPair   btc_usd:比特币    ltc_usd :莱特币
	 */
	GetFutureIndex(currencyPair CurrencyPair) (float64, error)

	/**
	 *全仓账户
	 */
	GetFutureUserinfo() (*FutureAccount, error)

	/**
	 * 期货下单
	 * @param currencyPair   btc_usd:比特币    ltc_usd :莱特币
	 * @param contractType   合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
	 * @param price  价格
	 * @param amount  委托数量
	 * @param openType   1:开多   2:开空   3:平多   4:平空
	 * @param matchPrice  是否为对手价 0:不是    1:是   ,当取值为1时,price无效
	 */
	PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string , openType, matchPrice , leverRate int) (string, error)

	/**
	 * 取消订单
	 * @param symbol   btc_usd:比特币    ltc_usd :莱特币
	 * @param contractType    合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
	 * @param orderId   订单ID

	 */
	FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error)

	/**
	 * 用户持仓查询
	 * @param symbol   btc_usd:比特币    ltc_usd :莱特币
	 * @param contractType   合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
	 * @return
	 */
	GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error)

	/**
	 *获取订单信息
	 */
	GetFutureOrders(orderId int64, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error)

	/**
	 *获取未完成订单信息
	 */
	GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error)
}
