package geoc

import (
	"time"
)

// Delayed is a geocoder that limits outgoing requests
// based on a time.Duration
type Delayed struct {
	gc     Geocoder
	ticker *time.Ticker
}

func NewDelayed(gc Geocoder, d time.Duration) Geocoder {
	return &Delayed{gc, time.NewTicker(d)}
}

func (d *Delayed) Geocode(query string) (lat, long float64, err error) {
	<-d.ticker.C
	return d.gc.Geocode(query)
}

func (d *Delayed) Close() error {
	d.ticker.Stop()
	return nil
}
