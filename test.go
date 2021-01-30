package main

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"time"
)

func verifyCorrectness() {
	rand.Seed(time.Now().UnixNano())
	randomList := rand.Perm(10000)
	randomListCounter := 0
	hasWrong := false

	hasEdge := make(map[int64]bool)

	codeList := []string{}

	g := Graph{}
	_ = g

	vertexCount := 1000
	edgeCount := 10000
	for i := 0; i < vertexCount; i++ {
		g.AddVertex(int64(i + 1))
		codeList = append(codeList, fmt.Sprintf("g.AddVertex(int64(%d))", i))
	}

	for i := 0; i < edgeCount; i++ {
		a := int64(rand.Intn(vertexCount) + 1)
		b := int64(rand.Intn(vertexCount) + 1)
		c := float64(randomList[randomListCounter])
		randomListCounter++
		if hasEdge[a*10000+b] {
			continue
		}
		hasEdge[a*10000+b] = true
		hasEdge[b*10000+a] = true
		g.AddEdge(a, b, c)
		g.AddEdge(b, a, c)
		codeList = append(codeList, fmt.Sprintf("g.AddEdge(%d, %d, %f)", a, b, c))
		codeList = append(codeList, fmt.Sprintf("g.AddEdge(%d, %d, %f)", b, a, c))
	}
	g.ComputeContractions()
	g.ComputeTNR(50)

	for i := 0; i < vertexCount; i++ {
		for j := i + 1; j < vertexCount; j++ {
			if i == j {
				continue
			}

			d1, p1 := g.ShortestPathWithoutTNR(int64(i+1), int64(j+1))
			d2, p2 := g.ShortestPath(int64(i+1), int64(j+1))
			d := g.Dijkstra(int64(i+1), int64(j+1))
			d_ := d
			// check if Dijkstra is wrong :(
			if d1 < d || d2 < d {
				d = 0
				for k := 0; k < len(p1)-1; k++ {
					for m := 0; m < len(g.vertices[g.mapping[p1[k]]].outwardEdges); m++ {
						if g.vertices[g.mapping[p1[k]]].outwardEdges[m].to.id == g.mapping[p1[k+1]] {
							d += g.vertices[g.mapping[p1[k]]].outwardEdges[m].weight
							break
						}
					}
				}
			}

			if d != d1 {
				fmt.Printf("\n")
				fmt.Println("==================================================error==================================================")
				fmt.Printf("%d to %d: ShortestPathWithoutTNR wrong, it gives %f but the asnwer is %f\n", i, j, d1, d_)
				fmt.Printf("path: %v\n", p1)
				hasWrong = true
			}
			if d != d2 {
				fmt.Printf("\n")
				fmt.Println("==================================================error==================================================")
				fmt.Printf("%d to %d: ShortestPath wrong, it gives %f but the asnwer is %f\n", i, j, d2, d_)
				fmt.Printf("path: %v\n", p1)
				hasWrong = true
			}
			if !reflect.DeepEqual(p1, p2) {
				fmt.Printf("\n")
				fmt.Println("==================================================error=================================================")
				fmt.Printf("%d to %d: path different\n", i, j)
				fmt.Printf("ShortestPathWithoutTNR: %v, ShortestPath: %v\n", p1, p2)
				hasWrong = true
			} else if d != -math.MaxFloat64 {
				// fmt.Printf("%d to %d length: %f, path: %v\n", i, j, d1, p1)
			}
		}
	}
	if hasWrong {
		for i := 0; i < len(codeList); i++ {
			fmt.Println(codeList[i])
		}
	} else {
		fmt.Println("everything all good")
	}
}
