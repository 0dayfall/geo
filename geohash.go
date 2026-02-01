package geo

import (
	"strings"
)

const (
	// base32 characters used in geohash encoding
	base32 = "0123456789bcdefghjkmnpqrstuvwxyz"
)

// Geohash encodes a geographic coordinate (latitude, longitude) into a geohash string.
// The precision parameter determines the length of the resulting geohash string.
func Geohash(lat, lon float64, precision int) string {
	if precision <= 0 {
		precision = 12 // default precision
	}

	latRange := [2]float64{-90.0, 90.0}
	lonRange := [2]float64{-180.0, 180.0}

	var geohash strings.Builder
	var bits uint
	var bit uint
	even := true

	for geohash.Len() < precision {
		if even {
			// longitude
			mid := (lonRange[0] + lonRange[1]) / 2
			if lon > mid {
				bits |= (1 << (4 - bit))
				lonRange[0] = mid
			} else {
				lonRange[1] = mid
			}
		} else {
			// latitude
			mid := (latRange[0] + latRange[1]) / 2
			if lat > mid {
				bits |= (1 << (4 - bit))
				latRange[0] = mid
			} else {
				latRange[1] = mid
			}
		}

		even = !even
		bit++

		if bit == 5 {
			geohash.WriteByte(base32[bits])
			bits = 0
			bit = 0
		}
	}

	return geohash.String()
}

// GeohashDecode decodes a geohash string into latitude and longitude coordinates.
// Returns the center point of the geohash cell and the error bounds.
func GeohashDecode(geohash string) (lat, lon, latErr, lonErr float64) {
	latRange := [2]float64{-90.0, 90.0}
	lonRange := [2]float64{-180.0, 180.0}

	even := true

	for i := 0; i < len(geohash); i++ {
		char := geohash[i]
		idx := strings.IndexByte(base32, char)
		if idx == -1 {
			// invalid character, return current position
			break
		}

		for mask := 4; mask >= 0; mask-- {
			if even {
				// longitude
				mid := (lonRange[0] + lonRange[1]) / 2
				if (idx & (1 << mask)) != 0 {
					lonRange[0] = mid
				} else {
					lonRange[1] = mid
				}
			} else {
				// latitude
				mid := (latRange[0] + latRange[1]) / 2
				if (idx & (1 << mask)) != 0 {
					latRange[0] = mid
				} else {
					latRange[1] = mid
				}
			}
			even = !even
		}
	}

	lat = (latRange[0] + latRange[1]) / 2
	lon = (lonRange[0] + lonRange[1]) / 2
	latErr = (latRange[1] - latRange[0]) / 2
	lonErr = (lonRange[1] - lonRange[0]) / 2

	return
}

// GeohashNeighbors returns the 8 neighboring geohashes around the given geohash.
// Returns neighbors in order: N, NE, E, SE, S, SW, W, NW
func GeohashNeighbors(geohash string) [8]string {
	lat, lon, latErr, lonErr := GeohashDecode(geohash)
	precision := len(geohash)

	// Calculate the 8 neighbors
	neighbors := [8]string{
		Geohash(lat+2*latErr, lon, precision),          // N
		Geohash(lat+2*latErr, lon+2*lonErr, precision), // NE
		Geohash(lat, lon+2*lonErr, precision),          // E
		Geohash(lat-2*latErr, lon+2*lonErr, precision), // SE
		Geohash(lat-2*latErr, lon, precision),          // S
		Geohash(lat-2*latErr, lon-2*lonErr, precision), // SW
		Geohash(lat, lon-2*lonErr, precision),          // W
		Geohash(lat+2*latErr, lon-2*lonErr, precision), // NW
	}

	return neighbors
}
