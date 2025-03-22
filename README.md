### Introduction
* Unified and standardize the interfaces of various cryptocurrency trading platforms.

![goex](goex_struct.png)

### Donate
* [BTC] 1GoEXwVvXG7kNdQSFaUNF35A3izHojLGxP
* [USDT-TRC20] TGoExC6xvzE4wSA9cYZnwcPaXEjibA5Vtc

### Example

```golang
package main

import (
	goexv2 "github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/logger"
	"github.com/nntaoli-project/goex/v2/model"
	"github.com/nntaoli-project/goex/v2/options"
	"log"
)

func main() {
	logger.SetLevel(logger.DEBUG)                             // set log level
	//goexv2.DefaultHttpCli.SetProxy("socks5://127.0.0.1:1080") //socks5 proxy
	goexv2.DefaultHttpCli.SetTimeout(5)                       // 5 second

	_, _, err := goexv2.OKx.Spot.GetExchangeInfo() //must call
	if err != nil {
		panic(err)
	}
	
	//This CurrencyPair structure contains various information about the trading pair, such as price_precision and qty_precision,min_qty ...
	btcUSDTCurrencyPair, err := goexv2.OKx.Spot.NewCurrencyPair(model.BTC, model.USDT)
	if err != nil {
		panic(err)
	}
		
	//public api
	log.Println(goexv2.OKx.Spot.GetTicker(btcUSDTCurrencyPair))

	//private api
	okxPrvApi := goexv2.OKx.Spot.NewPrvApi(
		options.WithApiKey(""), 
		options.WithApiSecretKey(""), 
		options.WithPassphrase(""))
	
	//create order api
	order, _, err := okxPrvApi.CreateOrder(btcUSDTCurrencyPair, 0.01, 18000, model.Spot_Buy, model.OrderType_Limit)
	log.Println(err)
	log.Println(order)
}
```

### FAQ
#### 1. okx simulated trading
See ([issues/238](https://github.com/nntaoli-project/goex/issues/238))

#### 2. Filled The OrderClientID 

```
ord, resp, err := okxPrvApi.CreateOrder(btcUSDTCurrencyPair, 0.01, 23000,
		model.Spot_Buy, model.OrderType_Limit,
		model.OptionParameter{}.OrderClientID("goex123027892")) //client id: goex123027892
```

### Thanks
<a href="https://www.jetbrains.com/?from=goex"><img src="https://account.jetbrains.com/static/images/jetbrains-logo-inv.svg" height="120" alt="JetBrains"/></a>

