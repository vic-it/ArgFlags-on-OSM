package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/qedus/osmpbf"
)

func ReadPBF(path string) basic {
	rawCoastlines := extractCoastlinesFromPBF(path)
	return rawCoastlines
}

func extractCoastlinesFromPBF(path string) basic {
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
					fmt.Println("add way with tag natural = " + v.Tags["natural"])
					fmt.Printf("and number of nodes in way: %d\n\n", len(v.NodeIDs))
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

	fmt.Printf("Nodes: %d, Ways: %d, Relations: %d\n", nc, wc, rc)
	return basic{nodes: nodes, ways: ways}
}
