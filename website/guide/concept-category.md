---
title: 类别 (Category)
outline: [2, 3]
---

# 📂 类别 (Category)

**类别（Category）**是 CWE 提供的一种**非层级分组**机制：它把若干 CWE 按「共同主题」聚到一起，但成员之间**不一定有父子关系**。一个 CWE 可以同时属于多个类别。

::: tip Category vs View
两者都是「分组」，但侧重点不同：
- **Category**：按**主题/特征**分组（如「内存分配不当」），扁平、非层级。
- **View**：按**视角/受众**分组（如「软件开发视角」），可能是层级图（Graph）或切片。
详见 [视图 (View)](./concept-view)。
:::

---

## 📦 Category 结构

```go
type Category struct {
    ID              int           // 类别数字标识符
    Name            string        // 类别名称
    Status          Status        // 状态
    Description     string        // 描述
    Members         []int         // 成员 CWE ID 列表
    Relationships   []Relationship// 关系（含 MemberOf/Has_Member）
    References      []Reference   // 参考文献
    ContentHistory  *ContentHistory
}
```

```go
cat, ok := registry.GetCategory(789)
if ok {
    fmt.Println(cat.Name, "成员数:", len(cat.Members))
}
```

---

## 🔗 成员关系

类别与成员之间通过 `MemberOf` / `Has_Member` 关系相连（属于[层级关系](./concept-relationship)的一种，但不构成父子树）：

| 关系 | 方向 | 含义 |
|------|------|------|
| `MemberOf` | 弱点 A MemberOf 类别 B | A 是 B 的成员 |
| `Has_Member` | 类别 B Has_Member 弱点 A | B 包含成员 A |

```text
  [类别: 内存分配不当]
     Has_Member   Has_Member   Has_Member
        │             │            │
     [CWE-x]      [CWE-y]      [CWE-z]   ← 成员间无父子关系
```

::: info MemberOf/Has_Member 是层级关系
在 `RelationshipNature` 分类里，`MemberOf`/`Has_Member` 属于 `IsHierarchical()` 为 true 的层级关系，但它们表达的是「归属」而非「父子继承」。这就是为什么 [关系类型](./concept-relationship) 把它们归入层级类。
:::

---

## 🔍 查询类别成员

```go
// 取某类别的全部成员 CWE ID
memberIDs := registry.GetCategoryMembers(789)

// 取某条目所属的全部类别/视图
memberOfIDs := registry.GetMemberOfIDs(79)
```

```go
// 导航器视角
nav := cweskills.NewNavigator(registry)
// 注意：MemberOf/Has_Member 通过 registry 的成员索引查询
// nav 主要面向父子/对等/链式/依赖导航
```

---

## 💻 CLI 操作

```bash
# 列出所有类别
cwe registry list-categories --xml cwec_v4.15.xml

# 查看某类别的成员
cwe registry category-members 789 --xml cwec_v4.15.xml

# 查看某 CWE 属于哪些类别/视图
cwe registry member-of CWE-79 --xml cwec_v4.15.xml
```

::: warning 离线专属
类别成员查询依赖 XML 目录里的 `MemberOf`/`Has_Member` 关系，**在线 API 通常不返回**。需要时走 [离线 XML](../sdk/xml-parser)。
:::

---

## 🎯 用途

### 1. 主题化归类

类别把分散的 CWE 按主题聚拢，便于专题分析。例如「输入校验相关」「密码学相关」「并发相关」等主题类别，能快速定位一类问题。

### 2. 报告分组

安全报告按类别分组展示，比平铺一长串 CWE ID 更易读：

```go
for _, cat := range registry.GetAllCategories() {
    members := registry.GetCategoryMembers(cat.ID)
    fmt.Printf("## %s（%d 项）\n", cat.Name, len(members))
    for _, id := range members {
        if w, ok := registry.Get(id); ok {
            fmt.Printf("- %s %s\n", w.CWEID(), w.Name)
        }
    }
}
```

### 3. 多维标签

一个 CWE 可属多个类别，相当于给它打了多个「主题标签」。结合 [视图](./concept-view)、[抽象层级](./concept-abstraction) 等维度，可做多维交叉分析。

```go
// CWE-79 属于哪些类别？
for _, catID := range registry.GetMemberOfIDs(79) {
    if cat, ok := registry.GetCategory(catID); ok {
        fmt.Println(cat.Name)
    }
}
```

---

## ⚠️ 注意事项

::: warning Category 不建树
类别是扁平分组，**不要用 `BuildTree` 建树**。要建树请用 [视图](./concept-view) 的 Graph 类型或指定根 CWE。
:::

::: info Category ID 与 Weakness ID 共享编号空间
CWE 的编号空间是统一的，类别、视图、弱点、复合元素都从同一池子里取 ID。`registry.Get(id)` 只查弱点集合；查类别用 `registry.GetCategory(id)`，视图用 `GetView`，复合元素用 `GetCompoundElement`。
:::

---

## 📖 相关文档

- [视图 (View)](./concept-view)
- [关系类型](./concept-relationship)
- [注册表 API](../sdk/registry)
- [CWE 是什么](./concept-cwe)
