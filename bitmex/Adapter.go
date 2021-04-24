package bitmex

import (
	"fmt"
	. "github.com/Jameslu041/goex"
	"strings"
)

func AdaptCurrencyPairToSymbol(pair CurrencyPair, contract string) string {
	if contract == "" || contract == SWAP_CONTRACT {
		if pair.CurrencyA.Eq(BTC) {
			pair = NewCurrencyPair(XBT, USD)
		}
		if pair.CurrencyB.Eq(BTC) {
			pair = NewCurrencyPair(pair.CurrencyA, XBT)
		}
		return pair.AdaptUsdtToUsd().ToSymbol("")
	}

	coin := pair.CurrencyA.Symbol
	if pair.CurrencyA.Eq(BTC) {
		coin = XBT.Symbol
	}
	return fmt.Sprintf("%s%s", coin, strings.ToUpper(contract))
}

func AdaptWsSymbol(symbol string) (pair CurrencyPair, contract string) {
	symbol = strings.ToUpper(symbol)

	if symbol == "XBTCUSD" {
		return BTC_USD, SWAP_CONTRACT
	}

	if symbol == "BCHUSD" {
		return BCH_USD, SWAP_CONTRACT
	}

	if symbol == "ETHUSD" {
		return ETH_USD, SWAP_CONTRACT
	}

	if symbol == "LTCUSD" {
		return LTC_USD, SWAP_CONTRACT
	}

	if symbol == "LINKUSDT" {
		return NewCurrencyPair2("LINK_USDT"), SWAP_CONTRACT
	}

	pair = NewCurrencyPair(NewCurrency(symbol[0:3], ""), USDT)
	contract = symbol[3:]
	if pair.CurrencyA.Eq(XBT) {
		return NewCurrencyPair(BTC, USDT), contract
	}

	return pair, contract
}
