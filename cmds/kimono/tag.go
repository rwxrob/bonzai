package kimono

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rwxrob/bonzai/fn/each"
	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/run"
)

type VerPart int

const (
	Major VerPart = 0
	Minor VerPart = 1
	Patch VerPart = 2
)

// TagBump identifies the current module path, identifies the latest
// version tag and tags the repo with the bumped version.
func TagBump(part VerPart) error {
	return nil
}

// TagList returns the list of tags for the current module.
func TagList() []string {
	prefix := modulePrefix()
	tags := run.Out(`git`, `tag`, `-l`, `--no-column`)
	out := make([]string, 0)
	each.Do(strings.Split(tags, "\n"), func(tag string) {
		if strings.HasPrefix(tag, prefix) ||
			!strings.Contains(tag, "/") {
			out = append(
				out,
				strings.TrimPrefix(tag, prefix),
			)
		}
	})
	return out
}

func modulePrefix() string {
	root, err := futil.HereOrAbove(".git")
	if err != nil {
		return ""
	}
	module, err := futil.HereOrAbove("go.mod")
	if err != nil {
		return ""
	}
	outprefix, err := filepath.Rel(
		filepath.Dir(root),
		filepath.Dir(module),
	)
	if err != nil {
		return ""
	}
	if outprefix == "." {
		return ""
	}
	return outprefix + "/"
}
