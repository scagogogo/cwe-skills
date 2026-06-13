package cwe

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// cweIDRegex 用于匹配CWE ID的正则表达式，支持多种格式：
// "CWE-79", "cwe-79", "CWE79", "79" 等
var cweIDRegex = regexp.MustCompile(`(?i)CWE[-\s]?(\d+)`)

// FormatCWEID 将CWE ID格式化为标准格式 "CWE-NNN"。
//
// 该函数接受各种常见的CWE ID格式，统一输出为大写的 "CWE-NNN" 格式。
// 输入可以是纯数字、"CWE-79"、"cwe-79"、"CWE79" 等格式。
//
// 参数：
//   - id: 需要格式化的CWE ID字符串
//
// 返回值：
//   - string: 格式化后的标准CWE ID字符串，如 "CWE-79"
//   - error: 如果输入不是有效的CWE ID格式，返回InvalidCWEIDError
//
// 示例：
//
//	FormatCWEID("79")     // 返回 "CWE-79", nil
//	FormatCWEID("cwe-79") // 返回 "CWE-79", nil
//	FormatCWEID("CWE79")  // 返回 "CWE-79", nil
//	FormatCWEID("")       // 返回 "", InvalidCWEIDError
func FormatCWEID(id string) (string, error) {
	num, err := ParseCWEID(id)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("CWE-%d", num), nil
}

// ParseCWEID 从CWE ID字符串中提取数字部分。
//
// 支持的输入格式：
//   - 纯数字: "79", "079"
//   - 标准格式: "CWE-79", "CWE-079"
//   - 无连字符: "CWE79"
//   - 大小写不敏感: "cwe-79", "Cwe-79"
//
// 参数：
//   - id: 需要解析的CWE ID字符串
//
// 返回值：
//   - int: 提取出的数字ID
//   - error: 如果输入不是有效的CWE ID格式，返回InvalidCWEIDError
//
// 示例：
//
//	ParseCWEID("CWE-79")  // 返回 79, nil
//	ParseCWEID("79")      // 返回 79, nil
//	ParseCWEID("cwe-079") // 返回 79, nil
//	ParseCWEID("")        // 返回 0, InvalidCWEIDError
//	ParseCWEID("abc")     // 返回 0, InvalidCWEIDError
func ParseCWEID(id string) (int, error) {
	if id == "" {
		return 0, NewInvalidCWEIDError(id)
	}

	id = strings.TrimSpace(id)

	// 尝试匹配 "CWE-NNN" 或 "CWE NNN" 或 "CWENNN" 格式
	matches := cweIDRegex.FindStringSubmatch(id)
	if len(matches) >= 2 {
		num, err := strconv.Atoi(matches[1])
		if err != nil {
			return 0, NewInvalidCWEIDError(id)
		}
		if num <= 0 {
			return 0, NewInvalidCWEIDError(id)
		}
		return num, nil
	}

	// 尝试作为纯数字解析
	num, err := strconv.Atoi(id)
	if err != nil {
		return 0, NewInvalidCWEIDError(id)
	}
	if num <= 0 {
		return 0, NewInvalidCWEIDError(id)
	}

	return num, nil
}

// FormatCWEIDFromInt 将整数ID格式化为标准CWE ID字符串 "CWE-NNN"。
//
// 参数：
//   - id: 整数形式的CWE ID
//
// 返回值：
//   - string: 格式化后的标准CWE ID字符串
//
// 示例：
//
//	FormatCWEIDFromInt(79)  // 返回 "CWE-79"
//	FormatCWEIDFromInt(1000) // 返回 "CWE-1000"
func FormatCWEIDFromInt(id int) string {
	return fmt.Sprintf("CWE-%d", id)
}

// IsCWEID 检查给定的字符串是否为有效的CWE ID格式。
//
// 该函数仅检查格式是否正确，不验证该CWE ID是否实际存在于MITRE数据库中。
//
// 参数：
//   - text: 需要检查的字符串
//
// 返回值：
//   - bool: 如果是有效的CWE ID格式返回true，否则返回false
//
// 示例：
//
//	IsCWEID("CWE-79")  // 返回 true
//	IsCWEID("79")      // 返回 true
//	IsCWEID("abc")     // 返回 false
//	IsCWEID("")        // 返回 false
func IsCWEID(text string) bool {
	_, err := ParseCWEID(text)
	return err == nil
}

