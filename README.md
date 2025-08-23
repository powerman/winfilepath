# winfilepath is a copy of path/filepath with Windows behavior on all platforms

[![License MIT](https://img.shields.io/badge/license-MIT-royalblue.svg)](LICENSE)
[![Go version](https://img.shields.io/github/go-mod/go-version/powerman/winfilepath?color=blue)](https://go.dev/)
[![Test](https://img.shields.io/github/actions/workflow/status/powerman/winfilepath/test.yml?label=test)](https://github.com/powerman/winfilepath/actions/workflows/test.yml)
[![Coverage Status](https://raw.githubusercontent.com/powerman/winfilepath/gh-badges/coverage.svg)](https://github.com/powerman/winfilepath/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/powerman/winfilepath)](https://goreportcard.com/report/github.com/powerman/winfilepath)
[![Release](https://img.shields.io/github/v/release/powerman/winfilepath?color=blue)](https://github.com/powerman/winfilepath/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/powerman/winfilepath.svg)](https://pkg.go.dev/github.com/powerman/winfilepath)

![Linux | amd64 arm64 armv7 ppc64le s390x riscv64](https://img.shields.io/badge/Linux-amd64%20arm64%20armv7%20ppc64le%20s390x%20riscv64-royalblue)
![macOS | amd64 arm64](https://img.shields.io/badge/macOS-amd64%20arm64-royalblue)
![Windows | amd64 arm64](https://img.shields.io/badge/Windows-amd64%20arm64-royalblue)

Most of the code is copied from Go 1.25.0 path/filepath package
with minimal modification needed to make it work on all OSes in same way as on Windows,
except for actual file access functions which panics on non-Windows.

It is mostly useful to process Windows file paths on non-Windows OSes,
e.g. test your code which runs on Windows on your native OS.
