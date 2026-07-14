package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/scagogogo/cwe-skills"
	"github.com/spf13/cobra"
)

// runCmd 执行一个 cobra 命令的 RunE，捕获输出，返回输出和错误
func runCmd(cmd *cobra.Command, args []string) (string, error) {
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	// 重置 flags 以避免测试间污染（cobra 命令是全局变量）
	cmd.SetArgs(args)
	err := cmd.RunE(cmd, args)
	return buf.String(), err
}

func withOutputFormat(format string, fn func()) {
	old := outputFormat
	defer func() { outputFormat = old }()
	outputFormat = format
	fn()
}

// ==================== versionCmd ====================

func TestVersionCmd_Text(t *testing.T) {
	oldV, oldC, oldD := cliVersion, cliGitCommit, cliBuildDate
	defer func() { cliVersion, cliGitCommit, cliBuildDate = oldV, oldC, oldD }()
	cliVersion = "1.2.3"
	cliGitCommit = "abc123"
	cliBuildDate = "2026-01-01"

	withOutputFormat("text", func() {
		out, err := runCmd(versionCmd, nil)
		if err != nil {
			t.Fatalf("versionCmd error: %v", err)
		}
		if !strings.Contains(out, "CWE CLI:") {
			t.Errorf("expected CWE CLI, got %q", out)
		}
		if !strings.Contains(out, "1.2.3") {
			t.Errorf("expected version 1.2.3, got %q", out)
		}
		if !strings.Contains(out, "Git Commit:  abc123") {
			t.Errorf("expected git commit, got %q", out)
		}
		if !strings.Contains(out, "Build Date:  2026-01-01") {
			t.Errorf("expected build date, got %q", out)
		}
	})
}

func TestVersionCmd_UnknownDefaults(t *testing.T) {
	// gitCommit/BuildDate 为 unknown 时不打印
	oldV, oldC, oldD := cliVersion, cliGitCommit, cliBuildDate
	defer func() { cliVersion, cliGitCommit, cliBuildDate = oldV, oldC, oldD }()
	cliVersion = "dev"
	cliGitCommit = "unknown"
	cliBuildDate = "unknown"

	withOutputFormat("text", func() {
		out, err := runCmd(versionCmd, nil)
		if err != nil {
			t.Fatalf("versionCmd error: %v", err)
		}
		if strings.Contains(out, "Git Commit") {
			t.Errorf("should not print Git Commit when unknown, got %q", out)
		}
		if strings.Contains(out, "Build Date") {
			t.Errorf("should not print Build Date when unknown, got %q", out)
		}
	})
}

func TestVersionCmd_JSON(t *testing.T) {
	oldV, oldC, oldD := cliVersion, cliGitCommit, cliBuildDate
	defer func() { cliVersion, cliGitCommit, cliBuildDate = oldV, oldC, oldD }()
	cliVersion = "1.2.3"
	cliGitCommit = "abc"
	cliBuildDate = "2026-01-01"

	withOutputFormat("json", func() {
		out, err := runCmd(versionCmd, nil)
		if err != nil {
			t.Fatalf("versionCmd error: %v", err)
		}
		var got map[string]interface{}
		if err := json.Unmarshal([]byte(out), &got); err != nil {
			t.Fatalf("expected JSON: %v\noutput: %s", err, out)
		}
		if got["cli"] != "1.2.3" {
			t.Errorf("expected cli 1.2.3, got %v", got["cli"])
		}
	})
}

// ==================== apiVersionCmd ====================

func TestApiVersionCmd_Text(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"Data": map[string]interface{}{
				"version":     "1.0",
				"releaseDate": "2026-01-01",
				"name":        "CWE API",
			},
		})
	}))
	defer srv.Close()

	apiVersionBaseURL = srv.URL

	withOutputFormat("text", func() {
		out, err := runCmd(apiVersionCmd, nil)
		if err != nil {
			t.Fatalf("apiVersionCmd error: %v", err)
		}
		if !strings.Contains(out, "1.0") {
			t.Errorf("expected version 1.0, got %q", out)
		}
		if !strings.Contains(out, "2026-01-01") {
			t.Errorf("expected release date, got %q", out)
		}
		if !strings.Contains(out, "CWE API") {
			t.Errorf("expected name, got %q", out)
		}
	})
}

func TestApiVersionCmd_JSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"Data": map[string]interface{}{
				"version": "1.0",
			},
		})
	}))
	defer srv.Close()

	apiVersionBaseURL = srv.URL

	withOutputFormat("json", func() {
		out, err := runCmd(apiVersionCmd, nil)
		if err != nil {
			t.Fatalf("apiVersionCmd error: %v", err)
		}
		var got map[string]interface{}
		if err := json.Unmarshal([]byte(out), &got); err != nil {
			t.Fatalf("expected JSON: %v", err)
		}
	})
}

