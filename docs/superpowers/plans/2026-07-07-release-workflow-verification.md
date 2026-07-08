# 跑通 Release Workflow 端到端验证计划

> **For agentic workers:** REQUIRED SUB-SKILL: `superpowers:subagent-driven-development`
> Steps use checkbox (`- [ ]`) syntax.

**Goal:** 通过打 `v0.1.0` git tag 触发 Release workflow，在真实 GitHub Actions 环境端到端验证 goreleaser 构建发布流程跑通，确认公开 Release 成功创建并包含全部产物。

**Architecture:** 推送 tag `v0.1.0` → GitHub 触发 `release.yml` → GoReleaser 在 ubuntu-latest 跑 `goreleaser release --clean` → 构建两个二进制（cwe 30+ 平台 / cwe-mcp 6 平台）→ 打包 archives + nfpms(deb/rpm/apk) + checksums → 推送 homebrew formula 到 homebrew-tap 仓库、scoop manifest 到 scoop-bucket 仓库 → 创建公开 GitHub Release。数据流：tag 推送 = 唯一触发源，`{{.Version}}` 从 tag 派生为 `0.1.0` 注入 CLI/MCP 二进制。

**Tech Stack:** GitHub Actions, GoReleaser v2 (latest), Go 1.25, goreleaser-action@v6

**Risks:**
- Task 1 打 tag 是面向外部不可逆动作 → 缓解：用户已明确授权"公开 Release"，本地 `goreleaser release --snapshot` 已跑通，配置校验通过
- homebrew-tap / scoop-bucket 仓库可能不存在或无写权限 → 缓解：Task 2 监控 workflow 日志，若推送 formula 失败，记录错误但不阻断（Release 本身仍会创建）
- GitHub Actions 环境 go 1.25 是否可用 → 缓解：CI 的 build-cli/build-mcp 作业已用 1.25 成功，release.yml 同版本
- **nfpms 包名碰撞（已修复）**：goreleaser v2.17+ 把 `linux_ppc64`/`aix_ppc64` 自动拆成 `*_power8` 变体，nfpms 模板的 `.Arch` 把它们归一化为 `ppc64`，导致同名 `cwe-skills_0.1.0_ppc64.rpm` 上传碰撞 404。同理 arm6/arm7 被 `ConventionalFileName` 映射成 armhf 碰撞 422。→ 缓解：(1) arm 用自定义 `{{ .Arch }}v{{ .Arm }}` 命名；(2) 直接移除大端 `linux_ppc64`/`aix_ppc64` target（现代 Linux 默认 ppc64le 小端），保留 `linux_ppc64le`。提交 a0ef106 + 831ddae。

---

### Task 1: 推送 v0.1.0 tag 触发 Release workflow

**Depends on:** None
**Files:**
- Modify: git tag（本地）→ origin（远程）

- [ ] **Step 1: 打 v0.1.0 tag — 在当前 main HEAD 标记发布点**

Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && git tag -a v0.1.0 -m "Release v0.1.0

Initial release of CWE Skills CLI and MCP server.

- CLI: 40+ subcommands, text/JSON dual output
- MCP server: 20 tools over stdio/SSE
- Go SDK: full CWE data model
- 30+ platform pre-built binaries"`

Expected:
- Exit code: 0
- `git tag` 输出含 `v0.1.0`

- [ ] **Step 2: 推送 tag 到远程 — 触发 Release workflow**

Run: `git push origin v0.1.0`

Expected:
- Exit code: 0
- Output contains: `* [new tag]`

- [ ] **Step 3: 确认 workflow 已触发 — GitHub 收到 tag push 事件**

Run: `sleep 10 && gh run list --workflow Release --limit 1 --json status,headSha,event,createdAt,databaseId`

Expected:
- Exit code: 0
- Output contains: `"event": "push"` and `"status": "in_progress"` or `"queued"`

---

### Task 2: 监控 Release workflow 运行直至完成

**Depends on:** Task 1
**Files:** None（只读监控）

- [ ] **Step 1: 监控 workflow 直至终态 — 等待 success 或 failure**

Run: `RID=$(gh run list --workflow Release --limit 1 --json databaseId --jq '.[0].databaseId') && gh run watch $RID --exit-status`

Expected:
- Exit code: 0（成功）
- Output contains: `success`
- 若失败：exit code 非 0，记录失败步骤

- [ ] **Step 2: 查看作业详情 — 确认 goreleaser 各步骤通过**

Run: `RID=$(gh run list --workflow Release --limit 1 --json databaseId --jq '.[0].databaseId') && gh run view $RID --json conclusion,jobs --jq '{conclusion, jobs:[.jobs[]|{name, conclusion, steps:[.steps[]|{name, conclusion}]}]}'`

Expected:
- Exit code: 0
- conclusion: `success`
- goreleaser job: `success`

