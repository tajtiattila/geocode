package geocode

import (
	"time"

	"github.com/kellydunn/golang-geo"
)

// Google returns a geocoder using the Google geocode API.
//
// It panics if apikey is empty.
func Google(apikey string) Geocoder {
	if apikey == "" {
		panic("apikey must be specified")
	}
	geo.GoogleAPIKey = apikey
	return Geo(new(geo.GoogleGeocoder))
}

// StdGoogle returns a geocoder
// using delayed calls to the Google geocode API.
//
// It panics if apikey is empty.
func StdGoogle(apikey string) Geocoder {
	gc := Google(apikey)

	return Delay(gc, 150*time.Millisecond)
}
