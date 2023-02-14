package backend

import (
	"fmt"
	"math"
	"runtime"
	"sort"
	"time"
)

// checks if a point, given coordinates, is on land (false) or in water (true)
func IsPointInWater(node []float64, coastline Coastline) bool {
	//guaranteed edges are the edges that are definitely in the way and maybeedges are edges that are not guaranteed in the way but have exactly one node above our to check node
	guaranteedEdges, maybeEdges := GetRelevantEdges(node, coastline)
	guaranteedCount := len(guaranteedEdges)
	//checks edges of which only one node was above our node and the other was below (which makes the edge possibly diagonally above OR below)
	for _, e := range maybeEdges {
		firstNodeLon := coastline.Nodes[coastline.Edges[e][0]][0]
		firstNodeLat := coastline.Nodes[coastline.Edges[e][0]][1]
		secondNodeLon := coastline.Nodes[coastline.Edges[e][1]][0]
		secondNodeLat := coastline.Nodes[coastline.Edges[e][1]][1]
		guaranteedCount += isEdgeInTheWay(node, [][]float64{{firstNodeLon, firstNodeLat}, {secondNodeLon, secondNodeLat}})
	}
	return guaranteedCount%2 == 0
}

// returns the distance between two points in kilometers (float) with standard formula
func dist(src []float64, dest []float64) float64 {
	PI := math.Pi

	srcLat := src[1]
	srcLon := src[0]
	destLat := dest[1]
	destLon := dest[0]
	//radius in km
	earthRadius := 6371.0
	phi1 := PI * srcLat / 180.0
	phi2 := PI * destLat / 180.0
	deltaPhi := PI * (srcLat - destLat) / 180.0
	deltaLambda := PI * (srcLon - destLon) / 180.0
	a := math.Sin(deltaPhi/2.0)*math.Sin(deltaPhi/2.0) + math.Cos(phi1)*math.Cos(phi2)*math.Sin(deltaLambda/2.0)*math.Sin(deltaLambda/2.0)
	c := 2.0 * math.Atan2(math.Sqrt(a), math.Sqrt(1.0-a))
	meters := earthRadius * c

	return meters
}

// return 1 if edge is in the way (to north pole), 0 else
func isEdgeInTheWay(p []float64, e [][]float64) int {
	firstLon := e[1][0]
	firstLat := e[0][1]
	secondLon := e[1][0]
	secondLat := e[1][1]
	if firstLat+(secondLat-firstLat)*(secondLon-firstLon)/(p[0]-firstLon) > p[1] {
		return 1
	}
	return 0
}

// takes a whole basic format map as input, checks if one way starts where another ends -> merges them
// not in use right now?
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

// simple function used to print progress into console with current value out of maximum value and a unit
func PrintProgress(current int, max int, unit string, startTime time.Time) {
	progress := float64(current) / float64(max)

	fmt.Printf("%.3fs - Water-checking nodes |  Progress: %2.2f%s%d%s%d %s\n\r", time.Since(startTime).Seconds(), 100*progress, "% - ", current, " out of ", max, unit)
}

// for input coordinates: estimate indes of closest node on the grid, then breadth search until it finds a node that is on the grid and in water
func GetClosestGridNode(lon float64, lat float64, graph Graph) int {
	pi := math.Pi
	// more or less reverses equidistant point generation algorithm
	a := 4.0 * pi / float64(graph.intendedNodeQuantity)
	d := math.Sqrt(a)
	Mv := math.Round(pi / d)
	dv := pi / Mv
	dp := a / dv
	//get radians
	p, v := DegToRad(lon, lat)
	// derive from v := pi * (m + 0.5) / Mv
	// v * Mv / pi = m + 0.5
	m := (v * Mv / pi) - 0.5
	startingLatPos := int(math.Round(m))
	Mp := math.Round(2.0 * pi * math.Sin(v) / dp)
	// p := 2.0 * pi * n / Mp
	n := p * Mp / (2.0 * pi)
	startingLon := int(math.Round(n))
	return getClosestValidNode(startingLatPos, startingLon, graph)
}

// FIFO breadth search for nearest node that is in water
func getClosestValidNode(startingLat int, startingLon int, graph Graph) int {
	pointMatrix := graph.NodeMatrix
	nodeInWaterList := graph.NodeInWaterMatrix
	startingLat = (startingLat + (5 * len(pointMatrix))) % len(pointMatrix)
	startingLon = (startingLon + (5 * len(pointMatrix[startingLat]))) % len(pointMatrix[startingLat])

	hasBeenChecked := make(map[int]bool)
	nodesToCheck := [][]int{{startingLat, startingLon}}
	currentIndex := 0
	for {
		curLat := nodesToCheck[currentIndex][0]
		curLon := nodesToCheck[currentIndex][1]
		if nodeInWaterList[curLat][curLon] {
			// clean up for better ram usage
			hasBeenChecked = map[int]bool{}
			runtime.GC()
			return pointMatrix[curLat][curLon]
		}
		currentIndex++
		nodesToCheck = append(nodesToCheck, getNeighbors(curLat, curLon, pointMatrix, hasBeenChecked)...)
	}
}

