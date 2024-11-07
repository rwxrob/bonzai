package kimono

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/mod/semver"

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

func TagBump(part VerPart, mustPush bool) error {
	versions := TagList(true)
	latest := ``
	if len(versions) == 0 {
		latest = `v0.0.0`
	} else {
		latest = versions[len(versions)-1]
	}
	prefix := modulePrefix()
	newVersion := fmt.Sprintf(
		`%s%s`,
		prefix,
		versionBump(latest, part),
	)
	fmt.Println(newVersion)
	if err := run.Exec(`git`, `tag`, newVersion); err != nil {
		return err
	}
	if mustPush {
		if err := run.Exec(`git`, `push`, `origin`, newVersion); err != nil {
			return err
		}
	}
	return nil
}

// versionBump increases the given part of the version.
func versionBump(version string, part VerPart) string {
	leading := ``
	versionN := version
	if leading == `v` {
		leading = `v`
		versionN = version[1:]
	}
	versionParts := strings.Split(versionN, `.`)

	// Bump the specified version part
	switch part {
	case Major:
		versionParts[0] = fmt.Sprintf(
			`%d`,
			1+parseInt(versionParts[0]),
		)
	case Minor:
		versionParts[1] = fmt.Sprintf(
			`%d`,
			1+parseInt(versionParts[1]),
		)
	case Patch:
		versionParts[2] = fmt.Sprintf(
			`%d`,
			1+parseInt(versionParts[2]),
		)
	}

	return fmt.Sprint(leading, strings.Join(versionParts, `.`))
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

// TagList returns the list of tags for the current module.
func TagList(shorten bool) []string {
	prefix := modulePrefix()
	tags := run.Out(`git`, `tag`, `-l`, `--no-column`)
	out := make([]string, 0)
	each.Do(strings.Split(tags, "\n"), func(tag string) {
		if isValidTag(tag, prefix) {
			if shorten {
				tag = strings.TrimPrefix(tag, prefix)
			}
			out = append(out, tag)
		}
	})
	semver.Sort(out)
	return out
}

func isValidTag(tag, prefix string) bool {
	return (len(prefix) > 0 && strings.HasPrefix(tag, prefix)) ||
		(len(prefix) == 0 && semver.IsValid(tag))
}

func modulePrefix() string {
	root, err := futil.HereOrAbove(`.git`)
	if err != nil {
		return ``
	}
	module, err := futil.HereOrAbove(`go.mod`)
	if err != nil {
		return ``
	}
	outprefix, err := filepath.Rel(
		filepath.Dir(root),
		filepath.Dir(module),
	)
	if err != nil {
		return ``
	}
	if outprefix == `.` {
		return ``
	}
	return outprefix + `/`
}
