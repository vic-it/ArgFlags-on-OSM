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

	//test here
	// intersection probably a bottle neck -> only look for maxlat minlat after the longitude test
	// testpoints := [][]float64{{51, -78}, {-152, -68}, {153, -88}, {-54, -88}, {-62.8771, -64.6543}}
	// for _, node := range testpoints {
	// 	relEdges, maybe := util.GetRelevantEdges(node, coastline)
	// 	fmt.Printf("----------------------\n(%f/%f) crosses %d edges upward:\n", node[0], node[1], len(relEdges))
	// 	for _, id := range relEdges {
	// 		fmt.Printf("lon: (%f to %f)\nlat: (%f to %f)\n-\n", nodes[edges[id][0]][0], nodes[edges[id][1]][0], nodes[edges[id][0]][1], nodes[edges[id][1]][1])
	// 	}
	// 	fmt.Printf("%d 'maybe' edges\n", len(maybe))
	// }

	//util.BASICtoGEOJSONFile(util.PBFtoBASIC(path))
}

func main() {
	testPointInWater()
	//start()
}

func start() {

	//"../../data/antarctica.osm.pbf"
	//"../../data/central-america.osm.pbf"
	// "../../data/global.sec" THIS IS THE BIG ONE FROM ILIAS (renamed, takes up ~11GB of RAM!)
	//points, indexMatrix, pointInWaterMatrix := util.GenerateGraphPoints(500, coastline)
	//util.PointsToGEOJSONFile(points)
	//util.PrintEdgesToGEOJSON(util.GenerateEdges(points, indexMatrix, pointInWaterMatrix))
	//util.CalcEdgeDistances(util.GenerateEdges(points, indexMatrix))
	startServer()
}

func testPointInWater() {
	path := "../../data/antarctica.osm.pbf"
	coastline := util.GetCoastline(path)
	//generates grid around globe
	var testPoints [][]float64
	testPoints = append(testPoints, []float64{-180, -85.5})
	testPoints = append(testPoints, []float64{-70, -15.5})
	testPoints = append(testPoints, []float64{50, 111})
	testPoints = append(testPoints, []float64{120, -65.5})
	for _, p := range testPoints {
		util.IsPointInWater(p, coastline)
	}

	// points, indexMatrix, pointInWaterMatrix := util.GenerateGraphPoints(5000, coastline)
	// util.PointsToGEOJSONFile(points)
	// p, src, dest := util.GenerateEdges(points, indexMatrix, pointInWaterMatrix)
	// util.CalcEdgeDistances(p, src, dest)
	// util.PrintEdgesToGEOJSON(p, src, dest)
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
