package cweskills

// Abstraction 表示CWE条目的抽象层级。
//
// CWE按照抽象程度从高到低分为四个层级：
// Pillar（柱石）> Class（类）> Base（基础）> Variant（变体）。
// 抽象层级越高，描述越通用；层级越低，描述越具体。
//
// 其中 Base 级别是映射到漏洞根因的首选级别。
type Abstraction string

const (
	// AbstractionPillar 表示柱石级别，最高抽象层级，代表一个主题
	// 例如：CWE-664 不正确的资源生命周期控制
	AbstractionPillar Abstraction = "Pillar"

	// AbstractionClass 表示类级别，通常与特定语言或技术无关
	// 例如：CWE-74 输出中特殊元素的不当中和（注入）
	AbstractionClass Abstraction = "Class"

	// AbstractionBase 表示基础级别，足够具体以推断检测/预防方法
	// 例如：CWE-79 XSS, CWE-89 SQL注入
	// 这是映射到漏洞根因的首选级别
	AbstractionBase Abstraction = "Base"

	// AbstractionVariant 表示变体级别，特定于某资源、技术或上下文
	// 例如：CWE-83 网页属性中脚本的不当中和
	AbstractionVariant Abstraction = "Variant"
)

// String 返回抽象层级的字符串表示。
func (a Abstraction) String() string { return string(a) }

// IsValid 检查抽象层级是否为有效值。
func (a Abstraction) IsValid() bool {
	switch a {
	case AbstractionPillar, AbstractionClass, AbstractionBase, AbstractionVariant:
		return true
	default:
		return false
	}
}

// ParseAbstraction 从字符串解析抽象层级。
//
// 参数：
//   - s: 需要解析的字符串
//
// 返回值：
//   - Abstraction: 解析出的抽象层级
//   - error: 如果字符串不是有效的抽象层级，返回错误
func ParseAbstraction(s string) (Abstraction, error) {
	a := Abstraction(s)
	if !a.IsValid() {
		return "", NewValidationError("Abstraction", s)
	}
	return a, nil
}

// AllAbstractionValues 返回所有有效的抽象层级值。
func AllAbstractionValues() []Abstraction {
	return []Abstraction{AbstractionPillar, AbstractionClass, AbstractionBase, AbstractionVariant}
}

// AbstractionOrder 返回抽象层级的排序权重，层级越高值越大。
// Pillar=4, Class=3, Base=2, Variant=1, 未知=0
func (a Abstraction) AbstractionOrder() int {
	switch a {
	case AbstractionPillar:
		return 4
	case AbstractionClass:
		return 3
	case AbstractionBase:
		return 2
	case AbstractionVariant:
		return 1
	default:
		return 0
	}
}

// Structure 表示CWE条目的结构类型。
//
// 结构类型描述了弱点之间的组成关系：
//   - Simple: 单一弱点，不依赖其他弱点
//   - Chain: 链式弱点，多个弱点必须按顺序可达才能产生漏洞
//   - Composite: 复合弱点，多个弱点必须同时存在才能产生漏洞
type Structure string

const (
	// StructureSimple 表示简单弱点，不依赖其他弱点的存在
	StructureSimple Structure = "Simple"

	// StructureChain 表示链式弱点，弱点必须按顺序可达
	// 例如：CWE-680 整数溢出到缓冲区溢出
	StructureChain Structure = "Chain"

	// StructureComposite 表示复合弱点，多个弱点必须同时存在
	// 例如：CWE-352 CSRF需要多个弱点同时存在
	StructureComposite Structure = "Composite"
)

// String 返回结构类型的字符串表示。
func (s Structure) String() string { return string(s) }

// IsValid 检查结构类型是否为有效值。
func (s Structure) IsValid() bool {
	switch s {
	case StructureSimple, StructureChain, StructureComposite:
		return true
	default:
		return false
	}
}

// ParseStructure 从字符串解析结构类型。
func ParseStructure(s string) (Structure, error) {
	st := Structure(s)
	if !st.IsValid() {
		return "", NewValidationError("Structure", s)
	}
	return st, nil
}

