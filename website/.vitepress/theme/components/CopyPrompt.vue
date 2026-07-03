<script setup lang="ts">
import { ref } from 'vue'
import { withBase } from 'vitepress'

// AI Agent 提示词 —— 引导 Codex 等 AI 客户端安装并使用 cwe 工具
const prompt = `## CWE Skills — AI-Native CWE Integration

You have access to the \`cwe\` CLI tool for CWE (Common Weakness Enumeration) operations.

### Install
\`\`\`bash
# Linux / macOS
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_linux_x86_64.tar.gz | tar xz && sudo mv cwe /usr/local/bin/
# or from source:
git clone https://github.com/scagogogo/cwe-skills.git && cd cwe-skills && go build -o cwe ./cmd/cwe/ && sudo mv cwe /usr/local/bin/
\`\`\`
Verify: \`cwe version\`

### Core Commands
| Command | Purpose |
|---------|---------|
| \`cwe parse CWE-79\` | Parse a CWE ID |
| \`cwe validate CWE-79\` | Validate CWE ID format |
| \`cwe format 79 89 352\` | Format to standard CWE-NNN |
| \`cwe extract "<text>"\` | Extract CWE IDs from text |
| \`cwe wellknown check CWE-79\` | Check Top 25 / OWASP / SANS membership |
| \`cwe wellknown top25\` | List CWE Top 25 (2024) |
| \`cwe enum abstraction\` | List enumeration values |
| \`cwe show CWE-79\` | Fetch weakness details from MITRE API |
| \`cwe relations parents CWE-79\` | Query parent weaknesses via API |
| \`cwe api-version\` | Check MITRE API version |
| \`cwe search --xml <file> --keyword Injection\` | Search offline XML catalog |
| \`cwe filter --xml <file> --abstraction Base --status Stable\` | Multi-criteria filter |
| \`cwe registry get CWE-79 --xml <file>\` | Get entry from local registry |
| \`cwe nav ancestors CWE-79 --xml <file>\` | Navigate relationships offline |
| \`cwe nav shortest-path CWE-79 CWE-1 --xml <file>\` | Find shortest path |
| \`cwe tree build CWE-1 --xml <file>\` | Build hierarchy tree |
| \`cwe stats --xml <file>\` | XML catalog statistics |

### Output Format
Every command supports \`-o json\` for structured JSON output. Example: \`cwe parse CWE-79 -o json\`

### Go SDK
\`\`\`go
import cweskills "github.com/scagogogo/cwe-skills"

id, _ := cweskills.ParseCWEID("CWE-79")
cweskills.IsInTop25(79)        // true
client := cweskills.NewAPIClient()
weakness, _ := client.GetWeakness(ctx, 79)
\`\`\`
Install: \`go get github.com/scagogogo/cwe-skills\`

### Documentation
- Full docs: https://scagogogo.github.io/cwe-skills/
- Skills index: https://scagogogo.github.io/cwe-skills/skills/
- SDK reference: https://scagogogo.github.io/cwe-skills/sdk/overview
- CLI reference: https://scagogogo.github.io/cwe-skills/cli/overview`

const copied = ref(false)

async function copyPrompt() {
  try {
    await navigator.clipboard.writeText(prompt)
  } catch {
    // fallback for non-secure contexts
    const ta = document.createElement('textarea')
    ta.value = prompt
    ta.style.position = 'fixed'
    ta.style.opacity = '0'
    document.body.appendChild(ta)
    ta.select()
    document.execCommand('copy')
    document.body.removeChild(ta)
  }
  copied.value = true
  setTimeout(() => (copied.value = false), 2000)
}
</script>

<template>
  <div class="copy-prompt">
    <div class="copy-prompt-header">
      <span class="cp-icon">🦾</span>
      <div class="cp-title">
        <strong>AI Agent 提示词</strong>
        <small>复制粘贴到 Codex / Claude / GPT 等 AI 客户端，即可安装使用 CWE Skills</small>
      </div>
      <button class="cp-button" :class="{ copied }" @click="copyPrompt">
        <span v-if="!copied">📋 一键复制</span>
        <span v-else>✅ 已复制</span>
      </button>
    </div>
    <details class="cp-details">
      <summary>查看提示词内容</summary>
      <pre><code>{{ prompt }}</code></pre>
    </details>
  </div>
</template>

<style scoped>
.copy-prompt {
  margin: 24px auto;
  max-width: 720px;
  border: 1px solid var(--vp-c-brand-1);
  border-radius: 12px;
  overflow: hidden;
  background: var(--vp-c-bg);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
}
.copy-prompt-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  background: linear-gradient(120deg, var(--vp-c-brand-soft), transparent);
}
.cp-icon { font-size: 1.6em; }
.cp-title { flex: 1; display: flex; flex-direction: column; gap: 2px; }
.cp-title strong { font-size: 1.05em; }
.cp-title small { color: var(--vp-c-text-2); font-size: 0.85em; }
.cp-button {
  border: 1px solid var(--vp-c-brand-1);
  background: var(--vp-c-brand-1);
  color: #fff;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.9em;
  font-weight: 600;
  transition: all 0.2s;
  white-space: nowrap;
}
.cp-button:hover { background: var(--vp-c-brand-2); }
.cp-button.copied { background: #16a34a; border-color: #16a34a; }
.cp-details { padding: 0 20px 16px; }
.cp-details summary {
  cursor: pointer;
  color: var(--vp-c-text-2);
  font-size: 0.9em;
  padding: 8px 0;
}
.cp-details pre {
  margin: 8px 0 0;
  padding: 14px;
  background: var(--vp-c-bg-soft);
  border-radius: 8px;
  overflow-x: auto;
  font-size: 0.78em;
  line-height: 1.5;
  max-height: 360px;
  overflow-y: auto;
}
.cp-details code { font-family: var(--vp-font-family-mono); }
</style>
