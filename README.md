# License Management System

这是一个完整的许可证管理系统，使用Go语言开发，提供以下功能：

1. License文件生成
   - 支持用户信息嵌入
   - 有效期设置
   - 多级别授权支持
   - RSA加密签名

2. License文件验证
   - 签名验证
   - 有效期检查
   - 设备绑定验证
   - 授权级别验证

3. License续期功能
   - 自动续期支持
   - 旧License导入

4. 授权信息查询
   - 查看当前授权状态
   - 过期提醒

5. 安全特性
   - RSA非对称加密
   - 防篡改机制
   - 防重放攻击

## 项目结构

```
license_demo/
├── cmd/                    # 命令行工具
│   ├── generate/          # 生成密钥和许可证
│   ├── validate/          # 验证许可证
│   └── info/              # 查看许可证信息
├── internal/              # 内部包
│   ├── crypto/            # 加密相关
│   ├── license/           # 许可证处理
│   └── utils/             # 工具函数
├── pkg/                   # 公共包
│   └── models/            # 数据模型
└── configs/               # 配置文件
```

## 安装

确保已安装Go 1.21或更高版本，然后执行：

```bash
go mod download
```

## 使用方法

1. 首先生成密钥对：
```bash
go run cmd/main.go keys generate
```

2. 生成新的License：
```bash
go run cmd/main.go license generate --user "用户名" --device-id "设备ID" --level "专业版" --days 365
```

3. 验证License：
```bash
go run cmd/main.go license validate --file "license.dat"
```

4. 查看License信息：
```bash
go run cmd/main.go license info --file "license.dat"
```

## 编译

```bash
go build -o license-tool cmd/main.go
```

## 注意事项

- 请妥善保管私钥文件
- License文件不可跨设备使用
- 建议定期备份License文件
