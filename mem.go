package geocode

type mem map[string]cacheEntry

// Memory returns a new Cache backed by process memory.
func Memory() QueryCache {
	return make(mem)
}

func (m mem) Load(query string) (Result, error) {
	if e, ok := m[query]; ok {
		return e.Res, nil
	}
	return Result{}, ErrCacheMiss
}

func (m mem) Store(query string, res Result, err error) error {
	m[query] = cacheEntry{res, err}
	return nil
}

func (mem) Close() error { return nil }
