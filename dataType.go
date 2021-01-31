package main

// Vertex 123
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

	forwardReachableVertex  map[int64]bool
	backwardReachableVertex map[int64]bool

	forwardAccessNodeDistance  map[int64]float64
	backwardAccessNodeDistance map[int64]float64

	forwardAccessNodePath  map[int64][]int64
	backwardAccessNodePath map[int64][]int64
}

// QueryVertex 123
type QueryVertex struct {
	id               int64
	isTransitNode    bool
	forwardDistance  float64
	backwardDistance float64
}

// Distance 123
type Distance struct {
	distance         float64
	contractID       int64
	sourceID         int64
	forwardDistance  float64
	backwardDistance float64
}

// Edge 123
type Edge struct {
	id         int64
	from       *Vertex
	to         *Vertex
	weight     float64
	isShortcut bool
}

// Graph 123
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
