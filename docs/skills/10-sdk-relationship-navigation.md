# Skill: SDK — Relationship Navigation

## Description

Navigate CWE relationships using the Navigator. Provides high-level methods for traversing parent/child hierarchies, finding paths between CWE entries, and determining relationship depth.

Requires a Registry with indexes built (`registry.BuildIndexes()`).

## When to Use

- Find all ancestors or descendants of a weakness
- Determine if two CWE entries are related
- Calculate the distance between CWE entries
- Find the shortest path through the CWE graph
- List siblings, peers, or chain members

## Creating a Navigator

```go
registry := cwepkg.NewRegistry()
// ... register entries ...
registry.BuildIndexes()

nav := cwepkg.NewNavigator(registry)
```

## Core Navigation Methods

### Parents & Children

```go
parents := nav.Parents(79)      // []*CWE — direct parents
children := nav.Children(74)    // []*CWE — direct children
```

### Ancestors & Descendants

```go
ancestors := nav.Ancestors(79)    // []*CWE — all ancestors (transitive)
descendants := nav.Descendants(74) // []*CWE — all descendants (transitive)
```

### Siblings

```go
siblings := nav.Siblings(79)      // []*CWE — same parent, different ID
```

### Peer Relationships

```go
peers := nav.Peers(79)            // []*CWE — PeerOf and CanAlsoBe
canPrecede := nav.CanPrecede(79)   // []*CWE — what 79 can precede
canFollow := nav.CanFollow(79)     // []*CWE — what 79 can follow
requires := nav.Requires(79)      // []*CWE — what 79 requires
requiredBy := nav.RequiredBy(79)  // []*CWE — what requires 79
canAlsoBe := nav.CanAlsoBe(79)    // []*CWE — what 79 can also be
```

### Compound Elements

```go
chainMembers := nav.ChainMembers(680)        // []*CWE — members of a chain
compositeMembers := nav.CompositeMembers(680) // []*CWE — members of a composite
```

## Path Finding

### ShortestPath

Find the shortest path between two CWE entries through any relationship type.

```go
path := nav.ShortestPath(79, 1) // []int{79, 74, 1}
```

Returns the sequence of CWE IDs forming the shortest path, or `nil` if no path exists.

### Relationship Checks

```go
// Is CWE-1 an ancestor of CWE-79?
yes := nav.IsAncestorOf(1, 79)    // true

// Is CWE-79 a descendant of CWE-1?
yes := nav.IsDescendantOf(79, 1)  // true

// Are CWE-79 and CWE-89 related in any way?
yes := nav.IsRelated(79, 89)      // true if any path exists
```

### RelationshipDepth

Calculate the minimum number of hops between two CWE entries.

```go
depth := nav.RelationshipDepth(79, 1)  // 2 (79 -> 74 -> 1)
depth := nav.RelationshipDepth(79, 79) // 0 (same entry)
```

Returns `-1` if no relationship exists.

## String Representation

```go
s := nav.String() // Human-readable summary of all relationships for the registry
```

## Method Reference

| Method | Returns | Description |
|--------|---------|-------------|
| `Parents(id)` | `[]*CWE` | Direct parents (ChildOf) |
| `Children(id)` | `[]*CWE` | Direct children (ParentOf) |
| `Ancestors(id)` | `[]*CWE` | All ancestors (transitive) |
| `Descendants(id)` | `[]*CWE` | All descendants (transitive) |
| `Siblings(id)` | `[]*CWE` | Same parent, different ID |
| `Peers(id)` | `[]*CWE` | PeerOf + CanAlsoBe |
| `CanPrecede(id)` | `[]*CWE` | CanPrecede relationships |
| `CanFollow(id)` | `[]*CWE` | CanFollow relationships |
| `Requires(id)` | `[]*CWE` | Requires relationships |
| `RequiredBy(id)` | `[]*CWE` | RequiredBy relationships |
| `CanAlsoBe(id)` | `[]*CWE` | CanAlsoBe relationships |
| `ChainMembers(id)` | `[]*CWE` | Members of a chain |
| `CompositeMembers(id)` | `[]*CWE` | Members of a composite |
| `ShortestPath(from, to)` | `[]int` | Shortest path or nil |
| `IsAncestorOf(a, b)` | `bool` | Is `a` an ancestor of `b`? |
| `IsDescendantOf(a, b)` | `bool` | Is `a` a descendant of `b`? |
| `IsRelated(a, b)` | `bool` | Any relationship between `a` and `b`? |
| `RelationshipDepth(a, b)` | `int` | Min hops, or -1 if unrelated |

## Important Notes

- All navigation methods require indexes to be built first
- `ShortestPath` uses BFS and traverses all relationship types (hierarchical, peer, sequential, dependency)
- `GetAncestorIDs`/`GetDescendantIDs` on Registry use DFS; `Ancestors`/`Descendants` on Navigator also use DFS
- `IsAncestorOf` and `IsDescendantOf` only consider hierarchical (ChildOf/ParentOf) relationships
- `IsRelated` and `ShortestPath` consider all relationship types
- For nil registry, all methods return nil/zero/false safely