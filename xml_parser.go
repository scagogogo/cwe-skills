package cweskills

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

// XMLParser 解析MITRE CWE XML目录格式。
//
// 该解析器支持解析MITRE官方提供的CWE XML下载文件，
// 将XML数据转换为Registry中的CWE条目。
// 这是实现离线模式的关键组件。
//
// 支持的XML格式为MITRE CWE Schema 7.x版本。
//
// 示例：
//
//	parser := cwe.NewXMLParser()
//	registry, err := parser.ParseFile("cwec_v4.10.xml")
//	if err != nil {
//	    log.Fatal(err)
//	}
type XMLParser struct{}

// NewXMLParser 创建一个新的XML解析器。
func NewXMLParser() *XMLParser {
	return &XMLParser{}
}

// Parse 从io.Reader解析CWE XML目录。
//
// 参数：
//   - reader: 包含XML数据的io.Reader
//
// 返回值：
//   - *Registry: 解析后的CWE注册表
//   - error: 解析失败时返回错误
func (p *XMLParser) Parse(reader io.Reader) (*Registry, error) {
	if reader == nil {
		return nil, NewValidationError("reader", "nil")
	}

	var catalog xmlWeaknessCatalog
	decoder := xml.NewDecoder(reader)
	if err := decoder.Decode(&catalog); err != nil {
		return nil, NewParseError(fmt.Sprintf("XML解析失败: %v", err), 0)
	}

	registry := NewRegistry()

	// 解析弱点
	for i := range catalog.Weaknesses {
		w := &catalog.Weaknesses[i]
		cwe := p.convertWeakness(w)
		if cwe != nil {
			_ = registry.Register(cwe) // 忽略重复注册错误
		}
	}

	// 解析类别
	for i := range catalog.Categories {
		cat := &catalog.Categories[i]
		category := p.convertCategory(cat)
		if category != nil {
			_ = registry.RegisterCategory(category)
		}
	}

	// 解析视图
	for i := range catalog.Views {
		v := &catalog.Views[i]
		view := p.convertView(v)
		if view != nil {
			_ = registry.RegisterView(view)
		}
	}

	// 解析复合元素
	for i := range catalog.CompoundElements {
		ce := &catalog.CompoundElements[i]
		compound := p.convertCompoundElement(ce)
		if compound != nil {
			_ = registry.RegisterCompoundElement(compound)
		}
	}

	return registry, nil
}

// ParseFile 从文件解析CWE XML目录。
//
// 参数：
//   - path: XML文件路径
//
// 返回值：
//   - *Registry: 解析后的CWE注册表
//   - error: 解析失败时返回错误
func (p *XMLParser) ParseFile(path string) (*Registry, error) {
	if path == "" {
		return nil, NewValidationError("path", "empty")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, NewParseError(fmt.Sprintf("打开文件失败: %v", err), 0)
	}
	defer file.Close()

	return p.Parse(file)
}

// ParseBytes 从字节数组解析CWE XML目录。
//
// 参数：
//   - data: XML数据的字节数组
//
// 返回值：
//   - *Registry: 解析后的CWE注册表
//   - error: 解析失败时返回错误
func (p *XMLParser) ParseBytes(data []byte) (*Registry, error) {
	if len(data) == 0 {
		return nil, NewValidationError("data", "empty")
	}

	return p.Parse(newByteReader(data))
}

// newByteReader 创建一个从字节切片读取的io.Reader。
func newByteReader(data []byte) *byteReader {
	return &byteReader{data: data}
}

type byteReader struct {
	data   []byte
	offset int
}

func (r *byteReader) Read(p []byte) (n int, err error) {
	if r.offset >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.offset:])
	r.offset += n
	return n, nil
}

// ========== XML数据结构定义 ==========

// xmlWeaknessCatalog 对应XML的 <Weakness_Catalog> 根元素
type xmlWeaknessCatalog struct {
	XMLName          xml.Name              `xml:"Weakness_Catalog"`
	Name             string                `xml:"Name,attr"`
	Version          string                `xml:"Version,attr"`
	Date             string                `xml:"Date,attr"`
	Weaknesses       []xmlWeakness         `xml:"Weaknesses>Weakness"`
	Categories       []xmlCategory         `xml:"Categories>Category"`
	Views            []xmlView             `xml:"Views>View"`
	CompoundElements []xmlCompoundElement  `xml:"Compound_Elements>Compound_Element"`
}

// xmlWeakness 对应XML的 <Weakness> 元素
type xmlWeakness struct {
	ID                  int    `xml:"ID,attr"`
	Name                string `xml:"Name,attr"`
	Abstraction         string `xml:"Abstraction,attr"`
	Structure           string `xml:"Structure,attr"`
	Status              string `xml:"Status,attr"`
	Description         string `xml:"Description"`
	ExtendedDescription string `xml:"Extended_Description"`
	LikelihoodOfExploit string `xml:"Likelihood_Of_Exploit"`
	Relationships       []xmlRelationship `xml:"Relationships>Relationship"`
}

// xmlCategory 对应XML的 <Category> 元素
type xmlCategory struct {
	ID          int    `xml:"ID,attr"`
	Name        string `xml:"Name,attr"`
	Status      string `xml:"Status,attr"`
	Summary     string `xml:"Summary"`
	Relationships []xmlRelationship `xml:"Relationships>Relationship"`
}

