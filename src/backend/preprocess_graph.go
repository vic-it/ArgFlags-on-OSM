package backend

import (
	"fmt"
	"math"
	"runtime"
	"sort"
	"time"
)

// math from here https://www.cmu.edu/biolphys/deserno/pdf/sphere_equi.pdf
// generates a graph with equidistant nodes over a sphere with a coastline as basis
func GenerateGraph(numberOfNodes int, coastline Coastline) Graph {
	startTime := time.Now()
	startTimeTotal := time.Now()
	pointInWaterCount := 0
	println("generating points...")
	// simple list of all points
	var points [][]float64
	//list but as 2d-ish grid for edge creation -> with index of points in "points[]"
	var pointMatrix [][]int
	// a bool matrix corresponding to the point matrix which just saves which nodes are in water
	var isPointInWaterMatrix [][]bool
	pi := math.Pi
	//count of nodes
	count := 0
	a := 4.0 * pi / float64(numberOfNodes)
	d := math.Sqrt(a)
	Mv := math.Round(pi / d)
	dv := pi / Mv
	dp := a / dv
	// iterate over lats
	for m := 0.0; m < Mv; m++ {
		var latList []int
		var waterLatList []bool
		v := pi * (m + 0.5) / Mv
		Mp := math.Round(2.0 * pi * math.Sin(v) / dp)
		// iterate over longs
		for n := 0.0; n < Mp; n++ {
			//generate point?
			p := 2.0 * pi * n / Mp
			var point []float64
			// p -> lon
			// v -> lat
			lon, lat := RadToDeg(p, v)
			point = append(point, lon)
			point = append(point, lat)
			points = append(points, point)
			latList = append(latList, len(points)-1)
			//CHECKS FOR EACH POINT IF IT IS IN WATER
			z := IsPointInWater(point, coastline)
			if z {
				pointInWaterCount++
			}
			waterLatList = append(waterLatList, z)
			count++
			if count%5000 == 0 {
				PrintProgress(count, numberOfNodes, "nodes checked.", startTime)
				startTime = time.Now()
				runtime.GC()
			}
		}
		isPointInWaterMatrix = append(isPointInWaterMatrix, waterLatList)
		pointMatrix = append(pointMatrix, latList)
	}
	fmt.Printf("Time to generate grid points: %.3fs\n", time.Since(startTimeTotal).Seconds())
	fmt.Printf("%d points created\n", count)
	//some garbage collection stuff because my RAM is not that good :'(
	coastline.Edges = [][]int64{}
	coastline.Nodes = nil
	coastline.SortedLonEdgeList = []EdgeCoordinate{}
	coastline = Coastline{}
	runtime.GC()
	// takes the generated nodes and adds edges too the graph
	return GenerateEdges(points, pointMatrix, isPointInWaterMatrix, numberOfNodes, pointInWaterCount)
}