func TestApiVersionCmd_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	apiVersionBaseURL = srv.URL

	_, err := runCmd(apiVersionCmd, nil)
	if err == nil {
		t.Fatal("expected error on API failure, got nil")
	}
}

// ==================== formatCmd ====================

func TestFormatCmd_Text(t *testing.T) {
	withOutputFormat("text", func() {
		out, err := runCmd(formatCmd, []string{"79", "cwe-89", "CWE-352"})
		if err != nil {
			t.Fatalf("formatCmd error: %v", err)
		}
		if !strings.Contains(out, "CWE-79") || !strings.Contains(out, "CWE-89") || !strings.Contains(out, "CWE-352") {
			t.Errorf("expected formatted IDs, got %q", out)
		}
	})
}

func TestFormatCmd_JSON(t *testing.T) {
	withOutputFormat("json", func() {
		out, err := runCmd(formatCmd, []string{"79"})
		if err != nil {
			t.Fatalf("formatCmd error: %v", err)
		}
		var results []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &results); err != nil {
			t.Fatalf("expected JSON array: %v", err)
		}
		if len(results) != 1 || results[0]["output"] != "CWE-79" {
			t.Errorf("unexpected: %v", results)
		}
	})
}

func TestFormatCmd_InvalidID(t *testing.T) {
	withOutputFormat("text", func() {
		out, err := runCmd(formatCmd, []string{"abc"})
		if err != nil {
			t.Fatalf("formatCmd error: %v", err)
		}
		if !strings.Contains(out, "错误") {
			t.Errorf("expected error in output, got %q", out)
		}
	})
}

func TestFormatCmd_NoArgs(t *testing.T) {
	withOutputFormat("text", func() {
		_, err := runCmd(formatCmd, nil)
		if err == nil {
			t.Fatal("expected error for no args, got nil")
		}
	})
}

func TestFormatCmd_Mixed(t *testing.T) {
	withOutputFormat("text", func() {
		out, err := runCmd(formatCmd, []string{"79", "abc"})
		if err != nil {
			t.Fatalf("formatCmd error: %v", err)
		}
		if !strings.Contains(out, "CWE-79") {
			t.Errorf("expected CWE-79, got %q", out)
		}
		if !strings.Contains(out, "错误") {
			t.Errorf("expected error for abc, got %q", out)
		}
	})
}

// ==================== extractCmd ====================

func TestExtractCmd_Text(t *testing.T) {
	withOutputFormat("text", func() {
		out, err := runCmd(extractCmd, []string{"受CWE-79和CWE-89影响"})
		if err != nil {
			t.Fatalf("extractCmd error: %v", err)
		}
		if !strings.Contains(out, "CWE-79") || !strings.Contains(out, "CWE-89") {
			t.Errorf("expected CWE-79 and CWE-89, got %q", out)
		}
	})
}

func TestExtractCmd_JSON(t *testing.T) {
	withOutputFormat("json", func() {
		out, err := runCmd(extractCmd, []string{"CWE-79"})
		if err != nil {
			t.Fatalf("extractCmd error: %v", err)
		}
		var got map[string]interface{}
		if err := json.Unmarshal([]byte(out), &got); err != nil {
			t.Fatalf("expected JSON: %v", err)
		}
		if got["count"].(float64) != 1 {
			t.Errorf("expected count 1, got %v", got["count"])
		}
	})
}

func TestExtractCmd_NoIDs(t *testing.T) {
	withOutputFormat("text", func() {
		out, err := runCmd(extractCmd, []string{"no cwe here"})
		if err != nil {
			t.Fatalf("extractCmd error: %v", err)
		}
		if !strings.Contains(out, "未找到") {
			t.Errorf("expected 未找到, got %q", out)
		}
	})
}

func TestExtractCmd_NoArgs(t *testing.T) {
	withOutputFormat("text", func() {
		_, err := runCmd(extractCmd, nil)
		if err == nil {
			t.Fatal("expected error for no args, got nil")
		}
	})
}

// ==================== compareCmd ====================

func TestCompareCmd_Text(t *testing.T) {
	withOutputFormat("text", func() {
		out, err := runCmd(compareCmd, []string{"CWE-79", "CWE-89"})
		if err != nil {
			t.Fatalf("compareCmd error: %v", err)
		}
		if !strings.Contains(out, "less than") {
			t.Errorf("expected less than, got %q", out)
		}
	})
}

