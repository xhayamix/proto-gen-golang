package randstr

import (
	"crypto/rand"
)

const (
	// ランダム文字列生成用文字
	defaultLetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	letterIdxMask      = 0x3F // 63 0b111111
)

// RandomString 指定した文字列長のランダム文字列を取得する
func RandomString(n int) (string, error) {
	bytes, err := RandomBytes(n)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// RandomString 指定した文字列長のランダム文字列のbyteスライスを取得する
func RandomBytes(n int) ([]byte, error) {
	return RandomBytesByLetterBytes(defaultLetterBytes, n)
}

// RandomString 指定した文字列長のランダム文字列のbyteスライスを取得する
func RandomBytesByLetterBytes(letterBytes string, n int) ([]byte, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}
	for i := 0; i < n; {
		idx := int(buf[i] & letterIdxMask)
		if idx < len(letterBytes) {
			buf[i] = letterBytes[idx]
			i++
		} else if _, err := rand.Read(buf[i : i+1]); err != nil {
			return nil, err
		}
	}
	return buf, nil
}
