// Package geocode is a simple interface to github.com/kellydunn/golang-geo.
package geocode

import (
	"io"
	"log"
)

// Geocode geographical locations.
type Geocoder interface {
	Geocode(query string) (lat, long float64, err error)
	Close() error
}

type printlner interface {
	Println(v ...interface{})
}

var logger printlner

// EnableLogging enables error logging on w.
func EnableLogging(w io.Writer) {
	logger = log.New(w, "geocoder", log.LstdFlags)
}

func reportError(what string, err error) {
	if logger == nil || err == nil {
		return
	}
	logger.Println(what+":", err)
}
