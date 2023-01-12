### Introduction

* 统一并标准化各个数字资产交易平台的接口
* 某些功能组件做到可插拔化，方便开发者二次开发

![goex](goex_struct.png)

### donate
* [BTC] 1GoEXwVvXG7kNdQSFaUNF35A3izHojLGxP
* [USDT-TRC20] TGoExC6xvzE4wSA9cYZnwcPaXEjibA5Vtc

### example

```
package main

import (
	goexv2 "github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/logger"
	"github.com/nntaoli-project/goex/v2/model"
	"github.com/nntaoli-project/goex/v2/options"
	"log"
)

func main() {
	logger.SetLevel(logger.DEBUG)                             //设置日志输出级别
	goexv2.DefaultHttpCli.SetProxy("socks5://127.0.0.1:1080") //socks代理
	goexv2.DefaultHttpCli.SetTimeout(5)                       // 5 second

	//共有api调用
	log.Println(goexv2.OKx.Spot.GetTicker(model.CurrencyPair{Symbol: "BTC-USDT"}))
	log.Println(goexv2.OKx.Futures.GetTicker(model.CurrencyPair{Symbol: "BTC-USDT-SWAP"}))

	//私有API调用
	okxPrvApi := goexv2.OKx.Spot.NewPrvApi(options.WithApiKey(""), options.WithApiSecretKey(""), options.WithPassphrase(""))
	order, _, err := okxPrvApi.CreateOrder(model.CurrencyPair{Symbol: "BTC-USDT"}, 0.01, 18000, model.Spot_Buy, model.OrderType_Limit)
	log.Println(err)
	log.Println(order)
}
```

