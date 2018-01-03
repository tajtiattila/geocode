package geocode

import (
	"time"
)

type delayed struct {
	gc     Geocoder
	ticker *time.Ticker
}

// NewDelayed delays requests sent to gc.
//
// It ensures requests are at least d apart.
func NewDelayed(gc Geocoder, d time.Duration) Geocoder {
	return &delayed{gc, time.NewTicker(d)}
}

func (d *delayed) Geocode(query string) (lat, long float64, err error) {
	<-d.ticker.C
	return d.gc.Geocode(query)
}

func (d *delayed) Close() error {
	d.ticker.Stop()
	return d.gc.Close()
}
