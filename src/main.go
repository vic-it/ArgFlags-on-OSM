package main

import (
	"github.com/vic-it/OSM-FMI/src/backend/util"
)

func fetchWorld(path string) {
	//loads whole osm.pbf world from filepath into basic format
	//basicWorld := util.PBFtoBASIC(path)
	util.PBFtoBASIC(path)
}

func main() {
	PBFpath := "../../data/antarctica.osm.pbf"
	fetchWorld(PBFpath)
}
