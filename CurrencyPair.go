package goex

import "strings"

type Currency struct {
	Symbol string
	Desc   string
}

func (c Currency) String() string {
	return c.Symbol
}

func (c Currency) Eq(c2 Currency) bool {
	return c.Symbol == c2.Symbol
}

// A->B(A兑换为B)
type CurrencyPair struct {
	CurrencyA      Currency
	CurrencyB      Currency
	AmountTickSize int // 下单量精度
	PriceTickSize  int //交易对价格精度
}

var (
	UNKNOWN = Currency{"UNKNOWN", ""}
	CNY     = Currency{"CNY", ""}
	USD     = Currency{"USD", ""}
	USDT    = Currency{"USDT", ""}
	PAX     = Currency{"PAX", "https://www.paxos.com/"}
	USDC    = Currency{"USDC", "https://www.centre.io/"}
	EUR     = Currency{"EUR", ""}
	KRW     = Currency{"KRW", ""}
	JPY     = Currency{"JPY", ""}
	BTC     = Currency{"BTC", "https://bitcoin.org/"}
	XBT     = Currency{"XBT", ""}
	BCC     = Currency{"BCC", ""}
	BCH     = Currency{"BCH", ""}
	BCX     = Currency{"BCX", ""}
	LTC     = Currency{"LTC", ""}
	ETH     = Currency{"ETH", ""}
	ETC     = Currency{"ETC", ""}
	EOS     = Currency{"EOS", ""}
	BTS     = Currency{"BTS", ""}
	QTUM    = Currency{"QTUM", ""}
	SC      = Currency{"SC", ""}
	ANS     = Currency{"ANS", ""}
	ZEC     = Currency{"ZEC", ""}
	DCR     = Currency{"DCR", ""}
	XRP     = Currency{"XRP", ""}
	BTG     = Currency{"BTG", ""}
	BCD     = Currency{"BCD", ""}
	NEO     = Currency{"NEO", ""}
	HSR     = Currency{"HSR", ""}
	BSV     = Currency{"BSV", ""}
	OKB     = Currency{"OKB", "OKB is a global utility token issued by OK Blockchain Foundation"}
	HT      = Currency{"HT", "HuoBi Token"}
	BNB     = Currency{"BNB", "BNB, or Binance Coin, is a cryptocurrency created by Binance."}
	TRX     = Currency{"TRX", ""}
	GBP     = Currency{"GBP", ""}
	XLM     = Currency{"XLM", ""}
	DOT     = Currency{"DOT", ""}
	DASH    = Currency{"DASH", ""}
	CRV     = Currency{"CRV", ""}
	ALGO    = Currency{"ALGO", ""}
	DOGE    = Currency{"DOGE", ""}

	//currency pair
	BTC_KRW = CurrencyPair{CurrencyA: BTC, CurrencyB: KRW, AmountTickSize: 2, PriceTickSize: 1}
	ETH_KRW = CurrencyPair{CurrencyA: ETH, CurrencyB: KRW, AmountTickSize: 2, PriceTickSize: 2}
	ETC_KRW = CurrencyPair{CurrencyA: ETC, CurrencyB: KRW, AmountTickSize: 2, PriceTickSize: 2}
	LTC_KRW = CurrencyPair{CurrencyA: LTC, CurrencyB: KRW, AmountTickSize: 2, PriceTickSize: 2}
	BCH_KRW = CurrencyPair{CurrencyA: BCH, CurrencyB: KRW, AmountTickSize: 2, PriceTickSize: 2}

	BTC_USD  = CurrencyPair{CurrencyA: BTC, CurrencyB: USD, AmountTickSize: 2, PriceTickSize: 1}
	LTC_USD  = CurrencyPair{CurrencyA: LTC, CurrencyB: USD, AmountTickSize: 2, PriceTickSize: 2}
	ETH_USD  = CurrencyPair{CurrencyA: ETH, CurrencyB: USD, AmountTickSize: 2, PriceTickSize: 2}
	ETC_USD  = CurrencyPair{CurrencyA: ETC, CurrencyB: USD, AmountTickSize: 2, PriceTickSize: 2}
	BCH_USD  = CurrencyPair{CurrencyA: BCH, CurrencyB: USD, AmountTickSize: 2, PriceTickSize: 2}
	XRP_USD  = CurrencyPair{CurrencyA: XRP, CurrencyB: USD, AmountTickSize: 3, PriceTickSize: 3}
	BCD_USD  = CurrencyPair{CurrencyA: BCD, CurrencyB: USD, AmountTickSize: 2, PriceTickSize: 3}
	EOS_USD  = CurrencyPair{CurrencyA: EOS, CurrencyB: USD, AmountTickSize: 2, PriceTickSize: 2}
	BTG_USD  = CurrencyPair{CurrencyA: BTG, CurrencyB: USD, AmountTickSize: 2, PriceTickSize: 2}
	BSV_USD  = CurrencyPair{CurrencyA: BSV, CurrencyB: USD, AmountTickSize: 2, PriceTickSize: 2}
	DOT_USD  = CurrencyPair{CurrencyA: DOT, CurrencyB: USD, AmountTickSize: 3, PriceTickSize: 2}
	DASH_USD = CurrencyPair{CurrencyA: DASH, CurrencyB: USD, AmountTickSize: 2, PriceTickSize: 2}
	CRV_USD  = CurrencyPair{CurrencyA: CRV, CurrencyB: USD, AmountTickSize: 4, PriceTickSize: 3}
	ALGO_USD = CurrencyPair{CurrencyA: ALGO, CurrencyB: USD, AmountTickSize: 4, PriceTickSize: 4}

	BTC_USDT  = CurrencyPair{CurrencyA: BTC, CurrencyB: USDT, AmountTickSize: 2, PriceTickSize: 1}
	LTC_USDT  = CurrencyPair{CurrencyA: LTC, CurrencyB: USDT, AmountTickSize: 2, PriceTickSize: 2}
	BCH_USDT  = CurrencyPair{CurrencyA: BCH, CurrencyB: USDT, AmountTickSize: 2, PriceTickSize: 2}
	ETC_USDT  = CurrencyPair{CurrencyA: ETC, CurrencyB: USDT, AmountTickSize: 2, PriceTickSize: 3}
	ETH_USDT  = CurrencyPair{CurrencyA: ETH, CurrencyB: USDT, AmountTickSize: 2, PriceTickSize: 2}
	BCD_USDT  = CurrencyPair{CurrencyA: BCD, CurrencyB: USDT, AmountTickSize: 2, PriceTickSize: 2}
	NEO_USDT  = CurrencyPair{CurrencyA: NEO, CurrencyB: USDT, AmountTickSize: 2, PriceTickSize: 2}
	EOS_USDT  = CurrencyPair{CurrencyA: EOS, CurrencyB: USDT, AmountTickSize: 2, PriceTickSize: 2}
	XRP_USDT  = CurrencyPair{CurrencyA: XRP, CurrencyB: USDT, AmountTickSize: 2, PriceTickSize: 2}
	HSR_USDT  = CurrencyPair{CurrencyA: HSR, CurrencyB: USDT, AmountTickSize: 2, PriceTickSize: 2}
	BSV_USDT  = CurrencyPair{CurrencyA: BSV, CurrencyB: USDT, AmountTickSize: 2, PriceTickSize: 2}
	OKB_USDT  = CurrencyPair{CurrencyA: OKB, CurrencyB: USDT, AmountTickSize: 2, PriceTickSize: 2}
	HT_USDT   = CurrencyPair{CurrencyA: HT, CurrencyB: USDT, AmountTickSize: 2, PriceTickSize: 4}
	BNB_USDT  = CurrencyPair{CurrencyA: BNB, CurrencyB: USDT, AmountTickSize: 2, PriceTickSize: 2}
	PAX_USDT  = CurrencyPair{CurrencyA: PAX, CurrencyB: USDT, AmountTickSize: 2, PriceTickSize: 3}
	TRX_USDT  = CurrencyPair{CurrencyA: TRX, CurrencyB: USDT, AmountTickSize: 2, PriceTickSize: 3}
	DOT_USDT  = CurrencyPair{CurrencyA: DOT, CurrencyB: USDT, AmountTickSize: 3, PriceTickSize: 2}
	DASH_USDT = CurrencyPair{CurrencyA: DASH, CurrencyB: USDT, AmountTickSize: 2, PriceTickSize: 2}
	CRV_USDT  = CurrencyPair{CurrencyA: CRV, CurrencyB: USDT, AmountTickSize: 3, PriceTickSize: 3}
	ALGO_USDT = CurrencyPair{CurrencyA: ALGO, CurrencyB: USDT, AmountTickSize: 3, PriceTickSize: 4}
	DOGE_USDT = CurrencyPair{CurrencyA: DOGE, CurrencyB: USDT, AmountTickSize: 3, PriceTickSize: 4}

	XRP_EUR = CurrencyPair{CurrencyA: XRP, CurrencyB: EUR, AmountTickSize: 2, PriceTickSize: 4}

	BTC_JPY = CurrencyPair{CurrencyA: BTC, CurrencyB: JPY, AmountTickSize: 2, PriceTickSize: 0}
	LTC_JPY = CurrencyPair{CurrencyA: LTC, CurrencyB: JPY, AmountTickSize: 2, PriceTickSize: 0}
	ETH_JPY = CurrencyPair{CurrencyA: ETH, CurrencyB: JPY, AmountTickSize: 2, PriceTickSize: 0}
	ETC_JPY = CurrencyPair{CurrencyA: ETC, CurrencyB: JPY, AmountTickSize: 2, PriceTickSize: 0}
	BCH_JPY = CurrencyPair{CurrencyA: BCH, CurrencyB: JPY, AmountTickSize: 2, PriceTickSize: 0}

	LTC_BTC = CurrencyPair{CurrencyA: LTC, CurrencyB: BTC, AmountTickSize: 2, PriceTickSize: 4}
	ETH_BTC = CurrencyPair{CurrencyA: ETH, CurrencyB: BTC, AmountTickSize: 2, PriceTickSize: 4}
	ETC_BTC = CurrencyPair{CurrencyA: ETC, CurrencyB: BTC, AmountTickSize: 2, PriceTickSize: 4}
	BCC_BTC = CurrencyPair{CurrencyA: BCC, CurrencyB: BTC, AmountTickSize: 2, PriceTickSize: 4}
	BCH_BTC = CurrencyPair{CurrencyA: BCH, CurrencyB: BTC, AmountTickSize: 2, PriceTickSize: 4}
	DCR_BTC = CurrencyPair{CurrencyA: DCR, CurrencyB: BTC, AmountTickSize: 2, PriceTickSize: 4}
	XRP_BTC = CurrencyPair{CurrencyA: XRP, CurrencyB: BTC, AmountTickSize: 2, PriceTickSize: 6}
	BTG_BTC = CurrencyPair{CurrencyA: BTG, CurrencyB: BTC, AmountTickSize: 2, PriceTickSize: 4}
	BCD_BTC = CurrencyPair{CurrencyA: BCD, CurrencyB: BTC, AmountTickSize: 2, PriceTickSize: 4}
	NEO_BTC = CurrencyPair{CurrencyA: NEO, CurrencyB: BTC, AmountTickSize: 2, PriceTickSize: 4}
	EOS_BTC = CurrencyPair{CurrencyA: EOS, CurrencyB: BTC, AmountTickSize: 2, PriceTickSize: 5}
	HSR_BTC = CurrencyPair{CurrencyA: HSR, CurrencyB: BTC, AmountTickSize: 2, PriceTickSize: 4}
	BSV_BTC = CurrencyPair{CurrencyA: BSV, CurrencyB: BTC, AmountTickSize: 2, PriceTickSize: 4}
	OKB_BTC = CurrencyPair{CurrencyA: OKB, CurrencyB: BTC, AmountTickSize: 2, PriceTickSize: 6}
	HT_BTC  = CurrencyPair{CurrencyA: HT, CurrencyB: BTC, AmountTickSize: 2, PriceTickSize: 7}
	BNB_BTC = CurrencyPair{CurrencyA: BNB, CurrencyB: BTC, AmountTickSize: 2, PriceTickSize: 6}
	TRX_BTC = CurrencyPair{CurrencyA: TRX, CurrencyB: BTC, AmountTickSize: 2, PriceTickSize: 7}
	DOT_BTC = CurrencyPair{CurrencyA: DOT, CurrencyB: BTC, AmountTickSize: 3, PriceTickSize: 6}

	ETC_ETH = CurrencyPair{CurrencyA: ETC, CurrencyB: ETH, AmountTickSize: 2, PriceTickSize: 4}
	EOS_ETH = CurrencyPair{CurrencyA: EOS, CurrencyB: ETH, AmountTickSize: 2, PriceTickSize: 4}
	ZEC_ETH = CurrencyPair{CurrencyA: ZEC, CurrencyB: ETH, AmountTickSize: 2, PriceTickSize: 4}
	NEO_ETH = CurrencyPair{CurrencyA: NEO, CurrencyB: ETH, AmountTickSize: 2, PriceTickSize: 4}
	HSR_ETH = CurrencyPair{CurrencyA: HSR, CurrencyB: ETH, AmountTickSize: 2, PriceTickSize: 4}
	LTC_ETH = CurrencyPair{CurrencyA: LTC, CurrencyB: ETH, AmountTickSize: 2, PriceTickSize: 4}

	UNKNOWN_PAIR = CurrencyPair{CurrencyA: UNKNOWN, CurrencyB: UNKNOWN}
)

