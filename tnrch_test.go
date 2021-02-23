package main

import (
	"math/rand"
	"reflect"
	"testing"
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
func TestCorrectness(t *testing.T) {
	randomList := rand.Perm(100000)
	randomListCounter := 0

	hasEdge := make(map[int64]bool)

	g := Graph{}
	_ = g

	vertexCount := 100
	edgeCount := 150
	tnrCount := 5
	for i := 0; i < vertexCount; i++ {
		g.AddVertex(int64(i))
	}

	for i := 0; i < edgeCount; i++ {
		a := int64(rand.Intn(vertexCount))
		b := min(max(int64(0), a+(int64(rand.Intn(100))%10-5)), int64(vertexCount-1))
		c := float64(randomList[randomListCounter])
		d := float64(randomList[randomListCounter+1])
		randomListCounter += 2
		if hasEdge[a*10000+b] || a == b {
			continue
		}
		hasEdge[a*10000+b] = true
		hasEdge[b*10000+a] = true
		g.AddEdge(a, b, c)
		g.AddEdge(b, a, d)
	}

	g.ComputeContractions()
	g.ComputeTNR(tnrCount)

	for i := 0; i < vertexCount; i++ {
		for j := i + 1; j < vertexCount; j++ {
			if i == j {
				continue
			}

			distanceCH, pathCH := g.ShortestPathWithoutTNR(int64(i), int64(j))
			distanceTNR, pathTNR := g.ShortestPath(int64(i), int64(j))
			distanceDijkstra, pathDijkstra := g.Dijkstra(int64(i), int64(j))

			if distanceCH != distanceDijkstra {
				t.Errorf("%d to %d: ShortestPathWithoutTNR wrong, it gives %f but the answer is %f", i, j, distanceCH, distanceDijkstra)
			}

			if distanceTNR != distanceDijkstra {
				t.Errorf("%d to %d: ShortestPath wrong, it gives %f but the answer is %f", i, j, distanceTNR, distanceDijkstra)
			}

			if !reflect.DeepEqual(pathCH, pathDijkstra) {
				// sometimes different path isn't an issue as long as both are the shortest
				distanceCHTemp := CalculatePathDistance(pathCH, g)
				distanceDijkstraTemp := CalculatePathDistance(pathDijkstra, g)

				if distanceCHTemp != distanceDijkstraTemp {
					t.Errorf("%d to %d: path wrong, ShortestPathWithoutTNR gives %v but the answer is %v", i, j, pathCH, pathDijkstra)
				}
			}

			if !reflect.DeepEqual(pathTNR, pathDijkstra) {
				distanceTNRTemp := CalculatePathDistance(pathTNR, g)
				distanceDijkstraTemp := CalculatePathDistance(pathDijkstra, g)

				if distanceTNRTemp != distanceDijkstraTemp {
					t.Errorf("%d to %d: path wrong, ShortestPath gives %v but the answer is %v", i, j, pathTNR, pathDijkstra)
				}
			}
		}
	}
}

// CalculatePathDistance calculate the distance given a path
func CalculatePathDistance(path []int64, g Graph) float64 {
	distance := float64(0)
	for k := 0; k < len(path)-1; k++ {
		for m := 0; m < len(g.vertices[g.mapping[path[k]]].outwardEdges); m++ {
			if g.vertices[g.mapping[path[k]]].outwardEdges[m].isShortcut {
				continue
			}
			if g.vertices[g.mapping[path[k]]].outwardEdges[m].to.id == g.mapping[path[k+1]] {
				distance += g.vertices[g.mapping[path[k]]].outwardEdges[m].weight
				break
			}
		}
	}
	return distance
}
