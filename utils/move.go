package utils

import (
	"fmt"
	"math"
	"reflect"
)

// A simple check.
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// Move all ants towards destination and returns time (turns)
func MoveAll(source *Room, dest *Room, paths []*Path) int {
	var antsTab []*Ant
	var ok, count int

	for i := 0; i < source.Ants; i++ {
		antsTab = append(antsTab, &Ant{Name: i + 1, Location: paths[0].Nodes[0]})
	}

	antsRepartition := BalanceAntsRepartition(paths)
	AssignPaths(antsTab, antsRepartition)

	for ok != len(antsTab) {
		for i := range antsTab {
			if antsTab[i].Path != nil {
				if LegalMove(antsTab[i]) {
					if antsTab[i].Location.Name != dest.Name {
						MoveAnt(antsTab[i])
						if reflect.DeepEqual(antsTab[i].Location, paths[0].Nodes[len(paths[0].Nodes)-1]) {
							ok++
						}
					}
				}
			}
		}
		DisplayRooms(paths)
		count++
		if count == 51 {
			break
		}
	}

	return count
}

// Assign paths to ants, based on antsRepartition
func AssignPaths(ants []*Ant, repartition []int) {
}

// Balance ants repartition for each of all paths
func BalanceAntsRepartition(paths []*Path) []int {
	equations := make([][]Rational, len(paths))
	repartition := make([]float64, len(paths))

	for idx := range paths {
		if idx == 0 {
			equations[idx] = InitEq(paths)
		} else {
			equations[idx] = OtherEq(len(paths), idx, Rational{(int64(len(paths[idx].Nodes) - len(paths[0].Nodes))), 1})
		}
	}

	res, gausErr := SolveGaussian(equations, false)
	Check(gausErr)

	for i := range res {
		for j := range res[i] {
			tmp := float64(float64(res[i][j].numerator) / float64(res[i][j].denominator))
			repartition[i] = tmp
		}
	}
	finalRepartition := RepartLastAnts(repartition)
	return finalRepartition
}

// Balance repartition tweak in case we do not get integers only
func RepartLastAnts(repartition []float64) []int {
	var roundedRep []int
	for _, v := range repartition {
		roundedRep = append(roundedRep, int(v))
	}

	var difIntFloatRep []float64
	for k, v := range repartition {
		difIntFloatRep = append(difIntFloatRep, v-float64(roundedRep[k]))
	}

	var antsLeft float64
	for _, v := range difIntFloatRep {
		antsLeft += v
	}

	intAntLeft := math.Round(antsLeft)

	if intAntLeft == 1 {
		roundedRep[0]++
	} else if intAntLeft == 2 {
		roundedRep[0]++
		roundedRep[1]++
	}
	return roundedRep
}

// Set equation : set coeff of the eq[coeffIdx] to -1
func OtherEq(numOfPaths, coeffIdx int, rightTerm Rational) []Rational {
	var found bool

	res := make([]Rational, numOfPaths+1)
	res[0] = Rational{1, 1}
	res[numOfPaths] = rightTerm

	for i := 1; i < numOfPaths; i++ {
		if i%coeffIdx == 0 {
			res[i] = Rational{-1, 1}
			found = true
		}
		if found {
			break
		}
	}
	for i := range res {
		res[i].denominator = 1
	}
	return res
}

// Initializes the first equation of our linear system : x0 + ... + xn = total ants
func InitEq(paths []*Path) []Rational {
	eq := make([]Rational, len(paths)+1)
	for i := range paths {
		eq[i] = Rational{1, 1}
	}
	eq[len(eq)-1] = Rational{int64(paths[0].Nodes[0].Ants), 1}
	return eq
}

// Returns the shortest available path, nil if none is found
// func shortestAvailable(paths []*Path) *Path {
// 	dest := paths[0].Nodes[len(paths[0].Nodes)-1]
// 	for i := range paths {
// 		if paths[i].Nodes[1].Ants < 1 || reflect.DeepEqual(paths[i].Nodes[1].Name, dest.Name) {
// 			return paths[i]
// 		}
// 	}
// 	return nil
// }

// Returns true if ant has legal move
func LegalMove(a *Ant) bool {
	for i := range a.Path.Nodes[:len(a.Path.Nodes)-1] {
		if reflect.DeepEqual(a.Path.Nodes[i], a.Location) {
			if a.Path.Nodes[i+1].Ants > 0 && a.Path.Nodes[i+1].Name != a.Path.Nodes[len(a.Path.Nodes)-1].Name {
				return false
			} else {
				return true
			}
		}
	}
	return false
}

// Displays on standard output the locations of every ant
func DisplayAntsLocations(antsTab []*Ant) {
	for _, ant := range antsTab {
		fmt.Printf("Ant %d is in room %s\n", ant.Name, ant.Location.Name)
	}
}

// Displays on standard output the number of ants in every room
func DisplayRooms(paths []*Path) {
	antCounts := make(map[string]int)

	for _, path := range paths {
		for _, room := range path.Nodes {
			antCounts[room.Name] = room.Ants
		}
	}

	fmt.Println("Ant counts in each room:")
	for roomName, count := range antCounts {
		fmt.Printf("Room %s: %d ants\n", roomName, count)
	}
}

// Moves an Ant along its dedicated path
func MoveAnt(a *Ant) {
	dest := a.Path.Nodes[len(a.Path.Nodes)-1]
	for i := range a.Path.Nodes {
		if reflect.DeepEqual(a.Path.Nodes[i], a.Location) && a.Location.Name != dest.Name {
			a.Path.Nodes[i].Ants--
			a.Location = a.Path.Nodes[i+1]
			a.Path.Nodes[i+1].Ants++
			break
		}
	}
}
