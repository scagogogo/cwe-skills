package cweskills

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetWeakness 获取指定ID的CWE弱点详情。
//
// 该方法调用 MITRE CWE REST API 的 GET /cwe/weakness/{id} 端点，
// 返回包含完整信息的弱点数据。
//
// 参数：
//   - ctx: 请求上下文，用于超时和取消控制
//   - id: CWE ID数字，例如 79
//
// 返回值：
//   - *CWE: 弱点详情
//   - error: 请求失败或解析失败时返回错误
//
// 示例：
//
//	client := cwe.NewAPIClient()
//	weakness, err := client.GetWeakness(context.Background(), 79)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(weakness.Name) // 输出: Improper Neutralization of Input During Web Page Generation ('Cross-site Scripting')
func (c *APIClient) GetWeakness(ctx context.Context, id int) (*CWE, error) {
	if id <= 0 {
		return nil, NewInvalidCWEIDError(fmt.Sprintf("%d", id))
	}

	path := fmt.Sprintf("/cwe/weakness/%d", id)

	var response struct {
		Data json.RawMessage `json:"Data"`
	}

	if err := c.httpClient.Get(ctx, path, &response); err != nil {
		return nil, fmt.Errorf("获取CWE弱点 %d 失败: %w", id, err)
	}

	var weaknesses []CWE
	if err := json.Unmarshal(response.Data, &weaknesses); err != nil {
		// 尝试作为单个弱点解析
		var weakness CWE
		if err2 := json.Unmarshal(response.Data, &weakness); err2 != nil {
			return nil, NewParseError(fmt.Sprintf("解析弱点数据失败: %v", err), 0)
		}
		weakness.CWEType = "weakness"
		return &weakness, nil
	}

	if len(weaknesses) == 0 {
		return nil, NewCWENotFoundError(id)
	}

	weaknesses[0].CWEType = "weakness"
	return &weaknesses[0], nil
}

// GetCategory 获取指定ID的CWE类别详情。
//
// 该方法调用 MITRE CWE REST API 的 GET /cwe/category/{id} 端点。
//
// 参数：
//   - ctx: 请求上下文
//   - id: 类别ID数字
//
// 返回值：
//   - *Category: 类别详情
//   - error: 请求失败时返回错误
func (c *APIClient) GetCategory(ctx context.Context, id int) (*Category, error) {
	if id <= 0 {
		return nil, NewInvalidCWEIDError(fmt.Sprintf("%d", id))
	}

	path := fmt.Sprintf("/cwe/category/%d", id)

	var response struct {
		Data json.RawMessage `json:"Data"`
	}

	if err := c.httpClient.Get(ctx, path, &response); err != nil {
		return nil, fmt.Errorf("获取CWE类别 %d 失败: %w", id, err)
	}

	var categories []Category
	if err := json.Unmarshal(response.Data, &categories); err != nil {
		var category Category
		if err2 := json.Unmarshal(response.Data, &category); err2 != nil {
			return nil, NewParseError(fmt.Sprintf("解析类别数据失败: %v", err), 0)
		}
		return &category, nil
	}

	if len(categories) == 0 {
		return nil, NewCWENotFoundError(id)
	}

	return &categories[0], nil
}

// GetView 获取指定ID的CWE视图详情。
//
// 该方法调用 MITRE CWE REST API 的 GET /cwe/view/{id} 端点。
//
// 参数：
//   - ctx: 请求上下文
//   - id: 视图ID数字，例如 1000
//
// 返回值：
//   - *View: 视图详情
//   - error: 请求失败时返回错误
func (c *APIClient) GetView(ctx context.Context, id int) (*View, error) {
	if id <= 0 {
		return nil, NewInvalidCWEIDError(fmt.Sprintf("%d", id))
	}

	path := fmt.Sprintf("/cwe/view/%d", id)

	var response struct {
		Data json.RawMessage `json:"Data"`
	}

	if err := c.httpClient.Get(ctx, path, &response); err != nil {
		return nil, fmt.Errorf("获取CWE视图 %d 失败: %w", id, err)
	}

	var views []View
	if err := json.Unmarshal(response.Data, &views); err != nil {
		var view View
		if err2 := json.Unmarshal(response.Data, &view); err2 != nil {
			return nil, NewParseError(fmt.Sprintf("解析视图数据失败: %v", err), 0)
		}
		return &view, nil
	}

	if len(views) == 0 {
		return nil, NewCWENotFoundError(id)
	}

	return &views[0], nil
}

// GetCWEs 批量获取多个CWE弱点。
//
// 该方法调用 MITRE CWE REST API 的 GET /cwe/{ids} 端点，
// 可以一次获取多个CWE弱点的信息。
//
// 参数：
//   - ctx: 请求上下文
//   - ids: CWE ID数字列表
//
// 返回值：
//   - map[string]*CWE: 以CWE ID字符串为键的弱点映射
//   - error: 请求失败时返回错误
func (c *APIClient) GetCWEs(ctx context.Context, ids []int) (map[string]*CWE, error) {
	if len(ids) == 0 {
		return make(map[string]*CWE), nil
	}

	// 构建逗号分隔的ID列表
	idStrs := make([]string, len(ids))
	for i, id := range ids {
		if id <= 0 {
			return nil, NewInvalidCWEIDError(fmt.Sprintf("%d", id))
		}
		idStrs[i] = fmt.Sprintf("%d", id)
	}

	path := fmt.Sprintf("/cwe/%s", joinIDs(idStrs))

	var response struct {
		Data json.RawMessage `json:"Data"`
	}

	if err := c.httpClient.Get(ctx, path, &response); err != nil {
		return nil, fmt.Errorf("批量获取CWE失败: %w", err)
	}

	var result map[string]*CWE
	if err := json.Unmarshal(response.Data, &result); err != nil {
		return nil, NewParseError(fmt.Sprintf("解析CWE数据失败: %v", err), 0)
	}

	for _, cwe := range result {
		if cwe != nil {
			cwe.CWEType = "weakness"
		}
	}

	return result, nil
}

// joinIDs 将ID列表用逗号连接。
func joinIDs(ids []string) string {
	result := ids[0]
	for i := 1; i < len(ids); i++ {
		result += "," + ids[i]
	}
	return result
}
