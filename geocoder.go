package geoc

import (
	"io"
	"log"
)

type Geocoder interface {
	Geocode(query string) (lat, long float64, err error)
	Close() error
}

type Printlner interface {
	Println(v ...interface{})
}

var Log Printlner

func EnableLogging(w io.Writer) {
	Log = log.New(w, "geocoder", log.LstdFlags)
}

func reportError(what string, err error) {
	if Log == nil || err == nil {
		return
	}
	Log.Println(what+":", err)
}
