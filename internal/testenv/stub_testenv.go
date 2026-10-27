// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE-go file.

// Package testenv is a stub of internal/testenv for winfilepath tests.
package testenv

import (
	"os"
	"os/exec"
	"runtime"
	"testing"
)

// MustHaveSymlink checks if symlinks are supported.
func MustHaveSymlink(tb testing.TB) {
	tb.Helper()
	if runtime.GOOS == "plan9" {
		tb.Skip("skipping test: symlinks not supported on plan9")
	}
	// For other platforms, assume symlinks work
	// In a real implementation, this would do proper testing
}

// SkipFlaky skips flaky tests.
func SkipFlaky(t *testing.T, _ int) {
	t.Helper()
	// Skip flaky tests
	t.Skip("skipping flaky test")
}

// CleanCmdEnv cleans command environment.
func CleanCmdEnv(env []string) []string {
	return env
}

// MustHaveExecPath checks if the executable path is available.
func MustHaveExecPath(tb testing.TB, _ string) {
	tb.Helper()
	// Simplified stub - assume all executables are available
}

// Command is a stub for [exec.Command] in test environment.
func Command(tb testing.TB, name string, arg ...string) *exec.Cmd {
	tb.Helper()
	// Not using tb.Context because callers may run it from t.Cleanup, when it is canceled.
	return exec.Command(name, arg...) //nolint:gosec,noctx // See above.
}

// GOROOT returns GOROOT path.
func GOROOT(tb testing.TB) string {
	tb.Helper()
	goroot := os.Getenv("GOROOT")
	if goroot == "" {
		tb.Skip("GOROOT not set")
	}
	return goroot
}
