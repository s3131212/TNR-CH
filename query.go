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
		isTransitNode:    sourceVertex.isTransitNode,
	})
	heap.Push(backwardSearchHeap, &QueryVertex{
		id:               targetVertex.id,
		forwardDistance:  0,
		backwardDistance: 0,
		isTransitNode:    targetVertex.isTransitNode,
	})

	forwardBacktrace := make(map[int64]int64)
	backwardBacktrace := make(map[int64]int64)

	var forwardVisitedTransitNodes []int64
	var backwardVisitedTransitNodes []int64

	isLocal := false

	for forwardSearchHeap.Len() != 0 || backwardSearchHeap.Len() != 0 {
		if forwardSearchHeap.Len() != 0 {
			queryVertex := heap.Pop(forwardSearchHeap).(*QueryVertex)

			// record the visited transit nodes
			if queryVertex.isTransitNode && !forwardVisited[queryVertex.id] {
				forwardVisitedTransitNodes = append(forwardVisitedTransitNodes, queryVertex.id)
			}

			forwardVisited[queryVertex.id] = true

			// relax
			if !queryVertex.isTransitNode {
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
								isTransitNode:    outEdge.to.isTransitNode,
							})
						}
					}
				}
			}

			// check if two search merged
			if backwardVisited[queryVertex.id] && !queryVertex.isTransitNode {
				isLocal = true
				break
			}
		}
		if backwardSearchHeap.Len() != 0 {
			queryVertex := heap.Pop(backwardSearchHeap).(*QueryVertex)

			// record the visited transit nodes
			if queryVertex.isTransitNode && !backwardVisited[queryVertex.id] {
				backwardVisitedTransitNodes = append(backwardVisitedTransitNodes, queryVertex.id)
			}

			backwardVisited[queryVertex.id] = true

			// relax
			if !queryVertex.isTransitNode {
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
								isTransitNode:    inEdge.from.isTransitNode,
							})
						}
					}
				}
			}

			// check if two search merged
			if forwardVisited[queryVertex.id] && !queryVertex.isTransitNode {
				isLocal = true
				break
			}
		}
	}

	_ = isLocal
	// is local search, use ShortestPathWithoutTNR instead
	if isLocal || len(forwardVisitedTransitNodes) == 0 || len(backwardVisitedTransitNodes) == 0 {
		//fmt.Printf("local path (%d to %d) detected, use CH instead\n", source, target)
		return -1, nil
	}

	// find the suitable access node
	forwardVisitedTransitNodesMask := make([]bool, len(forwardVisitedTransitNodes), len(forwardVisitedTransitNodes))
	for i := 0; i < len(forwardVisitedTransitNodes); i++ {
		for j := 0; j < len(forwardVisitedTransitNodes); j++ {
			if i == j {
				continue
			}
			d1, _ := graph.ShortestPathWithoutTNR(sourceVertex.name, graph.vertices[forwardVisitedTransitNodes[i]].name)
			d2, _ := graph.ShortestPathWithoutTNR(sourceVertex.name, graph.vertices[forwardVisitedTransitNodes[j]].name)
			//fmt.Printf("deleting access node, comparing %d and %d: %f + %f <= %f\n", forwardVisitedTransitNodes[i], forwardVisitedTransitNodes[j], d1, graph.tnrDistance[forwardVisitedTransitNodes[i]][forwardVisitedTransitNodes[j]], d2)
			if d1+graph.tnrDistance[forwardVisitedTransitNodes[i]][forwardVisitedTransitNodes[j]] <= d2 {
				forwardVisitedTransitNodesMask[j] = true // mask j since it won't be the solution
			}
		}
	}
	//fmt.Printf("forwardVisitedTransitNodes for %d to %d: %v, %v\n", source, target, forwardVisitedTransitNodes, forwardVisitedTransitNodesMask)
	sourceAccessNodes := []int64{}
	for i := 0; i < len(forwardVisitedTransitNodes); i++ {
		if !forwardVisitedTransitNodesMask[i] {
			sourceAccessNodes = append(sourceAccessNodes, forwardVisitedTransitNodes[i])
		}
	}

	backwardVisitedTransitNodesMask := make([]bool, len(backwardVisitedTransitNodes), len(backwardVisitedTransitNodes))
	for i := 0; i < len(backwardVisitedTransitNodes); i++ {
		for j := i + 1; j < len(backwardVisitedTransitNodes); j++ {
			d1, _ := graph.ShortestPathWithoutTNR(graph.vertices[backwardVisitedTransitNodes[i]].name, targetVertex.name)
			d2, _ := graph.ShortestPathWithoutTNR(graph.vertices[backwardVisitedTransitNodes[j]].name, targetVertex.name)
			if d1+graph.tnrDistance[backwardVisitedTransitNodes[i]][backwardVisitedTransitNodes[j]] <= d2 {
				backwardVisitedTransitNodesMask[j] = true // mask j since it won't be the solution
			}
		}
	}
	//fmt.Printf("backwardVisitedTransitNodes for %d to %d: %v, %v\n", source, target, backwardVisitedTransitNodes, backwardVisitedTransitNodesMask)
	targetAccessNodes := []int64{}
	for i := 0; i < len(backwardVisitedTransitNodes); i++ {
		if !backwardVisitedTransitNodesMask[i] {
			targetAccessNodes = append(targetAccessNodes, backwardVisitedTransitNodes[i])
		}
	}

	//fmt.Printf("access node: %v, %v\n", sourceAccessNodes, targetAccessNodes)

	if len(sourceAccessNodes) == 0 || len(targetAccessNodes) == 0 {
		//fmt.Printf("failed to find access node for path (%d to %d), use CH instead\n", source, target)
		return -1, nil
	}

	// compute distance and path
	bestDistance := math.MaxFloat64
	bestPath := []int64{}
	for i := 0; i < len(sourceAccessNodes); i++ {
		for j := 0; j < len(targetAccessNodes); j++ {
			distanceFromSource, pathFromSource := graph.ShortestPathWithoutTNR(sourceVertex.name, graph.vertices[sourceAccessNodes[i]].name)
			distanceToTarget, pathToTarget := graph.ShortestPathWithoutTNR(graph.vertices[targetAccessNodes[j]].name, targetVertex.name)
			distanceBetweenAccessNodes := graph.tnrDistance[sourceAccessNodes[i]][targetAccessNodes[j]]
			pathBetweenAccessNodes := graph.tnrPath[sourceAccessNodes[i]][targetAccessNodes[j]]

			if bestDistance > distanceFromSource+distanceBetweenAccessNodes+distanceToTarget {
				bestDistance = distanceFromSource + distanceBetweenAccessNodes + distanceToTarget
				bestPath = []int64{}
				bestPath = append(bestPath, pathFromSource[:len(pathFromSource)-1]...)
				if sourceAccessNodes[i] != targetAccessNodes[j] {
					bestPath = append(bestPath, pathBetweenAccessNodes[:len(pathBetweenAccessNodes)-1]...)
				}
				bestPath = append(bestPath, pathToTarget...)
				//fmt.Printf("selected access node: %d, %d\n", sourceAccessNodes[i], targetAccessNodes[j])
			}
		}
	}

	if bestDistance == math.MaxFloat64 {
		return -math.MaxFloat64, nil
	}

	return bestDistance, bestPath
}
