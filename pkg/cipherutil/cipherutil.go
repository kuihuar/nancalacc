package cipherutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

var (
	ErrEmptyPlaintext = errors.New("empty plaintext")
	ErrInvalidPadding = errors.New("invalid padding")
	ErrInvalidKey     = errors.New("invalid key")
	ErrIVGeneration   = errors.New("failed to generate initialization vector")
	ErrDecryptFailed  = errors.New("decrypt failed")
	ErrMD5WriteFailed = errors.New("md5 write failed")
	ErrInvalidInput   = errors.New("invalid input")
	ErrAesCipher      = errors.New("aes cipher creation failed")
)

func DecryptByAes(content string, key string) (string, error) {
	if len(content) < 24 {
		return "", ErrEmptyPlaintext
	}
	if len(key) == 0 {
		return "", ErrInvalidKey
	}
	// 使用MD5将应用SK转换为32位十六进制字符串作为AES密钥
	h := md5.New()
	h.Write([]byte(key))
	if _, err := h.Write([]byte(key)); err != nil {
		return "", ErrMD5WriteFailed
	}
	akey := hex.EncodeToString(h.Sum(nil))
	// Base64解码
	enDataFromBase64, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", err
	}
	// 3. 验证密文长度
	if len(enDataFromBase64) == 0 || len(enDataFromBase64)%aes.BlockSize != 0 {
		return "", ErrInvalidInput
	}

	// 创建AES加密块
	block, err := aes.NewCipher([]byte(akey))
	if err != nil {
		return "", ErrAesCipher
	}
	// 使用密钥前16字节作为初始化向量(IV)
	iv := []byte(akey)[:aes.BlockSize]
	decrypter := cipher.NewCBCDecrypter(block, iv)
	// 执行解密
	dst := make([]byte, len(enDataFromBase64))
	decrypter.CryptBlocks(dst, enDataFromBase64)
	// PKCS7解填充处理
	length := len(dst)
	unpadding := int(dst[length-1])
	if unpadding < 1 || unpadding > aes.BlockSize {
		return "", ErrInvalidPadding
	}

	if length < unpadding {
		return "", ErrInvalidPadding
	}
	// 输出解密结果
	return string(dst[:(length - unpadding)]), nil
}

func AesEncryptGcmByKey(content string, key string) (string, error) {
	// Validate key length
	if len(key) != 32 {
		return "", ErrInvalidKey
	}

	// Generate a random IV (nonce) - 12 bytes recommended for GCM
	iv := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", ErrIVGeneration
	}

	// Create cipher block
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("aes cipher creation failed: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("gcm mode creation failed: %w", err)
	}

	// Encrypt and authenticate the data
	encrypted := gcm.Seal(nil, iv, []byte(content), nil)

	// Pre-allocate buffer with exact size needed
	bufferResult := make([]byte, len(iv)+len(encrypted))
	copy(bufferResult[:len(iv)], iv)
	copy(bufferResult[len(iv):], encrypted)

	return base64.StdEncoding.EncodeToString(bufferResult), nil
}
