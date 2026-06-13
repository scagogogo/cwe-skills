package cwe

import (
	"strings"
	"testing"
)

const validXML = `<?xml version="1.0" encoding="UTF-8"?>
<Weakness_Catalog Name="CWE" Version="4.10" Date="2024-02-29">
  <Weaknesses>
    <Weakness ID="79" Name="Cross-site Scripting" Abstraction="Base" Structure="Simple" Status="Stable">
      <Description>XSS vulnerability</Description>
      <Extended_Description>Extended XSS info</Extended_Description>
      <Likelihood_Of_Exploit>High</Likelihood_Of_Exploit>
      <Relationships>
        <Relationship Nature="ChildOf" CWE_ID="74" View_ID="1000"/>
      </Relationships>
    </Weakness>
    <Weakness ID="89" Name="SQL Injection" Abstraction="Base" Structure="Simple" Status="Stable">
      <Description>SQL injection vulnerability</Description>
      <Relationships>
        <Relationship Nature="ChildOf" CWE_ID="74" View_ID="1000"/>
      </Relationships>
    </Weakness>
  </Weaknesses>
  <Categories>
    <Category ID="1" Name="Category1" Status="Stable">
      <Summary>A test category</Summary>
      <Relationships>
        <Relationship Nature="Has_Member" CWE_ID="79" View_ID="1000"/>
      </Relationships>
    </Category>
  </Categories>
  <Views>
    <View ID="1000" Name="Research Concepts" Type="Graph" Status="Stable">
      <Objective>Research view</Objective>
      <Members>
        <Has_Member CWE_ID="79" View_ID="1000" Direct="true"/>
        <Has_Member CWE_ID="89" View_ID="1000" Direct="true"/>
      </Members>
      <Audience>
        <Type>Researcher</Type>
      </Audience>
    </View>
  </Views>
  <Compound_Elements>
    <Compound_Element ID="680" Name="Integer Overflow to Buffer Overflow" Structure="Chain" Status="Draft">
      <Description>Chain of integer overflow to buffer overflow</Description>
      <Relationships>
        <Relationship Nature="CanFollow" CWE_ID="190" View_ID="1000"/>
        <Relationship Nature="CanPrecede" CWE_ID="119" View_ID="1000"/>
      </Relationships>
    </Compound_Element>
  </Compound_Elements>
</Weakness_Catalog>`

func TestNewXMLParser(t *testing.T) {
	parser := NewXMLParser()
	if parser == nil {
		t.Fatal("expected non-nil parser")
	}
}

func TestXMLParser_Parse(t *testing.T) {
	t.Run("valid XML with all types", func(t *testing.T) {
		parser := NewXMLParser()
		registry, err := parser.Parse(strings.NewReader(validXML))
		if err != nil {
			t.Fatalf("Parse() error: %v", err)
		}

		// Check weaknesses
		if registry.Size() != 2 {
			t.Errorf("expected 2 weaknesses, got %d", registry.Size())
		}
		cwe79, ok := registry.Get(79)
		if !ok {
			t.Fatal("expected CWE-79")
		}
		if cwe79.Name != "Cross-site Scripting" {
			t.Errorf("expected name 'Cross-site Scripting', got %q", cwe79.Name)
		}
		if cwe79.Abstraction != AbstractionBase {
			t.Errorf("expected Abstraction Base, got %v", cwe79.Abstraction)
		}
		if cwe79.Structure != StructureSimple {
			t.Errorf("expected Structure Simple, got %v", cwe79.Structure)
		}
		if cwe79.Status != StatusStable {
			t.Errorf("expected Status Stable, got %v", cwe79.Status)
		}
		if cwe79.Description != "XSS vulnerability" {
			t.Errorf("expected Description 'XSS vulnerability', got %q", cwe79.Description)
		}
		if cwe79.ExtendedDescription != "Extended XSS info" {
			t.Errorf("expected ExtendedDescription, got %q", cwe79.ExtendedDescription)
		}
		if cwe79.LikelihoodOfExploit != LikelihoodHigh {
			t.Errorf("expected Likelihood High, got %v", cwe79.LikelihoodOfExploit)
		}
		if cwe79.CWEType != "weakness" {
			t.Errorf("expected CWEType 'weakness', got %q", cwe79.CWEType)
		}
		if len(cwe79.Relationships) != 1 {
			t.Errorf("expected 1 relationship, got %d", len(cwe79.Relationships))
		}
		if cwe79.Relationships[0].Nature != RelationshipChildOf {
			t.Errorf("expected ChildOf relationship, got %v", cwe79.Relationships[0].Nature)
		}

		// Check categories
		if registry.CategoryCount() != 1 {
			t.Errorf("expected 1 category, got %d", registry.CategoryCount())
		}
		cat1, ok := registry.GetCategory(1)
		if !ok {
			t.Fatal("expected category 1")
		}
		if cat1.Name != "Category1" {
			t.Errorf("expected name 'Category1', got %q", cat1.Name)
		}
		if cat1.Description != "A test category" {
			t.Errorf("expected description 'A test category', got %q", cat1.Description)
		}

		// Check views
		if registry.ViewCount() != 1 {
			t.Errorf("expected 1 view, got %d", registry.ViewCount())
		}
		view1000, ok := registry.GetView(1000)
		if !ok {
			t.Fatal("expected view 1000")
		}
		if view1000.Name != "Research Concepts" {
			t.Errorf("expected name 'Research Concepts', got %q", view1000.Name)
		}
		if view1000.Type != ViewTypeGraph {
			t.Errorf("expected Type Graph, got %v", view1000.Type)
		}
		if len(view1000.Members) != 2 {
			t.Errorf("expected 2 members, got %d", len(view1000.Members))
		}

		// Check compound elements
		if registry.CompoundElementCount() != 1 {
			t.Errorf("expected 1 compound element, got %d", registry.CompoundElementCount())
		}
		ce680, ok := registry.GetCompoundElement(680)
		if !ok {
			t.Fatal("expected compound element 680")
		}
		if ce680.Name != "Integer Overflow to Buffer Overflow" {
			t.Errorf("expected name 'Integer Overflow to Buffer Overflow', got %q", ce680.Name)
		}
		if ce680.Structure != StructureChain {
			t.Errorf("expected Structure Chain, got %v", ce680.Structure)
		}
		if len(ce680.Relationships) != 2 {
			t.Errorf("expected 2 relationships, got %d", len(ce680.Relationships))
		}
	})

	t.Run("nil reader", func(t *testing.T) {
		parser := NewXMLParser()
		_, err := parser.Parse(nil)
		if err == nil {
			t.Fatal("expected error for nil reader")
		}
	})

	t.Run("malformed XML", func(t *testing.T) {
		parser := NewXMLParser()
		_, err := parser.Parse(strings.NewReader("<invalid xml"))
		if err == nil {
			t.Fatal("expected error for malformed XML")
		}
	})

	t.Run("empty catalog", func(t *testing.T) {
		xml := `<?xml version="1.0" encoding="UTF-8"?>
<Weakness_Catalog Name="CWE" Version="4.10" Date="2024-02-29">
</Weakness_Catalog>`
		parser := NewXMLParser()
		registry, err := parser.Parse(strings.NewReader(xml))
		if err != nil {
			t.Fatalf("Parse() error: %v", err)
		}
		if registry.Size() != 0 {
			t.Errorf("expected 0 weaknesses for empty catalog, got %d", registry.Size())
		}
	})
}

