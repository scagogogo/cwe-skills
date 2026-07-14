---
title: Introduction 引入方式
outline: [2, 3]
---

# 🚪 Introduction — 弱点引入方式

`Introduction` 描述弱点是在哪个生命周期阶段被引入到软件中的，是 `CWE.ModesOfIntroduction` 字段的元素类型。

## 📋 结构体定义

```go
type Introduction struct {
    Phase       IntroductionPhase `json:"phase" xml:"Phase"`
    Description string            `json:"description,omitempty" xml:"Description,omitempty"`
}
```

## 📝 字段说明

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `Phase` | `IntroductionPhase` | 引入阶段 |
| `Description` | `string` | 引入方式的描述 |

::: tip IntroductionPhase 取值
枚举定义在 `enums.go`，共 7 个阶段：
- `Architecture and Design`（架构与设计）
- `Implementation`（实现）
- `Build and Compilation`（构建与编译）
- `Operation`（运行）
- `System Configuration`（系统配置）
- `Installation`（安装）
- `Policy`（策略）

配套 `ParseIntroductionPhase` / `AllIntroductionPhaseValues`。
:::

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	cwe := cweskills.NewCWE(79, "XSS")
	cwe.ModesOfIntroduction = []cweskills.Introduction{
		{
			Phase:       cweskills.PhaseArchitectureAndDesign,
			Description: "未在架构层规定输出编码策略。",
		},
		{
			Phase:       cweskills.PhaseImplementation,
			Description: "直接拼接用户输入到 HTML 输出。",
		},
	}

	for _, intro := range cwe.ModesOfIntroduction {
		fmt.Printf("[%s] %s\n", intro.Phase, intro.Description)
	}
	// [Architecture and Design] 未在架构层规定输出编码策略。
	// [Implementation] 直接拼接用户输入到 HTML 输出。
}
```

## 🎯 典型用途

<Badge type="tip" text="SDLC 对齐" /> 把弱点映射到软件开发生命周期阶段，明确在哪一环修复
<Badge type="info" text="培训" /> 按阶段组织安全培训内容
<Badge type="warning" text="分类" /> 区分「设计期可避免」vs「实现期引入」的弱点

## ⚠️ 注意事项

::: warning Phase 是必填字段
`Phase` 的 JSON/XML 标签都无 `omitempty`，序列化时一定会输出。这与 `Mitigation.Phase`（可空）的约定不同。
:::

::: details 与 MitigationPhase 的区别
两者取值集合高度重叠（都有 Architecture and Design、Implementation 等），但语义相反：
- `IntroductionPhase`：弱点**何时被引入**
- `MitigationPhase`：缓解措施**何时施加**

一个在「实现期引入」的弱点，仍可在「运行期」通过配置缓解。
:::

## 🔗 相关链接

- 宿主字段：`CWE.ModesOfIntroduction`，见 [CWE 弱点](./cwe-struct)
- 阶段枚举：`IntroductionPhase`（见 `enums.go`）
- 对应的缓解阶段：`MitigationPhase`，见 [Mitigation](./mitigation)
- 模型概览：[model.go](./model)
- 源文件：[`model.go`](https://github.com/scagogogo/cwe-skills/blob/main/model.go)
