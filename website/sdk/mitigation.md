---
title: Mitigation 缓解措施
outline: [2, 3]
---

# 🛡️ Mitigation — 缓解措施

`Mitigation` 描述如何减少或消除某个弱点的风险，是 `CWE.PotentialMitigations` 字段的元素类型。每条缓解措施关联一个生命周期阶段、可选策略、详细描述与有效性评估。

## 📋 结构体定义

```go
type Mitigation struct {
    Phase         MitigationPhase `json:"phase,omitempty" xml:"Phase,omitempty"`
    Strategy      string          `json:"strategy,omitempty" xml:"Strategy,omitempty"`
    Description   string          `json:"description" xml:"Description"`
    Effectiveness Effectiveness   `json:"effectiveness,omitempty" xml:"Effectiveness,omitempty"`
}
```

## 📝 字段说明

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `Phase` | `MitigationPhase` | 缓解措施适用的生命周期阶段 |
| `Strategy` | `string` | 缓解策略名称（如 "Input Validation"） |
| `Description` | `string` | 缓解措施的详细描述 |
| `Effectiveness` | `Effectiveness` | 该措施的有效性评估 |

::: tip 两个枚举
- `MitigationPhase` 取值见 `enums.go`，包括 `Architecture and Design`、`Implementation`、`Operation`、`Build and Compilation`、`System Configuration`、`Installation`、`Policy`，配套 `ParseMitigationPhase` / `AllMitigationPhaseValues`。
- `Effectiveness` 取值包括 `High`、`Moderate`、`Limited`、`Defense in Depth`、`SOAR Partial`、`Unknown`，配套 `ParseEffectiveness` / `AllEffectivenessValues`。
:::

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	cwe := cweskills.NewCWE(79, "XSS")
	cwe.PotentialMitigations = []cweskills.Mitigation{
		{
			Phase:         cweskills.MitigationPhaseImplementation,
			Strategy:      "Input Validation",
			Description:   "对所有用户输入进行严格校验与输出编码。",
			Effectiveness: cweskills.EffectivenessHigh,
		},
		{
			Phase:         cweskills.MitigationPhaseArchitectureAndDesign,
			Strategy:      "Output Encoding",
			Description:   "使用上下文感知的输出编码库。",
			Effectiveness: cweskills.EffectivenessDefenseInDepth,
		},
	}

	for i, m := range cwe.PotentialMitigations {
		fmt.Printf("#%d [%s] %s -> %s\n",
			i+1, m.Phase, m.Strategy, m.Effectiveness)
	}
	// #1 [Implementation] Input Validation -> High
	// #2 [Architecture and Design] Output Encoding -> Defense in Depth
}
```

## 🎯 典型用途

<Badge type="tip" text="报告" /> 在安全报告里输出「如何修复」章节
<Badge type="info" text="过滤" /> 按 `Phase` 筛选设计期可落地的措施
<Badge type="warning" text="排序" /> 按 `Effectiveness` 优先展示高效措施

## ⚠️ 注意事项

::: warning Strategy 是自由文本
`Strategy` 没有对应枚举，是 MITRE 提供的自由文本字符串（如 `"Input Validation"`、`"Output Encoding"`）。不要假设它取自固定集合，做匹配时按字符串相等处理。
:::

## 🔗 相关链接

- 宿主字段：`CWE.PotentialMitigations`，见 [CWE 弱点](./cwe-struct)
- 阶段枚举：`MitigationPhase`（见 `enums.go`）
- 有效性枚举：`Effectiveness`（见 `enums.go`）
- 模型概览：[model.go](./model)
- 源文件：[`model.go`](https://github.com/scagogogo/cwe-skills/blob/main/model.go)
