//go:build windows

// Package winfilepath implements utility routines for manipulating filename paths
// in a way compatible with Windows file paths.
//
// On Windows, this package simply forwards to the standard path/filepath package.
package winfilepath

import (
	"io/fs"
	"os"
	stdfilepath "path/filepath"
)

// Re-export all constants from standard library
const (
	Separator     = stdfilepath.Separator
	ListSeparator = stdfilepath.ListSeparator
)

// Re-export all variables from standard library
var (
	SkipDir = stdfilepath.SkipDir
	SkipAll = stdfilepath.SkipAll
)

// Forward all functions to standard library
func Abs(path string) (string, error)                 { return stdfilepath.Abs(path) }
func Base(path string) string                         { return stdfilepath.Base(path) }
func Clean(path string) string                        { return stdfilepath.Clean(path) }
func Dir(path string) string                          { return stdfilepath.Dir(path) }
func EvalSymlinks(path string) (string, error)        { return stdfilepath.EvalSymlinks(path) }
func Ext(path string) string                          { return stdfilepath.Ext(path) }
func FromSlash(path string) string                    { return stdfilepath.FromSlash(path) }
func HasPrefix(p, prefix string) bool                 { return stdfilepath.HasPrefix(p, prefix) }
func IsAbs(path string) bool                          { return stdfilepath.IsAbs(path) }
func IsLocal(path string) bool                        { return stdfilepath.IsLocal(path) }
func Join(elem ...string) string                      { return stdfilepath.Join(elem...) }
func Rel(basepath, targpath string) (string, error)   { return stdfilepath.Rel(basepath, targpath) }
func Split(path string) (dir, file string)            { return stdfilepath.Split(path) }
func SplitList(path string) []string                  { return stdfilepath.SplitList(path) }
func ToSlash(path string) string                      { return stdfilepath.ToSlash(path) }
func VolumeName(path string) string                   { return stdfilepath.VolumeName(path) }
func Walk(root string, fn stdfilepath.WalkFunc) error { return stdfilepath.Walk(root, fn) }
func WalkDir(root string, fn fs.WalkDirFunc) error    { return stdfilepath.WalkDir(root, fn) }

// Type aliases
type WalkFunc = stdfilepath.WalkFunc

// For testing
var lstat = os.Lstat

func Localize(path string) (string, error) { return stdfilepath.Localize(path) }
