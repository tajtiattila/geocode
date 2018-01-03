package geoc

import (
	"path/filepath"
	"time"

	"github.com/kellydunn/golang-geo"
	"github.com/tajtiattila/basedir"
)

func NewGoogle(apikey string) Geocoder {
	if apikey == "" {
		panic("apikey must be specified")
	}
	geo.GoogleAPIKey = apikey
	return NewGeoGeocoder(new(geo.GoogleGeocoder))
}

func NewStdGoogle(apikey string) (Geocoder, error) {
	gc := NewGoogle(apikey)

	cachedir, err := basedir.Cache.EnsureDir("geocode", 0777)
	if err != nil {
		return nil, err
	}

	return NewCached(NewDelayed(gc, 150*time.Millisecond),
		filepath.Join(cachedir, "geocode-cache.leveldb"))
}
