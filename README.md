### 介绍
* 统一并标准化各个数字资产交易平台的接口
* 某些功能组件做到可插拔化，方便开发者二次开发

![goex](goex_struct.png)
### 示例

```go
package main

import (
	"github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/binance"
	"github.com/nntaoli-project/goex/v2/logger"
	. "github.com/nntaoli-project/goex/v2/model"
)

func main() {
	logger.SetLevel(logger.DEBUG)
	goex.DefaultHttpCli.SetTimeout(15)
	goex.DefaultHttpCli.SetProxy("socks5://127.0.0.1:2220")
	
	baSpot := binance.Spot.MarketApi()
	ticker, err := baSpot.GetTicker(CurrencyPair{Symbol: "BTCUSDT"})
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Infof("%+v", ticker)
}
```