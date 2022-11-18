package util

import (
	"fmt"
	"math"
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
func GenerateGraphPoints(numberOfNodes int) [][]float64 {
	var points [][]float64
	pi := math.Pi
	//count of nodes
	count := 0
	a := 4.0 * pi / float64(numberOfNodes)
	d := math.Sqrt(a)
	Mv := math.Round(pi / d)
	dv := pi / Mv
	dp := a / dv
	for m := 0.0; m < Mv; m++ {
		v := pi * (m + 0.5) / Mv
		Mp := math.Round(2.0 * pi * math.Sin(v) / dp)
		for n := 0.0; n < Mp; n++ {
			//generate point?
			p := 2.0 * pi * n / Mp
			var point []float64
			lon, lat := radToDeg(v, p)
			point = append(point, lon)
			point = append(point, lat)
			points = append(points, point)
			count++
		}
	}
	fmt.Printf("%d points created\n", count)
	return points
}

func radToDeg(theta float64, phi float64) (float64, float64) {
	lon := (360.0 * phi / (math.Pi * 2.0)) - 180.0
	lat := (theta * 180.0 / math.Pi) - 90.0
	return lon, lat
}