func (pair CurrencyPair) String() string {
	return pair.ToSymbol("_")
}

func (pair CurrencyPair) Eq(c2 CurrencyPair) bool {
	return pair.String() == c2.String()
}

func (c Currency) AdaptBchToBcc() Currency {
	if c.Symbol == "BCH" || c.Symbol == "bch" {
		return BCC
	}
	return c
}

func (c Currency) AdaptBccToBch() Currency {
	if c.Symbol == "BCC" || c.Symbol == "bcc" {
		return BCH
	}
	return c
}

func NewCurrency(symbol, desc string) Currency {
	switch symbol {
	case "cny", "CNY":
		return CNY
	case "usdt", "USDT":
		return USDT
	case "usd", "USD":
		return USD
	case "usdc", "USDC":
		return USDC
	case "pax", "PAX":
		return PAX
	case "jpy", "JPY":
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
	case "okb", "OKB":
		return OKB
	case "ht", "HT":
		return HT
	case "bnb", "BNB":
		return BNB
	case "trx", "TRX":
		return TRX
	case "dot", "DOT":
		return DOT
	default:
		return Currency{strings.ToUpper(symbol), desc}
	}
}

func NewCurrencyPair(currencyA Currency, currencyB Currency) CurrencyPair {
	return CurrencyPair{CurrencyA: currencyA, CurrencyB: currencyB}
}

