package main

import (
	"container/heap"
)

// ComputeContractions 123
func (graph *Graph) ComputeContractions() {
	if graph.contracted {
		return
	}

	importanceHeap := graph.computeImportance()
	orderCounter := 0
	contractionOrder := make([]int64, len(graph.vertices))

	for importanceHeap.Len() != 0 {
		vertex := heap.Pop(importanceHeap).(*Vertex)
		vertex.importance = len(vertex.inwardEdges)*len(vertex.outwardEdges) - len(vertex.inwardEdges) - len(vertex.outwardEdges)

		// check if the vertex is indeed the least important one
		if importanceHeap.Len() != 0 && vertex.importance > importanceHeap.Peek().(*Vertex).importance {
			heap.Push(importanceHeap, vertex)
			continue
		}
		contractionOrder[orderCounter] = vertex.id
		vertex.contractionOrder = orderCounter
		graph.contractVertex(vertex, int64(orderCounter))

		orderCounter = orderCounter + 1
	}

	graph.contracted = true
}

// contractVertex 123
func (graph *Graph) contractVertex(vertex *Vertex, contractID int64) {
	vertex.contracted = true

	inMax := 0.0
	outMax := 0.0

	for i := 0; i < len(vertex.inwardEdges); i++ {
		if !vertex.inwardEdges[i].from.contracted && inMax < vertex.inwardEdges[i].weight {
			inMax = vertex.inwardEdges[i].weight
		}
	}

	for i := 0; i < len(vertex.outwardEdges); i++ {
		if !vertex.outwardEdges[i].to.contracted && outMax < vertex.outwardEdges[i].weight {
			outMax = vertex.outwardEdges[i].weight
		}
	}

	maxCost := inMax + outMax

	for i := 0; i < len(vertex.inwardEdges); i++ {
		inEdge := vertex.inwardEdges[i]
		if inEdge.from.contracted {
			continue
		}

		graph.relaxEdges(inEdge.from, maxCost, contractID, int64(i))

		for j := 0; j < len(vertex.outwardEdges); j++ {
			outEdge := vertex.outwardEdges[j]
			if outEdge.to.contracted || inEdge.from == outEdge.to {
				continue
			}
			if outEdge.to.distance.contractID != contractID ||
				outEdge.to.distance.sourceID != int64(i) ||
				outEdge.to.distance.distance > inEdge.weight+outEdge.weight {
				if _, ok := graph.contracts[inEdge.from.id]; !ok {
					graph.contracts[inEdge.from.id] = make(map[int64]int64)
				}
				graph.contracts[inEdge.from.id][outEdge.to.id] = vertex.id

				newEdge := &Edge{
					id:         -1,
					from:       inEdge.from,
					to:         outEdge.to,
					weight:     inEdge.weight + outEdge.weight,
					isShortcut: true,
				}
				inEdge.from.outwardEdges = append(inEdge.from.outwardEdges, newEdge)
				outEdge.to.inwardEdges = append(outEdge.to.inwardEdges, newEdge)
			}
		}
	}
}

func (graph *Graph) relaxEdges(source *Vertex, maxCost float64, contractID int64, sourceID int64) {
	distanceHeap := &distanceHeap{}
	heap.Init(distanceHeap)
	heap.Push(distanceHeap, source)

	source.distance.distance = 0
	source.distance.contractID = contractID
	source.distance.sourceID = sourceID

	for distanceHeap.Len() != 0 {
		vertex := heap.Pop(distanceHeap).(*Vertex)
		if vertex.distance.distance > maxCost {
			return
		}
		for i := 0; i < len(vertex.outwardEdges); i++ {
			outEdge := vertex.outwardEdges[i]
			if outEdge.to.contracted {
				continue
			}

			if vertex.distance.contractID != outEdge.to.distance.contractID ||
				vertex.distance.sourceID != outEdge.to.distance.sourceID ||
				outEdge.to.distance.distance > vertex.distance.distance+outEdge.weight {
				// relax
				outEdge.to.distance.distance = vertex.distance.distance + outEdge.weight
				outEdge.to.distance.contractID = contractID
				outEdge.to.distance.sourceID = sourceID
				heap.Push(distanceHeap, outEdge.to)
			}
		}
	}
}

// computeImportance 123
func (graph *Graph) computeImportance() *importanceHeap {
	importanceHeap := &importanceHeap{}
	heap.Init(importanceHeap)
	for i := 0; i < len(graph.vertices); i++ {
		graph.vertices[i].importance = len(graph.vertices[i].inwardEdges)*len(graph.vertices[i].outwardEdges) - len(graph.vertices[i].inwardEdges) - len(graph.vertices[i].outwardEdges)
		importanceHeap.Push(graph.vertices[i])
	}
	return importanceHeap
}
