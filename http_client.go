package cweskills

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// DefaultBaseURL 是MITRE CWE REST API的默认基础URL
const DefaultBaseURL = "https://cwe-api.mitre.org/api/v1"

// DefaultTimeout 是HTTP请求的默认超时时间
const DefaultTimeout = 30 * time.Second

// DefaultUserAgent 是HTTP请求的默认User-Agent
const DefaultUserAgent = "cwe-sdk-go/" + Version

// HTTPClient 提供带有重试、速率限制和超时控制的HTTP客户端。
//
// 该客户端封装了标准的http.Client，添加了以下功能：
//   - 自动重试：对5xx错误自动重试
//   - 速率限制：通过RateLimiter控制请求频率
//   - 超时控制：通过context控制请求超时
//   - User-Agent：自动添加SDK版本标识
//
// 示例：
//
//	client := cwe.NewHTTPClient("https://cwe-api.mitre.org/api/v1",
//	    cwe.WithRetry(3, 5*time.Second),
//	    cwe.WithHTTPRateLimiter(1.0, 5),
//	    cwe.WithHTTPTimeout(30*time.Second),
//	)
//	defer client.Close()
type HTTPClient struct {
	client      *http.Client
	baseURL     string
	userAgent   string
	maxRetries  int
	retryDelay  time.Duration
	rateLimiter *RateLimiter
}

// HTTPClientOption 是HTTPClient的配置选项函数类型
type HTTPClientOption func(*HTTPClient)

// WithRetry 设置HTTP请求的重试次数和重试间隔。
//
// 参数：
//   - maxRetries: 最大重试次数（不包括首次请求），0表示不重试
//   - delay: 重试间隔时间
func WithRetry(maxRetries int, delay time.Duration) HTTPClientOption {
	return func(c *HTTPClient) {
		c.maxRetries = maxRetries
		c.retryDelay = delay
	}
}

// WithHTTPRateLimiter 设置HTTP请求的速率限制器。
//
// 参数：
//   - rate: 每秒允许的请求数
//   - burst: 允许的突发请求数
func WithHTTPRateLimiter(rate float64, burst int) HTTPClientOption {
	return func(c *HTTPClient) {
		c.rateLimiter = NewRateLimiter(rate, burst)
	}
}

// WithHTTPTimeout 设置HTTP请求的超时时间。
//
// 参数：
//   - timeout: 请求超时时间
func WithHTTPTimeout(timeout time.Duration) HTTPClientOption {
	return func(c *HTTPClient) {
		c.client.Timeout = timeout
	}
}

// WithUserAgent 设置HTTP请求的User-Agent头。
//
// 参数：
//   - ua: User-Agent字符串
func WithUserAgent(ua string) HTTPClientOption {
	return func(c *HTTPClient) {
		c.userAgent = ua
	}
}

// WithHTTPClient 设置自定义的http.Client。
//
// 参数：
//   - client: 自定义的http.Client实例
func WithHTTPClient(client *http.Client) HTTPClientOption {
	return func(c *HTTPClient) {
		c.client = client
	}
}

// NewHTTPClient 创建一个新的HTTPClient实例。
//
// 参数：
//   - baseURL: API基础URL
//   - opts: 可选的配置选项
//
// 返回值：
//   - *HTTPClient: 新创建的HTTPClient实例
func NewHTTPClient(baseURL string, opts ...HTTPClientOption) *HTTPClient {
	client := &HTTPClient{
		client:     &http.Client{Timeout: DefaultTimeout},
		baseURL:    baseURL,
		userAgent:  DefaultUserAgent,
		maxRetries: 0,
		retryDelay: 1 * time.Second,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Get 发送GET请求并解析JSON响应。
//
// 该方法自动处理速率限制、重试和错误响应。
// 响应体会被自动解析到result参数指向的结构体中。
//
// 参数：
//   - ctx: 请求上下文，用于取消和超时控制
//   - path: 请求路径（相对于baseURL）
//   - result: 用于存储解析结果的指针
//
// 返回值：
//   - error: 请求失败时返回APIError或其他错误
func (c *HTTPClient) Get(ctx context.Context, path string, result interface{}) error {
	respBody, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return NewParseError(fmt.Sprintf("JSON解析失败: %v", err), 0)
		}
	}

	return nil
}

// GetRaw 发送GET请求并返回原始响应体。
//
// 参数：
//   - ctx: 请求上下文
//   - path: 请求路径
//
// 返回值：
//   - []byte: 响应体
//   - error: 请求失败时返回错误
func (c *HTTPClient) GetRaw(ctx context.Context, path string) ([]byte, error) {
	return c.doRequest(ctx, http.MethodGet, path, nil)
}

