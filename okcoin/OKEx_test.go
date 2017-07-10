package okcoin

import (
	. "github.com/nntaoli-project/GoEx"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var (
	okex = NewOKEx(http.DefaultClient, "", "")
)

func TestOKEx_GetFutureDepth(t *testing.T) {
	dep, err := okex.GetFutureDepth(BTC_USD, QUARTER_CONTRACT, 1)
	assert.Nil(t, err)
	t.Log(dep)
}
