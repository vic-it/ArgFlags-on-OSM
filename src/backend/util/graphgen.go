package util

import (
	"fmt"
	"math"
	"time"
)

// transforms a node into a point (x,y) coordinates w.r.t meters
func NodeToPoint(n node) point {
	return point{}
}

// transforms a point (x,y) into a node (lat,lon) w.r.t degrees
// not sure if we need this
func PointToNode(p point) node {
	return node{}
}

// math from here https://www.cmu.edu/biolphys/deserno/pdf/sphere_equi.pdf
func GenerateGraphPoints(numberOfNodes int, coastline Coastline) ([][]float64, [][]int, [][]bool) {
	println("generating points...")
	startTime := time.Now()
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
			z := IsPointInWater(point, coastline)
			waterLatList = append(waterLatList, z)
			count++
			if count%500 == 0 {
				PrintProgress(count, numberOfNodes, "nodes")
			}
		}
		isPointInWaterMatrix = append(isPointInWaterMatrix, waterLatList)
		pointMatrix = append(pointMatrix, latList)
	}
	fmt.Printf("%d points created\n", count)
	endTime := time.Now()
	println(int(endTime.Sub(startTime).Seconds()))
	return points, pointMatrix, isPointInWaterMatrix
}

func GenerateEdges(points [][]float64, indexMatrix [][]int, pointsInWaterMatrix [][]bool) ([][]float64, []int, []int) {
	println("creating edges from points...")
	var edgeSource []int
	var edgeDest []int

	for y := 0; y < len(indexMatrix); y++ {
		latList := indexMatrix[y]
		for x := 0; x < len(latList); x++ {
			//add left edge
			if pointsInWaterMatrix[y][x] {
				if x == 0 {
					if pointsInWaterMatrix[y][len(latList)-1] {
						edgeSource = append(edgeSource, indexMatrix[y][x])
						edgeDest = append(edgeDest, indexMatrix[y][len(latList)-1])
					}
				} else {
					if pointsInWaterMatrix[y][x-1] {
						edgeSource = append(edgeSource, indexMatrix[y][x])
						edgeDest = append(edgeDest, indexMatrix[y][x-1])
					}
				}
				//add right edge
				if x == len(latList)-1 {
					if pointsInWaterMatrix[y][0] {
						edgeSource = append(edgeSource, indexMatrix[y][x])
						edgeDest = append(edgeDest, indexMatrix[y][0])
					}
				} else {
					if pointsInWaterMatrix[y][x+1] {
						edgeSource = append(edgeSource, indexMatrix[y][x])
						edgeDest = append(edgeDest, indexMatrix[y][x+1])
					}
				}
				//add below edge
				if y != 0 {
					position := int(math.Round(float64(x * len(indexMatrix[y-1]) / len(indexMatrix[y]))))
					if pointsInWaterMatrix[y-1][position] {
						edgeSource = append(edgeSource, indexMatrix[y][x])
						edgeDest = append(edgeDest, indexMatrix[y-1][position])
					}
				}
				//add above edge
				if y != len(indexMatrix)-1 {
					position := int(math.Round(float64(x * len(indexMatrix[y+1]) / len(indexMatrix[y]))))
					if pointsInWaterMatrix[y+1][position] {
						edgeSource = append(edgeSource, indexMatrix[y][x])
						edgeDest = append(edgeDest, indexMatrix[y+1][position])
					}
				}
			}
		}
	}
	fmt.Printf("%d edges generated", len(edgeDest))
	return points, edgeSource, edgeDest
}

func CalcEdgeDistances(points [][]float64, src []int, dest []int) []float64 {
	var distances []float64
	for i, d := range src {
		distance := dist(points[d], points[dest[i]])
		// fmt.Printf("(%f/%f) to (%f/%f)\n", points[d][0], points[d][1], points[dest[i]][0], points[dest[i]][1])
		// fmt.Printf("points %d - %d have a distance of %.2fm\n", src[i], dest[i], distance)
		distances = append(distances, distance)
	}
	return distances
}
