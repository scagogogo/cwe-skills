<p align="center">
  <img src="docs/logo.svg" alt="CWE Skills" width="160" />
</p>

<h1 align="center">CWE Skills — AI-Native CWE Integration</h1>

<p align="center">
  <a href="https://pkg.go.dev/github.com/scagogogo/cwe-skills"><img src="https://pkg.go.dev/badge/github.com/scagogogo/cwe-skills.svg" alt="Go Reference" /></a>
  <a href="https://github.com/scagogogo/cwe-skills/actions/workflows/ci.yml"><img src="https://github.com/scagogogo/cwe-skills/actions/workflows/ci.yml/badge.svg" alt="CI" /></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT" /></a>
  <a href="https://scagogogo.github.io/cwe-skills/"><img src="https://img.shields.io/badge/docs-online-3c6c8f.svg" alt="Docs" /></a>
</p>

<p align="center"><b>AI-native integration layer for <a href="https://cwe.mitre.org/">CWE (Common Weakness Enumeration)</a></b> — four ways to connect: <b>Skills</b>, Go SDK, CLI, and MCP.</p>

<p align="center">🇨🇳 <a href="README.zh.md">中文文档</a> · 📖 <a href="https://scagogogo.github.io/cwe-skills/">Documentation</a></p>

---

## 🚀 Four Ways to Integrate

```mermaid
flowchart TB
    subgraph APP["应用层"]
        AI["AI 代理\nClaude/GPT"]
        GO["Go 应用"]
        SH["Shell 脚本"]
        MCP["MCP 工具"]
    end

    subgraph SDK["CWE Skills 集成层"]
        SK["🦾 Skills"]
        GS["🔧 Go SDK"]
        CLI["💻 CLI"]
        MC["MCP Server"]
    end

    subgraph DATA["数据源"]
        API["MITRE REST API"]
        XML["XML 目录"]
        LIST["内置列表\nTop25/OWASP/SANS"]
    end

    AI --> SK --> API & XML & LIST
    GO --> GS --> API & XML & LIST
    SH --> CLI --> API & XML & LIST
    MCP --> MC --> API & XML & LIST

    classDef app fill:#dcfce7,stroke:#16a34a,color:#166534
    classDef sdk fill:#e8f1f8,stroke:#3c6c8f,color:#1d3a4f
    classDef data fill:#ffedd5,stroke:#ea580c,color:#9a3412
    class APP app
    class SDK sdk
    class DATA data
```

