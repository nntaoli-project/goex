
# Websocket

```go

// new fcoin websocket client
var fcws = NewFCoinWs(&http.Client{
	Transport: &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse("socks5://127.0.0.1:1080")
			return nil, nil
		},
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
	},
	Timeout: 10 * time.Second,
})

func init() {
	fcws.ProxyUrl("socks5://127.0.0.1:1080") // set proxy if necessary
	fcws.SetCallbacks(printfTicker, printfDepth, printfTrade, nil) // set callback for your subscribe
}

// ticker handler
func printfTicker(ticker *goex.Ticker) {
	log.Println("ticker:", ticker)
}

// depth handler
func printfDepth(depth *goex.Depth) {
	log.Println("depth:", depth)
}
// trade handler
func printfTrade(trade *goex.Trade) {
	log.Println("trade:", trade)
}
// kline handler
func printfKline(kline *goex.Kline, period int) {
	log.Println("kline:", kline)
}

func main()  {
 	fcws.SubscribeTicker(goex.BTC_USDT)
	fcws.SubscribeDepth(goex.BTC_USDT, 20)
	fcws.SubscribeTrade(goex.BTC_USDT)
	time.Sleep(time.Second * 120) // sleep to watch the print from handler
}
```

# FCOIN SPOT

## Rest API

```go
// new fcoin rest client
var ft = NewFCoin(&http.Client{
	Transport: &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			//return url.Parse("socks5://127.0.0.1:1080") // set proxy if necessary
			return nil, nil
		},
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
	},
	Timeout: 10 * time.Second,
}, "yourApiKey", "yourSecretKey")

func main () {
    ft.GetTicker(goex.NewCurrencyPair2("BTC_USDT")) 
}

```

# FCOIN MARGIN

```go
// new fcoin rest client
var fm = &FCoinMargin{ft}

func main () {
    fm.GetTicker(goex.NewCurrencyPair2("BTC_USDT")) 
}

```