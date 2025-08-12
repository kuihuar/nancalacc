package cipherutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"nancalacc/internal/conf"
)

var (
	ErrEmptyPlaintext          = errors.New("empty plaintext")
	ErrInvalidPadding          = errors.New("invalid padding")
	ErrInvalidPaddingBlockSize = errors.New("invalid padding block size")
	ErrInvalidKey              = errors.New("invalid key")
	ErrIVGeneration            = errors.New("failed to generate initialization vector")
	ErrDecryptFailed           = errors.New("decrypt failed")
	ErrMD5WriteFailed          = errors.New("md5 write failed")
	ErrInvalidInput            = errors.New("invalid input")
	ErrAesCipher               = errors.New("aes cipher creation failed")

	ErrNonceGeneration   = errors.New("failed to generate nonce")
	ErrCipherCreation    = errors.New("failed to create AES cipher")
	ErrGCMCreation       = errors.New("failed to create GCM instance")
	ErrInvalidCiphertext = errors.New("ciphertext too short or invalid")
	ErrDecryptionFailed  = errors.New("decryption failed")
)

// GetAppUID 获取应用UID，优先从环境变量获取，否则使用默认值
func GetAppUID() string {
	uid := conf.GetEnvWithDefault("APP_UID", "nancalacc-426614174000")
	return uid
}

func DecryptByAes(content string, key string) (string, error) {
	fmt.Printf("DecryptByAes.content: %s\n", content)
	fmt.Printf("DecryptByAes.key: %s\n", key)

	if len(content) < 24 {
		return "", ErrEmptyPlaintext
	}
	if len(key) == 0 {
		return "", ErrInvalidKey
	}
	h := md5.New()
	if _, err := h.Write([]byte(key)); err != nil {
		return "", ErrMD5WriteFailed
	}
	akey := hex.EncodeToString(h.Sum(nil))
	enDataFromBase64, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", err
	}
	if len(enDataFromBase64) == 0 || len(enDataFromBase64)%aes.BlockSize != 0 {
		return "", ErrInvalidInput
	}

	block, err := aes.NewCipher([]byte(akey))
	if err != nil {
		return "", ErrAesCipher
	}
	iv := []byte(akey)[:aes.BlockSize]
	decrypter := cipher.NewCBCDecrypter(block, iv)
	dst := make([]byte, len(enDataFromBase64))
	decrypter.CryptBlocks(dst, enDataFromBase64)
	length := len(dst)
	unpadding := int(dst[length-1])
	if length < unpadding {
		return "", ErrInvalidPadding
	}
	res := string(dst[:(length - unpadding)])
	return res, nil
}

func EncryptByAes(content string, key string) (string, error) {
	fmt.Printf("EncryptByAes.content: %s\n", content)
	fmt.Printf("EncryptByAes.key: %s\n", key)

	if len(content) == 0 {
		return "", ErrEmptyPlaintext
	}
	if len(key) == 0 {
		return "", ErrInvalidKey
	}

	h := md5.New()
	if _, err := h.Write([]byte(key)); err != nil {
		return "", ErrMD5WriteFailed
	}
	akey := hex.EncodeToString(h.Sum(nil))

	block, err := aes.NewCipher([]byte(akey))
	if err != nil {
		return "", ErrAesCipher
	}

	iv := []byte(akey)[:aes.BlockSize]

	plaintext := []byte(content)
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	padtext := append(plaintext, bytes.Repeat([]byte{byte(padding)}, padding)...)

	ciphertext := make([]byte, len(padtext))
	encrypter := cipher.NewCBCEncrypter(block, iv)
	encrypter.CryptBlocks(ciphertext, padtext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Encrypt(plaintext string, key []byte) (string, error) {
	if len(key) != 32 {
		return "", ErrInvalidKey
	}
	if plaintext == "" {
		return "", ErrEmptyPlaintext
	}

	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return "", ErrNonceGeneration
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", ErrCipherCreation
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", ErrGCMCreation
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	encrypted := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func Decrypt(ciphertext string, key []byte) (string, error) {
	if len(key) != 32 {
		return "", ErrInvalidKey
	}

	encryptedData, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", ErrInvalidCiphertext
	}

	if len(encryptedData) < 12 {
		return "", ErrInvalidCiphertext
	}

	nonce := encryptedData[:12]
	data := encryptedData[12:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", ErrCipherCreation
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", ErrGCMCreation
	}

	plaintext, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	return string(plaintext), nil
}

func GenerateKey(uid, salt string) string {
	h := sha256.New()
	h.Write([]byte(uid + salt))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash[:32]
}

func EncryptValueWithEnvSalt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	salt, err := conf.GetEnv("ENCRYPTION_SALT")
	if err != nil {
		return "", fmt.Errorf("failed to get encryption salt: %w", err)
	}

	uid := GetAppUID()
	envKey := GenerateKey(uid, salt)
	if len(envKey) != 32 {
		return "", fmt.Errorf("generated key must be 32 bytes, got %d", len(envKey))
	}

	encrypted, err := Encrypt(plaintext, []byte(envKey))
	if err != nil {
		return "", fmt.Errorf("encryption failed: %w", err)
	}
	return encrypted, nil
}

func DecryptValueWithEnvSalt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	salt, err := conf.GetEnv("ENCRYPTION_SALT")
	if err != nil {
		return "", fmt.Errorf("failed to get encryption salt: %w", err)
	}

	uid := GetAppUID()
	envKey := GenerateKey(uid, salt)
	if len(envKey) != 32 {
		return "", fmt.Errorf("generated key must be 32 bytes, got %d", len(envKey))
	}

	decrypted, err := Decrypt(ciphertext, []byte(envKey))
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}
	return decrypted, nil
}
