package edit

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rwxrob/bonzai/run"
)

var EditorPriority = []string{
	os.Getenv(`VISUAL`),
	os.Getenv(`EDITOR`),
	`code`, `nvim`, `vim`, `vi`, `nvi`, `emacs`, `nano`,
}

// Editor returns the first editor found from [EditorPriority] or blank
// if none found.
func Editor() string {
	for _, name := range EditorPriority {
		if len(name) == 0 {
			continue
		}
		editor, _ := exec.LookPath(name)
		if len(editor) > 0 {
			return editor
		}
	}
	return ""
}

// Files opens files with the first editor found from [EditorPriority]
// list.
func Files(files ...string) error {
	editor := Editor()
	if len(editor) == 0 {
		return fmt.Errorf("unable to find editor")
	}
	args := append([]string{editor}, files...)
	return run.Exec(args...)
}
