package geo

import "testing"

var (
	sinkFloat float64
	sinkStr   string
	sinkSlice []int
)

func BenchmarkGreatCircleDistance(b *testing.B) {
	lat1, lon1 := 40.7128, -74.0060
	lat2, lon2 := 51.5074, -0.1278
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkFloat = GreatCircleDistance(lat1, lon1, lat2, lon2)
	}
}

func BenchmarkRhumbLineDistance(b *testing.B) {
	lat1, lon1 := 40.7128, -74.0060
	lat2, lon2 := 51.5074, -0.1278
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkFloat = RhumbLineDistance(lat1, lon1, lat2, lon2)
	}
}

func BenchmarkGreatCircleIntermediatePoint(b *testing.B) {
	lat1, lon1 := 34.0522, -118.2437
	lat2, lon2 := 51.5074, -0.1278
	fraction := 0.35
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		lat, lon := GreatCircleIntermediatePoint(lat1, lon1, lat2, lon2, fraction)
		sinkFloat = lat + lon
	}
}

func BenchmarkGreatCirclePointAtSpeed(b *testing.B) {
	lat1, lon1 := 34.0522, -118.2437
	lat2, lon2 := 51.5074, -0.1278
	speedKmh := 900.0
	durationHours := 2.5
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		lat, lon := GreatCirclePointAtSpeed(lat1, lon1, lat2, lon2, speedKmh, durationHours)
		sinkFloat = lat + lon
	}
}

func BenchmarkGreatCircleProject(b *testing.B) {
	lat1, lon1 := 0.0, 0.0
	lat2, lon2 := 0.0, 90.0
	latP, lonP := 10.0, 45.0
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		lat, lon, cross, along := GreatCircleProject(lat1, lon1, lat2, lon2, latP, lonP)
		sinkFloat = lat + lon + cross + along
	}
}

func BenchmarkGreatCircleProjectToSegment(b *testing.B) {
	lat1, lon1 := 0.0, 0.0
	lat2, lon2 := 0.0, 30.0
	latP, lonP := 0.0, 60.0
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		lat, lon, cross, along := GreatCircleProjectToSegment(lat1, lon1, lat2, lon2, latP, lonP)
		sinkFloat = lat + lon + cross + along
	}
}

func BenchmarkGeohashEncode(b *testing.B) {
	lat, lon := 37.7749, -122.4194
	precision := 9
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkStr = Geohash(lat, lon, precision)
	}
}

func BenchmarkGeohashDecode(b *testing.B) {
	hash := Geohash(37.7749, -122.4194, 9)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		lat, lon, latErr, lonErr := GeohashDecode(hash)
		sinkFloat = lat + lon + latErr + lonErr
	}
}

func BenchmarkGeohashNeighbors(b *testing.B) {
	hash := Geohash(37.7749, -122.4194, 9)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		neighbors := GeohashNeighbors(hash)
		sinkStr = neighbors[0]
	}
}

func BenchmarkDijkstra(b *testing.B) {
	const n = 1000
	graph := NewGraph(n)
	for i := 0; i < n-1; i++ {
		graph.AddBidirectionalEdge(i, i+1, 1.0)
	}
	graph.AddBidirectionalEdge(0, n-1, 1.0)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		result := graph.Dijkstra(0)
		sinkFloat = result.Distances[n-1]
	}
}

func BenchmarkTSPNearestNeighbor(b *testing.B) {
	coords := []struct{ lat, lon float64 }{
		{40.7128, -74.0060},
		{34.0522, -118.2437},
		{41.8781, -87.6298},
		{29.7604, -95.3698},
		{33.4484, -112.0740},
		{39.7392, -104.9903},
		{47.6062, -122.3321},
		{25.7617, -80.1918},
		{32.7767, -96.7970},
		{38.9072, -77.0369},
		{37.7749, -122.4194},
		{42.3601, -71.0589},
	}
	n := len(coords)
	matrix := make([][]float64, n)
	for i := range matrix {
		matrix[i] = make([]float64, n)
		for j := range matrix[i] {
			if i != j {
				matrix[i][j] = GreatCircleDistance(
					coords[i].lat, coords[i].lon,
					coords[j].lat, coords[j].lon,
				)
			}
		}
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		result := TSPNearestNeighbor(matrix, 0)
		sinkSlice = result.Tour
		sinkFloat = result.Distance
	}
}
