# 文案与 import 风格统一为 cwe-skills 标准计划

> **For agentic workers:** REQUIRED SUB-SKILL: `superpowers:subagent-driven-development`
> Steps use checkbox (`- [ ]`) syntax.

**Goal:** 把全仓 Go import 示例的显式别名统一为「无别名」标准风格（`"github.com/scagogogo/cwe-skills"` + `cweskills.` 调用），与上一轮已统一于 cmd 源码的方向对齐；同时修复 1 处 `scagogogo` 拼写错误。品牌名 `CWE Skills`（带空格、首字母大写）是项目有意的产品名，**保留不动**；本计划只统一 import 别名风格与拼写。

**Architecture:** 仓库命名体系分三层：品牌名 `CWE Skills`（人类可读文案，README 标题 / 文档站 / 叙述性正文，保留）/ 技术名 `cwe-skills`（连字符，仓库路径 / module / 二进制 / import 路径 / dist 资产，保留）/ 包名 `cweskills`（去连字符，package 声明 / `cweskills.` 调用前缀，保留）。因 module 路径末段 `cwe-skills` 去连字符即等于包名 `cweskills`，Go 默认推导出的包名就是 `cweskills`，故 import `"github.com/scagogogo/cwe-skills"` 无需任何显式别名即可用 `cweskills.` 调用——这正是上一轮已统一于 cmd 源码、README、docs/skills 的标准风格。但 website/sdk/ 文档示例存在两派非标准写法：A 派（30 文件）`cwe "..."` 别名 + `cwe.` 调用；B 派（多文件）`cweskills "..."` 冗余别名。本计划把两派都收敛到无别名标准，并修复 `website/sdk/new-api-client.md` 的 `scaggogo` 笔误。

**Tech Stack:** Markdown 文档（website/sdk 示例）、Go 1.25.0（SKILLS.md/zh import 示例）、sed/精确 Edit 批量替换

**Risks:**
- A 派 30 文件里 `cwe.X` 替换可能误伤局部变量 → 缓解：已用 `grep -E '\bcwe\.[a-z]'` 与 `grep 'cwe :='` 验证 A 派文件中**无** `cwe` 局部变量、所有 `cwe.X` 后均为大写导出符号（`cwe.NewAPIClient`/`cwe.APIError`/`cwe.WithAPIBaseURL`/`cwe.Relationships` 等包前缀），替换安全。
- 替换后 import 行若残留多余空格 → 缓解：用精确字符串替换 `cwe "github.com/scagogogo/cwe-skills"` → `"github.com/scagogogo/cwe-skills"`，保留行首 tab/空格缩进。
- 文档示例不参与编译，无法用 `go build` 验证 → 缓解：用 grep 验证「无残留别名 + 无残留 `cwe.` 包前缀调用」作为正确性证据，并人工抽查关键文件。
- B 派文件删别名后调用前缀 `cweskills.` 是否仍合法 → 缓解：默认包名就是 `cweskills`，删冗余别名后 `cweskills.X` 仍合法，调用点零改动。

---

### 调研结论（import 风格现状）

