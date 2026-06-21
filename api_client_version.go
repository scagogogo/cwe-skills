package cweskills

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetVersion 获取当前CWE数据的版本信息。
//
// 该方法调用 MITRE CWE REST API 的 GET /version 端点，
// 返回当前CWE数据库的版本号和发布日期。
//
// 参数：
//   - ctx: 请求上下文
//
// 返回值：
//   - *VersionResponse: 版本信息
//   - error: 请求失败时返回错误
//
// 示例：
//
//	client := cwe.NewAPIClient()
//	version, err := client.GetVersion(context.Background())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("CWE Version: %s, Release Date: %s\n", version.Version, version.ReleaseDate)
func (c *APIClient) GetVersion(ctx context.Context) (*VersionResponse, error) {
	var response struct {
		Data json.RawMessage `json:"Data"`
	}

	if err := c.httpClient.Get(ctx, "/version", &response); err != nil {
		return nil, fmt.Errorf("获取CWE版本信息失败: %w", err)
	}

	var version VersionResponse
	if err := json.Unmarshal(response.Data, &version); err != nil {
		return nil, NewParseError(fmt.Sprintf("解析版本数据失败: %v", err), 0)
	}

	return &version, nil
}
