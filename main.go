package main

import (
	"fmt"
	"os"

	"lem-in/utils"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("ERROR: no given file")
		return
	}

	raw, farm, err := utils.ParseFarm(os.Args[1])
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return
	}

	paths := utils.FindPaths(farm)
	if paths == nil {
		fmt.Println("ERROR: invalid data format")
		return
	}

	assign := utils.AssignAnts(farm.Ants, paths)

	// Print input back exactly
	fmt.Print(raw)
	if len(raw) > 0 && raw[len(raw)-1] != '\n' {
		fmt.Print("\n")
	}
	fmt.Println()

	utils.SimulateAndPrint(farm, paths, assign)
}