func TestXMLParser_ParseFile(t *testing.T) {
	t.Run("non-existent file", func(t *testing.T) {
		parser := NewXMLParser()
		_, err := parser.ParseFile("/nonexistent/path/file.xml")
		if err == nil {
			t.Fatal("expected error for non-existent file")
		}
	})

	t.Run("empty path", func(t *testing.T) {
		parser := NewXMLParser()
		_, err := parser.ParseFile("")
		if err == nil {
			t.Fatal("expected error for empty path")
		}
	})
}

func TestXMLParser_ParseBytes(t *testing.T) {
	t.Run("valid XML", func(t *testing.T) {
		parser := NewXMLParser()
		registry, err := parser.ParseBytes([]byte(validXML))
		if err != nil {
			t.Fatalf("ParseBytes() error: %v", err)
		}
		if registry.Size() != 2 {
			t.Errorf("expected 2 weaknesses, got %d", registry.Size())
		}
	})

	t.Run("empty data", func(t *testing.T) {
		parser := NewXMLParser()
		_, err := parser.ParseBytes([]byte{})
		if err == nil {
			t.Fatal("expected error for empty data")
		}
	})
}

func TestXMLParser_ConvertWeakness_InvalidID(t *testing.T) {
	// Test that a weakness with ID=0 is skipped
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<Weakness_Catalog Name="CWE" Version="4.10" Date="2024-02-29">
  <Weaknesses>
    <Weakness ID="0" Name="Invalid" Abstraction="Base" Structure="Simple" Status="Stable">
      <Description>Invalid ID</Description>
    </Weakness>
    <Weakness ID="79" Name="Valid" Abstraction="Base" Structure="Simple" Status="Stable">
      <Description>Valid</Description>
    </Weakness>
  </Weaknesses>
</Weakness_Catalog>`
	parser := NewXMLParser()
	registry, err := parser.Parse(strings.NewReader(xml))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if registry.Size() != 1 {
		t.Errorf("expected 1 weakness (invalid ID skipped), got %d", registry.Size())
	}
}

func TestXMLParser_ConvertCategory_InvalidID(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<Weakness_Catalog Name="CWE" Version="4.10" Date="2024-02-29">
  <Categories>
    <Category ID="0" Name="Invalid" Status="Stable">
      <Summary>Invalid</Summary>
    </Category>
  </Categories>
</Weakness_Catalog>`
	parser := NewXMLParser()
	registry, err := parser.Parse(strings.NewReader(xml))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if registry.CategoryCount() != 0 {
		t.Errorf("expected 0 categories (invalid ID skipped), got %d", registry.CategoryCount())
	}
}

