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
abstr, err := cweskills.ParseAbstraction("Base")      // AbstractionBase
status, err := cweskills.ParseStatus("Stable")         // StatusStable
nature, err := cweskills.ParseRelationshipNature("ChildOf") // RelationshipChildOf
structure, err := cweskills.ParseStructure("Chain")     // StructureChain
likelihood, err := cweskills.ParseLikelihoodOfExploit("High") // LikelihoodHigh
scope, err := cweskills.ParseConsequenceScope("Confidentiality") // ScopeConfidentiality
impact, err := cweskills.ParseConsequenceImpact("High") // ImpactHigh
viewType, err := cweskills.ParseViewType("Graph")       // ViewTypeGraph
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
cweskills.AllAbstractionValues()          // []Abstraction
cweskills.AllStructureValues()            // []Structure
cweskills.AllStatusValues()               // []Status
cweskills.AllLikelihoodOfExploitValues()  // []LikelihoodOfExploit
cweskills.AllRelationshipNatureValues()   // []RelationshipNature
cweskills.AllConsequenceScopeValues()     // []ConsequenceScope
cweskills.AllConsequenceImpactValues()    // []ConsequenceImpact
cweskills.AllViewTypeValues()             // []ViewType
```

### Ordering Functions

Some enumerations have natural ordering:

```go
order := cweskills.AbstractionOrder(cweskills.AbstractionPillar)   // 0
order := cweskills.AbstractionOrder(cweskills.AbstractionClass)    // 1
order := cweskills.AbstractionOrder(cweskills.AbstractionBase)     // 2
order := cweskills.AbstractionOrder(cweskills.AbstractionVariant)  // 3

order := cweskills.LikelihoodOrder(cweskills.LikelihoodLow)    // 0
order := cweskills.LikelihoodOrder(cweskills.LikelihoodMedium) // 1
order := cweskills.LikelihoodOrder(cweskills.LikelihoodHigh)   // 2

order := cweskills.ImpactOrder(cweskills.ImpactLow)    // 0
order := cweskills.ImpactOrder(cweskills.ImpactMedium) // 1
order := cweskills.ImpactOrder(cweskills.ImpactHigh)   // 2
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