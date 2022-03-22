package util

import (
	"log"

	"github.com/rwxrob/json"
)

// NeedsUpdate compares the current isosec integer (second at GMT) to
// that retrieved from the URL (usually pointing to a file called
// UPDATED) which must return nothing but a single isosec integer (which
// is unmarshaled as a JSON number). Returns 1 if needed, 0 if not
// needed, and -1 if unable to determine (and will log any error
// encountered).
func NeedsUpdate(current int, url string) int {
	var remote int
	if err := json.Req(`GET`, url, nil, nil, &remote); err != nil {
		log.Print(err)
		return -1
	}
	if remote > current {
		return 1
	}
	return 0
}
