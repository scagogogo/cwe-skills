package cweskills

import "sort"

// Statistics 包含CWE注册表的统计信息。
type Statistics struct {
	// TotalCount CWE弱点条目总数
	TotalCount int `json:"total_count"`
	// WeaknessCount 弱点数量
	WeaknessCount int `json:"weakness_count"`
	// CategoryCount 类别数量
	CategoryCount int `json:"category_count"`
	// ViewCount 视图数量
	ViewCount int `json:"view_count"`
	// CompoundElementCount 复合元素数量
	CompoundElementCount int `json:"compound_element_count"`
	// ByAbstraction 按抽象层级统计
	ByAbstraction map[Abstraction]int `json:"by_abstraction"`
	// ByStatus 按状态统计
	ByStatus map[Status]int `json:"by_status"`
	// ByStructure 按结构类型统计
	ByStructure map[Structure]int `json:"by_structure"`
	// ByLikelihood 按利用可能性统计
	ByLikelihood map[LikelihoodOfExploit]int `json:"by_likelihood"`
	// TopScopes 影响范围统计（按出现次数排序）
	TopScopes []ConsequenceScopeCount `json:"top_scopes"`
}

// ConsequenceScopeCount 表示后果范围的统计计数。
type ConsequenceScopeCount struct {
	// Scope 后果范围
	Scope ConsequenceScope `json:"scope"`
	// Count 出现次数
	Count int `json:"count"`
}

// ComputeStatistics 计算CWE注册表的完整统计信息。
//
// 参数：
//   - r: CWE注册表
//
// 返回值：
//   - *Statistics: 统计信息
func ComputeStatistics(r *Registry) *Statistics {
	if r == nil {
		return &Statistics{
			ByAbstraction: make(map[Abstraction]int),
			ByStatus:      make(map[Status]int),
			ByStructure:   make(map[Structure]int),
			ByLikelihood:  make(map[LikelihoodOfExploit]int),
		}
	}

	stats := &Statistics{
		TotalCount:           r.Size(),
		WeaknessCount:        r.Size(),
		CategoryCount:        r.CategoryCount(),
		ViewCount:            r.ViewCount(),
		CompoundElementCount: r.CompoundElementCount(),
		ByAbstraction:        make(map[Abstraction]int),
		ByStatus:             make(map[Status]int),
		ByStructure:          make(map[Structure]int),
		ByLikelihood:         make(map[LikelihoodOfExploit]int),
	}

	scopeCount := make(map[ConsequenceScope]int)

	for _, cwe := range r.GetAll() {
		// 按抽象层级统计
		if cwe.Abstraction != "" {
			stats.ByAbstraction[cwe.Abstraction]++
		}

		// 按状态统计
		if cwe.Status != "" {
			stats.ByStatus[cwe.Status]++
		}

		// 按结构类型统计
		if cwe.Structure != "" {
			stats.ByStructure[cwe.Structure]++
		}

		// 按利用可能性统计
		if cwe.LikelihoodOfExploit != "" {
			stats.ByLikelihood[cwe.LikelihoodOfExploit]++
		}

		// 统计后果范围
		for _, consequence := range cwe.CommonConsequences {
			for _, scope := range consequence.Scopes {
				scopeCount[scope]++
			}
		}
	}

	// 将范围统计转换为排序后的切片
	stats.TopScopes = make([]ConsequenceScopeCount, 0, len(scopeCount))
	for scope, count := range scopeCount {
		stats.TopScopes = append(stats.TopScopes, ConsequenceScopeCount{
			Scope: scope,
			Count: count,
		})
	}
	sort.Slice(stats.TopScopes, func(i, j int) bool {
		return stats.TopScopes[i].Count > stats.TopScopes[j].Count
	})

	return stats
}

// CountByAbstraction 统计指定抽象层级的CWE条目数量。
func CountByAbstraction(r *Registry, a Abstraction) int {
	if r == nil {
		return 0
	}
	return len(FindByAbstraction(r, a))
}

// CountByStatus 统计指定状态的CWE条目数量。
func CountByStatus(r *Registry, s Status) int {
	if r == nil {
		return 0
	}
	return len(FindByStatus(r, s))
}

// CountByLikelihood 统计指定利用可能性的CWE条目数量。
func CountByLikelihood(r *Registry, l LikelihoodOfExploit) int {
	if r == nil {
		return 0
	}
	return len(FindByLikelihood(r, l))
}

// CountByScope 统计具有指定后果范围的CWE条目数量。
func CountByScope(r *Registry, scope ConsequenceScope) int {
	if r == nil {
		return 0
	}
	return len(FindByConsequenceScope(r, scope))
}
