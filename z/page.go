package Z

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rwxrob/to"
)

// NoPager disables all paging.
var NoPager bool

// FindPager returns a full path to a pager binary if it can find one on
// the system:
//
//     * $PAGER
//     * pager (in PATH)
//
// If neither is found returns an empty string.
func FindPager() string {
	if NoPager {
		return ""
	}
	path := os.Getenv(`PAGER`)
	if path == "" {
		path, _ = exec.LookPath(path)
	}
	if path == "" {
		path, _ = exec.LookPath(`pager`)
	}
	return path
}

// PageFile looks up the system pager and passes the first argument as
// a file argument to it. Just prints file if no pager can be found.
func PageFile(path string) error {
	pager := FindPager()
	if pager == "" || NoPager {
		buf, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fmt.Print(string(buf))
		return nil
	}
	return Exec(pager, path)
}

// Page writes the buf to a temporary file and passes it as first
// argument to the detected system pager. Anything that to.String
// accepts can be passed. Just prints output without paging if a pager
// cannot be found.
func Page[T any](buf T) error {
	pager := FindPager()
	if pager == "" || NoPager {
		fmt.Print(to.String(buf))
		return nil
	}
	f, err := os.CreateTemp("", `bonzai-page-*`)
	if err != nil {
		return err
	}
	name := f.Name()
	_, err = f.WriteString(to.String(buf))
	defer f.Close()
	defer os.Remove(name)
	if err != nil {
		return err
	}
	return Exec(pager, name)
}
