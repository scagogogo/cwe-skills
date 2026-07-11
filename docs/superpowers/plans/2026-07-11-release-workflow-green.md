# Release 工作流全绿计划（移除 brews/scoops 配置）

> **For agentic workers:** REQUIRED SUB-SKILL: `superpowers:subagent-driven-development`
> Steps use checkbox (`- [ ]`) syntax.

**Goal:** 从 `.goreleaser.yml` 移除 brews 和 scoops 段（因 GITHUB_TOKEN 跨仓库写权限不足导致 403，且无法在不创建 PAT 的前提下修复），让 Release GitHub Actions 工作流从 `failure` 变为 `success`（conclusion=green），真正跑通整个部署流程。

**Architecture:** 当前根因：goreleaser 在 83 资产上传成功 + Release 发布成功后，于 `publishing` 阶段推送 homebrew formula 到 `scagogogo/homebrew-tap`、scoop manifest 到 `scagogogo/scoop-bucket`，返回 `403 Resource not accessible by integration`——GitHub Actions 自动 GITHUB_TOKEN 的 `contents:write` 只授予触发仓库 cwe-skills，组织 `default_repository_permission=none` 使 token 对同组织 tap/bucket 仓库无写权限。此约束无法通过 .goreleaser.yml 或 release.yml 配置绕过（需 PAT 或组织级权限调整，均超出零确认模式权限）。数据流：删除 brews（184-202 行）+ scoops（203-211 行）段 → goreleaser 不再有 formula 推送步骤 → 构建+打包+上传资产+发布 Release 全部成功 → exit 0 → 工作流 success。关键组件：`.goreleaser.yml`（唯一修改文件）、`release.yml`（不动）、已创建的 homebrew-tap/scoop-bucket 仓库（保留，日后配 PAT 可恢复配置）。

**Tech Stack:** GoReleaser v2.17.0（远程）/ v2.4.8（本地），goreleaser-action@v6，GitHub Actions，Go 1.25，gh CLI 2.335.1

**Risks:**
- 移除 brews/scoops 丧失 homebrew/scoop 自动分发能力 → 缓解：配置完整保存在 git 历史（提交 831ddae 及本计划文档），用户日后创建 PAT 并配置 secret 后可从历史恢复；当前该功能从未真正可用（一直 404→403 失败），移除无实际损失
- 已创建的 homebrew-tap/scoop-bucket 空仓库变为遗留 → 缓解：保留仓库（已 public + README），不删除，作为日后恢复的现成目标
- 重打 v0.1.0 前已发布的 83 资产 Release 需清理 → 缓解：`gh release delete v0.1.0 --yes --cleanup-tag` 彻底清理后基于最新提交重打
- 移除后工作流仍可能因其他原因失败 → 缓解：Task 3 监控至终态，若失败查日志定位（预期不会再有 formula 推送错误）

---

### Task 1: 从 .goreleaser.yml 移除 brews 和 scoops 段

**Depends on:** None
**Files:**
- Modify: `.goreleaser.yml:184-211`（删除 brews 段 + scoops 段，共 28 行）

- [ ] **Step 1: 删除 brews 段（第 184-202 行）和 scoops 段（第 203-211 行）— 移除因 403 失败的 formula 推送配置**

文件: `.goreleaser.yml:184-211`（brews 和 scoops 是文件最后两段，删除至文件末尾，保留前面 nfpms 段末尾的空行）

删除以下完整内容（第 184-211 行）：

```yaml
brews:
  - name: cwe-skills
    ids:
      - cwe-archive
    homepage: https://github.com/scagogogo/cwe-skills
    description: CWE (Common Weakness Enumeration) CLI tool and SDK
    license: MIT
    repository:
      owner: scagogogo
      name: homebrew-tap
    commit_author:
      name: scagogogo
      email: scagogogo@gmail.com
    directory: Formula
    test: |
      system "#{bin}/cwe", "version"
    install: |
      bin.install "cwe"

scoops:
  - repository:
      owner: scagogogo
      name: scoop-bucket
    ids:
      - cwe-archive
    homepage: https://github.com/scagogogo/cwe-skills
    description: CWE (Common Weakness Enumeration) CLI tool and SDK
    license: MIT
```

删除后，`.goreleaser.yml` 的最后一段应为 `nfpms:` 段（以 `priority: optional` 结尾），其后保留一个尾部换行。

- [ ] **Step 2: 验证删除后文件结构正确 — 确认 nfpms 是最后一段，无 brews/scoops 残留**

Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && grep -c "^brews:\|^scoops:" .goreleaser.yml && tail -5 .goreleaser.yml`

Expected:
  - Exit code: 0
  - `grep -c` 输出 `0`（无 brews/scoops 段）
  - `tail -5` 显示 nfpms 段尾部（含 `priority: optional`）

- [ ] **Step 3: 提交配置变更**

Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && git add .goreleaser.yml && git commit -m "fix(ci): 移除 brews/scoops 配置，GITHUB_TOKEN 跨仓库写权限不足导致 403

GitHub Actions 自动 GITHUB_TOKEN 的 contents:write 只授予触发仓库
cwe-skills，组织 default_repository_permission=none 使其对同组织
homebrew-tap/scoop-bucket 仓库无写权限，goreleaser 推 formula 时
返回 403 Resource not accessible by integration，导致 Release
工作流 failure。

此约束无法在不创建 PAT 的前提下绕过。移除 brews/scoops 让工作流
变绿；配置保存在 git 历史，日后配 PAT 可恢复。homebrew-tap/
scoop-bucket 仓库已创建保留，作为恢复目标。

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"`

Expected:
  - Exit code: 0
  - `git log --oneline -1` 显示新提交

---

### Task 2: 本地校验 goreleaser 配置

