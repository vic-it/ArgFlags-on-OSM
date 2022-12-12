package util

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"sort"

	"github.com/qedus/osmpbf"
)

// THIS FILE IS NOT IN USE YET -> PREPARATION FOR MORE EFFICIENT GRAPH GENERATION/POINT IN WATER TEST
func PBFtoBASIC(path string) []way {
	ways := make(map[int64]way) // -> (ID of way, [list of node IDs in way])
	nodesPlaceHolder := make(map[int64]node)

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
				// save all nodes
				nodesPlaceHolder[v.ID] = node{lat: v.Lat, lon: v.Lon}
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
	fmt.Printf("Read: %d Nodes and %d Ways\n", len(nodesPlaceHolder), len(ways))
	for {
		if (MergeWays(Basic{Nodes: nodesPlaceHolder, ways: ways}) == 0) {
			break
		}
	}
	generatePolygons(Basic{Nodes: nodesPlaceHolder, ways: ways})
	return []way{}
}

func generatePolygons(input Basic) []polygon {
	var polygons []polygon

	for _, way := range input.ways {
		polygonToAdd := polygon{}
		maxEdgeWidth := 0.
		//add edges to polygon
		polygonToAdd.sortedEdges = append(polygonToAdd.sortedEdges, EdgeCoordinate{edgeID: 0, coordinate: input.Nodes[way.nodes[0]].lon})
		polygonToAdd.sortedEdges = append(polygonToAdd.sortedEdges, EdgeCoordinate{edgeID: 0, coordinate: input.Nodes[way.nodes[len(way.nodes)-1]].lon})
		maxEdgeWidth = math.Max(maxEdgeWidth, CalcLonDiff(input.Nodes[way.nodes[0]].lon, input.Nodes[way.nodes[len(way.nodes)-1]].lon))
		for i := 0; i < len(way.nodes)-1; i++ {
			lon1 := input.Nodes[way.nodes[i]].lon
			lon2 := input.Nodes[way.nodes[i+1]].lon
			polygonToAdd.sortedEdges = append(polygonToAdd.sortedEdges, EdgeCoordinate{edgeID: i + 1, coordinate: lon1})
			polygonToAdd.sortedEdges = append(polygonToAdd.sortedEdges, EdgeCoordinate{edgeID: i + 1, coordinate: lon2})
			maxEdgeWidth = math.Max(maxEdgeWidth, CalcLonDiff(input.Nodes[way.nodes[i]].lon, input.Nodes[way.nodes[i+1]].lon))
		}
		//CALCULATE POLYGON BOUNDARIES!
		sort.Sort(ByCoordinate(polygonToAdd.sortedEdges))
		polygons = append(polygons, polygonToAdd)

	}

	return polygons
}
