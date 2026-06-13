# AI Skills: CWE CLI

This directory contains progressive AI skill documentation for the `cwe` CLI tool. Each skill document describes one capability area in a structured format that AI agents can use to understand and invoke the tool.

Skills are ordered from simple to advanced — start with #1 and progress as needed.

## Available Skills

| # | Skill | Command | Description |
|---|-------|---------|-------------|
| 1 | CWE ID Parsing & Validation | `parse`, `validate`, `format` | Parse, validate, and format CWE IDs |
| 2 | CWE ID Extraction & Comparison | `extract`, `compare` | Extract CWE IDs from text, compare IDs |
| 3 | Well-Known Lists | `wellknown` | Query CWE Top 25, OWASP Top 10, SANS Top 25 |
| 4 | Enumeration Types | `enum` | List CWE enumeration values (abstraction, status, etc.) |
| 5 | API: Get Weakness Details | `show` | Fetch CWE weakness details from MITRE API |
| 6 | API: Relationship Queries | `relations` | Query parent/child/ancestor/descendant relationships |
| 7 | API: Version Check | `api-version` | Check MITRE CWE API version |
| 8 | Local: Search & Statistics | `search`, `stats` | Search and analyze offline XML catalogs |
| 9 | SDK: In-Memory Registry | — | Store, index, and query CWE entries locally |
| 10 | SDK: Relationship Navigation | — | Navigate CWE hierarchies and relationships |
| 11 | SDK: Tree Construction | — | Build and traverse CWE hierarchy trees |
| 12 | SDK: Serialization | — | JSON, XML, CSV import/export |

## Global Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--output` | `-o` | `text` | Output format: `text` or `json` |

## Installation

### From GitHub Release (Recommended)

Download the latest binary for your platform from [GitHub Releases](https://github.com/scagogogo/cwe-skills/releases/latest):

```bash
# Linux (amd64)
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_linux_x86_64.tar.gz | tar xz
sudo mv cwe /usr/local/bin/

# Linux (arm64)
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_linux_aarch64.tar.gz | tar xz
sudo mv cwe /usr/local/bin/

# macOS (Apple Silicon)
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_darwin_aarch64.tar.gz | tar xz
sudo mv cwe /usr/local/bin/

# macOS (Intel)
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_darwin_x86_64.tar.gz | tar xz
sudo mv cwe /usr/local/bin/

# Windows (amd64) — PowerShell
Invoke-WebRequest -Uri https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_windows_x86_64.zip -OutFile cwe.zip
Expand-Archive cwe.zip
```

### Specific Version

Replace `latest` with a specific tag (e.g., `v0.1.0`) in the download URL:

```bash
curl -sL https://github.com/scagogogo/cwe-skills/releases/download/v0.1.0/cwe-skills_v0.1.0_linux_x86_64.tar.gz | tar xz
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

# DEB package (Ubuntu/Debian)
sudo dpkg -i cwe-skills_*_linux_amd64.deb

# RPM package (CentOS/RHEL/Fedora)
sudo rpm -i cwe-skills-*_linux_amd64.rpm

# APK package (Alpine)
sudo apk add cwe-skills_*_linux_amd64.apk
```

### SDK Installation

```bash
go get github.com/scagogogo/cwe-skills
```

### Verify Installation

```bash
cwe version
# Output:
# CWE CLI:     v0.1.0
# CWE SDK:     v0.0.1
# Go Version:  go1.25.0
```

## Supported Platforms

Pre-built binaries are available for:

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
# Parse a CWE ID
cwe parse CWE-79

# Validate CWE IDs
cwe validate CWE-79 CWE-89

# Format to standard form
cwe format 79 89 352

# Extract from text
cwe extract "Affected by CWE-79 and CWE-89"

# Compare two IDs
cwe compare CWE-79 CWE-89

# Check well-known lists
cwe wellknown top25
cwe wellknown check CWE-79

# List enumeration values
cwe enum abstraction

# Fetch weakness from MITRE API
cwe show CWE-79

# Query relationships
cwe relations parents CWE-79

# Search offline XML
cwe search --xml cwec_latest.xml --keyword Injection

# Get API version
cwe api-version
```

## JSON Output Schema

All commands support `-o json` for structured output. Common patterns:

**Success response** — each command returns its own typed structure (see individual skill docs).

**Error response** — errors appear inline in the result array for batch operations, or as CLI errors:

```json
{"error": "string"}
```

## SDK Quick Start

```go
import cwepkg "github.com/scagogogo/cwe-skills"

// Parse CWE ID
id, _ := cwepkg.ParseCWEID("CWE-79")

// Query MITRE API
client := cwepkg.NewAPIClient()
weakness, _ := client.GetWeakness(ctx, 79)

// Local registry
registry := cwepkg.NewRegistry()
registry.Register(&cwepkg.CWE{ID: 79, Name: "XSS", Abstraction: cwepkg.AbstractionBase})
registry.BuildIndexes()

// Check well-known lists
cwepkg.IsInTop25(79) // true
```
