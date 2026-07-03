---
title: 实战 — CI 流水线集成
outline: [2, 3]
---

# ⚙️ 实战 — CI 流水线集成

在 GitHub Actions 中用 `cwe` CLI + JSON + jq 检查每次提交是否引入了 CWE Top 25 弱点，命中则让流水线失败。

<Badge type="tip" text="CI/CD 实战"/>
<Badge type="info" text="内置列表"/>

---

## 🎬 场景

团队想在 CI 里加一道门禁：扫描器输出的 CWE 列表里若有任何 Top 25 项，PR 检查标红，阻止合入。

---

## 📋 前置准备

- GitHub 仓库已启用 Actions
- 一个能输出 CWE 列表的扫描步骤（SAST/SBOM 等），结果写到 `scan_cwes.txt`，每行一个 `CWE-NNN`

---

## 💻 GitHub Actions 工作流

`.github/workflows/cwe-gate.yml`：

```yaml
name: CWE Top 25 门禁

on:
  pull_request:
    paths:
      - '**/*.go'
      - '**/*.py'
      - '**/*.js'

jobs:
  cwe-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: 安装 cwe CLI
        run: |
          curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_linux_x86_64.tar.gz | tar xz
          sudo mv cwe /usr/local/bin/
          cwe version

      - name: 运行扫描器（示例占位）
        run: |
          # 实际替换为你的 SAST/SBOM 扫描命令
          # 这里用模拟输出演示
          echo "CWE-79" > scan_cwes.txt
          echo "CWE-89" >> scan_cwes.txt
          echo "CWE-778" >> scan_cwes.txt

      - name: 检查是否引入 Top 25
        id: check
        run: |
          CWES=$(tr '\n' ' ' < scan_cwes.txt)
          cwe wellknown check $CWES -o json > wellknown_result.json
          # 提取命中 Top 25 的 CWE
          HITS=$(jq -r '[.[] | select(.in_list | index("Top 25")) | .cwe_id] | join(",")' wellknown_result.json)
          echo "hits=$HITS" >> $GITHUB_OUTPUT
          if [[ -n "$HITS" ]]; then
            echo "::error::检测到 Top 25 弱点: $HITS"
            exit 1
          fi

      - name: 上传结果
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: cwe-wellknown-result
          path: wellknown_result.json
```

---

## ▶️ 运行步骤

```bash
# 本地预演（不依赖 Actions）
echo -e "CWE-79\nCWE-89\nCWE-778" > scan_cwes.txt
CWES=$(tr '\n' ' ' < scan_cwes.txt)
cwe wellknown check $CWES -o json | jq '[.[] | select(.in_list | index("Top 25")) | .cwe_id]'
# 输出: ["CWE-79","CWE-89"]

# 提交 PR 触发流水线
git push origin feature-branch
```

---

## 📤 输出示例

命中 Top 25 时，Actions 日志：

```text
Run CWES=$(tr '\n' ' ' < scan_cwes.txt)
Run cwe wellknown check $CWES -o json > wellknown_result.json
Run HITS=$(jq -r '...')
Run echo "hits=$HITS" >> $GITHUB_OUTPUT
Run if [[ -n "$HITS" ]]; then
Run   echo "::error::检测到 Top 25 弱点: CWE-79,CWE-89"
Run   exit 1
Error: 检测到 Top 25 弱点: CWE-79,CWE-89
##[error]Process completed with exit code 1.
```

未命中时（`HITS` 为空）检查通过，PR 可合入。

---

## 🧩 扩展思路

- **分级门禁**：Top 25 阻断合入，OWASP/SANS 仅告警不阻断——用 `jq` 分别提取并设置不同退出码。
- **白名单**：对已知接受的 CWE 维护一份白名单，`jq` 过滤后再判断。
- **PR 评论**：用 `gh pr comment` 把命中详情发到 PR 评论，比仅标红更直观。
- **历史趋势**：把每次 `wellknown_result.json` 存档，跟踪 Top 25 命中数随时间的变化。

::: tip ::error 高亮
GitHub Actions 的 `::error::` 注解会让对应步骤标红并在 PR Checks 里醒目显示，适合做门禁。
:::

---

## 📖 相关文档

- [技能 03 — 知名列表](../skills/03-well-known-lists)
- [CLI: wellknown check](../cli/wellknown-check) · [输出格式](../guide/output-format)
- [返回示例总览](./)
