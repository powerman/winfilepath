// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE-go file.

// Windows path separator functions for non-Windows platforms

package winos

// IsPathSeparator reports whether c is a directory separator character.
// On Windows, this is backslash (\) and forward slash (/).
func IsPathSeparator(c uint8) bool {
	return c == '\\' || c == '/'
}

// PathSeparator is the Windows path separator.
const PathSeparator = '\\'

// PathListSeparator is the Windows path list separator.
const PathListSeparator = ';'
