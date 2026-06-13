# Skill: CWE ID Parsing & Validation

## Description

Parse, validate, and format CWE IDs. This is the foundational skill — every other CWE operation starts with a valid CWE ID.

CWE IDs follow the format `CWE-NNN` where NNN is a positive integer. The SDK accepts multiple input formats (e.g., `79`, `CWE-79`, `cwe-79`) and always normalizes to the canonical `CWE-NNN` form.

## CLI Commands

### parse

Parse CWE IDs and extract the numeric part.

```bash
cwe parse CWE-79 89 cwe-352
```

Output:
```
CWE-79 -> CWE-79 (ID: 79)
89 -> CWE-79 (ID: 89)
cwe-352 -> CWE-352 (ID: 352)
```

JSON output (`-o json`):
```json
[
  {"input": "CWE-79", "id": 79, "format": "CWE-79", "valid": true},
  {"input": "89", "id": 89, "format": "CWE-89", "valid": true}
]
```

Invalid inputs are reported but not fatal:
```json
{"input": "abc", "id": 0, "valid": false, "error": "cwe: [INVALID_CWE_ID] CWE ID格式无效: 输入值: \"abc\""}
```

### validate

Validate whether input strings are valid CWE IDs.

```bash
cwe validate CWE-79 CWE-89 abc
```

Output:
```
CWE-79 ✓ 有效
CWE-89 ✓ 有效
abc ✗ 无效: cwe: [INVALID_CWE_ID] CWE ID格式无效: 输入值: "abc"
部分CWE ID无效
```

Exit code 0 if all valid, 1 if any invalid.

### format

Format CWE IDs to the canonical `CWE-NNN` form.

```bash
cwe format 79 cwe-89 CWE-352
```

Output:
```
CWE-79
CWE-89
CWE-352
```

## SDK API

### ParseCWEID

```go
id, err := cwepkg.ParseCWEID("CWE-79")
// id = 79, err = nil

id, err := cwepkg.ParseCWEID("abc")
// id = 0, err = InvalidCWEIDError
```

Accepts: `"CWE-79"`, `"cwe-79"`, `"79"` — case-insensitive, optional prefix.

### FormatCWEID / FormatCWEIDFromInt

```go
formatted, err := cwepkg.FormatCWEID("79")     // "CWE-79"
formatted := cwepkg.FormatCWEIDFromInt(79)       // "CWE-79"
```

### IsCWEID / ValidateCWEID

```go
if cwepkg.IsCWEID("CWE-89") {  // true — quick check
    // valid
}

err := cwepkg.ValidateCWEID("abc")  // returns InvalidCWEIDError
err := cwepkg.ValidateCWEID("CWE-79")  // nil — valid
```

`IsCWEID` is a boolean quick-check. `ValidateCWEID` returns a structured error with details.

### CompareCWEIDs

```go
result, err := cwepkg.CompareCWEIDs("CWE-79", "CWE-89")
// result < 0 means CWE-79 < CWE-89
// result == 0 means equal
// result > 0 means greater
```

## Error Types

| Error | Condition |
|-------|-----------|
| `InvalidCWEIDError` | Input doesn't match CWE ID pattern (empty, no digits, negative) |

## Important Notes

- `ParseCWEID` is case-insensitive — `"CWE-79"` and `"cwe-79"` both return `79`
- Pure numeric input `"79"` is accepted
- Zero and negative IDs are rejected (`CWE-0`, `CWE--1`)
- `IsCWEID` performs a quick regex check; `ValidateCWEID` provides detailed error info