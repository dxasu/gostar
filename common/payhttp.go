package common

import (
	"bytes"
	"io/ioutil"
	"net/http"

	log "github.com/dxasu/gostar/util/glog"
)

type HEAD map[string]string

// "Content-Type": "application/x-www-form-urlencoded",
// "Accept": "application/json",

// GetRequestMsg
func GetRequestMsg(method, url string, content []byte, header HEAD) (string, error) {
	log.Infof("getRequestMsg NewRequest method:%s url:%s\n", method, url)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(content))

	if err != nil {
		return "", err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	return DoRequest(req)
}

// DoRequest
func DoRequest(req *http.Request) (string, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	msg := string(body)
	log.Infoln("status:", resp.Status, "response:", msg)

	return msg, err
}
