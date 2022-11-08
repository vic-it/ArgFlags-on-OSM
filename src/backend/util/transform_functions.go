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

func BASICtoGEOJSON() {

}

func BASICtoGRAPH() {

}

func GRAPHtoBASIC() {

}

// this function takes in a path (as string) to a PBF file, reads it, extracts all coastlines and transforms them into the basic format
// this code is mostly taken from https://pkg.go.dev/github.com/qedus/osmpbf#section-readme
func PBFtoBASIC(path string) basic {
	nodes := make(map[int64]node)
	ways := make(map[int64]way)

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
			switch v := v.(type) {
			case *osmpbf.Node:
				// Process Node v.
				nodes[v.ID] = node{lat: v.Lat, lon: v.Lon}
				nc++
			case *osmpbf.Way:
				// Process Way v.
				if v.Tags["natural"] == "coastline" {
					ways[v.ID] = way{nodes: v.NodeIDs}
				}
				wc++
			case *osmpbf.Relation:
				// Process Relation v.
				rc++
			default:
				log.Fatalf("unknown type %T\n", v)
			}
		}
	}

	fmt.Printf("Read: %d Nodes and %d Ways\n", nc, wc)
	return basic{nodes: nodes, ways: ways}
}
