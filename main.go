package main

import (
	"fmt"
	"os"
	"strings"

	"lem-in/utils"
)

// main validates the CLI input, parses the farm, finds usable paths,
// redistributes ants across those paths, then prints the simulation output.
func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR: invalid data format")
		return
	}

	fmt.Println("Starting Parsing Farm") // logger
	raw, farm, err := utils.ParseFarm(os.Args[1])
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return
	}

	fmt.Println("Finding Paths") // logger
	paths := utils.FindDisjointPaths(farm)
	if paths == nil {
		fmt.Println("ERROR: invalid data format")
		return
	}

	///////////////////
	fmt.Println("Founded paths:", len(paths))
	for i, p := range paths {
		names := make([]string, 0, len(p))
		for _, r := range p {
			names = append(names, r.Name)
		}
		fmt.Printf("P%d: %s (edges=%d)\n", i+1, strings.Join(names, " -> "), len(p)-1)
	}
	///////////////////
	fmt.Println("Start Assigning ants") // logger
	assign := utils.AssignAnts(farm.Ants, paths)

	// Print input back exactly
	fmt.Print(raw)
	if len(raw) > 0 && raw[len(raw)-1] != '\n' {
		fmt.Print("\n")
	}
	fmt.Println()

	fmt.Println("Starting Simulation") // logger
	utils.SimulateAndPrint(farm, paths, assign)
}
