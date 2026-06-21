package cweskills

// 知名视图ID常量
const (
	// CWEViewResearchConcepts 研究概念视图 (CWE-1000)
	// 按照抽象概念组织CWE条目的层次结构
	CWEViewResearchConcepts = 1000

	// CWEViewDevelopmentConcepts 软件开发视图 (CWE-699)
	// 按照软件开发活动组织CWE条目
	CWEViewDevelopmentConcepts = 699

	// CWEViewHardwareDesign 硬件设计视图 (CWE-1199)
	// 按照硬件设计活动组织CWE条目
	CWEViewHardwareDesign = 1199

	// CWEViewCWECrossSection CWE横截面视图 (CWE-888)
	// 提供CWE条目的横截面视图
	CWEViewCWECrossSection = 888

	// CWEViewComprehensiveDictionary 综合CWE字典 (CWE-1400)
	// 包含所有CWE条目的综合视图
	CWEViewComprehensiveDictionary = 1400
)

// CWETop25 包含CWE Top 25最危险软件弱点列表（2024版）。
//
// 该列表由MITRE基于NVD数据的频率分析和CVSS评分计算得出，
// 代表了对软件最严重的安全威胁。
//
// 使用场景：
//   - 安全工具优先级排序
//   - 开发者安全培训重点
//   - 漏洞管理风险评估
var CWETop25 = []int{
	79,   // Cross-site Scripting (XSS)
	89,   // SQL Injection
	352,  // Cross-Site Request Forgery (CSRF)
	862,  // Missing Authorization
	787,  // Out-of-bounds Write
	22,   // Path Traversal
	416,  // Use After Free
	125,  // Out-of-bounds Read
	78,   // OS Command Injection
	94,   // Code Injection
	120,  // Buffer Copy without Checking Size of Input
	434,  // Unrestricted Upload of File with Dangerous Type
	476,  // NULL Pointer Dereference
	121,  // Stack-based Buffer Overflow
	502,  // Deserialization of Untrusted Data
	122,  // Heap-based Buffer Overflow
	863,  // Incorrect Authorization
	20,   // Improper Input Validation
	284,  // Improper Access Control
	200,  // Exposure of Sensitive Information
	306,  // Missing Authentication for Critical Function
	918,  // Server-Side Request Forgery (SSRF)
	77,   // Command Injection
	639,  // Authorization Bypass Through User-Controlled Key
	770,  // Allocation of Resources Without Limits or Throttling
}

// OWASPTop10 包含OWASP Top 10（2021版）到CWE ID的映射。
//
// OWASP Top 10是Web应用安全风险的行业标准列表，
// 每个类别映射到相关的CWE ID。
//
// 使用场景：
//   - Web应用安全评估
//   - 合规性检查
//   - 安全测试用例设计
var OWASPTop10 = map[string][]int{
	"A01:2021-Broken Access Control": {
		22, 23, 35, 59, 78, 94, 200, 201, 219, 255, 269, 276,
		284, 285, 287, 306, 346, 639, 651, 668, 862, 863, 922,
	},
	"A02:2021-Cryptographic Failures": {
		260, 261, 295, 310, 311, 312, 319, 325, 326, 327, 328,
		329, 330, 337, 338, 340, 347, 522, 757, 759, 760, 780,
	},
	"A03:2021-Injection": {
		20, 74, 75, 77, 78, 79, 80, 83, 87, 88, 89, 90, 91,
		94, 95, 96, 97, 98, 99, 100, 113, 116, 138, 141, 147,
		150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160,
		161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171,
	},
	"A04:2021-Insecure Design": {
		209, 235, 256, 267, 284, 285, 287, 311, 326, 384, 393, 664, 863,
	},
	"A05:2021-Security Misconfiguration": {
		2, 5, 11, 13, 15, 16, 260, 315, 520, 526, 537, 540, 544, 546,
		547, 548, 611, 613, 614, 759, 760, 1021,
	},
	"A06:2021-Vulnerable and Outdated Components": {},
	"A07:2021-Identification and Authentication Failures": {
		255, 256, 258, 259, 260, 287, 288, 290, 294, 295, 297,
		306, 307, 346, 384, 521, 522, 523, 613, 620, 640, 798,
	},
	"A08:2021-Software and Data Integrity Failures": {
		311, 345, 353, 426, 494, 502, 565, 610, 653, 754, 829, 912,
	},
	"A09:2021-Security Logging and Monitoring Failures": {
		117, 223, 532, 778,
	},
	"A10:2021-Server-Side Request Forgery (SSRF)": {
		918, 1021,
	},
}

