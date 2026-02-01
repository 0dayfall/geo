# geo

Geo is a lightweight toolkit for spatial routing and optimization. It provides great-circle distance, geohash indexing, graph search with Dijkstra, and heuristics for TSP. Designed for logistics, maritime and mapping use, it enables fast proximity queries, clustering, and path planning on real-world networks.

## Features

- **Distance Calculations**
  - Great Circle Distance (Haversine formula) - shortest distance on a sphere
  - Rhumb Line Distance - constant bearing path distance
  - Distance outputs in kilometers, meters, and nautical miles

- **Geohash**
  - Encode geographic coordinates into geohash strings
  - Decode geohash strings back to coordinates
  - Find neighboring geohashes

- **Graph Algorithms**
  - Dijkstra's shortest path algorithm
  - Weighted directed/undirected graphs

- **Route Optimization**
  - Traveling Salesman Problem (TSP) solver
  - Nearest Neighbor heuristic
  - 2-Opt local search improvement
  - Simulated Annealing metaheuristic

## Installation

```bash
go get github.com/0dayfall/geo
```

## Usage

### Distance Calculations

```go
import "github.com/0dayfall/geo"

// Great circle distance (shortest path on sphere)
distance := geo.GreatCircleDistance(40.7128, -74.0060, 51.5074, -0.1278)
fmt.Printf("NY to London: %.2f km\n", distance)

// Additional units
distanceMeters := geo.GreatCircleDistanceMeters(40.7128, -74.0060, 51.5074, -0.1278)
distanceNM := geo.GreatCircleDistanceNauticalMiles(40.7128, -74.0060, 51.5074, -0.1278)
fmt.Printf("NY to London: %.0f m (%.2f NM)\n", distanceMeters, distanceNM)

// Rhumb line distance (constant bearing)
rhumb := geo.RhumbLineDistance(40.7128, -74.0060, 51.5074, -0.1278)
fmt.Printf("Rhumb distance: %.2f km\n", rhumb)

// Additional units
rhumbMeters := geo.RhumbLineDistanceMeters(40.7128, -74.0060, 51.5074, -0.1278)
rhumbNM := geo.RhumbLineDistanceNauticalMiles(40.7128, -74.0060, 51.5074, -0.1278)
fmt.Printf("Rhumb distance: %.0f m (%.2f NM)\n", rhumbMeters, rhumbNM)

// Intermediate point along great circle (fraction 0..1)
midLat, midLon := geo.GreatCircleIntermediatePoint(
    40.7128, -74.0060,
    51.5074, -0.1278,
    0.5,
)
fmt.Printf("Midpoint: %.4f, %.4f\n", midLat, midLon)

// Point along path given speed (km/h) and duration (hours)
travelLat, travelLon := geo.GreatCirclePointAtSpeed(
    40.7128, -74.0060,
    51.5074, -0.1278,
    800.0, 2.0,
)
fmt.Printf("After 2h at 800 km/h: %.4f, %.4f\n", travelLat, travelLon)
```

### Geohash

```go
// Encode coordinates to geohash
hash := geo.Geohash(37.7749, -122.4194, 9) // San Francisco
fmt.Println(hash) // "9q8yyk8yt"

// Decode geohash back to coordinates
lat, lon, latErr, lonErr := geo.GeohashDecode(hash)

// Find neighboring geohashes
neighbors := geo.GeohashNeighbors(hash) // Returns [8]string
```

### Dijkstra's Algorithm

```go
// Create a graph with 4 nodes
graph := geo.NewGraph(4)
graph.AddBidirectionalEdge(0, 1, 5.0)
graph.AddEdge(0, 2, 3.0)
graph.AddEdge(2, 3, 1.0)

// Find shortest paths from node 0
result := graph.Dijkstra(0)
fmt.Println("Distance to node 3:", result.Distances[3])

// Get the actual path
path := result.GetPath(3)
fmt.Println("Path:", path) // [0 2 3]
```

### Traveling Salesman Problem

```go
// Distance matrix between 4 cities
distanceMatrix := [][]float64{
    {0, 10, 15, 20},
    {10, 0, 35, 25},
    {15, 35, 0, 30},
    {20, 25, 30, 0},
}

// Solve TSP using nearest neighbor
solution := geo.TSPNearestNeighbor(distanceMatrix, 0)
fmt.Println("Tour:", solution.Tour)
fmt.Println("Distance:", solution.Distance)

// Improve solution with 2-opt
improved := geo.TSP2Opt(distanceMatrix, solution.Tour, 100)
fmt.Println("Improved distance:", improved.Distance)

// Or use simulated annealing for larger problems
sa := geo.TSPSimulatedAnnealing(distanceMatrix, 0, 1000, 100.0, 0.95)
```

## Examples

See the [examples](examples/) directory for complete working examples.

Run the example:

```bash
cd examples
go run main.go
```

## Testing

```bash
go test ./...
```

## License

MIT
