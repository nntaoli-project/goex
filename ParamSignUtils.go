package coinapi

import (
	"crypto/md5"
	"encoding/hex"
	"crypto/sha256"
	"crypto/hmac"
	"crypto/sha512"
	"crypto/sha1"
)

/**
 *md5签名,okcoin和huobi适用
 */
func GetParamMD5Sign(secret, params string) (string, error) {
	hash := md5.New();
	_, err := hash.Write([]byte(params));

	if err != nil {
		return "", err;
	}

	return hex.EncodeToString(hash.Sum(nil)), nil;
}

func GetSHA(text string) (string, error) {
	sha := sha1.New();
	_, err := sha.Write([]byte(text))
	if err != nil {
		return "", err;
	}
	return hex.EncodeToString(sha.Sum(nil)), nil;
}

func GetParamHmacSHA256Sign(secret, params string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret));
	_, err := mac.Write([]byte(params));
	if err != nil {
		return "", err;
	}
	return hex.EncodeToString(mac.Sum(nil)), nil;
}

func GetParamHmacSHA512Sign(secret, params string) (string, error) {
	mac := hmac.New(sha512.New, []byte(secret));
	_, err := mac.Write([]byte(params));
	if err != nil {
		return "", err;
	}
	return hex.EncodeToString(mac.Sum(nil)), nil;
}

func GetParamHmacSHA1Sign(secret, params string) (string, error) {
	mac := hmac.New(sha1.New, []byte(secret));
	_, err := mac.Write([]byte(params));
	if err != nil {
		return "", err;
	}
	return hex.EncodeToString(mac.Sum(nil)), nil;
}

func GetParamHmacMD5Sign(secret, params string) (string, error) {
	mac := hmac.New(md5.New, []byte(secret));
	_, err := mac.Write([]byte(params));
	if err != nil {
		return "", err;
	}
	return hex.EncodeToString(mac.Sum(nil)), nil;
}