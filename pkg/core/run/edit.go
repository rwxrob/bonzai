package run

import (
	"fmt"
	"os"
	"os/exec"
)

// Edit opens the file at the given path for editing searching for an
// editor on the system using the following (in order of priority):
//
// * VISUAL
// * EDITOR
// * code
// * nvim
// * vim
// * vi
// * nvi
// * emacs
// * nano
func Edit(path string) error {
	ed := os.Getenv("VISUAL")
	if ed != "" {
		return Exec(ed, path)
	}
	ed = os.Getenv("EDITOR")
	if ed != "" {
		return Exec(ed, path)
	}
	ed, _ = exec.LookPath("code")
	if ed != "" {
		return Exec(ed, path)
	}
	ed, _ = exec.LookPath("nvim")
	if ed != "" {
		return Exec(ed, path)
	}
	ed, _ = exec.LookPath("vim")
	if ed != "" {
		return Exec(ed, path)
	}
	ed, _ = exec.LookPath("vi")
	if ed != "" {
		return Exec(ed, path)
	}
	ed, _ = exec.LookPath("nvi")
	if ed != "" {
		return Exec(ed, path)
	}
	ed, _ = exec.LookPath("emacs")
	if ed != "" {
		return Exec(ed, path)
	}
	ed, _ = exec.LookPath("nano")
	if ed != "" {
		return Exec(ed, path)
	}
	return fmt.Errorf("unable to find editor")
}
