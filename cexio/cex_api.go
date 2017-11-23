package cexio

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	API_BASE_URL = "https://cex.io/api/"
)

type Cex struct {
	httpClient *http.Client
	Username   string
	Api_key    string
	Api_secret string
}

func (cex *Cex) Nonce() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func (cex *Cex) ToHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

func (cex *Cex) Signature() (string, string) {
	nonce := cex.Nonce()
	message := nonce + cex.Username + cex.Api_key
	signature := cex.ToHmac256(message, cex.Api_secret)
	return signature, nonce
}

func (cex *Cex) GetMethod(u string) ([]byte, error) {
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (cex *Cex) PostMethod(u string, v url.Values) ([]byte, error) {
	res, err := cex.httpClient.PostForm(u, v)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (cex *Cex) ApiCall(method string, id string, param map[string]string, opt string) []byte {
	var data []byte
	u := API_BASE_URL + method + "/"
	if len(opt) != 0 {
		u = u + opt + "/"
	}
	w := API_BASE_URL + "ghash.io/" + method
	// Post method for private method
	signature, nonce := cex.Signature()
	v := url.Values{}
	v.Set("key", cex.Api_key)
	v.Add("signature", signature)
	v.Add("nonce", nonce)
	// Place order param
	if len(param) != 0 {
		v.Add("type", param["ordertype"])
		v.Add("amount", param["amount"])
		v.Add("price", param["price"])
	}
	// Cancel order id
	if len(id) != 0 {
		v.Add("id", id)
	}
	v.Encode()
	if method == "workers" || method == "hashrate" {
		// Ghash.io post method
		data, _ = cex.PostMethod(w, v) // url ghash.io , param
	} else {
		// Cex.io post method
		data, _ = cex.PostMethod(u, v) // url cex.io, param
	}
	return data
}

func (cex *Cex) ApiCallPublic(method string, id string, param map[string]string, opt string) []byte {
	var data []byte
	u := API_BASE_URL + method + "/"
	if len(opt) != 0 {
		u = u + opt + "/"
	}
	// Get method for public method
	data, _ = cex.GetMethod(u)
	return data
}

// Public functions
func (cex *Cex) Ticker(opt string) []byte {
	return cex.ApiCallPublic("ticker", "", map[string]string{}, opt)
}

func (cex *Cex) OrderBook(size int, opt string) []byte {
	var data []byte
	u := API_BASE_URL + "order_book" + "/"
	if len(opt) != 0 {
		u = u + opt + "/?depth=" + strconv.Itoa(size)
	}
	data, _ = cex.GetMethod(u)
	return data
}

func (cex *Cex) TradeHistory(opt string) []byte {
	return cex.ApiCallPublic("trade_history", "", map[string]string{}, opt)
}

// Private functions
func (cex *Cex) Balance() []byte {
	return cex.ApiCall("balance", "", map[string]string{}, "")
}

func (cex *Cex) OpenOrders(opt string) []byte {
	return cex.ApiCall("open_orders", "", map[string]string{}, opt)
}

// Orders functions
func (cex *Cex) PlaceOrder(ordertype string, amount string, price string, opt string) []byte {
	var param = map[string]string{
		"ordertype": ordertype,
		"amount":    amount,
		"price":     price}
	return cex.ApiCall("place_order", "", param, opt)
}

func (cex *Cex) CancelOrderRaw(id string) []byte {
	return cex.ApiCall("cancel_order", id, map[string]string{}, "")
}

// Workers functions
func (cex *Cex) Hashrate() []byte {
	return cex.ApiCall("hashrate", "", map[string]string{}, "")
}

func (cex *Cex) Workers() []byte {
	return cex.ApiCall("workers", "", map[string]string{}, "")
}
