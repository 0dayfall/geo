package geo

import (
	"math"
	"math/rand"
)

// TSPResult contains the result of a TSP solution
type TSPResult struct {
	Tour     []int   // order of nodes to visit
	Distance float64 // total distance of the tour
}

// TSPNearestNeighbor solves the TSP using the nearest neighbor heuristic.
// distanceMatrix[i][j] represents the distance from node i to node j.
// Returns a tour starting from the specified start node.
func TSPNearestNeighbor(distanceMatrix [][]float64, start int) *TSPResult {
	n := len(distanceMatrix)
	if n == 0 || start < 0 || start >= n {
		return nil
	}

	visited := make([]bool, n)
	tour := []int{start}
	visited[start] = true
	totalDistance := 0.0
	current := start

	// Visit all nodes
	for len(tour) < n {
		nearest := -1
		minDist := math.Inf(1)

		// Find nearest unvisited neighbor
		for j := 0; j < n; j++ {
			if !visited[j] && distanceMatrix[current][j] < minDist {
				minDist = distanceMatrix[current][j]
				nearest = j
			}
		}

		if nearest == -1 {
			break
		}

		tour = append(tour, nearest)
		visited[nearest] = true
		totalDistance += minDist
		current = nearest
	}

	// Return to start
	if len(tour) == n {
		totalDistance += distanceMatrix[current][start]
	}

	return &TSPResult{
		Tour:     tour,
		Distance: totalDistance,
	}
}

// TSP2Opt improves a TSP tour using the 2-opt local search heuristic.
// This algorithm iteratively improves the tour by removing crossing edges.
func TSP2Opt(distanceMatrix [][]float64, initialTour []int, maxIterations int) *TSPResult {
	n := len(distanceMatrix)
	if n == 0 || len(initialTour) == 0 {
		return nil
	}

	tour := make([]int, len(initialTour))
	copy(tour, initialTour)

	// Calculate initial distance
	distance := calculateTourDistance(distanceMatrix, tour)

	improved := true
	iteration := 0

	for improved && (maxIterations <= 0 || iteration < maxIterations) {
		improved = false
		iteration++

		for i := 0; i < n-1; i++ {
			for j := i + 2; j < n; j++ {
				// Try swapping edges (i, i+1) and (j, j+1)
				// Calculate change in distance
				delta := -distanceMatrix[tour[i]][tour[i+1]] -
					distanceMatrix[tour[j]][tour[(j+1)%n]]
				delta += distanceMatrix[tour[i]][tour[j]] +
					distanceMatrix[tour[i+1]][tour[(j+1)%n]]

				if delta < -1e-10 { // improvement found
					// Reverse the segment between i+1 and j
					reverse(tour, i+1, j)
					distance += delta
					improved = true
				}
			}
		}
	}

	return &TSPResult{
		Tour:     tour,
		Distance: distance,
	}
}

// TSPSimulatedAnnealing solves TSP using simulated annealing metaheuristic.
// This is more robust for larger instances but slower.
func TSPSimulatedAnnealing(distanceMatrix [][]float64, start int, iterations int, temperature float64, coolingRate float64) *TSPResult {
	n := len(distanceMatrix)
	if n == 0 || start < 0 || start >= n {
		return nil
	}

	// Create initial tour using nearest neighbor
	current := TSPNearestNeighbor(distanceMatrix, start)
	if current == nil {
		return nil
	}

	best := &TSPResult{
		Tour:     make([]int, len(current.Tour)),
		Distance: current.Distance,
	}
	copy(best.Tour, current.Tour)

	temp := temperature
	rng := rand.New(rand.NewSource(42))

	for iter := 0; iter < iterations; iter++ {
		// Generate neighbor solution by swapping two random cities
		i := rng.Intn(n)
		j := rng.Intn(n)
		if i == j {
			continue
		}
		if i > j {
			i, j = j, i
		}

		// Create new tour by reversing segment
		newTour := make([]int, len(current.Tour))
		copy(newTour, current.Tour)
		reverse(newTour, i, j)

		newDistance := calculateTourDistance(distanceMatrix, newTour)
		delta := newDistance - current.Distance

		// Accept or reject the new solution
		if delta < 0 || rng.Float64() < math.Exp(-delta/temp) {
			current.Tour = newTour
			current.Distance = newDistance

			// Update best solution
			if newDistance < best.Distance {
				best.Tour = make([]int, len(newTour))
				copy(best.Tour, newTour)
				best.Distance = newDistance
			}
		}

		// Cool down
		temp *= coolingRate
	}

	return best
}

// calculateTourDistance computes the total distance of a tour
func calculateTourDistance(distanceMatrix [][]float64, tour []int) float64 {
	distance := 0.0
	for i := 0; i < len(tour)-1; i++ {
		distance += distanceMatrix[tour[i]][tour[i+1]]
	}
	// Return to start
	if len(tour) > 0 {
		distance += distanceMatrix[tour[len(tour)-1]][tour[0]]
	}
	return distance
}

// reverse reverses a segment of the tour between indices i and j (inclusive)
func reverse(tour []int, i, j int) {
	for i < j {
		tour[i], tour[j] = tour[j], tour[i]
		i++
		j--
	}
}
