/*
Copyright 2021 Robert S. Muhlestein.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"os"
	"os/exec"
	"testing"
)

// go coverage detection is fucked for this sort of stuff, oh well, we
// did the test even if coverage falsely reports 50%
func TestExec(t *testing.T) {
	if os.Getenv("TESTING_EXEC") == "1" {
		err := Exec("go", "version")
		if err != nil {
			t.Fatal(err)
		}
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestExec")
	cmd.Env = append(os.Environ(), "TESTING_EXEC=1")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("process exited with %v", err)
	}
}

func TestExec_noargs(t *testing.T) {
	err := Exec()
	if err == nil {
		t.Error("should have failed since no command")
	}
}

func TestRun_noargs(t *testing.T) {
	err := Run()
	if err == nil {
		t.Error("should have failed since no command")
	}
}

func TestRun_nocmd(t *testing.T) {
	err := Run("__inoexist")
	if err == nil {
		t.Error("should have failed since no command")
	}
}

func TestExec_nocmd(t *testing.T) {
	err := Exec("__inoexist")
	if err == nil {
		t.Error("should have failed since no command")
	}
}

func TestRun(t *testing.T) {
	err := Run("go", "version")
	if err != nil {
		t.Error(err)
	}
}
