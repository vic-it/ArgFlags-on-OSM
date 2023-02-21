package backend

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

	fmt.Printf("Time to write graph to file: %.1fs\n", time.Since(startTime).Seconds())
}

// imports a graph from a .graph file
func FileToGraph(path string) Graph {
	startTime := time.Now()
	graph := Graph{Nodes: [][]float64{}, Sources: []int{}, Targets: []int{}, Weights: []int{}, Offsets: []int{}, NodeMatrix: [][]int{}, NodeInWaterMatrix: [][]bool{}, intendedNodeQuantity: 0}
	countOfNodesInWater := 0
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
				fmt.Printf("Importing graph with: %d nodes and %d edges...\n", no, ed)
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
					if isInWater {
						countOfNodesInWater++
					}
					row = append(row, isInWater)
				}
				graph.NodeInWaterMatrix = append(graph.NodeInWaterMatrix, row)
			case 6: // numberOfNodesIntendedToCreate
				graph.intendedNodeQuantity, _ = strconv.Atoi(line)
			}
		}
	}
	graph.Offsets = append(graph.Offsets, len(graph.Targets))
	graph.countOfWaterNodes = countOfNodesInWater
	//END
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Time to read graph from file: %.1fs\n", time.Since(startTime).Seconds())
	return graph
}

// writes a arcflag data a file with the following format:
//
//	(ctr = 0)
//
// numOfPartitions
// x (ctr = 1)
// nodeID0ofRow,nodeID1ofRow,... (for all nodes in this PartitionNodeMatrix - repeat for all rows)
// x (ctr = 2)
// Arcflags1,...			(analogous to nodeid matrix)
// z (end)
func ArcFlagsToFile(arcData ArcData, path string) {
	startTime := time.Now()
	println("WRITING arc flags TO FILE...")

	f, err := os.Create(path)
	check(err)
	defer f.Close()
	w := bufio.NewWriter(f)

	//WRITE HERE

	w.WriteString(fmt.Sprintf("%d\n", arcData.NumberOfPartitions))
	w.WriteString("x\n")

	//WRITE NODE PARTITION MATRIX
	for _, row := range arcData.NodePartitionMatrix {
		rowString := ""
		for i, partitionID := range row {
			if i < len(row)-1 {
				rowString += fmt.Sprintf("%d,", partitionID)
			} else {
				rowString += fmt.Sprintf("%d\n", partitionID)
			}
		}
		w.WriteString(rowString)
	}
	w.WriteString("x\n")
	//WRITE ARC FLAGS
	for _, row := range arcData.ArcFlags {
		rowString := ""
		for i, flag := range row {
			if i < len(row)-1 {
				rowString += fmt.Sprintf("%t,", flag)
			} else {
				rowString += fmt.Sprintf("%t\n", flag)
			}
		}
		w.WriteString(rowString)
	}
	//END
	w.WriteString("z")
	w.Flush()

	fmt.Printf("Time to write arc flags to file: %.1fs\n", time.Since(startTime).Seconds())
}

// imports a graph from a .graph file
func FileToArcFlags(path string) ArcData {
	startTime := time.Now()
	arcData := ArcData{}
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
				arcData.NumberOfPartitions, _ = strconv.Atoi(line)
				fmt.Printf("Importing arc flag data with: %d partitions...\n", arcData.NumberOfPartitions)
			case 1: // partitionOfNode1, partitionOfNode2,...
				list := strings.Split(line, ",")
				row := []int{}
				for _, stringID := range list {
					partitionID, _ := strconv.Atoi(stringID)
					row = append(row, partitionID)
				}
				arcData.NodePartitionMatrix = append(arcData.NodePartitionMatrix, row)
			case 2: // arcflags
				list := strings.Split(line, ",")
				row := []bool{}
				for _, boolAsString := range list {
					flag, _ := strconv.ParseBool(boolAsString)

					row = append(row, flag)
				}
				arcData.ArcFlags = append(arcData.ArcFlags, row)
			}
		}
	}
	//END
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Time to read arc flag data from file: %.1fs\n", time.Since(startTime).Seconds())
	return arcData
}

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

func PrintPointsToGEOJSON2(graph Graph, nodePartitionMatrix [][]int) {
	println("WRITING NODES TO GEOJSON")

	startTime := time.Now()
	points := [][]float64{}
	for rowID, row := range nodePartitionMatrix {
		for colID, partID := range row {
			if partID%5 == 0 {
				points = append(points, graph.Nodes[graph.NodeMatrix[rowID][colID]])
			}
		}
	}
	fc := geojson.NewMultiPointFeature(points...)
	fc.SetProperty("x", "y")
	rawJSON, _ := fc.MarshalJSON()
	err := os.WriteFile("../../data/pointgrid.json", rawJSON, 0644)
	if err != nil {
		panic(err)
	}
	rawJSON = nil
	fmt.Printf("Time to read in coast lines: %.1fs\n", time.Since(startTime).Seconds())
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

	for i := 0; i < latGranularity; i++ {
		maxEdgeWidths = append(maxEdgeWidths, 0)
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

	fmt.Printf("Time to read in coast lines: %.1fs\n", time.Since(startTime).Seconds())
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
