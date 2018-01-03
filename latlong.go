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

func (ll *latLongCoder) Geocode(a string) (lat, long float64, err error) {
	var rest string
	if _, err := fmt.Sscanf(a, "%f,%f%s", &lat, &long, &rest); err == nil && rest == "" {
		return lat, long, nil
	}
	if ll.gc == nil {
		return 0, 0, fmt.Errorf("unrecognized lat/long %q", a)
	}
	return ll.gc.Geocode(a)
}

func (ll *latLongCoder) Close() error { return ll.gc.Close() }
