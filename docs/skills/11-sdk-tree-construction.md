# Skill: SDK — Tree Construction

## Description

Build and traverse hierarchical trees from CWE relationships. Trees represent the CWE hierarchy visually and support various traversal patterns.

## When to Use

- Visualize the CWE hierarchy
- Find all leaf or root nodes
- Calculate tree depth
- Find the path from root to a specific node
- Build a forest of disconnected sub-trees

## Creating Trees

### BuildTree — Single Root

Build a tree rooted at a specific CWE entry.

```go
tree := cwepkg.BuildTree(registry, 1) // tree rooted at CWE-1
```

### BuildForest — Multiple Pillars

Build a forest with one tree per pillar (top-level) entry.

```go
forest := cwepkg.BuildForest(registry) // []*TreeNode
```

### BuildViewTree — From a View

Build a tree from a specific CWE view's members.

```go
tree := cwepkg.BuildViewTree(registry, 1000) // Research Concepts view
```

## TreeNode Operations

### Structure

```go
type TreeNode struct {
    CWE      *cwepkg.CWE
    Children []*TreeNode
    Parent   *TreeNode
    Depth    int
}
```

### Traversal

```go
// DFS (Depth-First Search)
tree.Walk(func(node *cwepkg.TreeNode) bool {
    fmt.Printf("%s%s\n", strings.Repeat("  ", node.Depth), node.CWE.Name)
    return true // continue; return false to stop
})

// BFS (Breadth-First Search)
tree.WalkBFS(func(node *cwepkg.TreeNode) bool {
    fmt.Printf("Depth %d: %s\n", node.Depth, node.CWE.Name)
    return true
})
```

### Search

```go
// Find a node by CWE ID
node := tree.Find(79) // *TreeNode or nil

// Find a node and get its path from root
path := tree.Path(79) // []*TreeNode from root to CWE-79
```

### Properties

```go
leaves := tree.LeafNodes()   // []*TreeNode — all leaf nodes
maxDepth := tree.MaxDepth()  // int — maximum depth
count := tree.Count()         // int — total node count
isLeaf := tree.IsLeaf()      // bool — no children
isRoot := tree.IsRoot()      // bool — no parent
```

### String Representation

```go
fmt.Println(tree.String())
// Output:
// CWE-1: Root
//   CWE-2: Mid
//     CWE-4: Leaf
//   CWE-3: Mid2
```

### AddChild

```go
child := cwepkg.NewTreeNode(&cwepkg.CWE{ID: 99, Name: "New"})
tree.AddChild(child)
// child.Parent is now set to tree
// child.Depth is now tree.Depth + 1
```

## Use Cases

1. **Visualization**: Print the CWE hierarchy as an indented tree
2. **Scope Analysis**: Find all leaf weaknesses under a pillar (use `LeafNodes()`)
3. **Impact Radius**: Calculate `Count()` of a subtree to understand how many weaknesses are affected
4. **Path Tracing**: Use `Path()` to show the full chain from root to a specific weakness
5. **View Filtering**: Use `BuildViewTree()` to show only the weaknesses relevant to a specific view

## Important Notes

- `BuildTree` returns a single root node; `BuildForest` returns multiple root nodes
- Trees are built using the hierarchical relationships (ChildOf/ParentOf) only
- `BuildViewTree` uses view membership to select which CWE entries to include
- The `Walk` and `WalkBFS` callbacks receive each node; return `false` to stop traversal early
- `Find()` does a DFS search; for large trees, consider using indexed lookups on the Registry instead
- Tree nodes hold pointers to the original CWE structs — modifying them affects the Registry