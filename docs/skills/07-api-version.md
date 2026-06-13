# Skill: API — Version Check

## Description

Check the version of the MITRE CWE REST API. Useful for verifying API connectivity and determining the current CWE data version.

## CLI Commands

### api-version

```bash
cwe api-version
cwe api-version --base-url https://cwe-api.mitre.org/api
```

Text output:
```
MITRE CWE API版本: 4.15
发布日期: 2024-11-19
版本名称: CWE v4.15
```

JSON output:
```json
{
  "version": "4.15",
  "releaseDate": "2024-11-19",
  "name": "CWE v4.15"
}
```

## SDK API

### GetVersion

```go
version, err := client.GetVersion(ctx)
if err != nil {
    // Handle APIError
}

fmt.Println(version.Version)     // "4.15"
fmt.Println(version.ReleaseDate) // "2024-11-19"
fmt.Println(version.Name)        // "CWE v4.15"
```

### VersionResponse Struct

```go
type VersionResponse struct {
    Version     string `json:"version"`
    ReleaseDate string `json:"releaseDate"`
    Name        string `json:"name"`
}
```

## Important Notes

- This is a lightweight call — good for connectivity checks
- No rate limiting concern for occasional version checks