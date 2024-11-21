package main

import (
	"fmt"
	"os"
	"path/filepath"

	"license-cli/internal/license"
	"license-cli/pkg/models"

	"github.com/spf13/cobra"
)

var (
	// 全局标志
	configDir string
	userName  string
	deviceID  string
	level     string
	days      int
	filePath  string
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	configDir = filepath.Join(homeDir, ".license-demo")

	// 生成密钥命令
	keysCmd.Flags().StringVarP(&configDir, "config-dir", "c", configDir, "配置目录路径")
	rootCmd.AddCommand(keysCmd)

	// 生成许可证命令
	generateCmd.Flags().StringVarP(&userName, "user", "u", "", "用户名")
	generateCmd.Flags().StringVarP(&deviceID, "device-id", "d", "", "设备ID")
	generateCmd.Flags().StringVarP(&level, "level", "l", "basic", "授权级别 (basic/professional/enterprise)")
	generateCmd.Flags().IntVarP(&days, "days", "t", 365, "有效期天数")
	generateCmd.Flags().StringVarP(&filePath, "output", "o", "license.dat", "输出文件路径")
	rootCmd.AddCommand(generateCmd)

	// 验证许可证命令
	verifyCmd.Flags().StringVarP(&deviceID, "device-id", "d", "", "设备ID")
	verifyCmd.Flags().StringVarP(&filePath, "file", "f", "license.dat", "许可证文件路径")
	rootCmd.AddCommand(verifyCmd)

	// 查看许可证信息命令
	infoCmd.Flags().StringVarP(&filePath, "file", "f", "license.dat", "许可证文件路径")
	rootCmd.AddCommand(infoCmd)

	// 续期许可证命令
	renewCmd.Flags().StringVarP(&filePath, "file", "f", "license.dat", "许可证文件路径")
	renewCmd.Flags().IntVarP(&days, "days", "t", 365, "续期天数")
	rootCmd.AddCommand(renewCmd)
}

var rootCmd = &cobra.Command{
	Use:   "license-tool",
	Short: "许可证管理工具",
	Long:  `一个用于生成、验证和管理软件许可证的命令行工具。`,
}

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "生成RSA密钥对",
	Run: func(cmd *cobra.Command, args []string) {
		privateKeyPath := filepath.Join(configDir, "private.pem")
		publicKeyPath := filepath.Join(configDir, "public.pem")

		manager := license.NewManager(privateKeyPath, publicKeyPath)
		if err := manager.GenerateKeys(); err != nil {
			fmt.Printf("生成密钥对失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("密钥对已生成:\n私钥: %s\n公钥: %s\n", privateKeyPath, publicKeyPath)
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "生成新的许可证",
	Run: func(cmd *cobra.Command, args []string) {
		if userName == "" || deviceID == "" {
			fmt.Println("用户名和设备ID不能为空")
			os.Exit(1)
		}

		privateKeyPath := filepath.Join(configDir, "private.pem")
		publicKeyPath := filepath.Join(configDir, "public.pem")
		manager := license.NewManager(privateKeyPath, publicKeyPath)

		licenseLevel := models.LicenseLevel(level)
		lic, err := manager.GenerateLicense(userName, deviceID, licenseLevel, days)
		if err != nil {
			fmt.Printf("生成许可证失败: %v\n", err)
			os.Exit(1)
		}

		if err := manager.SaveLicense(lic, filePath); err != nil {
			fmt.Printf("保存许可证失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("许可证已生成: %s\n", filePath)
	},
}

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "验证许可证",
	Run: func(cmd *cobra.Command, args []string) {
		if deviceID == "" {
			fmt.Println("设备ID不能为空")
			os.Exit(1)
		}

		privateKeyPath := filepath.Join(configDir, "private.pem")
		publicKeyPath := filepath.Join(configDir, "public.pem")
		manager := license.NewManager(privateKeyPath, publicKeyPath)

		lic, err := manager.LoadLicense(filePath)
		if err != nil {
			fmt.Printf("加载许可证失败: %v\n", err)
			os.Exit(1)
		}

		if err := manager.VerifyLicense(lic, deviceID); err != nil {
			fmt.Printf("许可证无效: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("许可证有效")
		if days := lic.DaysUntilExpiration(); days > 0 {
			fmt.Printf("剩余有效期: %d 天\n", days)
		}
	},
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "显示许可证信息",
	Run: func(cmd *cobra.Command, args []string) {
		privateKeyPath := filepath.Join(configDir, "private.pem")
		publicKeyPath := filepath.Join(configDir, "public.pem")
		manager := license.NewManager(privateKeyPath, publicKeyPath)

		lic, err := manager.LoadLicense(filePath)
		if err != nil {
			fmt.Printf("加载许可证失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("许可证信息:\n")
		fmt.Printf("ID: %s\n", lic.ID)
		fmt.Printf("用户名: %s\n", lic.UserName)
		fmt.Printf("设备ID: %s\n", lic.DeviceID)
		fmt.Printf("授权级别: %s\n", lic.Level)
		fmt.Printf("颁发时间: %s\n", lic.IssuedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("过期时间: %s\n", lic.ExpiresAt.Format("2006-01-02 15:04:05"))
		if lic.IsExpired() {
			fmt.Println("状态: 已过期")
		} else {
			fmt.Printf("状态: 有效（剩余 %d 天）\n", lic.DaysUntilExpiration())
		}
	},
}

var renewCmd = &cobra.Command{
	Use:   "renew",
	Short: "续期许可证",
	Run: func(cmd *cobra.Command, args []string) {
		privateKeyPath := filepath.Join(configDir, "private.pem")
		publicKeyPath := filepath.Join(configDir, "public.pem")
		manager := license.NewManager(privateKeyPath, publicKeyPath)

		lic, err := manager.LoadLicense(filePath)
		if err != nil {
			fmt.Printf("加载许可证失败: %v\n", err)
			os.Exit(1)
		}

		lic, err = manager.RenewLicense(lic, days)
		if err != nil {
			fmt.Printf("续期许可证失败: %v\n", err)
			os.Exit(1)
		}

		if err := manager.SaveLicense(lic, filePath); err != nil {
			fmt.Printf("保存许可证失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("许可证已续期 %d 天\n", days)
		fmt.Printf("新的过期时间: %s\n", lic.ExpiresAt.Format("2006-01-02 15:04:05"))
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
