package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"lem-in/utils"
)

func main() {

	// === DEBUG ===
	// dat, err := os.ReadFile("../colonies/example05.txt")
	// utils.Check(err)

	arg := os.Args[1]

	dat, err := os.ReadFile("./colonies/" + arg + ".txt")
	utils.Check(err)

	source := utils.FindStart(dat)
	sink := utils.FindEnd(dat)
	pipes, _ := utils.GetPipes(dat)
	howMany := len(utils.FindNeighbors(&source, pipes))

	_, Paths := utils.FindPaths_BFS(dat, &source, &sink, howMany)
	utils.Checker(dat, Paths)

	fmt.Println(string(dat))
	fmt.Println()

	sort.Slice(Paths, comparePaths(Paths))

	var filteredPaths []*utils.Path

	if source.Ants < 5 {
		Paths = append(filteredPaths, Paths[0])
	} else if source.Ants < 10 && len(pipes) < 15 {
		Paths = append(filteredPaths, Paths[0], Paths[1])
	}

	// 	scores = append(scores, utils.MoveAll(&source, &sink, paths))
	// 	if len(paths) > 1 {
	// 		howMany--
	// 	} else {
	// 		break
	// 	}

	// }

	Repartition := utils.BalanceAntsRepartition(Paths)
	tab := []float64{}
	for _, v := range Repartition {
		tab = append(tab, float64(v))
	}

	// for i := 0; i < len(Paths); i++ {
	// 	for k := range Paths[i].Nodes {
	// 		fmt.Println(Paths[i].Nodes[k])

	// 	}
	// }

	//TODO:
	// Fair un lien entre Les Paths et la répartition example [4 3 6] 4 a assigner au premier path, 3 pour le deuxiemme, 6 pour le troisiemme si il y a trois paths
	// Faire un boucle infinie que l'on breack quand toutes les fourmil sont passer
	// dans la boucle faire une autre boucle par rapport a la len de path pour gere l'avancement des fourmil dans leur chemain
	// si la salle juste après la salle de départ et prise après déplacement ne plus bouger et passer au prochain paths jusqu'a la fin du tour
	// au tour suivant faire avancer la fourmil dans dans sa voie et la remplace par une autre
	// stoper d'utiliser cette path quand dans Repartition il n'y a plus de fourmil assigner se path

	// ON FAIT A L'ENVERT POUR FAIRE AVANCER LES FOURMIL ET PAS LES EMPACTER
	// REMARQUE ON PEUT LE FAIRE A L'ENDROIT COMME SA SI ELLE SONT DEUX SONT DEGAGE LA PLUS ANCIENNT QUI ES ARRIVER DES DEUX

	// for i, path := range Paths {
	// 	fmt.Print("Path ", i+1, ": ")
	// 	for _, room := range path.Nodes {
	// 		fmt.Print(room.Name + " ")
	// 	}
	// 	fmt.Println()
	// }

	//TODO: - de 5 un chemain et - de 10 mais plus de 5 juste deux

	var trueAntsRep = AntsByPaths(Repartition, &source)

	var toto bool

	chuiPassey := false

	var turn string

	for toto == false {

		// fmt.Println("/////////////////////////IDEE DE BASSE A GARDER /////////////////////////")
		// for u, path := range Paths {
		// 	fmt.Print("Path ", u+1, ": ")
		// 	for _, room := range path.Nodes[1:] {
		// 		fmt.Printf(room.Name + " ")
		// 	}
		// 	fmt.Println()
		// }
		// fmt.Println("/////////////////////////IDEE DE BASSE A GARDER /////////////////////////")

		for i := len(Paths) - 1; i >= 0; i-- {
			//FIXME: fmt.Println("Path ", i, ": ")
			for k := len(Paths[i].Nodes) - 1; k >= 1; k-- {
				if Paths[i].Nodes[k-1].IsStart == true && !chuiPassey {
					for _, v := range trueAntsRep {
						Paths[i].Nodes[k-1].Ant = append(Paths[i].Nodes[k-1].Ant, v)

					}
					chuiPassey = true

				}

				//FIXME: example trois ligne 2 du result le L1-3 se repete alor qu'il ne devrait pas
				//FIXME: faire en sorte que en fonction du nombre de ants l'on utilise seulement 1 chemain ou les 2 plus cours

				if Paths[i].Nodes[k-1].Ant != nil {

					Paths[i].Nodes[k].Ant = make([][]utils.AntByPath, 1)

					Paths[i].Nodes[k].Ant[0] = make([]utils.AntByPath, 1)

					if len(Paths[i].Nodes[k-1].Ant) > 1 {
						if len(Paths[i].Nodes[k-1].Ant[i]) == 0 {
						} else {
							Paths[i].Nodes[k].Ant[0][0] = Paths[i].Nodes[k-1].Ant[i][0]
						}
					} else {
						if len(Paths[i].Nodes[k-1].Ant[0]) != 0 {
							Paths[i].Nodes[k].Ant[0][0] = Paths[i].Nodes[k-1].Ant[0][0]
						}
					}

					Paths[i].Nodes[k].Ant[0][0] = utils.AntByPath{Name: Paths[i].Nodes[k].Ant[0][0].Name, Location: Paths[i].Nodes[k]}

					if len(Paths[i].Nodes[k-1].Ant) != 1 {
						if len(Paths[i].Nodes[k-1].Ant[i]) == 0 {
						} else {
							Paths[i].Nodes[k-1].Ant[i] = append(Paths[i].Nodes[k-1].Ant[i][:0], Paths[i].Nodes[k-1].Ant[i][1:]...)
						}
					} else {
						if len(Paths[i].Nodes[k-1].Ant[0]) != 0 {

							Paths[i].Nodes[k-1].Ant[0] = append(Paths[i].Nodes[k-1].Ant[0][:0], Paths[i].Nodes[k-1].Ant[0][1:]...)
						}

					}

				}

				// if len(Paths[i].Nodes[k].Ant) > 0 && len(Paths[i].Nodes[k].Ant[0]) == 0 {
				// 	// Le tableau est vide, nous allons le supprimer
				// 	Paths[i].Nodes[k].Ant = append(Paths[i].Nodes[k].Ant[:0], Paths[i].Nodes[k].Ant[1:]...)
				// }

				if Paths[i].Nodes[k].Ant != nil {
					if Paths[i].Nodes[k].Ant[0][0].Name != 0 {
						// fmt.Printf("L%v-%s ", Paths[i].Nodes[k].Ant[0][0].Name, Paths[i].Nodes[k].Ant[0][0].Location.Name)
						if strings.Contains(turn, "L"+strconv.Itoa(Paths[i].Nodes[k].Ant[0][0].Name)+"-") {
							turn += "L" + strconv.Itoa(Paths[i].Nodes[k].Ant[0][0].Name) + "-" + Paths[i].Nodes[k].Ant[0][0].Location.Name + " "
							turn += "\n"

						} else {
							turn += "L" + strconv.Itoa(Paths[i].Nodes[k].Ant[0][0].Name) + "-" + Paths[i].Nodes[k].Ant[0][0].Location.Name + " "

						}

					}
				}
				if Paths[i].Nodes[k].IsEnd {
					Paths[i].Nodes[k].Ant = nil
				}

			}

		}
		//FIXME: fmt.Println("/////////////////////////DO A BARREL ROLL/////////////////////////")

		if turn == "" {
			toto = true
		}

		fmt.Println(turn)

		turn = ""
	}

}

func AntsByPaths(Rep []int, firstRoom *utils.Room) [][]utils.AntByPath {
	result := make([][]utils.AntByPath, len(Rep)) // Crée une slice de slices

	count := 1

	for i := 0; i < len(result); i++ {
		for k := 1; k < Rep[i]+1; k++ {
			Ants := utils.AntByPath{Name: count, Location: firstRoom}
			result[i] = append(result[i], Ants)
			count++
		}
	}

	return result
}

func comparePaths(paths []*utils.Path) func(int, int) bool {
	return func(i, j int) bool {
		return len(paths[i].Nodes) < len(paths[j].Nodes)
	}
}