// ValidateCWEID 验证CWE ID格式并返回详细的错误信息。
//
// 该函数对CWE ID进行完整验证，包括格式检查和基本的有效性检查。
// 与IsCWEID不同，该函数返回详细的错误信息，方便调用方进行错误处理。
//
// 参数：
//   - text: 需要验证的CWE ID字符串
//
// 返回值：
//   - error: 如果验证通过返回nil，否则返回InvalidCWEIDError
//
// 示例：
//
//	err := ValidateCWEID("CWE-79")  // 返回 nil
//	err := ValidateCWEID("")        // 返回 InvalidCWEIDError
func ValidateCWEID(text string) error {
	if text == "" {
		return NewInvalidCWEIDError(text)
	}

	_, err := ParseCWEID(text)
	if err != nil {
		return err
	}

	return nil
}

// ExtractCWEIDs 从文本中提取所有CWE ID。
//
// 该函数在给定的文本中搜索所有匹配CWE ID格式的子串，
// 返回所有找到的标准格式CWE ID列表。
//
// 参数：
//   - text: 需要搜索的文本
//
// 返回值：
//   - []string: 找到的所有CWE ID列表，格式为 "CWE-NNN"。如果没有找到，返回空切片。
//
// 示例：
//
//	ExtractCWEIDs("See CWE-79 and CWE-89 for details")
//	// 返回 ["CWE-79", "CWE-89"]
//
//	ExtractCWEIDs("No CWE IDs here")
//	// 返回 []
func ExtractCWEIDs(text string) []string {
	if text == "" {
		return []string{}
	}

	matches := cweIDRegex.FindAllStringSubmatch(text, -1)
	result := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) >= 2 {
			num, err := strconv.Atoi(match[1])
			if err == nil && num > 0 {
				result = append(result, fmt.Sprintf("CWE-%d", num))
			}
		}
	}

	return result
}

// ExtractFirstCWEID 从文本中提取第一个CWE ID。
//
// 参数：
//   - text: 需要搜索的文本
//
// 返回值：
//   - string: 找到的第一个CWE ID，格式为 "CWE-NNN"。如果没有找到，返回空字符串。
//
// 示例：
//
//	ExtractFirstCWEID("See CWE-79 and CWE-89")
//	// 返回 "CWE-79"
//
//	ExtractFirstCWEID("No CWE IDs here")
//	// 返回 ""
func ExtractFirstCWEID(text string) string {
	if text == "" {
		return ""
	}

	match := cweIDRegex.FindStringSubmatch(text)
	if len(match) >= 2 {
		num, err := strconv.Atoi(match[1])
		if err == nil && num > 0 {
			return fmt.Sprintf("CWE-%d", num)
		}
	}

	return ""
}

// CompareCWEIDs 比较两个CWE ID的数值大小。
//
// 参数：
//   - a: 第一个CWE ID字符串
//   - b: 第二个CWE ID字符串
//
// 返回值：
//   - int: 如果a < b返回-1，a == b返回0，a > b返回1
//   - error: 如果任一CWE ID格式无效，返回错误
//
// 示例：
//
//	CompareCWEIDs("CWE-79", "CWE-89")   // 返回 -1, nil
//	CompareCWEIDs("CWE-79", "CWE-79")   // 返回 0, nil
//	CompareCWEIDs("CWE-89", "CWE-79")   // 返回 1, nil
func CompareCWEIDs(a, b string) (int, error) {
	numA, err := ParseCWEID(a)
	if err != nil {
		return 0, fmt.Errorf("比较失败，第一个CWE ID无效: %w", err)
	}

	numB, err := ParseCWEID(b)
	if err != nil {
		return 0, fmt.Errorf("比较失败，第二个CWE ID无效: %w", err)
	}

	if numA < numB {
		return -1, nil
	}
	if numA > numB {
		return 1, nil
	}
	return 0, nil
}
