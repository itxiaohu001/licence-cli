package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

const (
	// 默认RSA密钥长度
	DefaultKeySize = 2048
)

// GenerateKeyPair 生成RSA密钥对
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// SavePrivateKey 保存私钥到文件
func SavePrivateKey(key *rsa.PrivateKey, filename string) error {
	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return pem.Encode(file, block)
}

// SavePublicKey 保存公钥到文件
func SavePublicKey(key *rsa.PublicKey, filename string) error {
	keyBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return err
	}
	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: keyBytes,
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return pem.Encode(file, block)
}

// LoadPrivateKey 从文件加载私钥
func LoadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// LoadPublicKey 从文件加载公钥
func LoadPublicKey(filename string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("invalid public key type")
	}
}

// Sign 使用私钥对数据进行签名
func Sign(privateKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	hashed := sha256.Sum256(data)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
}

// Verify 使用公钥验证签名
func Verify(publicKey *rsa.PublicKey, data []byte, signature []byte) error {
	hashed := sha256.Sum256(data)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
}
