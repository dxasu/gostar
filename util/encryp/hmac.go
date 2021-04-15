package encryp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// VerifyHMAC Hmac
func VerifyHMAC(key, text, HMAC string) bool {

	HMACBytes := []byte(HMAC)

	nowHMAC := GenerateHMAC(key, text)
	nowHMACBytes := []byte(nowHMAC)

	return hmac.Equal(HMACBytes, nowHMACBytes)
}

// GenerateHMAC text key
func GenerateHMAC(key, text string) string {

	keyBytes := []byte(key)

	hash := hmac.New(sha256.New, keyBytes)

	hash.Write([]byte(text))

	result := hash.Sum(nil)

	return hex.EncodeToString(result)
}
