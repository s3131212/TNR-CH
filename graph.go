package main

import (
	"container/heap"
	"math"
)

// NewGraph 123
func NewGraph() Graph {
	return Graph{
		contracted: false,
		TNRed:      false,
	}
}

// AddVertex 123
func (graph *Graph) AddVertex(name int64) {
	if _, ok := graph.mapping[name]; ok {
		return
	}

	newVertex := &Vertex{
		id:            int64(len(graph.vertices)),
		name:          name,
		distance:      NewDistance(),
		contracted:    false,
		isTransitNode: false,
	}
	if graph.mapping == nil {
		graph.mapping = make(map[int64]int64)
	}
	if graph.contracts == nil {
		graph.contracts = make(map[int64]map[int64]int64)
	}

	graph.mapping[name] = int64(len(graph.vertices))
	graph.vertices = append(graph.vertices, newVertex)
}

// AddEdge 123
func (graph *Graph) AddEdge(from int64, to int64, weight float64) {
	fromVertex := graph.vertices[graph.mapping[from]]
	toVertex := graph.vertices[graph.mapping[to]]

	newEdge := &Edge{
		from:       fromVertex,
		to:         toVertex,
		weight:     weight,
		isShortcut: false,
	}
	fromVertex.outwardEdges = append(fromVertex.outwardEdges, newEdge)
	toVertex.inwardEdges = append(toVertex.inwardEdges, newEdge)
}

// NewDistance 123
func NewDistance() *Distance {
	return &Distance{
		distance:   math.MaxFloat64,
		contractID: -1,
		sourceID:   -1,
	}
}

// Dijkstra Dijkstra's algorithm
func (graph *Graph) Dijkstra(from int64, to int64) float64 {
	fromVertex := graph.vertices[graph.mapping[from]]
	toVertex := graph.vertices[graph.mapping[to]]

	for i := 0; i < len(graph.vertices); i++ {
		graph.vertices[i].distance.distance = math.MaxFloat64
	}

	distanceHeap := &distanceHeap{}
	visited := make(map[int64]bool)
	fromVertex.distance.distance = 0
	heap.Init(distanceHeap)
	heap.Push(distanceHeap, fromVertex)

	for distanceHeap.Len() != 0 {
		vertex := heap.Pop(distanceHeap).(*Vertex)

		if visited[vertex.id] {
			continue
		}
		visited[vertex.id] = true

		for i := 0; i < len(vertex.outwardEdges); i++ {
			if vertex.outwardEdges[i].isShortcut {
				continue
			}
			if vertex.distance.distance+vertex.outwardEdges[i].weight < vertex.outwardEdges[i].to.distance.distance {
				vertex.outwardEdges[i].to.distance.distance = vertex.distance.distance + vertex.outwardEdges[i].weight
				heap.Push(distanceHeap, vertex.outwardEdges[i].to)
			}
		}
	}
	if toVertex.distance.distance == math.MaxFloat64 {
		return -math.MaxFloat64
	}
	return toVertex.distance.distance
}
