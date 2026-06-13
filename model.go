package cwe

import "fmt"

// CWE 表示一个CWE弱点条目，是SDK的核心类型。
//
// CWE（Common Weakness Enumeration，通用弱点枚举）条目描述了一种特定的安全弱点，
// 包括其描述、抽象层级、结构类型、后果、缓解措施等信息。
type CWE struct {
	// ID CWE条目的数字标识符
	ID int `json:"id" xml:"ID,attr"`
	// Name CWE条目的名称
	Name string `json:"name" xml:"Name"`
	// Abstraction 抽象层级（Pillar、Class、Base、Variant）
	Abstraction Abstraction `json:"abstraction,omitempty" xml:"Abstraction,omitempty"`
	// Structure 结构类型（Simple、Chain、Composite）
	Structure Structure `json:"structure,omitempty" xml:"Structure,omitempty"`
	// Status 状态（Stable、Usable、Draft、Incomplete、Obsolete、Deprecated）
	Status Status `json:"status,omitempty" xml:"Status,omitempty"`
	// Description 弱点描述
	Description string `json:"description" xml:"Description"`
	// ExtendedDescription 扩展描述，提供更详细的说明
	ExtendedDescription string `json:"extended_description,omitempty" xml:"Extended_Description,omitempty"`
	// LikelihoodOfExploit 被利用的可能性
	LikelihoodOfExploit LikelihoodOfExploit `json:"likelihood_of_exploit,omitempty" xml:"LikelihoodOfExploit,omitempty"`
	// CommonConsequences 常见后果列表
	CommonConsequences []Consequence `json:"common_consequences,omitempty" xml:"CommonConsequences>Consequence,omitempty"`
	// PotentialMitigations 潜在缓解措施列表
	PotentialMitigations []Mitigation `json:"potential_mitigations,omitempty" xml:"PotentialMitigations>Mitigation,omitempty"`
	// DemonstrativeExamples 示范性示例列表
	DemonstrativeExamples []DemonstrativeExample `json:"demonstrative_examples,omitempty" xml:"DemonstrativeExamples>DemonstrativeExample,omitempty"`
	// ObservedExamples 观察到的示例列表
	ObservedExamples []ObservedExample `json:"observed_examples,omitempty" xml:"ObservedExamples>ObservedExample,omitempty"`
	// References 参考文献列表
	References []Reference `json:"references,omitempty" xml:"References>Reference,omitempty"`
	// Relationships 关系列表
	Relationships []Relationship `json:"relationships,omitempty" xml:"Relationships>Relationship,omitempty"`
	// ApplicablePlatforms 适用平台信息
	ApplicablePlatforms *ApplicablePlatforms `json:"applicable_platforms,omitempty" xml:"ApplicablePlatforms,omitempty"`
	// ModesOfIntroduction 引入方式列表
	ModesOfIntroduction []Introduction `json:"modes_of_introduction,omitempty" xml:"ModesOfIntroduction>Introduction,omitempty"`
	// AlternateTerms 备用术语列表
	AlternateTerms []AlternateTerm `json:"alternate_terms,omitempty" xml:"AlternateTerms>AlternateTerm,omitempty"`
	// Notes 备注信息
	Notes string `json:"notes,omitempty" xml:"Notes,omitempty"`
	// ContentHistory 内容历史
	ContentHistory *ContentHistory `json:"content_history,omitempty" xml:"ContentHistory,omitempty"`
	// CWEType CWE条目类型（weakness、category、view、compound_element）
	CWEType string `json:"cwe_type" xml:"-"`
	// URL CWE条目的URL地址
	URL string `json:"url,omitempty" xml:"-"`
}

// CWEID 返回标准格式的CWE ID字符串。
//
// 返回值：
//   - string: 格式为"CWE-NNN"的标准CWE ID
func (c *CWE) CWEID() string {
	return FormatCWEIDFromInt(c.ID)
}

