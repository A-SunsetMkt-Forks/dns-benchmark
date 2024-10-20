//go:build release

package main

import (
	_ "embed"

	"github.com/oschwald/geoip2-golang"
)

//go:embed data/GeoLite2-Country.mmdb
var GeoData []byte

func GetGeoData() (*geoip2.Reader, error) {
	return geoip2.FromBytes(GeoData)
}
