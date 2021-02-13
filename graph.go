package main

import (
	"container/heap"
	"container/list"
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
func (graph *Graph) Dijkstra(from int64, to int64) (float64, []int64) {
	fromVertex := graph.vertices[graph.mapping[from]]
	toVertex := graph.vertices[graph.mapping[to]]

	distance := make([]float64, len(graph.vertices), len(graph.vertices))
	for i := 0; i < len(graph.vertices); i++ {
		distance[i] = math.MaxFloat64
	}

	distanceHeap := &forwardSearchHeap{}
	parent := make(map[int64]int64)

	distance[fromVertex.id] = 0

	heap.Init(distanceHeap)
	heap.Push(distanceHeap, &QueryVertex{
		id:               fromVertex.id,
		forwardDistance:  0,
		backwardDistance: 0,
	})

	for distanceHeap.Len() != 0 {
		queryVertex := heap.Pop(distanceHeap).(*QueryVertex)

		if distance[queryVertex.id] < queryVertex.forwardDistance {
			continue
		}

		if queryVertex.id == toVertex.id {
			break
		}

		for i := 0; i < len(graph.vertices[queryVertex.id].outwardEdges); i++ {
			outEdge := graph.vertices[queryVertex.id].outwardEdges[i]
			if outEdge.isShortcut {
				continue
			}
			if distance[queryVertex.id]+outEdge.weight < distance[outEdge.to.id] {
				distance[outEdge.to.id] = distance[queryVertex.id] + outEdge.weight
				parent[outEdge.to.id] = queryVertex.id
				heap.Push(distanceHeap, &QueryVertex{
					id:               outEdge.to.id,
					forwardDistance:  distance[outEdge.to.id],
					backwardDistance: 0,
				})
			}
		}
	}
	if distance[toVertex.id] == math.MaxFloat64 {
		return -math.MaxFloat64, nil
	}

	pathList := list.New()
	v := to
	ok := false
	pathList.PushFront(v)
	for {
		if v, ok = parent[v]; ok {
			pathList.PushFront(v)
		} else {
			break
		}
	}
	// list to slice
	var path []int64
	for e := pathList.Front(); e != nil; e = e.Next() {
		path = append(path, graph.vertices[e.Value.(int64)].name)
	}

	return distance[toVertex.id], path
}
