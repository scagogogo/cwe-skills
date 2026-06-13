# Skill: SDK — In-Memory Registry

## Description

The Registry is the core in-memory data store for CWE entries. It supports storing weaknesses, categories, views, and compound elements, with automatic relationship indexing and concurrency-safe access.

This is an SDK-only skill — there is no direct CLI command for registry operations (use `search`/`stats` with XML files for CLI-level local operations).

## When to Use

- Store CWE entries locally for fast repeated queries
- Build relationship indexes for navigation
- Import/export CWE data (JSON)
- Cache API results for offline use

## Creating and Populating

```go
registry := cwepkg.NewRegistry()

// Register a weakness
err := registry.Register(&cwepkg.CWE{
    ID:          79,
    Name:        "Cross-site Scripting (XSS)",
    Abstraction: cwepkg.AbstractionBase,
    Status:      cwepkg.StatusStable,
    Relationships: []cwepkg.Relationship{
        {Nature: cwepkg.RelationshipChildOf, CWEID: 74},
    },
})

// Register a category
err := registry.RegisterCategory(&cwepkg.Category{
    ID:   1,
    Name: "Deprecated",
    Relationships: []cwepkg.Relationship{
        {Nature: cwepkg.RelationshipHasMember, CWEID: 79},
    },
})

// Register a view
err := registry.RegisterView(&cwepkg.View{
    ID:   1000,
    Name: "Research Concepts",
    Type: cwepkg.ViewTypeGraph,
    Members: []cwepkg.ViewMember{
        {CWEID: 79, ViewID: 1000, Direct: true},
    },
})

// Register a compound element
err := registry.RegisterCompoundElement(&cwepkg.CompoundElement{
    ID:        680,
    Name:      "Integer Overflow to Buffer Overflow",
    Structure: cwepkg.StructureChain,
})
```

## Building Indexes

After registering entries, call `BuildIndexes()` to create relationship indexes:

```go
registry.BuildIndexes()
```

This populates indexes for:
- `parentIndex` — ChildOf → parent mapping
- `childIndex` — ParentOf → child mapping
- `peerIndex` — PeerOf/CanAlsoBe → peer mapping
- `memberIndex` — HasMember → member mapping
- `memberOfIndex` — MemberOf → group mapping

Indexes are automatically deduplicated. Registering new entries invalidates indexes (call `BuildIndexes()` again).

## Querying

### Get by ID

```go
cwe, ok := registry.Get(79)                    // *CWE
cat, ok := registry.GetCategory(1)             // *Category
view, ok := registry.GetView(1000)             // *View
ce, ok := registry.GetCompoundElement(680)      // *CompoundElement
```

### Get all entries

```go
allCWEs := registry.GetAll()              // []*CWE
allCats := registry.GetAllCategories()     // []*Category
allViews := registry.GetAllViews()         // []*View
```

### Counts and membership

```go
count := registry.Size()                  // number of weaknesses
count := registry.CategoryCount()         // number of categories
count := registry.ViewCount()             // number of views
count := registry.CompoundElementCount()   // number of compound elements
exists := registry.Contains(79)           // true if CWE-79 exists
```

### Relationship queries (after BuildIndexes)

```go
parentIDs := registry.GetParentIDs(79)       // []int — direct parents
childIDs := registry.GetChildIDs(74)        // []int — direct children
peerIDs := registry.GetPeerIDs(79)           // []int — peers
ancestorIDs := registry.GetAncestorIDs(79)   // []int — all ancestors (transitive)
descendantIDs := registry.GetDescendantIDs(74) // []int — all descendants (transitive)
viewMembers := registry.GetViewMembers(1000)  // []int — members of view
catMembers := registry.GetCategoryMembers(1)  // []int — members of category
memberOfIDs := registry.GetMemberOfIDs(79)    // []int — what 79 belongs to
```

## Modification

### Remove entries

```go
err := registry.Remove(79)             // remove a weakness
err := registry.RemoveCategory(1)      // remove a category
err := registry.RemoveView(1000)        // remove a view
```

### Clear all

```go
registry.Clear()  // removes all entries and indexes
```

### Check index status

```go
if registry.IndexesBuilt() {
    // safe to query relationships
}
```

## Import/Export

### JSON round-trip

```go
data, err := registry.ExportJSON()

newRegistry := cwepkg.NewRegistry()
err = newRegistry.ImportJSON(data)
```

### CSV export

```go
csvData, err := registry.ExportCSV()
```

## Concurrency

The Registry is safe for concurrent read access after `BuildIndexes()`. Write operations (Register, Remove, Clear) acquire a write lock. For concurrent read-heavy workloads:

```go
// Build indexes once, then read concurrently
registry.BuildIndexes()

// Multiple goroutines can safely call:
//   registry.Get(), registry.GetParentIDs(), etc.
```

## Error Handling

| Error | Condition |
|-------|-----------|
| `InvalidCWEIDError` | ID is 0 or negative |
| `CWENotFoundError` | Remove called on non-existent ID |
| — | Duplicate ID (Register returns an error for existing IDs) |

## Important Notes

- Always call `BuildIndexes()` after registering entries before querying relationships
- Registering new entries invalidates indexes — call `BuildIndexes()` again
- The registry uses value copies for relationship index queries (safe to modify returned slices)
- `GetAncestorIDs` and `GetDescendantIDs` traverse the full tree (may be slow for very deep hierarchies)
- Import/Export JSON preserves all entry types (weaknesses, categories, views, compound elements)