func TestCompareCmd_GreaterThan(t *testing.T) {
	withOutputFormat("text", func() {
		out, err := runCmd(compareCmd, []string{"CWE-89", "CWE-79"})
		if err != nil {
			t.Fatalf("compareCmd error: %v", err)
		}
		if !strings.Contains(out, "greater than") {
			t.Errorf("expected greater than, got %q", out)
		}
	})
}

func TestCompareCmd_Equal(t *testing.T) {
	withOutputFormat("text", func() {
		out, err := runCmd(compareCmd, []string{"CWE-79", "CWE-79"})
		if err != nil {
			t.Fatalf("compareCmd error: %v", err)
		}
		if !strings.Contains(out, "equal to") {
			t.Errorf("expected equal to, got %q", out)
		}
	})
}

func TestCompareCmd_JSON(t *testing.T) {
	withOutputFormat("json", func() {
		out, err := runCmd(compareCmd, []string{"CWE-79", "CWE-89"})
		if err != nil {
			t.Fatalf("compareCmd error: %v", err)
		}
		var got map[string]interface{}
		if err := json.Unmarshal([]byte(out), &got); err != nil {
			t.Fatalf("expected JSON: %v", err)
		}
		if got["comparison"] != "less than" {
			t.Errorf("expected less than, got %v", got["comparison"])
		}
	})
}

func TestCompareCmd_InvalidID(t *testing.T) {
	withOutputFormat("text", func() {
		_, err := runCmd(compareCmd, []string{"abc", "CWE-79"})
		if err == nil {
			t.Fatal("expected error for invalid ID, got nil")
		}
	})
}

// ==================== compareIntCmd ====================

func TestCompareIntCmd_Text(t *testing.T) {
	withOutputFormat("text", func() {
		out, err := runCmd(compareIntCmd, []string{"79", "89"})
		if err != nil {
			t.Fatalf("compareIntCmd error: %v", err)
		}
		if !strings.Contains(out, "less than") {
			t.Errorf("expected less than, got %q", out)
		}
	})
}

func TestCompareIntCmd_JSON(t *testing.T) {
	withOutputFormat("json", func() {
		out, err := runCmd(compareIntCmd, []string{"89", "79"})
		if err != nil {
			t.Fatalf("compareIntCmd error: %v", err)
		}
		var got map[string]interface{}
		if err := json.Unmarshal([]byte(out), &got); err != nil {
			t.Fatalf("expected JSON: %v", err)
		}
		if got["comparison"] != "greater than" {
			t.Errorf("expected greater than, got %v", got["comparison"])
		}
	})
}

func TestCompareIntCmd_Equal(t *testing.T) {
	withOutputFormat("text", func() {
		out, err := runCmd(compareIntCmd, []string{"79", "79"})
		if err != nil {
			t.Fatalf("compareIntCmd error: %v", err)
		}
		if !strings.Contains(out, "equal to") {
			t.Errorf("expected equal to, got %q", out)
		}
	})
}

func TestCompareIntCmd_InvalidInt(t *testing.T) {
	withOutputFormat("text", func() {
		_, err := runCmd(compareIntCmd, []string{"abc", "79"})
		if err == nil {
			t.Fatal("expected error for non-int, got nil")
		}
	})
	// 第二个参数无效
	withOutputFormat("text", func() {
		_, err := runCmd(compareIntCmd, []string{"79", "xyz"})
		if err == nil {
			t.Fatal("expected error for second non-int, got nil")
		}
	})
}

// ==================== parseCmd ====================

func TestParseCmd_Text(t *testing.T) {
	withOutputFormat("text", func() {
		out, err := runCmd(parseCmd, []string{"CWE-79", "89"})
		if err != nil {
			t.Fatalf("parseCmd error: %v", err)
		}
		if !strings.Contains(out, "CWE-79") || !strings.Contains(out, "ID: 79") {
			t.Errorf("expected parse output, got %q", out)
		}
	})
}

func TestParseCmd_JSON(t *testing.T) {
	withOutputFormat("json", func() {
		out, err := runCmd(parseCmd, []string{"CWE-79"})
		if err != nil {
			t.Fatalf("parseCmd error: %v", err)
		}
		var results []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &results); err != nil {
			t.Fatalf("expected JSON array: %v", err)
		}
		if len(results) != 1 || !results[0]["valid"].(bool) {
			t.Errorf("unexpected: %v", results)
		}
	})
}

