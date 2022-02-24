<div align="center">
<img width="409" heigth="205" src="https://upload-images.jianshu.io/upload_images/6760989-dec7dc747846880e.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240"  alt="goex">
<img width="198" height="205" src="https://upload-images.jianshu.io/upload_images/6760989-81f29f7a5dbd9bb6.jpg" >
</div>

### goex目标

goex项目是为了统一并标准化各个数字资产交易平台的接口而设计，同一个策略可以随时切换到任意一个交易平台，而不需要更改任何代码。

[English](https://github.com/nntaoli-project/goex/blob/dev/README_en.md)

### wiki文档

[文档](https://github.com/nntaoli-project/goex/wiki)

### goex已支持交易所 `23+`

| 交易所 | 行情接口 | 交易接口 | 版本号 |   
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

### 安装goex库  
> go get

``` go get github.com/nntaoli-project/goex ```

>建议go mod 管理依赖
``` 
require (
          github.com/nntaoli-project/goex latest
)
```

### 注意事项

1. 推荐使用GoLand开发。
2. 推荐关闭自动格式化功能,代码请使用go fmt 格式化.
3. 不建议对现已存在的文件进行重新格式化，这样会导致commit特别糟糕。
4. 请用OrderID2这个字段代替OrderID
5. 请不要使用deprecated关键字标注的方法和字段，后面版本可能随时删除的
-----------------

donate
-----------------
BTC:bc1qm4fz8vg78yr75syclch9zfnwa0x00efy95972s

LTC:MPgBjHmecACXDKH3KdmLR6X8mmp5YgAEZS
 
ETH:0xDEa8C4B6B36294B00b05B6A6A5b63d8aA52A1bF7

### 欢迎为作者付一碗面钱

![微信](wx_pay.JPG) ![支付宝](IMG_1177.jpg)  