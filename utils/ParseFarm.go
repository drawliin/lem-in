package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var ErrInvalid = errors.New("invalid data format")

func invalidData(reason string) error {
	return fmt.Errorf("%w, %s", ErrInvalid, reason)
}

func ParseFarm(filename string) (raw string, farm *Farm, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", nil, fmt.Errorf("file can't be read")
	}
	raw = string(data)

	f := &Farm{Rooms: make(map[string]*Room)}

	sc := bufio.NewScanner(strings.NewReader(raw))
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	lineNo := 0
	gotAnts := false
	linksMode := false

	nextIsStart := false
	nextIsEnd := false

	// track duplicate links: store "a|b" with sorted names
	seenLink := make(map[string]bool)

	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}

		// comments
		if strings.HasPrefix(line, "#") {
			// commands
			if line == "##start" {
				nextIsStart = true
				nextIsEnd = false
			} else if line == "##end" {
				nextIsEnd = true
				nextIsStart = false
			}
			continue
		}

		// first meaningful line must be ants count
		if !gotAnts {
			n, convErr := strconv.Atoi(line)
			if convErr != nil || n <= 0 {
				return raw, nil, invalidData("invalid number of Ants")
			}
			f.Ants = n
			gotAnts = true
			continue
		}

		// decide if this is a room or a link
		if isRoomLine(line) {
			if linksMode {
				return raw, nil, invalidData(fmt.Sprintf("room found after link declaration at line %d", lineNo))
			}
			r, roomErr := parseRoom(line)
			if roomErr != nil {
				return raw, nil, invalidData(fmt.Sprintf("invalid room definition at line %d", lineNo))
			}
			if _, exists := f.Rooms[r.Name]; exists {
				return raw, nil, invalidData(fmt.Sprintf("duplicate room name at line %d", lineNo))
			}
			f.Rooms[r.Name] = r

			if nextIsStart {
				if f.Start != nil {
					return raw, nil, invalidData("multiple start rooms found")
				}
				f.Start = r
				nextIsStart = false
			} else if nextIsEnd {
				if f.End != nil {
					return raw, nil, invalidData("multiple end rooms found")
				}
				f.End = r
				nextIsEnd = false
			}
			continue
		}

		// link line
		if isLinkLine(line) {
			linksMode = true
			a, b, linkErr := parseLink(line)
			if linkErr != nil {
				return raw, nil, invalidData(fmt.Sprintf("invalid link definition at line %d", lineNo))
			}
			ra := f.Rooms[a]
			rb := f.Rooms[b]
			if ra == nil || rb == nil {
				return raw, nil, invalidData(fmt.Sprintf("link references unknown room at line %d", lineNo))
			}
			if a == b {
				return raw, nil, invalidData(fmt.Sprintf("room links to itself at line %d", lineNo))
			}

			key := linkKey(a, b)
			if seenLink[key] {
				return raw, nil, invalidData(fmt.Sprintf("duplicate link at line %d", lineNo))
			}
			seenLink[key] = true

			ra.Links = append(ra.Links, rb)
			rb.Links = append(rb.Links, ra)
			continue
		}

		// unknown line format
		return raw, nil, invalidData(fmt.Sprintf("invalid line format at line %d", lineNo))
	}

	if scErr := sc.Err(); scErr != nil {
		return raw, nil, invalidData("line too long or scan failed")
	}

	// final validations
	if !gotAnts {
		return raw, nil, invalidData("invalid number of Ants")
	}
	if f.Start == nil {
		return raw, nil, invalidData("no start room found")
	}
	if f.End == nil {
		return raw, nil, invalidData("no end room found")
	}

	// commands without a following room should fail
	if nextIsStart {
		return raw, nil, invalidData("##start must be followed by a room")
	}
	if nextIsEnd {
		return raw, nil, invalidData("##end must be followed by a room")
	}

	return raw, f, nil
}
