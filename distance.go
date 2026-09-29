// Package distance calculates geographical distance between GPS coordinates,
// checks radius bounds, finds the nearest location, and sorts location lists.
//
// It is an idiomatic Go port of the PHP/Laravel package
// alkauni/distance-coordinate-calculator. The numeric results match the PHP
// implementation exactly, including its unit-conversion factors.
//
// Two ways to use it:
//
//	// Package-level functions (quick use, default unit = KM):
//	km := distance.Calculate(-6.2, 106.816666, -6.917464, 107.619123, "")
//
//	// Calculator value (carries a default unit, like the Laravel config):
//	c := distance.Calculator{DefaultUnit: distance.Miles}
//	mi := c.Calculate(-6.2, 106.816666, -6.917464, 107.619123, "")
package distance

import (
	"math"
	"sort"
	"strconv"
	"strings"
)

// Calculate returns the distance between two points (lat1,lon1)-(lat2,lon2)
// in the given unit. An empty unit falls back to DefaultUnit (KM).
//
// This reproduces the source package's formula precisely: it is the
// spherical law-of-cosines variant scaled through miles, NOT the textbook
// Haversine, so outputs are bit-for-bit comparable with the PHP version.
func Calculate(lat1, lon1, lat2, lon2 float64, unit Unit) float64 {
	unit = resolveUnit(unit)

	theta := lon1 - lon2
	dist := math.Sin(deg2rad(lat1))*math.Sin(deg2rad(lat2)) +
		math.Cos(deg2rad(lat1))*math.Cos(deg2rad(lat2))*math.Cos(deg2rad(theta))
	dist = math.Acos(dist)
	dist = rad2deg(dist)
	miles := dist * 60 * 1.1515

	return milesTo(miles, unit)
}

// CalculateFormatted returns the distance as a string rounded to decimals,
// with the resolved unit code appended, e.g. "118.44 KM".
func CalculateFormatted(lat1, lon1, lat2, lon2 float64, unit Unit, decimals int) string {
	unit = resolveUnit(unit)
	d := Calculate(lat1, lon1, lat2, lon2, unit)
	if decimals < 0 {
		decimals = 0
	}
	return strconv.FormatFloat(d, 'f', decimals, 64) + " " + strings.ToUpper(string(unit))
}

// IsWithinRadius reports whether point 2 is within radius of point 1,
// measured in the given unit.
func IsWithinRadius(lat1, lon1, lat2, lon2, radius float64, unit Unit) bool {
	return Calculate(lat1, lon1, lat2, lon2, unit) <= radius
}

// IsValidCoordinate reports whether lat/lon fall within valid geographic
// bounds: latitude [-90, 90], longitude [-180, 180].
func IsValidCoordinate(lat, lon float64) bool {
	return lat >= -90.0 && lat <= 90.0 && lon >= -180.0 && lon <= 180.0
}

// FindNearest returns the location closest to (lat, lon), with its Distance
// field populated. It returns ok=false when the list is empty.
func FindNearest(lat, lon float64, locations []Location, unit Unit) (nearest Location, ok bool) {
	sorted := SortByDistance(lat, lon, locations, unit, Asc)
	if len(sorted) == 0 {
		return Location{}, false
	}
	return sorted[0], true
}

// SortByDistance returns a copy of locations sorted by distance from
// (lat, lon), each element's Distance field populated in the given unit.
// direction Asc orders nearest-first (default), Desc farthest-first.
// The input slice is not mutated.
func SortByDistance(lat, lon float64, locations []Location, unit Unit, direction SortDirection) []Location {
	result := make([]Location, len(locations))
	copy(result, locations)

	for i := range result {
		result[i].Distance = Calculate(lat, lon, result[i].Lat, result[i].Lon, unit)
	}

	desc := strings.EqualFold(string(direction), string(Desc))
	sort.SliceStable(result, func(i, j int) bool {
		if desc {
			return result[i].Distance > result[j].Distance
		}
		return result[i].Distance < result[j].Distance
	})

	return result
}

// Calculator carries a default unit, mirroring the Laravel config's
// default_unit. Methods with an empty unit argument fall back to
// c.DefaultUnit (and then to the package DefaultUnit if that is empty too).
type Calculator struct {
	DefaultUnit Unit
}

func (c Calculator) unit(u Unit) Unit {
	if strings.TrimSpace(string(u)) != "" {
		return u
	}
	if strings.TrimSpace(string(c.DefaultUnit)) != "" {
		return c.DefaultUnit
	}
	return DefaultUnit
}

// Calculate is the Calculator-scoped form of the package Calculate.
func (c Calculator) Calculate(lat1, lon1, lat2, lon2 float64, unit Unit) float64 {
	return Calculate(lat1, lon1, lat2, lon2, c.unit(unit))
}

// CalculateFormatted is the Calculator-scoped form of CalculateFormatted.
func (c Calculator) CalculateFormatted(lat1, lon1, lat2, lon2 float64, unit Unit, decimals int) string {
	return CalculateFormatted(lat1, lon1, lat2, lon2, c.unit(unit), decimals)
}

// IsWithinRadius is the Calculator-scoped form of IsWithinRadius.
func (c Calculator) IsWithinRadius(lat1, lon1, lat2, lon2, radius float64, unit Unit) bool {
	return IsWithinRadius(lat1, lon1, lat2, lon2, radius, c.unit(unit))
}

// FindNearest is the Calculator-scoped form of FindNearest.
func (c Calculator) FindNearest(lat, lon float64, locations []Location, unit Unit) (Location, bool) {
	return FindNearest(lat, lon, locations, c.unit(unit))
}

// SortByDistance is the Calculator-scoped form of SortByDistance.
func (c Calculator) SortByDistance(lat, lon float64, locations []Location, unit Unit, direction SortDirection) []Location {
	return SortByDistance(lat, lon, locations, c.unit(unit), direction)
}

func deg2rad(d float64) float64 { return d * math.Pi / 180 }
func rad2deg(r float64) float64 { return r * 180 / math.Pi }
