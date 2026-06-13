package cwe

import (
	"fmt"
	"time"
)

// CWEError 是所有CWE SDK错误的基础类型。
//
// 该类型提供了统一的错误结构，包含错误码、消息和详细信息，
// 方便上层应用进行错误分类和处理。
type CWEError struct {
	// Code 错误码
	Code string
	// Message 错误消息
	Message string
	// Detail 详细信息
	Detail string
	// Err 被包装的内部错误
	Err error
}

// Error 实现error接口，返回格式化的错误消息。
func (e *CWEError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("cwe: [%s] %s: %s: %v", e.Code, e.Message, e.Detail, e.Err)
	}
	if e.Detail != "" {
		return fmt.Sprintf("cwe: [%s] %s: %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("cwe: [%s] %s", e.Code, e.Message)
}

// Unwrap 返回被包装的内部错误，支持errors.Is/As链式查找。
func (e *CWEError) Unwrap() error {
	return e.Err
}

// InvalidCWEIDError 表示CWE ID格式无效的错误。
//
// 当CWE ID不符合 "CWE-NNN" 格式或数字部分无效时返回此错误。
type InvalidCWEIDError struct {
	*CWEError
	// ID 无效的CWE ID输入
	ID string
}

// NewInvalidCWEIDError 创建一个新的InvalidCWEIDError。
//
// 参数：
//   - id: 无效的CWE ID字符串
func NewInvalidCWEIDError(id string) *InvalidCWEIDError {
	return &InvalidCWEIDError{
		CWEError: &CWEError{
			Code:    "INVALID_CWE_ID",
			Message: "CWE ID格式无效",
			Detail:  fmt.Sprintf("输入值: %q", id),
		},
		ID: id,
	}
}

// CWENotFoundError 表示CWE条目未找到的错误。
//
// 当在注册表中查找不存在的CWE条目时返回此错误。
type CWENotFoundError struct {
	*CWEError
	// ID 未找到的CWE ID
	ID int
}

// NewCWENotFoundError 创建一个新的CWENotFoundError。
//
// 参数：
//   - id: 未找到的CWE ID
func NewCWENotFoundError(id int) *CWENotFoundError {
	return &CWENotFoundError{
		CWEError: &CWEError{
			Code:    "CWE_NOT_FOUND",
			Message: "CWE条目未找到",
			Detail:  fmt.Sprintf("CWE ID: %d", id),
		},
		ID: id,
	}
}

// APIError 表示CWE API调用失败的错误。
//
// 当对MITRE CWE REST API的请求返回错误状态码时返回此错误。
type APIError struct {
	*CWEError
	// StatusCode HTTP状态码
	StatusCode int
	// URL 请求的URL
	URL string
	// Method HTTP方法
	Method string
}

// NewAPIError 创建一个新的APIError。
//
// 参数：
//   - statusCode: HTTP状态码
//   - url: 请求的URL
//   - method: HTTP方法
func NewAPIError(statusCode int, url, method string) *APIError {
	return &APIError{
		CWEError: &CWEError{
			Code:    "API_ERROR",
			Message: "CWE API调用失败",
			Detail:  fmt.Sprintf("HTTP %d %s %s", statusCode, method, url),
		},
		StatusCode: statusCode,
		URL:        url,
		Method:     method,
	}
}

// RateLimitError 表示请求速率超限的错误。
//
// 当请求频率超过API速率限制时返回此错误。
type RateLimitError struct {
	*CWEError
	// RetryAfter 建议等待的时间
	RetryAfter time.Duration
}

// NewRateLimitError 创建一个新的RateLimitError。
//
// 参数：
//   - retryAfter: 建议等待的时间
func NewRateLimitError(retryAfter time.Duration) *RateLimitError {
	return &RateLimitError{
		CWEError: &CWEError{
			Code:    "RATE_LIMIT",
			Message: "请求速率超限",
			Detail:  fmt.Sprintf("建议等待: %v", retryAfter),
		},
		RetryAfter: retryAfter,
	}
}

// ValidationError 表示模型验证失败的错误。
//
// 当CWE条目的字段值不符合约束条件时返回此错误。
type ValidationError struct {
	*CWEError
	// Field 验证失败的字段名
	Field string
	// Value 验证失败的值
	Value string
}

// NewValidationError 创建一个新的ValidationError。
//
// 参数：
//   - field: 验证失败的字段名
//   - value: 验证失败的值
func NewValidationError(field, value string) *ValidationError {
	return &ValidationError{
		CWEError: &CWEError{
			Code:    "VALIDATION_ERROR",
			Message: "模型验证失败",
			Detail:  fmt.Sprintf("字段 %q 的值 %q 无效", field, value),
		},
		Field: field,
		Value: value,
	}
}

// ParseError 表示解析失败的错误。
//
// 当XML、JSON或其他格式的数据解析失败时返回此错误。
type ParseError struct {
	*CWEError
	// Offset 解析失败的位置偏移量
	Offset int64
}

// NewParseError 创建一个新的ParseError。
//
// 参数：
//   - detail: 解析失败的详细描述
//   - offset: 解析失败的位置偏移量
func NewParseError(detail string, offset int64) *ParseError {
	return &ParseError{
		CWEError: &CWEError{
			Code:    "PARSE_ERROR",
			Message: "数据解析失败",
			Detail:  detail,
		},
		Offset: offset,
	}
}

// RelationshipError 表示关系操作失败的错误。
//
// 当尝试建立无效的CWE关系时返回此错误。
type RelationshipError struct {
	*CWEError
	// From 源CWE ID
	From string
	// To 目标CWE ID
	To string
	// Nature 关系类型
	Nature RelationshipNature
}

// NewRelationshipError 创建一个新的RelationshipError。
//
// 参数：
//   - from: 源CWE ID
//   - to: 目标CWE ID
//   - nature: 关系类型
func NewRelationshipError(from, to string, nature RelationshipNature) *RelationshipError {
	return &RelationshipError{
		CWEError: &CWEError{
			Code:    "RELATIONSHIP_ERROR",
			Message: "关系操作失败",
			Detail:  fmt.Sprintf("无法建立 %s -> %s (类型: %s) 的关系", from, to, nature),
		},
		From:   from,
		To:     to,
		Nature: nature,
	}
}
