# AI Skills: CWE CLI

Progressive AI skill documentation for the `cwe` CLI tool — a comprehensive command-line interface for [CWE (Common Weakness Enumeration)](https://cwe.mitre.org/).

Skills are ordered from simple to advanced. Each document covers CLI commands, SDK API, and examples.

## Available Skills

| # | Skill | Command | Description |
|---|-------|---------|-------------|
| 1 | [CWE ID Parsing & Validation](01-cwe-id-parsing-validation.md) | `parse`, `validate`, `format` | Parse, validate, format CWE IDs |
| 2 | [CWE ID Extraction & Comparison](02-cwe-id-extraction-comparison.md) | `extract`, `compare` | Extract from text, compare IDs |
| 3 | [Well-Known Lists](03-well-known-lists.md) | `wellknown` | CWE Top 25, OWASP Top 10, SANS Top 25 |
| 4 | [Enumeration Types](04-enumeration-types.md) | `enum` | Abstraction, Status, Relationship types |
| 5 | [API: Get Weakness Details](05-api-show-weakness.md) | `show` | Fetch from MITRE API |
| 6 | [API: Relationship Queries](06-api-relationships.md) | `relations` | Parent/child/ancestor/descendant via API |
| 7 | [API: Version Check](07-api-version.md) | `api-version` | Check MITRE API version |
| 8 | [Local: Search & Filter](08-local-search-filter.md) | `search`, `filter` | Search & multi-criteria filter offline |
| 9 | [Local: Registry Operations](09-local-registry.md) | `registry` | Load XML, query, export local data |
| 10 | [Local: Relationship Navigation](10-local-navigation.md) | `nav` | Navigate relationships offline |
| 11 | [Local: Tree Construction](11-local-tree.md) | `tree` | Build & traverse hierarchy trees |
| 12 | [SDK: Serialization](12-sdk-serialization.md) | — | JSON, XML, CSV import/export |

## Installation

### From GitHub Release (Recommended)

```bash
# Linux (amd64)
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_linux_x86_64.tar.gz | tar xz
sudo mv cwe /usr/local/bin/

# macOS (Apple Silicon)
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_darwin_aarch64.tar.gz | tar xz
sudo mv cwe /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_windows_x86_64.zip -OutFile cwe.zip
Expand-Archive cwe.zip
```

### From Source

```bash
git clone https://github.com/scagogogo/cwe-skills.git
cd cwe-skills
go build -o cwe ./cmd/cwe/
sudo mv cwe /usr/local/bin/
```

### From Go Install

```bash
go install github.com/scagogogo/cwe-skills/cmd/cwe@latest
```

### From Package Managers

```bash
# Homebrew (macOS/Linux)
brew install scagogogo/tap/cwe-skills

# Scoop (Windows)
scoop bucket add scagogogo https://github.com/scagogogo/scoop-bucket
scoop install cwe-skills
```

### SDK

```bash
go get github.com/scagogogo/cwe-skills
```

### Verify

```bash
cwe version
```

## Supported Platforms

| OS | Architectures |
|----|---------------|
| Linux | amd64, 386, arm64, arm (v5/v6/v7), mips, mipsle, mips64, mips64le, ppc64, ppc64le, s390x, riscv64 |
| macOS | amd64 (Intel), arm64 (Apple Silicon) |
| Windows | amd64, 386, arm64, arm (v5/v6/v7) |
| FreeBSD | amd64, 386, arm64, arm (v5/v6/v7) |
| NetBSD | amd64, 386, arm (v6/v7) |
| OpenBSD | amd64, 386, arm64, arm (v7) |
| AIX | ppc64 |
| Illumos | amd64 |
| Solaris | amd64 |

## Quick Start

```bash
# Parse & validate
cwe parse CWE-79 89
cwe validate CWE-79

# Format & extract
cwe format 79 89 352
cwe extract "CWE-79 and CWE-89 affected"

# Check well-known lists
cwe wellknown top25
cwe wellknown check CWE-79

# Query MITRE API
cwe show CWE-79
cwe relations parents CWE-79

# Local operations (requires XML catalog)
cwe search --xml cwec_latest.xml --keyword Injection
cwe filter --xml cwec_latest.xml --abstraction Base --status Stable
cwe registry get CWE-79 --xml cwec_latest.xml
cwe nav siblings CWE-79 --xml cwec_latest.xml
cwe tree build CWE-1 --xml cwec_latest.xml

# All commands support JSON output
cwe parse CWE-79 -o json
```

## Global Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--output` | `-o` | `text` | Output format: `text` or `json` |

## SDK Quick Start

```go
import "github.com/scagogogo/cwe-skills"

id, _ := cweskills.ParseCWEID("CWE-79")

client := cweskills.NewAPIClient()
weakness, _ := client.GetWeakness(ctx, 79)

registry := cweskills.NewRegistry()
registry.Register(&cweskills.CWE{ID: 79, Name: "XSS", Abstraction: cweskills.AbstractionBase})
registry.BuildIndexes()

cweskills.IsInTop25(79) // true
```

## MCP Server — Alternative to CLI

If your AI agent runs in a sandbox (no shell access) or you prefer structured tool calls over CLI text parsing, use the MCP server instead of the CLI.

The `cwe-mcp` server exposes **20 tools** (parse, validate, extract, get_weakness, get_ancestors, build_tree, search_keyword, filter_cwes, etc.) over stdio/SSE for MCP-compatible clients like Claude Desktop.

```bash
# Build & run
go build -o cwe-mcp ./cmd/cwe-mcp/
./cwe-mcp --xml cwec_v4.15.xml              # stdio (local clients)
./cwe-mcp --transport http --addr :8080      # SSE (remote)
```

Configure Claude Desktop (`claude_desktop_config.json`):
```json
{
  "mcpServers": {
    "cwe-skills": {
      "command": "/path/to/cwe-mcp",
      "args": ["--xml", "/path/to/cwec_v4.15.xml"]
    }
  }
}
```

→ Full MCP guide: https://scagogogo.github.io/cwe-skills/guide/integration-mcp

### Skills (CLI) vs MCP — which to use?

| Scenario | Use |
|----------|-----|
| AI can run shell, want zero infra | **Skills** (CLI prompt) |
| AI in sandbox, no shell | **MCP** (tool calls) |
| Need structured JSON + schema | **MCP** |
| Claude Desktop / Cursor integration | **MCP** |
| Shell scripts / CI pipelines | **CLI** directly |

