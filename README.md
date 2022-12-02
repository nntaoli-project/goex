### goex目标

goex项目是为了统一并标准化各个数字资产交易平台的接口而设计，同一个策略可以随时切换到任意一个交易平台，而不需要更改任何代码。

![goex](goex_struct.png)
### 示例

```go
package main

import (
	. "github.com/nntaoli-project/goex/v2/model"
	"github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/huobi"
	"github.com/nntaoli-project/goex/v2/logger"
)

func main() {
	logger.SetLevel(logger.DEBUG)
	goex.DefaultHttpCli.SetTimeout(15)
	goex.DefaultHttpCli.SetProxy("socks5://127.0.0.1:2220")

	baSpot := huobi.Spot.MarketApi()
	ticker, err := baSpot.GetTicker(CurrencyPair{Symbol: "btcusdt"})

	if err != nil {
		logger.Error(err)
		return
	}

	logger.Infof("%+v", ticker)
}

```