package token_extractor

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

var (
	ErrInvalidCredentials = errors.New("无效的登录凭证")
	ErrLoginFailed        = errors.New("登录失败")
	ErrNetworkError       = errors.New("网络连接错误")
	ErrBrowserError       = errors.New("浏览器初始化失败")
	ErrTimeout            = errors.New("操作超时")
	ErrNoHeaders          = errors.New("未能捕获到请求头")
)

// Extractor token提取器接口
type Extractor interface {
	// Extract 执行提取操作
	Extract(ctx context.Context, req LoginRequest) (*ExtractResult, error)

	// Close 关闭资源
	Close() error
}

// ChromeExtractor 基于Chrome的实现
type ChromeExtractor struct {
	allocCtx   context.Context
	allocCancel context.CancelFunc
}

// NewChromeExtractor 创建新的提取器
func NewChromeExtractor() (*ChromeExtractor, error) {
	// 创建浏览器上下文
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-images", true),
		chromedp.Flag("disable-javascript", false),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)

	return &ChromeExtractor{
		allocCtx:    allocCtx,
		allocCancel: allocCancel,
	}, nil
}

// Extract 实现提取逻辑
func (e *ChromeExtractor) Extract(ctx context.Context, req LoginRequest) (*ExtractResult, error) {
	// 验证请求
	if err := req.Validate(); err != nil {
		return &ExtractResult{
			Success:   false,
			Timestamp: time.Now(),
			Error:     err.Error(),
		}, err
	}

	// 创建浏览器上下文
	browserCtx, cancel := chromedp.NewContext(e.allocCtx)
	defer cancel()

	// 设置超时
	timeoutCtx, timeoutCancel := context.WithTimeout(browserCtx, 90*time.Second)
	defer timeoutCancel()

	// 存储捕获的请求头
	capturedHeaders := make(map[string]string)
	var headersMutex = make(chan struct{}, 1)
	headersMutex <- struct{}{}

	// 监听网络请求
	chromedp.ListenTarget(timeoutCtx, func(ev interface{}) {
		if ev, ok := ev.(*network.EventRequestWillBeSent); ok {
			// 只捕获目标域名的请求
			if strings.Contains(ev.Request.URL, "ankersolix-professional-ci.anker.com") {
				<-headersMutex
				for name, value := range ev.Request.Headers {
					if strValue, ok := value.(string); ok {
						capturedHeaders[name] = strValue
					}
				}
				headersMutex <- struct{}{}
			}
		}
	})

	// 执行登录流程
	err := chromedp.Run(timeoutCtx,
		network.Enable(),
		chromedp.Navigate(req.TargetURL),
		chromedp.Sleep(2*time.Second),

		// 等待登录表单加载
		chromedp.WaitVisible(`input[type="text"], input[type="email"], input[name="username"], input[placeholder*="账号"], input[placeholder*="用户名"], input[placeholder*="邮箱"]`, chromedp.ByQuery),

		// 填写用户名
		chromedp.SendKeys(`input[type="text"], input[type="email"], input[name="username"], input[placeholder*="账号"], input[placeholder*="用户名"], input[placeholder*="邮箱"]`, req.Username, chromedp.ByQuery),

		// 填写密码
		chromedp.SendKeys(`input[type="password"]`, req.Password, chromedp.ByQuery),

		chromedp.Sleep(1*time.Second),

		// 点击登录按钮
		chromedp.Click(`button[type="submit"], button:contains("登录"), button:contains("Login"), .login-button`, chromedp.ByQuery),

		// 等待登录完成（等待页面跳转或特定元素出现）
		chromedp.Sleep(5*time.Second),
	)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return &ExtractResult{
				Success:   false,
				Timestamp: time.Now(),
				Error:     ErrTimeout.Error(),
			}, ErrTimeout
		}

		return &ExtractResult{
			Success:   false,
			Timestamp: time.Now(),
			Error:     fmt.Sprintf("登录失败: %v", err),
		}, ErrLoginFailed
	}

	// 等待一下确保请求被捕获
	time.Sleep(2 * time.Second)

	// 检查是否捕获到头部
	<-headersMutex
	if len(capturedHeaders) == 0 {
		headersMutex <- struct{}{}
		return &ExtractResult{
			Success:   false,
			Timestamp: time.Now(),
			Error:     ErrNoHeaders.Error(),
		}, ErrNoHeaders
	}
	headersMutex <- struct{}{}

	// 构建结果
	var headers []HeaderInfo
	for name, value := range capturedHeaders {
		headers = append(headers, HeaderInfo{
			Name:  name,
			Value: value,
			IsKey: IsKeyHeader(name),
		})
	}

	return &ExtractResult{
		Success:   true,
		Timestamp: time.Now(),
		Headers:   headers,
	}, nil
}

// Close 清理资源
func (e *ChromeExtractor) Close() error {
	if e.allocCancel != nil {
		e.allocCancel()
	}
	return nil
}
