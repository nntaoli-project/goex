package btc38

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	. "github.com/nntaoli/crypto_coin_api"
)

//因服务器有防CC攻击策略，每60秒内调用次数不可超过120次，超过部分将被防火墙拦截。
const (
	EXCHANGE_NAME = "btc38.com"

	API_BASE_URL = "http://api.btc38.com/"
	API_V1       = API_BASE_URL + "v1/"

	//	BASE_URL    = "http://api.btc38.com/v1/"
	TICKER_URI  = "ticker.php?c=%s&mk_type=%s"
	DEPTH_URI   = "depth.php?c=%s&mk_type=%s"
	ACCOUNT_URI = "getMyBalance.php"

	TRADE_URI   = "submitOrder.php"
	CANCEL_URI  = "cancelOrder.php"
	ORDERS_INFO = "getOrderList.php"
)

type Btc38 struct {
	accessKey,
	secretKey,
	accountId string
	httpClient *http.Client
}

func New(client *http.Client, accessKey, secretKey, accountId string) *Btc38 {
	return &Btc38{accessKey, secretKey, accountId, client}
}

func (btc38 *Btc38) GetExchangeName() string {
	return EXCHANGE_NAME
}

func (btc38 *Btc38) buildPostForm(postForm *url.Values) error {
	postForm.Set("created", fmt.Sprintf("%d", time.Now().Unix()))
	postForm.Set("access_key", btc38.accessKey)
	postForm.Set("secret_key", btc38.secretKey)
	sign, err := GetParamMD5Sign(btc38.secretKey, postForm.Encode())
	if err != nil {
		return err
	}
	postForm.Set("sign", sign)
	postForm.Del("secret_key")
	return nil
}

func convertCurrencyPair(currencyPair CurrencyPair) (string, string) {
	switch currencyPair {
	case BTC_CNY:
		return "btc", "cny"
	case LTC_CNY:
		return "ltc", "cny"
	case LTC_BTC:
		return "ltc", "btc"
	case ETH_CNY:
		return "eth", "cny"
	case ETH_BTC:
		return "eth", "btc"
	case ETC_CNY:
		return "etc", "cny"
	case ETC_BTC:
		return "etc", "btc"
	case XRP_CNY:
		return "xrp", "cny"
	case XRP_BTC:
		return "xrp", "btc"
	case DOGE_CNY:
		return "doge", "cny"
	case DOGE_BTC:
		return "doge", "btc"
	case BLK_CNY:
		return "blk", "cny"
	case BLK_BTC:
		return "blk", "btc"
	case LSK_CNY:
		return "lsk", "cny"
	case LSK_BTC:
		return "lsk", "btc"
	case GAME_CNY:
		return "game", "cny"
	case GAME_BTC:
		return "game", "btc"
	case SC_CNY:
		return "sc", "cny"
	case SC_BTC:
		return "sc", "btc"
	case BTS_CNY:
		return "bts", "cny"
	case BTS_BTC:
		return "bts", "btc"
	case HLB_CNY:
		return "hlb", "cny"
	case HLB_BTC:
		return "hlb", "btc"
	case XPM_CNY:
		return "xpm", "cny"
	case XPM_BTC:
		return "xpm", "btc"
	case RIC_CNY:
		return "ric", "cny"
	case RIC_BTC:
		return "ric", "btc"
	case XEM_CNY:
		return "xem", "cny"
	case XEM_BTC:
		return "xem", "btc"
	case EAC_CNY:
		return "ea", "cny"
	case EAC_BTC:
		return "eac", "btc"
	case PPC_CNY:
		return "ppc", "cny"
	case PPC_BTC:
		return "ppc", "btc"
	case VTC_CNY:
		return "vtc", "cny"
	case VTC_BTC:
		return "vtc", "btc"
	case VRC_CNY:
		return "vrc", "cny"
	case VRC_BTC:
		return "vrc", "btc"
	case NXT_CNY:
		return "nxt", "cny"
	case NXT_BTC:
		return "nxt", "btc"
	case ZCC_CNY:
		return "zcc", "cny"
	case ZCC_BTC:
		return "zcc", "btc"
	case WDC_CNY:
		return "wdc", "cny"
	case WDC_BTC:
		return "wdc", "btc"
	case SYS_CNY:
		return "sys", "cny"
	case SYS_BTC:
		return "sys", "btc"
	case DASH_CNY:
		return "dash", "cny"
	case DASH_BTC:
		return "dash", "btc"
	default:
		return "err", "err"
	}
}