| # | Method | Best For | One-Line Setup |
|---|--------|----------|----------------|
| 1 | **Skills** | AI agents (Claude, GPT, etc.) | Copy the prompt below |
| 2 | **Go SDK** | Go applications & libraries | `go get github.com/scagogogo/cwe-skills` |
| 3 | **CLI** | Shell scripts & dev workflows | Download from [Releases](https://github.com/scagogogo/cwe-skills/releases/latest) |
| 4 | **MCP** | MCP-compatible AI tools | `go build ./cmd/cwe-mcp` |

---

## 1. Skills — AI Agent Integration

Copy and paste this block into your AI agent's system prompt or skill configuration:

```markdown
## CWE Skills

You have access to the `cwe` CLI tool for CWE (Common Weakness Enumeration) operations.

### Install
```bash
# Download pre-built binary (Linux/macOS/Windows)
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_linux_x86_64.tar.gz | tar xz && sudo mv cwe /usr/local/bin/
# Or build from source:
git clone https://github.com/scagogogo/cwe-skills.git && cd cwe-skills && go build -o cwe ./cmd/cwe/ && sudo mv cwe /usr/local/bin/
```

### Core Commands
| Command | What it does |
|---------|-------------|
| `cwe parse CWE-79` | Parse a CWE ID |
| `cwe validate CWE-79` | Validate CWE ID format |
| `cwe show CWE-79` | Fetch weakness details from MITRE API |
| `cwe wellknown check CWE-79` | Check if in Top 25 / OWASP / SANS lists |
| `cwe enum abstraction` | List valid enumeration values |
| `cwe search --xml <file> --keyword Injection` | Search offline XML catalog |
| `cwe filter --xml <file> --abstraction Base --status Stable` | Multi-criteria filter |
| `cwe registry get CWE-79 --xml <file>` | Get entry from local registry |
| `cwe nav ancestors CWE-79 --xml <file>` | Navigate relationships offline |
| `cwe nav shortest-path CWE-79 CWE-1 --xml <file>` | Find shortest path between two CWEs |
| `cwe tree build CWE-1 --xml <file>` | Build hierarchy tree |
| `cwe stats --xml <file>` | Statistics from XML catalog |

### Output
All commands support `-o json` for structured JSON output. Example: `cwe parse CWE-79 -o json`

### Go SDK
```go
import "github.com/scagogogo/cwe-skills"
id, _ := cweskills.ParseCWEID("CWE-79")
cweskills.IsInTop25(79) // true
client := cweskills.NewAPIClient()
weakness, _ := client.GetWeakness(ctx, 79)
```

### Skill Docs
Progressive capability docs: https://github.com/scagogogo/cwe-skills/tree/main/docs/skills
```

---

## 2. Go SDK

```mermaid
flowchart LR
    subgraph SRC["数据源"]
        API["MITRE API"]
        XML["XML 目录"]
    end
    subgraph PARSE["解析"]
        AC["APIClient"]
        XP["XMLParser"]
    end
    REG["Registry\n+ 索引"]
    subgraph USE["消费"]
        NAV["Navigator"]
        TREE["BuildTree"]
        SRCH["FindBy*/Filter"]
        SER["MarshalJSON/CSV"]
    end
    API --> AC --> REG
    XML --> XP --> REG
    REG --> NAV & TREE & SRCH & SER

    classDef online fill:#dbeafe,stroke:#2563eb,color:#1e40af
    classDef offline fill:#ffedd5,stroke:#ea580c,color:#9a3412
    classDef core fill:#e8f1f8,stroke:#3c6c8f,color:#1d3a4f
    classDef local fill:#dcfce7,stroke:#16a34a,color:#166534
    class API,AC online
    class XML,XP offline
    class REG core
    class USE local
```

```go
import (
    "context"
    "github.com/scagogogo/cwe-skills"
)

// Parse & validate CWE IDs
id, _ := cweskills.ParseCWEID("CWE-79")
if cweskills.IsCWEID("CWE-89") { /* valid */ }

// Query MITRE REST API
client := cweskills.NewAPIClient()
defer client.Close()
weakness, _ := client.GetWeakness(context.Background(), 79)
parents, _ := client.GetParents(context.Background(), 79)

// Local registry from XML
registry, _ := cweskills.NewXMLParser().ParseFile("cwec_v4.15.xml")
registry.BuildIndexes()

// Navigate relationships
nav := cweskills.NewNavigator(registry)
ancestors := nav.Ancestors(79)
path := nav.ShortestPath(79, 1)

// Build hierarchy tree
tree := cweskills.BuildTree(registry, 1)
leaves := tree.LeafNodes()

// Search & filter
results := cweskills.FindByKeyword(registry, "Injection")
filtered := cweskills.Filter(results, cweskills.FilterOption{
    Abstraction: cweskills.AbstractionBase,
    Status:      cweskills.StatusStable,
})

// Well-known lists
cweskills.IsInTop25(79)       // true
cweskills.IsInOWASPTop10(79)  // true
cweskills.IsInSANSTop25(79)   // true

// Serialization
jsonData, _ := registry.ExportJSON()
csvData, _ := registry.ExportCSV()
```

**Install**: `go get github.com/scagogogo/cwe-skills`

---

## 3. CLI

```mermaid
mindmap
  root((cwe CLI))
    🆔 ID 工具
      parse
      validate
      format
      extract
      compare
    📚 枚举
      enum abstraction
      enum status
      enum relationship
    🏆 知名列表
      wellknown top25
      wellknown owasp
      wellknown sans
      wellknown check
    🌐 API
      show
      relations
      api-version
    🔍 搜索过滤
      search
      filter
      stats
    🗃️ 注册表
      registry load/get
      registry export
    🧭 导航
      nav parents/children
      nav ancestors/descendants
      nav shortest-path
    🌳 树
      tree build
      tree forest
      tree path
```

### Install

**From Release** (recommended):
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

**From Source**:
```bash
git clone https://github.com/scagogogo/cwe-skills.git
cd cwe-skills && go build -o cwe ./cmd/cwe/
```

**From Package Managers**:
```bash
brew install scagogogo/tap/cwe-skills          # Homebrew
scoop install cwe-skills                         # Scoop (Windows)
go install github.com/scagogogo/cwe-skills/cmd/cwe@latest  # Go
```

### Quick Examples

```bash
# CWE ID operations
cwe parse CWE-79 89 cwe-352
cwe validate CWE-79 CWE-89
cwe format 79 89 352
cwe extract "Affected by CWE-79 and CWE-89"
cwe compare CWE-79 CWE-89

# Well-known lists
cwe wellknown top25
cwe wellknown owasp
cwe wellknown check CWE-79

# MITRE API (online)
cwe show CWE-79
cwe relations parents CWE-79
cwe api-version

# Local search & filter (offline)
cwe search --xml cwec_latest.xml --keyword Injection --sort name
cwe filter --xml cwec_latest.xml --abstraction Base --status Stable --likelihood High

# Local registry (offline)
cwe registry load --xml cwec_latest.xml
cwe registry get CWE-79 --xml cwec_latest.xml
cwe registry ancestors CWE-79 --xml cwec_latest.xml
cwe registry export --xml cwec_latest.xml --format json

# Local navigation (offline)
cwe nav siblings CWE-79 --xml cwec_latest.xml
cwe nav peers CWE-79 --xml cwec_latest.xml
cwe nav shortest-path CWE-79 CWE-1 --xml cwec_latest.xml
cwe nav is-ancestor CWE-1 CWE-79 --xml cwec_latest.xml
cwe nav depth CWE-79 CWE-1 --xml cwec_latest.xml

# Tree operations (offline)
cwe tree build CWE-1 --xml cwec_latest.xml
cwe tree forest --xml cwec_latest.xml
cwe tree path CWE-79 --xml cwec_latest.xml
cwe tree leaves CWE-1 --xml cwec_latest.xml

# Enumeration types
cwe enum abstraction
cwe enum status
cwe enum relationship

# JSON output on every command
cwe parse CWE-79 -o json
cwe wellknown check CWE-79 -o json
```

### Command Reference

| Command | Description |
|---------|-------------|
| `cwe version` | Show version info |
| `cwe parse/validate/format/extract/compare` | CWE ID utilities |
| `cwe enum <type>` | List enumeration values |
| `cwe wellknown top25/owasp/sans/check` | Well-known lists |
| `cwe show [IDs...]` | Fetch from MITRE API |
| `cwe relations parents/children/ancestors/descendants` | API relationships |
| `cwe api-version` | Check MITRE API version |
| `cwe search --xml <file> [flags]` | Search offline XML |
| `cwe filter --xml <file> [flags]` | Multi-criteria filter |
| `cwe stats --xml <file>` | Statistics |
| `cwe registry <subcmd> --xml <file>` | Registry operations |
| `cwe nav <subcmd> --xml <file>` | Relationship navigation |
| `cwe tree <subcmd> --xml <file>` | Tree operations |

---

## 4. MCP

The `cwe-mcp` server exposes 15 tools (parse, validate, extract, get_weakness, get_ancestors, build_tree, search_keyword, etc.) over stdio/SSE for MCP-compatible AI tools like Claude Desktop.

```bash
# Build
go build -o cwe-mcp ./cmd/cwe-mcp/

# Run (stdio for local clients)
./cwe-mcp --xml cwec_v4.15.xml

# Run (SSE for remote)
./cwe-mcp --transport http --addr :8080
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

→ **[MCP integration guide](https://scagogogo.github.io/cwe-skills/guide/integration-mcp)**

---

## Skills Documentation

Progressive skill documentation for AI agents and developers — from simple to advanced:

| # | Skill | Description |
|---|-------|-------------|
| 1 | [CWE ID Parsing & Validation](docs/skills/01-cwe-id-parsing-validation.md) | Parse, validate, format CWE IDs |
| 2 | [CWE ID Extraction & Comparison](docs/skills/02-cwe-id-extraction-comparison.md) | Extract from text, compare IDs |
| 3 | [Well-Known Lists](docs/skills/03-well-known-lists.md) | CWE Top 25, OWASP Top 10, SANS Top 25 |
| 4 | [Enumeration Types](docs/skills/04-enumeration-types.md) | Abstraction, Status, Relationship types |
| 5 | [API: Get Weakness Details](docs/skills/05-api-show-weakness.md) | Fetch from MITRE API |
| 6 | [API: Relationship Queries](docs/skills/06-api-relationships.md) | Parent/child/ancestor/descendant via API |
| 7 | [API: Version Check](docs/skills/07-api-version.md) | Check MITRE API version |
| 8 | [Local: Search & Filter](docs/skills/08-local-search-filter.md) | Search & multi-criteria filter |
| 9 | [Local: Registry Operations](docs/skills/09-local-registry.md) | Load, query, export local data |
| 10 | [Local: Relationship Navigation](docs/skills/10-local-navigation.md) | Navigate relationships offline |
| 11 | [Local: Tree Construction](docs/skills/11-local-tree.md) | Build & traverse hierarchy trees |
| 12 | [SDK: Serialization](docs/skills/12-sdk-serialization.md) | JSON, XML, CSV import/export |

→ **[Full Skills Index](docs/skills/README.md)**

## Supported Platforms

Pre-built binaries for 30+ platforms: Linux (amd64/386/arm64/arm/mips/ppc64/s390x/riscv64), macOS (Intel/Apple Silicon), Windows (amd64/386/arm64), FreeBSD, NetBSD, OpenBSD, AIX, Illumos, Solaris.

## Features

- **Complete CWE Data Model**: Weaknesses, Categories, Views, Compound Elements
- **Typed Enumerations**: Abstraction, Status, Relationship, Consequence, View types
- **CWE ID Utilities**: Parse, format, validate, extract, compare
- **Well-Known Lists**: CWE Top 25, OWASP Top 10, SANS Top 25
- **MITRE REST API Client**: Rate limiting, retry, structured errors
- **XML Catalog Parser**: Offline MITRE XML parsing
- **In-Memory Registry**: Store, index, query with relationship indexes
- **Search & Filter**: Keyword, abstraction, status, likelihood, scope, sort, group
- **Relationship Navigation**: Parents, children, ancestors, descendants, siblings, peers, chains, composites, shortest path, relationship depth
- **Tree Construction**: Build, traverse, find paths, list leaves
- **Serialization**: JSON, XML, CSV import/export
- **40+ CLI Subcommands**: Text/JSON dual output
- **Zero Dependencies**: Core SDK uses only Go standard library

## License

MIT License - see [LICENSE](LICENSE) for details.
