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

func (o *olcCoder) Geocode(a string) (Result, error) {
	code, ref, ok := splitCode(a)
	if ok {
		if ref != "" {
			res, err := o.gc.Geocode(ref)
			if err != nil {
				return res, err
			}

			code, err = olc.RecoverNearest(code, res.Lat, res.Long)
		}

		ca, err := olc.Decode(code)
		if err != nil {
			return Result{}, err
		}
		return rectResult(ca.LatHi, ca.LngHi, ca.LatLo, ca.LngLo), nil
	}

	return o.gc.Geocode(a)
}

func splitCode(query string) (code, ref string, ok bool) {
	query = trim(query)

	if olc.CheckFull(query) == nil {
		return query, "", true
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
