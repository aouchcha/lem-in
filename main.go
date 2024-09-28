package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Graph struct {
	vertices        []*Vertex
	Number_of_rooms int
	Ants            int
}

type Vertex struct {
	vesited  int
	Etat     string
	value    string
	adjecent []*Vertex
}

type DFS struct {
	Paths        [][]string
	Unique_Paths [][][]string
}

func (g *Graph) AddVertix(v string, stat string) {
	g.Number_of_rooms++
	for _, vertices := range g.vertices {
		if vertices.value == v {
			fmt.Fprintf(os.Stderr, "You have already a vertix with the value %s\n", v)
			return
		}
	}
	g.vertices = append(g.vertices, &Vertex{value: v, Etat: stat})
}

func (g *Graph) GetVertex(value string) *Vertex {
	for _, vertice := range g.vertices {
		if vertice.value == value || vertice.Etat == value {
			return vertice
		}
	}
	return nil
}

func (g *Graph) AddEdges(from, to string) {
	if from == to {
		return
	}
	fromv := g.GetVertex(from)
	tov := g.GetVertex(to)

	if fromv == nil || tov == nil {
		if fromv == nil {
			fmt.Fprintf(os.Stderr, "This vertex %s doesn't exist", from)
			os.Exit(1)
		} else {
			fmt.Fprintf(os.Stderr, "This vertex %s doesn't exist", to)
			os.Exit(1)
		}
	}

	for _, jiran := range fromv.adjecent {
		if jiran.value == to {
			fmt.Fprintf(os.Stderr, "There is already a relation between %s and %s\n", from, to)
			return
		}
	}

	fromv.adjecent = append(fromv.adjecent, tov)
	tov.adjecent = append(tov.adjecent, fromv)
}

func (g *Graph) Print() {
	// for _, vertice := range g.vertices {
	// 	fmt.Fprintf(os.Stdout, "you add a vertex with value %s to the graph\n", vertice.value)
	// }
	for _, vertice := range g.vertices {
		fmt.Print(vertice.value, "(", vertice.Etat, ")")
		for i := range vertice.adjecent {
			fmt.Print("->", vertice.adjecent[i].value, "(", vertice.adjecent[i].Etat, ")")
		}
		fmt.Println()
	}
}

func main() {
	gr := &Graph{}
	paths := &DFS{}
	CreatRoomsAndPaths(gr, TraitData())
	// gr.Print()
	var path []string
	// fmt.Println(gr.Number_of_rooms)
	SearchInTheGraph(gr.GetVertex("start"), paths, path)
	paths.SortPaths()
	ChooseUniquePaths(paths)
	fmt.Println("PATHS :")
	for i, p := range paths.Paths {
		fmt.Println(i, p)
	}
	fmt.Println("-----------------------------------------------------------------------")
	fmt.Println("UNIQUE PATHS :")
	for i, p := range paths.Unique_Paths {
		fmt.Println(i, p)
	}
}

func TraitData() []string {
	content, err := os.ReadFile("test.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "We can't read the file")
		os.Exit(1)
	}
	lines := strings.Split(string(content), "\n")
	return lines
}

func CreatRoomsAndPaths(gr *Graph, lines []string) {
	NmAnts, err := strconv.Atoi(lines[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalide ants number")
		os.Exit(1)
	}
	gr.Ants = NmAnts
	start := 0
	end := 0
	var roomes []string
	var links []string
	for i := 1; i < len(lines); i++ {
		if strings.Contains(lines[i], " ") || strings.HasPrefix(lines[i], "##") {
			roomes = append(roomes, lines[i])
		} else if strings.Contains(lines[i], "-") {
			links = append(links, lines[i])
		}
	}

	for i := 0; i < len(roomes); i++ {
		if strings.HasPrefix(roomes[i], "#") && (roomes[i] != "##start" && roomes[i] != "##end") {
			continue
		} else if roomes[i] == "##start" || roomes[i] == "##end" {
			sli := Check_Coord(roomes[i+1], i+1)
			if roomes[i] == "##start" {
				start++
				gr.AddVertix(sli[0], "start")
			} else {
				end++
				gr.AddVertix(sli[0], "end")
			}
			i++
		} else {
			sli := Check_Coord(roomes[i], i)
			gr.AddVertix(sli[0], "standard")
		}
	}
	if start != 1 || end != 1 {
		fmt.Fprintln(os.Stderr, "You give more than one start or end check your file !!!")
		os.Exit(1)
	}

	for i := 0; i < len(links); i++ {

		sli := strings.Split(links[i], "-")
		if len(sli) != 2 {
			fmt.Fprintf(os.Stderr, "Bad Data on links in line %s\n", links[i])
			continue
		}
		gr.AddEdges(sli[0], sli[1])

	}
}

func Check_Coord(str string, i int) []string {
	sli := strings.Split(str, " ")
	if len(sli) != 3 || sli[0] == " " || sli[0] == "#" || sli[0] == "L" {
		if sli[0] == "##start" || sli[0] == "##end" {
			fmt.Fprintln(os.Stderr, "You give more than one start or end check your file !!!")
			os.Exit(1)
		} else {
			fmt.Fprintf(os.Stderr, "Bad Data in line %d\n", i+2)
			os.Exit(1)
		}
	}
	_, err1 := strconv.Atoi(sli[1])
	_, err2 := strconv.Atoi(sli[2])
	if err1 != nil || err2 != nil {
		fmt.Fprintf(os.Stderr, "Bad coord in the line %d\n", i+1)
		os.Exit(1)
	}
	return sli
}

func SearchInTheGraph(current *Vertex, paths *DFS, path []string) {
	Temp_Path := make([]string, len(path))
	copy(Temp_Path, path)
	current.vesited = 1
	for _, jar := range current.adjecent {
		if current.Etat == "end" {
			Temp_Path = append(Temp_Path, current.value)
			paths.Paths = append(paths.Paths, Temp_Path)
			break
		} else if jar.vesited == 1 {
			continue
		} else {
			if len(Temp_Path) > 0 {
				if Temp_Path[len(Temp_Path)-1] == current.value {
					Temp_Path = Temp_Path[:len(Temp_Path)-1]
				}
			}
			Temp_Path = append(Temp_Path, current.value)
			SearchInTheGraph(jar, paths, Temp_Path)
		}
	}
	current.vesited = 0
}

func CheckRepition(arr1 [][]string, arr2 []string) bool {
	element := make(map[string]string)
	for i := 0; i < len(arr1); i++ {
		for j := 1; j < len(arr1[i])-1; j++ {
			element[arr1[i][j]] = "y"
		}
	}

	for i := 1; i < len(arr2)-1; i++ {
		if _, exist := element[arr2[i]]; exist {
			return true
		}
	}
	return false
}

func (paths *DFS) SortPaths() {
	for i := 0; i < len(paths.Paths)-1; i++ {
		for j := i + 1; j < len(paths.Paths); j++ {
			if len(paths.Paths[i]) > len(paths.Paths[j]) {
				paths.Paths[i], paths.Paths[j] = paths.Paths[j], paths.Paths[i]
			}
		}
	}
}

func ChooseUniquePaths(paths *DFS) {

	for i := 0; i < len(paths.Paths); i++ {
		unique := [][]string{}
		unique = append(unique, paths.Paths[i])

		for j := 0; j < len(paths.Paths); j++ {
			if i != j && !CheckRepition(unique, paths.Paths[j])  {
				unique = append(unique, paths.Paths[j])
			}
		}
		paths.Unique_Paths = append(paths.Unique_Paths, unique)
	}
}
