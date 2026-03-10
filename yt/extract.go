package yt

import (
	"net/url"
	"regexp"
	"strings"
)

var ytID = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

// ExtractID returns the YouTube video ID from either a full URL or an ID.
// Returns empty string if nothing valid is found.
func ExtractID(s string) string {

	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	// If it already looks like an ID, accept it.
	if ytID.MatchString(s) && !strings.Contains(s, "/") {
		return s
	}

	u, err := url.Parse(s)
	if err != nil {
		return ""
	}

	host := strings.ToLower(u.Host)

	switch {

	// youtube.com/watch?v=ID
	case strings.Contains(host, "youtube.com"):
		q := u.Query().Get("v")
		if ytID.MatchString(q) {
			return q
		}

		// /embed/ID
		// /shorts/ID
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) >= 2 && ytID.MatchString(parts[1]) {
			return parts[1]
		}

	// youtu.be/ID
	case strings.Contains(host, "youtu.be"):
		id := strings.Trim(u.Path, "/")
		if ytID.MatchString(id) {
			return id
		}
	}

	return ""
}

func NormalizeURL(s string) string {
	id := ExtractID(s)
	if id == "" {
		return ""
	}
	return "https://youtube.com/watch?v=" + id
}
