package goex

import (
	"fmt"
	"strings"
)

type CurrencyPair int

func (c CurrencyPair) String() string {
	if c == 0 {
		return "nil"
	}
	return currencyPairSymbol[c-1]
}
func (c CurrencyPair) DeleteUnderLineString() string {
	if c == 0 {
		return "nil"
	}
	cur := c.String()
	s := strings.Replace(cur, "_", "", -1)
	return s
}
func (c CurrencyPair) SwitchString() string {
	if c == 0 {
		return "nil"
	}
	cur := c.String()
	s := strings.Split(cur, "_")
	ss := fmt.Sprintf("%s_%s", strings.ToUpper(s[1]), strings.ToUpper(s[0]))
	return ss
}
func (c CurrencyPair) CrossedLineString() string {
	if c == 0 {
		return "nil"
	}
	cur := c.String()
	s := strings.Split(cur, "_")
	ss := fmt.Sprintf("%s-%s", strings.ToUpper(s[1]), strings.ToUpper(s[0]))
	return ss
}

type Currency int

func (c Currency) String() string {
	if c == 0 {
		return "nil"
	}
	return currencySymbol[c-1]
}

type TradeSide int

func (ts TradeSide) String() string {
	switch ts {
	case 1:
		return "BUY"
	case 2:
		return "SELL"
	case 3:
		return "BUY_MARKET"
	case 4:
		return "SELL_MARKET"
	default:
		return "UNKNOWN"
	}
}

type TradeStatus int

func (ts TradeStatus) String() string {
	return orderStatusSymbol[ts]
}

/*currencySymbol array*/
var currencySymbol = [...]string{"cny", "usd", "btc", "ltc", "eth", "etc", "zec", "sc","eos","xpm", "xrp", "zcc", "mec", "anc", "bec", "ppc",
	"src", "tag", "bts", "wdc", "xlm", "dgc", "qrk", "doge", "ybc", "ric", "bost", "nxt", "blk", "nrs", "med", "ncs", "eac", "xcn", "sys", "xem",
	"vash", "dash", "emc", "hlb", "ardr", "xzc", "mgc", "tmc", "bns", "corg", "neos", "xst", "1cr", "bdc", "drkc", "frac", "srcc", "cc", "dao",
	"etok", "nav", "trust", "aur", "dime", "exp", "game", "ioc", "blu", "fac", "gemz", "cyc", "emo", "jlh", "xbc", "xdp", "gap", "smc", "xhc",
	"btcd", "grcx", "xusd", "mil", "lgc", "piggy", "xcp", "burst", "gns", "hiro", "huge", "lc", "fldc", "index", "leaf", "myr", "spa", "cure",
	"flo", "naut", "sjcx", "twe", "mon", "block", "cha", "giar", "hz", "ifc", "dgb", "mrc", "via", "btcs", "gold", "mmnxt", "xlb", "xmg", "balls",
	"hot", "note", "sync", "dns", "frk", "mast", "clam", "nobl", "xxc", "c2", "nxc", "q2c", "wiki", "xsv", "aero", "fzn", "mint", "qbk", "vox",
	"burn", "ltbc", "qcn", "dice", "flt", "omni", "ac", "aph", "bdg", "bitcny", "crypt", "opal", "rzr", "shibe", "sql", "sum", "bank", "con",
	"jug", "meth", "util", "vtc", "love", "mcn", "pot", "cinni", "ecc", "gdn", "grs", "key", "shopx", "xap", "steem", "yin", "amp", "nmc",
	"srg", "xdn", "yang", "xai", "ccn", "cga", "maid", "uro", "x13", "vrc", "xch", "hyp", "mrs", "plx", "qora", "usdt", "slr", "comm", "dis",
	"fvz", "ixc", "lbc", "gml", "ltcx", "nas", "axis", "cnl", "flap", "fox", "qtl", "rep", "rads", "sbd", "bcc", "cnote", "fz", "gpc", "sxc",
	"voot", "bitusd", "cai", "diem", "xsi", "ach", "cnmt", "max", "nbt", "nsr", "xmr", "emc2", "pawn", "efl", "grc", "rdd", "strat", "tac", "btm",
	"jpc", "kdc", "mts", "n5x", "prc", "unity", "bones", "fct", "silk", "gpuc", "sun", "bcn", "mzc", "uis", "geo", "lol", "dcr", "ntx", "pmc", "dvk",
	"gnt", "pand", "yc", "gue", "lcl", "bbr", "nl", "prt", "dsh", "pts", "ultc", "wc", "xcr", "noxt", "utc", "air", "bcy", "enc", "lsk", "mmxiv",
	"sdc", "soc", "tor", "ssd", "uvc", "wolf", "bbl", "glb", "mmc", "mnta", "rby", "adn", "bela", "gno", "swarm", "bits", "hvc", "itc", "usde", "aeon",
	"exe", "xc", "aby", "cach", "ebt", "min", "nxti", "fcn", "lqd", "mun", "xvc", "arch", "h2o", "drm", "str", "yacc", "fibre", "huc", "pasc", "frq",
	"pink", "xpb"}

