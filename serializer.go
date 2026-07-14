package cweskills

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
)

// ========== JSON序列化 ==========

// MarshalJSON 将CWE条目序列化为JSON格式。
//
// 参数：
//   - cwe: 要序列化的CWE条目
//
// 返回值：
//   - []byte: JSON数据
//   - error: 序列化失败时返回错误
func MarshalJSON(cwe *CWE) ([]byte, error) {
	if cwe == nil {
		return nil, NewValidationError("CWE", "nil")
	}
	return json.MarshalIndent(cwe, "", "  ")
}

// UnmarshalJSON 从JSON数据反序列化CWE条目。
//
// 参数：
//   - data: JSON数据
//
// 返回值：
//   - *CWE: 反序列化的CWE条目
//   - error: 反序列化失败时返回错误
func UnmarshalJSON(data []byte) (*CWE, error) {
	if len(data) == 0 {
		return nil, NewValidationError("data", "empty")
	}

	var cwe CWE
	if err := json.Unmarshal(data, &cwe); err != nil {
		return nil, NewParseError(fmt.Sprintf("JSON反序列化失败: %v", err), 0)
	}
	return &cwe, nil
}

// MarshalJSONList 将CWE条目列表序列化为JSON格式。
func MarshalJSONList(cwes []*CWE) ([]byte, error) {
	if cwes == nil {
		return []byte("[]"), nil
	}
	return json.MarshalIndent(cwes, "", "  ")
}

// UnmarshalJSONList 从JSON数据反序列化CWE条目列表。
func UnmarshalJSONList(data []byte) ([]*CWE, error) {
	if len(data) == 0 {
		return nil, NewValidationError("data", "empty")
	}

	var cwes []*CWE
	if err := json.Unmarshal(data, &cwes); err != nil {
		return nil, NewParseError(fmt.Sprintf("JSON反序列化列表失败: %v", err), 0)
	}
	return cwes, nil
}

// ========== XML序列化 ==========

// xmlMarshalIndenter 抽象 xml.MarshalIndent，便于测试注入错误。
// 默认指向 xml.MarshalIndent，对合法 safeCWE 不会返回错误；
// 测试可替换为返回错误的实现以覆盖错误分支。
var xmlMarshalIndenter = xml.MarshalIndent

// MarshalXML 将CWE条目序列化为XML格式。
func MarshalXML(cwe *CWE) ([]byte, error) {
	if cwe == nil {
		return nil, NewValidationError("CWE", "nil")
	}

	// 使用SafeCWE避免循环引用
	safe := toSafeCWE(cwe)

	output, err := xmlMarshalIndenter(safe, "", "  ")
	if err != nil {
		return nil, NewParseError(fmt.Sprintf("XML序列化失败: %v", err), 0)
	}

	return append([]byte(xml.Header), output...), nil
}

// UnmarshalXML 从XML数据反序列化CWE条目。
func UnmarshalXML(data []byte) (*CWE, error) {
	if len(data) == 0 {
		return nil, NewValidationError("data", "empty")
	}

	var safe safeCWE
	if err := xml.Unmarshal(data, &safe); err != nil {
		return nil, NewParseError(fmt.Sprintf("XML反序列化失败: %v", err), 0)
	}

	return fromSafeCWE(&safe), nil
}

// safeCWE 用于XML序列化时避免循环引用
type safeCWE struct {
	XMLName             xml.Name           `xml:"CWE"`
	ID                  int                `xml:"ID,attr"`
	Name                string             `xml:"Name"`
	Abstraction         Abstraction        `xml:"Abstraction,omitempty"`
	Structure           Structure          `xml:"Structure,omitempty"`
	Status              Status             `xml:"Status,omitempty"`
	Description         string             `xml:"Description"`
	ExtendedDescription string             `xml:"Extended_Description,omitempty"`
	LikelihoodOfExploit LikelihoodOfExploit `xml:"LikelihoodOfExploit,omitempty"`
	Relationships       []Relationship     `xml:"Relationships>Relationship,omitempty"`
	URL                 string             `xml:"URL,omitempty"`
}

// toSafeCWE 将CWE转换为SafeCWE（去除循环引用字段）
func toSafeCWE(cwe *CWE) *safeCWE {
	return &safeCWE{
		ID:                  cwe.ID,
		Name:                cwe.Name,
		Abstraction:         cwe.Abstraction,
		Structure:           cwe.Structure,
		Status:              cwe.Status,
		Description:         cwe.Description,
		ExtendedDescription: cwe.ExtendedDescription,
		LikelihoodOfExploit: cwe.LikelihoodOfExploit,
		Relationships:       cwe.Relationships,
		URL:                 cwe.URL,
	}
}