---

### Task 3: 验证 Release 产物完整可访问

**Depends on:** Task 2
**Files:** None（只读验证）

- [ ] **Step 1: 确认公开 Release 已创建 — 查 GitHub Releases**

Run: `gh release view v0.1.0 --json tagName,name,isDraft,isPrerelease,url,assets --jq '{tag: .tagName, name: .name, draft: .isDraft, prerelease: .isPrerelease, url: .url, asset_count: (.assets|length)}'`

Expected:
- Exit code: 0
- draft: `false`
- prerelease: `false`（v0.1.0 非 prerelease 格式）
- asset_count: >= 45（39 cwe archive + 6 cwe-mcp archive + checksums + 包等）

- [ ] **Step 2: 抽检产物可下载 — 下载一个 cwe 和一个 cwe-mcp archive**

Run: `gh release download v0.1.0 --pattern '*linux_x86_64.tar.gz' --dir /tmp/release-check && ls -la /tmp/release-check/`

Expected:
- Exit code: 0
- 至少含 cwe-skills_0.1.0_linux_x86_64.tar.gz 和 cwe-skills-mcp_0.1.0_linux_x86_64.tar.gz

- [ ] **Step 3: 验证下载的二进制能运行**

Run: `cd /tmp/release-check && tar xzf cwe-skills_0.1.0_linux_x86_64.tar.gz && ./cwe version && tar xzf cwe-skills-mcp_0.1.0_linux_x86_64.tar.gz && ./cwe-mcp --version`

Expected:
- Exit code: 0
- `./cwe version` 输出含 `0.1.0`
- `./cwe-mcp --version` 输出含 `0.1.0`

- [ ] **Step 4: 清理下载的验证产物**

Run: `rm -rf /tmp/release-check`

Expected:
- Exit code: 0

---

### Task 4: 验证 Homebrew formula 与 Scoop manifest 推送成功

**Depends on:** Task 2
**Files:** None（只读验证 homebrew-tap / scoop-bucket 仓库）

- [ ] **Step 1: 检查 homebrew-tap 仓库是否收到 formula — goreleaser 推送的 Formula/cwe-skills.rb**

Run: `gh api repos/scagogogo/homebrew-tap/contents/Formula/cwe-skills.rb --jq '.name' 2>&1 || echo "FORMULA_NOT_FOUND"`

Expected:
- 若仓库存在且有写权限：输出 `cwe-skills.rb`
- 若仓库不存在或无权限：输出 `FORMULA_NOT_FOUND`（记录但不阻断，这是可选附属产物）

- [ ] **Step 2: 检查 scoop-bucket 仓库是否收到 manifest**

Run: `gh api repos/scagogogo/scoop-bucket/contents/cwe-skills.json --jq '.name' 2>&1 || echo "MANIFEST_NOT_FOUND"`

Expected:
- 若仓库存在：输出 `cwe-skills.json`
- 若不存在：输出 `MANIFEST_NOT_FOUND`（记录但不阻断）

- [ ] **Step 3: 提交本次验证的记录到项目仓库（如有 workflow 文档更新需要）**

Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && git status --short`

Expected:
- Exit code: 0
- 工作区干净（除非发现需修复的配置，则单独提交）

---

## 验证总表

部署流程的完整闭环验证状态：

| Workflow | 触发方式 | 验证状态 |
|----------|---------|---------|
| CI | push main | ✅ 已验证（4 作业全 success） |
| Deploy Website | push website/ | ✅ 已验证（HTTP 200 + 最新内容） |
| Release | push tag v0.1.0 | ✅ 核心已验证（v0.1.0 公开，83 资产已上传，cwe/cwe-mcp 二进制可运行输出 0.1.0） |
| (go-test.yml) | 已删除 | ✅ 移除遗留 |

### Release 工作流已知后续项

- **homebrew/scoop formula 推送失败**：`scagogogo/homebrew-tap` 和 `scagogogo/scoop-bucket` 仓库不存在（404），goreleaser 在所有资产上传成功后推送 formula 时失败。用户决策：保留现状不改，核心 Release 已成功。
- **修复历史**：本次 Release 经历 3 次失败迭代——
  1. 422 armhf.deb 碰撞（`ConventionalFileName` 把 arm6/arm7 都映射 armhf）→ 提交 a0ef106 改用 `{{ .Arch }}v{{ .Arm }}` 命名
  2. 404 ppc64.rpm 碰撞（goreleaser v2.17+ 把 `linux_ppc64`/`aix_ppc64` 拆 power8 变体，`.Arch` 都归一化 ppc64）→ 提交 831ddae 移除大端 ppc64/aix target，保留 ppc64le
  3. 404 homebrew-tap/scoop-bucket 仓库不存在 → 用户选择保留现状，不阻断核心部署
