package kimono

import (
	"fmt"
	"path/filepath"
	"slices"
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
	newVer, err := versionBump(latest, part)
	if err != nil {
		return fmt.Errorf(`failed to bump version: %w`, err)
	}
	newVerStr := fmt.Sprintf(`%s%s`, prefix, newVer)
	fmt.Println(newVerStr)
	if err := run.Exec(`git`, `tag`, newVerStr); err != nil {
		return err
	}
	if mustPush {
		if err := run.Exec(`git`, `push`, `origin`, newVerStr); err != nil {
			return err
		}
	}
	return nil
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

// TagDelete deletes the given tag from local git repository.
func TagDelete(tag string, remote bool) error {
	tags := TagList(false)
	if !slices.Contains(tags, tag) {
		return fmt.Errorf("tag '%s' not found", tag)
	}
	if err := run.Exec(`git`, `tag`, `-d`, tag); err != nil {
		return fmt.Errorf(
			"failed to delete local tag '%s': %w",
			tag,
			err,
		)
	}
	fmt.Println(`Deleted local tag:`, tag)
	return nil
}

// versionBump increases the given part of the version.
func versionBump(version string, part VerPart) (string, error) {
	leading := version[:1]
	versionN := version[1:]
	if leading != `v` {
		leading = `v`
		versionN = version
	}
	versionParts := strings.Split(versionN, `.`)
	switch part {
	case Major:
		major, err := strconv.Atoi(versionParts[0])
		if err != nil {
			return ``, err
		}
		versionParts[0] = strconv.Itoa(major + 1)
		versionParts[1] = `0`
		versionParts[2] = `0`
	case Minor:
		minor, err := strconv.Atoi(versionParts[1])
		if err != nil {
			return ``, err
		}
		versionParts[1] = strconv.Itoa(minor + 1)
		versionParts[2] = `0`
	case Patch:
		patch, err := strconv.Atoi(versionParts[2])
		if err != nil {
			return ``, err
		}
		versionParts[2] = strconv.Itoa(patch + 1)
	}

	return fmt.Sprint(leading, strings.Join(versionParts, `.`)), nil
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