// AllStructureValues 返回所有有效的结构类型值。
func AllStructureValues() []Structure {
	return []Structure{StructureSimple, StructureChain, StructureComposite}
}

// Status 表示CWE条目的状态。
//
// 状态值描述了CWE条目的成熟度和稳定性：
//   - Stable: 所有重要元素已验证，不太可能发生显著变化
//   - Usable: 已经过深入审查，关键元素已验证
//   - Draft: 所有重要元素已填写，可能仍有问题或空缺
//   - Incomplete: 并非所有重要元素都已填写，无质量保证
//   - Obsolete: 仍然有效但不再相关，已被更新的实体取代
//   - Deprecated: 已从CWE中移除，是重复或错误创建的
type Status string

const (
	StatusStable     Status = "Stable"
	StatusUsable     Status = "Usable"
	StatusDraft      Status = "Draft"
	StatusIncomplete Status = "Incomplete"
	StatusObsolete   Status = "Obsolete"
	StatusDeprecated Status = "Deprecated"
)

// String 返回状态的字符串表示。
func (s Status) String() string { return string(s) }

// IsValid 检查状态是否为有效值。
func (s Status) IsValid() bool {
	switch s {
	case StatusStable, StatusUsable, StatusDraft, StatusIncomplete, StatusObsolete, StatusDeprecated:
		return true
	default:
		return false
	}
}

// ParseStatus 从字符串解析状态。
func ParseStatus(s string) (Status, error) {
	st := Status(s)
	if !st.IsValid() {
		return "", NewValidationError("Status", s)
	}
	return st, nil
}

// AllStatusValues 返回所有有效的状态值。
func AllStatusValues() []Status {
	return []Status{StatusStable, StatusUsable, StatusDraft, StatusIncomplete, StatusObsolete, StatusDeprecated}
}

// LikelihoodOfExploit 表示漏洞被利用的可能性。
type LikelihoodOfExploit string

const (
	LikelihoodHigh    LikelihoodOfExploit = "High"
	LikelihoodMedium  LikelihoodOfExploit = "Medium"
	LikelihoodLow     LikelihoodOfExploit = "Low"
	LikelihoodUnknown LikelihoodOfExploit = "Unknown"
)

// String 返回利用可能性的字符串表示。
func (l LikelihoodOfExploit) String() string { return string(l) }

// IsValid 检查利用可能性是否为有效值。
func (l LikelihoodOfExploit) IsValid() bool {
	switch l {
	case LikelihoodHigh, LikelihoodMedium, LikelihoodLow, LikelihoodUnknown:
		return true
	default:
		return false
	}
}

// ParseLikelihoodOfExploit 从字符串解析利用可能性。
func ParseLikelihoodOfExploit(s string) (LikelihoodOfExploit, error) {
	l := LikelihoodOfExploit(s)
	if !l.IsValid() {
		return "", NewValidationError("LikelihoodOfExploit", s)
	}
	return l, nil
}

// AllLikelihoodOfExploitValues 返回所有有效的利用可能性值。
func AllLikelihoodOfExploitValues() []LikelihoodOfExploit {
	return []LikelihoodOfExploit{LikelihoodHigh, LikelihoodMedium, LikelihoodLow, LikelihoodUnknown}
}

// LikelihoodOrder 返回利用可能性的排序权重。
// High=4, Medium=3, Low=2, Unknown=1, 未知=0
func (l LikelihoodOfExploit) LikelihoodOrder() int {
	switch l {
	case LikelihoodHigh:
		return 4
	case LikelihoodMedium:
		return 3
	case LikelihoodLow:
		return 2
	case LikelihoodUnknown:
		return 1
	default:
		return 0
	}
}

