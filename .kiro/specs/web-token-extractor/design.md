# Design Document - Web Token Extractor

## Overview

Web Token Extractor是一个集成到现有Fyne应用中的功能模块，用于自动化登录Anker Solix专业版网站并提取HTTP请求头信息（特别是认证token）。该功能将作为一个新的Tab页面添加到主应用中。

### 技术栈
- **UI框架**: Fyne v2.7.0
- **HTTP客户端**: Go标准库 net/http + chromedp（用于浏览器自动化）
- **存储**: 可选的本地JSON存储（用于保存历史记录）
- **语言**: Go 1.24

### 核心功能
1. 用户输入账号密码
2. 使用headless Chrome自动登录目标网站
3. 捕获认证后的HTTP请求头
4. 展示关键头部信息（特别是token相关）
5. 提供复制功能

## Architecture

### 系统架构图

```
┌─────────────────────────────────────────────────────────┐
│                    Main Application                      │
│                     (main.go)                           │
└────────────────────┬────────────────────────────────────┘
                     │
                     │ 添加新Tab
                     ▼
┌─────────────────────────────────────────────────────────┐
│              Token Extractor Module                      │
│            (token_extractor package)                     │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │   UI Layer   │  │ Service Layer│  │ Storage Layer│ │
│  │   (ui.go)    │  │ (extractor.go│  │ (storage.go) │ │
│  │              │  │               │  │              │ │
│  │ - 输入表单    │  │ - 浏览器控制  │  │ - 历史记录    │ │
│  │ - 结果展示    │  │ - 请求拦截    │  │ - 配置存储    │ │
│  │ - 复制功能    │  │ - 头部提取    │  │              │ │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘ │
│         │                 │                 │          │
│         └────────┬────────┘                 │          │
│                  │                          │          │
│                  ▼                          ▼          │
│         ┌─────────────────┐       ┌─────────────────┐ │
│         │   Model Layer   │       │   JSON Files    │ │
│         │   (model.go)    │       │                 │ │
│         │                 │       │ - history.json  │ │
│         │ - LoginRequest  │       │                 │ │
│         │ - HeaderInfo    │       │                 │ │
│         │ - ExtractResult │       │                 │ │
│         └─────────────────┘       └─────────────────┘ │
└─────────────────────────────────────────────────────────┘
                     │
                     │ 使用chromedp
                     ▼
┌─────────────────────────────────────────────────────────┐
│              Headless Chrome Browser                     │
│                                                          │
│  - 自动化登录                                             │
│  - 网络请求拦截                                           │
│  - HTTP头部捕获                                          │
└─────────────────────────────────────────────────────────┘
```

### 模块职责

#### 1. UI Layer (ui.go)
- 创建用户界面组件
- 处理用户输入
- 显示提取结果
- 提供复制到剪贴板功能
- 显示加载状态和错误信息

#### 2. Service Layer (extractor.go)
- 初始化headless浏览器
- 执行自动登录流程
- 拦截和捕获HTTP请求
- 提取关键请求头信息
- 错误处理和重试逻辑

#### 3. Model Layer (model.go)
- 定义数据结构
- 数据验证逻辑
- 业务规则

#### 4. Storage Layer (storage.go)
- 可选：保存提取历史
- 可选：保存用户配置（不包括密码）

## Components and Interfaces

### 1. 数据模型 (model.go)

```go
package token_extractor

import "time"

// LoginRequest 登录请求
type LoginRequest struct {
    Username string
    Password string
    TargetURL string
}

// Validate 验证登录请求
func (r *LoginRequest) Validate() error

// HeaderInfo HTTP头部信息
type HeaderInfo struct {
    Name  string
    Value string
    IsKey bool // 是否为关键头部（如token）
}

// ExtractResult 提取结果
type ExtractResult struct {
    Success   bool
    Timestamp time.Time
    Headers   []HeaderInfo
    Error     string
}

// HistoryRecord 历史记录
type HistoryRecord struct {
    ID        string
    Timestamp time.Time
    Username  string
    Success   bool
    KeyHeaders map[string]string // 仅保存关键头部
}
```

### 2. 提取服务接口 (extractor.go)

