package wps

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"

	"strings"
	"time"
)

type Wps4Auth struct {
	AccessKey string
	SecretKey string
}

var (
	Wps4AuthSign    = "Wps-Docs-Authorization"
	contentTypeSign = "application/json"
	wps4DateSign    = "Wps-Docs-Date"
	//connectionSign =

)

func NewWps4Auth(accessKey, secretKey string) *Wps4Auth {
	auth := &Wps4Auth{}
	auth.AccessKey = accessKey
	auth.SecretKey = secretKey
	return auth
}

func (a *Wps4Auth) BuildWps4Headers(method string, url *url.URL, data []byte, contentType string) *http.Header {
	if strings.TrimSpace(contentType) == "" {
		contentType = "application/json"
	}

	header := http.Header{}
	auth, date := a.prepare(method, url, data, contentType)
	header.Set(Wps4AuthSign, auth)
	header.Set(contentTypeSign, contentType)
	header.Set(wps4DateSign, date)
	//header.Set(connectionSign, "keep-alive")
	return &header
}

func (a *Wps4Auth) prepare(method string, url *url.URL, data []byte, contentType string) (auth, date string) {
	path := url.Path
	if url.RawQuery != "" {
		path += fmt.Sprintf("?%s", url.RawQuery)
	}

	var content []byte
	if data != nil {
		content = data
	}

	date = time.Now().UTC().Format(http.TimeFormat)
	sig := a.sign(method, path, contentType, date, content)
	auth = fmt.Sprintf("WPS-4 %s:%s", a.AccessKey, sig)

	return
}

func (a *Wps4Auth) sign(method, uri, contentType, date string, content []byte) (sign string) {
	bodySha := ""
	if content != nil {
		bodySha = a.getSha256(content)
	}

	mac := hmac.New(sha256.New, []byte(a.SecretKey))
	mac.Write([]byte("WPS-4"))
	mac.Write([]byte(method))
	mac.Write([]byte(uri))
	mac.Write([]byte(contentType))
	mac.Write([]byte(date))
	mac.Write([]byte(bodySha))

	return hex.EncodeToString(mac.Sum(nil))
}

func (a *Wps4Auth) getSha256(ontent []byte) string {
	h := sha256.New()
	h.Write(ontent)
	return hex.EncodeToString(h.Sum(nil))
}
