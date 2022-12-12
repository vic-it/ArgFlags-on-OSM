package util

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	geojson "github.com/paulmach/go.geojson"
	"github.com/qedus/osmpbf"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// writes a graph into a file with the following format:
//
//	(ctr = 0)
//
// numOfNodes,numOfEdges
// x (ctr = 1)
// lon,lat		(of all nodes)
// x (ctr = 2)
// src,target,weight  (for each edge)
// x (ctr = 3)
// offset				(for all nodes)
// x (ctr = 4)
// nodeID0ofRow,nodeID1ofRow,... (for all nodes in this NodeMatrixRow - repeat for all rows)
// x (ctr = 5)
// isInWater1,...			(analogous to nodeid matrix)
// x (ctr = 6)
// numberOfNodesIntendedToCreate
// z (end)
func GraphToFile(graph Graph, path string) {
	startTime := time.Now()
	println("WRITING GRAPH TO FILE...")
	nodesToWrite := len(graph.Nodes)
	edgesToWrite := len(graph.Targets)

	f, err := os.Create(path)
	check(err)
	defer f.Close()
	w := bufio.NewWriter(f)

	//WRITE HERE

	w.WriteString(fmt.Sprintf("%d,%d\n", nodesToWrite, edgesToWrite))
	w.WriteString("x\n")
	//WRITE NODES
	for _, node := range graph.Nodes {
		w.WriteString(fmt.Sprintf("%f,%f\n", node[0], node[1]))
	}
	w.WriteString("x\n")
	//WRITE EDGES
	for i := 0; i < len(graph.Targets); i++ {
		w.WriteString(fmt.Sprintf("%d,%d,%d\n", graph.Sources[i], graph.Targets[i], graph.Weights[i]))
	}
	w.WriteString("x\n")
	//WRITE OFFSETS
	for _, offset := range graph.Offsets {
		w.WriteString(fmt.Sprintf("%d\n", offset))
	}
	w.WriteString("x\n")
	//WRITE NODE MATRIX
	for _, row := range graph.NodeMatrix {
		rowString := ""
		for i, nodeID := range row {
			if i < len(row)-1 {
				rowString += fmt.Sprintf("%d,", nodeID)
			} else {
				rowString += fmt.Sprintf("%d\n", nodeID)
			}
		}
		w.WriteString(rowString)
	}
	w.WriteString("x\n")
	//WRITE WATER MATRIX
	for _, row := range graph.NodeInWaterMatrix {
		rowString := ""
		for i, isInWater := range row {
			if i < len(row)-1 {
				rowString += fmt.Sprintf("%t,", isInWater)
			} else {
				rowString += fmt.Sprintf("%t\n", isInWater)
			}
		}
		w.WriteString(rowString)
	}
	w.WriteString("x\n")
	//WRITE INTENDED NUM OF NODES
	w.WriteString(fmt.Sprintf("%d\n", graph.intendedNodeQuantity))
	//END
	w.WriteString("z")
	w.Flush()

	fmt.Printf("Time to write graph to file: %.3fs\n", time.Since(startTime).Seconds())
}

// imports a graph from a .graph file
func FileToGraph(path string) Graph {
	startTime := time.Now()
	graph := Graph{Nodes: [][]float64{}, Sources: []int{}, Targets: []int{}, Weights: []int{}, Offsets: []int{}, NodeMatrix: [][]int{}, NodeInWaterMatrix: [][]bool{}, intendedNodeQuantity: 0}
	println("IMPORTING GRAPH FROM FILE...")
	f, err := os.Open(path)
	check(err)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	ctr := 0
	//READ GRAPH HERE
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\n")
		if line == "x" {
			ctr++
		} else if line == "z" {
			break
		} else {
			switch ctr {
			case 0: //READ NUM OF NODES/EDGES
				list := strings.Split(line, ",")
				no, _ := strconv.Atoi(list[0])
				ed, _ := strconv.Atoi(list[1])
				fmt.Printf("Reading graph with: %d nodes and %d edges\n", no, ed)
			case 1: // lon,lat				(of all nodes)
				list := strings.Split(line, ",")
				lon, _ := strconv.ParseFloat(list[0], 64)
				lat, _ := strconv.ParseFloat(list[1], 64)
				graph.Nodes = append(graph.Nodes, []float64{lon, lat})
			case 2: // src,target,weight 	(for each edge)
				list := strings.Split(line, ",")
				src, _ := strconv.Atoi(list[0])
				trgt, _ := strconv.Atoi(list[1])
				weight, _ := strconv.Atoi(list[2])
				graph.Sources = append(graph.Sources, src)
				graph.Targets = append(graph.Targets, trgt)
				graph.Weights = append(graph.Weights, weight)
			case 3: // offset				(for all nodes)
				list := strings.Split(line, ",")
				offset, _ := strconv.Atoi(list[0])
				graph.Offsets = append(graph.Offsets, offset)
			case 4: // nodeID0ofRow,nodeID1ofRow,..
				list := strings.Split(line, ",")
				row := []int{}
				for _, stringID := range list {
					nodeID, _ := strconv.Atoi(stringID)
					row = append(row, nodeID)
				}
				graph.NodeMatrix = append(graph.NodeMatrix, row)
			case 5: // isInWater1,...
				list := strings.Split(line, ",")
				row := []bool{}
				for _, boolAsString := range list {
					isInWater, _ := strconv.ParseBool(boolAsString)
					row = append(row, isInWater)
				}
				graph.NodeInWaterMatrix = append(graph.NodeInWaterMatrix, row)
			case 6: // numberOfNodesIntendedToCreate
				graph.intendedNodeQuantity, _ = strconv.Atoi(line)
			}
		}
	}
	//END
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Time to read graph from file: %.3fs\n", time.Since(startTime).Seconds())
	return graph
}

