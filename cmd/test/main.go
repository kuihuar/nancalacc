package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

type KsoSign struct {
	accessKey string
	secretKey string
}

type Out struct {
	Date          string // X-Kso-Date
	Authorization string // X-Kso-Authorization
}

func NewKsoSign(accessKey, secretKey string) (*KsoSign, error) {
	if accessKey == "" || secretKey == "" {
		return nil, errors.New("NewKsoSign error: AccessKey/SecretKey can not be empty")
	}
	return &KsoSign{
		accessKey: accessKey,
		secretKey: secretKey,
	}, nil
}

func (k *KsoSign) getKso1Signature(secretKey, method, uri, ksoDate, contentType string, requestBody []byte) string {
	sha256Hex := ""
	if len(requestBody) > 0 {
		s := sha256.New()
		s.Write(requestBody)
		sha256Hex = hex.EncodeToString(s.Sum(nil))
	}

	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte("KSO-1" + method + uri + contentType + ksoDate + sha256Hex))
	return hex.EncodeToString(mac.Sum(nil))
}

func (k *KsoSign) KSO1Sign(method, uri, contentType, ksoDate string, body []byte) (*Out, error) {

	fmt.Println()
	fmt.Println("KSO1Sign start:")
	fmt.Printf("accessKey: %s\n", k.accessKey)
	fmt.Printf("secretKey: %s\n", k.secretKey)
	fmt.Printf("method: %s\n", method)
	fmt.Printf("signPath: %s\n", uri)
	fmt.Printf("contentType: %s\n", contentType)
	fmt.Printf("ksoDate: %s\n", ksoDate)
	fmt.Printf("body: %s\n", string(body))
	fmt.Println("KSO1Sign end:")
	fmt.Println()

	ksoSignature := k.getKso1Signature(k.secretKey, method, uri, ksoDate, contentType, body)
	authorization := fmt.Sprintf("%s %s:%s", "KSO-1", k.accessKey, ksoSignature)
	return &Out{
		Date:          ksoDate,
		Authorization: authorization,
	}, nil
}

func main() {
	accessKey := "AK123456"
	secretKey := "sk098765"
	method := "POST"
	uri := "/v7/test/body"
	contentType := "application/json"
	ksoDate := "Mon, 02 Jan 2006 15:04:05 GMT"
	body := `{"key":"value"}` // 注意 json 格式，会影响到签名计算

	sign, err := NewKsoSign(accessKey, secretKey)
	if err != nil {
		panic(err)
	}

	out, err := sign.KSO1Sign(method, uri, contentType, ksoDate, []byte(body))
	if err != nil {
		panic(err)
	}
	fmt.Printf("out: %v\n", out)
	// 输出：out: &{Mon, 02 Jan 2006 15:04:05 GMT KSO-1 AK123456:c46e6c988130818ecba2484d51ac685948fbbef6814602c7874d6bfc41dc17b3}
}