const (
	CNY = 1 + iota
	USD
	BTC
	LTC
	ETH
	ETC
	ZEC
	SC
	EOS
	XPM
	XRP
	ZCC
	MEC
	ANC
	BEC
	PPC
	SRC
	TAG
	BTS
	WDC
	XLM
	DGC
	QRK
	DOGE
	YBC
	RIC
	BOST
	NXT
	BLK
	NRS
	MED
	NCS
	EAC
	XCN
	SYS
	XEM
	VASH
	DASH
	EMC
	HLB
	ARDR
	XZC
	MGC
	TMC
	BNS
	//	BTS
	CORG
	NEOS
	XST
	OneCR
	BDC
	DRKC
	FRAC
	SRCC
	CC
	DAO
	eTOK
	NAV
	TRUST
	AUR
	DIME
	EXP
	GAME
	IOC
	BLU
	FAC
	GEMZ
	CYC
	EMO
	JLH
	XBC
	XDP
	//	DASH
	GAP
	SMC
	XHC
	BTCD
	GRCX
	XUSD
	MIL
	LGC
	PIGGY
	XCP
	BURST
	GNS
	HIRO
	HUGE
	LC
	FLDC
	INDEX
	LEAF
	MYR
	SPA
	CURE
	FLO
	NAUT
	SJCX
	TWE
	MON
	BLOCK
	CHA
	GIAR
	HZ
	IFC
	DGB
	MRC
	VIA
	BTCS
	GOLD
	MMNXT
	XLB
	XMG
	BALLS
	HOT
	NOTE
	SYNC
	//	ARDR
	DNS
	//	ETC
	FRK
	MAST
	CLAM
	NOBL
	XXC
	C2
	NXC
	Q2C
	WIKI
	XSV
	AERO
	FZN
	MINT
	QBK
	VOX
	BURN
	LTBC
	QCN
	//	XEM
	DICE
	FLT
	OMNI
	AC
	APH
	BDG
	BITCNY
	CRYPT
	//	NXT
	OPAL
	RZR
	SHIBE
	SQL
	SUM
	BANK
	CON
	JUG
	METH
	//	SC
	UTIL
	VTC
	LOVE
	MCN
	POT
	CINNI
	ECC
	GDN
	GRS
	KEY
	SHOPX
	XAP
	STEEM
	YIN
	AMP
	NMC
	SRG
	XDN
	YANG
	XAI
	CCN
	CGA
	MAID
	URO
	X13
	VRC
	XCH
	HYP
	MRS
	PLX
	QORA
	USDT
	SLR
	COMM
	DIS
	FVZ
	IXC
	LBC
	GML
	LTCX
	NAS
	AXIS
	CNL
	//	ETH
	FLAP
	FOX
	QTL
	REP
	RADS
	//	RIC
	SBD
	BCC
	CNOTE
	FZ
	GPC
	//	MEC
	SXC
	VOOT
	BITUSD
	CAI
	DIEM
	XSI
	ACH
	CNMT
	MAX
	NBT
	NSR
	XMR
	EMC2
	PAWN
	//	SYS
	//	BOST
	EFL
	GRC
	RDD
	STRAT
	TAC
	BTM
	JPC
	KDC
	MTS
	N5X
	//	BTC
	PRC
	UNITY
	BONES
	//	EAC
	FCT
	SILK
	GPUC
	SUN
	//	XCN
	BCN
	MZC
	UIS
	//	XRP
	GEO
	LOL
	DCR
	NTX
	//	ZEC
	PMC
	DVK
	GNT
	//	LTC
	PAND
	YC
	GUE
	LCL
	BBR
	NL
	PRT
	//	XPM
	DSH
	PTS
	ULTC
	WC
	XCR
	NOXT
	UTC
	AIR
	BCY
	ENC
	LSK
	MMXIV
	SDC
	SOC
	TOR
	SSD
	UVC
	WOLF
	BBL
	GLB
	MMC
	MNTA
	RBY
	ADN
	BELA
	//	DOGE
	GNO
	SWARM
	BITS
	HVC
	ITC
	USDE
	AEON
	EXE
	XC
	ABY
	CACH
	EBT
	MIN
	NXTI
	FCN
	LQD
	MUN
	//	WDC
	XVC
	ARCH
	H2O
	DRM
	STR
	YACC
	//	BLK
	FIBRE
	HUC
	//	NRS
	PASC
	FRQ
	PINK
	//	PPC
	XPB
)

