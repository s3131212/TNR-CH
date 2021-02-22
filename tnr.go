package main

import (
	"container/heap"
	"fmt"
	"math"
)

// ComputeTNR Compute Transit Node Routing
func (graph *Graph) ComputeTNR(transitCnt int) {
	if !graph.contracted {
		fmt.Println("The graph has not contracted, run ComputeContractions first.")
		return
	}
	if graph.TNRed {
		fmt.Println("The graph has already calculated the TNR.")
		return
	}
	if len(graph.vertices) < transitCnt {
		fmt.Println("Too many transit nodes")
		return
	}

	graph.SelectTransitNodes(transitCnt)
	graph.ComputeDistanceTable(transitCnt)
	graph.ComputeVoronoiRegion()
	graph.ComputeLocalFilter()

	graph.TNRed = true
}

// SelectTransitNodes Select the Transit Nodes by contraction orders
func (graph *Graph) SelectTransitNodes(transitCnt int) {
	vertexCnt := len(graph.vertices)
	for i := 0; i < vertexCnt; i++ {
		//fmt.Printf("id: %d, contractionOrder: %d\n", graph.vertices[i].id, graph.vertices[i].contractionOrder)
		if graph.vertices[i].contractionOrder >= vertexCnt-transitCnt {
			graph.vertices[i].isTransitNode = true
			graph.transitNodes = append(graph.transitNodes, graph.vertices[i])
			// fmt.Printf("select transit %d\n", graph.vertices[i].id)
		}
		graph.vertices[i].transitPath = make(map[int64]*Vertex)
	}
}

// ComputeDistanceTable Select the transit nodes and compute the Distance Table
func (graph *Graph) ComputeDistanceTable(transitCnt int) {

	if graph.tnrDistance == nil {
		graph.tnrDistance = make(map[int64]map[int64]float64)
	}

	for i := 0; i < transitCnt; i++ {
		for j := 0; j < transitCnt; j++ {
			if _, ok := graph.tnrDistance[graph.transitNodes[i].id]; !ok {
				graph.tnrDistance[graph.transitNodes[i].id] = make(map[int64]float64)
			}

			if i == j {
				graph.tnrDistance[graph.transitNodes[i].id][graph.transitNodes[j].id] = 0
				continue
			}

			distance, path := graph.ShortestPathWithoutTNR(graph.transitNodes[i].name, graph.transitNodes[j].name)

			graph.tnrDistance[graph.transitNodes[i].id][graph.transitNodes[j].id] = distance

			for k := 0; k < len(path)-1; k++ {
				if _, ok := graph.vertices[graph.mapping[path[k]]].transitPath[graph.transitNodes[j].id]; ok {
					break
				}
				graph.vertices[graph.mapping[path[k]]].transitPath[graph.transitNodes[j].id] = graph.vertices[graph.mapping[path[k+1]]]
			}

			//fmt.Printf("between transit nodes %d and %d: %f, %v\n", graph.transitNodes[i].id, graph.transitNodes[j].id, distance, path)
		}
	}
}

// ComputeVoronoiRegion 123
func (graph *Graph) ComputeVoronoiRegion() {
	for i := 0; i < len(graph.vertices); i++ {
		graph.vertices[i].distance.distance = math.MaxFloat64
		graph.vertices[i].voronoiRegionID = -1
	}

	distanceHeap := &distanceHeap{}
	visited := make(map[int64]bool)

	heap.Init(distanceHeap)
	for i := 0; i < len(graph.transitNodes); i++ {
		graph.transitNodes[i].distance.distance = 0
		graph.transitNodes[i].voronoiRegionID = graph.transitNodes[i].id
		heap.Push(distanceHeap, graph.transitNodes[i])
	}

	for distanceHeap.Len() != 0 {
		vertex := heap.Pop(distanceHeap).(*Vertex)

		if visited[vertex.id] {
			continue
		}
		visited[vertex.id] = true

		for i := 0; i < len(vertex.inwardEdges); i++ {
			if !vertex.inwardEdges[i].isShortcut {
				if vertex.distance.distance+vertex.inwardEdges[i].weight < vertex.inwardEdges[i].from.distance.distance {
					vertex.inwardEdges[i].from.distance.distance = vertex.distance.distance + vertex.inwardEdges[i].weight
					heap.Push(distanceHeap, vertex.inwardEdges[i].from)
					vertex.inwardEdges[i].from.voronoiRegionID = vertex.voronoiRegionID
				}
			}
		}
	}
	/*
		for i := 0; i < len(graph.vertices); i++ {
			fmt.Printf("vertex %d is assigned to voronoi region %d\n", graph.vertices[i].id, graph.vertices[i].voronoiRegionID)
		}
	*/
}