// RelationshipNature 表示CWE条目之间关系的类型。
//
// CWE定义了以下关系类型：
//   - ChildOf: 此弱点是目标弱点的子项（更具体）
//   - ParentOf: 此弱点是目标弱点的父项（更通用）
//   - CanPrecede: 此弱点可以创建使目标弱点成为可能的条件（链式前驱）
//   - CanFollow: 此弱点可以跟随目标弱点（链式后继）
//   - Requires: 此复合弱点需要目标弱点存在
//   - RequiredBy: 此弱点被目标复合弱点所需要
//   - CanAlsoBe: 此弱点在适当上下文中也可以被视为目标弱点
//   - PeerOf: 与目标弱点有相似性，但不适合其他关系类型
//   - MemberOf: 此条目是目标类别/视图的成员
//   - HasMember: 此类别/视图包含目标条目作为成员
type RelationshipNature string

const (
	RelationshipChildOf    RelationshipNature = "ChildOf"
	RelationshipParentOf   RelationshipNature = "ParentOf"
	RelationshipCanPrecede RelationshipNature = "CanPrecede"
	RelationshipCanFollow  RelationshipNature = "CanFollow"
	RelationshipRequires   RelationshipNature = "Requires"
	RelationshipRequiredBy RelationshipNature = "RequiredBy"
	RelationshipCanAlsoBe  RelationshipNature = "CanAlsoBe"
	RelationshipPeerOf     RelationshipNature = "PeerOf"
	RelationshipMemberOf   RelationshipNature = "MemberOf"
	RelationshipHasMember  RelationshipNature = "Has_Member"
)

// String 返回关系类型的字符串表示。
func (r RelationshipNature) String() string { return string(r) }

// IsValid 检查关系类型是否为有效值。
func (r RelationshipNature) IsValid() bool {
	switch r {
	case RelationshipChildOf, RelationshipParentOf, RelationshipCanPrecede,
		RelationshipCanFollow, RelationshipRequires, RelationshipRequiredBy,
		RelationshipCanAlsoBe, RelationshipPeerOf, RelationshipMemberOf,
		RelationshipHasMember:
		return true
	default:
		return false
	}
}

// ParseRelationshipNature 从字符串解析关系类型。
func ParseRelationshipNature(s string) (RelationshipNature, error) {
	r := RelationshipNature(s)
	if !r.IsValid() {
		return "", NewValidationError("RelationshipNature", s)
	}
	return r, nil
}

// AllRelationshipNatureValues 返回所有有效的关系类型值。
func AllRelationshipNatureValues() []RelationshipNature {
	return []RelationshipNature{
		RelationshipChildOf, RelationshipParentOf, RelationshipCanPrecede,
		RelationshipCanFollow, RelationshipRequires, RelationshipRequiredBy,
		RelationshipCanAlsoBe, RelationshipPeerOf, RelationshipMemberOf,
		RelationshipHasMember,
	}
}

// IsHierarchical 检查关系是否为层级关系（ChildOf, ParentOf, MemberOf, HasMember）。
func (r RelationshipNature) IsHierarchical() bool {
	switch r {
	case RelationshipChildOf, RelationshipParentOf, RelationshipMemberOf, RelationshipHasMember:
		return true
	default:
		return false
	}
}

// IsSequential 检查关系是否为顺序关系（CanPrecede, CanFollow）。
func (r RelationshipNature) IsSequential() bool {
	switch r {
	case RelationshipCanPrecede, RelationshipCanFollow:
		return true
	default:
		return false
	}
}

// IsDependency 检查关系是否为依赖关系（Requires, RequiredBy）。
func (r RelationshipNature) IsDependency() bool {
	switch r {
	case RelationshipRequires, RelationshipRequiredBy:
		return true
	default:
		return false
	}
}

// IsPeer 检查关系是否为对等关系（PeerOf, CanAlsoBe）。
func (r RelationshipNature) IsPeer() bool {
	switch r {
	case RelationshipPeerOf, RelationshipCanAlsoBe:
		return true
	default:
		return false
	}
}

// ConsequenceScope 表示后果的影响范围。
type ConsequenceScope string

const (
	ScopeConfidentiality ConsequenceScope = "Confidentiality"
	ScopeIntegrity       ConsequenceScope = "Integrity"
	ScopeAvailability    ConsequenceScope = "Availability"
	ScopeAccessControl   ConsequenceScope = "Access Control"
	ScopeAccountability  ConsequenceScope = "Accountability"
	ScopeAuthentication  ConsequenceScope = "Authentication"
	ScopeAuthorization   ConsequenceScope = "Authorization"
	ScopeNonRepudiation  ConsequenceScope = "Non-Repudiation"
)

