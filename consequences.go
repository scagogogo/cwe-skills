package cwe

// Consequence 表示CWE条目可能造成的后果。
//
// 每个后果包含影响范围、影响严重程度、利用可能性和备注信息。
// 一个CWE条目可以有多种不同的后果，每种后果可能影响不同的范围。
type Consequence struct {
	// Scopes 受影响的安全范围列表，如机密性、完整性、可用性等
	Scopes []ConsequenceScope `json:"scopes" xml:"Scopes>Scope"`
	// Impacts 影响严重程度列表，如高、中、低
	Impacts []ConsequenceImpact `json:"impacts,omitempty" xml:"Impacts>Impact,omitempty"`
	// Likelihood 此后果被利用的可能性
	Likelihood LikelihoodOfExploit `json:"likelihood,omitempty" xml:"Likelihood,omitempty"`
	// Note 关于后果的补充说明
	Note string `json:"note,omitempty" xml:"Note,omitempty"`
}

// HasScope 检查后果是否包含指定的安全范围。
//
// 参数：
//   - scope: 需要检查的安全范围
//
// 返回值：
//   - bool: 如果后果中包含该安全范围返回true，否则返回false
func (c *Consequence) HasScope(scope ConsequenceScope) bool {
	for _, s := range c.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// HasImpact 检查后果是否包含指定的影响严重程度。
//
// 参数：
//   - impact: 需要检查的影响严重程度
//
// 返回值：
//   - bool: 如果后果中包含该影响严重程度返回true，否则返回false
func (c *Consequence) HasImpact(impact ConsequenceImpact) bool {
	for _, i := range c.Impacts {
		if i == impact {
			return true
		}
	}
	return false
}

// MaxImpact 返回后果中最高的影响严重程度。
//
// 影响严重程度的排序为：High > Medium > Low > Unknown。
// 如果Impacts列表为空，返回ImpactUnknown。
//
// 返回值：
//   - ConsequenceImpact: 最高的影响严重程度
func (c *Consequence) MaxImpact() ConsequenceImpact {
	if len(c.Impacts) == 0 {
		return ImpactUnknown
	}

	max := c.Impacts[0]
	for _, impact := range c.Impacts[1:] {
		if impact.ImpactOrder() > max.ImpactOrder() {
			max = impact
		}
	}
	return max
}

// Validate 验证后果的有效性。
//
// 检查条件：
//   - Scopes列表至少包含一个范围
//
// 返回值：
//   - error: 如果验证失败返回ValidationError，否则返回nil
func (c *Consequence) Validate() error {
	if len(c.Scopes) == 0 {
		return NewValidationError("Scopes", "空列表")
	}
	return nil
}