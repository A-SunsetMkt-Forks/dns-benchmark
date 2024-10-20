//go:build !release

package main

import "github.com/oschwald/geoip2-golang"

const geoDataPath = "./data/GeoLite2-Country.mmdb"

func GetGeoData() (*geoip2.Reader, error) {
	return geoip2.Open(geoDataPath)
}
