#使用特定交易所的API可以获取更多的接口实现
```golang
var okex = NewOKEx(&goex.APIConfig{
                    	Endpoint: "https://www.okex.com",
                    	//HttpClient: &http.Client{
                    	//	Transport: &http.Transport{
                    	//		Proxy: func(req *http.Request) (*url.URL, error) {
                    	//			return &url.URL{
                    	//				Scheme: "socks5",
                    	//				Host:   "127.0.0.1:1080"}, nil
                    	//		},
                    	//	},
                    	//},
                    	ApiKey:        "",
                    	ApiSecretKey:  "",
                    	ApiPassphrase: "",
                    })
 var (
   okexSpot = okex.OKExSpot
   okexSwap = okex.OKExSwap   //永续合约实现
   okexFuture=okex.OKExFuture //交割合约实现
   okexWallet =okex.OKExWallet //资金账户（钱包）操作
   )
 
  //接口调用,更多接口调用请看代码
  log.Prinitln(okexSpot.GetAccount()) //获取账户资产信息
  //okexSpot.BatchPlaceOrders([]goex.Order{...}) //批量下单,单个交易对同时最大只能下10笔
  log.Println(okexSwap.GetFutureUserinfo()) //获取账户权益信息
  log.Println(okexFuture.GetFutureUserinfo())//获取账户权益信息
  
```