func TestXMLParser_ConvertView_InvalidID(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<Weakness_Catalog Name="CWE" Version="4.10" Date="2024-02-29">
  <Views>
    <View ID="0" Name="Invalid" Type="Graph" Status="Stable">
      <Objective>Invalid</Objective>
    </View>
  </Views>
</Weakness_Catalog>`
	parser := NewXMLParser()
	registry, err := parser.Parse(strings.NewReader(xml))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if registry.ViewCount() != 0 {
		t.Errorf("expected 0 views (invalid ID skipped), got %d", registry.ViewCount())
	}
}

func TestXMLParser_ConvertCompoundElement_InvalidID(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<Weakness_Catalog Name="CWE" Version="4.10" Date="2024-02-29">
  <Compound_Elements>
    <Compound_Element ID="0" Name="Invalid" Structure="Chain" Status="Draft">
      <Description>Invalid</Description>
    </Compound_Element>
  </Compound_Elements>
</Weakness_Catalog>`
	parser := NewXMLParser()
	registry, err := parser.Parse(strings.NewReader(xml))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if registry.CompoundElementCount() != 0 {
		t.Errorf("expected 0 compound elements (invalid ID skipped), got %d", registry.CompoundElementCount())
	}
}

func TestXMLParser_ViewMemberDirectAttribute(t *testing.T) {
	tests := []struct {
		name   string
		xml    string
		direct bool
	}{
		{
			"direct true",
			`<?xml version="1.0" encoding="UTF-8"?>
<Weakness_Catalog Name="CWE" Version="4.10" Date="2024-02-29">
  <Views>
    <View ID="1000" Name="Test" Type="Graph" Status="Stable">
      <Objective>Test</Objective>
      <Members>
        <Has_Member CWE_ID="79" View_ID="1000" Direct="true"/>
      </Members>
    </View>
  </Views>
</Weakness_Catalog>`,
			true,
		},
		{
			"direct 1",
			`<?xml version="1.0" encoding="UTF-8"?>
<Weakness_Catalog Name="CWE" Version="4.10" Date="2024-02-29">
  <Views>
    <View ID="1000" Name="Test" Type="Graph" Status="Stable">
      <Objective>Test</Objective>
      <Members>
        <Has_Member CWE_ID="79" View_ID="1000" Direct="1"/>
      </Members>
    </View>
  </Views>
</Weakness_Catalog>`,
			true,
		},
		{
			"direct false",
			`<?xml version="1.0" encoding="UTF-8"?>
<Weakness_Catalog Name="CWE" Version="4.10" Date="2024-02-29">
  <Views>
    <View ID="1000" Name="Test" Type="Graph" Status="Stable">
      <Objective>Test</Objective>
      <Members>
        <Has_Member CWE_ID="79" View_ID="1000" Direct="false"/>
      </Members>
    </View>
  </Views>
</Weakness_Catalog>`,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewXMLParser()
			registry, err := parser.Parse(strings.NewReader(tt.xml))
			if err != nil {
				t.Fatalf("Parse() error: %v", err)
			}
			view, ok := registry.GetView(1000)
			if !ok {
				t.Fatal("expected view 1000")
			}
			if len(view.Members) != 1 {
				t.Fatalf("expected 1 member, got %d", len(view.Members))
			}
			if view.Members[0].Direct != tt.direct {
				t.Errorf("expected Direct=%v, got %v", tt.direct, view.Members[0].Direct)
			}
		})
	}
}

func TestXMLParser_UnknownEnums(t *testing.T) {
	// Test that unknown enum values are handled gracefully
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<Weakness_Catalog Name="CWE" Version="4.10" Date="2024-02-29">
  <Weaknesses>
    <Weakness ID="79" Name="Test" Abstraction="UnknownAbstraction" Structure="UnknownStructure" Status="UnknownStatus">
      <Description>Test</Description>
      <Likelihood_Of_Exploit>UnknownLikelihood</Likelihood_Of_Exploit>
      <Relationships>
        <Relationship Nature="UnknownNature" CWE_ID="74"/>
      </Relationships>
    </Weakness>
  </Weaknesses>
</Weakness_Catalog>`
	parser := NewXMLParser()
	registry, err := parser.Parse(strings.NewReader(xml))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	cwe, ok := registry.Get(79)
	if !ok {
		t.Fatal("expected CWE-79")
	}
	// Unknown enum values should result in empty/zero values (ParseXxx returns "" on error)
	if cwe.Abstraction != "" {
		t.Errorf("expected empty Abstraction for unknown value, got %v", cwe.Abstraction)
	}
	if cwe.Structure != "" {
		t.Errorf("expected empty Structure for unknown value, got %v", cwe.Structure)
	}
	if cwe.Status != "" {
		t.Errorf("expected empty Status for unknown value, got %v", cwe.Status)
	}
	if cwe.LikelihoodOfExploit != "" {
		t.Errorf("expected empty LikelihoodOfExploit for unknown value, got %v", cwe.LikelihoodOfExploit)
	}
}
