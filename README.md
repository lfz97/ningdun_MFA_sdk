# ND_Dkey

一个用于在宁盾(宁盾统一身份)平台租户环境中：
1. 登录后台
2. 设置当前租户上下文
3. 同步外部数据源 (External Identity Store)
4. 通过邮箱发送一次性令牌 / MFA 邮件

封装底层接口调用（使用 `resty`），提供更简洁的 Go 调用方式。主要适合脚本自动化、批处理、或将宁盾身份相关动作集成到你自己的服务中。

## ✨ 特性
- 简单初始化：一行代码完成登录与租户上下文设置
- 数据源批量同步：支持多数据源循环同步，失败即返回错误
- 用户邮箱检索 + 令牌发送：链路封装，减少重复样板代码
- 轻量无额外依赖：除 `resty` 与标准库外无其他第三方依赖
- 可扩展：保留客户端指针 `client_ptr` 以便后续继续封装其它 API

## 📦 安装

将本目录放入你的 Go module 中（例如 `MyPackages/ND_Dkey`），在代码中：
```go
import "your-module-path/MyPackages/ND_Dkey"
```

目前未发布到公共仓库（如 GitHub / Go proxy），需要通过本地或私有仓库引用。

## 🚀 快速开始
1. 修改配置文件 `ND_DkeyConf.go` 中的占位符：
   - `NDhost`、`adminLoginName`、`adminPassword`、`tenantId` 等
   - `sync_tenant_config` 中的 `dataSource1` / `dataSource2` 为实际数据源 ID
2. 在代码中初始化：
```go
cli, err := ND_Dkey.NDInit()
if err != nil {
	panic(err)
}

// 同步数据源
if err = cli.DatasourceSync(); err != nil {
	panic(err)
}

// 发送 MFA 邮件令牌（示例：有效期 5 天）
if err = cli.SendMFA("user@example.com", "5"); err != nil {
	panic(err)
}
```

## ⚙️ 配置说明 (`ND_DkeyConf.go`)
```go
var sync_tenant_config = map[string]string{
	"dataSource1": "<dataSource1>",
	"dataSource2": "<dataSource2>",
}

const (
	NDhost            = "<http://NDhost:port>"
	adminLoginName    = "<adminLoginName>"
	adminPassword     = "<adminPassword>"
	tenantName        = "<tenantName>"          // 当前租户名（如需要）
	identityStoreName = "<identityStoreName>"   // 外部身份源名称（如需要）
	tenantId          = "<tenantId>"            // 关键：外部用户页面触发接口获取
)
```
获取 `tenantId` 与数据源 ID 的方式：在宁盾后台点击 “外部用户” 页面，触发接口：
```
/am/controller/tenant/api/2/identity/externalIdentityStore/list?tenant=<tenant>
```

## 🧩 API 概览
### 1. 初始化
```go
func NDInit() (*Client, error)
```
动作：创建 Resty 客户端 -> 登录 -> 设置当前租户。失败时返回错误。

### 2. 同步数据源
```go
func (c *Client) DatasourceSync() error
```
遍历 `sync_tenant_config` 的各个数据源，依次调用同步接口。任一失败即返回。

### 3. 发送 MFA 邮件令牌
```go
func (c *Client) SendMFA(mail string, expireDays string) error
```
流程：
1. 搜索用户（邮箱字段）
2. 取第一个匹配用户 ID 与邮箱地址
3. 调用令牌投递接口，失败时返回详细响应内容。

## 📘 使用示例（完整）
```go
package main

import (
	"fmt"
	"your-module-path/MyPackages/ND_Dkey"
)

func main() {
	cli, err := ND_Dkey.NDInit()
	if err != nil {
		panic(err)
	}

	if err = cli.DatasourceSync(); err != nil {
		panic(err)
	}

	if err = cli.SendMFA("user@example.com", "5"); err != nil {
		panic(err)
	}

	fmt.Println("Done.")
}
```

