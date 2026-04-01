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

func FindPaths(f *Farm) [][]*Room {
	if f == nil || f.Start == nil || f.End == nil {
		return nil
	}

	dist := distanceToEnd(f)
	// if start cannot reach end.. return
	if dist[f.Start] < 0 {
		return nil
	}

	all := collectAllPaths(f, dist, maxCandidatePaths, pathDepthSlack)
	if len(all) == 0 {
		return nil
	}

	best := chooseBestPaths(all, f.Ants)
	if len(best) == 0 {
		return nil
	}
	return best
}

// computes how far each room is from end
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

func collectAllPaths(f *Farm, dist map[*Room]int, maxCount, slack int) [][]*Room {
	minLen := dist[f.Start]
	if minLen < 0 {
		return nil
	}

	all := make([][]*Room, 0, maxCount)

	// open means paths we still want to explore.
	open := []PathState{{
		path:    []*Room{f.Start},
		visited: map[*Room]bool{f.Start: true},
	}}

	for len(open) > 0 && len(all) < maxCount {
		sort.Slice(open, func(i, j int) bool {
			si, sj := open[i], open[j]
			ei := EstimatePath(si.path, dist)
			ej := EstimatePath(sj.path, dist)
			if ei != ej {
				return ei < ej
			}
			if len(si.path) != len(sj.path) {
				return len(si.path) < len(sj.path)
			}
			return PathSignature(si.path) < PathSignature(sj.path)
		})

		curState := open[0]
		open = open[1:]
		cur := curState.path[len(curState.path)-1]

		if cur == f.End {
			all = append(all, curState.path)
			continue
		}
		if dist[cur] < 0 {
			continue
		}

		edgesSoFar := len(curState.path) - 1
		bestPossible := edgesSoFar + dist[cur]
		if bestPossible > minLen+slack {
			continue
		}

		nexts := append([]*Room(nil), cur.Links...)
		sort.Slice(nexts, func(i, j int) bool {
			di, dj := dist[nexts[i]], dist[nexts[j]]
			if di < 0 {
				return false
			}
			if dj < 0 {
				return true
			}
			if di == dj {
				return nexts[i].Name < nexts[j].Name
			}
			return di < dj
		})

		for _, nxt := range nexts {
			if curState.visited[nxt] {
				continue
			}

			nextPath := make([]*Room, len(curState.path)+1)
			copy(nextPath, curState.path)
			nextPath[len(curState.path)] = nxt

			nextVisited := CloneVisited(curState.visited)
			nextVisited[nxt] = true

			open = append(open, PathState{
				path:    nextPath,
				visited: nextVisited,
			})
		}

		if len(open) > maxCount*maxCount {
			sort.Slice(open, func(i, j int) bool {
				si, sj := open[i], open[j]
				ei := EstimatePath(si.path, dist)
				ej := EstimatePath(sj.path, dist)
				if ei != ej {
					return ei < ej
				}
				if len(si.path) != len(sj.path) {
					return len(si.path) < len(sj.path)
				}
				return PathSignature(si.path) < PathSignature(sj.path)
			})
			open = open[:maxCount*maxCount]
		}
	}

	sort.Slice(all, func(i, j int) bool {
		if len(all[i]) != len(all[j]) {
			return len(all[i]) < len(all[j])
		}
		return PathSignature(all[i]) < PathSignature(all[j])
	})
	return all
}

func chooseBestPaths(all [][]*Room, ants int) [][]*Room {
	bestTurns := math.MaxInt
	bestLenSum := math.MaxInt
	var best [][]*Room

	used := make(map[string]bool)
	cur := make([][]*Room, 0, len(all))

	var helper func(i int)
	helper = func(i int) {
		if i == len(all) {
			if len(cur) == 0 {
				return
			}
			turns := MinTurnsForSet(ants, cur)
			lenSum := SumPathEdges(cur)
			if turns < bestTurns || (turns == bestTurns && lenSum < bestLenSum) {
				bestTurns = turns
				bestLenSum = lenSum
				best = ClonePaths(cur)
			}
			return
		}

		// Skip current path
		helper(i + 1)

		// Take current path if still disjoint by intermediate rooms
		p := all[i]
		if !CanUsePath(p, used) {
			return
		}
		MarkPath(p, used, true)
		cur = append(cur, p)
		helper(i + 1)
		cur = cur[:len(cur)-1]
		MarkPath(p, used, false)
	}

	helper(0)
	return best
}
