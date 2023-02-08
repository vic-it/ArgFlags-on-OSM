package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"

	"github.com/vic-it/OSM-FMI/src/backend/util"
)

// the graph!
var graph util.Graph
var arcData util.ArcData
var nodePartitionList []int

func main() {
	fmt.Printf("Starting")
	// either creates a new graph entirely (can take some time) or imports a preprocessed graph (very fast)
	initialize()
	// starts server and web interface -> reachable at localhost:8080
	testStuff()
	startServer()
}

func startServer() {
	println("Starting HTTP server...")
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/getpoint", getPointHandler())
	http.HandleFunc("/getroute", getRouteHandler())
	http.HandleFunc("/directory", getDirectoryHandler())
	http.HandleFunc("/preprocess", getPreprocessingHandler())
	println("GUI available on: localhost:8080")
	println("READY!\n-----------")
	http.ListenAndServe(":8080", nil)

}

// checks for nearest valid neighbor node in water and returns this to the web interface
func getPointHandler() http.HandlerFunc {
	pointHandler := func(writer http.ResponseWriter, request *http.Request) {
		urlQuery := request.URL.Query()
		inputLon, _ := strconv.ParseFloat(urlQuery["lon"][0], 64)
		inputLat, _ := strconv.ParseFloat(urlQuery["lat"][0], 64)

		// fmt.Printf("click lon: %f\n", inputLon)
		// fmt.Printf("click lat: %f\n", inputLat)
		nodeID := util.GetClosestGridNode(inputLon, inputLat, graph)

		outputString := fmt.Sprintf("%fx%fx%d", graph.Nodes[nodeID][0], graph.Nodes[nodeID][1], nodeID)
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(outputString))
	}

	return pointHandler
}

// calculates the shortest path between two nodes via dijkstra
func getRouteHandler() http.HandlerFunc {
	routingHandler := func(writer http.ResponseWriter, request *http.Request) {
		urlQuery := request.URL.Query()
		src, _ := strconv.ParseInt(urlQuery["src"][0], 0, 64)
		dest, _ := strconv.ParseInt(urlQuery["dest"][0], 0, 64)
		mode, _ := strconv.ParseInt(urlQuery["mode"][0], 0, 64)

		// fmt.Printf("source node ID: %d\n", src)
		// fmt.Printf("destination node ID: %d\n", dest)

		writer.WriteHeader(http.StatusOK)
		//calc route below
		writer.Write([]byte(getRouteString(src, dest, mode)))
	}
	return routingHandler
}

// should return a string in the following format:
// if no route found -> -1y0
// if route found -> distanceylon1zlat1xlon2zlat2x...
func getRouteString(src int64, dest int64, mode int64) string {
	var distance int
	var path []int
	var nodesPopped int
	var initTime float64
	var searchTime float64
	if mode == 0 {
		distance, path, initTime, searchTime, nodesPopped = util.CalculateDijkstra(graph, int(src), int(dest))
	} else {
		distance, path, initTime, searchTime, nodesPopped = util.CalculateArcFlagDijkstra(graph, int(src), int(dest), arcData, nodePartitionList)
	}
	output := fmt.Sprintf("%dy%dy%.3fy%.3fy", distance, nodesPopped, initTime, searchTime)
	for i, nodeID := range path {
		if i == len(path)-1 {
			output = fmt.Sprintf("%s%fz%f", output, graph.Nodes[nodeID][0], graph.Nodes[nodeID][1])
		} else {
			output = fmt.Sprintf("%s%fz%fx", output, graph.Nodes[nodeID][0], graph.Nodes[nodeID][1])
		}
	}
	return output
}

// NOT IN USE!
func getPreprocessingHandler() http.HandlerFunc {
	proprocessHandler := func(writer http.ResponseWriter, request *http.Request) {
		urlQuery := request.URL.Query()
		processType, _ := strconv.ParseFloat(urlQuery["process"][0], 64)
		processOption, _ := strconv.ParseFloat(urlQuery["option"][0], 64)

		fmt.Printf("process type: %f\n", processType)
		fmt.Printf("option: %f\n--\n", processOption)

		outputString := "lol hier text"
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(outputString))
	}
	return proprocessHandler
}

