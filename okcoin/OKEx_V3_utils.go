package okcoin

import (
	"fmt"
	"strconv"
	"time"

	. "github.com/nntaoli-project/GoEx"
)

type IContractIDProvider interface {
	GetContractID(CurrencyPair, string) (string, error)
	ParseContractID(string) (CurrencyPair, string, error)
}

type OKExV3DataParser struct {
	contractIDProvider IContractIDProvider
}

func NewOKExV3DataParser(contractIDProvider IContractIDProvider) *OKExV3DataParser {
	return &OKExV3DataParser{contractIDProvider: contractIDProvider}
}

func (okV3dp *OKExV3DataParser) ParseFutureTicker(data interface{}) (*FutureTicker, error) {
	var fallback *FutureTicker
	switch v := data.(type) {
	case map[string]interface{}:
		contractID := v["instrument_id"].(string)
		currencyPair, contractType, err := okV3dp.contractIDProvider.ParseContractID(contractID)
		if err != nil {
			return fallback, err
		}
		t := new(Ticker)
		t.Pair = currencyPair
		timestamp, _ := timeStringToInt64(v["timestamp"].(string))
		t.Date = uint64(timestamp)
		t.Buy, _ = strconv.ParseFloat(v["best_ask"].(string), 64)
		t.Sell, _ = strconv.ParseFloat(v["best_bid"].(string), 64)
		t.Last, _ = strconv.ParseFloat(v["last"].(string), 64)
		t.High, _ = strconv.ParseFloat(v["high_24h"].(string), 64)
		t.Low, _ = strconv.ParseFloat(v["low_24h"].(string), 64)
		t.Vol, _ = strconv.ParseFloat(v["volume_24h"].(string), 64)
		ticker := new(FutureTicker)
		ticker.ContractType = contractType
		ticker.Ticker = t
		return ticker, nil
	}

	return fallback, fmt.Errorf("unknown FutureTicker data: %v", data)
}

func (okV3dp *OKExV3DataParser) ParseDepth(depth *Depth, data interface{}, size int) (*Depth, error) {
	var fallback *Depth
	if depth == nil {
		depth = new(Depth)
	}
	switch v := data.(type) {
	case map[string]interface{}:
		if !okV3dp.checkContractInfo(depth.Pair, depth.ContractType) {
			if v["instrument_id"] != nil {
				contractID := v["instrument_id"].(string)
				currencyPair, contractType, err := okV3dp.contractIDProvider.ParseContractID(contractID)
				if err != nil {
					return fallback, err
				}
				depth.Pair = currencyPair
				depth.ContractType = contractType
			}
		}
		var err error
		var timeStr string
		//name of timestamp field is different between swap and future api
		if v["time"] != nil {
			timeStr = v["time"].(string)
		} else if v["timestamp"] != nil {
			timeStr = v["timestamp"].(string)
		} else {
			return fallback, fmt.Errorf("no time field in %v", v)
		}

		depth.UTime, err = time.Parse(time.RFC3339, timeStr)
		if err != nil {
			return fallback, err
		}

		size2 := len(v["asks"].([]interface{}))
		skipSize := 0
		if size < size2 {
			skipSize = size2 - size
		}

		for _, v := range v["asks"].([]interface{}) {
			if skipSize > 0 {
				skipSize--
				continue
			}

			var dr DepthRecord
			for i, vv := range v.([]interface{}) {
				switch i {
				case 0:
					dr.Price, err = strconv.ParseFloat(vv.(string), 64)
					if err != nil {
						return fallback, err
					}
				case 1:
					dr.Amount, err = strconv.ParseFloat(vv.(string), 64)
					if err != nil {
						return fallback, err
					}
				}
			}
			depth.AskList = append(depth.AskList, dr)
		}

		for _, v := range v["bids"].([]interface{}) {
			var dr DepthRecord
			for i, vv := range v.([]interface{}) {
				switch i {
				case 0:
					dr.Price, err = strconv.ParseFloat(vv.(string), 64)
					if err != nil {
						return fallback, err
					}
				case 1:
					dr.Amount, err = strconv.ParseFloat(vv.(string), 64)
					if err != nil {
						return fallback, err
					}
				}
			}
			depth.BidList = append(depth.BidList, dr)

			size--
			if size == 0 {
				break
			}
		}
		return depth, nil
	}

	return fallback, fmt.Errorf("unknown Depth data: %v", data)
}

var emptyPair = CurrencyPair{}

