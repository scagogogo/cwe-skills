// Package cweskills 提供了对CWE（Common Weakness Enumeration，通用缺陷枚举）的完整支持，
// 包括CWE条目的解析、验证、搜索、过滤、关系导航、树构建、序列化等功能。
//
// 本包旨在作为网络安全产品的底层SDK，支持SAST/DAST工具、漏洞管理平台、
// 合规检查系统等上层应用基于此构建。
//
// 主要功能：
//   - CWE ID的格式化、解析、验证和提取
//   - 完整的枚举类型定义（抽象层级、状态、关系类型等）
//   - 结构化错误类型
//   - 知名CWE列表（Top 25、OWASP Top 10等）
//   - 核心数据模型（Weakness、Category、View、CompoundElement）
//   - 关系导航（父/子/祖先/后代/链/组合）
//   - 内存注册表与索引
//   - 搜索与过滤
//   - MITRE CWE REST API客户端
//   - XML目录解析（离线模式）
//   - JSON/XML/CSV序列化
//   - 速率限制的HTTP客户端
//
// 示例：
//
//	// 解析CWE ID
//	id, err := cwe.ParseCWEID("CWE-79")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(id) // 输出: 79
//
//	// 格式化CWE ID
//	formatted := cwe.FormatCWEID("79")
//	fmt.Println(formatted) // 输出: CWE-79
//
//	// 检查是否为Top 25
//	if cwe.IsInTop25(79) {
//	    fmt.Println("CWE-79 is in CWE Top 25")
//	}
package cweskills

// Version 表示本SDK的版本号
const Version = "v0.0.1"