// NOT IN USE!
func getDirectoryHandler() http.HandlerFunc {
	directoryHandler := func(writer http.ResponseWriter, request *http.Request) {
		urlQuery := request.URL.Query()
		processType, _ := strconv.ParseFloat(urlQuery["process"][0], 64)
		processOption, _ := strconv.ParseFloat(urlQuery["option"][0], 64)

		fmt.Printf("process type: %f\n", processType)
		fmt.Printf("option: %f\n--\n", processOption)

		outputString := "lol hier text"
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(outputString))
	}
	return directoryHandler
}

// NOT IN USE!
func analyzeDirectory(basePath string) {
	println("Reading directory...")

}

// creates a graph with the coastlines from the path and roughly the number of nodes
func createGraph(pathToCoastlinesPBF string, numberOfNodes int) {
	graph = util.Graph{}
	runtime.GC()
	graph = util.GenerateGraph(numberOfNodes, util.GetCoastline(pathToCoastlinesPBF))
}

func preparePartitionList() {
	for _, row := range arcData.NodePartitionMatrix {
		nodePartitionList = append(nodePartitionList, row...)
	}
}

// initializes a graph either by importing it from a file or by creating one (creating can take a long time)
func initialize() {
	// relevant paths
	graphPath := "../../data/graph.graph"
	arcFlagPath := "../../data/arc.flags"
	antarctica := "../../data/antarctica.osm.pbf"
	global := "../../data/global.sec"
	// prints "..." so we dont have to comment/uncomment all paths because go is weird like that
	fmt.Printf("%s%s%s%s\n", antarctica[0:1], global[0:1], graphPath[0:1], arcFlagPath[0:1])

	// CREATE NEW GRAPH BY UNCOMMENTING BELOW:
	//-----------------------------------------------------
	//createGraph(antarctica, 100000)
	//util.GraphToFile(graph, graphPath)
	//-----------------------------------------------------

	// IMPORT GRAPH BY UNCOMMENTING BELOW:
	//-----------------------------------------------------
	graph = util.FileToGraph(graphPath)
	//-----------------------------------------------------

	// PRINT TO GEOJSON BY UNCOMMENTING BELOW:
	//-----------------------------------------------------
	//util.PrintPointsToGEOJSON(graph)
	//util.PrintEdgesToGEOJSON(graph)
	//-----------------------------------------------------

	//arcFlagStuff
	// GENERATE NEW ARCFLAGS BY UNCOMMENTING BELOW
	// 7 - 3 generates roughly square partitions (64 of them)
	arcData = util.PreprocessArcFlags(graph, 7, 3)
	util.ArcFlagsToFile(arcData, arcFlagPath)

	// IMPORT ARCFLAGS BY UNCOMMENTING BELOW:
	//arcData = util.FileToArcFlags(arcFlagPath)
	//this speeds up arc flag since it doesnt have to calculate row/col of nodepartitionmatrix anymore
	preparePartitionList()
}

func testStuff() {
	//util.PrintPointsToGEOJSON2(graph, arcData.NodePartitionMatrix)
	// for _, row := range graph.NodeMatrix {
	// 	fmt.Printf("first lon: %3.3f - second lon: %3.3f\n", graph.Nodes[row[0]][0], graph.Nodes[row[1]][0])
	// }
	// for _, line := range arcFlags {
	// 	fmt.Printf("%v\n", line)
	// }
	util.TestAlgorithms(graph, arcData, 3000, nodePartitionList)
	// for _, row := range  util.PreprocessArcFlags(graph, 8, 1){
	// 	fmt.Printf("[")
	// 	for _, val := range row {
	// 		fmt.Printf("%t, ", val)
	// 	}
	// 	fmt.Printf("]\n")
	// }
}
