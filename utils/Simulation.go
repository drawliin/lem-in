package utils

import (
	"fmt"
	"sort"
	"strings"
)

// SimulateAndPrint advances ants one step per turn, respecting room occupancy,
// and prints each turn in the required `L<id>-<room>` format.
func SimulateAndPrint(f *Farm, paths [][]*Room, assign []int) {
	ants := make([]*AntState, 0, f.Ants)
	for id := 1; id <= f.Ants; id++ {
		ants = append(ants, &AntState{
			ID:      id,
			PathIdx: assign[id],
			Pos:     0,
		})
	}

	// occupancy for rooms except start/end
	occupied := make(map[string]int)

	doneCount := 0
	for doneCount < f.Ants {
		moves := make([]string, 0)

		// Move ants that are further ahead first
		sort.Slice(ants, func(i, j int) bool {
			// higher position first
			if ants[i].Pos != ants[j].Pos {
				return ants[i].Pos > ants[j].Pos
			}
			// tie: smaller ID first (stable-ish)
			return ants[i].ID < ants[j].ID
		})

		for _, a := range ants {
			if a.Done {
				continue
			}

			path := paths[a.PathIdx]
			if a.Pos >= len(path)-1 {
				// already at end
				a.Done = true
				continue
			}

			nextPos := a.Pos + 1
			nextRoom := path[nextPos]

			// Can always enter end
			if nextRoom != f.End {
				if _, ok := occupied[nextRoom.Name]; ok {
					continue // room occupied
				}
			}

			// leaving current room: free it if not start/end
			curRoom := path[a.Pos]
			if curRoom != f.Start && curRoom != f.End {
				delete(occupied, curRoom.Name)
			}

			// entering next room: occupy if not start/end
			if nextRoom != f.Start && nextRoom != f.End {
				occupied[nextRoom.Name] = a.ID
			}

			a.Pos = nextPos
			if nextRoom == f.End {
				a.Done = true
				doneCount++
			}

			moves = append(moves, fmt.Sprintf("L%d-%s", a.ID, nextRoom.Name))
		}

		if len(moves) == 0 {
			// No moves possible -> should not happen with valid disjoint paths,
			// but avoid infinite loop.
			break
		}

		fmt.Println(strings.Join(moves, " "))
	}
}
