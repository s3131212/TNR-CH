# TNR-CH
Implementation of Transit Node Routing + Contraction Hierarchies, inspired by [Transit Node Routing Reconsidered](https://arxiv.org/abs/1302.5611). This implementation includes Contraction Hierarchies, transit node selection, and Search Space Based Locality Filter, and Graph Voronoi Label Compression.

## Usage
Let the code speak.
```go
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
```

## Benchmark
Query 1,000 randomly chosen path with 1,000 transit nodes.

On Taipei's roadmap:
```
vertexCount: 48753
edgeCount: 62157
Compute Contraction Hierarchies took 4.464321111s
Compute TNR took 13m14.049223101s
Query using Dijkstra took 12.859151479s
Query using Contraction Hierarchies took 820.613105ms
Query using TNR took 264.495739ms
```

## Reference
The code of Contraction Hierarchies is inspired by [LdDl/ch](https://github.com/LdDl/ch). The method is from [Arz et al](https://arxiv.org/abs/1302.5611).

## License
MIT License

## TODO
1. The TNR Computation can be faster since the method now is quite naive and it does many redundant calculation.
2. Memory efficiency.
3. Change `int64` to `int`.
4. Import and export computed graph.
5. `go benchmark`, currently you can call `ComparePerformace` to do the benchmark.