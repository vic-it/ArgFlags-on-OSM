package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/qedus/osmpbf"
)

func BASICtoFMI() {

}

func FMItoBASIC() {

}

func BASICtoGEOJSON(basicData basic) {
	//save basic data in geojson format -> as file (.json)
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