func (btc38 *Btc38) GetTicker(currency CurrencyPair) (*Ticker, error) {

	cur, money := convertCurrencyPair(currency)
	if cur == "err" {
		log.Println("Unsupport The CurrencyPair")
		return nil, errors.New("Unsupport The CurrencyPair")
	}
	tickerUri := API_V1 + fmt.Sprintf(TICKER_URI, cur, money)
	fmt.Println("tickerUri:", tickerUri)
	timestamp := time.Now().Unix()
	bodyDataMap, err := HttpGet(btc38.httpClient, tickerUri)
	fmt.Println(err)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Println("bodyDataMap:", bodyDataMap)
	var tickerMap map[string]interface{}
	var ticker Ticker

	switch bodyDataMap["ticker"].(type) {
	case map[string]interface{}:
		tickerMap = bodyDataMap["ticker"].(map[string]interface{})
		fmt.Println("tickerMap:", tickerMap)

	default:
		return nil, errors.New(fmt.Sprintf("Type Convert Error ? \n %s", bodyDataMap))
	}

	ticker.Date = uint64(timestamp)
	ticker.Last = tickerMap["last"].(float64)
	ticker.Buy = tickerMap["buy"].(float64)
	ticker.Sell = tickerMap["sell"].(float64)
	ticker.Low = tickerMap["low"].(float64)
	ticker.High = tickerMap["high"].(float64)
	ticker.Vol = tickerMap["vol"].(float64)

	return &ticker, nil
}

func (btc38 *Btc38) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	var depthUri string
	cur, money := convertCurrencyPair(currency)
	if cur == "err" {
		log.Println("Unsupport The CurrencyPair")
		return nil, errors.New("Unsupport The CurrencyPair")
	}
	depthUri = fmt.Sprintf(API_V1+DEPTH_URI, cur, money)

	bodyDataMap, err := HttpGet(btc38.httpClient, depthUri)

	if err != nil {
		return nil, err
	}

	if bodyDataMap["code"] != nil {
		log.Println(bodyDataMap)
		return nil, errors.New(fmt.Sprintf("%s", bodyDataMap))
	}

	var depth Depth

	asks, isOK := bodyDataMap["asks"].([]interface{})
	if !isOK {
		return nil, errors.New("asks assert error")
	}

	i := len(asks) - 1

	for ; i >= 0; i-- {
		ask := asks[i]
		var dr DepthRecord
		for i, vv := range ask.([]interface{}) {
			switch i {
			case 0:
				dr.Price = vv.(float64)
			case 1:
				dr.Amount = vv.(float64)
			}
		}
		depth.AskList = append(depth.AskList, dr)
	}

	for _, v := range bodyDataMap["bids"].([]interface{}) {
		var dr DepthRecord
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price = vv.(float64)
			case 1:
				dr.Amount = vv.(float64)
			}
		}
		depth.BidList = append(depth.BidList, dr)
	}

	return &depth, nil
}

