package bonzai

import (
	"log"

	"github.com/rwxrob/json"
	"golang.org/x/mod/semver"
)

// CompareUpdated compares the current isosec integer (second at GMT) to
// that retrieved from the URL (usually pointing to a file called
// UPDATED) which must return nothing but a single isosec integer (which
// is unmarshaled as a JSON number). Returns 1 if current is more
// recent, 0 if they are equal, -1 if current is older, and -2 if unable
// to determine (which will also log any error encountered).
func CompareUpdated(current int, remoteURL string) int {
	var remote int
	if err := json.Req(`GET`, remoteURL, nil, nil, &remote); err != nil {
		log.Print(err)
		return -2
	}
	if remote == current {
		return 0
	}
	if remote > current {
		return 1
	}
	return -1
}

// CompareVersions compares the current semantic version string to that
// retrieved from the URL (usually pointing to a file called VERSION)
// which must return nothing but a single semantic version. Like
// semver.Compare returns 0 if they are the same, 1 if current is more
// recent, -1 if current is behind remote, and -2 if unable to retrieve
// the remote semver string (in which case an error is also logged).
func CompareVersions(current, remoteURL string) int {
	var remote string
	if err := json.Req(`GET`, remoteURL, nil, nil, &remote); err != nil {
		log.Print(err)
		return -2
	}
	return semver.Compare(current, remote)
}
