# Lem-in

**Lem-in** is a Go implementation of the classic ant farm pathfinding project. It parses a farm description, validates rooms and links, finds efficient non-overlapping paths from `##start` to `##end`, assigns ants across those paths, and prints the turn-by-turn simulation.

This repository is useful for anyone looking for a **Go lem-in solver**, **ant colony pathfinding project**, **graph traversal example in Go**, or a **42/1337 lem-in style implementation**.

## Small Description

This project reads an ant farm input file, builds a graph of rooms and tunnels, computes valid paths between the start and end rooms, distributes ants to reduce total turns, and outputs the movements in the expected `Lx-room` format.

## Topics

Suggested topics for GitHub SEO and discovery:

`go`, `golang`, `lem-in`, `pathfinding`, `graph-algorithm`, `bfs`, `ant-colony`, `simulation`, `42-school`, `1337-school`, `algorithm`, `cli`

## Features

- Parses farm files from disk
- Validates ant count, rooms, links, `##start`, and `##end`
- Rejects malformed input and duplicate room/link definitions
- Computes candidate paths through the farm graph
- Chooses compatible paths that do not overlap on intermediate rooms
- Assigns ants to the most efficient paths
- Simulates movements turn by turn while respecting room and tunnel constraints
- Prints the original input followed by the generated solution

## Project Structure

- `main.go`: CLI entry point
- `utils/ParseFarm.go`: input parsing and validation
- `utils/FindPaths.go`: path discovery and best path-set selection
- `utils/AssignAnts.go`: ant-to-path distribution
- `utils/Simulation.go`: movement simulation and output printing
- `utils/types.go`: core data structures

## Requirements

- Go `1.23.2` or later

## Run

```bash
go run . path/to/farm.txt
```

You can also build the binary first:

```bash
go build -o lem-in .
./lem-in path/to/farm.txt
```

## Input Format

The program expects a farm description file with:

- The number of ants on the first meaningful line
- Room definitions in the format: `name x y`
- `##start` followed by the start room
- `##end` followed by the end room
- Links in the format: `room1-room2`
- Optional comments beginning with `#`

Example:

```txt
4
##start
start 0 0
room1 1 0
room2 2 0
##end
end 3 0
start-room1
room1-room2
room2-end
```

## Output

The program prints:

1. The original input
2. A blank line
3. The ant movements turn by turn

Example movement line:

```txt
L1-room1 L2-room1
```

## How It Works

1. The parser reads and validates the farm file.
2. The graph is built from rooms and bidirectional links.
3. The solver estimates distances from the end room using breadth-first search.
4. Candidate paths are collected and ranked.
5. A compatible set of non-conflicting paths is selected.
6. Ants are distributed across paths using a simple path-length-plus-load score.
7. The simulation prints legal moves until all ants reach the end.

## Error Handling

The program reports errors such as:

- unreadable input files
- invalid ant counts
- missing start or end rooms
- duplicate rooms or links
- malformed room or link lines
- links that reference unknown rooms
- invalid overall farm structure

## SEO Description

A Go `lem-in` project that parses ant farm maps, validates graph input, finds efficient non-overlapping paths, assigns ants intelligently, and simulates movement from start to end through tunnels.

## License

Add your preferred license here if you plan to publish or share the project publicly.
