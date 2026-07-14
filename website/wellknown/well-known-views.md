---
title: 知名视图
outline: [2, 3]
---

# 🔍 知名视图（Well-Known Views）

CWE 通过「视图（View）」组织条目。MITRE 维护若干官方视图，按不同维度（研究、开发、硬件、横截面、综合字典）裁剪 CWE 集合。SDK 在 `wellknown_ids.go` 中内置 5 个知名视图 ID 常量，配合 `IsInWellKnownView` 判断。

## 📊 背景来源

视图（View）是 CWE 体系的核心组织单元，定义在 CWE XML 的 `<View>` 元素中。每个视图：

- 拥有独立 ID（如 1000、699）
- 通过 `Membership`/`Has_Member` 等关系筛选出一组相关 CWE
- 从特定视角组织条目，便于不同受众检索

MITRE 维护的官方知名视图涵盖研究、软件开发、硬件设计等领域，是 CWE 数据集的「目录索引」。

## 📋 知名视图常量

| 常量 | 视图 ID | 视图名称 | 用途 |
| --- | ---: | --- | --- |
| `CWEViewResearchConcepts` | 1000 | 研究概念视图 | 按抽象概念层次组织所有 CWE，适合研究/分类 |
| `CWEViewDevelopmentConcepts` | 699 | 软件开发视图 | 按软件开发活动组织，适合开发者检索 |
| `CWEViewHardwareDesign` | 1199 | 硬件设计视图 | 按硬件设计活动组织，面向硬件安全 |
| `CWEViewCWECrossSection` | 888 | CWE 横截面视图 | 提供 CWE 条目的横截面视角 |
| `CWEViewComprehensiveDictionary` | 1400 | 综合 CWE 字典 | 包含所有 CWE 条目的综合视图 |

<Badge type="tip" text="常量即 ID" /> 这些常量值就是 MITRE 官方视图 ID，可直接传入任何接受 `viewID` 的 SDK 函数。

## ✅ SDK 用法示例

```go
package main

import (
    "fmt"
    "github.com/scagogogo/cwe-skills"
)

func main() {
    // 判断是否为知名视图
    fmt.Println(cweskills.IsInWellKnownView(cweskills.CWEViewResearchConcepts)) // true
    fmt.Println(cweskills.IsInWellKnownView(699))  // true（软件开发视图）
    fmt.Println(cweskills.IsInWellKnownView(1234)) // false（非知名视图）

    // 用常量而非魔法数字，代码更自文档化
    if cweskills.IsInWellKnownView(cweskills.CWEViewHardwareDesign) {
        fmt.Println("当前涉及硬件设计视图，应启用硬件相关检查")
    }

    // 遍历所有知名视图常量
    views := []struct {
        name string
        id   int
    }{
        {"研究概念", cweskills.CWEViewResearchConcepts},
        {"软件开发", cweskills.CWEViewDevelopmentConcepts},
        {"硬件设计", cweskills.CWEViewHardwareDesign},
        {"横截面", cweskills.CWEViewCWECrossSection},
        {"综合字典", cweskills.CWEViewComprehensiveDictionary},
    }
    for _, v := range views {
        fmt.Printf("%s 视图 ID = %d, 知名 = %v\n", v.name, v.id, cweskills.IsInWellKnownView(v.id))
    }
}
```

::: warning 常量集中定义
这些视图 ID 是 MITRE 官方编号，**不会变动**，因此 SDK 用常量集中管理。请始终使用常量（如 `CWEViewResearchConcepts`）而非裸数字 `1000`，以保证代码可读性与可维护性。
:::

## 🖥️ CLI 用法

视图常量主要用于 SDK 编程。CLI 层面，视图相关的查询通过 `registry` 与 `tree-view` 等命令消费视图 ID：

```bash
# 列出注册表中所有视图（含知名视图）
cwe registry list-views

# 按视图构建树（此处 699 = 软件开发视图）
cwe tree view 699

# 按视图导出条目
cwe registry export --view 699 -o json
```

```text
视图 CWE-699（软件开发视图）构建中...
  ├─ CWE-20  Improper Input Validation
  ├─ CWE-79  Cross-site Scripting (XSS)
  ...
```

::: tip CLI 中的视图 ID
`cwe tree view`、`cwe show view` 等命令接受的 `viewID` 即为这些常量值。`wellknown` 子命令本身不直接处理视图，视图常量在 SDK 与 registry/tree 命令中发挥作用。
:::

## 🎯 使用场景

<Badge type="tip" text="按视角检索" /> 用软件开发视图（699）为开发者定制 CWE 子集
<Badge type="info" text="硬件安全" /> 用硬件设计视图（1199）聚焦硬件相关弱点
<Badge type="warning" text="全量分析" /> 用综合字典视图（1400）获取全部 CWE，做全集统计
<Badge type="info" text="研究分类" /> 用研究概念视图（1000）按抽象层次研究弱点分类学
<Badge type="info" text="视图有效性校验" /> 调用 `IsInWellKnownView` 校验外部传入的视图 ID 是否为官方知名视图

## 🔗 相关链接

- [知名列表总览](./overview)
- [CWE Top 25 (2024)](./cwe-top-25)
- 概念：[视图概念](../guide/concept-view)
- SDK：[IsInWellKnownView](../sdk/is-in-wellknown-view)、[wellknown-ids 概览](../sdk/wellknown-ids)
- CLI：[registry list-views](../cli/registry-list-views)、[tree view](../cli/tree-view)、[show view](../cli/show-view)
- 源文件：[`wellknown_ids.go`](https://github.com/scagogogo/cwe-skills/blob/main/wellknown_ids.go)
