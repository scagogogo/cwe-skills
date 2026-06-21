# Skill: Local — Search & Filter

## Description

Search and filter CWE entries from an offline MITRE XML catalog. No network required.

Download XML from [MITRE](https://cwe.mitre.org/data/xml.html).

## CLI Commands

### search — Single-criterion search

```bash
cwe search --xml cwec_latest.xml [flags]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--xml` | `-x` | **(required)** Path to CWE XML catalog |
| `--keyword` | `-k` | Search by keyword (name + description) |
| `--abstraction` | `-a` | Filter by abstraction (Pillar/Class/Base/Variant) |
| `--status` | `-s` | Filter by status (Stable/Draft/Deprecated) |
| `--likelihood` | `-l` | Filter by likelihood (High/Medium/Low) |
| `--structure` | `-t` | Filter by structure (Simple/Chain/Composite) |
| `--scope` | | Filter by consequence scope |
| `--top-level` | | Show only pillar (top-level) weaknesses |
| `--base-weaknesses` | | Show only base weaknesses |
| `--chains` | | Show only chain-type weaknesses |
| `--composites` | | Show only composite weaknesses |
| `--sort` | | Sort by: id, name, abstraction |
| `--group-by` | | Group by: abstraction, status, likelihood |
| `--dedup` | | Remove duplicates |

### filter — Multi-criteria filter (AND logic)

```bash
cwe filter --xml cwec_latest.xml --abstraction Base --status Stable --keyword Injection
cwe filter --xml cwec_latest.xml --likelihood High --scope Confidentiality --sort name
```

### stats — Statistics

```bash
cwe stats --xml cwec_latest.xml
```

## SDK API

```go
results := cweskills.FindByKeyword(registry, "Injection")
results := cweskills.FindByAbstraction(registry, cweskills.AbstractionBase)
results := cweskills.FindByConsequenceScope(registry, cweskills.ScopeConfidentiality)
results := cweskills.FindTopLevel(registry)
results := cweskills.FindBaseWeaknesses(registry)

filtered := cweskills.Filter(all, cweskills.FilterOption{
    Abstraction: cweskills.AbstractionBase,
    Status:      cweskills.StatusStable,
    Keyword:     "Injection",
})

cweskills.SortByID(results)
cweskills.SortByName(results)
cweskills.GroupByAbstraction(results)
cweskills.Deduplicate(results)
```

## Installation & Building from Source

```bash
# Download pre-built binary
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_linux_x86_64.tar.gz | tar xz

# Or build from source
git clone https://github.com/scagogogo/cwe-skills.git
cd cwe-skills && go build -o cwe ./cmd/cwe/
```
