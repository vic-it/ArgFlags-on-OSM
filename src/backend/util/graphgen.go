package util

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

// math from here https://www.cmu.edu/biolphys/deserno/pdf/sphere_equi.pdf
func GenerateGraph(numberOfNodes int, coastline Coastline) Graph {
	startTime := time.Now()
	startTimeTotal := time.Now()
	println("generating points...")
	// simple list of all points
	var points [][]float64
	//list but as 2d-ish grid for edge creation -> with index of points in "points[]"
	var pointMatrix [][]int
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
			//bottle neck here
			z := IsPointInWater(point, coastline)
			waterLatList = append(waterLatList, z)
			count++
			if count%500 == 0 {
				PrintProgress(count, numberOfNodes, "nodes checked.", startTime)
				startTime = time.Now()
			}
		}
		isPointInWaterMatrix = append(isPointInWaterMatrix, waterLatList)
		pointMatrix = append(pointMatrix, latList)
	}
	fmt.Printf("Time to generate grid points: %.3fs\n", time.Since(startTimeTotal).Seconds())
	fmt.Printf("%d points created\n", count)
	coastline.Edges = [][]int64{}
	coastline.Nodes = nil
	coastline.SortedLonEdgeList = []EdgeCoordinate{}
	coastline = Coastline{}
	runtime.GC()
	return GenerateEdges(points, pointMatrix, isPointInWaterMatrix, numberOfNodes)
}

func GenerateEdges(points [][]float64, indexMatrix [][]int, pointsInWaterMatrix [][]bool, numOfNodes int) Graph {
	println("creating edges from points...")
	startTime := time.Now()
	var edgeSource []int
	var edgeDest []int
	var offsetList = make([]int, len(points))
	var distanceList []int
	var totalOffset = 0
	offsetList[0] = 0
	for y := 0; y < len(indexMatrix); y++ {
		latList := indexMatrix[y]
		for x := 0; x < len(latList); x++ {
			//add left edge
			offsetList[indexMatrix[y][x]] = totalOffset
			if pointsInWaterMatrix[y][x] {

				if x == 0 {
					if pointsInWaterMatrix[y][len(latList)-1] {
						totalOffset++
						edgeSource = append(edgeSource, indexMatrix[y][x])
						edgeDest = append(edgeDest, indexMatrix[y][len(latList)-1])
					}
				} else {
					if pointsInWaterMatrix[y][x-1] {

						totalOffset++
						edgeSource = append(edgeSource, indexMatrix[y][x])
						edgeDest = append(edgeDest, indexMatrix[y][x-1])
					}
				}
				//add right edge
				if x == len(latList)-1 {
					if pointsInWaterMatrix[y][0] {

						totalOffset++
						edgeSource = append(edgeSource, indexMatrix[y][x])
						edgeDest = append(edgeDest, indexMatrix[y][0])
					}
				} else {
					if pointsInWaterMatrix[y][x+1] {

						totalOffset++
						edgeSource = append(edgeSource, indexMatrix[y][x])
						edgeDest = append(edgeDest, indexMatrix[y][x+1])
					}
				}
				//add below edge
				if y != 0 {
					position := int(math.Round(float64(x * len(indexMatrix[y-1]) / len(indexMatrix[y]))))
					if pointsInWaterMatrix[y-1][position] {

						totalOffset++
						edgeSource = append(edgeSource, indexMatrix[y][x])
						edgeDest = append(edgeDest, indexMatrix[y-1][position])
					}
				}
				//add above edge
				if y != len(indexMatrix)-1 {
					position := int(math.Round(float64(x * len(indexMatrix[y+1]) / len(indexMatrix[y]))))
					if pointsInWaterMatrix[y+1][position] {

						totalOffset++
						edgeSource = append(edgeSource, indexMatrix[y][x])
						edgeDest = append(edgeDest, indexMatrix[y+1][position])
					}
				}
			}

		}
	}

	fmt.Printf("Time to generate edges of grid: %.3fs\n", time.Since(startTime).Seconds())
	startTime = time.Now()
	distanceList = CalcEdgeDistances(points, edgeSource, edgeDest)

	fmt.Printf("Time to calculate edge distances: %.3fs\n", time.Since(startTime).Seconds())
	return Graph{Nodes: points, Sources: edgeSource, Targets: edgeDest, Weights: distanceList, Offsets: offsetList, NodeMatrix: indexMatrix, NodeInWaterMatrix: pointsInWaterMatrix, intendedNodeQuantity: numOfNodes}
}

func CalcEdgeDistances(points [][]float64, src []int, dest []int) []int {
	var distances []int
	for i, d := range src {
		distance := int(dist(points[d], points[dest[i]]))
		distances = append(distances, distance)
	}
	return distances
}
