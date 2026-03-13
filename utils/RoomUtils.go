package utils

import (
	"strconv"
	"strings"
)

// isRoomLine checks whether a line looks like `name x y`.
func isRoomLine(s string) bool {
	// room: 3 tokens, no '-' token style
	parts := strings.Fields(s)
	return len(parts) == 3 && !strings.Contains(parts[0], "-")
}

// parseRoom converts a room line into a Room after validating its name and coordinates.
func parseRoom(s string) (*Room, error) {
	parts := strings.Fields(s)
	if len(parts) != 3 {
		return nil, ErrInvalid
	}
	name := parts[0]
	if name == "" || strings.HasPrefix(name, "L") || strings.HasPrefix(name, "#") || strings.ContainsAny(name, " \t") {
		return nil, ErrInvalid
	}
	x, err1 := strconv.Atoi(parts[1])
	y, err2 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil {
		return nil, ErrInvalid
	}
	return &Room{Name: name, X: x, Y: y}, nil
}
