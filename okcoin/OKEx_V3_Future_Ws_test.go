package okcoin

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"sync"
	"github.com/nntaoli-project/GoEx"
)


func newOKExV3FutureWs()*OKExV3FutureWs{
	okV3Ws := NewOKExV3FutureWs(okexV3)
	okV3Ws.WsUrl("wss://okexcomreal.bafang.com:10442/ws/v3")
	okV3Ws.ErrorHandleFunc(func(err error) {
		log.Println(err)
	})
	return okV3Ws
}

var (
	okV3Ws = newOKExV3FutureWs()
)

func TestOKExV3FutureWsTickerCallback(t *testing.T) {
	n := 10
	tickers := make([]goex.FutureTicker, 0)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	okV3Ws.TickerCallback(func(ticker *goex.FutureTicker) {
		t.Log(ticker, ticker.Ticker)
		if len(tickers) <= n {
			tickers = append(tickers, *ticker)
		}
		if len(tickers) == n {
			wg.Done()
		}
	})
	okV3Ws.SubscribeTicker(goex.EOS_USD, goex.QUARTER_CONTRACT)
	okV3Ws.SubscribeTicker(goex.EOS_USD, goex.SWAP_CONTRACT)
	wg.Wait()
}

func TestOKExV3FutureWsDepthCallback(t *testing.T) {
	n := 10
	depths := make([]goex.Depth, 0)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	okV3Ws.DepthCallback(func(depth *goex.Depth) {
		if len(depths) <= n {
			t.Log(depth)
			depths = append(depths, *depth)
		} 
		if len(depths) == n {
			wg.Done()
		}
	})
	okV3Ws.SubscribeDepth(goex.EOS_USD, goex.QUARTER_CONTRACT, 5)
	okV3Ws.SubscribeDepth(goex.EOS_USD, goex.SWAP_CONTRACT, 5)
	wg.Wait()
}

func TestOKExV3FutureWsTradeCallback(t *testing.T) {
	n := 10
	trades := make([]goex.Trade, 0)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	okV3Ws.TradeCallback(func(trade *goex.Trade, contractType string) {
		if len(trades) <= n {
			t.Log(contractType, trade)
			trades = append(trades, *trade)
		} 
		if len(trades) == n {
			wg.Done()
		}
	})
	okV3Ws.SubscribeTrade(goex.EOS_USD, goex.QUARTER_CONTRACT)
	okV3Ws.SubscribeTrade(goex.EOS_USD, goex.SWAP_CONTRACT)
	wg.Wait()
}

func TestOKExV3FutureWsLogin(t *testing.T) {
	if authed {
		okV3Ws := newOKExV3FutureWs()
		err := okV3Ws.Login("", apiSecretKey, passphrase) // fail
		assert.True(t, err != nil)
		okV3Ws = newOKExV3FutureWs()
		err = okV3Ws.Login(apiKey, apiSecretKey, passphrase) //succes
		assert.Nil(t, err)
		err = okV3Ws.Login(apiKey, apiSecretKey, passphrase) //duplicate login
		assert.Nil(t, err)
	} else {
		t.Log("not authed, skip test websocket login")
	}
}

func placeAndCancel(currencyPair goex.CurrencyPair, contractType string) {
	leverage := 20
	depth := 100
	dep, _ := okexV3.GetFutureDepth(currencyPair, contractType, depth)
	price := fmt.Sprintf("%f", dep.BidList[depth-1].Price)
	symbol, _ := okexV3.GetContract(currencyPair, contractType)
	amount := symbol.getSizeIncrement()
	orderID, err := okexV3.PlaceFutureOrder(
		currencyPair, contractType, price, amount, goex.OPEN_BUY, 0, leverage)
	if err != nil {
		log.Println(err)
	}
	_, err = okexV3.FutureCancelOrder(currencyPair, contractType, orderID)
	if err != nil {
		log.Println(err)
	}
}

func TestOKExV3FutureWsOrderCallback(t *testing.T) {
	if authed {
		err := okV3Ws.Login(apiKey, apiSecretKey, passphrase)
		assert.Nil(t, err)
		n := 4
		orders := make([]goex.FutureOrder, 0)
		wg := &sync.WaitGroup{}
		wg.Add(1)
		okV3Ws.OrderCallback(func(order *goex.FutureOrder, contractType string) {
			if len(orders) <= n {
				t.Log(contractType, order)
				orders = append(orders, *order)
			} 
			if len(orders) == n {
				wg.Done()
			}
		})
		err = okV3Ws.SubscribeOrder(goex.EOS_USD, goex.QUARTER_CONTRACT)
		assert.Nil(t, err)
		err = okV3Ws.SubscribeOrder(goex.EOS_USD, goex.SWAP_CONTRACT)
		assert.Nil(t, err)
		placeAndCancel(goex.EOS_USD, goex.QUARTER_CONTRACT)
		placeAndCancel(goex.EOS_USD, goex.SWAP_CONTRACT)
		wg.Wait()
	}
}