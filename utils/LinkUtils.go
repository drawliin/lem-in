package utils

import "strings"

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

func linkKey(a, b string) string {
	if a < b {
		return a + "|" + b
	}
	return b + "|" + a
}
