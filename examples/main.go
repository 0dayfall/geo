package main

import (
	"fmt"
	"github.com/0dayfall/geo"
)

func main() {
	fmt.Println("=== Geo Library Examples ===\n")

	// 1. Great Circle Distance
	fmt.Println("1. Great Circle Distance (Haversine)")
	fmt.Println("   Calculating distance from New York to London:")
	nyLat, nyLon := 40.7128, -74.0060
	londonLat, londonLon := 51.5074, -0.1278
	gcDist := geo.GreatCircleDistance(nyLat, nyLon, londonLat, londonLon)
	fmt.Printf("   Distance: %.2f km\n\n", gcDist)

	// 2. Rhumb Line Distance
	fmt.Println("2. Rhumb Line Distance (constant bearing)")
	rhumbDist := geo.RhumbLineDistance(nyLat, nyLon, londonLat, londonLon)
	fmt.Printf("   Distance: %.2f km\n", rhumbDist)
	fmt.Printf("   Difference from great circle: %.2f km\n\n", rhumbDist-gcDist)

	// 3. Geohash
	fmt.Println("3. Geohash Encoding")
	lat, lon := 37.7749, -122.4194 // San Francisco
	hash := geo.Geohash(lat, lon, 9)
	fmt.Printf("   San Francisco (%.4f, %.4f)\n", lat, lon)
	fmt.Printf("   Geohash: %s\n", hash)

	decodedLat, decodedLon, _, _ := geo.GeohashDecode(hash)
	fmt.Printf("   Decoded: (%.4f, %.4f)\n\n", decodedLat, decodedLon)

	// 4. Geohash Neighbors
	fmt.Println("4. Geohash Neighbors")
	neighbors := geo.GeohashNeighbors(hash)
	fmt.Printf("   Neighbors of %s:\n", hash)
	directions := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	for i, neighbor := range neighbors {
		fmt.Printf("   %s: %s\n", directions[i], neighbor)
	}
	fmt.Println()

	// 5. Dijkstra's Algorithm
	fmt.Println("5. Dijkstra's Shortest Path")
	fmt.Println("   Creating a graph with 5 nodes...")
	graph := geo.NewGraph(5)
	graph.AddBidirectionalEdge(0, 1, 4.0)
	graph.AddBidirectionalEdge(0, 2, 1.0)
	graph.AddBidirectionalEdge(1, 3, 1.0)
	graph.AddBidirectionalEdge(2, 3, 5.0)
	graph.AddBidirectionalEdge(3, 4, 3.0)

	result := graph.Dijkstra(0)
	fmt.Printf("   Shortest distances from node 0:\n")
	for i, dist := range result.Distances {
		fmt.Printf("   To node %d: %.1f\n", i, dist)
	}
	
	path := result.GetPath(4)
	fmt.Printf("   Path from 0 to 4: %v\n\n", path)

	// 6. Traveling Salesman Problem
	fmt.Println("6. Traveling Salesman Problem (TSP)")
	fmt.Println("   Finding optimal tour for 4 cities...")
	
	// Create distance matrix
	distanceMatrix := [][]float64{
		{0, 10, 15, 20},
		{10, 0, 35, 25},
		{15, 35, 0, 30},
		{20, 25, 30, 0},
	}

	// Nearest Neighbor solution
	nn := geo.TSPNearestNeighbor(distanceMatrix, 0)
	fmt.Printf("   Nearest Neighbor tour: %v\n", nn.Tour)
	fmt.Printf("   Distance: %.1f\n", nn.Distance)

	// 2-opt improvement
	improved := geo.TSP2Opt(distanceMatrix, nn.Tour, 100)
	fmt.Printf("   2-Opt improved tour: %v\n", improved.Tour)
	fmt.Printf("   Distance: %.1f\n", improved.Distance)
	fmt.Printf("   Improvement: %.1f%%\n\n", (nn.Distance-improved.Distance)/nn.Distance*100)

	// 7. TSP with Geographic Locations
	fmt.Println("7. TSP with Real Geographic Locations")
	locations := []struct {
		name string
		lat  float64
		lon  float64
	}{
		{"New York", 40.7128, -74.0060},
		{"Los Angeles", 34.0522, -118.2437},
		{"Chicago", 41.8781, -87.6298},
		{"Houston", 29.7604, -95.3698},
		{"Phoenix", 33.4484, -112.0740},
	}

	// Build distance matrix using great circle distance
	n := len(locations)
	geoDistMatrix := make([][]float64, n)
	for i := range geoDistMatrix {
		geoDistMatrix[i] = make([]float64, n)
		for j := range geoDistMatrix[i] {
			if i != j {
				geoDistMatrix[i][j] = geo.GreatCircleDistance(
					locations[i].lat, locations[i].lon,
					locations[j].lat, locations[j].lon,
				)
			}
		}
	}

	tspResult := geo.TSPNearestNeighbor(geoDistMatrix, 0)
	fmt.Printf("   Optimal tour starting from %s:\n", locations[0].name)
	for i, cityIdx := range tspResult.Tour {
		fmt.Printf("   %d. %s\n", i+1, locations[cityIdx].name)
	}
	fmt.Printf("   Total distance: %.2f km\n", tspResult.Distance)

	// Improve with 2-opt
	tspImproved := geo.TSP2Opt(geoDistMatrix, tspResult.Tour, 100)
	fmt.Printf("   Improved distance: %.2f km\n", tspImproved.Distance)
}
