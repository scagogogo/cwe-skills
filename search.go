package cweskills

import (
	"strings"
)

// FindByID 在注册表中根据ID查找CWE弱点条目。
//
// 参数：
//   - r: CWE注册表
//   - id: 要查找的CWE ID
//
// 返回值：
//   - *CWE: 找到的CWE条目，如果未找到返回nil
//   - bool: 是否找到
func FindByID(r *Registry, id int) (*CWE, bool) {
	if r == nil {
		return nil, false
	}
	return r.Get(id)
}

// FindByKeyword 在注册表中根据关键字搜索CWE条目。
//
// 搜索范围包括名称和描述字段，匹配是不区分大小写的。
//
// 参数：
//   - r: CWE注册表
//   - keyword: 搜索关键字
//
// 返回值：
//   - []*CWE: 匹配的CWE条目列表
func FindByKeyword(r *Registry, keyword string) []*CWE {
	if r == nil || keyword == "" {
		return nil
	}

	keywordLower := strings.ToLower(keyword)
	var result []*CWE

	for _, cwe := range r.GetAll() {
		if strings.Contains(strings.ToLower(cwe.Name), keywordLower) ||
			strings.Contains(strings.ToLower(cwe.Description), keywordLower) {
			result = append(result, cwe)
		}
	}

	return result
}

// FindByAbstraction 在注册表中查找指定抽象层级的所有CWE条目。
//
// 参数：
//   - r: CWE注册表
//   - abstraction: 抽象层级
//
// 返回值：
//   - []*CWE: 匹配的CWE条目列表
func FindByAbstraction(r *Registry, abstraction Abstraction) []*CWE {
	if r == nil {
		return nil
	}

	var result []*CWE
	for _, cwe := range r.GetAll() {
		if cwe.Abstraction == abstraction {
			result = append(result, cwe)
		}
	}

	return result
}

// FindByStatus 在注册表中查找指定状态的所有CWE条目。
//
// 参数：
//   - r: CWE注册表
//   - status: CWE状态
//
// 返回值：
//   - []*CWE: 匹配的CWE条目列表
func FindByStatus(r *Registry, status Status) []*CWE {
	if r == nil {
		return nil
	}

	var result []*CWE
	for _, cwe := range r.GetAll() {
		if cwe.Status == status {
			result = append(result, cwe)
		}
	}

	return result
}

// FindByLikelihood 在注册表中查找指定利用可能性的所有CWE条目。
//
// 参数：
//   - r: CWE注册表
//   - likelihood: 利用可能性
//
// 返回值：
//   - []*CWE: 匹配的CWE条目列表
func FindByLikelihood(r *Registry, likelihood LikelihoodOfExploit) []*CWE {
	if r == nil {
		return nil
	}

	var result []*CWE
	for _, cwe := range r.GetAll() {
		if cwe.LikelihoodOfExploit == likelihood {
			result = append(result, cwe)
		}
	}

	return result
}

// FindByConsequenceScope 在注册表中查找具有指定后果范围的所有CWE条目。
//
// 参数：
//   - r: CWE注册表
//   - scope: 后果范围
//
// 返回值：
//   - []*CWE: 匹配的CWE条目列表
func FindByConsequenceScope(r *Registry, scope ConsequenceScope) []*CWE {
	if r == nil {
		return nil
	}

	var result []*CWE
	for _, cwe := range r.GetAll() {
		if cwe.HasConsequenceScope(scope) {
			result = append(result, cwe)
		}
	}

	return result
}

// FindByStructure 在注册表中查找指定结构类型的所有CWE条目。
//
// 参数：
//   - r: CWE注册表
//   - structure: 结构类型
//
// 返回值：
//   - []*CWE: 匹配的CWE条目列表
func FindByStructure(r *Registry, structure Structure) []*CWE {
	if r == nil {
		return nil
	}

	var result []*CWE
	for _, cwe := range r.GetAll() {
		if cwe.Structure == structure {
			result = append(result, cwe)
		}
	}

	return result
}

// FindTopLevel 查找所有顶层（Pillar抽象级别）的CWE条目。
//
// 参数：
//   - r: CWE注册表
//
// 返回值：
//   - []*CWE: 所有顶层CWE条目
func FindTopLevel(r *Registry) []*CWE {
	return FindByAbstraction(r, AbstractionPillar)
}

// FindBaseWeaknesses 查找所有基础级别的CWE弱点。
//
// 参数：
//   - r: CWE注册表
//
// 返回值：
//   - []*CWE: 所有基础级别CWE条目
func FindBaseWeaknesses(r *Registry) []*CWE {
	return FindByAbstraction(r, AbstractionBase)
}

// FindChains 查找所有链式CWE条目。
//
// 参数：
//   - r: CWE注册表
//
// 返回值：
//   - []*CWE: 所有链式CWE条目
func FindChains(r *Registry) []*CWE {
	return FindByStructure(r, StructureChain)
}

// FindComposites 查找所有复合CWE条目。
//
// 参数：
//   - r: CWE注册表
//
// 返回值：
//   - []*CWE: 所有复合CWE条目
func FindComposites(r *Registry) []*CWE {
	return FindByStructure(r, StructureComposite)
}
