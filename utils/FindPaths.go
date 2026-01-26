package utils

import (
	"container/list"
	"fmt"
)

// BFSShortestPath returns the shortest path from start to end (inclusive).
// If no path exists, it returns nil.
func BFSShortestPath(f *Farm, blocked map[string]bool) []*Room {
	start := f.Start
	end := f.End
	if start == nil || end == nil {
		return nil
	}

	visited := make(map[*Room]bool, len(f.Rooms))
	parent := make(map[*Room]*Room, len(f.Rooms))

	q := list.New()
	q.PushBack(start)
	visited[start] = true

	found := false

	for q.Len() > 0 && !found {
		front := q.Front()
		q.Remove(front)
		cur := front.Value.(*Room)

		for _, nxt := range cur.Links {
			if visited[nxt] {
				continue
			}

			if nxt != start && nxt != end && blocked[nxt.Name] {
				continue
			}

			visited[nxt] = true
			parent[nxt] = cur

			if nxt == end {
				found = true
				break
			}
			q.PushBack(nxt)
		}
	}

	if !visited[end] {
		return nil
	}

	// Reconstruct path from end -> start
	var rev []*Room
	for node := end; node != nil; node = parent[node] {
		rev = append(rev, node)
		if node == start {
			break
		}
	}

	// reverse to get start -> end
	for i, j := 0, len(rev)-1; i < j; i, j = i+1, j-1 {
		rev[i], rev[j] = rev[j], rev[i]
	}
	return rev
}

// FindDisjointPaths repeatedly finds shortest paths while blocking intermediate rooms
// from previously chosen paths. Returns paths including start and end.
func FindDisjointPaths(f *Farm) [][]*Room {
	blocked := make(map[string]bool)
	var paths [][]*Room

	for {
		p := BFSShortestPath(f, blocked)
		if p == nil {
			fmt.Println("Breaked")
			break
		}

		paths = append(paths, p)

		// block intermediate rooms to force vertex-disjoint paths
		for i := 1; i < len(p)-1; i++ {
			blocked[p[i].Name] = true
		}
	}

	return paths
}
