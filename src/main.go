package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/vic-it/OSM-FMI/src/backend/util"
)

func fetchWorld(path string) {
	//loads whole osm.pbf world from filepath into basic format
	//basicWorld := util.PBFtoBASIC(path)
	nodes, edges, lonList, maxLat, minLat, maxDiff := util.GetWorld(path)
	fmt.Printf("Number of edges left: %d\n", len(nodes))

	//test here
	testpoints := [][]float64{{51, -78}, {-152, -68}, {153, -88}, {-54, -88}, {85, -88}}
	for _, node := range testpoints {
		relEdges, _ := util.GetRelevantEdges(node, lonList, maxLat, minLat, maxDiff)
		fmt.Printf("----------------------\n(%f/%f) crosses %d edges upward:\n", node[0], node[1], len(relEdges))
		for _, id := range relEdges {
			fmt.Printf("lon: (%f to %f)\nlat: (%f to %f)\n-\n", nodes[edges[id][0]][0], nodes[edges[id][1]][0], nodes[edges[id][0]][1], nodes[edges[id][1]][1])
		}
	}

	//util.BASICtoGEOJSONFile(util.PBFtoBASIC(path))
}

func main() {
	//"../../data/antarctica.osm.pbf"
	//"../../data/central-america.osm.pbf"
	// "../../data/global.sec" THIS IS THE BIG ONE FROM ILIAS (renamed, takes up ~11GB of RAM!)
	PBFpath := "../../data/antarctica.osm.pbf"
	fetchWorld(PBFpath)

	//generates grid around globe
	points, indexMatrix := util.GenerateGraphPoints(100)
	util.PointsToGEOJSONFile(points)
	//util.PrintEdgesToGEOJSON(util.GenerateEdges(points, indexMatrix))
	util.CalcEdgeDistances(util.GenerateEdges(points, indexMatrix))
	startServer()
}

func startServer() {
	println("\nstarting http server...")
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/getpoint", getPointHandler())
	http.ListenAndServe(":8080", nil)
}

func getPointHandler() http.HandlerFunc {
	pointHandler := func(writer http.ResponseWriter, request *http.Request) {
		urlQuery := request.URL.Query()
		inputLon, _ := strconv.ParseFloat(urlQuery["lon"][0], 64)
		inputLat, _ := strconv.ParseFloat(urlQuery["lat"][0], 64)

		fmt.Printf("click lon: %f\n", inputLon)
		fmt.Printf("click lat: %f\n", inputLat)

		//TODO GET CLOSEST POINT HERE
		lon, lat := util.GetClosestGridNode(inputLon, inputLat)

		outputString := fmt.Sprintf("%fx%f", lon, lat)
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(outputString))
	}

	return pointHandler
}

func getRouteHandler() {

}

// left[(lon1, edgeID1), (lon2, edgeID2)]
// right[(lon3, edgeID3), (lon4, edgeID2)]
// -> edgeID2 is a possibly rfelevant edge