// xmlView 对应XML的 <View> 元素
type xmlView struct {
	ID       int    `xml:"ID,attr"`
	Name     string `xml:"Name,attr"`
	Type     string `xml:"Type,attr"`
	Status   string `xml:"Status,attr"`
	Objective string `xml:"Objective"`
	Members  []xmlViewMember `xml:"Members>Has_Member"`
	Audience []xmlAudience   `xml:"Audience"`
}

// xmlViewMember 对应XML视图中的成员
type xmlViewMember struct {
	CWEID     int    `xml:"CWE_ID,attr"`
	ViewID    int    `xml:"View_ID,attr"`
	Ordinal   string `xml:"Ordinal,attr,omitempty"`
	Direct    string `xml:"Direct,attr,omitempty"`
}

// xmlAudience 对应XML中的 <Audience> 元素
type xmlAudience struct {
	Type        string `xml:"Type"`
	Description string `xml:"Description,omitempty"`
}

// xmlCompoundElement 对应XML的 <Compound_Element> 元素
type xmlCompoundElement struct {
	ID            int    `xml:"ID,attr"`
	Name          string `xml:"Name,attr"`
	Structure     string `xml:"Structure,attr"`
	Status        string `xml:"Status,attr"`
	Description   string `xml:"Description"`
	Relationships []xmlRelationship `xml:"Relationships>Relationship"`
}

// xmlRelationship 对应XML中的关系元素
type xmlRelationship struct {
	Nature  string `xml:"Nature,attr"`
	CWEID   int    `xml:"CWE_ID,attr"`
	ViewID  int    `xml:"View_ID,attr,omitempty"`
	Ordinal string `xml:"Ordinal,attr,omitempty"`
}

// ========== 转换方法 ==========

// convertWeakness 将XML弱点转换为CWE结构体
func (p *XMLParser) convertWeakness(w *xmlWeakness) *CWE {
	if w == nil || w.ID <= 0 {
		return nil
	}

	abstraction, _ := ParseAbstraction(w.Abstraction)
	structure, _ := ParseStructure(w.Structure)
	status, _ := ParseStatus(w.Status)
	likelihood, _ := ParseLikelihoodOfExploit(w.LikelihoodOfExploit)

	cwe := &CWE{
		ID:                  w.ID,
		Name:                w.Name,
		Abstraction:         abstraction,
		Structure:           structure,
		Status:              status,
		Description:         w.Description,
		ExtendedDescription: w.ExtendedDescription,
		LikelihoodOfExploit: likelihood,
		CWEType:             "weakness",
		URL:                 fmt.Sprintf("https://cwe.mitre.org/data/definitions/%d.html", w.ID),
	}

	// 转换关系
	for _, rel := range w.Relationships {
		nature, _ := ParseRelationshipNature(rel.Nature)
		cwe.Relationships = append(cwe.Relationships, Relationship{
			Nature:  nature,
			CWEID:   rel.CWEID,
			ViewID:  rel.ViewID,
			Ordinal: rel.Ordinal,
		})
	}

	return cwe
}

// convertCategory 将XML类别转换为Category结构体
func (p *XMLParser) convertCategory(cat *xmlCategory) *Category {
	if cat == nil || cat.ID <= 0 {
		return nil
	}

	status, _ := ParseStatus(cat.Status)

	category := &Category{
		ID:          cat.ID,
		Name:        cat.Name,
		Status:      status,
		Description: cat.Summary,
	}

	// 转换关系
	for _, rel := range cat.Relationships {
		nature, _ := ParseRelationshipNature(rel.Nature)
		category.Relationships = append(category.Relationships, Relationship{
			Nature:  nature,
			CWEID:   rel.CWEID,
			ViewID:  rel.ViewID,
			Ordinal: rel.Ordinal,
		})
	}

	return category
}

// convertView 将XML视图转换为View结构体
func (p *XMLParser) convertView(v *xmlView) *View {
	if v == nil || v.ID <= 0 {
		return nil
	}

	viewType, _ := ParseViewType(v.Type)
	status, _ := ParseStatus(v.Status)

	view := &View{
		ID:          v.ID,
		Name:        v.Name,
		Type:        viewType,
		Status:      status,
		Description: v.Objective,
	}

	// 转换成员
	for _, m := range v.Members {
		direct := m.Direct == "true" || m.Direct == "1"
		view.Members = append(view.Members, ViewMember{
			CWEID:  m.CWEID,
			ViewID: m.ViewID,
			Direct: direct,
		})
	}

	return view
}

// convertCompoundElement 将XML复合元素转换为CompoundElement结构体
func (p *XMLParser) convertCompoundElement(ce *xmlCompoundElement) *CompoundElement {
	if ce == nil || ce.ID <= 0 {
		return nil
	}

	structure, _ := ParseStructure(ce.Structure)
	status, _ := ParseStatus(ce.Status)

	compound := &CompoundElement{
		ID:          ce.ID,
		Name:        ce.Name,
		Structure:   structure,
		Status:      status,
		Description: ce.Description,
	}

	// 转换关系
	for _, rel := range ce.Relationships {
		nature, _ := ParseRelationshipNature(rel.Nature)
		compound.Relationships = append(compound.Relationships, Relationship{
			Nature:  nature,
			CWEID:   rel.CWEID,
			ViewID:  rel.ViewID,
			Ordinal: rel.Ordinal,
		})
	}

	return compound
}
