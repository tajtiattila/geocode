package geocode

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Google returns a geocoder using the Google geocode API.
//
// The apikey argument may be empty.
func Google(apikey string) Geocoder {
	if apikey == "" {
		panic("apikey must be specified")
	}
	return &google{
		apikey: apikey,
	}
}

// StdGoogle returns a geocoder
// using delayed calls to the Google geocode API.
//
// The apikey argument may be empty.
func StdGoogle(apikey string) Geocoder {
	gc := Google(apikey)

	return Delay(gc, 150*time.Millisecond)
}

type google struct {
	apikey string
}

func (g *google) Close() error { return nil }

func (g *google) Geocode(query string) (Result, error) {
	u, err := url.Parse("https://maps.googleapis.com/maps/api/geocode/json")
	if err != nil {
		panic("impossible")
	}
	v := make(url.Values)

	v.Set("address", query)
	v.Set("sensor", "false")
	if g.apikey != "" {
		v.Set("key", g.apikey)
	}
	u.RawQuery = v.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	return decodeGoogleResponse(resp.Body)
}

type googleResTop struct {
	Results []*googleResult `json:"results"`
}

type googleResult struct {
	Geometry *googleGeometry `json:"geometry"`
}

type googleGeometry struct {
	Location *googleLocation `json:"location"`
	Bounds   *googleRect     `json:"bounds"`
	Viewport *googleRect     `json:"viewport"`
}

type googleRect struct {
	NE googleLocation `json:"northeast"`
	SW googleLocation `json:"southwest"`
}

type googleLocation struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"lng"`
}

func decodeGoogleResponse(r io.Reader) (Result, error) {
	// logic from perkeep.org/internal/geocode/geocode.go
	var resTop googleResTop
	if err := json.NewDecoder(r).Decode(&resTop); err != nil {
		return Result{}, err
	}

	for _, res := range resTop.Results {
		if res.Geometry != nil && res.Geometry.Bounds != nil {
			r := res.Geometry.Bounds
			if r.NE.Lat == 90 && r.NE.Long == 180 &&
				r.SW.Lat == -90 && r.SW.Long == -180 {
				// Google sometimes returns a "whole world" rect for large addresses (like "USA")
				// so instead use the viewport in that case.
				return googleRes(res.Geometry.Location, *res.Geometry.Viewport)
			} else {
				return googleRes(res.Geometry.Location, *r)
			}
		}
	}

	return Result{}, errors.New("geocode/google: empty result set")
}

func googleRes(l *googleLocation, r googleRect) (Result, error) {
	x := Result{
		North: r.NE.Lat,
		East:  r.NE.Long,
		South: r.SW.Lat,
		West:  r.SW.Long,
	}
	if l != nil {
		x.Lat, x.Long = l.Lat, l.Long
	} else {
		x.Lat, x.Long = rectMidPt(x.North, x.East, x.South, x.West)
	}
	return x, nil
}
