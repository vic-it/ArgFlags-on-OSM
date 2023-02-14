package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"

	"github.com/vic-it/OSM-FMI/src/backend"
)

// the graph!
var graph backend.Graph
var arcData backend.ArcData
var nodePartitionList []int
var runningActiveTests bool

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
	http.HandleFunc("/testalgorithms", getTestHandler())
	http.HandleFunc("/querytestprogress", getProgressHandler())
	http.HandleFunc("/abort", getAbortHandler())
	println("GUI available on: localhost:8080")
	println("READY!\n-----------")
	http.ListenAndServe(":8080", nil)

}

func getAbortHandler() http.HandlerFunc {
	abortHandler := func(writer http.ResponseWriter, request *http.Request) {
		if runningActiveTests {
			runningActiveTests = false
			backend.AbortTests()
		}
	}
	return abortHandler
}
func getTestHandler() http.HandlerFunc {
	testHandler := func(writer http.ResponseWriter, request *http.Request) {
		if !runningActiveTests {
			runningActiveTests = true
			urlQuery := request.URL.Query()
			numOfTests, _ := strconv.ParseInt(urlQuery["num"][0], 0, 64)
			backend.TestAlgorithms(graph, arcData, int(numOfTests), nodePartitionList)
			runningActiveTests = false
		}
	}
	return testHandler
}

func getProgressHandler() http.HandlerFunc {
	progressHandler := func(writer http.ResponseWriter, request *http.Request) {

		outputString := backend.GetTestStatusAsString()
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(outputString))
	}

	return progressHandler
}

// checks for nearest valid neighbor node in water and returns this to the web interface
func getPointHandler() http.HandlerFunc {
	pointHandler := func(writer http.ResponseWriter, request *http.Request) {
		urlQuery := request.URL.Query()
		inputLon, _ := strconv.ParseFloat(urlQuery["lon"][0], 64)
		inputLat, _ := strconv.ParseFloat(urlQuery["lat"][0], 64)

		// fmt.Printf("click lon: %f\n", inputLon)
		// fmt.Printf("click lat: %f\n", inputLat)
		nodeID := backend.GetClosestGridNode(inputLon, inputLat, graph)

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
		distance, path, initTime, searchTime, nodesPopped = backend.CalculateDijkstra(graph, int(src), int(dest))
	} else {
		distance, path, initTime, searchTime, nodesPopped = backend.CalculateArcFlagDijkstra(graph, int(src), int(dest), arcData, nodePartitionList)
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

// creates a graph with the coastlines from the path and roughly the number of nodes
func createGraph(pathToCoastlinesPBF string, numberOfNodes int) {
	graph = backend.Graph{}
	runtime.GC()
	graph = backend.GenerateGraph(numberOfNodes, backend.GetCoastline(pathToCoastlinesPBF))
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
	// createGraph(antarctica, 100000)
	// backend.GraphToFile(graph, graphPath)
	//-----------------------------------------------------

	// IMPORT GRAPH BY UNCOMMENTING BELOW:
	//-----------------------------------------------------
	graph = backend.FileToGraph(graphPath)
	//-----------------------------------------------------

	// PRINT TO GEOJSON BY UNCOMMENTING BELOW:
	//-----------------------------------------------------
	// backend.PrintPointsToGEOJSON(graph)
	// backend.PrintEdgesToGEOJSON(graph)
	//-----------------------------------------------------

	//arcFlagStuff
	// GENERATE NEW ARCFLAGS BY UNCOMMENTING BELOW
	// 7 - 3 generates roughly square partitions (64 of them)
	// arcData = backend.PreprocessArcFlags(graph, 7, 3)
	// backend.ArcFlagsToFile(arcData, arcFlagPath)

	// IMPORT ARCFLAGS BY UNCOMMENTING BELOW:
	arcData = backend.FileToArcFlags(arcFlagPath)
	// this speeds up arc flag since it doesnt have to calculate row/col of nodepartitionmatrix anymore
	preparePartitionList()
}

func testStuff() {
	//backend.PrintPointsToGEOJSON2(graph, arcData.NodePartitionMatrix)
}