func NewCurrencyPair2(currencyPairSymbol string) CurrencyPair {
	return NewCurrencyPair3(currencyPairSymbol, "_")
}

func NewCurrencyPair3(currencyPairSymbol string, sep string) CurrencyPair {
	currencys := strings.Split(currencyPairSymbol, sep)
	if len(currencys) >= 2 {
		return CurrencyPair{CurrencyA: NewCurrency(currencys[0], ""),
			CurrencyB: NewCurrency(currencys[1], "")}
	}
	return UNKNOWN_PAIR
}

func (pair *CurrencyPair) SetAmountTickSize(tickSize int) CurrencyPair {
	pair.AmountTickSize = tickSize
	return *pair
}

func (pair *CurrencyPair) SetPriceTickSize(tickSize int) CurrencyPair {
	pair.PriceTickSize = tickSize
	return *pair
}

func (pair CurrencyPair) ToSymbol(joinChar string) string {
	return strings.Join([]string{pair.CurrencyA.Symbol, pair.CurrencyB.Symbol}, joinChar)
}

func (pair CurrencyPair) ToSymbol2(joinChar string) string {
	return strings.Join([]string{pair.CurrencyB.Symbol, pair.CurrencyA.Symbol}, joinChar)
}

func (pair CurrencyPair) AdaptUsdtToUsd() CurrencyPair {
	CurrencyB := pair.CurrencyB
	if pair.CurrencyB.Eq(USDT) {
		CurrencyB = USD
	}
	pair.CurrencyB = CurrencyB
	return pair
}

func (pair CurrencyPair) AdaptUsdToUsdt() CurrencyPair {
	CurrencyB := pair.CurrencyB
	if pair.CurrencyB.Eq(USD) {
		CurrencyB = USDT
	}
	pair.CurrencyB = CurrencyB
	return pair
}

//for to symbol lower , Not practical '==' operation method
func (pair CurrencyPair) ToLower() CurrencyPair {
	return CurrencyPair{CurrencyA: Currency{Symbol: strings.ToLower(pair.CurrencyA.Symbol), Desc: pair.CurrencyA.Desc},
		CurrencyB: Currency{Symbol: strings.ToLower(pair.CurrencyB.Symbol), Desc: pair.CurrencyB.Desc}}
}

func (pair CurrencyPair) Reverse() CurrencyPair {
	return CurrencyPair{CurrencyA: pair.CurrencyB, CurrencyB: pair.CurrencyA,
		AmountTickSize: pair.AmountTickSize, PriceTickSize: pair.PriceTickSize}
}
