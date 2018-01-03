package geocode

import "errors"

// QueryCache caches geocode results.
type QueryCache interface {
	Load(query string) (lat, long float64, err error)
	Store(query string, lat, long float64, err error) error

	Close() error
}

// ErrCacheMiss is returned by Cache.Load when
// a query string is not cached yet.
var ErrCacheMiss = errors.New("cache miss")

type cacheEntry struct {
	Lat  float64
	Long float64
	Err  error `json:",omitempty"`
}

// Cache returns a geocoder that caches results from gc.
//
// The cache is stored at path.
func Cache(gc Geocoder, qc QueryCache) Geocoder {
	return &cached{gc, qc}
}

type cached struct {
	gc Geocoder
	qc QueryCache
}

func (g *cached) Close() error {
	err1 := g.gc.Close()
	err2 := g.qc.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

func (g *cached) Geocode(a string) (lat, long float64, err error) {
	lat, long, err = g.qc.Load(a)
	if err != ErrCacheMiss {
		return lat, long, err
	}

	lat, long, err = g.gc.Geocode(a)
	g.qc.Store(a, lat, long, err)
	return lat, long, err
}
