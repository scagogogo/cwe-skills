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

```bash
go install github.com/scagogogo/cwe/cmd/cwe@latest
```

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
import cwepkg "github.com/scagogogo/cwe"

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
