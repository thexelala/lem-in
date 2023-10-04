package utils

type Room struct {
	Name           string
	IsStart, IsEnd bool
	Ant            [][]AntByPath
	Ants, Count    int
	Neighbors      []*Room
}

type Pipe struct {
	From Room
	To   Room
}

type Ant struct {
	Name     int
	Location *Room
	Path     *Path
}

type AntByPath struct {
	Name     int
	Location *Room
}

type Path struct {
	Name  string
	Nodes []*Room
}
