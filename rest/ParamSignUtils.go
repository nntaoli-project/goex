package rest

import (
	"crypto/md5"
	"encoding/hex"
)

/**
 *md5签名,okcoin和huobi适用
 */
func GetParamMD5Sign(secret, params string) (string, error) {
	_params := params +
	"&secret_key=" + secret;

	hash := md5.New();
	_, err := hash.Write([]byte(_params));

	if err != nil {
		return "", err;
	}

	return hex.EncodeToString(hash.Sum(nil)), nil;
}
