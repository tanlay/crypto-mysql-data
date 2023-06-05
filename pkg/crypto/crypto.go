package crypto

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/tjfoc/gmsm/sm4"
)

// SM4 ECB加密
// @param hexKey 16进制key 长度32位
// @param raw 待加密内容

func SM4ECBEncrypt(hexKey, raw string) (string, error) {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return "", err
	}
	out, err := sm4.Sm4Ecb(key, []byte(raw), true)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(out), nil
}

// SM4 ECB解密
// @param hexKey 16进制key 长度32位
// @param base64Raw 加密内容 base64格式

func SM4ECBDecrypt(hexKey, base64Raw string) (string, error) {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return "", err
	}

	raw, err := base64.StdEncoding.DecodeString(base64Raw)
	if err != nil {
		return "", err
	}

	out, err := sm4.Sm4Ecb(key, raw, false)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
