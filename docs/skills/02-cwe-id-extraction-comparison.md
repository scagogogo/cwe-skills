# Skill: CWE ID Extraction & Comparison

## Description

Extract CWE IDs from free-form text and compare CWE IDs. Useful for processing vulnerability reports, security advisories, and compliance documents.

## CLI Commands

### extract

Extract all CWE IDs from text.

```bash
cwe extract "This system is affected by CWE-79 and CWE-89"
```

Output:
```
找到 2 个CWE ID:
  CWE-79
  CWE-89
```

JSON output:
```json
{
  "text": "This system is affected by CWE-79 and CWE-89",
  "ids": ["CWE-79", "CWE-89"],
  "count": 2
}
```

### compare

Compare two CWE IDs.

```bash
cwe compare CWE-79 CWE-89
cwe compare CWE-79 CWE-79
cwe compare CWE-89 CWE-79
```

Output:
```
CWE-79 is less than CWE-89
CWE-79 is equal to CWE-79
CWE-89 is greater than CWE-79
```

## SDK API

### ExtractCWEIDs

```go
ids := cwepkg.ExtractCWEIDs("Affected by CWE-79, CWE-89, and CWE-352")
// ids = ["CWE-79", "CWE-89", "CWE-352"]
```

Returns all `CWE-NNN` patterns found in the text. Case-insensitive matching.

### ExtractFirstCWEID

```go
id := cwepkg.ExtractFirstCWEID("See CWE-79 and CWE-89")
// id = "CWE-79"
```

Returns the first CWE ID found, or empty string if none.

### CompareCWEIDs

```go
result, err := cwepkg.CompareCWEIDs("CWE-79", "CWE-89")
// result = -1 (CWE-79 < CWE-89)

result, err := cwepkg.CompareCWEIDs("CWE-79", "CWE-79")
// result = 0 (equal)

result, err := cwepkg.CompareCWEIDs("CWE-89", "CWE-79")
// result = 1 (CWE-89 > CWE-79)
```

Returns `-1`, `0`, or `1` for less than, equal, or greater than.

## Use Cases

1. **Vulnerability Report Parsing**: Extract all referenced CWE IDs from security advisories
2. **Compliance Checking**: Compare CWE IDs to determine ordering or grouping
3. **Data Normalization**: Extract and normalize CWE references from free-text fields

## Important Notes

- `ExtractCWEIDs` uses regex pattern matching for `CWE-\d+` (case-insensitive)
- `CompareCWEIDs` compares the numeric parts — `CWE-79` < `CWE-89` because 79 < 89
- Both functions accept the standard CWE ID formats (`CWE-79`, `79`, `cwe-79`)