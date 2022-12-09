package goex

import (
	"github.com/nntaoli-project/goex/v2/httpcli"
	"github.com/nntaoli-project/goex/v2/logger"
	"reflect"
)

var (
	DefaultHttpCli = httpcli.Cli
)

func SetDefaultHttpCli(cli httpcli.IHttpClient) {
	logger.Infof("use new http client implement: %s", reflect.TypeOf(cli).Elem().String())
	httpcli.Cli = cli
}
