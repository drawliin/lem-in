package utils

import (
	"maps"
	"math"
	"strings"
)

func EstimatePath(path []*Room, dist map[*Room]int) int {
	cur := path[len(path)-1]
	if dist[cur] < 0 {
		return math.MaxInt
	}
	return len(path) - 1 + dist[cur]
}

func CloneVisited(src map[*Room]bool) map[*Room]bool {
	dst := make(map[*Room]bool, len(src)+1)
	maps.Copy(dst, src)
	return dst
}

func CanUsePath(p []*Room, used map[string]bool) bool {
	for i := 1; i < len(p)-1; i++ {
		if used[p[i].Name] {
			return false
		}
	}
	return true
}

func MarkPath(p []*Room, used map[string]bool, on bool) {
	for i := 1; i < len(p)-1; i++ {
		if on {
			used[p[i].Name] = true
		} else {
			delete(used, p[i].Name)
		}
	}
}

func MinTurnsForSet(ants int, paths [][]*Room) int {
	if ants <= 0 || len(paths) == 0 {
		return 0
	}

	minL := math.MaxInt
	maxL := 0
	for _, p := range paths {
		e := len(p) - 1
		if e < minL {
			minL = e
		}
		if e > maxL {
			maxL = e
		}
	}

	lo, hi := minL, maxL+ants
	for lo < hi {
		midT := (lo + hi) / 2
		cap := 0
		for _, p := range paths {
			e := len(p) - 1
			if midT >= e {
				cap += midT - e + 1  // ants = Turns - edges + 1
			}
		}
		if cap >= ants {
			hi = midT
		} else {
			lo = midT + 1
		}
	}
	return lo
}

func SumPathEdges(paths [][]*Room) int {
	total := 0
	for _, p := range paths {
		total += len(p) - 1
	}
	return total
}

func ClonePaths(src [][]*Room) [][]*Room {
	out := make([][]*Room, len(src))
	copy(out, src)
	return out
}

func PathSignature(p []*Room) string {
	var s strings.Builder
	for i, r := range p {
		if i > 0 {
			s.WriteString("->")
		}
		s.WriteString(r.Name)
	}
	return s.String()
}
