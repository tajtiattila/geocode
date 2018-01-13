package geocode

import "fmt"

// LatLong returns a geocoder that simply decodes
// strings with latitude and longitude values.
//
// It calls gc when not nil and it cannot decode the latitude and longitude.
func LatLong(gc Geocoder) Geocoder {
	return &latLongCoder{gc}
}

type latLongCoder struct {
	gc Geocoder
}

func (ll *latLongCoder) Geocode(a string) (Result, error) {
	var lat, long float64
	var rest string
	if n, _ := fmt.Sscanf(a, "%f,%f%s", &lat, &long, &rest); n == 2 {
		return pointResult(lat, long), nil
	}
	if ll.gc == nil {
		return Result{}, fmt.Errorf("unrecognized lat/long %q", a)
	}
	return ll.gc.Geocode(a)
}

func (ll *latLongCoder) Close() error { return ll.gc.Close() }
