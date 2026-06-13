# Skill: Local — Relationship Navigation

## Description

Navigate CWE relationships offline using the Navigator. Provides richer queries than the API-based `cwe relations`, including siblings, peers, shortest path, and relationship depth.

## CLI Commands

All nav commands require `--xml <file>`.

```bash
# Hierarchical
cwe nav parents CWE-79 --xml <file>
cwe nav children CWE-74 --xml <file>
cwe nav ancestors CWE-79 --xml <file>
cwe nav descendants CWE-74 --xml <file>
cwe nav siblings CWE-79 --xml <file>

# Peer & sequential
cwe nav peers CWE-79 --xml <file>
cwe nav precede CWE-89 --xml <file>
cwe nav follow CWE-79 --xml <file>

# Dependency
cwe nav requires CWE-79 --xml <file>
cwe nav required-by CWE-79 --xml <file>
cwe nav can-also-be CWE-79 --xml <file>

# Compound
cwe nav chain-members 680 --xml <file>
cwe nav composite-members 680 --xml <file>

# Path queries
cwe nav shortest-path CWE-79 CWE-1 --xml <file>
cwe nav is-ancestor CWE-1 CWE-79 --xml <file>
cwe nav is-related CWE-79 CWE-89 --xml <file>
cwe nav depth CWE-79 CWE-1 --xml <file>
```

| Flag | Short | Description |
|------|-------|-------------|
| `--xml` | `-x` | **(required)** Path to CWE XML catalog |

## SDK API

```go
registry, _ := cwepkg.NewXMLParser().ParseFile("cwec_latest.xml")
registry.BuildIndexes()
nav := cwepkg.NewNavigator(registry)

parents := nav.Parents(79)
children := nav.Children(74)
ancestors := nav.Ancestors(79)
descendants := nav.Descendants(74)
siblings := nav.Siblings(79)
peers := nav.Peers(79)
path := nav.ShortestPath(79, 1)  // []int or nil
isAncestor := nav.IsAncestorOf(1, 79)  // bool
isRelated := nav.IsRelated(79, 89)  // bool
depth := nav.RelationshipDepth(79, 1)  // int (-1 if unrelated)
```

## Installation & Building from Source

```bash
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_linux_x86_64.tar.gz | tar xz
# Or: git clone https://github.com/scagogogo/cwe-skills.git && cd cwe-skills && go build -o cwe ./cmd/cwe/
```
