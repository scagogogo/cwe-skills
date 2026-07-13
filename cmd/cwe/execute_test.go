package main

import (
	"bytes"
	"testing"
)

// TestExecuteRoot_HelpSuccess 覆盖 cobra 执行成功路径（return 0）
func TestExecuteRoot_HelpSuccess(t *testing.T) {
	rootCmd.SetArgs([]string{"--help"})
	t.Cleanup(func() { rootCmd.SetArgs([]string{}) })

	code := executeRoot()
	if code != 0 {
		t.Errorf("executeRoot --help: want 0, got %d", code)
	}
}

// TestExecuteRoot_UnknownCommandFailure 覆盖 cobra 执行失败路径（return 1 + stderr）
func TestExecuteRoot_UnknownCommandFailure(t *testing.T) {
	var buf bytes.Buffer
	orig := osStderr
	osStderr = &buf
	t.Cleanup(func() { osStderr = orig })

	rootCmd.SetArgs([]string{"nonexistent-subcommand-xyz"})
	t.Cleanup(func() { rootCmd.SetArgs([]string{}) })

	code := executeRoot()
	if code != 1 {
		t.Errorf("executeRoot unknown cmd: want 1, got %d", code)
	}
	if buf.Len() == 0 {
		t.Error("expected error output on stderr")
	}
}
