package geo

import (
	"math"
	"testing"
)

func TestNewGraph(t *testing.T) {
	g := NewGraph(5)
	if g.Nodes != 5 {
		t.Errorf("Expected 5 nodes, got %d", g.Nodes)
	}
	if len(g.Edges) != 5 {
		t.Errorf("Expected 5 edge lists, got %d", len(g.Edges))
	}
}

func TestAddEdge(t *testing.T) {
	g := NewGraph(3)
	g.AddEdge(0, 1, 5.0)
	g.AddEdge(0, 2, 3.0)

	if len(g.Edges[0]) != 2 {
		t.Errorf("Expected 2 edges from node 0, got %d", len(g.Edges[0]))
	}

	if g.Edges[0][0].To != 1 || g.Edges[0][0].Weight != 5.0 {
		t.Errorf("Edge 0->1 not correctly added")
	}
}

func TestAddBidirectionalEdge(t *testing.T) {
	g := NewGraph(2)
	g.AddBidirectionalEdge(0, 1, 5.0)

	if len(g.Edges[0]) != 1 || len(g.Edges[1]) != 1 {
		t.Errorf("Bidirectional edge not added correctly")
	}

	if g.Edges[0][0].To != 1 || g.Edges[1][0].To != 0 {
		t.Errorf("Bidirectional edge destinations incorrect")
	}
}

func TestDijkstraSimple(t *testing.T) {
	// Create a simple graph:
	//   0 --5--> 1
	//   |        |
	//   3        2
	//   |        |
	//   v        v
	//   2 --1--> 3
	g := NewGraph(4)
	g.AddEdge(0, 1, 5.0)
	g.AddEdge(0, 2, 3.0)
	g.AddEdge(1, 3, 2.0)
	g.AddEdge(2, 3, 1.0)

	result := g.Dijkstra(0)

	if result == nil {
		t.Fatal("Dijkstra returned nil")
	}

	// Check distances
	expectedDistances := []float64{0, 5, 3, 4}
	for i, expected := range expectedDistances {
		if math.Abs(result.Distances[i]-expected) > 1e-9 {
			t.Errorf("Distance to node %d = %v, want %v", i, result.Distances[i], expected)
		}
	}

	// Check path to node 3
	path := result.GetPath(3)
	expectedPath := []int{0, 2, 3}
	if !equalPath(path, expectedPath) {
		t.Errorf("Path to node 3 = %v, want %v", path, expectedPath)
	}
}

func TestDijkstraDisconnected(t *testing.T) {
	// Create a graph with disconnected components
	g := NewGraph(4)
	g.AddEdge(0, 1, 1.0)
	g.AddEdge(2, 3, 1.0)

	result := g.Dijkstra(0)

	// Nodes 2 and 3 should be unreachable from 0
	if !math.IsInf(result.Distances[2], 1) {
		t.Errorf("Node 2 should be unreachable, got distance %v", result.Distances[2])
	}
	if !math.IsInf(result.Distances[3], 1) {
		t.Errorf("Node 3 should be unreachable, got distance %v", result.Distances[3])
	}
}

func TestDijkstraComplexGraph(t *testing.T) {
	// More complex graph
	//     1
	//   /   \
	//  4     2
	//  |   / | \
	//  | 1   3  1
	//  |/    |   \
	//  0     5    3
	//   \   /
	//    2-4
	g := NewGraph(6)
	g.AddBidirectionalEdge(0, 1, 4.0)
	g.AddBidirectionalEdge(0, 2, 1.0)
	g.AddBidirectionalEdge(0, 4, 2.0)
	g.AddBidirectionalEdge(1, 3, 2.0)
	g.AddBidirectionalEdge(2, 3, 1.0)
	g.AddBidirectionalEdge(2, 5, 5.0)
	g.AddBidirectionalEdge(3, 5, 1.0)
	g.AddBidirectionalEdge(4, 5, 4.0)

	result := g.Dijkstra(0)

	// Check some specific distances
	tests := []struct {
		node     int
		expected float64
	}{
		{0, 0.0},
		{1, 4.0},
		{2, 1.0},
		{3, 2.0},
		{4, 2.0},
		{5, 3.0},
	}

	for _, tt := range tests {
		if math.Abs(result.Distances[tt.node]-tt.expected) > 1e-9 {
			t.Errorf("Distance to node %d = %v, want %v",
				tt.node, result.Distances[tt.node], tt.expected)
		}
	}
}

func TestGetPathNoPath(t *testing.T) {
	g := NewGraph(3)
	g.AddEdge(0, 1, 1.0)
	// Node 2 is disconnected

	result := g.Dijkstra(0)
	path := result.GetPath(2)

	if path != nil {
		t.Errorf("Expected nil path to unreachable node, got %v", path)
	}
}

func equalPath(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
