package cwe

import "encoding/json"

// APIResponse 是MITRE CWE REST API的基础响应结构。
//
// 所有API响应都包含Data字段和可选的Message字段。
type APIResponse struct {
	// Data 包含API返回的数据
	Data json.RawMessage `json:"Data"`
	// Message 可选的消息字段
	Message string `json:"Message,omitempty"`
}

// VersionResponse 是版本API的响应结构。
type VersionResponse struct {
	// Version CWE数据版本号
	Version string `json:"version"`
	// ReleaseDate 发布日期
	ReleaseDate string `json:"releaseDate"`
	// Name 版本名称
	Name string `json:"name"`
}

// WeaknessesResponse 是弱点查询API的响应结构。
type WeaknessesResponse struct {
	// Data 包含弱点数据的原始JSON
	Data json.RawMessage `json:"Data"`
	// Weaknesses 弱点列表
	Weaknesses []CWE `json:"weaknesses,omitempty"`
}

// CategoriesResponse 是类别查询API的响应结构。
type CategoriesResponse struct {
	// Data 包含类别数据的原始JSON
	Data json.RawMessage `json:"Data,omitempty"`
	// Categories 类别列表
	Categories []Category `json:"categories,omitempty"`
}

// ViewsResponse 是视图查询API的响应结构。
type ViewsResponse struct {
	// Data 包含视图数据的原始JSON
	Data json.RawMessage `json:"Data,omitempty"`
	// Views 视图列表
	Views []View `json:"views,omitempty"`
}

// RelationsResponse 是关系查询API的响应结构。
type RelationsResponse struct {
	// Data 包含关系数据的原始JSON
	Data json.RawMessage `json:"Data,omitempty"`
}

// CWEsResponse 是批量CWE查询API的响应结构。
type CWEsResponse struct {
	// Data 包含CWE数据的原始JSON
	Data json.RawMessage `json:"Data"`
	// Weaknesses 弱点映射表，键为CWE ID
	Weaknesses map[string]*CWE `json:"weaknesses,omitempty"`
}
