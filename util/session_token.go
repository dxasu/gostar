package util

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	log "github.com/dxasu/gostar/util/glog"

	"google.golang.org/grpc/metadata"
)

const (
	DEVICE_TYPE      = "1"
	RPC_ACCESS_TOKEN = "1"
	FIXED_BITS       = 16
)

var (
	key = []byte("\x2e\xe2\x53\xba\x33\x14\x59\x48\xa0\xa4\x4e\x3c\x56\x3c\xa7\xb6")
	iv  = []byte("\x38\x37\xf5\x98\x84\xf7\x41\x0c\x2f\x05\xa3\x15\x79\x86\xd1\x5a")
)

func pKCS7Padding(ciphertext []byte) []byte {
	padding := FIXED_BITS - len(ciphertext)%FIXED_BITS
	padtext := bytes.Repeat([]byte{byte('\n')}, padding)
	return append(ciphertext, padtext...)
}

func pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func aesCBCEncrypt(rawData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	rawData = PKCS5Padding(rawData, 16)
	cipherText := make([]byte, len(rawData))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, rawData)
	return cipherText, nil
}

func aesCBCDncrypt(encryptData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	if len(encryptData)%FIXED_BITS != 0 {
		panic("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptData, encryptData)
	encryptData = pKCS7UnPadding(encryptData)
	return encryptData, nil
}

func encrypt(rawData []byte) (string, error) {
	data, err := aesCBCEncrypt(rawData)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}

func decrypt(encrytedStr string) (string, error) {
	data, err := hex.DecodeString(encrytedStr)
	if err != nil {
		return "", err
	}
	dnData, err := aesCBCDncrypt(data)
	if err != nil {
		return "", err
	}
	return string(dnData), nil
}

func EncryptTicket(uuid int) (string, error) {
	expireSeconds := time.Now().Add(time.Hour * 24 * 2).Unix()
	//expireSeconds := time.Now().Add(time.Minute * 2).Unix()

	plainText := fmt.Sprintf("%s:%s:%s:%s", strconv.Itoa(uuid), DEVICE_TYPE, strconv.FormatInt(expireSeconds, 10), RPC_ACCESS_TOKEN)

	log.Infoln("plainText: ", plainText)
	return encrypt([]byte(plainText))
}

func DecryptTicket(encrytedStr string) (int64, int64, error) {
	plainText, err := decrypt(encrytedStr)
	if err != nil {
		return -1, -1, err
	}
	arrays := strings.Split(plainText, ":")
	if len(arrays) != 4 {
		return -1, -1, fmt.Errorf("invalid plainText(%v)", plainText)
	}

	uid, err := strconv.ParseInt(arrays[0], 10, 64)
	if err != nil {
		return -1, -1, fmt.Errorf("invalid uid from plainText(%v)", plainText)
	}

	expiredTime, err := strconv.ParseInt(arrays[2], 10, 64)
	if err != nil {
		return -1, -1, fmt.Errorf("invalid expiredtime from plainText(%v)", plainText)
	}
	return uid, expiredTime, nil
}

func VerifyToken(uid int64, ticket string) error {
	ticketUID, expiredTime, err := DecryptTicket(ticket)
	if err != nil {
		return err
	}

	if uid != ticketUID {
		return fmt.Errorf("invalid uid ")
	}

	if expiredTime < time.Now().Unix() {
		return fmt.Errorf("invalid expiredtime")
	}

	return nil
}

func CheckTokenAndFetchId(ctx context.Context) (int64, error) {
	uid, ticket, err := fetchTokenFromHeader(ctx)
	if err != nil {
		return 0, err
	}

	return uid, VerifyToken(uid, ticket)
}

func fetchTokenFromHeader(ctx context.Context) (int64, string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Error("No header info in request")
		return 0, "", fmt.Errorf("no header info in reques")
	}
	vlus := md.Get("x-auth-token")
	if len(vlus) == 0 || vlus[0] == "" {
		log.Error("No token in request.header")
		return 0, "", fmt.Errorf("no token in request.header")
	}
	strUid := md.Get("uid")
	if len(vlus) == 0 || vlus[0] == "" {
		log.Error("No token in request.header")
		return 0, "", fmt.Errorf("no token in request.header")
	}

	uid, err := strconv.ParseInt(strUid[0], 10, 64)
	if err != nil {
		return 0, "", fmt.Errorf("invalid uid")
	}

	return uid, vlus[0], nil
}

//----------------------------------------------tcp token --------------------------------------------------------------

const ks = "abcdefghijkmlnop"
const alphanum = "ABCDEFGHJKLMNPQRSTUWXYZabcdefghjklmnpqrstuwxyz2345789"

func init() {
	rand.Seed(time.Now().Unix())

}

// GenRandKey gen key for app  sign handshake req
func GenRandKey() string {
	curLen := 0
	// 16 bytes
	key := []byte("0123456789123456")
	rnd := rand.Int63()

	for {
		for rnd > 0 {
			if curLen == 16 {
				return string(key)
			}
			key[curLen] = alphanum[rnd%53]
			curLen++
			rnd /= 53
		}
		rnd = rand.Int63()
	}

}

// AESBase64Encrypt  gen tcp token
func AESBase64Encrypt(originData []byte) (base64Result string, err error) {
	key := md5.Sum([]byte(ks))
	iv := bytes.Repeat([]byte{byte(0)}, 16)
	var block cipher.Block
	if block, err = aes.NewCipher(key[:]); err != nil {
		log.Infoln(err)
		return
	}
	encrypt := cipher.NewCBCEncrypter(block, iv)
	var source []byte = PKCS5Padding(originData, 16)
	var dst []byte = make([]byte, len(source))
	encrypt.CryptBlocks(dst, source)
	base64Result = base64.StdEncoding.EncodeToString(dst)
	return
}

// AESBase64Decrypt descrypt tcp token
func AESBase64Decrypt(encryptedData string) (originData []byte, err error) {
	key := md5.Sum([]byte(ks))
	iv := bytes.Repeat([]byte{byte(0)}, 16)
	var block cipher.Block
	if block, err = aes.NewCipher(key[:]); err != nil {
		log.Infoln(err)
		return
	}
	encrypt := cipher.NewCBCDecrypter(block, iv)

	var source []byte
	if source, err = base64.StdEncoding.DecodeString(encryptedData); err != nil {
		log.Infoln(err)
		return
	}
	var dst []byte = make([]byte, len(source))
	encrypt.CryptBlocks(dst, source)
	originData = PKCS5Unpadding(dst)
	return
}

// PKCS5Padding padding according to PKCS5
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS5Unpadding remove tail padding
func PKCS5Unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// VerSign check app has right key, and now sign algthrim
func VerSign(plaintext, key, digest []byte) bool {
	mac := hmac.New(md5.New, key)
	mac.Write(plaintext)
	return bytes.Equal(mac.Sum(nil), digest)
}
