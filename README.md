# CWE Skills — AI-Native CWE SDK for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/scagogogo/cwe-skills.svg)](https://pkg.go.dev/github.com/scagogogo/cwe-skills)
[![CI](https://github.com/scagogogo/cwe-skills/actions/workflows/ci.yml/badge.svg)](https://github.com/scagogogo/cwe-skills/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**An AI-native SDK and CLI for [CWE (Common Weakness Enumeration)](https://cwe.mitre.org/)** — providing complete API support for building cybersecurity products, SAST/DAST tools, vulnerability management platforms, and AI-powered security agents.

> 🇨🇳 [中文文档](README.zh.md)

## Why "AI-Native"?

CWE Skills is designed from the ground up for AI agent integration:

- **Structured JSON output** on every CLI command (`-o json`) — machines can parse results directly
- **Zero interactive prompts** — every operation is non-interactive and scriptable
- **Progressive skill docs** — AI agents can discover capabilities incrementally
- **Complete SDK API** — 100+ public functions with full test coverage (97.4%)
- **Offline-first** — load MITRE XML catalogs for air-gapped environments
- **Multi-format serialization** — JSON/XML/CSV for data pipeline integration

## Features

- **Complete CWE Data Model**: Weaknesses, Categories, Views, and Compound Elements
- **Typed Enumerations**: Abstraction levels, Status values, Relationship types, Consequence scopes
- **CWE ID Utilities**: Parse, format, validate, and extract CWE IDs from text
- **Well-Known Lists**: CWE Top 25, OWASP Top 10, SANS Top 25 with membership checks
- **MITRE REST API Client**: Full access with rate limiting and retry
- **XML Catalog Parser**: Offline parsing of MITRE's official XML downloads
- **In-Memory Registry**: Store, index, and query CWE entries with relationship indexes
- **Search & Filter**: By keyword, abstraction, status, likelihood, consequence scope, and more
- **Relationship Navigation**: Parents, children, ancestors, descendants, siblings, peers, chains, composites
- **Tree Construction**: Build and traverse hierarchical trees from CWE relationships
- **Serialization**: JSON, XML, and CSV import/export
- **Cobra CLI**: 40+ subcommands with text/JSON dual output
- **Zero Dependencies**: Core SDK uses only the Go standard library

## Installation

### SDK

```bash
go get github.com/scagogogo/cwe-skills
```

### CLI — From GitHub Release (Recommended)

Download from [Releases](https://github.com/scagogogo/cwe-skills/releases/latest):

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

### CLI — From Source

```bash
git clone https://github.com/scagogogo/cwe-skills.git
cd cwe-skills
go build -o cwe ./cmd/cwe/
```

### CLI — From Package Managers

```bash
# Homebrew
brew install scagogogo/tap/cwe-skills

# Scoop (Windows)
scoop bucket add scagogogo https://github.com/scagogogo/scoop-bucket
scoop install cwe-skills

# Go Install
go install github.com/scagogogo/cwe-skills/cmd/cwe@latest
```

### Verify

```bash
cwe version
```

## Quick Start

### Using the CLI

```bash
# Parse & validate CWE IDs
cwe parse CWE-79 89 cwe-352
cwe validate CWE-79 CWE-89

# Format & extract
cwe format 79 89 352
cwe extract "Affected by CWE-79 and CWE-89"

# Check well-known lists
cwe wellknown top25
cwe wellknown check CWE-79

# Query MITRE API (online)
cwe show CWE-79
cwe relations parents CWE-79

# Local search & filter (offline, requires XML catalog)
cwe search --xml cwec_latest.xml --keyword Injection
cwe filter --xml cwec_latest.xml --abstraction Base --status Stable

# Local registry operations (offline)
cwe registry get CWE-79 --xml cwec_latest.xml
cwe registry parents CWE-79 --xml cwec_latest.xml
cwe registry export --xml cwec_latest.xml --format json

# Local relationship navigation (offline)
cwe nav siblings CWE-79 --xml cwec_latest.xml
cwe nav shortest-path CWE-79 CWE-1 --xml cwec_latest.xml
cwe nav is-ancestor CWE-1 CWE-79 --xml cwec_latest.xml

# Tree operations (offline)
cwe tree build CWE-1 --xml cwec_latest.xml
cwe tree forest --xml cwec_latest.xml
cwe tree path CWE-79 --xml cwec_latest.xml

# All commands support JSON output
cwe parse CWE-79 -o json
```

### Using the Go SDK

```go
package main

import (
    "fmt"
    "context"
    cwepkg "github.com/scagogogo/cwe-skills"
)

func main() {
    // Parse CWE ID
    id, _ := cwepkg.ParseCWEID("CWE-79")
    fmt.Println(id) // 79

    // Query MITRE API
    client := cwepkg.NewAPIClient()
    defer client.Close()
    weakness, _ := client.GetWeakness(context.Background(), 79)

    // Local registry
    registry := cwepkg.NewRegistry()
    registry.Register(&cwepkg.CWE{ID: 79, Name: "XSS", Abstraction: cwepkg.AbstractionBase})
    registry.BuildIndexes()

    // Navigate relationships
    nav := cwepkg.NewNavigator(registry)
    parents := nav.Parents(79)
    ancestors := nav.Ancestors(79)

    // Build tree
    tree := cwepkg.BuildTree(registry, 1)
    leaves := tree.LeafNodes()

    // Check well-known lists
    if cwepkg.IsInTop25(79) {
        fmt.Println("CWE-79 is in the Top 25!")
    }
}
```

### AI Agent Integration

CWE Skills is designed for AI agent integration. Copy and paste this prompt into your AI agent's configuration:

```markdown
## CWE Skills Integration

You have access to the `cwe` CLI tool for CWE (Common Weakness Enumeration) operations.

### Installation
Download from: https://github.com/scagogogo/cwe-skills/releases/latest
Or build from source: `git clone https://github.com/scagogogo/cwe-skills.git && cd cwe-skills && go build -o cwe ./cmd/cwe/`

### Key Commands
- `cwe parse CWE-79` — Parse a CWE ID
- `cwe validate CWE-79` — Validate a CWE ID format
- `cwe show CWE-79` — Fetch weakness details from MITRE API
- `cwe wellknown check CWE-79` — Check if in Top 25/OWASP/SANS lists
- `cwe search --xml <file> --keyword <term>` — Search offline XML catalog
- `cwe nav ancestors CWE-79 --xml <file>` — Navigate relationships offline
- `cwe tree build CWE-1 --xml <file>` — Build hierarchy tree

### Output Format
All commands support `-o json` for structured JSON output.

### SDK (Go)
```go
import cwepkg "github.com/scagogogo/cwe-skills"
id, _ := cwepkg.ParseCWEID("CWE-79")
cwepkg.IsInTop25(79) // true
```

### Documentation
Full skill docs: https://github.com/scagogogo/cwe-skills/tree/main/docs/skills
```

## CLI Command Reference

| Command | Description |
|---------|-------------|
| `cwe version` | Show version info |
| `cwe parse [IDs...]` | Parse CWE IDs |
| `cwe validate [IDs...]` | Validate CWE ID format |
| `cwe format [IDs...]` | Format to CWE-NNN |
| `cwe extract [text...]` | Extract CWE IDs from text |
| `cwe compare <ID1> <ID2>` | Compare two CWE IDs |
| `cwe enum <type>` | List enumeration values |
| `cwe wellknown top25` | CWE Top 25 list |
| `cwe wellknown owasp` | OWASP Top 10 list |
| `cwe wellknown sans` | SANS Top 25 list |
| `cwe wellknown check [IDs...]` | Check list membership |
| `cwe show [IDs...]` | Fetch from MITRE API |
| `cwe relations parents/children/ancestors/descendants [ID]` | API relationship queries |
| `cwe api-version` | Check MITRE API version |
| `cwe search --xml <file> [flags]` | Search offline XML |
| `cwe filter --xml <file> [flags]` | Multi-criteria filter |
| `cwe stats --xml <file>` | Statistics from XML |
| `cwe registry load/get/contains/... --xml <file>` | Registry operations |
| `cwe nav parents/children/siblings/peers/... --xml <file>` | Local navigation |
| `cwe tree build/forest/view/path/leaves --xml <file>` | Tree operations |

## Skills Documentation

Progressive skill documentation for AI agents and developers:

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

## License

MIT License - see [LICENSE](LICENSE) for details.
