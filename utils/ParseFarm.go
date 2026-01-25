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

func ParseFarm(filename string) (raw string, farm *Farm, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", nil, fmt.Errorf("Error in file")
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
			// other comments / unknown ## ignored
			continue
		}

		// first meaningful line must be ants count
		if !gotAnts {
			n, convErr := strconv.Atoi(line)
			if convErr != nil || n <= 0 {
				return raw, nil, fmt.Errorf("convert err: "+line)
			}
			f.Ants = n
			gotAnts = true
			continue
		}

		// decide if this is a room or a link
		if isRoomLine(line) {
			if linksMode {
				return raw, nil, fmt.Errorf("Found Room After Link") // rooms after links => invalid
			}
			r, roomErr := parseRoom(line)
			if roomErr != nil {
				return raw, nil, fmt.Errorf("Err Parsing")
			}
			if _, exists := f.Rooms[r.Name]; exists {
				return raw, nil, fmt.Errorf("Room Already Exist")
			}
			f.Rooms[r.Name] = r

			if nextIsStart {
				if f.Start != nil {
					return raw, nil, fmt.Errorf("Duplicate Start")
				}
				f.Start = r
				nextIsStart = false
			} else if nextIsEnd {
				if f.End != nil {
					return raw, nil, fmt.Errorf("Duplicate End")
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
				return raw, nil, fmt.Errorf("Err in link")
			}
			ra := f.Rooms[a]
			rb := f.Rooms[b]
			if ra == nil || rb == nil {
				return raw, nil, fmt.Errorf("Room Not Found "+a+" or "+b)
			}
			if a == b {
				return raw, nil, fmt.Errorf("Link To Same Room "+a)
			}

			key := linkKey(a, b)
			if seenLink[key] {
				return raw, nil, fmt.Errorf("duplicate link "+key)
			}
			seenLink[key] = true

			ra.Links = append(ra.Links, rb)
			rb.Links = append(rb.Links, ra)
			continue
		}

		// unknown line format
		return raw, nil, fmt.Errorf("Invalid Format")
	}

	if scErr := sc.Err(); scErr != nil {
		return raw, nil, fmt.Errorf("Err While Scanning")
	}

	// final validations
	if !gotAnts || f.Start == nil || f.End == nil {
		return raw, nil, fmt.Errorf("Final Validation")
	}

	// commands without a following room should fail
	if nextIsStart || nextIsEnd {
		return raw, nil, fmt.Errorf("ErrInvalid")
	}

	return raw, f, nil
}
