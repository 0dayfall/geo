package geo

import (
	"math"
	"testing"
)

func TestLineStringPointAtDistance(t *testing.T) {
	line := NewLineString([]Position{
		{0, 0},
		{90, 0},
	})
	total := GreatCircleDistance(0, 0, 0, 90)

	pt, err := LineStringPointAtDistance(line, total/2)
	if err != nil {
		t.Fatalf("LineStringPointAtDistance() error = %v", err)
	}
	if math.Abs(pt.Coordinates[0]-45.0) > 1e-6 || math.Abs(pt.Coordinates[1]-0.0) > 1e-6 {
		t.Errorf("point = (%v, %v), want (45, 0)", pt.Coordinates[0], pt.Coordinates[1])
	}
}

func TestGeoJSONBearing(t *testing.T) {
	bearingNorth := GeoJSONBearing(NewPoint(0, 0), NewPoint(0, 10))
	if math.Abs(bearingNorth-0.0) > 1e-6 {
		t.Errorf("bearing north = %v, want 0", bearingNorth)
	}

	bearingEast := GeoJSONBearing(NewPoint(0, 0), NewPoint(10, 0))
	if math.Abs(bearingEast-90.0) > 1e-6 {
		t.Errorf("bearing east = %v, want 90", bearingEast)
	}
}

func TestGeoJSONCenter(t *testing.T) {
	fc := NewFeatureCollection([]Feature{
		NewFeature(NewPoint(0, 0)),
		NewFeature(NewPoint(10, 10)),
	})
	center, err := GeoJSONCenter(fc)
	if err != nil {
		t.Fatalf("GeoJSONCenter() error = %v", err)
	}
	if math.Abs(center.Coordinates[0]-5.0) > 1e-9 || math.Abs(center.Coordinates[1]-5.0) > 1e-9 {
		t.Errorf("center = (%v, %v), want (5, 5)", center.Coordinates[0], center.Coordinates[1])
	}
}

func TestGeoJSONCenterOfMassPolygon(t *testing.T) {
	poly := NewPolygon([][]Position{
		{
			{0, 0},
			{2, 0},
			{2, 2},
			{0, 2},
			{0, 0},
		},
	})
	center, err := GeoJSONCenterOfMass(poly)
	if err != nil {
		t.Fatalf("GeoJSONCenterOfMass() error = %v", err)
	}
	if math.Abs(center.Coordinates[0]-1.0) > 1e-9 || math.Abs(center.Coordinates[1]-1.0) > 1e-9 {
		t.Errorf("center = (%v, %v), want (1, 1)", center.Coordinates[0], center.Coordinates[1])
	}
}

func TestGreatCircleGeoJSON(t *testing.T) {
	geom, err := GreatCircleGeoJSON(NewPoint(179, 0), NewPoint(-179, 0), 5)
	if err != nil {
		t.Fatalf("GreatCircleGeoJSON() error = %v", err)
	}
	if _, ok := geom.(MultiLineString); !ok {
		t.Errorf("expected MultiLineString for antimeridian crossing")
	}

	geom2, err := GreatCircleGeoJSON(NewPoint(0, 0), NewPoint(0, 0), 3)
	if err != nil {
		t.Fatalf("GreatCircleGeoJSON() error = %v", err)
	}
	ls, ok := geom2.(LineString)
	if !ok {
		t.Fatalf("expected LineString for identical points")
	}
	if len(ls.Coordinates) != 3 {
		t.Errorf("linestring length = %v, want 3", len(ls.Coordinates))
	}
}

func TestGreatCircleGeoJSONByDistance(t *testing.T) {
	geom, err := GreatCircleGeoJSONByDistance(NewPoint(179, 0), NewPoint(-179, 0), 200)
	if err != nil {
		t.Fatalf("GreatCircleGeoJSONByDistance() error = %v", err)
	}
	if _, ok := geom.(MultiLineString); !ok {
		t.Errorf("expected MultiLineString for antimeridian crossing")
	}

	geom2, err := GreatCircleGeoJSONByDistance(NewPoint(0, 0), NewPoint(90, 0), 5000)
	if err != nil {
		t.Fatalf("GreatCircleGeoJSONByDistance() error = %v", err)
	}
	ls, ok := geom2.(LineString)
	if !ok {
		t.Fatalf("expected LineString")
	}
	if len(ls.Coordinates) < 3 {
		t.Errorf("linestring length = %v, want at least 3", len(ls.Coordinates))
	}
	if len(ls.Coordinates) >= 2 {
		near := ls.Coordinates[1]
		if math.Abs(near[0]-45.0) > 2.0 || math.Abs(near[1]-0.0) > 1e-6 {
			t.Errorf("sample point approx = (%v, %v), want near (45, 0)", near[0], near[1])
		}
	}
}

func TestLinePointDistance(t *testing.T) {
	line := NewLineString([]Position{
		{0, 0},
		{90, 0},
	})
	point := NewPoint(45, 10)
	dist, err := LinePointDistance(line, point)
	if err != nil {
		t.Fatalf("LinePointDistance() error = %v", err)
	}
	expected := EarthRadiusKm * toRadians(10.0)
	if math.Abs(dist-expected) > 1e-3 {
		t.Errorf("distance = %v, want %v", dist, expected)
	}
}

func TestPolygonPointDistance(t *testing.T) {
	poly := NewPolygon([][]Position{
		{
			{0, 0},
			{2, 0},
			{2, 2},
			{0, 2},
			{0, 0},
		},
	})
	point := NewPoint(1, 1)
	dist, err := PolygonPointDistance(poly, point)
	if err != nil {
		t.Fatalf("PolygonPointDistance() error = %v", err)
	}
	expected := GreatCircleDistance(1, 1, 0, 1)
	if dist >= 0 || math.Abs(math.Abs(dist)-expected) > 0.05 {
		t.Errorf("distance = %v, want negative approx %v", dist, expected)
	}
}