// String 返回影响范围的字符串表示。
func (s ConsequenceScope) String() string { return string(s) }

// IsValid 检查影响范围是否为有效值。
func (s ConsequenceScope) IsValid() bool {
	switch s {
	case ScopeConfidentiality, ScopeIntegrity, ScopeAvailability,
		ScopeAccessControl, ScopeAccountability, ScopeAuthentication,
		ScopeAuthorization, ScopeNonRepudiation:
		return true
	default:
		return false
	}
}

// ParseConsequenceScope 从字符串解析影响范围。
func ParseConsequenceScope(s string) (ConsequenceScope, error) {
	sc := ConsequenceScope(s)
	if !sc.IsValid() {
		return "", NewValidationError("ConsequenceScope", s)
	}
	return sc, nil
}

// AllConsequenceScopeValues 返回所有有效的影响范围值。
func AllConsequenceScopeValues() []ConsequenceScope {
	return []ConsequenceScope{
		ScopeConfidentiality, ScopeIntegrity, ScopeAvailability,
		ScopeAccessControl, ScopeAccountability, ScopeAuthentication,
		ScopeAuthorization, ScopeNonRepudiation,
	}
}

// ConsequenceImpact 表示后果的影响严重程度。
type ConsequenceImpact string

const (
	ImpactHigh    ConsequenceImpact = "High"
	ImpactMedium  ConsequenceImpact = "Medium"
	ImpactLow     ConsequenceImpact = "Low"
	ImpactUnknown ConsequenceImpact = "Unknown"
)

// String 返回影响严重程度的字符串表示。
func (i ConsequenceImpact) String() string { return string(i) }

// IsValid 检查影响严重程度是否为有效值。
func (i ConsequenceImpact) IsValid() bool {
	switch i {
	case ImpactHigh, ImpactMedium, ImpactLow, ImpactUnknown:
		return true
	default:
		return false
	}
}

// ParseConsequenceImpact 从字符串解析影响严重程度。
func ParseConsequenceImpact(s string) (ConsequenceImpact, error) {
	imp := ConsequenceImpact(s)
	if !imp.IsValid() {
		return "", NewValidationError("ConsequenceImpact", s)
	}
	return imp, nil
}

// AllConsequenceImpactValues 返回所有有效的影响严重程度值。
func AllConsequenceImpactValues() []ConsequenceImpact {
	return []ConsequenceImpact{ImpactHigh, ImpactMedium, ImpactLow, ImpactUnknown}
}

// ImpactOrder 返回影响严重程度的排序权重。
// High=4, Medium=3, Low=2, Unknown=1, 未知=0
func (i ConsequenceImpact) ImpactOrder() int {
	switch i {
	case ImpactHigh:
		return 4
	case ImpactMedium:
		return 3
	case ImpactLow:
		return 2
	case ImpactUnknown:
		return 1
	default:
		return 0
	}
}

// ViewType 表示CWE视图的类型。
type ViewType string

const (
	// ViewTypeGraph 表示图类型视图，具有层次化的关系表示
	// 例如：CWE-1000 研究概念，CWE-699 软件开发
	ViewTypeGraph ViewType = "Graph"

	// ViewTypeExplicitSlice 表示显式切片视图，通过外部因素相关的扁平列表
	// 例如：CWE Top 25, OWASP Top Ten
	ViewTypeExplicitSlice ViewType = "Explicit Slice"

	// ViewTypeImplicitSlice 表示隐式切片视图，通过过滤器/属性定义的扁平列表
	// 例如：所有草稿状态的条目
	ViewTypeImplicitSlice ViewType = "Implicit Slice"
)

// String 返回视图类型的字符串表示。
func (v ViewType) String() string { return string(v) }

// IsValid 检查视图类型是否为有效值。
func (v ViewType) IsValid() bool {
	switch v {
	case ViewTypeGraph, ViewTypeExplicitSlice, ViewTypeImplicitSlice:
		return true
	default:
		return false
	}
}

