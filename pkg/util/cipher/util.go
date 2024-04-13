package cipher

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5" //nolint:gosec // 暗号化ではなくワンタイムのハッシュにのみ使っている
	"errors"
)

// padding バイト配列へパディングを付与する.
func padding(src []byte, blockSize int) []byte {
	srcLen := len(src)
	padLen := blockSize - (srcLen % blockSize)
	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(src, padText...)
}

// unpadding バイト配列からパディングを除去する.
func unpadding(src []byte, blockSize int) ([]byte, error) {
	srcLen := len(src)
	paddingLen := int(src[srcLen-1])

	if paddingLen >= srcLen || paddingLen > blockSize {
		return nil, errors.New("padding size error")
	}
	return src[:srcLen-paddingLen], nil
}

// Encrypt plain を AES/CBC/PKCS#7 で暗号化する。
func Encrypt(plain, secret, cryptKey []byte) ([]byte, error) {
	block, iv, err := getCipher(secret, cryptKey)
	if err != nil {
		return nil, err
	}
	encrypter := cipher.NewCBCEncrypter(block, iv)

	// PKCS#5 に沿ってパディングを付与
	padded := padding(plain, encrypter.BlockSize())
	// 暗号化
	encrypted := make([]byte, len(padded))
	encrypter.CryptBlocks(encrypted, padded)
	return encrypted, nil
}

// Decrypt encrypted を AES/CBC/PKCS#7 で復号する
func Decrypt(encrypted, secret, cryptKey []byte) ([]byte, error) {
	block, iv, err := getCipher(secret, cryptKey)
	if err != nil {
		return nil, err
	}
	decrypter := cipher.NewCBCDecrypter(block, iv)

	plain := make([]byte, len(encrypted))
	decrypter.CryptBlocks(plain, encrypted)
	// パディングを除去
	return unpadding(plain, decrypter.BlockSize())
}

func getCipher(secret, cryptKey []byte) (cipher.Block, []byte, error) {
	//nolint:gosec // keyBlock
	hasher := md5.New()
	_, err := hasher.Write(secret)
	if err != nil {
		return nil, nil, err
	}
	block, _ := aes.NewCipher(hasher.Sum(nil))

	// iv
	keyByte := make([]byte, 0)
	keyByte = append(keyByte, secret...)
	keyByte = append(keyByte, cryptKey...)
	//nolint:gosec // iv
	hasher = md5.New()
	_, err = hasher.Write(keyByte)
	if err != nil {
		return nil, nil, err
	}
	iv := hasher.Sum(nil)[:aes.BlockSize]

	return block, iv, nil
}