func SymbolCurrency(s string) int {
	for index, v := range currencySymbol {
		if 0 == strings.Compare(strings.ToLower(s), v) {
			return index + 1
		}
	}
	return -1
}
func SymbolPairCurrency(s string) int {
	for index, v := range currencyPairSymbol {
		if 0 == strings.Compare(s, v) {
			return index + 1
		}
	}
	return -1
}

var currencyPairSymbol = [...]string{"btc_cny", "btc_usdt", "btc_usd", "btc_jpy", "fx_btc_jpy", "ltc_cny", "ltc_usdt", "ltc_usd", "ltc_btc",
	"eth_cny", "eth_usdt", "eth_usd", "eth_btc", "etc_cny", "etc_usdt", "etc_usd", "etc_btc", "etc_eth", "zec_cny", "zec_usd", "zec_btc",
	"rep_cny", "rep_eth", "rep_btc", "xrp_cny", "xrp_usdt", "xrp_usd", "xrp_btc", "doge_cny", "doge_usd",
	"doge_btc", "blk_cny", "blk_usd", "blk_btc", "lsk_cny", "lsk_usd", "lsk_btc", "game_cny", "game_usd", "game_btc", "sc_cny", "sc_usd",
	"sc_btc", "gnt_btc", "gnt_cny", "bts_cny", "bts_usd", "bts_btc", "hlb_cny", "hlb_usd", "hlb_btc", "xpm_cny", "xpm_usd", "xpm_btc", "ric_cny", "ric_usd",
	"ric_btc", "xem_cny", "xem_usd", "xem_btc", "eac_cny", "eac_usd", "eac_btc", "ppc_cny", "ppc_usd", "ppc_btc", "plc_cny", "plc_usd", "plc_btc",
	"vtc_cny", "vtc_usd", "vtc_btc", "vrc_cny", "vrc_usd", "vrc_btc", "nxt_cny", "nxt_usd", "nxt_btc", "zcc_cny", "zcc_usd", "zcc_btc",
	"wdc_cny", "wdc_usd", "wdc_btc", "sys_cny", "sys_usd", "sys_btc", "dash_cny", "dash_usd", "dash_btc", "dsh_usd", "dash_usdt", "ybc_cny", "ybc_usd", "ybc_btc", "xcn_btc",
	"ans_cny", "ans_usd", "ans_btc", "xmr_cny", "xmr_usd", "xmr_btc", "iota_cny", "iota_usd", "iota_btc", "iota_eth", "eos_cny", "eos_usd", "eos_btc", "eos_eth", "eos_usdt",
	"bcc_cny", "bcc_usd", "bcc_btc", "bcu_cny", "bcu_usd", "bcu_btc"}