// ParseViewType 从字符串解析视图类型。
func ParseViewType(s string) (ViewType, error) {
	vt := ViewType(s)
	if !vt.IsValid() {
		return "", NewValidationError("ViewType", s)
	}
	return vt, nil
}

// AllViewTypeValues 返回所有有效的视图类型值。
func AllViewTypeValues() []ViewType {
	return []ViewType{ViewTypeGraph, ViewTypeExplicitSlice, ViewTypeImplicitSlice}
}

// PlatformType 表示适用平台的类型。
type PlatformType string

const (
	PlatformLanguage     PlatformType = "Language"
	PlatformOperatingSystem PlatformType = "Operating System"
	PlatformArchitecture PlatformType = "Architecture"
	PlatformTechnology   PlatformType = "Technology"
)

// String 返回平台类型的字符串表示。
func (p PlatformType) String() string { return string(p) }

// IsValid 检查平台类型是否为有效值。
func (p PlatformType) IsValid() bool {
	switch p {
	case PlatformLanguage, PlatformOperatingSystem, PlatformArchitecture, PlatformTechnology:
		return true
	default:
		return false
	}
}

// ParsePlatformType 从字符串解析平台类型。
func ParsePlatformType(s string) (PlatformType, error) {
	pt := PlatformType(s)
	if !pt.IsValid() {
		return "", NewValidationError("PlatformType", s)
	}
	return pt, nil
}

// AllPlatformTypeValues 返回所有有效的平台类型值。
func AllPlatformTypeValues() []PlatformType {
	return []PlatformType{PlatformLanguage, PlatformOperatingSystem, PlatformArchitecture, PlatformTechnology}
}

// Prevalence 表示平台的使用普遍程度。
type Prevalence string

const (
	PrevalenceOften       Prevalence = "Often"
	PrevalenceSometimes   Prevalence = "Sometimes"
	PrevalenceRarely      Prevalence = "Rarely"
	PrevalenceUndetermined Prevalence = "Undetermined"
)

// String 返回普遍程度的字符串表示。
func (p Prevalence) String() string { return string(p) }

// IsValid 检查普遍程度是否为有效值。
func (p Prevalence) IsValid() bool {
	switch p {
	case PrevalenceOften, PrevalenceSometimes, PrevalenceRarely, PrevalenceUndetermined:
		return true
	default:
		return false
	}
}

// ParsePrevalence 从字符串解析普遍程度。
func ParsePrevalence(s string) (Prevalence, error) {
	p := Prevalence(s)
	if !p.IsValid() {
		return "", NewValidationError("Prevalence", s)
	}
	return p, nil
}

// AllPrevalenceValues 返回所有有效的普遍程度值。
func AllPrevalenceValues() []Prevalence {
	return []Prevalence{PrevalenceOften, PrevalenceSometimes, PrevalenceRarely, PrevalenceUndetermined}
}

// IntroductionPhase 表示弱点引入的阶段。
type IntroductionPhase string

const (
	PhaseArchitectureAndDesign IntroductionPhase = "Architecture and Design"
	PhaseImplementation        IntroductionPhase = "Implementation"
	PhaseBuildAndCompilation   IntroductionPhase = "Build and Compilation"
	PhaseOperation             IntroductionPhase = "Operation"
	PhaseSystemConfiguration   IntroductionPhase = "System Configuration"
	PhaseInstallation          IntroductionPhase = "Installation"
	PhasePolicy                IntroductionPhase = "Policy"
)

// String 返回引入阶段的字符串表示。
func (p IntroductionPhase) String() string { return string(p) }

// IsValid 检查引入阶段是否为有效值。
func (p IntroductionPhase) IsValid() bool {
	switch p {
	case PhaseArchitectureAndDesign, PhaseImplementation, PhaseBuildAndCompilation,
		PhaseOperation, PhaseSystemConfiguration, PhaseInstallation, PhasePolicy:
		return true
	default:
		return false
	}
}

