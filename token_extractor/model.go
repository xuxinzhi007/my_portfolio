package token_extractor

import (
	"errors"
	"net/url"
	"time"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username  string
	Password  string
	TargetURL string
}

// Validate 验证登录请求
func (r *LoginRequest) Validate() error {
	if r.Username == "" {
		return errors.New("用户名不能为空")
	}
	if len(r.Username) > 100 {
		return errors.New("用户名长度不能超过100个字符")
	}

	if r.Password == "" {
		return errors.New("密码不能为空")
	}
	if len(r.Password) > 100 {
		return errors.New("密码长度不能超过100个字符")
	}

	if r.TargetURL == "" {
		return errors.New("目标URL不能为空")
	}

	// 验证URL格式
	parsedURL, err := url.Parse(r.TargetURL)
	if err != nil {
		return errors.New("无效的URL格式")
	}

	if parsedURL.Scheme != "https" {
		return errors.New("仅支持HTTPS协议")
	}

	return nil
}

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
	ID         string
	Timestamp  time.Time
	Username   string
	Success    bool
	KeyHeaders map[string]string // 仅保存关键头部
}

// IsKeyHeader 判断是否为关键头部
func IsKeyHeader(name string) bool {
	keyHeaders := []string{
		"X-Auth-Token",
		"X-Auth-Ts",
		"Gtoken",
		"Authorization",
		"x-auth-token",
		"x-auth-ts",
		"gtoken",
		"authorization",
	}

	for _, key := range keyHeaders {
		if name == key {
			return true
		}
	}

	return false
}