// func BASICtoGEOJSONFile(basicData Basic) {
// 	//save basic data in geojson format -> as file (.json)
// 	var polygonList [][][]float64
// 	wayCtr := 0
// 	nodeCtr := 0
// 	for _, wayx := range basicData.ways {
// 		//increment way counter
// 		wayCtr++
// 		//store this way as polygon
// 		var polygon [][]float64
// 		for _, nodex := range wayx.nodes {
// 			//increment node counter
// 			nodeCtr++
// 			var nodeAsArray []float64
// 			nodeAsArray = append(nodeAsArray, basicData.Nodes[nodex].lon)
// 			nodeAsArray = append(nodeAsArray, basicData.Nodes[nodex].lat)
// 			// prepare node s.t. garbage collection will clean it up
// 			polygon = append(polygon, nodeAsArray)
// 			//basicData.nodes[nodex] = node{}

// 			// force garbage collection -> else memory overruns
// 			if nodeCtr%10000 == 0 {
// 				runtime.GC()
// 			}

// 		}
// 		polygonList = append(polygonList, polygon)
// 		//prepare way s.t. garbage collection will clean it
// 		wayx = way{}
// 		// print geojson progress aswell as force garbage collection
// 		if wayCtr%10000 == 0 {
// 			PrintProgress(wayCtr, len(basicData.ways), "ways")
// 			runtime.GC()
// 		}
// 	}
// 	g := geojson.NewMultiPolygonGeometry(polygonList)
// 	rawJSON, _ := g.MarshalJSON()
// 	err := os.WriteFile("../../data/geojson.json", rawJSON, 0644)
// 	println("geojson file written to: '../../data/geojson.json'")
// 	fmt.Printf("%d out of %d nodes were processed\n", nodeCtr, len(basicData.Nodes))
// 	if err != nil {
// 		panic(err)
// 	}
// 	rawJSON = nil
// }

func PrintPointsToGEOJSON(graph Graph) {
	println("WRITING NODES TO GEOJSON")

	startTime := time.Now()

	points := graph.Nodes
	fc := geojson.NewMultiPointFeature(points...)
	fc.SetProperty("x", "y")
	rawJSON, _ := fc.MarshalJSON()
	err := os.WriteFile("../../data/pointgrid.json", rawJSON, 0644)
	if err != nil {
		panic(err)
	}
	rawJSON = nil
	fmt.Printf("Time to read in coast lines: %.3fs\n", time.Since(startTime).Seconds())
}

func PrintEdgesToGEOJSON(graph Graph) {
	println("WRITING EDGES TO GEOJSON")
	points := graph.Nodes
	src := graph.Sources
	dest := graph.Targets
	var lineList [][][]float64
	for i := 0; i < len(src); i++ {
		line := [][]float64{points[src[i]], points[dest[i]]}
		lineList = append(lineList, line)
	}
	fc := geojson.NewMultiLineStringFeature(lineList...)
	fc.SetProperty("x", "y")
	rawJSON, _ := fc.MarshalJSON()
	err := os.WriteFile("../../data/edgegrid.json", rawJSON, 0644)
	if err != nil {
		panic(err)
	}
	rawJSON = nil
}