func (btc38 *Btc38) GetAccount() (*Account, error) {
	postData := url.Values{}
	postData.Set("key", btc38.accessKey)
	timeNow := fmt.Sprintf("%d", time.Now().Unix())
	postData.Set("time", timeNow)

	mdt := fmt.Sprintf("%s_%s_%s_%s", btc38.accessKey, btc38.accountId, btc38.secretKey, timeNow)
	sign, _ := GetParamMD5Sign(btc38.secretKey, mdt)
	postData.Set("md5", sign)
	fmt.Println("postData:", postData)
	accountUri := fmt.Sprintf(API_V1 + ACCOUNT_URI)
	fmt.Println(accountUri)
	bodyData, err := HttpPostForm(btc38.httpClient, accountUri, postData)
	if err != nil {
		fmt.Println("err:", err)
		return nil, err
	}
	bodyDataS := string(bodyData)
	fmt.Println("bodyDataS:", bodyDataS)
	var bodyDataMap map[string]interface{}

	err = json.Unmarshal(bodyData, &bodyDataMap)
	if err != nil {
		println(string(bodyData))
		fmt.Println("err:", err)
		return nil, err
	}
	fmt.Println("bodyDataMap:", bodyDataMap)
	if bodyDataMap["code"] != nil {
		return nil, errors.New(fmt.Sprintf("%s", bodyDataMap))
	}

	account := new(Account)
	account.Exchange = btc38.GetExchangeName()

	account.SubAccounts = make(map[Currency]SubAccount, 50)

	var btcSubAccount SubAccount
	btcSubAccount.Currency = BTC
	btcSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["btc_balance"].(string), 64)
	btcSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["btc_balance_lock"].(string), 64)
	account.SubAccounts[btcSubAccount.Currency] = btcSubAccount

	var ltcSubAccount SubAccount
	ltcSubAccount.Currency = LTC
	ltcSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["ltc_balance"].(string), 64)
	ltcSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["ltc_balance_lock"].(string), 64)
	account.SubAccounts[ltcSubAccount.Currency] = ltcSubAccount

	var cnySubAccount SubAccount
	cnySubAccount.Currency = CNY
	cnySubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["cny_balance"].(string), 64)
	cnySubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["cny_balance_lock"].(string), 64)
	account.SubAccounts[cnySubAccount.Currency] = cnySubAccount

	var xpmSubAccount SubAccount
	xpmSubAccount.Currency = XPM
	xpmSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["xpm_balance"].(string), 64)
	xpmSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["xpm_balance_lock"].(string), 64)
	account.SubAccounts[xpmSubAccount.Currency] = xpmSubAccount

	var xrpSubAccount SubAccount
	xrpSubAccount.Currency = XRP
	xrpSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["xrp_balance"].(string), 64)
	xrpSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["xrp_balance_lock"].(string), 64)
	account.SubAccounts[xrpSubAccount.Currency] = xrpSubAccount

	var zccSubAccount SubAccount
	zccSubAccount.Currency = ZCC
	zccSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["zcc_balance"].(string), 64)
	zccSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["zcc_balance_lock"].(string), 64)
	account.SubAccounts[zccSubAccount.Currency] = zccSubAccount

	var mecSubAccount SubAccount
	mecSubAccount.Currency = MEC
	mecSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["mec_balance"].(string), 64)
	mecSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["mec_balance_lock"].(string), 64)
	account.SubAccounts[mecSubAccount.Currency] = mecSubAccount

	var ancSubAccount SubAccount
	ancSubAccount.Currency = ANC
	ancSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["anc_balance"].(string), 64)
	ancSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["anc_balance_lock"].(string), 64)
	account.SubAccounts[ancSubAccount.Currency] = ancSubAccount

	var becSubAccount SubAccount
	becSubAccount.Currency = BEC
	becSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["bec_balance"].(string), 64)
	becSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["bec_balance_lock"].(string), 64)
	account.SubAccounts[becSubAccount.Currency] = becSubAccount

	var ppcSubAccount SubAccount
	ppcSubAccount.Currency = PPC
	ppcSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["ppc_balance"].(string), 64)
	ppcSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["ppc_balance_lock"].(string), 64)
	account.SubAccounts[ppcSubAccount.Currency] = ppcSubAccount

	var srcSubAccount SubAccount
	srcSubAccount.Currency = SRC
	srcSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["src_balance"].(string), 64)
	srcSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["src_balance_lock"].(string), 64)
	account.SubAccounts[srcSubAccount.Currency] = srcSubAccount

	var tagSubAccount SubAccount
	tagSubAccount.Currency = TAG
	tagSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["tag_balance"].(string), 64)
	tagSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["tag_balance_lock"].(string), 64)
	account.SubAccounts[tagSubAccount.Currency] = tagSubAccount

	var btsSubAccount SubAccount
	btsSubAccount.Currency = BTS
	btsSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["bts_balance"].(string), 64)
	btsSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["bts_balance_lock"].(string), 64)
	account.SubAccounts[btsSubAccount.Currency] = btsSubAccount

	var wdcSubAccount SubAccount
	wdcSubAccount.Currency = WDC
	wdcSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["wdc_balance"].(string), 64)
	wdcSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["wdc_balance_lock"].(string), 64)
	account.SubAccounts[wdcSubAccount.Currency] = wdcSubAccount

	var xlmSubAccount SubAccount
	xlmSubAccount.Currency = XLM
	xlmSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["xlm_balance"].(string), 64)
	xlmSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["xlm_balance_lock"].(string), 64)
	account.SubAccounts[xlmSubAccount.Currency] = xlmSubAccount

	var dgcSubAccount SubAccount
	dgcSubAccount.Currency = DGC
	dgcSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["dgc_balance"].(string), 64)
	dgcSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["dgc_balance_lock"].(string), 64)
	account.SubAccounts[dgcSubAccount.Currency] = dgcSubAccount

	var qrkSubAccount SubAccount
	qrkSubAccount.Currency = QRK
	qrkSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["qrk_balance"].(string), 64)
	qrkSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["qrk_balance_lock"].(string), 64)
	account.SubAccounts[qrkSubAccount.Currency] = qrkSubAccount

	var dogeSubAccount SubAccount
	dogeSubAccount.Currency = DOGE
	dogeSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["doge_balance"].(string), 64)
	dogeSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["doge_balance_lock"].(string), 64)
	account.SubAccounts[dogeSubAccount.Currency] = dogeSubAccount

	var ybcSubAccount SubAccount
	ybcSubAccount.Currency = YBC
	ybcSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["ybc_balance"].(string), 64)
	ybcSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["ybc_balance_lock"].(string), 64)
	account.SubAccounts[ybcSubAccount.Currency] = ybcSubAccount

	var ricSubAccount SubAccount
	ricSubAccount.Currency = RIC
	ricSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["ric_balance"].(string), 64)
	ricSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["ric_balance_lock"].(string), 64)
	account.SubAccounts[ricSubAccount.Currency] = ricSubAccount

	var bostSubAccount SubAccount
	bostSubAccount.Currency = BOST
	bostSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["bost_balance"].(string), 64)
	bostSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["bost_balance_lock"].(string), 64)
	account.SubAccounts[bostSubAccount.Currency] = bostSubAccount

	var nxtSubAccount SubAccount
	nxtSubAccount.Currency = NXT
	nxtSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["nxt_balance"].(string), 64)
	nxtSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["nxt_balance_lock"].(string), 64)
	account.SubAccounts[nxtSubAccount.Currency] = nxtSubAccount

	var blkSubAccount SubAccount
	blkSubAccount.Currency = BLK
	blkSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["blk_balance"].(string), 64)
	blkSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["blk_balance_lock"].(string), 64)
	account.SubAccounts[blkSubAccount.Currency] = blkSubAccount

	var nrsSubAccount SubAccount
	nrsSubAccount.Currency = NRS
	nrsSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["nrs_balance"].(string), 64)
	nrsSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["nrs_balance_lock"].(string), 64)
	account.SubAccounts[nrsSubAccount.Currency] = nrsSubAccount

	var medSubAccount SubAccount
	medSubAccount.Currency = MED
	medSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["med_balance"].(string), 64)
	medSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["med_balance_lock"].(string), 64)
	account.SubAccounts[medSubAccount.Currency] = medSubAccount

	var ncsSubAccount SubAccount
	ncsSubAccount.Currency = NCS
	ncsSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["ncs_balance"].(string), 64)
	ncsSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["ncs_balance_lock"].(string), 64)
	account.SubAccounts[ncsSubAccount.Currency] = ncsSubAccount

	var eacSubAccount SubAccount
	eacSubAccount.Currency = EAC
	eacSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["eac_balance"].(string), 64)
	eacSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["eac_balance_lock"].(string), 64)
	account.SubAccounts[eacSubAccount.Currency] = eacSubAccount

	var xcnSubAccount SubAccount
	xcnSubAccount.Currency = XCN
	xcnSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["xcn_balance"].(string), 64)
	xcnSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["xcn_balance_lock"].(string), 64)
	account.SubAccounts[xcnSubAccount.Currency] = xcnSubAccount

	var sysSubAccount SubAccount
	sysSubAccount.Currency = SYS
	sysSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["sys_balance"].(string), 64)
	sysSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["sys_balance_lock"].(string), 64)
	account.SubAccounts[sysSubAccount.Currency] = sysSubAccount

	var xemSubAccount SubAccount
	xemSubAccount.Currency = XEM
	xemSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["xem_balance"].(string), 64)
	xemSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["xem_balance_lock"].(string), 64)
	account.SubAccounts[xemSubAccount.Currency] = xemSubAccount

	var vashSubAccount SubAccount
	vashSubAccount.Currency = VASH
	vashSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["vash_balance"].(string), 64)
	vashSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["vash_balance_lock"].(string), 64)
	account.SubAccounts[vashSubAccount.Currency] = vashSubAccount

	var dashSubAccount SubAccount
	dashSubAccount.Currency = DASH
	dashSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["dash_balance"].(string), 64)
	dashSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["dash_balance_lock"].(string), 64)
	account.SubAccounts[dashSubAccount.Currency] = dashSubAccount

	var emcSubAccount SubAccount
	emcSubAccount.Currency = EMC
	emcSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["emc_balance"].(string), 64)
	emcSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["emc_balance_lock"].(string), 64)
	account.SubAccounts[emcSubAccount.Currency] = emcSubAccount

	var hlbSubAccount SubAccount
	hlbSubAccount.Currency = HLB
	hlbSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["hlb_balance"].(string), 64)
	hlbSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["hlb_balance_lock"].(string), 64)
	account.SubAccounts[hlbSubAccount.Currency] = hlbSubAccount

	var ardrSubAccount SubAccount
	ardrSubAccount.Currency = ARDR
	ardrSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["ardr_balance"].(string), 64)
	ardrSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["ardr_balance_lock"].(string), 64)
	account.SubAccounts[ardrSubAccount.Currency] = ardrSubAccount

	var xzcSubAccount SubAccount
	xzcSubAccount.Currency = XZC
	xzcSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["xzc_balance"].(string), 64)
	xzcSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["xzc_balance_lock"].(string), 64)
	account.SubAccounts[xzcSubAccount.Currency] = xzcSubAccount

	var mgcSubAccount SubAccount
	mgcSubAccount.Currency = MGC
	mgcSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["mgc_balance"].(string), 64)
	mgcSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["mgc_balance_lock"].(string), 64)
	account.SubAccounts[mgcSubAccount.Currency] = mgcSubAccount

	var tmcSubAccount SubAccount
	tmcSubAccount.Currency = TMC
	tmcSubAccount.Amount, _ = strconv.ParseFloat(bodyDataMap["tmc_balance"].(string), 64)
	tmcSubAccount.ForzenAmount, _ = strconv.ParseFloat(bodyDataMap["tmc_balance_lock"].(string), 64)
	account.SubAccounts[tmcSubAccount.Currency] = tmcSubAccount
	return account, nil
}