// IsWeakness 检查CWE条目是否为弱点类型。
//
// 返回值：
//   - bool: 如果CWEType为"weakness"返回true，否则返回false
func (c *CWE) IsWeakness() bool {
	return c.CWEType == "weakness"
}

// IsCategory 检查CWE条目是否为类别类型。
//
// 返回值：
//   - bool: 如果CWEType为"category"返回true，否则返回false
func (c *CWE) IsCategory() bool {
	return c.CWEType == "category"
}

// IsView 检查CWE条目是否为视图类型。
//
// 返回值：
//   - bool: 如果CWEType为"view"返回true，否则返回false
func (c *CWE) IsView() bool {
	return c.CWEType == "view"
}

// IsCompoundElement 检查CWE条目是否为复合元素类型。
//
// 返回值：
//   - bool: 如果CWEType为"compound_element"返回true，否则返回false
func (c *CWE) IsCompoundElement() bool {
	return c.CWEType == "compound_element"
}

// IsPillar 检查CWE条目是否为柱石级别。
//
// 返回值：
//   - bool: 如果Abstraction为AbstractionPillar返回true，否则返回false
func (c *CWE) IsPillar() bool {
	return c.Abstraction == AbstractionPillar
}

// IsBase 检查CWE条目是否为基础级别。
//
// 返回值：
//   - bool: 如果Abstraction为AbstractionBase返回true，否则返回false
func (c *CWE) IsBase() bool {
	return c.Abstraction == AbstractionBase
}

// IsVariant 检查CWE条目是否为变体级别。
//
// 返回值：
//   - bool: 如果Abstraction为AbstractionVariant返回true，否则返回false
func (c *CWE) IsVariant() bool {
	return c.Abstraction == AbstractionVariant
}

// IsChain 检查CWE条目是否为链式结构。
//
// 返回值：
//   - bool: 如果Structure为StructureChain返回true，否则返回false
func (c *CWE) IsChain() bool {
	return c.Structure == StructureChain
}

// IsComposite 检查CWE条目是否为复合结构。
//
// 返回值：
//   - bool: 如果Structure为StructureComposite返回true，否则返回false
func (c *CWE) IsComposite() bool {
	return c.Structure == StructureComposite
}

// IsStable 检查CWE条目是否为稳定状态。
//
// 返回值：
//   - bool: 如果Status为StatusStable返回true，否则返回false
func (c *CWE) IsStable() bool {
	return c.Status == StatusStable
}

// IsDeprecated 检查CWE条目是否已废弃。
//
// 返回值：
//   - bool: 如果Status为StatusDeprecated返回true，否则返回false
func (c *CWE) IsDeprecated() bool {
	return c.Status == StatusDeprecated
}

// GetParentIDs 获取此弱点的所有父级CWE ID。
//
// 通过遍历Relationships中Nature为ChildOf的关系来获取父级ID。
//
// 返回值：
//   - []int: 父级CWE ID列表
func (c *CWE) GetParentIDs() []int {
	var ids []int
	for _, r := range c.Relationships {
		if r.Nature == RelationshipChildOf {
			ids = append(ids, r.CWEID)
		}
	}
	return ids
}

// GetChildIDs 获取此弱点的所有子级CWE ID。
//
// 通过遍历Relationships中Nature为ParentOf的关系来获取子级ID。
//
// 返回值：
//   - []int: 子级CWE ID列表
func (c *CWE) GetChildIDs() []int {
	var ids []int
	for _, r := range c.Relationships {
		if r.Nature == RelationshipParentOf {
			ids = append(ids, r.CWEID)
		}
	}
	return ids
}

// GetPeerIDs 获取此弱点的所有对等CWE ID。
//
// 通过遍历Relationships中Nature为PeerOf或CanAlsoBe的关系来获取对等ID。
//
// 返回值：
//   - []int: 对等CWE ID列表
func (c *CWE) GetPeerIDs() []int {
	var ids []int
	for _, r := range c.Relationships {
		if r.Nature == RelationshipPeerOf || r.Nature == RelationshipCanAlsoBe {
			ids = append(ids, r.CWEID)
		}
	}
	return ids
}