func TestParseCmd_Invalid(t *testing.T) {
	withOutputFormat("text", func() {
		out, err := runCmd(parseCmd, []string{"abc"})
		if err != nil {
			t.Fatalf("parseCmd error: %v", err)
		}
		if !strings.Contains(out, "无效") {
			t.Errorf("expected 无效, got %q", out)
		}
	})
}

func TestParseCmd_NoArgs(t *testing.T) {
	withOutputFormat("text", func() {
		_, err := runCmd(parseCmd, nil)
		if err == nil {
			t.Fatal("expected error for no args, got nil")
		}
	})
}

// ==================== validateCmd ====================

func TestValidateCmd_Text(t *testing.T) {
	withOutputFormat("text", func() {
		out, err := runCmd(validateCmd, []string{"CWE-79"})
		if err != nil {
			t.Fatalf("validateCmd error: %v", err)
		}
		if !strings.Contains(out, "有效") {
			t.Errorf("expected 有效, got %q", out)
		}
	})
}

func TestValidateCmd_JSON(t *testing.T) {
	withOutputFormat("json", func() {
		out, err := runCmd(validateCmd, []string{"CWE-79"})
		if err != nil {
			t.Fatalf("validateCmd error: %v", err)
		}
		var results []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &results); err != nil {
			t.Fatalf("expected JSON array: %v", err)
		}
		if len(results) != 1 || !results[0]["valid"].(bool) {
			t.Errorf("unexpected: %v", results)
		}
	})
}

func TestValidateCmd_Invalid(t *testing.T) {
	withOutputFormat("text", func() {
		out, err := runCmd(validateCmd, []string{"abc"})
		if err == nil {
			t.Fatal("expected error for invalid ID, got nil")
		}
		if !strings.Contains(out, "无效") {
			t.Errorf("expected 无效 in output, got %q", out)
		}
	})
}

func TestValidateCmd_Mixed(t *testing.T) {
	withOutputFormat("text", func() {
		out, err := runCmd(validateCmd, []string{"CWE-79", "abc"})
		if err == nil {
			t.Fatal("expected error for mixed valid/invalid, got nil")
		}
		if !strings.Contains(out, "有效") || !strings.Contains(out, "无效") {
			t.Errorf("expected both 有效 and 无效, got %q", out)
		}
	})
}

func TestValidateCmd_NoArgs(t *testing.T) {
	withOutputFormat("text", func() {
		_, err := runCmd(validateCmd, nil)
		if err == nil {
			t.Fatal("expected error for no args, got nil")
		}
	})
}

// ==================== Execute (root.go) ====================
// Execute() 在错误时调用 os.Exit(1)，无法安全测试其错误分支。
// 测试 rootCmd.Execute() 直接执行（覆盖 cobra 调度逻辑），不经过 Execute() 包装。

func TestRootCmd_ExecuteVersion(t *testing.T) {
	// 直接调 rootCmd.Execute()，避免 Execute() 的 os.Exit
	oldV, oldFmt := cliVersion, outputFormat
	defer func() { cliVersion, outputFormat = oldV, oldFmt }()
	cliVersion = "test-version"
	outputFormat = "text"

	// 在 versionCmd 上设置输出 buffer（子命令执行时用自己的 Out）
	var buf bytes.Buffer
	versionCmd.SetOut(&buf)

	rootCmd.SetArgs([]string{"version"})

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("rootCmd.Execute panicked: %v", r)
		}
	}()

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("rootCmd.Execute error: %v", err)
	}
	if !strings.Contains(buf.String(), "CWE CLI") {
		t.Errorf("expected CWE CLI in output, got %q", buf.String())
	}
}

// 确保未使用的 import 被引用
var _ = context.Background
var _ = cweskills.Version
var _ = &cobra.Command{}

// ==================== load*Registry 函数 ====================

// testXMLContent 用于 load*Registry 测试的最小 CWE XML
const testXMLContent = `<?xml version="1.0" encoding="UTF-8"?>
<Weakness_Catalog Name="Test" Version="7.0" Date="2026-01-01">
  <Weaknesses>
    <Weakness ID="79" Name="XSS" Abstraction="Variant" Structure="Simple" Status="Stable">
      <Description>Test desc</Description>
      <Relationships>
        <Relationship Nature="ChildOf" CWE_ID="74" View_ID="1000"/>
      </Relationships>
    </Weakness>
    <Weakness ID="74" Name="Injection" Abstraction="Base" Structure="Simple" Status="Stable">
      <Description>Parent</Description>
    </Weakness>
  </Weaknesses>
</Weakness_Catalog>`