//func (btc38 *Btc38) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
//	postData := url.Values{}
//	postData.Set("method", "order_info")
//	postData.Set("id", orderId)

//	switch currency {
//	case BTC_CNY:
//		postData.Set("coin_type", "1")
//	case LTC_CNY:
//		postData.Set("coin_type", "2")
//	}

//	btc38.buildPostForm(&postData)

//	bodyData, err := HttpPostForm(btc38.httpClient, TRADE_API_V3, postData)
//	if err != nil {
//		log.Println(err)
//		return nil, err
//	}

//	var bodyDataMap map[string]interface{}
//	err = json.Unmarshal(bodyData, &bodyDataMap)
//	if err != nil {
//		println(string(bodyData))
//		return nil, err
//	}

//	if bodyDataMap["code"] != nil {
//		return nil, errors.New(string(bodyData))
//	}

//	//fmt.Println(bodyDataMap);
//	order := new(Order)
//	order.Currency = currency
//	order.OrderID, _ = strconv.Atoi(orderId)
//	order.Side = TradeSide(bodyDataMap["type"].(float64))
//	order.Amount, _ = strconv.ParseFloat(bodyDataMap["order_amount"].(string), 64)
//	order.DealAmount, _ = strconv.ParseFloat(bodyDataMap["processed_amount"].(string), 64)
//	order.Price, _ = strconv.ParseFloat(bodyDataMap["order_price"].(string), 64)
//	order.AvgPrice, _ = strconv.ParseFloat(bodyDataMap["processed_price"].(string), 64)
//	order.Fee, _ = strconv.ParseFloat(bodyDataMap["fee"].(string), 64)

