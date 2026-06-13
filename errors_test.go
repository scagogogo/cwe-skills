package cwe

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewInvalidCWEIDError(t *testing.T) {
	err := NewInvalidCWEIDError("abc")
	if err.ID != "abc" {
		t.Errorf("InvalidCWEIDError.ID = %q, want %q", err.ID, "abc")
	}
	if err.Code != "INVALID_CWE_ID" {
		t.Errorf("InvalidCWEIDError.Code = %q, want %q", err.Code, "INVALID_CWE_ID")
	}
	if err.Message != "CWE ID格式无效" {
		t.Errorf("InvalidCWEIDError.Message = %q, want %q", err.Message, "CWE ID格式无效")
	}
	if !strings.Contains(err.Detail, "abc") {
		t.Errorf("InvalidCWEIDError.Detail should contain %q, got %q", "abc", err.Detail)
	}
	errStr := err.Error()
	if !strings.Contains(errStr, "INVALID_CWE_ID") {
		t.Errorf("InvalidCWEIDError.Error() should contain code, got %q", errStr)
	}
	if !strings.Contains(errStr, "CWE ID格式无效") {
		t.Errorf("InvalidCWEIDError.Error() should contain message, got %q", errStr)
	}
}

func TestNewCWENotFoundError(t *testing.T) {
	err := NewCWENotFoundError(79)
	if err.ID != 79 {
		t.Errorf("CWENotFoundError.ID = %d, want %d", err.ID, 79)
	}
	if err.Code != "CWE_NOT_FOUND" {
		t.Errorf("CWENotFoundError.Code = %q, want %q", err.Code, "CWE_NOT_FOUND")
	}
	if err.Message != "CWE条目未找到" {
		t.Errorf("CWENotFoundError.Message = %q, want %q", err.Message, "CWE条目未找到")
	}
	errStr := err.Error()
	if !strings.Contains(errStr, "CWE_NOT_FOUND") {
		t.Errorf("CWENotFoundError.Error() should contain code, got %q", errStr)
	}
	if !strings.Contains(errStr, "79") {
		t.Errorf("CWENotFoundError.Error() should contain ID, got %q", errStr)
	}
}

func TestNewAPIError(t *testing.T) {
	err := NewAPIError(404, "https://example.com/api", "GET")
	if err.StatusCode != 404 {
		t.Errorf("APIError.StatusCode = %d, want %d", err.StatusCode, 404)
	}
	if err.URL != "https://example.com/api" {
		t.Errorf("APIError.URL = %q, want %q", err.URL, "https://example.com/api")
	}
	if err.Method != "GET" {
		t.Errorf("APIError.Method = %q, want %q", err.Method, "GET")
	}
	if err.Code != "API_ERROR" {
		t.Errorf("APIError.Code = %q, want %q", err.Code, "API_ERROR")
	}
	errStr := err.Error()
	if !strings.Contains(errStr, "API_ERROR") {
		t.Errorf("APIError.Error() should contain code, got %q", errStr)
	}
	if !strings.Contains(errStr, "404") {
		t.Errorf("APIError.Error() should contain status code, got %q", errStr)
	}
	if !strings.Contains(errStr, "GET") {
		t.Errorf("APIError.Error() should contain method, got %q", errStr)
	}
	if !strings.Contains(errStr, "https://example.com/api") {
		t.Errorf("APIError.Error() should contain URL, got %q", errStr)
	}
}

func TestNewRateLimitError(t *testing.T) {
	retryAfter := 5 * time.Second
	err := NewRateLimitError(retryAfter)
	if err.RetryAfter != retryAfter {
		t.Errorf("RateLimitError.RetryAfter = %v, want %v", err.RetryAfter, retryAfter)
	}
	if err.Code != "RATE_LIMIT" {
		t.Errorf("RateLimitError.Code = %q, want %q", err.Code, "RATE_LIMIT")
	}
	errStr := err.Error()
	if !strings.Contains(errStr, "RATE_LIMIT") {
		t.Errorf("RateLimitError.Error() should contain code, got %q", errStr)
	}
	if !strings.Contains(errStr, "5s") {
		t.Errorf("RateLimitError.Error() should contain retry duration, got %q", errStr)
	}
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("Field1", "BadValue")
	if err.Field != "Field1" {
		t.Errorf("ValidationError.Field = %q, want %q", err.Field, "Field1")
	}
	if err.Value != "BadValue" {
		t.Errorf("ValidationError.Value = %q, want %q", err.Value, "BadValue")
	}
	if err.Code != "VALIDATION_ERROR" {
		t.Errorf("ValidationError.Code = %q, want %q", err.Code, "VALIDATION_ERROR")
	}
	errStr := err.Error()
	if !strings.Contains(errStr, "VALIDATION_ERROR") {
		t.Errorf("ValidationError.Error() should contain code, got %q", errStr)
	}
	if !strings.Contains(errStr, "Field1") {
		t.Errorf("ValidationError.Error() should contain field, got %q", errStr)
	}
	if !strings.Contains(errStr, "BadValue") {
		t.Errorf("ValidationError.Error() should contain value, got %q", errStr)
	}
}

func TestNewParseError(t *testing.T) {
	err := NewParseError("unexpected EOF", 42)
	if err.Offset != 42 {
		t.Errorf("ParseError.Offset = %d, want %d", err.Offset, 42)
	}
	if err.Code != "PARSE_ERROR" {
		t.Errorf("ParseError.Code = %q, want %q", err.Code, "PARSE_ERROR")
	}
	if err.Detail != "unexpected EOF" {
		t.Errorf("ParseError.Detail = %q, want %q", err.Detail, "unexpected EOF")
	}
	errStr := err.Error()
	if !strings.Contains(errStr, "PARSE_ERROR") {
		t.Errorf("ParseError.Error() should contain code, got %q", errStr)
	}
	if !strings.Contains(errStr, "unexpected EOF") {
		t.Errorf("ParseError.Error() should contain detail, got %q", errStr)
	}
}

