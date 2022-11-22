package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	geojson "github.com/paulmach/go.geojson"
	"github.com/qedus/osmpbf"
)

func BASICtoFMI() {

}

func FMItoBASIC() {

}

func BASICtoGEOJSONFile(basicData basic) {
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
			nodeAsArray = append(nodeAsArray, basicData.nodes[nodex].lon)
			nodeAsArray = append(nodeAsArray, basicData.nodes[nodex].lat)
			// prepare node s.t. garbage collection will clean it up
			basicData.nodes[nodex] = node{}

			polygon = append(polygon, nodeAsArray)
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
	fmt.Printf("%d out of %d nodes were processed", nodeCtr, len(basicData.nodes))
	if err != nil {
		panic(err)
	}
	rawJSON = nil
}

func PointsToGEOJSONFile(points [][]float64) {

	fc := geojson.NewMultiPointFeature(points...)
	fc.SetProperty("x", "y")
	rawJSON, _ := fc.MarshalJSON()
	err := os.WriteFile("../../data/pointgrid.json", rawJSON, 0644)
	if err != nil {
		panic(err)
	}
	rawJSON = nil
}

func BASICtoGRAPH() {

}

func GRAPHtoBASIC() {

}

// this function takes in a path (as string) to a PBF file, reads it, extracts all coastlines and transforms them into the basic format
// this code is mostly taken from https://pkg.go.dev/github.com/qedus/osmpbf#section-readme
func PBFtoBASIC(path string) basic {
	nodes := make(map[int64]node) //(key,value) -> (ID of node, {latitude, longitude})
	ways := make(map[int64]way)   // -> (ID of way, [list of node IDs in way])

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
				// save all nodes
				nodes[v.ID] = node{lat: v.Lat, lon: v.Lon}
				nc++
			case *osmpbf.Way:
				// only save ways with the coastline tag
				if v.Tags["natural"] == "coastline" {
					ways[v.NodeIDs[0]] = way{nodes: v.NodeIDs, lastNodeID: v.NodeIDs[len(v.NodeIDs)-1]}
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
	//merges touching ways
	fmt.Printf("Read: %d Nodes and %d Ways\n", len(nodes), len(ways))
	for {
		if (MergeWays(basic{nodes: nodes, ways: ways}) == 0) {
			break
		}
	}
	//checks if all ways left are polygons (can be deleted)
	ctr := 0
	for _, way := range ways {
		if way.nodes[0] == way.nodes[len(way.nodes)-1] {
			ctr++
		}
	}
	fmt.Printf("Number of ways that are closed polygons: %d\n", ctr)
	return basic{nodes: nodes, ways: ways}
}

func PrintEdgesToGEOJSON(points [][]float64, src []int, dest []int) {
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
