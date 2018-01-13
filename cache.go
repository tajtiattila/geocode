package geocode

import "errors"

// QueryCache caches geocode results.
type QueryCache interface {
	Load(query string) (Result, error)
	Store(query string, res Result, err error) error

	Close() error
}

// ErrCacheMiss is returned by Cache.Load when
// a query string is not cached yet.
var ErrCacheMiss = errors.New("cache miss")

type cacheEntry struct {
	Res Result `json:"result"`
	Err error  `json:",omitempty"`
}

// Cache returns a geocoder that caches results from gc.
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

func (g *cached) Geocode(a string) (Result, error) {
	r, err := g.qc.Load(a)
	if err != ErrCacheMiss {
		return r, err
	}

	r, err = g.gc.Geocode(a)
	g.qc.Store(a, r, err)
	return r, err
}
