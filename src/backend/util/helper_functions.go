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
func MergeWays(inputMap Basic) int {
	mergeCounter := 0
	ways := inputMap.ways
	//iterate over all ways
	for firstWayID, way := range ways {
		endWay, exists := ways[way.lastNodeID]
		if exists && firstWayID != way.lastNodeID {
			newWay := mergeTwoWays(way, endWay)
			inputMap.ways[firstWayID] = newWay
			delete(inputMap.ways, endWay.nodes[0])
			mergeCounter++
		}
	}
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

func GetClosestGridNode(lon float64, lat float64) (float64, float64) {
	return lon, lat
}

func GetRelevantEdges(node []float64, sortedLonList []EdgeCoordinate, maxLatList []EdgeCoordinate, minLatList []EdgeCoordinate, maxLonDiff float64) (int, []int) {
	var leftList []EdgeCoordinate
	var rightList []EdgeCoordinate
	// numberOfEdgesGaruanteedInTheWay
	var n int

	//regular case
	if math.Abs(node[0])+maxLonDiff < 180 {
		// left side: lon-maxdiff to lon
		rawLeftStart := BinarySearchForID(node[0]-maxLonDiff, sortedLonList)
		rawLeftEnd := BinarySearchForID(node[0], sortedLonList)
		// right side: lon to lon+maxdiff
		rawRighStart := BinarySearchForID(node[0], sortedLonList)
		rawRightEnd := BinarySearchForID(node[0]+maxLonDiff, sortedLonList)
		// MAKE CLEANED START INDEX (1 extra element from each direction, but not out of bounds)
		leftStart := int(math.Max(0, float64(rawLeftStart-1.0)))
		leftEnd := int(math.Min(float64((len(sortedLonList) - 1)), float64(rawLeftEnd+1)))
		rightStart := int(math.Max(0, float64(rawRighStart-1.0)))
		rightEnd := int(math.Min(float64((len(sortedLonList) - 1)), float64(rawRightEnd+1)))
		// make slices
		leftList = sortedLonList[leftStart:leftEnd]
		rightList = sortedLonList[rightStart:rightEnd]
		//case we are too close to 180, coming from left side
	} else if node[0]+maxLonDiff >= 180 {
		// left side from lon-maxdiff to node
		rawLeftStart := BinarySearchForID(node[0]-maxLonDiff, sortedLonList)
		rawLeftEnd := BinarySearchForID(node[0], sortedLonList)
		// right side from node to 180
		rawRightStart1 := BinarySearchForID(node[0], sortedLonList)
		rawRightEnd1 := BinarySearchForID(180.0, sortedLonList)
		// right side from -180 to rest of right nodes (e.g. to -175)
		rawRightStart2 := BinarySearchForID(-180, sortedLonList)
		rawRightEnd2 := BinarySearchForID(node[0]+maxLonDiff-360.0, sortedLonList)
		// MAKE CLEANED START INDEX (1 extra element from each direction, but not out of bounds)
		leftStart := int(math.Max(0, float64(rawLeftStart-1.0)))
		leftEnd := int(math.Min(float64((len(sortedLonList) - 1)), float64(rawLeftEnd+1)))
		rightStart1 := int(math.Max(0, float64(rawRightStart1-1.0)))
		rightEnd1 := int(math.Min(float64((len(sortedLonList) - 1)), float64(rawRightEnd1+1)))
		rightStart2 := int(math.Max(0, float64(rawRightStart2-1.0)))
		rightEnd2 := int(math.Min(float64((len(sortedLonList) - 1)), float64(rawRightEnd2+1)))
		//make slices
		leftList = sortedLonList[leftStart:leftEnd]
		rightList = sortedLonList[rightStart1:rightEnd1]
		rightList = append(rightList, sortedLonList[rightStart2:rightEnd2]...)
		//case we are too close to -180 coming from right side
	} else {
		// left side from -180 to lon
		rawLeftStart1 := BinarySearchForID(-180.0, sortedLonList)
		rawLeftEnd1 := BinarySearchForID(node[0], sortedLonList)
		// remainder of left side (e.g. from 175) to 180
		rawLeftStart2 := BinarySearchForID(node[0]-maxLonDiff+360.0, sortedLonList)
		rawLeftEnd2 := BinarySearchForID(180, sortedLonList)
		// right side from lon to lon+diff
		rightStart := BinarySearchForID(node[0], sortedLonList)
		rightEnd := BinarySearchForID(node[0]+maxLonDiff, sortedLonList)
		// make slices
		leftList = sortedLonList[int(math.Max(0, float64(rawLeftStart1-1.0))):int(math.Min(float64((len(sortedLonList)-1)), float64(rawLeftEnd1+1)))]
		leftList = append(leftList, sortedLonList[int(math.Max(0, float64(rawLeftStart2-1.0))):int(math.Min(float64((len(sortedLonList)-1)), float64(rawLeftEnd2+1)))]...)
		rightList = sortedLonList[int(math.Max(0, float64(rightStart-1.0))):int(math.Min(float64((len(sortedLonList)-1)), float64(rightEnd+1)))]
	}

	//compute relevant latitudes
	relevantMaxLat := maxLatList[BinarySearchForID(node[1], maxLatList):]
	relevantMinLat := minLatList[BinarySearchForID(node[1], minLatList):]

	relevantLonEdges := mergeEdgeCoordinateLists(leftList, rightList)
	defAboveList := mergeEdgeCoordinateLists(relevantMaxLat, relevantMinLat)
	//elements definitely in the way
	n = len(mergeIDLists(relevantLonEdges, defAboveList))
	edgesWhereOnePointIsBelow := secondListMinusSecondList(relevantMinLat, relevantMaxLat)
	maybeRelevantEdges := mergeIDLists(edgesWhereOnePointIsBelow, relevantLonEdges)
	//get list of edges _maybe_ in the way, not guaranteed

	//intersection of relevantMinLat and the relevant longitude wise edges -> number of guaranteed edges
	//intersection of relevantMaxLat and the relevant longitude wise edges -> list of maybe in the way edges +  guaranteed in the way edges

	// maybe relevant: maxLat > node[1], minLat < node [1]
	// ===> and one lon on left side and one lon on right side

	return n, maybeRelevantEdges
}

func CalcLonDiff(lon1 float64, lon2 float64) float64 {
	abs := math.Abs(lon1 - lon2)
	if abs > 180.0 {
		return 360.0 - abs
	}
	return abs
}

// returns low (or high), if it returns -1 -> threshhold out of list
func BinarySearchForID(threshhold float64, list []EdgeCoordinate) int {
	//index of first value ABOVE threshhold
	low := 0
	//index of first value BELOW threshhold
	high := len(list) - 1

	for low <= high {
		median := (low + high) / 2

		if list[median].coordinate < threshhold {
			low = median + 1
		} else {
			high = median - 1
		}
		if high < 0 || low < 0 || high > len(list)-1 || low > len(list)-1 {
			println("threshhold is out of list bounds")
			return -1
		}
		//fmt.Printf("threshhold: %f\nvalue at low(%d): %f\nvalue at high(%d): %f\n", threshhold, low, list[low].coordinate, high, list[high].coordinate)
		//10° -> 175 -> left side in 165-175, right side 175-180 and -180 to -175
		//10° -> 135 -> left side in 125-135, right side 135-145
	}

	return low
}

func mergeEdgeCoordinateLists(l1 []EdgeCoordinate, l2 []EdgeCoordinate) []int {
	m := make(map[int]bool)
	var c []int

	for _, item := range l1 {
		m[item.edgeID] = true
	}

	for _, item := range l2 {
		if _, ok := m[item.edgeID]; ok {
			c = append(c, item.edgeID)
		}
	}
	return c
}

func mergeIDLists(l1 []int, l2 []int) []int {
	m := make(map[int]bool)
	var c []int

	for _, item := range l1 {
		m[item] = true
	}

	for _, item := range l2 {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}
	return c
}

func secondListMinusSecondList(l1 []EdgeCoordinate, l2 []EdgeCoordinate) []int {
	m := make(map[int]bool)
	var c []int

	for _, item := range l1 {
		m[item.edgeID] = true
	}

	for _, item := range l2 {
		if _, ok := m[item.edgeID]; !ok {
			c = append(c, item.edgeID)
		}
	}
	return c
}
