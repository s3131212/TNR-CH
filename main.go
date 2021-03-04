package main

import "fmt"

func main() {
	g := Graph{}

	// g.AddVertex(<vertex name>)
	g.AddVertex(int64(0))
	g.AddVertex(int64(1))
	g.AddVertex(int64(2))
	g.AddVertex(int64(3))
	g.AddVertex(int64(4))

	// g.AddEdge(<source vertex name>, <target vertex name>, <length>)
	g.AddEdge(int64(0), int64(1), 1.0) // Note that AddEdge only add an uni-directional edge
	g.AddEdge(int64(1), int64(0), 1.0) // If one needs bi-directional edges, simply add another direction
	g.AddEdge(int64(1), int64(2), 2.0)
	g.AddEdge(int64(2), int64(3), 3.0)
	g.AddEdge(int64(4), int64(2), 4.0)
	g.AddEdge(int64(2), int64(4), 2.0)

	g.ComputeContractions() // Compute the contractions before using contraction hierarchies
	g.ComputeTNR(2)         // Compute transit nodes before using TNR algorithm, 2 stands for the amount of transit nodes

	distanceCH, pathCH := g.ShortestPathWithoutTNR(int64(1), int64(4)) // Compute shortest paths without using TNR
	distanceTNR, pathTNR := g.ShortestPath(int64(1), int64(4))         // Compute shortest path using TNR if possible, fallback to CH for local paths
	distanceDijkstra, pathDijkstra := g.Dijkstra(int64(1), int64(4))   // Naive Dijkstra

	fmt.Printf("Shortest path using CH: %v, %f\n", pathCH, distanceCH)
	fmt.Printf("Shortest path using TNR+CH: %v, %f\n", pathTNR, distanceTNR)
	fmt.Printf("Shortest path using Dijkstra: %v, %f\n", pathDijkstra, distanceDijkstra)
}
