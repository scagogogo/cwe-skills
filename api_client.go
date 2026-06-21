package cweskills

import "time"

// APIClient 是MITRE CWE REST API的客户端。
//
// 该客户端封装了HTTP请求逻辑，提供了对CWE REST API的完整访问能力。
// 支持自定义基础URL、超时时间、速率限制等配置。
//
// 示例：
//
//	client := cwe.NewAPIClient()
//	weakness, err := client.GetWeakness(context.Background(), 79)
//
//	// 自定义配置
//	client := cwe.NewAPIClient(
//	    cwe.WithAPITimeout(60*time.Second),
//	    cwe.WithAPIRateLimit(0.5, 1),
//	)
type APIClient struct {
	httpClient *HTTPClient
	baseURL    string
}

// APIClientOption 是APIClient的配置选项函数类型
type APIClientOption func(*APIClient)

// WithAPIBaseURL 设置API的基础URL。
//
// 参数：
//   - url: 基础URL，默认为 "https://cwe-api.mitre.org/api/v1"
func WithAPIBaseURL(url string) APIClientOption {
	return func(c *APIClient) {
		c.baseURL = url
	}
}

// WithAPITimeout 设置API请求的超时时间。
//
// 参数：
//   - timeout: 超时时间
func WithAPITimeout(timeout time.Duration) APIClientOption {
	return func(c *APIClient) {
		c.httpClient.client.Timeout = timeout
	}
}

// WithAPIRateLimit 设置API请求的速率限制。
//
// 参数：
//   - rate: 每秒允许的请求数
//   - burst: 允许的突发请求数
func WithAPIRateLimit(rate float64, burst int) APIClientOption {
	return func(c *APIClient) {
		c.httpClient.rateLimiter = NewRateLimiter(rate, burst)
	}
}

// WithAPIRetry 设置API请求的重试策略。
//
// 参数：
//   - maxRetries: 最大重试次数
//   - delay: 重试间隔
func WithAPIRetry(maxRetries int, delay time.Duration) APIClientOption {
	return func(c *APIClient) {
		c.httpClient.maxRetries = maxRetries
		c.httpClient.retryDelay = delay
	}
}

// WithAPIHTTPClient 设置自定义的HTTP客户端选项。
//
// 参数：
//   - opts: HTTPClient的配置选项
func WithAPIHTTPClient(opts ...HTTPClientOption) APIClientOption {
	return func(c *APIClient) {
		for _, opt := range opts {
			opt(c.httpClient)
		}
	}
}

// NewAPIClient 创建一个新的APIClient实例。
//
// 默认使用MITRE CWE REST API的基础URL和30秒超时。
//
// 参数：
//   - opts: 可选的配置选项
//
// 返回值：
//   - *APIClient: 新创建的APIClient实例
func NewAPIClient(opts ...APIClientOption) *APIClient {
	client := &APIClient{
		baseURL: DefaultBaseURL,
		httpClient: NewHTTPClient(DefaultBaseURL,
			WithHTTPRateLimiter(0.1, 1), // MITRE API默认限制：每10秒1个请求
		),
	}

	for _, opt := range opts {
		opt(client)
	}

	// 确保HTTP客户端的baseURL与APIClient一致
	client.httpClient.baseURL = client.baseURL

	return client
}

// GetBaseURL 获取API的基础URL。
func (c *APIClient) GetBaseURL() string {
	return c.baseURL
}

// SetBaseURL 设置API的基础URL。
func (c *APIClient) SetBaseURL(url string) {
	c.baseURL = url
	c.httpClient.SetBaseURL(url)
}

// GetHTTPClient 获取底层的HTTPClient实例。
func (c *APIClient) GetHTTPClient() *HTTPClient {
	return c.httpClient
}

// SetHTTPClient 设置底层的HTTPClient实例。
func (c *APIClient) SetHTTPClient(client *HTTPClient) {
	c.httpClient = client
}

// GetRateLimiter 获取速率限制器。
func (c *APIClient) GetRateLimiter() *RateLimiter {
	return c.httpClient.GetRateLimiter()
}

// SetRateLimiter 设置速率限制器。
func (c *APIClient) SetRateLimiter(limiter *RateLimiter) {
	c.httpClient.SetRateLimiter(limiter)
}

// Close 关闭API客户端，释放资源。
func (c *APIClient) Close() {
	c.httpClient.Close()
}
