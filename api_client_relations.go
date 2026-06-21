package cweskills

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetParents 获取指定CWE条目的父级关系。
//
// 该方法调用 MITRE CWE REST API 的 GET /cwe/{id}/parents 端点，
// 返回在指定视图中此CWE条目的所有父级关系。
//
// 参数：
//   - ctx: 请求上下文
//   - id: CWE ID数字
//   - viewID: 可选的视图ID，用于限定关系范围
//
// 返回值：
//   - []Relationship: 父级关系列表
//   - error: 请求失败时返回错误
func (c *APIClient) GetParents(ctx context.Context, id int, viewID ...int) ([]Relationship, error) {
	if id <= 0 {
		return nil, NewInvalidCWEIDError(fmt.Sprintf("%d", id))
	}

	path := fmt.Sprintf("/cwe/%d/parents", id)
	if len(viewID) > 0 && viewID[0] > 0 {
		path = fmt.Sprintf("%s?view=%d", path, viewID[0])
	}

	return c.getRelations(ctx, path, id, "parents")
}

// GetChildren 获取指定CWE条目的子级关系。
//
// 该方法调用 MITRE CWE REST API 的 GET /cwe/{id}/children 端点，
// 返回在指定视图中此CWE条目的所有子级关系。
//
// 参数：
//   - ctx: 请求上下文
//   - id: CWE ID数字
//   - viewID: 可选的视图ID，用于限定关系范围
//
// 返回值：
//   - []Relationship: 子级关系列表
//   - error: 请求失败时返回错误
func (c *APIClient) GetChildren(ctx context.Context, id int, viewID ...int) ([]Relationship, error) {
	if id <= 0 {
		return nil, NewInvalidCWEIDError(fmt.Sprintf("%d", id))
	}

	path := fmt.Sprintf("/cwe/%d/children", id)
	if len(viewID) > 0 && viewID[0] > 0 {
		path = fmt.Sprintf("%s?view=%d", path, viewID[0])
	}

	return c.getRelations(ctx, path, id, "children")
}

// GetAncestors 获取指定CWE条目的所有祖先关系。
//
// 该方法调用 MITRE CWE REST API 的 GET /cwe/{id}/ancestors 端点，
// 返回此CWE条目的所有祖先（递归的父级）关系。
//
// 参数：
//   - ctx: 请求上下文
//   - id: CWE ID数字
//
// 返回值：
//   - []Relationship: 祖先关系列表
//   - error: 请求失败时返回错误
func (c *APIClient) GetAncestors(ctx context.Context, id int) ([]Relationship, error) {
	if id <= 0 {
		return nil, NewInvalidCWEIDError(fmt.Sprintf("%d", id))
	}

	path := fmt.Sprintf("/cwe/%d/ancestors", id)
	return c.getRelations(ctx, path, id, "ancestors")
}

// GetDescendants 获取指定CWE条目的所有后代关系。
//
// 该方法调用 MITRE CWE REST API 的 GET /cwe/{id}/descendants 端点，
// 返回此CWE条目的所有后代（递归的子级）关系。
//
// 参数：
//   - ctx: 请求上下文
//   - id: CWE ID数字
//
// 返回值：
//   - []Relationship: 后代关系列表
//   - error: 请求失败时返回错误
func (c *APIClient) GetDescendants(ctx context.Context, id int) ([]Relationship, error) {
	if id <= 0 {
		return nil, NewInvalidCWEIDError(fmt.Sprintf("%d", id))
	}

	path := fmt.Sprintf("/cwe/%d/descendants", id)
	return c.getRelations(ctx, path, id, "descendants")
}

// getRelations 通用关系查询方法。
func (c *APIClient) getRelations(ctx context.Context, path string, id int, relType string) ([]Relationship, error) {
	var response struct {
		Data json.RawMessage `json:"Data"`
	}

	if err := c.httpClient.Get(ctx, path, &response); err != nil {
		return nil, fmt.Errorf("获取CWE %d 的%s关系失败: %w", id, relType, err)
	}

	var relations []Relationship
	if err := json.Unmarshal(response.Data, &relations); err != nil {
		// 尝试作为对象数组解析（不同API版本格式可能不同）
		var rawRelations []struct {
			Nature string `json:"nature"`
			CWEID  int    `json:"cweId"`
			ViewID int    `json:"viewId,omitempty"`
		}
		if err2 := json.Unmarshal(response.Data, &rawRelations); err2 != nil {
			return nil, NewParseError(fmt.Sprintf("解析关系数据失败: %v", err), 0)
		}
		for _, r := range rawRelations {
			nature, parseErr := ParseRelationshipNature(r.Nature)
			if parseErr != nil {
				nature = RelationshipNature(r.Nature)
			}
			relations = append(relations, Relationship{
				Nature: nature,
				CWEID:  r.CWEID,
				ViewID: r.ViewID,
			})
		}
	}

	return relations, nil
}
