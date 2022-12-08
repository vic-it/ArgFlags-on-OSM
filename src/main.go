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

func main() {
	println("\nStart!")
	randomTestFunction(20000)
	startServer()
}

func startServer() {
	println("starting http server...")
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/getpoint", getPointHandler())
	http.HandleFunc("/getroute", getRouteHandler())
	http.HandleFunc("/directory", getDirectoryHandler())
	http.HandleFunc("/preprocess", getPreprocessingHandler())
	http.ListenAndServe(":8080", nil)
}

func getPointHandler() http.HandlerFunc {
	pointHandler := func(writer http.ResponseWriter, request *http.Request) {
		urlQuery := request.URL.Query()
		inputLon, _ := strconv.ParseFloat(urlQuery["lon"][0], 64)
		inputLat, _ := strconv.ParseFloat(urlQuery["lat"][0], 64)

		fmt.Printf("click lon: %f\n", inputLon)
		fmt.Printf("click lat: %f\n", inputLat)
		nodeID := util.GetClosestGridNode(inputLon, inputLat, graph)

		outputString := fmt.Sprintf("%fx%fx%d", graph.Nodes[nodeID][0], graph.Nodes[nodeID][1], nodeID)
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(outputString))
	}

	return pointHandler
}

func getRouteHandler() http.HandlerFunc {
	routingHandler := func(writer http.ResponseWriter, request *http.Request) {
		urlQuery := request.URL.Query()
		src, _ := strconv.ParseInt(urlQuery["src"][0], 0, 64)
		dest, _ := strconv.ParseInt(urlQuery["dest"][0], 0, 64)

		fmt.Printf("source node ID: %d\n", src)
		fmt.Printf("destination node ID: %d\n--\n", dest)

		writer.WriteHeader(http.StatusOK)
		//calc route below
		writer.Write([]byte(getRouteString(src, dest)))
	}
	return routingHandler
}

// should return a string in the following format:
// if no route found -> -1y0
// if route found -> distanceylon1zlat1xlon2zlat2x...
func getRouteString(src int64, dest int64) string {
	distance, path := util.CalculateDijkstra(graph, int(src), int(dest))
	output := fmt.Sprintf("%dy", distance)
	for i, nodeID := range path {
		if i == len(path)-1 {
			output = fmt.Sprintf("%s%fz%f", output, graph.Nodes[nodeID][0], graph.Nodes[nodeID][1])
		} else {
			output = fmt.Sprintf("%s%fz%fx", output, graph.Nodes[nodeID][0], graph.Nodes[nodeID][1])
		}
	}
	return output
}

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

func analyzeDirectory(basePath string) {
	println("Reading directory...")

}

// creates a graph with the coastlines from the path and roughly the number of nodes
func createGraph(pathToCoastlinesPBF string, numberOfNodes int) {
	graph = util.Graph{}
	runtime.GC()
	graph = util.GenerateGraph(numberOfNodes, util.GetCoastline(pathToCoastlinesPBF))
}

func readGraphFromFMI(path string) {

}

func printGraphToFMI(path string) {

}

func randomTestFunction(numOfNodes int) {

	// "../../data/antarctica.osm.pbf"
	// "../../data/central-america.osm.pbf"
	// "../../data/global.sec" THIS IS THE BIG ONE FROM ILIAS (renamed, takes up ~11GB of RAM!)
	path := "../../data/antarctica.osm.pbf"
	coastline := util.GetCoastline(path)
	//generates grid around globes
	graph = util.GenerateGraph(numOfNodes, coastline)
	util.PrintPointsToGEOJSON(graph.Nodes)
	util.PrintEdgesToGEOJSON(graph)
}