// returns a list of neighbor points of an input point, INCLUDING POINTS ON LAND!
func getNeighbors(lat int, lon int, pointMatrix [][]int, hasBeenChecked map[int]bool) [][]int {
	var neighborList [][]int

	//add above node to check
	if lat != len(pointMatrix)-1 {
		position := int(math.Round(float64(lon * len(pointMatrix[lat+1]) / len(pointMatrix[lat]))))
		if !hasBeenChecked[pointMatrix[lat+1][position]] {
			neighborList = append(neighborList, []int{lat + 1, position})
		}
	}
	//add below node to check
	if lat != 0 {
		position := int(math.Round(float64(lon * len(pointMatrix[lat-1]) / len(pointMatrix[lat]))))
		if !hasBeenChecked[pointMatrix[lat-1][position]] {
			neighborList = append(neighborList, []int{lat - 1, position})
		}
	}
	//add left node to check
	if lon == 0 {
		if !hasBeenChecked[pointMatrix[lat][len(pointMatrix[lat])-1]] {
			neighborList = append(neighborList, []int{lat, len(pointMatrix[lat]) - 1})
		}
	} else {
		if !hasBeenChecked[pointMatrix[lat][lon-1]] {
			neighborList = append(neighborList, []int{lat, lon - 1})
		}
	}
	//add right node to check
	if lon == len(pointMatrix[lat])-1 {
		if !hasBeenChecked[pointMatrix[lat][0]] {
			neighborList = append(neighborList, []int{lat, 0})
		}
	} else {
		if !hasBeenChecked[pointMatrix[lat][lon+1]] {
			neighborList = append(neighborList, []int{lat, lon + 1})
		}
	}
	for _, i := range neighborList {
		hasBeenChecked[pointMatrix[i[0]][i[1]]] = true
	}
	return neighborList
}