// GetChainIDs 获取与此弱点相关的链式CWE ID。
//
// 通过遍历Relationships中Nature为CanPrecede或CanFollow的关系来获取链式ID。
//
// 返回值：
//   - []int: 链式CWE ID列表
func (c *CWE) GetChainIDs() []int {
	var ids []int
	for _, r := range c.Relationships {
		if r.Nature == RelationshipCanPrecede || r.Nature == RelationshipCanFollow {
			ids = append(ids, r.CWEID)
		}
	}
	return ids
}

// HasConsequenceScope 检查CWE条目是否包含指定的后果范围。
//
// 遍历所有CommonConsequences，检查是否有任何后果包含指定的安全范围。
//
// 参数：
//   - scope: 需要检查的安全范围
//
// 返回值：
//   - bool: 如果任何后果包含该安全范围返回true，否则返回false
func (c *CWE) HasConsequenceScope(scope ConsequenceScope) bool {
	for _, cons := range c.CommonConsequences {
		if cons.HasScope(scope) {
			return true
		}
	}
	return false
}

// Validate 验证CWE条目的有效性。
//
// 检查条件：
//   - ID必须大于0
//   - Name不能为空
//
// 返回值：
//   - error: 如果验证失败返回ValidationError，否则返回nil
func (c *CWE) Validate() error {
	if c.ID <= 0 {
		return NewValidationError("ID", fmt.Sprintf("%d", c.ID))
	}
	if c.Name == "" {
		return NewValidationError("Name", "空字符串")
	}
	return nil
}

// NewCWE 创建一个新的CWE弱点条目。
//
// 创建的CWE条目默认CWEType为"weakness"。
//
// 参数：
//   - id: CWE条目的数字标识符
//   - name: CWE条目的名称
//
// 返回值：
//   - *CWE: 新创建的CWE实例
func NewCWE(id int, name string) *CWE {
	return &CWE{
		ID:      id,
		Name:    name,
		CWEType: "weakness",
	}
}

// Mitigation 表示缓解措施。
//
// 缓解措施描述了如何减少或消除特定弱点的风险。
type Mitigation struct {
	// Phase 缓解措施适用的阶段
	Phase MitigationPhase `json:"phase,omitempty" xml:"Phase,omitempty"`
	// Strategy 缓解策略名称
	Strategy string `json:"strategy,omitempty" xml:"Strategy,omitempty"`
	// Description 缓解措施的详细描述
	Description string `json:"description" xml:"Description"`
	// Effectiveness 缓解措施的有效性
	Effectiveness Effectiveness `json:"effectiveness,omitempty" xml:"Effectiveness,omitempty"`
}

// DemonstrativeExample 表示示范性示例。
//
// 示范性示例展示了该弱点在实际中可能出现的情况。
type DemonstrativeExample struct {
	// IntroText 示例的介绍文本
	IntroText string `json:"intro_text,omitempty" xml:"IntroText,omitempty"`
	// BodyText 示例的主体内容
	BodyText string `json:"body_text,omitempty" xml:"BodyText,omitempty"`
}

// ObservedExample 表示观察到的真实示例。
//
// 观察到的示例来自真实的漏洞报告或安全事件。
type ObservedExample struct {
	// Reference 参考编号
	Reference string `json:"reference,omitempty" xml:"Reference,omitempty"`
	// Description 示例描述
	Description string `json:"description" xml:"Description"`
	// Link 相关链接
	Link string `json:"link,omitempty" xml:"Link,omitempty"`
}