func writeTestXML(t *testing.T) string {
	t.Helper()
	path := t.TempDir() + "/test.xml"
	if err := os.WriteFile(path, []byte(testXMLContent), 0644); err != nil {
		t.Fatalf("write xml: %v", err)
	}
	return path
}

func TestLoadRegistry_EmptyPath(t *testing.T) {
	old := registryXMLPath
	defer func() { registryXMLPath = old }()
	registryXMLPath = ""
	_, err := loadRegistry()
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestLoadRegistry_Valid(t *testing.T) {
	old := registryXMLPath
	defer func() { registryXMLPath = old }()
	registryXMLPath = writeTestXML(t)
	r, err := loadRegistry()
	if err != nil {
		t.Fatalf("loadRegistry error: %v", err)
	}
	if r.Size() != 2 {
		t.Errorf("expected 2 weaknesses, got %d", r.Size())
	}
}

func TestLoadRegistry_InvalidFile(t *testing.T) {
	old := registryXMLPath
	defer func() { registryXMLPath = old }()
	registryXMLPath = "/nonexistent/file.xml"
	_, err := loadRegistry()
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
}

func TestLoadNavRegistry_EmptyPath(t *testing.T) {
	old := navXMLPath
	defer func() { navXMLPath = old }()
	navXMLPath = ""
	_, _, err := loadNavRegistry()
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestLoadNavRegistry_Valid(t *testing.T) {
	old := navXMLPath
	defer func() { navXMLPath = old }()
	navXMLPath = writeTestXML(t)
	r, nav, err := loadNavRegistry()
	if err != nil {
		t.Fatalf("loadNavRegistry error: %v", err)
	}
	if r == nil || nav == nil {
		t.Fatal("expected non-nil registry and navigator")
	}
}

func TestLoadTreeRegistry_EmptyPath(t *testing.T) {
	old := treeXMLPath
	defer func() { treeXMLPath = old }()
	treeXMLPath = ""
	_, err := loadTreeRegistry()
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestLoadTreeRegistry_Valid(t *testing.T) {
	old := treeXMLPath
	defer func() { treeXMLPath = old }()
	treeXMLPath = writeTestXML(t)
	r, err := loadTreeRegistry()
	if err != nil {
		t.Fatalf("loadTreeRegistry error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

// ==================== enums 子命令 RunE ====================

func TestEnumSubCommands(t *testing.T) {
	// 遍历 enumCmd 的所有子命令，text 和 json 各执行一次
	for _, sub := range enumCmd.Commands() {
		t.Run(sub.Name()+"_text", func(t *testing.T) {
			withOutputFormat("text", func() {
				out, err := runCmd(sub, nil)
				if err != nil {
					t.Fatalf("enum %s error: %v", sub.Name(), err)
				}
				if len(out) == 0 {
					t.Errorf("expected non-empty output for %s", sub.Name())
				}
			})
		})
		t.Run(sub.Name()+"_json", func(t *testing.T) {
			withOutputFormat("json", func() {
				out, err := runCmd(sub, nil)
				if err != nil {
					t.Fatalf("enum %s json error: %v", sub.Name(), err)
				}
				var got interface{}
				if err := json.Unmarshal([]byte(out), &got); err != nil {
					t.Errorf("expected valid JSON for %s, got %q: %v", sub.Name(), out, err)
				}
			})
		})
	}
}

// ==================== load*Registry 无效 XML ====================

func TestLoadRegistry_InvalidXML(t *testing.T) {
	old := registryXMLPath
	defer func() { registryXMLPath = old }()
	// 写一个无效 XML 文件
	path := t.TempDir() + "/bad.xml"
	os.WriteFile(path, []byte("<not valid xml"), 0644)
	registryXMLPath = path
	_, err := loadRegistry()
	if err == nil {
		t.Fatal("expected error for invalid XML, got nil")
	}
}

func TestLoadNavRegistry_InvalidXML(t *testing.T) {
	old := navXMLPath
	defer func() { navXMLPath = old }()
	path := t.TempDir() + "/bad.xml"
	os.WriteFile(path, []byte("<not valid xml"), 0644)
	navXMLPath = path
	_, _, err := loadNavRegistry()
	if err == nil {
		t.Fatal("expected error for invalid XML, got nil")
	}
}

func TestLoadTreeRegistry_InvalidXML(t *testing.T) {
	old := treeXMLPath
	defer func() { treeXMLPath = old }()
	path := t.TempDir() + "/bad.xml"
	os.WriteFile(path, []byte("<not valid xml"), 0644)
	treeXMLPath = path
	_, err := loadTreeRegistry()
	if err == nil {
		t.Fatal("expected error for invalid XML, got nil")
	}
}