// Reads in the coastline files and preprocesses them into a Coastline format (see datatypes.go)
func GetCoastline(path string) Coastline {
	println("READING IN COASTLINES... (can take some time)")
	startTime := time.Now()
	nodes := make(map[int64][]float64) //(ID,[lon, lat])
	//isUsefulNode := make(map[int64]bool)
	var edges [][]int64 // -> (ID of EDGE, [ID node 1, ID node 2])
	//all sorted from min -> max
	var sortedLonList []EdgeCoordinate
	//maximum width of an edge
	var maxEdgeWidth float64
	var maxEdgeWidths []float64
	//higher -> possibly better performance -> diminishing returns at some point?
	latGranularity := 360

	var placeholder [][]EdgeCoordinate
	for i := 0; i < latGranularity; i++ {
		maxEdgeWidths = append(maxEdgeWidths, 0)
		placeholder = append(placeholder, []EdgeCoordinate{})
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)

	// use more memory from the start, it is faster
	d.SetBufferSize(osmpbf.MaxBlobSize)

	// start decoding with several goroutines, it is faster
	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		log.Fatal(err)
	}

	garbageCollectorCounter := 0
	var nc, wc, rc uint64
	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			// Here you can change what happens to nodes/ways when read
			switch v := v.(type) {
			case *osmpbf.Node:
				// save all nodes as [lon, lat]
				nodes[v.ID] = []float64{v.Lon, v.Lat}
				nc++
				v.Tags = nil
				v = nil
				garbageCollectorCounter++
				if garbageCollectorCounter%5000000 == 0 {
					runtime.GC()
				}
			case *osmpbf.Way:
				// only save ways with the coastline tag
				if v.Tags["natural"] == "coastline" && v.Tags["coastline"] != "bogus" {
					//add all edges
					for i := 0; i < len(v.NodeIDs)-1; i++ {
						edges = append(edges, []int64{v.NodeIDs[i], v.NodeIDs[i+1]})
					}
					//re-enable for deleting useless nodes -> takes 30s extra time for global
					// for _, id := range v.NodeIDs {
					// 	isUsefulNode[id] = true
					// }
					wc++
					v.NodeIDs = []int64{}
					v.Tags = nil
					v = nil
					garbageCollectorCounter++
					if garbageCollectorCounter%5000000 == 0 {
						runtime.GC()
					}
				}
			case *osmpbf.Relation:
				// dont save any relations for now
				rc++
			default:
				log.Fatalf("unknown type %T\n", v)
			}
		}
	}

	fmt.Printf("Time to read in coast lines: %.3fs\n", time.Since(startTime).Seconds())
	f.Close()
	runtime.GC()
	startTime = time.Now()
	// fill [to be sorted] lists
	for id, edge := range edges {
		maxEdgeWidth = math.Max(maxEdgeWidth, CalcLonDiff(nodes[edge[0]][0], nodes[edge[1]][0]))
		sortedLonList = append(sortedLonList, EdgeCoordinate{edgeID: id, coordinate: nodes[edge[0]][0]})
		sortedLonList = append(sortedLonList, EdgeCoordinate{edgeID: id, coordinate: nodes[edge[1]][0]})
	}

	fmt.Printf("Maximum lat diff: %.6f\n", maxEdgeWidth)

	// sort lists by coordinate
	//functions for sorting algorithm
	println("Sorting longitude list")
	sort.Sort(ByCoordinate(sortedLonList))

	fmt.Printf("Time to sort longitude list: %.3fs\n", time.Since(startTime).Seconds())
	fmt.Printf("Read: %d Nodes and %d edges\n", len(nodes), len(edges))

	coastline := Coastline{Nodes: nodes, Edges: edges, SortedLonEdgeList: sortedLonList, MaxLonDiffs: maxEdgeWidths, maxLonDiff: maxEdgeWidth}
	return coastline
}

func checkMaxLon(maxLonDiffList []float64, lonDiff float64, lat1 float64, lat2 float64) {
	n := len(maxLonDiffList)

	index1 := GetLonDiffIndex(n, lat1)

	index2 := GetLonDiffIndex(n, lat2)

	maxLonDiffList[index1] = math.Max(maxLonDiffList[index1], lonDiff)
	maxLonDiffList[index2] = math.Max(maxLonDiffList[index2], lonDiff)
}

func GetLonDiffIndex(n int, lat float64) int {
	return int(math.Round((lat + 90) * (float64(n) / 180)))
}