func TestNewRelationshipError(t *testing.T) {
	err := NewRelationshipError("CWE-79", "CWE-74", RelationshipChildOf)
	if err.From != "CWE-79" {
		t.Errorf("RelationshipError.From = %q, want %q", err.From, "CWE-79")
	}
	if err.To != "CWE-74" {
		t.Errorf("RelationshipError.To = %q, want %q", err.To, "CWE-74")
	}
	if err.Nature != RelationshipChildOf {
		t.Errorf("RelationshipError.Nature = %q, want %q", err.Nature, RelationshipChildOf)
	}
	if err.Code != "RELATIONSHIP_ERROR" {
		t.Errorf("RelationshipError.Code = %q, want %q", err.Code, "RELATIONSHIP_ERROR")
	}
	errStr := err.Error()
	if !strings.Contains(errStr, "RELATIONSHIP_ERROR") {
		t.Errorf("RelationshipError.Error() should contain code, got %q", errStr)
	}
	if !strings.Contains(errStr, "CWE-79") {
		t.Errorf("RelationshipError.Error() should contain From, got %q", errStr)
	}
	if !strings.Contains(errStr, "CWE-74") {
		t.Errorf("RelationshipError.Error() should contain To, got %q", errStr)
	}
	if !strings.Contains(errStr, "ChildOf") {
		t.Errorf("RelationshipError.Error() should contain Nature, got %q", errStr)
	}
}

func TestCWEError_WithErr(t *testing.T) {
	innerErr := errors.New("inner error")
	e := &CWEError{
		Code:    "TEST",
		Message: "test message",
		Detail:  "test detail",
		Err:     innerErr,
	}
	errStr := e.Error()
	if !strings.Contains(errStr, "inner error") {
		t.Errorf("CWEError.Error() with Err should contain wrapped error, got %q", errStr)
	}
	if !strings.Contains(errStr, "TEST") {
		t.Errorf("CWEError.Error() with Err should contain code, got %q", errStr)
	}
	if !strings.Contains(errStr, "test message") {
		t.Errorf("CWEError.Error() with Err should contain message, got %q", errStr)
	}
	if !strings.Contains(errStr, "test detail") {
		t.Errorf("CWEError.Error() with Err should contain detail, got %q", errStr)
	}
}

func TestCWEError_WithoutErr(t *testing.T) {
	e := &CWEError{
		Code:    "TEST",
		Message: "test message",
		Detail:  "test detail",
	}
	errStr := e.Error()
	if strings.Contains(errStr, "inner error") {
		t.Errorf("CWEError.Error() without Err should not contain wrapped error, got %q", errStr)
	}
	if !strings.Contains(errStr, "TEST") {
		t.Errorf("CWEError.Error() without Err should contain code, got %q", errStr)
	}
	if !strings.Contains(errStr, "test message") {
		t.Errorf("CWEError.Error() without Err should contain message, got %q", errStr)
	}
	if !strings.Contains(errStr, "test detail") {
		t.Errorf("CWEError.Error() without Err should contain detail, got %q", errStr)
	}
}

func TestCWEError_WithoutDetail(t *testing.T) {
	e := &CWEError{
		Code:    "TEST",
		Message: "test message",
	}
	errStr := e.Error()
	if !strings.Contains(errStr, "TEST") {
		t.Errorf("CWEError.Error() without Detail should contain code, got %q", errStr)
	}
	if !strings.Contains(errStr, "test message") {
		t.Errorf("CWEError.Error() without Detail should contain message, got %q", errStr)
	}
	// Should be shorter - no detail segment
	if strings.Contains(errStr, "test detail") {
		t.Errorf("CWEError.Error() without Detail should not contain detail, got %q", errStr)
	}
}

func TestCWEError_AllFields(t *testing.T) {
	innerErr := errors.New("base")
	e := &CWEError{
		Code:    "FULL",
		Message: "full msg",
		Detail:  "full detail",
		Err:     innerErr,
	}
	errStr := e.Error()
	if !strings.Contains(errStr, "FULL") {
		t.Errorf("expected code in error, got %q", errStr)
	}
	if !strings.Contains(errStr, "full msg") {
		t.Errorf("expected message in error, got %q", errStr)
	}
	if !strings.Contains(errStr, "full detail") {
		t.Errorf("expected detail in error, got %q", errStr)
	}
	if !strings.Contains(errStr, "base") {
		t.Errorf("expected wrapped error in error, got %q", errStr)
	}
}

func TestCWEError_Unwrap(t *testing.T) {
	innerErr := errors.New("inner")
	e := &CWEError{
		Code:    "TEST",
		Message: "msg",
		Err:     innerErr,
	}
	unwrapped := e.Unwrap()
	if unwrapped != innerErr {
		t.Errorf("CWEError.Unwrap() = %v, want %v", unwrapped, innerErr)
	}
	// Test with no wrapped error
	e2 := &CWEError{
		Code:    "TEST",
		Message: "msg",
	}
	unwrapped2 := e2.Unwrap()
	if unwrapped2 != nil {
		t.Errorf("CWEError.Unwrap() with no Err should return nil, got %v", unwrapped2)
	}
}

func TestCWEError_UnwrapWithErrorsIs(t *testing.T) {
	innerErr := errors.New("inner")
	e := &CWEError{
		Code:    "TEST",
		Message: "msg",
		Err:     innerErr,
	}
	if !errors.Is(e, innerErr) {
		t.Error("errors.Is should find wrapped error")
	}
}
