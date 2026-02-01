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

func TestGreatCircleDistanceConversions(t *testing.T) {
	lat1, lon1 := 40.7128, -74.0060
	lat2, lon2 := 51.5074, -0.1278

	km := GreatCircleDistance(lat1, lon1, lat2, lon2)
	meters := GreatCircleDistanceMeters(lat1, lon1, lat2, lon2)
	nm := GreatCircleDistanceNauticalMiles(lat1, lon1, lat2, lon2)

	if math.Abs(meters-km*MetersPerKm) > 1e-6 {
		t.Errorf("GreatCircleDistanceMeters() = %v, want %v", meters, km*MetersPerKm)
	}
	if math.Abs(nm-km/KmPerNauticalMile) > 1e-6 {
		t.Errorf("GreatCircleDistanceNauticalMiles() = %v, want %v", nm, km/KmPerNauticalMile)
	}
}

func TestRhumbLineDistanceConversions(t *testing.T) {
	lat1, lon1 := 40.7128, -74.0060
	lat2, lon2 := 51.5074, -0.1278

	km := RhumbLineDistance(lat1, lon1, lat2, lon2)
	meters := RhumbLineDistanceMeters(lat1, lon1, lat2, lon2)
	nm := RhumbLineDistanceNauticalMiles(lat1, lon1, lat2, lon2)

	if math.Abs(meters-km*MetersPerKm) > 1e-6 {
		t.Errorf("RhumbLineDistanceMeters() = %v, want %v", meters, km*MetersPerKm)
	}
	if math.Abs(nm-km/KmPerNauticalMile) > 1e-6 {
		t.Errorf("RhumbLineDistanceNauticalMiles() = %v, want %v", nm, km/KmPerNauticalMile)
	}
}

func TestGreatCircleIntermediatePoint(t *testing.T) {
	t.Run("fraction endpoints", func(t *testing.T) {
		lat1, lon1 := 10.0, -20.0
		lat2, lon2 := -5.0, 40.0

		latStart, lonStart := GreatCircleIntermediatePoint(lat1, lon1, lat2, lon2, 0.0)
		if math.Abs(latStart-lat1) > 1e-9 || math.Abs(lonStart-lon1) > 1e-9 {
			t.Errorf("fraction 0 = (%v, %v), want (%v, %v)", latStart, lonStart, lat1, lon1)
		}

		latEnd, lonEnd := GreatCircleIntermediatePoint(lat1, lon1, lat2, lon2, 1.0)
		if math.Abs(latEnd-lat2) > 1e-9 || math.Abs(lonEnd-lon2) > 1e-9 {
			t.Errorf("fraction 1 = (%v, %v), want (%v, %v)", latEnd, lonEnd, lat2, lon2)
		}
	})

	t.Run("equator midpoint", func(t *testing.T) {
		lat, lon := GreatCircleIntermediatePoint(0.0, 0.0, 0.0, 90.0, 0.5)
		if math.Abs(lat-0.0) > 1e-9 || math.Abs(lon-45.0) > 1e-9 {
			t.Errorf("midpoint = (%v, %v), want (0, 45)", lat, lon)
		}
	})

	t.Run("crosses equator", func(t *testing.T) {
		lat, lon := GreatCircleIntermediatePoint(-10.0, 0.0, 10.0, 0.0, 0.5)
		if math.Abs(lat-0.0) > 1e-9 || math.Abs(lon-0.0) > 1e-9 {
			t.Errorf("midpoint = (%v, %v), want (0, 0)", lat, lon)
		}
	})

	t.Run("distance fraction", func(t *testing.T) {
		lat1, lon1 := 34.0522, -118.2437
		lat2, lon2 := 51.5074, -0.1278
		f := 0.25

		latMid, lonMid := GreatCircleIntermediatePoint(lat1, lon1, lat2, lon2, f)
		total := GreatCircleDistance(lat1, lon1, lat2, lon2)
		part := GreatCircleDistance(lat1, lon1, latMid, lonMid)
		if math.Abs(part-total*f) > 1e-6*total {
			t.Errorf("distance fraction = %v, want %v", part/total, f)
		}
	})
}

