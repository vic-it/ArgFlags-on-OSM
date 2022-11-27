package util

import (
	"fmt"
	"math"
	"sort"
)

// checks if a point, given coordinates, is on land (false) or in water (true)
func IsPointInWater(node []float64, coastline Coastline) bool {
	//guaranteed edges are the edges that are definitely in the way and maybeedges are edges that are not guaranteed in the way but have exactly one node above our to check node
	guaranteedEdges, maybeEdges := GetRelevantEdges(node, coastline)
	guaranteedCount := len(guaranteedEdges)
	fmt.Printf("--------------------\n(%f/%f) (lat/lon) has the following guaranteed edges in the way:\n", node[1], node[0])
	for _, e := range guaranteedEdges {
		firstNodeLon := coastline.Nodes[coastline.Edges[e][0]][0]
		firstNodeLat := coastline.Nodes[coastline.Edges[e][0]][1]
		secondNodeLon := coastline.Nodes[coastline.Edges[e][1]][0]
		secondNodeLat := coastline.Nodes[coastline.Edges[e][1]][1]
		fmt.Printf("lat: [%f to %f]\n", firstNodeLat, secondNodeLat)
		fmt.Printf("lat: [%f to %f]\n-\n", firstNodeLon, secondNodeLon)
	}
	x := guaranteedCount
	fmt.Printf("count: %d\n", guaranteedCount)
	fmt.Printf("(%f/%f) (lat/lon) has the following maybe edges in the way:\n", node[1], node[0])
	for _, e := range maybeEdges {
		firstNodeLon := coastline.Nodes[coastline.Edges[e][0]][0]
		firstNodeLat := coastline.Nodes[coastline.Edges[e][0]][1]
		secondNodeLon := coastline.Nodes[coastline.Edges[e][1]][0]
		secondNodeLat := coastline.Nodes[coastline.Edges[e][1]][1]
		guaranteedCount += isEdgeInTheWay(node, [][]float64{{firstNodeLon, firstNodeLat}, {secondNodeLon, secondNodeLat}})
		fmt.Printf("lat: [%f to %f]\n", firstNodeLat, secondNodeLat)
		fmt.Printf("lat: [%f to %f]\n-\n", firstNodeLon, secondNodeLon)
	}
	fmt.Printf("count: %d\n", guaranteedCount-x)

	fmt.Printf("(%f/%f) in water: %t\n", node[0], node[1], guaranteedCount%2 == 0)
	return guaranteedCount%2 == 0
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

// return 1 if edge is in the way, 0 else
func isEdgeInTheWay(p []float64, e [][]float64) int {
	firstLon := e[1][0]
	firstLat := e[0][1]
	secondLon := e[1][0]
	secondLat := e[1][1]
	var lat_cross float64
	lat_cross = firstLat + (secondLat-firstLat)*(secondLon-firstLon)/(p[0]-firstLon)
	if lat_cross > p[1] {
		return 1
	}
	return 0
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

func GetRelevantEdges(node []float64, coastline Coastline) ([]int, []int) {
	var leftList []EdgeCoordinate
	var rightList []EdgeCoordinate
	nodes := coastline.Nodes
	edges := coastline.Edges
	sortedLonList := coastline.SortedLonEdgeList
	maxLonDiff := coastline.MaxLonDiff

	//regular case
	if math.Abs(node[0])+maxLonDiff <= 180 {
		// left side: lon-maxdiff to lon
		rawLeftStart := BinarySearchForID(node[0]-maxLonDiff, sortedLonList)
		rawLeftEnd := BinarySearchForID(node[0], sortedLonList)
		// right side: lon to lon+maxdiff
		rawRightStart := BinarySearchForID(node[0], sortedLonList)
		rawRightEnd := BinarySearchForID(node[0]+maxLonDiff, sortedLonList)
		// make slices
		leftList = sortedLonList[rawLeftStart:rawLeftEnd]
		rightList = sortedLonList[rawRightStart:rawRightEnd]
		//case we are too close to 180, coming from left side
	} else if node[0]+maxLonDiff > 180 {
		// left side from lon-maxdiff to node
		rawLeftStart := BinarySearchForID(node[0]-maxLonDiff, sortedLonList)
		rawLeftEnd := BinarySearchForID(node[0], sortedLonList)
		// right side from node to 180
		rawRightStart1 := BinarySearchForID(node[0], sortedLonList)
		rawRightEnd1 := BinarySearchForID(180.0, sortedLonList)
		// right side from -180 to rest of right nodes (e.g. to -175)
		rawRightStart2 := BinarySearchForID(-180, sortedLonList)
		rawRightEnd2 := BinarySearchForID(node[0]+maxLonDiff-360.0, sortedLonList)
		//make slices
		leftList = sortedLonList[rawLeftStart:rawLeftEnd]
		rightList = sortedLonList[rawRightStart1:rawRightEnd1]
		rightList = append(rightList, sortedLonList[rawRightStart2:rawRightEnd2]...)
		//case we are too close to -180 coming from right side
	} else {
		// left side from -180 to lon
		rawLeftStart1 := BinarySearchForID(-180, sortedLonList)
		rawLeftEnd1 := BinarySearchForID(node[0], sortedLonList)
		// remainder of left side (e.g. from 175) to 180
		rawLeftStart2 := BinarySearchForID(node[0]-maxLonDiff+360.0, sortedLonList)
		rawLeftEnd2 := BinarySearchForID(180, sortedLonList)
		// right side from lon to lon+diff
		rawRightStart := BinarySearchForID(node[0], sortedLonList)
		rawRightEnd := BinarySearchForID(node[0]+maxLonDiff, sortedLonList)
		// make slices
		leftList = sortedLonList[rawLeftStart1:rawLeftEnd1]
		leftList = append(leftList, sortedLonList[rawLeftStart2:rawLeftEnd2]...)
		rightList = sortedLonList[rawRightStart:rawRightEnd]
	}

	relevantLonEdges := mergeEdgeCoordinateLists(leftList, rightList)
	// generate lat list out of longitude-relevant edges
	var maxLatList []EdgeCoordinate
	var minLatList []EdgeCoordinate
	for _, index := range relevantLonEdges {
		maxLatList = append(maxLatList, EdgeCoordinate{edgeID: index, coordinate: math.Max(nodes[edges[index][0]][1], nodes[edges[index][1]][1])})
		minLatList = append(maxLatList, EdgeCoordinate{edgeID: index, coordinate: math.Min(nodes[edges[index][0]][1], nodes[edges[index][1]][1])})
	}
	println("before sort:")
	for _, x := range maxLatList {
		fmt.Printf("lat coordinate: %f\n", x.coordinate)
	}
	sort.Sort(ByCoordinate(maxLatList))
	sort.Sort(ByCoordinate(minLatList))

	println("after sort:")
	for _, x := range maxLatList {
		fmt.Printf("lat coordinate: %f\n", x.coordinate)
	}

	//compute relevant latitudes
	idOfBiggerThanMaxLat := BinarySearchForID(node[1], maxLatList)
	idOfBiggerThanMinLat := BinarySearchForID(node[1], minLatList)
	relevantMaxLat := []EdgeCoordinate{}
	relevantMinLat := []EdgeCoordinate{}
	if idOfBiggerThanMaxLat >= 0 {
		relevantMaxLat = maxLatList[idOfBiggerThanMaxLat:]
	}
	if idOfBiggerThanMinLat >= 0 {
		relevantMinLat = maxLatList[idOfBiggerThanMinLat:]
	}
	println("max lat above:")
	for _, x := range relevantMaxLat {
		e := x.edgeID
		fmt.Printf("lon edge: lat[%f - %f] lon[%f - %f]\n", nodes[edges[e][0]][1], nodes[edges[e][1]][1], nodes[edges[e][0]][0], nodes[edges[e][1]][0])
	}
	defAboveList := mergeEdgeCoordinateLists(relevantMaxLat, relevantMinLat)
	//elements definitely in the way
	defRelevantEdges := mergeIDLists(relevantLonEdges, defAboveList)
	edgesWhereOnePointIsBelow := secondListMinusFirstList(relevantMinLat, relevantMaxLat)
	maybeRelevantEdges := mergeIDLists(edgesWhereOnePointIsBelow, relevantLonEdges)
	//get list of edges _maybe_ in the way, not guaranteed

	//intersection of relevantMinLat and the relevant longitude wise edges -> number of guaranteed edges
	//intersection of relevantMaxLat and the relevant longitude wise edges -> list of maybe in the way edges +  guaranteed in the way edges

	// maybe relevant: maxLat > node[1], minLat < node [1]
	// ===> and one lon on left side and one lon on right side

	return defRelevantEdges, maybeRelevantEdges
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
			if median == len(list)-1 {
				return median
			}
			low = median + 1
		} else {
			if median == 0 {
				return median
			}
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
	//fmt.Printf("index of list: %d - list size: %d elements\n", low, len(list))
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

func secondListMinusFirstList(l1 []EdgeCoordinate, l2 []EdgeCoordinate) []int {
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
