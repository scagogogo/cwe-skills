# Skill: SDK — Serialization

## Description

Import and export CWE data in JSON, XML, and CSV formats. Supports both individual entries and bulk operations through the Registry.

## When to Use

- Persist CWE data to disk
- Exchange data with other systems
- Import MITRE XML catalogs for offline processing
- Export filtered results for reporting

## JSON Serialization

### Individual CWE

```go
cwe := &cwepkg.CWE{ID: 79, Name: "XSS", Abstraction: cwepkg.AbstractionBase}

// Marshal
data, err := cwepkg.MarshalJSON(cwe)

// Unmarshal
parsed, err := cwepkg.UnmarshalJSON(data)
```

### CWE List

```go
cwes := []*cwepkg.CWE{cwe1, cwe2}

// Marshal list
data, err := cwepkg.MarshalJSONList(cwes)

// Unmarshal list
parsed, err := cwepkg.UnmarshalJSONList(data)
```

### Registry JSON Round-Trip

```go
// Export entire registry to JSON
data, err := registry.ExportJSON()

// Import JSON into a new registry
newRegistry := cwepkg.NewRegistry()
err = newRegistry.ImportJSON(data)
```

Export JSON preserves all entry types: weaknesses, categories, views, and compound elements.

## XML Serialization

### Individual CWE

```go
// Marshal to XML
data, err := cwepkg.MarshalXML(cwe)

// Unmarshal from XML
parsed, err := cwepkg.UnmarshalXML(data)
```

The XML format includes the standard `<?xml version="1.0" encoding="UTF-8"?>` header.

### MITRE XML Catalog

```go
parser := cwepkg.NewXMLParser()
registry, err := parser.ParseFile("cwec_v4.15.xml")
// Or from a reader:
registry, err := parser.Parse(reader)
// Or from bytes:
registry, err := parser.ParseBytes(xmlData)
```

The XML parser handles the official MITRE CWE catalog format, converting all entry types (Weaknesses, Categories, Views, Compound_Elements) into Registry entries.

## CSV Serialization

### CWE List

```go
cwes := []*cwepkg.CWE{cwe1, cwe2, cwe3}

// Marshal to CSV
data, err := cwepkg.MarshalCSV(cwes)
// Returns: "ID,Name,Abstraction,Status,Structure,Description,CWEType\n79,XSS,Base,...\n"

// Unmarshal from CSV
parsed, err := cwepkg.UnmarshalCSV(data)
```

### Registry CSV Export

```go
csvData, err := registry.ExportCSV()
```

### CSV Format

```csv
ID,Name,Abstraction,Status,Structure,Description,CWEType
79,Cross-site Scripting,Base,Stable,Simple,The product does not...,weakness
89,SQL Injection,Base,Stable,Simple,The product constructs...,weakness
```

- Fields with commas or quotes are properly escaped
- Variable column counts are handled (missing columns default to empty)
- Non-numeric IDs in CSV rows are skipped during unmarshal

## Consequence Struct

```go
consequence := cwepkg.Consequence{
    Scopes:  []cwepkg.ConsequenceScope{cwepkg.ScopeConfidentiality, cwepkg.ScopeIntegrity},
    Impacts: []cwepkg.ConsequenceImpact{cwepkg.ImpactHigh},
    Note:    "Details about the impact",
}

// Check if a scope is present
hasConf := consequence.HasScope(cwepkg.ScopeConfidentiality) // true

// Check if an impact is present
hasHigh := consequence.HasImpact(cwepkg.ImpactHigh) // true

// Get maximum impact level
maxImpact := consequence.MaxImpact() // ImpactHigh
```

## Error Handling

| Error | Condition |
|-------|-----------|
| `ParseError` | Invalid format during unmarshal |
| `ValidationError` | Invalid field values |
| — | `ImportJSON` returns error for malformed JSON |

## Important Notes

- JSON serialization handles `interface{}` fields (e.g., in show results) gracefully
- XML unmarshal requires the `CWE` XML tag structure; for MITRE catalog XML, use `XMLParser`
- CSV is lossy — only basic fields are preserved (ID, Name, Abstraction, Status, Structure, Description, CWEType)
- CSV unmarshal is lenient: rows with fewer columns than the header are still parsed (missing fields default to empty)
- Registry JSON export/import preserves all data types; CSV export only includes weaknesses
- All serialization functions return byte slices — write to file using `os.WriteFile`