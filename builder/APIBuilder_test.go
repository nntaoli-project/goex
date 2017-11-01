package builder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var builder = NewAPIBuilder()

func TestAPIBuilder_Build(t *testing.T) {
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build("okcoin.cn").GetExchangeName(), "okcoin.cn")
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build("okcoin.com").GetExchangeName(), "okcoin.com")
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build("huobi.com").GetExchangeName(), "huobi.com")
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build("chbtc.com").GetExchangeName(), "chbtc.com")
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build("yunbi.com").GetExchangeName(), "yunbi.com")
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build("poloniex.com").GetExchangeName(), "poloniex.com")
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build("coincheck.com").GetExchangeName(), "coincheck.com")
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build("zaif.jp").GetExchangeName(), "zaif.jp")
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build("huobi.pro").GetExchangeName(), "huobi.pro")
}
