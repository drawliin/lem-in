package utils

import (
	"container/list"
	"math"
	"sort"
)

const (
	maxCandidatePaths = 18
	pathDepthSlack    = 8
)

// FindDisjointPaths builds several simple start->end candidates, then picks
// the vertex-disjoint subset that minimizes total turns for all ants.
func FindDisjointPaths(f *Farm) [][]*Room {
	if f == nil || f.Start == nil || f.End == nil {
		return nil
	}

	dist := distanceToEnd(f)
	if dist[f.Start] < 0 {
		return nil
	}

	all := collectCandidatePaths(f, dist, maxCandidatePaths, pathDepthSlack)
	if len(all) == 0 {
		return nil
	}

	best := chooseBestDisjointSet(all, f.Ants)
	if len(best) == 0 {
		return nil
	}
	return best
}

// distanceToEnd computes shortest edge-distance from every room to end.
func distanceToEnd(f *Farm) map[*Room]int {
	dist := make(map[*Room]int, len(f.Rooms))
	for _, r := range f.Rooms {
		dist[r] = -1
	}

	q := list.New()
	dist[f.End] = 0
	q.PushBack(f.End)

	for q.Len() > 0 {
		node := q.Remove(q.Front()).(*Room)
		for _, nxt := range node.Links {
			if dist[nxt] != -1 {
				continue
			}
			dist[nxt] = dist[node] + 1
			q.PushBack(nxt)
		}
	}
	return dist
}

// collectCandidatePaths enumerates near-shortest simple paths and keeps the
// shortest candidates for subset evaluation.
func collectCandidatePaths(f *Farm, dist map[*Room]int, maxCount, slack int) [][]*Room {
	minLen := dist[f.Start]
	if minLen < 0 {
		return nil
	}

	visited := map[*Room]bool{f.Start: true}
	path := []*Room{f.Start}
	all := make([][]*Room, 0, maxCount)

	var dfs func(cur *Room)
	dfs = func(cur *Room) {
		if len(all) >= maxCount {
			return
		}
		if cur == f.End {
			cp := make([]*Room, len(path))
			copy(cp, path)
			all = append(all, cp)
			return
		}

		if dist[cur] < 0 {
			return
		}

		edgesSoFar := len(path) - 1
		bestPossible := edgesSoFar + dist[cur]
		if bestPossible > minLen+slack {
			return
		}

		nexts := append([]*Room(nil), cur.Links...)
		sort.Slice(nexts, func(i, j int) bool {
			di, dj := dist[nexts[i]], dist[nexts[j]]
			if di == dj {
				return nexts[i].Name < nexts[j].Name
			}
			if di < 0 {
				return false
			}
			if dj < 0 {
				return true
			}
			return di < dj
		})

		for _, nxt := range nexts {
			if visited[nxt] {
				continue
			}
			visited[nxt] = true
			path = append(path, nxt)
			dfs(nxt)
			path = path[:len(path)-1]
			visited[nxt] = false
			if len(all) >= maxCount {
				return
			}
		}
	}

	dfs(f.Start)

	sort.Slice(all, func(i, j int) bool {
		if len(all[i]) != len(all[j]) {
			return len(all[i]) < len(all[j])
		}
		return pathSignature(all[i]) < pathSignature(all[j])
	})
	return all
}

func chooseBestDisjointSet(all [][]*Room, ants int) [][]*Room {
	bestTurns := math.MaxInt
	bestLenSum := math.MaxInt
	var best [][]*Room

	used := make(map[string]bool)
	cur := make([][]*Room, 0, len(all))

	var bt func(i int)
	bt = func(i int) {
		if i == len(all) {
			if len(cur) == 0 {
				return
			}
			turns := minTurnsForSet(ants, cur)
			lenSum := sumPathEdges(cur)
			if turns < bestTurns || (turns == bestTurns && lenSum < bestLenSum) {
				bestTurns = turns
				bestLenSum = lenSum
				best = clonePaths(cur)
			}
			return
		}

		// Skip current path
		bt(i + 1)

		// Take current path if still disjoint by intermediate rooms
		p := all[i]
		if !canUsePath(p, used) {
			return
		}
		markPath(p, used, true)
		cur = append(cur, p)
		bt(i + 1)
		cur = cur[:len(cur)-1]
		markPath(p, used, false)
	}

	bt(0)
	return best
}

func canUsePath(p []*Room, used map[string]bool) bool {
	for i := 1; i < len(p)-1; i++ {
		if used[p[i].Name] {
			return false
		}
	}
	return true
}

func markPath(p []*Room, used map[string]bool, on bool) {
	for i := 1; i < len(p)-1; i++ {
		if on {
			used[p[i].Name] = true
		} else {
			delete(used, p[i].Name)
		}
	}
}

func minTurnsForSet(ants int, paths [][]*Room) int {
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
		mid := (lo + hi) / 2
		cap := 0
		for _, p := range paths {
			e := len(p) - 1
			if mid >= e {
				cap += mid - e + 1
			}
		}
		if cap >= ants {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

func sumPathEdges(paths [][]*Room) int {
	total := 0
	for _, p := range paths {
		total += len(p) - 1
	}
	return total
}

func clonePaths(src [][]*Room) [][]*Room {
	out := make([][]*Room, len(src))
	copy(out, src)
	return out
}

func pathSignature(p []*Room) string {
	s := ""
	for i, r := range p {
		if i > 0 {
			s += "->"
		}
		s += r.Name
	}
	return s
}
