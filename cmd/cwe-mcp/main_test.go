package main

import (
	"bytes"
	"testing"
)

// TestRunMain_Version 覆盖 --version 分支（showVer 为真 → return 0）
func TestRunMain_Version(t *testing.T) {
	code := runMain([]string{"-version"})
	if code != 0 {
		t.Errorf("runMain -version: want exit 0, got %d", code)
	}
}

// TestRunMain_UnknownTransport 覆盖 default 分支（未知 transport → return 2 + stderr）
func TestRunMain_UnknownTransport(t *testing.T) {
	var buf bytes.Buffer
	orig := osStderr
	osStderr = &buf
	t.Cleanup(func() { osStderr = orig })

	code := runMain([]string{"-transport", "weird"})
	if code != 2 {
		t.Errorf("runMain unknown transport: want exit 2, got %d", code)
	}
	if !bytes.Contains(buf.Bytes(), []byte("未知传输方式")) {
		t.Errorf("expected stderr to mention unknown transport, got: %s", buf.String())
	}
}

// TestRunMain_FlagParseError 覆盖 fs.Parse 失败分支（return 2）
func TestRunMain_FlagParseError(t *testing.T) {
	// -version=notabool 让 flag 包报 "invalid boolean value"
	code := runMain([]string{"-version=notabool"})
	if code != 2 {
		t.Errorf("runMain bad bool flag: want exit 2, got %d", code)
	}
}
