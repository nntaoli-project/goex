package huobi

import (
	"fmt"
	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/internal/logger"
	"sort"
	"strings"
)

func ParseDepthFromResponse(r DepthResponse) goex.Depth {
	var dep goex.Depth
	for _, bid := range r.Bids {
		dep.BidList = append(dep.BidList, goex.DepthRecord{Price: bid[0], Amount: bid[1]})
	}

	for _, ask := range r.Asks {
		dep.AskList = append(dep.AskList, goex.DepthRecord{Price: ask[0], Amount: ask[1]})
	}

	sort.Sort(sort.Reverse(dep.BidList))
	sort.Sort(sort.Reverse(dep.AskList))
	return dep
}

func ParseCurrencyPairFromSpotWsCh(ch string) goex.CurrencyPair {
	meta := strings.Split(ch, ".")
	if len(meta) < 2 {
		logger.Errorf("parse error, ch=%s", ch)
		return goex.UNKNOWN_PAIR
	}

	currencyPairStr := meta[1]
	if strings.HasSuffix(currencyPairStr, "usdt") {
		currencyA := strings.TrimSuffix(currencyPairStr, "usdt")
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_usdt", currencyA))
	}

	if strings.HasSuffix(currencyPairStr, "btc") {
		currencyA := strings.TrimSuffix(currencyPairStr, "btc")
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_btc", currencyA))
	}

	if strings.HasSuffix(currencyPairStr, "eth") {
		currencyA := strings.TrimSuffix(currencyPairStr, "eth")
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_eth", currencyA))
	}

	if strings.HasSuffix(currencyPairStr, "husd") {
		currencyA := strings.TrimSuffix(currencyPairStr, "husd")
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_husd", currencyA))
	}

	if strings.HasSuffix(currencyPairStr, "ht") {
		currencyA := strings.TrimSuffix(currencyPairStr, "ht")
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_ht", currencyA))
	}

	if strings.HasSuffix(currencyPairStr, "trx") {
		currencyA := strings.TrimSuffix(currencyPairStr, "trx")
		return goex.NewCurrencyPair2(fmt.Sprintf("%s_trx", currencyA))
	}

	return goex.UNKNOWN_PAIR
}
