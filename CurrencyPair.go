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
	USDT    = Currency{"USDT", "UST dollar"}
	EUR     = Currency{"EUR", ""}
	AUD     = Currency{"AUD", "AUD dollar"}
	BTC     = Currency{"BTC", "bitcoin.org"}
	BCC     = Currency{"BCC", "bitcoin-abc"}
	LTC     = Currency{"LTC", "litecoin.org"}
	ETH     = Currency{"ETH", ""}
	ETC     = Currency{"ETC", ""}
	XRP     = Currency{"XRP", ""}
	DOGE    = Currency{"DOGE", ""}
	BLK     = Currency{"BLK", ""}
	LSK     = Currency{"LSK", ""}
	GAME    = Currency{"GAME", ""}
	SC      = Currency{"SC", ""}
	BTS     = Currency{"BTS", ""}
	XPM     = Currency{"XPM", ""}
	XEM     = Currency{"XEM", ""}
	PPC     = Currency{"PPC", ""}
	VTC     = Currency{"VTC", ""}
	VRC     = Currency{"VRC", ""}
	NXT     = Currency{"NXT", ""}
	SYS     = Currency{"SYS", ""}
	DASH    = Currency{"DASH", ""}
	ANS     = Currency{"ANS", "www.antshares.org"}
	NEO     = Currency{"NEO", "neo.org"}
	ZEC     = Currency{"ZEC", ""}
	RIC     = Currency{"RIC", ""}
	ZCC     = Currency{"ZCC", ""}
	WDC     = Currency{"WDC", ""}
	EAC     = Currency{"EAC", ""}
	HLB     = Currency{"HLB", ""}
	QTUM    = Currency{"QTUM", ""}
	DGD     = Currency{"DGD", ""}
	REP     = Currency{"REP", ""}
	FIRST   = Currency{"1ST", ""}
	EOS     = Currency{"EOS", ""}
	BAT     = Currency{"BAT", ""}
	SNT     = Currency{"SNT", ""}
	PAY     = Currency{"PAY", ""}
	OMG     = Currency{"OMG", ""}
	NMC     = Currency{"NMC", ""}
	XCP     = Currency{"XCP", ""}
	XMR     = Currency{"XMR", ""}
	ARDR    = Currency{"ARDR", ""}
	GNT     = Currency{"GNT", ""}
	XLM     = Currency{"XLM", ""}
	EMC     = Currency{"EMC", ""}
	XZC     = Currency{"XZC", ""}
	BNT     = Currency{"BNT", ""}
	FTC     = Currency{"FTC", ""}
	BTM     = Currency{"BTM", ""}
	CVC     = Currency{"CVC", ""}

	//currency pair
	//cny pair
	BTC_CNY   = CurrencyPair{BTC, CNY}
	LTC_CNY   = CurrencyPair{LTC, CNY}
	BCC_CNY   = CurrencyPair{BCC, CNY}
	ETH_CNY   = CurrencyPair{ETH, CNY}
	ETC_CNY   = CurrencyPair{ETC, CNY}
	XRP_CNY   = CurrencyPair{XRP, CNY}
	DOGE_CNY  = CurrencyPair{DOGE, CNY}
	BLK_CNY   = CurrencyPair{BLK, CNY}
	LSK_CNY   = CurrencyPair{LSK, CNY}
	GAME_CNY  = CurrencyPair{GAME, CNY}
	SC_CNY    = CurrencyPair{SC, CNY}
	BTS_CNY   = CurrencyPair{BTS, CNY}
	XPM_CNY   = CurrencyPair{XPM, CNY}
	XEM_CNY   = CurrencyPair{XEM, CNY}
	PPC_CNY   = CurrencyPair{PPC, CNY}
	VTC_CNY   = CurrencyPair{VTC, CNY}
	VRC_CNY   = CurrencyPair{VRC, CNY}
	NXT_CNY   = CurrencyPair{NXT, CNY}
	SYS_CNY   = CurrencyPair{SYS, CNY}
	DASH_CNY  = CurrencyPair{DASH, CNY}
	ANS_CNY   = CurrencyPair{ANS, CNY}
	NEO_CNY   = CurrencyPair{NEO, CNY}
	ZEC_CNY   = CurrencyPair{ZEC, CNY}
	RIC_CNY   = CurrencyPair{RIC, CNY}
	ZCC_CNY   = CurrencyPair{ZCC, CNY}
	WDC_CNY   = CurrencyPair{WDC, CNY}
	EAC_CNY   = CurrencyPair{EAC, CNY}
	HLB_CNY   = CurrencyPair{HLB, CNY}
	QTUM_CNY  = CurrencyPair{QTUM, CNY}
	DGD_CNY   = CurrencyPair{DGD, CNY}
	REP_CNY   = CurrencyPair{REP, CNY}
	FIRST_CNY = CurrencyPair{FIRST, CNY}
	EOS_CNY   = CurrencyPair{EOS, CNY}
	BAT_CNY   = CurrencyPair{BAT, CNY}
	SNT_CNY   = CurrencyPair{SNT, CNY}
	PAY_CNY   = CurrencyPair{PAY, CNY}
	OMG_CNY   = CurrencyPair{OMG, CNY}
	NMC_CNY   = CurrencyPair{NMC, CNY}
	XCP_CNY   = CurrencyPair{XCP, CNY}
	XMR_CNY   = CurrencyPair{XMR, CNY}
	ARDR_CNY  = CurrencyPair{ARDR, CNY}
	GNT_CNY   = CurrencyPair{GNT, CNY}
	XLM_CNY   = CurrencyPair{XLM, CNY}
	EMC_CNY   = CurrencyPair{EMC, CNY}
	XZC_CNY   = CurrencyPair{XZC, CNY}
	BNT_CNY   = CurrencyPair{BNT, CNY}
	FTC_CNY   = CurrencyPair{FTC, CNY}
	BTM_CNY   = CurrencyPair{BTM, CNY}
	CVC_CNY   = CurrencyPair{CVC, CNY}

	//usd pair
	BTC_USD   = CurrencyPair{BTC, USD}
	LTC_USD   = CurrencyPair{LTC, USD}
	BCC_USD   = CurrencyPair{BCC, USD}
	ETH_USD   = CurrencyPair{ETH, USD}
	ETC_USD   = CurrencyPair{ETC, USD}
	XRP_USD   = CurrencyPair{XRP, USD}
	DOGE_USD  = CurrencyPair{DOGE, USD}
	BLK_USD   = CurrencyPair{BLK, USD}
	LSK_USD   = CurrencyPair{LSK, USD}
	GAME_USD  = CurrencyPair{GAME, USD}
	SC_USD    = CurrencyPair{SC, USD}
	BTS_USD   = CurrencyPair{BTS, USD}
	XPM_USD   = CurrencyPair{XPM, USD}
	XEM_USD   = CurrencyPair{XEM, USD}
	PPC_USD   = CurrencyPair{PPC, USD}
	VTC_USD   = CurrencyPair{VTC, USD}
	VRC_USD   = CurrencyPair{VRC, USD}
	NXT_USD   = CurrencyPair{NXT, USD}
	SYS_USD   = CurrencyPair{SYS, USD}
	DASH_USD  = CurrencyPair{DASH, USD}
	ANS_USD   = CurrencyPair{ANS, USD}
	NEO_USD   = CurrencyPair{NEO, USD}
	ZEC_USD   = CurrencyPair{ZEC, USD}
	RIC_USD   = CurrencyPair{RIC, USD}
	ZCC_USD   = CurrencyPair{ZCC, USD}
	WDC_USD   = CurrencyPair{WDC, USD}
	EAC_USD   = CurrencyPair{EAC, USD}
	HLB_USD   = CurrencyPair{HLB, USD}
	QTUM_USD  = CurrencyPair{QTUM, USD}
	DGD_USD   = CurrencyPair{DGD, USD}
	REP_USD   = CurrencyPair{REP, USD}
	FIRST_USD = CurrencyPair{FIRST, USD}
	EOS_USD   = CurrencyPair{EOS, USD}
	BAT_USD   = CurrencyPair{BAT, USD}
	SNT_USD   = CurrencyPair{SNT, USD}
	PAY_USD   = CurrencyPair{PAY, USD}
	OMG_USD   = CurrencyPair{OMG, USD}
	NMC_USD   = CurrencyPair{NMC, USD}
	XCP_USD   = CurrencyPair{XCP, USD}
	XMR_USD   = CurrencyPair{XMR, USD}
	ARDR_USD  = CurrencyPair{ARDR, USD}
	GNT_USD   = CurrencyPair{GNT, USD}
	XLM_USD   = CurrencyPair{XLM, USD}
	EMC_USD   = CurrencyPair{EMC, USD}
	XZC_USD   = CurrencyPair{XZC, USD}
	BNT_USD   = CurrencyPair{BNT, USD}
	FTC_USD   = CurrencyPair{FTC, USD}
	BTM_USD   = CurrencyPair{BTM, USD}

	//usdt pair
	BTC_USDT   = CurrencyPair{BTC, USDT}
	LTC_USDT   = CurrencyPair{LTC, USDT}
	BCC_USDT   = CurrencyPair{BCC, USDT}
	ETH_USDT   = CurrencyPair{ETH, USDT}
	ETC_USDT   = CurrencyPair{ETC, USDT}
	XRP_USDT   = CurrencyPair{XRP, USDT}
	DOGE_USDT  = CurrencyPair{DOGE, USDT}
	BLK_USDT   = CurrencyPair{BLK, USDT}
	LSK_USDT   = CurrencyPair{LSK, USDT}
	GAME_USDT  = CurrencyPair{GAME, USDT}
	SC_USDT    = CurrencyPair{SC, USDT}
	BTS_USDT   = CurrencyPair{BTS, USDT}
	XPM_USDT   = CurrencyPair{XPM, USDT}
	XEM_USDT   = CurrencyPair{XEM, USDT}
	PPC_USDT   = CurrencyPair{PPC, USDT}
	VTC_USDT   = CurrencyPair{VTC, USDT}
	VRC_USDT   = CurrencyPair{VRC, USDT}
	NXT_USDT   = CurrencyPair{NXT, USDT}
	SYS_USDT   = CurrencyPair{SYS, USDT}
	DASH_USDT  = CurrencyPair{DASH, USDT}
	ANS_USDT   = CurrencyPair{ANS, USDT}
	NEO_USDT   = CurrencyPair{NEO, USDT}
	ZEC_USDT   = CurrencyPair{ZEC, USDT}
	RIC_USDT   = CurrencyPair{RIC, USDT}
	ZCC_USDT   = CurrencyPair{ZCC, USDT}
	WDC_USDT   = CurrencyPair{WDC, USDT}
	EAC_USDT   = CurrencyPair{EAC, USDT}
	HLB_USDT   = CurrencyPair{HLB, USDT}
	QTUM_USDT  = CurrencyPair{QTUM, USDT}
	DGD_USDT   = CurrencyPair{DGD, USDT}
	REP_USDT   = CurrencyPair{REP, USDT}
	FIRST_USDT = CurrencyPair{FIRST, USDT}
	EOS_USDT   = CurrencyPair{EOS, USDT}
	BAT_USDT   = CurrencyPair{BAT, USDT}
	SNT_USDT   = CurrencyPair{SNT, USDT}
	PAY_USDT   = CurrencyPair{PAY, USDT}
	OMG_USDT   = CurrencyPair{OMG, USDT}
	NMC_USDT   = CurrencyPair{NMC, USDT}
	XCP_USDT   = CurrencyPair{XCP, USDT}
	XMR_USDT   = CurrencyPair{XMR, USDT}
	ARDR_USDT  = CurrencyPair{ARDR, USDT}
	GNT_USDT   = CurrencyPair{GNT, USDT}
	XLM_USDT   = CurrencyPair{XLM, USDT}
	EMC_USDT   = CurrencyPair{EMC, USDT}
	XZC_USDT   = CurrencyPair{XZC, USDT}
	BNT_USDT   = CurrencyPair{BNT, USDT}
	FTC_USDT   = CurrencyPair{FTC, USDT}
	BTM_USDT   = CurrencyPair{BTM, USDT}

	//btc pair
	BTC_BTC   = CurrencyPair{BTC, BTC}
	LTC_BTC   = CurrencyPair{LTC, BTC}
	BCC_BTC   = CurrencyPair{BCC, BTC}
	ETH_BTC   = CurrencyPair{ETH, BTC}
	ETC_BTC   = CurrencyPair{ETC, BTC}
	XRP_BTC   = CurrencyPair{XRP, BTC}
	DOGE_BTC  = CurrencyPair{DOGE, BTC}
	BLK_BTC   = CurrencyPair{BLK, BTC}
	LSK_BTC   = CurrencyPair{LSK, BTC}
	GAME_BTC  = CurrencyPair{GAME, BTC}
	SC_BTC    = CurrencyPair{SC, BTC}
	BTS_BTC   = CurrencyPair{BTS, BTC}
	XPM_BTC   = CurrencyPair{XPM, BTC}
	XEM_BTC   = CurrencyPair{XEM, BTC}
	PPC_BTC   = CurrencyPair{PPC, BTC}
	VTC_BTC   = CurrencyPair{VTC, BTC}
	VRC_BTC   = CurrencyPair{VRC, BTC}
	NXT_BTC   = CurrencyPair{NXT, BTC}
	SYS_BTC   = CurrencyPair{SYS, BTC}
	DASH_BTC  = CurrencyPair{DASH, BTC}
	ANS_BTC   = CurrencyPair{ANS, BTC}
	NEO_BTC   = CurrencyPair{NEO, BTC}
	ZEC_BTC   = CurrencyPair{ZEC, BTC}
	RIC_BTC   = CurrencyPair{RIC, BTC}
	ZCC_BTC   = CurrencyPair{ZCC, BTC}
	WDC_BTC   = CurrencyPair{WDC, BTC}
	EAC_BTC   = CurrencyPair{EAC, BTC}
	HLB_BTC   = CurrencyPair{HLB, BTC}
	QTUM_BTC  = CurrencyPair{QTUM, BTC}
	DGD_BTC   = CurrencyPair{DGD, BTC}
	REP_BTC   = CurrencyPair{REP, BTC}
	FIRST_BTC = CurrencyPair{FIRST, BTC}
	EOS_BTC   = CurrencyPair{EOS, BTC}
	BAT_BTC   = CurrencyPair{BAT, BTC}
	SNT_BTC   = CurrencyPair{SNT, BTC}
	PAY_BTC   = CurrencyPair{PAY, BTC}
	OMG_BTC   = CurrencyPair{OMG, BTC}
	NMC_BTC   = CurrencyPair{NMC, BTC}
	XCP_BTC   = CurrencyPair{XCP, BTC}
	XMR_BTC   = CurrencyPair{XMR, BTC}
	ARDR_BTC  = CurrencyPair{ARDR, BTC}
	GNT_BTC   = CurrencyPair{GNT, BTC}
	XLM_BTC   = CurrencyPair{XLM, BTC}
	EMC_BTC   = CurrencyPair{EMC, BTC}
	XZC_BTC   = CurrencyPair{XZC, BTC}
	BNT_BTC   = CurrencyPair{BNT, BTC}
	FTC_BTC   = CurrencyPair{FTC, BTC}
	BTM_BTC   = CurrencyPair{BTM, BTC}
	CVC_BTC   = CurrencyPair{CVC, BTC}
	//ETH pair
	//BTC_ETH   = CurrencyPair{BTC, ETH}

	ETC_ETH = CurrencyPair{ETC, ETH}
	EOS_ETH = CurrencyPair{EOS, ETH}

	UNKNOWN_PAIR = CurrencyPair{UNKNOWN, UNKNOWN}
)

func (c CurrencyPair) String() string {
	return c.ToSymbol("_")
}

func NewCurrency(symbol, desc string) Currency {
	return Currency{symbol, desc}
}

func NewCurrencyPair(currencyA Currency, currencyB Currency) CurrencyPair {
	return CurrencyPair{currencyA, currencyB}
}

func NewCurrencyPair2(currencyPairSymbol string) CurrencyPair {
	currencys := strings.Split(currencyPairSymbol, "_")
	if len(currencys) == 2 {
		return CurrencyPair{NewCurrency(currencys[0], ""),
			NewCurrency(currencys[1], "")}
	}
	return UNKNOWN_PAIR
}

func (pair CurrencyPair) ToSymbol(joinChar string) string {
	return strings.Join([]string{pair.CurrencyA.Symbol, pair.CurrencyB.Symbol}, joinChar)
}

func (pair CurrencyPair) ToSymbol2(joinChar string) string {
	return strings.Join([]string{pair.CurrencyB.Symbol, pair.CurrencyA.Symbol}, joinChar)
}
