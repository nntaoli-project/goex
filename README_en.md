<div align="center">
<img width="409" heigth="205" src="https://upload-images.jianshu.io/upload_images/6760989-dec7dc747846880e.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240"  alt="goex">
</div>

### goex

goex project is designed to unify and standardize the interfaces of each digital asset trading platform. The same strategy can be switched to any trading platform at any time without changing any code.

[中文](https://github.com/nntaoli-project/goex/blob/dev/README.md)

### Exchanges are supported by goex `22+`
| Exchange | Market API | Order API | Version |   
| ---   | ---     | ---     | ---   |  
| hbg.com | Y | Y | 1 |   
| hbdm.com | Y (REST / WS)| Y |  1 |  
| okex.com (spot/future)| Y (REST / WS) | Y | 1 |  
| okex.com (swap future) | Y | Y | 2 |
| binance.com | Y | Y | 1 |  
| bitstamp.net | Y | Y | 1 |  
| bitfinex.com | Y | Y | 1 |  
| zb.com | Y | Y | 1 |  
| kraken.com | Y | Y | * |  
| poloniex.com | Y | Y | * |  
| aacoin.com | Y | Y | 1 |   
| allcoin.ca | Y | Y | * |  
| big.one | Y | Y | 2\|3 | 
| fcoin.com | Y (REST / WS) | Y | 2 |  
| hitbtc.com | Y | Y | * |
| coinex.com | Y | Y | 1 |
| exx.com | Y | Y | 1 |
| bithumb.com | Y | Y | * |
| gate.io | Y | N | 1 |
| btcbox.co.jp | Y | N | * |
| bittrex.com | Y | N | 1.1 |
| btcchina.com | Y | Y | 1 |
| coinbig.com | Y | Y | * |

### Install goex
``` go get github.com/nntaoli-project/goex ```

### Example
```golang

   package main
   
   import (
   	"github.com/nntaoli-project/goex"
   	"github.com/nntaoli-project/goex/builder"
   	"log"
   	"time"
   )
   
   func main() {
   	apiBuilder := builder.NewAPIBuilder().HttpTimeout(5 * time.Second)
   	//apiBuilder := builder.NewAPIBuilder().HttpTimeout(5 * time.Second).HttpProxy("socks5://127.0.0.1:1080")
   	
   	//build spot api
   	//api := apiBuilder.APIKey("").APISecretkey("").ClientID("123").Build(goex.BITSTAMP)
   	api := apiBuilder.APIKey("").APISecretkey("").Build(goex.HUOBI_PRO)
   	log.Println(api.GetExchangeName())
   	log.Println(api.GetTicker(goex.BTC_USD))
   	log.Println(api.GetDepth(2, goex.BTC_USD))
   	//log.Println(api.GetAccount())
   	//log.Println(api.GetUnfinishOrders(goex.BTC_USD))
   
   	//build future api
   	futureApi := apiBuilder.APIKey("").APISecretkey("").BuildFuture(goex.HBDM)
   	log.Println(futureApi.GetExchangeName())
   	log.Println(futureApi.GetFutureTicker(goex.BTC_USD, goex.QUARTER_CONTRACT))
   	log.Println(futureApi.GetFutureDepth(goex.BTC_USD, goex.QUARTER_CONTRACT, 5))
   	//log.Println(futureApi.GetFutureUserinfo()) // account
   	//log.Println(futureApi.GetFuturePosition(goex.BTC_USD , goex.QUARTER_CONTRACT))//position info
   }

```

### websocket Example
```golang
import (
	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/huobi"
	//"github.com/nntaoli-project/goex/okcoin"
	"log"
)

func main() {

	//ws := okcoin.NewOKExFutureWs() //ok future
	ws := huobi.NewHbdmWs() //huobi future
	//setup callback
	ws.SetCallbacks(func(ticker *goex.FutureTicker) {
		log.Println(ticker)
	}, func(depth *goex.Depth) {
		log.Println(depth)
	}, func(trade *goex.Trade, contract string) {
		log.Println(contract, trade)
	})
	//subscribe
	ws.SubscribeTrade(goex.BTC_USDT, goex.NEXT_WEEK_CONTRACT)
	ws.SubscribeDepth(goex.BTC_USDT, goex.QUARTER_CONTRACT, 5)
	ws.SubscribeTicker(goex.BTC_USDT, goex.QUARTER_CONTRACT)
}  

```

### More Detail

[goex.TOP](https://goex.top)

# Highly Recommended(IMPORTANCE)
1. use GoLand development.
2. turn off the auto format function.
3. DONOT reformat existing files, which will result in a particularly bad commit.
4. use the OrderID2 field instead of the OrderID

### How to find us
Join QQ group: [574829125](#)

-----------------

### Buy me a Coffe

<img src="https://raw.githubusercontent.com/nntaoli-project/goex/dev/wx_pay.JPG" width="250" alt="Buy me a Coffe">&nbsp;&nbsp;&nbsp;<img src="https://raw.githubusercontent.com/nntaoli-project/goex/dev/IMG_1177.jpg" width="250" alt="Buy me a Coffe">
