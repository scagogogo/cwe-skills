# 🌐 cwe api-version

查询 MITRE CWE REST API 的版本信息。

<Badge type="info" text="在线"/>

## 语法

```bash
cwe api-version [flags]
```

## 参数

本命令无位置参数。

## Flags

| Flag | 默认值 | 说明 |
| --- | --- | --- |
| `--base-url` | `https://cwe-api.mitre.org/api` | MITRE API 基础 URL |

::: tip 注意
`api-version` 使用的是 `Flags`（局部）而非 `PersistentFlags`，参数仅在 `cwe api-version` 自身后有效。
:::

## 示例

### text 输出

```bash
cwe api-version
```

```text
MITRE CWE API版本: 2.1
发布日期: 2024-xx-xx
版本名称: ...
```

### JSON 输出

```bash
cwe api-version -o json
```

```json
{
  "version": "2.1",
  "release_date": "2024-xx-xx",
  "name": "..."
}
```

### 指定自定义 API 地址

```bash
cwe api-version --base-url https://my-mirror.example.com/api
```

## 使用场景

- 调试 API 连通性与可达性。
- 确认 MITRE API 当前版本，以便对齐 [`show`](./show)、[`relations`](./relations) 等在线命令的行为。
- 使用自建镜像或代理时验证端点。

::: warning 连接失败
若返回“查询API版本失败”，请检查网络、代理设置及 `--base-url` 是否正确。在线命令均依赖此端点可达。
:::

## 下一步

- [show](./show) — 在线获取弱点详情。
- [relations](./relations) — 在线关系查询。
- [全局参数](./global-flags)

## 相关文档

- [SDK GetVersion](../sdk/api-get-version)
- [SDK API 客户端](../sdk/api-client)
