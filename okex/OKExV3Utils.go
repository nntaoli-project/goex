package okex

import "time"

//
import (
	. "github.com/nntaoli-project/goex"
)

func adaptKLinePeriod(period KlinePeriod) int {
	granularity := -1
	switch period {
	case KLINE_PERIOD_1MIN:
		granularity = 60
	case KLINE_PERIOD_3MIN:
		granularity = 180
	case KLINE_PERIOD_5MIN:
		granularity = 300
	case KLINE_PERIOD_15MIN:
		granularity = 900
	case KLINE_PERIOD_30MIN:
		granularity = 1800
	case KLINE_PERIOD_1H, KLINE_PERIOD_60MIN:
		granularity = 3600
	case KLINE_PERIOD_2H:
		granularity = 7200
	case KLINE_PERIOD_4H:
		granularity = 14400
	case KLINE_PERIOD_6H:
		granularity = 21600
	case KLINE_PERIOD_1DAY:
		granularity = 86400
	case KLINE_PERIOD_1WEEK:
		granularity = 604800
	}
	return granularity
}

func adaptSecondsToKlinePeriod(seconds int) KlinePeriod {
	var p KlinePeriod
	switch seconds {
	case 60:
		p = KLINE_PERIOD_1MIN
	case 180:
		p = KLINE_PERIOD_3MIN
	case 300:
		p = KLINE_PERIOD_5MIN
	case 900:
		p = KLINE_PERIOD_15MIN
	case 1800:
		p = KLINE_PERIOD_30MIN
	case 3600:
		p = KLINE_PERIOD_1H
	case 7200:
		p = KLINE_PERIOD_2H
	case 14400:
		p = KLINE_PERIOD_4H
	case 21600:
		p = KLINE_PERIOD_6H
	case 86400:
		p = KLINE_PERIOD_1DAY
	case 604800:
		p = KLINE_PERIOD_1WEEK
	}
	return p
}

func timeStringToInt64(t string) (int64, error) {
	timestamp, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return 0, err
	}
	return timestamp.UnixNano() / int64(time.Millisecond), nil
}