//	tradeStatus := TradeStatus(bodyDataMap["status"].(float64))
//	switch tradeStatus {
//	case 0:
//		order.Status = ORDER_UNFINISH
//	case 1:
//		order.Status = ORDER_PART_FINISH
//	case 2:
//		order.Status = ORDER_FINISH
//	case 3:
//		order.Status = ORDER_CANCEL
//	}
//	//fmt.Println(order)
//	return order, nil
//}

func (btc38 *Btc38) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	_, money := convertCurrencyPair(currency)
	postData := url.Values{}
	postData.Set("key", btc38.accessKey)
	timeNow := fmt.Sprintf("%d", time.Now().Unix())
	postData.Set("time", timeNow)

	mdt := fmt.Sprintf("%s_%s_%s_%s", btc38.accessKey, btc38.accountId, btc38.secretKey, timeNow)
	sign, _ := GetParamMD5Sign(btc38.secretKey, mdt)
	postData.Set("md5", sign)

	postData.Set("mk_type", strings.ToUpper(money))
	//	postData.Set("coinname", money)
	//	postData.Set("mk_type", "CNY")
	postData.Set("coinname", "")

	orderList := fmt.Sprintf(API_V1 + ORDERS_INFO)
	bodyData, err := HttpPostForm(btc38.httpClient, orderList, postData)
	if err != nil {
		fmt.Println("err:", err)
		return nil, err
	}

	if strings.Contains(string(bodyData), "code") {
		return nil, errors.New(string(bodyData))
	}
	fmt.Println("bodyData:", string(bodyData))

	var bodyDataMap []map[string]interface{}
	err = json.Unmarshal(bodyData, &bodyDataMap)
	if err != nil {
		return nil, err
	}

	var orders []Order

	for _, v := range bodyDataMap {
		order := Order{}
		cur := fmt.Sprintf("%s_cny", v["coinname"].(string))
		//		fmt.Println("cur:", cur)
		//		fmt.Println("SymbolPairCurrency(cur):", SymbolPairCurrency(cur))
		order.Currency = (CurrencyPair)(SymbolPairCurrency(cur))
		order.Amount, _ = strconv.ParseFloat(v["amount"].(string), 64)
		order.Price, _ = strconv.ParseFloat(v["price"].(string), 64)
		t1 := v["time"].(string)
		t2, _ := time.Parse("2006-01-02 15:04:05", t1)
		order.OrderTime = (int)(t2.Unix())
		id, _ := strconv.ParseInt(v["id"].(string), 10, 64)
		order.OrderID = (int)(id)

		types, _ := strconv.ParseInt(v["type"].(string), 10, 64)
		order.Side = TradeSide(types)
		orders = append(orders, order)
		fmt.Println("order:", order)
	}

	return orders, nil
}