func (okV3dp *OKExV3DataParser) checkContractInfo(currencyPair CurrencyPair, contractType string) bool {
	if currencyPair.Eq(emptyPair) || currencyPair.Eq(UNKNOWN_PAIR) {
		return false
	}
	if contractType == "" {
		return false
	}
	return true
}

func (okV3dp *OKExV3DataParser) ParseFutureOrder(data interface{}) (*FutureOrder, string, error) {
	var fallback *FutureOrder
	switch v := data.(type) {
	case map[string]interface{}:
		contractID := v["instrument_id"].(string)
		currencyPair, contractType, err := okV3dp.contractIDProvider.ParseContractID(contractID)
		if err != nil {
			return fallback, "", err
		}

		futureOrder := &FutureOrder{}
		// swap orderID is not in int format, so just skip this error
		futureOrder.OrderID, _ = strconv.ParseInt(v["order_id"].(string), 10, 64)
		futureOrder.OrderID2 = v["order_id"].(string)
		futureOrder.Amount, err = strconv.ParseFloat(v["size"].(string), 64)
		if err != nil {
			return fallback, "", err
		}
		futureOrder.Price, err = strconv.ParseFloat(v["price"].(string), 64)
		if err != nil {
			return fallback, "", err
		}
		futureOrder.AvgPrice, err = strconv.ParseFloat(v["price_avg"].(string), 64)
		if err != nil {
			return fallback, "", err
		}
		futureOrder.DealAmount, err = strconv.ParseFloat(v["filled_qty"].(string), 64)
		if err != nil {
			return fallback, "", err
		}
		futureOrder.Fee, err = strconv.ParseFloat(v["fee"].(string), 64)
		if err != nil {
			return fallback, "", err
		}
		if i, err := strconv.ParseInt(v["type"].(string), 10, 64); err == nil {
			futureOrder.OType = int(i)
		} else {
			return fallback, "", err
		}
		futureOrder.OrderTime, err = timeStringToInt64(v["timestamp"].(string))
		if err != nil {
			return fallback, "", err
		}
		// leverage not appear in swap
		if v["leverage"] != nil {
			i, err := strconv.ParseInt(v["leverage"].(string), 10, 64)
			if err != nil {
				return fallback, "", err
			}
			futureOrder.LeverRate = int(i)
		}
		futureOrder.ContractName = v["instrument_id"].(string)
		futureOrder.Currency = currencyPair

		state, err := strconv.ParseInt(v["state"].(string), 10, 64)
		if err != nil {
			return fallback, "", err
		}
		switch state {
		case 0:
			futureOrder.Status = ORDER_UNFINISH
		case 1:
			futureOrder.Status = ORDER_PART_FINISH
		case 2:
			futureOrder.Status = ORDER_FINISH
		case 4:
			futureOrder.Status = ORDER_CANCEL_ING
		case -1:
			futureOrder.Status = ORDER_CANCEL
		case 3:
			futureOrder.Status = ORDER_UNFINISH
		case -2:
			futureOrder.Status = ORDER_REJECT
		default:
			return fallback, "", fmt.Errorf("unknown order status: %v", v)
		}
		return futureOrder, contractType, nil
	}

	return fallback, "", fmt.Errorf("unknown FutureOrder data: %v", data)
}

func (okV3dp *OKExV3DataParser) ParseTrade(trade *Trade, contractType string, data interface{}) (*Trade, string, error) {
	var fallback *Trade
	if trade == nil {
		trade = new(Trade)
	}
	switch v := data.(type) {
	case map[string]interface{}:
		if !okV3dp.checkContractInfo(trade.Pair, contractType) {
			if v["instrument_id"] != nil {
				contractID := v["instrument_id"].(string)
				currencyPair, _contractType, err := okV3dp.contractIDProvider.ParseContractID(contractID)
				if err != nil {
					return fallback, "", err
				}
				trade.Pair = currencyPair
				contractType = _contractType
			}
		}

		tid, _ := strconv.ParseInt(v["trade_id"].(string), 10, 64)
		direction := v["side"].(string)
		var amountStr string
		// wtf api
		if v["qty"] != nil {
			amountStr = v["qty"].(string)
		} else if v["size"] != nil {
			amountStr = v["size"].(string)
		}
		amount, _ := strconv.ParseFloat(amountStr, 64)
		price, _ := strconv.ParseFloat(v["price"].(string), 64)
		time, _ := timeStringToInt64(v["timestamp"].(string))
		trade.Tid = tid
		trade.Type = AdaptTradeSide(direction)
		trade.Amount = amount
		trade.Price = price
		trade.Date = time
		return trade, contractType, nil
	}
	return fallback, "", fmt.Errorf("unknown Trade data: %v", data)
}
