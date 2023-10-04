package utils

import (
	"bytes"
	"log"
	"strconv"
	"strings"
)

// TODO :
// 	- A function that, given a colony and all valid paths in it, computes the best way to maximize the flow within the colony.

// NOT FULLY IMPLEMENTED YET.
// This function checks if input is a valid colony.
// Some examples of invalid input :
//   - A colony which contains at least one room linked to itself is considered invalid.
//   - A colony without ant is considered invalid.
//   - A colony which contains at least one duplicated room is invalid.
//   - A colony which does not contain any start and end room is invalid.
//
// Input : an ant colony, composed of rooms, tunnels and ants.
// Output : True if colony is valid, False otherwise.
func CheckColony(colony []Room) (isValid bool) {
	return true
}

// Evaluate all possible combinations based on overlap and length criteria, weights; and returns the best
func EvaluatePathsCombinations(paths []*Path, numPathsToSelectMax int, overlapWeight, lengthWeight float64) (bestPaths []*Path) {
	var bestCombination []*Path
	bestScore := -1.0

	combinations := generateCombinations(paths, numPathsToSelectMax)

	for _, combination := range combinations {
		totalUniqueNodes := make(map[*Room]bool)
		totalLength := 0
		for _, path := range combination {
			for idx, room := range path.Nodes {
				if idx != 0 && idx != len(path.Nodes)-1 {
					totalUniqueNodes[room] = true
				}
			}
			totalLength += len(path.Nodes[1 : len(path.Nodes)-1])
		}

		uniqueNodeRate := float64(len(totalUniqueNodes)) / float64(totalLength)
		lengthRate := float64(1) / float64(totalLength)

		score := (overlapWeight * uniqueNodeRate) + (lengthWeight * lengthRate)

		if score > bestScore {
			bestCombination = combination
			bestScore = score
		}
	}

	return bestCombination
}

// Generates all combinations of numPathsToSelect-uplets among paths (recursive)
func generateCombinations(paths []*Path, numPathsToSelect int) [][]*Path {
	if numPathsToSelect == 0 {
		return [][]*Path{{}}
	}

	if len(paths) == 0 {
		return nil
	}

	head, tail := paths[0], paths[1:]

	// Generate combinations by including/excluding first path
	combosIncluding := generateCombinations(tail, numPathsToSelect-1)
	for i := range combosIncluding {
		combosIncluding[i] = append([]*Path{head}, combosIncluding[i]...)
	}

	combosExcluding := generateCombinations(tail, numPathsToSelect)

	// Merge both combination sets
	combosIncluding = append(combosIncluding, combosExcluding...)

	return combosIncluding
}

// This function returns maximum number of unique paths in colony
func FindPaths_BFS(data []byte, s *Room, d *Room, howManyPathsMin int) (paths [][]string, paths_struct []*Path) {
	queue := [][]*Room{{s}}
	pipes, _ := GetPipes(data)
	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		node := path[len(path)-1]
		if node.Name == d.Name {
			var pathStr []string
			for _, room := range path {
				pathStr = append(pathStr, room.Name)
			}
			paths = append(paths, pathStr)
		} else {
			for idx := range pipes {
				if pipes[idx].From.Name == node.Name {
					if !visitedInPath(&pipes[idx].To, path) {
						newPath := make([]*Room, len(path))
						copy(newPath, path)
						newPath = append(newPath, &pipes[idx].To)
						queue = append(queue, newPath)
					}
				} else if pipes[idx].To.Name == node.Name {
					if !visitedInPath(&pipes[idx].From, path) {
						newPath := make([]*Room, len(path))
						copy(newPath, path)
						newPath = append(newPath, &pipes[idx].From)
						queue = append(queue, newPath)
					}
				}
			}
		}
	}

	paths_struct = buildPathsFromNames(paths, s)

	if len(paths) > 20 {
		firstHalf := paths_struct[:len(paths_struct)/2]
		paths_struct = EvaluatePathsCombinations(firstHalf, howManyPathsMin, 50, 50)
	} else {
		paths_struct = EvaluatePathsCombinations(paths_struct, howManyPathsMin, 50, 50)
	}

	return paths, paths_struct
}

