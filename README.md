# CWE SDK for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/scagogogo/cwe.svg)](https://pkg.go.dev/github.com/scagogogo/cwe)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/cwe)](https://goreportcard.com/report/github.com/scagogogo/cwe)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A comprehensive Go SDK for [CWE (Common Weakness Enumeration)](https://cwe.mitre.org/), providing complete API support for building cybersecurity products.

## Features

- **Complete CWE Data Model**: Weaknesses, Categories, Views, and Compound Elements with full field coverage
- **Typed Enumerations**: Abstraction levels, Status values, Relationship types, Consequence scopes, etc.
- **CWE ID Utilities**: Parse, format, validate, and extract CWE IDs from text
- **Well-Known Lists**: CWE Top 25, OWASP Top 10, SANS Top 25 with membership checks
- **MITRE REST API Client**: Full access to the CWE API with rate limiting and retry
- **XML Catalog Parser**: Offline parsing of MITRE's official XML downloads
- **In-Memory Registry**: Store, index, and query CWE entries with relationship indexes
- **Search & Filter**: By keyword, abstraction, status, likelihood, consequence scope, and more
- **Relationship Navigation**: Parents, children, ancestors, descendants, peers, chains, composites
- **Tree Construction**: Build hierarchical trees from CWE relationships
- **Serialization**: JSON, XML, and CSV import/export
- **Zero Dependencies**: Uses only the Go standard library

## Installation

```bash
go get github.com/scagogogo/cwe
```

## Quick Start

### Parse and Validate CWE IDs

```go
package main

import (
    "fmt"
    "github.com/scagogogo/cwe"
)

func main() {
    // Parse CWE ID
    id, _ := cwe.ParseCWEID("CWE-79")
    fmt.Println(id) // 79

    // Format CWE ID
    formatted, _ := cwe.FormatCWEID("79")
    fmt.Println(formatted) // CWE-79

    // Validate
    if cwe.IsCWEID("CWE-89") {
        fmt.Println("Valid CWE ID")
    }

    // Extract from text
    ids := cwe.ExtractCWEIDs("See CWE-79 and CWE-89 for details")
    fmt.Println(ids) // [CWE-79 CWE-89]
}
```

### Query the MITRE CWE API

```go
client := cwe.NewAPIClient()

// Get a weakness
weakness, err := client.GetWeakness(ctx, 79)

// Get version
version, err := client.GetVersion(ctx)

// Get relationships
parents, err := client.GetParents(ctx, 79)
children, err := client.GetChildren(ctx, 79, 1000) // with view ID
```

### Use the Registry for Local Operations

```go
registry := cwe.NewRegistry()

// Register CWE entries
registry.Register(&cwe.CWE{
    ID:          79,
    Name:        "Cross-site Scripting (XSS)",
    Abstraction: cwe.AbstractionBase,
    Status:      cwe.StatusStable,
})

// Build relationship indexes
registry.BuildIndexes()

// Search and filter
results := cwe.FindByAbstraction(registry, cwe.AbstractionBase)
filtered := cwe.Filter(results, cwe.FilterOption{Status: cwe.StatusStable})

// Navigate relationships
nav := cwe.NewNavigator(registry)
parents := nav.Parents(79)
ancestors := nav.Ancestors(79)
```

### Parse Offline XML Catalog

```go
parser := cwe.NewXMLParser()
registry, err := parser.ParseFile("cwec_v4.10.xml")
```

### Check Well-Known Lists

```go
if cwe.IsInTop25(79) {
    fmt.Println("CWE-79 is in the Top 25!")
}

category := cwe.GetOWASPCategory(79)
fmt.Println(category) // A03:2021-Injection
```

## API Reference

### Core Types

| Type | Description |
|------|-------------|
| `CWE` | Core weakness entry with all CWE fields |
| `Category` | CWE category with members |
| `View` | CWE view with members and type |
| `CompoundElement` | Chain or composite weakness |
| `Relationship` | Relationship between CWE entries |
| `Consequence` | Impact consequence with scope and impact |

### Enumerations

| Type | Values |
|------|--------|
| `Abstraction` | Pillar, Class, Base, Variant |
| `Structure` | Simple, Chain, Composite |
| `Status` | Stable, Usable, Draft, Incomplete, Obsolete, Deprecated |
| `LikelihoodOfExploit` | High, Medium, Low, Unknown |
| `RelationshipNature` | ChildOf, ParentOf, CanPrecede, CanFollow, Requires, RequiredBy, CanAlsoBe, PeerOf, MemberOf, HasMember |
| `ConsequenceScope` | Confidentiality, Integrity, Availability, Access Control, etc. |
| `ViewType` | Graph, Explicit Slice, Implicit Slice |

### Key Functions

| Function | Description |
|----------|-------------|
| `ParseCWEID(s)` | Parse CWE ID from string |
| `FormatCWEID(s)` | Format to "CWE-NNN" |
| `IsCWEID(s)` | Check if valid CWE ID |
| `IsInTop25(id)` | Check if in CWE Top 25 |
| `IsInOWASPTop10(id)` | Check if in OWASP Top 10 |
| `FindByKeyword(r, kw)` | Search by keyword |
| `Filter(cwes, opts)` | Filter by multiple criteria |
| `BuildTree(r, id)` | Build hierarchy tree |
| `MarshalJSON/UnmarshalJSON` | JSON serialization |
| `MarshalXML/UnmarshalXML` | XML serialization |
| `MarshalCSV/UnmarshalCSV` | CSV serialization |

## License

MIT License - see [LICENSE](LICENSE) for details.
