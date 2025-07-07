package cipherutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func AesEncryptGCMByKey(content string, key []byte) (string, error) {
	// 创建 AES GCM 加密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 生成随机初始向量
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密内容
	ciphertext := gcm.Seal(nonce, nonce, []byte(content), nil)

	// 返回 base64 编码的结果
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
