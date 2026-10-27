// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE-go file.

// Package bytealg is a subset of internal/bytealg with a stub of internal/cpu.
package bytealg

// Minimal cpu features struct to satisfy bytealg requirements.
// It is referenced only by unused offset constants in the vendored bytealg.go.
var cpu struct { //nolint:gochecknoglobals,unused // Required to compile vendored bytealg.go.
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
	RISCV64 struct {
		HasV bool
	}
}
