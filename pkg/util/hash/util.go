package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256(base string, salts ...string) string {
	// salt指定がない場合はそのままハッシュ化
	if len(salts) == 0 {
		bytes := sha256.Sum256([]byte(base))
		return hex.EncodeToString(bytes[:])
	}
	hashed := base
	// 指定されたsalt回数分ハッシュ化を繰り返す
	for _, salt := range salts {
		baseBytes := []byte(hashed)
		saltBytes := []byte(salt)
		hashBytes := make([]byte, 0, len(baseBytes)+len(saltBytes))
		hashBytes = append(hashBytes, baseBytes...)
		hashBytes = append(hashBytes, saltBytes...)
		bytes := sha256.Sum256(hashBytes)
		hashed = hex.EncodeToString(bytes[:])
	}
	return hashed
}
