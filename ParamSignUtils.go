package goex

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
)

/**
 *md5签名,okcoin和huobi适用
 */
func GetParamMD5Sign(secret, params string) (string, error) {
	hash := md5.New()
	_, err := hash.Write([]byte(params))

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func GetSHA(text string) (string, error) {
	sha := sha1.New()
	_, err := sha.Write([]byte(text))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sha.Sum(nil)), nil
}

func GetParamHmacSHA256Sign(secret, params string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func GetParamHmacSHA512Sign(secret, params string) (string, error) {
	mac := hmac.New(sha512.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func GetParamHmacSHA1Sign(secret, params string) (string, error) {
	mac := hmac.New(sha1.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func GetParamHmacMD5Sign(secret, params string) (string, error) {
	mac := hmac.New(md5.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

/*for bitfinex.com*/
func GetParamHmacSha384Sign(secret, params string) (string, error) {
	mac := hmac.New(sha512.New384, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", nil
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func GetParamHmacSHA256Base64Sign(secret, params string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	signByte := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(signByte), nil
}

func GetParamHmacSHA512Base64Sign(hmac_key string, hmac_data string) string {
	hmh := hmac.New(sha512.New, []byte(hmac_key))
	hmh.Write([]byte(hmac_data))

	hex_data := hex.EncodeToString(hmh.Sum(nil))
	hash_hmac_bytes := []byte(hex_data)
	hmh.Reset()

	return base64.StdEncoding.EncodeToString(hash_hmac_bytes)
}