//func (btc38 *Btc38) placeOrder(method, amount, price string, currency CurrencyPair) (*Order, error) {
//	postData := url.Values{}
//	postData.Set("method", method)

//	switch method {
//	case "buy", "sell":
//		postData.Set("amount", amount)
//		postData.Set("price", price)
//	case "buy_market":
//		postData.Set("amount", price)
//	case "sell_market":
//		postData.Set("amount", amount)
//	}

//	switch currency {
//	case BTC_CNY:
//		postData.Set("coin_type", "1")
//	case LTC_CNY:
//		postData.Set("coin_type", "2")
//	}

//	btc38.buildPostForm(&postData)

//	bodyData, err := HttpPostForm(btc38.httpClient, TRADE_API_V3, postData)
//	if err != nil {
//		return nil, err
//	}

//	//{"result":"success","id":1321475746}
//	//println(string(bodyData))

//	var bodyDataMap map[string]interface{}
//	err = json.Unmarshal(bodyData, &bodyDataMap)

//	if err != nil {
//		return nil, err
//	}

//	if bodyDataMap["code"] != nil {
//		return nil, errors.New(string(bodyData))
//	}

//	ret := bodyDataMap["result"].(string)

//	if strings.Compare(ret, "success") == 0 {
//		order := new(Order)
//		order.OrderID = int(bodyDataMap["id"].(float64))
//		order.Price, _ = strconv.ParseFloat(price, 64)
//		order.Amount, _ = strconv.ParseFloat(amount, 64)
//		order.Currency = currency
//		order.Status = ORDER_UNFINISH
//		return order, nil
//	}

