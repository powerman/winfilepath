// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Stub godebug package for winfilepath tests
package godebug

// New creates a new debug setting stub
func New(name string) *Setting {
	return &Setting{name: name}
}

// Setting represents a debug setting stub
type Setting struct {
	name string
}

// Value returns empty string stub
func (s *Setting) Value() string {
	return ""
}
