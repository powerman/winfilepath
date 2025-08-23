// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package winfilepath_test

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"runtime/debug"
	"slices"
	"strings"
	"testing"

	filepath "github.com/powerman/winfilepath"
	"github.com/powerman/winfilepath/internal/godebug"
	"github.com/powerman/winfilepath/internal/testenv"
)

func TestWinSplitListTestsAreValid(t *testing.T) {
	comspec := os.Getenv("ComSpec")
	if comspec == "" {
		t.Fatal("%ComSpec% must be set")
	}

	for ti, tt := range winsplitlisttests {
		testWinSplitListTestIsValid(t, ti, tt, comspec)
	}
}

func testWinSplitListTestIsValid(t *testing.T, ti int, tt SplitListTest,
	comspec string,
) {
	const (
		cmdfile             = `printdir.cmd`
		perm    fs.FileMode = 0o700
	)

	tmp := t.TempDir()
	for i, d := range tt.result {
		if d == "" {
			continue
		}
		if cd := filepath.Clean(d); filepath.VolumeName(cd) != "" ||
			cd[0] == '\\' || cd == ".." || (len(cd) >= 3 && cd[0:3] == `..\`) {
			t.Errorf("%d,%d: %#q refers outside working directory", ti, i, d)
			return
		}
		dd := filepath.Join(tmp, d)
		if _, err := os.Stat(dd); err == nil {
			t.Errorf("%d,%d: %#q already exists", ti, i, d)
			return
		}
		if err := os.MkdirAll(dd, perm); err != nil {
			t.Errorf("%d,%d: MkdirAll(%#q) failed: %v", ti, i, dd, err)
			return
		}
		fn, data := filepath.Join(dd, cmdfile), []byte("@echo "+d+"\r\n")
		if err := os.WriteFile(fn, data, perm); err != nil {
			t.Errorf("%d,%d: WriteFile(%#q) failed: %v", ti, i, fn, err)
			return
		}
	}

	// on some systems, SystemRoot is required for cmd to work
	systemRoot := os.Getenv("SystemRoot")

	for i, d := range tt.result {
		if d == "" {
			continue
		}
		exp := []byte(d + "\r\n")
		cmd := &exec.Cmd{
			Path: comspec,
			Args: []string{`/c`, cmdfile},
			Env:  []string{`Path=` + systemRoot + "/System32;" + tt.list, `SystemRoot=` + systemRoot},
			Dir:  tmp,
		}
		out, err := cmd.CombinedOutput()
		switch {
		case err != nil:
			t.Errorf("%d,%d: execution error %v\n%q", ti, i, err, out)
			return
		case !slices.Equal(out, exp):
			t.Errorf("%d,%d: expected %#q, got %#q", ti, i, exp, out)
			return
		default:
			// unshadow cmdfile in next directory
			err = os.Remove(filepath.Join(tmp, d, cmdfile))
			if err != nil {
				t.Fatalf("Remove test command failed: %v", err)
			}
		}
	}
}

func TestWindowsEvalSymlinks(t *testing.T) {
	testenv.MustHaveSymlink(t)

	tmpDir := tempDirCanonical(t)

	if len(tmpDir) < 3 {
		t.Fatalf("tmpDir path %q is too short", tmpDir)
	}
	if tmpDir[1] != ':' {
		t.Fatalf("tmpDir path %q must have drive letter in it", tmpDir)
	}
	test := EvalSymlinksTest{"test/linkabswin", tmpDir[:3]}

	// Create the symlink farm using relative paths.
	testdirs := append(EvalSymlinksTestDirs, test)
	for _, d := range testdirs {
		var err error
		path := simpleJoin(tmpDir, d.path)
		if d.dest == "" {
			err = os.Mkdir(path, 0o755)
		} else {
			err = os.Symlink(d.dest, path)
		}
		if err != nil {
			t.Fatal(err)
		}
	}

	path := simpleJoin(tmpDir, test.path)

	testEvalSymlinks(t, path, test.dest)

	testEvalSymlinksAfterChdir(t, path, ".", test.dest)

	testEvalSymlinksAfterChdir(t,
		path,
		filepath.VolumeName(tmpDir)+".",
		test.dest)

	testEvalSymlinksAfterChdir(t,
		simpleJoin(tmpDir, "test"),
		simpleJoin("..", test.path),
		test.dest)

	testEvalSymlinksAfterChdir(t, tmpDir, test.path, test.dest)
}

// TestEvalSymlinksCanonicalNames verify that EvalSymlinks
// returns "canonical" path names on windows.
func TestEvalSymlinksCanonicalNames(t *testing.T) {
	ctmp := tempDirCanonical(t)
	dirs := []string{
		"test",
		"test/dir",
		"testing_long_dir",
		"TEST2",
	}

	for _, d := range dirs {
		dir := filepath.Join(ctmp, d)
		err := os.Mkdir(dir, 0o755)
		if err != nil {
			t.Fatal(err)
		}
		cname, err := filepath.EvalSymlinks(dir)
		if err != nil {
			t.Errorf("EvalSymlinks(%q) error: %v", dir, err)
			continue
		}
		if dir != cname {
			t.Errorf("EvalSymlinks(%q) returns %q, but should return %q", dir, cname, dir)
			continue
		}
		// test non-canonical names
		test := strings.ToUpper(dir)
		p, err := filepath.EvalSymlinks(test)
		if err != nil {
			t.Errorf("EvalSymlinks(%q) error: %v", test, err)
			continue
		}
		if p != cname {
			t.Errorf("EvalSymlinks(%q) returns %q, but should return %q", test, p, cname)
			continue
		}
		// another test
		test = strings.ToLower(dir)
		p, err = filepath.EvalSymlinks(test)
		if err != nil {
			t.Errorf("EvalSymlinks(%q) error: %v", test, err)
			continue
		}
		if p != cname {
			t.Errorf("EvalSymlinks(%q) returns %q, but should return %q", test, p, cname)
			continue
		}
	}
}

// checkVolume8dot3Setting runs "fsutil 8dot3name query c:" command
// (where c: is vol parameter) to discover "8dot3 name creation state".
// The state is combination of 2 flags. The global flag controls if it
// is per volume or global setting:
//
//	0 - Enable 8dot3 name creation on all volumes on the system
//	1 - Disable 8dot3 name creation on all volumes on the system
//	2 - Set 8dot3 name creation on a per volume basis
//	3 - Disable 8dot3 name creation on all volumes except the system volume
//
// If global flag is set to 2, then per-volume flag needs to be examined:
//
//	0 - Enable 8dot3 name creation on this volume
//	1 - Disable 8dot3 name creation on this volume
//
// checkVolume8dot3Setting verifies that "8dot3 name creation" flags
// are set to 2 and 0, if enabled parameter is true, or 2 and 1, if enabled
// is false. Otherwise checkVolume8dot3Setting returns error.
func checkVolume8dot3Setting(vol string, enabled bool) error {
	// It appears, on some systems "fsutil 8dot3name query ..." command always
	// exits with error. Ignore exit code, and look at fsutil output instead.
	out, _ := exec.Command("fsutil", "8dot3name", "query", vol).CombinedOutput()
	// Check that system has "Volume level setting" set.
	expected := "The registry state of NtfsDisable8dot3NameCreation is 2, the default (Volume level setting)"
	if !strings.Contains(string(out), expected) {
		// Windows 10 version of fsutil has different output message.
		expectedWindow10 := "The registry state is: 2 (Per volume setting - the default)"
		if !strings.Contains(string(out), expectedWindow10) {
			return fmt.Errorf("fsutil output should contain %q, but is %q", expected, string(out))
		}
	}
	// Now check the volume setting.
	expected = "Based on the above two settings, 8dot3 name creation is %s on %s"
	if enabled {
		expected = fmt.Sprintf(expected, "enabled", vol)
	} else {
		expected = fmt.Sprintf(expected, "disabled", vol)
	}
	if !strings.Contains(string(out), expected) {
		return fmt.Errorf("unexpected fsutil output: %q", string(out))
	}
	return nil
}

func setVolume8dot3Setting(vol string, enabled bool) error {
	cmd := []string{"fsutil", "8dot3name", "set", vol}
	if enabled {
		cmd = append(cmd, "0")
	} else {
		cmd = append(cmd, "1")
	}
	// It appears, on some systems "fsutil 8dot3name set ..." command always
	// exits with error. Ignore exit code, and look at fsutil output instead.
	out, _ := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if string(out) != "\r\nSuccessfully set 8dot3name behavior.\r\n" {
		// Windows 10 version of fsutil has different output message.
		expectedWindow10 := "Successfully %s 8dot3name generation on %s\r\n"
		if enabled {
			expectedWindow10 = fmt.Sprintf(expectedWindow10, "enabled", vol)
		} else {
			expectedWindow10 = fmt.Sprintf(expectedWindow10, "disabled", vol)
		}
		if string(out) != expectedWindow10 {
			return fmt.Errorf("%v command failed: %q", cmd, string(out))
		}
	}
	return nil
}

var runFSModifyTests = flag.Bool("run_fs_modify_tests", false, "run tests which modify filesystem parameters")

// This test assumes registry state of NtfsDisable8dot3NameCreation is 2,
// the default (Volume level setting).
func TestEvalSymlinksCanonicalNamesWith8dot3Disabled(t *testing.T) {
	if !*runFSModifyTests {
		t.Skip("skipping test that modifies file system setting; enable with -run_fs_modify_tests")
	}
	tempVol := filepath.VolumeName(os.TempDir())
	if len(tempVol) != 2 {
		t.Fatalf("unexpected temp volume name %q", tempVol)
	}

	err := checkVolume8dot3Setting(tempVol, true)
	if err != nil {
		t.Fatal(err)
	}
	err = setVolume8dot3Setting(tempVol, false)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := setVolume8dot3Setting(tempVol, true)
		if err != nil {
			t.Fatal(err)
		}
		err = checkVolume8dot3Setting(tempVol, true)
		if err != nil {
			t.Fatal(err)
		}
	}()
	err = checkVolume8dot3Setting(tempVol, false)
	if err != nil {
		t.Fatal(err)
	}
	TestEvalSymlinksCanonicalNames(t)
}

func TestUNC(t *testing.T) {
	// Test that this doesn't go into an infinite recursion.
	// See golang.org/issue/15879.
	defer debug.SetMaxStack(debug.SetMaxStack(1e6))
	filepath.Glob(`\\?\c:\*`)
}

func testWalkMklink(t *testing.T, linktype string) {
	output, _ := exec.Command("cmd", "/c", "mklink", "/?").Output()
	if !strings.Contains(string(output), fmt.Sprintf(" /%s ", linktype)) {
		t.Skipf(`skipping test; mklink does not supports /%s parameter`, linktype)
	}
	testWalkSymlink(t, func(target, link string) error {
		output, err := exec.Command("cmd", "/c", "mklink", "/"+linktype, link, target).CombinedOutput()
		if err != nil {
			return fmt.Errorf(`"mklink /%s %v %v" command failed: %v\n%v`, linktype, link, target, err, string(output))
		}
		return nil
	})
}

func TestWalkDirectoryJunction(t *testing.T) {
	testenv.MustHaveSymlink(t)
	testWalkMklink(t, "J")
}

func TestWalkDirectorySymlink(t *testing.T) {
	testenv.MustHaveSymlink(t)
	testWalkMklink(t, "D")
}

func createMountPartition(t *testing.T, vhd string, args string) []byte {
	testenv.MustHaveExecPath(t, "powershell")
	t.Cleanup(func() {
		cmd := testenv.Command(t, "powershell", "-Command", fmt.Sprintf("Dismount-VHD %q", vhd))
		out, err := cmd.CombinedOutput()
		if err != nil {
			if t.Skipped() {
				// Probably failed to dismount because we never mounted it in
				// the first place. Log the error, but ignore it.
				t.Logf("%v: %v (skipped)\n%s", cmd, err, out)
			} else {
				// Something went wrong, and we don't want to leave dangling VHDs.
				// Better to fail the test than to just log the error and continue.
				t.Errorf("%v: %v\n%s", cmd, err, out)
			}
		}
	})

	script := filepath.Join(t.TempDir(), "test.ps1")
	cmd := strings.Join([]string{
		"$ErrorActionPreference = \"Stop\"",
		fmt.Sprintf("$vhd = New-VHD -Path %q -SizeBytes 3MB -Fixed", vhd),
		"$vhd | Mount-VHD",
		fmt.Sprintf("$vhd = Get-VHD %q", vhd),
		"$vhd | Get-Disk | Initialize-Disk -PartitionStyle GPT",
		"$part = $vhd | Get-Disk | New-Partition -UseMaximumSize -AssignDriveLetter:$false",
		"$vol = $part | Format-Volume -FileSystem NTFS",
		args,
	}, "\n")

	err := os.WriteFile(script, []byte(cmd), 0o666)
	if err != nil {
		t.Fatal(err)
	}
	output, err := testenv.Command(t, "powershell", "-File", script).CombinedOutput()
	if err != nil {
		// This can happen if Hyper-V is not installed or enabled.
		t.Skip("skipping test because failed to create VHD: ", err, string(output))
	}
	return output
}

var (
	winsymlink        = godebug.New("winsymlink")
	winreadlinkvolume = godebug.New("winreadlinkvolume")
)

func TestEvalSymlinksJunctionToVolumeID(t *testing.T) {
	// Test that EvalSymlinks resolves a directory junction which
	// is mapped to volumeID (instead of drive letter). See go.dev/issue/39786.
	if winsymlink.Value() == "0" {
		t.Skip("skipping test because winsymlink is not enabled")
	}
	t.Parallel()

	output, _ := exec.Command("cmd", "/c", "mklink", "/?").Output()
	if !strings.Contains(string(output), " /J ") {
		t.Skip("skipping test because mklink command does not support junctions")
	}

	tmpdir := tempDirCanonical(t)
	vhd := filepath.Join(tmpdir, "Test.vhdx")
	output = createMountPartition(t, vhd, "Write-Host $vol.Path -NoNewline")
	vol := string(output)

	dirlink := filepath.Join(tmpdir, "dirlink")
	output, err := testenv.Command(t, "cmd", "/c", "mklink", "/J", dirlink, vol).CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run mklink %v %v: %v %q", dirlink, vol, err, output)
	}
	got, err := filepath.EvalSymlinks(dirlink)
	if err != nil {
		t.Fatal(err)
	}
	if got != dirlink {
		t.Errorf(`EvalSymlinks(%q): got %q, want %q`, dirlink, got, dirlink)
	}
}

func TestEvalSymlinksMountPointRecursion(t *testing.T) {
	// Test that EvalSymlinks doesn't follow recursive mount points.
	// See go.dev/issue/40176.
	if winsymlink.Value() == "0" {
		t.Skip("skipping test because winsymlink is not enabled")
	}
	t.Parallel()

	tmpdir := tempDirCanonical(t)
	dirlink := filepath.Join(tmpdir, "dirlink")
	err := os.Mkdir(dirlink, 0o755)
	if err != nil {
		t.Fatal(err)
	}

	vhd := filepath.Join(tmpdir, "Test.vhdx")
	createMountPartition(t, vhd, fmt.Sprintf("$part | Add-PartitionAccessPath -AccessPath %q\n", dirlink))

	got, err := filepath.EvalSymlinks(dirlink)
	if err != nil {
		t.Fatal(err)
	}
	if got != dirlink {
		t.Errorf(`EvalSymlinks(%q): got %q, want %q`, dirlink, got, dirlink)
	}
}

func TestNTNamespaceSymlink(t *testing.T) {
	output, _ := exec.Command("cmd", "/c", "mklink", "/?").Output()
	if !strings.Contains(string(output), " /J ") {
		t.Skip("skipping test because mklink command does not support junctions")
	}

	tmpdir := tempDirCanonical(t)

	vol := filepath.VolumeName(tmpdir)
	output, err := exec.Command("cmd", "/c", "mountvol", vol, "/L").CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run mountvol %v /L: %v %q", vol, err, output)
	}
	target := strings.Trim(string(output), " \n\r")

	dirlink := filepath.Join(tmpdir, "dirlink")
	output, err = exec.Command("cmd", "/c", "mklink", "/J", dirlink, target).CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run mklink %v %v: %v %q", dirlink, target, err, output)
	}

	got, err := filepath.EvalSymlinks(dirlink)
	if err != nil {
		t.Fatal(err)
	}
	var want string
	if winsymlink.Value() == "0" {
		if winreadlinkvolume.Value() == "0" {
			want = vol + `\`
		} else {
			want = target
		}
	} else {
		want = dirlink
	}
	if got != want {
		t.Errorf(`EvalSymlinks(%q): got %q, want %q`, dirlink, got, want)
	}

	// Make sure we have sufficient privilege to run mklink command.
	testenv.MustHaveSymlink(t)

	file := filepath.Join(tmpdir, "file")
	err = os.WriteFile(file, []byte(""), 0o666)
	if err != nil {
		t.Fatal(err)
	}

	target = filepath.Join(target, file[len(filepath.VolumeName(file)):])

	filelink := filepath.Join(tmpdir, "filelink")
	output, err = exec.Command("cmd", "/c", "mklink", filelink, target).CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run mklink %v %v: %v %q", filelink, target, err, output)
	}

	got, err = filepath.EvalSymlinks(filelink)
	if err != nil {
		t.Fatal(err)
	}

	if winreadlinkvolume.Value() == "0" {
		want = file
	} else {
		want = target
	}
	if got != want {
		t.Errorf(`EvalSymlinks(%q): got %q, want %q`, filelink, got, want)
	}
}

func TestIssue52476(t *testing.T) {
	tests := []struct {
		lhs, rhs string
		want     string
	}{
		{`..\.`, `C:`, `..\C:`},
		{`..`, `C:`, `..\C:`},
		{`.`, `:`, `.\:`},
		{`.`, `C:`, `.\C:`},
		{`.`, `C:/a/b/../c`, `.\C:\a\c`},
		{`.`, `\C:`, `.\C:`},
		{`C:\`, `.`, `C:\`},
		{`C:\`, `C:\`, `C:\C:`},
		{`C`, `:`, `C\:`},
		{`\.`, `C:`, `\C:`},
		{`\`, `C:`, `\C:`},
	}

	for _, test := range tests {
		got := filepath.Join(test.lhs, test.rhs)
		if got != test.want {
			t.Errorf(`Join(%q, %q): got %q, want %q`, test.lhs, test.rhs, got, test.want)
		}
	}
}

func TestAbsWindows(t *testing.T) {
	for _, test := range []struct {
		path string
		want string
	}{
		{`C:\foo`, `C:\foo`},
		{`\\host\share\foo`, `\\host\share\foo`},
		{`\\host`, `\\host`},
		{`\\.\NUL`, `\\.\NUL`},
		{`NUL`, `\\.\NUL`},
		{`COM1`, `\\.\COM1`},
		{`a/NUL`, `\\.\NUL`},
	} {
		got, err := filepath.Abs(test.path)
		if err != nil || got != test.want {
			t.Errorf("Abs(%q) = %q, %v; want %q, nil", test.path, got, err, test.want)
		}
	}
}
