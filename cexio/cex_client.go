package cexio

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
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

// Public functions
func (cex *Cex) ticker(pair string) (map[string]interface{}, error) {
	path := API_BASE_URL + "ticker" + "/" + pair
	return cex.doRequest(path)
}

func (cex *Cex) orderBook(size int, pair string) (map[string]interface{}, error) {
	u := API_BASE_URL + "order_book" + "/" + pair + "/?depth=" + strconv.Itoa(size)
	return cex.doRequest(u)
}

func (cex *Cex) tradeHistory(opt string) (map[string]interface{}, error) {
	u := API_BASE_URL + "trade_history" + "/" + opt
	return cex.doRequest(u)
}

// Private functions
func (cex *Cex) balance() (map[string]interface{}, error) {
	u := API_BASE_URL + "balance" + "/"
	return cex.doAuthenticatedRequest(u, map[string]string{})
}

func (cex *Cex) openOrders(pair string) (map[string]interface{}, error) {
	u := API_BASE_URL + "open_orders" + "/" + pair
	return cex.doAuthenticatedRequest(u, map[string]string{})
}

// Orders functions
func (cex *Cex) placeOrder(ordertype string, amount string, price string, pair string) (map[string]interface{}, error) {
	u := API_BASE_URL + "place_order" + "/" + pair
	var param = map[string]string{
		"ordertype": ordertype,
		"amount":    amount,
		"price":     price}
	return cex.doAuthenticatedRequest(u, param)
}

func (cex *Cex) cancelOrder(id string, pair string) (map[string]interface{}, error) {
	u := API_BASE_URL + "cancel_order" + "/" + pair
	var param = map[string]string{}
	if id != "" {
		param["id"] = id
	}
	return cex.doAuthenticatedRequest(u, param)
}

func (cex *Cex) nonce() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func (cex *Cex) toHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

func (cex *Cex) signature() (string, string) {
	nonce := cex.nonce()
	message := nonce + cex.userId + cex.Api_key
	signature := cex.toHmac256(message, cex.Api_secret)
	return signature, nonce
}

func (cex *Cex) doRequest(u string) (map[string]interface{}, error) {
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	fmt.Println(string(data))
	var ret map[string]interface{}
	err = json.Unmarshal(data, &ret)
	return ret, err
}

func (cex *Cex) doAuthenticatedRequest(path string, param map[string]string) (map[string]interface{}, error) {
	signature, nonce := cex.signature()
	v := url.Values{}
	v.Set("key", cex.Api_key)
	v.Add("signature", signature)
	v.Add("nonce", nonce)
	for k, val := range param {
		v.Add(k, val)
	}
	res, err := cex.httpClient.PostForm(path, v)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	var ret map[string]interface{}
	err = json.Unmarshal(data, &ret)
	return ret, err
}
