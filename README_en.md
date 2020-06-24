<div align="center">
<img width="409" heigth="205" src="https://upload-images.jianshu.io/upload_images/6760989-dec7dc747846880e.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240"  alt="goex">
</div>

### goex

goex project is designed to unify and standardize the interfaces of each digital asset trading platform. The same strategy can be switched to any trading platform at any time without changing any code.

[中文](https://github.com/nntaoli-project/goex/blob/dev/README.md)

### Exchanges are supported by goex `23+`
| Exchange | Market API | Order API | Version |   
| ---   | ---     | ---     | ---   |  
| huobi.pro | Y | Y | 1 |   
| hbdm.com | Y (REST / WS)| Y |  1 |  
| okex.com (spot/future)| Y (REST / WS) | Y | 1 |  
| okex.com (swap future) | Y | Y | 2 |
| binance.com | Y | Y | 1 |  
| kucoin.com | Y | Y | 1 |
| bitstamp.net | Y | Y | 1 |  
| bitfinex.com | Y | Y | 1 |  
| zb.com | Y | Y | 1 |  
| kraken.com | Y | Y | * |  
| poloniex.com | Y | Y | * |   
| big.one | Y | Y | 2\|3 | 
| hitbtc.com | Y | Y | * |
| coinex.com | Y | Y | 1 |
| exx.com | Y | Y | 1 |
| bithumb.com | Y | Y | * |
| gate.io | Y | N | 1 |
| bittrex.com | Y | N | 1.1 |

### Install goex
> go get   

``` go get github.com/nntaoli-project/goex ```
> go mod

``` 
require (
          github.com/nntaoli-project/goex latest
)
```

#Recommended(IMPORTANCE)
1. use GoLand development.
2. turn off the auto format function.
3. DONOT reformat existing files, which will result in a particularly bad commit.
4. use the OrderID2 field instead of the OrderID
5. can't use the deprecated field or method

donate & Buy a cup of Coffee for author
-----------------
BTC:13cBHLk6B7t3Uj7caJbCwv1UaiuiA6Qx8z

LTC:LVxM7y1K2dnpuNBU42ei3dKzPySf4VAm1H
 
ETH:0x98573ddb33cdddce480c3bc1f9279ccd88ca1e93