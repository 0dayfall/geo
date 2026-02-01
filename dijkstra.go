package geo

import (
	"container/heap"
	"math"
)

// Edge represents a weighted edge in a graph
type Edge struct {
	To     int     // destination node
	Weight float64 // edge weight (distance, cost, etc.)
}

// Graph represents a weighted directed graph
type Graph struct {
	Nodes int      // number of nodes
	Edges [][]Edge // adjacency list
}

// NewGraph creates a new graph with the specified number of nodes
func NewGraph(nodes int) *Graph {
	return &Graph{
		Nodes: nodes,
		Edges: make([][]Edge, nodes),
	}
}

// AddEdge adds a directed edge from 'from' to 'to' with the given weight
func (g *Graph) AddEdge(from, to int, weight float64) {
	g.Edges[from] = append(g.Edges[from], Edge{To: to, Weight: weight})
}

// AddBidirectionalEdge adds edges in both directions
func (g *Graph) AddBidirectionalEdge(from, to int, weight float64) {
	g.AddEdge(from, to, weight)
	g.AddEdge(to, from, weight)
}

// DijkstraResult contains the results of Dijkstra's algorithm
type DijkstraResult struct {
	Distances []float64 // shortest distances from source
	Previous  []int     // previous node in shortest path (-1 if none)
}

// priorityQueueItem represents an item in the priority queue
type priorityQueueItem struct {
	node     int
	distance float64
	index    int
}

// priorityQueue implements heap.Interface
type priorityQueue []*priorityQueueItem

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].distance < pq[j].distance
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*priorityQueueItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// Dijkstra computes the shortest paths from a source node to all other nodes
// using Dijkstra's algorithm.
func (g *Graph) Dijkstra(source int) *DijkstraResult {
	if source < 0 || source >= g.Nodes {
		return nil
	}

	// Initialize distances and previous nodes
	distances := make([]float64, g.Nodes)
	previous := make([]int, g.Nodes)
	for i := range distances {
		distances[i] = math.Inf(1)
		previous[i] = -1
	}
	distances[source] = 0

	// Initialize priority queue
	pq := make(priorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &priorityQueueItem{
		node:     source,
		distance: 0,
	})

	visited := make([]bool, g.Nodes)

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*priorityQueueItem)
		u := current.node

		if visited[u] {
			continue
		}
		visited[u] = true

		// Explore neighbors
		for _, edge := range g.Edges[u] {
			v := edge.To
			if visited[v] {
				continue
			}

			alt := distances[u] + edge.Weight
			if alt < distances[v] {
				distances[v] = alt
				previous[v] = u
				heap.Push(&pq, &priorityQueueItem{
					node:     v,
					distance: alt,
				})
			}
		}
	}

	return &DijkstraResult{
		Distances: distances,
		Previous:  previous,
	}
}

// GetPath reconstructs the shortest path from source to target
func (r *DijkstraResult) GetPath(target int) []int {
	// Check if target is unreachable (infinite distance)
	if math.IsInf(r.Distances[target], 1) {
		return nil // no path exists
	}

	path := []int{}
	for u := target; u != -1; u = r.Previous[u] {
		path = append([]int{u}, path...)
		if r.Previous[u] == -1 {
			break
		}
	}

	return path
}
