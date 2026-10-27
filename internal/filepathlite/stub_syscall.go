// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE-go file.

// This file contains stub implementations of Windows syscall functions
// needed by the filepathlite package to work on non-Windows platforms.

package filepathlite

import "syscall"

// UTF16PtrFromString stub for non-Windows platforms.
func UTF16PtrFromString(_ string) (*uint16, error) {
	// This is a simplified stub that just returns an error
	// The real implementation would convert string to UTF16
	return nil, syscall.EINVAL
}

// RtlIsDosDeviceName_U stub - simplified logic for non-Windows.
func RtlIsDosDeviceName_U(p *uint16) uint32 { //nolint:revive // Name is used by vendored code.
	// For simplified implementation, just return 0 (not a DOS device)
	// This will make reserved names with extensions be considered reserved
	// which is more conservative than Windows 11 behavior but safer
	return 0
}
