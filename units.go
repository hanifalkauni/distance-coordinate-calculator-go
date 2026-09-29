package distance

import "strings"

// Unit is a geographic distance unit supported by the calculator.
type Unit string

const (
	// Kilometers is the default unit.
	Kilometers Unit = "KM"
	// Miles is the base unit of the underlying formula.
	Miles Unit = "Miles"
	// NauticalMiles is nautical miles.
	NauticalMiles Unit = "N"
	// Meters is meters.
	Meters Unit = "M"
)

// DefaultUnit is used when an empty unit is supplied, matching the PHP
// package's config fallback (`default_unit => 'KM'`).
const DefaultUnit = Kilometers

// milesTo converts a distance expressed in miles (the formula's base unit)
// into the requested unit. It mirrors the PHP switch exactly so results are
// numerically identical, including its non-standard nautical-mile factor.
func milesTo(miles float64, unit Unit) float64 {
	switch strings.ToUpper(string(unit)) {
	case "KM":
		return miles * 1.609344
	case "N":
		return miles * 0.8684
	case "M":
		return miles * 1.609344 * 1000
	case "MILES":
		return miles
	default:
		return miles
	}
}

// resolveUnit returns the given unit, or DefaultUnit when it is empty,
// mirroring the PHP resolveUnit() fallback behaviour.
func resolveUnit(unit Unit) Unit {
	if strings.TrimSpace(string(unit)) == "" {
		return DefaultUnit
	}
	return unit
}
