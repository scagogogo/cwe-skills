# Skill: Well-Known Lists

## Description

Query and check membership in well-known CWE lists: CWE Top 25 Most Dangerous Software Weaknesses, OWASP Top 10 (2021), and SANS Top 25 Most Dangerous Software Errors. These lists are critical for security prioritization and compliance.

## CLI Commands

### wellknown top25

List the CWE Top 25 Most Dangerous Software Weaknesses.

```bash
cwe wellknown top25
```

Output:
```
CWE Top 25 Most Dangerous Software Weaknesses (25 项):

  1. CWE-79
  2. CWE-89
  3. CWE-352
  ...
```

### wellknown owasp

List OWASP Top 10 (2021) categories and their associated CWE IDs.

```bash
cwe wellknown owasp
```

Output:
```
OWASP Top 10 (2021):

  A01:2021-Broken Access Control:
    CWE-862
    ...
  A03:2021-Injection:
    CWE-79
    CWE-89
    ...
```

### wellknown sans

List SANS Top 25 Most Dangerous Software Errors.

```bash
cwe wellknown sans
```

### wellknown check

Check whether specific CWE IDs belong to any well-known list.

```bash
cwe wellknown check CWE-79 CWE-89 CWE-999
```

Output:
```
CWE-79: [Top 25 OWASP Top 10 (A03:2021-Injection) SANS Top 25]
CWE-89: [Top 25 OWASP Top 10 (A03:2021-Injection) SANS Top 25]
CWE-999: 不在任何知名列表中
```

JSON output:
```json
[
  {"cwe_id": "CWE-79", "in_list": ["Top 25", "OWASP Top 10 (A03:2021-Injection)", "SANS Top 25"]},
  {"cwe_id": "CWE-999", "in_list": []}
]
```

## SDK API

### Membership Checks

```go
cwepkg.IsInTop25(79)       // true
cwepkg.IsInOWASPTop10(79)  // true
cwepkg.IsInSANSTop25(79)   // true
```

### OWASP Category

```go
category := cwepkg.GetOWASPCategory(79)
// "A03:2021-Injection"

categories := cwepkg.GetOWASPCategories(79)
// ["A03:2021-Injection"]
```

### View Membership

```go
cwepkg.IsInWellKnownView(1000)  // true — View 1000 is a well-known research view
```

### Pre-Defined Lists

```go
// Direct access to the ID slices
top25 := cwepkg.CWETop25       // []int{79, 89, 352, 862, ...}
sans25 := cwepkg.SANSTop25      // []int{119, 20, 79, ...}

// OWASP Top 10 mapping
owasp := cwepkg.OWASPTop10      // map[string][]int
// Key example: "A03:2021-Injection" -> []int{79, 89, ...}
```

## Use Cases

1. **Security Prioritization**: Focus remediation on Top 25 weaknesses
2. **Compliance Reporting**: Map vulnerabilities to OWASP Top 10 categories
3. **Risk Scoring**: Weight findings higher if they appear in well-known lists
4. **Filtering**: Automatically flag high-priority CWE entries

## Important Notes

- All checks accept integer CWE IDs (not strings)
- A single CWE ID can belong to multiple lists simultaneously (e.g., CWE-79 is in all three)
- OWASP Top 10 uses the 2021 edition
- `GetOWASPCategory` returns the first matching category; `GetOWASPCategories` returns all