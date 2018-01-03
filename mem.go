package geocode

type mem map[string]cacheEntry

// Memory returns a new Cache backed by process memory.
func Memory() QueryCache {
	return make(mem)
}

func (m mem) Load(query string) (lat, long float64, err error) {
	if e, ok := m[query]; ok {
		return e.Lat, e.Long, nil
	}
	return 0, 0, ErrCacheMiss
}

func (m mem) Store(query string, lat, long float64, err error) error {
	m[query] = cacheEntry{lat, long, err}
	return nil
}

func (mem) Close() error { return nil }
