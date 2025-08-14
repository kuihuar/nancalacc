package wps

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type KsoSign struct {
	accessKey string
	secretKey string
}

type Out struct {
	Date          string // X-Kso-Date
	Authorization string // X-Kso-Authorization
}

const (
	ContentType   = "Content-Type"
	KsoAuthHeader = "X-Kso-Authorization"
	KsoDateHeader = "X-Kso-Date"
	Kso1Version   = "KSO-1"
)

func NewKsoSign(accessKey, secretKey string) (*KsoSign, error) {
	if accessKey == "" || secretKey == "" {
		return nil, errors.New("NewKsoSign error: AccessKey/SecretKey can not be empty")
	}
	return &KsoSign{
		accessKey: accessKey,
		secretKey: secretKey,
	}, nil
}

// signUri := strings.TrimLeft(req.URL.Path, openApiPath)
func (k *KsoSign) getKso1Signature(secretKey, method, uri, ksoDate, contentType string, requestBody []byte) string {
	sha256Hex := ""
	if len(requestBody) > 0 {
		s := sha256.New()
		s.Write(requestBody)
		sha256Hex = hex.EncodeToString(s.Sum(nil))
	}

	mac := hmac.New(sha256.New, []byte(secretKey))

	signatureByte := []byte("KSO-1" + method + uri + contentType + ksoDate + sha256Hex)

	// fmt.Printf("signature origin: %s\n", string(signatureByte))
	// fmt.Printf("signature sha256Hex body: %s\n", sha256Hex)
	mac.Write(signatureByte)

	return hex.EncodeToString(mac.Sum(nil))
}

func (k *KsoSign) KSO1Sign(method, signPath, contentType, ksoDate string, body []byte) (*Out, error) {

	//fmt.Printf("[KSO1Sign] method: %s, signPath: %s, contentType: %s, ksoDate: %s, body: %s\n", method, signPath, contentType, ksoDate, string(body))
	ksoSignature := k.getKso1Signature(k.secretKey, method, signPath, ksoDate, contentType, body)
	authorization := fmt.Sprintf("%s %s:%s", "KSO-1", k.accessKey, ksoSignature)
	//fmt.Printf("[KSO1Sign authorization]: %s\n", authorization)
	return &Out{
		Date:          ksoDate,
		Authorization: authorization,
	}, nil
}

func (k *KsoSign) validDate(ksoDate string) (tm time.Time, err error) {
	if tm, err = time.Parse(time.RFC1123, ksoDate); err == nil {
		return
	}
	if tm, err = time.Parse(time.RFC1123Z, ksoDate); err == nil {
		return
	}

	// 或者使用星期是全拼的非标准RFC1123
	RFC1123 := "Monday, 02 Jan 2006 15:04:05 MST"
	if tm, err = time.Parse(RFC1123, ksoDate); err == nil {
		return
	}

	return tm, errors.New("kso1sign check error: invalid kso-date header")
}

func (k *KsoSign) extractSign(authHeader string) (string, string, string, error) {
	spaceParts := strings.SplitN(authHeader, " ", 2)
	if len(spaceParts) != 2 {
		return "", "", "", errors.New("kso1sign check error: invalid authorization header")
	}

	colonParts := strings.SplitN(spaceParts[1], ":", 2)
	if len(colonParts) != 2 {
		return "", "", "", errors.New("kso1sign check error: invalid authorization header")
	}

	return spaceParts[0], colonParts[0], colonParts[1], nil
}

func (k *KsoSign) validSign(req *http.Request, ksoSignature, secretKey, ksoDate, contentType string, requestBody []byte) error {
	sha256Hex := ""
	if len(requestBody) > 0 {
		s := sha256.New()
		s.Write(requestBody)
		sha256Hex = hex.EncodeToString(s.Sum(nil))
	}

	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte("KSO-1" + req.Method + req.URL.RequestURI() + contentType + ksoDate + sha256Hex))
	sign := hex.EncodeToString(mac.Sum(nil))

	if sign == ksoSignature {
		return nil
	}

	return errors.New("kso1sign check error: invalid signature")
}

func (k *KsoSign) Kso1SignCheck(req *http.Request, requestBody []byte, SKGetter func(string) (string, error)) error {
	ksoAuthHeader := req.Header.Get(KsoAuthHeader)
	ksoDate := req.Header.Get(KsoDateHeader)
	contentType := req.Header.Get(ContentType)

	_, err := k.validDate(ksoDate)
	if err != nil {
		return err
	}

	signVersion, accessKey, ksoSignature, err := k.extractSign(ksoAuthHeader)
	if err != nil {
		return err
	}
	if signVersion != Kso1Version {
		return errors.New("kso1sign check error: unknown authorization version")
	}

	// 获取 SK 及校验 AK 合法性
	secretKey, err := SKGetter(accessKey)
	if err != nil {
		return err
	}

	return k.validSign(req, ksoSignature, secretKey, ksoDate, contentType, requestBody)
}
