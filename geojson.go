package geo

import (
	"errors"
	"fmt"
	"math"
)

// Position represents a GeoJSON coordinate [longitude, latitude].
type Position [2]float64

// Point is a GeoJSON Point geometry.
type Point struct {
	Type        string   `json:"type"`
	Coordinates Position `json:"coordinates"`
}

// LineString is a GeoJSON LineString geometry.
type LineString struct {
	Type        string     `json:"type"`
	Coordinates []Position `json:"coordinates"`
}

// Polygon is a GeoJSON Polygon geometry.
type Polygon struct {
	Type        string       `json:"type"`
	Coordinates [][]Position `json:"coordinates"`
}

// MultiLineString is a GeoJSON MultiLineString geometry.
type MultiLineString struct {
	Type        string       `json:"type"`
	Coordinates [][]Position `json:"coordinates"`
}

// MultiPolygon is a GeoJSON MultiPolygon geometry.
type MultiPolygon struct {
	Type        string         `json:"type"`
	Coordinates [][][]Position `json:"coordinates"`
}

// Feature is a GeoJSON Feature.
type Feature struct {
	Type       string                 `json:"type"`
	Geometry   interface{}            `json:"geometry"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// FeatureCollection is a GeoJSON FeatureCollection.
type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

// NewPoint creates a GeoJSON Point.
func NewPoint(lon, lat float64) Point {
	return Point{Type: "Point", Coordinates: Position{lon, lat}}
}

// NewLineString creates a GeoJSON LineString.
func NewLineString(coords []Position) LineString {
	return LineString{Type: "LineString", Coordinates: coords}
}

// NewPolygon creates a GeoJSON Polygon.
func NewPolygon(coords [][]Position) Polygon {
	return Polygon{Type: "Polygon", Coordinates: coords}
}

// NewMultiLineString creates a GeoJSON MultiLineString.
func NewMultiLineString(coords [][]Position) MultiLineString {
	return MultiLineString{Type: "MultiLineString", Coordinates: coords}
}

// NewMultiPolygon creates a GeoJSON MultiPolygon.
func NewMultiPolygon(coords [][][]Position) MultiPolygon {
	return MultiPolygon{Type: "MultiPolygon", Coordinates: coords}
}

// NewFeature creates a GeoJSON Feature.
func NewFeature(geom interface{}) Feature {
	return Feature{Type: "Feature", Geometry: geom}
}

// NewFeatureCollection creates a GeoJSON FeatureCollection.
func NewFeatureCollection(features []Feature) FeatureCollection {
	return FeatureCollection{Type: "FeatureCollection", Features: features}
}

func positionLatLon(p Position) (lat, lon float64) {
	return p[1], p[0]
}

func pointFromLatLon(lat, lon float64) Point {
	return NewPoint(lon, lat)
}

// LineStringPointAtDistance returns a Point at a specified distance along the LineString.
// Distance is in kilometers. If distance is <= 0, the start point is returned.
// If distance exceeds the line length, the end point is returned.
func LineStringPointAtDistance(line LineString, distanceKm float64) (Point, error) {
	if len(line.Coordinates) < 2 {
		return Point{}, errors.New("linestring must have at least 2 coordinates")
	}
	if distanceKm <= 0 {
		return pointFromLatLon(positionLatLon(line.Coordinates[0])), nil
	}

	remaining := distanceKm
	for i := 0; i < len(line.Coordinates)-1; i++ {
		start := line.Coordinates[i]
		end := line.Coordinates[i+1]
		lat1, lon1 := positionLatLon(start)
		lat2, lon2 := positionLatLon(end)
		seg := GreatCircleDistance(lat1, lon1, lat2, lon2)
		if remaining <= seg {
			f := remaining / seg
			lat, lon := GreatCircleIntermediatePoint(lat1, lon1, lat2, lon2, f)
			return pointFromLatLon(lat, lon), nil
		}
		remaining -= seg
	}

	last := line.Coordinates[len(line.Coordinates)-1]
	return pointFromLatLon(positionLatLon(last)), nil
}

// GeoJSONBearing returns the great-circle bearing between two GeoJSON Points.
// Bearing is in degrees from true north, in the range [0, 360).
func GeoJSONBearing(start, end Point) float64 {
	lat1, lon1 := positionLatLon(start.Coordinates)
	lat2, lon2 := positionLatLon(end.Coordinates)
	return Bearing(lat1, lon1, lat2, lon2)
}

// GeoJSONRhumbBearing returns the rhumb line bearing between two GeoJSON Points.
// Bearing is in degrees from true north, in the range [0, 360).
func GeoJSONRhumbBearing(start, end Point) float64 {
	lat1, lon1 := positionLatLon(start.Coordinates)
	lat2, lon2 := positionLatLon(end.Coordinates)
	return RhumbLineBearing(lat1, lon1, lat2, lon2)
}

// GeoJSONRhumbDestination returns the destination Point after traveling along a rhumb line.
// Distance is in kilometers, bearing is in degrees from true north.
func GeoJSONRhumbDestination(start Point, distanceKm, bearingDeg float64) Point {
	lat1, lon1 := positionLatLon(start.Coordinates)
	lat2, lon2 := RhumbLineDestination(lat1, lon1, distanceKm, bearingDeg)
	return NewPoint(lon2, lat2)
}

// GeoJSONRhumbDistance returns rhumb line distance between two Points in the requested unit.
func GeoJSONRhumbDistance(start, end Point, unit DistanceUnit) float64 {
	lat1, lon1 := positionLatLon(start.Coordinates)
	lat2, lon2 := positionLatLon(end.Coordinates)
	return RhumbLineDistanceUnits(lat1, lon1, lat2, lon2, unit)
}

// GeoJSONCenter returns the bbox center of all coordinates in a Feature or FeatureCollection.
func GeoJSONCenter(obj interface{}) (Point, error) {
	positions, err := collectPositions(obj)
	if err != nil {
		return Point{}, err
	}
	if len(positions) == 0 {
		return Point{}, errors.New("no coordinates found")
	}

	minLon, maxLon := positions[0][0], positions[0][0]
	minLat, maxLat := positions[0][1], positions[0][1]
	for _, p := range positions[1:] {
		if p[0] < minLon {
			minLon = p[0]
		}
		if p[0] > maxLon {
			maxLon = p[0]
		}
		if p[1] < minLat {
			minLat = p[1]
		}
		if p[1] > maxLat {
			maxLat = p[1]
		}
	}

	return NewPoint((minLon+maxLon)/2, (minLat+maxLat)/2), nil
}

// GeoJSONCenterOfMass returns a center-of-mass point.
// For polygons, this uses the centroid of polygon formula.
// For lines, it uses the midpoint along length.
// For points, it uses the average of points.
func GeoJSONCenterOfMass(obj interface{}) (Point, error) {
	acc := massAccumulator{}
	if err := acc.add(obj); err != nil {
		return Point{}, err
	}

	switch {
	case acc.areaSum > 0:
		return NewPoint(acc.areaLonSum/acc.areaSum, acc.areaLatSum/acc.areaSum), nil
	case acc.lengthSum > 0:
		return NewPoint(acc.lengthLonSum/acc.lengthSum, acc.lengthLatSum/acc.lengthSum), nil
	case acc.pointCount > 0:
		return NewPoint(acc.pointLonSum/float64(acc.pointCount), acc.pointLatSum/float64(acc.pointCount)), nil
	default:
		return Point{}, errors.New("no coordinates found")
	}
}

// GeoJSONPointOnSurface returns a Point guaranteed to lie on the feature's surface.
func GeoJSONPointOnSurface(obj interface{}) (Point, error) {
	switch g := obj.(type) {
	case Point:
		return g, nil
	case *Point:
		if g == nil {
			return Point{}, errors.New("nil point")
		}
		return *g, nil
	case LineString:
		return lineMidpoint(g)
	case *LineString:
		if g == nil {
			return Point{}, errors.New("nil linestring")
		}
		return lineMidpoint(*g)
	case Polygon:
		return polygonPointOnSurface(g)
	case *Polygon:
		if g == nil {
			return Point{}, errors.New("nil polygon")
		}
		return polygonPointOnSurface(*g)
	case MultiLineString:
		return multiLinePointOnSurface(g)
	case *MultiLineString:
		if g == nil {
			return Point{}, errors.New("nil multilinestring")
		}
		return multiLinePointOnSurface(*g)
	case MultiPolygon:
		return multiPolygonPointOnSurface(g)
	case *MultiPolygon:
		if g == nil {
			return Point{}, errors.New("nil multipolygon")
		}
		return multiPolygonPointOnSurface(*g)
	case Feature:
		return GeoJSONPointOnSurface(g.Geometry)
	case *Feature:
		if g == nil {
			return Point{}, errors.New("nil feature")
		}
		return GeoJSONPointOnSurface(g.Geometry)
	case FeatureCollection:
		return featureCollectionPointOnSurface(g)
	case *FeatureCollection:
		if g == nil {
			return Point{}, errors.New("nil featurecollection")
		}
		return featureCollectionPointOnSurface(*g)
	default:
		return Point{}, fmt.Errorf("unsupported geojson type %T", obj)
	}
}

// GreatCircleGeoJSON returns a great-circle route as a LineString or MultiLineString.
// If the path crosses the antimeridian, a MultiLineString is returned.
// If start and end are the same, a LineString with duplicate coordinates is returned.
func GreatCircleGeoJSON(start, end Point, npoints int) (interface{}, error) {
	if npoints <= 0 {
		npoints = 2
	}

	startPos := start.Coordinates
	endPos := end.Coordinates

	if startPos == endPos {
		coords := make([]Position, npoints)
		for i := 0; i < npoints; i++ {
			coords[i] = startPos
		}
		return NewLineString(coords), nil
	}

	lat1, lon1 := positionLatLon(startPos)
	lat2, lon2 := positionLatLon(endPos)

	coords := make([]Position, npoints)
	for i := 0; i < npoints; i++ {
		f := float64(i) / float64(npoints-1)
		lat, lon := GreatCircleIntermediatePoint(lat1, lon1, lat2, lon2, f)
		coords[i] = Position{lon, lat}
	}

	var lines [][]Position
	current := []Position{coords[0]}
	for i := 1; i < len(coords); i++ {
		prev := coords[i-1]
		curr := coords[i]
		if math.Abs(curr[0]-prev[0]) > 180.0 {
			lines = append(lines, current)
			current = []Position{curr}
		} else {
			current = append(current, curr)
		}
	}
	lines = append(lines, current)

	if len(lines) == 1 {
		return NewLineString(lines[0]), nil
	}
	return NewMultiLineString(lines), nil
}

// LinePointDistance returns the distance between a point and the nearest point on a line.
// Distance is returned in kilometers.
func LinePointDistance(line LineString, point Point) (float64, error) {
	if len(line.Coordinates) < 2 {
		return 0, errors.New("linestring must have at least 2 coordinates")
	}

	latP, lonP := positionLatLon(point.Coordinates)
	minDist := math.Inf(1)

	for i := 0; i < len(line.Coordinates)-1; i++ {
		start := line.Coordinates[i]
		end := line.Coordinates[i+1]
		lat1, lon1 := positionLatLon(start)
		lat2, lon2 := positionLatLon(end)
		_, _, crossTrackKm, _ := GreatCircleProjectToSegment(lat1, lon1, lat2, lon2, latP, lonP)
		dist := math.Abs(crossTrackKm)
		if dist < minDist {
			minDist = dist
		}
	}

	return minDist, nil
}

// PolygonPointDistance returns signed distance from a point to the edges of a polygon or multipolygon.
// Distances are in kilometers. Negative values indicate the point is inside the polygon.
// A hole is treated as exterior.
func PolygonPointDistance(obj interface{}, point Point) (float64, error) {
	switch g := obj.(type) {
	case Polygon:
		return polygonPointDistance(g, point)
	case *Polygon:
		if g == nil {
			return 0, errors.New("nil polygon")
		}
		return polygonPointDistance(*g, point)
	case MultiPolygon:
		return multiPolygonPointDistance(g, point)
	case *MultiPolygon:
		if g == nil {
			return 0, errors.New("nil multipolygon")
		}
		return multiPolygonPointDistance(*g, point)
	case Feature:
		return PolygonPointDistance(g.Geometry, point)
	case *Feature:
		if g == nil {
			return 0, errors.New("nil feature")
		}
		return PolygonPointDistance(g.Geometry, point)
	case FeatureCollection:
		return polygonDistanceFromCollection(g, point)
	case *FeatureCollection:
		if g == nil {
			return 0, errors.New("nil featurecollection")
		}
		return polygonDistanceFromCollection(*g, point)
	default:
		return 0, fmt.Errorf("unsupported geojson type %T", obj)
	}
}

// ---------------- Helpers ----------------

func collectPositions(obj interface{}) ([]Position, error) {
	var positions []Position
	if err := collectPositionsInto(obj, &positions); err != nil {
		return nil, err
	}
	return positions, nil
}

func collectPositionsInto(obj interface{}, positions *[]Position) error {
	switch g := obj.(type) {
	case Point:
		*positions = append(*positions, g.Coordinates)
	case *Point:
		if g == nil {
			return errors.New("nil point")
		}
		*positions = append(*positions, g.Coordinates)
	case LineString:
		*positions = append(*positions, g.Coordinates...)
	case *LineString:
		if g == nil {
			return errors.New("nil linestring")
		}
		*positions = append(*positions, g.Coordinates...)
	case Polygon:
		for _, ring := range g.Coordinates {
			*positions = append(*positions, ring...)
		}
	case *Polygon:
		if g == nil {
			return errors.New("nil polygon")
		}
		for _, ring := range g.Coordinates {
			*positions = append(*positions, ring...)
		}
	case MultiLineString:
		for _, line := range g.Coordinates {
			*positions = append(*positions, line...)
		}
	case *MultiLineString:
		if g == nil {
			return errors.New("nil multilinestring")
		}
		for _, line := range g.Coordinates {
			*positions = append(*positions, line...)
		}
	case MultiPolygon:
		for _, poly := range g.Coordinates {
			for _, ring := range poly {
				*positions = append(*positions, ring...)
			}
		}
	case *MultiPolygon:
		if g == nil {
			return errors.New("nil multipolygon")
		}
		for _, poly := range g.Coordinates {
			for _, ring := range poly {
				*positions = append(*positions, ring...)
			}
		}
	case Feature:
		return collectPositionsInto(g.Geometry, positions)
	case *Feature:
		if g == nil {
			return errors.New("nil feature")
		}
		return collectPositionsInto(g.Geometry, positions)
	case FeatureCollection:
		for i := range g.Features {
			if err := collectPositionsInto(g.Features[i], positions); err != nil {
				return err
			}
		}
	case *FeatureCollection:
		if g == nil {
			return errors.New("nil featurecollection")
		}
		for i := range g.Features {
			if err := collectPositionsInto(g.Features[i], positions); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unsupported geojson type %T", obj)
	}
	return nil
}

type massAccumulator struct {
	areaSum      float64
	areaLonSum   float64
	areaLatSum   float64
	lengthSum    float64
	lengthLonSum float64
	lengthLatSum float64
	pointCount   int
	pointLonSum  float64
	pointLatSum  float64
}

func (m *massAccumulator) add(obj interface{}) error {
	switch g := obj.(type) {
	case Point:
		m.addPoint(g.Coordinates)
	case *Point:
		if g == nil {
			return errors.New("nil point")
		}
		m.addPoint(g.Coordinates)
	case LineString:
		m.addLine(g)
	case *LineString:
		if g == nil {
			return errors.New("nil linestring")
		}
		m.addLine(*g)
	case Polygon:
		m.addPolygon(g)
	case *Polygon:
		if g == nil {
			return errors.New("nil polygon")
		}
		m.addPolygon(*g)
	case MultiLineString:
		for _, line := range g.Coordinates {
			m.addLine(LineString{Coordinates: line})
		}
	case *MultiLineString:
		if g == nil {
			return errors.New("nil multilinestring")
		}
		for _, line := range g.Coordinates {
			m.addLine(LineString{Coordinates: line})
		}
	case MultiPolygon:
		for _, poly := range g.Coordinates {
			m.addPolygon(Polygon{Coordinates: poly})
		}
	case *MultiPolygon:
		if g == nil {
			return errors.New("nil multipolygon")
		}
		for _, poly := range g.Coordinates {
			m.addPolygon(Polygon{Coordinates: poly})
		}
	case Feature:
		return m.add(g.Geometry)
	case *Feature:
		if g == nil {
			return errors.New("nil feature")
		}
		return m.add(g.Geometry)
	case FeatureCollection:
		for i := range g.Features {
			if err := m.add(g.Features[i]); err != nil {
				return err
			}
		}
	case *FeatureCollection:
		if g == nil {
			return errors.New("nil featurecollection")
		}
		for i := range g.Features {
			if err := m.add(g.Features[i]); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unsupported geojson type %T", obj)
	}
	return nil
}

func (m *massAccumulator) addPoint(p Position) {
	m.pointCount++
	m.pointLonSum += p[0]
	m.pointLatSum += p[1]
}

func (m *massAccumulator) addLine(line LineString) {
	if len(line.Coordinates) < 2 {
		return
	}
	length, mid, err := lineMidpointWithLength(line)
	if err != nil || length == 0 {
		return
	}
	m.lengthSum += length
	m.lengthLonSum += mid[0] * length
	m.lengthLatSum += mid[1] * length
}

func (m *massAccumulator) addPolygon(poly Polygon) {
	centroid, area, ok := polygonCentroidArea(poly)
	if !ok || area == 0 {
		return
	}
	m.areaSum += area
	m.areaLonSum += centroid[0] * area
	m.areaLatSum += centroid[1] * area
}

func lineMidpoint(line LineString) (Point, error) {
	length, mid, err := lineMidpointWithLength(line)
	if err != nil {
		return Point{}, err
	}
	if length == 0 {
		return pointFromLatLon(positionLatLon(line.Coordinates[0])), nil
	}
	return NewPoint(mid[0], mid[1]), nil
}

func lineMidpointWithLength(line LineString) (float64, Position, error) {
	if len(line.Coordinates) < 2 {
		return 0, Position{}, errors.New("linestring must have at least 2 coordinates")
	}
	length, err := lineStringLengthKm(line)
	if err != nil || length == 0 {
		return length, Position{}, err
	}
	mid, err := LineStringPointAtDistance(line, length/2)
	if err != nil {
		return 0, Position{}, err
	}
	return length, mid.Coordinates, nil
}

func lineStringLengthKm(line LineString) (float64, error) {
	if len(line.Coordinates) < 2 {
		return 0, errors.New("linestring must have at least 2 coordinates")
	}
	var total float64
	for i := 0; i < len(line.Coordinates)-1; i++ {
		start := line.Coordinates[i]
		end := line.Coordinates[i+1]
		lat1, lon1 := positionLatLon(start)
		lat2, lon2 := positionLatLon(end)
		total += GreatCircleDistance(lat1, lon1, lat2, lon2)
	}
	return total, nil
}

func polygonPointOnSurface(poly Polygon) (Point, error) {
	if len(poly.Coordinates) == 0 || len(poly.Coordinates[0]) == 0 {
		return Point{}, errors.New("polygon has no coordinates")
	}
	centroid, _, ok := polygonCentroidArea(poly)
	if ok && pointInPolygon(centroid, poly) {
		return NewPoint(centroid[0], centroid[1]), nil
	}
	return NewPoint(poly.Coordinates[0][0][0], poly.Coordinates[0][0][1]), nil
}

func multiLinePointOnSurface(ml MultiLineString) (Point, error) {
	var best LineString
	var bestLen float64
	for _, line := range ml.Coordinates {
		length, err := lineStringLengthKm(LineString{Coordinates: line})
		if err == nil && length > bestLen {
			bestLen = length
			best = LineString{Coordinates: line}
		}
	}
	if bestLen == 0 && len(ml.Coordinates) > 0 {
		best = LineString{Coordinates: ml.Coordinates[0]}
	}
	if len(best.Coordinates) == 0 {
		return Point{}, errors.New("multilinestring has no coordinates")
	}
	return lineMidpoint(best)
}

func multiPolygonPointOnSurface(mp MultiPolygon) (Point, error) {
	var best Polygon
	var bestArea float64
	for _, poly := range mp.Coordinates {
		centroid, area, ok := polygonCentroidArea(Polygon{Coordinates: poly})
		if ok && area > bestArea {
			bestArea = area
			best = Polygon{Coordinates: poly}
			_ = centroid
		}
	}
	if bestArea == 0 && len(mp.Coordinates) > 0 {
		best = Polygon{Coordinates: mp.Coordinates[0]}
	}
	if len(best.Coordinates) == 0 {
		return Point{}, errors.New("multipolygon has no coordinates")
	}
	return polygonPointOnSurface(best)
}

func featureCollectionPointOnSurface(fc FeatureCollection) (Point, error) {
	var bestPoly Polygon
	var bestArea float64
	var bestLine LineString
	var bestLineLen float64
	var firstPoint *Point

	for i := range fc.Features {
		switch g := fc.Features[i].Geometry.(type) {
		case Point:
			if firstPoint == nil {
				p := g
				firstPoint = &p
			}
		case LineString:
			length, err := lineStringLengthKm(g)
			if err == nil && length > bestLineLen {
				bestLineLen = length
				bestLine = g
			}
		case Polygon:
			_, area, ok := polygonCentroidArea(g)
			if ok && area > bestArea {
				bestArea = area
				bestPoly = g
			}
		case MultiLineString:
			p, err := multiLinePointOnSurface(g)
			if err == nil && bestLineLen == 0 {
				return p, nil
			}
		case MultiPolygon:
			p, err := multiPolygonPointOnSurface(g)
			if err == nil && bestArea == 0 {
				return p, nil
			}
		}
	}

	if bestArea > 0 {
		return polygonPointOnSurface(bestPoly)
	}
	if bestLineLen > 0 {
		return lineMidpoint(bestLine)
	}
	if firstPoint != nil {
		return *firstPoint, nil
	}
	return Point{}, errors.New("featurecollection has no supported geometries")
}

func polygonPointDistance(poly Polygon, point Point) (float64, error) {
	if len(poly.Coordinates) == 0 {
		return 0, errors.New("polygon has no coordinates")
	}
	pt := point.Coordinates
	minDist := math.Inf(1)

	for _, ring := range poly.Coordinates {
		dist, err := ringDistance(ring, point)
		if err != nil {
			continue
		}
		if dist < minDist {
			minDist = dist
		}
	}

	if math.IsInf(minDist, 1) {
		return 0, errors.New("unable to compute distance to polygon edges")
	}

	if pointInPolygon(pt, poly) {
		return -minDist, nil
	}
	return minDist, nil
}

func multiPolygonPointDistance(mp MultiPolygon, point Point) (float64, error) {
	minDist := math.Inf(1)
	inside := false

	for _, poly := range mp.Coordinates {
		polygon := Polygon{Coordinates: poly}
		dist, err := polygonPointDistance(polygon, point)
		if err != nil {
			continue
		}
		if math.Abs(dist) < minDist {
			minDist = math.Abs(dist)
		}
		if dist < 0 {
			inside = true
		}
	}

	if math.IsInf(minDist, 1) {
		return 0, errors.New("multipolygon has no valid rings")
	}
	if inside {
		return -minDist, nil
	}
	return minDist, nil
}

func polygonDistanceFromCollection(fc FeatureCollection, point Point) (float64, error) {
	minDist := math.Inf(1)
	inside := false

	for i := range fc.Features {
		switch g := fc.Features[i].Geometry.(type) {
		case Polygon:
			dist, err := polygonPointDistance(g, point)
			if err == nil {
				if math.Abs(dist) < minDist {
					minDist = math.Abs(dist)
				}
				if dist < 0 {
					inside = true
				}
			}
		case MultiPolygon:
			dist, err := multiPolygonPointDistance(g, point)
			if err == nil {
				if math.Abs(dist) < minDist {
					minDist = math.Abs(dist)
				}
				if dist < 0 {
					inside = true
				}
			}
		}
	}

	if math.IsInf(minDist, 1) {
		return 0, errors.New("featurecollection contains no polygons")
	}
	if inside {
		return -minDist, nil
	}
	return minDist, nil
}

func ringDistance(ring []Position, point Point) (float64, error) {
	if len(ring) < 2 {
		return 0, errors.New("ring must have at least 2 coordinates")
	}
	coords := ring
	if ring[0] != ring[len(ring)-1] {
		coords = append(coords, ring[0])
	}
	return LinePointDistance(LineString{Coordinates: coords}, point)
}

func pointInPolygon(pt Position, poly Polygon) bool {
	if len(poly.Coordinates) == 0 {
		return false
	}
	if !pointInRing(pt, poly.Coordinates[0]) {
		return false
	}
	for i := 1; i < len(poly.Coordinates); i++ {
		if pointInRing(pt, poly.Coordinates[i]) {
			return false
		}
	}
	return true
}

func pointInRing(pt Position, ring []Position) bool {
	n := len(ring)
	if n < 3 {
		return false
	}

	// Treat boundary as inside.
	for i := 0; i < n-1; i++ {
		if pointOnSegment(pt, ring[i], ring[i+1]) {
			return true
		}
	}
	if ring[0] != ring[n-1] && pointOnSegment(pt, ring[n-1], ring[0]) {
		return true
	}

	inside := false
	j := n - 1
	x := pt[0]
	y := pt[1]
	for i := 0; i < n; i++ {
		xi := ring[i][0]
		yi := ring[i][1]
		xj := ring[j][0]
		yj := ring[j][1]

		intersect := ((yi > y) != (yj > y)) &&
			(x < (xj-xi)*(y-yi)/(yj-yi)+xi)
		if intersect {
			inside = !inside
		}
		j = i
	}
	return inside
}

func pointOnSegment(p, a, b Position) bool {
	const eps = 1e-12
	ax, ay := a[0], a[1]
	bx, by := b[0], b[1]
	px, py := p[0], p[1]

	cross := (px-ax)*(by-ay) - (py-ay)*(bx-ax)
	if math.Abs(cross) > eps {
		return false
	}
	dot := (px-ax)*(bx-ax) + (py-ay)*(by-ay)
	if dot < -eps {
		return false
	}
	sqLen := (bx-ax)*(bx-ax) + (by-ay)*(by-ay)
	if dot-sqLen > eps {
		return false
	}
	return true
}

func polygonCentroidArea(poly Polygon) (Position, float64, bool) {
	if len(poly.Coordinates) == 0 {
		return Position{}, 0, false
	}
	outer := poly.Coordinates[0]
	outerArea, outerCx, outerCy := ringAreaCentroid(outer)
	if outerArea == 0 {
		return Position{}, 0, false
	}

	areaSum := math.Abs(outerArea)
	lonSum := outerCx * math.Abs(outerArea)
	latSum := outerCy * math.Abs(outerArea)

	for i := 1; i < len(poly.Coordinates); i++ {
		area, cx, cy := ringAreaCentroid(poly.Coordinates[i])
		if area == 0 {
			continue
		}
		absArea := math.Abs(area)
		areaSum -= absArea
		lonSum -= cx * absArea
		latSum -= cy * absArea
	}

	if areaSum <= 0 {
		return Position{}, 0, false
	}
	return Position{lonSum / areaSum, latSum / areaSum}, areaSum, true
}

func ringAreaCentroid(ring []Position) (float64, float64, float64) {
	n := len(ring)
	if n < 3 {
		return 0, 0, 0
	}

	var area float64
	var cx float64
	var cy float64
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		x0 := ring[i][0]
		y0 := ring[i][1]
		x1 := ring[j][0]
		y1 := ring[j][1]
		cross := x0*y1 - x1*y0
		area += cross
		cx += (x0 + x1) * cross
		cy += (y0 + y1) * cross
	}
	area *= 0.5
	if area == 0 {
		return 0, 0, 0
	}
	cx /= 6 * area
	cy /= 6 * area
	return area, cx, cy
}