//	return nil, errors.New(fmt.Sprintf("Place Limit %s Order Fail.", method))
//}

//func (btc38 *Btc38) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
//	order, err := btc38.placeOrder("buy", amount, price, currency)

//	if err != nil {
//		return nil, err
//	}

//	order.Side = BUY

//	return order, nil
//}

//func (btc38 *Btc38) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
//	order, err := btc38.placeOrder("sell", amount, price, currency)

//	if err != nil {
//		return nil, err
//	}

//	order.Side = SELL

//	return order, nil
//}

//func (btc38 *Btc38) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
//	order, err := btc38.placeOrder("buy_market", amount, price, currency)
//	if err != nil {
//		log.Println(err)
//		return nil, err
//	}
//	order.Side = BUY
//	return order, nil
//}

//func (btc38 *Btc38) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
//	order, err := btc38.placeOrder("sell_market", amount, price, currency)
//	if err != nil {
//		log.Println(err)
//		return nil, err
//	}
//	order.Side = SELL
//	return order, nil
//}

//func (btc38 *Btc38) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
//	//1321490762
//	postData := url.Values{}
//	postData.Set("method", "cancel_order")
//	postData.Set("id", orderId)

//	switch currency {
//	case BTC_CNY:
//		postData.Set("coin_type", "1")
//	case LTC_CNY:
//		postData.Set("coin_type", "2")
//	}

//	btc38.buildPostForm(&postData)

//	bodyData, err := HttpPostForm(btc38.httpClient, TRADE_API_V3, postData)
//	if err != nil {
//		return false, err
//	}

//	//{"result":"success"}
//	//{"code":42,"msg":"该委托已经取消, 不能取消或修改","message":"该委托已经取消, 不能取消或修改"}
//	//println(string(bodyData))

//	var bodyDataMap map[string]interface{}
//	err = json.Unmarshal(bodyData, &bodyDataMap)
//	if err != nil {
//		return false, err
//	}

//	if bodyDataMap["code"] != nil {
//		return false, errors.New(string(bodyData))
//	}

//	ret := bodyDataMap["result"].(string)
//	return (strings.Compare(ret, "success") == 0), nil
//}

