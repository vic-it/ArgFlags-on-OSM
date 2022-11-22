package main

import (
	"github.com/vic-it/OSM-FMI/src/backend/util"
)

func fetchWorld(path string) {
	//loads whole osm.pbf world from filepath into basic format
	//basicWorld := util.PBFtoBASIC(path)
	util.BASICtoGEOJSONFile(util.PBFtoBASIC(path))
}

func main() {
	//"../../data/antarctica.osm.pbf"
	//"../../data/central-america.osm.pbf"
	// "../../data/global.sec" THIS IS THE BIG ONE FROM ILIAS (renamed, takes up ~11GB of RAM!)
	//PBFpath := "../../data/antarctica.osm.pbf"
	//fetchWorld(PBFpath)

	//generates grid around globe
	points, indexMatrix := util.GenerateGraphPoints(500)
	util.PointsToGEOJSONFile(points)
	//util.PrintEdgesToGEOJSON(util.GenerateEdges(points, indexMatrix))
	util.CalcEdgeDistances(util.GenerateEdges(points, indexMatrix))
}
