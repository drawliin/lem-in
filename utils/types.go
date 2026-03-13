package utils

// Room stores one node in the graph along with its coordinates and neighbors.
type Room struct {
	Name  string
	X, Y  int
	Links []*Room
}

// Farm is the parsed input: ant count, all rooms, and the designated endpoints.
type Farm struct {
	Ants  int
	Rooms map[string]*Room
	Start *Room
	End   *Room
}

// AntState tracks one ant during the turn-by-turn simulation.
type AntState struct {
	ID      int
	PathIdx int
	Pos     int // index inside path slice
	Done    bool
}