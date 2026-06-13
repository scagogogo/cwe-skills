# Skill: Local — Registry Operations

## Description

Load MITRE XML catalogs and query the local Registry. All operations work offline.

## CLI Commands

All registry commands require `--xml <file>`.

```bash
cwe registry load --xml cwec_latest.xml              # Load and show summary
cwe registry get CWE-79 --xml cwec_latest.xml         # Get entry detail
cwe registry contains CWE-79 CWE-89 --xml <file>      # Check existence
cwe registry list-views --xml <file>                   # List all views
cwe registry list-categories --xml <file>              # List all categories
cwe registry parents CWE-79 --xml <file>               # Query parent IDs
cwe registry children CWE-74 --xml <file>              # Query child IDs
cwe registry ancestors CWE-79 --xml <file>             # All ancestors
cwe registry descendants CWE-74 --xml <file>           # All descendants
cwe registry peers CWE-79 --xml <file>                 # Peer IDs
cwe registry view-members 1000 --xml <file>            # View members
cwe registry category-members 1 --xml <file>           # Category members
cwe registry member-of CWE-79 --xml <file>             # What 79 belongs to
cwe registry export --xml <file> --format json         # Export as JSON
cwe registry export --xml <file> --format csv          # Export as CSV
```

| Flag | Short | Description |
|------|-------|-------------|
| `--xml` | `-x` | **(required)** Path to CWE XML catalog |
| `--format` | | Export format: json or csv (for `export`) |
| `--output-file` | | Write to file instead of stdout (for `export`) |

## SDK API

```go
registry, _ := cwepkg.NewXMLParser().ParseFile("cwec_latest.xml")
registry.BuildIndexes()

cwe, ok := registry.Get(79)
views := registry.GetAllViews()
cats := registry.GetAllCategories()
parentIDs := registry.GetParentIDs(79)
childIDs := registry.GetChildIDs(74)
ancestorIDs := registry.GetAncestorIDs(79)
descIDs := registry.GetDescendantIDs(74)
peerIDs := registry.GetPeerIDs(79)
viewMembers := registry.GetViewMembers(1000)
catMembers := registry.GetCategoryMembers(1)
memberOf := registry.GetMemberOfIDs(79)

json, _ := registry.ExportJSON()
csv, _ := registry.ExportCSV()
```

## Installation & Building from Source

```bash
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_linux_x86_64.tar.gz | tar xz
# Or: git clone https://github.com/scagogogo/cwe-skills.git && cd cwe-skills && go build -o cwe ./cmd/cwe/
```