// Reference 表示参考文献。
//
// 参考文献提供了关于该弱点的更多信息来源。
type Reference struct {
	// ID 参考文献的标识符
	ID int `json:"id,omitempty" xml:"Reference_ID,omitempty"`
	// Author 作者
	Author string `json:"author,omitempty" xml:"Author,omitempty"`
	// Title 标题
	Title string `json:"title" xml:"Title"`
	// URL 链接地址
	URL string `json:"url,omitempty" xml:"URL,omitempty"`
}

// ApplicablePlatforms 表示CWE条目适用的平台信息。
//
// 适用平台包括编程语言、操作系统、架构和技术等方面。
type ApplicablePlatforms struct {
	// Languages 适用的编程语言列表
	Languages []PlatformEntry `json:"languages,omitempty" xml:"Languages>Language,omitempty"`
	// OperatingSystems 适用的操作系统列表
	OperatingSystems []PlatformEntry `json:"operating_systems,omitempty" xml:"Operating_Systems>Operating_System,omitempty"`
	// Architectures 适用的架构列表
	Architectures []PlatformEntry `json:"architectures,omitempty" xml:"Architectures>Architecture,omitempty"`
	// Technologies 适用的技术列表
	Technologies []PlatformEntry `json:"technologies,omitempty" xml:"Technologies>Technology,omitempty"`
}

// PlatformEntry 表示平台条目。
//
// 每个平台条目包含名称和使用普遍程度。
type PlatformEntry struct {
	// Name 平台名称
	Name string `json:"name" xml:"Name,attr"`
	// Prevalence 使用普遍程度
	Prevalence Prevalence `json:"prevalence,omitempty" xml:"Prevalence,attr,omitempty"`
}

// Introduction 表示弱点的引入方式。
//
// 引入方式描述了弱点在哪个阶段被引入到软件中。
type Introduction struct {
	// Phase 引入阶段
	Phase IntroductionPhase `json:"phase" xml:"Phase"`
	// Description 引入方式的描述
	Description string `json:"description,omitempty" xml:"Description,omitempty"`
}

// AlternateTerm 表示备用术语。
//
// 备用术语提供了该CWE条目的其他常用名称或术语。
type AlternateTerm struct {
	// Term 备用术语名称
	Term string `json:"term" xml:"Term"`
	// Description 备用术语的描述
	Description string `json:"description,omitempty" xml:"Description,omitempty"`
}

// ContentHistory 表示CWE条目的内容历史。
//
// 内容历史记录了CWE条目的提交和修改记录。
type ContentHistory struct {
	// Submission 提交信息
	Submission *HistoryEntry `json:"submission,omitempty" xml:"Submission,omitempty"`
	// Modifications 修改记录列表
	Modifications []HistoryEntry `json:"modifications,omitempty" xml:"Modifications>Modification,omitempty"`
}

// HistoryEntry 表示历史条目。
//
// 历史条目记录了CWE条目的提交者或修改者信息。
type HistoryEntry struct {
	// Name 提交者或修改者姓名
	Name string `json:"name,omitempty" xml:"Name,omitempty"`
	// Organization 所属组织
	Organization string `json:"organization,omitempty" xml:"Organization,omitempty"`
	// Date 日期
	Date string `json:"date,omitempty" xml:"Date,omitempty"`
	// Comment 注释
	Comment string `json:"comment,omitempty" xml:"Comment,omitempty"`
}

// Category 表示CWE类别条目。
//
// 类别是一种CWE条目类型，用于将相关的弱点分组归类。
type Category struct {
	// ID 类别的数字标识符
	ID int `json:"id" xml:"ID,attr"`
	// Name 类别名称
	Name string `json:"name" xml:"Name"`
	// Status 状态
	Status Status `json:"status,omitempty" xml:"Status,omitempty"`
	// Description 类别描述
	Description string `json:"description" xml:"Description"`
	// Relationships 关系列表
	Relationships []Relationship `json:"relationships,omitempty" xml:"Relationships>Relationship,omitempty"`
	// Notes 备注信息
	Notes string `json:"notes,omitempty" xml:"Notes,omitempty"`
	// References 参考文献列表
	References []Reference `json:"references,omitempty" xml:"References>Reference,omitempty"`
	// ContentHistory 内容历史
	ContentHistory *ContentHistory `json:"content_history,omitempty" xml:"ContentHistory,omitempty"`
}