```go
package token_extractor

import "context"

// Extractor token提取器接口
type Extractor interface {
    // Extract 执行提取操作
    Extract(ctx context.Context, req LoginRequest) (*ExtractResult, error)
    
    // Close 关闭资源
    Close() error
}

// ChromeExtractor 基于Chrome的实现
type ChromeExtractor struct {
    // chromedp相关字段
}

// NewChromeExtractor 创建新的提取器
func NewChromeExtractor() (*ChromeExtractor, error)

// Extract 实现提取逻辑
func (e *ChromeExtractor) Extract(ctx context.Context, req LoginRequest) (*ExtractResult, error)

// Close 清理资源
func (e *ChromeExtractor) Close() error
```

### 3. UI组件 (ui.go)

```go
package token_extractor

import "fyne.io/fyne/v2"

// TokenExtractorUI token提取器UI
type TokenExtractorUI struct {
    window    fyne.Window
    extractor Extractor
    storage   Storage
    
    // UI组件
    usernameEntry *widget.Entry
    passwordEntry *widget.Entry
    extractButton *widget.Button
    statusLabel   *widget.Label
    resultList    *widget.List
    progressBar   *widget.ProgressBarInfinite
}

// NewTokenExtractorUI 创建UI实例
func NewTokenExtractorUI(window fyne.Window) *TokenExtractorUI

// MakeUI 构建UI界面
func (ui *TokenExtractorUI) MakeUI() fyne.CanvasObject

// handleExtract 处理提取操作
func (ui *TokenExtractorUI) handleExtract()

// displayResult 显示提取结果
func (ui *TokenExtractorUI) displayResult(result *ExtractResult)

// copyToClipboard 复制到剪贴板
func (ui *TokenExtractorUI) copyToClipboard(text string)
```

### 4. 存储接口 (storage.go)

```go
package token_extractor

// Storage 存储接口
type Storage interface {
    // SaveHistory 保存历史记录
    SaveHistory(record HistoryRecord) error
    
    // GetHistory 获取历史记录
    GetHistory(limit int) ([]HistoryRecord, error)
    
    // ClearHistory 清空历史
    ClearHistory() error
}

// JSONStorage JSON文件存储实现
type JSONStorage struct {
    filePath string
}

// NewJSONStorage 创建JSON存储
func NewJSONStorage(filePath string) *JSONStorage
```

## Data Models

### LoginRequest
```go
type LoginRequest struct {
    Username  string `json:"username"`
    Password  string `json:"password"`
    TargetURL string `json:"target_url"`
}
```

**验证规则**:
- Username: 非空，长度1-100
- Password: 非空，长度1-100
- TargetURL: 必须是有效的HTTPS URL

### HeaderInfo
```go
type HeaderInfo struct {
    Name  string `json:"name"`
    Value string `json:"value"`
    IsKey bool   `json:"is_key"`
}
```

**关键头部列表**:
- X-Auth-Token
- X-Auth-Ts
- Gtoken
- Authorization
- Cookie (如果包含token)

### ExtractResult
```go
type ExtractResult struct {
    Success   bool         `json:"success"`
    Timestamp time.Time    `json:"timestamp"`
    Headers   []HeaderInfo `json:"headers"`
    Error     string       `json:"error,omitempty"`
}
```

## Error Handling

### 错误类型

```go
var (
    ErrInvalidCredentials = errors.New("无效的登录凭证")
    ErrLoginFailed       = errors.New("登录失败")
    ErrNetworkError      = errors.New("网络连接错误")
    ErrBrowserError      = errors.New("浏览器初始化失败")
    ErrTimeout           = errors.New("操作超时")
    ErrNoHeaders         = errors.New("未能捕获到请求头")
)
```

### 错误处理策略

1. **输入验证错误**: 立即返回，显示具体错误信息
2. **网络错误**: 提供重试选项，最多重试2次
3. **浏览器错误**: 显示详细错误，建议检查系统环境
4. **超时错误**: 默认超时60秒，可配置

### 用户反馈

- 使用Fyne的dialog.ShowError显示错误
- 在状态栏显示当前操作状态
- 使用进度条显示长时间操作

