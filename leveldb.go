package geocode

import (
	"encoding/json"

	"github.com/syndtr/goleveldb/leveldb"
)

type levelDB struct {
	db *leveldb.DB
}

// LevelDB returns a QueryCache stored at path.
func LevelDB(path string) (QueryCache, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &levelDB{db}, nil
}

func (l *levelDB) Load(query string) (lat, long float64, err error) {
	data, err := l.db.Get([]byte(query), nil)
	switch {
	case err == nil:
		var ce cacheEntry
		if err = json.Unmarshal(data, &ce); err == nil {
			return ce.Lat, ce.Long, ce.Err
		}
		reportError("unmarshal", err)
	case err == leveldb.ErrNotFound:
		// pass
	default:
		reportError("leveldb get", err)
	}
	return 0, 0, ErrCacheMiss
}

func (l *levelDB) Store(query string, lat, long float64, err error) error {
	data, err := json.Marshal(cacheEntry{lat, long, err})
	if err != nil {
		panic("impossible")
	}
	err = l.db.Put([]byte(query), data, nil)
	if err != nil {
		reportError("leveldb put", err)
	}
	return err
}

func (l *levelDB) Close() error {
	return l.db.Close()
}
