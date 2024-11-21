package models

import (
	"time"

	"github.com/google/uuid"
)

// LicenseLevel 定义授权级别
type LicenseLevel string

const (
	LicenseBasic       LicenseLevel = "basic"
	LicenseProfessional LicenseLevel = "professional"
	LicenseEnterprise  LicenseLevel = "enterprise"
)

// License 定义许可证结构
type License struct {
	ID            string       `json:"id" yaml:"id"`                         // 许可证唯一标识
	UserName      string       `json:"userName" yaml:"userName"`             // 用户名
	DeviceID      string       `json:"deviceId" yaml:"deviceId"`            // 设备ID
	Level         LicenseLevel `json:"level" yaml:"level"`                  // 授权级别
	IssuedAt      time.Time    `json:"issuedAt" yaml:"issuedAt"`           // 颁发时间
	ExpiresAt     time.Time    `json:"expiresAt" yaml:"expiresAt"`         // 过期时间
	Features      []string     `json:"features" yaml:"features"`            // 功能列表
	Signature     string       `json:"signature" yaml:"signature"`          // 数字签名
	IssuerID      string       `json:"issuerId" yaml:"issuerId"`           // 颁发者ID
	Version       string       `json:"version" yaml:"version"`              // 许可证版本
}

// NewLicense 创建新的许可证
func NewLicense(userName, deviceID string, level LicenseLevel, validDays int) *License {
	now := time.Now()
	return &License{
		ID:        uuid.New().String(),
		UserName:  userName,
		DeviceID:  deviceID,
		Level:     level,
		IssuedAt:  now,
		ExpiresAt: now.AddDate(0, 0, validDays),
		Version:   "1.0",
	}
}

// IsExpired 检查许可证是否过期
func (l *License) IsExpired() bool {
	return time.Now().After(l.ExpiresAt)
}

// IsValid 检查许可证是否有效
func (l *License) IsValid(deviceID string) bool {
	if l.IsExpired() {
		return false
	}
	if l.DeviceID != deviceID {
		return false
	}
	return true
}

// DaysUntilExpiration 获取剩余有效天数
func (l *License) DaysUntilExpiration() int {
	if l.IsExpired() {
		return 0
	}
	return int(l.ExpiresAt.Sub(time.Now()).Hours() / 24)
}

// Renew 续期许可证
func (l *License) Renew(days int) {
	if l.IsExpired() {
		l.ExpiresAt = time.Now().AddDate(0, 0, days)
	} else {
		l.ExpiresAt = l.ExpiresAt.AddDate(0, 0, days)
	}
}
