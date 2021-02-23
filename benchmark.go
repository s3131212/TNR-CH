package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// ComparePerformace Benchmark the performance of CH and TNR.
func ComparePerformace() {
	rand.Seed(time.Now().UnixNano())
	randomList := rand.Perm(10000000)
	randomListCounter := 0

	hasEdge := make(map[int64]bool)

	g := Graph{}
	_ = g

	vertexCount := 5000
	edgeCount := 7000
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
	BmDijkstra(&g, tryPath)
	BmCH(&g, tryPath)
	BmTNR(&g, tryPath)
}

func ComparePerformanceOnRealWorldRoadMap() {
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
	BmDijkstra(&g, tryPath)
	BmCH(&g, tryPath)
	BmTNR(&g, tryPath)
}

// elapsed Benchmark helper
func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}

// BmDijkstra Benchmark Dijkstra
func BmDijkstra(g *Graph, tryPath [][]int64) {
	defer elapsed("Query using Dijkstra")()
	for i := 0; i < len(tryPath); i++ {
		_, _ = g.Dijkstra(tryPath[i][0], tryPath[i][1])
	}
}

// BmCH Benchmark Contraction Hierarchies
func BmCH(g *Graph, tryPath [][]int64) {
	defer elapsed("Query using Contraction Hierarchies")()
	for i := 0; i < len(tryPath); i++ {
		_, _ = g.ShortestPathWithoutTNR(tryPath[i][0], tryPath[i][1])
		//fmt.Printf("CH, %d to %d: lenght = %f, path = %v\n", tryPath[i][0], tryPath[i][1], d, p)
	}
}

// BmTNR Benchmark Transit Node Routing
func BmTNR(g *Graph, tryPath [][]int64) {
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
