# Distance Coordinate Calculator (Go)

An idiomatic Go port of the PHP/Laravel package
[`alkauni/distance-coordinate-calculator`](https://github.com/hanifalkauni/distance-coordinate-calculator-package).
Calculate geographical distance between GPS coordinates, check radius bounds,
find the nearest location, and sort location lists.

Numeric results match the original PHP package exactly, including its
unit-conversion factors.

## Install

```bash
go get github.com/hanifalkauni/distance-coordinate-calculator-go
```

```go
import distance "github.com/hanifalkauni/distance-coordinate-calculator-go"
```

## Units

| Constant                 | Code      | Description             |
|--------------------------|-----------|-------------------------|
| `distance.Kilometers`    | `KM`      | Kilometers *(default)*  |
| `distance.Miles`         | `Miles`   | Miles                   |
| `distance.NauticalMiles` | `N`       | Nautical miles          |
| `distance.Meters`        | `M`       | Meters                  |

Passing an empty unit (`""`) falls back to `KM`.

## Usage

### 1. Basic distance

```go
km := distance.Calculate(-6.200000, 106.816666, -6.917464, 107.619123, distance.Kilometers)
mi := distance.Calculate(40.748817, -73.985428, 34.052235, -118.243683, distance.Miles)
```

### 2. Formatted output

```go
s := distance.CalculateFormatted(-6.2, 106.816666, -6.917464, 107.619123, distance.Kilometers, 2)
// "119.25 KM"
```

### 3. Radius check (geofencing)

```go
ok := distance.IsWithinRadius(-6.2, 106.816666, -6.205, 106.818, 1.0, distance.Kilometers)
// true
```

### 4. Find nearest

```go
branches := []distance.Location{
    {Name: "Surabaya", Lat: -7.257472, Lon: 112.752088},
    {Name: "Bandung", Lat: -6.917464, Lon: 107.619123},
    {Name: "Medan", Lat: 3.595196, Lon: 98.672223},
}

nearest, ok := distance.FindNearest(-6.2, 106.816666, branches, distance.Kilometers)
// nearest.Name == "Bandung", nearest.Distance populated, ok == true
```

### 5. Sort by distance

```go
sorted := distance.SortByDistance(-6.2, 106.816666, branches, distance.Kilometers, distance.Asc)
// nearest first; each element's Distance field is populated.
// Use distance.Desc for farthest first. The input slice is not mutated.
```

### 6. Validate coordinates

```go
distance.IsValidCoordinate(-6.2, 106.816666) // true
distance.IsValidCoordinate(95.0, 200.0)      // false
```

### Calculator with a default unit

The Laravel package reads a default unit from config. In Go, use a
`Calculator` value instead:

```go
c := distance.Calculator{DefaultUnit: distance.Miles}
mi := c.Calculate(40.748817, -73.985428, 34.052235, -118.243683, "") // uses Miles
km := c.Calculate(40.748817, -73.985428, 34.052235, -118.243683, distance.Kilometers) // explicit unit wins
```

## Porting notes (Laravel → Go)

| Laravel concept                         | Go equivalent                              |
|-----------------------------------------|--------------------------------------------|
| `Distance` Facade / ServiceProvider     | Package-level functions                    |
| `config('distance.default_unit')`       | `Calculator{DefaultUnit: ...}`             |
| Associative array `['lat'=>,'lon'=>]`   | `Location` struct                          |
| `findNearest` returns `null`            | `(Location, bool)` — `ok=false` when empty |
| exceptions                              | idiomatic `bool`/`error` returns           |

## Testing

```bash
go test ./...
```

## License

MIT. See [LICENSE](LICENSE).
