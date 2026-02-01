package geo

// DistanceUnit represents unit conversions for distance values.
type DistanceUnit int

const (
	UnitKilometers DistanceUnit = iota
	UnitMeters
	UnitMiles
	UnitNauticalMiles
)

const (
	// KmPerMile converts miles to kilometers.
	KmPerMile = 1.609344
)

// ConvertDistanceFromKm converts a kilometer value to the requested unit.
func ConvertDistanceFromKm(km float64, unit DistanceUnit) float64 {
	switch unit {
	case UnitMeters:
		return km * MetersPerKm
	case UnitMiles:
		return km / KmPerMile
	case UnitNauticalMiles:
		return km / KmPerNauticalMile
	case UnitKilometers:
		fallthrough
	default:
		return km
	}
}

// ConvertDistanceToKm converts a distance from the requested unit to kilometers.
func ConvertDistanceToKm(value float64, unit DistanceUnit) float64 {
	switch unit {
	case UnitMeters:
		return value / MetersPerKm
	case UnitMiles:
		return value * KmPerMile
	case UnitNauticalMiles:
		return value * KmPerNauticalMile
	case UnitKilometers:
		fallthrough
	default:
		return value
	}
}