## Testing Strategy

### 单元测试

1. **Model层测试**
   - 测试数据验证逻辑
   - 测试数据序列化/反序列化

2. **Storage层测试**
   - 测试历史记录保存和读取
   - 测试文件操作错误处理

### 集成测试

1. **Extractor测试**
   - 使用mock服务器测试登录流程
   - 测试请求头捕获逻辑
   - 测试错误场景

### 手动测试

1. **UI测试**
   - 测试所有用户交互
   - 测试复制功能
   - 测试不同屏幕尺寸

2. **端到端测试**
   - 使用测试账号进行完整流程测试
   - 验证提取的token有效性

## Implementation Details

### 浏览器自动化流程

1. **初始化浏览器**
```go
// 创建chromedp上下文
ctx, cancel := chromedp.NewContext(context.Background())
defer cancel()

// 设置超时
ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
defer cancel()
```

2. **登录流程**
```go
// 导航到登录页面
chromedp.Navigate(targetURL)

// 等待登录表单加载
chromedp.WaitVisible(`input[name="username"]`)

// 填写表单
chromedp.SendKeys(`input[name="username"]`, username)
chromedp.SendKeys(`input[name="password"]`, password)

// 提交表单
chromedp.Click(`button[type="submit"]`)

// 等待登录成功（根据实际页面调整）
chromedp.WaitVisible(`某个登录后才有的元素`)
```

3. **捕获请求头**
```go
// 使用chromedp的网络事件监听
chromedp.ListenTarget(ctx, func(ev interface{}) {
    if ev, ok := ev.(*network.EventRequestWillBeSent); ok {
        // 捕获请求头
        headers := ev.Request.Headers
        // 存储关键头部
    }
})
```

### UI布局设计

```
┌──────────────────────────────────────────────────────────────────┐
│                    🔐 网页Token提取器                             │
├──────────────────────────────────────────────────────────────────┤
│  URL: [https://ankersolix-professional-ci.anker.com/...       ] │
│  账号: [账号输入框]  密码: [密码输入框]  [开始提取]              │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│  状态: 等待操作...                                               │
├──────────────────────────────────────────────────────────────────┤
│  📊 提取结果:                                                    │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │ ⭐ X-Auth-Token  e27e7e6d6b7e720d0a23c54...      [复制]    │ │
│  │ ⭐ X-Auth-Ts     1762942176                      [复制]    │ │
│  │ ⭐ Gtoken        c82231545079c65fe4a28...        [复制]    │ │
│  │ 📋 User-Agent    Mozilla/5.0...                  [复制]    │ │
│  │ � Con tent-Type  application/json                [复制]    │ │
│  └────────────────────────────────────────────────────────────┘ │
│  [复制所有关键Token]  [清空结果]                                 │
└──────────────────────────────────────────────────────────────────┘
```

### 依赖项

需要添加到go.mod:
```
github.com/chromedp/chromedp v0.9.3
```

### 安全考虑

1. **密码处理**
   - 密码仅在内存中临时存储
   - 不保存到任何文件
   - 使用完立即清除

2. **Token存储**
   - 历史记录中的token应该脱敏（只显示前后几位）
   - 提供清空历史功能
   - 文件权限设置为仅当前用户可读

3. **网络安全**
   - 仅支持HTTPS连接
   - 验证SSL证书
   - 不信任自签名证书（除非用户明确允许）

## Performance Considerations

1. **浏览器启动优化**
   - 使用headless模式减少资源消耗
   - 禁用不必要的浏览器功能（图片、CSS等）
   - 复用浏览器实例（如果需要多次提取）

2. **超时设置**
   - 页面加载超时: 30秒
   - 登录操作超时: 60秒
   - 总体操作超时: 90秒

3. **资源清理**
   - 操作完成后立即关闭浏览器
   - 清理临时文件
   - 释放内存

## Future Enhancements

1. **多站点支持**: 支持配置多个不同的目标网站
2. **自动刷新**: 定期自动刷新token
3. **导出功能**: 导出为JSON、环境变量格式等
4. **代理支持**: 支持通过代理访问
5. **自定义选择器**: 允许用户自定义登录表单选择器