// SANSTop25 包含SANS Top 25最危险软件错误列表。
//
// SANS Top 25由SANS Institute与MITRE合作编制，
// 侧重于可被攻击者利用来获取系统控制权的编程错误。
var SANSTop25 = []int{
	89,  // SQL Injection
	78,  // OS Command Injection
	79,  // Cross-site Scripting (XSS)
	20,  // Improper Input Validation
	22,  // Path Traversal
	352, // Cross-Site Request Forgery (CSRF)
	416, // Use After Free
	787, // Out-of-bounds Write
	125, // Out-of-bounds Read
	94,  // Code Injection
	190, // Integer Overflow or Wraparound
	434, // Unrestricted Upload of File with Dangerous Type
	862, // Missing Authorization
	287, // Improper Authentication
	306, // Missing Authentication for Critical Function
	863, // Incorrect Authorization
	798, // Use of Hard-coded Credentials
	502, // Deserialization of Untrusted Data
	77,  // Command Injection
	119, // Improper Restriction of Operations within the Bounds of a Memory Buffer
	639, // Authorization Bypass Through User-Controlled Key
	770, // Allocation of Resources Without Limits or Throttling
	918, // Server-Side Request Forgery (SSRF)
	476, // NULL Pointer Dereference
	200, // Exposure of Sensitive Information
}

// IsInTop25 检查给定的CWE ID是否在CWE Top 25列表中。
//
// 参数：
//   - cweID: 需要检查的CWE ID数字
//
// 返回值：
//   - bool: 如果在Top 25中返回true，否则返回false
func IsInTop25(cweID int) bool {
	for _, id := range CWETop25 {
		if id == cweID {
			return true
		}
	}
	return false
}

// IsInOWASPTop10 检查给定的CWE ID是否在OWASP Top 10映射中。
//
// 参数：
//   - cweID: 需要检查的CWE ID数字
//
// 返回值：
//   - bool: 如果在OWASP Top 10中返回true，否则返回false
func IsInOWASPTop10(cweID int) bool {
	for _, ids := range OWASPTop10 {
		for _, id := range ids {
			if id == cweID {
				return true
			}
		}
	}
	return false
}

// IsInSANSTop25 检查给定的CWE ID是否在SANS Top 25列表中。
//
// 参数：
//   - cweID: 需要检查的CWE ID数字
//
// 返回值：
//   - bool: 如果在SANS Top 25中返回true，否则返回false
func IsInSANSTop25(cweID int) bool {
	for _, id := range SANSTop25 {
		if id == cweID {
			return true
		}
	}
	return false
}

// GetOWASPCategory 获取给定CWE ID所属的OWASP Top 10类别。
//
// 如果CWE ID属于多个类别，只返回第一个匹配的类别。
//
// 参数：
//   - cweID: 需要查询的CWE ID数字
//
// 返回值：
//   - string: OWASP类别名称（如 "A01:2021-Broken Access Control"），如果不在任何类别中返回空字符串
func GetOWASPCategory(cweID int) string {
	for category, ids := range OWASPTop10 {
		for _, id := range ids {
			if id == cweID {
				return category
			}
		}
	}
	return ""
}

// GetOWASPCategories 获取给定CWE ID所属的所有OWASP Top 10类别。
//
// 与GetOWASPCategory不同，该函数返回所有匹配的类别。
//
// 参数：
//   - cweID: 需要查询的CWE ID数字
//
// 返回值：
//   - []string: 所有匹配的OWASP类别名称列表，如果没有匹配返回空切片
func GetOWASPCategories(cweID int) []string {
	var categories []string
	for category, ids := range OWASPTop10 {
		for _, id := range ids {
			if id == cweID {
				categories = append(categories, category)
				break
			}
		}
	}
	return categories
}

// IsInWellKnownView 检查给定的视图ID是否为知名视图。
//
// 知名视图包括：CWE-1000（研究概念）、CWE-699（软件开发）、
// CWE-1199（硬件设计）、CWE-888（横截面）、CWE-1400（综合字典）。
//
// 参数：
//   - viewID: 需要检查的视图ID数字
//
// 返回值：
//   - bool: 如果是知名视图返回true，否则返回false
func IsInWellKnownView(viewID int) bool {
	switch viewID {
	case CWEViewResearchConcepts, CWEViewDevelopmentConcepts,
		CWEViewHardwareDesign, CWEViewCWECrossSection,
		CWEViewComprehensiveDictionary:
		return true
	default:
		return false
	}
}