const (
	BTC_CNY = 1 + iota
	BTC_USDT
	BTC_USD
	BTC_JPY
	FX_BTC_JPY

	LTC_CNY
	LTC_USDT
	LTC_USD
	LTC_BTC

	ETH_CNY
	ETH_USDT
	ETH_USD
	ETH_BTC

	ETC_CNY
	ETC_USDT
	ETC_USD
	ETC_BTC
	ETC_ETH

	ZEC_CNY
	ZEC_USD
	ZEC_BTC
	REP_CNY
	REP_ETH
	REP_BTC

	XRP_CNY
	XRP_USDT
	XRP_USD
	XRP_BTC

	DOGE_CNY
	DOGE_USD
	DOGE_BTC

	BLK_CNY
	BLK_USD
	BLK_BTC

	LSK_CNY
	LSK_USD
	LSK_BTC

	GAME_CNY
	GAME_USD
	GAME_BTC

	SC_CNY
	SC_USD
	SC_BTC
	GNT_BTC
	GNT_CNY

	BTS_CNY
	BTS_USD
	BTS_BTC

	HLB_CNY
	HLB_USD
	HLB_BTC

	XPM_CNY
	XPM_USD
	XPM_BTC

	RIC_CNY
	RIC_USD
	RIC_BTC

	XEM_CNY
	XEM_USD
	XEM_BTC

	EAC_CNY
	EAC_USD
	EAC_BTC

	PPC_CNY
	PPC_USD
	PPC_BTC

	PLC_CNY
	PLC_USD
	PLC_BTC

	VTC_CNY
	VTC_USD
	VTC_BTC

	VRC_CNY
	VRC_USD
	VRC_BTC

	NXT_CNY
	NXT_USD
	NXT_BTC

	ZCC_CNY
	ZCC_USD
	ZCC_BTC

	WDC_CNY
	WDC_USD
	WDC_BTC

	SYS_CNY
	SYS_USD
	SYS_BTC

	DASH_CNY
	DASH_USD
	DASH_BTC
	DSH_USD
	DASH_USDT

	YBC_CNY
	YBC_USD
	YBC_BTC

	XCN_BTC

	ANS_CNY
	ANS_USD
	ANS_BTC

	XMR_CNY
	XMR_USD
	XMR_BTC

	IOTA_CNY
	IOTA_USD
	IOTA_BTC
	IOTA_ETH

	EOS_CNY
	EOS_USD
	EOS_BTC
	EOS_ETH
	EOS_USDT

	BCC_CNY
	BCC_USD
	BCC_BTC

	BCU_CNY
	BCU_USD
	BCU_BTC
)

const (
	BUY = 1 + iota
	SELL
	BUY_MARKET
	SELL_MARKET
)

var orderStatusSymbol = [...]string{"UNFINISH", "PART_FINISH", "FINISH", "CANCEL", "REJECT", "CANCEL_ING"}

const (
	ORDER_UNFINISH = iota
	ORDER_PART_FINISH
	ORDER_FINISH
	ORDER_CANCEL
	ORDER_REJECT
	ORDER_CANCEL_ING
)

const (
	OPEN_BUY   = 1 + iota //开多
	OPEN_SELL             //开空
	CLOSE_BUY             //平多
	CLOSE_SELL            //平空
)

var CurrencyPairSymbol = map[CurrencyPair]string{
	BTC_CNY: "btc_cny",
	BTC_USD: "btc_usd",
	LTC_CNY: "ltc_cny",
	LTC_USD: "ltc_usd",
	ETH_CNY: "eth_cny",
	ETH_USD: "eth_usd",
	ETH_BTC: "eth_btc",
	ETC_CNY: "etc_cny",
	ETC_USD: "etc_usd",
	ETC_BTC: "etc_btc",
	BTS_CNY: "bts_cny",
	SC_CNY:  "sc_cny",
	EOS_CNY: "eos_cny"}

var (
	THIS_WEEK_CONTRACT = "this_week" //周合约
	NEXT_WEEK_CONTRACT = "next_week" //次周合约
	QUARTER_CONTRACT   = "quarter"   //季度合约
)
