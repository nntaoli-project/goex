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
   )
 
  //接口调用,更多接口调用请看代码
  log.Prinitln(okexSpot.GetAccount()) //获取账户资产信息
  log.Println(okexSwap.GetFutureUserinfo()) //获取账户权益信息
  log.Println(okexFuture.GetFutureUserinfo())//获取账户权益信息
  
```
