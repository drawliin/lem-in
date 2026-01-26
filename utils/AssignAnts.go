package utils

// AssignAnts returns slice where index is antID (1..ants), value is path index.
func AssignAnts(ants int, paths [][]*Room) []int {
	assign := make([]int, ants+1) // ignore 0
	used := make([]int, len(paths))

	for ant := 1; ant <= ants; ant++ {
		best := 0
		bestScore := score(paths[0], used[0])

		for i := 1; i < len(paths); i++ {
			s := score(paths[i], used[i])
			if s < bestScore {
				bestScore = s
				best = i
			}
		}

		assign[ant] = best
		used[best]++
	}
	return assign
}

func score(p []*Room, assigned int) int {
	// smaller is better
	// p length includes start+end, so edges = len(p)-1
	return len(p) + assigned
}
