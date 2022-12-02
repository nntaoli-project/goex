### goex目标

goex项目是为了统一并标准化各个数字资产交易平台的接口而设计，同一个策略可以随时切换到任意一个交易平台，而不需要更改任何代码。

![goex](goex_struct.png)
### 示例

```go
    import (
        "github.com/nntaoli-projects/goex/v2"
    )
    goex.SetHttpTimeout(5)
    //goex.SetHttpProxy("socks5://127.0.0.1:2220")
    goex.SetupDefaultLibs()  //must need call
    marketApi := binance.Spot.MarketApi()
    tk,err := marketApi.GetTicker(goex.CurrencyPair{Symbol: "btcusdt"})
```