| 区域 | 当前风格 | 是否标准 | 处理 |
|---|---|---|---|
| cmd 源码（上一轮已统一） | 无别名 `"..."` + `cweskills.` | ✅ 标准 | 不动 |
| README.md / README.zh.md | 无别名 `"..."` + `cweskills.` | ✅ 标准 | 不动 |
| docs/skills/*.md | 无别名 `cweskills.` 调用 | ✅ 标准 | 不动 |
| website/sdk/ A 派（30 文件） | `cwe "..."` 别名 + `cwe.` 调用 | ❌ 非标准 | Task 1 改 |
| website/sdk/ B 派（多文件） | `cweskills "..."` 冗余别名 | ❌ 非标准 | Task 2 改 |
| SKILLS.md / SKILLS.zh.md:119 | `cweskills "..."` 冗余别名 | ❌ 非标准 | Task 3 改 |
| website/sdk/new-api-client.md:91 | `scaggogo`（拼写错误，少 o） | ❌ bug | Task 3 改 |
| 品牌名 `CWE Skills`（README/website/config.ts 叙述） | 带空格首字母大写 | ✅ 产品名 | 保留不动 |

**A 派文件清单（30 个，import 行 `cwe "github.com/scagogogo/cwe-skills"`，全部位于 `website/sdk/`）：**
api-ancestors-descendants、api-error、api-get-category、api-get-cwes、api-get-view、api-get-version、api-client、api-parents-children、api-response、basic-fetcher、cwe-error、cwe-not-found-error、data-fetcher、http-client、http-client-option、http-methods、http-retry、invalid-cwe-id-error、multiple-fetcher、new-api-client（含拼写错误，Task 3 单独处理）、new-xml-parser、parse-error、rate-limit-error、rate-limiter、rate-limiter-api、relationship-error、tree-fetcher、validation-error、xml-parse、xml-parser。

**B 派文件清单（97 个，import 行 `cweskills "github.com/scagogogo/cwe-skills"`，全部位于 `website/sdk/`）：**
alternate-term、build-forest、build-indexes、build-tree、build-view-tree、category、compare-cwe-ids、compute-statistics、compound-element、consequence、content-history、count-by-abstraction、count-by-likelihood、count-by-scope、count-by-status、cwe-relationship-methods、cwe-struct、cwe-top-25、cwe-type-methods、cwe-utils、deduplicate、demonstrative-example、enum-abstraction、enum-consequence-impact、enum-consequence-scope、enum-likelihood、enum-platform-type、enum-relationship-nature、enum-status、enum-structure、enum-view-type、export-csv、extract-cwe-ids、extract-first-cwe-id、filter、filter-option、find-by-abstraction、find-by-consequence-scope、find-by-id、find-by-keyword、find-by-likelihood、find-by-status、find-by-structure、find-chains-composites、find-top-level、format-cwe-id、format-cwe-id-from-int、get-owasp-categories、get-owasp-category、group-by、introduction、is-cwe-id、is-in-owasp-top-10、is-in-sans-top-25、is-in-top-25、is-in-wellknown-view、marshal-csv、marshal-json、marshal-json-list、marshal-xml、mitigation、nav-ancestor-related、nav-ancestors-descendants、nav-can-also-be、nav-chain-composite、nav-parents-children、nav-precede-follow、nav-relationship-depth、nav-requires、nav-shortest-path、nav-siblings-peers、navigator、observed-example、owasp-top-10、overview、parse-cwe-id、reference、registry、registry-json、registry-operations、relationship-indexes、sans-top-25、search、serializer、sort、stats、tree、tree-depth-count、tree-leaf-nodes、tree-node、tree-path、tree-walk、validate-cwe-id、view、wellknown-ids。（以 grep 实际清单为准，执行用 `grep -rln` 动态取）

---

### Task 1: 统一 website/sdk A 派 29 文件为无别名 — 删 `cwe` 别名并把 `cwe.X` 调用改为 `cweskills.X`

**Depends on:** None
**Files:**
- Modify: `website/sdk/` 下 A 派 30 文件中除 `new-api-client.md` 外的 29 个（new-api-client.md 含拼写错误，Task 3 单独处理）：api-ancestors-descendants、api-error、api-get-category、api-get-cwes、api-get-view、api-get-version、api-client、api-parents-children、api-response、basic-fetcher、cwe-error、cwe-not-found-error、data-fetcher、http-client、http-client-option、http-methods、http-retry、invalid-cwe-id-error、multiple-fetcher、new-xml-parser、parse-error、rate-limit-error、rate-limiter、rate-limiter-api、relationship-error、tree-fetcher、validation-error、xml-parse、xml-parser

- [ ] **Step 1: 用 sed 删除 A 派 29 文件的 `cwe ` 别名前缀 — 让 import 行变为无别名**

对 `website/sdk/` 下所有含 `cwe "github.com/scagogogo/cwe-skills"` 的文件（除 new-api-client.md），把 import 行：

```go
// 替换前
    cwe "github.com/scagogogo/cwe-skills"
// 替换后
    "github.com/scagogogo/cwe-skills"
```

执行命令（逐文件精确替换，保留行首缩进）：

```bash
cd /home/cc11001100/github/scagogogo/cwe-skills
for f in api-ancestors-descendants api-error api-get-category api-get-cwes api-get-view api-get-version api-client api-parents-children api-response basic-fetcher cwe-error cwe-not-found-error data-fetcher http-client http-client-option http-methods http-retry invalid-cwe-id-error multiple-fetcher new-xml-parser parse-error rate-limit-error rate-limiter rate-limiter-api relationship-error tree-fetcher validation-error xml-parse xml-parser; do
  sed -i 's|cwe "github.com/scagogogo/cwe-skills"|"github.com/scagogogo/cwe-skills"|' "website/sdk/${f}.md"
done
```

- [ ] **Step 2: 把 A 派 29 文件里的 `cwe.X` 包前缀调用改为 `cweskills.X`**

已验证：A 派文件中所有 `cwe.X` 后均跟大写字母（包前缀调用，非局部变量字段）。用 sed 把 `\bcwe.` 后跟大写字母的替换为 `cweskills.`：

```bash
cd /home/cc11001100/github/scagogogo/cwe-skills
for f in api-ancestors-descendants api-error api-get-category api-get-cwes api-get-view api-get-version api-client api-parents-children api-response basic-fetcher cwe-error cwe-not-found-error data-fetcher http-client http-client-option http-methods http-retry invalid-cwe-id-error multiple-fetcher new-xml-parser parse-error rate-limit-error rate-limiter rate-limiter-api relationship-error tree-fetcher validation-error xml-parse xml-parser; do
  sed -i -E 's|\bcwe\.([A-Z])|cweskills.\1|g' "website/sdk/${f}.md"
done
```

- [ ] **Step 3: 验证 A 派 29 文件无残留 `cwe ` 别名与无残留 `cwe.` 包前缀调用**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && grep -rln 'cwe "github.com/scagogogo/cwe-skills"' website/sdk/*.md`
Expected:
  - 输出仅剩 `website/sdk/new-api-client.md`（该文件 Task 3 处理，含拼写错误故本轮 sed 未覆盖）
  - 其余 29 文件不再出现在结果中

Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && grep -rnE '\bcwe\.[A-Z]' website/sdk/api-client.md website/sdk/api-error.md website/sdk/basic-fetcher.md`
Expected:
  - 无输出（A 派文件的 `cwe.X` 包前缀已全部改为 `cweskills.X`）

- [ ] **Step 4: 抽查 A 派代表文件 import + 调用一致性**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && sed -n '55,75p' website/sdk/api-client.md && echo '---' && sed -n '59,75p' website/sdk/api-error.md`
Expected:
  - import 行为 `    "github.com/scagogogo/cwe-skills"`（无 `cwe` 别名）
  - 调用为 `cweskills.NewAPIClient()`、`cweskills.WithAPIBaseURL(...)`、`cweskills.APIError` 等

- [ ] **Step 5: 提交**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && git add website/sdk && git commit -m "docs(website): 统一 SDK 示例 import 为无别名风格（A 派 cwe 别名）

website/sdk 下 29 个示例文件用了 cwe \"...\" 显式别名 + cwe. 调用前缀，
与上一轮 cmd 源码统一的「无别名 + cweskills.」标准不一致。
cwe-skills 末段去连字符即包名 cweskills，默认推导一致，别名冗余。
删别名、cwe.X → cweskills.X，统一为无别名标准。

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"`

---

### Task 2: 统一 website/sdk B 派 97 文件为无别名 — 删除冗余 `cweskills` 显式别名

**Depends on:** Task 1
**Files:**
- Modify: `website/sdk/` 下全部 97 个含 `cweskills "github.com/scagogogo/cwe-skills"` 的 .md 文件（完整清单见上方「B 派文件清单」），以 `grep -rln` 动态结果为准

- [ ] **Step 1: 列出 B 派文件实际清单**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && grep -rln 'cweskills "github.com/scagogogo/cwe-skills"' website/sdk/*.md | wc -l`
Expected:
  - 输出 `97`（B 派文件总数）

- [ ] **Step 2: 用 sed 删除 B 派文件的冗余 `cweskills ` 别名前缀**

把 import 行：

```go
// 替换前
	cweskills "github.com/scagogogo/cwe-skills"
// 替换后
	"github.com/scagogogo/cwe-skills"
```

执行命令（按 Step 1 实际清单，对所有 B 派文件执行同一替换）：

```bash
cd /home/cc11001100/github/scagogogo/cwe-skills
grep -rln 'cweskills "github.com/scagogogo/cwe-skills"' website/sdk/*.md | while read -r f; do
  sed -i 's|cweskills "github.com/scagogogo/cwe-skills"|"github.com/scagogogo/cwe-skills"|' "$f"
done
```

说明：B 派文件的调用前缀本就是 `cweskills.`（如 `cweskills.NewCWE`、`cweskills.AbstractionBase`），删别名后默认包名仍是 `cweskills`，调用点零改动。

- [ ] **Step 3: 验证 B 派文件无残留冗余别名且调用前缀不变**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && grep -rln 'cweskills "github.com/scagogogo/cwe-skills"' website/sdk/*.md`
Expected:
  - 无输出（B 派冗余别名全部删除）

Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && sed -n '35,45p' website/sdk/alternate-term.md && echo '---' && sed -n '56,78p' website/sdk/cwe-type-methods.md`
Expected:
  - import 行为 `"github.com/scagogogo/cwe-skills"`（无别名）
  - 调用仍为 `cweskills.NewCWE(...)`、`cweskills.AbstractionBase` 等（前缀不变）

- [ ] **Step 4: 提交**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && git add website/sdk && git commit -m "docs(website): 统一 SDK 示例 import 为无别名风格（B 派 cweskills 别名）

website/sdk 下 B 派示例文件用了 cweskills \"...\" 冗余显式别名。
默认包名即 cweskills，别名冗余。删别名，调用前缀 cweskills. 保持不变。

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"`

---

### Task 3: 修复 new-api-client.md 拼写错误并统一其 import — 同时处理 SKILLS.md/zh 冗余别名

**Depends on:** Task 2
**Files:**
- Modify: `website/sdk/new-api-client.md:91`
- Modify: `SKILLS.md:119`
- Modify: `SKILLS.zh.md:119`

- [ ] **Step 1: 修复 new-api-client.md 的 import 别名 + 拼写错误**

文件: `website/sdk/new-api-client.md:91`

该行同时有两个问题：`cwe` 显式别名 + `scaggogo` 拼写错误（少一个 o）。替换为无别名且拼写正确的标准行：

```go
// 替换前
    cwe "github.com/scaggogo/cwe-skills"
// 替换后
    "github.com/scagogogo/cwe-skills"
```

执行命令：

```bash
cd /home/cc11001100/github/scagogogo/cwe-skills
sed -i 's|cwe "github.com/scaggogo/cwe-skills"|"github.com/scagogogo/cwe-skills"|' website/sdk/new-api-client.md
```

- [ ] **Step 2: 把 new-api-client.md 里的 `cwe.X` 包前缀调用改为 `cweskills.X`**

执行命令（同 Task 1 Step 2 规则，仅替换 `cwe.` 后跟大写字母的包前缀）：

```bash
cd /home/cc11001100/github/scagogogo/cwe-skills
sed -i -E 's|\bcwe\.([A-Z])|cweskills.\1|g' website/sdk/new-api-client.md
```

- [ ] **Step 3: 删除 SKILLS.md / SKILLS.zh.md 的冗余 `cweskills` 别名**

文件: `SKILLS.md:119`、`SKILLS.zh.md:119`

把 import 行：

```go
// 替换前
import cweskills "github.com/scagogogo/cwe-skills"
// 替换后
import "github.com/scagogogo/cwe-skills"
```

执行命令（两个文件该行内容相同）：

```bash
cd /home/cc11001100/github/scagogogo/cwe-skills
sed -i 's|import cweskills "github.com/scagogogo/cwe-skills"|import "github.com/scagogogo/cwe-skills"|' SKILLS.md SKILLS.zh.md
```

说明：SKILLS.md/zh 的调用前缀本就是 `cweskills.`（`cweskills.ParseCWEID`、`cweskills.IsInTop25`、`cweskills.NewAPIClient`），删别名后默认包名仍是 `cweskills`，调用点零改动。

- [ ] **Step 4: 验证拼写错误已修复 + 无残留 scaggogo + SKILLS.md import 已统一**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && grep -rn 'scaggogo' . --include="*.md" --include="*.go" --include="*.ts" 2>/dev/null | grep -v './.git/' | grep -v './dist/' | grep -v node_modules`
Expected:
  - 无输出（`scaggogo` 拼写错误已全部修复，无残留）

Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && grep -n 'import' SKILLS.md SKILLS.zh.md | grep cwe-skills`
Expected:
  - 两文件均输出 `import "github.com/scagogogo/cwe-skills"`（无 `cweskills` 别名）

Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && sed -n '88,95p' website/sdk/new-api-client.md`
Expected:
  - import 行为 `    "github.com/scagogogo/cwe-skills"`（无别名、拼写正确）

- [ ] **Step 5: 提交**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && git add website/sdk/new-api-client.md SKILLS.md SKILLS.zh.md && git commit -m "docs: 修复 scaggogo 拼写错误并统一 new-api-client/SKILLS import

- website/sdk/new-api-client.md: 修复 scaggogo→scagogogo 拼写错误（少一个 o），
  并把 cwe 别名 import 统一为无别名 + cweskills. 调用
- SKILLS.md / SKILLS.zh.md: 删除 cweskills 冗余别名，统一为无别名 import

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"`

---

### Task 4: 全仓 import 风格统一终验 — 确认零残留别名、零拼写错误、零 cwe. 包前缀调用

**Depends on:** Task 3
**Files:** 无（仅验证）

- [ ] **Step 1: 验证全仓无任何显式 import 别名（cwe 或 cweskills）**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && grep -rnE '(cwe|cweskills) "github.com/scagogogo/cwe-skills"' . --include="*.md" --include="*.go" 2>/dev/null | grep -v './.git/' | grep -v './dist/' | grep -v 'docs/superpowers/plans/'`
Expected:
  - 无输出（全仓 import 别名已全部统一为无别名 `"github.com/scagogogo/cwe-skills"`）

- [ ] **Step 2: 验证全仓无 scaggogo 拼写错误**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && grep -rn 'scaggogo' . --include="*.md" --include="*.go" --include="*.ts" --include="*.yml" 2>/dev/null | grep -v './.git/' | grep -v './dist/' | grep -v node_modules`
Expected:
  - 无输出（拼写错误已全部修复）

- [ ] **Step 3: 验证 website/sdk 无残留 cwe. 包前缀调用**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && grep -rnE '\bcwe\.[A-Z]' website/sdk/*.md`
Expected:
  - 无输出（A 派 `cwe.X` 包前缀已全部改为 `cweskills.X`）

- [ ] **Step 4: 验证 Go 源码与测试仍全部通过（确认未误伤源码 import）**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && go build ./... 2>&1 | tail -5 && go test -count=1 ./... 2>&1 | tail -8`
Expected:
  - Exit code: 0
  - 三个包均 "ok"，无 FAIL（源码 import 未受文档改动影响）

- [ ] **Step 5: 验证品牌名 CWE Skills 保留不动（确认未误改产品名）**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && grep -rn 'CWE Skills' README.md README.zh.md website/.vitepress/config.ts website/guide/what-is-cwe-skills.md 2>/dev/null | head -10`
Expected:
  - 仍有大量 `CWE Skills` 匹配（品牌名保留，未被误改）

---

## 验证总表（本计划完成后）

| 检查项 | 结果 |
|---|---|
| 全仓 import 无显式别名（cwe/cweskills） | ✅ 无别名标准 |
| `scaggogo` 拼写错误 | ✅ 已修复为 `scagogogo` |
| website/sdk 无 `cwe.` 包前缀调用 | ✅ 全部 `cweskills.` |
| Go 源码 build / test | ✅ 全绿（未受影响） |
| 品牌名 `CWE Skills` | ✅ 保留不动（产品名） |

## 失败回退

- 若 Task 1/2 的 sed 误伤了局部变量 `cwe` 字段访问 → `grep -rnE '\bcwe\.[a-z]' website/sdk/*.md` 复查；若发现误替换，`git checkout -- website/sdk/` 回退该文件，改用逐行精确 Edit。
- 若 Task 3 拼写修复后发现 new-api-client.md 还有别的 `scaggogo`（非 import 行）→ 单独 Edit 修正。
- 若 go test 失败 → 说明文档 sed 误改了源码 .go 文件（不应发生，因 sed 仅作用于 website/sdk/*.md 与 SKILLS.md），`git checkout` 回退 .go 文件。
