package distance

// Location is a named geographic point. It replaces the PHP package's
// associative-array items (['name'=>..., 'lat'=>..., 'lon'=>...]) with a
// typed struct, which is the idiomatic Go equivalent.
//
// After FindNearest / SortByDistance, the Distance field is populated with
// the computed distance from the reference point, in the requested unit.
type Location struct {
	Name     string
	Lat      float64
	Lon      float64
	Distance float64
}

// SortDirection controls the ordering of SortByDistance.
type SortDirection string

const (
	// Asc orders from nearest to farthest (default).
	Asc SortDirection = "asc"
	// Desc orders from farthest to nearest.
	Desc SortDirection = "desc"
)
