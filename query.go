package main

import (
	"container/heap"
	"container/list"
	"fmt"
	"math"
)

// ShortestPath 123
func (graph *Graph) ShortestPath(source int64, target int64) (float64, []int64) {
	if graph.TNRed {
		d, p := graph.ShortestPathWithTNR(source, target)
		if d == -1 {
			//fmt.Printf("%d to %d: local path or something wrong, fallback to CH\n", source, target)
			return graph.ShortestPathWithoutTNR(source, target)
		}
		return d, p
	}
	return graph.ShortestPathWithoutTNR(source, target)
}

// ShortestPathWithoutTNR 123
func (graph *Graph) ShortestPathWithoutTNR(source int64, target int64) (float64, []int64) {
	sourceVertex := graph.vertices[graph.mapping[source]]
	targetVertex := graph.vertices[graph.mapping[target]]

	forwardDistance := make([]float64, len(graph.vertices), len(graph.vertices))
	backwardDistance := make([]float64, len(graph.vertices), len(graph.vertices))
	forwardVisited := make([]bool, len(graph.vertices), len(graph.vertices))
	backwardVisited := make([]bool, len(graph.vertices), len(graph.vertices))

	for i := 0; i < len(graph.vertices); i++ {
		forwardDistance[i] = math.MaxFloat64
		backwardDistance[i] = math.MaxFloat64
	}

	forwardVisited[sourceVertex.id] = true
	forwardDistance[sourceVertex.id] = 0
	backwardVisited[targetVertex.id] = true
	backwardDistance[targetVertex.id] = 0

	forwardSearchHeap := &forwardSearchHeap{}
	backwardSearchHeap := &backwardSearchHeap{}
	heap.Init(forwardSearchHeap)
	heap.Init(backwardSearchHeap)

	heap.Push(forwardSearchHeap, &QueryVertex{
		id:               sourceVertex.id,
		forwardDistance:  0,
		backwardDistance: 0,
	})
	heap.Push(backwardSearchHeap, &QueryVertex{
		id:               targetVertex.id,
		forwardDistance:  0,
		backwardDistance: 0,
	})

	forwardBacktrace := make(map[int64]int64)
	backwardBacktrace := make(map[int64]int64)

	estimateDistance := math.MaxFloat64
	intersectVertex := int64(-1)

	for forwardSearchHeap.Len() != 0 || backwardSearchHeap.Len() != 0 {
		if forwardSearchHeap.Len() != 0 {
			queryVertex := heap.Pop(forwardSearchHeap).(*QueryVertex)
			if queryVertex.forwardDistance <= estimateDistance {
				forwardVisited[queryVertex.id] = true

				// relax
				for i := 0; i < len(graph.vertices[queryVertex.id].outwardEdges); i++ {
					outEdge := graph.vertices[queryVertex.id].outwardEdges[i]
					if graph.vertices[queryVertex.id].contractionOrder < outEdge.to.contractionOrder {
						if forwardDistance[outEdge.to.id] > forwardDistance[queryVertex.id]+outEdge.weight {
							forwardDistance[outEdge.to.id] = forwardDistance[queryVertex.id] + outEdge.weight
							forwardBacktrace[outEdge.to.id] = queryVertex.id
							heap.Push(forwardSearchHeap, &QueryVertex{
								id:               outEdge.to.id,
								forwardDistance:  forwardDistance[queryVertex.id] + outEdge.weight,
								backwardDistance: 0,
							})
						}
					}
				}
			}

			// check if is intersection point
			if backwardVisited[queryVertex.id] {
				if queryVertex.forwardDistance+backwardDistance[queryVertex.id] < estimateDistance {
					intersectVertex = queryVertex.id
					estimateDistance = queryVertex.forwardDistance + backwardDistance[queryVertex.id]
				}
			}
		}
		if backwardSearchHeap.Len() != 0 {
			queryVertex := heap.Pop(backwardSearchHeap).(*QueryVertex)
			if queryVertex.backwardDistance <= estimateDistance {
				backwardVisited[queryVertex.id] = true

				// relax
				for i := 0; i < len(graph.vertices[queryVertex.id].inwardEdges); i++ {
					inEdge := graph.vertices[queryVertex.id].inwardEdges[i]
					if graph.vertices[queryVertex.id].contractionOrder < inEdge.from.contractionOrder {
						if backwardDistance[inEdge.from.id] > backwardDistance[queryVertex.id]+inEdge.weight {
							backwardDistance[inEdge.from.id] = backwardDistance[queryVertex.id] + inEdge.weight
							backwardBacktrace[inEdge.from.id] = queryVertex.id
							heap.Push(backwardSearchHeap, &QueryVertex{
								id:               inEdge.from.id,
								forwardDistance:  0,
								backwardDistance: backwardDistance[queryVertex.id] + inEdge.weight,
							})
						}
					}
				}
			}

			// check if is intersection point
			if forwardVisited[queryVertex.id] {
				if queryVertex.backwardDistance+forwardDistance[queryVertex.id] < estimateDistance {
					intersectVertex = queryVertex.id
					estimateDistance = queryVertex.backwardDistance + forwardDistance[queryVertex.id]
				}
			}
		}
	}

	//fmt.Printf("intersection: %d\n", intersectVertex)
	//fmt.Printf("length: %f\n", estimateDistance)

	if estimateDistance == math.MaxFloat64 {
		return -math.MaxFloat64, nil
	}

	return estimateDistance, graph.RetrievePath(forwardBacktrace, backwardBacktrace, intersectVertex)
}

