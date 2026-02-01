package geo

import (
	"math"
)

const (
	// EarthRadiusKm is the Earth's radius in kilometers
	EarthRadiusKm = 6371.0
	// EarthRadiusMiles is the Earth's radius in miles
	EarthRadiusMiles = 3959.0
)

// toRadians converts degrees to radians
func toRadians(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// toDegrees converts radians to degrees
func toDegrees(rad float64) float64 {
	return rad * 180.0 / math.Pi
}

// GreatCircleDistance calculates the great circle distance between two points
// using the Haversine formula. Coordinates are in degrees (latitude, longitude).
// Returns distance in kilometers.
func GreatCircleDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert to radians
	φ1 := toRadians(lat1)
	φ2 := toRadians(lat2)
	Δφ := toRadians(lat2 - lat1)
	Δλ := toRadians(lon2 - lon1)

	// Haversine formula
	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*
			math.Sin(Δλ/2)*math.Sin(Δλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadiusKm * c
}

// RhumbLineDistance calculates the rhumb line (loxodrome) distance between two points.
// A rhumb line is a path of constant bearing. Coordinates are in degrees (latitude, longitude).
// Returns distance in kilometers.
func RhumbLineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert to radians
	φ1 := toRadians(lat1)
	φ2 := toRadians(lat2)
	Δφ := φ2 - φ1
	Δλ := toRadians(lon2 - lon1)

	// Handle crossing antimeridian
	if math.Abs(Δλ) > math.Pi {
		if Δλ > 0 {
			Δλ = -(2*math.Pi - Δλ)
		} else {
			Δλ = 2*math.Pi + Δλ
		}
	}

	// Calculate Δψ (distance along parallel)
	Δψ := math.Log(math.Tan(φ2/2+math.Pi/4) / math.Tan(φ1/2+math.Pi/4))

	// Handle case when E-W line (course of 90° or 270°)
	var q float64
	if math.Abs(Δψ) > 1e-12 {
		q = Δφ / Δψ
	} else {
		q = math.Cos(φ1)
	}

	// Distance in radians
	δ := math.Sqrt(Δφ*Δφ + q*q*Δλ*Δλ)

	return δ * EarthRadiusKm
}