// fromSafeCWE 将SafeCWE转换回CWE
func fromSafeCWE(safe *safeCWE) *CWE {
	return &CWE{
		ID:                  safe.ID,
		Name:                safe.Name,
		Abstraction:         safe.Abstraction,
		Structure:           safe.Structure,
		Status:              safe.Status,
		Description:         safe.Description,
		ExtendedDescription: safe.ExtendedDescription,
		LikelihoodOfExploit: safe.LikelihoodOfExploit,
		Relationships:       safe.Relationships,
		URL:                 safe.URL,
		CWEType:             "weakness",
	}
}

// ========== CSV序列化 ==========

// csvHeader CSV文件的表头
var csvHeader = []string{"ID", "Name", "Abstraction", "Structure", "Status", "Description", "LikelihoodOfExploit"}

// csvSink 是 csv.Writer 的底层写入目标。默认返回新 *bytes.Buffer（与原行为一致）；
// 测试可替换为「写 N 次后返回 error」的 fake io.Writer，使 csv.Writer 在 Flush
// 时底层写失败、writer.Error() 非 nil，从而覆盖 MarshalCSV 的 Flush 错误分支。
var csvSink = func() io.Writer { return new(bytes.Buffer) }

// MarshalCSV 将CWE条目列表序列化为CSV格式。
func MarshalCSV(cwes []*CWE) ([]byte, error) {
	if cwes == nil {
		return []byte{}, nil
	}

	sink := csvSink()
	buf, _ := sink.(*bytes.Buffer)
	if buf == nil {
		// 测试注入了非 *bytes.Buffer 的 sink（如错误 writer）：
		// 无法收集字节，但仍走完整写入流程以触发错误分支。
		buf = new(bytes.Buffer)
	}
	writer := csv.NewWriter(sink)

	// 写入表头
	if err := writer.Write(csvHeader); err != nil {
		return nil, NewParseError(fmt.Sprintf("CSV写入表头失败: %v", err), 0)
	}

	// 写入数据行
	for _, cwe := range cwes {
		record := []string{
			fmt.Sprintf("%d", cwe.ID),
			cwe.Name,
			string(cwe.Abstraction),
			string(cwe.Structure),
			string(cwe.Status),
			cwe.Description,
			string(cwe.LikelihoodOfExploit),
		}
		if err := writer.Write(record); err != nil {
			return nil, NewParseError(fmt.Sprintf("CSV写入数据失败: %v", err), 0)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, NewParseError(fmt.Sprintf("CSV刷新失败: %v", err), 0)
	}

	return buf.Bytes(), nil
}

// UnmarshalCSV 从CSV数据反序列化CWE条目列表。
func UnmarshalCSV(data []byte) ([]*CWE, error) {
	if len(data) == 0 {
		return nil, NewValidationError("data", "empty")
	}

	reader := csv.NewReader(bytes.NewReader(data))
	reader.FieldsPerRecord = -1 // 允许列数不一致

	// 读取表头
	header, err := reader.Read()
	if err != nil {
		return nil, NewParseError(fmt.Sprintf("CSV读取表头失败: %v", err), 0)
	}

	// 验证表头
	if len(header) < 2 {
		return nil, NewParseError("CSV表头列数不足", 0)
	}

	var cwes []*CWE
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, NewParseError(fmt.Sprintf("CSV读取数据失败: %v", err), 0)
		}

		if len(record) < 2 {
			continue
		}

		id, err := parseCSVInt(record[0])
		if err != nil {
			continue // 跳过无效行
		}

		cwe := &CWE{
			ID:          id,
			Name:        record[1],
			CWEType:     "weakness",
		}

		if len(record) > 2 {
			cwe.Abstraction = Abstraction(record[2])
		}
		if len(record) > 3 {
			cwe.Structure = Structure(record[3])
		}
		if len(record) > 4 {
			cwe.Status = Status(record[4])
		}
		if len(record) > 5 {
			cwe.Description = record[5]
		}
		if len(record) > 6 {
			cwe.LikelihoodOfExploit = LikelihoodOfExploit(record[6])
		}

		cwes = append(cwes, cwe)
	}

	return cwes, nil
}

// parseCSVInt 解析CSV中的整数值
func parseCSVInt(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("空字符串")
	}
	var result int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("无效字符: %c", c)
		}
		result = result*10 + int(c-'0')
	}
	return result, nil
}

// ========== Registry级别的序列化 ==========

// ExportCSV 将注册表中的CWE弱点导出为CSV格式。
func (r *Registry) ExportCSV() ([]byte, error) {
	return MarshalCSV(r.GetAll())
}
