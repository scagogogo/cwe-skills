# Skill: Local — Search & Statistics

## Description

Search and analyze CWE entries from an offline MITRE XML catalog. No network required — works entirely with local data.

Download the XML catalog from [MITRE](https://cwe.mitre.org/data/xml.html).

## CLI Commands

### search

Search CWE entries from an XML catalog file.

```bash
# Search by keyword
cwe search --xml cwec_latest.xml --keyword "Injection"

# Search by abstraction level
cwe search --xml cwec_latest.xml --abstraction Base

# Search by status
cwe search --xml cwec_latest.xml --status Stable

# Search by likelihood
cwe search --xml cwec_latest.xml --likelihood High

# Search by structure type
cwe search --xml cwec_latest.xml --structure Chain

# No filter — list all entries
cwe search --xml cwec_latest.xml
```

Text output:
```
找到 2 个CWE条目:

  CWE-89 - SQL Injection [Base, Stable]
  CWE-74 - Injection [Class, Stable]
```

### stats

Compute statistics from an XML catalog file.

```bash
cwe stats --xml cwec_latest.xml
```

Text output:
```
CWE数据统计:

  总条目数:     1293
  类别数:       337
  视图数:       54

抽象层级分布:
  Pillar: 10
  Class: 69
  Base: 364
  Variant: 850

状态分布:
  Stable: 922
  Draft: 317
  Deprecated: 54

利用可能性分布:
  High: 45
  Medium: 89
  Low: 23
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--xml` | `-x` | (required) | Path to CWE XML catalog file |
| `--keyword` | `-k` | | Search by keyword (name + description) |
| `--abstraction` | `-a` | | Filter by abstraction level |
| `--status` | `-s` | | Filter by status |
| `--likelihood` | `-l` | | Filter by likelihood of exploit |
| `--structure` | `-t` | | Filter by structure type |

## SDK API

### XMLParser

```go
parser := cwepkg.NewXMLParser()

// Parse from file
registry, err := parser.ParseFile("cwec_v4.15.xml")

// Parse from io.Reader
registry, err := parser.Parse(reader)

// Parse from bytes
registry, err := parser.ParseBytes(xmlData)
```

### Search Functions

```go
// By keyword (searches name and description)
results := cwepkg.FindByKeyword(registry, "Injection")

// By abstraction level
results := cwepkg.FindByAbstraction(registry, cwepkg.AbstractionBase)

// By status
results := cwepkg.FindByStatus(registry, cwepkg.StatusStable)

// By likelihood of exploit
results := cwepkg.FindByLikelihood(registry, cwepkg.LikelihoodHigh)

// By consequence scope
results := cwepkg.FindByConsequenceScope(registry, cwepkg.ScopeConfidentiality)

// By structure type
results := cwepkg.FindByStructure(registry, cwepkg.StructureChain)

// Find top-level (pillar) weaknesses
results := cwepkg.FindTopLevel(registry)

// Find base weaknesses only
results := cwepkg.FindBaseWeaknesses(registry)

// Find chain-type weaknesses
results := cwepkg.FindChains(registry)

// Find composite-type weaknesses
results := cwepkg.FindComposites(registry)
```

### Statistics

```go
stats := cwepkg.ComputeStatistics(registry)
fmt.Println(stats.TotalCount)       // total entries
fmt.Println(stats.CategoryCount)    // categories
fmt.Println(stats.ViewCount)        // views
fmt.Println(stats.ByAbstraction)    // map[Abstraction]int
fmt.Println(stats.ByStatus)         // map[Status]int
fmt.Println(stats.ByLikelihood)     // map[LikelihoodOfExploit]int
fmt.Println(stats.ByScope)          // map[ConsequenceScope]int
```

### Filter

```go
// Multi-criteria filter
filtered := cwepkg.Filter(allCWEs, cwepkg.FilterOption{
    Abstraction:      cwepkg.AbstractionBase,
    Status:           cwepkg.StatusStable,
    LikelihoodOfExploit: cwepkg.LikelihoodHigh,
    MinID:            1,
    MaxID:            999,
    Keyword:          "Injection",
    ConsequenceScope: cwepkg.ScopeConfidentiality,
})
```

### Sort & Group

```go
cwepkg.SortByID(results)           // sort by CWE ID ascending
cwepkg.SortByName(results)         // sort by name alphabetically
cwepkg.SortByAbstraction(results)  // sort by abstraction level order

groups := cwepkg.GroupByAbstraction(registry)  // map[Abstraction][]*CWE
groups := cwepkg.GroupByStatus(registry)       // map[Status][]*CWE
groups := cwepkg.GroupByLikelihood(registry)   // map[LikelihoodOfExploit][]*CWE
```

## Important Notes

- The XML file can be large (~10MB); parsing may take a few seconds
- All search functions return `[]*CWE` slices
- `FindByKeyword` is case-insensitive and searches both name and description
- `Filter` supports combining multiple criteria — all must match (AND logic)
- Use `Deduplicate()` to remove duplicate entries from result sets