package main

import (
	"github.com/vic-it/OSM-FMI/src/backend/util"
)

func readFile(path string) {
	util.ReadPBF(path)
}

func main() {
	PBFpath := "../../data/antarctica.osm.pbf"
	readFile(PBFpath)
}
