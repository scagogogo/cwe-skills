# Skill: API — Relationship Queries

## Description

Query CWE relationships (parent/child/ancestor/descendant) from the MITRE CWE REST API. Understanding these relationships is essential for navigating the CWE hierarchy.

## CLI Commands

### relations parents

Query the parent weaknesses of a CWE entry.

```bash
cwe relations parents CWE-79
cwe relations parents CWE-79 --view-id 1000
```

Output:
```
CWE-79 的 父级弱点 (1 项):
  ChildOf -> CWE-74 (View: 1000)
```

### relations children

Query the child weaknesses of a CWE entry.

```bash
cwe relations children CWE-74
```

### relations ancestors

Query all ancestor weaknesses (transitive parents up to root).

```bash
cwe relations ancestors CWE-79
```

### relations descendants

Query all descendant weaknesses (transitive children down to leaves).

```bash
cwe relations descendants CWE-74
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--base-url` | `https://cwe-api.mitre.org/api` | MITRE API base URL |
| `--view-id` | `0` | Filter by view ID (only for parents/children) |

## SDK API

### GetParents

```go
// Without view filter
parents, err := client.GetParents(ctx, 79)

// With view filter
parents, err := client.GetParents(ctx, 79, 1000)
```

### GetChildren

```go
children, err := client.GetChildren(ctx, 74)
children, err := client.GetChildren(ctx, 74, 1000) // with view
```

### GetAncestors

```go
ancestors, err := client.GetAncestors(ctx, 79)
```

Note: `GetAncestors` does not support the `viewID` parameter.

### GetDescendants

```go
descendants, err := client.GetDescendants(ctx, 74)
```

Note: `GetDescendants` does not support the `viewID` parameter.

### Relationship Struct

```go
type Relationship struct {
    Nature RelationshipNature  // ChildOf, ParentOf, etc.
    CWEID  int                 // The related CWE ID
    ViewID int                 // View context (0 if none)
}
```

### RelationshipNature Categories

| Category | Natures | Description |
|----------|---------|-------------|
| Hierarchical | `ChildOf`, `ParentOf` | Parent-child in the CWE tree |
| Sequential | `CanPrecede`, `CanFollow` | Temporal ordering |
| Dependency | `Requires`, `RequiredBy` | Functional dependency |
| Peer | `PeerOf`, `CanAlsoBe` | Related alternatives |
| Membership | `MemberOf`, `HasMember` | Category/view membership |

## Use Cases

1. **Impact Analysis**: Find all descendants of a base weakness to understand scope
2. **Root Cause Analysis**: Walk ancestors to find the root category
3. **View Filtering**: Use `--view-id 1000` (Research Concepts) for the standard hierarchy
4. **Compliance Mapping**: Trace relationships between specific weaknesses and categories

## Important Notes

- `GetParents`/`GetChildren` support optional `viewID` via variadic parameter
- `GetAncestors`/`GetDescendants` do not support view filtering
- The relationship nature in the response tells you the type (ChildOf, ParentOf, etc.)
- Ancestors are transitive — they include all levels up to the root