// Post 发送POST请求并解析JSON响应。
//
// 参数：
//   - ctx: 请求上下文
//   - path: 请求路径
//   - body: 请求体（会被序列化为JSON）
//   - result: 用于存储解析结果的指针
//
// 返回值：
//   - error: 请求失败时返回错误
func (c *HTTPClient) Post(ctx context.Context, path string, body interface{}, result interface{}) error {
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("序列化请求体失败: %w", err)
		}
	}

	respBody, err := c.doRequest(ctx, http.MethodPost, path, bodyBytes)
	if err != nil {
		return err
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return NewParseError(fmt.Sprintf("JSON解析失败: %v", err), 0)
		}
	}

	return nil
}

// PostForm 发送POST表单请求。
//
// 参数：
//   - ctx: 请求上下文
//   - path: 请求路径
//   - data: 表单数据
//   - result: 用于存储解析结果的指针
//
// 返回值：
//   - error: 请求失败时返回错误
func (c *HTTPClient) PostForm(ctx context.Context, path string, data url.Values, result interface{}) error {
	reqURL := c.buildURL(path)

	// 速率限制
	if c.rateLimiter != nil {
		if err := c.rateLimiter.Wait(ctx); err != nil {
			return err
		}
	}

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(c.retryDelay):
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBufferString(data.Encode()))
		if err != nil {
			return fmt.Errorf("创建请求失败: %w", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("User-Agent", c.userAgent)

		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode >= 500 && attempt < c.maxRetries {
			lastErr = NewAPIError(resp.StatusCode, reqURL, http.MethodPost)
			continue
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return NewAPIError(resp.StatusCode, reqURL, http.MethodPost)
		}

		if result != nil {
			if err := json.Unmarshal(respBody, result); err != nil {
				return NewParseError(fmt.Sprintf("JSON解析失败: %v", err), 0)
			}
		}

		return nil
	}

	return lastErr
}

// doRequest 执行HTTP请求，处理重试和速率限制。
func (c *HTTPClient) doRequest(ctx context.Context, method, path string, body []byte) ([]byte, error) {
	reqURL := c.buildURL(path)

	// 速率限制
	if c.rateLimiter != nil {
		if err := c.rateLimiter.Wait(ctx); err != nil {
			return nil, err
		}
	}

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.retryDelay):
			}
		}

		var bodyReader io.Reader
		if body != nil {
			bodyReader = bytes.NewReader(body)
		}

		req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("创建请求失败: %w", err)
		}

		req.Header.Set("User-Agent", c.userAgent)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}

		// 5xx错误自动重试
		if resp.StatusCode >= 500 && attempt < c.maxRetries {
			lastErr = NewAPIError(resp.StatusCode, reqURL, method)
			continue
		}

		// 非2xx状态码返回错误
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, NewAPIError(resp.StatusCode, reqURL, method)
		}

		return respBody, nil
	}

	return nil, lastErr
}

// buildURL 构建完整的请求URL。
func (c *HTTPClient) buildURL(path string) string {
	return c.baseURL + path
}

// Close 关闭HTTP客户端，释放资源。
func (c *HTTPClient) Close() {
	c.client.CloseIdleConnections()
}

// GetBaseURL 获取基础URL。
func (c *HTTPClient) GetBaseURL() string {
	return c.baseURL
}

// SetBaseURL 设置基础URL。
func (c *HTTPClient) SetBaseURL(url string) {
	c.baseURL = url
}

// GetMaxRetries 获取最大重试次数。
func (c *HTTPClient) GetMaxRetries() int {
	return c.maxRetries
}

// SetMaxRetries 设置最大重试次数。
func (c *HTTPClient) SetMaxRetries(maxRetries int) {
	c.maxRetries = maxRetries
}

// GetRetryDelay 获取重试间隔时间。
func (c *HTTPClient) GetRetryDelay() time.Duration {
	return c.retryDelay
}

// SetRetryDelay 设置重试间隔时间。
func (c *HTTPClient) SetRetryDelay(delay time.Duration) {
	c.retryDelay = delay
}

// GetRateLimiter 获取速率限制器。
func (c *HTTPClient) GetRateLimiter() *RateLimiter {
	return c.rateLimiter
}

// SetRateLimiter 设置速率限制器。
func (c *HTTPClient) SetRateLimiter(limiter *RateLimiter) {
	c.rateLimiter = limiter
}

// GetHTTPClient 获取底层的http.Client。
func (c *HTTPClient) GetHTTPClient() *http.Client {
	return c.client
}

// SetHTTPClient 设置底层的http.Client。
func (c *HTTPClient) SetHTTPClient(client *http.Client) {
	c.client = client
}
