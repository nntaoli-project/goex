package goex

import "strings"

type Currency struct {
	Symbol string
	Desc   string
}

func (c Currency) String() string {
	return c.Symbol
}

// A->B(A兑换为B)
type CurrencyPair struct {
	CurrencyA Currency
	CurrencyB Currency
}

var (
	UNKNOWN = Currency{"UNKNOWN", ""}
	CNY     = Currency{"CNY", "rmb （China Yuan)"}
	USD     = Currency{"USD", "USA dollar"}
	USDT    = Currency{"USDT","http://tether.io"}
	EUR     = Currency{"EUR", ""}
	KRW     = Currency{"KRW", ""}
	JPY     = Currency{"JPY", "japanese yen"}
	BTC     = Currency{"BTC", "bitcoin.org"}
	XBT     = Currency{"XBT", "bitcoin.org"}
	BCC     = Currency{"BCC", "bitcoin-abc"}
	BCH     = Currency{"BCH","bitcoin-abc"}
	BCX     = Currency{"BCX",""}
	LTC     = Currency{"LTC", "litecoin.org"}
	ETH     = Currency{"ETH", ""}
	ETC     = Currency{"ETC", ""}
	EOS     = Currency{"EOS", ""}
	BTS     = Currency{"BTS", ""}
	QTUM    = Currency{"QTUM", ""}
	SC      = Currency{"SC", "sia.tech"}
	ANS     = Currency{"ANS", "www.antshares.org"}
	ZEC     = Currency{"ZEC", ""}
	DCR     = Currency{"DCR", ""}
	XRP     = Currency{"XRP", ""}
	BTG     = Currency{"BTG",""}
	BCD     = Currency{"BCD",""}
	NEO     = Currency{"NEO","neo.org"}
	HSR     = Currency{"HSR",""}

	//currency pair

	BTC_CNY  = CurrencyPair{BTC, CNY}
	LTC_CNY  = CurrencyPair{LTC, CNY}
	BCC_CNY  = CurrencyPair{BCC, CNY}
	ETH_CNY  = CurrencyPair{ETH, CNY}
	ETC_CNY  = CurrencyPair{ETC, CNY}
	EOS_CNY  = CurrencyPair{EOS, CNY}
	BTS_CNY  = CurrencyPair{BTS, CNY}
	QTUM_CNY = CurrencyPair{QTUM, CNY}
	SC_CNY   = CurrencyPair{SC, CNY}
	ANS_CNY  = CurrencyPair{ANS, CNY}
	ZEC_CNY  = CurrencyPair{ZEC, CNY}

	BTC_KRW  = CurrencyPair{BTC,KRW}
	ETH_KRW  = CurrencyPair{ETH,KRW}
	ETC_KRW  = CurrencyPair{ETC,KRW}
	LTC_KRW  = CurrencyPair{LTC,KRW}
	BCH_KRW  = CurrencyPair{BCH,KRW}

	BTC_USD = CurrencyPair{BTC, USD}
	LTC_USD = CurrencyPair{LTC, USD}
	ETH_USD = CurrencyPair{ETH, USD}
	ETC_USD = CurrencyPair{ETC, USD}
	BCH_USD = CurrencyPair{BCH, USD}
	BCC_USD = CurrencyPair{BCC, USD}
	XRP_USD = CurrencyPair{XRP, USD}
	BCD_USD = CurrencyPair{BCD,USD}

	BTC_USDT = CurrencyPair{BTC,USDT}
	LTC_USDT = CurrencyPair{LTC, USDT}
	BCH_USDT = CurrencyPair{BCH, USDT}
	BCC_USDT = CurrencyPair{BCC,USDT}
	ETC_USDT = CurrencyPair{ETC,USDT}
	ETH_USDT = CurrencyPair{ETH,USDT}
	BCD_USDT = CurrencyPair{BCD,USDT}
	NEO_USDT = CurrencyPair{NEO,USDT}
	EOS_USDT = CurrencyPair{EOS,USDT}
	XRP_USDT = CurrencyPair{XRP,USDT}
	HSR_USDT = CurrencyPair{HSR,USDT}

	XRP_EUR = CurrencyPair{XRP , EUR}

	BTC_JPY = CurrencyPair{BTC,JPY}
	LTC_JPY = CurrencyPair{LTC,JPY}
	ETH_JPY = CurrencyPair{ETH,JPY}
	ETC_JPY = CurrencyPair{ETC,JPY}
	BCH_JPY = CurrencyPair{BCH,JPY}

	LTC_BTC = CurrencyPair{LTC, BTC}
	ETH_BTC = CurrencyPair{ETH, BTC}
	ETC_BTC = CurrencyPair{ETC, BTC}
	BCC_BTC = CurrencyPair{BCC, BTC}
	BCH_BTC = CurrencyPair{BCH, BTC}
	DCR_BTC = CurrencyPair{DCR, BTC}
	XRP_BTC = CurrencyPair{XRP ,BTC}
	BTG_BTC = CurrencyPair{BTG,BTC}
	BCD_BTC = CurrencyPair{BCD,BTC}
	NEO_BTC = CurrencyPair{NEO,BTC}
	EOS_BTC = CurrencyPair{EOS,BTC}
	HSR_BTC = CurrencyPair{HSR,BTC}

	ETC_ETH = CurrencyPair{ETC, ETH}
	EOS_ETH = CurrencyPair{EOS, ETH}
	ZEC_ETH = CurrencyPair{ZEC,ETH}
	NEO_ETH = CurrencyPair{NEO,ETH}
	HSR_ETH = CurrencyPair{HSR,ETH}

	UNKNOWN_PAIR = CurrencyPair{UNKNOWN,UNKNOWN}
)

func (c CurrencyPair) String() string {
	return c.ToSymbol("_")
}

func NewCurrency(symbol, desc string) Currency {
	switch symbol {
	case "cny" , "CNY":
		return CNY
	case "usdt", "USDT":
		return USDT
	case "usd", "USD":
		return USD
	case "jpy" , "JPY":
		return JPY
	case "krw", "KRW":
		return KRW
	case "eur", "EUR":
		return EUR
	case "btc", "BTC":
		return BTC
	case "xbt", "XBT":
		return XBT
	case "bch", "BCH":
		return BCH
	case "bcc", "BCC":
		return BCC
	case "ltc", "LTC":
		return LTC
	case "sc", "SC":
		return SC
	case "ans", "ANS":
		return ANS
	case "neo", "NEO":
		return NEO
	default:
		return Currency{strings.ToUpper(symbol) , desc}
	}
}

func NewCurrencyPair(currencyA Currency, currencyB Currency) CurrencyPair {
	return CurrencyPair{currencyA, currencyB}
}

func NewCurrencyPair2(currencyPairSymbol string) CurrencyPair  {
	currencys := strings.Split(currencyPairSymbol , "_")
	if len(currencys) == 2 {
		return CurrencyPair{NewCurrency(currencys[0] , "") ,
		NewCurrency(currencys[1] , "")}
	}
	return UNKNOWN_PAIR
}

func (pair CurrencyPair) ToSymbol(joinChar string) string {
	return strings.Join([]string{pair.CurrencyA.Symbol, pair.CurrencyB.Symbol}, joinChar)
}

func (pair CurrencyPair) ToSymbol2(joinChar string) string {
	return strings.Join([]string{pair.CurrencyB.Symbol, pair.CurrencyA.Symbol}, joinChar)
}
