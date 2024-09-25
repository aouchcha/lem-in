package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Graph struct {
	vertices []*Vertix
}

type Vertix struct {
	vesited  int
	Etat     string
	value    string
	adjecent []*Vertix
}

func (g *Graph) AddVertix(v string, stat string) {
	for _, vertices := range g.vertices {
		if vertices.value == v {
			fmt.Fprintf(os.Stderr, "You have already a vertix with the value %s\n", v)
			return
		}
	}
	g.vertices = append(g.vertices, &Vertix{value: v, Etat: stat})
}

func (g *Graph) AddPaths(from, to string) {
	foundf := false
	foundt := false
	for _, vertice := range g.vertices {
		if vertice.value == from {
			foundf = true
			for _, vertice2 := range vertice.adjecent {
				if vertice2.value == to {
					foundt = true
					fmt.Fprintf(os.Stderr, "There is already a relation between %s and %s\n", from, to)
					return
				}
			}
		}
	}

	for _, vertice := range g.vertices {
		if vertice.value == from {
			foundf = true
			for _, vertice2 := range g.vertices {
				if vertice2.value == to {
					foundt = true
					vertice.adjecent = append(vertice.adjecent, vertice2)
					vertice2.adjecent = append(vertice2.adjecent, vertice)
				}
			}
		}
	}

	if !foundf {
		fmt.Fprintf(os.Stderr, "There is no vertix with the name %s\n", from)
		return
	} else if !foundt {
		fmt.Fprintf(os.Stderr, "There is no vertix with the name %s\n", to)
		return
	}
}

func (g *Graph) Print() {
	for _, vertice := range g.vertices {
		fmt.Fprintf(os.Stdout, "you add a vertex with value %s to the graph\n", vertice.value)
	}
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

	CreatRoomsAndPaths(gr, TraitData())

	gr.Print()
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
	var roomes []string
	var links []string
	for i := 1; i < len(lines); i++ {
		if strings.Contains(lines[i], " ") || strings.HasPrefix(lines[i], "##") {
			roomes = append(roomes, lines[i])
		} else if strings.Contains(lines[i], "-") {
			links = append(links, lines[i])
		}
	}
	// for i := range roomes {
	// 	fmt.Println(roomes[i])
	// }
	// fmt.Println(links)
	for i := 0; i < len(roomes); i++ {
		// fmt.Println(roomes[i])
		if strings.HasPrefix(roomes[i], "#") && (roomes[i] != "##start" && roomes[i] != "##end") {
			continue
		} else if roomes[i] == "##start" || roomes[i] == "##end" {
			sli := Check_Coord(roomes[i+1], i+1)
			if len(sli) != 3 || sli[0] == " " || sli[0] == "#" || sli[0] == "L" {
				fmt.Fprintf(os.Stderr, "Bad Data in line %d", i+1)
			}
			if roomes[i] == "##start" {
				gr.AddVertix(sli[0], "start")
			} else {
				gr.AddVertix(sli[0], "end")
			}
			i++
		} else {
			sli := Check_Coord(roomes[i], i)
			gr.AddVertix(sli[0], "standard")
		}
	}

	for i := 0; i < len(links); i++ {

		sli := strings.Split(links[i], "-")
		if len(sli) != 2 {
			fmt.Fprintf(os.Stderr, "Bad Data on links in line %s\n", links[i])
			continue
		}
		gr.AddPaths(sli[0], sli[1])

	}
}

func Check_Coord(str string, i int) []string {
	sli := strings.Split(str, " ")
	if len(sli) != 3 || sli[0] == " " || sli[0] == "#" || sli[0] == "L" {
		fmt.Fprintf(os.Stderr, "Bad Data on normal cordinat in line %d\n", i+1)
		os.Exit(1)
	}
	_, err1 := strconv.Atoi(sli[1])
	_, err2 := strconv.Atoi(sli[2])
	if err1 != nil || err2 != nil {
		fmt.Fprintf(os.Stderr, "Bad coord in the line %d\n", i+1)
		os.Exit(1)
	}
	return sli
}
