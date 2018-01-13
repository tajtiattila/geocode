package geocode

import (
	"time"
)

type delayed struct {
	gc     Geocoder
	ticker *time.Ticker
}

// Delay returns a Geocoder that delays requests sent to gc.
//
// It ensures requests are at least d apart.
func Delay(gc Geocoder, d time.Duration) Geocoder {
	if d <= 0 {
		return gc
	}
	return &delayed{gc, time.NewTicker(d)}
}

func (d *delayed) Geocode(query string) (Result, error) {
	<-d.ticker.C
	return d.gc.Geocode(query)
}

func (d *delayed) Close() error {
	d.ticker.Stop()
	return d.gc.Close()
}
