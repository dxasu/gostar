package encryp

import "testing"

func Test_Hmac(t *testing.T) {
	text := "你好 世界 hello world"
	key := "abcabc123123"

	HMAC := GenerateHMAC(text, key)
	t.Log(HMAC)

	verifyHMAC := VerifyHMAC(HMAC, text, key)
	t.Log(verifyHMAC)
}
