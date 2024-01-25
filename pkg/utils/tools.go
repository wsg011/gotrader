package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"hash"
	"net/url"
	"strconv"

	"github.com/spf13/cast"
)

// Sha512Str
func Sha512Str(src string) string {
	h := sha512.New()
	h.Write([]byte(src)) //
	return hex.EncodeToString(h.Sum(nil))
}

// HmacSha256 计算Hmac Sha256签名串
func HmacSha256(message string, secret string) hash.Hash {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return h
}

// HmacSha512 计算Hmac Sha512签名串
func HmacSha512(message string, secret string) hash.Hash {
	key := []byte(secret)
	h := hmac.New(sha512.New, key)
	h.Write([]byte(message))
	return h
}

// GenHexDigest 返回加密后的16进制字符串
func GenHexDigest(h hash.Hash) string {
	return hex.EncodeToString(h.Sum(nil))
}

// GenBase64Digest 返回加密后的base64编码的字符串
func GenBase64Digest(h hash.Hash) string {
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// UrlEncodeParams 对query参数进行url编码，不要求参数顺序
func UrlEncodeParams(paramsDict map[string]interface{}) string {
	if paramsDict == nil {
		return ""
	}
	param := url.Values{}
	for key, val := range paramsDict {
		param.Add(key, cast.ToString(val))
	}
	return param.Encode()
}

// UrlEncodeParamsByKeys 对query参数进行url编码，根据传入的keys顺序进行处理
func UrlEncodeParamsByKeys(paramsDict map[string]interface{}, keys []string) string {
	param := url.Values{}
	for _, k := range keys {
		val := cast.ToString(paramsDict[k])
		param.Add(k, val)
	}
	return param.Encode()
}

func ParseInt(str string) (int64, error) {
	numberInt, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return numberInt, nil
}

func ParseFloat(str string) (float64, error) {
	numberFloat, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}
	return numberFloat, nil
}
