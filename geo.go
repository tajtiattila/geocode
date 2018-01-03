package geoc

import "github.com/kellydunn/golang-geo"

func NewGeoGeocoder(gc geo.Geocoder) Geocoder {
	return &geoGeocoder{gc}
}

type geoGeocoder struct {
	gc geo.Geocoder
}

func (g *geoGeocoder) Geocode(query string) (lat, long float64, err error) {
	p, err := g.gc.Geocode(query)
	if err != nil || p == nil {
		return 0, 0, err
	}
	return p.Lat(), p.Lng(), err
}

func (*geoGeocoder) Close() error {
	return nil
}