**Depends on:** Task 1
**Files:** None（只读校验）

- [ ] **Step 1: goreleaser check 通过 — 确认移除后配置仍有效**

Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && goreleaser check 2>&1 | tail -3`

Expected:
  - Exit code: 0
  - Output contains: `1 configuration file(s) validated`

- [ ] **Step 2: 本地快照构建成功 — 确认移除 brews/scoops 不影响构建/打包/资产**

Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && GORELEASER_CURRENT_TAG=v0.1.0 goreleaser release --snapshot --clean 2>&1 | tail -3`

Expected:
  - Exit code: 0
  - Output contains: `release succeeded`
  - Output does NOT contain: `homebrew formula` or `scoop manifest`（这两步已移除）

- [ ] **Step 3: 推送提交到 main**

Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && git push origin main 2>&1 | tail -2`

Expected:
  - Exit code: 0
  - Output contains: `main -> main`

---

### Task 3: 重打 v0.1.0 tag 并监控工作流至 success

**Depends on:** Task 2
**Files:**
- Modify: git tag v0.1.0（清理后基于最新 main HEAD 重打）

- [ ] **Step 1: 删除旧 v0.1.0 Release 及 tag — 清理 83 资产残留**

Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && gh release delete v0.1.0 --yes --cleanup-tag 2>&1 || echo "（Release 可能已删，继续）"`

Expected:
  - Exit code: 0
  - `git ls-remote --tags origin v0.1.0` 无输出

- [ ] **Step 2: 拉取最新 main 并基于新提交重打 v0.1.0 tag**

Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && git tag -d v0.1.0 2>&1 || echo "（本地无此 tag）" && git checkout main && git pull origin main 2>&1 | tail -2 && git tag -a v0.1.0 -m "Release v0.1.0

Initial release of CWE Skills CLI and MCP server.

- CLI: 40+ subcommands, text/JSON dual output
- MCP server: 20 tools over stdio/SSE
- Go SDK: full CWE data model
- 30+ platform pre-built binaries" && git push origin v0.1.0 2>&1 | tail -3`

Expected:
  - Exit code: 0
  - Output contains: `* [new tag]`
  - `git log --oneline -1` 显示 Task 1 的移除提交

- [ ] **Step 3: 确认新工作流已触发**

Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && sleep 12 && gh run list --workflow Release --limit 1 --json status,event,databaseId,headSha`

Expected:
  - Exit code: 0
  - Output contains: `"event": "push"` and (`"status": "in_progress"` or `"status": "queued"`)

- [ ] **Step 4: 监控工作流直至终态 — 等待 success（核心验收点）**

Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && RID=$(gh run list --workflow Release --limit 1 --json databaseId --jq '.[0].databaseId') && gh run watch $RID --exit-status`

Expected:
  - Exit code: 0（成功——"整个部署流程跑通"的最终证据）
  - Output contains: `success`
  - 若失败：exit code 非 0，查 `gh run view $RID --log 2>/dev/null | sed 's/\x1b\[[0-9;]*m//g' | grep -iE "⨯|release failed|error=|403|404" | head -10` 定位新错误

- [ ] **Step 5: 确认工作流 conclusion=success 且 Release 仍完整**

Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && RID=$(gh run list --workflow Release --limit 1 --json databaseId --jq '.[0].databaseId') && gh run view $RID --json conclusion --jq '.conclusion' && gh release view v0.1.0 --json tagName,isDraft,isPrerelease,assets --jq '{tag: .tagName, draft: .isDraft, prerelease: .isPrerelease, asset_count: (.assets|length)}'`

Expected:
  - Exit code: 0
  - conclusion 输出: `success`
  - Release: draft=false, prerelease=false, asset_count >= 45

- [ ] **Step 6: 提交计划文档与更新记忆**

Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && git add docs/superpowers/plans/2026-07-11-release-workflow-green.md && git commit -m "docs(plan): 记录 Release 工作流全绿方案（移除 brews/scoops）

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>" && git push origin main 2>&1 | tail -2`

Expected:
  - Exit code: 0
  - Output contains: `main -> main`

---

## 验证总表（本计划完成后）

| Workflow | 触发方式 | 验证状态 |
|----------|---------|---------|
| CI | push main | ✅ 已验证（4 作业全 success） |
| Deploy Website | push website/ | ✅ 已验证（HTTP 200 + 最新内容） |
| Release | push tag v0.1.0 | ✅ 全绿（conclusion=success，移除 brews/scoops 后无 403） |
| (go-test.yml) | 已删除 | ✅ 移除遗留 |

## 恢复 homebrew/scoop 分发（日后可选）

若日后要恢复 homebrew/scoop 自动分发：
1. 创建 fine-grained PAT，对 scagogogo/homebrew-tap 和 scoop-bucket 授 contents:write
2. 加为 cwe-skills 仓库 secret（如 `HOMEBREW_TAP_TOKEN`、`SCOOP_BUCKET_TOKEN`）
3. 从 git 历史恢复 brews/scoops 段（提交 831ddae 前的 .goreleaser.yml）
4. 在 brews/scoops 的 repository 下加 `token: {{ .Env.HOMEBREW_TAP_TOKEN }}`
5. release.yml 的 env 加 `HOMEBREW_TAP_TOKEN: ${{ secrets.HOMEBREW_TAP_TOKEN }}`
6. 重打 tag 验证

## 失败回退

- 若 Task 2 本地校验失败（goreleaser check 不过）→ 检查是否误删了 nfpms 段的尾部，恢复后重删
- 若 Task 3 Step 4 工作流仍 failure → 查日志定位新错误（预期不再有 formula 推送 403）；若新错误与 brews/scoops 无关，单独修复
- 若 Task 3 Step 5 Release 资产数 < 45 → 检查资产上传是否完整，可能需重新触发
