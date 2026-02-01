package geo

import (
	"math"
)

const (
	// MetersPerKm converts kilometers to meters
	MetersPerKm = 1000.0
	// KmPerNauticalMile converts nautical miles to kilometers
	KmPerNauticalMile = 1.852

	// EarthRadiusKm is the Earth's radius in kilometers
	EarthRadiusKm = 6371.0
	// EarthRadiusMiles is the Earth's radius in miles
	EarthRadiusMiles = 3959.0
	// EarthRadiusMeters is the Earth's radius in meters
	EarthRadiusMeters = EarthRadiusKm * MetersPerKm
	// EarthRadiusNauticalMiles is the Earth's radius in nautical miles
	EarthRadiusNauticalMiles = EarthRadiusKm / KmPerNauticalMile
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

// GreatCircleDistanceMeters returns the great circle distance in meters.
func GreatCircleDistanceMeters(lat1, lon1, lat2, lon2 float64) float64 {
	return GreatCircleDistance(lat1, lon1, lat2, lon2) * MetersPerKm
}

// GreatCircleDistanceNauticalMiles returns the great circle distance in nautical miles.
func GreatCircleDistanceNauticalMiles(lat1, lon1, lat2, lon2 float64) float64 {
	return GreatCircleDistance(lat1, lon1, lat2, lon2) / KmPerNauticalMile
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

// RhumbLineDistanceMeters returns the rhumb line distance in meters.
func RhumbLineDistanceMeters(lat1, lon1, lat2, lon2 float64) float64 {
	return RhumbLineDistance(lat1, lon1, lat2, lon2) * MetersPerKm
}

// RhumbLineDistanceNauticalMiles returns the rhumb line distance in nautical miles.
func RhumbLineDistanceNauticalMiles(lat1, lon1, lat2, lon2 float64) float64 {
	return RhumbLineDistance(lat1, lon1, lat2, lon2) / KmPerNauticalMile
}
