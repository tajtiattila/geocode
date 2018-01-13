// Package geocode is a simple interface to github.com/kellydunn/golang-geo.
package geocode

import (
	"io"
	"log"
)

// Geocoder looks up geographical locations.
type Geocoder interface {
	Geocode(query string) (Result, error)
	Close() error
}

// Result is a geocode result.
type Result struct {
	// Lat, Long is the most relevant position (when available)
	// or the centre of the boundary rectangle.
	Lat, Long float64

	// boundary rectangle
	North float64
	East  float64
	South float64
	West  float64
}

// Center returns the center of the boundary rectangle in r.
func (r Result) Center() (lat, long float64) {
	return rectMidPt(r.North, r.East, r.South, r.West)
}

// EmptyArea reports if r has zero area.
func (r Result) EmptyArea() bool {
	return r.North == r.South || r.East == r.West
}

func pointResult(lat, long float64) Result {
	return Result{
		Lat:   lat,
		North: lat,
		South: lat,

		Long: long,
		East: long,
		West: long,
	}
}

func rectResult(n, e, s, w float64) Result {
	lat, long := rectMidPt(n, e, s, w)
	return Result{
		Lat:  lat,
		Long: long,

		North: n,
		East:  e,
		South: s,
		West:  w,
	}
}

func rectMidPt(n, e, s, w float64) (lat, long float64) {
	if e < w {
		// spanning through ±180° longitude (international date line)
		e += 360
	}

	lat = (n + s) / 2
	long = (e + w) / 2

	if long > 180 {
		long -= 180
	}

	return
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
