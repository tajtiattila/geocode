package geoc

import (
	"encoding/json"

	"github.com/syndtr/goleveldb/leveldb"
)

// geocoder is a caching geocoder
type Cached struct {
	db *leveldb.DB
	gc Geocoder
}

func NewCached(gc Geocoder, cachepath string) (*Cached, error) {
	db, err := leveldb.OpenFile(cachepath, nil)
	if err != nil {
		return nil, err
	}
	return &Cached{db, gc}, nil
}

func (g *Cached) Close() error {
	err1 := g.db.Close()
	err2 := g.gc.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

func (g *Cached) Geocode(a string) (lat, long float64, err error) {
	type cacheEnt struct {
		Lat  float64
		Long float64
		Err  error `json:",omitempty"`
	}

	data, err := g.db.Get([]byte(a), nil)
	switch {
	case err == nil:
		var ce cacheEnt
		if err = json.Unmarshal(data, &ce); err == nil {
			return ce.Lat, ce.Long, ce.Err
		}
		reportError("unmarshal", err)
	case err != leveldb.ErrNotFound:
		reportError("leveldb get", err)
	}

	lat, long, geoerr := g.gc.Geocode(a)

	data, err = json.Marshal(cacheEnt{lat, long, geoerr})
	if err != nil {
		panic("impossible")
	}
	if err = g.db.Put([]byte(a), data, nil); err != nil {
		reportError("leveldb put", err)
	}

	return lat, long, geoerr
}
