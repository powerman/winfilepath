// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE-go file.

// Stub testenv package for winfilepath tests
package testenv

import (
	"os"
	"os/exec"
	"runtime"
	"testing"
)

// MustHaveSymlink checks if symlinks are supported
func MustHaveSymlink(t testing.TB) {
	if runtime.GOOS == "plan9" {
		t.Helper()
		t.Skip("skipping test: symlinks not supported on plan9")
	}
	// For other platforms, assume symlinks work
	// In a real implementation, this would do proper testing
}

// SkipFlaky skips flaky tests
func SkipFlaky(t *testing.T, issue int) {
	// Skip flaky tests
	t.Skip("skipping flaky test")
}

// CleanCmdEnv cleans command environment
func CleanCmdEnv(env []string) []string {
	return env
}

// MustHaveExecPath checks if the executable path is available
func MustHaveExecPath(t testing.TB, exe string) {
	t.Helper()
	// Simplified stub - assume all executables are available
}

// Command is a stub for exec.Command in test environment
func Command(t testing.TB, name string, arg ...string) *exec.Cmd {
	t.Helper()
	return exec.Command(name, arg...)
}

// GOROOT returns GOROOT path
func GOROOT(t testing.TB) string {
	t.Helper()
	goroot := os.Getenv("GOROOT")
	if goroot == "" {
		t.Skip("GOROOT not set")
	}
	return goroot
}