// computes all edges which are relevant for the point in polygon test (e.g. only points within a certain longitude range)
func GetRelevantEdges(node []float64, coastline Coastline) ([]int, []int) {
	if node[0] < -179.5 {
		node[0] = -179.5
	}
	//list of edges with at least one point on the left side
	var leftList []EdgeCoordinate
	//list of edges with at least  one point on the right side
	var rightList []EdgeCoordinate
	//coordinates of the nodes in the edges
	nodes := coastline.Nodes
	//IDs of the nodes (coordinates) of the edges
	edges := coastline.Edges
	//all points sorted by longitude, stored together with the ID of the respective edge
	sortedLonList := coastline.SortedLonEdgeList
	//maximum longitude difference between two points of the same edge
	maxLonDiff := coastline.maxLonDiff
	//get relevant longitudes when point to check is not close to 180°/-180° longitude (within relevant longitude range - see "maxLonDiff")
	if math.Abs(node[0])+maxLonDiff < 180 {
		// left side: lon-maxdiff to lon
		rawLeftStart := BinarySearchForID(node[0]-maxLonDiff, sortedLonList)
		rawLeftEnd := BinarySearchForID(node[0], sortedLonList)
		// right side: lon to lon+maxdiff
		rawRightStart := BinarySearchForID(node[0], sortedLonList)
		rawRightEnd := BinarySearchForID(node[0]+maxLonDiff, sortedLonList)
		// make slices
		leftList = sortedLonList[rawLeftStart:rawLeftEnd]
		rightList = sortedLonList[rawRightStart:rawRightEnd]
		//case we are too close to 180°, coming from positive side
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
		//make slices
		leftList = sortedLonList[rawLeftStart:rawLeftEnd]
		rightList = sortedLonList[rawRightStart1:rawRightEnd1]
		rightList = append(rightList, sortedLonList[rawRightStart2:rawRightEnd2]...)
		//case we are too close to -180° coming from negative side
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
	// intersects left and right list, s.t. only edges are left which intersect the longitude of the node to check
	relevantLonEdges := edgeIntersectionOfCoordinatesIntoIDs(leftList, rightList)
	var maxLatList []EdgeCoordinate
	var minLatList []EdgeCoordinate
	for _, index := range relevantLonEdges {
		maxLatList = append(maxLatList, EdgeCoordinate{edgeID: index, coordinate: math.Max(nodes[edges[index][0]][1], nodes[edges[index][1]][1])})
		minLatList = append(minLatList, EdgeCoordinate{edgeID: index, coordinate: math.Min(nodes[edges[index][0]][1], nodes[edges[index][1]][1])})
	}
	// sorts all relevant edges by latitude
	sort.Sort(ByCoordinate(maxLatList))
	sort.Sort(ByCoordinate(minLatList))

	//compute relevant latitudes
	//uses other binary search function for lat because here boundaries need to be exact!
	idOfBiggerThanMaxLat := BinarySearchForLatID(node[1], maxLatList)
	idOfBiggerThanMinLat := BinarySearchForLatID(node[1], minLatList)
	relevantMaxLat := []EdgeCoordinate{}
	relevantMinLat := []EdgeCoordinate{}
	if idOfBiggerThanMaxLat >= 0 {
		relevantMaxLat = maxLatList[idOfBiggerThanMaxLat:]
	}
	if idOfBiggerThanMinLat >= 0 {
		relevantMinLat = maxLatList[idOfBiggerThanMinLat:]
	}
	//elements definitely in the way
	var defAboveList []int
	for _, e := range relevantMinLat {
		// all relevant edges where both nodes are above our node
		defAboveList = append(defAboveList, e.edgeID)
	}
	// a list of all edges where exactly one node is above our point and one is below, so it has to be checked manually
	var maybeAboveList []int
	if len(relevantMinLat) != len(relevantMaxLat) {
		maybeAboveList = secondListMinusFirstList(relevantMinLat, relevantMaxLat)
	}
	return defAboveList, maybeAboveList
}

// calculates the longitude differences on a sphere between two points
func CalcLonDiff(lon1 float64, lon2 float64) float64 {
	abs := math.Abs(lon1 - lon2)
	if abs > 180.0 {
		return 360.0 - abs
	}
	return abs
}

// converts an input node (lon/lat) from degrees to radians (lat  theta, lon - phi)
func RadToDeg(phi float64, theta float64) (float64, float64) {
	lon := (360.0 * phi / (math.Pi * 2.0)) - 180.0
	lat := (theta * 180.0 / math.Pi) - 90.0
	return lon, lat
}

// converts an input node (lon/lat) from degrees to radians (lat  theta, lon - phi)
func DegToRad(lon float64, lat float64) (float64, float64) {
	phi := (lon + 180.0) * (math.Pi * 2.0) / 360
	theta := (lat + 90) * math.Pi / 180
	return phi, theta
}

// returns index of first element above the threshold via binary search
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
			//threshhold is out of list bounds
			return -1
		}
	}
	return low
}

// same as binary search for ID except without the return median line if median == len(list)-1
func BinarySearchForLatID(threshhold float64, list []EdgeCoordinate) int {
	low := 0
	//index of first value BELOW threshhold
	high := len(list) - 1

	for low <= high {
		median := (low + high) / 2

		if list[median].coordinate < threshhold {

			low = median + 1
		} else {
			if median == 0 {
				return median
			}
			high = median - 1
		}
		if high < 0 || low < 0 || high > len(list)-1 || low > len(list)-1 {
			//threshhold out of bounds
			return -1
		}
	}
	return low
}

// takes two lists of edge coordinate objects and gives back the intersection as an ID list
func edgeIntersectionOfCoordinatesIntoIDs(l1 []EdgeCoordinate, l2 []EdgeCoordinate) []int {
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

// takes two lists of edge ID's and gives back the intersection
func edgeIntersectionOfIDs(l1 []int, l2 []int) []int {
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

// returns an array of integers which are the IDs of the edges of the 2nd list, which dont appear in the first list
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

// stuff for sorting algorithm
type ByCoordinate []EdgeCoordinate

func (a ByCoordinate) Len() int {
	return len(a)
}
func (a ByCoordinate) Less(i, j int) bool {
	return a[i].coordinate < a[j].coordinate
}
func (a ByCoordinate) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func getReverseEdgeID(graph Graph, edgeID int) int {
	sourceNodeID := graph.Sources[edgeID]
	targetNodeID := graph.Targets[edgeID]
	reverseOffsetStart := graph.Offsets[targetNodeID]
	for i := 0; i <= 5; i++ {
		if reverseOffsetStart+i > len(graph.Sources)-1 {

			println("no back edge?")
			return -1
		}
		if graph.Targets[reverseOffsetStart+i] == sourceNodeID && graph.Sources[reverseOffsetStart+i] == targetNodeID {
			// fmt.Printf("input edge: %d to %d, reverse edge: %d to %d\n", sourceNodeID, targetNodeID, graph.Sources[reverseOffsetStart+i], graph.Targets[reverseOffsetStart+i])
			return reverseOffsetStart + i
		}
	}
	println("huh")
	return -1
}