// ParseIntroductionPhase 从字符串解析引入阶段。
func ParseIntroductionPhase(s string) (IntroductionPhase, error) {
	ip := IntroductionPhase(s)
	if !ip.IsValid() {
		return "", NewValidationError("IntroductionPhase", s)
	}
	return ip, nil
}

// AllIntroductionPhaseValues 返回所有有效的引入阶段值。
func AllIntroductionPhaseValues() []IntroductionPhase {
	return []IntroductionPhase{
		PhaseArchitectureAndDesign, PhaseImplementation, PhaseBuildAndCompilation,
		PhaseOperation, PhaseSystemConfiguration, PhaseInstallation, PhasePolicy,
	}
}

// MitigationPhase 表示缓解措施的阶段。
type MitigationPhase string

const (
	MitigationPhaseArchitectureAndDesign MitigationPhase = "Architecture and Design"
	MitigationPhaseBuildAndCompilation   MitigationPhase = "Build and Compilation"
	MitigationPhaseImplementation        MitigationPhase = "Implementation"
	MitigationPhaseOperation             MitigationPhase = "Operation"
	MitigationPhaseSystemConfiguration   MitigationPhase = "System Configuration"
	MitigationPhaseInstallation          MitigationPhase = "Installation"
	MitigationPhasePolicy                MitigationPhase = "Policy"
)

// String 返回缓解措施阶段的字符串表示。
func (m MitigationPhase) String() string { return string(m) }

// IsValid 检查缓解措施阶段是否为有效值。
func (m MitigationPhase) IsValid() bool {
	switch m {
	case MitigationPhaseArchitectureAndDesign, MitigationPhaseBuildAndCompilation,
		MitigationPhaseImplementation, MitigationPhaseOperation,
		MitigationPhaseSystemConfiguration, MitigationPhaseInstallation,
		MitigationPhasePolicy:
		return true
	default:
		return false
	}
}

// ParseMitigationPhase 从字符串解析缓解措施阶段。
func ParseMitigationPhase(s string) (MitigationPhase, error) {
	mp := MitigationPhase(s)
	if !mp.IsValid() {
		return "", NewValidationError("MitigationPhase", s)
	}
	return mp, nil
}

// AllMitigationPhaseValues 返回所有有效的缓解措施阶段值。
func AllMitigationPhaseValues() []MitigationPhase {
	return []MitigationPhase{
		MitigationPhaseArchitectureAndDesign, MitigationPhaseBuildAndCompilation,
		MitigationPhaseImplementation, MitigationPhaseOperation,
		MitigationPhaseSystemConfiguration, MitigationPhaseInstallation,
		MitigationPhasePolicy,
	}
}

// Effectiveness 表示缓解措施或检测方法的有效性。
type Effectiveness string

const (
	EffectivenessHigh          Effectiveness = "High"
	EffectivenessModerate      Effectiveness = "Moderate"
	EffectivenessLimited       Effectiveness = "Limited"
	EffectivenessDefenseInDepth Effectiveness = "Defense in Depth"
	EffectivenessSOARPartial   Effectiveness = "SOAR Partial"
	EffectivenessUnknown       Effectiveness = "Unknown"
)

// String 返回有效性的字符串表示。
func (e Effectiveness) String() string { return string(e) }

// IsValid 检查有效性是否为有效值。
func (e Effectiveness) IsValid() bool {
	switch e {
	case EffectivenessHigh, EffectivenessModerate, EffectivenessLimited,
		EffectivenessDefenseInDepth, EffectivenessSOARPartial, EffectivenessUnknown:
		return true
	default:
		return false
	}
}

// ParseEffectiveness 从字符串解析有效性。
func ParseEffectiveness(s string) (Effectiveness, error) {
	eff := Effectiveness(s)
	if !eff.IsValid() {
		return "", NewValidationError("Effectiveness", s)
	}
	return eff, nil
}

// AllEffectivenessValues 返回所有有效的有效性值。
func AllEffectivenessValues() []Effectiveness {
	return []Effectiveness{
		EffectivenessHigh, EffectivenessModerate, EffectivenessLimited,
		EffectivenessDefenseInDepth, EffectivenessSOARPartial, EffectivenessUnknown,
	}
}