func TestGreatCirclePointAtDistance(t *testing.T) {
	t.Run("equator half distance", func(t *testing.T) {
		lat1, lon1 := 0.0, 0.0
		lat2, lon2 := 0.0, 90.0
		total := GreatCircleDistance(lat1, lon1, lat2, lon2)

		lat, lon := GreatCirclePointAtDistance(lat1, lon1, lat2, lon2, total/2)
		if math.Abs(lat-0.0) > 1e-9 || math.Abs(lon-45.0) > 1e-6 {
			t.Errorf("point = (%v, %v), want (0, 45)", lat, lon)
		}
	})

	t.Run("clamps to endpoints", func(t *testing.T) {
		lat1, lon1 := 10.0, -20.0
		lat2, lon2 := -5.0, 40.0
		total := GreatCircleDistance(lat1, lon1, lat2, lon2)

		latStart, lonStart := GreatCirclePointAtDistance(lat1, lon1, lat2, lon2, -10.0)
		if math.Abs(latStart-lat1) > 1e-9 || math.Abs(lonStart-lon1) > 1e-9 {
			t.Errorf("start clamp = (%v, %v), want (%v, %v)", latStart, lonStart, lat1, lon1)
		}

		latEnd, lonEnd := GreatCirclePointAtDistance(lat1, lon1, lat2, lon2, total*2)
		if math.Abs(latEnd-lat2) > 1e-9 || math.Abs(lonEnd-lon2) > 1e-9 {
			t.Errorf("end clamp = (%v, %v), want (%v, %v)", latEnd, lonEnd, lat2, lon2)
		}
	})
}

func TestGreatCirclePointAtSpeed(t *testing.T) {
	lat1, lon1 := 34.0522, -118.2437
	lat2, lon2 := 51.5074, -0.1278
	speedKmh := 900.0
	durationHours := 2.5

	latSpeed, lonSpeed := GreatCirclePointAtSpeed(lat1, lon1, lat2, lon2, speedKmh, durationHours)
	latDist, lonDist := GreatCirclePointAtDistance(lat1, lon1, lat2, lon2, speedKmh*durationHours)

	if math.Abs(latSpeed-latDist) > 1e-9 || math.Abs(lonSpeed-lonDist) > 1e-9 {
		t.Errorf("point at speed = (%v, %v), want (%v, %v)", latSpeed, lonSpeed, latDist, lonDist)
	}
}

func TestGreatCircleProject(t *testing.T) {
	t.Run("equator projection", func(t *testing.T) {
		lat1, lon1 := 0.0, 0.0
		lat2, lon2 := 0.0, 90.0
		latP, lonP := 10.0, 45.0

		projLat, projLon, crossTrackKm, alongTrackKm := GreatCircleProject(lat1, lon1, lat2, lon2, latP, lonP)
		if math.Abs(projLat-0.0) > 1e-9 || math.Abs(projLon-45.0) > 1e-6 {
			t.Errorf("projected = (%v, %v), want (0, 45)", projLat, projLon)
		}

		expectedCross := EarthRadiusKm * toRadians(10.0)
		if math.Abs(math.Abs(crossTrackKm)-expectedCross) > 1e-3 {
			t.Errorf("cross-track = %v, want ±%v", crossTrackKm, expectedCross)
		}

		expectedAlong := EarthRadiusKm * toRadians(45.0)
		if math.Abs(alongTrackKm-expectedAlong) > 1e-3 {
			t.Errorf("along-track = %v, want %v", alongTrackKm, expectedAlong)
		}
	})

	t.Run("outside segment", func(t *testing.T) {
		lat1, lon1 := 0.0, 0.0
		lat2, lon2 := 0.0, 30.0
		latP, lonP := 0.0, 60.0

		_, _, crossTrackKm, alongTrackKm := GreatCircleProject(lat1, lon1, lat2, lon2, latP, lonP)
		if math.Abs(crossTrackKm) > 1e-9 {
			t.Errorf("cross-track = %v, want 0", crossTrackKm)
		}
		if alongTrackKm <= GreatCircleDistance(lat1, lon1, lat2, lon2) {
			t.Errorf("along-track = %v, want > total", alongTrackKm)
		}
	})
}