// RetrievePath 123
func (graph *Graph) RetrievePath(forwardBacktrace map[int64]int64, backwardBacktrace map[int64]int64, intersectVertex int64) []int64 {
	pathList := list.New()

	// pull from backtracing
	pathList.PushBack(intersectVertex)
	ok := true
	v := intersectVertex
	for {
		if v, ok = forwardBacktrace[v]; ok {
			pathList.PushFront(v)
		} else {
			break
		}
	}

	v = intersectVertex
	for {
		if v, ok = backwardBacktrace[v]; ok {
			pathList.PushBack(v)
		} else {
			break
		}
	}

	// expanding contraction
	allExpanded := false
	for !allExpanded {
		allExpanded = true
		for e := pathList.Front(); e.Next() != nil; e = e.Next() {
			if contractedVertex, ok := graph.contracts[e.Value.(int64)][e.Next().Value.(int64)]; ok {
				allExpanded = false
				pathList.InsertAfter(contractedVertex, e)
			}
		}
	}

	// list to slice
	var path []int64
	for e := pathList.Front(); e != nil; e = e.Next() {
		path = append(path, graph.vertices[e.Value.(int64)].name)
	}

	return path
}

// ShortestPathWithTNR 123
func (graph *Graph) ShortestPathWithTNR(source int64, target int64) (float64, []int64) {
	if !graph.TNRed {
		fmt.Println("not TNRed")
		return -1, nil
	}

	sourceVertex := graph.vertices[graph.mapping[source]]
	targetVertex := graph.vertices[graph.mapping[target]]

	// check if there are access nodes
	if len(sourceVertex.forwardAccessNodeDistance) == 0 || len(targetVertex.backwardAccessNodeDistance) == 0 {
		//fmt.Println("fallback to CH because no access node available")
		return -1, nil
	}

	// check if is local search
	for k := range sourceVertex.forwardReachableVertex {
		if _, ok := targetVertex.backwardReachableVertex[k]; ok {
			//fmt.Println("fallback to CH because local search")
			return -1, nil
		}
	}
	//fmt.Println("use TNR")

	// compute distance and path
	bestDistance := math.MaxFloat64
	var bestSourceAccessNode int64
	var bestTargetAccessNode int64
	bestPath := []int64{}

	for k1, d1 := range sourceVertex.forwardAccessNodeDistance {
		for k2, d2 := range targetVertex.backwardAccessNodeDistance {
			// two transit nodes are not reachable
			if graph.tnrDistance[k1][k2] == -math.MaxFloat64 {
				continue
			}
			if bestDistance > d1+graph.tnrDistance[k1][k2]+d2 {
				bestDistance = d1 + graph.tnrDistance[k1][k2] + d2
				bestSourceAccessNode = k1
				bestTargetAccessNode = k2
			}
		}
	}

	if bestDistance == math.MaxFloat64 {
		return -math.MaxFloat64, nil
	}

	pathFromSource := sourceVertex.forwardAccessNodePath[bestSourceAccessNode]
	pathToTarget := targetVertex.backwardAccessNodePath[bestTargetAccessNode]
	pathBetweenAccessNodes := graph.tnrPath[bestSourceAccessNode][bestTargetAccessNode]

	bestPath = []int64{}
	bestPath = append(bestPath, pathFromSource[:len(pathFromSource)-1]...)
	if bestSourceAccessNode != bestTargetAccessNode && len(pathBetweenAccessNodes) > 1 {
		bestPath = append(bestPath, pathBetweenAccessNodes[:len(pathBetweenAccessNodes)-1]...)
	}
	bestPath = append(bestPath, pathToTarget...)

	return bestDistance, bestPath
}
