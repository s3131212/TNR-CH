package main

// Vertex the vertex / node on the graph.
type Vertex struct {
	name int64
	id   int64

	inwardEdges  []*Edge
	outwardEdges []*Edge
	distance     *Distance

	importance       int
	contractionOrder int

	contracted    bool
	isTransitNode bool

	voronoiRegionID     int64
	forwardSearchSpace  map[int64]bool
	backwardSearchSpace map[int64]bool

	forwardAccessNodeDistance  map[int64]float64
	backwardAccessNodeDistance map[int64]float64

	forwardTNRed  bool
	backwardTNRed bool
}

// QueryVertex a temporary and simplified vertex data structure, only for using heaps.
type QueryVertex struct {
	id               int64
	isTransitNode    bool
	forwardDistance  float64
	backwardDistance float64
}

// Distance the (current) optimal distance of a vertex.
type Distance struct {
	distance         float64
	contractID       int64
	sourceID         int64
	forwardDistance  float64
	backwardDistance float64
}

// Edge the edge / road on the graph.
type Edge struct {
	id         int64
	from       *Vertex
	to         *Vertex
	weight     float64
	isShortcut bool
}

// Graph the graph / road network.
type Graph struct {
	vertices []*Vertex
	heap     *minHeap

	mapping   map[int64]int64
	contracts map[int64]map[int64]int64

	transitNodes []*Vertex
	tnrDistance  map[int64]map[int64]float64
	tnrPath      map[int64]map[int64][]int64

	contracted bool
	TNRed      bool
}
