package main

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"time"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func verifyCorrectness() {
	rand.Seed(time.Now().UnixNano())
	randomList := rand.Perm(100000)
	randomListCounter := 0
	hasWrong := false

	hasEdge := make(map[int64]bool)

	codeList := []string{}

	g := Graph{}
	_ = g

	vertexCount := 100
	edgeCount := 250
	tnrCount := 10
	for i := 0; i < vertexCount; i++ {
		g.AddVertex(int64(i))
		codeList = append(codeList, fmt.Sprintf("g.AddVertex(int64(%d))", i))
	}

	for i := 0; i < edgeCount; i++ {
		a := int64(rand.Intn(vertexCount))
		b := min(max(int64(0), a+(int64(rand.Intn(100))%10-5)), int64(vertexCount-1))
		//b := int64(rand.Intn(vertexCount))
		c := float64(randomList[randomListCounter])
		randomListCounter++
		d := float64(randomList[randomListCounter])
		randomListCounter++
		if hasEdge[a*10000+b] || a == b {
			continue
		}
		hasEdge[a*10000+b] = true
		hasEdge[b*10000+a] = true
		g.AddEdge(a, b, c)
		g.AddEdge(b, a, d)
		codeList = append(codeList, fmt.Sprintf("g.AddEdge(%d, %d, %f)", a, b, c))
		codeList = append(codeList, fmt.Sprintf("g.AddEdge(%d, %d, %f)", b, a, d))
	}

	// make sure the graph is connected
	/*
		for i := 0; i < vertexCount-1; i++ {
			a := int64(i)
			b := int64(i + 1)
			c := float64(randomList[randomListCounter])
			randomListCounter++
			d := float64(randomList[randomListCounter])
			randomListCounter++
			if hasEdge[a*10000+b] {
				continue
			}
			hasEdge[a*10000+b] = true
			hasEdge[b*10000+a] = true
			g.AddEdge(a, b, c)
			g.AddEdge(b, a, d)
			codeList = append(codeList, fmt.Sprintf("g.AddEdge(%d, %d, %f)", a, b, c))
			codeList = append(codeList, fmt.Sprintf("g.AddEdge(%d, %d, %f)", b, a, d))
		}
	*/

	g.ComputeContractions()
	g.ComputeTNR(tnrCount)

	for i := 0; i < vertexCount; i++ {
		for j := i + 1; j < vertexCount; j++ {
			if i == j {
				continue
			}

			d1, p1 := g.ShortestPathWithoutTNR(int64(i), int64(j))
			d2, p2 := g.ShortestPath(int64(i), int64(j))
			d := g.Dijkstra(int64(i), int64(j))
			d_ := d
			// check if Dijkstra is wrong :(
			if d1 < d || d2 < d {
				d = 0
				for k := 0; k < len(p1)-1; k++ {
					for m := 0; m < len(g.vertices[g.mapping[p1[k]]].outwardEdges); m++ {
						if g.vertices[g.mapping[p1[k]]].outwardEdges[m].isShortcut {
							continue
						}
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
				// sometimes different path isn't an issue as long as both are the shortest
				d1_ := float64(0)
				for k := 0; k < len(p1)-1; k++ {
					for m := 0; m < len(g.vertices[g.mapping[p1[k]]].outwardEdges); m++ {
						if g.vertices[g.mapping[p1[k]]].outwardEdges[m].isShortcut {
							continue
						}
						if g.vertices[g.mapping[p1[k]]].outwardEdges[m].to.id == g.mapping[p1[k+1]] {
							d1_ += g.vertices[g.mapping[p1[k]]].outwardEdges[m].weight
							break
						}
					}
				}

				d2_ := float64(0)
				for k := 0; k < len(p2)-1; k++ {
					for m := 0; m < len(g.vertices[g.mapping[p2[k]]].outwardEdges); m++ {
						if g.vertices[g.mapping[p2[k]]].outwardEdges[m].isShortcut {
							continue
						}
						if g.vertices[g.mapping[p2[k]]].outwardEdges[m].to.id == g.mapping[p2[k+1]] {
							d2_ += g.vertices[g.mapping[p2[k]]].outwardEdges[m].weight
							break
						}
					}
				}

				if d1_ == d2_ {
					continue
				}

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

// Benchmark 123
func Benchmark() {
	rand.Seed(time.Now().UnixNano())
	randomList := rand.Perm(100000)
	randomListCounter := 0

	hasEdge := make(map[int64]bool)

	g := Graph{}
	_ = g

	vertexCount := 6000
	edgeCount := 14000
	tryCount := 1000
	tnrCount := 100
	for i := 0; i < vertexCount; i++ {
		g.AddVertex(int64(i))
	}

	for i := 0; i < edgeCount; i++ {
		a := int64(rand.Intn(vertexCount))
		b := min(max(int64(0), a+(int64(rand.Intn(100))%10-5)), int64(vertexCount-1))
		c := float64(randomList[randomListCounter])
		randomListCounter++
		if hasEdge[a*10000+b] {
			continue
		}
		hasEdge[a*10000+b] = true
		hasEdge[b*10000+a] = true
		g.AddEdge(a, b, c)
		g.AddEdge(b, a, c)
	}

	tryPath := make([][]int64, 0)
	for i := 0; i < tryCount; i++ {
		a := rand.Intn(vertexCount)
		b := rand.Intn(vertexCount)
		for a == b {
			a = rand.Intn(vertexCount)
			b = rand.Intn(vertexCount)
		}
		tryPath = append(tryPath, []int64{int64(a), int64(b)})
	}

	ComputeContractions(&g)
	ComputeTNR(&g, tnrCount)
	TestDijkstra(&g, tryPath)
	TestCH(&g, tryPath)
	TestTNR(&g, tryPath)
}

func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}

// TestDijkstra 123
func TestDijkstra(g *Graph, tryPath [][]int64) {
	defer elapsed("Query using Dijkstra")()
	for i := 0; i < len(tryPath); i++ {
		_ = g.Dijkstra(tryPath[i][0], tryPath[i][1])
	}
}

// TestCH 123
func TestCH(g *Graph, tryPath [][]int64) {
	defer elapsed("Query using Contraction Hierarchies")()
	for i := 0; i < len(tryPath); i++ {
		_, _ = g.ShortestPathWithoutTNR(tryPath[i][0], tryPath[i][1])
		//fmt.Printf("CH, %d to %d: lenght = %f, path = %v\n", tryPath[i][0], tryPath[i][1], d, p)
	}
}

// TestTNR 123
func TestTNR(g *Graph, tryPath [][]int64) {
	defer elapsed("Query using TNR")()
	for i := 0; i < len(tryPath); i++ {
		_, _ = g.ShortestPath(tryPath[i][0], tryPath[i][1])
		//fmt.Printf("TNR, %d to %d: lenght = %f, path = %v\n", tryPath[i][0], tryPath[i][1], d, p)
	}
}

// ComputeContractions 123
func ComputeContractions(g *Graph) {
	defer elapsed("Compute Contraction Hierarchies")()
	g.ComputeContractions()
}

// ComputeTNR 123
func ComputeTNR(g *Graph, n int) {
	defer elapsed("Compute TNR")()
	g.ComputeTNR(n)
}
