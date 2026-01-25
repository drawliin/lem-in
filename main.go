package main

import (
	"fmt"
	"lem-in/utils"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR: invalid data format")
		return
	}

	_, farm, err := utils.ParseFarm(os.Args[1])
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return
	}

	path := utils.BFSShortestPath(farm)
	if path == nil {
		fmt.Println("ERROR: invalid data format")
		return
	}

	// Print path as names (debug)
	names := make([]string, 0, len(path))
	for _, r := range path {
		names = append(names, r.Name)
	}
	fmt.Println("Shortest path:", strings.Join(names, " -> "))
	fmt.Println("Edges:", len(path)-1)

}
