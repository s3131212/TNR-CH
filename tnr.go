package main

import "fmt"

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

			distance, path := graph.ShortestPath(graph.transitNodes[i].name, graph.transitNodes[j].name)

			graph.tnrDistance[graph.transitNodes[i].id][graph.transitNodes[j].id] = distance
			graph.tnrPath[graph.transitNodes[i].id][graph.transitNodes[j].id] = path

			//fmt.Printf("between transit nodes %d and %d: %f, %v\n", graph.transitNodes[i].id, graph.transitNodes[j].id, distance, path)
		}
	}
	graph.TNRed = true
}
