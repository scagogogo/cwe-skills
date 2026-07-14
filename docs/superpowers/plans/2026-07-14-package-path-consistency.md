# Go 包名与 GitHub 仓库路径一致性整理计划

> **For agentic workers:** REQUIRED SUB-SKILL: `superpowers:subagent-driven-development`
> Steps use checkbox (`- [ ]`) syntax.

**Goal:** 校验并统一 Go 包名/导入路径与 GitHub 仓库 `scagogogo/cwe-skills` 的一致性——确认 module 路径与仓库路径一致、package 名合法，并消除 cmd 包内冗余的显式 import 别名，使全仓 import 风格统一为"无别名 + `cweskills.` 调用"。

**Architecture:** GitHub 仓库 `scagogogo/cwe-skills` → go.mod module `github.com/scagogogo/cwe-skills`（路径末段 `cwe-skills` 含连字符）→ 根 package 去连字符命名为 `cweskills`（合法标识符，且 Go 工具链默认推导出的包名正是 `cweskills`，与 package 声明一致）。因此 import `"github.com/scagogogo/cwe-skills"` 无需显式别名即可用 `cweskills.` 调用——这正是 README 示例与 cmd/cwe 实现文件已有的地道写法。cmd/cwe-mcp 全包及 cmd/cwe 测试文件冗余地写了显式别名 `cweskills "..."`，与默认包名重复。整理方向：删除冗余别名，全仓统一为无别名 import。

**Tech Stack:** Go 1.25.0, goimports 风格

**Risks:**
- 删别名后若某文件调用前缀不是 `cweskills.`（如笔误成别的）会编译失败 → 缓解：已用 grep 确认 cmd/cwe 全部 36 处调用均为 `cweskills.` 前缀，cmd/cwe-mcp 同理用显式别名调用 `cweskills.`，删别名后默认名仍是 `cweskills`，调用点零改动。
- 改动跨 11 个文件但每文件仅删 import 行的 `cweskills ` 前缀（4 个字符）→ 缓解：逐文件 Edit 精确替换，go build + go test 全量验证。

---

### 调研结论（一致性状态）

| 检查项 | 期望 | 实际 | 状态 |
|---|---|---|---|
| go.mod module 路径 | `github.com/scagogogo/cwe-skills` | `github.com/scagogogo/cwe-skills` | ✅ 一致 |
| GitHub 仓库路径 | `scagogogo/cwe-skills` | `scagogogo/cwe-skills` | ✅ 一致 |
| 根 package 名 | 末段去连字符 = `cweskills` | `cweskills` | ✅ 合法且与默认推导一致 |
| cmd/cwe import 路径 | `github.com/scagogogo/cwe-skills/cmd/cwe` | 一致 | ✅ |
| cmd/cwe-mcp import 路径 | `github.com/scagogogo/cwe-skills/cmd/cwe-mcp` | 一致 | ✅ |
| import 别名风格 | 全仓统一 | cmd/cwe 实现无别名 / cmd/cwe-mcp 全包+cmd/cwe 测试有冗余别名 | ❌ 待统一 |

**需修改文件（删除冗余 `cweskills ` 别名前缀）：**
1. `cmd/cwe-mcp/main.go:26`
2. `cmd/cwe-mcp/tools_api.go:7`
3. `cmd/cwe-mcp/tools_offline.go:7`
4. `cmd/cwe-mcp/tools_id.go:7`
5. `cmd/cwe-mcp/tools_wellknown.go:7`
6. `cmd/cwe-mcp/tools_extra.go:7`
7. `cmd/cwe-mcp/helpers_test.go:8`
8. `cmd/cwe-mcp/register_test.go:12`
9. `cmd/cwe-mcp/tools_api_test.go:8`
10. `cmd/cwe/helpers_test.go:10`
11. `cmd/cwe/commands_test.go:13`

---

### Task 1: 统一 SDK import 为无别名风格 — 删除 11 个文件的冗余 `cweskills` 别名

