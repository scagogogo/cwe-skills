---
title: Status 状态枚举
outline: [2, 3]
---

# ✅ Status — 状态枚举

`Status` 表示 CWE 条目的成熟度与稳定性，反映该条目在 CWE 体系中的生命周期阶段。

## 🧬 类型定义

```go
type Status string
```

## 📋 全部取值

| 值 | 常量名 | 含义 |
| --- | --- | --- |
| `"Stable"` | `StatusStable` | 所有重要元素已验证，不太可能发生显著变化 |
| `"Usable"` | `StatusUsable` | 已经过深入审查，关键元素已验证 |
| `"Draft"` | `StatusDraft` | 所有重要元素已填写，可能仍有问题或空缺 |
| `"Incomplete"` | `StatusIncomplete` | 并非所有重要元素都已填写，无质量保证 |
| `"Obsolete"` | `StatusObsolete` | 仍然有效但不再相关，已被更新的实体取代 |
| `"Deprecated"` | `StatusDeprecated` | 已从 CWE 中移除，是重复或错误创建的 |

```go
const (
	StatusStable     Status = "Stable"
	StatusUsable     Status = "Usable"
	StatusDraft      Status = "Draft"
	StatusIncomplete Status = "Incomplete"
	StatusObsolete   Status = "Obsolete"
	StatusDeprecated Status = "Deprecated"
)
```

::: warning Obsolete 与 Deprecated 的区别
**Obsolete** 条目本身仍有效，只是被更新的条目取代、不再推荐使用；**Deprecated** 条目已从 CWE 中移除，属于重复或错误创建。做历史溯源时两者都可能遇到，但前者仍有参考价值，后者应直接弃用。
:::

## ✅ 通用方法（四件套）

| 方法 / 函数 | 签名 |
| --- | --- |
| `String` | `func (s Status) String() string` |
| `IsValid` | `func (s Status) IsValid() bool` |
| `ParseStatus` | `func ParseStatus(s string) (Status, error)` |
| `AllStatusValues` | `func AllStatusValues() []Status` |

```go
s, err := cweskills.ParseStatus("Stable")
fmt.Println(s, err)                            // Stable <nil>
fmt.Println(s.String())                        // Stable
fmt.Println(cweskills.Status("Foo").IsValid()) // false
fmt.Println(cweskills.AllStatusValues())       // [Stable Usable Draft Incomplete Obsolete Deprecated]
```

::: tip 严格匹配
`ParseStatus` 区分大小写：`ParseStatus("stable")` 返回错误，必须传入 `"Stable"`。处理外部输入时需先做规范化（如 `strings.Title` 或自建映射）。
:::

## 🔄 典型用法

```go
// 过滤出可用于生产映射的稳定条目
stable, _ := registry.FindByStatus(cweskills.StatusStable)
usable, _ := registry.FindByStatus(cweskills.StatusUsable)
fmt.Println("可用条目:", len(stable)+len(usable))

// 排除已废弃条目，避免误导
for _, c := range all {
	if c.Status == cweskills.StatusDeprecated || c.Status == cweskills.StatusObsolete {
		continue
	}
	process(c)
}
```

## 💻 CLI 对应命令

```bash
cwe enum status
```

输出全部合法取值，详见 [CLI enum status](../cli/enum-status)。

## 🔗 相关链接

- SDK 视角：[Status 状态枚举](../sdk/enum-status)
- 按状态筛选：[`FindByStatus`](../sdk/find-by-status)
- 按状态统计：[`CountByStatus`](../sdk/count-by-status)
- 源文件：[`enums.go`](https://github.com/scagogogo/cwe-skills/blob/main/enums.go)
