# Skill: API — Get Weakness Details

## Description

Fetch detailed CWE weakness information from the MITRE CWE REST API. Requires network connectivity.

## CLI Commands

### show

Fetch one or more CWE weaknesses from the MITRE API.

```bash
cwe show CWE-79
cwe show 79 89 352
cwe show --base-url https://cwe-api.mitre.org/api CWE-79
cwe show --timeout 60 CWE-79
```

Text output:
```
=== CWE-79 ===
  名称:     Improper Neutralization of Input During Web Page Generation ('Cross-site Scripting')
  抽象层级: Base
  状态:     Stable
  描述:     The product does not neutralize...
  结构:     Simple
  关系:     4 项
```

JSON output:
```json
[
  {
    "cwe_id": "CWE-79",
    "detail": {
      "id": 79,
      "name": "Improper Neutralization of Input...",
      "abstraction": "Base",
      "status": "Stable",
      ...
    }
  }
]
```

### show category

Fetch CWE category details.

```bash
cwe show category 1
```

### show view

Fetch CWE view details.

```bash
cwe show view 1000
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--base-url` | `https://cwe-api.mitre.org/api` | MITRE API base URL |
| `--timeout` | `30` | Request timeout in seconds |

## SDK API

### NewAPIClient

```go
// Default client
client := cwepkg.NewAPIClient()
defer client.Close()

// With options
client := cwepkg.NewAPIClient(
    cwepkg.WithAPIBaseURL("https://cwe-api.mitre.org/api"),
    cwepkg.WithAPITimeout(30 * time.Second),
    cwepkg.WithAPIRateLimit(10, time.Second),
    cwepkg.WithAPIRetry(3),
)
```

### GetWeakness

```go
weakness, err := client.GetWeakness(ctx, 79)
if err != nil {
    // Handle APIError, CWENotFoundError, or RateLimitError
}

fmt.Println(weakness.Name)         // "Improper Neutralization of Input..."
fmt.Println(weakness.Abstraction)  // "Base"
fmt.Println(weakness.Status)       // "Stable"
fmt.Println(weakness.Description)  // Full description text
```

### GetCategory

```go
category, err := client.GetCategory(ctx, 1)
fmt.Println(category.Name)         // "Deprecated"
fmt.Println(category.Description)
```

### GetView

```go
view, err := client.GetView(ctx, 1000)
fmt.Println(view.Name)             // "Research Concepts"
fmt.Println(view.Type)             // "Simple"
```

### Client Options

| Option | Signature | Description |
|--------|-----------|-------------|
| `WithAPIBaseURL` | `(url string)` | Override the default API base URL |
| `WithAPITimeout` | `(timeout time.Duration)` | Set HTTP request timeout |
| `WithAPIRateLimit` | `(rate int, interval time.Duration)` | Rate limit: N requests per interval |
| `WithAPIRetry` | `(maxRetries int)` | Retry on transient errors |
| `WithAPIHTTPClient` | `(client *http.Client)` | Use a custom HTTP client |

### CWE Struct Fields

```go
type CWE struct {
    ID                  int
    Name                string
    Abstraction         Abstraction
    Structure           Structure
    Status              Status
    Description         string
    ExtendedDescription string
    LikelihoodOfExploit LikelihoodOfExploit
    CommonConsequences  []Consequence
    Relationships       []Relationship
    Mitigations         []Mitigation
    References          []Reference
    ApplicablePlatforms []PlatformEntry
    ModesOfIntroduction []Introduction
    // ... and more
}
```

## Error Handling

| Error Type | Condition |
|------------|-----------|
| `CWENotFoundError` | CWE ID does not exist in the API |
| `APIError` | HTTP error (5xx, network failure) |
| `RateLimitError` | Rate limit exceeded |
| `InvalidCWEIDError` | Invalid CWE ID format |

## Important Notes

- API calls require network connectivity to `cwe-api.mitre.org`
- Always `defer client.Close()` to release resources
- Use rate limiting for bulk queries to avoid hitting MITRE's API limits
- The API returns structured data that maps directly to the `CWE` struct