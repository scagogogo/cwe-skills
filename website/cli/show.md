# 🌐 cwe show

通过 MITRE CWE REST API 获取弱点（Weakness）的详细信息。

<Badge type="info" text="在线"/> 需要访问 MITRE API（默认 `https://cwe-api.mitre.org/api`）。

## 语法

```bash
cwe show [CWE-ID...] [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `CWE-ID...` | 是 | 一个或多个待查询的弱点 ID |

## Flags

| Flag | 默认值 | 说明 |
| --- | --- | --- |
| `--base-url` | `https://cwe-api.mitre.org/api` | MITRE API 基础 URL |
| `--timeout` | `30` | API 请求超时时间（秒） |

本命令的 flags 为 `PersistentFlags`，对其子命令 [`category`](./show-category)、[`view`](./show-view) 同样生效。

## 示例

### text 输出

```bash
cwe show CWE-79
```

```text
=== CWE-79 ===
  名称:     Improper Neutralization of Input During Web Page Generation
  抽象层级: Base
  状态:     Stable
  描述:     ...
  关系:     N 项
```

批量查询：

```bash
cwe show 79 89 352
```

### JSON 输出

```bash
cwe show CWE-79 -o json
```

```json
[
  {
    "cwe_id": "CWE-79",
    "detail": {
      "id": 79,
      "name": "Improper Neutralization of Input During Web Page Generation",
      "abstraction": "Base",
      "status": "Stable",
      "description": "..."
    }
  }
]
```

## 使用场景

- 在线查阅某个弱点的官方定义、描述与关系。
- 编写漏洞报告时引用权威描述。
- 验证某 CWE 编号是否真实存在（不存在时返回错误）。

::: tip 离线替代
若已下载 XML 目录，可用 [`registry get`](./registry-get) 离线获取同样信息且无需联网，速度更快。
:::

::: warning 网络与速率
MITRE API 有速率限制，批量查询大量 ID 时建议适当间隔，或改用离线 [registry](./registry) 命令。
:::

## 下一步

- [show category](./show-category) — 获取类别详情。
- [show view](./show-view) — 获取视图详情。
- [relations](./relations) — 查询关系。
- [registry get](./registry-get) — 离线替代。

## 相关文档

- [SDK API 客户端](../sdk/api-client)
- [SDK GetWeakness](../sdk/api-get-weakness)
