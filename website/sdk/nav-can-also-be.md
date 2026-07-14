---
title: 可作为导航
outline: [2, 3]
---

# 🧭 可作为导航

`CanAlsoBe` 返回与当前弱点**语义可互换**或**可作为等价表现**的弱点。MITRE 用此关系标注「同一个根因在不同语境下的不同表述」。

## 📐 方法签名

```go
func (n *Navigator) CanAlsoBe(id int) []*CWE
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `id` | `int` | 弱点 ID |
| 返回 | `[]*CWE` | 当前弱点「也可表现为」的弱点列表 |

::: tip 方向语义
设 A 的关系里写了 `CanAlsoBe B`，意为「A 也可表现为 B」。`CanAlsoBe(A)` 返回 B。
该关系在 MITRE 数据中通常**双向**声明，但不保证——单向声明时 `CanAlsoBe(B)` 未必返回 A。
:::

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	r := cweskills.NewRegistry()
	// 79 也可表现为 80
	c79 := cweskills.NewCWE(79, "XSS")
	c79.Relationships = []cweskills.Relationship{
		{CWEID: 80, Nature: cweskills.RelationshipCanAlsoBe},
	}
	r.Register(c79)
	r.Register(cweskills.NewCWE(80, "Buffer Overflow"))

	r.BuildIndexes()
	nav := cweskills.NewNavigator(r)

	for _, c := range nav.CanAlsoBe(79) {
		fmt.Println("79 也可表现为:", c.CWEID()) // CWE-80
	}
}
```

## 🆚 与 Peers 的关系

| 维度 | `CanAlsoBe` | `Peers` |
| --- | --- | --- |
| 数据范围 | 仅 `CanAlsoBe` 边 | `PeerOf` + `CanAlsoBe` |
| 语义焦点 | 语义可互换 | 泛对等 |
| 适用 | 找等价表述 | 找平级关联 |

`Peers` 是 `CanAlsoBe` 的**超集**（额外含 `PeerOf`）。只关心可互换表述时用 `CanAlsoBe` 更精准。

## ⚠️ 注意事项

::: warning 不保证双向
若仅 A 声明 `CanAlsoBe B`，`CanAlsoBe(B)` 不会自动返回 A——除非 B 也声明了对 A 的 `CanAlsoBe`。需要双向互查时建议两侧都遍历。
:::

## 🔗 相关链接

- 泛对等查询：[兄弟与对等](./nav-siblings-peers)
- 关系类型枚举：[RelationshipNature](./enum-relationship-nature)
- 本地关系读取：[CWE 关系获取方法](./cwe-relationship-methods)
- 源文件：[`navigator.go`](https://github.com/scagogogo/cwe-skills/blob/main/navigator.go)
