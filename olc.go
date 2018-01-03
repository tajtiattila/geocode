package geocode

import (
	"strings"

	"github.com/google/open-location-code/go"
)

// OpenLocationCode returns a Geocoder
// that understands Open Location Codes.
//
// It uses gc to find the reference of short open location codes,
// and to decode query strings that are not location codes.
func OpenLocationCode(gc Geocoder) Geocoder {
	return &olcCoder{
		gc: gc,
	}
}

type olcCoder struct {
	gc Geocoder
}

func (o *olcCoder) Close() error {
	return o.gc.Close()
}

func (o *olcCoder) Geocode(a string) (lat, long float64, err error) {
	code, ref, ok := splitCode(a)
	if ok {
		if ref != "" {
			rlat, rlong, err := o.gc.Geocode(ref)
			if err != nil {
				return rlat, rlong, err
			}

			code, err = olc.RecoverNearest(code, rlat, rlong)
		}

		ca, err := olc.Decode(code)
		if err != nil {
			return 0, 0, err
		}
		lat = (ca.LatLo + ca.LatHi) / 2
		long = (ca.LngLo + ca.LngHi) / 2
		return lat, long, nil
	}

	return o.gc.Geocode(a)
}

func splitCode(query string) (code, ref string, ok bool) {
	query = trim(query)

	if olc.CheckFull(query) == nil {
		return code, "", true
	}

	// check for code at beginning, eg.
	//  MQPQ+QG Kibera
	//  MQPQ+QG, Kibera
	i := strings.Index(code, " ")
	if i == -1 {
		return "", "", false
	}

	if code := trim(query[:i]); olc.Check(code) == nil {
		return code, trim(query[i:]), true
	}

	// check for code at end, eg.
	//  Belo Horizonte 22WM+PW
	//  Belo Horizonte, 22WM+PW
	i = strings.LastIndex(query, " ")
	if code := trim(query[i:]); olc.Check(code) == nil {
		return code, trim(query[:i]), true
	}

	return "", "", false
}

func trim(s string) string {
	s = strings.TrimSpace(s)
	return strings.TrimRightFunc(s, func(r rune) bool {
		return r == ',' || r == ';'
	})
}
