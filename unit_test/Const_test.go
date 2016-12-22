package unit

import (
	. "../"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Const(t *testing.T) {
	assert.Equal(t, BTC_CNY, 1);
	assert.Equal(t, BTC_USD, 2);

	var currency CurrencyPair;
	currency = BTC_CNY;
	assert.True(t, currency == BTC_CNY); //值比较
	assert.EqualValues(t, currency, BTC_CNY); //值比较
	assert.NotEqual(t, currency, BTC_CNY); //类型不同

	var side TradeSide;
	side = SELL;
	t.Logf("side = %d", side);
	assert.True(t, side == SELL);
}
