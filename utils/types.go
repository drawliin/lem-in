package utils

type Room struct {
	Name  string
	X, Y  int
	Links []*Room
}

type Farm struct {
	Ants  int
	Rooms map[string]*Room
	Start *Room
	End   *Room
}

type AntState struct {
	ID      int
	PathIdx int
	Pos     int // index inside path slice
	Done    bool
}

type PathState struct {
	path    []*Room
	visited map[*Room]bool
}