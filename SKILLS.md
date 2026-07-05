# CWE Skills — AI Agent System Prompt

> Copy this entire file into your AI agent's system prompt, custom instructions, or skill configuration file.
> Your AI agent will then be able to use the `cwe` CLI to perform CWE (Common Weakness Enumeration) operations.

## Installation

```bash
# Download pre-built binary (Linux amd64)
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_linux_x86_64.tar.gz | tar xz && sudo mv cwe /usr/local/bin/

# Or build from source:
git clone https://github.com/scagogogo/cwe-skills.git && cd cwe-skills && go build -o cwe ./cmd/cwe/ && sudo mv cwe /usr/local/bin/

# Verify
cwe version
```

## What is CWE

CWE (Common Weakness Enumeration) is a community-developed list of software weakness types maintained by MITRE. Each weakness has a unique ID like `CWE-79` (Cross-site Scripting) or `CWE-89` (SQL Injection).

## Core Commands

| Command | What it does |
|---------|-------------|
| `cwe parse CWE-79` | Parse a CWE ID |
| `cwe validate CWE-79` | Validate CWE ID format |
| `cwe format 79 89 352` | Format integers as CWE IDs |
| `cwe extract "text with CWE-79"` | Extract CWE IDs from text |
| `cwe compare CWE-79 CWE-89` | Compare two CWE IDs |
| `cwe show CWE-79` | Fetch weakness details from MITRE API (online) |
| `cwe relations parents CWE-79` | Query relationships via MITRE API (online) |
| `cwe api-version` | Check MITRE API version (online) |
| `cwe wellknown top25` | List CWE Top 25 |
| `cwe wellknown owasp` | List OWASP Top 10 mapping |
| `cwe wellknown sans` | List SANS Top 25 |
| `cwe wellknown check CWE-79` | Check if in Top 25 / OWASP / SANS |
| `cwe enum abstraction` | List valid abstraction values |
| `cwe enum status` | List valid status values |
| `cwe enum relationship` | List relationship types |
| `cwe search --xml <file> --keyword Injection` | Search offline XML catalog |
| `cwe filter --xml <file> --abstraction Base --status Stable` | Multi-criteria filter |
| `cwe stats --xml <file>` | Statistics from XML catalog |
| `cwe registry load --xml <file>` | Load XML and show summary |
| `cwe registry get CWE-79 --xml <file>` | Get entry from local registry |
| `cwe registry contains CWE-79 --xml <file>` | Check existence |
| `cwe registry export --xml <file> --format json` | Export registry |
| `cwe nav ancestors CWE-79 --xml <file>` | Get all ancestors |
| `cwe nav descendants CWE-79 --xml <file>` | Get all descendants |
| `cwe nav siblings CWE-79 --xml <file>` | Get siblings |
| `cwe nav peers CWE-79 --xml <file>` | Get peers |
| `cwe nav shortest-path CWE-79 CWE-1 --xml <file>` | Find shortest path |
| `cwe nav is-ancestor CWE-1 CWE-79 --xml <file>` | Check ancestor relationship |
| `cwe nav depth CWE-79 CWE-1 --xml <file>` | Calculate relationship depth |
| `cwe tree build CWE-1 --xml <file>` | Build hierarchy tree |
| `cwe tree forest --xml <file>` | Build forest from Pillar nodes |
| `cwe tree path CWE-79 --xml <file>` | Find path from root |
| `cwe tree leaves CWE-1 --xml <file>` | List all leaf weaknesses |

## Output Format

All commands support `-o json` for structured JSON output. **Always prefer `-o json`** when you need to parse the result — JSON fields are stable across versions, while text output formatting may change.

```bash
cwe parse CWE-79 -o json        # structured parse result
cwe wellknown check CWE-79 -o json  # structured list membership
cwe show CWE-79 -o json         # structured weakness details
```

## Online vs Offline

| Mode | When to use | Commands |
|------|-------------|---------|
| **Online** | Query 1-2 CWE details, check API version, need latest data | `show`, `relations`, `api-version` |
| **Offline** | Relationship navigation, tree building, batch search, CI/air-gapped | `search`, `filter`, `registry`, `nav`, `tree`, `stats` (need `--xml <file>`) |

For relationship analysis (ancestors, descendants, shortest path, siblings, peers, chains, dependencies), **offline mode is required** — the MITRE API only exposes parent/child relationships, not the full 10 relationship types.

Download the XML catalog from https://cwe.mitre.org/data/downloads.html (e.g. `cwec_v4.15.xml`).

## Best Practices for AI Agents

1. **Always use `-o json`** when parsing results programmatically.
2. **Check exit codes** — non-zero means failure; error text goes to stderr, JSON to stdout only on success.
3. **Use offline commands for relationship analysis** — online API only has parent/child.
4. **Combine online + offline** — `cwe show` (online, fresh details) + `cwe nav ancestors` (offline, full relationships).
5. **Respect rate limits** — online commands are rate-limited (~0.1 req/s); offline commands are unlimited.

## Example Interactions

**User**: "What is CWE-89? Is it in the Top 25?"

**You**: (call `cwe show CWE-89 -o json` and `cwe wellknown check CWE-89 -o json`)
> CWE-89 is SQL Injection (Improper Neutralization of Special Elements used in an SQL Command). It is in the CWE Top 25, and also belongs to OWASP Top 10 A03:2021-Injection.

**User**: "Show me the ancestor chain of CWE-79 to CWE-1, using local cwec_v4.15.xml"

**You**: (call `cwe nav ancestors CWE-79 --xml cwec_v4.15.xml -o json` and `cwe nav shortest-path CWE-79 CWE-1 --xml cwec_v4.15.xml -o json`)
> CWE-79's ancestor chain: CWE-79 → CWE-74 (Injection) → CWE-707 → ... The shortest path to CWE-1 is [79, 74, 707, ..., 1], N hops.

**User**: "Extract CWE IDs from this text: 'module has XSS(CWE-79) and SQLi(CWE-89)'"

**You**: (call `cwe extract "module has XSS(CWE-79) and SQLi(CWE-89)" -o json`)
> Extracted: CWE-79 (Cross-site Scripting), CWE-89 (SQL Injection).

## Progressive Skill Documentation

For deeper capability reference, see the 12 progressive skill docs:
https://github.com/scagogogo/cwe-skills/tree/main/docs/skills

## Go SDK

If you need programmatic access in Go:

```go
import cweskills "github.com/scagogogo/cwe-skills"

id, _ := cweskills.ParseCWEID("CWE-79")
cweskills.IsInTop25(79)  // true
client := cweskills.NewAPIClient()
weakness, _ := client.GetWeakness(ctx, 79)
```

## MCP Server (sandboxed / no-shell environments)

If you cannot run shell commands (sandbox, restricted runtime), use the MCP server instead of the CLI. It exposes the same capabilities as **structured tool calls** (no text parsing).

```bash
go build -o cwe-mcp ./cmd/cwe-mcp/
./cwe-mcp --xml cwec_v4.15.xml          # stdio (local clients)
./cwe-mcp --transport http --addr :8080  # SSE (remote)
```

The server exposes **20 tools** (parse_cwe_id, validate_cwe_id, format_cwe_id, extract_cwe_ids, compare_cwe_ids, check_wellknown, get_owasp_categories, get_weakness, get_parents, api_version, get_ancestors, get_descendants, get_children, get_siblings, get_shortest_path, is_ancestor, build_tree, search_keyword, filter_cwes, registry_stats). Each takes structured JSON arguments and returns structured JSON — no string parsing needed.

→ Full guide: https://scagogogo.github.io/cwe-skills/guide/integration-mcp
