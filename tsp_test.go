package geo

import (
	"math"
	"testing"
)

func TestTSPNearestNeighbor(t *testing.T) {
	// Simple 4-city problem
	distanceMatrix := [][]float64{
		{0, 10, 15, 20},
		{10, 0, 35, 25},
		{15, 35, 0, 30},
		{20, 25, 30, 0},
	}

	result := TSPNearestNeighbor(distanceMatrix, 0)

	if result == nil {
		t.Fatal("TSPNearestNeighbor returned nil")
	}

	// Check that all cities are visited
	if len(result.Tour) != 4 {
		t.Errorf("Tour should visit 4 cities, got %d", len(result.Tour))
	}

	// Check that tour starts at the specified city
	if result.Tour[0] != 0 {
		t.Errorf("Tour should start at city 0, got %d", result.Tour[0])
	}

	// Check that all cities are unique
	visited := make(map[int]bool)
	for _, city := range result.Tour {
		if visited[city] {
			t.Errorf("City %d visited multiple times", city)
		}
		visited[city] = true
	}

	// Check that distance is positive
	if result.Distance <= 0 {
		t.Errorf("Distance should be positive, got %v", result.Distance)
	}
}

func TestTSP2Opt(t *testing.T) {
	// Create a simple distance matrix
	distanceMatrix := [][]float64{
		{0, 2, 9, 10},
		{2, 0, 6, 4},
		{9, 6, 0, 8},
		{10, 4, 8, 0},
	}

	// Start with a suboptimal tour
	initialTour := []int{0, 2, 1, 3}

	result := TSP2Opt(distanceMatrix, initialTour, 100)

	if result == nil {
		t.Fatal("TSP2Opt returned nil")
	}

	// Check tour length
	if len(result.Tour) != 4 {
		t.Errorf("Expected tour of length 4, got %d", len(result.Tour))
	}

	// The optimal tour for this matrix should be better than the initial
	initialDistance := calculateTourDistance(distanceMatrix, initialTour)
	if result.Distance > initialDistance {
		t.Errorf("2-opt should not increase distance: initial=%v, result=%v",
			initialDistance, result.Distance)
	}
}

func TestTSPSimulatedAnnealing(t *testing.T) {
	// Small symmetric distance matrix
	distanceMatrix := [][]float64{
		{0, 10, 15, 20},
		{10, 0, 35, 25},
		{15, 35, 0, 30},
		{20, 25, 30, 0},
	}

	result := TSPSimulatedAnnealing(distanceMatrix, 0, 1000, 100.0, 0.95)

	if result == nil {
		t.Fatal("TSPSimulatedAnnealing returned nil")
	}

	// Check tour validity
	if len(result.Tour) != 4 {
		t.Errorf("Tour should have 4 cities, got %d", len(result.Tour))
	}

	// Verify all cities are in the tour
	visited := make(map[int]bool)
	for _, city := range result.Tour {
		visited[city] = true
	}
	if len(visited) != 4 {
		t.Errorf("Tour should visit all 4 cities, visited %d", len(visited))
	}

	// Distance should be reasonable
	if result.Distance <= 0 || math.IsInf(result.Distance, 0) {
		t.Errorf("Invalid distance: %v", result.Distance)
	}
}

func TestCalculateTourDistance(t *testing.T) {
	distanceMatrix := [][]float64{
		{0, 1, 2, 3},
		{1, 0, 4, 5},
		{2, 4, 0, 6},
		{3, 5, 6, 0},
	}

	tour := []int{0, 1, 2, 3}
	// Distance: 0->1 (1) + 1->2 (4) + 2->3 (6) + 3->0 (3) = 14
	expected := 14.0

	result := calculateTourDistance(distanceMatrix, tour)

	if math.Abs(result-expected) > 1e-9 {
		t.Errorf("calculateTourDistance() = %v, want %v", result, expected)
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name     string
		tour     []int
		i        int
		j        int
		expected []int
	}{
		{
			name:     "Reverse middle segment",
			tour:     []int{0, 1, 2, 3, 4},
			i:        1,
			j:        3,
			expected: []int{0, 3, 2, 1, 4},
		},
		{
			name:     "Reverse entire tour",
			tour:     []int{0, 1, 2, 3},
			i:        0,
			j:        3,
			expected: []int{3, 2, 1, 0},
		},
		{
			name:     "Reverse two elements",
			tour:     []int{0, 1, 2, 3},
			i:        1,
			j:        2,
			expected: []int{0, 2, 1, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tourCopy := make([]int, len(tt.tour))
			copy(tourCopy, tt.tour)
			reverse(tourCopy, tt.i, tt.j)

			if !equalIntSlice(tourCopy, tt.expected) {
				t.Errorf("reverse() = %v, want %v", tourCopy, tt.expected)
			}
		})
	}
}

func TestTSPWithGeographicDistances(t *testing.T) {
	// Test with actual geographic coordinates
	locations := []struct {
		name string
		lat  float64
		lon  float64
	}{
		{"New York", 40.7128, -74.0060},
		{"Los Angeles", 34.0522, -118.2437},
		{"Chicago", 41.8781, -87.6298},
		{"Houston", 29.7604, -95.3698},
	}

	// Build distance matrix using great circle distance
	n := len(locations)
	distanceMatrix := make([][]float64, n)
	for i := range distanceMatrix {
		distanceMatrix[i] = make([]float64, n)
		for j := range distanceMatrix[i] {
			if i != j {
				distanceMatrix[i][j] = GreatCircleDistance(
					locations[i].lat, locations[i].lon,
					locations[j].lat, locations[j].lon,
				)
			}
		}
	}

	// Test nearest neighbor
	result := TSPNearestNeighbor(distanceMatrix, 0)
	if result == nil {
		t.Fatal("TSP failed with geographic distances")
	}

	// Verify tour is valid
	if len(result.Tour) != n {
		t.Errorf("Expected tour of length %d, got %d", n, len(result.Tour))
	}

	// Improve with 2-opt
	improved := TSP2Opt(distanceMatrix, result.Tour, 100)
	if improved.Distance > result.Distance {
		t.Errorf("2-opt should not worsen the solution")
	}
}

func equalIntSlice(a, b []int) bool {
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
