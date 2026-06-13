# Skill: Local — Tree Construction

## Description

Build and traverse CWE hierarchy trees from offline XML data. Visualize the full weakness taxonomy.

## CLI Commands

All tree commands require `--xml <file>`.

```bash
# Build a tree from a root node
cwe tree build CWE-1 --xml <file>

# Build forest of all pillar nodes
cwe tree forest --xml <file>

# Build tree from a specific view
cwe tree view 1000 --xml <file>

# Find path from root to a CWE
cwe tree path CWE-79 --xml <file>
cwe tree path CWE-79 --xml <file> --root 1

# List all leaf nodes under a root
cwe tree leaves CWE-1 --xml <file>
```

| Flag | Short | Description |
|------|-------|-------------|
| `--xml` | `-x` | **(required)** Path to CWE XML catalog |
| `--root` | | Root node ID for `path` (auto-detected if omitted) |

## SDK API

```go
registry, _ := cwepkg.NewXMLParser().ParseFile("cwec_latest.xml")
registry.BuildIndexes()

tree := cwepkg.BuildTree(registry, 1)        // *TreeNode
forest := cwepkg.BuildForest(registry)         // []*TreeNode
viewTree := cwepkg.BuildViewTree(registry, 1000) // *TreeNode

// Traversal
tree.Walk(func(node *cwepkg.TreeNode) bool {
    fmt.Printf("%s%s\n", strings.Repeat("  ", node.Depth), node.CWE.Name)
    return true
})

// Queries
path := tree.Find(79).Path()  // []*TreeNode from root to 79
leaves := tree.LeafNodes()     // []*TreeNode
maxDepth := tree.MaxDepth()    // int
count := tree.Count()           // int
isLeaf := tree.IsLeaf()        // bool
```

## Installation & Building from Source

```bash
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_linux_x86_64.tar.gz | tar xz
# Or: git clone https://github.com/scagogogo/cwe-skills.git && cd cwe-skills && go build -o cwe ./cmd/cwe/
```
