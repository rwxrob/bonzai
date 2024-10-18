// Copyright 2022 Robert S. Muhlestein
// SPDX-License-Identifier: Apache-2.0

package term

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

// coverage will never catch this test

func TestIsInteractive_false(t *testing.T) {
	if os.Getenv("TEST_ISNOTTERM") == "1" {
		fmt.Println("out")
		if !IsInteractive() {
			os.Exit(20)
		}
		os.Exit(1)
	}
	exe := os.Args[0]
	cmd := exec.Command(exe, "-test.run=TestIsInteractive_false")
	cmd.Env = append(os.Environ(), "TEST_ISNOTTERM=1")
	cmd.StdoutPipe() // just enough to push into background
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok {
		t.Log(e.ExitCode())
		if e.ExitCode() != 20 {
			t.Errorf("exit %v: still a terminal", e.ExitCode())
		}
	}
}

func TestIsInteractive_true(t *testing.T) {
	if !IsInteractive() {
		t.Error("terminal not connected")
	}
}
