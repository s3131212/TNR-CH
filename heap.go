package main

type minHeap []*Vertex

func (h minHeap) Len() int {
	return len(h)
}
func (h minHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *minHeap) Push(x interface{}) {
	*h = append(*h, x.(*Vertex))
}
func (h *minHeap) Pop() interface{} {
	n := len(*h)
	lastNode := (*h)[n-1]
	*h = (*h)[0 : n-1]
	return lastNode
}
func (h *minHeap) Peek() interface{} {
	return (*h)[0]
}

type importanceHeap struct {
	minHeap
}

func (h importanceHeap) Less(i int, j int) bool {
	return h.minHeap[i].importance < h.minHeap[j].importance
}

type distanceHeap struct {
	minHeap
}

func (h distanceHeap) Less(i int, j int) bool {
	return h.minHeap[i].distance.distance < h.minHeap[j].distance.distance
}

type minHeapQ []*QueryVertex

func (h minHeapQ) Len() int {
	return len(h)
}
func (h minHeapQ) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *minHeapQ) Push(x interface{}) {
	*h = append(*h, x.(*QueryVertex))
}
func (h *minHeapQ) Pop() interface{} {
	n := len(*h)
	lastNode := (*h)[n-1]
	*h = (*h)[0 : n-1]
	return lastNode
}
func (h *minHeapQ) Peek() interface{} {
	return (*h)[0]
}

type forwardSearchHeap struct {
	minHeapQ
}

func (h forwardSearchHeap) Less(i int, j int) bool {
	return h.minHeapQ[i].forwardDistance < h.minHeapQ[j].forwardDistance
}

type backwardSearchHeap struct {
	minHeapQ
}

func (h backwardSearchHeap) Less(i int, j int) bool {
	return h.minHeapQ[i].backwardDistance < h.minHeapQ[j].backwardDistance
}
