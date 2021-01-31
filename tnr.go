package main

import (
	"container/heap"
	"fmt"
	"math"
)

// ComputeTNR 123
func (graph *Graph) ComputeTNR(transitCnt int) {
	if graph.TNRed || !graph.contracted {
		fmt.Println("Not contracted or is TNRed")
		return
	}
	vertexCnt := len(graph.vertices)
	if vertexCnt < transitCnt {
		fmt.Println("Too many transit nodes")
		return
	}

	if graph.tnrDistance == nil {
		graph.tnrDistance = make(map[int64]map[int64]float64)
	}
	if graph.tnrPath == nil {
		graph.tnrPath = make(map[int64]map[int64][]int64)
	}

	for i := 0; i < vertexCnt; i++ {
		//fmt.Printf("id: %d, contractionOrder: %d\n", graph.vertices[i].id, graph.vertices[i].contractionOrder)
		if graph.vertices[i].contractionOrder >= vertexCnt-transitCnt {
			graph.vertices[i].isTransitNode = true
			graph.transitNodes = append(graph.transitNodes, graph.vertices[i])
			//fmt.Printf("select transit %d\n", graph.vertices[i].id)
		}
	}
	for i := 0; i < transitCnt; i++ {
		for j := 0; j < transitCnt; j++ {
			if _, ok := graph.tnrDistance[graph.transitNodes[i].id]; !ok {
				graph.tnrDistance[graph.transitNodes[i].id] = make(map[int64]float64)
			}
			if _, ok := graph.tnrPath[graph.transitNodes[i].id]; !ok {
				graph.tnrPath[graph.transitNodes[i].id] = make(map[int64][]int64)
			}

			if i == j {
				graph.tnrDistance[graph.transitNodes[i].id][graph.transitNodes[j].id] = 0
				graph.tnrPath[graph.transitNodes[i].id][graph.transitNodes[j].id] = []int64{}
				continue
			}

			distance, path := graph.ShortestPathWithoutTNR(graph.transitNodes[i].name, graph.transitNodes[j].name)

			graph.tnrDistance[graph.transitNodes[i].id][graph.transitNodes[j].id] = distance
			graph.tnrPath[graph.transitNodes[i].id][graph.transitNodes[j].id] = path

			//fmt.Printf("between transit nodes %d and %d: %f, %v\n", graph.transitNodes[i].id, graph.transitNodes[j].id, distance, path)
		}
	}

	// Preprocess forward search
	for v := 0; v < len(graph.vertices); v++ {
		sourceVertex := graph.vertices[v]
		if sourceVertex.forwardReachableVertex == nil {
			sourceVertex.forwardReachableVertex = make(map[int64]bool)
		}
		if sourceVertex.forwardAccessNodeDistance == nil {
			sourceVertex.forwardAccessNodeDistance = make(map[int64]float64)
		}
		if sourceVertex.forwardAccessNodePath == nil {
			sourceVertex.forwardAccessNodePath = make(map[int64][]int64)
		}

		// find access node
		searchHeap := &forwardSearchHeap{}
		heap.Init(searchHeap)
		heap.Push(searchHeap, &QueryVertex{
			id:               sourceVertex.id,
			forwardDistance:  0,
			backwardDistance: 0,
			isTransitNode:    sourceVertex.isTransitNode,
		})

		distance := make([]float64, len(graph.vertices), len(graph.vertices))
		for i := 0; i < len(graph.vertices); i++ {
			distance[i] = math.MaxFloat64
		}
		distance[sourceVertex.id] = 0

		for searchHeap.Len() != 0 {
			queryVertex := heap.Pop(searchHeap).(*QueryVertex)

			// relax
			if !queryVertex.isTransitNode {
				sourceVertex.forwardReachableVertex[queryVertex.id] = true
				for i := 0; i < len(graph.vertices[queryVertex.id].outwardEdges); i++ {
					outEdge := graph.vertices[queryVertex.id].outwardEdges[i]
					if graph.vertices[queryVertex.id].contractionOrder < outEdge.to.contractionOrder {
						if distance[outEdge.to.id] > distance[queryVertex.id]+outEdge.weight {
							distance[outEdge.to.id] = distance[queryVertex.id] + outEdge.weight
							heap.Push(searchHeap, &QueryVertex{
								id:               outEdge.to.id,
								forwardDistance:  distance[queryVertex.id] + outEdge.weight,
								backwardDistance: 0,
								isTransitNode:    outEdge.to.isTransitNode,
							})
						}
					}
				}
			} else {
				sourceVertex.forwardAccessNodeDistance[queryVertex.id] = -1
			}
		}

		for k := range sourceVertex.forwardAccessNodeDistance {
			sourceVertex.forwardAccessNodeDistance[k], sourceVertex.forwardAccessNodePath[k] = graph.ShortestPathWithoutTNR(sourceVertex.name, graph.vertices[k].name)
		}

		// delete invalid access node
		accessNodeMask := make(map[int64]bool)
		for k1, d1 := range sourceVertex.forwardAccessNodeDistance {
			for k2, d2 := range sourceVertex.forwardAccessNodeDistance {
				if k1 == k2 {
					continue
				}
				if d1+graph.tnrDistance[k1][k2] <= d2 {
					accessNodeMask[k2] = true // mask j since it won't be the solution
				}
			}
		}
		for k := range accessNodeMask {
			delete(sourceVertex.forwardAccessNodeDistance, k)
			delete(sourceVertex.forwardAccessNodePath, k)
		}
	}

	// Preprocess backward search
	for v := 0; v < len(graph.vertices); v++ {
		sourceVertex := graph.vertices[v]
		if sourceVertex.backwardReachableVertex == nil {
			sourceVertex.backwardReachableVertex = make(map[int64]bool)
		}
		if sourceVertex.backwardAccessNodeDistance == nil {
			sourceVertex.backwardAccessNodeDistance = make(map[int64]float64)
		}
		if sourceVertex.backwardAccessNodePath == nil {
			sourceVertex.backwardAccessNodePath = make(map[int64][]int64)
		}

		// find access node
		searchHeap := &backwardSearchHeap{}
		heap.Init(searchHeap)
		heap.Push(searchHeap, &QueryVertex{
			id:               sourceVertex.id,
			forwardDistance:  0,
			backwardDistance: 0,
			isTransitNode:    sourceVertex.isTransitNode,
		})

		distance := make([]float64, len(graph.vertices), len(graph.vertices))
		for i := 0; i < len(graph.vertices); i++ {
			distance[i] = math.MaxFloat64
		}
		distance[sourceVertex.id] = 0

		for searchHeap.Len() != 0 {
			queryVertex := heap.Pop(searchHeap).(*QueryVertex)

			// relax
			if !queryVertex.isTransitNode {
				sourceVertex.backwardReachableVertex[queryVertex.id] = true
				for i := 0; i < len(graph.vertices[queryVertex.id].inwardEdges); i++ {
					inEdge := graph.vertices[queryVertex.id].inwardEdges[i]
					if graph.vertices[queryVertex.id].contractionOrder < inEdge.from.contractionOrder {
						if distance[inEdge.from.id] > distance[queryVertex.id]+inEdge.weight {
							distance[inEdge.from.id] = distance[queryVertex.id] + inEdge.weight
							heap.Push(searchHeap, &QueryVertex{
								id:               inEdge.from.id,
								forwardDistance:  0,
								backwardDistance: distance[queryVertex.id] + inEdge.weight,
								isTransitNode:    inEdge.from.isTransitNode,
							})
						}
					}
				}
			} else {
				sourceVertex.backwardAccessNodeDistance[queryVertex.id] = -1
			}
		}
		for k := range sourceVertex.backwardAccessNodeDistance {
			sourceVertex.backwardAccessNodeDistance[k], sourceVertex.backwardAccessNodePath[k] = graph.ShortestPathWithoutTNR(graph.vertices[k].name, sourceVertex.name)
		}

		// delete invalid access node
		accessNodeMask := make(map[int64]bool)
		for k1, d1 := range sourceVertex.backwardAccessNodeDistance {
			for k2, d2 := range sourceVertex.backwardAccessNodeDistance {
				if k1 == k2 {
					continue
				}
				if d1+graph.tnrDistance[k2][k1] <= d2 {
					accessNodeMask[k2] = true // mask j since it won't be the solution
				}
			}
		}
		for k := range accessNodeMask {
			delete(sourceVertex.backwardAccessNodeDistance, k)
			delete(sourceVertex.backwardAccessNodePath, k)
		}
	}
	graph.TNRed = true
}
