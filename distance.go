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

// normalizeLongitude keeps longitude in the [-180, 180] range.
func normalizeLongitude(lon float64) float64 {
	lon = math.Mod(lon+180.0, 360.0)
	if lon < 0 {
		lon += 360.0
	}
	return lon - 180.0
}

// initialBearingRad returns the initial bearing from point 1 to point 2 in radians.
func initialBearingRad(lat1, lon1, lat2, lon2 float64) float64 {
	φ1 := toRadians(lat1)
	φ2 := toRadians(lat2)
	Δλ := toRadians(lon2 - lon1)

	y := math.Sin(Δλ) * math.Cos(φ2)
	x := math.Cos(φ1)*math.Sin(φ2) - math.Sin(φ1)*math.Cos(φ2)*math.Cos(Δλ)
	return math.Atan2(y, x)
}

// angularDistanceRad returns the central angle between two points in radians.
func angularDistanceRad(lat1, lon1, lat2, lon2 float64) float64 {
	φ1 := toRadians(lat1)
	φ2 := toRadians(lat2)
	Δφ := φ2 - φ1
	Δλ := toRadians(lon2 - lon1)

	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*
			math.Sin(Δλ/2)*math.Sin(Δλ/2)
	return 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
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

// GreatCircleProject projects a point onto the great circle path between two coordinates.
// Returns the projected point (lat, lon), cross-track distance (km), and along-track
// distance from the start (km). Along-track can be negative or exceed total distance,
// indicating the perpendicular projection falls outside the segment.
// Cross-track is signed: positive means the point is to the right of the path from
// start to end, negative to the left.
func GreatCircleProject(lat1, lon1, lat2, lon2, latP, lonP float64) (float64, float64, float64, float64) {
	totalAngle := angularDistanceRad(lat1, lon1, lat2, lon2)
	if totalAngle == 0 {
		return lat1, normalizeLongitude(lon1),
			GreatCircleDistance(lat1, lon1, latP, lonP),
			0
	}

	δ13 := angularDistanceRad(lat1, lon1, latP, lonP)
	θ13 := initialBearingRad(lat1, lon1, latP, lonP)
	θ12 := initialBearingRad(lat1, lon1, lat2, lon2)

	δxt := math.Asin(math.Sin(δ13) * math.Sin(θ13-θ12))
	crossTrackKm := δxt * EarthRadiusKm

	δat := math.Atan2(math.Sin(δ13)*math.Cos(θ13-θ12), math.Cos(δ13))
	alongTrackKm := δat * EarthRadiusKm

	fraction := δat / totalAngle
	projLat, projLon := GreatCircleIntermediatePoint(lat1, lon1, lat2, lon2, fraction)

	return projLat, projLon, crossTrackKm, alongTrackKm
}

// GreatCircleIntermediatePoint returns the point at the given fraction along the
// great circle path between two coordinates. Fraction 0 returns the start point,
// fraction 1 returns the end point. Coordinates are in degrees (latitude, longitude).
func GreatCircleIntermediatePoint(lat1, lon1, lat2, lon2, fraction float64) (float64, float64) {
	φ1 := toRadians(lat1)
	λ1 := toRadians(lon1)
	φ2 := toRadians(lat2)
	λ2 := toRadians(lon2)

	Δφ := φ2 - φ1
	Δλ := λ2 - λ1

	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*
			math.Sin(Δλ/2)*math.Sin(Δλ/2)
	δ := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	if δ == 0 {
		return lat1, normalizeLongitude(lon1)
	}

	aCoef := math.Sin((1-fraction)*δ) / math.Sin(δ)
	bCoef := math.Sin(fraction*δ) / math.Sin(δ)

	x := aCoef*math.Cos(φ1)*math.Cos(λ1) + bCoef*math.Cos(φ2)*math.Cos(λ2)
	y := aCoef*math.Cos(φ1)*math.Sin(λ1) + bCoef*math.Cos(φ2)*math.Sin(λ2)
	z := aCoef*math.Sin(φ1) + bCoef*math.Sin(φ2)

	φi := math.Atan2(z, math.Sqrt(x*x+y*y))
	λi := math.Atan2(y, x)

	return toDegrees(φi), normalizeLongitude(toDegrees(λi))
}

// GreatCirclePointAtDistance returns the point at a given distance (in kilometers)
// along the great circle path between two coordinates. Distance is clamped to [0, total].
func GreatCirclePointAtDistance(lat1, lon1, lat2, lon2, distanceKm float64) (float64, float64) {
	total := GreatCircleDistance(lat1, lon1, lat2, lon2)
	if total == 0 {
		return lat1, normalizeLongitude(lon1)
	}
	if distanceKm <= 0 {
		return lat1, normalizeLongitude(lon1)
	}
	if distanceKm >= total {
		return lat2, normalizeLongitude(lon2)
	}
	return GreatCircleIntermediatePoint(lat1, lon1, lat2, lon2, distanceKm/total)
}

// GreatCirclePointAtSpeed returns the point after traveling at speedKmh for durationHours
// along the great circle path between two coordinates.
func GreatCirclePointAtSpeed(lat1, lon1, lat2, lon2, speedKmh, durationHours float64) (float64, float64) {
	return GreatCirclePointAtDistance(lat1, lon1, lat2, lon2, speedKmh*durationHours)
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
