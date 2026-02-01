package geo

import (
	"testing"
)

func TestGeohash(t *testing.T) {
	tests := []struct {
		name      string
		lat       float64
		lon       float64
		precision int
		expected  string
	}{
		{
			name:      "Eiffel Tower",
			lat:       48.8584,
			lon:       2.2945,
			precision: 9,
			expected:  "u09tunquc",
		},
		{
			name:      "Statue of Liberty",
			lat:       40.6892,
			lon:       -74.0445,
			precision: 9,
			expected:  "dr5r7p4ry",
		},
		{
			name:      "Sydney Opera House",
			lat:       -33.8568,
			lon:       151.2153,
			precision: 9,
			expected:  "r3gx2ux9g",
		},
		{
			name:      "Precision 5",
			lat:       37.7749,
			lon:       -122.4194,
			precision: 5,
			expected:  "9q8yy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Geohash(tt.lat, tt.lon, tt.precision)
			if result != tt.expected {
				t.Errorf("Geohash() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGeohashDecode(t *testing.T) {
	tests := []struct {
		name    string
		geohash string
		lat     float64
		lon     float64
		epsilon float64
	}{
		{
			name:    "Eiffel Tower area",
			geohash: "u09tunquc",
			lat:     48.8584,
			lon:     2.2945,
			epsilon: 0.001,
		},
		{
			name:    "Statue of Liberty area",
			geohash: "dr5r7p4ry",
			lat:     40.6892,
			lon:     -74.0445,
			epsilon: 0.001,
		},
		{
			name:    "Short geohash",
			geohash: "9q8yy",
			lat:     37.7749,
			lon:     -122.4194,
			epsilon: 0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lat, lon, _, _ := GeohashDecode(tt.geohash)
			if abs(lat-tt.lat) > tt.epsilon || abs(lon-tt.lon) > tt.epsilon {
				t.Errorf("GeohashDecode() = (%v, %v), want (~%v, ~%v)", lat, lon, tt.lat, tt.lon)
			}
		})
	}
}

func TestGeohashRoundTrip(t *testing.T) {
	testCases := []struct {
		lat       float64
		lon       float64
		precision int
		epsilon   float64
	}{
		{40.7128, -74.0060, 9, 0.001},
		{51.5074, -0.1278, 9, 0.001},
		{35.6762, 139.6503, 9, 0.001},
		{-33.8688, 151.2093, 9, 0.001},
		{0.0, 0.0, 9, 0.001},
	}

	for _, tc := range testCases {
		geohash := Geohash(tc.lat, tc.lon, tc.precision)
		lat, lon, _, _ := GeohashDecode(geohash)

		if abs(lat-tc.lat) > tc.epsilon || abs(lon-tc.lon) > tc.epsilon {
			t.Errorf("Round trip failed: input (%v, %v), got (%v, %v)",
				tc.lat, tc.lon, lat, lon)
		}
	}
}

func TestGeohashNeighbors(t *testing.T) {
	geohash := "9q8yy"
	neighbors := GeohashNeighbors(geohash)

	// Should get 8 neighbors
	if len(neighbors) != 8 {
		t.Errorf("Expected 8 neighbors, got %d", len(neighbors))
	}

	// Each neighbor should have the same precision
	for i, neighbor := range neighbors {
		if len(neighbor) != len(geohash) {
			t.Errorf("Neighbor %d has length %d, expected %d", i, len(neighbor), len(geohash))
		}
	}

	// Decode center and neighbors to verify they're actually neighbors
	centerLat, centerLon, _, _ := GeohashDecode(geohash)

	for i, neighbor := range neighbors {
		nLat, nLon, _, _ := GeohashDecode(neighbor)

		// Neighbors should be relatively close
		dist := GreatCircleDistance(centerLat, centerLon, nLat, nLon)
		if dist > 200 { // 200km seems reasonable for precision 5
			t.Errorf("Neighbor %d is too far away: %v km", i, dist)
		}
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
