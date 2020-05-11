package builder

import (
	"github.com/nntaoli-project/goex"
	"github.com/stretchr/testify/assert"
	"testing"
)

var builder = NewAPIBuilder()

func TestAPIBuilder_Build(t *testing.T) {
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(goex.OKCOIN_COM).GetExchangeName(), goex.OKCOIN_COM)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(goex.HUOBI_PRO).GetExchangeName(), goex.HUOBI_PRO)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(goex.ZB).GetExchangeName(), goex.ZB)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(goex.BIGONE).GetExchangeName(), goex.BIGONE)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(goex.OKEX).GetExchangeName(), goex.OKEX)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(goex.POLONIEX).GetExchangeName(), goex.POLONIEX)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(goex.KRAKEN).GetExchangeName(), goex.KRAKEN)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(goex.FCOIN_MARGIN).GetExchangeName(), goex.FCOIN_MARGIN)
	assert.Equal(t, builder.APIKey("").APISecretkey("").BuildFuture(goex.HBDM).GetExchangeName(), goex.HBDM)
}
