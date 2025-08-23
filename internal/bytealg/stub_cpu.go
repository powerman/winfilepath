// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Stub implementation of internal/cpu for bytealg compatibility
package bytealg

// Minimal cpu features struct to satisfy bytealg requirements
var cpu struct {
	X86 struct {
		HasSSE42  bool
		HasAVX2   bool
		HasPOPCNT bool
	}
	Loong64 struct {
		HasLSX  bool
		HasLASX bool
	}
	S390X struct {
		HasVX bool
	}
	PPC64 struct {
		IsPOWER9 bool
	}
}