// NewCategory 创建一个新的CWE类别条目。
//
// 参数：
//   - id: 类别的数字标识符
//   - name: 类别名称
//
// 返回值：
//   - *Category: 新创建的类别实例
func NewCategory(id int, name string) *Category {
	return &Category{
		ID:   id,
		Name: name,
	}
}

// View 表示CWE视图条目。
//
// 视图是一种CWE条目类型，提供了从特定角度查看和组织弱点的方式。
type View struct {
	// ID 视图的数字标识符
	ID int `json:"id" xml:"ID,attr"`
	// Name 视图名称
	Name string `json:"name" xml:"Name"`
	// Type 视图类型（Graph、Explicit Slice、Implicit Slice）
	Type ViewType `json:"type,omitempty" xml:"Type,omitempty"`
	// Status 状态
	Status Status `json:"status,omitempty" xml:"Status,omitempty"`
	// Description 视图描述
	Description string `json:"description" xml:"Description"`
	// Members 视图成员列表
	Members []ViewMember `json:"members,omitempty" xml:"Members>ViewMember,omitempty"`
	// References 参考文献列表
	References []Reference `json:"references,omitempty" xml:"References>Reference,omitempty"`
	// ContentHistory 内容历史
	ContentHistory *ContentHistory `json:"content_history,omitempty" xml:"ContentHistory,omitempty"`
}

// ViewMember 表示视图成员。
//
// 视图成员描述了某个CWE条目在特定视图中的归属关系。
type ViewMember struct {
	// CWEID 成员的CWE ID
	CWEID int `json:"cwe_id" xml:"CWE_ID"`
	// ViewID 所属视图ID
	ViewID int `json:"view_id" xml:"View_ID"`
	// Direct 是否为直接成员
	Direct bool `json:"direct" xml:"Direct"`
	// Predicate 谓词，可选
	Predicate string `json:"predicate,omitempty" xml:"Predicate,omitempty"`
}

// NewView 创建一个新的CWE视图条目。
//
// 参数：
//   - id: 视图的数字标识符
//   - name: 视图名称
//   - viewType: 视图类型
//
// 返回值：
//   - *View: 新创建的视图实例
func NewView(id int, name string, viewType ViewType) *View {
	return &View{
		ID:   id,
		Name: name,
		Type: viewType,
	}
}

// CompoundElement 表示CWE复合元素条目。
//
// 复合元素是一种CWE条目类型，描述了由多个弱点组合而成的复合弱点，
// 包括链式弱点和复合弱点。
type CompoundElement struct {
	// ID 复合元素的数字标识符
	ID int `json:"id" xml:"ID,attr"`
	// Name 复合元素名称
	Name string `json:"name" xml:"Name"`
	// Structure 结构类型（Chain或Composite）
	Structure Structure `json:"structure" xml:"Structure"`
	// Status 状态
	Status Status `json:"status,omitempty" xml:"Status,omitempty"`
	// Description 复合元素描述
	Description string `json:"description" xml:"Description"`
	// Relationships 关系列表
	Relationships []Relationship `json:"relationships,omitempty" xml:"Relationships>Relationship,omitempty"`
}

// NewCompoundElement 创建一个新的CWE复合元素条目。
//
// 参数：
//   - id: 复合元素的数字标识符
//   - name: 复合元素名称
//   - structure: 结构类型
//
// 返回值：
//   - *CompoundElement: 新创建的复合元素实例
func NewCompoundElement(id int, name string, structure Structure) *CompoundElement {
	return &CompoundElement{
		ID:        id,
		Name:      name,
		Structure: structure,
	}
}