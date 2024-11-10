package depends

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/rwxrob/bonzai/term"
)

// OnWithHook checks if each dependency in dependencies exists in the
// system PATH. If any dependency is missing, it triggers onError, which
// defaults to defaultOnError if onError is nil.
func OnWithHook(onError func(error), dependencies ...string) {
	for _, dep := range dependencies {
		if !IsInPath(dep) {
			if onError == nil {
				onError = defaultOnError
			}
			onError(fmt.Errorf("program depends on %s", dep))
		}
	}
}

// On checks if each dependency in dependencies exists in the system
// PATH and sends [syscall.SIGTERM] to the PPID if it's not found.
func On(dependencies ...string) {
	OnWithHook(nil, dependencies...)
}

func defaultOnError(err error) {
	fmt.Printf("%s not found in PATH.\n", err)
	if !term.IsInteractive() {
		syscall.Kill(os.Getppid(), syscall.SIGTERM)
	}
}

// IsInPath returns true if the given dependency [dep] is found in the
// system PATH; otherwise, it returns false.
func IsInPath(dep string) bool {
	_, err := exec.LookPath(dep)
	return err == nil
}
