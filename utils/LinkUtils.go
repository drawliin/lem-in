package utils

import "strings"

// isLinkLine checks whether a line matches the `roomA-roomB` link format.
func isLinkLine(s string) bool {
	// exactly one '-', not at ends, and no spaces
	if strings.ContainsAny(s, " \t") {
		return false
	}
	if strings.Count(s, "-") != 1 {
		return false
	}
	if strings.HasPrefix(s, "-") || strings.HasSuffix(s, "-") {
		return false
	}
	return true
}

// parseLink splits a validated link line and rejects illegal room-name prefixes.
func parseLink(s string) (string, string, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", ErrInvalid
	}
	// room names in links can't start with L/# either (safer)
	for _, p := range parts {
		if strings.HasPrefix(p, "L") || strings.HasPrefix(p, "#") {
			return "", "", ErrInvalid
		}
	}
	return parts[0], parts[1], nil
}

// linkKey builds a stable unordered key so `a-b` and `b-a` count as the same link.
func linkKey(a, b string) string {
	if a < b {
		return a + "|" + b
	}
	return b + "|" + a
}
