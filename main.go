package main

import (
	"fmt"
	"math/rand"
	"time"
)

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
		_, _ = g.ShortestPathWithoutTNR(tryPath[i][0], tryPath[i][1])
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

func main() {
	rand.Seed(time.Now().UnixNano())
	randomList := rand.Perm(10000)
	randomListCounter := 0

	hasEdge := make(map[int64]bool)

	codeList := []string{}

	g := Graph{}
	_ = g

	vertexCount := 1000
	edgeCount := 8000
	tryCount := 1000
	tnrCount := 50
	for i := 0; i < vertexCount; i++ {
		g.AddVertex(int64(i))
		codeList = append(codeList, fmt.Sprintf("g.AddVertex(int64(%d))", i))
	}

	for i := 0; i < edgeCount; i++ {
		a := int64(rand.Intn(vertexCount))
		b := int64(rand.Intn(vertexCount))
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
