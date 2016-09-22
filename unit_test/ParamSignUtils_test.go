package unit_test

import (
	"testing"
	. "../"
	"github.com/stretchr/testify/assert"
	"strings"
)

func Test_appendParamSign(t *testing.T) {
	params := "a=a&b=b";
	secret := "secret";

	sign, _ := GetParamMD5Sign(secret, params);
	assert.Equal(t, strings.ToUpper(sign), "C984E55DA4789949626EA9B26C4869F6", "md5 sign fail");

	hmacSha256Sign, _ := GetParamHmacSHA256Sign(secret, params);
	assert.Equal(t, hmacSha256Sign, "c5b19a95474fc71f7e054d596452b8fe7c3ff517b5a05c12804da4ba6f987664", "hmac sha256 sign fail");

	hmacSha1Sign, _ := GetParamHmacSHA1Sign(secret, params);
	assert.Equal(t, hmacSha1Sign, "1c0a1aabae2dd88e728f6d84c46656b1b07d0fce", "hmac sha1 sign fail");

	hmacSha512Sign, _ := GetParamHmacSHA512Sign(secret, params);
	assert.Equal(t, hmacSha512Sign,
		"46b7c47269ead767a7ce2d4b385cba52f0450e2ea6ffb68c11ce620453a7a3ab62e36082775c94bd6859fccba1f9819a83cc94dbfb416c7c8c510882463ce3a2",
		"hmac sha512 sign fail");

	secretSHA, _ := GetSHA(secret);
	hmacMD5Sign, _ := GetParamHmacMD5Sign(secretSHA, params);
	assert.Equal(t, "f60000cefb83b9848464666e82311306", hmacMD5Sign)
}