// takes nodes and the above mentioned point matrices as input and creates valid edges for points that are indeed in water
func GenerateEdges(points [][]float64, indexMatrix [][]int, pointsInWaterMatrix [][]bool, numOfNodes int, pointInWaterCount int) Graph {
	println("creating edges from points...")
	startTime := time.Now()
	var fwEdgeSources []int
	var fwEdgeDest []int
	//initalize bwedges with two empty lists
	var bwEdges [][]int
	var distanceList []int
	for y := 0; y < len(indexMatrix); y++ {
		latList := indexMatrix[y]
		for x := 0; x < len(latList); x++ {
			//add left edge - does not need seperate back edge creation as they are created equally for all
			if pointsInWaterMatrix[y][x] {

				if x == 0 {
					if pointsInWaterMatrix[y][len(latList)-1] {
						fwEdgeSources = append(fwEdgeSources, indexMatrix[y][x])
						fwEdgeDest = append(fwEdgeDest, indexMatrix[y][len(latList)-1])
					}
				} else {
					if pointsInWaterMatrix[y][x-1] {

						fwEdgeSources = append(fwEdgeSources, indexMatrix[y][x])
						fwEdgeDest = append(fwEdgeDest, indexMatrix[y][x-1])
					}
				}
				//add right edge  - does not need seperate back edge creation as they are created equally for all
				if x == len(latList)-1 {
					if pointsInWaterMatrix[y][0] {

						fwEdgeSources = append(fwEdgeSources, indexMatrix[y][x])
						fwEdgeDest = append(fwEdgeDest, indexMatrix[y][0])
					}
				} else {
					if pointsInWaterMatrix[y][x+1] {

						fwEdgeSources = append(fwEdgeSources, indexMatrix[y][x])
						fwEdgeDest = append(fwEdgeDest, indexMatrix[y][x+1])
					}
				}

				//add above edge as well as a backwards edge "below edge" for the node above to this one
				if y != len(indexMatrix)-1 {
					position := int(math.Round(float64(x * len(indexMatrix[y+1]) / len(indexMatrix[y]))))
					if pointsInWaterMatrix[y+1][position] {
						fwEdgeSources = append(fwEdgeSources, indexMatrix[y][x])
						fwEdgeDest = append(fwEdgeDest, indexMatrix[y+1][position])
						bwEdges = append(bwEdges, []int{indexMatrix[y+1][position], indexMatrix[y][x]})
					}
				}
			}

		}
	}
	// since the backedges are not necessarily in correct order for the offset list, we must first sort it and merge it to one (two) big list(s)
	mergedEdgeSources, mergedEdgeDest, offsetList := mergeEdges(fwEdgeSources, fwEdgeDest, bwEdges, len(points))

	fmt.Printf("Time to generate edges of grid: %.3fs\n", time.Since(startTime).Seconds())
	startTime = time.Now()
	//calculates the distance between two points (an edge)
	distanceList = CalcEdgeDistances(points, mergedEdgeSources, mergedEdgeDest)

	fmt.Printf("Time to calculate edge distances: %.3fs\n", time.Since(startTime).Seconds())
	return Graph{Nodes: points, Sources: mergedEdgeSources, Targets: mergedEdgeDest, Weights: distanceList, Offsets: offsetList, NodeMatrix: indexMatrix, NodeInWaterMatrix: pointsInWaterMatrix, intendedNodeQuantity: numOfNodes, countOfWaterNodes: pointInWaterCount}
}

// takes forward and backward edge lists and merges them into one "forward" edge list as well as an offset list
func mergeEdges(fwsrc []int, fwdest []int, bwEdges [][]int, n int) ([]int, []int, []int) {
	println("before sort")
	// sort backward edge list ("below edges") by their source node ID, s.t. the offset can easily be calculated
	sort.Sort(bySourceID(bwEdges))
	var offsetList = make([]int, n)
	var mergedSrc []int
	var mergedDest []int
	fwdctr := 0
	bwctr := 0
	// i = current source node IDX of which we want to add all edges
	for i := 0; i < n; i++ {
		//increment offset counter
		offsetList[i] = len(mergedSrc)
		//add normal forward edges in order
		for {
			//break if list exhausted or next element would be bigger than current source node idx
			if fwdctr == len(fwsrc) || fwsrc[fwdctr] > i {
				break
			}
			mergedSrc = append(mergedSrc, fwsrc[fwdctr])
			mergedDest = append(mergedDest, fwdest[fwdctr])
			fwdctr++
		}
		//add backward edges in order
		for {
			//break if list exhausted or next element would be bigger than current source node idx
			if bwctr == len(bwEdges) || bwEdges[bwctr][0] > i {
				break
			}
			mergedSrc = append(mergedSrc, bwEdges[bwctr][0])
			mergedDest = append(mergedDest, bwEdges[bwctr][1])
			bwctr++
		}
	}

	return mergedSrc, mergedDest, offsetList
}

// self explanatory... really...
func CalcEdgeDistances(points [][]float64, src []int, dest []int) []int {
	var distances []int
	for i, d := range src {
		distance := int(dist(points[d], points[dest[i]]))
		distances = append(distances, distance)
	}
	return distances
}

// sorting stuff for backward edges
type bySourceID [][]int

func (edges bySourceID) Len() int {
	return len(edges)
}
func (edges bySourceID) Swap(i, j int) {
	println("swap")
	edges[i][0], edges[j][0] = edges[j][0], edges[i][0]
	edges[i][1], edges[j][1] = edges[j][1], edges[i][1]
}
func (edges bySourceID) Less(i, j int) bool {
	return edges[i][0] < edges[j][0]
}
