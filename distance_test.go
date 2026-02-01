package geo

import (
	"math"
	"testing"
)

func TestGreatCircleDistance(t *testing.T) {
	tests := []struct {
		name     string
		lat1     float64
		lon1     float64
		lat2     float64
		lon2     float64
		expected float64
		epsilon  float64
	}{
		{
			name:     "New York to London",
			lat1:     40.7128,
			lon1:     -74.0060,
			lat2:     51.5074,
			lon2:     -0.1278,
			expected: 5570.0, // approximately 5570 km
			epsilon:  10.0,
		},
		{
			name:     "Same location",
			lat1:     0.0,
			lon1:     0.0,
			lat2:     0.0,
			lon2:     0.0,
			expected: 0.0,
			epsilon:  0.001,
		},
		{
			name:     "Equator half way around",
			lat1:     0.0,
			lon1:     0.0,
			lat2:     0.0,
			lon2:     180.0,
			expected: 20015.0, // approximately half Earth's circumference
			epsilon:  20.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GreatCircleDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			if math.Abs(result-tt.expected) > tt.epsilon {
				t.Errorf("GreatCircleDistance() = %v, want %v (±%v)", result, tt.expected, tt.epsilon)
			}
		})
	}
}

func TestRhumbLineDistance(t *testing.T) {
	tests := []struct {
		name     string
		lat1     float64
		lon1     float64
		lat2     float64
		lon2     float64
		expected float64
		epsilon  float64
	}{
		{
			name:     "New York to London (rhumb)",
			lat1:     40.7128,
			lon1:     -74.0060,
			lat2:     51.5074,
			lon2:     -0.1278,
			expected: 5794.0, // rhumb line is longer than great circle
			epsilon:  10.0,
		},
		{
			name:     "Same location",
			lat1:     0.0,
			lon1:     0.0,
			lat2:     0.0,
			lon2:     0.0,
			expected: 0.0,
			epsilon:  0.001,
		},
		{
			name:     "Along equator",
			lat1:     0.0,
			lon1:     0.0,
			lat2:     0.0,
			lon2:     90.0,
			expected: 10007.5, // quarter of equator
			epsilon:  10.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RhumbLineDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			if math.Abs(result-tt.expected) > tt.epsilon {
				t.Errorf("RhumbLineDistance() = %v, want %v (±%v)", result, tt.expected, tt.epsilon)
			}
		})
	}
}

func TestDistanceComparison(t *testing.T) {
	// Great circle should always be shorter than or equal to rhumb line
	testCases := []struct {
		lat1, lon1, lat2, lon2 float64
	}{
		{40.7128, -74.0060, 51.5074, -0.1278},    // NY to London
		{37.7749, -122.4194, 34.0522, -118.2437}, // SF to LA
		{-33.8688, 151.2093, -37.8136, 144.9631}, // Sydney to Melbourne
	}

	for _, tc := range testCases {
		gc := GreatCircleDistance(tc.lat1, tc.lon1, tc.lat2, tc.lon2)
		rhumb := RhumbLineDistance(tc.lat1, tc.lon1, tc.lat2, tc.lon2)

		if gc > rhumb {
			t.Errorf("Great circle distance (%v) should not be greater than rhumb line distance (%v)",
				gc, rhumb)
		}
	}
}