// ComputeLocalFilter Calculate the local filter (access nodes + sub-transit-node sets)
func (graph *Graph) ComputeLocalFilter() {
	contractionMaxHeap := &contractionMaxHeap{}
	heap.Init(contractionMaxHeap)
	for v := 0; v < len(graph.vertices); v++ {
		graph.vertices[v].forwardSearchSpace = nil
		graph.vertices[v].forwardAccessNodeDistance = nil
		graph.vertices[v].forwardTNRed = false
		graph.vertices[v].backwardSearchSpace = nil
		graph.vertices[v].backwardAccessNodeDistance = nil
		graph.vertices[v].backwardTNRed = false

		heap.Push(contractionMaxHeap, graph.vertices[v])
	}

	for contractionMaxHeap.Len() != 0 {
		sourceVertex := heap.Pop(contractionMaxHeap).(*Vertex)
		if !sourceVertex.forwardTNRed {
			//fmt.Printf("forward TNR vertex %d\n", sourceVertex.id)
			sourceVertex.forwardSearchSpace = make(map[int64]bool)
			sourceVertex.forwardAccessNodeDistance = make(map[int64]float64)

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
					sourceVertex.forwardSearchSpace[graph.vertices[queryVertex.id].voronoiRegionID] = true

					// check if visited this node
					if graph.vertices[queryVertex.id].forwardTNRed {
						//fmt.Printf("met forward-TNRed vertex %d\n", queryVertex.id)
						for k := range graph.vertices[queryVertex.id].forwardSearchSpace {
							sourceVertex.forwardSearchSpace[k] = true
						}
						for k := range graph.vertices[queryVertex.id].forwardAccessNodeDistance {
							sourceVertex.forwardAccessNodeDistance[k] = -1
						}
					} else {
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
					}
				} else {
					sourceVertex.forwardAccessNodeDistance[queryVertex.id] = -1
				}
			}

			for k := range sourceVertex.forwardAccessNodeDistance {
				sourceVertex.forwardAccessNodeDistance[k], _ = graph.ShortestPathWithoutTNR(sourceVertex.name, graph.vertices[k].name)
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
			}
			sourceVertex.forwardTNRed = true
		}

		if !sourceVertex.backwardTNRed {
			//fmt.Printf("backward TNR vertex %d\n", sourceVertex.id)
			sourceVertex.backwardSearchSpace = make(map[int64]bool)
			sourceVertex.backwardAccessNodeDistance = make(map[int64]float64)

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
					sourceVertex.backwardSearchSpace[graph.vertices[queryVertex.id].voronoiRegionID] = true

					// check if visited this node
					if graph.vertices[queryVertex.id].backwardTNRed {
						//fmt.Printf("met backward-TNRed vertex %d\n", queryVertex.id)
						for k := range graph.vertices[queryVertex.id].backwardSearchSpace {
							sourceVertex.backwardSearchSpace[k] = true
						}
						for k := range graph.vertices[queryVertex.id].backwardAccessNodeDistance {
							sourceVertex.backwardAccessNodeDistance[k] = -1
						}
					} else {
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
					}
				} else {
					sourceVertex.backwardAccessNodeDistance[queryVertex.id] = -1
				}
			}
			for k := range sourceVertex.backwardAccessNodeDistance {
				sourceVertex.backwardAccessNodeDistance[k], _ = graph.ShortestPathWithoutTNR(graph.vertices[k].name, sourceVertex.name)
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
			}
		}
		sourceVertex.backwardTNRed = true
	}
}
