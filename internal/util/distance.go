package util

import "math"

type Unit int

const (
	Meters Unit = iota + 1
	Kilometers
	Miles
	NM // ???
)

func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

func rad2deg(rad float64) float64 {
	return rad * 180.0 / math.Pi
}

func Distance(lat1, lon1, lat2, lon2 float64, unit Unit) float64 {
	lat1 = deg2rad(lat1)
	lon1 = deg2rad(lon1)
	lat2 = deg2rad(lat2)
	lon2 = deg2rad(lon2)

	// dist in radians
	dist := math.Acos(
		math.Sin(lat1)*math.Sin(lat2) +
			math.Cos(lat1)*math.Cos(lat2)*math.Cos(lon1-lon2),
	)

	// dist in miles
	dist = rad2deg(dist) * 60 * 1.1515

	// convert if necessary
	switch unit {
	case Kilometers:
		dist = dist * 1.609344
	case Meters:
		dist = dist * 1609.344
	case NM:
		dist = dist * 0.8684
	}

	return dist
}
