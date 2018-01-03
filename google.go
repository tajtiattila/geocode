package geocode

import (
	"path/filepath"
	"time"

	"github.com/kellydunn/golang-geo"
	"github.com/tajtiattila/basedir"
)

// NewGoogle returns a geocoder using the Google geocode API.
//
// It panics if apikey is empty.
func NewGoogle(apikey string) Geocoder {
	if apikey == "" {
		panic("apikey must be specified")
	}
	geo.GoogleAPIKey = apikey
	return NewGeoGeocoder(new(geo.GoogleGeocoder))
}

// NewGoogle returns a caching geocoder
// using the delayed calls to the Google geocode API.
//
// It panics if apikey is empty.
func NewStdGoogle(apikey string) (Geocoder, error) {
	gc := NewGoogle(apikey)

	cachedir, err := basedir.Cache.EnsureDir("geocode", 0777)
	if err != nil {
		return nil, err
	}

	return NewCached(
		NewDelayed(gc, 150*time.Millisecond),
		filepath.Join(cachedir, "geocode-cache.leveldb"))
}
