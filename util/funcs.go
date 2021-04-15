package util

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	log "github.com/dxasu/gostar/util/glog"
)

const (
	// len: 23
	APREFIX = "ABCDEFGHJKLMNPQRSTUWXYZ"
	// len: 7
	NPREFIX = "2345789"
	BASETS  = 1578412800
	BITXOR  = 0x35D3C5D
)

// NewOrderId uid chid
func NewOrderId(uid uint64, chid int64) string {
	apre := APREFIX[uid%uint64(len(APREFIX))]
	bpre := APREFIX[(uid*101693)%uint64(len(APREFIX))]
	npre := NPREFIX[uid%uint64(len(NPREFIX))]
	sec := time.Now().Unix()
	maskSec := (sec - BASETS) ^ BITXOR
	num := rand.Intn(999983)
	for i := 1; i < 5; i++ {
		if num < 100_000 {
			num = rand.Intn(999983)
		} else {
			break
		}
	}
	return fmt.Sprintf("%c%c%c%do%do%d", apre, bpre, npre, chid,
		maskSec, num)
}

func abs(num int64) int64 {
	if num < 0 {
		num = -num
	}
	return num
}

//随机字符串
func RandomString(length int) string {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//去掉换行,空格,回车,制表
func TrimRightSpace(s string) string {
	return strings.TrimRight(string(s), "\r\n\t ")
}

func GetDgit(text string) string {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		log.Warningln("regexp.Compile err", err)
		return ""
	}
	return reg.ReplaceAllString(text, "")

}

func URandom() []byte {
	f, _ := os.Open("/dev/urandom")
	b := make([]byte, 16)
	f.Read(b)
	f.Close()

	return b
}

//获取当前执行文件的目录
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	//返回绝对路径
	filepath.Dir(os.Args[0])
	if err != nil {
		log.Fatalln(err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

//获取执行文件当前的目录 最后带/
func GetRootDir() string {
	// 文件不存在获取执行路径
	file, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		file = fmt.Sprintf(".%s", string(os.PathSeparator))
	} else {
		file = fmt.Sprintf("%s%s", file, string(os.PathSeparator))
	}
	return file
}

//获取执行文件的路径
func GetExecFilePath() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		file = fmt.Sprintf(".%s", string(os.PathSeparator))
	} else {
		file, err = filepath.Abs(file)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	return file
}
