// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package run_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/rwxrob/bonzai/pkg/run"
)

// go coverage detection is fucked for this sort of stuff, oh well, we
// did the test even if coverage falsely reports 50%
func TestSysExe(t *testing.T) {
	if os.Getenv("TESTING_EXEC") == "1" {
		err := run.Exe("go", "version")
		if err != nil {
			t.Fatal(err)
		}
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestSysExe")
	cmd.Env = append(os.Environ(), "TESTING_EXEC=1")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("process exited with %v", err)
	}
}

func TestSysExe_noargs(t *testing.T) {
	err := run.SysExe()
	if err == nil {
		t.Error("should have failed since no command")
	}
}

func TestExe_noargs(t *testing.T) {
	err := run.Exe()
	if err == nil {
		t.Error("should have failed since no command")
	}
}

func TestExe_nocmd(t *testing.T) {
	err := run.Exe("__inoexist")
	if err == nil {
		t.Error("should have failed since no command")
	}
}

func TestExe(t *testing.T) {
	err := run.Exe("go", "version")
	if err != nil {
		t.Error(err)
	}
}