// Returns an array containing all names of nodes in path
func buildPathsFromNames(n [][]string, s *Room) (p []*Path) {
	var paths []*Path
	visitedRooms := make(map[string]*Room)
	visited := make(map[string]bool)

	for _, names := range n {
		var path Path

		for i, name := range names {
			if visited[name] {
				path.Nodes = append(path.Nodes, visitedRooms[name])
			} else {
				r := &Room{Name: name}
				if i == 0 {
					r.IsStart = true
					r.Ants = s.Ants
				} else if i == len(names)-1 {
					r.IsEnd = true
				}

				path.Nodes = append(path.Nodes, r)
				visited[name] = true
				visitedRooms[name] = r
			}
		}

		paths = append(paths, &path)
	}

	return paths
}

// Returns a room's neighbors
func FindNeighbors(room *Room, pipes []Pipe) (res []*Room) {
	for i, pipe := range pipes {
		if pipe.To.Name == room.Name && !roomInRooms(&pipes[i].From, res) {
			res = append(res, &pipes[i].From)
		} else if pipe.From.Name == room.Name && !roomInRooms(&pipes[i].To, res) {
			res = append(res, &pipes[i].To)
		}
	}
	return res
}

// Checks for instance of Room in []Room
func roomInRooms(node *Room, path []*Room) bool {
	for _, n := range path {
		if n.Name == node.Name {
			return true
		}
	}
	return false
}

// Checks for instance of node in path
func visitedInPath(node *Room, path []*Room) bool {
	for _, n := range path {
		if n.Name == node.Name {
			return true
		}
	}
	return false
}

// This function checks for the start Room of colony.
func FindStart(data []byte) (start Room) {
	dataBis := ReadData(data)
	for i, v := range dataBis {
		if v == "##start" {
			start.Name = strings.Split(dataBis[i+1][:len(dataBis[i+1])], " ")[0]
			tmp, err := strconv.Atoi(dataBis[0])
			Check(err)
			start.Ants = tmp
			start.IsStart = true
		}
	}
	return start
}

// This function checks for the end Room of colony.
func FindEnd(data []byte) (end Room) {
	dataBis := ReadData(data)
	for i, v := range dataBis {
		if v == "##end" {
			end.Name = strings.Split(dataBis[i+1][:len(dataBis[i+1])], " ")[0]
			end.IsEnd = true
		}
	}
	return end
}

// A simple function that let us manipulate strings instead of bytes.
func ReadData(data []byte) (result []string) {
	var sep []byte
	sep = append(sep, 10)
	tmp := bytes.Split(data, sep)
	for _, v := range tmp {
		result = append(result, string(v))
	}
	return result
}

// This function looks at the colony and returns all pipes which compose it.
func GetPipes(data []byte) (pipes []Pipe, checked bool) {
	checked = true
	dataStr := ReadData(data)
	tmp := Pipe{}
	for _, v := range dataStr {
		if v[0] != '#' && strings.Contains(v, "-") {
			tmp.From.Name = strings.Split(v, "-")[0]
			tmp.To.Name = strings.Split(v, "-")[1]
			if tmp.From.Name == tmp.To.Name {
				checked = false
			}
			pipes = append(pipes, tmp)
		}
	}
	return pipes, checked
}

// This function looks at the colony and returns all rooms which compose it.
func GetRooms(data []byte) []*Room {
	dataStr := ReadData(data)
	rooms := make([]*Room, 0)

	for _, line := range dataStr {
		if strings.HasPrefix(line, "##") || strings.Contains(line, "-") {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) >= 3 {
			room := &Room{
				Name: parts[0],
			}
			rooms = append(rooms, room)
		}
	}

	return rooms
}

// Checks if there is any Start/End Room in the Colony
func IsThereAnyStartEnd(data []byte) bool {
	dataStr := ReadData(data)
	str := ""
	for _, v := range dataStr {
		str += v
	}
	if !(strings.Contains(str, "##end")) || !(strings.Contains(str, "##start")) {
		return false
	}
	return true
}

// This checks if there is a path available from start to end.
func IsThereAPath(paths []*Path) bool {
	if len(paths) != 0 {
		return true
	}
	return false
}

// This checks wether or not the Colony is valid
func Checker(data []byte, paths []*Path) bool {
	dataBis := ReadData(data)
	if nb, _ := strconv.Atoi(dataBis[0]); nb <= 0 {
		log.Fatal("ERROR: invalid data format, invalid number of Ants")
	}
	_, check := GetPipes(data)
	if !check {
		log.Fatal("ERROR: invalid data format, infinite loop detected (room linked to itself)")
	}
	startEndCheck := IsThereAnyStartEnd(data)
	if startEndCheck == false {
		log.Fatal("ERROR: invalid data format, no start or end room found")
	}
	if !IsThereAPath(paths) {
		log.Fatal("ERROR: invalid data format, no path available to reach end from start")
	}
	return true
}
