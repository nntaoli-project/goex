package unit_test

import (
	"testing"
	"github.com/crypto_coin_api/rest"
	"github.com/stretchr/testify/assert"
	"strings"
)

func Test_appendParamSign(t *testing.T) {
	params := "a=a&b=b";
	secret := "secret";

	sign, _ := rest.GetParamMD5Sign(secret, params);
	assert.Equal(t, strings.ToUpper(sign), "E0AEBA5156E0FEE8BDBC3B3A963C7968");
}