//func (btc38 *Btc38) GetKlineRecords(currency CurrencyPair, period string, size, since int) ([]Kline, error) {
//	klineUri := API_BASE_URL + KLINE_URI

//	switch currency {
//	case BTC_CNY:
//		klineUri = fmt.Sprintf(klineUri, "btc", period, size)
//	case LTC_CNY:
//		klineUri = fmt.Sprintf(klineUri, "ltc", period, size)
//	default:
//		return nil, errors.New("Unsupport " + CurrencyPairSymbol[currency])
//	}
//	//println(klineUri)
//	resp, err := http.Get(klineUri)

//	if err != nil {
//		return nil, err
//	}

//	defer resp.Body.Close()

//	body, err := ioutil.ReadAll(resp.Body)

//	var klines [][]interface{}

//	err = json.Unmarshal(body, &klines)

//	if err != nil {
//		return nil, err
//	}

//	loc, _ := time.LoadLocation("Local") //获取本地时区
//	var klineRecords []Kline

//	for _, record := range klines {
//		r := Kline{}
//		for i, e := range record {
//			switch i {
//			case 0:
//				d := e.(string)
//				if len(d) >= 12 {
//					t, _ := time.ParseInLocation("200601021504", d[0:12], loc)
//					r.Timestamp = t.Unix()
//				}
//			case 1:
//				r.Open = e.(float64)
//			case 2:
//				r.High = e.(float64)
//			case 3:
//				r.Low = e.(float64)
//			case 4:
//				r.Close = e.(float64)
//			case 5:
//				r.Vol = e.(float64)
//			}
//		}

//		if r.Timestamp < int64(since/1000) {
//			continue
//		}

//		klineRecords = append(klineRecords, r)
//	}

//	return klineRecords, nil
//}

//func (btc38 *Btc38) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
//	return nil, nil
//}

///**
// * 获取全站最近的交易记录
// */
//func (btc38 *Btc38) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
//	tradeUrl := API_BASE_URL + trade_url
//	switch currencyPair {
//	case BTC_CNY:
//		tradeUrl = fmt.Sprintf(tradeUrl, "btc")
//	case LTC_CNY:
//		tradeUrl = fmt.Sprintf(tradeUrl, "ltc")
//	default:
//		return nil, errors.New("unsupport " + currencyPair.String())
//	}

//	var respmap map[string]interface{}

//	resp, err := http.Get(tradeUrl)
//	if err != nil {
//		return nil, err
//	}

//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	err = json.Unmarshal(body, &respmap)
//	if err != nil {
//		return nil, err
//	}

//	tradesmap, isOK := respmap["trades"].([]interface{})
//	if !isOK {
//		return nil, errors.New("assert error")
//	}

//	now := time.Now()
//	var trades []Trade
//	for _, t := range tradesmap {
//		tr := t.(map[string]interface{})
//		trade := Trade{}
//		trade.Amount = tr["amount"].(float64)
//		trade.Price = tr["price"].(float64)
//		trade.Type = tr["type"].(string)
//		timeStr := tr["time"].(string)
//		timeMeta := strings.Split(timeStr, ":")
//		h, _ := strconv.Atoi(timeMeta[0])
//		m, _ := strconv.Atoi(timeMeta[1])
//		s, _ := strconv.Atoi(timeMeta[2])
//		//临界点处理
//		if now.Hour() == 0 {
//			if h <= 23 && h >= 20 {
//				pre := now.AddDate(0, 0, -1)
//				trade.Date = time.Date(pre.Year(), pre.Month(), pre.Day(), h, m, s, 0, time.Local).Unix() * 1000
//			} else if h == 0 {
//				trade.Date = time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, time.Local).Unix() * 1000
//			}
//		} else {
//			trade.Date = time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, time.Local).Unix() * 1000
//		}
//		//fmt.Println(time.Unix(trade.Date/1000 , 0))
//		trades = append(trades, trade)
//	}

//	//fmt.Println(tradesmap)

//	return trades, nil
//}
