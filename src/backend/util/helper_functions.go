package util

import (
	"fmt"
	"math"
)

// checks if a point, given coordinates, is on land (false) or in water (true)
func IsPointInWater(p point) bool {
	//input should be: numberOfEdgesInTheWayToNorthpole, []edge
	//SQL command for relevant edges
	// numberOfEdgesGaruanteedInTheWay, questionableEdges := GetTestSQLResult(p)
	// for every edge e call isEdgeInTheWay(p, e)
	// if true -> number += 1

	return false
}

func dist(src []float64, dest []float64) float64 {
	const PI float64 = math.Pi

	srcLat := src[1]
	srcLon := src[0]
	destLat := src[1]
	destLon := src[0]

	radlat1 := float64(PI * srcLat / 180.0)
	radlat2 := float64(PI * destLat / 180.0)

	theta := float64(srcLon - destLon)
	radtheta := float64(PI * theta / 180.0)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	// K - 1.609344
	// N -  0.8684

	dist = dist * 0.8684
	return dist
}

func threeD_coord(lon float64, lat float64) point_threeD {
	rad := float64(6378137.0)
	// Radius of the Earth (in meters)
	cosLat := math.Cos(lat)
	sinLat := math.Sin(lat)
	cosLon := math.Cos(lon)
	sinLon := math.Sin(lon)

	x := rad * cosLon * sinLat
	y := rad * sinLon * sinLat
	z := rad * cosLat

	anspoint := point_threeD{x, y, z}

	return anspoint

}

func isEdgeInTheWay(p point, e edge) bool {
	//check if edge is in the way
	//point : {lat, lon}
	//p.lat
	//edge : {point1, point2}
	//edge.point.lat
	return true
}

// takes a whole basic format map as input, checks if one way starts where another ends -> merges them
func MergeWays(inputMap basic) int {
	mergeCounter := 0
	ways := inputMap.ways
	for firstWayID, way := range ways {
		endWay, exists := ways[way.lastNodeID]
		if exists && firstWayID != way.lastNodeID {
			newWay := mergeTwoWays(way, endWay)
			inputMap.ways[firstWayID] = newWay
			delete(inputMap.ways, endWay.nodes[0])
			mergeCounter++
		}
	}
	fmt.Println("--------------attempt to merge touching ways--------------")
	fmt.Printf("ways merged: %d\n", mergeCounter)
	fmt.Printf("ways left: %d\n", len(inputMap.ways))
	return mergeCounter
}

// takes two connected ways and merges them into one
func mergeTwoWays(startWay way, endWay way) way {
	newNodes := append(startWay.nodes, endWay.nodes[1:]...)
	newWay := way{nodes: newNodes, lastNodeID: endWay.lastNodeID}
	return newWay
}

func PrintProgress(current int, max int, unit string) {
	progress := float64(current) / float64(max)
	fmt.Printf("Generating geojson file. Progress: %2.2f%s%d%s%d %s\n\r", 100*progress, "%... - ", current, " out of ", max, unit)
}
