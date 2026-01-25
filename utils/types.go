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