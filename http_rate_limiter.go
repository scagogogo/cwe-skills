package cwe

import (
	"context"
	"sync"
	"time"
)

// RateLimiter 实现了一个基于令牌桶算法的速率限制器。
//
// 该速率限制器使用令牌桶算法控制请求频率，
// 适用于对MITRE CWE API等外部服务进行速率限制。
// 本实现仅使用Go标准库，不依赖任何第三方包。
//
// 使用场景：
//   - 限制对MITRE CWE API的请求频率
//   - 防止因请求过快而被外部服务拒绝
//   - 控制并发请求的速率
//
// 示例：
//
//	limiter := cwe.NewRateLimiter(1.0, 5)  // 每秒1个请求，突发5个
//	if limiter.Allow() {
//	    // 执行请求
//	} else {
//	    // 等待
//	    limiter.Wait(context.Background())
//	}
type RateLimiter struct {
	mu         sync.Mutex
	rate       float64       // 每秒允许的令牌数
	burst      int           // 桶容量（最大突发数）
	tokens     float64       // 当前令牌数
	lastRefill time.Time     // 上次填充令牌的时间
	interval   time.Duration // 请求间隔（用于兼容旧接口）
	lastReq    time.Time     // 上次请求时间（用于兼容旧接口）
}

// NewRateLimiter 创建一个新的速率限制器。
//
// 参数：
//   - rate: 每秒允许的请求数。例如1.0表示每秒1个请求，0.1表示每10秒1个请求。
//   - burst: 允许的突发请求数。当令牌桶满时，可以一次性发送burst个请求。
//
// 返回值：
//   - *RateLimiter: 新创建的速率限制器
//
// 示例：
//
//	// 每秒1个请求，最多突发5个
//	limiter := cwe.NewRateLimiter(1.0, 5)
//
//	// 每10秒1个请求（MITRE API默认限制）
//	limiter := cwe.NewRateLimiter(0.1, 1)
func NewRateLimiter(rate float64, burst int) *RateLimiter {
	now := time.Now()
	return &RateLimiter{
		rate:       rate,
		burst:      burst,
		tokens:     float64(burst), // 初始满桶
		lastRefill: now,
		interval:   time.Duration(float64(time.Second) / rate),
		lastReq:    time.Time{}, // 零值表示尚未使用
	}
}

// refill 填充令牌桶。
// 必须在持有锁的情况下调用。
func (r *RateLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(r.lastRefill)

	// 根据经过的时间计算应添加的令牌数
	newTokens := elapsed.Seconds() * r.rate
	r.tokens += newTokens

	// 令牌数不超过桶容量
	if r.tokens > float64(r.burst) {
		r.tokens = float64(r.burst)
	}

	r.lastRefill = now
}

// Allow 检查是否允许立即执行请求（非阻塞）。
//
// 如果令牌桶中有可用令牌，消耗一个令牌并返回true。
// 如果没有可用令牌，返回false，不消耗令牌。
//
// 返回值：
//   - bool: 如果允许请求返回true，否则返回false
func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.refill()

	if r.tokens >= 1.0 {
		r.tokens -= 1.0
		r.lastReq = time.Now()
		return true
	}

	return false
}

// Wait 阻塞等待直到可以执行请求，或上下文被取消。
//
// 如果令牌桶中没有可用令牌，该函数会等待直到有令牌可用。
// 可以通过context来取消等待。
//
// 参数：
//   - ctx: 用于取消等待的上下文
//
// 返回值：
//   - error: 如果等待被取消返回context的错误，否则返回nil
func (r *RateLimiter) Wait(ctx context.Context) error {
	for {
		r.mu.Lock()
		r.refill()

		if r.tokens >= 1.0 {
			r.tokens -= 1.0
			r.lastReq = time.Now()
			r.mu.Unlock()
			return nil
		}

		// 计算需要等待的时间
		waitTime := time.Duration((1.0-r.tokens)/r.rate) * time.Second
		if waitTime < time.Millisecond {
			waitTime = time.Millisecond
		}
		r.mu.Unlock()

		// 等待或取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
			// 继续尝试
		}
	}
}

// WaitForRequest 阻塞等待直到可以执行请求（兼容旧接口）。
//
// 该方法使用请求间隔模式，确保两次请求之间至少间隔指定的时间。
// 如果这是第一次请求，立即返回。
func (r *RateLimiter) WaitForRequest() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.lastReq.IsZero() {
		r.lastReq = time.Now()
		return
	}

	elapsed := time.Since(r.lastReq)
	if elapsed < r.interval {
		waitTime := r.interval - elapsed
		r.mu.Unlock()
		time.Sleep(waitTime)
		r.mu.Lock()
	}

	r.lastReq = time.Now()
}

// GetInterval 获取请求间隔时间。
//
// 返回值：
//   - time.Duration: 两次请求之间的最小间隔时间
func (r *RateLimiter) GetInterval() time.Duration {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.interval
}

// SetInterval 设置请求间隔时间。
//
// 参数：
//   - interval: 两次请求之间的最小间隔时间
func (r *RateLimiter) SetInterval(interval time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.interval = interval
	if interval > 0 {
		r.rate = float64(time.Second) / float64(interval)
	}
}

// GetRate 获取每秒允许的请求数。
//
// 返回值：
//   - float64: 每秒允许的请求数
func (r *RateLimiter) GetRate() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rate
}

// GetBurst 获取允许的突发请求数。
//
// 返回值：
//   - int: 允许的突发请求数
func (r *RateLimiter) GetBurst() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.burst
}

// ResetLastRequest 重置上次请求时间。
//
// 调用此方法后，下一次WaitForRequest将立即返回，
// 而不需要等待间隔时间。
func (r *RateLimiter) ResetLastRequest() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastReq = time.Time{}
	r.lastRefill = time.Now()
	r.tokens = float64(r.burst)
}

// Tokens 返回当前可用的令牌数。
//
// 返回值：
//   - float64: 当前可用的令牌数
func (r *RateLimiter) Tokens() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.refill()
	return r.tokens
}
