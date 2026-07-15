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

// TestExecute_InjectsOSExit 覆盖 Execute 的 os.Exit 调用点（root.go）：
// 注入 osExit 不退出，断言 Execute 调用了 osExit 且收到 executeRoot 的退出码。
func TestExecute_InjectsOSExit(t *testing.T) {
	orig := osExit
	var gotCode int
	called := false
	osExit = func(code int) { gotCode = code; called = true }
	t.Cleanup(func() { osExit = orig })

	rootCmd.SetArgs([]string{"--help"})
	t.Cleanup(func() { rootCmd.SetArgs([]string{}) })

	Execute()
	if !called {
		t.Fatal("Execute: expected osExit to be called")
	}
	if gotCode != 0 {
		t.Errorf("Execute --help: want exit code 0, got %d", gotCode)
	}
}

// TestMain_InjectsOSExit 覆盖 main 的 osExit 调用点（main.go）：
// 注入 osExit 不退出，断言 main 调用了 osExit。
func TestMain_InjectsOSExit(t *testing.T) {
	orig := osExit
	var gotCode int
	called := false
	osExit = func(code int) { gotCode = code; called = true }
	t.Cleanup(func() { osExit = orig })

	rootCmd.SetArgs([]string{"--help"})
	t.Cleanup(func() { rootCmd.SetArgs([]string{}) })

	main()
	if !called {
		t.Fatal("main: expected osExit to be called")
	}
	_ = gotCode
}
