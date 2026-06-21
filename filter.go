package cweskills

import (
	"sort"
	"strings"
)

// FilterOption 定义CWE条目的过滤选项。
//
// 多个FilterOption可以组合使用，所有非零值的选项都会作为过滤条件，
// 各条件之间是AND（与）关系。
type FilterOption struct {
	// Abstraction 按抽象层级过滤
	Abstraction Abstraction
	// Status 按状态过滤
	Status Status
	// Structure 按结构类型过滤
	Structure Structure
	// Likelihood 按利用可能性过滤
	Likelihood LikelihoodOfExploit
	// MinID 按最小CWE ID过滤（包含）
	MinID int
	// MaxID 按最大CWE ID过滤（包含）
	MaxID int
	// Keyword 按关键字过滤（匹配名称和描述，不区分大小写）
	Keyword string
	// Scope 按后果范围过滤
	Scope ConsequenceScope
}

// Filter 根据过滤选项过滤CWE条目列表。
//
// 所有非零值的选项都会作为过滤条件，条件之间是AND关系。
// 零值的选项被忽略。
//
// 参数：
//   - cwes: 要过滤的CWE条目列表
//   - opts: 过滤选项（可以指定多个，但只使用第一个）
//
// 返回值：
//   - []*CWE: 过滤后的CWE条目列表
func Filter(cwes []*CWE, opts ...FilterOption) []*CWE {
	if len(cwes) == 0 {
		return cwes
	}

	if len(opts) == 0 {
		return cwes
	}

	opt := opts[0]
	result := make([]*CWE, 0, len(cwes))

	for _, cwe := range cwes {
		if matchesFilter(cwe, opt) {
			result = append(result, cwe)
		}
	}

	return result
}

// matchesFilter 检查CWE条目是否匹配过滤选项。
func matchesFilter(cwe *CWE, opt FilterOption) bool {
	// 按抽象层级过滤
	if opt.Abstraction != "" && cwe.Abstraction != opt.Abstraction {
		return false
	}

	// 按状态过滤
	if opt.Status != "" && cwe.Status != opt.Status {
		return false
	}

	// 按结构类型过滤
	if opt.Structure != "" && cwe.Structure != opt.Structure {
		return false
	}

	// 按利用可能性过滤
	if opt.Likelihood != "" && cwe.LikelihoodOfExploit != opt.Likelihood {
		return false
	}

	// 按最小ID过滤
	if opt.MinID > 0 && cwe.ID < opt.MinID {
		return false
	}

	// 按最大ID过滤
	if opt.MaxID > 0 && cwe.ID > opt.MaxID {
		return false
	}

	// 按关键字过滤
	if opt.Keyword != "" && !cweMatchesKeyword(cwe, opt.Keyword) {
		return false
	}

	// 按后果范围过滤
	if opt.Scope != "" && !cwe.HasConsequenceScope(opt.Scope) {
		return false
	}

	return true
}

// cweMatchesKeyword 检查CWE条目是否匹配关键字。
func cweMatchesKeyword(cwe *CWE, keyword string) bool {
	kw := lower(keyword)
	return lowerContains(cwe.Name, kw) || lowerContains(cwe.Description, kw)
}

// SortByID 按CWE ID升序排序。
func SortByID(cwes []*CWE) []*CWE {
	if len(cwes) == 0 {
		return cwes
	}
	sort.Slice(cwes, func(i, j int) bool {
		return cwes[i].ID < cwes[j].ID
	})
	return cwes
}

// SortByName 按CWE名称升序排序。
func SortByName(cwes []*CWE) []*CWE {
	if len(cwes) == 0 {
		return cwes
	}
	sort.Slice(cwes, func(i, j int) bool {
		return cwes[i].Name < cwes[j].Name
	})
	return cwes
}

// SortByAbstraction 按抽象层级从高到低排序。
func SortByAbstraction(cwes []*CWE) []*CWE {
	if len(cwes) == 0 {
		return cwes
	}
	sort.Slice(cwes, func(i, j int) bool {
		return cwes[i].Abstraction.AbstractionOrder() > cwes[j].Abstraction.AbstractionOrder()
	})
	return cwes
}

// GroupByAbstraction 按抽象层级分组。
func GroupByAbstraction(cwes []*CWE) map[Abstraction][]*CWE {
	result := make(map[Abstraction][]*CWE)
	for _, cwe := range cwes {
		result[cwe.Abstraction] = append(result[cwe.Abstraction], cwe)
	}
	return result
}

// GroupByStatus 按状态分组。
func GroupByStatus(cwes []*CWE) map[Status][]*CWE {
	result := make(map[Status][]*CWE)
	for _, cwe := range cwes {
		result[cwe.Status] = append(result[cwe.Status], cwe)
	}
	return result
}

// GroupByLikelihood 按利用可能性分组。
func GroupByLikelihood(cwes []*CWE) map[LikelihoodOfExploit][]*CWE {
	result := make(map[LikelihoodOfExploit][]*CWE)
	for _, cwe := range cwes {
		result[cwe.LikelihoodOfExploit] = append(result[cwe.LikelihoodOfExploit], cwe)
	}
	return result
}

// Deduplicate 对CWE条目列表去重（按ID）。
func Deduplicate(cwes []*CWE) []*CWE {
	if len(cwes) == 0 {
		return cwes
	}

	seen := make(map[int]bool)
	result := make([]*CWE, 0, len(cwes))
	for _, cwe := range cwes {
		if !seen[cwe.ID] {
			seen[cwe.ID] = true
			result = append(result, cwe)
		}
	}
	return result
}

// lower 是strings.ToLower的简写。
func lower(s string) string {
	return strings.ToLower(s)
}

// lowerContains 检查s是否包含substr（不区分大小写）。
func lowerContains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), substr)
}