**Depends on:** None
**Files:**
- Modify: `cmd/cwe-mcp/main.go:26`、`cmd/cwe-mcp/tools_api.go:7`、`cmd/cwe-mcp/tools_offline.go:7`、`cmd/cwe-mcp/tools_id.go:7`、`cmd/cwe-mcp/tools_wellknown.go:7`、`cmd/cwe-mcp/tools_extra.go:7`、`cmd/cwe-mcp/helpers_test.go:8`、`cmd/cwe-mcp/register_test.go:12`、`cmd/cwe-mcp/tools_api_test.go:8`、`cmd/cwe/helpers_test.go:10`、`cmd/cwe/commands_test.go:13`

- [ ] **Step 1: 删除 cmd/cwe-mcp 9 个文件的冗余别名 — 把 `cweskills "github.com/scagogogo/cwe-skills"` 改为 `"github.com/scagogogo/cwe-skills"`**

每个文件的单行 import 替换（精确字符串）：

文件: `cmd/cwe-mcp/main.go:26`
```go
// 替换前
	cweskills "github.com/scagogogo/cwe-skills"
// 替换后
	"github.com/scagogogo/cwe-skills"
```

文件: `cmd/cwe-mcp/tools_api.go:7`、`cmd/cwe-mcp/tools_offline.go:7`、`cmd/cwe-mcp/tools_id.go:7`、`cmd/cwe-mcp/tools_wellknown.go:7`、`cmd/cwe-mcp/tools_extra.go:7`、`cmd/cwe-mcp/helpers_test.go:8`、`cmd/cwe-mcp/register_test.go:12`、`cmd/cwe-mcp/tools_api_test.go:8`

对上述 8 个文件，执行同一替换：把行
```go
	cweskills "github.com/scagogogo/cwe-skills"
```
替换为
```go
	"github.com/scagogogo/cwe-skills"
```

（每个文件该行内容完全相同，按文件逐个 Edit。注意保留行首的 tab 缩进。）

- [ ] **Step 2: 删除 cmd/cwe 2 个测试文件的冗余别名**

文件: `cmd/cwe/helpers_test.go:10`、`cmd/cwe/commands_test.go:13`

把行
```go
	cweskills "github.com/scagogogo/cwe-skills"
```
替换为
```go
	"github.com/scagogogo/cwe-skills"
```

- [ ] **Step 3: 验证编译与测试通过**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && go build ./... 2>&1 | tail -5`
Expected:
  - Exit code: 0
  - 无输出（编译通过，调用点 `cweskills.` 仍合法）

Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && go test -count=1 ./... 2>&1 | tail -8`
Expected:
  - Exit code: 0
  - 三个包均 "ok"，无 FAIL

- [ ] **Step 4: 验证无残留冗余别名**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && grep -rn 'cweskills "github.com/scagogogo/cwe-skills"' . --include="*.go"`
Expected:
  - Exit code: 1（grep 无匹配即 exit 1）
  - 无任何输出（全部冗余别名已删除）

Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && go vet ./... 2>&1 | tail -5`
Expected:
  - Exit code: 0
  - 无问题

- [ ] **Step 5: 提交**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && git add cmd/cwe-mcp cmd/cwe && git commit -m "refactor: 统一 SDK import 为无别名风格

cwe-skills 末段去连字符即为 package 名 cweskills，默认推导一致，
显式别名冗余。全仓统一为 \"github.com/scagogogo/cwe-skills\" 无别名 import，
与 README 示例及 cmd/cwe 实现文件风格一致。

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"`

---

## 验证总表（本计划完成后）

| 检查项 | 结果 |
|---|---|
| module 路径 = 仓库路径 | ✅ 一致（本就如此，无需改） |
| 根 package 名合法 | ✅ cweskills（本就如此） |
| import 别名风格全仓统一 | ✅ 无别名（本次整理） |
| go build / go test / go vet | ✅ 全绿 |

## 失败回退

- 若删别名后某文件编译失败（说明该文件调用前缀不是 `cweskills.`）→ 查看编译错误指向的文件，确认其 package 声明或调用前缀，单独修正该文件而非回退整个 Task。
- 若 go test 失败 → 立即 `git checkout -- cmd/` 回退，重新逐文件核对 import 行。
