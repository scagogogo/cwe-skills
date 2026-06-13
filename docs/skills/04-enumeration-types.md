# Skill: Enumeration Types

## Description

List and understand CWE enumeration values. The CWE specification defines typed enumerations for abstraction levels, status values, relationship types, consequence scopes, and more. These are used across all SDK APIs.

## CLI Commands

### enum

List all valid values for a CWE enumeration type.

```bash
cwe enum abstraction     # Abstraction levels
cwe enum structure        # Structure types
cwe enum status           # Status values
cwe enum likelihood       # Likelihood of exploit
cwe enum relationship     # Relationship nature types
cwe enum scope            # Consequence scopes
cwe enum impact           # Consequence impact levels
cwe enum viewtype         # View types
```

Example output:
```
抽象层级 (Class/Base/Variant/Pillar) (4 项):
  - Pillar
  - Class
  - Base
  - Variant
```

JSON output:
```json
["Pillar", "Class", "Base", "Variant"]
```

## SDK API

### Parse Functions

Each enumeration has a `Parse*` function that converts a string to the typed value:

```go
abstr, err := cwepkg.ParseAbstraction("Base")      // AbstractionBase
status, err := cwepkg.ParseStatus("Stable")         // StatusStable
nature, err := cwepkg.ParseRelationshipNature("ChildOf") // RelationshipChildOf
structure, err := cwepkg.ParseStructure("Chain")     // StructureChain
likelihood, err := cwepkg.ParseLikelihoodOfExploit("High") // LikelihoodHigh
scope, err := cwepkg.ParseConsequenceScope("Confidentiality") // ScopeConfidentiality
impact, err := cwepkg.ParseConsequenceImpact("High") // ImpactHigh
viewType, err := cwepkg.ParseViewType("Graph")       // ViewTypeGraph
```

Returns an error for invalid values.

### Validation

```go
abstr.IsValid()    // bool — check if the value is a known valid value
abstr.String()     // string — canonical string representation
```

All enumeration types implement `IsValid()` and `String()`.

### All Values

```go
cwepkg.AllAbstractionValues()          // []Abstraction
cwepkg.AllStructureValues()            // []Structure
cwepkg.AllStatusValues()               // []Status
cwepkg.AllLikelihoodOfExploitValues()  // []LikelihoodOfExploit
cwepkg.AllRelationshipNatureValues()   // []RelationshipNature
cwepkg.AllConsequenceScopeValues()     // []ConsequenceScope
cwepkg.AllConsequenceImpactValues()    // []ConsequenceImpact
cwepkg.AllViewTypeValues()             // []ViewType
```

### Ordering Functions

Some enumerations have natural ordering:

```go
order := cwepkg.AbstractionOrder(cwepkg.AbstractionPillar)   // 0
order := cwepkg.AbstractionOrder(cwepkg.AbstractionClass)    // 1
order := cwepkg.AbstractionOrder(cwepkg.AbstractionBase)     // 2
order := cwepkg.AbstractionOrder(cwepkg.AbstractionVariant)  // 3

order := cwepkg.LikelihoodOrder(cwepkg.LikelihoodLow)    // 0
order := cwepkg.LikelihoodOrder(cwepkg.LikelihoodMedium) // 1
order := cwepkg.LikelihoodOrder(cwepkg.LikelihoodHigh)   // 2

order := cwepkg.ImpactOrder(cwepkg.ImpactLow)    // 0
order := cwepkg.ImpactOrder(cwepkg.ImpactMedium) // 1
order := cwepkg.ImpactOrder(cwepkg.ImpactHigh)   // 2
```

### Relationship Nature Categories

```go
nature.IsHierarchical()   // true for ChildOf, ParentOf
nature.IsSequential()     // true for CanPrecede, CanFollow
nature.IsDependency()     // true for Requires, RequiredBy
nature.IsPeer()           // true for PeerOf, CanAlsoBe
```

## Enumeration Reference

### Abstraction

| Value | Description |
|-------|-------------|
| `Pillar` | Top-level abstraction (e.g., CWE-664) |
| `Class` | Mid-level category (e.g., CWE-74 Injection) |
| `Base` | Specific weakness type (e.g., CWE-79 XSS) |
| `Variant` | Platform/language-specific variant |

### Structure

| Value | Description |
|-------|-------------|
| `Simple` | Single weakness |
| `Chain` | Sequential chain of weaknesses |
| `Composite` | Simultaneous combination of weaknesses |

### Status

| Value | Description |
|-------|-------------|
| `Stable` | Fully reviewed and accepted |
| `Usable` | Sufficient detail for use |
| `Draft` | Under development |
| `Incomplete` | Partially defined |
| `Obsolete` | No longer relevant |
| `Deprecated` | Replaced or removed |

### LikelihoodOfExploit

| Value | Order |
|-------|-------|
| `Low` | 0 |
| `Medium` | 1 |
| `High` | 2 |
| `Unknown` | -1 |

### RelationshipNature

| Value | Category |
|-------|----------|
| `ChildOf` | Hierarchical |
| `ParentOf` | Hierarchical |
| `CanPrecede` | Sequential |
| `CanFollow` | Sequential |
| `Requires` | Dependency |
| `RequiredBy` | Dependency |
| `CanAlsoBe` | Peer |
| `PeerOf` | Peer |
| `MemberOf` | Membership |
| `HasMember` | Membership |