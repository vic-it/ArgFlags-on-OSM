package util

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"sort"

	geojson "github.com/paulmach/go.geojson"
	"github.com/qedus/osmpbf"
)

func BASICtoFMI() {

}

func FMItoBASIC() {

}

func BASICtoGEOJSONFile(basicData Basic) {
	//save basic data in geojson format -> as file (.json)
	var polygonList [][][]float64
	wayCtr := 0
	nodeCtr := 0
	for _, wayx := range basicData.ways {
		//increment way counter
		wayCtr++
		//store this way as polygon
		var polygon [][]float64
		for _, nodex := range wayx.nodes {
			//increment node counter
			nodeCtr++
			var nodeAsArray []float64
			nodeAsArray = append(nodeAsArray, basicData.Nodes[nodex].lon)
			nodeAsArray = append(nodeAsArray, basicData.Nodes[nodex].lat)
			// prepare node s.t. garbage collection will clean it up
			polygon = append(polygon, nodeAsArray)
			//basicData.nodes[nodex] = node{}

			// force garbage collection -> else memory overruns
			if nodeCtr%10000 == 0 {
				runtime.GC()
			}

		}
		polygonList = append(polygonList, polygon)
		//prepare way s.t. garbage collection will clean it
		wayx = way{}
		// print geojson progress aswell as force garbage collection
		if wayCtr%10000 == 0 {
			PrintProgress(wayCtr, len(basicData.ways), "ways")
			runtime.GC()
		}
	}
	g := geojson.NewMultiPolygonGeometry(polygonList)
	rawJSON, _ := g.MarshalJSON()
	err := os.WriteFile("../../data/geojson.json", rawJSON, 0644)
	println("geojson file written to: '../../data/geojson.json'")
	fmt.Printf("%d out of %d nodes were processed\n", nodeCtr, len(basicData.Nodes))
	if err != nil {
		panic(err)
	}
	rawJSON = nil
}

func PrintPointsToGEOJSON(points [][]float64) {

	fc := geojson.NewMultiPointFeature(points...)
	fc.SetProperty("x", "y")
	rawJSON, _ := fc.MarshalJSON()
	err := os.WriteFile("../../data/pointgrid.json", rawJSON, 0644)
	if err != nil {
		panic(err)
	}
	rawJSON = nil
}

func PrintEdgesToGEOJSON(graph Graph) {
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

// gets the pbf file from the path and outputs a list of all edges and 3 lists of edge id's sorted by e.g. max lat
func GetCoastline(path string) Coastline {
	nodes := make(map[int64][]float64) //(ID,[lon, lat])
	//isUsefulNode := make(map[int64]bool)
	var edges [][]int64 // -> (ID of EDGE, [ID node 1, ID node 2])
	//all sorted from min -> max
	var sortedLonList []EdgeCoordinate
	//maximum width of an edge
	var maxEdgeWidth float64

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
				}
			case *osmpbf.Relation:
				// dont save any relations for now
				rc++
			default:
				log.Fatalf("unknown type %T\n", v)
			}
		}
	}

	// fill [to be sorted] lists
	for id, edge := range edges {
		maxEdgeWidth = math.Max(maxEdgeWidth, CalcLonDiff(nodes[edge[0]][0], nodes[edge[1]][0]))
		sortedLonList = append(sortedLonList, EdgeCoordinate{edgeID: id, coordinate: nodes[edge[0]][0]})
		sortedLonList = append(sortedLonList, EdgeCoordinate{edgeID: id, coordinate: nodes[edge[1]][0]})
	}
	fmt.Printf("Maximum lat diff: %.6f\n", maxEdgeWidth)

	// sort lists by coordinate
	//functions for sorting algorithm
	sort.Sort(ByCoordinate(sortedLonList))

	fmt.Printf("Read: %d Nodes and %d edges\n", len(nodes), len(edges))

	coastline := Coastline{Nodes: nodes, Edges: edges, SortedLonEdgeList: sortedLonList, MaxLonDiff: maxEdgeWidth}
	return coastline
}
