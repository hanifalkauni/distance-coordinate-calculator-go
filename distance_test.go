package distance

import (
	"math"
	"strings"
	"testing"
)

// NYC and LA reference points, reused from the PHP test suite.
const (
	nycLat, nycLon = 40.748817, -73.985428
	laLat, laLon   = 34.052235, -118.243683
)

func TestCalculateKilometers(t *testing.T) {
	got := Calculate(nycLat, nycLon, laLat, laLon, Kilometers)
	if math.Abs(math.Round(got)-3940) > 10 {
		t.Fatalf("KM: got %.2f, want ~3940", got)
	}
}

func TestCalculateMiles(t *testing.T) {
	got := Calculate(nycLat, nycLon, laLat, laLon, Miles)
	if math.Abs(math.Round(got)-2448) > 10 {
		t.Fatalf("Miles: got %.2f, want ~2448", got)
	}
}

func TestCalculateMeters(t *testing.T) {
	got := Calculate(nycLat, nycLon, laLat, laLon, Meters)
	if math.Abs(math.Round(got)-3940000) > 10000 {
		t.Fatalf("Meters: got %.2f, want ~3940000", got)
	}
}

func TestCalculateDefaultUnitFallback(t *testing.T) {
	// Empty unit must behave like KM.
	got := Calculate(nycLat, nycLon, laLat, laLon, "")
	if math.Abs(math.Round(got)-3940) > 10 {
		t.Fatalf("default (KM): got %.2f, want ~3940", got)
	}
}

func TestCalculateFormatted(t *testing.T) {
	got := CalculateFormatted(nycLat, nycLon, laLat, laLon, Kilometers, 1)
	if !strings.HasSuffix(got, "KM") {
		t.Fatalf("formatted %q must end with KM", got)
	}
	if !strings.HasPrefix(got, "3937") {
		t.Fatalf("formatted %q must start with 3937", got)
	}
}

func TestIsWithinRadius(t *testing.T) {
	if !IsWithinRadius(nycLat, nycLon, laLat, laLon, 4000, Kilometers) {
		t.Error("NYC->LA should be within 4000 KM")
	}
	if IsWithinRadius(nycLat, nycLon, laLat, laLon, 1000, Kilometers) {
		t.Error("NYC->LA should NOT be within 1000 KM")
	}
}

func TestIsValidCoordinate(t *testing.T) {
	if !IsValidCoordinate(-6.200000, 106.816666) {
		t.Error("Jakarta should be valid")
	}
	if IsValidCoordinate(95.0, 106.816666) {
		t.Error("lat 95 should be invalid")
	}
	if IsValidCoordinate(-6.2, 200.0) {
		t.Error("lon 200 should be invalid")
	}
}

func jakartaBranches() []Location {
	return []Location{
		{Name: "Surabaya", Lat: -7.257472, Lon: 112.752088},
		{Name: "Bandung", Lat: -6.917464, Lon: 107.619123},
		{Name: "Medan", Lat: 3.595196, Lon: 98.672223},
	}
}

func TestFindNearest(t *testing.T) {
	n, ok := FindNearest(-6.200000, 106.816666, jakartaBranches(), Kilometers)
	if !ok {
		t.Fatal("expected a nearest result")
	}
	if n.Name != "Bandung" {
		t.Fatalf("nearest = %q, want Bandung", n.Name)
	}
	if n.Distance <= 0 {
		t.Errorf("nearest Distance should be populated, got %v", n.Distance)
	}
}

func TestFindNearestEmpty(t *testing.T) {
	if _, ok := FindNearest(0, 0, nil, Kilometers); ok {
		t.Error("empty list should return ok=false")
	}
}

func TestSortByDistanceAsc(t *testing.T) {
	sorted := SortByDistance(-6.200000, 106.816666, jakartaBranches(), Kilometers, Asc)
	want := []string{"Bandung", "Surabaya", "Medan"}
	for i, w := range want {
		if sorted[i].Name != w {
			t.Errorf("asc[%d] = %q, want %q", i, sorted[i].Name, w)
		}
	}
}

func TestSortByDistanceDesc(t *testing.T) {
	sorted := SortByDistance(-6.200000, 106.816666, jakartaBranches(), Kilometers, Desc)
	if sorted[0].Name != "Medan" {
		t.Errorf("desc[0] = %q, want Medan", sorted[0].Name)
	}
}

func TestSortByDistanceDoesNotMutateInput(t *testing.T) {
	in := jakartaBranches()
	_ = SortByDistance(-6.2, 106.8, in, Kilometers, Asc)
	if in[0].Name != "Surabaya" || in[0].Distance != 0 {
		t.Error("input slice must not be mutated")
	}
}

func TestCalculatorDefaultUnit(t *testing.T) {
	c := Calculator{DefaultUnit: Miles}
	got := c.Calculate(nycLat, nycLon, laLat, laLon, "")
	if math.Abs(math.Round(got)-2448) > 10 {
		t.Fatalf("Calculator default Miles: got %.2f, want ~2448", got)
	}
	// Explicit unit overrides the default.
	km := c.Calculate(nycLat, nycLon, laLat, laLon, Kilometers)
	if math.Abs(math.Round(km)-3940) > 10 {
		t.Fatalf("Calculator explicit KM: got %.2f, want ~3940", km)
	}
}
