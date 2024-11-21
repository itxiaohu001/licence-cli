package license

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"license-cli/internal/crypto"
	"license-cli/pkg/models"
)

// Manager 许可证管理器
type Manager struct {
	privateKeyPath string
	publicKeyPath  string
}

// NewManager 创建新的许可证管理器
func NewManager(privateKeyPath, publicKeyPath string) *Manager {
	return &Manager{
		privateKeyPath: privateKeyPath,
		publicKeyPath:  publicKeyPath,
	}
}

// GenerateKeys 生成新的密钥对
func (m *Manager) GenerateKeys() error {
	privateKey, publicKey, err := crypto.GenerateKeyPair(crypto.DefaultKeySize)
	if err != nil {
		return err
	}

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(m.privateKeyPath), 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(m.publicKeyPath), 0755); err != nil {
		return err
	}

	// 保存密钥
	if err := crypto.SavePrivateKey(privateKey, m.privateKeyPath); err != nil {
		return err
	}
	return crypto.SavePublicKey(publicKey, m.publicKeyPath)
}

// GenerateLicense 生成新的许可证
func (m *Manager) GenerateLicense(userName, deviceID string, level models.LicenseLevel, validDays int) (*models.License, error) {
	// 创建许可证
	license := models.NewLicense(userName, deviceID, level, validDays)

	// 加载私钥
	privateKey, err := crypto.LoadPrivateKey(m.privateKeyPath)
	if err != nil {
		return nil, err
	}

	// 序列化许可证数据（不包括签名字段）
	license.Signature = ""
	licenseData, err := json.Marshal(license)
	if err != nil {
		return nil, err
	}

	// 生成签名
	signature, err := crypto.Sign(privateKey, licenseData)
	if err != nil {
		return nil, err
	}

	// 设置签名
	license.Signature = base64.StdEncoding.EncodeToString(signature)
	return license, nil
}

// SaveLicense 保存许可证到文件
func (m *Manager) SaveLicense(license *models.License, filePath string) error {
	data, err := json.MarshalIndent(license, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// LoadLicense 从文件加载许可证
func (m *Manager) LoadLicense(filePath string) (*models.License, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var license models.License
	if err := json.Unmarshal(data, &license); err != nil {
		return nil, err
	}
	return &license, nil
}

// VerifyLicense 验证许可证
func (m *Manager) VerifyLicense(license *models.License, deviceID string) error {
	// 检查基本有效性
	if !license.IsValid(deviceID) {
		return errors.New("license is not valid for this device or has expired")
	}

	// 加载公钥
	publicKey, err := crypto.LoadPublicKey(m.publicKeyPath)
	if err != nil {
		return err
	}

	// 获取签名
	signature, err := base64.StdEncoding.DecodeString(license.Signature)
	if err != nil {
		return err
	}

	// 准备验证数据
	tempSignature := license.Signature
	license.Signature = ""
	licenseData, err := json.Marshal(license)
	if err != nil {
		return err
	}
	license.Signature = tempSignature

	// 验证签名
	return crypto.Verify(publicKey, licenseData, signature)
}

// RenewLicense 续期许可证
func (m *Manager) RenewLicense(license *models.License, days int) (*models.License, error) {
	// 更新过期时间
	license.Renew(days)

	// 重新生成签名
	privateKey, err := crypto.LoadPrivateKey(m.privateKeyPath)
	if err != nil {
		return nil, err
	}

	license.Signature = ""
	licenseData, err := json.Marshal(license)
	if err != nil {
		return nil, err
	}

	signature, err := crypto.Sign(privateKey, licenseData)
	if err != nil {
		return nil, err
	}

	license.Signature = base64.StdEncoding.EncodeToString(signature)
	return license, nil
}
