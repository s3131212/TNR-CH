package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
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

// TestCorrectness Verify the correctness of CH and TNR, comparing with the naive Dijkstra's algorithm.
func TestCorrectness() {
	rand.Seed(time.Now().UnixNano())
	randomList := rand.Perm(100000)
	randomListCounter := 0
	hasWrong := false

	hasEdge := make(map[int64]bool)

	codeList := []string{}

	g := Graph{}
	_ = g

	vertexCount := 100
	edgeCount := 200
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

// Benchmark Benchmark the performance of CH and TNR.
func Benchmark() {
	rand.Seed(time.Now().UnixNano())
	randomList := rand.Perm(10000000)
	randomListCounter := 0

	hasEdge := make(map[int64]bool)

	g := Graph{}
	_ = g

	vertexCount := 3000
	edgeCount := 5000
	tryCount := 3000
	tnrCount := 50
	for i := 0; i < vertexCount; i++ {
		g.AddVertex(int64(i))
	}

	for i := 0; i < edgeCount; i++ {
		a := int64(rand.Intn(vertexCount))
		b := min(max(int64(0), a+(int64(rand.Intn(100))%10-5)), int64(vertexCount-1))
		c := float64(randomList[randomListCounter])
		randomListCounter++
		if hasEdge[a*1000000+b] {
			continue
		}
		hasEdge[a*1000000+b] = true
		hasEdge[b*1000000+a] = true
		g.AddEdge(a, b, c)
		g.AddEdge(b, a, c)
	}

	for i := 0; i < vertexCount-1; i++ {
		a := int64(i)
		b := int64(i + 1)
		c := float64(randomList[randomListCounter])
		randomListCounter++
		d := float64(randomList[randomListCounter])
		randomListCounter++
		if hasEdge[a*1000000+b] {
			continue
		}
		hasEdge[a*1000000+b] = true
		hasEdge[b*1000000+a] = true
		g.AddEdge(a, b, c)
		g.AddEdge(b, a, d)
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

func TestRealWorldRoadMap() {
	vertexCount := -1
	edgeCount := -1
	tnrCount := 1000
	g := Graph{}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == '#' {
			continue
		}
		s := strings.Split(line, " ")
		if len(s) == 3 {
			id, _ := strconv.Atoi(s[0])
			g.AddVertex(int64(id))
		}
		if len(s) == 6 {
			source, _ := strconv.Atoi(s[0])
			target, _ := strconv.Atoi(s[1])
			length, _ := strconv.ParseFloat(s[2], 64)
			g.AddEdge(int64(source), int64(target), length)

			if s[5] == "1" {
				g.AddEdge(int64(target), int64(source), length)
			}
		}
		if vertexCount == -1 {
			i, _ := strconv.Atoi(line)
			vertexCount = i
			continue
		}
		if edgeCount == -1 {
			i, _ := strconv.Atoi(line)
			edgeCount = i
			continue
		}
	}

	fmt.Printf("vertexCount: %d\n", vertexCount)
	fmt.Printf("edgeCount: %d\n", edgeCount)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	tryPath := make([][]int64, 0)
	tryCount := 1000
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

// elapsed Benchmark helper
func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}

// TestDijkstra Benchmark Dijkstra
func TestDijkstra(g *Graph, tryPath [][]int64) {
	defer elapsed("Query using Dijkstra")()
	for i := 0; i < len(tryPath); i++ {
		_ = g.Dijkstra(tryPath[i][0], tryPath[i][1])
	}
}

// TestCH Benchmark Contraction Hierarchies
func TestCH(g *Graph, tryPath [][]int64) {
	defer elapsed("Query using Contraction Hierarchies")()
	for i := 0; i < len(tryPath); i++ {
		_, _ = g.ShortestPathWithoutTNR(tryPath[i][0], tryPath[i][1])
		//fmt.Printf("CH, %d to %d: lenght = %f, path = %v\n", tryPath[i][0], tryPath[i][1], d, p)
	}
}

// TestTNR Benchmark Transit Node Routing
func TestTNR(g *Graph, tryPath [][]int64) {
	defer elapsed("Query using TNR")()
	for i := 0; i < len(tryPath); i++ {
		_, _ = g.ShortestPath(tryPath[i][0], tryPath[i][1])
		//fmt.Printf("TNR, %d to %d: lenght = %f, path = %v\n", tryPath[i][0], tryPath[i][1], d, p)
	}
}

// ComputeContractions Benchmark ComputeContractions
func ComputeContractions(g *Graph) {
	defer elapsed("Compute Contraction Hierarchies")()
	g.ComputeContractions()
}

// ComputeTNR Benchmark ComputeTNR
func ComputeTNR(g *Graph, n int) {
	defer elapsed("Compute TNR")()
	g.ComputeTNR(n)
}
