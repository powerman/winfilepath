// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package winsyscall

import "os"

// FullPath stub for non-Windows platforms
func FullPath(path string) (string, error) {
	// Simple implementation that mimics Windows FullPath behavior
	// Note: This requires importing the parent winfilepath package,
	// but since this is a stub that panics on filesystem access anyway,
	// we'll use a basic implementation
	if len(path) > 0 && (path[0] == '\\' || len(path) > 1 && path[1] == ':') {
		// Looks like absolute Windows path
		return path, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	// Simple join with backslash for Windows-style paths
	if wd[len(wd)-1] != '\\' {
		wd += "\\"
	}
	return wd + path, nil
}
