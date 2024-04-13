package cipher

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/xhayamix/proto-gen-golang/pkg/cerrors"
)

func ParsePrivateKey(pemData []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
	return key, nil
}

func DecryptPKCS1(priv *rsa.PrivateKey, cipherText []byte) ([]byte, error) {
	rawData, err := rsa.DecryptPKCS1v15(rand.Reader, priv, cipherText)
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	return rawData, nil
}
